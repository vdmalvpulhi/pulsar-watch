// Package partition provides utilities for tracking and balancing
// message consumption across Pulsar topic partitions.
package partition

import (
	"errors"
	"fmt"
	"sync"
)

// Stats holds per-partition message counters.
type Stats struct {
	Partition int
	Seen      int64
	Matched   int64
}

// Tracker records message activity per partition index.
type Tracker struct {
	mu   sync.RWMutex
	data map[int]*Stats
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{
		data: make(map[int]*Stats),
	}
}

// Record increments the seen counter for the given partition.
// If matched is true the matched counter is also incremented.
// Returns an error when partition is negative.
func (t *Tracker) Record(partition int, matched bool) error {
	if partition < 0 {
		return errors.New("partition index must be non-negative")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	s, ok := t.data[partition]
	if !ok {
		s = &Stats{Partition: partition}
		t.data[partition] = s
	}
	s.Seen++
	if matched {
		s.Matched++
	}
	return nil
}

// Snapshot returns a copy of all partition stats sorted by partition index.
func (t *Tracker) Snapshot() []Stats {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Stats, 0, len(t.data))
	for _, s := range t.data {
		out = append(out, *s)
	}
	// simple insertion sort — partition count is typically small
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].Partition < out[j-1].Partition; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// Summary returns a human-readable one-line summary of all partitions.
func (t *Tracker) Summary() string {
	snap := t.Snapshot()
	if len(snap) == 0 {
		return "no partitions recorded"
	}
	var total, matched int64
	for _, s := range snap {
		total += s.Seen
		matched += s.Matched
	}
	return fmt.Sprintf("partitions=%d total_seen=%d total_matched=%d",
		len(snap), total, matched)
}

// Reset clears all recorded data.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.data = make(map[int]*Stats)
}
