package store

import (
	"testing"
	"time"
)

func TestListAllKeys(t *testing.T) {
	s := New()
	s.Set("foo", "1", 0)
	s.Set("bar", "2", 0)
	s.Set("baz", "3", 0)

	keys := s.List("")
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
}

func TestListWithPrefix(t *testing.T) {
	s := New()
	s.Set("user:1", "alice", 0)
	s.Set("user:2", "bob", 0)
	s.Set("session:abc", "xyz", 0)

	keys := s.List("user:")
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}

func TestListSkipsExpired(t *testing.T) {
	s := New()
	s.Set("alive", "yes", 0)
	s.Set("dead", "no", 1*time.Millisecond)

	time.Sleep(10 * time.Millisecond)

	keys := s.List("")
	if len(keys) != 1 || keys[0] != "alive" {
		t.Fatalf("expected only 'alive', got %v", keys)
	}
}

func TestListEmpty(t *testing.T) {
	s := New()
	keys := s.List("")
	if len(keys) != 0 {
		t.Fatalf("expected empty list, got %v", keys)
	}
}

func TestListPrefixNoMatch(t *testing.T) {
	s := New()
	s.Set("foo", "bar", 0)

	keys := s.List("zzz")
	if len(keys) != 0 {
		t.Fatalf("expected no matches, got %v", keys)
	}
}
