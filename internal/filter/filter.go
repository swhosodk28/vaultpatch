// Package filter provides key-based filtering for Vault secret diffs.
package filter

import (
	"path"
	"strings"

	"github.com/youorg/vaultpatch/internal/diff"
)

// Options controls which diff entries are included.
type Options struct {
	// Patterns is a list of glob patterns matched against secret keys.
	Patterns []string
	// Types restricts results to specific change types: "added", "removed", "modified".
	Types []string
}

// Apply returns a filtered copy of entries based on the given Options.
// An empty Patterns slice matches all keys. An empty Types slice matches all types.
func Apply(entries []diff.Diff, opts Options) []diff.Diff {
	var out []diff.Diff
	for _, e := range entries {
		if !matchesType(e, opts.Types) {
			continue
		}
		if !matchesPattern(e.Key, opts.Patterns) {
			continue
		}
		out = append(out, e)
	}
	return out
}

func matchesType(e diff.Diff, types []string) bool {
	if len(types) == 0 {
		return true
	}
	for _, t := range types {
		switch strings.ToLower(t) {
		case "added":
			if e.Old == "" && e.New != "" {
				return true
			}
		case "removed":
			if e.Old != "" && e.New == "" {
				return true
			}
		case "modified":
			if e.Old != "" && e.New != "" && e.Old != e.New {
				return true
			}
		}
	}
	return false
}

func matchesPattern(key string, patterns []string) bool {
	if len(patterns) == 0 {
		return true
	}
	for _, p := range patterns {
		matched, err := path.Match(p, key)
		if err == nil && matched {
			return true
		}
		// Also support substring match for convenience.
		if strings.Contains(key, p) {
			return true
		}
	}
	return false
}
