package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newRateLimitServer(limit int, window time.Duration) *httptest.Server {
	h := NewRateLimitHandler(limit, window)
	mux := http.NewServeMux()
	mux.HandleFunc("/ratelimit/check", h.CheckHandler)
	mux.HandleFunc("/ratelimit/reset", h.ResetHandler)
	return httptest.NewServer(mux)
}

func TestRateLimitHTTPAllow(t *testing.T) {
	ts := newRateLimitServer(5, time.Second)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/ratelimit/check?key=user1")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var body map[string]any
	json.NewDecoder(resp.Body).Decode(&body)
	if body["allowed"] != true {
		t.Fatal("expected allowed true")
	}
}

func TestRateLimitHTTPDeny(t *testing.T) {
	ts := newRateLimitServer(1, time.Second)
	defer ts.Close()
	http.Get(ts.URL + "/ratelimit/check?key=user2")
	resp, _ := http.Get(ts.URL + "/ratelimit/check?key=user2")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", resp.StatusCode)
	}
}

func TestRateLimitHTTPReset(t *testing.T) {
	ts := newRateLimitServer(1, time.Second)
	defer ts.Close()
	http.Get(ts.URL + "/ratelimit/check?key=user3")
	req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/ratelimit/reset?key=user3", nil)
	http.DefaultClient.Do(req)
	resp, _ := http.Get(ts.URL + "/ratelimit/check?key=user3")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 after reset, got %d", resp.StatusCode)
	}
}

func TestRateLimitMethodNotAllowed(t *testing.T) {
	ts := newRateLimitServer(5, time.Second)
	defer ts.Close()
	resp, _ := http.Post(ts.URL+"/ratelimit/check?key=x", "", nil)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", resp.StatusCode)
	}
}
