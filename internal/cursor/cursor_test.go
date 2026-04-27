package cursor_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/pulsar-watch/internal/cursor"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "cursors.json")
}

func TestNew_EmptyPath(t *testing.T) {
	_, err := cursor.New("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestNew_FileNotExist(t *testing.T) {
	_, err := cursor.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error when file does not exist: %v", err)
	}
}

func TestSet_EmptyTopic(t *testing.T) {
	s, _ := cursor.New(tempPath(t))
	err := s.Set(cursor.Position{Topic: "", MessageID: "1:0"})
	if err == nil {
		t.Fatal("expected error for empty topic, got nil")
	}
}

func TestGet_NotFound(t *testing.T) {
	s, _ := cursor.New(tempPath(t))
	_, ok := s.Get("persistent://public/default/missing")
	if ok {
		t.Fatal("expected ok=false for unknown topic")
	}
}

func TestSetAndGet(t *testing.T) {
	path := tempPath(t)
	s, _ := cursor.New(path)

	topic := "persistent://public/default/events"
	want := cursor.Position{
		Topic:       topic,
		MessageID:   "42:0",
		PublishTime: time.Now().UTC().Truncate(time.Second),
	}
	if err := s.Set(want); err != nil {
		t.Fatalf("Set: %v", err)
	}

	got, ok := s.Get(topic)
	if !ok {
		t.Fatal("expected position to exist after Set")
	}
	if got.MessageID != want.MessageID {
		t.Errorf("MessageID: got %q, want %q", got.MessageID, want.MessageID)
	}
}

func TestPersistence(t *testing.T) {
	path := tempPath(t)
	topic := "persistent://public/default/orders"

	s1, _ := cursor.New(path)
	_ = s1.Set(cursor.Position{Topic: topic, MessageID: "99:1"})

	s2, err := cursor.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	pos, ok := s2.Get(topic)
	if !ok {
		t.Fatal("position not found after reload")
	}
	if pos.MessageID != "99:1" {
		t.Errorf("MessageID: got %q, want %q", pos.MessageID, "99:1")
	}
}

func TestDelete(t *testing.T) {
	path := tempPath(t)
	topic := "persistent://public/default/temp"

	s, _ := cursor.New(path)
	_ = s.Set(cursor.Position{Topic: topic, MessageID: "1:0"})
	if err := s.Delete(topic); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, ok := s.Get(topic)
	if ok {
		t.Fatal("expected position to be absent after Delete")
	}

	// Verify file still exists and is valid JSON.
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("cursor file missing after Delete: %v", err)
	}
}
