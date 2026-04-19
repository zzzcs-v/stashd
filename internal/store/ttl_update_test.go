package store

import (
	"testing"
	"time"
)

func TestTouchResetsExpiry(t *testing.T) {
	s := New()
	s.Set("key", "value", 100*time.Millisecond)

	time.Sleep(50 * time.Millisecond)
	ok := s.Touch("key", 500*time.Millisecond)
	if !ok {
		t.Fatal("expected Touch to return true")
	}

	time.Sleep(200 * time.Millisecond)
	val, found := s.Get("key")
	if !found {
		t.Fatal("expected key to still exist after TTL reset")
	}
	if val != "value" {
		t.Errorf("expected 'value', got %v", val)
	}
}

func TestTouchMissingKey(t *testing.T) {
	s := New()
	ok := s.Touch("missing", time.Second)
	if ok {
		t.Fatal("expected Touch to return false for missing key")
	}
}

func TestTouchExpiredKey(t *testing.T) {
	s := New()
	s.Set("key", "value", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	ok := s.Touch("key", time.Second)
	if ok {
		t.Fatal("expected Touch to return false for expired key")
	}
}

func TestTTLReturnsRemaining(t *testing.T) {
	s := New()
	s.Set("key", "value", 500*time.Millisecond)

	ttl := s.TTL("key")
	if ttl <= 0 || ttl > 500*time.Millisecond {
		t.Errorf("unexpected TTL: %v", ttl)
	}
}

func TestTTLNoExpiry(t *testing.T) {
	s := New()
	s.Set("key", "value", 0)

	if s.TTL("key") != -1 {
		t.Error("expected -1 for key with no expiry")
	}
}

func TestTTLMissingKey(t *testing.T) {
	s := New()
	if s.TTL("ghost") != -2 {
		t.Error("expected -2 for missing key")
	}
}
