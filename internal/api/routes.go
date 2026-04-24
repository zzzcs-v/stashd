package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"stashd/internal/store"
)

// NewRouter builds the full mux.Router with all feature routes registered.
func NewRouter(
	s *store.Store,
	ps *store.PubSub,
	wm *store.WatchManager,
	lm *store.LockManager,
	rl *store.RateLimiter,
	q *store.Queue,
	lb *store.LeaderboardManager,
	gm *store.GeoManager,
	hll *store.HyperLogLogManager,
	jm *store.JSONDocManager,
) http.Handler {
	r := mux.NewRouter()

	NewHandler(r, s)
	NewPubSubHandler(r, ps)
	r.HandleFunc("/watch", watchHandler(wm))
	NewLockHandler(r, lm)
	NewRateLimitHandler(r, rl)
	NewQueueHandler(r, q)
	NewLeaderboardHandler(r, lb)
	NewGeoHandler(r, gm)
	NewHyperLogLogHandler(r, hll)
	NewJSONDocHandler(r, jm)

	return r
}
