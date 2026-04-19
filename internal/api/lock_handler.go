package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/nicholasgasior/stashd/internal/store"
)

type LockHandler struct {
	lm    *store.LockManager
	held  map[string]func()
	mu    sync.Mutex
}

func NewLockHandler(lm *store.LockManager) *LockHandler {
	return &LockHandler{lm: lm, held: make(map[string]func())}
}

func (h *LockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/lock/")
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		h.mu.Lock()
		if _, exists := h.held[key]; exists {
			h.mu.Unlock()
			http.Error(w, "already locked", http.StatusConflict)
			return
		}
		unlock := h.lm.Lock(key)
		h.held[key] = unlock
		h.mu.Unlock()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "locked", "key": key})

	case http.MethodDelete:
		h.mu.Lock()
		unlock, exists := h.held[key]
		if !exists {
			h.mu.Unlock()
			http.Error(w, "not locked", http.StatusNotFound)
			return
		}
		delete(h.held, key)
		h.mu.Unlock()
		unlock()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "unlocked", "key": key})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
