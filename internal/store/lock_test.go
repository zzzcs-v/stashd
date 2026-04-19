package store

import (
	"sync"
	"testing"
	"time"
)

func TestLockBasic(t *testing.T) {
	lm := NewLockManager()
	unlock := lm.Lock("foo")
	if lm.Len() != 1 {
		t.Fatalf("expected 1 lock, got %d", lm.Len())
	}
	unlock()
	if lm.Len() != 0 {
		t.Fatalf("expected 0 locks after release, got %d", lm.Len())
	}
}

func TestLockMutualExclusion(t *testing.T) {
	lm := NewLockManager()
	var order []int
	var mu sync.Mutex

	unlock1 := lm.Lock("bar")

	started := make(chan struct{})
	done := make(chan struct{})
	go func() {
		close(started)
		unlock2 := lm.Lock("bar")
		mu.Lock()
		order = append(order, 2)
		mu.Unlock()
		unlock2()
		close(done)
	}()

	<-started
	time.Sleep(20 * time.Millisecond)
	mu.Lock()
	order = append(order, 1)
	mu.Unlock()
	unlock1()
	<-done

	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Fatalf("expected order [1 2], got %v", order)
	}
}

func TestLockDifferentKeys(t *testing.T) {
	lm := NewLockManager()
	u1 := lm.Lock("a")
	u2 := lm.Lock("b")
	if lm.Len() != 2 {
		t.Fatalf("expected 2 locks, got %d", lm.Len())
	}
	u1()
	u2()
	if lm.Len() != 0 {
		t.Fatalf("expected 0 locks, got %d", lm.Len())
	}
}

func TestLockCleanup(t *testing.T) {
	lm := NewLockManager()
	for i := 0; i < 10; i++ {
		unlock := lm.Lock("key")
		unlock()
	}
	if lm.Len() != 0 {
		t.Fatalf("expected no lingering locks, got %d", lm.Len())
	}
}
