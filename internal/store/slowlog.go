package store

import (
	"sync"
	"time"
)

// SlowLogEntry represents a single slow command entry.
type SlowLogEntry struct {
	ID        int64
	Timestamp time.Time
	Duration  time.Duration
	Command   string
	Args      []string
}

// SlowLog records commands that exceed a configurable threshold.
type SlowLog struct {
	mu        sync.Mutex
	entries   []SlowLogEntry
	maxLen    int
	threshold time.Duration
	nextID    int64
}

// NewSlowLog creates a new SlowLog with the given max length and threshold.
func NewSlowLog(maxLen int, threshold time.Duration) *SlowLog {
	if maxLen <= 0 {
		maxLen = 128
	}
	return &SlowLog{
		maxLen:    maxLen,
		threshold: threshold,
		entries:   make([]SlowLogEntry, 0, maxLen),
	}
}

// Record adds an entry if the duration exceeds the threshold.
func (s *SlowLog) Record(command string, args []string, duration time.Duration) {
	if duration < s.threshold {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := SlowLogEntry{
		ID:        s.nextID,
		Timestamp: time.Now(),
		Duration:  duration,
		Command:   command,
		Args:      args,
	}
	s.nextID++
	s.entries = append([]SlowLogEntry{entry}, s.entries...)
	if len(s.entries) > s.maxLen {
		s.entries = s.entries[:s.maxLen]
	}
}

// Get returns up to n recent slow log entries. If n <= 0, returns all.
func (s *SlowLog) Get(n int) []SlowLogEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	if n <= 0 || n > len(s.entries) {
		n = len(s.entries)
	}
	result := make([]SlowLogEntry, n)
	copy(result, s.entries[:n])
	return result
}

// Reset clears all slow log entries.
func (s *SlowLog) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = s.entries[:0]
}

// Len returns the current number of entries.
func (s *SlowLog) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.entries)
}
