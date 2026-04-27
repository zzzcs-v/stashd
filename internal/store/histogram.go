package store

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

// HistogramManager manages named histograms for tracking value distributions.
type HistogramManager struct {
	mu   sync.RWMutex
	hists map[string][]float64
}

func NewHistogramManager() *HistogramManager {
	return &HistogramManager{
		hists: make(map[string][]float64),
	}
}

// Observe adds a value to the named histogram.
func (h *HistogramManager) Observe(key string, value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.hists[key] = append(h.hists[key], value)
}

// Quantile returns the q-th quantile (0.0–1.0) of the named histogram.
func (h *HistogramManager) Quantile(key string, q float64) (float64, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	vals, ok := h.hists[key]
	if !ok || len(vals) == 0 {
		return 0, fmt.Errorf("histogram %q not found or empty", key)
	}
	sorted := make([]float64, len(vals))
	copy(sorted, vals)
	sort.Float64s(sorted)
	idx := int(math.Ceil(q*float64(len(sorted)))) - 1
	if idx < 0 {
		idx = 0
	}
	return sorted[idx], nil
}

// Summary returns count, sum, min, max for a histogram.
func (h *HistogramManager) Summary(key string) (map[string]float64, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	vals, ok := h.hists[key]
	if !ok || len(vals) == 0 {
		return nil, fmt.Errorf("histogram %q not found or empty", key)
	}
	min, max, sum := vals[0], vals[0], 0.0
	for _, v := range vals {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return map[string]float64{
		"count": float64(len(vals)),
		"sum":   sum,
		"min":   min,
		"max":   max,
		"mean":  sum / float64(len(vals)),
	}, nil
}

// Delete removes a histogram.
func (h *HistogramManager) Delete(key string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.hists, key)
}
