package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/user/stashd/internal/store"
)

// NewSlowLogHandler returns an http.Handler for the slow log endpoints.
// GET  /slowlog?count=N  — retrieve up to N entries (0 = all)
// DELETE /slowlog        — reset the slow log
func NewSlowLogHandler(sl *store.SlowLog) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			slowLogGetHandler(w, r, sl)
		case http.MethodDelete:
			slowLogResetHandler(w, sl)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func slowLogGetHandler(w http.ResponseWriter, r *http.Request, sl *store.SlowLog) {
	n := 0
	if raw := r.URL.Query().Get("count"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 0 {
			http.Error(w, "invalid count", http.StatusBadRequest)
			return
		}
		n = parsed
	}
	entries := sl.Get(n)
	type responseEntry struct {
		ID        int64    `json:"id"`
		Timestamp int64    `json:"timestamp_unix"`
		DurationUS int64   `json:"duration_us"`
		Command   string   `json:"command"`
		Args      []string `json:"args"`
	}
	out := make([]responseEntry, len(entries))
	for i, e := range entries {
		out[i] = responseEntry{
			ID:         e.ID,
			Timestamp:  e.Timestamp.Unix(),
			DurationUS: e.Duration.Microseconds(),
			Command:    e.Command,
			Args:       e.Args,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"entries": out, "count": len(out)})
}

func slowLogResetHandler(w http.ResponseWriter, sl *store.SlowLog) {
	sl.Reset()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
