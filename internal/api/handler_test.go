package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/stashd/internal/api"
	"github.com/user/stashd/internal/store"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	s := store.New()
	router := api.NewRouter(s)
	return httptest.NewServer(router)
}

func TestSetAndGetHTTP(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	body, _ := json.Marshal(map[string]any{"value": "hello"})
	resp, err := http.NewRequest(http.MethodPut, srv.URL+"/keys/foo", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{}
	res, _ := client.Do(resp)
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", res.StatusCode)
	}

	getRes, _ := http.Get(srv.URL + "/keys/foo")
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getRes.StatusCode)
	}
	var out map[string]string
	json.NewDecoder(getRes.Body).Decode(&out)
	if out["value"] != "hello" {
		t.Errorf("expected 'hello', got %q", out["value"])
	}
}

func TestGetMissingHTTP(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	res, _ := http.Get(srv.URL + "/keys/missing")
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", res.StatusCode)
	}
}

func TestDeleteHTTP(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	body, _ := json.Marshal(map[string]any{"value": "bye"})
	req, _ := http.NewRequest(http.MethodPut, srv.URL+"/keys/bar", bytes.NewReader(body))
	client := &http.Client{}
	client.Do(req)

	del, _ := http.NewRequest(http.MethodDelete, srv.URL+"/keys/bar", nil)
	client.Do(del)

	res, _ := http.Get(srv.URL + "/keys/bar")
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", res.StatusCode)
	}
}
