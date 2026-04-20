package store

import (
	"testing"
)

func TestSetAddAndMembers(t *testing.T) {
	s := New()
	added := s.SetAdd("myset", "a", "b", "c")
	if added != 3 {
		t.Fatalf("expected 3 added, got %d", added)
	}
	members, ok := s.SetMembers("myset")
	if !ok {
		t.Fatal("expected set to exist")
	}
	if len(members) != 3 {
		t.Fatalf("expected 3 members, got %d", len(members))
	}
}

func TestSetAddDuplicates(t *testing.T) {
	s := New()
	s.SetAdd("myset", "a", "b")
	added := s.SetAdd("myset", "b", "c")
	if added != 1 {
		t.Fatalf("expected 1 new member added, got %d", added)
	}
	members, _ := s.SetMembers("myset")
	if len(members) != 3 {
		t.Fatalf("expected 3 total members, got %d", len(members))
	}
}

func TestSetRemove(t *testing.T) {
	s := New()
	s.SetAdd("myset", "a", "b", "c")
	removed := s.SetRemove("myset", "b")
	if removed != 1 {
		t.Fatalf("expected 1 removed, got %d", removed)
	}
	members, _ := s.SetMembers("myset")
	if len(members) != 2 {
		t.Fatalf("expected 2 members after remove, got %d", len(members))
	}
}

func TestSetRemoveMissingKey(t *testing.T) {
	s := New()
	removed := s.SetRemove("ghost", "x")
	if removed != 0 {
		t.Fatalf("expected 0 removed from missing key, got %d", removed)
	}
}

func TestSetIsMember(t *testing.T) {
	s := New()
	s.SetAdd("myset", "hello", "world")
	if !s.SetIsMember("myset", "hello") {
		t.Error("expected 'hello' to be a member")
	}
	if s.SetIsMember("myset", "nope") {
		t.Error("expected 'nope' to not be a member")
	}
}

func TestSetMembersMissingKey(t *testing.T) {
	s := New()
	_, ok := s.SetMembers("missing")
	if ok {
		t.Error("expected false for missing key")
	}
}
