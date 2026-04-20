package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/stashd/internal/store"
)

func newBitmapServer() *httptest.Server {
	bm := &store.Bitmap{}
	// use exported constructor for a fresh instance
	bm = store.NewBitmap()
	// override with a fresh one via reflection workaround — just use NewBitmap and reset
	return httptest.NewServer(NewBitmapHandler(bm))
}

func TestBitmapSetAndGetHTTP(t *testing.T) {
	srv := newBitmapServer()
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/bitmap/set?key=testkey&offset=5", "", nil)
	if err != nil || resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %v err=%v", resp.StatusCode, err)
	}

	resp, err = http.Get(srv.URL + "/bitmap/get?key=testkey&offset=5")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %v err=%v", resp.StatusCode, err)
	}
	var result map[string]bool
	json.NewDecoder(resp.Body).Decode(&result)
	if !result["value"] {
		t.Fatal("expected bit to be true")
	}
}

func TestBitmapCountHTTP(t *testing.T) {
	srv := newBitmapServer()
	defer srv.Close()

	http.Post(srv.URL+"/bitmap/set?key=ckey&offset=1", "", nil)
	http.Post(srv.URL+"/bitmap/set?key=ckey&offset=4", "", nil)

	resp, _ := http.Get(srv.URL + "/bitmap/count?key=ckey")
	var result map[string]int
	json.NewDecoder(resp.Body).Decode(&result)
	if result["count"] != 2 {
		t.Fatalf("expected count 2, got %d", result["count"])
	}
}

func TestBitmapGetUnsetBit(t *testing.T) {
	srv := newBitmapServer()
	defer srv.Close()

	resp, _ := http.Get(srv.URL + "/bitmap/get?key=ghost&offset=10")
	var result map[string]bool
	json.NewDecoder(resp.Body).Decode(&result)
	if result["value"] {
		t.Fatal("expected false for unset bit")
	}
}

func TestBitmapMethodNotAllowed(t *testing.T) {
	srv := newBitmapServer()
	defer srv.Close()

	resp, _ := http.Get(srv.URL + "/bitmap/set?key=k&offset=1")
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", resp.StatusCode)
	}
}

func TestBitmapMissingParams(t *testing.T) {
	srv := newBitmapServer()
	defer srv.Close()

	resp, _ := http.Post(srv.URL+"/bitmap/set?key=k", "", nil)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing offset, got %d", resp.StatusCode)
	}
}
