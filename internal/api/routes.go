package api

import (
	"net/http"

	"github.com/user/stashd/internal/store"
)

func NewRouter(s *store.Store) http.Handler {
	h := NewHandler(s)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /keys/{key}", h.Get)
	mux.HandleFunc("PUT /keys/{key}", h.Set)
	mux.HandleFunc("DELETE /keys/{key}", h.Delete)
	mux.HandleFunc("GET /keys/{key}/ttl", h.TTL)

	return mux
}
