package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the vaultpatch configuration.
type Config struct {
	Vault    VaultConfig    `yaml:"vault"`
	Environments []Environment `yaml:"environments"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address   string `yaml:"address"`
	TokenEnv  string `yaml:"token_env"`
	Namespace string `yaml:"namespace"`
}

// Environment represents a named Vault environment with a secret path prefix.
type Environment struct {
	Name   string `yaml:"name"`
	Prefix string `yaml:"prefix"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// Validate checks that required fields are present.
func (c *Config) Validate() error {
	if c.Vault.Address == "" {
		return fmt.Errorf("vault.address is required")
	}
	if c.Vault.TokenEnv == "" {
		c.Vault.TokenEnv = "VAULT_TOKEN"
	}
	if len(c.Environments) == 0 {
		return fmt.Errorf("at least one environment must be defined")
	}
	for i, env := range c.Environments {
		if env.Name == "" {
			return fmt.Errorf("environments[%d].name is required", i)
		}
		if env.Prefix == "" {
			return fmt.Errorf("environments[%d].prefix is required for %q", i, env.Name)
		}
	}
	return nil
}

// Token resolves the Vault token from the configured environment variable.
func (c *Config) Token() string {
	return os.Getenv(c.Vault.TokenEnv)
}
