package api

import (
	"encoding/json"
	"net/http"

	"github.com/andrebq/stashd/internal/store"
)

type scriptRequest struct {
	Script string `json:"script"`
}

type scriptResponse struct {
	Outputs []string `json:"outputs"`
	Error   string   `json:"error,omitempty"`
}

// NewScriptHandler returns an http.HandlerFunc that executes a DSL script atomically.
// POST /script  — body: {"script": "SET k v\nGET k"}
func NewScriptHandler(eng *store.ScriptEngine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req scriptRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		if req.Script == "" {
			http.Error(w, "script is required", http.StatusBadRequest)
			return
		}

		result := eng.Exec(req.Script)

		w.Header().Set("Content-Type", "application/json")
		if result.Error != "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		json.NewEncoder(w).Encode(scriptResponse{
			Outputs: result.Outputs,
			Error:   result.Error,
		})
	}
}
