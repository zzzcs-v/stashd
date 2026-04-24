package store

import (
	"testing"
)

func TestJSONSetAndGet(t *testing.T) {
	m := NewJSONDocManager()
	err := m.JSONSet("user", []byte(`{"name":"alice","age":30}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val, ok := m.JSONGet("user", ".")
	if !ok {
		t.Fatal("expected document to exist")
	}
	doc, ok := val.(map[string]interface{})
	if !ok || doc["name"] != "alice" {
		t.Fatalf("unexpected doc: %v", val)
	}
}

func TestJSONGetField(t *testing.T) {
	m := NewJSONDocManager()
	_ = m.JSONSet("user", []byte(`{"name":"bob","score":99}`))
	val, ok := m.JSONGet("user", "name")
	if !ok || val != "bob" {
		t.Fatalf("expected 'bob', got %v", val)
	}
}

func TestJSONGetNestedField(t *testing.T) {
	m := NewJSONDocManager()
	_ = m.JSONSet("obj", []byte(`{"meta":{"version":2}}`))
	val, ok := m.JSONGet("obj", "meta.version")
	if !ok {
		t.Fatal("expected nested field")
	}
	if val.(float64) != 2 {
		t.Fatalf("expected 2, got %v", val)
	}
}

func TestJSONGetMissingKey(t *testing.T) {
	m := NewJSONDocManager()
	_, ok := m.JSONGet("nope", ".")
	if ok {
		t.Fatal("expected miss")
	}
}

func TestJSONGetMissingField(t *testing.T) {
	m := NewJSONDocManager()
	_ = m.JSONSet("x", []byte(`{"a":1}`))
	_, ok := m.JSONGet("x", "b")
	if ok {
		t.Fatal("expected miss for missing field")
	}
}

func TestJSONDel(t *testing.T) {
	m := NewJSONDocManager()
	_ = m.JSONSet("doc", []byte(`{"k":"v"}`))
	if !m.JSONDel("doc") {
		t.Fatal("expected true on delete")
	}
	_, ok := m.JSONGet("doc", ".")
	if ok {
		t.Fatal("expected miss after delete")
	}
}

func TestJSONSetInvalidJSON(t *testing.T) {
	m := NewJSONDocManager()
	err := m.JSONSet("bad", []byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
