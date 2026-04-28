package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpatch/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "vaultpatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	yaml := `
vault:
  address: "https://vault.example.com"
  token_env: "MY_VAULT_TOKEN"
environments:
  - name: staging
    prefix: secret/staging
  - name: production
    prefix: secret/production
`
	path := writeTemp(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "https://vault.example.com" {
		t.Errorf("expected vault address, got %q", cfg.Vault.Address)
	}
	if len(cfg.Environments) != 2 {
		t.Errorf("expected 2 environments, got %d", len(cfg.Environments))
	}
	if cfg.Environments[0].Name != "staging" {
		t.Errorf("expected staging, got %q", cfg.Environments[0].Name)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_MissingVaultAddress(t *testing.T) {
	yaml := `
vault:
  address: ""
environments:
  - name: dev
    prefix: secret/dev
`
	path := writeTemp(t, yaml)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing vault address")
	}
}

func TestLoad_NoEnvironments(t *testing.T) {
	yaml := `
vault:
  address: "https://vault.example.com"
environments: []
`
	path := writeTemp(t, yaml)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty environments")
	}
}

func TestLoad_DefaultTokenEnv(t *testing.T) {
	yaml := `
vault:
  address: "https://vault.example.com"
environments:
  - name: dev
    prefix: secret/dev
`
	path := writeTemp(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.TokenEnv != "VAULT_TOKEN" {
		t.Errorf("expected default token_env VAULT_TOKEN, got %q", cfg.Vault.TokenEnv)
	}
}
