package store

import (
	"testing"
)

func TestZAddAndZScore(t *testing.T) {
	sm := NewSortedSetManager()
	sm.ZAdd("scores", "alice", 10.0)
	sm.ZAdd("scores", "bob", 20.0)

	score, err := sm.ZScore("scores", "alice")
	if err != nil || score != 10.0 {
		t.Fatalf("expected 10.0, got %v (err: %v)", score, err)
	}
}

func TestZAddUpdateScore(t *testing.T) {
	sm := NewSortedSetManager()
	sm.ZAdd("scores", "alice", 5.0)
	sm.ZAdd("scores", "alice", 99.0)

	score, err := sm.ZScore("scores", "alice")
	if err != nil || score != 99.0 {
		t.Fatalf("expected updated score 99.0, got %v", score)
	}
	if sm.ZCard("scores") != 1 {
		t.Fatal("expected only one member after update")
	}
}

func TestZRange(t *testing.T) {
	sm := NewSortedSetManager()
	sm.ZAdd("board", "c", 30.0)
	sm.ZAdd("board", "a", 10.0)
	sm.ZAdd("board", "b", 20.0)

	result := sm.ZRange("board", 0, 2)
	if len(result) != 3 {
		t.Fatalf("expected 3 members, got %d", len(result))
	}
	if result[0].Member != "a" || result[1].Member != "b" || result[2].Member != "c" {
		t.Fatalf("unexpected order: %+v", result)
	}
}

func TestZRem(t *testing.T) {
	sm := NewSortedSetManager()
	sm.ZAdd("k", "x", 1.0)
	sm.ZAdd("k", "y", 2.0)

	removed := sm.ZRem("k", "x")
	if !removed {
		t.Fatal("expected ZRem to return true")
	}
	if sm.ZCard("k") != 1 {
		t.Fatal("expected 1 member after removal")
	}
	_, err := sm.ZScore("k", "x")
	if err == nil {
		t.Fatal("expected error for removed member")
	}
}

func TestZScoreMissing(t *testing.T) {
	sm := NewSortedSetManager()
	_, err := sm.ZScore("nokey", "nobody")
	if err == nil {
		t.Fatal("expected error for missing member")
	}
}

func TestZRangeEmpty(t *testing.T) {
	sm := NewSortedSetManager()
	result := sm.ZRange("empty", 0, 10)
	if result != nil {
		t.Fatal("expected nil for empty key")
	}
}
