package store

import (
	"errors"
	"sync"
)

// HashMap stores a map of field->value under a single key.
type HashMap struct {
	mu   sync.RWMutex
	maps map[string]map[string]string
}

func NewHashMap() *HashMap {
	return &HashMap{
		maps: make(map[string]map[string]string),
	}
}

// HSet sets a field in the hash stored at key.
func (h *HashMap) HSet(key, field, value string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.maps[key]; !ok {
		h.maps[key] = make(map[string]string)
	}
	h.maps[key][field] = value
}

// HGet returns the value of a field in the hash stored at key.
func (h *HashMap) HGet(key, field string) (string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	m, ok := h.maps[key]
	if !ok {
		return "", errors.New("key not found")
	}
	v, ok := m[field]
	if !ok {
		return "", errors.New("field not found")
	}
	return v, nil
}

// HGetAll returns all field-value pairs for a key.
func (h *HashMap) HGetAll(key string) (map[string]string, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	m, ok := h.maps[key]
	if !ok {
		return nil, errors.New("key not found")
	}
	copy := make(map[string]string, len(m))
	for k, v := range m {
		copy[k] = v
	}
	return copy, nil
}

// HDel removes a field from the hash stored at key.
func (h *HashMap) HDel(key, field string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	m, ok := h.maps[key]
	if !ok {
		return errors.New("key not found")
	}
	delete(m, field)
	if len(m) == 0 {
		delete(h.maps, key)
	}
	return nil
}

// HExists checks whether a field exists in the hash at key.
func (h *HashMap) HExists(key, field string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	m, ok := h.maps[key]
	if !ok {
		return false
	}
	_, ok = m[field]
	return ok
}
