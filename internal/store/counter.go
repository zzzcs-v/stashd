package store

import "fmt"

// IncrBy increments a key's integer value by delta. If the key doesn't exist,
// it is initialized to 0 before incrementing. Returns the new value.
func (s *Store) IncrBy(key string, delta int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var current int64
	if item, ok := s.items[key]; ok && !item.isExpired() {
		var parsed int64
		_, err := fmt.Sscanf(item.Value, "%d", &parsed)
		if err != nil {
			return 0, fmt.Errorf("value at key %q is not an integer", key)
		}
		current = parsed
	}

	newVal := current + delta
	newStr := fmt.Sprintf("%d", newVal)

	var ttl int64
	if item, ok := s.items[key]; ok {
		ttl = item.TTL
	}
	s.setLocked(key, newStr, ttl)
	return newVal, nil
}

// Incr increments a key's integer value by 1.
func (s *Store) Incr(key string) (int64, error) {
	return s.IncrBy(key, 1)
}

// Decr decrements a key's integer value by 1.
func (s *Store) Decr(key string) (int64, error) {
	return s.IncrBy(key, -1)
}
