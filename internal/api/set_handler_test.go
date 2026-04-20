package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/stashd/internal/store"
)

func newSetServer(t *testing.T) (*httptest.Server, *SetHandler) {
	t.Helper()
	s := store.New()
	h := NewSetHandler(s)
	mux := http.NewServeMux()
	mux.HandleFunc("/set/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/add") {
			h.AddHandler(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/remove") {
			h.RemoveHandler(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/members") {
			h.MembersHandler(w, r)
		}
	})
	return httptest.NewServer(mux), h
}

func TestSetAddHTTP(t *testing.T) {
	s := store.New()
	h := NewSetHandler(s)
	body, _ := json.Marshal(map[string][]string{"members": {"a", "b", "c"}})
	req := httptest.NewRequest(http.MethodPost, "/set/fruits/add", bytes.NewReader(body))
	req.URL.Path = "/set/fruits/add"
	w := httptest.NewRecorder()
	h.AddHandler(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]int
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["added"] != 3 {
		t.Fatalf("expected 3 added, got %d", resp["added"])
	}
}

func TestSetMembersHTTP(t *testing.T) {
	s := store.New()
	h := NewSetHandler(s)
	s.SetAdd("colors", "red", "green", "blue")
	req := httptest.NewRequest(http.MethodGet, "/set/colors/members", nil)
	w := httptest.NewRecorder()
	h.MembersHandler(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string][]string
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp["members"]) != 3 {
		t.Fatalf("expected 3 members, got %d", len(resp["members"]))
	}
}

func TestSetMembersMissingHTTP(t *testing.T) {
	s := store.New()
	h := NewSetHandler(s)
	req := httptest.NewRequest(http.MethodGet, "/set/ghost/members", nil)
	w := httptest.NewRecorder()
	h.MembersHandler(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestSetRemoveHTTP(t *testing.T) {
	s := store.New()
	h := NewSetHandler(s)
	s.SetAdd("items", "x", "y", "z")
	body, _ := json.Marshal(map[string][]string{"members": {"x"}})
	req := httptest.NewRequest(http.MethodPost, "/set/items/remove", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.RemoveHandler(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]int
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["removed"] != 1 {
		t.Fatalf("expected 1 removed, got %d", resp["removed"])
	}
}

func TestSetMethodNotAllowed(t *testing.T) {
	s := store.New()
	h := NewSetHandler(s)
	req := httptest.NewRequest(http.MethodDelete, "/set/x/members", nil)
	w := httptest.NewRecorder()
	h.MembersHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}
