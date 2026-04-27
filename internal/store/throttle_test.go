package store

import (
	"testing"
	"time"
)

func TestThrottleAllow(t *testing.T) {
	tm := NewThrottleManager()
	allowed, remaining, _ := tm.Allow("user:1", 3, time.Minute)
	if !allowed {
		t.Fatal("expected first request to be allowed")
	}
	if remaining != 2 {
		t.Fatalf("expected 2 remaining, got %d", remaining)
	}
}

func TestThrottleDeny(t *testing.T) {
	tm := NewThrottleManager()
	for i := 0; i < 3; i++ {
		tm.Allow("user:2", 3, time.Minute)
	}
	allowed, remaining, _ := tm.Allow("user:2", 3, time.Minute)
	if allowed {
		t.Fatal("expected request to be denied after limit")
	}
	if remaining != 0 {
		t.Fatalf("expected 0 remaining, got %d", remaining)
	}
}

func TestThrottleWindowExpiry(t *testing.T) {
	tm := NewThrottleManager()
	for i := 0; i < 2; i++ {
		tm.Allow("user:3", 2, 50*time.Millisecond)
	}
	allowed, _, _ := tm.Allow("user:3", 2, 50*time.Millisecond)
	if allowed {
		t.Fatal("expected deny before window expires")
	}
	time.Sleep(60 * time.Millisecond)
	allowed, _, _ = tm.Allow("user:3", 2, 50*time.Millisecond)
	if !allowed {
		t.Fatal("expected allow after window reset")
	}
}

func TestThrottleReset(t *testing.T) {
	tm := NewThrottleManager()
	tm.Allow("user:4", 1, time.Minute)
	if err := tm.Reset("user:4"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	allowed, _, _ := tm.Allow("user:4", 1, time.Minute)
	if !allowed {
		t.Fatal("expected allow after reset")
	}
}

func TestThrottleResetMissing(t *testing.T) {
	tm := NewThrottleManager()
	if err := tm.Reset("ghost"); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestThrottleStatus(t *testing.T) {
	tm := NewThrottleManager()
	tm.Allow("user:5", 5, time.Minute)
	tm.Allow("user:5", 5, time.Minute)
	count, _, active := tm.Status("user:5")
	if !active {
		t.Fatal("expected active throttle")
	}
	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func TestThrottleStatusMissing(t *testing.T) {
	tm := NewThrottleManager()
	_, _, active := tm.Status("nobody")
	if active {
		t.Fatal("expected inactive for unknown key")
	}
}
