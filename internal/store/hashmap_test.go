package store

import (
	"testing"
)

func TestHSetAndHGet(t *testing.T) {
	h := NewHashMap()
	h.HSet("user:1", "name", "alice")
	v, err := h.HGet("user:1", "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "alice" {
		t.Errorf("expected alice, got %s", v)
	}
}

func TestHGetMissingKey(t *testing.T) {
	h := NewHashMap()
	_, err := h.HGet("nokey", "field")
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestHGetMissingField(t *testing.T) {
	h := NewHashMap()
	h.HSet("k", "f1", "v1")
	_, err := h.HGet("k", "missing")
	if err == nil {
		t.Error("expected error for missing field")
	}
}

func TestHGetAll(t *testing.T) {
	h := NewHashMap()
	h.HSet("obj", "a", "1")
	h.HSet("obj", "b", "2")
	all, err := h.HGetAll("obj")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 fields, got %d", len(all))
	}
}

func TestHDel(t *testing.T) {
	h := NewHashMap()
	h.HSet("k", "f", "v")
	err := h.HDel("k", "f")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.HExists("k", "f") {
		t.Error("field should have been deleted")
	}
}

func TestHDelRemovesEmptyKey(t *testing.T) {
	h := NewHashMap()
	h.HSet("k", "f", "v")
	h.HDel("k", "f")
	_, err := h.HGetAll("k")
	if err == nil {
		t.Error("expected key to be removed after last field deleted")
	}
}

func TestHExists(t *testing.T) {
	h := NewHashMap()
	h.HSet("k", "field", "val")
	if !h.HExists("k", "field") {
		t.Error("expected field to exist")
	}
	if h.HExists("k", "other") {
		t.Error("expected field to not exist")
	}
}
