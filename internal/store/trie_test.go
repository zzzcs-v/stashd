package store

import (
	"sort"
	"testing"
)

func TestTrieInsertAndSearch(t *testing.T) {
	tr := NewTrie()
	tr.Insert("apple")
	tr.Insert("app")
	tr.Insert("application")
	tr.Insert("banana")

	results := tr.Search("app")
	sort.Strings(results)

	expected := []string{"app", "apple", "application"}
	if len(results) != len(expected) {
		t.Fatalf("expected %d results, got %d: %v", len(expected), len(results), results)
	}
	for i, v := range expected {
		if results[i] != v {
			t.Errorf("expected %s at index %d, got %s", v, i, results[i])
		}
	}
}

func TestTrieSearchNoMatch(t *testing.T) {
	tr := NewTrie()
	tr.Insert("hello")

	results := tr.Search("world")
	if len(results) != 0 {
		t.Errorf("expected no results, got %v", results)
	}
}

func TestTrieHasPrefix(t *testing.T) {
	tr := NewTrie()
	tr.Insert("stash")

	if !tr.HasPrefix("sta") {
		t.Error("expected HasPrefix to return true for 'sta'")
	}
	if tr.HasPrefix("xyz") {
		t.Error("expected HasPrefix to return false for 'xyz'")
	}
}

func TestTrieDelete(t *testing.T) {
	tr := NewTrie()
	tr.Insert("key")
	tr.Insert("keystone")

	deleted := tr.Delete("key")
	if !deleted {
		t.Error("expected Delete to return true")
	}

	results := tr.Search("key")
	if len(results) != 1 || results[0] != "keystone" {
		t.Errorf("expected only 'keystone' after deletion, got %v", results)
	}
}

func TestTrieDeleteMissingKey(t *testing.T) {
	tr := NewTrie()
	tr.Insert("exist")

	if tr.Delete("missing") {
		t.Error("expected Delete to return false for missing key")
	}
}

func TestTrieSearchEmptyPrefix(t *testing.T) {
	tr := NewTrie()
	tr.Insert("a")
	tr.Insert("b")
	tr.Insert("c")

	results := tr.Search("")
	if len(results) != 3 {
		t.Errorf("expected 3 results for empty prefix, got %d", len(results))
	}
}
