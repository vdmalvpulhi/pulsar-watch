// Package cursor manages persistent read position tracking for Pulsar topics.
package cursor

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Position holds the last successfully processed message position for a topic.
type Position struct {
	Topic       string    `json:"topic"`
	MessageID   string    `json:"message_id"`
	PublishTime time.Time `json:"publish_time"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Store persists and retrieves cursor positions.
type Store struct {
	mu   sync.RWMutex
	path string
	data map[string]Position
}

// New creates a new Store backed by the given file path.
// If the file does not yet exist, an empty store is returned without error.
func New(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("cursor: store path must not be empty")
	}
	s := &Store{
		path: path,
		data: make(map[string]Position),
	}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

// Get returns the stored Position for a topic, and whether it exists.
func (s *Store) Get(topic string) (Position, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.data[topic]
	return p, ok
}

// Set updates the Position for a topic and persists to disk.
func (s *Store) Set(pos Position) error {
	if pos.Topic == "" {
		return errors.New("cursor: position topic must not be empty")
	}
	pos.UpdatedAt = time.Now().UTC()
	s.mu.Lock()
	s.data[pos.Topic] = pos
	s.mu.Unlock()
	return s.save()
}

// Delete removes the stored position for a topic and persists to disk.
// If the topic does not exist, Delete is a no-op and returns nil.
func (s *Store) Delete(topic string) error {
	s.mu.Lock()
	_, exists := s.data[topic]
	if exists {
		delete(s.data, topic)
	}
	s.mu.Unlock()
	if !exists {
		return nil
	}
	return s.save()
}

// Topics returns a slice of all topic names that have a stored position.
func (s *Store) Topics() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	topics := make([]string, 0, len(s.data))
	for t := range s.data {
		topics = append(topics, t)
	}
	return topics
}

func (s *Store) load() error {
	b, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return json.Unmarshal(b, &s.data)
}

func (s *Store) save() error {
	s.mu.RLock()
	b, err := json.MarshalIndent(s.data, "", "  ")
	s.mu.RUnlock()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}
