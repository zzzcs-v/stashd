package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"stashd/internal/store"
)

type LeaderboardHandler struct {
	lm *store.LeaderboardManager
}

func NewLeaderboardHandler(lm *store.LeaderboardManager) *LeaderboardHandler {
	return &LeaderboardHandler{lm: lm}
}

func (h *LeaderboardHandler) scoreHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]
	member := vars["member"]

	var body struct {
		Score float64 `json:"score"`
		Mode  string  `json:"mode"` // "add" or "set"
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.Mode == "set" {
		h.lm.Set(board, member, body.Score)
	} else {
		h.lm.Add(board, member, body.Score)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *LeaderboardHandler) topHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"]
	n := 10
	if nStr := r.URL.Query().Get("n"); nStr != "" {
		if parsed, err := strconv.Atoi(nStr); err == nil {
			n = parsed
		}
	}
	entries := h.lm.Top(board, n)
	w.Header().Set("Content-Type", "application/json")
	jEncode(entries)
}

func (h *LeaderboardHandler) rankHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board := vars["board"["member"]

	rank, score, err := h.lm.Rank(board, member)
	if err != nil {
		http.Error(w, "member not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"rank": rank, "score": score})
}
