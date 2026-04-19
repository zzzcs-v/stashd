package api

import "net/http"

// NewRouter sets up and returns the HTTP mux with all routes registered.
func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/get", h.handleGet)
	mux.HandleFunc("/set", h.handleSet)
	mux.HandleFunc("/delete", h.handleDelete)
	mux.HandleFunc("/snapshot", h.handleSnapshot)
	mux.HandleFunc("/stats", h.handleStats)
	return mux
}
