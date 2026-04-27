package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/radovskyb/stashd/internal/store"
)

// NewHistogramHandler wires up histogram routes onto mux.
func NewHistogramHandler(mux *http.ServeMux, hm *store.HistogramManager) {
	mux.HandleFunc("/histogram/observe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Key   string  `json:"key"`
			Value float64 `json:"value"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Key == "" {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		hm.Observe(req.Key, req.Value)
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("/histogram/quantile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		key := r.URL.Query().Get("key")
		qs := r.URL.Query().Get("q")
		q, err := strconv.ParseFloat(qs, 64)
		if key == "" || err != nil {
			http.Error(w, "missing key or q", http.StatusBadRequest)
			return
		}
		val, err := hm.Quantile(key, q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]float64{"quantile": val})
	})

	mux.HandleFunc("/histogram/summary", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "missing key", http.StatusBadRequest)
			return
		}
		summary, err := hm.Summary(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(summary)
	})
}
