package store

import "time"

// StartEviction launches a background goroutine that periodically removes
// expired keys from the store at the given interval.
func (s *Store) StartEviction(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.evictExpired()
			case <-s.quit:
				return
			}
		}
	}()
}

// StopEviction signals the eviction goroutine to stop.
func (s *Store) StopEviction() {
	close(s.quit)
}

// evictExpired removes all keys whose TTL has passed.
func (s *Store) evictExpired() {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range s.data {
		if !v.Expiry.IsZero() && now.After(v.Expiry) {
			delete(s.data, k)
		}
	}
}
