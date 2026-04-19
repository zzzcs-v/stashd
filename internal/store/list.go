package store

import "strings"

// List returns all non-expired keys, optionally filtered by prefix.
func (s *Store) List(prefix string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := s.now()
	keys := make([]string, 0)

	for k, entry := range s.data {
		if !entry.Expiry.IsZero() && now.After(entry.Expiry) {
			continue
		}
		if prefix == "" || strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}

	return keys
}
