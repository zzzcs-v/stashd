package store

import (
	"testing"
)

func TestIncrNewKey(t *testing.T) {
	s := New()
	val, err := s.Incr("hits")
	if err != nil {
		t.Fatal(err)
	}
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}
}

func TestIncrExistingKey(t *testing.T) {
	s := New()
	s.Set("hits", "10", 0)
	val, err := s.Incr("hits")
	if err != nil {
		t.Fatal(err)
	}
	if val != 11 {
		t.Errorf("expected 11, got %d", val)
	}
}

func TestIncrByDelta(t *testing.T) {
	s := New()
	s.Set("score", "5", 0)
	val, err := s.IncrBy("score", 10)
	if err != nil {
		t.Fatal(err)
	}
	if val != 15 {
		t.Errorf("expected 15, got %d", val)
	}
}

func TestDecrKey(t *testing.T) {
	s := New()
	s.Set("count", "3", 0)
	val, err := s.Decr("count")
	if err != nil {
		t.Fatal(err)
	}
	if val != 2 {
		t.Errorf("expected 2, got %d", val)
	}
}

func TestIncrNonInteger(t *testing.T) {
	s := New()
	s.Set("key", "hello", 0)
	_, err := s.Incr("key")
	if err == nil {
		t.Error("expected error for non-integer value")
	}
}
