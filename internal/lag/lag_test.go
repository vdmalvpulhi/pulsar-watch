package lag

import (
	"testing"
)

func TestNew_EmptyTracker(t *testing.T) {
	tr := New()
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
	if got := tr.Snapshot(); len(got) != 0 {
		t.Fatalf("expected empty snapshot, got %d entries", len(got))
	}
}

func TestRecord_EmptyTopic(t *testing.T) {
	tr := New()
	if err := tr.Record("", "sub", 10); err == nil {
		t.Fatal("expected error for empty topic")
	}
}

func TestRecord_EmptySubscription(t *testing.T) {
	tr := New()
	if err := tr.Record("topic", "", 10); err == nil {
		t.Fatal("expected error for empty subscription")
	}
}

func TestRecord_NegativeBacklog(t *testing.T) {
	tr := New()
	if err := tr.Record("topic", "sub", -1); err == nil {
		t.Fatal("expected error for negative backlog")
	}
}

func TestRecord_Valid(t *testing.T) {
	tr := New()
	if err := tr.Record("topic-a", "sub-1", 42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := tr.Get("topic-a", "sub-1")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Backlog != 42 {
		t.Fatalf("expected backlog 42, got %d", e.Backlog)
	}
}

func TestGet_NotFound(t *testing.T) {
	tr := New()
	_, ok := tr.Get("missing", "sub")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestSnapshot_MultipleEntries(t *testing.T) {
	tr := New()
	_ = tr.Record("t1", "s1", 5)
	_ = tr.Record("t2", "s2", 10)
	snap := tr.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	tr := New()
	_ = tr.Record("t1", "s1", 100)
	tr.Reset()
	if got := tr.Snapshot(); len(got) != 0 {
		t.Fatalf("expected empty snapshot after reset, got %d", len(got))
	}
}

func TestRecord_UpdatesExisting(t *testing.T) {
	tr := New()
	_ = tr.Record("t1", "s1", 10)
	_ = tr.Record("t1", "s1", 99)
	e, _ := tr.Get("t1", "s1")
	if e.Backlog != 99 {
		t.Fatalf("expected updated backlog 99, got %d", e.Backlog)
	}
}
