// Package rollback provides functionality to restore Vault secrets
// to a previously captured snapshot state.
package rollback

import (
	"context"
	"fmt"

	"github.com/youorg/vaultpatch/internal/snapshot"
)

// Writer defines the interface for writing secrets back to Vault.
type Writer interface {
	WriteSecret(ctx context.Context, path string, data map[string]interface{}) error
	DeleteSecret(ctx context.Context, path string) error
}

// Result holds the outcome of a rollback operation.
type Result struct {
	Restored []string
	Deleted  []string
	Errors   []error
}

// Summary returns a human-readable summary of the rollback result.
func (r Result) Summary() string {
	return fmt.Sprintf("rollback: %d restored, %d deleted, %d errors",
		len(r.Restored), len(r.Deleted), len(r.Errors))
}

// Options configures rollback behaviour.
type Options struct {
	DryRun bool
}

// Apply restores Vault secrets to the state captured in the given snapshot.
// If DryRun is true, no writes are performed.
func Apply(ctx context.Context, w Writer, snap *snapshot.Snapshot, current map[string]map[string]interface{}, opts Options) (Result, error) {
	if snap == nil {
		return Result{}, fmt.Errorf("rollback: snapshot must not be nil")
	}

	var res Result

	// Restore secrets present in the snapshot.
	for path, data := range snap.Secrets {
		if opts.DryRun {
			res.Restored = append(res.Restored, path)
			continue
		}
		if err := w.WriteSecret(ctx, path, data); err != nil {
			res.Errors = append(res.Errors, fmt.Errorf("write %s: %w", path, err))
			continue
		}
		res.Restored = append(res.Restored, path)
	}

	// Delete secrets that exist now but were absent in the snapshot.
	for path := range current {
		if _, ok := snap.Secrets[path]; !ok {
			if opts.DryRun {
				res.Deleted = append(res.Deleted, path)
				continue
			}
			if err := w.DeleteSecret(ctx, path); err != nil {
				res.Errors = append(res.Errors, fmt.Errorf("delete %s: %w", path, err))
				continue
			}
			res.Deleted = append(res.Deleted, path)
		}
	}

	return res, nil
}
