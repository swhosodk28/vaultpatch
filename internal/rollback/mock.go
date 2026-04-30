package rollback

import (
	"context"
	"fmt"
	"sync"
)

// MockWriter is an in-memory Writer implementation for use in tests.
type MockWriter struct {
	mu      sync.Mutex
	Secrets map[string]map[string]interface{}
	Deleted []string

	// WriteErr, if set, is returned by every WriteSecret call.
	WriteErr error
	// DeleteErr, if set, is returned by every DeleteSecret call.
	DeleteErr error
}

// NewMockWriter returns a ready-to-use MockWriter.
func NewMockWriter() *MockWriter {
	return &MockWriter{
		Secrets: make(map[string]map[string]interface{}),
	}
}

// WriteSecret stores data at path or returns WriteErr if set.
func (m *MockWriter) WriteSecret(_ context.Context, path string, data map[string]interface{}) error {
	if m.WriteErr != nil {
		return m.WriteErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	copy := make(map[string]interface{}, len(data))
	for k, v := range data {
		copy[k] = v
	}
	m.Secrets[path] = copy
	return nil
}

// DeleteSecret removes path from the in-memory store or returns DeleteErr if set.
func (m *MockWriter) DeleteSecret(_ context.Context, path string) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.Secrets[path]; !ok {
		return fmt.Errorf("mock: path %q not found", path)
	}
	delete(m.Secrets, path)
	m.Deleted = append(m.Deleted, path)
	return nil
}
