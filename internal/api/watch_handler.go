package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/user/stashd/internal/store"
)

// watchHandler streams key events to the client via Server-Sent Events.
func watchHandler(wm *store.WatchManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}

		var keys []string
		if q := r.URL.Query().Get("keys"); q != "" {
			keys = strings.Split(q, ",")
		}

		watcher := wm.Subscribe(keys)
		defer wm.Unsubscribe(watcher)

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ctx := r.Context()
		for {
			select {
			case <-ctx.Done():
				return
			case ev, ok := <-watcher.Ch:
				if !ok {
					return
				}
				data, err := json.Marshal(ev)
				if err != nil {
					continue
				}
				_, _ = w.Write([]byte("data: "))
				_, _ = w.Write(data)
				_, _ = w.Write([]byte("\n\n"))
				flusher.Flush()
			}
		}
	}
}
