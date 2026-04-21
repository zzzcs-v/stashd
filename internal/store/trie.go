package store

import "sync"

// TrieNode represents a node in the trie.
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

// Trie is a concurrent prefix tree for key autocomplete and prefix search.
type Trie struct {
	mu   sync.RWMutex
	root *TrieNode
}

// NewTrie creates a new empty Trie.
func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{children: make(map[rune]*TrieNode)},
	}
}

// Insert adds a key to the trie.
func (t *Trie) Insert(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	node := t.root
	for _, ch := range key {
		if _, ok := node.children[ch]; !ok {
			node.children[ch] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

// Delete removes a key from the trie. Returns true if the key existed.
func (t *Trie) Delete(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.delete(t.root, []rune(key), 0)
}

func (t *Trie) delete(node *TrieNode, runes []rune, depth int) bool {
	if depth == len(runes) {
		if !node.isEnd {
			return false
		}
		node.isEnd = false
		return true
	}
	ch := runes[depth]
	child, ok := node.children[ch]
	if !ok {
		return false
	}
	found := t.delete(child, runes, depth+1)
	if found && !child.isEnd && len(child.children) == 0 {
		delete(node.children, ch)
	}
	return found
}

// HasPrefix returns true if any inserted key starts with the given prefix.
func (t *Trie) HasPrefix(prefix string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	node := t.root
	for _, ch := range prefix {
		if _, ok := node.children[ch]; !ok {
			return false
		}
		node = node.children[ch]
	}
	return true
}

// Search returns all keys in the trie that start with the given prefix.
func (t *Trie) Search(prefix string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	node := t.root
	for _, ch := range prefix {
		child, ok := node.children[ch]
		if !ok {
			return nil
		}
		node = child
	}
	var results []string
	t.collect(node, prefix, &results)
	return results
}

func (t *Trie) collect(node *TrieNode, prefix string, results *[]string) {
	if node.isEnd {
		*results = append(*results, prefix)
	}
	for ch, child := range node.children {
		t.collect(child, prefix+string(ch), results)
	}
}
