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

// postScore is a test helper that submits a score for a member on a given board.
func postScore(t *testing.T, srvURL, board, member string, score float64, mode string) {
	t.Helper()
	body, _ := json.Marshal(map[string]interface{}{"score": score, "mode": mode})
	resp, err := http.Post(srvURL+"/leaderboard/"+board+"/score/"+member, "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected 204, got %d", resp.StatusCode)
	}
}

func TestLeaderboardScoreAndTop(t *testing.T) {
	srv := newLeaderboardServer()
	defer srv.Close()

	postScore(t, srv.URL, "game", "alice", 100, "add")
	postScore(t, srv.URL, "game", "bob", 200, "set")
	postScore(t, srv.URL, "game", "alice", 50, "add")

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

	postScore(t, srv.URL, "lb", "carol", 300.0, "set")

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
