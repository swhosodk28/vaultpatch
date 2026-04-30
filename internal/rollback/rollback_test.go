package rollback_test

import (
	"context"
	"errors"
	"testing"

	"github.com/youorg/vaultpatch/internal/rollback"
	"github.com/youorg/vaultpatch/internal/snapshot"
)

type mockWriter struct {
	written map[string]map[string]interface{}
	deleted []string
	writeErr error
	deleteErr error
}

func (m *mockWriter) WriteSecret(_ context.Context, path string, data map[string]interface{}) error {
	if m.writeErr != nil {
		return m.writeErr
	}
	m.written[path] = data
	return nil
}

func (m *mockWriter) DeleteSecret(_ context.Context, path string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.deleted = append(m.deleted, path)
	return nil
}

func newSnap(secrets map[string]map[string]interface{}) *snapshot.Snapshot {
	return &snapshot.Snapshot{Secrets: secrets}
}

func TestApply_NilSnapshot(t *testing.T) {
	w := &mockWriter{written: map[string]map[string]interface{}{}}
	_, err := rollback.Apply(context.Background(), w, nil, nil, rollback.Options{})
	if err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestApply_DryRun(t *testing.T) {
	snap := newSnap(map[string]map[string]interface{}{
		"secret/a": {"key": "val"},
	})
	current := map[string]map[string]interface{}{
		"secret/b": {"key": "other"},
	}
	w := &mockWriter{written: map[string]map[string]interface{}{}}
	res, err := rollback.Apply(context.Background(), w, snap, current, rollback.Options{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(w.written) != 0 {
		t.Error("dry run should not write secrets")
	}
	if len(res.Restored) != 1 || len(res.Deleted) != 1 {
		t.Errorf("expected 1 restored and 1 deleted, got %d/%d", len(res.Restored), len(res.Deleted))
	}
}

func TestApply_RestoresAndDeletes(t *testing.T) {
	snap := newSnap(map[string]map[string]interface{}{
		"secret/a": {"x": "1"},
	})
	current := map[string]map[string]interface{}{
		"secret/a": {"x": "2"},
		"secret/c": {"y": "3"},
	}
	w := &mockWriter{written: map[string]map[string]interface{}{}}
	res, err := rollback.Apply(context.Background(), w, snap, current, rollback.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Restored) != 1 {
		t.Errorf("expected 1 restored, got %d", len(res.Restored))
	}
	if len(res.Deleted) != 1 {
		t.Errorf("expected 1 deleted, got %d", len(res.Deleted))
	}
	if len(res.Errors) != 0 {
		t.Errorf("expected no errors, got %v", res.Errors)
	}
}

func TestApply_WriteError(t *testing.T) {
	snap := newSnap(map[string]map[string]interface{}{
		"secret/a": {"k": "v"},
	})
	w := &mockWriter{written: map[string]map[string]interface{}{}, writeErr: errors.New("vault unavailable")}
	res, err := rollback.Apply(context.Background(), w, snap, nil, rollback.Options{})
	if err != nil {
		t.Fatalf("unexpected hard error: %v", err)
	}
	if len(res.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(res.Errors))
	}
}

func TestResult_Summary(t *testing.T) {
	r := rollback.Result{
		Restored: []string{"a", "b"},
		Deleted:  []string{"c"},
		Errors:   []error{errors.New("oops")},
	}
	got := r.Summary()
	expected := "rollback: 2 restored, 1 deleted, 1 errors"
	if got != expected {
		t.Errorf("Summary() = %q, want %q", got, expected)
	}
}
