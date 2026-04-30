package store

import (
	"testing"
	"time"
)

func TestFrequencyHitIncrementsCount(t *testing.T) {
	fm := NewFrequencyManager(5 * time.Second)

	count := fm.Hit("page:home")
	if count != 1 {
		t.Fatalf("expected 1, got %d", count)
	}
	count = fm.Hit("page:home")
	if count != 2 {
		t.Fatalf("expected 2, got %d", count)
	}
}

func TestFrequencyCountReturnsZeroForMissing(t *testing.T) {
	fm := NewFrequencyManager(5 * time.Second)

	if c := fm.Count("missing"); c != 0 {
		t.Fatalf("expected 0, got %d", c)
	}
}

func TestFrequencyWindowExpiry(t *testing.T) {
	fm := NewFrequencyManager(50 * time.Millisecond)

	fm.Hit("key")
	fm.Hit("key")

	time.Sleep(80 * time.Millisecond)

	if c := fm.Count("key"); c != 0 {
		t.Fatalf("expected 0 after window expiry, got %d", c)
	}

	// Hit again after expiry should reset to 1
	count := fm.Hit("key")
	if count != 1 {
		t.Fatalf("expected 1 after reset, got %d", count)
	}
}

func TestFrequencyReset(t *testing.T) {
	fm := NewFrequencyManager(5 * time.Second)

	fm.Hit("event")
	fm.Hit("event")

	if err := fm.Reset("event"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c := fm.Count("event"); c != 0 {
		t.Fatalf("expected 0 after reset, got %d", c)
	}
}

func TestFrequencyResetMissingKey(t *testing.T) {
	fm := NewFrequencyManager(5 * time.Second)

	if err := fm.Reset("nope"); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestFrequencyTTL(t *testing.T) {
	fm := NewFrequencyManager(2 * time.Second)

	fm.Hit("k")
	ttl, ok := fm.TTL("k")
	if !ok {
		t.Fatal("expected TTL to be present")
	}
	if ttl <= 0 || ttl > 2*time.Second {
		t.Fatalf("unexpected TTL value: %v", ttl)
	}
}

func TestFrequencyTTLMissingKey(t *testing.T) {
	fm := NewFrequencyManager(2 * time.Second)

	_, ok := fm.TTL("ghost")
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
}
