package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/user/stashd/internal/store"
)

// NewExpireAtHandler returns an http.Handler for EXPIREAT and EXPIRETIME operations.
//
// POST /expireat/{key}?ts=<unix_sec>          — set absolute expiry (seconds)
// POST /pexpireat/{key}?ts=<unix_ms>          — set absolute expiry (milliseconds)
// GET  /expiretime/{key}                      — get absolute expiry as unix seconds
func NewExpireAtHandler(s *store.Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/expireat/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		key := r.URL.Path[len("/expireat/"):]
		tsStr := r.URL.Query().Get("ts")
		ts, err := strconv.ParseInt(tsStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid ts parameter", http.StatusBadRequest)
			return
		}
		if err := s.ExpireAt(key, ts); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/pexpireat/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		key := r.URL.Path[len("/pexpireat/"):]
		tsStr := r.URL.Query().Get("ts")
		ts, err := strconv.ParseInt(tsStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid ts parameter", http.StatusBadRequest)
			return
		}
		if err := s.PExpireAt(key, ts); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/expiretime/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		key := r.URL.Path[len("/expiretime/"):]
		ts := s.ExpireTime(key)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int64{"expiretime": ts})
	})

	return mux
}
