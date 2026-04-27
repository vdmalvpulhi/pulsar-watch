package watcher_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/pulsar-watch/internal/exporter"
	"github.com/user/pulsar-watch/internal/filter"
	"github.com/user/pulsar-watch/internal/output"
	"github.com/user/pulsar-watch/internal/stats"
	"github.com/user/pulsar-watch/internal/watcher"
)

// stubMessage satisfies consumer.Message for testing.
type stubMessage struct {
	key     string
	payload []byte
}

func (m *stubMessage) Key() string       { return m.key }
func (m *stubMessage) Payload() []byte   { return m.payload }
func (m *stubMessage) Properties() map[string]string { return nil }
func (m *stubMessage) PublishTime() time.Time        { return time.Now() }

// stubConsumer delivers a fixed set of messages then blocks.
type stubConsumer struct {
	msgs []*stubMessage
	pos  int
	acked int
}

func (c *stubConsumer) Receive(ctx context.Context) (interface{ Key() string; Payload() []byte }, error) {
	if c.pos >= len(c.msgs) {
		<-ctx.Done()
		return nil, errors.New("context done")
	}
	m := c.msgs[c.pos]
	c.pos++
	return m, nil
}
func (c *stubConsumer) Ack(msg interface{}) { c.acked++ }
func (c *stubConsumer) Close()              {}

func TestNew_MissingConsumer(t *testing.T) {
	_, err := watcher.New(watcher.Config{
		Filter: &filter.Filter{},
		Stats:  stats.New(),
		Output: output.New(),
	})
	if err == nil {
		t.Fatal("expected error for nil consumer")
	}
}

func TestNew_MissingFilter(t *testing.T) {
	_, err := watcher.New(watcher.Config{
		Consumer: &stubConsumer{},
		Stats:    stats.New(),
		Output:   output.New(),
	})
	if err == nil {
		t.Fatal("expected error for nil filter")
	}
}

func TestRun_StopsOnContextCancel(t *testing.T) {
	f, _ := filter.New(filter.Options{})
	w, err := watcher.New(watcher.Config{
		Consumer: &stubConsumer{},
		Filter:   f,
		Stats:    stats.New(),
		Output:   output.New(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	if err := w.Run(ctx); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
}

func TestRun_ExporterNilIsAllowed(t *testing.T) {
	f, _ := filter.New(filter.Options{})
	w, err := watcher.New(watcher.Config{
		Consumer: &stubConsumer{},
		Filter:   f,
		Exporter: (*exporter.Exporter)(nil),
		Stats:    stats.New(),
		Output:   output.New(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	w.Run(ctx)
}
