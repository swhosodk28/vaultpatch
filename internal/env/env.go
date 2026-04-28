package env

import (
	"fmt"
	"strings"
)

// Environment represents a named Vault environment with its path prefix.
type Environment struct {
	Name   string
	Prefix string
}

// Resolver resolves secret paths for a given environment.
type Resolver struct {
	envs map[string]Environment
}

// NewResolver creates a Resolver from a slice of environments.
func NewResolver(envs []Environment) *Resolver {
	m := make(map[string]Environment, len(envs))
	for _, e := range envs {
		m[strings.ToLower(e.Name)] = e
	}
	return &Resolver{envs: m}
}

// Resolve returns the full Vault path for a secret key in the given environment.
func (r *Resolver) Resolve(envName, secretKey string) (string, error) {
	e, ok := r.envs[strings.ToLower(envName)]
	if !ok {
		return "", fmt.Errorf("unknown environment: %q", envName)
	}
	prefix := strings.TrimRight(e.Prefix, "/")
	key := strings.TrimLeft(secretKey, "/")
	return prefix + "/" + key, nil
}

// List returns all registered environment names.
func (r *Resolver) List() []string {
	names := make([]string, 0, len(r.envs))
	for name := range r.envs {
		names = append(names, name)
	}
	return names
}

// Has reports whether the named environment is registered.
func (r *Resolver) Has(envName string) bool {
	_, ok := r.envs[strings.ToLower(envName)]
	return ok
}
