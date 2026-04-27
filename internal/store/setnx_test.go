package store

import (
	"testing"
	"time"
)

func TestSetNXNewKey(t *testing.T) {
	s := New()
	set := s.SetNX("foo", "bar", 0)
	if !set {
		t.Fatal("expected SetNX to return true for new key")
	}
	val, ok := s.Get("foo")
	if !ok || val != "bar" {
		t.Fatalf("expected 'bar', got '%s'", val)
	}
}

func TestSetNXExistingKey(t *testing.T) {
	s := New()
	s.Set("foo", "original", 0)
	set := s.SetNX("foo", "new", 0)
	if set {
		t.Fatal("expected SetNX to return false for existing key")
	}
	val, _ := s.Get("foo")
	if val != "original" {
		t.Fatalf("expected 'original', got '%s'", val)
	}
}

func TestSetNXExpiredKey(t *testing.T) {
	s := New()
	s.Set("foo", "old", 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	set := s.SetNX("foo", "new", 0)
	if !set {
		t.Fatal("expected SetNX to succeed on expired key")
	}
	val, ok := s.Get("foo")
	if !ok || val != "new" {
		t.Fatalf("expected 'new', got '%s'", val)
	}
}

func TestGetSet(t *testing.T) {
	s := New()
	s.Set("foo", "old", 0)
	old, found := s.GetSet("foo", "new")
	if !found || old != "old" {
		t.Fatalf("expected old='old' found=true, got old='%s' found=%v", old, found)
	}
	val, _ := s.Get("foo")
	if val != "new" {
		t.Fatalf("expected 'new', got '%s'", val)
	}
}

func TestGetSetMissingKey(t *testing.T) {
	s := New()
	old, found := s.GetSet("missing", "val")
	if found || old != "" {
		t.Fatalf("expected empty/false for missing key, got '%s'/%v", old, found)
	}
}

func TestSetXXExistingKey(t *testing.T) {
	s := New()
	s.Set("foo", "old", 0)
	updated := s.SetXX("foo", "new", 0)
	if !updated {
		t.Fatal("expected SetXX to return true for existing key")
	}
	val, _ := s.Get("foo")
	if val != "new" {
		t.Fatalf("expected 'new', got '%s'", val)
	}
}

func TestSetXXMissingKey(t *testing.T) {
	s := New()
	updated := s.SetXX("ghost", "val", 0)
	if updated {
		t.Fatal("expected SetXX to return false for missing key")
	}
	_, ok := s.Get("ghost")
	if ok {
		t.Fatal("expected key to not exist")
	}
}

func TestSetXXExpiredKey(t *testing.T) {
	s := New()
	s.Set("foo", "old", 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	updated := s.SetXX("foo", "new", 0)
	if updated {
		t.Fatal("expected SetXX to return false for expired key")
	}
	_, ok := s.Get("foo")
	if ok {
		t.Fatal("expected expired key to not exist after SetXX")
	}
}
