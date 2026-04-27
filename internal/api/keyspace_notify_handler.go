package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/user/stashd/internal/store"
)

// KeyspaceNotifyHandler streams keyspace events over SSE (Server-Sent Events).
// GET /keyspace/notify?key=<key>  — use key=* for all keys.
func KeyspaceNotifyHandler(kn *store.KeyspaceNotifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "missing key parameter", http.StatusBadRequest)
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ch := kn.Subscribe(key)
		defer kn.Unsubscribe(key, ch)

		timeout := time.After(30 * time.Second)
		for {
			select {
			case ev, ok := <-ch:
				if !ok {
					return
				}
				data, _ := json.Marshal(ev)
				_, _ = w.Write([]byte("data: "))
				_, _ = w.Write(data)
				_, _ = w.Write([]byte("\n\n"))
				flusher.Flush()
			case <-timeout:
				return
			case <-r.Context().Done():
				return
			}
		}
	}
}
