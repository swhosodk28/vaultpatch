package snapshot_test

import (
	"errors"
	"testing"
	"time"

	"github.com/your-org/vaultpatch/internal/snapshot"
)

// stubReader implements snapshot.Reader for testing.
type stubReader struct {
	secrets map[string]string
	err     error
}

func (s *stubReader) ReadSecrets(_ string) (map[string]string, error) {
	return s.secrets, s.err
}

func TestCapture_Success(t *testing.T) {
	r := &stubReader{secrets: map[string]string{"key": "value", "db": "postgres"}}
	snap, err := snapshot.Capture("prod", "secret/prod", r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Environment != "prod" {
		t.Errorf("expected environment 'prod', got %q", snap.Environment)
	}
	if snap.Path != "secret/prod" {
		t.Errorf("expected path 'secret/prod', got %q", snap.Path)
	}
	if len(snap.Secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(snap.Secrets))
	}
	if snap.CapturedAt.IsZero() {
		t.Error("expected CapturedAt to be set")
	}
}

func TestCapture_EmptyEnvironment(t *testing.T) {
	r := &stubReader{secrets: map[string]string{}}
	_, err := snapshot.Capture("", "secret/prod", r)
	if err == nil {
		t.Fatal("expected error for empty environment")
	}
}

func TestCapture_EmptyPath(t *testing.T) {
	r := &stubReader{secrets: map[string]string{}}
	_, err := snapshot.Capture("prod", "", r)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestCapture_ReaderError(t *testing.T) {
	r := &stubReader{err: errors.New("vault unavailable")}
	_, err := snapshot.Capture("prod", "secret/prod", r)
	if err == nil {
		t.Fatal("expected error from reader")
	}
}

func TestMarshalUnmarshal_RoundTrip(t *testing.T) {
	original := &snapshot.Snapshot{
		Environment: "staging",
		Path:        "secret/staging",
		CapturedAt:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Secrets:     map[string]string{"foo": "bar"},
	}

	data, err := original.Marshal()
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	restored, err := snapshot.Unmarshal(data)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if restored.Environment != original.Environment {
		t.Errorf("environment mismatch: got %q want %q", restored.Environment, original.Environment)
	}
	if restored.Secrets["foo"] != "bar" {
		t.Errorf("secret mismatch: got %q want %q", restored.Secrets["foo"], "bar")
	}
}

func TestUnmarshal_InvalidJSON(t *testing.T) {
	_, err := snapshot.Unmarshal([]byte("not-json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
