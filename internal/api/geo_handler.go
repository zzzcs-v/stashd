package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/user/stashd/internal/store"
)

type GeoHandler struct {
	geo *store.GeoManager
}

func NewGeoHandler(geo *store.GeoManager) *GeoHandler {
	return &GeoHandler{geo: geo}
}

func (h *GeoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/geo/add":
		h.handleAdd(w, r)
	case "/geo/dist":
		h.handleDist(w, r)
	case "/geo/nearby":
		h.handleNearby(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *GeoHandler) handleAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Namespace string  `json:"namespace"`
		Member    string  `json:"member"`
		Lat       float64 `json:"lat"`
		Lon       float64 `json:"lon"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	h.geo.Add(req.Namespace, req.Member, req.Lat, req.Lon)
	w.WriteHeader(http.StatusNoContent)
}

func (h *GeoHandler) handleDist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	ns, a, b := q.Get("ns"), q.Get("a"), q.Get("b")
	dist, err := h.geo.Distance(ns, a, b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]float64{"distance_km": dist})
}

func (h *GeoHandler) handleNearby(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	ns := q.Get("ns")
	lat, _ := strconv.ParseFloat(q.Get("lat"), 64)
	lon, _ := strconv.ParseFloat(q.Get("lon"), 64)
	radius, _ := strconv.ParseFloat(q.Get("radius"), 64)
	results := h.geo.Nearby(ns, lat, lon, radius)
	if results == nil {
		results = []store.GeoPoint{}
	}
	json.NewEncoder(w).Encode(results)
}
