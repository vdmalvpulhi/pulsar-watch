// Package topicmap tracks per-topic message counts and metadata
// across multiple topics observed during a watch session.
package topicmap

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// Entry holds aggregated stats for a single topic.
type Entry struct {
	Topic     string
	Seen      int64
	Matched   int64
	LastSeen  time.Time
}

// TopicMap is a concurrent map of topic name to Entry.
type TopicMap struct {
	mu      sync.Mutex
	entries map[string]*Entry
	maxKeys int
}

// New creates a TopicMap. maxKeys limits the number of tracked topics;
// pass 0 for unlimited.
func New(maxKeys int) (*TopicMap, error) {
	if maxKeys < 0 {
		return nil, fmt.Errorf("topicmap: maxKeys must be >= 0, got %d", maxKeys)
	}
	return &TopicMap{
		entries: make(map[string]*Entry),
		maxKeys: maxKeys,
	}, nil
}

// RecordSeen increments the seen counter for topic.
func (t *TopicMap) RecordSeen(topic string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e := t.getOrCreate(topic)
	if e == nil {
		return
	}
	e.Seen++
	e.LastSeen = time.Now()
}

// RecordMatched increments the matched counter for topic.
func (t *TopicMap) RecordMatched(topic string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e := t.getOrCreate(topic)
	if e == nil {
		return
	}
	e.Matched++
}

// Snapshot returns a sorted copy of all entries by Seen descending.
func (t *TopicMap) Snapshot() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, *e)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Seen != out[j].Seen {
			return out[i].Seen > out[j].Seen
		}
		return out[i].Topic < out[j].Topic
	})
	return out
}

// Len returns the number of tracked topics.
func (t *TopicMap) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.entries)
}

// getOrCreate returns the entry for topic, creating it if necessary.
// Returns nil if maxKeys would be exceeded.
func (t *TopicMap) getOrCreate(topic string) *Entry {
	if e, ok := t.entries[topic]; ok {
		return e
	}
	if t.maxKeys > 0 && len(t.entries) >= t.maxKeys {
		return nil
	}
	e := &Entry{Topic: topic}
	t.entries[topic] = e
	return e
}
