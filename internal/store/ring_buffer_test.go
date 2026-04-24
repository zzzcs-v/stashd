package store

import (
	"testing"
)

func TestRingBufferPushAndPop(t *testing.T) {
	m := NewRingBufferManager()
	if err := m.Push("rb", "a", 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m.Push("rb", "b", 3)
	val, err := m.Pop("rb")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "a" {
		t.Errorf("expected 'a', got %q", val)
	}
}

func TestRingBufferFull(t *testing.T) {
	m := NewRingBufferManager()
	m.Push("rb", "x", 2)
	m.Push("rb", "y", 2)
	err := m.Push("rb", "z", 2)
	if err != ErrRingBufferFull {
		t.Errorf("expected ErrRingBufferFull, got %v", err)
	}
}

func TestRingBufferPopEmpty(t *testing.T) {
	m := NewRingBufferManager()
	_, err := m.Pop("missing")
	if err != ErrRingBufferEmpty {
		t.Errorf("expected ErrRingBufferEmpty, got %v", err)
	}
}

func TestRingBufferLen(t *testing.T) {
	m := NewRingBufferManager()
	if m.Len("rb") != 0 {
		t.Error("expected len 0 for missing key")
	}
	m.Push("rb", "a", 4)
	m.Push("rb", "b", 4)
	if m.Len("rb") != 2 {
		t.Errorf("expected len 2, got %d", m.Len("rb"))
	}
}

func TestRingBufferFIFOOrder(t *testing.T) {
	m := NewRingBufferManager()
	vals := []string{"first", "second", "third"}
	for _, v := range vals {
		m.Push("rb", v, 5)
	}
	for _, expected := range vals {
		got, err := m.Pop("rb")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	}
}

func TestRingBufferCapacity(t *testing.T) {
	m := NewRingBufferManager()
	if m.Capacity("rb") != 0 {
		t.Error("expected capacity 0 for missing key")
	}
	m.Push("rb", "v", 10)
	if m.Capacity("rb") != 10 {
		t.Errorf("expected capacity 10, got %d", m.Capacity("rb"))
	}
}
