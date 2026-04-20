package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"stashd/internal/store"
)

func newLeaderboardServer() *httptest.Server {
	lm := store.NewLeaderboardManager()
	h := NewLeaderboardHandler(lm)
	r := mux.NewRouter()
	r.HandleFunc("/leaderboard/{board}/score/{member}", h.scoreHandler).Methods(http.MethodPost)
	r.HandleFunc("/leaderboard/{board}/top", h.topHandler).Methods(http.MethodGet)
	r.HandleFunc("/leaderboard/{board}/rank/{member}", h.rankHandler).Methods(http.MethodGet)
	return httptest.NewServer(r)
}

func TestLeaderboardScoreAndTop(t *testing.T) {
	srv := newLeaderboardServer()
	defer srv.Close()

	postScore := func(board, member string, score float64, mode string) {
		body, _ := json.Marshal(map[string]interface{}{"score": score, "mode": mode})
		resp, err := http.Post(srv.URL+"/leaderboard/"+board+"/score/"+member, "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("expected 204, got %d", resp.StatusCode)
		}
	}

	postScore("game", "alice", 100, "add")
	postScore("game", "bob", 200, "set")
	postScore("game", "alice", 50, "add")

	resp, err := http.Get(srv.URL + "/leaderboard/game/top?n=2")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var entries []store.LeaderboardEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Member != "bob" {
		t.Errorf("expected bob at top, got %s", entries[0].Member)
	}
}

func TestLeaderboardRankHTTP(t *testing.T) {
	srv := newLeaderboardServer()
	defer srv.Close()

	body, _ := json.Marshal(map[string]interface{}{"score": 300.0, "mode": "set"})
	http.Post(srv.URL+"/leaderboard/lb/score/carol", "application/json", bytes.NewReader(body))

	resp, err := http.Get(srv.URL + "/leaderboard/lb/rank/carol")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if result["rank"].(float64) != 1 {
		t.Errorf("expected rank 1, got %v", result["rank"])
	}
}

func TestLeaderboardRankMissingHTTP(t *testing.T) {
	srv := newLeaderboardServer()
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/leaderboard/lb/rank/nobody")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}
