package vault

import (
	"context"
	"fmt"
)

// MockClient implements a simple in-memory Vault client for use in tests.
type MockClient struct {
	store map[string]map[string]string
	ReadErr  error
	WriteErr error
}

// NewMockClient returns an initialised MockClient.
func NewMockClient() *MockClient {
	return &MockClient{
		store: make(map[string]map[string]string),
	}
}

// Seed pre-populates the mock store at the given path.
func (m *MockClient) Seed(path string, data map[string]string) {
	copy := make(map[string]string, len(data))
	for k, v := range data {
		copy[k] = v
	}
	m.store[path] = copy
}

// ReadSecrets returns the in-memory secrets for the given path.
func (m *MockClient) ReadSecrets(_ context.Context, _, path string) (map[string]string, error) {
	if m.ReadErr != nil {
		return nil, m.ReadErr
	}
	data, ok := m.store[path]
	if !ok {
		return map[string]string{}, nil
	}
	copy := make(map[string]string, len(data))
	for k, v := range data {
		copy[k] = v
	}
	return copy, nil
}

// WriteSecrets stores the provided secrets at the given path.
func (m *MockClient) WriteSecrets(_ context.Context, _, path string, data map[string]string) error {
	if m.WriteErr != nil {
		return fmt.Errorf("mock write error: %w", m.WriteErr)
	}
	copy := make(map[string]string, len(data))
	for k, v := range data {
		copy[k] = v
	}
	m.store[path] = copy
	return nil
}
