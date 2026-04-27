package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/calebdoxsey/stashd/internal/store"
)

// ThrottleHandler handles HTTP requests for the throttle feature.
type ThrottleHandler struct {
	tm *store.ThrottleManager
}

// NewThrottleHandler creates a ThrottleHandler with the given manager.
func NewThrottleHandler(tm *store.ThrottleManager) *ThrottleHandler {
	return &ThrottleHandler{tm: tm}
}

func (h *ThrottleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost:
		limitStr := r.URL.Query().Get("limit")
		windowStr := r.URL.Query().Get("window")
		if limitStr == "" || windowStr == "" {
			http.Error(w, "missing limit or window", http.StatusBadRequest)
			return
		}
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		windowSec, err := strconv.Atoi(windowStr)
		if err != nil || windowSec <= 0 {
			http.Error(w, "invalid window", http.StatusBadRequest)
			return
		}
		allowed, remaining, resetAt := h.tm.Allow(key, limit, time.Duration(windowSec)*time.Second)
		status := http.StatusOK
		if !allowed {
			status = http.StatusTooManyRequests
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]any{
			"allowed":   allowed,
			"remaining": remaining,
			"reset_at":  resetAt.Unix(),
		})

	case http.MethodDelete:
		if err := h.tm.Reset(key); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
