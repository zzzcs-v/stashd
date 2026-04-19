package store

import (
	"testing"
	"time"
)

func TestPubSubSubscribeAndPublish(t *testing.T) {
	ps := NewPubSub()
	ch := ps.Subscribe("alerts")
	ps.Publish("alerts", "hello")

	select {
	case msg := <-ch:
		if msg != "hello" {
			t.Fatalf("expected 'hello', got %q", msg)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for message")
	}
}

func TestPubSubNoSubscribers(t *testing.T) {
	ps := NewPubSub()
	// should not panic
	ps.Publish("ghost", "nobody home")
}

func TestPubSubUnsubscribe(t *testing.T) {
	ps := NewPubSub()
	ch := ps.Subscribe("events")
	ps.Unsubscribe("events", ch)

	// channel should be closed
	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("expected channel to be closed")
		}
	default:
		t.Fatal("expected closed channel to be readable")
	}
}

func TestPubSubMultipleSubscribers(t *testing.T) {
	ps := NewPubSub()
	ch1 := ps.Subscribe("news")
	ch2 := ps.Subscribe("news")
	ps.Publish("news", "update")

	for _, ch := range []chan string{ch1, ch2} {
		select {
		case msg := <-ch:
			if msg != "update" {
				t.Fatalf("expected 'update', got %q", msg)
			}
		case <-time.After(time.Second):
			t.Fatal("timed out")
		}
	}
}

func TestPubSubTopics(t *testing.T) {
	ps := NewPubSub()
	ps.Subscribe("a")
	ps.Subscribe("b")
	topics := ps.Topics()
	if len(topics) != 2 {
		t.Fatalf("expected 2 topics, got %d", len(topics))
	}
}
