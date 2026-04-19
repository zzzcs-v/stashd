package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) incrHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	key := mux.Vars(r)["key"]
	deltaStr := r.URL.Query().Get("delta")
	var delta int64 = 1
	if deltaStr != "" {
		d, err := strconv.ParseInt(deltaStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid delta", http.StatusBadRequest)
			return
		}
		delta = d
	}

	newVal, err := h.store.IncrBy(key, delta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"value": newVal})
}

func (h *Handler) decrHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	key := mux.Vars(r)["key"]
	newVal, err := h.store.Decr(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"value": newVal})
}
