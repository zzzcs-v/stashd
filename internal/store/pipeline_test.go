package store

import (
	"testing"
)

func TestPipelineSetAndGet(t *testing.T) {
	s := New()

	ops := []PipelineOp{
		{Type: "set", Key: "foo", Value: "bar"},
		{Type: "get", Key: "foo"},
	}

	results := s.ExecPipeline(ops)

	if !results[0].OK {
		t.Fatalf("expected set to succeed")
	}
	if !results[1].OK || results[1].Value != "bar" {
		t.Fatalf("expected get to return 'bar', got %q", results[1].Value)
	}
}

func TestPipelineGetMissing(t *testing.T) {
	s := New()

	ops := []PipelineOp{
		{Type: "get", Key: "missing"},
	}

	results := s.ExecPipeline(ops)

	if results[0].OK {
		t.Fatal("expected get to fail for missing key")
	}
	if results[0].Error != "not found" {
		t.Fatalf("expected 'not found' error, got %q", results[0].Error)
	}
}

func TestPipelineDelete(t *testing.T) {
	s := New()
	s.Set("x", "1", 0)

	ops := []PipelineOp{
		{Type: "delete", Key: "x"},
		{Type: "get", Key: "x"},
	}

	results := s.ExecPipeline(ops)

	if !results[0].OK {
		t.Fatal("expected delete to succeed")
	}
	if results[1].OK {
		t.Fatal("expected get after delete to fail")
	}
}

func TestPipelineUnknownOp(t *testing.T) {
	s := New()

	ops := []PipelineOp{
		{Type: "noop", Key: "k"},
	}

	results := s.ExecPipeline(ops)

	if results[0].OK {
		t.Fatal("expected unknown op to fail")
	}
	if results[0].Error == "" {
		t.Fatal("expected an error message for unknown op")
	}
}

func TestPipelineMultipleOps(t *testing.T) {
	s := New()

	ops := []PipelineOp{
		{Type: "set", Key: "a", Value: "1"},
		{Type: "set", Key: "b", Value: "2"},
		{Type: "get", Key: "a"},
		{Type: "get", Key: "b"},
		{Type: "delete", Key: "a"},
		{Type: "get", Key: "a"},
	}

	results := s.ExecPipeline(ops)

	if len(results) != 6 {
		t.Fatalf("expected 6 results, got %d", len(results))
	}
	if results[2].Value != "1" {
		t.Errorf("expected '1', got %q", results[2].Value)
	}
	if results[3].Value != "2" {
		t.Errorf("expected '2', got %q", results[3].Value)
	}
	if results[5].OK {
		t.Error("expected deleted key to be missing")
	}
}
