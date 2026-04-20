package store

import (
	"fmt"
	"sort"
	"sync"
)

type LeaderboardEntry struct {
	Member string  `json:"member"`
	Score  float64 `json:"score"`
}

type Leaderboard struct {
	mu      sync.RWMutex
	scores  map[string]float64
}

type LeaderboardManager struct {
	mu     sync.RWMutex
	boards map[string]*Leaderboard
}

func NewLeaderboardManager() *LeaderboardManager {
	return &LeaderboardManager{
		boards: make(map[string]*Leaderboard),
	}
}

func (lm *LeaderboardManager) board(name string) *Leaderboard {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	if _, ok := lm.boards[name]; !ok {
		lm.boards[name] = &Leaderboard{scores: make(map[string]float64)}
	}
	return lm.boards[name]
}

func (lm *LeaderboardManager) Add(name, member string, score float64) {
	b := lm.board(name)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.scores[member] += score
}

func (lm *LeaderboardManager) Set(name, member string, score float64) {
	b := lm.board(name)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.scores[member] = score
}

func (lm *LeaderboardManager) Remove(name, member string) error {
	b := lm.board(name)
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.scores[member]; !ok {
		return fmt.Errorf("member %q not found", member)
	}
	delete(b.scores, member)
	return nil
}

func (lm *LeaderboardManager) Top(name string, n int) []LeaderboardEntry {
	b := lm.board(name)
	b.mu.RLock()
	defer b.mu.RUnlock()

	entries := make([]LeaderboardEntry, 0, len(b.scores))
	for member, score := range b.scores {
		entries = append(entries, LeaderboardEntry{Member: member, Score: score})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})
	if n > 0 && n < len(entries) {
		return entries[:n]
	}
	return entries
}

func (lm *LeaderboardManager) Rank(name, member string) (int, float64, error) {
	top := lm.Top(name, 0)
	for i, e := range top {
		if e.Member == member {
			return i + 1, e.Score, nil
		}
	}
	return 0, 0, fmt.Errorf("member %q not found", member)
}
