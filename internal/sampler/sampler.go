package sampler

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

// Sampler decides whether a given message should be processed based on a
// configured sampling rate between 0.0 (drop all) and 1.0 (keep all).
type Sampler struct {
	mu   sync.Mutex
	rng  *rand.Rand
	rate float64
	seen uint64
	kept uint64
}

// New creates a Sampler with the given rate. Rate must be in [0.0, 1.0].
func New(rate float64) (*Sampler, error) {
	if rate < 0.0 || rate > 1.0 {
		return nil, errors.New("sampler: rate must be between 0.0 and 1.0")
	}
	return &Sampler{
		rng:  rand.New(rand.NewSource(time.Now().UnixNano())),
		rate: rate,
	}, nil
}

// Sample returns true if the message should be kept according to the rate.
func (s *Sampler) Sample() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seen++
	if s.rate >= 1.0 {
		s.kept++
		return true
	}
	if s.rate <= 0.0 {
		return false
	}
	if s.rng.Float64() < s.rate {
		s.kept++
		return true
	}
	return false
}

// Stats returns the total seen and kept counts.
func (s *Sampler) Stats() (seen, kept uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.seen, s.kept
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() float64 {
	return s.rate
}
