package store

import (
	"encoding/json"
	"os"
	"time"
)

type snapshotEntry struct {
	Value     string     `json:"value"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// SaveSnapshot writes the current store contents to a file.
func (s *Store) SaveSnapshot(path string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := make(map[string]snapshotEntry, len(s.data))
	for k, v := range s.data {
		entry := snapshotEntry{Value: v.value}
		if v.expiresAt != nil {
			t := *v.expiresAt
			entry.ExpiresAt = &t
		}
		data[k] = entry
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(data)
}

// LoadSnapshot reads a snapshot file and restores store contents.
func (s *Store) LoadSnapshot(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	var data map[string]snapshotEntry
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return err
	}

	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, v := range data {
		if v.ExpiresAt != nil && v.ExpiresAt.Before(now) {
			continue // skip already expired entries
		}
		s.data[k] = entry{value: v.Value, expiresAt: v.ExpiresAt}
	}

	return nil
}
