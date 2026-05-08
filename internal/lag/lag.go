package lag

import (
	"errors"
	"sync"
	"time"
)

// Entry holds lag information for a single topic/subscription pair.
type Entry struct {
	Topic        string
	Subscription string
	Backlog      int64
	LastUpdated  time.Time
}

// Tracker tracks consumer lag (backlog) per topic and subscription.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New creates a new Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]*Entry),
	}
}

func key(topic, subscription string) string {
	return topic + "|" + subscription
}

// Record updates the backlog for a topic/subscription pair.
func (t *Tracker) Record(topic, subscription string, backlog int64) error {
	if topic == "" {
		return errors.New("lag: topic must not be empty")
	}
	if subscription == "" {
		return errors.New("lag: subscription must not be empty")
	}
	if backlog < 0 {
		return errors.New("lag: backlog must not be negative")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[key(topic, subscription)] = &Entry{
		Topic:        topic,
		Subscription: subscription,
		Backlog:      backlog,
		LastUpdated:  time.Now(),
	}
	return nil
}

// Get returns the lag entry for a topic/subscription pair.
func (t *Tracker) Get(topic, subscription string) (Entry, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[key(topic, subscription)]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// Snapshot returns all current lag entries.
func (t *Tracker) Snapshot() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, *e)
	}
	return out
}

// Reset clears all tracked entries.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[string]*Entry)
}
