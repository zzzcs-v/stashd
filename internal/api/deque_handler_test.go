package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/stashd/internal/store"
)

func newDequeServer() *httptest.Server {
	dm := store.NewDequeManager()
	h := NewDequeHandler(dm)
	return httptest.NewServer(h)
}

func TestDequePushFrontHTTP(t *testing.T) {
	srv := newDequeServer()
	defer srv.Close()

	body, _ := json.Marshal(map[string]string{"value": "hello"})
	resp, err := http.Post(srv.URL+"/deque/mylist/pushfront", "application/json", bytes.NewReader(body))
	if err != nil || resp.StatusCode != http.StatusNoContent {
		t.Fatalf("pushfront failed: %v %v", err, resp.StatusCode)
	}
}

func TestDequePushBackAndRangeHTTP(t *testing.T) {
	srv := newDequeServer()
	defer srv.Close()

	for _, v := range []string{"a", "b", "c"} {
		body, _ := json.Marshal(map[string]string{"value": v})
		http.Post(srv.URL+"/deque/lst/pushback", "application/json", bytes.NewReader(body))
	}

	resp, _ := http.Get(srv.URL + "/deque/lst/range")
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if result["len"].(float64) != 3 {
		t.Fatalf("expected len 3, got %v", result["len"])
	}
}

func TestDequePopFrontHTTP(t *testing.T) {
	srv := newDequeServer()
	defer srv.Close()

	body, _ := json.Marshal(map[string]string{"value": "first"})
	http.Post(srv.URL+"/deque/q/pushback", "application/json", bytes.NewReader(body))

	resp, _ := http.Get(srv.URL + "/deque/q/popfront")
	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	if result["value"] != "first" {
		t.Fatalf("expected first, got %s", result["value"])
	}
}

func TestDequePopEmptyHTTP(t *testing.T) {
	srv := newDequeServer()
	defer srv.Close()

	resp, _ := http.Get(srv.URL + "/deque/empty/popback")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestDequeMethodNotAllowed(t *testing.T) {
	srv := newDequeServer()
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/deque/q/range", nil)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", resp.StatusCode)
	}
}
