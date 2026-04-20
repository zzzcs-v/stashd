package store

import "time"

// SetAdd adds one or more members to a set stored at key.
func (s *Store) SetAdd(key string, members ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	raw, ok := s.data[key]
	var set map[string]struct{}
	if ok {
		set, ok = raw.Value.(map[string]struct{})
		if !ok {
			set = make(map[string]struct{})
		}
	} else {
		set = make(map[string]struct{})
	}

	added := 0
	for _, m := range members {
		if _, exists := set[m]; !exists {
			set[m] = struct{}{}
			added++
		}
	}

	entry := s.data[key]
	entry.Value = set
	s.data[key] = entry
	return added
}

// SetRemove removes one or more members from a set stored at key.
func (s *Store) SetRemove(key string, members ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	raw, ok := s.data[key]
	if !ok {
		return 0
	}
	set, ok := raw.Value.(map[string]struct{})
	if !ok {
		return 0
	}

	removed := 0
	for _, m := range members {
		if _, exists := set[m]; exists {
			delete(set, m)
			removed++
		}
	}
	return removed
}

// SetMembers returns all members of the set stored at key.
func (s *Store) SetMembers(key string) ([]string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	raw, ok := s.data[key]
	if !ok {
		return nil, false
	}
	if raw.Expiry != nil && time.Now().After(*raw.Expiry) {
		return nil, false
	}
	set, ok := raw.Value.(map[string]struct{})
	if !ok {
		return nil, false
	}

	members := make([]string, 0, len(set))
	for m := range set {
		members = append(members, m)
	}
	return members, true
}

// SetIsMember checks if a value is a member of the set stored at key.
func (s *Store) SetIsMember(key, member string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	raw, ok := s.data[key]
	if !ok {
		return false
	}
	if raw.Expiry != nil && time.Now().After(*raw.Expiry) {
		return false
	}
	set, ok := raw.Value.(map[string]struct{})
	if !ok {
		return false
	}
	_, exists := set[member]
	return exists
}
