// Package filter provides utilities for narrowing diff results by key pattern
// and change type before rendering or applying them.
//
// Usage:
//
//	filtered := filter.Apply(diffs, filter.Options{
//		Patterns: []string{"db/*"},
//		Types:    []string{"added", "modified"},
//	})
//
// Patterns support standard Go path.Match glob syntax as well as simple
// substring matching for convenience. Types must be one of "added",
// "removed", or "modified" (case-insensitive).
package filter
