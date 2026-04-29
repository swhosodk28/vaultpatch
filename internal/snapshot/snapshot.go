// Package snapshot provides functionality to capture and restore
// Vault secret state at a point in time for a given environment.
package snapshot

import (
	"encoding/json"
	"fmt"
	"time"
)

// Snapshot represents the state of secrets in a Vault environment
// at a specific point in time.
type Snapshot struct {
	Environment string            `json:"environment"`
	Path        string            `json:"path"`
	CapturedAt  time.Time         `json:"captured_at"`
	Secrets     map[string]string `json:"secrets"`
}

// Reader is the interface for reading secrets from a Vault path.
type Reader interface {
	ReadSecrets(path string) (map[string]string, error)
}

// Capture reads secrets from the given path using the provided reader
// and returns a Snapshot representing the current state.
func Capture(env, path string, r Reader) (*Snapshot, error) {
	if env == "" {
		return nil, fmt.Errorf("snapshot: environment must not be empty")
	}
	if path == "" {
		return nil, fmt.Errorf("snapshot: path must not be empty")
	}

	secrets, err := r.ReadSecrets(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: failed to read secrets: %w", err)
	}

	return &Snapshot{
		Environment: env,
		Path:        path,
		CapturedAt:  time.Now().UTC(),
		Secrets:     secrets,
	}, nil
}

// Marshal serialises the snapshot to JSON bytes.
func (s *Snapshot) Marshal() ([]byte, error) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("snapshot: marshal failed: %w", err)
	}
	return data, nil
}

// Unmarshal deserialises JSON bytes into a Snapshot.
func Unmarshal(data []byte) (*Snapshot, error) {
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}
	return &s, nil
}
