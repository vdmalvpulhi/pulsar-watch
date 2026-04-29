package tail_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/pulsar-watch/internal/consumer"
	"github.com/user/pulsar-watch/internal/output"
	"github.com/user/pulsar-watch/internal/tail"
)

// stubMessage satisfies consumer.Message for testing.
type stubMessage struct {
	key     string
	payload []byte
	pubTime time.Time
}

func (m *stubMessage) Key() string            { return m.key }
func (m *stubMessage) Payload() []byte        { return m.payload }
func (m *stubMessage) PublishTime() time.Time { return m.pubTime }
func (m *stubMessage) Properties() map[string]string { return nil }

// stubConsumer feeds a fixed list of messages then blocks until ctx is done.
type stubConsumer struct {
	msgs []*stubMessage
	pos  int
}

func (c *stubConsumer) Receive(ctx context.Context) (consumer.Message, error) {
	if c.pos < len(c.msgs) {
		m := c.msgs[c.pos]
		c.pos++
		return m, nil
	}
	<-ctx.Done()
	return nil, ctx.Err()
}
func (c *stubConsumer) Ack(_ consumer.Message) {}
func (c *stubConsumer) Close() {}

func newOut(buf *bytes.Buffer) *output.Output {
	out, _ := output.NewWithWriter(buf, false)
	return out
}

func TestNew_NilConsumer(t *testing.T) {
	_, err := tail.New(tail.Options{Output: newOut(new(bytes.Buffer))})
	if err == nil {
		t.Fatal("expected error for nil consumer")
	}
}

func TestNew_NilOutput(t *testing.T) {
	_, err := tail.New(tail.Options{Consumer: &consumer.Consumer{}})
	if err == nil {
		t.Fatal("expected error for nil output")
	}
}

func TestRun_MaxMessages(t *testing.T) {
	buf := new(bytes.Buffer)
	sc := &stubConsumer{
		msgs: []*stubMessage{
			{key: "k1", payload: []byte("hello")},
			{key: "k2", payload: []byte("world")},
			{key: "k3", payload: []byte("extra")},
		},
	}
	tl, err := tail.New(tail.Options{
		Consumer:    sc,
		Output:      newOut(buf),
		MaxMessages: 2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := tl.Run(context.Background()); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if sc.pos != 2 {
		t.Errorf("expected 2 messages consumed, got %d", sc.pos)
	}
}

func TestRun_ContextCancel(t *testing.T) {
	buf := new(bytes.Buffer)
	sc := &stubConsumer{msgs: []*stubMessage{}}
	tl, _ := tail.New(tail.Options{
		Consumer: sc,
		Output:   newOut(buf),
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := tl.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestRun_ShowTimestamp(t *testing.T) {
	buf := new(bytes.Buffer)
	ts := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	sc := &stubConsumer{
		msgs: []*stubMessage{{key: "k", payload: []byte("msg"), pubTime: ts}},
	}
	tl, _ := tail.New(tail.Options{
		Consumer:      sc,
		Output:        newOut(buf),
		MaxMessages:   1,
		ShowTimestamp: true,
	})
	tl.Run(context.Background())
	if !bytes.Contains(buf.Bytes(), []byte("2024-01-15")) {
		t.Errorf("expected timestamp in output, got: %s", buf.String())
	}
}
