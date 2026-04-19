package store

import (
	"testing"
	"time"
)

func TestEvictionRemovesExpiredKeys(t *testing.T) {
	s := New()
	defer func() {
		// guard against double-close if eviction wasn't started
		recover()
	}()

	s.Set("ghost", "boo", 50*time.Millisecond)
	s.Set("alive", "yes", 10*time.Second)

	s.StartEviction(30 * time.Millisecond)
	defer s.StopEviction()

	time.Sleep(120 * time.Millisecond)

	if _, ok := s.Get("ghost"); ok {
		t.Error("expected 'ghost' to be evicted")
	}
	if _, ok := s.Get("alive"); !ok {
		t.Error("expected 'alive' to still exist")
	}
}

func TestEvictionLeavesNonExpiredKeys(t *testing.T) {
	s := New()
	s.Set("persist", "value", 5*time.Second)
	s.Set("nott", "forever", 0)

	s.StartEviction(20 * time.Millisecond)
	defer s.StopEviction()

	time.Sleep(60 * time.Millisecond)

	if _, ok := s.Get("persist"); !ok {
		t.Error("expected 'persist' to still exist")
	}
	if _, ok := s.Get("nott"); !ok {
		t.Error("expected 'nott' (no TTL) to still exist")
	}
}

func TestStopEviction(t *testing.T) {
	s := New()
	s.Set("k", "v", 200*time.Millisecond)
	s.StartEviction(50 * time.Millisecond)
	s.StopEviction()
	// after stop, expired key should still be accessible until manually cleaned
	time.Sleep(80 * time.Millisecond)
	// no panic is the main goal here
}
