package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/pulsar-watch/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoints.json")
}

func TestNew_EmptyPath(t *testing.T) {
	_, err := checkpoint.New("")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestNew_FileNotExist(t *testing.T) {
	cp, err := checkpoint.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cp == nil {
		t.Fatal("expected non-nil checkpoint")
	}
}

func TestSet_EmptyTopic(t *testing.T) {
	cp, _ := checkpoint.New(tempPath(t))
	if err := cp.Set("", "msg:1:0"); err == nil {
		t.Fatal("expected error for empty topic")
	}
}

func TestGet_NotFound(t *testing.T) {
	cp, _ := checkpoint.New(tempPath(t))
	_, ok := cp.Get("persistent://public/default/missing")
	if ok {
		t.Fatal("expected ok=false for unknown topic")
	}
}

func TestSetAndGet(t *testing.T) {
	path := tempPath(t)
	cp, _ := checkpoint.New(path)
	topic := "persistent://public/default/events"
	msgID := "msg:42:0"

	if err := cp.Set(topic, msgID); err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	got, ok := cp.Get(topic)
	if !ok {
		t.Fatal("expected entry to exist after Set")
	}
	if got != msgID {
		t.Fatalf("expected %q, got %q", msgID, got)
	}
}

func TestPersistence_ReloadFromDisk(t *testing.T) {
	path := tempPath(t)
	topic := "persistent://public/default/orders"
	msgID := "msg:99:1"

	cp1, _ := checkpoint.New(path)
	_ = cp1.Set(topic, msgID)

	cp2, err := checkpoint.New(path)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	got, ok := cp2.Get(topic)
	if !ok {
		t.Fatal("expected entry after reload")
	}
	if got != msgID {
		t.Fatalf("expected %q, got %q", msgID, got)
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	path := tempPath(t)
	cp, _ := checkpoint.New(path)
	topic := "persistent://public/default/temp"

	_ = cp.Set(topic, "msg:1:0")
	if err := cp.Delete(topic); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, ok := cp.Get(topic)
	if ok {
		t.Fatal("expected entry to be removed after Delete")
	}

	// Confirm removal persists on disk.
	cp2, _ := checkpoint.New(path)
	_, ok = cp2.Get(topic)
	if ok {
		t.Fatal("expected entry absent after reload")
	}
}

func TestNew_InvalidJSON(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := checkpoint.New(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON file")
	}
}
