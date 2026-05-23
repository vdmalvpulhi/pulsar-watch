package timeline

import (
	"errors"
	"sync"
	"time"
)

// Entry represents a single timestamped event recorded in the timeline.
type Entry struct {
	Topic     string
	Key       string
	Payload   string
	RecordedAt time.Time
}

// Timeline tracks message events in chronological order with a bounded capacity.
type Timeline struct {
	mu      sync.RWMutex
	entries []Entry
	maxSize int
}

// New creates a new Timeline with the given maximum number of entries.
// Returns an error if maxSize is less than 1.
func New(maxSize int) (*Timeline, error) {
	if maxSize < 1 {
		return nil, errors.New("timeline: maxSize must be at least 1")
	}
	return &Timeline{
		entries: make([]Entry, 0, maxSize),
		maxSize: maxSize,
	}, nil
}

// Add appends an entry to the timeline. If the timeline is full, the oldest
// entry is evicted to make room for the new one.
func (t *Timeline) Add(topic, key, payload string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	entry := Entry{
		Topic:      topic,
		Key:        key,
		Payload:    payload,
		RecordedAt: time.Now(),
	}
	if len(t.entries) >= t.maxSize {
		t.entries = t.entries[1:]
	}
	t.entries = append(t.entries, entry)
}

// Entries returns a snapshot of all current entries in chronological order.
func (t *Timeline) Entries() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	snap := make([]Entry, len(t.entries))
	copy(snap, t.entries)
	return snap
}

// Len returns the current number of entries in the timeline.
func (t *Timeline) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.entries)
}

// Clear removes all entries from the timeline.
func (t *Timeline) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = t.entries[:0]
}
