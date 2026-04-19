package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIncrHTTP(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/keys/hits/incr", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var body map[string]int64
	json.NewDecoder(resp.Body).Decode(&body)
	if body["value"] != 1 {
		t.Errorf("expected value 1, got %d", body["value"])
	}
}

func TestIncrHTTPWithDelta(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/keys/score/incr?delta=5", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var body map[string]int64
	json.NewDecoder(resp.Body).Decode(&body)
	if body["value"] != 5 {
		t.Errorf("expected 5, got %d", body["value"])
	}
}

func TestDecrHTTP(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	http.Post(ts.URL+"/keys/count/incr?delta=3", "", nil)

	resp, err := http.Post(ts.URL+"/keys/count/decr", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var body map[string]int64
	json.NewDecoder(resp.Body).Decode(&body)
	if body["value"] != 2 {
		t.Errorf("expected 2, got %d", body["value"])
	}
}

func TestIncrMethodNotAllowed(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/keys/hits/incr", nil)
	rec := httptest.NewRecorder()
	_ = rec
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}
