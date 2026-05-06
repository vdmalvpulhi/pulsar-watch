// Package window provides a sliding time-window counter for tracking
// message rates and burst activity over a configurable duration.
package window

import (
	"errors"
	"sync"
	"time"
)

// Window is a sliding time-window that counts events within a duration.
type Window struct {
	mu       sync.Mutex
	buckets  []int64
	size     int
	duration time.Duration
	bucket   time.Duration
	last     time.Time
}

// New creates a new sliding Window with the given duration and number of buckets.
// Returns an error if duration is non-positive or buckets is less than 1.
func New(duration time.Duration, buckets int) (*Window, error) {
	if duration <= 0 {
		return nil, errors.New("window: duration must be positive")
	}
	if buckets < 1 {
		return nil, errors.New("window: buckets must be at least 1")
	}
	return &Window{
		buckets:  make([]int64, buckets),
		size:     buckets,
		duration: duration,
		bucket:   duration / time.Duration(buckets),
		last:     time.Now(),
	}, nil
}

// Add records n events at the current time.
func (w *Window) Add(n int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.advance(time.Now())
	w.buckets[0] += n
}

// Count returns the total number of events recorded within the window duration.
func (w *Window) Count() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.advance(time.Now())
	var total int64
	for _, v := range w.buckets {
		total += v
	}
	return total
}

// Rate returns the average events-per-second over the window duration.
func (w *Window) Rate() float64 {
	count := w.Count()
	return float64(count) / w.duration.Seconds()
}

// Reset clears all bucket counts.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for i := range w.buckets {
		w.buckets[i] = 0
	}
	w.last = time.Now()
}

// advance shifts buckets forward based on elapsed time since last update.
func (w *Window) advance(now time.Time) {
	elapsed := now.Sub(w.last)
	if elapsed < w.bucket {
		return
	}
	shift := int(elapsed / w.bucket)
	if shift >= w.size {
		for i := range w.buckets {
			w.buckets[i] = 0
		}
	} else {
		for i := w.size - 1; i >= shift; i-- {
			w.buckets[i] = w.buckets[i-shift]
		}
		for i := 0; i < shift; i++ {
			w.buckets[i] = 0
		}
	}
	w.last = now
}
