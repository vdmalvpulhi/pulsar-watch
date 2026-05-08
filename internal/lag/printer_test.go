package lag

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestPrint_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	p.Print(nil)
	if !strings.Contains(buf.String(), "no lag data") {
		t.Fatalf("expected no-data message, got: %q", buf.String())
	}
}

func TestPrint_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	p.Print([]Entry{
		{Topic: "t1", Subscription: "s1", Backlog: 5, LastUpdated: time.Now()},
	})
	out := buf.String()
	for _, h := range []string{"TOPIC", "SUBSCRIPTION", "BACKLOG", "LAST UPDATED"} {
		if !strings.Contains(out, h) {
			t.Fatalf("expected header %q in output: %q", h, out)
		}
	}
}

func TestPrint_BacklogValue(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	p.Print([]Entry{
		{Topic: "orders", Subscription: "worker", Backlog: 123, LastUpdated: time.Now()},
	})
	out := buf.String()
	if !strings.Contains(out, "123") {
		t.Fatalf("expected backlog value 123 in output: %q", out)
	}
	if !strings.Contains(out, "orders") {
		t.Fatalf("expected topic name in output: %q", out)
	}
}

func TestPrint_SortedByTopic(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf)
	p.Print([]Entry{
		{Topic: "zzz", Subscription: "s", Backlog: 1, LastUpdated: time.Now()},
		{Topic: "aaa", Subscription: "s", Backlog: 2, LastUpdated: time.Now()},
	})
	out := buf.String()
	idxA := strings.Index(out, "aaa")
	idxZ := strings.Index(out, "zzz")
	if idxA > idxZ {
		t.Fatal("expected entries sorted by topic ascending")
	}
}
