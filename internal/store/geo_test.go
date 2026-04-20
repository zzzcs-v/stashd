package store

import (
	"testing"
)

func TestGeoAddAndGet(t *testing.T) {
	gm := NewGeoManager()
	gm.Add("cities", "london", 51.5074, -0.1278)
	pt, err := gm.Get("cities", "london")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pt.Name != "london" {
		t.Errorf("expected london, got %s", pt.Name)
	}
}

func TestGeoGetMissing(t *testing.T) {
	gm := NewGeoManager()
	_, err := gm.Get("cities", "nowhere")
	if err == nil {
		t.Error("expected error for missing member")
	}
}

func TestGeoDistance(t *testing.T) {
	gm := NewGeoManager()
	gm.Add("cities", "london", 51.5074, -0.1278)
	gm.Add("cities", "paris", 48.8566, 2.3522)
	dist, err := gm.Distance("cities", "london", "paris")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// roughly 340 km
	if dist < 300 || dist > 400 {
		t.Errorf("unexpected distance: %.2f km", dist)
	}
}

func TestGeoDistanceMissingMember(t *testing.T) {
	gm := NewGeoManager()
	gm.Add("cities", "london", 51.5074, -0.1278)
	_, err := gm.Distance("cities", "london", "ghost")
	if err == nil {
		t.Error("expected error for missing member")
	}
}

func TestGeoNearby(t *testing.T) {
	gm := NewGeoManager()
	gm.Add("places", "a", 40.7128, -74.0060) // NYC
	gm.Add("places", "b", 40.6501, -73.9496) // Brooklyn
	gm.Add("places", "c", 34.0522, -118.2437) // LA

	results := gm.Nearby("places", 40.7128, -74.0060, 20.0)
	if len(results) != 2 {
		t.Errorf("expected 2 nearby results, got %d", len(results))
	}
}

func TestGeoNearbyEmptyNamespace(t *testing.T) {
	gm := NewGeoManager()
	results := gm.Nearby("empty", 0, 0, 100)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
