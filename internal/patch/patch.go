package patch

import (
	"fmt"

	"github.com/vaultpatch/internal/diff"
)

// VaultWriter defines the interface for writing secrets to Vault.
type VaultWriter interface {
	WriteSecret(path string, data map[string]interface{}) error
	DeleteSecret(path string) error
}

// Result holds the outcome of applying a single diff entry.
type Result struct {
	Path    string
	Action  string
	Success bool
	Err     error
}

// Apply applies a slice of diff.Diff entries to Vault using the provided writer.
// It returns a slice of Results and a non-nil error only if a fatal issue occurs.
func Apply(diffs []diff.Diff, writer VaultWriter, dryRun bool) ([]Result, error) {
	if writer == nil {
		return nil, fmt.Errorf("vault writer must not be nil")
	}

	results := make([]Result, 0, len(diffs))

	for _, d := range diffs {
		r := Result{
			Path:   d.Path,
			Action: string(d.Type),
		}

		if dryRun {
			r.Success = true
			results = append(results, r)
			continue
		}

		switch d.Type {
		case diff.Added, diff.Modified:
			if err := writer.WriteSecret(d.Path, d.NewValue); err != nil {
				r.Success = false
				r.Err = fmt.Errorf("write %s: %w", d.Path, err)
			} else {
				r.Success = true
			}
		case diff.Removed:
			if err := writer.DeleteSecret(d.Path); err != nil {
				r.Success = false
				r.Err = fmt.Errorf("delete %s: %w", d.Path, err)
			} else {
				r.Success = true
			}
		}

		results = append(results, r)
	}

	return results, nil
}

// Summary returns a human-readable summary of patch results.
func Summary(results []Result) string {
	var succeeded, failed int
	for _, r := range results {
		if r.Success {
			succeeded++
		} else {
			failed++
		}
	}
	return fmt.Sprintf("patch complete: %d succeeded, %d failed", succeeded, failed)
}
