package store

import (
	"testing"
	"time"
)

func newCASManager() (*CASManager, *Store) {
	s := New()
	return NewCASManager(s), s
}

func TestCASSuccess(t *testing.T) {
	cas, s := newCASManager()
	s.Set("k", "old", 0)

	if err := cas.CompareAndSwap("k", "old", "new", 0); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	val, ok := s.Get("k")
	if !ok || val != "new" {
		t.Fatalf("expected 'new', got %q ok=%v", val, ok)
	}
}

func TestCASConflict(t *testing.T) {
	cas, s := newCASManager()
	s.Set("k", "current", 0)

	err := cas.CompareAndSwap("k", "wrong", "new", 0)
	if err != ErrCASConflict {
		t.Fatalf("expected ErrCASConflict, got %v", err)
	}
	val, _ := s.Get("k")
	if val != "current" {
		t.Fatalf("value should be unchanged, got %q", val)
	}
}

func TestCASMissingKey(t *testing.T) {
	cas, _ := newCASManager()

	err := cas.CompareAndSwap("missing", "x", "y", 0)
	if err != ErrCASMissing {
		t.Fatalf("expected ErrCASMissing, got %v", err)
	}
}

func TestCASWithTTL(t *testing.T) {
	cas, s := newCASManager()
	s.Set("k", "v", 0)

	if err := cas.CompareAndSwap("k", "v", "v2", 50*time.Millisecond); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	time.Sleep(80 * time.Millisecond)
	_, ok := s.Get("k")
	if ok {
		t.Fatal("expected key to have expired after TTL")
	}
}

func TestCASDeleteSuccess(t *testing.T) {
	cas, s := newCASManager()
	s.Set("k", "val", 0)

	if err := cas.CompareAndDelete("k", "val"); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	_, ok := s.Get("k")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestCASDeleteConflict(t *testing.T) {
	cas, s := newCASManager()
	s.Set("k", "val", 0)

	err := cas.CompareAndDelete("k", "wrong")
	if err != ErrCASConflict {
		t.Fatalf("expected ErrCASConflict, got %v", err)
	}
	_, ok := s.Get("k")
	if !ok {
		t.Fatal("key should still exist after conflict")
	}
}

func TestCASDeleteMissing(t *testing.T) {
	cas, _ := newCASManager()

	err := cas.CompareAndDelete("ghost", "x")
	if err != ErrCASMissing {
		t.Fatalf("expected ErrCASMissing, got %v", err)
	}
}
