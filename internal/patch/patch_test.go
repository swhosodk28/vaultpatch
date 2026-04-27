package patch_test

import (
	"errors"
	"testing"

	"github.com/vaultpatch/internal/diff"
	"github.com/vaultpatch/internal/patch"
)

// mockWriter implements patch.VaultWriter for testing.
type mockWriter struct {
	written map[string]map[string]interface{}
	deleted []string
	failOn  string
}

func newMockWriter() *mockWriter {
	return &mockWriter{written: make(map[string]map[string]interface{})}
}

func (m *mockWriter) WriteSecret(path string, data map[string]interface{}) error {
	if m.failOn == path {
		return errors.New("simulated write error")
	}
	m.written[path] = data
	return nil
}

func (m *mockWriter) DeleteSecret(path string) error {
	if m.failOn == path {
		return errors.New("simulated delete error")
	}
	m.deleted = append(m.deleted, path)
	return nil
}

func TestApply_DryRun(t *testing.T) {
	w := newMockWriter()
	diffs := []diff.Diff{
		{Path: "secret/foo", Type: diff.Added, NewValue: map[string]interface{}{"key": "val"}},
	}
	results, err := patch.Apply(diffs, w, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Success {
		t.Errorf("expected 1 successful dry-run result")
	}
	if len(w.written) != 0 {
		t.Errorf("dry run should not write to vault")
	}
}

func TestApply_AddAndDelete(t *testing.T) {
	w := newMockWriter()
	diffs := []diff.Diff{
		{Path: "secret/new", Type: diff.Added, NewValue: map[string]interface{}{"a": "1"}},
		{Path: "secret/old", Type: diff.Removed},
	}
	results, err := patch.Apply(diffs, w, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if _, ok := w.written["secret/new"]; !ok {
		t.Errorf("expected secret/new to be written")
	}
	if len(w.deleted) != 1 || w.deleted[0] != "secret/old" {
		t.Errorf("expected secret/old to be deleted")
	}
}

func TestApply_WriteFailure(t *testing.T) {
	w := newMockWriter()
	w.failOn = "secret/bad"
	diffs := []diff.Diff{
		{Path: "secret/bad", Type: diff.Modified, NewValue: map[string]interface{}{"x": "y"}},
	}
	results, err := patch.Apply(diffs, w, false)
	if err != nil {
		t.Fatalf("unexpected fatal error: %v", err)
	}
	if results[0].Success {
		t.Errorf("expected failure result for bad path")
	}
}

func TestApply_NilWriter(t *testing.T) {
	_, err := patch.Apply(nil, nil, false)
	if err == nil {
		t.Error("expected error for nil writer")
	}
}

func TestSummary(t *testing.T) {
	results := []patch.Result{
		{Success: true},
		{Success: true},
		{Success: false},
	}
	got := patch.Summary(results)
	expected := "patch complete: 2 succeeded, 1 failed"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}
