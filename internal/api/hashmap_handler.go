package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/user/stashd/internal/store"
)

type HashMapHandler struct {
	hm *store.HashMap
}

func NewHashMapHandler(hm *store.HashMap) *HashMapHandler {
	return &HashMapHandler{hm: hm}
}

// POST /hash/{key}/{field} — set a field
// GET  /hash/{key}/{field} — get a field
// GET  /hash/{key}         — get all fields
// DELETE /hash/{key}/{field} — delete a field
func (h *HashMapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/hash/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}
	key := parts[0]

	if len(parts) == 1 {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		all, err := h.hm.HGetAll(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(all)
		return
	}

	field := parts[1]
	switch r.Method {
	case http.MethodGet:
		v, err := h.hm.HGet(key, field)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"value": v})
	case http.MethodPost:
		var body struct {
			Value string `json:"value"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		h.hm.HSet(key, field, body.Value)
		w.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		if err := h.hm.HDel(key, field); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
