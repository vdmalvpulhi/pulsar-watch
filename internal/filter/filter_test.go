package filter

import (
	"testing"
)

func TestNew_InvalidKeyPattern(t *testing.T) {
	_, err := New(Options{KeyPattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid key pattern, got nil")
	}
}

func TestNew_InvalidPayloadPattern(t *testing.T) {
	_, err := New(Options{PayloadPattern: "(unclosed"})
	if err == nil {
		t.Fatal("expected error for invalid payload pattern, got nil")
	}
}

func TestMatch_NoFilters(t *testing.T) {
	f, err := New(Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msg := Message{Key: "any-key", Payload: "any-payload"}
	if !f.Match(msg) {
		t.Error("expected match with no filters set")
	}
}

func TestMatch_KeyPattern(t *testing.T) {
	f, _ := New(Options{KeyPattern: "^order-"})

	if !f.Match(Message{Key: "order-123"}) {
		t.Error("expected key 'order-123' to match pattern '^order-'")
	}
	if f.Match(Message{Key: "invoice-456"}) {
		t.Error("expected key 'invoice-456' to NOT match pattern '^order-'")
	}
}

func TestMatch_PayloadPattern(t *testing.T) {
	f, _ := New(Options{PayloadPattern: `"status":"active"`})

	if !f.Match(Message{Payload: `{"id":1,"status":"active"}`}) {
		t.Error("expected payload to match")
	}
	if f.Match(Message{Payload: `{"id":2,"status":"inactive"}`}) {
		t.Error("expected payload to NOT match")
	}
}

func TestMatch_PropertyKeyOnly(t *testing.T) {
	f, _ := New(Options{PropertyKey: "region"})

	if !f.Match(Message{Properties: map[string]string{"region": "us-east"}}) {
		t.Error("expected message with 'region' property to match")
	}
	if f.Match(Message{Properties: map[string]string{"env": "prod"}}) {
		t.Error("expected message without 'region' property to NOT match")
	}
}

func TestMatch_PropertyKeyAndValue(t *testing.T) {
	f, _ := New(Options{PropertyKey: "env", PropertyValue: "prod"})

	if !f.Match(Message{Properties: map[string]string{"env": "PROD"}}) {
		t.Error("expected case-insensitive property value match")
	}
	if f.Match(Message{Properties: map[string]string{"env": "staging"}}) {
		t.Error("expected 'staging' to NOT match 'prod'")
	}
}

func TestMatch_CombinedFilters(t *testing.T) {
	f, _ := New(Options{
		KeyPattern:    "^evt-",
		PayloadPattern: "error",
		PropertyKey:   "source",
		PropertyValue: "api",
	})

	match := Message{
		Key:        "evt-99",
		Payload:    "critical error occurred",
		Properties: map[string]string{"source": "api"},
	}
	if !f.Match(match) {
		t.Error("expected combined filter to match")
	}

	noMatch := Message{
		Key:        "evt-99",
		Payload:    "all systems ok",
		Properties: map[string]string{"source": "api"},
	}
	if f.Match(noMatch) {
		t.Error("expected combined filter to NOT match when payload doesn't contain 'error'")
	}
}
