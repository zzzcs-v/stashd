package store

import (
	"fmt"
	"sync"
	"time"
)

// TypeRegistry tracks the declared type of each key so mixed-type
// operations can be rejected before they corrupt data.

type ValueType string

const (
	TypeString    ValueType = "string"
	TypeList      ValueType = "list"
	TypeSet       ValueType = "set"
	TypeHash      ValueType = "hash"
	TypeSortedSet ValueType = "zset"
	TypeCounter   ValueType = "counter"
)

type typeEntry struct {
	vt      ValueType
	expiry  time.Time
	hasExp  bool
}

// TypeRegistry stores key → ValueType associations.
type TypeRegistry struct {
	mu      sync.RWMutex
	entries map[string]typeEntry
}

func NewTypeRegistry() *TypeRegistry {
	return &TypeRegistry{entries: make(map[string]typeEntry)}
}

// Set registers or updates the type for a key.
func (r *TypeRegistry) Set(key string, vt ValueType, ttl time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	e := typeEntry{vt: vt}
	if ttl > 0 {
		e.expiry = time.Now().Add(ttl)
		e.hasExp = true
	}
	r.entries[key] = e
}

// Get returns the type for a key, or an error if missing / expired.
func (r *TypeRegistry) Get(key string) (ValueType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[key]
	if !ok {
		return "", fmt.Errorf("key %q not found in type registry", key)
	}
	if e.hasExp && time.Now().After(e.expiry) {
		return "", fmt.Errorf("key %q has expired", key)
	}
	return e.vt, nil
}

// Assert returns nil if the key matches the expected type (or does not
// exist yet), otherwise returns a type-mismatch error.
func (r *TypeRegistry) Assert(key string, expected ValueType) error {
	vt, err := r.Get(key)
	if err != nil {
		// key absent or expired — allow the operation
		return nil
	}
	if vt != expected {
		return fmt.Errorf("WRONGTYPE: key %q holds a %s value, not %s", key, vt, expected)
	}
	return nil
}

// Delete removes a key from the registry.
func (r *TypeRegistry) Delete(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, key)
}

// Keys returns all non-expired registered keys.
func (r *TypeRegistry) Keys() map[string]ValueType {
	r.mu.RLock()
	defer r.mu.RUnlock()
	now := time.Now()
	out := make(map[string]ValueType, len(r.entries))
	for k, e := range r.entries {
		if e.hasExp && now.After(e.expiry) {
			continue
		}
		out[k] = e.vt
	}
	return out
}
