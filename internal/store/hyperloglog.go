package store

import (
	"math"
	"sync"
)

// HyperLogLog provides approximate cardinality estimation using a simple
// hash-based sketch. Not a full HLL implementation but good enough for
// stashd's purposes.

const hllRegisters = 128

type hllSketch struct {
	regs [hllRegisters]uint8
}

func (h *hllSketch) add(val string) {
	hash := fnv32a(val)
	idx := hash % hllRegisters
	lz := leadingZeros(hash>>7) + 1
	if lz > h.regs[idx] {
		h.regs[idx] = lz
	}
}

func (h *hllSketch) estimate() int64 {
	var sum float64
	for _, v := range h.regs {
		sum += math.Pow(2, -float64(v))
	}
	alpha := 0.7213 / (1.0 + 1.079/float64(hllRegisters))
	est := alpha * float64(hllRegisters) * float64(hllRegisters) / sum
	return int64(est)
}

func fnv32a(s string) uint32 {
	var h uint32 = 2166136261
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

func leadingZeros(x uint32) uint8 {
	if x == 0 {
		return 32
	}
	var n uint8
	for x&0x80000000 == 0 {
		n++
		x <<= 1
	}
	return n
}

type HyperLogLogManager struct {
	mu     sync.Mutex
	sketches map[string]*hllSketch
}

func NewHyperLogLogManager() *HyperLogLogManager {
	return &HyperLogLogManager{
		sketches: make(map[string]*hllSketch),
	}
}

func (m *HyperLogLogManager) Add(key string, values ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.sketches[key]; !ok {
		m.sketches[key] = &hllSketch{}
	}
	for _, v := range values {
		m.sketches[key].add(v)
	}
}

func (m *HyperLogLogManager) Count(key string) (int64, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sketches[key]
	if !ok {
		return 0, false
	}
	return s.estimate(), true
}

func (m *HyperLogLogManager) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sketches, key)
}
