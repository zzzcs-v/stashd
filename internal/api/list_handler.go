package api

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	prefix := r.URL.Query().Get("prefix")
	keys := h.store.List(prefix)

	// Return an empty array instead of null when there are no keys
	if keys == nil {
		keys = []string{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string][]string{"keys": keys}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
