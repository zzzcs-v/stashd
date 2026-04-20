package store

import (
	"fmt"
	"sync"
)

// Bitmap stores a map of string key -> bit array (as []bool)
type Bitmap struct {
	mu   sync.RWMutex
	bits map[string][]bool
}

var globalBitmap = &Bitmap{
	bits: make(map[string][]bool),
}

func ensureCapacity(bits []bool, offset int) []bool {
	for len(bits) <= offset {
		bits = append(bits, false)
	}
	return bits
}

// BitSet sets the bit at offset to 1 for the given key.
func (b *Bitmap) BitSet(key string, offset int) error {
	if offset < 0 {
		return fmt.Errorf("offset must be non-negative")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.bits[key] = ensureCapacity(b.bits[key], offset)
	b.bits[key][offset] = true
	return nil
}

// BitGet returns the bit value at offset for the given key.
func (b *Bitmap) BitGet(key string, offset int) (bool, error) {
	if offset < 0 {
		return false, fmt.Errorf("offset must be non-negative")
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	bits, ok := b.bits[key]
	if !ok || offset >= len(bits) {
		return false, nil
	}
	return bits[offset], nil
}

// BitCount returns the number of set bits for the given key.
func (b *Bitmap) BitCount(key string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	count := 0
	for _, v := range b.bits[key] {
		if v {
			count++
		}
	}
	return count
}

// BitClear sets the bit at offset to 0 for the given key.
func (b *Bitmap) BitClear(key string, offset int) error {
	if offset < 0 {
		return fmt.Errorf("offset must be non-negative")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	bits, ok := b.bits[key]
	if !ok || offset >= len(bits) {
		return nil
	}
	bits[offset] = false
	return nil
}

// NewBitmap returns the global bitmap instance.
func NewBitmap() *Bitmap {
	return globalBitmap
}
