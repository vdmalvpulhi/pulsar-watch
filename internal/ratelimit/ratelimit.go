// Package ratelimit provides a simple token-bucket rate limiter for
// controlling the rate at which Pulsar messages are processed or replayed.
package ratelimit

import (
	"context"
	"fmt"
	"time"
)

// Limiter controls the rate of message processing using a token bucket.
type Limiter struct {
	rate     int
	ticker   *time.Ticker
	tokens   chan struct{}
	stopCh   chan struct{}
}

// New creates a Limiter that allows up to ratePerSecond messages per second.
// A rate of 0 means unlimited.
func New(ratePerSecond int) (*Limiter, error) {
	if ratePerSecond < 0 {
		return nil, fmt.Errorf("ratelimit: rate must be >= 0, got %d", ratePerSecond)
	}
	l := &Limiter{
		rate:   ratePerSecond,
		stopCh: make(chan struct{}),
	}
	if ratePerSecond > 0 {
		interval := time.Second / time.Duration(ratePerSecond)
		l.ticker = time.NewTicker(interval)
		l.tokens = make(chan struct{}, ratePerSecond)
		go l.produce()
	}
	return l, nil
}

// Wait blocks until a token is available or the context is cancelled.
// If the limiter is unlimited (rate == 0), Wait returns immediately.
func (l *Limiter) Wait(ctx context.Context) error {
	if l.rate == 0 {
		return nil
	}
	select {
	case <-l.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop releases resources held by the Limiter.
func (l *Limiter) Stop() {
	if l.ticker != nil {
		l.ticker.Stop()
		close(l.stopCh)
	}
}

// Rate returns the configured rate per second.
func (l *Limiter) Rate() int {
	return l.rate
}

func (l *Limiter) produce() {
	for {
		select {
		case <-l.ticker.C:
			select {
			case l.tokens <- struct{}{}:
			default:
				// bucket full, drop token
			}
		case <-l.stopCh:
			return
		}
	}
}
