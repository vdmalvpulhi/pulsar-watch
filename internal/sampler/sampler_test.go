package sampler

import (
	"testing"
)

func TestNew_RateTooLow(t *testing.T) {
	_, err := New(-0.1)
	if err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestNew_RateTooHigh(t *testing.T) {
	_, err := New(1.1)
	if err == nil {
		t.Fatal("expected error for rate > 1.0")
	}
}

func TestNew_ValidRate(t *testing.T) {
	s, err := New(0.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Rate() != 0.5 {
		t.Errorf("expected rate 0.5, got %f", s.Rate())
	}
}

func TestSample_ZeroRateDropsAll(t *testing.T) {
	s, _ := New(0.0)
	for i := 0; i < 100; i++ {
		if s.Sample() {
			t.Fatal("expected all messages dropped at rate 0.0")
		}
	}
	seen, kept := s.Stats()
	if seen != 100 {
		t.Errorf("expected seen=100, got %d", seen)
	}
	if kept != 0 {
		t.Errorf("expected kept=0, got %d", kept)
	}
}

func TestSample_FullRateKeepsAll(t *testing.T) {
	s, _ := New(1.0)
	for i := 0; i < 50; i++ {
		if !s.Sample() {
			t.Fatal("expected all messages kept at rate 1.0")
		}
	}
	seen, kept := s.Stats()
	if seen != 50 || kept != 50 {
		t.Errorf("expected seen=50 kept=50, got seen=%d kept=%d", seen, kept)
	}
}

func TestSample_PartialRate_Approximate(t *testing.T) {
	s, _ := New(0.5)
	const n = 10000
	var accepted int
	for i := 0; i < n; i++ {
		if s.Sample() {
			accepted++
		}
	}
	ratio := float64(accepted) / float64(n)
	if ratio < 0.40 || ratio > 0.60 {
		t.Errorf("expected ~50%% acceptance, got %.2f%%", ratio*100)
	}
}

func TestStats_Consistency(t *testing.T) {
	s, _ := New(0.7)
	for i := 0; i < 200; i++ {
		s.Sample()
	}
	seen, kept := s.Stats()
	if seen != 200 {
		t.Errorf("expected seen=200, got %d", seen)
	}
	if kept > seen {
		t.Errorf("kept (%d) should never exceed seen (%d)", kept, seen)
	}
}
