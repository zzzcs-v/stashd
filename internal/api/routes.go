package api

import "net/http"

// NewRouter wires up all routes for the stashd HTTP API.
func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/get", h.handleGet)
	mux.HandleFunc("/set", h.handleSet)
	mux.HandleFunc("/delete", h.handleDelete)
	mux.HandleFunc("/snapshot", h.handleSnapshot)
	return mux
}
