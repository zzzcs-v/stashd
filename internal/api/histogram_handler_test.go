package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/radovskyb/stashd/internal/store"
)

func newHistogramServer() *httptest.Server {
	mux := http.NewServeMux()
	hm := store.NewHistogramManager()
	NewHistogramHandler(mux, hm)
	return httptest.NewServer(mux)
}

func TestHistogramObserveHTTP(t *testing.T) {
	srv := newHistogramServer()
	defer srv.Close()

	body, _ := json.Marshal(map[string]interface{}{"key": "rtt", "value": 42.5})
	resp, err := http.Post(srv.URL+"/histogram/observe", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", resp.StatusCode)
	}
}

func TestHistogramQuantileHTTP(t *testing.T) {
	srv := newHistogramServer()
	defer srv.Close()

	for i := 1; i <= 10; i++ {
		body, _ := json.Marshal(map[string]interface{}{"key": "lat", "value": float64(i * 10)})
		http.Post(srv.URL+"/histogram/observe", "application/json", bytes.NewReader(body))
	}

	resp, err := http.Get(srv.URL + "/histogram/quantile?key=lat&q=0.5")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var result map[string]float64
	json.NewDecoder(resp.Body).Decode(&result)
	if result["quantile"] != 50.0 {
		t.Errorf("expected p50=50, got %v", result["quantile"])
	}
}

func TestHistogramSummaryHTTP(t *testing.T) {
	srv := newHistogramServer()
	defer srv.Close()

	for _, v := range []float64{10, 20, 30} {
		body, _ := json.Marshal(map[string]interface{}{"key": "dur", "value": v})
		http.Post(srv.URL+"/histogram/observe", "application/json", bytes.NewReader(body))
	}

	resp, err := http.Get(srv.URL + "/histogram/summary?key=dur")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var s map[string]float64
	json.NewDecoder(resp.Body).Decode(&s)
	if s["count"] != 3 || s["mean"] != 20 {
		t.Errorf("unexpected summary: %v", s)
	}
}

func TestHistogramMethodNotAllowed(t *testing.T) {
	srv := newHistogramServer()
	defer srv.Close()

	resp, _ := http.Get(srv.URL + "/histogram/observe")
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func TestHistogramQuantileMissingKey(t *testing.T) {
	srv := newHistogramServer()
	defer srv.Close()

	resp, _ := http.Get(srv.URL + "/histogram/quantile?key=ghost&q=0.5")
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}
