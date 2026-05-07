package topicmap

import (
	"testing"
	"time"
)

func TestNew_NegativeMaxKeys(t *testing.T) {
	_, err := New(-1)
	if err == nil {
		t.Fatal("expected error for negative maxKeys")
	}
}

func TestNew_ZeroMaxKeys(t *testing.T) {
	tm, err := New(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tm == nil {
		t.Fatal("expected non-nil TopicMap")
	}
}

func TestNew_PositiveMaxKeys(t *testing.T) {
	tm, err := New(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tm.maxKeys != 5 {
		t.Fatalf("expected maxKeys=5, got %d", tm.maxKeys)
	}
}

func TestRecordSeen_IncrementsCount(t *testing.T) {
	tm, _ := New(0)
	tm.RecordSeen("topic-a")
	tm.RecordSeen("topic-a")
	snap := tm.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snap))
	}
	if snap[0].Seen != 2 {
		t.Fatalf("expected Seen=2, got %d", snap[0].Seen)
	}
}

func TestRecordMatched_IncrementsMatched(t *testing.T) {
	tm, _ := New(0)
	tm.RecordSeen("topic-b")
	tm.RecordMatched("topic-b")
	snap := tm.Snapshot()
	if snap[0].Matched != 1 {
		t.Fatalf("expected Matched=1, got %d", snap[0].Matched)
	}
}

func TestRecordSeen_UpdatesLastSeen(t *testing.T) {
	tm, _ := New(0)
	before := time.Now()
	tm.RecordSeen("topic-c")
	after := time.Now()
	snap := tm.Snapshot()
	ls := snap[0].LastSeen
	if ls.Before(before) || ls.After(after) {
		t.Fatalf("LastSeen %v not between %v and %v", ls, before, after)
	}
}

func TestMaxKeys_DropsExcess(t *testing.T) {
	tm, _ := New(2)
	tm.RecordSeen("t1")
	tm.RecordSeen("t2")
	tm.RecordSeen("t3") // should be dropped
	if tm.Len() != 2 {
		t.Fatalf("expected 2 topics, got %d", tm.Len())
	}
}

func TestSnapshot_SortedBySeenDesc(t *testing.T) {
	tm, _ := New(0)
	tm.RecordSeen("low")
	for i := 0; i < 5; i++ {
		tm.RecordSeen("high")
	}
	snap := tm.Snapshot()
	if snap[0].Topic != "high" {
		t.Fatalf("expected high first, got %s", snap[0].Topic)
	}
}

func TestSnapshot_TieBreakByName(t *testing.T) {
	tm, _ := New(0)
	tm.RecordSeen("zebra")
	tm.RecordSeen("alpha")
	snap := tm.Snapshot()
	if snap[0].Topic != "alpha" {
		t.Fatalf("expected alpha first on tie, got %s", snap[0].Topic)
	}
}
