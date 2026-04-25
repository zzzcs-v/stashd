package store

import "time"

// SetNX sets a key only if it does not already exist.
// Returns true if the key was set, false if it already existed.
func (s *Store) SetNX(key, value string, ttl time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if item, ok := s.data[key]; ok {
		// Key exists and has not expired
		if item.Expiry.IsZero() || item.Expiry.After(time.Now()) {
			return false
		}
	}

	entry := Item{Value: value}
	if ttl > 0 {
		entry.Expiry = time.Now().Add(ttl)
	}
	s.data[key] = entry
	return true
}

// GetSet atomically sets key to value and returns the old value.
// Returns the old value and true if the key existed, or empty string and false if not.
func (s *Store) GetSet(key, value string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var old string
	found := false

	if item, ok := s.data[key]; ok {
		if item.Expiry.IsZero() || item.Expiry.After(time.Now()) {
			old = item.Value
			found = true
		}
	}

	s.data[key] = Item{Value: value}
	return old, found
}

// SetXX sets a key only if it already exists.
// Returns true if the key was updated, false if it did not exist.
func (s *Store) SetXX(key, value string, ttl time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.data[key]
	if !ok {
		return false
	}
	if !item.Expiry.IsZero() && !item.Expiry.After(time.Now()) {
		return false
	}

	entry := Item{Value: value}
	if ttl > 0 {
		entry.Expiry = time.Now().Add(ttl)
	}
	s.data[key] = entry
	return true
}
