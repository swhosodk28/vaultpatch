package audit

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestLog_WritesJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	event := Event{
		Timestamp:   time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Type:        EventDiff,
		Environment: "staging",
		Path:        "secret/app",
	}

	if err := logger.Log(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := strings.TrimSpace(buf.String())
	var got Event
	if err := json.Unmarshal([]byte(line), &got); err != nil {
		t.Fatalf("failed to parse output JSON: %v", err)
	}

	if got.Type != EventDiff {
		t.Errorf("expected type %q, got %q", EventDiff, got.Type)
	}
	if got.Environment != "staging" {
		t.Errorf("expected environment %q, got %q", "staging", got.Environment)
	}
}

func TestLog_SetsTimestampIfZero(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	if err := logger.Log(Event{Type: EventRead}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("failed to parse output JSON: %v", err)
	}
	if got.Timestamp.IsZero() {
		t.Error("expected timestamp to be set automatically")
	}
}

func TestLogDiff_WritesDetails(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	if err := logger.LogDiff("prod", "secret/db", 2, 1, 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("failed to parse output JSON: %v", err)
	}

	if got.Details["added"] != "2" {
		t.Errorf("expected added=2, got %q", got.Details["added"])
	}
	if got.Details["modified"] != "3" {
		t.Errorf("expected modified=3, got %q", got.Details["modified"])
	}
}

func TestLogApply_WithError(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf)

	applyErr := errors.New("permission denied")
	if err := logger.LogApply("dev", "secret/svc", false, applyErr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got Event
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &got); err != nil {
		t.Fatalf("failed to parse output JSON: %v", err)
	}

	if got.Error != "permission denied" {
		t.Errorf("expected error field %q, got %q", "permission denied", got.Error)
	}
	if got.Details["dry_run"] != "false" {
		t.Errorf("expected dry_run=false, got %q", got.Details["dry_run"])
	}
}

func TestNewLogger_NilUsesStdout(t *testing.T) {
	logger := NewLogger(nil)
	if logger.out == nil {
		t.Error("expected non-nil writer when nil passed")
	}
}
