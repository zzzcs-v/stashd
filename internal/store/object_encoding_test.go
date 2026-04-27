package store

import (
	"testing"
)

func TestObjectEncodingInt(t *testing.T) {
	om := NewObjectEncodingManager()
	om.Set("counter", "42")

	enc, err := om.Encoding("counter")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enc != EncodingInt {
		t.Errorf("expected %q, got %q", EncodingInt, enc)
	}
}

func TestObjectEncodingFloat(t *testing.T) {
	om := NewObjectEncodingManager()
	om.Set("pi", "3.14")

	enc, err := om.Encoding("pi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enc != EncodingFloat {
		t.Errorf("expected %q, got %q", EncodingFloat, enc)
	}
}

func TestObjectEncodingString(t *testing.T) {
	om := NewObjectEncodingManager()
	om.Set("name", "stashd")

	enc, err := om.Encoding("name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enc != EncodingString {
		t.Errorf("expected %q, got %q", EncodingString, enc)
	}
}

func TestObjectEncodingMissingKey(t *testing.T) {
	om := NewObjectEncodingManager()

	enc, err := om.Encoding("ghost")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
	if enc != EncodingNone {
		t.Errorf("expected %q, got %q", EncodingNone, enc)
	}
}

func TestObjectEncodingDelete(t *testing.T) {
	om := NewObjectEncodingManager()
	om.Set("temp", "123")
	om.Delete("temp")

	_, err := om.Encoding("temp")
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestObjectEncodingGetSet(t *testing.T) {
	om := NewObjectEncodingManager()
	om.Set("key", "hello")

	v, ok := om.Get("key")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if v != "hello" {
		t.Errorf("expected %q, got %q", "hello", v)
	}
}

func TestObjectEncodingNegativeInt(t *testing.T) {
	om := NewObjectEncodingManager()
	om.Set("neg", "-99")

	enc, err := om.Encoding("neg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enc != EncodingInt {
		t.Errorf("expected %q, got %q", EncodingInt, enc)
	}
}
