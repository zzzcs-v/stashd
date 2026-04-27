package store

import (
	"container/list"
	"sync"
)

type lruEntry struct {
	key   string
	value string
}

// LRUCache is a thread-safe least-recently-used cache with a fixed capacity.
type LRUCache struct {
	mu       sync.Mutex
	cap      int
	ll       *list.List
	items    map[string]*list.Element
}

// NewLRUCache creates a new LRUCache with the given capacity.
func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		capacity = 1
	}
	return &LRUCache{
		cap:   capacity,
		ll:    list.New(),
		items: make(map[string]*list.Element),
	}
}

// Set inserts or updates a key in the cache. Evicts the LRU entry if at capacity.
func (c *LRUCache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.items[key]; ok {
		c.ll.MoveToFront(el)
		el.Value.(*lruEntry).value = value
		return
	}

	if c.ll.Len() >= c.cap {
		c.evict()
	}

	el := c.ll.PushFront(&lruEntry{key: key, value: value})
	c.items[key] = el
}

// Get retrieves a value by key. Returns the value and whether it was found.
func (c *LRUCache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	el, ok := c.items[key]
	if !ok {
		return "", false
	}
	c.ll.MoveToFront(el)
	return el.Value.(*lruEntry).value, true
}

// Delete removes a key from the cache.
func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.items[key]; ok {
		c.removeElement(el)
	}
}

// Len returns the current number of items in the cache.
func (c *LRUCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ll.Len()
}

func (c *LRUCache) evict() {
	el := c.ll.Back()
	if el != nil {
		c.removeElement(el)
	}
}

func (c *LRUCache) removeElement(el *list.Element) {
	c.ll.Remove(el)
	delete(c.items, el.Value.(*lruEntry).key)
}
