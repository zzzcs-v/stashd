package store

import (
	"testing"
	"time"
)

func TestBatchSet(t *testing.T) {
	s := New()
	items := []BatchSetItem{
		{Key: "a", Value: "1", TTL: 0},
		{Key: "b", Value: "2", TTL: 0},
		{Key: "c", Value: "3", TTL: 0},
	}
	s.BatchSet(items)
	for _, item := range items {
		val, ok := s.Get(item.Key)
		if !ok || val != item.Value {
			t.Errorf("expected %s=%s, got %s (found=%v)", item.Key, item.Value, val, ok)
		}
	}
}

func TestBatchGet(t *testing.T) {
	s := New()
	s.Set("x", "10", 0)
	s.Set("y", "20", 0)

	results := s.BatchGet([]string{"x", "y", "z"})
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if !results[0].Found || results[0].Value != "10" {
		t.Errorf("expected x=10")
	}
	if !results[1].Found || results[1].Value != "20" {
		t.Errorf("expected y=20")
	}
	if results[2].Found {
		t.Errorf("expected z to be missing")
	}
}

func TestBatchDelete(t *testing.T) {
	s := New()
	s.Set("p", "1", 0)
	s.Set("q", "2", 0)
	s.BatchDelete([]string{"p", "q"})
	if _, ok := s.Get("p"); ok {
		t.Error("expected p to be deleted")
	}
	if _, ok := s.Get("q"); ok {
		t.Error("expected q to be deleted")
	}
}

func TestBatchSetWithTTL(t *testing.T) {
	s := New()
	s.BatchSet([]BatchSetItem{
		{Key: "ttlkey", Value: "val", TTL: 50 * time.Millisecond},
	})
	if _, ok := s.Get("ttlkey"); !ok {
		t.Error("expected ttlkey to exist")
	}
	time.Sleep(100 * time.Millisecond)
	if _, ok := s.Get("ttlkey"); ok {
		t.Error("expected ttlkey to be expired")
	}
}
