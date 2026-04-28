package audit

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileLogger_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "logs", "audit.log")

	fl, err := NewFileLogger(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer fl.Close()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected log file to be created")
	}
}

func TestFileLogger_WritesAndReadsBack(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	fl, err := NewFileLogger(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := fl.LogDiff("staging", "secret/api", 1, 0, 2); err != nil {
		t.Fatalf("log error: %v", err)
	}

	if err := fl.Close(); err != nil {
		t.Fatalf("close error: %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open error: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("expected at least one line in log file")
	}

	var event Event
	if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if event.Environment != "staging" {
		t.Errorf("expected environment %q, got %q", "staging", event.Environment)
	}
	if event.Type != EventDiff {
		t.Errorf("expected type %q, got %q", EventDiff, event.Type)
	}
}

func TestFileLogger_AppendMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	for i := 0; i < 3; i++ {
		fl, err := NewFileLogger(path)
		if err != nil {
			t.Fatalf("open %d: %v", i, err)
		}
		if err := fl.LogApply("dev", "secret/x", true, nil); err != nil {
			t.Fatalf("log %d: %v", i, err)
		}
		if err := fl.Close(); err != nil {
			t.Fatalf("close %d: %v", i, err)
		}
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open error: %v", err)
	}
	defer f.Close()

	lines := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines++
	}
	if lines != 3 {
		t.Errorf("expected 3 log lines, got %d", lines)
	}
}
