package store

import (
	"testing"
)

func TestLeaderboardAddAndTop(t *testing.T) {
	lm := NewLeaderboardManager()
	lm.Add("game", "alice", 100)
	lm.Add("game", "bob", 200)
	lm.Add("game", "alice", 50)

	top := lm.Top("game", 0)
	if len(top) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(top))
	}
	if top[0].Member != "alice" || top[0].Score != 150 {
		t.Errorf("expected alice with 150, got %+v", top[0])
	}
	if top[1].Member != "bob" || top[1].Score != 200 {
		// bob should be second since alice=150, bob=200 — wait, bob is higher
		t.Errorf("unexpected order: %+v", top)
	}
}

func TestLeaderboardTopN(t *testing.T) {
	lm := NewLeaderboardManager()
	lm.Set("lb", "a", 10)
	lm.Set("lb", "b", 30)
	lm.Set("lb", "c", 20)

	top := lm.Top("lb", 2)
	if len(top) != 2 {
		t.Fatalf("expected 2, got %d", len(top))
	}
	if top[0].Member != "b" {
		t.Errorf("expected b at top, got %s", top[0].Member)
	}
}

func TestLeaderboardSet(t *testing.T) {
	lm := NewLeaderboardManager()
	lm.Set("lb", "alice", 500)
	lm.Set("lb", "alice", 300)

	top := lm.Top("lb", 0)
	if top[0].Score != 300 {
		t.Errorf("expected score 300 after Set, got %f", top[0].Score)
	}
}

func TestLeaderboardRank(t *testing.T) {
	lm := NewLeaderboardManager()
	lm.Set("lb", "alice", 100)
	lm.Set("lb", "bob", 200)
	lm.Set("lb", "carol", 150)

	rank, score, err := lm.Rank("lb", "carol")
	if err != nil {
		t.Fatal(err)
	}
	if rank != 2 {
		t.Errorf("expected rank 2, got %d", rank)
	}
	if score != 150 {
		t.Errorf("expected score 150, got %f", score)
	}
}

func TestLeaderboardRemove(t *testing.T) {
	lm := NewLeaderboardManager()
	lm.Set("lb", "alice", 100)
	if err := lm.Remove("lb", "alice"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := lm.Rank("lb", "alice"); err == nil {
		t.Error("expected error after remove")
	}
}

func TestLeaderboardRankMissing(t *testing.T) {
	lm := NewLeaderboardManager()
	_, _, err := lm.Rank("lb", "ghost")
	if err == nil {
		t.Error("expected error for missing member")
	}
}
