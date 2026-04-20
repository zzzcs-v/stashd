package store

import (
	"testing"
	"time"
)

func TestExpireAt(t *testing.T) {
	s := New()
	s.Set("foo", "bar", 0)

	future := time.Now().Add(5 * time.Second).Unix()
	if err := s.ExpireAt("foo", future); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := s.ExpireTime("foo")
	if got != future {
		t.Errorf("expected %d, got %d", future, got)
	}
}

func TestExpireAtMissingKey(t *testing.T) {
	s := New()
	err := s.ExpireAt("missing", time.Now().Add(5*time.Second).Unix())
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestPExpireAt(t *testing.T) {
	s := New()
	s.Set("foo", "bar", 0)

	futureMs := time.Now().Add(10*time.Second).UnixMilli()
	if err := s.PExpireAt("foo", futureMs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := s.ExpireTime("foo")
	expected := time.UnixMilli(futureMs).Unix()
	if got != expected {
		t.Errorf("expected %d, got %d", expected, got)
	}
}

func TestExpireTimeNoExpiry(t *testing.T) {
	s := New()
	s.Set("foo", "bar", 0)

	got := s.ExpireTime("foo")
	if got != -1 {
		t.Errorf("expected -1 for no expiry, got %d", got)
	}
}

func TestExpireTimeMissingKey(t *testing.T) {
	s := New()
	got := s.ExpireTime("ghost")
	if got != -2 {
		t.Errorf("expected -2 for missing key, got %d", got)
	}
}

func TestExpireAtAlreadyExpired(t *testing.T) {
	s := New()
	s.Set("foo", "bar", 0)

	past := time.Now().Add(-1 * time.Second).Unix()
	_ = s.ExpireAt("foo", past)

	got := s.ExpireTime("foo")
	if got != -2 {
		t.Errorf("expected -2 for expired key, got %d", got)
	}
}
