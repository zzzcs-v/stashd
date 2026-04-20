package api

import (
	"net/http"

	"github.com/user/stashd/internal/store"
)

// NewRouter wires up all API routes and returns the root mux.
func NewRouter(
	s *store.Store,
	wm *store.WatchManager,
	ps *store.PubSub,
	lm *store.LockManager,
	rl *store.RateLimiter,
	q *store.Queue,
	lb *store.LeaderboardManager,
	bm *store.Bitmap,
) http.Handler {
	mux := http.NewServeMux()

	h := NewHandler(s)
	mux.HandleFunc("/get", h.Get)
	mux.HandleFunc("/set", h.Set)
	mux.HandleFunc("/delete", h.Delete)
	mux.HandleFunc("/ttl", h.TTL)
	mux.HandleFunc("/touch", h.Touch)

	mux.HandleFunc("/snapshot", snapshotHandler(s))
	mux.HandleFunc("/stats", statsHandler(s))
	mux.HandleFunc("/list", listHandler(s))

	mux.HandleFunc("/watch", watchHandler(wm))

	mux.Handle("/ns/", namespaceHandler(s))

	pubsubH := NewPubSubHandler(ps)
	mux.HandleFunc("/publish", pubsubH.Publish)
	mux.HandleFunc("/subscribe", pubsubH.Subscribe)

	counterH := counterHandler(s)
	mux.HandleFunc("/incr", counterH)
	mux.HandleFunc("/decr", counterH)

	lockH := NewLockHandler(lm)
	mux.Handle("/lock/", lockH)

	mux.Handle("/batch/", batchHandler(s))

	rlH := NewRateLimitHandler(rl)
	mux.Handle("/ratelimit/", rlH)

	qH := NewQueueHandler(q)
	mux.Handle("/queue/", qH)

	lbH := NewLeaderboardHandler(lb)
	mux.Handle("/leaderboard/", lbH)

	setH := NewSetHandler(s)
	mux.Handle("/set/", setH)

	hmH := NewHashMapHandler(s)
	mux.Handle("/hash/", hmH)

	bitmapH := NewBitmapHandler(bm)
	mux.Handle("/bitmap/", bitmapH)

	return mux
}
