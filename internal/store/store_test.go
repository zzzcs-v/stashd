package store

import (
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	s := New()
	s.Set("hello", "world", 0)
	v, ok := s.Get("hello")
	if !ok || v != "world" {
		t.Fatalf("expected world, got %q ok=%v", v, ok)
	}
}

func TestGetMissing(t *testing.T) {
	s := New()
	_, ok := s.Get("nope")
	if ok {
		t.Fatal("expected miss")
	}
}

func TestDelete(t *testing.T) {
	s := New()
	s.Set("k", "v", 0)
	if !s.Delete("k") {
		t.Fatal("expected true on delete")
	}
	_, ok := s.Get("k")
	if ok {
		t.Fatal("expected miss after delete")
	}
	if s.Delete("k") {
		t.Fatal("expected false on second delete")
	}
}

func TestTTLExpiry(t *testing.T) {
	s := New()
	s.Set("temp", "val", 50*time.Millisecond)
	v, ok := s.Get("temp")
	if !ok || v != "val" {
		t.Fatal("expected hit before expiry")
	}
	time.Sleep(100 * time.Millisecond)
	_, ok = s.Get("temp")
	if ok {
		t.Fatal("expected miss after expiry")
	}
}

func TestNoTTL(t *testing.T) {
	s := New()
	s.Set("persist", "yes", 0)
	time.Sleep(20 * time.Millisecond)
	_, ok := s.Get("persist")
	if !ok {
		t.Fatal("key without TTL should persist")
	}
}
