package store

import (
	"testing"
)

func newScriptEngine() (*Store, *ScriptEngine) {
	s := New()
	return s, NewScriptEngine(s)
}

func TestScriptSetAndGet(t *testing.T) {
	_, eng := newScriptEngine()
	res := eng.Exec("SET foo bar\nGET foo")
	if res.Error != "" {
		t.Fatalf("unexpected error: %s", res.Error)
	}
	if len(res.Outputs) != 2 || res.Outputs[1] != "bar" {
		t.Fatalf("expected bar, got %v", res.Outputs)
	}
}

func TestScriptGetMissing(t *testing.T) {
	_, eng := newScriptEngine()
	res := eng.Exec("GET missing")
	if res.Error != "" {
		t.Fatalf("unexpected error: %s", res.Error)
	}
	if res.Outputs[0] != "(nil)" {
		t.Fatalf("expected (nil), got %s", res.Outputs[0])
	}
}

func TestScriptDel(t *testing.T) {
	s, eng := newScriptEngine()
	s.Set("x", "1", 0)
	res := eng.Exec("DEL x\nGET x")
	if res.Error != "" {
		t.Fatalf("unexpected error: %s", res.Error)
	}
	if res.Outputs[1] != "(nil)" {
		t.Fatalf("expected (nil) after del, got %s", res.Outputs[1])
	}
}

func TestScriptIncr(t *testing.T) {
	_, eng := newScriptEngine()
	res := eng.Exec("INCR counter\nINCR counter\nINCR counter")
	if res.Error != "" {
		t.Fatalf("unexpected error: %s", res.Error)
	}
	if res.Outputs[2] != "3" {
		t.Fatalf("expected 3, got %s", res.Outputs[2])
	}
}

func TestScriptUnknownCommand(t *testing.T) {
	_, eng := newScriptEngine()
	res := eng.Exec("FLORP key")
	if res.Error == "" {
		t.Fatal("expected error for unknown command")
	}
}

func TestScriptIgnoresComments(t *testing.T) {
	_, eng := newScriptEngine()
	res := eng.Exec("# this is a comment\nSET a b\n# another comment\nGET a")
	if res.Error != "" {
		t.Fatalf("unexpected error: %s", res.Error)
	}
	if len(res.Outputs) != 2 {
		t.Fatalf("expected 2 outputs, got %d", len(res.Outputs))
	}
}

func TestScriptIncrNonInteger(t *testing.T) {
	s, eng := newScriptEngine()
	s.Set("k", "notanumber", 0)
	res := eng.Exec("INCR k")
	if res.Error == "" {
		t.Fatal("expected error for non-integer incr")
	}
}
