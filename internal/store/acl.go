package store

import (
	"errors"
	"sync"
)

// Permission represents an allowed operation.
type Permission uint8

const (
	PermRead  Permission = 1 << iota // 1
	PermWrite                        // 2
	PermDelete                       // 4
	PermAdmin                        // 8
)

var ErrACLDenied = errors.New("acl: permission denied")
var ErrACLTokenNotFound = errors.New("acl: token not found")

// ACLManager stores token → permission mappings.
type ACLManager struct {
	mu     sync.RWMutex
	tokens map[string]Permission
}

func NewACLManager() *ACLManager {
	return &ACLManager{tokens: make(map[string]Permission)}
}

// SetToken creates or updates a token with the given permissions.
func (a *ACLManager) SetToken(token string, perm Permission) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.tokens[token] = perm
}

// RevokeToken removes a token entirely.
func (a *ACLManager) RevokeToken(token string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.tokens, token)
}

// Check returns nil if the token holds all of the requested permissions.
func (a *ACLManager) Check(token string, required Permission) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	perm, ok := a.tokens[token]
	if !ok {
		return ErrACLTokenNotFound
	}
	if perm&required != required {
		return ErrACLDenied
	}
	return nil
}

// GetPermissions returns the raw permission bits for a token.
func (a *ACLManager) GetPermissions(token string) (Permission, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	perm, ok := a.tokens[token]
	if !ok {
		return 0, ErrACLTokenNotFound
	}
	return perm, nil
}
