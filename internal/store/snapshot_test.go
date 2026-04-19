package store

import (
	"os"
	"testing"
	"time"
)

func TestSaveAndLoadSnapshot(t *testing.T) {
	s := New()
	s.Set("foo", "bar", 0)
	s.Set("hello", "world", 0)

	tmp, err := os.CreateTemp("", "stashd-snapshot-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	if err := s.SaveSnapshot(tmp.Name()); err != nil {
		t.Fatalf("SaveSnapshot failed: %v", err)
	}

	s2 := New()
	if err := s2.LoadSnapshot(tmp.Name()); err != nil {
		t.Fatalf("LoadSnapshot failed: %v", err)
	}

	if v, ok := s2.Get("foo"); !ok || v != "bar" {
		t.Errorf("expected foo=bar, got %q ok=%v", v, ok)
	}
	if v, ok := s2.Get("hello"); !ok || v != "world" {
		t.Errorf("expected hello=world, got %q ok=%v", v, ok)
	}
}

func TestLoadSnapshotSkipsExpired(t *testing.T) {
	s := New()
	s.Set("alive", "yes", 0)

	// manually insert an expired entry
	past := time.Now().Add(-time.Second)
	s.mu.Lock()
	s.data["dead"] = entry{value: "no", expiresAt: &past}
	s.mu.Unlock()

	tmp, err := os.CreateTemp("", "stashd-snapshot-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	s.SaveSnapshot(tmp.Name())

	s2 := New()
	s2.LoadSnapshot(tmp.Name())

	if _, ok := s2.Get("dead"); ok {
		t.Error("expected expired key 'dead' to be skipped")
	}
	if v, ok := s2.Get("alive"); !ok || v != "yes" {
		t.Errorf("expected alive=yes, got %q ok=%v", v, ok)
	}
}

func TestLoadSnapshotMissingFile(t *testing.T) {
	s := New()
	if err := s.LoadSnapshot("/tmp/stashd-nonexistent-file.json"); err != nil {
		t.Errorf("expected nil for missing file, got %v", err)
	}
}
