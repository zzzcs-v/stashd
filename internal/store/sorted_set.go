package store

import (
	"fmt"
	"sort"
	"sync"
)

type SortedSetMember struct {
	Member string
	Score  float64
}

type SortedSetManager struct {
	mu   sync.RWMutex
	sets map[string][]SortedSetMember
}

func NewSortedSetManager() *SortedSetManager {
	return &SortedSetManager{
		sets: make(map[string][]SortedSetMember),
	}
}

func (s *SortedSetManager) ZAdd(key, member string, score float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, m := range s.sets[key] {
		if m.Member == member {
			s.sets[key][i].Score = score
			s.sortKey(key)
			return
		}
	}
	s.sets[key] = append(s.sets[key], SortedSetMember{Member: member, Score: score})
	s.sortKey(key)
}

func (s *SortedSetManager) sortKey(key string) {
	sort.Slice(s.sets[key], func(i, j int) bool {
		return s.sets[key][i].Score < s.sets[key][j].Score
	})
}

func (s *SortedSetManager) ZScore(key, member string) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, m := range s.sets[key] {
		if m.Member == member {
			return m.Score, nil
		}
	}
	return 0, fmt.Errorf("member not found")
}

func (s *SortedSetManager) ZRange(key string, start, stop int) []SortedSetMember {
	s.mu.RLock()
	defer s.mu.RUnlock()
	members := s.sets[key]
	n := len(members)
	if n == 0 {
		return nil
	}
	if stop < 0 || stop >= n {
		stop = n - 1
	}
	if start < 0 {
		start = 0
	}
	if start > stop {
		return nil
	}
	result := make([]SortedSetMember, stop-start+1)
	copy(result, members[start:stop+1])
	return result
}

func (s *SortedSetManager) ZRem(key, member string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, m := range s.sets[key] {
		if m.Member == member {
			s.sets[key] = append(s.sets[key][:i], s.sets[key][i+1:]...)
			return true
		}
	}
	return false
}

func (s *SortedSetManager) ZCard(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.sets[key])
}
