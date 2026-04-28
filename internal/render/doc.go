// Package render formats diff and patch results for terminal output.
//
// It supports plain text and ANSI-colored output, and is used by the
// CLI layer to present changes to the operator before and after applying
// a patch to HashiCorp Vault.
//
// Example usage:
//
//	r := render.New(os.Stdout, true)
//	r.RenderDiff(diffs)
//	r.RenderSummary(added, removed, modified)
package render
