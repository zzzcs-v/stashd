package store

import (
	"testing"
)

func freshBitmap() *Bitmap {
	return &Bitmap{bits: make(map[string][]bool)}
}

func TestBitSetAndGet(t *testing.T) {
	b := freshBitmap()
	_ = b.BitSet("mykey", 7)
	val, err := b.BitGet("mykey", 7)
	if err != nil || !val {
		t.Fatalf("expected bit 7 to be set, got %v err=%v", val, err)
	}
}

func TestBitGetUnset(t *testing.T) {
	b := freshBitmap()
	val, err := b.BitGet("missing", 3)
	if err != nil || val {
		t.Fatalf("expected false for missing key, got %v err=%v", val, err)
	}
}

func TestBitCount(t *testing.T) {
	b := freshBitmap()
	_ = b.BitSet("k", 0)
	_ = b.BitSet("k", 3)
	_ = b.BitSet("k", 7)
	if c := b.BitCount("k"); c != 3 {
		t.Fatalf("expected count 3, got %d", c)
	}
}

func TestBitClear(t *testing.T) {
	b := freshBitmap()
	_ = b.BitSet("k", 2)
	_ = b.BitClear("k", 2)
	val, _ := b.BitGet("k", 2)
	if val {
		t.Fatal("expected bit to be cleared")
	}
}

func TestBitNegativeOffset(t *testing.T) {
	b := freshBitmap()
	if err := b.BitSet("k", -1); err == nil {
		t.Fatal("expected error for negative offset")
	}
	if _, err := b.BitGet("k", -1); err == nil {
		t.Fatal("expected error for negative offset")
	}
}

func TestBitCountEmpty(t *testing.T) {
	b := freshBitmap()
	if c := b.BitCount("nokey"); c != 0 {
		t.Fatalf("expected 0 for missing key, got %d", c)
	}
}
