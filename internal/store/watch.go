package store

import "sync"

// WatchEvent represents a change to a key in the store.
type WatchEvent struct {
	Key    string
	Value  string
	Action string // "set" or "delete"
}

// Watcher holds a channel that receives events for subscribed keys.
type Watcher struct {
	Ch     chan WatchEvent
	keys   map[string]struct{}
	mu     sync.RWMutex
}

func newWatcher(keys []string) *Watcher {
	km := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		km[k] = struct{}{}
	}
	return &Watcher{
		Ch:   make(chan WatchEvent, 16),
		keys: km,
	}
}

func (w *Watcher) matches(key string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if len(w.keys) == 0 {
		return true
	}
	_, ok := w.keys[key]
	return ok
}

// WatchManager manages watchers for the store.
type WatchManager struct {
	mu       sync.RWMutex
	watchers []*Watcher
}

func NewWatchManager() *WatchManager {
	return &WatchManager{}
}

// Subscribe registers a new watcher for the given keys (empty = all keys).
func (wm *WatchManager) Subscribe(keys []string) *Watcher {
	w := newWatcher(keys)
	wm.mu.Lock()
	wm.watchers = append(wm.watchers, w)
	wm.mu.Unlock()
	return w
}

// Unsubscribe removes a watcher.
func (wm *WatchManager) Unsubscribe(w *Watcher) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	for i, ww := range wm.watchers {
		if ww == w {
			wm.watchers = append(wm.watchers[:i], wm.watchers[i+1:]...)
			close(w.Ch)
			return
		}
	}
}

// Notify sends an event to all matching watchers.
func (wm *WatchManager) Notify(event WatchEvent) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	for _, w := range wm.watchers {
		if w.matches(event.Key) {
			select {
			case w.Ch <- event:
			default:
			}
		}
	}
}
