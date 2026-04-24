package store

import (
	"testing"
)

func TestAppendNewKey(t *testing.T) {
	s := New()
	n, err := s.Append("msg", "hello")
	if err != nil || n != 5 {
		t.Fatalf("expected 5, got %d err %v", n, err)
	}
	v, _ := s.Get("msg")
	if v != "hello" {
		t.Fatalf("expected hello, got %s", v)
	}
}

func TestAppendExistingKey(t *testing.T) {
	s := New()
	s.Set("msg", "hello", 0)
	n, _ := s.Append("msg", " world")
	if n != 11 {
		t.Fatalf("expected 11, got %d", n)
	}
	v, _ := s.Get("msg")
	if v != "hello world" {
		t.Fatalf("expected 'hello world', got %s", v)
	}
}

func TestGetRange(t *testing.T) {
	s := New()
	s.Set("k", "Hello, World!", 0)

	v, _ := s.GetRange("k", 0, 4)
	if v != "Hello" {
		t.Fatalf("expected Hello, got %s", v)
	}

	v, _ = s.GetRange("k", -6, -1)
	if v != "World!" {
		t.Fatalf("expected World!, got %s", v)
	}
}

func TestGetRangeMissingKey(t *testing.T) {
	s := New()
	v, _ := s.GetRange("missing", 0, 5)
	if v != "" {
		t.Fatalf("expected empty string")
	}
}

func TestStrLen(t *testing.T) {
	s := New()
	s.Set("k", "hello", 0)
	if s.StrLen("k") != 5 {
		t.Fatal("expected 5")
	}
	if s.StrLen("missing") != 0 {
		t.Fatal("expected 0 for missing key")
	}
}

func TestGetSet(t *testing.T) {
	s := New()
	s.Set("k", "old", 0)
	old, err := s.GetSet("k", "new")
	if err != nil || old != "old" {
		t.Fatalf("expected old, got %s err %v", old, err)
	}
	v, _ := s.Get("k")
	if v != "new" {
		t.Fatalf("expected new, got %s", v)
	}
}

func TestGetSetMissingKey(t *testing.T) {
	s := New()
	_, err := s.GetSet("missing", "val")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestSetNX(t *testing.T) {
	s := New()
	if !s.SetNX("k", "v1") {
		t.Fatal("expected true for new key")
	}
	if s.SetNX("k", "v2") {
		t.Fatal("expected false for existing key")
	}
	v, _ := s.Get("k")
	if v != "v1" {
		t.Fatalf("expected v1, got %s", v)
	}
}

func TestMSetNX(t *testing.T) {
	s := New()
	ok := s.MSetNX(map[string]string{"a": "1", "b": "2"})
	if !ok {
		t.Fatal("expected true")
	}
	// One key already exists — should fail
	ok = s.MSetNX(map[string]string{"a": "99", "c": "3"})
	if ok {
		t.Fatal("expected false")
	}
	v, _ := s.Get("a")
	if v != "1" {
		t.Fatalf("expected a=1, got %s", v)
	}
}
