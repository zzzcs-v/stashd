package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/user/stashd/internal/store"
)

type RateLimitHandler struct {
	rl *store.RateLimiter
}

func NewRateLimitHandler(limit int, window time.Duration) *RateLimitHandler {
	return &RateLimitHandler{rl: store.NewRateLimiter(limit, window)}
}

func (h *RateLimitHandler) CheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}
	allowed := h.rl.Allow(key)
	remaining := h.rl.Remaining(key)
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
	if !allowed {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]any{"allowed": false, "remaining": 0})
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"allowed": true, "remaining": remaining})
}

func (h *RateLimitHandler) ResetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}
	h.rl.Reset(key)
	w.WriteHeader(http.StatusNoContent)
}
