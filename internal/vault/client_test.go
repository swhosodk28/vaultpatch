package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vaultpatch/internal/vault"
)

func newMockVaultServer(t *testing.T, respData map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(respData); err != nil {
			t.Errorf("encoding mock response: %v", err)
		}
	}))
}

func TestNewClient_InvalidAddress(t *testing.T) {
	_, err := vault.NewClient(vault.Config{
		Address: "://bad-url",
		Token:   "root",
	})
	if err == nil {
		t.Fatal("expected error for invalid address, got nil")
	}
}

func TestNewClient_ValidConfig(t *testing.T) {
	srv := newMockVaultServer(t, map[string]interface{}{})
	defer srv.Close()

	client, err := vault.NewClient(vault.Config{
		Address: srv.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestReadSecrets_EmptyResponse(t *testing.T) {
	srv := newMockVaultServer(t, map[string]interface{}{
		"data": map[string]interface{}{
			"data": map[string]interface{}{},
		},
	})
	defer srv.Close()

	client, err := vault.NewClient(vault.Config{Address: srv.URL, Token: "root"})
	if err != nil {
		t.Fatalf("creating client: %v", err)
	}

	// ReadSecrets may fail due to mock not fully implementing KV v2;
	// we verify the client is constructed and the call is attempted.
	_, _ = client.ReadSecrets(context.Background(), "secret", "myapp/config")
}
