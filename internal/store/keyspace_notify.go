package store

import "sync"

// KeyspaceEvent represents an event on a key.
type KeyspaceEvent struct {
	Key    string
	Op     string // set, del, expire, expired
}

// KeyspaceNotifier dispatches keyspace notifications to subscribers.
type KeyspaceNotifier struct {
	mu          sync.RWMutex
	subscribers map[string][]chan KeyspaceEvent // key -> channels
	global      []chan KeyspaceEvent            // wildcard subscribers
}

// NewKeyspaceNotifier creates a new KeyspaceNotifier.
func NewKeyspaceNotifier() *KeyspaceNotifier {
	return &KeyspaceNotifier{
		subscribers: make(map[string][]chan KeyspaceEvent),
	}
}

// Subscribe returns a channel that receives events for the given key.
// Pass "*" to receive events for all keys.
func (kn *KeyspaceNotifier) Subscribe(key string) chan KeyspaceEvent {
	ch := make(chan KeyspaceEvent, 16)
	kn.mu.Lock()
	defer kn.mu.Unlock()
	if key == "*" {
		kn.global = append(kn.global, ch)
	} else {
		kn.subscribers[key] = append(kn.subscribers[key], ch)
	}
	return ch
}

// Unsubscribe removes a channel from notifications.
func (kn *KeyspaceNotifier) Unsubscribe(key string, ch chan KeyspaceEvent) {
	kn.mu.Lock()
	defer kn.mu.Unlock()
	if key == "*" {
		kn.global = removeChannel(kn.global, ch)
	} else {
		kn.subscribers[key] = removeChannel(kn.subscribers[key], ch)
	}
}

// Notify dispatches an event to all relevant subscribers.
func (kn *KeyspaceNotifier) Notify(key, op string) {
	event := KeyspaceEvent{Key: key, Op: op}
	kn.mu.RLock()
	defer kn.mu.RUnlock()
	for _, ch := range kn.subscribers[key] {
		select {
		case ch <- event:
		default:
		}
	}
	for _, ch := range kn.global {
		select {
		case ch <- event:
		default:
		}
	}
}

func removeChannel(channels []chan KeyspaceEvent, target chan KeyspaceEvent) []chan KeyspaceEvent {
	out := channels[:0]
	for _, ch := range channels {
		if ch != target {
			out = append(out, ch)
		}
	}
	return out
}
