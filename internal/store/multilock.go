package store

import (
	"fmt"
	"sort"
	"sync"
)

// MultiLockManager allows atomically acquiring locks on multiple keys at once,
// preventing deadlocks by always acquiring locks in sorted key order.
type MultiLockManager struct {
	mu    sync.Mutex
	locks map[string]*multiLockEntry
}

type multiLockEntry struct {
	mu      sync.Mutex
	holders int
}

func NewMultiLockManager() *MultiLockManager {
	return &MultiLockManager{
		locks: make(map[string]*multiLockEntry),
	}
}

func (m *MultiLockManager) getOrCreate(key string) *multiLockEntry {
	m.mu.Lock()
	defer m.mu.Unlock()
	if e, ok := m.locks[key]; ok {
		e.holders++
		return e
	}
	e := &multiLockEntry{holders: 1}
	m.locks[key] = e
	return e
}

func (m *MultiLockManager) release(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if e, ok := m.locks[key]; ok {
		e.holders--
		if e.holders == 0 {
			delete(m.locks, key)
		}
	}
}

// Lock acquires locks on all provided keys in sorted order to avoid deadlocks.
// Returns a token string and an unlock function.
func (m *MultiLockManager) Lock(keys []string) (string, func()) {
	if len(keys) == 0 {
		return "", func() {}
	}

	// Deduplicate and sort
	seen := make(map[string]struct{})
	uniq := keys[:0:0]
	for _, k := range keys {
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			uniq = append(uniq, k)
		}
	}
	sort.Strings(uniq)

	entries := make([]*multiLockEntry, len(uniq))
	for i, k := range uniq {
		entries[i] = m.getOrCreate(k)
	}
	for _, e := range entries {
		e.mu.Lock()
	}

	token := fmt.Sprintf("%v", uniq)
	return token, func() {
		for i := len(entries) - 1; i >= 0; i-- {
			entries[i].mu.Unlock()
		}
		for _, k := range uniq {
			m.release(k)
		}
	}
}
