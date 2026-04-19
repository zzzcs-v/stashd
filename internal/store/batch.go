package store

import "time"

// BatchSetItem represents a single item in a batch set operation.
type BatchSetItem struct {
	Key   string
	Value string
	TTL   time.Duration
}

// BatchGetResult holds the result for a single key in a batch get.
type BatchGetResult struct {
	Key   string
	Value string
	Found bool
}

// BatchSet sets multiple keys at once.
func (s *Store) BatchSet(items []BatchSetItem) {
	for _, item := range items {
		s.Set(item.Key, item.Value, item.TTL)
	}
}

// BatchGet retrieves multiple keys at once.
func (s *Store) BatchGet(keys []string) []BatchGetResult {
	results := make([]BatchGetResult, len(keys))
	for i, key := range keys {
		val, ok := s.Get(key)
		results[i] = BatchGetResult{Key: key, Value: val, Found: ok}
	}
	return results
}

// BatchDelete deletes multiple keys at once.
func (s *Store) BatchDelete(keys []string) {
	for _, key := range keys {
		s.Delete(key)
	}
}
