package groupby

import (
	"testing"
)

func TestNew_ZeroMaxKeys(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero maxKeys")
	}
}

func TestNew_NegativeMaxKeys(t *testing.T) {
	_, err := New(-1)
	if err == nil {
		t.Fatal("expected error for negative maxKeys")
	}
}

func TestNew_Valid(t *testing.T) {
	g, err := New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil Grouper")
	}
}

func TestRecordSeen_CreatesGroup(t *testing.T) {
	g, _ := New(10)
	g.RecordSeen("topic-a")
	snap := g.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 group, got %d", len(snap))
	}
	if snap[0].Seen != 1 {
		t.Errorf("expected Seen=1, got %d", snap[0].Seen)
	}
}

func TestRecordMatched_IncrementsMatched(t *testing.T) {
	g, _ := New(10)
	g.RecordSeen("topic-a")
	g.RecordMatched("topic-a")
	snap := g.Snapshot()
	if snap[0].Matched != 1 {
		t.Errorf("expected Matched=1, got %d", snap[0].Matched)
	}
}

func TestRecordSeen_MaxKeysEnforced(t *testing.T) {
	g, _ := New(2)
	g.RecordSeen("a")
	g.RecordSeen("b")
	g.RecordSeen("c") // should be dropped
	snap := g.Snapshot()
	if len(snap) != 2 {
		t.Errorf("expected 2 groups (cap), got %d", len(snap))
	}
}

func TestReset_ClearsGroups(t *testing.T) {
	g, _ := New(10)
	g.RecordSeen("x")
	g.Reset()
	snap := g.Snapshot()
	if len(snap) != 0 {
		t.Errorf("expected 0 groups after reset, got %d", len(snap))
	}
}

func TestRecordSeen_LastSeenUpdated(t *testing.T) {
	g, _ := New(10)
	g.RecordSeen("k")
	snap := g.Snapshot()
	if snap[0].LastSeen.IsZero() {
		t.Error("expected LastSeen to be set")
	}
}
