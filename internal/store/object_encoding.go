package store

import (
	"fmt"
	"strconv"
	"sync"
)

// Encoding represents the internal encoding type of a stored value.
type Encoding string

const (
	EncodingInt    Encoding = "int"
	EncodingFloat  Encoding = "float"
	EncodingString Encoding = "string"
	EncodingList   Encoding = "list"
	EncodingSet    Encoding = "set"
	EncodingHash   Encoding = "hash"
	EncodingNone   Encoding = "none"
)

// ObjectEncodingManager inspects values and reports their internal encoding.
type ObjectEncodingManager struct {
	mu    sync.RWMutex
	store map[string]string
}

// NewObjectEncodingManager creates a new ObjectEncodingManager.
func NewObjectEncodingManager() *ObjectEncodingManager {
	return &ObjectEncodingManager{
		store: make(map[string]string),
	}
}

// Set stores a raw string value under the given key.
func (o *ObjectEncodingManager) Set(key, value string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.store[key] = value
}

// Get retrieves the raw string value for a key.
func (o *ObjectEncodingManager) Get(key string) (string, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	v, ok := o.store[key]
	return v, ok
}

// Delete removes a key from the manager.
func (o *ObjectEncodingManager) Delete(key string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	delete(o.store, key)
}

// Encoding returns the detected encoding for the value stored at key.
func (o *ObjectEncodingManager) Encoding(key string) (Encoding, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	v, ok := o.store[key]
	if !ok {
		return EncodingNone, fmt.Errorf("key not found: %s", key)
	}

	if _, err := strconv.ParseInt(v, 10, 64); err == nil {
		return EncodingInt, nil
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return EncodingFloat, nil
	}
	return EncodingString, nil
}
