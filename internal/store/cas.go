package store

import (
	"errors"
	"sync"
	"time"
)

// ErrCASConflict is returned when the compare-and-swap value does not match.
var ErrCASConflict = errors.New("cas: value mismatch, swap not performed")

// ErrCASMissing is returned when the key does not exist during a CAS operation.
var ErrCASMissing = errors.New("cas: key does not exist")

// CASManager handles compare-and-swap operations on the underlying store.
type CASManager struct {
	mu    sync.Mutex
	store *Store
}

// NewCASManager creates a new CASManager backed by the given store.
func NewCASManager(s *Store) *CASManager {
	return &CASManager{store: s}
}

// CompareAndSwap atomically sets key to newVal only if the current value equals
// expectedVal. Returns ErrCASMissing if the key is absent or expired, and
// ErrCASConflict if the current value does not match expectedVal.
func (c *CASManager) CompareAndSwap(key, expectedVal, newVal string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	current, ok := c.store.Get(key)
	if !ok {
		return ErrCASMissing
	}
	if current != expectedVal {
		return ErrCASConflict
	}
	c.store.Set(key, newVal, ttl)
	return nil
}

// CompareAndDelete atomically deletes key only if the current value equals
// expectedVal. Returns ErrCASMissing if absent and ErrCASConflict on mismatch.
func (c *CASManager) CompareAndDelete(key, expectedVal string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	current, ok := c.store.Get(key)
	if !ok {
		return ErrCASMissing
	}
	if current != expectedVal {
		return ErrCASConflict
	}
	c.store.Delete(key)
	return nil
}
