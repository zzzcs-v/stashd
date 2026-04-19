package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/user/stashd/internal/store"
)

type Handler struct {
	store *store.Store
}

func NewHandler(s *store.Store) *Handler {
	return &Handler{store: s}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	val, ok := h.store.Get(key)
	if !ok {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"key": key, "value": val})
}

func (h *Handler) Set(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	var body struct {
		Value string `json:"value"`
		TTL   int    `json:"ttl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	var ttl time.Duration
	if body.TTL > 0 {
		ttl = time.Duration(body.TTL) * time.Second
	}
	h.store.Set(key, body.Value, ttl)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	h.store.Delete(key)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) TTL(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	ttl, ok := h.store.TTL(key)
	if !ok {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"key": key, "ttl": strconv.Itoa(int(ttl.Seconds()))})
}
