package store

import "sync"

// Deque is a double-ended queue supporting push/pop from both ends.
type Deque struct {
	mu    sync.Mutex
	items []string
}

type DequeManager struct {
	mu     sync.Mutex
	deques map[string]*Deque
}

func NewDequeManager() *DequeManager {
	return &DequeManager{deques: make(map[string]*Deque)}
}

func (dm *DequeManager) getOrCreate(key string) *Deque {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	if _, ok := dm.deques[key]; !ok {
		dm.deques[key] = &Deque{}
	}
	return dm.deques[key]
}

func (dm *DequeManager) PushFront(key, value string) {
	d := dm.getOrCreate(key)
	d.mu.Lock()
	defer d.mu.Unlock()
	d.items = append([]string{value}, d.items...)
}

func (dm *DequeManager) PushBack(key, value string) {
	d := dm.getOrCreate(key)
	d.mu.Lock()
	defer d.mu.Unlock()
	d.items = append(d.items, value)
}

func (dm *DequeManager) PopFront(key string) (string, bool) {
	d := dm.getOrCreate(key)
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.items) == 0 {
		return "", false
	}
	val := d.items[0]
	d.items = d.items[1:]
	return val, true
}

func (dm *DequeManager) PopBack(key string) (string, bool) {
	d := dm.getOrCreate(key)
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.items) == 0 {
		return "", false
	}
	n := len(d.items)
	val := d.items[n-1]
	d.items = d.items[:n-1]
	return val, true
}

func (dm *DequeManager) Len(key string) int {
	d := dm.getOrCreate(key)
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.items)
}

func (dm *DequeManager) Range(key string) []string {
	d := dm.getOrCreate(key)
	d.mu.Lock()
	defer d.mu.Unlock()
	result := make([]string, len(d.items))
	copy(result, d.items)
	return result
}
