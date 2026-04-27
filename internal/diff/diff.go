package diff

import "fmt"

// ChangeType represents the type of change for a secret key.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Change represents a single key-level difference between two secret maps.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Diff computes the diff between a source and target secret map.
// Values are compared as strings. Both maps must be non-nil.
func Diff(source, target map[string]string) []Change {
	changes := []Change{}

	for key, srcVal := range source {
		if tgtVal, ok := target[key]; ok {
			if srcVal != tgtVal {
				changes = append(changes, Change{
					Key:      key,
					Type:     Modified,
					OldValue: srcVal,
					NewValue: tgtVal,
				})
			} else {
				changes = append(changes, Change{
					Key:      key,
					Type:     Unchanged,
					OldValue: srcVal,
					NewValue: tgtVal,
				})
			}
		} else {
			changes = append(changes, Change{
				Key:      key,
				Type:     Removed,
				OldValue: srcVal,
				NewValue: "",
			})
		}
	}

	for key, tgtVal := range target {
		if _, ok := source[key]; !ok {
			changes = append(changes, Change{
				Key:      key,
				Type:     Added,
				OldValue: "",
				NewValue: tgtVal,
			})
		}
	}

	return changes
}

// Summary returns a human-readable summary line for a Change.
func Summary(c Change) string {
	switch c.Type {
	case Added:
		return fmt.Sprintf("+ %s = %q", c.Key, c.NewValue)
	case Removed:
		return fmt.Sprintf("- %s (was %q)", c.Key, c.OldValue)
	case Modified:
		return fmt.Sprintf("~ %s: %q -> %q", c.Key, c.OldValue, c.NewValue)
	default:
		return fmt.Sprintf("  %s = %q", c.Key, c.NewValue)
	}
}
