package store

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// JSONDocManager stores JSON documents and supports field-level get/set via dot paths.
type JSONDocManager struct {
	mu   sync.RWMutex
	docs map[string]map[string]interface{}
}

func NewJSONDocManager() *JSONDocManager {
	return &JSONDocManager{
		docs: make(map[string]map[string]interface{}),
	}
}

// JSONSet stores a JSON document under key.
func (j *JSONDocManager) JSONSet(key string, value []byte) error {
	var doc map[string]interface{}
	if err := json.Unmarshal(value, &doc); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	j.mu.Lock()
	defer j.mu.Unlock()
	j.docs[key] = doc
	return nil
}

// JSONGet retrieves the full document or a nested field via dot-path.
func (j *JSONDocManager) JSONGet(key, path string) (interface{}, bool) {
	j.mu.RLock()
	defer j.mu.RUnlock()
	doc, ok := j.docs[key]
	if !ok {
		return nil, false
	}
	if path == "" || path == "." {
		return doc, true
	}
	return resolvePath(doc, path)
}

// JSONDel removes a document by key.
func (j *JSONDocManager) JSONDel(key string) bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	_, ok := j.docs[key]
	delete(j.docs, key)
	return ok
}

// resolvePath walks a dot-separated path into a nested map.
func resolvePath(doc map[string]interface{}, path string) (interface{}, bool) {
	parts := strings.SplitN(path, ".", 2)
	val, ok := doc[parts[0]]
	if !ok {
		return nil, false
	}
	if len(parts) == 1 {
		return val, true
	}
	nested, ok := val.(map[string]interface{})
	if !ok {
		return nil, false
	}
	return resolvePath(nested, parts[1])
}
