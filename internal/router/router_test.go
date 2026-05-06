package router

import (
	"sync"
	"testing"
)

func TestRegister_EmptyName(t *testing.T) {
	r := New()
	err := r.Register("", ".*", func(_ string, _ []byte) {})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRegister_NilHandler(t *testing.T) {
	r := New()
	err := r.Register("h1", ".*", nil)
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestRegister_InvalidPattern(t *testing.T) {
	r := New()
	err := r.Register("h1", "[invalid", func(_ string, _ []byte) {})
	if err == nil {
		t.Fatal("expected error for invalid regex pattern")
	}
}

func TestRegister_DuplicateName(t *testing.T) {
	r := New()
	h := func(_ string, _ []byte) {}
	if err := r.Register("h1", ".*", h); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err := r.Register("h1", "orders.*", h)
	if err == nil {
		t.Fatal("expected error for duplicate handler name")
	}
}

func TestDispatch_NoMatch(t *testing.T) {
	r := New()
	_ = r.Register("h1", "^orders/", func(_ string, _ []byte) {})

	n := r.Dispatch("payments/pay-123", []byte("data"))
	if n != 0 {
		t.Fatalf("expected 0 dispatches, got %d", n)
	}
}

func TestDispatch_SingleMatch(t *testing.T) {
	r := New()
	var got string
	_ = r.Register("h1", "^orders/", func(topic string, _ []byte) { got = topic })

	n := r.Dispatch("orders/ord-42", []byte("payload"))
	if n != 1 {
		t.Fatalf("expected 1 dispatch, got %d", n)
	}
	if got != "orders/ord-42" {
		t.Fatalf("expected topic orders/ord-42, got %q", got)
	}
}

func TestDispatch_MultipleHandlers(t *testing.T) {
	r := New()
	var mu sync.Mutex
	called := map[string]bool{}

	for _, name := range []string{"h1", "h2"} {
		n := name
		_ = r.Register(n, ".*", func(_ string, _ []byte) {
			mu.Lock()
			called[n] = true
			mu.Unlock()
		})
	}

	n := r.Dispatch("any/topic", []byte{})
	if n != 2 {
		t.Fatalf("expected 2 dispatches, got %d", n)
	}
	if !called["h1"] || !called["h2"] {
		t.Fatal("expected both handlers to be called")
	}
}

func TestLen(t *testing.T) {
	r := New()
	if r.Len() != 0 {
		t.Fatal("expected 0 routes initially")
	}
	_ = r.Register("h1", ".*", func(_ string, _ []byte) {})
	_ = r.Register("h2", "orders/.*", func(_ string, _ []byte) {})
	if r.Len() != 2 {
		t.Fatalf("expected 2 routes, got %d", r.Len())
	}
}
