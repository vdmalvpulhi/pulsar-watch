package timeline

import (
	"strings"
	"testing"
)

func TestNewPrinter_NilWriter(t *testing.T) {
	// Should not panic when nil writer is provided (falls back to os.Stdout).
	p := NewPrinter(nil)
	if p == nil {
		t.Fatal("expected non-nil printer")
	}
}

func TestPrint_EmptyTimeline(t *testing.T) {
	tl, _ := New(10)
	var buf strings.Builder
	NewPrinter(&buf).Print(tl)
	out := buf.String()
	if !strings.Contains(out, "TIME") {
		t.Error("expected header row even for empty timeline")
	}
}

func TestPrint_EntryValues(t *testing.T) {
	tl, _ := New(10)
	tl.Add("persistent://public/default/orders", "order-99", `{"amount":42}`)
	var buf strings.Builder
	NewPrinter(&buf).Print(tl)
	out := buf.String()
	if !strings.Contains(out, "persistent://public/default/orders") {
		t.Error("expected topic in output")
	}
	if !strings.Contains(out, "order-99") {
		t.Error("expected key in output")
	}
	if !strings.Contains(out, `{"amount":42}`) {
		t.Error("expected payload in output")
	}
}

func TestPrint_MultipleEntries(t *testing.T) {
	tl, _ := New(10)
	tl.Add("t", "k1", "p1")
	tl.Add("t", "k2", "p2")
	tl.Add("t", "k3", "p3")
	var buf strings.Builder
	NewPrinter(&buf).Print(tl)
	out := buf.String()
	for _, key := range []string{"k1", "k2", "k3"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %q in output", key)
		}
	}
}
