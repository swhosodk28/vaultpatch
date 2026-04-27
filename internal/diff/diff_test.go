package diff

import (
	"testing"
)

func TestDiff_NoChanges(t *testing.T) {
	src := map[string]string{"key1": "val1", "key2": "val2"}
	tgt := map[string]string{"key1": "val1", "key2": "val2"}

	changes := Diff(src, tgt)
	for _, c := range changes {
		if c.Type != Unchanged {
			t.Errorf("expected Unchanged for key %q, got %s", c.Key, c.Type)
		}
	}
}

func TestDiff_Added(t *testing.T) {
	src := map[string]string{}
	tgt := map[string]string{"new_key": "new_val"}

	changes := Diff(src, tgt)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != Added {
		t.Errorf("expected Added, got %s", changes[0].Type)
	}
	if changes[0].Key != "new_key" {
		t.Errorf("expected key 'new_key', got %q", changes[0].Key)
	}
}

func TestDiff_Removed(t *testing.T) {
	src := map[string]string{"old_key": "old_val"}
	tgt := map[string]string{}

	changes := Diff(src, tgt)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != Removed {
		t.Errorf("expected Removed, got %s", changes[0].Type)
	}
	if changes[0].OldValue != "old_val" {
		t.Errorf("expected OldValue 'old_val', got %q", changes[0].OldValue)
	}
}

func TestDiff_Modified(t *testing.T) {
	src := map[string]string{"key": "v1"}
	tgt := map[string]string{"key": "v2"}

	changes := Diff(src, tgt)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != Modified {
		t.Errorf("expected Modified, got %s", changes[0].Type)
	}
	if changes[0].OldValue != "v1" || changes[0].NewValue != "v2" {
		t.Errorf("unexpected values: old=%q new=%q", changes[0].OldValue, changes[0].NewValue)
	}
}

func TestSummary(t *testing.T) {
	cases := []struct {
		change   Change
		expected string
	}{
		{Change{Key: "k", Type: Added, NewValue: "v"}, `+ k = "v"`},
		{Change{Key: "k", Type: Removed, OldValue: "v"}, `- k (was "v")`},
		{Change{Key: "k", Type: Modified, OldValue: "a", NewValue: "b"}, `~ k: "a" -> "b"`},
		{Change{Key: "k", Type: Unchanged, NewValue: "v"}, `  k = "v"`},
	}

	for _, tc := range cases {
		got := Summary(tc.change)
		if got != tc.expected {
			t.Errorf("Summary(%+v) = %q, want %q", tc.change, got, tc.expected)
		}
	}
}
