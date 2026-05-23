package offset

import (
	"os"
	"path/filepath"
	"testing"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "offsets.json")
}

func TestNew_EmptyPath(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNew_FileNotExist(t *testing.T) {
	_, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSet_EmptyTopic(t *testing.T) {
	tr, _ := New(tempPath(t))
	if err := tr.Set("", 1); err == nil {
		t.Fatal("expected error for empty topic")
	}
}

func TestGet_NotFound(t *testing.T) {
	tr, _ := New(tempPath(t))
	_, ok := tr.Get("persistent://public/default/missing")
	if ok {
		t.Fatal("expected false for unknown topic")
	}
}

func TestSet_Get_RoundTrip(t *testing.T) {
	tr, _ := New(tempPath(t))
	topic := "persistent://public/default/events"
	if err := tr.Set(topic, 42); err != nil {
		t.Fatalf("Set: %v", err)
	}
	v, ok := tr.Get(topic)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if v != 42 {
		t.Fatalf("expected 42, got %d", v)
	}
}

func TestSave_Persist(t *testing.T) {
	p := tempPath(t)
	tr, _ := New(p)
	topic := "persistent://public/default/orders"
	_ = tr.Set(topic, 99)
	if err := tr.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	// reload from disk
	tr2, err := New(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	v, ok := tr2.Get(topic)
	if !ok || v != 99 {
		t.Fatalf("expected 99 after reload, got %d (ok=%v)", v, ok)
	}
}

func TestSave_CreatesFile(t *testing.T) {
	p := tempPath(t)
	tr, _ := New(p)
	_ = tr.Set("persistent://public/default/t", 1)
	_ = tr.Save()
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}
