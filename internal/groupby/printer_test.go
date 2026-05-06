package groupby

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestPrint_EmptyGroups(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	p.Print(nil)
	if !strings.Contains(buf.String(), "no groups") {
		t.Errorf("expected 'no groups' message, got: %s", buf.String())
	}
}

func TestPrint_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	groups := []Group{
		{Key: "topic-a", Seen: 5, Matched: 3, LastSeen: time.Now()},
	}
	p.Print(groups)
	out := buf.String()
	for _, hdr := range []string{"KEY", "SEEN", "MATCHED", "LAST SEEN"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestPrint_SortedBySeen(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	groups := []Group{
		{Key: "low", Seen: 1, LastSeen: time.Now()},
		{Key: "high", Seen: 99, LastSeen: time.Now()},
	}
	p.Print(groups)
	out := buf.String()
	hiIdx := strings.Index(out, "high")
	loIdx := strings.Index(out, "low")
	if hiIdx > loIdx {
		t.Error("expected 'high' to appear before 'low' (sorted by Seen desc)")
	}
}

func TestPrint_KeyInOutput(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	groups := []Group{
		{Key: "my-topic", Seen: 7, Matched: 2, LastSeen: time.Now()},
	}
	p.Print(groups)
	if !strings.Contains(buf.String(), "my-topic") {
		t.Error("expected key 'my-topic' in output")
	}
}
