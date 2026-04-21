package store

import (
	"testing"
	"time"
)

func TestRateLimitAllow(t *testing.T) {
	rl := NewRateLimiter(3, time.Second)
	for i := 0; i < 3; i++ {
		if !rl.Allow("client1") {
			t.Fatalf("expected allow on request %d", i+1)
		}
	}
	if rl.Allow("client1") {
		t.Fatal("expected deny after limit exceeded")
	}
}

func TestRateLimitRemaining(t *testing.T) {
	rl := NewRateLimiter(5, time.Second)
	rl.Allow("x")
	rl.Allow("x")
	if got := rl.Remaining("x"); got != 3 {
		t.Fatalf("expected 3 remaining, got %d", got)
	}
}

func TestRateLimitReset(t *testing.T) {
	rl := NewRateLimiter(2, time.Second)
	rl.Allow("y")
	rl.Allow("y")
	if rl.Allow("y") {
		t.Fatal("expected deny")
	}
	rl.Reset("y")
	if !rl.Allow("y") {
		t.Fatal("expected allow after reset")
	}
}

func TestRateLimitWindowExpiry(t *testing.T) {
	rl := NewRateLimiter(1, 50*time.Millisecond)
	rl.Allow("z")
	if rl.Allow("z") {
		t.Fatal("expected deny within window")
	}
	time.Sleep(60 * time.Millisecond)
	if !rl.Allow("z") {
		t.Fatal("expected allow after window reset")
	}
}

func TestRateLimitIsolation(t *testing.T) {
	rl := NewRateLimiter(1, time.Second)
	rl.Allow("a")
	if !rl.Allow("b") {
		t.Fatal("different keys should be independent")
	}
}

func TestRateLimitRemainingUnknownKey(t *testing.T) {
	rl := NewRateLimiter(5, time.Second)
	// A key that has never been used should report the full limit as remaining.
	if got := rl.Remaining("never-used"); got != 5 {
		t.Fatalf("expected 5 remaining for unknown key, got %d", got)
	}
}
