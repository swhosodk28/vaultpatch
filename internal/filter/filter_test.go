package filter_test

import (
	"testing"

	"github.com/youorg/vaultpatch/internal/diff"
	"github.com/youorg/vaultpatch/internal/filter"
)

func entries() []diff.Diff {
	return []diff.Diff{
		{Key: "db/password", Old: "", New: "secret"},
		{Key: "db/user", Old: "admin", New: ""},
		{Key: "api/key", Old: "old", New: "new"},
		{Key: "api/timeout", Old: "30", New: "30"},
	}
}

func TestApply_NoOptions(t *testing.T) {
	got := filter.Apply(entries(), filter.Options{})
	if len(got) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(got))
	}
}

func TestApply_FilterByTypeAdded(t *testing.T) {
	got := filter.Apply(entries(), filter.Options{Types: []string{"added"}})
	if len(got) != 1 || got[0].Key != "db/password" {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestApply_FilterByTypeRemoved(t *testing.T) {
	got := filter.Apply(entries(), filter.Options{Types: []string{"removed"}})
	if len(got) != 1 || got[0].Key != "db/user" {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestApply_FilterByTypeModified(t *testing.T) {
	got := filter.Apply(entries(), filter.Options{Types: []string{"modified"}})
	if len(got) != 1 || got[0].Key != "api/key" {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestApply_FilterByPattern(t *testing.T) {
	got := filter.Apply(entries(), filter.Options{Patterns: []string{"api/*"}})
	if len(got) != 2 {
		t.Fatalf("expected 2 api entries, got %d", len(got))
	}
}

func TestApply_FilterByPatternAndType(t *testing.T) {
	got := filter.Apply(entries(), filter.Options{
		Patterns: []string{"api/*"},
		Types:    []string{"modified"},
	})
	if len(got) != 1 || got[0].Key != "api/key" {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestApply_NoMatch(t *testing.T) {
	got := filter.Apply(entries(), filter.Options{Patterns: []string{"nonexistent/*"}})
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestApply_SubstringPattern(t *testing.T) {
	got := filter.Apply(entries(), filter.Options{Patterns: []string{"password"}})
	if len(got) != 1 || got[0].Key != "db/password" {
		t.Fatalf("unexpected result: %+v", got)
	}
}
