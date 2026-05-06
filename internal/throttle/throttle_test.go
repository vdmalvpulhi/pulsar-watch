package throttle

import (
	"testing"
	"time"
)

func TestNew_NegativeLimit(t *testing.T) {
	_, err := New(-1, time.Second)
	if err == nil {
		t.Fatal("expected error for negative limit")
	}
}

func TestNew_ZeroLimit(t *testing.T) {
	_, err := New(0, time.Second)
	if err == nil {
		t.Fatal("expected error for zero limit")
	}
}

func TestNew_ZeroWindow(t *testing.T) {
	_, err := New(10, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_Valid(t *testing.T) {
	th, err := New(50, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th == nil {
		t.Fatal("expected non-nil throttle")
	}
}

func TestRecord_BelowLimit(t *testing.T) {
	th, _ := New(10, time.Second)
	now := time.Now()
	for i := 0; i < 5; i++ {
		if th.Record(now) {
			t.Errorf("iteration %d: expected not throttled below limit", i)
		}
	}
}

func TestRecord_ExceedsLimit(t *testing.T) {
	th, _ := New(5, time.Second)
	now := time.Now()
	var throttled bool
	for i := 0; i < 20; i++ {
		if th.Record(now) {
			throttled = true
			break
		}
	}
	if !throttled {
		t.Fatal("expected throttle to trigger when rate exceeds limit")
	}
}

func TestThrottled_ReflectsMostRecentRecord(t *testing.T) {
	th, _ := New(2, time.Second)
	now := time.Now()
	th.Record(now)
	th.Record(now)
	th.Record(now)
	if !th.Throttled() {
		t.Fatal("expected Throttled() to return true after exceeding limit")
	}
}

func TestRate_CountsWithinWindow(t *testing.T) {
	th, _ := New(100, time.Second)
	now := time.Now()
	for i := 0; i < 10; i++ {
		th.Record(now)
	}
	rate := th.Rate(now)
	if rate < 9 || rate > 11 {
		t.Errorf("expected rate ~10, got %.2f", rate)
	}
}

func TestRate_ExcludesExpiredEntries(t *testing.T) {
	th, _ := New(100, time.Second)
	old := time.Now().Add(-2 * time.Second)
	for i := 0; i < 50; i++ {
		th.Record(old)
	}
	now := time.Now()
	th.Record(now)
	rate := th.Rate(now)
	if rate > 5 {
		t.Errorf("expected old entries to be excluded, got rate %.2f", rate)
	}
}

func TestReset_ClearsState(t *testing.T) {
	th, _ := New(2, time.Second)
	now := time.Now()
	th.Record(now)
	th.Record(now)
	th.Record(now)
	th.Reset()
	if th.Throttled() {
		t.Fatal("expected Throttled() to return false after Reset")
	}
	if r := th.Rate(now); r != 0 {
		t.Errorf("expected rate 0 after Reset, got %.2f", r)
	}
}
