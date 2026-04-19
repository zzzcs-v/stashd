package store

import (
	"testing"
	"time"
)

func TestStatsEmpty(t *testing.T) {
	s := New()
	stats := s.Stats()
	if stats.TotalKeys != 0 {
		t.Errorf("expected 0 total keys, got %d", stats.TotalKeys)
	}
}

func TestStatsWithKeys(t *testing.T) {
	s := New()
	s.Set("a", "1", 0)
	s.Set("b", "2", 0)
	stats := s.Stats()
	if stats.TotalKeys != 2 {
		t.Errorf("expected 2 total keys, got %d", stats.TotalKeys)
	}
}

func TestStatsExpiredKeys(t *testing.T) {
	s := New()
	s.Set("x", "val", 1*time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	stats := s.Stats()
	if stats.ExpiredKeys != 1 {
		t.Errorf("expected 1 expired key, got %d", stats.ExpiredKeys)
	}
}

func TestStatsUptime(t *testing.T) {
	s := New()
	time.Sleep(10 * time.Millisecond)
	stats := s.Stats()
	if stats.UptimeSeconds < 0 {
		t.Error("uptime should be non-negative")
	}
}
