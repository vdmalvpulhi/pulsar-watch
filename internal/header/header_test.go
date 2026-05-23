package header

import (
	"testing"
)

func TestNew_BothEmpty(t *testing.T) {
	_, err := New("", "")
	if err == nil {
		t.Fatal("expected error for both empty patterns")
	}
}

func TestNew_InvalidKeyPattern(t *testing.T) {
	_, err := New("[invalid", "")
	if err == nil {
		t.Fatal("expected error for invalid key pattern")
	}
}

func TestNew_InvalidValuePattern(t *testing.T) {
	_, err := New("", "[invalid")
	if err == nil {
		t.Fatal("expected error for invalid value pattern")
	}
}

func TestNew_ValidKeyOnly(t *testing.T) {
	e, err := New("^x-", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil extractor")
	}
}

func TestMatch_KeyPatternOnly(t *testing.T) {
	e, _ := New("^x-", "")
	props := map[string]string{"x-trace-id": "abc", "content-type": "json"}
	if !e.Match(props) {
		t.Error("expected match on key prefix x-")
	}
}

func TestMatch_KeyPatternNoMatch(t *testing.T) {
	e, _ := New("^x-", "")
	props := map[string]string{"content-type": "json"}
	if e.Match(props) {
		t.Error("expected no match")
	}
}

func TestMatch_ValuePatternOnly(t *testing.T) {
	e, _ := New("", "json")
	props := map[string]string{"content-type": "application/json"}
	if !e.Match(props) {
		t.Error("expected match on value containing json")
	}
}

func TestMatch_BothPatterns_EntryMustSatisfyBoth(t *testing.T) {
	e, _ := New("^x-", "abc")
	// only x-trace-id has value abc
	props := map[string]string{"x-trace-id": "abc", "content-type": "abc"}
	if !e.Match(props) {
		t.Error("expected match when one entry satisfies both patterns")
	}
}

func TestMatch_BothPatterns_NoEntryMatchesBoth(t *testing.T) {
	e, _ := New("^x-", "abc")
	// x- key has different value; content-type has abc but wrong key
	props := map[string]string{"x-trace-id": "xyz", "content-type": "abc"}
	if e.Match(props) {
		t.Error("expected no match")
	}
}

func TestMatch_EmptyProps(t *testing.T) {
	e, _ := New("^x-", "")
	if e.Match(map[string]string{}) {
		t.Error("expected no match on empty props")
	}
}

func TestExtract_KeyPattern(t *testing.T) {
	e, _ := New("^x-", "")
	props := map[string]string{"x-trace-id": "abc", "content-type": "json"}
	out := e.Extract(props)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out["x-trace-id"] != "abc" {
		t.Errorf("unexpected value: %v", out)
	}
}

func TestExtract_NoKeyPattern_ReturnsAll(t *testing.T) {
	e, _ := New("", "json")
	props := map[string]string{"a": "1", "b": "2"}
	out := e.Extract(props)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}
