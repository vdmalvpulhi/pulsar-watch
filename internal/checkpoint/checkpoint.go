// Package checkpoint provides persistent tracking of the last processed
// message ID for a given topic, enabling resumable consumption across
// pulsar-watch restarts.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// Entry holds the checkpoint data for a single topic.
type Entry struct {
	Topic     string    `json:"topic"`
	MessageID string    `json:"message_id"`
	SavedAt   time.Time `json:"saved_at"`
}

// Checkpoint manages persistent message-ID checkpoints keyed by topic.
type Checkpoint struct {
	mu      sync.RWMutex
	path    string
	entries map[string]Entry
}

// New loads (or creates) a checkpoint file at path.
func New(path string) (*Checkpoint, error) {
	if path == "" {
		return nil, errors.New("checkpoint: path must not be empty")
	}
	cp := &Checkpoint{
		path:    path,
		entries: make(map[string]Entry),
	}
	if err := cp.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return cp, nil
}

// Set stores the latest message ID for the given topic and flushes to disk.
func (c *Checkpoint) Set(topic, messageID string) error {
	if topic == "" {
		return errors.New("checkpoint: topic must not be empty")
	}
	c.mu.Lock()
	c.entries[topic] = Entry{Topic: topic, MessageID: messageID, SavedAt: time.Now().UTC()}
	c.mu.Unlock()
	return c.flush()
}

// Get returns the stored message ID for topic, or an empty string if none exists.
func (c *Checkpoint) Get(topic string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[topic]
	if !ok {
		return "", false
	}
	return e.MessageID, true
}

// Delete removes the checkpoint entry for topic and flushes to disk.
func (c *Checkpoint) Delete(topic string) error {
	c.mu.Lock()
	delete(c.entries, topic)
	c.mu.Unlock()
	return c.flush()
}

func (c *Checkpoint) load() error {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}
	for _, e := range entries {
		c.entries[e.Topic] = e
	}
	return nil
}

func (c *Checkpoint) flush() error {
	c.mu.RLock()
	list := make([]Entry, 0, len(c.entries))
	for _, e := range c.entries {
		list = append(list, e)
	}
	c.mu.RUnlock()
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0o644)
}
