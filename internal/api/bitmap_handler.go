package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/user/stashd/internal/store"
)

// NewBitmapHandler returns an http.Handler for bitmap operations.
func NewBitmapHandler(bm *store.Bitmap) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/bitmap/set", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		key := r.URL.Query().Get("key")
		offsetStr := r.URL.Query().Get("offset")
		offset, err := strconv.Atoi(offsetStr)
		if key == "" || err != nil {
			http.Error(w, "key and valid offset required", http.StatusBadRequest)
			return
		}
		if err := bm.BitSet(key, offset); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("/bitmap/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		key := r.URL.Query().Get("key")
		offsetStr := r.URL.Query().Get("offset")
		offset, err := strconv.Atoi(offsetStr)
		if key == "" || err != nil {
			http.Error(w, "key and valid offset required", http.StatusBadRequest)
			return
		}
		val, err := bm.BitGet(key, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]bool{"value": val})
	})

	mux.HandleFunc("/bitmap/count", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "key required", http.StatusBadRequest)
			return
		}
		count := bm.BitCount(key)
		json.NewEncoder(w).Encode(map[string]int{"count": count})
	})

	return mux
}
