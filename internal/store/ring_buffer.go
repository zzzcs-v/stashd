package store

import (
	"errors"
	"sync"
)

var ErrRingBufferFull = errors.New("ring buffer is full")
var ErrRingBufferEmpty = errors.New("ring buffer is empty")

type RingBuffer struct {
	mu       sync.Mutex
	items    []string
	head     int
	tail     int
	size     int
	capacity int
}

type RingBufferManager struct {
	mu      sync.Mutex
	buffers map[string]*RingBuffer
}

func NewRingBufferManager() *RingBufferManager {
	return &RingBufferManager{
		buffers: make(map[string]*RingBuffer),
	}
}

func (m *RingBufferManager) getOrCreate(key string, capacity int) *RingBuffer {
	m.mu.Lock()
	defer m.mu.Unlock()
	if rb, ok := m.buffers[key]; ok {
		return rb
	}
	rb := &RingBuffer{
		items:    make([]string, capacity),
		capacity: capacity,
	}
	m.buffers[key] = rb
	return rb
}

func (m *RingBufferManager) Push(key, value string, capacity int) error {
	rb := m.getOrCreate(key, capacity)
	rb.mu.Lock()
	defer rb.mu.Unlock()
	if rb.size == rb.capacity {
		return ErrRingBufferFull
	}
	rb.items[rb.tail] = value
	rb.tail = (rb.tail + 1) % rb.capacity
	rb.size++
	return nil
}

func (m *RingBufferManager) Pop(key string) (string, error) {
	m.mu.Lock()
	rb, ok := m.buffers[key]
	m.mu.Unlock()
	if !ok || rb.size == 0 {
		return "", ErrRingBufferEmpty
	}
	rb.mu.Lock()
	defer rb.mu.Unlock()
	val := rb.items[rb.head]
	rb.head = (rb.head + 1) % rb.capacity
	rb.size--
	return val, nil
}

func (m *RingBufferManager) Len(key string) int {
	m.mu.Lock()
	rb, ok := m.buffers[key]
	m.mu.Unlock()
	if !ok {
		return 0
	}
	rb.mu.Lock()
	defer rb.mu.Unlock()
	return rb.size
}

func (m *RingBufferManager) Capacity(key string) int {
	m.mu.Lock()
	rb, ok := m.buffers[key]
	m.mu.Unlock()
	if !ok {
		return 0
	}
	return rb.capacity
}
