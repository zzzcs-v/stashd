package store

import (
	"testing"
	"time"
)

func TestSlowLogRecordAboveThreshold(t *testing.T) {
	sl := NewSlowLog(10, 10*time.Millisecond)
	sl.Record("SET", []string{"key", "val"}, 20*time.Millisecond)
	if sl.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", sl.Len())
	}
	entries := sl.Get(1)
	if entries[0].Command != "SET" {
		t.Errorf("expected command SET, got %s", entries[0].Command)
	}
}

func TestSlowLogRecordBelowThreshold(t *testing.T) {
	sl := NewSlowLog(10, 50*time.Millisecond)
	sl.Record("GET", []string{"key"}, 5*time.Millisecond)
	if sl.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", sl.Len())
	}
}

func TestSlowLogMaxLen(t *testing.T) {
	sl := NewSlowLog(3, time.Millisecond)
	for i := 0; i < 5; i++ {
		sl.Record("SET", []string{"k"}, 10*time.Millisecond)
	}
	if sl.Len() != 3 {
		t.Fatalf("expected max 3 entries, got %d", sl.Len())
	}
}

func TestSlowLogGetLimitedN(t *testing.T) {
	sl := NewSlowLog(10, time.Millisecond)
	for i := 0; i < 5; i++ {
		sl.Record("DEL", []string{"k"}, 5*time.Millisecond)
	}
	entries := sl.Get(2)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestSlowLogReset(t *testing.T) {
	sl := NewSlowLog(10, time.Millisecond)
	sl.Record("SET", []string{"a", "b"}, 10*time.Millisecond)
	sl.Reset()
	if sl.Len() != 0 {
		t.Fatalf("expected 0 entries after reset, got %d", sl.Len())
	}
}

func TestSlowLogIDIncrement(t *testing.T) {
	sl := NewSlowLog(10, time.Millisecond)
	sl.Record("GET", []string{"x"}, 5*time.Millisecond)
	sl.Record("SET", []string{"x", "1"}, 5*time.Millisecond)
	entries := sl.Get(0)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	// Most recent first
	if entries[0].ID <= entries[1].ID {
		t.Errorf("expected entries in descending ID order")
	}
}

func TestSlowLogGetAllWhenNIsZero(t *testing.T) {
	sl := NewSlowLog(10, time.Millisecond)
	sl.Record("INCR", []string{"c"}, 2*time.Millisecond)
	sl.Record("INCR", []string{"c"}, 2*time.Millisecond)
	sl.Record("INCR", []string{"c"}, 2*time.Millisecond)
	entries := sl.Get(0)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries with n=0, got %d", len(entries))
	}
}
