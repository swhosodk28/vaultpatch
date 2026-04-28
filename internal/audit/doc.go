// Package audit provides structured JSON audit logging for vaultpatch operations.
//
// It records diff, apply, and read events with timestamps, environment context,
// and optional error information. Events can be written to any io.Writer or
// persisted to a file using FileLogger.
//
// Basic usage:
//
//	logger := audit.NewLogger(os.Stdout)
//	logger.LogDiff("production", "secret/myapp", 2, 0, 1)
//
// File-backed usage:
//
//	fl, err := audit.NewFileLogger("/var/log/vaultpatch/audit.log")
//	if err != nil { ... }
//	defer fl.Close()
//	fl.LogApply("staging", "secret/myapp", false, nil)
package audit
