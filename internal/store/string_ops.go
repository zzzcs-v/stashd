package store

import (
	"errors"
	"strings"
)

// Append appends a value to an existing string key, or creates it if missing.
// Returns the new length of the value.
func (s *Store) Append(key, value string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing := ""
	if item, ok := s.data[key]; ok && !item.isExpired() {
		existing = item.Value
	}

	newVal := existing + value
	s.data[key] = item{Value: newVal, Expiry: s.data[key].Expiry}
	return len(newVal), nil
}

// GetRange returns a substring of the value stored at key.
// start and end are inclusive byte offsets. Negative indices count from the end.
func (s *Store) GetRange(key string, start, end int) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	it, ok := s.data[key]
	if !ok || it.isExpired() {
		return "", nil
	}

	v := it.Value
	n := len(v)
	if n == 0 {
		return "", nil
	}

	if start < 0 {
		start = n + start
	}
	if end < 0 {
		end = n + end
	}
	if start < 0 {
		start = 0
	}
	if end >= n {
		end = n - 1
	}
	if start > end {
		return "", nil
	}
	return v[start : end+1], nil
}

// StrLen returns the length of the string stored at key, or 0 if missing/expired.
func (s *Store) StrLen(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	it, ok := s.data[key]
	if !ok || it.isExpired() {
		return 0
	}
	return len(it.Value)
}

// GetSet sets a new value and returns the old one. Returns error if key missing.
func (s *Store) GetSet(key, value string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	old := ""
	if it, ok := s.data[key]; ok && !it.isExpired() {
		old = it.Value
	} else {
		return "", errors.New("key not found")
	}

	s.data[key] = item{Value: value}
	return old, nil
}

// SetNX sets a key only if it does not already exist (or is expired).
// Returns true if the key was set, false otherwise.
func (s *Store) SetNX(key, value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if it, ok := s.data[key]; ok && !it.isExpired() {
		return false
	}
	s.data[key] = item{Value: value}
	return true
}

// MSetNX sets multiple keys only if none of them exist.
// Returns true if all keys were set, false if any already existed.
func (s *Store) MSetNX(pairs map[string]string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k := range pairs {
		if it, ok := s.data[k]; ok && !it.isExpired() {
			return false
		}
	}
	for k, v := range pairs {
		s.data[k] = item{Value: strings.Clone(v)}
	}
	return true
}
