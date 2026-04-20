package store

import (
	"testing"
)

func TestQueuePushAndPop(t *testing.T) {
	q := NewQueue()
	q.Push("mylist", "a")
	q.Push("mylist", "b")

	v, err := q.Pop("mylist")
	if err != nil || v != "a" {
		t.Fatalf("expected 'a', got '%s' err=%v", v, err)
	}
	v, err = q.Pop("mylist")
	if err != nil || v != "b" {
		t.Fatalf("expected 'b', got '%s' err=%v", v, err)
	}
}

func TestQueuePopEmpty(t *testing.T) {
	q := NewQueue()
	_, err := q.Pop("missing")
	if err != ErrQueueEmpty {
		t.Fatalf("expected ErrQueueEmpty, got %v", err)
	}
}

func TestQueueLen(t *testing.T) {
	q := NewQueue()
	q.Push("k", "x")
	q.Push("k", "y")
	if q.Len("k") != 2 {
		t.Fatalf("expected length 2")
	}
	q.Pop("k")
	if q.Len("k") != 1 {
		t.Fatalf("expected length 1")
	}
}

func TestQueuePeek(t *testing.T) {
	q := NewQueue()
	q.Push("p", "hello")
	v, err := q.Peek("p")
	if err != nil || v != "hello" {
		t.Fatalf("expected 'hello', got '%s' err=%v", v, err)
	}
	if q.Len("p") != 1 {
		t.Fatal("peek should not remove item")
	}
}

func TestQueuePeekEmpty(t *testing.T) {
	q := NewQueue()
	_, err := q.Peek("nope")
	if err != ErrQueueEmpty {
		t.Fatalf("expected ErrQueueEmpty, got %v", err)
	}
}
