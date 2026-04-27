package api

import (
	"encoding/json"
	"net/http"

	"github.com/user/stashd/internal/store"
)

type txRequest struct {
	Ops []store.TxOp `json:"ops"`
}

type txResponse struct {
	Applied int    `json:"applied"`
	Error   string `json:"error,omitempty"`
}

// NewTransactionHandler returns an http.HandlerFunc that executes a transaction.
func NewTransactionHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req txRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if len(req.Ops) == 0 {
			http.Error(w, "no ops provided", http.StatusBadRequest)
			return
		}

		tx := store.NewTransaction()
		for _, op := range req.Ops {
			tx.Queue(op)
		}

		applied, err := tx.Exec(s)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(txResponse{Applied: applied, Error: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(txResponse{Applied: applied})
	}
}
