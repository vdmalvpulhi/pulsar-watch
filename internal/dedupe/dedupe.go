// Package dedupe provides message deduplication for Pulsar topic consumers.
// It tracks recently seen message IDs using a fixed-size ring buffer to avoid
// processing duplicate messages during replay or reconnect scenarios.
package dedupe

import (
	"sync"
	"time"
)

// Entry holds a message ID and the time it was first seen.
type Entry struct {
	ID      string
	SeenAt  time.Time
}

// Deduplicator tracks recently seen message IDs within a configurable window.
type Deduplicator struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	window  time.Duration
	maxSize int
}

// New creates a Deduplicator that forgets entries older than window.
// maxSize caps the number of tracked IDs; must be > 0.
func New(window time.Duration, maxSize int) (*Deduplicator, error) {
	if window <= 0 {
		return nil, ErrInvalidWindow
	}
	if maxSize <= 0 {
		return nil, ErrInvalidMaxSize
	}
	return &Deduplicator{
		seen:    make(map[string]time.Time, maxSize),
		window:  window,
		maxSize: maxSize,
	}, nil
}

// IsDuplicate returns true if id was seen within the deduplication window.
// It also records the id if it is new, evicting stale entries as needed.
func (d *Deduplicator) IsDuplicate(id string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	d.evict(now)

	if _, exists := d.seen[id]; exists {
		return true
	}

	// Enforce max size by dropping oldest entry when full.
	if len(d.seen) >= d.maxSize {
		d.dropOldest()
	}

	d.seen[id] = now
	return false
}

// Size returns the current number of tracked IDs.
func (d *Deduplicator) Size() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}

// evict removes entries older than the deduplication window. Must be called with lock held.
func (d *Deduplicator) evict(now time.Time) {
	for id, seenAt := range d.seen {
		if now.Sub(seenAt) > d.window {
			delete(d.seen, id)
		}
	}
}

// dropOldest removes the single oldest entry. Must be called with lock held.
func (d *Deduplicator) dropOldest() {
	var oldestID string
	var oldestTime time.Time
	for id, seenAt := range d.seen {
		if oldestID == "" || seenAt.Before(oldestTime) {
			oldestID = id
			oldestTime = seenAt
		}
	}
	if oldestID != "" {
		delete(d.seen, oldestID)
	}
}
