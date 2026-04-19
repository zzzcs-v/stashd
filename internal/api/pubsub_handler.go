package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/user/stashd/internal/store"
)

type pubSubHandler struct {
	ps *store.PubSub
}

func NewPubSubHandler(ps *store.PubSub) http.Handler {
	h := &pubSubHandler{ps: ps}
	mux := http.NewServeMux()
	mux.HandleFunc("/pubsub/publish", h.handlePublish)
	mux.HandleFunc("/pubsub/subscribe", h.handleSubscribe)
	return mux
}

// POST /pubsub/publish?topic=<topic>  body: plain text message
func (h *pubSubHandler) handlePublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		http.Error(w, "topic required", http.StatusBadRequest)
		return
	}
	var body struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Message == "" {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	h.ps.Publish(topic, body.Message)
	w.WriteHeader(http.StatusNoContent)
}

// GET /pubsub/subscribe?topic=<topic> — SSE stream
func (h *pubSubHandler) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		http.Error(w, "topic required", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := h.ps.Subscribe(topic)
	describe(topic, ch)
lusher, ok := w.(http.Fhttp.Error(w, "streamingupported", http.StatusInternalServerError)
		return
	}
	for {
		select {
		case msg, open := <-ch:
			if !open {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
