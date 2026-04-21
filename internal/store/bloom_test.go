package store

import (
	"fmt"
	"testing"
)

func TestBloomAddAndContains(t *testing.T) {
	bf := NewBloomFilter(100, 0.01)
	bf.Add("hello")
	bf.Add("world")

	if !bf.MayContain("hello") {
		t.Error("expected 'hello' to be found")
	}
	if !bf.MayContain("world") {
		t.Error("expected 'world' to be found")
	}
}

func TestBloomDefinitelyNotContains(t *testing.T) {
	bf := NewBloomFilter(100, 0.01)
	bf.Add("apple")

	// "banana" was never added; it should not be found (no false positive for this case)
	// We check a string that is very unlikely to collide.
	if bf.MayContain("zzzzzzzzzzzzzzzzzzzzzzzzz") {
		t.Log("false positive detected — acceptable but rare")
	}
	if !bf.MayContain("apple") {
		t.Error("expected 'apple' to be present")
	}
}

func TestBloomReset(t *testing.T) {
	bf := NewBloomFilter(100, 0.01)
	bf.Add("key1")
	bf.Reset()

	if bf.MayContain("key1") {
		t.Error("expected filter to be empty after reset")
	}
}

func TestBloomFalsePositiveRate(t *testing.T) {
	n := 1000
	bf := NewBloomFilter(n, 0.01)
	for i := 0; i < n; i++ {
		bf.Add(fmt.Sprintf("item-%d", i))
	}

	falsePositives := 0
	trials := 1000
	for i := n; i < n+trials; i++ {
		if bf.MayContain(fmt.Sprintf("item-%d", i)) {
			falsePositives++
		}
	}

	rate := float64(falsePositives) / float64(trials)
	if rate > 0.05 {
		t.Errorf("false positive rate too high: %.2f", rate)
	}
}

func TestBloomManagerGetOrCreate(t *testing.T) {
	bm := NewBloomManager()
	f1 := bm.GetOrCreate("myfilter", 100, 0.01)
	f1.Add("ping")

	f2 := bm.GetOrCreate("myfilter", 100, 0.01)
	if !f2.MayContain("ping") {
		t.Error("expected same filter instance to be returned")
	}
}

func TestBloomManagerDelete(t *testing.T) {
	bm := NewBloomManager()
	bm.GetOrCreate("temp", 50, 0.05)
	bm.Delete("temp")

	_, ok := bm.Get("temp")
	if ok {
		t.Error("expected filter to be deleted")
	}
}
