package store

import (
	"testing"
)

func TestTransactionExecSet(t *testing.T) {
	s := New()
	tx := NewTransaction()
	tx.Queue(TxOp{Op: "set", Key: "foo", Value: "bar", TTL: 0})
	tx.Queue(TxOp{Op: "set", Key: "baz", Value: "qux", TTL: 0})

	n, err := tx.Exec(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Fatalf("expected 2 ops applied, got %d", n)
	}

	if v, ok := s.Get("foo"); !ok || v != "bar" {
		t.Errorf("expected foo=bar, got %v %v", v, ok)
	}
	if v, ok := s.Get("baz"); !ok || v != "qux" {
		t.Errorf("expected baz=qux, got %v %v", v, ok)
	}
}

func TestTransactionExecDel(t *testing.T) {
	s := New()
	s.Set("key1", "val1", 0)

	tx := NewTransaction()
	tx.Queue(TxOp{Op: "del", Key: "key1"})

	n, err := tx.Exec(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 op applied, got %d", n)
	}
	if _, ok := s.Get("key1"); ok {
		t.Error("expected key1 to be deleted")
	}
}

func TestTransactionDiscard(t *testing.T) {
	s := New()
	tx := NewTransaction()
	tx.Queue(TxOp{Op: "set", Key: "x", Value: "1"})
	tx.Discard()

	_, err := tx.Exec(s)
	if err == nil {
		t.Error("expected error on empty transaction after discard")
	}
	if _, ok := s.Get("x"); ok {
		t.Error("key should not exist after discard")
	}
}

func TestTransactionEmptyExec(t *testing.T) {
	s := New()
	tx := NewTransaction()
	_, err := tx.Exec(s)
	if err == nil {
		t.Error("expected error for empty transaction")
	}
}

func TestTransactionUnknownOp(t *testing.T) {
	s := New()
	tx := NewTransaction()
	tx.Queue(TxOp{Op: "incr", Key: "counter"})
	_, err := tx.Exec(s)
	if err == nil {
		t.Error("expected error for unknown op")
	}
}
