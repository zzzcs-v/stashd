package store

import (
	"testing"
)

func TestACLSetAndCheck(t *testing.T) {
	acl := NewACLManager()
	acl.SetToken("tok1", PermRead|PermWrite)

	if err := acl.Check("tok1", PermRead); err != nil {
		t.Fatalf("expected read allowed, got %v", err)
	}
	if err := acl.Check("tok1", PermWrite); err != nil {
		t.Fatalf("expected write allowed, got %v", err)
	}
}

func TestACLDeniedMissingPermission(t *testing.T) {
	acl := NewACLManager()
	acl.SetToken("readonly", PermRead)

	if err := acl.Check("readonly", PermDelete); err != ErrACLDenied {
		t.Fatalf("expected ErrACLDenied, got %v", err)
	}
}

func TestACLTokenNotFound(t *testing.T) {
	acl := NewACLManager()

	if err := acl.Check("ghost", PermRead); err != ErrACLTokenNotFound {
		t.Fatalf("expected ErrACLTokenNotFound, got %v", err)
	}
}

func TestACLRevokeToken(t *testing.T) {
	acl := NewACLManager()
	acl.SetToken("tmp", PermRead|PermWrite)
	acl.RevokeToken("tmp")

	if err := acl.Check("tmp", PermRead); err != ErrACLTokenNotFound {
		t.Fatalf("expected token not found after revoke, got %v", err)
	}
}

func TestACLGetPermissions(t *testing.T) {
	acl := NewACLManager()
	acl.SetToken("admin", PermRead|PermWrite|PermDelete|PermAdmin)

	perm, err := acl.GetPermissions("admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if perm != PermRead|PermWrite|PermDelete|PermAdmin {
		t.Fatalf("unexpected permissions: %d", perm)
	}
}

func TestACLGetPermissionsMissing(t *testing.T) {
	acl := NewACLManager()
	_, err := acl.GetPermissions("nobody")
	if err != ErrACLTokenNotFound {
		t.Fatalf("expected ErrACLTokenNotFound, got %v", err)
	}
}

func TestACLAdminCanDoAll(t *testing.T) {
	acl := NewACLManager()
	acl.SetToken("superuser", PermRead|PermWrite|PermDelete|PermAdmin)

	for _, p := range []Permission{PermRead, PermWrite, PermDelete, PermAdmin} {
		if err := acl.Check("superuser", p); err != nil {
			t.Fatalf("admin denied permission %d: %v", p, err)
		}
	}
}
