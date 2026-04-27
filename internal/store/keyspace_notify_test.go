package store

import (
	"testing"
	"time"
)

func TestKeyspaceNotifySpecificKey(t *testing.T) {
	kn := NewKeyspaceNotifier()
	ch := kn.Subscribe("mykey")

	kn.Notify("mykey", "set")

	select {
	case ev := <-ch:
		if ev.Key != "mykey" || ev.Op != "set" {
			t.Fatalf("unexpected event: %+v", ev)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for event")
	}
}

func TestKeyspaceNotifyWildcard(t *testing.T) {
	kn := NewKeyspaceNotifier()
	ch := kn.Subscribe("*")

	kn.Notify("foo", "del")
	kn.Notify("bar", "expire")

	for i := 0; i < 2; i++ {
		select {
		case <-ch:
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("timed out on event %d", i+1)
		}
	}
}

func TestKeyspaceNotifyNoMatch(t *testing.T) {
	kn := NewKeyspaceNotifier()
	ch := kn.Subscribe("otherkey")

	kn.Notify("mykey", "set")

	select {
	case ev := <-ch:
		t.Fatalf("unexpected event received: %+v", ev)
	case <-time.After(50 * time.Millisecond):
		// expected: no event
	}
}

func TestKeyspaceNotifyUnsubscribe(t *testing.T) {
	kn := NewKeyspaceNotifier()
	ch := kn.Subscribe("mykey")
	kn.Unsubscribe("mykey", ch)

	kn.Notify("mykey", "set")

	select {
	case ev := <-ch:
		t.Fatalf("received event after unsubscribe: %+v", ev)
	case <-time.After(50 * time.Millisecond):
		// expected
	}
}

func TestKeyspaceNotifyMultipleSubscribers(t *testing.T) {
	kn := NewKeyspaceNotifier()
	ch1 := kn.Subscribe("k")
	ch2 := kn.Subscribe("k")

	kn.Notify("k", "expired")

	for _, ch := range []chan KeyspaceEvent{ch1, ch2} {
		select {
		case ev := <-ch:
			if ev.Op != "expired" {
				t.Fatalf("expected expired, got %s", ev.Op)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("timed out")
		}
	}
}
