package config

import (
	"fmt"
	"os"
)

// DefaultConfigPaths lists candidate locations for the vaultpatch config file.
var DefaultConfigPaths = []string{
	".vaultpatch.yaml",
	".vaultpatch.yml",
	"vaultpatch.yaml",
	"vaultpatch.yml",
}

// LoadDefault attempts to load config from well-known default paths.
// It returns an error if none of the paths exist or if parsing fails.
func LoadDefault() (*Config, error) {
	for _, p := range DefaultConfigPaths {
		if _, err := os.Stat(p); err == nil {
			return Load(p)
		}
	}
	return nil, fmt.Errorf(
		"no config file found; tried: %v", DefaultConfigPaths,
	)
}

// MustLoad loads config from path and panics on error.
// Intended for use in tests or CLI bootstrapping where failure is fatal.
func MustLoad(path string) *Config {
	cfg, err := Load(path)
	if err != nil {
		panic(fmt.Sprintf("vaultpatch: failed to load config: %v", err))
	}
	return cfg
}

// EnvForName returns the Environment matching the given name, or an error.
func (c *Config) EnvForName(name string) (*Environment, error) {
	for i := range c.Environments {
		if c.Environments[i].Name == name {
			return &c.Environments[i], nil
		}
	}
	return nil, fmt.Errorf("environment %q not found in config", name)
}
