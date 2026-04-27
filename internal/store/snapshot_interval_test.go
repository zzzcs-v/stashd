package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSnapshotSchedulerStartStop(t *testing.T) {
	s := New()
	s.Set("k", "v", 0)

	tmp := filepath.Join(t.TempDir(), "snap.db")
	sched := NewSnapshotScheduler(s, tmp, 50*time.Millisecond)

	if sched.IsRunning() {
		t.Fatal("expected scheduler to be stopped initially")
	}

	sched.Start()
	if !sched.IsRunning() {
		t.Fatal("expected scheduler to be running after Start")
	}

	// Wait long enough for at least one snapshot tick.
	time.Sleep(120 * time.Millisecond)

	sched.Stop()
	if sched.IsRunning() {
		t.Fatal("expected scheduler to be stopped after Stop")
	}

	if _, err := os.Stat(tmp); os.IsNotExist(err) {
		t.Fatal("expected snapshot file to exist after scheduled save")
	}
}

func TestSnapshotSchedulerDoubleStart(t *testing.T) {
	s := New()
	tmp := filepath.Join(t.TempDir(), "snap.db")
	sched := NewSnapshotScheduler(s, tmp, 1*time.Second)

	sched.Start()
	sched.Start() // should be a no-op, not panic

	if !sched.IsRunning() {
		t.Fatal("expected scheduler to still be running")
	}
	sched.Stop()
}

func TestSnapshotSchedulerDoubleStop(t *testing.T) {
	s := New()
	tmp := filepath.Join(t.TempDir(), "snap.db")
	sched := NewSnapshotScheduler(s, tmp, 1*time.Second)

	sched.Start()
	sched.Stop()
	sched.Stop() // should be a no-op, not panic or deadlock

	if sched.IsRunning() {
		t.Fatal("expected scheduler to be stopped")
	}
}

func TestSnapshotSchedulerWritesValidSnapshot(t *testing.T) {
	s := New()
	s.Set("hello", "world", 0)
	s.Set("foo", "bar", 0)

	tmp := filepath.Join(t.TempDir(), "snap.db")
	sched := NewSnapshotScheduler(s, tmp, 40*time.Millisecond)
	sched.Start()
	time.Sleep(100 * time.Millisecond)
	sched.Stop()

	s2 := New()
	if err := LoadSnapshot(s2, tmp); err != nil {
		t.Fatalf("failed to load snapshot: %v", err)
	}

	if v, ok := s2.Get("hello"); !ok || v != "world" {
		t.Errorf("expected hello=world, got %v %v", v, ok)
	}
	if v, ok := s2.Get("foo"); !ok || v != "bar" {
		t.Errorf("expected foo=bar, got %v %v", v, ok)
	}
}
