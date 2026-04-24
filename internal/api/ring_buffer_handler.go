package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/user/stashd/internal/store"
)

type RingBufferHandler struct {
	m *store.RingBufferManager
}

func NewRingBufferHandler(m *store.RingBufferManager) *RingBufferHandler {
	return &RingBufferHandler{m: m}
}

func (h *RingBufferHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		var body struct {
			Value    string `json:"value"`
			Capacity int    `json:"capacity"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Value == "" {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		cap := body.Capacity
		if cap <= 0 {
			cap = 16
		}
		if err := h.m.Push(key, body.Value, cap); err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case http.MethodGet:
		action := r.URL.Query().Get("action")
		if action == "pop" {
			val, err := h.m.Pop(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"value": val})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"len":      strconv.Itoa(h.m.Len(key)),
			"capacity": strconv.Itoa(h.m.Capacity(key)),
		})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
