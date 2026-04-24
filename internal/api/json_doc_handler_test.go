package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"stashd/internal/store"
)

func newJSONDocServer() *httptest.Server {
	r := mux.NewRouter()
	NewJSONDocHandler(r, store.NewJSONDocManager())
	return httptest.NewServer(r)
}

func TestJSONDocSetAndGet(t *testing.T) {
	srv := newJSONDocServer()
	defer srv.Close()

	body := []byte(`{"name":"alice","age":25}`)
	req, _ := http.NewRequest(http.MethodPut, srv.URL+"/json/user", bytes.NewReader(body))
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}

	resp, _ = http.Get(srv.URL + "/json/user")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var doc map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&doc)
	if doc["name"] != "alice" {
		t.Fatalf("unexpected doc: %v", doc)
	}
}

func TestJSONDocGetWithPath(t *testing.T) {
	srv := newJSONDocServer()
	defer srv.Close()

	body := []byte(`{"meta":{"v":7}}`)
	req, _ := http.NewRequest(http.MethodPut, srv.URL+"/json/obj", bytes.NewReader(body))
	http.DefaultClient.Do(req)

	resp, _ := http.Get(srv.URL + "/json/obj?path=meta.v")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var val float64
	json.NewDecoder(resp.Body).Decode(&val)
	if val != 7 {
		t.Fatalf("expected 7, got %v", val)
	}
}

func TestJSONDocGetMissing(t *testing.T) {
	srv := newJSONDocServer()
	defer srv.Close()
	resp, _ := http.Get(srv.URL + "/json/ghost")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestJSONDocDelete(t *testing.T) {
	srv := newJSONDocServer()
	defer srv.Close()

	body := []byte(`{"x":1}`)
	req, _ := http.NewRequest(http.MethodPut, srv.URL+"/json/tmp", bytes.NewReader(body))
	http.DefaultClient.Do(req)

	req, _ = http.NewRequest(http.MethodDelete, srv.URL+"/json/tmp", nil)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}

	resp, _ = http.Get(srv.URL + "/json/tmp")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", resp.StatusCode)
	}
}

func TestJSONDocSetBadJSON(t *testing.T) {
	srv := newJSONDocServer()
	defer srv.Close()

	body := []byte(`not valid json`)
	req, _ := http.NewRequest(http.MethodPut, srv.URL+"/json/bad", bytes.NewReader(body))
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}
