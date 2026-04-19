package store

import (
	"sync"
	"time"
)

type rateLimitEntry struct {
	count    int
	windowStart time.Time
}

// RateLimiter tracks request counts per key within a time window.
type RateLimiter struct {
	mu      sync.Mutex
	entries map[string]*rateLimitEntry
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		entries: make(map[string]*rateLimitEntry),
		limit:   limit,
		window:  window,
	}
}

// Allow returns true if the key is within the rate limit.
func (r *RateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	e, ok := r.entries[key]
	if !ok || now.Sub(e.windowStart) >= r.window {
		r.entries[key] = &rateLimitEntry{count: 1, windowStart: now}
		return true
	}
	if e.count >= r.limit {
		return false
	}
	e.count++
	return true
}

// Remaining returns how many requests are left in the current window.
func (r *RateLimiter) Remaining(key string) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	e, ok := r.entries[key]
	if !ok || now.Sub(e.windowStart) >= r.window {
		return r.limit
	}
	rem := r.limit - e.count
	if rem < 0 {
		return 0
	}
	return rem
}

// Reset clears the rate limit state for a key.
func (r *RateLimiter) Reset(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, key)
}
