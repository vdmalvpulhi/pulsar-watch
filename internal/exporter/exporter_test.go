package exporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/pulsar-watch/internal/exporter"
)

var sampleMsg = exporter.Message{
	Topic:     "persistent://public/default/test",
	Key:       "msg-key-1",
	Payload:   `{"event":"click"}`,
	Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	Headers:   map[string]string{"source": "web"},
}

func TestNewWithWriter_UnsupportedFormat(t *testing.T) {
	_, err := exporter.NewWithWriter("xml", &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	e, err := exporter.NewWithWriter(exporter.FormatJSON, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := e.Write(sampleMsg); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	var got exporter.Message
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got.Key != sampleMsg.Key {
		t.Errorf("key mismatch: got %q, want %q", got.Key, sampleMsg.Key)
	}
	if got.Payload != sampleMsg.Payload {
		t.Errorf("payload mismatch: got %q, want %q", got.Payload, sampleMsg.Payload)
	}
}

func TestWrite_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	e, err := exporter.NewWithWriter(exporter.FormatText, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := e.Write(sampleMsg); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, sampleMsg.Key) {
		t.Errorf("output missing key %q: %s", sampleMsg.Key, output)
	}
	if !strings.Contains(output, sampleMsg.Topic) {
		t.Errorf("output missing topic %q: %s", sampleMsg.Topic, output)
	}
	if !strings.Contains(output, "2024-01-15T10:00:00Z") {
		t.Errorf("output missing timestamp: %s", output)
	}
}

func TestWrite_MultipleMessages(t *testing.T) {
	var buf bytes.Buffer
	e, err := exporter.NewWithWriter(exporter.FormatJSON, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 0; i < 3; i++ {
		if err := e.Write(sampleMsg); err != nil {
			t.Fatalf("Write %d failed: %v", i, err)
		}
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}
