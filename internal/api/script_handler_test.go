package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andrebq/stashd/internal/store"
)

func newScriptServer() *httptest.Server {
	s := store.New()
	eng := store.NewScriptEngine(s)
	mux := http.NewServeMux()
	mux.HandleFunc("/script", NewScriptHandler(eng))
	return httptest.NewServer(mux)
}

func postScript(t *testing.T, srv *httptest.Server, script string) (int, map[string]interface{}) {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"script": script})
	resp, err := http.Post(srv.URL+"/script", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return resp.StatusCode, result
}

func TestScriptHTTPSetAndGet(t *testing.T) {
	srv := newScriptServer()
	defer srv.Close()

	code, res := postScript(t, srv, "SET hello world\nGET hello")
	if code != http.StatusOK {
		t.Fatalf("expected 200, got %d", code)
	}
	outputs := res["outputs"].([]interface{})
	if outputs[1].(string) != "world" {
		t.Fatalf("expected world, got %v", outputs[1])
	}
}

func TestScriptHTTPError(t *testing.T) {
	srv := newScriptServer()
	defer srv.Close()

	code, res := postScript(t, srv, "BADCMD key")
	if code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", code)
	}
	if res["error"] == "" {
		t.Fatal("expected error message in response")
	}
}

func TestScriptHTTPMethodNotAllowed(t *testing.T) {
	srv := newScriptServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/script")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", resp.StatusCode)
	}
}

func TestScriptHTTPEmptyScript(t *testing.T) {
	srv := newScriptServer()
	defer srv.Close()

	body, _ := json.Marshal(map[string]string{"script": ""})
	resp, _ := http.Post(srv.URL+"/script", "application/json", bytes.NewReader(body))
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}
