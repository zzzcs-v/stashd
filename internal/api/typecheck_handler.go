package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/user/stashd/internal/store"
)

// NewTypeCheckHandler returns an http.Handler that exposes the type
// registry over HTTP.
//
//   GET  /type/{key}          → returns the type of a key
//   POST /type/{key}?type=X   → registers / updates the type for a key
//   DELETE /type/{key}        → removes the key from the registry
func NewTypeCheckHandler(reg *store.TypeRegistry) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/type/", func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimPrefix(r.URL.Path, "/type/")
		if key == "" {
			http.Error(w, "missing key", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			vt, err := reg.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"key": key, "type": string(vt)})

		case http.MethodPost:
			typeName := r.URL.Query().Get("type")
			if typeName == "" {
				http.Error(w, "missing ?type= query param", http.StatusBadRequest)
				return
			}
			reg.Set(key, store.ValueType(typeName), 0)
			w.WriteHeader(http.StatusNoContent)

		case http.MethodDelete:
			reg.Delete(key)
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/types", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		keys := reg.Keys()
		out := make(map[string]string, len(keys))
		for k, vt := range keys {
			out[k] = string(vt)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
	})

	return mux
}
