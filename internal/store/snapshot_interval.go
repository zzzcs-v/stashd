package store

import (
	"sync"
	"time"
)

// SnapshotScheduler periodically saves a snapshot of the store to disk.
type SnapshotScheduler struct {
	store    *Store
	path     string
	interval time.Duration
	stopCh   chan struct{}
	wg       sync.WaitGroup
	mu       sync.Mutex
	running  bool
}

// NewSnapshotScheduler creates a new scheduler that will save snapshots of s
// to path every interval. Call Start to begin scheduling.
func NewSnapshotScheduler(s *Store, path string, interval time.Duration) *SnapshotScheduler {
	return &SnapshotScheduler{
		store:    s,
		path:     path,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

// Start begins the periodic snapshot loop. It is a no-op if already running.
func (ss *SnapshotScheduler) Start() {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if ss.running {
		return
	}
	ss.running = true
	ss.stopCh = make(chan struct{})
	ss.wg.Add(1)
	go func() {
		defer ss.wg.Done()
		ticker := time.NewTicker(ss.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = SaveSnapshot(ss.store, ss.path)
			case <-ss.stopCh:
				return
			}
		}
	}()
}

// Stop halts the scheduler and waits for the background goroutine to exit.
func (ss *SnapshotScheduler) Stop() {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	if !ss.running {
		return
	}
	close(ss.stopCh)
	ss.wg.Wait()
	ss.running = false
}

// IsRunning reports whether the scheduler is currently active.
func (ss *SnapshotScheduler) IsRunning() bool {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.running
}
