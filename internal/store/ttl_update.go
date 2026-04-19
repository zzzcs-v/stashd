package store

import "time"

// Touch resets the TTL of an existing key. Returns false if the key doesn't exist or is expired.
func (s *Store) Touch(key string, ttl time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[key]
	if !ok {
		return false
	}

	// Check if already expired
	if item.Expiry != nil && time.Now().After(*item.Expiry) {
		delete(s.items, key)
		return false
	}

	if ttl > 0 {
		expiry := time.Now().Add(ttl)
		item.Expiry = &expiry
	} else {
		item.Expiry = nil
	}

	s.items[key] = item
	return true
}

// TTL returns the remaining time-to-live for a key.
// Returns -1 if the key has no expiry, -2 if the key does not exist or is expired.
func (s *Store) TTL(key string) time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[key]
	if !ok {
		return -2
	}

	if item.Expiry == nil {
		return -1
	}

	remaining := time.Until(*item.Expiry)
	if remaining <= 0 {
		return -2
	}

	return remaining
}
