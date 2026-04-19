package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// listNamespaceHandler returns all keys in a namespace.
func (h *Handler) listNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	ns := vars["namespace"]
	keys := h.store.ListNamespace(ns)
	if keys == nil {
		keys = []string{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"namespace": ns, "keys": keys})
}

// deleteNamespaceHandler removes all keys in a namespace.
func (h *Handler) deleteNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	ns := vars["namespace"]
	count := h.store.DeleteNamespace(ns)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"namespace": ns, "deleted": count})
}

// RegisterNamespaceRoutes wires namespace endpoints onto a router.
func (h *Handler) RegisterNamespaceRoutes(r *mux.Router) {
	r.HandleFunc("/namespace/{namespace}", h.listNamespaceHandler).Methods(http.MethodGet)
	r.HandleFunc("/namespace/{namespace}", h.deleteNamespaceHandler).Methods(http.MethodDelete)
}
