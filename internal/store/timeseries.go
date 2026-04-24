package store

import (
	"errors"
	"sync"
	"time"
)

// TSPoint represents a single time-series data point.
type TSPoint struct {
	Timestamp time.Time
	Value     float64
}

type tsSeries struct {
	points []TSPoint
}

// TimeSeriesManager manages named time-series collections.
type TimeSeriesManager struct {
	mu     sync.RWMutex
	series map[string]*tsSeries
}

// NewTimeSeriesManager creates a new TimeSeriesManager.
func NewTimeSeriesManager() *TimeSeriesManager {
	return &TimeSeriesManager{
		series: make(map[string]*tsSeries),
	}
}

// Add appends a value with the current timestamp to the named series.
func (m *TimeSeriesManager) Add(key string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.series[key]; !ok {
		m.series[key] = &tsSeries{}
	}
	m.series[key].points = append(m.series[key].points, TSPoint{
		Timestamp: time.Now().UTC(),
		Value:     value,
	})
}

// Range returns all points in [from, to] for the named series.
func (m *TimeSeriesManager) Range(key string, from, to time.Time) ([]TSPoint, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.series[key]
	if !ok {
		return nil, errors.New("key not found")
	}
	var result []TSPoint
	for _, p := range s.points {
		if !p.Timestamp.Before(from) && !p.Timestamp.After(to) {
			result = append(result, p)
		}
	}
	return result, nil
}

// Latest returns the most recent n points for the named series.
func (m *TimeSeriesManager) Latest(key string, n int) ([]TSPoint, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.series[key]
	if !ok {
		return nil, errors.New("key not found")
	}
	pts := s.points
	if n >= len(pts) {
		return append([]TSPoint{}, pts...), nil
	}
	return append([]TSPoint{}, pts[len(pts)-n:]...), nil
}

// Delete removes a series entirely.
func (m *TimeSeriesManager) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.series, key)
}
