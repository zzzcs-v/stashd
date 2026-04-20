package store

import (
	"errors"
	"math"
	"sync"
)

const earthRadiusKm = 6371.0

type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Name string  `json:"name"`
}

type GeoManager struct {
	mu   sync.RWMutex
	data map[string]map[string]GeoPoint // namespace -> member -> point
}

func NewGeoManager() *GeoManager {
	return &GeoManager{
		data: make(map[string]map[string]GeoPoint),
	}
}

func (g *GeoManager) Add(namespace, member string, lat, lon float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.data[namespace]; !ok {
		g.data[namespace] = make(map[string]GeoPoint)
	}
	g.data[namespace][member] = GeoPoint{Lat: lat, Lon: lon, Name: member}
}

func (g *GeoManager) Get(namespace, member string) (GeoPoint, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	ns, ok := g.data[namespace]
	if !ok {
		return GeoPoint{}, errors.New("namespace not found")
	}
	pt, ok := ns[member]
	if !ok {
		return GeoPoint{}, errors.New("member not found")
	}
	return pt, nil
}

func (g *GeoManager) Distance(namespace, memberA, memberB string) (float64, error) {
	a, err := g.Get(namespace, memberA)
	if err != nil {
		return 0, err
	}
	b, err := g.Get(namespace, memberB)
	if err != nil {
		return 0, err
	}
	return haversine(a.Lat, a.Lon, b.Lat, b.Lon), nil
}

func (g *GeoManager) Nearby(namespace string, lat, lon, radiusKm float64) []GeoPoint {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var results []GeoPoint
	for _, pt := range g.data[namespace] {
		if haversine(lat, lon, pt.Lat, pt.Lon) <= radiusKm {
			results = append(results, pt)
		}
	}
	return results
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	return earthRadiusKm * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func toRad(deg float64) float64 {
	return deg * math.Pi / 180
}
