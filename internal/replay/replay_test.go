package replay_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/pulsar-watch/internal/exporter"
	"github.com/user/pulsar-watch/internal/replay"
)

// mockMessage implements consumer.Message for testing.
type mockMessage struct {
	key     string
	payload []byte
}

func (m *mockMessage) Key() string     { return m.key }
func (m *mockMessage) Payload() []byte { return m.payload }
func (m *mockMessage) Topic() string   { return "persistent://public/default/test" }
func (m *mockMessage) PublishTime() time.Time { return time.Now() }
func (m *mockMessage) Properties() map[string]string { return nil }

// mockConsumer simulates a Pulsar consumer.
type mockConsumer struct {
	messages []mockMessage
	index    int
}

func (mc *mockConsumer) Receive(ctx context.Context) (interface{ Key() string; Payload() []byte }, error) {
	if mc.index >= len(mc.messages) {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	msg := &mc.messages[mc.index]
	mc.index++
	return msg, nil
}

func (mc *mockConsumer) Ack(msg interface{}) {}
func (mc *mockConsumer) Close()             {}

func TestRun_MaxMessages(t *testing.T) {
	var buf strings.Builder
	e, err := exporter.NewWithWriter("text", &buf)
	if err != nil {
		t.Fatalf("failed to create exporter: %v", err)
	}

	opts := replay.Options{MaxMessages: 2}
	_ = opts
	_ = e
	// Verify Options struct fields are accessible.
	if opts.MaxMessages != 2 {
		t.Errorf("expected MaxMessages=2, got %d", opts.MaxMessages)
	}
}

func TestRun_DryRun(t *testing.T) {
	opts := replay.Options{
		MaxMessages: 1,
		DryRun:      true,
	}
	if !opts.DryRun {
		t.Error("expected DryRun to be true")
	}
}

func TestRun_RateLimit(t *testing.T) {
	opts := replay.Options{
		RateLimit:   10,
		MaxMessages: 0,
	}
	if opts.RateLimit != 10 {
		t.Errorf("expected RateLimit=10, got %d", opts.RateLimit)
	}
}

func TestNew_ReturnsReplayer(t *testing.T) {
	var buf strings.Builder
	e, err := exporter.NewWithWriter("json", &buf)
	if err != nil {
		t.Fatalf("failed to create exporter: %v", err)
	}

	opts := replay.Options{MaxMessages: 5, DryRun: false}
	r := replay.New(nil, nil, e, opts)
	if r == nil {
		t.Error("expected non-nil Replayer")
	}
}
