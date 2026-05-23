// Package retry provides a simple exponential back-off retry helper
// used when transient errors occur during message consumption or export.
package retry

import (
	"context"
	"errors"
	"time"
)

// ErrMaxAttempts is returned when all retry attempts have been exhausted.
var ErrMaxAttempts = errors.New("retry: max attempts reached")

// Config holds the parameters that control retry behaviour.
type Config struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// InitialDelay is the wait time before the second attempt.
	InitialDelay time.Duration
	// MaxDelay caps the exponential back-off ceiling.
	MaxDelay time.Duration
	// Multiplier is applied to the delay after each failure (e.g. 2.0).
	Multiplier float64
}

// Retryer executes a function with exponential back-off.
type Retryer struct {
	cfg Config
}

// New returns a Retryer or an error when cfg is invalid.
func New(cfg Config) (*Retryer, error) {
	if cfg.MaxAttempts < 1 {
		return nil, errors.New("retry: MaxAttempts must be >= 1")
	}
	if cfg.InitialDelay < 0 {
		return nil, errors.New("retry: InitialDelay must be >= 0")
	}
	if cfg.MaxDelay < cfg.InitialDelay {
		return nil, errors.New("retry: MaxDelay must be >= InitialDelay")
	}
	if cfg.Multiplier < 1 {
		return nil, errors.New("retry: Multiplier must be >= 1")
	}
	return &Retryer{cfg: cfg}, nil
}

// Do calls fn up to MaxAttempts times. It stops early when ctx is cancelled
// or fn returns nil. If all attempts fail it returns ErrMaxAttempts wrapping
// the last error.
func (r *Retryer) Do(ctx context.Context, fn func() error) error {
	delay := r.cfg.InitialDelay
	var lastErr error

	for attempt := 0; attempt < r.cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if attempt == r.cfg.MaxAttempts-1 {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
		delay = time.Duration(float64(delay) * r.cfg.Multiplier)
		if delay > r.cfg.MaxDelay {
			delay = r.cfg.MaxDelay
		}
	}
	return errors.Join(ErrMaxAttempts, lastErr)
}
