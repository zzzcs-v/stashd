package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bjarnemagnussen/stashd/internal/store"
)

type batchSetRequest struct {
	Items []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		TTL   int    `json:"ttl"` // seconds
	} `json:"items"`
}

type batchGetRequest struct {
	Keys []string `json:"keys"`
}

func batchSetHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req batchSetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		items := make([]store.BatchSetItem, len(req.Items))
		for i, it := range req.Items {
			items[i] = store.BatchSetItem{
				Key:   it.Key,
				Value: it.Value,
				TTL:   time.Duration(it.TTL) * time.Second,
			}
		}
		s.BatchSet(items)
		w.WriteHeader(http.StatusNoContent)
	}
}

func batchGetHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req batchGetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		results := s.BatchGet(req.Keys)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}
