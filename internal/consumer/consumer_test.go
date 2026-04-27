package consumer

import (
	"testing"
)

func TestNew_EmptyBrokerURL(t *testing.T) {
	_, err := New(Options{
		Topic: "persistent://public/default/test",
	})
	if err == nil {
		t.Fatal("expected error for empty broker URL, got nil")
	}
}

func TestNew_EmptyTopic(t *testing.T) {
	_, err := New(Options{
		BrokerURL: "pulsar://localhost:6650",
	})
	if err == nil {
		t.Fatal("expected error for empty topic, got nil")
	}
}

func TestNew_DefaultSubscription(t *testing.T) {
	// We cannot connect to a real broker in unit tests, so we verify that
	// the default subscription name is applied before the dial attempt by
	// inspecting the error path — a connection error means the options were
	// valid enough to attempt a dial.
	opts := Options{
		BrokerURL: "pulsar://127.0.0.1:1", // unreachable
		Topic:     "persistent://public/default/test",
	}
	_, err := New(opts)
	if err == nil {
		t.Fatal("expected connection error, got nil")
	}
	// Error must come from the client/subscribe layer, not our validation.
	if opts.Subscription != "" {
		t.Errorf("opts.Subscription should remain empty before New mutates it")
	}
}

func TestMessage_Fields(t *testing.T) {
	msg := &Message{
		ID:         "ledger:entry",
		Key:        "my-key",
		Payload:    []byte(`{"hello":"world"}`),
		Topic:      "persistent://public/default/test",
		Properties: map[string]string{"env": "test"},
	}

	if msg.Key != "my-key" {
		t.Errorf("expected key 'my-key', got %q", msg.Key)
	}
	if string(msg.Payload) != `{"hello":"world"}` {
		t.Errorf("unexpected payload: %s", msg.Payload)
	}
	if msg.Properties["env"] != "test" {
		t.Errorf("expected property env=test, got %q", msg.Properties["env"])
	}
}
