package dedupe

import (
	"fmt"
	"testing"
	"time"
)

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(0, 100)
	if err != ErrInvalidWindow {
		t.Fatalf("expected ErrInvalidWindow, got %v", err)
	}
}

func TestNew_InvalidMaxSize(t *testing.T) {
	_, err := New(time.Minute, 0)
	if err != ErrInvalidMaxSize {
		t.Fatalf("expected ErrInvalidMaxSize, got %v", err)
	}
}

func TestNew_Valid(t *testing.T) {
	d, err := New(time.Minute, 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Size() != 0 {
		t.Fatalf("expected empty deduplicator, got size %d", d.Size())
	}
}

func TestIsDuplicate_FirstSeen(t *testing.T) {
	d, _ := New(time.Minute, 100)
	if d.IsDuplicate("msg-1") {
		t.Fatal("expected false for first-seen message")
	}
	if d.Size() != 1 {
		t.Fatalf("expected size 1, got %d", d.Size())
	}
}

func TestIsDuplicate_Duplicate(t *testing.T) {
	d, _ := New(time.Minute, 100)
	d.IsDuplicate("msg-1")
	if !d.IsDuplicate("msg-1") {
		t.Fatal("expected true for duplicate message")
	}
}

func TestIsDuplicate_MaxSizeEvictsOldest(t *testing.T) {
	d, _ := New(time.Hour, 3)
	for i := 0; i < 3; i++ {
		d.IsDuplicate(fmt.Sprintf("msg-%d", i))
	}
	if d.Size() != 3 {
		t.Fatalf("expected size 3, got %d", d.Size())
	}
	// Adding a fourth should evict one.
	d.IsDuplicate("msg-new")
	if d.Size() != 3 {
		t.Fatalf("expected size to remain 3 after eviction, got %d", d.Size())
	}
}

func TestIsDuplicate_ExpiredEntryNotDuplicate(t *testing.T) {
	d, _ := New(1*time.Millisecond, 100)
	d.IsDuplicate("msg-expire")
	time.Sleep(5 * time.Millisecond)
	// After window expires the entry should be evicted and not considered duplicate.
	if d.IsDuplicate("msg-expire") {
		t.Fatal("expected false after window expiry")
	}
}

func TestIsDuplicate_ConcurrentSafe(t *testing.T) {
	d, _ := New(time.Minute, 1000)
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func(n int) {
			for j := 0; j < 50; j++ {
				d.IsDuplicate(fmt.Sprintf("msg-%d-%d", n, j))
			}
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}
