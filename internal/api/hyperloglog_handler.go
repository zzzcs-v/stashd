package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/user/stashd/internal/store"
)

type HyperLogLogHandler struct {
	hll *store.HyperLogLogManager
}

func NewHyperLogLogHandler(hll *store.HyperLogLogManager) *HyperLogLogHandler {
	return &HyperLogLogHandler{hll: hll}
}

func (h *HyperLogLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// /hll/{key}/add  POST
	// /hll/{key}/count GET
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	key := parts[1]
	op := parts[2]

	switch op {
	case "add":
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var body struct {
			Values []string `json:"values"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Values) == 0 {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		h.hll.Add(key, body.Values...)
		w.WriteHeader(http.StatusNoContent)

	case "count":
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		count, ok := h.hll.Count(key)
		if !ok {
			http.Error(w, "key not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int64{"count": count})

	default:
		http.Error(w, "unknown operation", http.StatusBadRequest)
	}
}
