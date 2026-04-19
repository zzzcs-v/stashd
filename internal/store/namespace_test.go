package store

import (
	"testing"
	"time"
)

func TestNamespacedKey(t *testing.T) {
	if got := NamespacedKey("ns", "key"); got != "ns:key" {
		t.Errorf("expected ns:key, got %s", got)
	}
	if got := NamespacedKey("", "key"); got != "key" {
		t.Errorf("expected key, got %s", got)
	}
}

func TestParseNamespacedKey(t *testing.T) {
	ns, key := ParseNamespacedKey("ns:key")
	if ns != "ns" || key != "key" {
		t.Errorf("expected ns/key, got %s/%s", ns, key)
	}
	ns, key = ParseNamespacedKey("barekey")
	if ns != "" || key != "barekey" {
		t.Errorf("expected empty ns, got %s/%s", ns, key)
	}
}

func TestListNamespace(t *testing.T) {
	s := New()
	s.Set(NamespacedKey("users", "alice"), "1", 0)
	s.Set(NamespacedKey("users", "bob"), "2", 0)
	s.Set(NamespacedKey("orders", "x"), "3", 0)

	keys := s.ListNamespace("users")
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

func TestListNamespaceSkipsExpired(t *testing.T) {
	s := New()
	s.Set(NamespacedKey("ns", "live"), "1", 0)
	s.Set(NamespacedKey("ns", "dead"), "2", 1*time.Millisecond)
	time.Sleep(10 * time.Millisecond)

	keys := s.ListNamespace("ns")
	if len(keys) != 1 || keys[0] != "live" {
		t.Errorf("expected only live key, got %v", keys)
	}
}

func TestDeleteNamespace(t *testing.T) {
	s := New()
	s.Set(NamespacedKey("ns", "a"), "1", 0)
	s.Set(NamespacedKey("ns", "b"), "2", 0)
	s.Set(NamespacedKey("other", "c"), "3", 0)

	count := s.DeleteNamespace("ns")
	if count != 2 {
		t.Errorf("expected 2 deleted, got %d", count)
	}
	if keys := s.ListNamespace("ns"); len(keys) != 0 {
		t.Errorf("expected empty namespace, got %v", keys)
	}
	if keys := s.ListNamespace("other"); len(keys) != 1 {
		t.Errorf("expected other namespace intact, got %v", keys)
	}
}
