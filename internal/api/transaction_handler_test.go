package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/stashd/internal/store"
)

func newTxServer(t *testing.T) (*httptest.Server, *store.Store) {
	t.Helper()
	s := store.New()
	mux := http.NewServeMux()
	mux.HandleFunc("/tx", NewTransactionHandler(s))
	return httptest.NewServer(mux), s
}

func TestTransactionHTTPSetAndGet(t *testing.T) {
	srv, s := newTxServer(t)
	defer srv.Close()

	body, _ := json.Marshal(map[string]any{
		"ops": []map[string]any{
			{"Op": "set", "Key": "hello", "Value": "world", "TTL": 0},
			{"Op": "set", "Key": "num", "Value": "42", "TTL": 0},
		},
	})
	resp, err := http.Post(srv.URL+"/tx", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var res map[string]any
	json.NewDecoder(resp.Body).Decode(&res)
	if res["applied"].(float64) != 2 {
		t.Errorf("expected applied=2, got %v", res["applied"])
	}

	if v, ok := s.Get("hello"); !ok || v != "world" {
		t.Errorf("expected hello=world, got %v %v", v, ok)
	}
}

func TestTransactionHTTPUnknownOp(t *testing.T) {
	srv, _ := newTxServer(t)
	defer srv.Close()

	body, _ := json.Marshal(map[string]any{
		"ops": []map[string]any{
			{"Op": "badop", "Key": "k"},
		},
	})
	resp, _ := http.Post(srv.URL+"/tx", "application/json", bytes.NewReader(body))
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", resp.StatusCode)
	}
}

func TestTransactionHTTPMethodNotAllowed(t *testing.T) {
	srv, _ := newTxServer(t)
	defer srv.Close()

	resp, _ := http.Get(srv.URL + "/tx")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func TestTransactionHTTPEmptyOps(t *testing.T) {
	srv, _ := newTxServer(t)
	defer srv.Close()

	body := bytes.NewBufferString(`{"ops":[]}`)
	resp, _ := http.Post(srv.URL+"/tx", "application/json", body)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}
