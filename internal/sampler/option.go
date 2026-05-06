package sampler

import "math/rand"

// Option is a functional option for Sampler.
type Option func(*Sampler)

// WithSeed sets a deterministic seed for the internal random number generator.
// Useful in tests or reproducible replay scenarios.
func WithSeed(seed int64) Option {
	return func(s *Sampler) {
		s.rng = rand.New(rand.NewSource(seed))
	}
}

// NewWithOptions creates a Sampler with the given rate and applies any
// provided options. Rate must be in [0.0, 1.0].
func NewWithOptions(rate float64, opts ...Option) (*Sampler, error) {
	s, err := New(rate)
	if err != nil {
		return nil, err
	}
	for _, o := range opts {
		o(s)
	}
	return s, nil
}
