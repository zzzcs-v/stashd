package api

import (
	"net/http"
	"time"

	"github.com/user/stashd/internal/store"
)

func NewRouter(s *store.Store, ps *store.PubSub, wm *store.WatchManager, lm *store.LockManager) http.Handler {
	h := NewHandler(s)
	snap := &snapshotHandler{store: s}
	stats := &statsHandler{store: s}
	ns := &namespaceHandler{store: s}
	list := &listHandler{store: s}
	counter := &counterHandler{store: s}
	batch := struct {
		set http.HandlerFunc
		get http.HandlerFunc
		del http.HandlerFunc
	}{batchSetHandler(s), batchGetHandler(s), batchDeleteHandler(s)}
	lockH := NewLockHandler(lm)
	pubsubH := NewPubSubHandler(ps)
	rlH := NewRateLimitHandler(100, time.Minute)

	mux := http.NewServeMux()
	mux.HandleFunc("/get", h.getHandler)
	mux.HandleFunc("/set", h.setHandler)
	mux.HandleFunc("/delete", h.deleteHandler)
	mux.HandleFunc("/snapshot", snap.handleSnapshot)
	mux.HandleFunc("/stats", stats.handleStats)
	mux.HandleFunc("/namespace/list", ns.listNamespaceHandler)
	mux.HandleFunc("/namespace/delete", ns.deleteNamespaceHandler)
	mux.HandleFunc("/list", list.handleList)
	mux.HandleFunc("/incr", counter.incrHandler)
	mux.HandleFunc("/decr", counter.decrHandler)
	mux.HandleFunc("/batch/set", batch.set)
	mux.HandleFunc("/batch/get", batch.get)
	mux.HandleFunc("/batch/delete", batch.del)
	mux.HandleFunc("/lock", lockH.HandleLock)
	mux.HandleFunc("/publish", pubsubH.PublishHandler)
	mux.HandleFunc("/subscribe", pubsubH.SubscribeHandler)
	mux.HandleFunc("/watch", watchHandler(wm))
	mux.HandleFunc("/ratelimit/check", rlH.CheckHandler)
	mux.HandleFunc("/ratelimit/reset", rlH.ResetHandler)
	return mux
}
