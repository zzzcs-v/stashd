package store

import (
	"sort"
	"strings"
	"time"
)

// PrefixScanResult holds a key-value pair from a prefix scan.
type PrefixScanResult struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// PrefixScan returns all non-expired keys matching the given prefix,
// along with their values, optionally limited to `limit` results.
// Pass limit <= 0 for no limit.
func (s *Store) PrefixScan(prefix string, limit int) []PrefixScanResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	results := make([]PrefixScanResult, 0)

	for k, entry := range s.data {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		if entry.Expiry != nil && entry.Expiry.Before(now) {
			continue
		}
		results = append(results, PrefixScanResult{Key: k, Value: entry.Value})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	if limit > 0 && len(results) > limit {
		return results[:limit]
	}
	return results
}
