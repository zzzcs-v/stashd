package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/theleeeo/stashd/internal/store"
)

// NewACLHandler returns an http.Handler that manages ACL tokens.
// Routes:
//   POST   /acl/token          — create/update a token  {"token":"x","perm":3}
//   DELETE /acl/token?token=X  — revoke a token
//   GET    /acl/check?token=X&perm=N — check a permission bitmask
func NewACLHandler(acl *store.ACLManager) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/acl/token", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var req struct {
				Token string `json:"token"`
				Perm  uint8  `json:"perm"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			acl.SetToken(req.Token, store.Permission(req.Perm))
			w.WriteHeader(http.StatusNoContent)

		case http.MethodDelete:
			token := r.URL.Query().Get("token")
			if token == "" {
				http.Error(w, "missing token", http.StatusBadRequest)
				return
			}
			acl.RevokeToken(token)
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/acl/check", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		token := r.URL.Query().Get("token")
		var p uint8
		if _, err := fmt.Sscanf(r.URL.Query().Get("perm"), "%d", &p); err != nil {
			http.Error(w, "invalid perm", http.StatusBadRequest)
			return
		}
		err := acl.Check(token, store.Permission(p))
		switch err {
		case store.ErrACLTokenNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case store.ErrACLDenied:
			http.Error(w, err.Error(), http.StatusForbidden)
		case nil:
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return mux
}
