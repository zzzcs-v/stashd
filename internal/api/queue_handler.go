package api

import (
	"encoding/json"
	"net/http"

	"github.com/radovskyb/stashd/internal/store"
)

type QueueHandler struct {
	queue *store.Queue
}

func NewQueueHandler(q *store.Queue) *QueueHandler {
	return &QueueHandler{queue: q}
}

func (h *QueueHandler) PushHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	key := r.URL.Query().Get("key")
	var body struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || key == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	h.queue.Push(key, body.Value)
	w.WriteHeader(http.StatusNoContent)
}

func (h *QueueHandler) PopHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	val, err := h.queue.Pop(key)
	if err == store.ErrQueueEmpty {
		http.Error(w, "queue empty", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"value": val,
		"remaining": h.queue.Len(key),
	})
}
