package store

import (
	"testing"
)

func TestDequePushFrontAndPopFront(t *testing.T) {
	dm := NewDequeManager()
	dm.PushFront("q", "a")
	dm.PushFront("q", "b")
	val, ok := dm.PopFront("q")
	if !ok || val != "b" {
		t.Fatalf("expected b, got %s", val)
	}
}

func TestDequePushBackAndPopBack(t *testing.T) {
	dm := NewDequeManager()
	dm.PushBack("q", "x")
	dm.PushBack("q", "y")
	val, ok := dm.PopBack("q")
	if !ok || val != "y" {
		t.Fatalf("expected y, got %s", val)
	}
}

func TestDequePopFrontEmpty(t *testing.T) {
	dm := NewDequeManager()
	_, ok := dm.PopFront("empty")
	if ok {
		t.Fatal("expected false for empty deque")
	}
}

func TestDequePopBackEmpty(t *testing.T) {
	dm := NewDequeManager()
	_, ok := dm.PopBack("empty")
	if ok {
		t.Fatal("expected false for empty deque")
	}
}

func TestDequeLen(t *testing.T) {
	dm := NewDequeManager()
	dm.PushBack("q", "a")
	dm.PushBack("q", "b")
	dm.PushFront("q", "c")
	if dm.Len("q") != 3 {
		t.Fatalf("expected len 3, got %d", dm.Len("q"))
	}
}

func TestDequeRange(t *testing.T) {
	dm := NewDequeManager()
	dm.PushBack("q", "1")
	dm.PushBack("q", "2")
	dm.PushBack("q", "3")
	items := dm.Range("q")
	if len(items) != 3 || items[0] != "1" || items[2] != "3" {
		t.Fatalf("unexpected range result: %v", items)
	}
}

func TestDequeMixedOps(t *testing.T) {
	dm := NewDequeManager()
	dm.PushBack("q", "middle")
	dm.PushFront("q", "front")
	dm.PushBack("q", "back")
	items := dm.Range("q")
	if items[0] != "front" || items[1] != "middle" || items[2] != "back" {
		t.Fatalf("unexpected order: %v", items)
	}
}
