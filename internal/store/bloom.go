package store

import (
	"math"
	"sync"
)

// BloomFilter is a probabilistic data structure for membership testing.
type BloomFilter struct {
	bits     []bool
	k        int // number of hash functions
	size     uint
	mu       sync.RWMutex
}

// NewBloomFilter creates a bloom filter sized for n expected items at false-positive rate p.
func NewBloomFilter(n int, p float64) *BloomFilter {
	size := optimalM(n, p)
	k := optimalK(size, n)
	return &BloomFilter{
		bits: make([]bool, size),
		k:    k,
		size: size,
	}
}

func optimalM(n int, p float64) uint {
	return uint(math.Ceil(-float64(n) * math.Log(p) / (math.Log(2) * math.Log(2))))
}

func optimalK(m uint, n int) int {
	return int(math.Round(float64(m) / float64(n) * math.Log(2)))
}

func (bf *BloomFilter) hashes(item string) []uint {
	h := make([]uint, bf.k)
	a := uint(14695981039346656037)
	b := uint(1099511628211)
	for i := 0; i < bf.k; i++ {
		hash := a
		for _, c := range []byte(item) {
			hash ^= uint(c)
			hash *= b
		}
		hash ^= uint(i) * 2654435761
		h[i] = hash % bf.size
	}
	return h
}

// Add inserts an item into the bloom filter.
func (bf *BloomFilter) Add(item string) {
	bf.mu.Lock()
	defer bf.mu.Unlock()
	for _, idx := range bf.hashes(item) {
		bf.bits[idx] = true
	}
}

// MayContain returns true if the item might be in the set (false positives possible).
func (bf *BloomFilter) MayContain(item string) bool {
	bf.mu.RLock()
	defer bf.mu.RUnlock()
	for _, idx := range bf.hashes(item) {
		if !bf.bits[idx] {
			return false
		}
	}
	return true
}

// Reset clears all bits in the filter.
func (bf *BloomFilter) Reset() {
	bf.mu.Lock()
	defer bf.mu.Unlock()
	for i := range bf.bits {
		bf.bits[i] = false
	}
}

// BloomManager manages named bloom filters.
type BloomManager struct {
	mu      sync.RWMutex
	filters map[string]*BloomFilter
}

func NewBloomManager() *BloomManager {
	return &BloomManager{filters: make(map[string]*BloomFilter)}
}

func (bm *BloomManager) GetOrCreate(name string, n int, p float64) *BloomFilter {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	if f, ok := bm.filters[name]; ok {
		return f
	}
	f := NewBloomFilter(n, p)
	bm.filters[name] = f
	return f
}

func (bm *BloomManager) Get(name string) (*BloomFilter, bool) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	f, ok := bm.filters[name]
	return f, ok
}

func (bm *BloomManager) Delete(name string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	delete(bm.filters, name)
}
