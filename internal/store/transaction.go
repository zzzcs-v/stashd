package store

import (
	"errors"
	"sync"
)

// TxOp represents a single operation in a transaction.
type TxOp struct {
	Op    string // "set", "del"
	Key   string
	Value string
	TTL   int // seconds, 0 = no expiry
}

// Transaction holds a queued list of operations to execute atomically.
type Transaction struct {
	mu  sync.Mutex
	ops []TxOp
}

// NewTransaction creates a new empty transaction.
func NewTransaction() *Transaction {
	return &Transaction{}
}

// Queue adds an operation to the transaction.
func (t *Transaction) Queue(op TxOp) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ops = append(t.ops, op)
}

// Discard clears all queued operations.
func (t *Transaction) Discard() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ops = nil
}

// Exec executes all queued operations against the store atomically.
// Returns the number of operations applied and any error.
func (t *Transaction) Exec(s *Store) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.ops) == 0 {
		return 0, errors.New("transaction is empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	applied := 0
	for _, op := range t.ops {
		switch op.Op {
		case "set":
			s.setLocked(op.Key, op.Value, op.TTL)
			applied++
		case "del":
			delete(s.data, op.Key)
			applied++
		default:
			return applied, errors.New("unknown op: " + op.Op)
		}
	}
	t.ops = nil
	return applied, nil
}
