package metrics_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/pulsar-watch/internal/metrics"
	"github.com/user/pulsar-watch/internal/stats"
)

func snapshot(seen, matched, exported, errors int64, last time.Time) stats.Snapshot {
	return stats.Snapshot{
		Seen:          seen,
		Matched:       matched,
		Exported:      exported,
		Errors:        errors,
		LastMessageAt: last,
	}
}

func TestPrint_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	p := metrics.NewWithWriter(&buf)

	if err := p.Print(snapshot(10, 5, 3, 1, time.Now())); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, header := range []string{"METRIC", "VALUE", "Seen", "Matched", "Exported", "Errors", "Last Message"} {
		if !strings.Contains(out, header) {
			t.Errorf("expected output to contain %q, got:\n%s", header, out)
		}
	}
}

func TestPrint_ZeroLastMessage(t *testing.T) {
	var buf bytes.Buffer
	p := metrics.NewWithWriter(&buf)

	if err := p.Print(snapshot(0, 0, 0, 0, time.Time{})); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "—") {
		t.Errorf("expected em-dash placeholder for zero time, got:\n%s", buf.String())
	}
}

func TestPrintDelta_ContainsRateColumn(t *testing.T) {
	var buf bytes.Buffer
	p := metrics.NewWithWriter(&buf)

	prev := snapshot(0, 0, 0, 0, time.Time{})
	curr := snapshot(100, 50, 20, 2, time.Now())

	if err := p.PrintDelta(prev, curr, time.Second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, col := range []string{"RATE/s", "TOTAL"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected output to contain %q, got:\n%s", col, out)
		}
	}
}

func TestPrintDelta_ZeroInterval(t *testing.T) {
	var buf bytes.Buffer
	p := metrics.NewWithWriter(&buf)

	prev := snapshot(0, 0, 0, 0, time.Time{})
	curr := snapshot(10, 5, 2, 0, time.Now())

	// zero interval should not panic
	if err := p.PrintDelta(prev, curr, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
