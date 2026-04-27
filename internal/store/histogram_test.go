package store

import (
	"testing"
)

func TestHistogramObserveAndQuantile(t *testing.T) {
	h := NewHistogramManager()
	for i := 1; i <= 100; i++ {
		h.Observe("latency", float64(i))
	}
	p50, err := h.Quantile("latency", 0.5)
	if err != nil {
		t.Fatal(err)
	}
	if p50 != 50.0 {
		t.Errorf("expected p50=50, got %v", p50)
	}
	p99, err := h.Quantile("latency", 0.99)
	if err != nil {
		t.Fatal(err)
	}
	if p99 != 99.0 {
		t.Errorf("expected p99=99, got %v", p99)
	}
}

func TestHistogramMissingKey(t *testing.T) {
	h := NewHistogramManager()
	_, err := h.Quantile("missing", 0.5)
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestHistogramSummary(t *testing.T) {
	h := NewHistogramManager()
	h.Observe("req", 10)
	h.Observe("req", 20)
	h.Observe("req", 30)
	s, err := h.Summary("req")
	if err != nil {
		t.Fatal(err)
	}
	if s["count"] != 3 {
		t.Errorf("expected count=3, got %v", s["count"])
	}
	if s["sum"] != 60 {
		t.Errorf("expected sum=60, got %v", s["sum"])
	}
	if s["min"] != 10 {
		t.Errorf("expected min=10, got %v", s["min"])
	}
	if s["max"] != 30 {
		t.Errorf("expected max=30, got %v", s["max"])
	}
	if s["mean"] != 20 {
		t.Errorf("expected mean=20, got %v", s["mean"])
	}
}

func TestHistogramDelete(t *testing.T) {
	h := NewHistogramManager()
	h.Observe("x", 1.0)
	h.Delete("x")
	_, err := h.Summary("x")
	if err == nil {
		t.Error("expected error after delete")
	}
}

func TestHistogramSummaryMissing(t *testing.T) {
	h := NewHistogramManager()
	_, err := h.Summary("nope")
	if err == nil {
		t.Error("expected error for missing histogram")
	}
}
