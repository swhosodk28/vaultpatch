package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// EventType represents the type of audit event.
type EventType string

const (
	EventDiff  EventType = "diff"
	EventApply EventType = "apply"
	EventRead  EventType = "read"
)

// Event represents a single audit log entry.
type Event struct {
	Timestamp   time.Time         `json:"timestamp"`
	Type        EventType         `json:"type"`
	Environment string            `json:"environment"`
	Path        string            `json:"path"`
	Details     map[string]string `json:"details,omitempty"`
	Error       string            `json:"error,omitempty"`
}

// Logger writes structured audit events to an output.
type Logger struct {
	out io.Writer
}

// NewLogger creates a new audit Logger writing to the given writer.
// Pass nil to use os.Stdout.
func NewLogger(out io.Writer) *Logger {
	if out == nil {
		out = os.Stdout
	}
	return &Logger{out: out}
}

// Log writes an audit event as a JSON line.
func (l *Logger) Log(event Event) error {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("audit: marshal event: %w", err)
	}
	_, err = fmt.Fprintf(l.out, "%s\n", data)
	return err
}

// LogDiff records a diff operation event.
func (l *Logger) LogDiff(env, path string, added, removed, modified int) error {
	return l.Log(Event{
		Type:        EventDiff,
		Environment: env,
		Path:        path,
		Details: map[string]string{
			"added":    fmt.Sprintf("%d", added),
			"removed":  fmt.Sprintf("%d", removed),
			"modified": fmt.Sprintf("%d", modified),
		},
	})
}

// LogApply records an apply operation event.
func (l *Logger) LogApply(env, path string, dryRun bool, err error) error {
	e := Event{
		Type:        EventApply,
		Environment: env,
		Path:        path,
		Details:     map[string]string{"dry_run": fmt.Sprintf("%v", dryRun)},
	}
	if err != nil {
		e.Error = err.Error()
	}
	return l.Log(e)
}
