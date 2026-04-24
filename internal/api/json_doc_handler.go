package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"stashd/internal/store"
)

// NewJSONDocHandler wires up JSON document routes onto r.
func NewJSONDocHandler(r *mux.Router, m *store.JSONDocManager) {
	r.HandleFunc("/json/{key}", jsonSetHandler(m)).Methods(http.MethodPut)
	r.HandleFunc("/json/{key}", jsonGetHandler(m)).Methods(http.MethodGet)
	r.HandleFunc("/json/{key}", jsonDelHandler(m)).Methods(http.MethodDelete)
}

func jsonSetHandler(m *store.JSONDocManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}
		if err := m.JSONSet(key, body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func jsonGetHandler(m *store.JSONDocManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]
		path := r.URL.Query().Get("path")
		val, ok := m.JSONGet(key, path)
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(val)
	}
}

func jsonDelHandler(m *store.JSONDocManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]
		if !m.JSONDel(key) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
