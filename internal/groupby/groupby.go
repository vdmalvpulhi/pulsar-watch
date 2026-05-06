package groupby

import (
	"errors"
	"sync"
	"time"
)

// Group holds aggregated message counts for a key.
type Group struct {
	Key      string
	Seen     int64
	Matched  int64
	LastSeen time.Time
}

// Grouper aggregates message statistics by a chosen field key.
type Grouper struct {
	mu     sync.Mutex
	groups map[string]*Group
	maxKeys int
}

// New creates a Grouper with the given maximum number of tracked keys.
// maxKeys must be greater than zero.
func New(maxKeys int) (*Grouper, error) {
	if maxKeys <= 0 {
		return nil, errors.New("groupby: maxKeys must be greater than zero")
	}
	return &Grouper{
		groups:  make(map[string]*Group, maxKeys),
		maxKeys: maxKeys,
	}, nil
}

// RecordSeen increments the seen counter for key.
// If the key is new and the max capacity is reached the call is a no-op.
func (g *Grouper) RecordSeen(key string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	grp := g.getOrCreate(key)
	if grp == nil {
		return
	}
	grp.Seen++
	grp.LastSeen = time.Now()
}

// RecordMatched increments the matched counter for key.
func (g *Grouper) RecordMatched(key string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	grp := g.getOrCreate(key)
	if grp == nil {
		return
	}
	grp.Matched++
	grp.LastSeen = time.Now()
}

// Snapshot returns a copy of all groups.
func (g *Grouper) Snapshot() []Group {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make([]Group, 0, len(g.groups))
	for _, grp := range g.groups {
		out = append(out, *grp)
	}
	return out
}

// Reset clears all accumulated groups.
func (g *Grouper) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.groups = make(map[string]*Group, g.maxKeys)
}

// getOrCreate returns an existing group or creates one if capacity allows.
// Caller must hold g.mu.
func (g *Grouper) getOrCreate(key string) *Group {
	if grp, ok := g.groups[key]; ok {
		return grp
	}
	if len(g.groups) >= g.maxKeys {
		return nil
	}
	grp := &Group{Key: key}
	g.groups[key] = grp
	return grp
}
