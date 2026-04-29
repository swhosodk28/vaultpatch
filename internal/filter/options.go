package filter

import (
	"fmt"
	"strings"
)

// ValidTypes lists the accepted change type strings.
var ValidTypes = []string{"added", "removed", "modified"}

// Validate checks that all Types in opts are recognised values.
// It returns an error listing any unrecognised types.
func (opts Options) Validate() error {
	valid := map[string]bool{}
	for _, t := range ValidTypes {
		valid[t] = true
	}
	var bad []string
	for _, t := range opts.Types {
		if !valid[strings.ToLower(t)] {
			bad = append(bad, t)
		}
	}
	if len(bad) > 0 {
		return fmt.Errorf("filter: unknown type(s): %s (valid: %s)",
			strings.Join(bad, ", "),
			strings.Join(ValidTypes, ", "),
		)
	}
	return nil
}

// FromFlags is a convenience constructor that builds Options from raw CLI
// flag values, trimming whitespace from each element.
func FromFlags(patterns, types []string) Options {
	trimAll := func(ss []string) []string {
		out := make([]string, 0, len(ss))
		for _, s := range ss {
			if v := strings.TrimSpace(s); v != "" {
				out = append(out, v)
			}
		}
		return out
	}
	return Options{
		Patterns: trimAll(patterns),
		Types:    trimAll(types),
	}
}

// IsEmpty reports whether the options contain no patterns and no type filters,
// meaning no filtering will be applied.
func (opts Options) IsEmpty() bool {
	return len(opts.Patterns) == 0 && len(opts.Types) == 0
}
