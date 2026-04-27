package store

import (
	"testing"
)

func TestLRUSetAndGet(t *testing.T) {
	c := NewLRUCache(3)
	c.Set("a", "1")
	c.Set("b", "2")

	v, ok := c.Get("a")
	if !ok || v != "1" {
		t.Fatalf("expected '1', got %q (ok=%v)", v, ok)
	}
}

func TestLRUGetMissing(t *testing.T) {
	c := NewLRUCache(3)
	_, ok := c.Get("missing")
	if ok {
		t.Fatal("expected miss for unknown key")
	}
}

func TestLRUEvictsLRUEntry(t *testing.T) {
	c := NewLRUCache(3)
	c.Set("a", "1")
	c.Set("b", "2")
	c.Set("c", "3")

	// Access "a" and "b" to make "c" the least recently used... wait, insertion order:
	// LRU is "a" (oldest). Access "b" to push "a" further back.
	// Actually after set order: front=c, b, a=back
	// Get "a" to move it to front: front=a, c, b=back
	c.Get("a")
	// Now insert "d" — should evict "b" (back)
	c.Set("d", "4")

	if c.Len() != 3 {
		t.Fatalf("expected len 3, got %d", c.Len())
	}
	_, ok := c.Get("b")
	if ok {
		t.Fatal("expected 'b' to be evicted")
	}
}

func TestLRUUpdateExisting(t *testing.T) {
	c := NewLRUCache(2)
	c.Set("x", "old")
	c.Set("x", "new")

	v, ok := c.Get("x")
	if !ok || v != "new" {
		t.Fatalf("expected 'new', got %q", v)
	}
	if c.Len() != 1 {
		t.Fatalf("expected len 1, got %d", c.Len())
	}
}

func TestLRUDelete(t *testing.T) {
	c := NewLRUCache(3)
	c.Set("a", "1")
	c.Delete("a")

	_, ok := c.Get("a")
	if ok {
		t.Fatal("expected key to be deleted")
	}
	if c.Len() != 0 {
		t.Fatalf("expected len 0, got %d", c.Len())
	}
}

func TestLRUCapacityOne(t *testing.T) {
	c := NewLRUCache(1)
	c.Set("a", "1")
	c.Set("b", "2")

	_, ok := c.Get("a")
	if ok {
		t.Fatal("expected 'a' to be evicted after capacity-1 insert")
	}
	v, ok := c.Get("b")
	if !ok || v != "2" {
		t.Fatalf("expected 'b'='2', got %q ok=%v", v, ok)
	}
}
