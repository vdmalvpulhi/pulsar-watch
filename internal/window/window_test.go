package window

import (
	"testing"
	"time"
)

func TestNew_NegativeDuration(t *testing.T) {
	_, err := New(-time.Second, 10)
	if err == nil {
		t.Fatal("expected error for negative duration")
	}
}

func TestNew_ZeroDuration(t *testing.T) {
	_, err := New(0, 10)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestNew_ZeroBuckets(t *testing.T) {
	_, err := New(time.Second, 0)
	if err == nil {
		t.Fatal("expected error for zero buckets")
	}
}

func TestNew_Valid(t *testing.T) {
	w, err := New(10*time.Second, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil window")
	}
}

func TestAdd_IncrementsCount(t *testing.T) {
	w, _ := New(10*time.Second, 10)
	w.Add(5)
	w.Add(3)
	if got := w.Count(); got != 8 {
		t.Fatalf("expected 8, got %d", got)
	}
}

func TestCount_EmptyWindow(t *testing.T) {
	w, _ := New(10*time.Second, 10)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestReset_ClearsCount(t *testing.T) {
	w, _ := New(10*time.Second, 10)
	w.Add(42)
	w.Reset()
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestRate_NonZero(t *testing.T) {
	w, _ := New(10*time.Second, 10)
	w.Add(100)
	rate := w.Rate()
	if rate <= 0 {
		t.Fatalf("expected positive rate, got %f", rate)
	}
}

func TestAdvance_ShiftsOldBuckets(t *testing.T) {
	w, _ := New(100*time.Millisecond, 10)
	w.Add(10)
	// Simulate time passing beyond the full window
	w.last = w.last.Add(-200 * time.Millisecond)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after window expiry, got %d", got)
	}
}

func TestAdd_ConcurrentSafe(t *testing.T) {
	w, _ := New(10*time.Second, 10)
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			w.Add(1)
			done <- struct{}{}
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	if got := w.Count(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}
