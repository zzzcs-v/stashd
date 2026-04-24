package store

import (
	"testing"
	"time"
)

func TestTimeSeriesAddAndLatest(t *testing.T) {
	m := NewTimeSeriesManager()
	m.Add("cpu", 0.5)
	m.Add("cpu", 0.7)
	m.Add("cpu", 0.9)

	pts, err := m.Latest("cpu", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pts) != 2 {
		t.Fatalf("expected 2 points, got %d", len(pts))
	}
	if pts[0].Value != 0.7 || pts[1].Value != 0.9 {
		t.Errorf("unexpected values: %+v", pts)
	}
}

func TestTimeSeriesLatestMoreThanStored(t *testing.T) {
	m := NewTimeSeriesManager()
	m.Add("mem", 1.0)

	pts, err := m.Latest("mem", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pts) != 1 {
		t.Errorf("expected 1 point, got %d", len(pts))
	}
}

func TestTimeSeriesRange(t *testing.T) {
	m := NewTimeSeriesManager()
	before := time.Now().UTC()
	m.Add("temp", 22.0)
	mid := time.Now().UTC()
	m.Add("temp", 23.5)
	after := time.Now().UTC()

	pts, err := m.Range("temp", mid, after)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pts) != 1 {
		t.Errorf("expected 1 point in range, got %d", len(pts))
	}
	_ = before
}

func TestTimeSeriesMissingKey(t *testing.T) {
	m := NewTimeSeriesManager()
	_, err := m.Latest("ghost", 5)
	if err == nil {
		t.Error("expected error for missing key")
	}
	_, err = m.Range("ghost", time.Now(), time.Now())
	if err == nil {
		t.Error("expected error for missing key in range")
	}
}

func TestTimeSeriesDelete(t *testing.T) {
	m := NewTimeSeriesManager()
	m.Add("disk", 100.0)
	m.Delete("disk")
	_, err := m.Latest("disk", 1)
	if err == nil {
		t.Error("expected error after delete")
	}
}
