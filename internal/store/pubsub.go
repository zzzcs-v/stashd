package store

import "sync"

// PubSub provides a simple publish/subscribe mechanism for key events.
type PubSub struct {
	mu          sync.RWMutex
	subscribers map[string][]chan string
}

func NewPubSub() *PubSub {
	return &PubSub{
		subscribers: make(map[string][]chan string),
	}
}

// Subscribe returns a channel that receives messages published to the given topic.
func (ps *PubSub) Subscribe(topic string) chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ch := make(chan string, 16)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch
}

// Unsubscribe removes the channel from the given topic.
func (ps *PubSub) Unsubscribe(topic string, ch chan string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	subs := ps.subscribers[topic]
	for i, s := range subs {
		if s == ch {
			ps.subscribers[topic] = append(subs[:i], subs[i+1:]...)
			close(ch)
			return
		}
	}
}

// Publish sends the given message to all subscribers of the given topic.
func (ps *PubSub) Publish(topic, message string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	for _, ch := range ps.subscribers[topic] {
		select {
		case ch <- message:
		default:
		}
	}
}

// Topics returns all active topics.
func (ps *PubSub) Topics() []string {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	topics := make([]string, 0, len(ps.subscribers))
	for t := range ps.subscribers {
		topics = append(topics, t)
	}
	return topics
}

// SubscriberCount returns the number of active subscribers for the given topic.
func (ps *PubSub) SubscriberCount(topic string) int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return len(ps.subscribers[topic])
}
