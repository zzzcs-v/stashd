package store

import "sync"

// LockManager provides named key-level locking.
type LockManager struct {
	mu    sync.Mutex
	locks map[string]*keyLock
}

type keyLock struct {
	mu      sync.Mutex
	holders int
}

func NewLockManager() *LockManager {
	return &LockManager{locks: make(map[string]*keyLock)}
}

func (lm *LockManager) acquire(key string) *keyLock {
	lm.mu.Lock()
	kl, ok := lm.locks[key]
	if !ok {
		kl = &keyLock{}
		lm.locks[key] = kl
	}
	kl.holders++
	lm.mu.Unlock()
	kl.mu.Lock()
	return kl
}

func (lm *LockManager) release(key string, kl *keyLock) {
	kl.mu.Unlock()
	lm.mu.Lock()
	kl.holders--
	if kl.holders == 0 {
		delete(lm.locks, key)
	}
	lm.mu.Unlock()
}

// Lock acquires the lock for the given key and returns a release function.
func (lm *LockManager) Lock(key string) func() {
	kl := lm.acquire(key)
	return func() { lm.release(key, kl) }
}

// Len returns the number of currently tracked key locks.
func (lm *LockManager) Len() int {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	return len(lm.locks)
}
