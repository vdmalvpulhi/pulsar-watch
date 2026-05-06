// Package throttle provides adaptive message throughput control
// based on observed processing rates and configurable thresholds.
package throttle

import (
	"errors"
	"sync"
	"time"
)

// Throttle tracks message throughput and signals when processing
// should be slowed based on a configured messages-per-second ceiling.
type Throttle struct {
	mu       sync.Mutex
	limit    float64
	window   time.Duration
	bucket   []time.Time
	paused   bool
}

// New creates a Throttle with the given messages-per-second limit and
// sliding-window duration. limit must be positive; window must be > 0.
func New(limit float64, window time.Duration) (*Throttle, error) {
	if limit <= 0 {
		return nil, errors.New("throttle: limit must be positive")
	}
	if window <= 0 {
		return nil, errors.New("throttle: window must be positive")
	}
	return &Throttle{
		limit:  limit,
		window: window,
		bucket: make([]time.Time, 0, 64),
	}, nil
}

// Record registers a message arrival and returns true when the current
// rate exceeds the configured limit, indicating the caller should pause.
func (t *Throttle) Record(now time.Time) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := now.Add(-t.window)
	valid := t.bucket[:0]
	for _, ts := range t.bucket {
		if ts.After(cutoff) {
			valid = append(valid, ts)
		}
	}
	valid = append(valid, now)
	t.bucket = valid

	rate := float64(len(t.bucket)) / t.window.Seconds()
	t.paused = rate > t.limit
	return t.paused
}

// Throttled returns true if the most recent Record call exceeded the limit.
func (t *Throttle) Throttled() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.paused
}

// Rate returns the current observed messages-per-second within the window.
func (t *Throttle) Rate(now time.Time) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := now.Add(-t.window)
	count := 0
	for _, ts := range t.bucket {
		if ts.After(cutoff) {
			count++
		}
	}
	return float64(count) / t.window.Seconds()
}

// Reset clears all recorded timestamps and resets the throttled state.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.bucket = t.bucket[:0]
	t.paused = false
}
