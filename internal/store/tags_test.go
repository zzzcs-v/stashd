package store

import (
	"testing"
)

func TestSetAndGetTags(t *testing.T) {
	s := New()
	s.Set("foo", "bar", 0)
	s.SetTags("foo", []string{"hot", "featured"})

	tags, ok := s.GetTags("foo")
	if !ok {
		t.Fatal("expected tags to exist")
	}
	if len(tags) != 2 || tags[0] != "hot" || tags[1] != "featured" {
		t.Fatalf("unexpected tags: %v", tags)
	}
}

func TestGetTagsMissing(t *testing.T) {
	s := New()
	_, ok := s.GetTags("nonexistent")
	if ok {
		t.Fatal("expected no tags for missing key")
	}
}

func TestDeleteTags(t *testing.T) {
	s := New()
	s.Set("foo", "bar", 0)
	s.SetTags("foo", []string{"hot"})
	s.DeleteTags("foo")

	_, ok := s.GetTags("foo")
	if ok {
		t.Fatal("expected tags to be deleted")
	}
}

func TestKeysByTag(t *testing.T) {
	s := New()
	s.Set("a", "1", 0)
	s.Set("b", "2", 0)
	s.Set("c", "3", 0)
	s.SetTags("a", []string{"sale", "new"})
	s.SetTags("b", []string{"sale"})
	s.SetTags("c", []string{"new"})

	keys := s.KeysByTag("sale")
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys with tag 'sale', got %d", len(keys))
	}
}

func TestKeysByTagNoMatch(t *testing.T) {
	s := New()
	s.Set("a", "1", 0)
	s.SetTags("a", []string{"hot"})

	keys := s.KeysByTag("cold")
	if len(keys) != 0 {
		t.Fatalf("expected 0 keys, got %d", len(keys))
	}
}
