package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/user/stashd/internal/store"
)

type DequeHandler struct {
	dm *store.DequeManager
}

func NewDequeHandler(dm *store.DequeManager) *DequeHandler {
	return &DequeHandler{dm: dm}
}

func (h *DequeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// /deque/{key}/{action}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.Error(w, "usage: /deque/{key}/{action}", http.StatusBadRequest)
		return
	}
	key := parts[1]
	action := parts[2]

	switch r.Method {
	case http.MethodPost:
		var body struct {
			Value string `json:"value"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Value == "" {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		switch action {
		case "pushfront":
			h.dm.PushFront(key, body.Value)
		case "pushback":
			h.dm.PushBack(key, body.Value)
		default:
			http.Error(w, "unknown action", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case http.MethodGet:
		switch action {
		case "popfront":
			val, ok := h.dm.PopFront(key)
			if !ok {
				http.Error(w, "empty", http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"value": val})
		case "popback":
			val, ok := h.dm.PopBack(key)
			if !ok {
				http.Error(w, "empty", http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(map[string]string{"value": val})
		case "range":
			items := h.dm.Range(key)
			json.NewEncoder(w).Encode(map[string]interface{}{"items": items, "len": len(items)})
		default:
			http.Error(w, "unknown action", http.StatusBadRequest)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
