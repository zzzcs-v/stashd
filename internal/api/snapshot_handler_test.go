package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/user/stashd/internal/store"
)

func TestSnapshotHTTP(t *testing.T) {
	s := store.New()
	s.Set("snap-key", "snap-val", 0)

	server := newTestServer(s)
	defer server.Close()

	tmp, err := os.CreateTemp("", "stashd-http-snapshot-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	resp, err := http.Post(server.URL+"/snapshot?path="+tmp.Name(), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	// verify snapshot file was created and is non-empty
	info, err := os.Stat(tmp.Name())
	if err != nil {
		t.Fatalf("snapshot file not found: %v", err)
	}
	if info.Size() == 0 {
		t.Error("snapshot file is empty")
	}
}

func TestSnapshotHTTPMethodNotAllowed(t *testing.T) {
	s := store.New()
	server := newTestServer(s)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/snapshot", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func newTestServerSnapshot(s *store.Store) *httptest.Server {
	h := NewHandler(s)
	return httptest.NewServer(NewRouter(h))
}
