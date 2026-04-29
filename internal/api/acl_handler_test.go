package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theleeeo/stashd/internal/store"
)

func newACLServer(t *testing.T) *httptest.Server {
	t.Helper()
	acl := store.NewACLManager()
	return httptest.NewServer(NewACLHandler(acl))
}

func TestACLHTTPSetAndCheck(t *testing.T) {
	srv := newACLServer(t)
	defer srv.Close()

	// Create token with read+write (perm=3)
	resp, err := http.Post(srv.URL+"/acl/token", "application/json",
		strings.NewReader(`{"token":"mytoken","perm":3}`))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}

	// Check read permission (perm=1) — should succeed
	resp, err = http.Get(srv.URL + "/acl/check?token=mytoken&perm=1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestACLHTTPDenied(t *testing.T) {
	srv := newACLServer(t)
	defer srv.Close()

	// Token with read-only (perm=1)
	http.Post(srv.URL+"/acl/token", "application/json",
		strings.NewReader(`{"token":"ro","perm":1}`))

	// Check delete (perm=4) — should be forbidden
	resp, err := http.Get(srv.URL + "/acl/check?token=ro&perm=4")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", resp.StatusCode)
	}
}

func TestACLHTTPTokenNotFound(t *testing.T) {
	srv := newACLServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/acl/check?token=ghost&perm=1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestACLHTTPRevoke(t *testing.T) {
	srv := newACLServer(t)
	defer srv.Close()

	http.Post(srv.URL+"/acl/token", "application/json",
		strings.NewReader(`{"token":"tmp","perm":3}`))

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/acl/token?token=tmp", nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}

	resp, _ = http.Get(srv.URL + "/acl/check?token=tmp&perm=1")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 after revoke, got %d", resp.StatusCode)
	}
}

func TestACLHTTPMethodNotAllowed(t *testing.T) {
	srv := newACLServer(t)
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodPatch, srv.URL+"/acl/token", nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", resp.StatusCode)
	}
}
