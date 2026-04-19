package api

import (
	"net/http"
)

// handleSnapshot handles POST /snapshot — triggers a store snapshot save.
func (h *Handler) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		path = "stashd.snapshot.json"
	}

	if err := h.store.SaveSnapshot(path); err != nil {
		http.Error(w, "failed to save snapshot: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("snapshot saved to " + path + "\n"))
}
