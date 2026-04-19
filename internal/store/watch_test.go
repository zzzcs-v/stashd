package store

import (
	"testing"
	"time"
)

func TestWatchAllKeys(t *testing.T) {
	wm := NewWatchManager()
	w := wm.Subscribe(nil)
	defer wm.Unsubscribe(w)

	wm.Notify(WatchEvent{Key: "foo", Value: "bar", Action: "set"})

	select {
	case ev := <-w.Ch:
		if ev.Key != "foo" || ev.Value != "bar" || ev.Action != "set" {
			t.Errorf("unexpected event: %+v", ev)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected event but got none")
	}
}

func TestWatchSpecificKey(t *testing.T) {
	wm := NewWatchManager()
	w := wm.Subscribe([]string{"mykey"})
	defer wm.Unsubscribe(w)

	wm.Notify(WatchEvent{Key: "other", Value: "v", Action: "set"})
	wm.Notify(WatchEvent{Key: "mykey", Value: "hello", Action: "set"})

	select {
	case ev := <-w.Ch:
		if ev.Key != "mykey" {
			t.Errorf("expected mykey, got %s", ev.Key)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected event but got none")
	}
}

func TestWatchUnsubscribe(t *testing.T) {
	wm := NewWatchManager()
	w := wm.Subscribe(nil)
	wm.Unsubscribe(w)

	// channel should be closed
	select {
	case _, ok := <-w.Ch:
		if ok {
			t.Fatal("expected closed channel")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected channel to be closed")
	}
}

func TestWatchDeleteEvent(t *testing.T) {
	wm := NewWatchManager()
	w := wm.Subscribe(nil)
	defer wm.Unsubscribe(w)

	wm.Notify(WatchEvent{Key: "gone", Action: "delete"})

	select {
	case ev := <-w.Ch:
		if ev.Action != "delete" || ev.Key != "gone" {
			t.Errorf("unexpected event: %+v", ev)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected delete event")
	}
}
