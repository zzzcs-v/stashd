package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/user/stashd/internal/store"
)

type SetHandler struct {
	store *store.Store
}

func NewSetHandler(s *store.Store) *SetHandler {
	return &SetHandler{store: s}
}

func (h *SetHandler) AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	key := strings.TrimPrefix(r.URL.Path, "/set/")
	key = strings.TrimSuffix(key, "/add")
	var body struct {
		Members []string `json:"members"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Members) == 0 {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	added := h.store.SetAdd(key, body.Members...)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"added": added})
}

func (h *SetHandler) RemoveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	key := strings.TrimPrefix(r.URL.Path, "/set/")
	key = strings.TrimSuffix(key, "/remove")
	var body struct {
		Members []string `json:"members"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Members) == 0 {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	removed := h.store.SetRemove(key, body.Members...)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"removed": removed})
}

func (h *SetHandler) MembersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	key := strings.TrimPrefix(r.URL.Path, "/set/")
	key = strings.TrimSuffix(key, "/members")
	members, ok := h.store.SetMembers(key)
	if !ok {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"members": members})
}
