package timeline

import (
	"strings"
	"testing"
)

func TestNew_InvalidMaxSize(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for maxSize=0")
	}
}

func TestNew_NegativeMaxSize(t *testing.T) {
	_, err := New(-1)
	if err == nil {
		t.Fatal("expected error for negative maxSize")
	}
}

func TestNew_Valid(t *testing.T) {
	tl, err := New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tl.Len() != 0 {
		t.Errorf("expected empty timeline, got %d entries", tl.Len())
	}
}

func TestAdd_AppendsEntry(t *testing.T) {
	tl, _ := New(10)
	tl.Add("topic-a", "key-1", "hello")
	if tl.Len() != 1 {
		t.Errorf("expected 1 entry, got %d", tl.Len())
	}
}

func TestAdd_EvictsOldestWhenFull(t *testing.T) {
	tl, _ := New(3)
	tl.Add("t", "k1", "first")
	tl.Add("t", "k2", "second")
	tl.Add("t", "k3", "third")
	tl.Add("t", "k4", "fourth")

	if tl.Len() != 3 {
		t.Errorf("expected 3 entries after eviction, got %d", tl.Len())
	}
	entries := tl.Entries()
	if entries[0].Key != "k2" {
		t.Errorf("expected oldest to be evicted, got key=%s", entries[0].Key)
	}
	if entries[2].Key != "k4" {
		t.Errorf("expected newest to be last, got key=%s", entries[2].Key)
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	tl, _ := New(5)
	tl.Add("t", "k", "p")
	snap := tl.Entries()
	snap[0].Key = "mutated"
	original := tl.Entries()
	if original[0].Key == "mutated" {
		t.Error("Entries should return a copy, not a reference")
	}
}

func TestClear_RemovesAllEntries(t *testing.T) {
	tl, _ := New(5)
	tl.Add("t", "k", "p")
	tl.Add("t", "k2", "p2")
	tl.Clear()
	if tl.Len() != 0 {
		t.Errorf("expected 0 entries after clear, got %d", tl.Len())
	}
}

func TestPrint_ContainsHeaders(t *testing.T) {
	tl, _ := New(5)
	tl.Add("my-topic", "key-x", `{"v":1}`)
	var buf strings.Builder
	NewPrinter(&buf).Print(tl)
	out := buf.String()
	for _, h := range []string{"TIME", "TOPIC", "KEY", "PAYLOAD"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestPrint_TruncatesLongPayload(t *testing.T) {
	tl, _ := New(5)
	long := strings.Repeat("x", 80)
	tl.Add("t", "k", long)
	var buf strings.Builder
	NewPrinter(&buf).Print(tl)
	if strings.Contains(buf.String(), long) {
		t.Error("expected long payload to be truncated")
	}
	if !strings.Contains(buf.String(), "...") {
		t.Error("expected truncation indicator '...'")
	}
}
