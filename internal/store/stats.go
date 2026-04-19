package store

import "time"

// Stats holds runtime statistics for the store.
type Stats struct {
	TotalKeys   int    `json:"total_keys"`
	ExpiredKeys int    `json:"expired_keys"`
	UptimeSeconds int64 `json:"uptime_seconds"`
}

// Stats returns a snapshot of current store statistics.
func (s *Store) Stats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := 0
	expired := 0
	now := time.Now()

	for _, entry := range s.data {
		total++
		if entry.Expiry != nil && entry.Expiry.Before(now) {
			expired++
		}
	}

	return Stats{
		TotalKeys:     total,
		ExpiredKeys:   expired,
		UptimeSeconds: int64(time.Since(s.startedAt).Seconds()),
	}
}
