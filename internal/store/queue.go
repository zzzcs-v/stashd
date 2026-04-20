package store

import (
	"errors"
	"sync"
)

var ErrQueueEmpty = errors.New("queue is empty")

type Queue struct {
	mu    sync.Mutex
	items map[string][]string
}

func NewQueue() *Queue {
	return &Queue{items: make(map[string][]string)}
}

func (q *Queue) Push(key, value string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items[key] = append(q.items[key], value)
}

func (q *Queue) Pop(key string) (string, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	vals := q.items[key]
	if len(vals) == 0 {
		return "", ErrQueueEmpty
	}
	val := vals[0]
	q.items[key] = vals[1:]
	return val, nil
}

func (q *Queue) Len(key string) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items[key])
}

func (q *Queue) Peek(key string) (string, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	vals := q.items[key]
	if len(vals) == 0 {
		return "", ErrQueueEmpty
	}
	return vals[0], nil
}
