package store

import (
	"fmt"
	"sync"
	"time"
)

// FrequencyEntry tracks how many times a key has been accessed within a window.
type FrequencyEntry struct {
	Count     int64
	WindowEnd time.Time
}

// FrequencyManager tracks access frequency for arbitrary keys with a sliding window.
type FrequencyManager struct {
	mu      sync.Mutex
	entries map[string]*FrequencyEntry
	window  time.Duration
}

// NewFrequencyManager creates a new FrequencyManager with the given window duration.
func NewFrequencyManager(window time.Duration) *FrequencyManager {
	return &FrequencyManager{
		entries: make(map[string]*FrequencyEntry),
		window:  window,
	}
}

// Hit records one access for the given key, resetting the window if expired.
func (f *FrequencyManager) Hit(key string) int64 {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()
	e, ok := f.entries[key]
	if !ok || now.After(e.WindowEnd) {
		f.entries[key] = &FrequencyEntry{Count: 1, WindowEnd: now.Add(f.window)}
		return 1
	}
	e.Count++
	return e.Count
}

// Count returns the current hit count for a key within its active window.
// Returns 0 if the key is missing or the window has expired.
func (f *FrequencyManager) Count(key string) int64 {
	f.mu.Lock()
	defer f.mu.Unlock()

	e, ok := f.entries[key]
	if !ok || time.Now().After(e.WindowEnd) {
		return 0
	}
	return e.Count
}

// Reset clears the frequency data for the given key.
func (f *FrequencyManager) Reset(key string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.entries[key]; !ok {
		return fmt.Errorf("frequency: key not found: %s", key)
	}
	delete(f.entries, key)
	return nil
}

// TTL returns the remaining time in the current window for the key.
func (f *FrequencyManager) TTL(key string) (time.Duration, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()

	e, ok := f.entries[key]
	if !ok {
		return 0, false
	}
	remaining := time.Until(e.WindowEnd)
	if remaining <= 0 {
		return 0, false
	}
	return remaining, true
}
