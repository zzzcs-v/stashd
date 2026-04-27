package store

import (
	"testing"
	"time"
)

func TestTypeRegistrySetAndGet(t *testing.T) {
	r := NewTypeRegistry()
	r.Set("mykey", TypeString, 0)
	vt, err := r.Get("mykey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vt != TypeString {
		t.Fatalf("expected %q, got %q", TypeString, vt)
	}
}

func TestTypeRegistryGetMissing(t *testing.T) {
	r := NewTypeRegistry()
	_, err := r.Get("ghost")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestTypeRegistryAssertMatch(t *testing.T) {
	r := NewTypeRegistry()
	r.Set("counter", TypeCounter, 0)
	if err := r.Assert("counter", TypeCounter); err != nil {
		t.Fatalf("unexpected mismatch: %v", err)
	}
}

func TestTypeRegistryAssertMismatch(t *testing.T) {
	r := NewTypeRegistry()
	r.Set("counter", TypeCounter, 0)
	err := r.Assert("counter", TypeString)
	if err == nil {
		t.Fatal("expected WRONGTYPE error")
	}
}

func TestTypeRegistryAssertMissingKeyAllowed(t *testing.T) {
	r := NewTypeRegistry()
	// absent key should not produce an error
	if err := r.Assert("nokey", TypeHash); err != nil {
		t.Fatalf("expected nil for absent key, got: %v", err)
	}
}

func TestTypeRegistryDelete(t *testing.T) {
	r := NewTypeRegistry()
	r.Set("k", TypeList, 0)
	r.Delete("k")
	_, err := r.Get("k")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestTypeRegistryExpiredKey(t *testing.T) {
	r := NewTypeRegistry()
	r.Set("temp", TypeSet, 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	_, err := r.Get("temp")
	if err == nil {
		t.Fatal("expected error for expired key")
	}
}

func TestTypeRegistryKeys(t *testing.T) {
	r := NewTypeRegistry()
	r.Set("a", TypeString, 0)
	r.Set("b", TypeHash, 0)
	r.Set("c", TypeSortedSet, 5*time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	keys := r.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 live keys, got %d", len(keys))
	}
	if _, ok := keys["c"]; ok {
		t.Fatal("expired key 'c' should not appear in Keys()")
	}
}
