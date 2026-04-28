// Package render provides formatted output for diff and patch results.
package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/vaultpatch/vaultpatch/internal/diff"
)

// ColorMode controls whether ANSI color codes are emitted.
type ColorMode int

const (
	ColorAuto ColorMode = iota
	ColorAlways
	ColorNever
)

const (
	ansiReset  = "\033[0m"
	ansiRed    = "\033[31m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiBold   = "\033[1m"
)

// Renderer writes human-readable output to a writer.
type Renderer struct {
	w     io.Writer
	color bool
}

// New creates a Renderer writing to w. color enables ANSI escape codes.
func New(w io.Writer, color bool) *Renderer {
	return &Renderer{w: w, color: color}
}

// RenderDiff prints a unified-style diff for a slice of diff.Diff entries.
func (r *Renderer) RenderDiff(diffs []diff.Diff) {
	if len(diffs) == 0 {
		fmt.Fprintln(r.w, "No changes detected.")
		return
	}
	for _, d := range diffs {
		switch d.Op {
		case diff.OpAdd:
			fmt.Fprintln(r.w, r.colorize(ansiGreen, fmt.Sprintf("+ [%s] %s = %s", d.Path, d.Key, d.NewValue)))
		case diff.OpRemove:
			fmt.Fprintln(r.w, r.colorize(ansiRed, fmt.Sprintf("- [%s] %s = %s", d.Path, d.Key, d.OldValue)))
		case diff.OpModify:
			fmt.Fprintln(r.w, r.colorize(ansiYellow, fmt.Sprintf("~ [%s] %s: %s -> %s", d.Path, d.Key, d.OldValue, d.NewValue)))
		}
	}
}

// RenderSummary prints a one-line summary of diff counts.
func (r *Renderer) RenderSummary(added, removed, modified int) {
	parts := []string{}
	if added > 0 {
		parts = append(parts, r.colorize(ansiGreen, fmt.Sprintf("%d added", added)))
	}
	if removed > 0 {
		parts = append(parts, r.colorize(ansiRed, fmt.Sprintf("%d removed", removed)))
	}
	if modified > 0 {
		parts = append(parts, r.colorize(ansiYellow, fmt.Sprintf("%d modified", modified)))
	}
	if len(parts) == 0 {
		fmt.Fprintln(r.w, "Summary: no changes.")
		return
	}
	fmt.Fprintf(r.w, "%sSummary:%s %s\n", ansiBold, ansiReset, strings.Join(parts, ", "))
}

func (r *Renderer) colorize(code, s string) string {
	if !r.color {
		return s
	}
	return code + s + ansiReset
}
