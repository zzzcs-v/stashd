package store

import (
	"fmt"
	"testing"
)

func TestHLLAddAndCount(t *testing.T) {
	m := NewHyperLogLogManager()
	m.Add("visitors", "user1", "user2", "user3")

	count, ok := m.Count("visitors")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if count <= 0 {
		t.Errorf("expected positive count, got %d", count)
	}
}

func TestHLLMissingKey(t *testing.T) {
	m := NewHyperLogLogManager()
	_, ok := m.Count("nonexistent")
	if ok {
		t.Error("expected false for missing key")
	}
}

func TestHLLDeduplication(t *testing.T) {
	m := NewHyperLogLogManager()
	// Add same value many times — estimate should stay low
	for i := 0; i < 100; i++ {
		m.Add("k", "same-value")
	}
	count, _ := m.Count("k")
	if count > 10 {
		t.Errorf("expected low estimate for single unique value, got %d", count)
	}
}

func TestHLLApproximateCardinality(t *testing.T) {
	m := NewHyperLogLogManager()
	for i := 0; i < 1000; i++ {
		m.Add("big", fmt.Sprintf("user-%d", i))
	}
	count, ok := m.Count("big")
	if !ok {
		t.Fatal("expected key to exist")
	}
	// Allow 30% error margin for this simple implementation
	if count < 700 || count > 1300 {
		t.Errorf("estimate %d too far from 1000", count)
	}
}

func TestHLLDelete(t *testing.T) {
	m := NewHyperLogLogManager()
	m.Add("temp", "a", "b", "c")
	m.Delete("temp")
	_, ok := m.Count("temp")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestHLLMultipleKeys(t *testing.T) {
	m := NewHyperLogLogManager()
	m.Add("k1", "a", "b")
	m.Add("k2", "x", "y", "z")

	c1, _ := m.Count("k1")
	c2, _ := m.Count("k2")
	if c1 <= 0 || c2 <= 0 {
		t.Error("expected positive counts for both keys")
	}
}
