package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/stashd/internal/store"
)

func newTestServerStats(t *testing.T) *httptest.Server {
	t.Helper()
	s := store.New()
	h := NewHandler(s)
	r := NewRouter(h)
	return httptest.NewServer(r)
}

func TestStatsHTTP(t *testing.T) {
	srv := newTestServerStats(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/stats")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal("failed to decode response:", err)
	}

	if _, ok := result["total_keys"]; !ok {
		t.Error("expected total_keys in response")
	}
	if _, ok := result["uptime_seconds"]; !ok {
		t.Error("expected uptime_seconds in response")
	}
}

func TestStatsHTTPMethodNotAllowed(t *testing.T) {
	srv := newTestServerStats(t)
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/stats", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", resp.StatusCode)
	}
}
