package api

import (
	"net/http"

	"github.com/radovskyb/stashd/internal/store"
)

// NewRouter builds and returns the main HTTP mux wiring all handlers.
func NewRouter(
	s *store.Store,
	pm *store.PubSub,
	wm *store.WatchManager,
	lm *store.LockManager,
	mlm *store.MultiLockManager,
	rl *store.RateLimiter,
	q *store.Queue,
	lb *store.LeaderboardManager,
	gm *store.GeoManager,
	hll *store.HyperLogLogManager,
	bm *store.BloomManager,
	dm *store.DequeManager,
	sm *store.SortedSetManager,
	tm *store.Trie,
	jm *store.JSONDocManager,
	rbm *store.RingBufferManager,
	tsm *store.TimeSeriesManager,
	se *store.ScriptEngine,
	tr *store.TypeRegistry,
	eat *store.ExpireAtManager,
	hist *store.HistogramManager,
) http.Handler {
	mux := http.NewServeMux()

	NewHandler(mux, s)
	NewPubSubHandler(mux, pm)
	NewLockHandler(mux, lm)
	NewRateLimitHandler(mux, rl)
	NewQueueHandler(mux, q)
	NewLeaderboardHandler(mux, lb)
	NewGeoHandler(mux, gm)
	NewHyperLogLogHandler(mux, hll)
	NewBitmapHandler(mux, s)
	NewDequeHandler(mux, dm)
	NewSetHandler(mux, s)
	NewHashMapHandler(mux, s)
	NewJSONDocHandler(mux, jm)
	NewRingBufferHandler(mux, rbm)
	NewScriptHandler(mux, se)
	NewTransactionHandler(mux, s)
	NewTypeCheckHandler(mux, tr)
	NewExpireAtHandler(mux, eat)
	NewHistogramHandler(mux, hist)

	return mux
}
