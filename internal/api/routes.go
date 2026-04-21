package api

import (
	"net/http"

	"github.com/user/stashd/internal/store"
)

func NewRouter(
	s *store.Store,
	snap *store.SnapshotManager,
	wm *store.WatchManager,
	ps *store.PubSub,
	lm *store.LockManager,
	rl *store.RateLimiter,
	q *store.Queue,
	lb *store.LeaderboardManager,
	gm *store.GeoManager,
	hll *store.HyperLogLogManager,
	bf *store.BloomManager,
	dm *store.DequeManager,
) http.Handler {
	mux := http.NewServeMux()

	h := NewHandler(s)
	mux.HandleFunc("/get", h.handleGet)
	mux.HandleFunc("/set", h.handleSet)
	mux.HandleFunc("/delete", h.handleDelete)
	mux.HandleFunc("/list", NewListHandler(s).ServeHTTP)
	mux.HandleFunc("/snapshot", NewSnapshotHandler(snap).ServeHTTP)
	mux.HandleFunc("/stats", NewStatsHandler(s).ServeHTTP)
	mux.HandleFunc("/watch", func(w http.ResponseWriter, r *http.Request) {
		watchHandler(wm, w, r)
	})
	mux.HandleFunc("/namespace/", NewNamespaceHandler(s).ServeHTTP)
	mux.HandleFunc("/ttl/", NewTTLHandler(s).ServeHTTP)
	mux.HandleFunc("/pubsub/", NewPubSubHandler(ps).ServeHTTP)
	mux.HandleFunc("/incr/", NewCounterHandler(s).ServeHTTP)
	mux.HandleFunc("/decr/", NewCounterHandler(s).ServeHTTP)
	mux.HandleFunc("/lock/", NewLockHandler(lm).ServeHTTP)
	mux.HandleFunc("/batch/set", batchSetHandler(s))
	mux.HandleFunc("/batch/get", batchGetHandler(s))
	mux.HandleFunc("/ratelimit/", NewRateLimitHandler(rl).ServeHTTP)
	mux.HandleFunc("/queue/", NewQueueHandler(q).ServeHTTP)
	mux.HandleFunc("/leaderboard/", NewLeaderboardHandler(lb).ServeHTTP)
	mux.HandleFunc("/set/", NewSetHandler(s).ServeHTTP)
	mux.HandleFunc("/hash/", NewHashMapHandler(store.NewHashMap()).ServeHTTP)
	mux.HandleFunc("/bitmap/", NewBitmapHandler(store.NewBitmapStore()).ServeHTTP)
	mux.HandleFunc("/geo/", NewGeoHandler(gm).ServeHTTP)
	mux.HandleFunc("/pipeline", NewPipelineHandler(s).ServeHTTP)
	mux.HandleFunc("/expireat/", NewExpireAtHandler(s).ServeHTTP)
	mux.HandleFunc("/hll/", NewHyperLogLogHandler(hll).ServeHTTP)
	mux.HandleFunc("/bloom/", NewBloomHandler(bf).ServeHTTP)
	mux.HandleFunc("/deque/", NewDequeHandler(dm).ServeHTTP)

	return mux
}
