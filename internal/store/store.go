package store

import (
	"sync"
	"time"
)

type entry struct {
	value     string
	expiresAt time.Time
	hasTTL    bool
}

type Store struct {
	mu   sync.RWMutex
	items map[string]entry
}

func New() *Store {
	s := &Store{
		items: make(map[string]entry),
	}
	go s.reap()
	return s
}

func (s *Store) Set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e := entry{value: value}
	if ttl > 0 {
		e.hasTTL = true
		e.expiresAt = time.Now().Add(ttl)
	}
	s.items[key] = e
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.items[key]
	if !ok {
		return "", false
	}
	if e.hasTTL && time.Now().After(e.expiresAt) {
		return "", false
	}
	return e.value, true
}

func (s *Store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.items[key]
	if ok {
		delete(s.items, key)
	}
	return ok
}

func (s *Store) reap() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		s.mu.Lock()
		for k, e := range s.items {
			if e.hasTTL && now.After(e.expiresAt) {
				delete(s.items, k)
			}
		}
		s.mu.Unlock()
	}
}
