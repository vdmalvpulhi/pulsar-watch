package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/pulsar-watch/internal/retry"
)

var baseCfg = retry.Config{
	MaxAttempts:  3,
	InitialDelay: time.Millisecond,
	MaxDelay:     10 * time.Millisecond,
	Multiplier:   2.0,
}

func TestNew_InvalidMaxAttempts(t *testing.T) {
	cfg := baseCfg
	cfg.MaxAttempts = 0
	_, err := retry.New(cfg)
	if err == nil {
		t.Fatal("expected error for MaxAttempts=0")
	}
}

func TestNew_NegativeInitialDelay(t *testing.T) {
	cfg := baseCfg
	cfg.InitialDelay = -1
	_, err := retry.New(cfg)
	if err == nil {
		t.Fatal("expected error for negative InitialDelay")
	}
}

func TestNew_MaxDelayLessThanInitial(t *testing.T) {
	cfg := baseCfg
	cfg.MaxDelay = 0
	_, err := retry.New(cfg)
	if err == nil {
		t.Fatal("expected error when MaxDelay < InitialDelay")
	}
}

func TestNew_MultiplierLessThanOne(t *testing.T) {
	cfg := baseCfg
	cfg.Multiplier = 0.5
	_, err := retry.New(cfg)
	if err == nil {
		t.Fatal("expected error for Multiplier < 1")
	}
}

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	r, _ := retry.New(baseCfg)
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	r, _ := retry.New(baseCfg)
	calls := 0
	sentinel := errors.New("boom")
	err := r.Do(context.Background(), func() error {
		calls++
		if calls < 3 {
			return sentinel
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	r, _ := retry.New(baseCfg)
	sentinel := errors.New("always fails")
	err := r.Do(context.Background(), func() error { return sentinel })
	if !errors.Is(err, retry.ErrMaxAttempts) {
		t.Fatalf("expected ErrMaxAttempts, got %v", err)
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected wrapped sentinel error, got %v", err)
	}
}

func TestDo_ContextCancelled(t *testing.T) {
	r, _ := retry.New(baseCfg)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := r.Do(ctx, func() error { return nil })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
