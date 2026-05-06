package partition

import (
	"testing"
)

func TestNew_EmptyTracker(t *testing.T) {
	tr := New()
	snap := tr.Snapshot()
	if len(snap) != 0 {
		t.Fatalf("expected empty snapshot, got %d entries", len(snap))
	}
}

func TestRecord_NegativePartition(t *testing.T) {
	tr := New()
	if err := tr.Record(-1, false); err == nil {
		t.Fatal("expected error for negative partition index")
	}
}

func TestRecord_SeenIncrement(t *testing.T) {
	tr := New()
	for i := 0; i < 5; i++ {
		if err := tr.Record(0, false); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	snap := tr.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 partition, got %d", len(snap))
	}
	if snap[0].Seen != 5 {
		t.Errorf("expected Seen=5, got %d", snap[0].Seen)
	}
	if snap[0].Matched != 0 {
		t.Errorf("expected Matched=0, got %d", snap[0].Matched)
	}
}

func TestRecord_MatchedIncrement(t *testing.T) {
	tr := New()
	_ = tr.Record(2, true)
	_ = tr.Record(2, false)
	_ = tr.Record(2, true)
	snap := tr.Snapshot()
	if snap[0].Seen != 3 {
		t.Errorf("expected Seen=3, got %d", snap[0].Seen)
	}
	if snap[0].Matched != 2 {
		t.Errorf("expected Matched=2, got %d", snap[0].Matched)
	}
}

func TestSnapshot_SortedByPartition(t *testing.T) {
	tr := New()
	for _, p := range []int{3, 1, 2, 0} {
		_ = tr.Record(p, false)
	}
	snap := tr.Snapshot()
	for i := 0; i < len(snap)-1; i++ {
		if snap[i].Partition > snap[i+1].Partition {
			t.Errorf("snapshot not sorted at index %d: %d > %d",
				i, snap[i].Partition, snap[i+1].Partition)
		}
	}
}

func TestSummary_NoPartitions(t *testing.T) {
	tr := New()
	if got := tr.Summary(); got != "no partitions recorded" {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestSummary_WithData(t *testing.T) {
	tr := New()
	_ = tr.Record(0, true)
	_ = tr.Record(1, false)
	got := tr.Summary()
	if got == "" {
		t.Fatal("expected non-empty summary")
	}
}

func TestReset_ClearsData(t *testing.T) {
	tr := New()
	_ = tr.Record(0, true)
	tr.Reset()
	if snap := tr.Snapshot(); len(snap) != 0 {
		t.Errorf("expected empty snapshot after reset, got %d entries", len(snap))
	}
}
