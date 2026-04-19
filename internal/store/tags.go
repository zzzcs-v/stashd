package store

import "strings"

// TagKey returns the internal key used to store tags for a given key.
func TagKey(key string) string {
	return "__tag__:" + key
}

// SetTags assigns a list of tags to a key.
func (s *Store) SetTags(key string, tags []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[TagKey(key)] = entry{value: strings.Join(tags, ",")}
}

// GetTags returns the tags associated with a key.
func (s *Store) GetTags(key string) ([]string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.data[TagKey(key)]
	if !ok || e.value == "" {
		return nil, false
	}
	return strings.Split(e.value, ","), true
}

// DeleteTags removes tags for a key.
func (s *Store) DeleteTags(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, TagKey(key))
}

// KeysByTag returns all keys that have the given tag.
func (s *Store) KeysByTag(tag string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var keys []string
	prefix := "__tag__:"
	for k, e := range s.data {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		for _, t := range strings.Split(e.value, ",") {
			if t == tag {
				keys = append(keys, strings.TrimPrefix(k, prefix))
				break
			}
		}
	}
	return keys
}
