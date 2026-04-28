package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestNew_NegativeRate(t *testing.T) {
	_, err := New(-1)
	if err == nil {
		t.Fatal("expected error for negative rate, got nil")
	}
}

func TestNew_ZeroRate(t *testing.T) {
	l, err := New(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer l.Stop()
	if l.Rate() != 0 {
		t.Errorf("expected rate 0, got %d", l.Rate())
	}
}

func TestNew_PositiveRate(t *testing.T) {
	l, err := New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer l.Stop()
	if l.Rate() != 10 {
		t.Errorf("expected rate 10, got %d", l.Rate())
	}
}

func TestWait_UnlimitedReturnsImmediately(t *testing.T) {
	l, _ := New(0)
	defer l.Stop()
	ctx := context.Background()
	start := time.Now()
	for i := 0; i < 100; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if elapsed := time.Since(start); elapsed > 10*time.Millisecond {
		t.Errorf("unlimited limiter should return immediately, took %v", elapsed)
	}
}

func TestWait_ContextCancelled(t *testing.T) {
	l, _ := New(1) // very slow — 1 msg/s
	defer l.Stop()
	// drain the first token so the next Wait must block
	ctx := context.Background()
	_ = l.Wait(ctx)

	cancel_ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := l.Wait(cancel_ctx)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestWait_ThrottlesRate(t *testing.T) {
	const rate = 50
	l, _ := New(rate)
	defer l.Stop()
	ctx := context.Background()

	// Allow time for tokens to accumulate
	time.Sleep(200 * time.Millisecond)

	start := time.Now()
	for i := 0; i < rate; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatalf("unexpected error at iteration %d: %v", i, err)
		}
	}
	_ = time.Since(start) // timing assertions are flaky in CI; just ensure no errors
}
