package stats_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourusername/pulsar-watch/internal/stats"
)

func TestNew(t *testing.T) {
	s := stats.New()
	if s == nil {
		t.Fatal("expected non-nil Stats")
	}
	if s.StartTime.IsZero() {
		t.Error("expected StartTime to be set")
	}
}

func TestRecordSeen(t *testing.T) {
	s := stats.New()
	s.RecordSeen()
	s.RecordSeen()
	snap := s.Snapshot()
	if snap.TotalSeen != 2 {
		t.Errorf("expected TotalSeen=2, got %d", snap.TotalSeen)
	}
}

func TestRecordMatched(t *testing.T) {
	s := stats.New()
	s.RecordSeen()
	s.RecordMatched()
	snap := s.Snapshot()
	if snap.TotalMatched != 1 {
		t.Errorf("expected TotalMatched=1, got %d", snap.TotalMatched)
	}
}

func TestRecordExported(t *testing.T) {
	s := stats.New()
	s.RecordExported()
	snap := s.Snapshot()
	if snap.TotalExported != 1 {
		t.Errorf("expected TotalExported=1, got %d", snap.TotalExported)
	}
}

func TestSnapshot_LastMessageAt(t *testing.T) {
	s := stats.New()
	before := time.Now()
	s.RecordSeen()
	after := time.Now()
	snap := s.Snapshot()
	if snap.LastMessageAt.Before(before) || snap.LastMessageAt.After(after) {
		t.Error("LastMessageAt not within expected range")
	}
}

func TestSnapshot_String_ContainsFields(t *testing.T) {
	s := stats.New()
	s.RecordSeen()
	s.RecordSeen()
	s.RecordMatched()
	s.RecordExported()
	str := s.Snapshot().String()
	for _, want := range []string{"seen=2", "matched=1", "exported=1", "match_rate=", "duration="} {
		if !strings.Contains(str, want) {
			t.Errorf("expected %q in snapshot string: %s", want, str)
		}
	}
}

func TestSnapshot_MatchRate_ZeroSeen(t *testing.T) {
	s := stats.New()
	str := s.Snapshot().String()
	if !strings.Contains(str, "match_rate=0.0%") {
		t.Errorf("expected match_rate=0.0%% when no messages seen, got: %s", str)
	}
}
