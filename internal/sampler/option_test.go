package sampler

import (
	"testing"
)

func TestWithSeed_Deterministic(t *testing.T) {
	const seed = 42
	const n = 100

	s1, err := NewWithOptions(0.5, WithSeed(seed))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s2, err := NewWithOptions(0.5, WithSeed(seed))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 0; i < n; i++ {
		a := s1.Sample()
		b := s2.Sample()
		if a != b {
			t.Fatalf("samplers with same seed diverged at iteration %d", i)
		}
	}
}

func TestNewWithOptions_InvalidRate(t *testing.T) {
	_, err := NewWithOptions(2.0, WithSeed(1))
	if err == nil {
		t.Fatal("expected error for invalid rate")
	}
}

func TestNewWithOptions_NoOptions(t *testing.T) {
	s, err := NewWithOptions(1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Sample() {
		t.Fatal("expected sample to return true at rate 1.0")
	}
}
