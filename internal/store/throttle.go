package store

import (
	"fmt"
	"sync"
	"time"
)

// ThrottleEntry tracks request count and window for a key.
type ThrottleEntry struct {
	Count     int
	WindowEnd time.Time
}

// ThrottleManager enforces a max number of requests per duration window.
type ThrottleManager struct {
	mu      sync.Mutex
	entries map[string]*ThrottleEntry
}

// NewThrottleManager creates a new ThrottleManager.
func NewThrottleManager() *ThrottleManager {
	return &ThrottleManager{
		entries: make(map[string]*ThrottleEntry),
	}
}

// Allow checks whether the given key is within its allowed limit per window.
// Returns (allowed bool, remaining int, resetAt time.Time).
func (t *ThrottleManager) Allow(key string, limit int, window time.Duration) (bool, int, time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	entry, ok := t.entries[key]
	if !ok || now.After(entry.WindowEnd) {
		entry = &ThrottleEntry{
			Count:     0,
			WindowEnd: now.Add(window),
		}
		t.entries[key] = entry
	}

	if entry.Count >= limit {
		return false, 0, entry.WindowEnd
	}

	entry.Count++
	remaining := limit - entry.Count
	return true, remaining, entry.WindowEnd
}

// Reset clears throttle state for the given key.
func (t *ThrottleManager) Reset(key string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, ok := t.entries[key]; !ok {
		return fmt.Errorf("throttle key not found: %s", key)
	}
	delete(t.entries, key)
	return nil
}

// Status returns current count and window end for a key without incrementing.
func (t *ThrottleManager) Status(key string) (int, time.Time, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	entry, ok := t.entries[key]
	if !ok || now.After(entry.WindowEnd) {
		return 0, time.Time{}, false
	}
	return entry.Count, entry.WindowEnd, true
}
