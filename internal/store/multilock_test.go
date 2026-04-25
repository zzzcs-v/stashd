package store

import (
	"sync"
	"testing"
	"time"
)

func TestMultiLockBasic(t *testing.T) {
	m := NewMultiLockManager()
	token, unlock := m.Lock([]string{"a", "b", "c"})
	if token == "" {
		t.Fatal("expected non-empty token")
	}
	unlock()
}

func TestMultiLockEmptyKeys(t *testing.T) {
	m := NewMultiLockManager()
	token, unlock := m.Lock([]string{})
	if token != "" {
		t.Errorf("expected empty token for empty keys, got %q", token)
	}
	unlock() // should not panic
}

func TestMultiLockDeduplicatesKeys(t *testing.T) {
	m := NewMultiLockManager()
	_, unlock := m.Lock([]string{"x", "x", "x"})
	done := make(chan struct{})
	go func() {
		// after unlock, another goroutine should be able to lock the same key
		unlock()
		_, u2 := m.Lock([]string{"x"})
		u2()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("deadlock detected")
	}
}

func TestMultiLockMutualExclusion(t *testing.T) {
	m := NewMultiLockManager()
	var counter int
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, unlock := m.Lock([]string{"shared", "resource"})
			defer unlock()
			counter++
		}()
	}
	wg.Wait()
	if counter != 10 {
		t.Errorf("expected counter=10, got %d", counter)
	}
}

func TestMultiLockSortedAcquisition(t *testing.T) {
	m := NewMultiLockManager()
	done := make(chan struct{})

	go func() {
		_, u1 := m.Lock([]string{"alpha", "beta", "gamma"})
		time.Sleep(10 * time.Millisecond)
		u1()
	}()

	go func() {
		time.Sleep(2 * time.Millisecond)
		_, u2 := m.Lock([]string{"gamma", "alpha", "beta"})
		u2()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("deadlock detected with sorted acquisition")
	}
}
