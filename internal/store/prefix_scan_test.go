package store

import (
	"testing"
	"time"
)

func TestPrefixScanBasic(t *testing.T) {
	s := New()
	s.Set("user:1", "alice", 0)
	s.Set("user:2", "bob", 0)
	s.Set("user:3", "carol", 0)
	s.Set("session:abc", "xyz", 0)

	results := s.PrefixScan("user:", 0)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Key != "user:1" || results[0].Value != "alice" {
		t.Errorf("unexpected first result: %+v", results[0])
	}
}

func TestPrefixScanNoMatch(t *testing.T) {
	s := New()
	s.Set("foo:1", "a", 0)

	results := s.PrefixScan("bar:", 0)
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestPrefixScanSkipsExpired(t *testing.T) {
	s := New()
	s.Set("item:live", "yes", 0)
	s.Set("item:dead", "no", 1) // 1ms TTL
	time.Sleep(5 * time.Millisecond)

	results := s.PrefixScan("item:", 0)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "item:live" {
		t.Errorf("expected item:live, got %s", results[0].Key)
	}
}

func TestPrefixScanWithLimit(t *testing.T) {
	s := New()
	s.Set("k:1", "a", 0)
	s.Set("k:2", "b", 0)
	s.Set("k:3", "c", 0)
	s.Set("k:4", "d", 0)

	results := s.PrefixScan("k:", 2)
	if len(results) != 2 {
		t.Fatalf("expected 2 results with limit, got %d", len(results))
	}
}

func TestPrefixScanSorted(t *testing.T) {
	s := New()
	s.Set("z:3", "c", 0)
	s.Set("z:1", "a", 0)
	s.Set("z:2", "b", 0)

	results := s.PrefixScan("z:", 0)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for i, expected := range []string{"z:1", "z:2", "z:3"} {
		if results[i].Key != expected {
			t.Errorf("index %d: expected %s, got %s", i, expected, results[i].Key)
		}
	}
}
