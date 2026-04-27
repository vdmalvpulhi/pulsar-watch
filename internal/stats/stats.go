// Package stats provides message consumption statistics tracking
// for pulsar-watch monitoring sessions.
package stats

import (
	"fmt"
	"sync"
	"time"
)

// Stats tracks message consumption metrics during a watch session.
type Stats struct {
	mu           sync.RWMutex
	StartTime    time.Time
	TotalSeen    int64
	TotalMatched int64
	TotalExported int64
	LastMessageAt time.Time
}

// New creates a new Stats instance with the start time set to now.
func New() *Stats {
	return &Stats{
		StartTime: time.Now(),
	}
}

// RecordSeen increments the total messages seen counter.
func (s *Stats) RecordSeen() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalSeen++
	s.LastMessageAt = time.Now()
}

// RecordMatched increments the total messages matched counter.
func (s *Stats) RecordMatched() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalMatched++
}

// RecordExported increments the total messages exported counter.
func (s *Stats) RecordExported() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalExported++
}

// Snapshot returns a consistent read of the current stats.
func (s *Stats) Snapshot() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return Snapshot{
		Duration:      time.Since(s.StartTime),
		TotalSeen:     s.TotalSeen,
		TotalMatched:  s.TotalMatched,
		TotalExported: s.TotalExported,
		LastMessageAt: s.LastMessageAt,
	}
}

// Snapshot holds a point-in-time copy of stats values.
type Snapshot struct {
	Duration      time.Duration
	TotalSeen     int64
	TotalMatched  int64
	TotalExported int64
	LastMessageAt time.Time
}

// String returns a human-readable summary of the snapshot.
func (s Snapshot) String() string {
	matchRate := 0.0
	if s.TotalSeen > 0 {
		matchRate = float64(s.TotalMatched) / float64(s.TotalSeen) * 100
	}
	return fmt.Sprintf(
		"duration=%s seen=%d matched=%d exported=%d match_rate=%.1f%%",
		s.Duration.Round(time.Second),
		s.TotalSeen,
		s.TotalMatched,
		s.TotalExported,
		matchRate,
	)
}
