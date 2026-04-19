package store

import "strings"

// NamespacedKey returns a key prefixed with the given namespace.
func NamespacedKey(namespace, key string) string {
	if namespace == "" {
		return key
	}
	return namespace + ":" + key
}

// ParseNamespacedKey splits a namespaced key into namespace and key parts.
// If there is no namespace prefix, namespace will be empty.
func ParseNamespacedKey(nsKey string) (namespace, key string) {
	parts := strings.SplitN(nsKey, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", nsKey
}

// ListNamespace returns all keys belonging to a given namespace.
func (s *Store) ListNamespace(namespace string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	prefix := namespace + ":"
	var keys []string
	for k, item := range s.data {
		if strings.HasPrefix(k, prefix) && !isExpired(item) {
			_, bare := ParseNamespacedKey(k)
			keys = append(keys, bare)
		}
	}
	return keys
}

// DeleteNamespace removes all keys belonging to a given namespace.
func (s *Store) DeleteNamespace(namespace string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	prefix := namespace + ":"
	count := 0
	for k := range s.data {
		if strings.HasPrefix(k, prefix) {
			delete(s.data, k)
			count++
		}
	}
	return count
}
