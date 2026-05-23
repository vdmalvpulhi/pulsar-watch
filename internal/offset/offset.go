// Package offset tracks per-topic message offsets (sequence IDs) for
// resumable consumption across pulsar-watch sessions.
package offset

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// Tracker persists and retrieves the last-seen sequence ID for each topic.
type Tracker struct {
	mu   sync.RWMutex
	path string
	data map[string]uint64
}

// New creates a Tracker backed by the file at path.
// If the file does not exist it is created on the first Save call.
func New(path string) (*Tracker, error) {
	if path == "" {
		return nil, errors.New("offset: path must not be empty")
	}
	t := &Tracker{path: path, data: make(map[string]uint64)}
	if err := t.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return t, nil
}

// Set records sequenceID as the latest offset for topic.
func (t *Tracker) Set(topic string, sequenceID uint64) error {
	if topic == "" {
		return errors.New("offset: topic must not be empty")
	}
	t.mu.Lock()
	t.data[topic] = sequenceID
	t.mu.Unlock()
	return nil
}

// Get returns the last-saved sequence ID for topic.
// The second return value is false when no offset has been recorded.
func (t *Tracker) Get(topic string) (uint64, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	v, ok := t.data[topic]
	return v, ok
}

// Save flushes all offsets to disk.
func (t *Tracker) Save() error {
	t.mu.RLock()
	defer t.mu.RUnlock()
	b, err := json.Marshal(t.data)
	if err != nil {
		return err
	}
	return os.WriteFile(t.path, b, 0o600)
}

// load reads offsets from disk into memory.
func (t *Tracker) load() error {
	b, err := os.ReadFile(t.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &t.data)
}
