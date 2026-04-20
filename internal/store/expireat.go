package store

import (
	"errors"
	"time"
)

var ErrKeyNotFound = errors.New("key not found")

// ExpireAt sets an absolute Unix timestamp (seconds) as the expiry for a key.
func (s *Store) ExpireAt(key string, unixSec int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.data[key]
	if !ok {
		return ErrKeyNotFound
	}

	t := time.Unix(unixSec, 0)
	item.Expiry = &t
	s.data[key] = item
	return nil
}

// PExpireAt sets an absolute Unix timestamp in milliseconds as the expiry for a key.
func (s *Store) PExpireAt(key string, unixMs int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.data[key]
	if !ok {
		return ErrKeyNotFound
	}

	t := time.UnixMilli(unixMs)
	item.Expiry = &t
	s.data[key] = item
	return nil
}

// ExpireTime returns the absolute expiry time of a key as a Unix timestamp in seconds.
// Returns -1 if the key has no expiry, -2 if the key does not exist.
func (s *Store) ExpireTime(key string) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.data[key]
	if !ok {
		return -2
	}
	if item.Expiry == nil {
		return -1
	}
	if time.Now().After(*item.Expiry) {
		return -2
	}
	return item.Expiry.Unix()
}
