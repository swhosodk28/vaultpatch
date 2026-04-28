package env_test

import (
	"sort"
	"testing"

	"github.com/your-org/vaultpatch/internal/env"
)

func makeResolver() *env.Resolver {
	return env.NewResolver([]env.Environment{
		{Name: "dev", Prefix: "secret/dev"},
		{Name: "staging", Prefix: "secret/staging/"},
		{Name: "prod", Prefix: "secret/prod"},
	})
}

func TestResolve_KnownEnvironment(t *testing.T) {
	r := makeResolver()
	path, err := r.Resolve("dev", "database/password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "secret/dev/database/password"
	if path != want {
		t.Errorf("got %q, want %q", path, want)
	}
}

func TestResolve_TrimsSlashes(t *testing.T) {
	r := makeResolver()
	path, err := r.Resolve("staging", "/api/key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "secret/staging/api/key"
	if path != want {
		t.Errorf("got %q, want %q", path, want)
	}
}

func TestResolve_CaseInsensitive(t *testing.T) {
	r := makeResolver()
	path, err := r.Resolve("PROD", "token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "secret/prod/token"
	if path != want {
		t.Errorf("got %q, want %q", path, want)
	}
}

func TestResolve_UnknownEnvironment(t *testing.T) {
	r := makeResolver()
	_, err := r.Resolve("qa", "some/key")
	if err == nil {
		t.Fatal("expected error for unknown environment, got nil")
	}
}

func TestList_ReturnsAllNames(t *testing.T) {
	r := makeResolver()
	names := r.List()
	sort.Strings(names)
	want := []string{"dev", "prod", "staging"}
	if len(names) != len(want) {
		t.Fatalf("got %v, want %v", names, want)
	}
	for i, n := range names {
		if n != want[i] {
			t.Errorf("names[%d] = %q, want %q", i, n, want[i])
		}
	}
}

func TestHas_KnownAndUnknown(t *testing.T) {
	r := makeResolver()
	if !r.Has("dev") {
		t.Error("expected Has(\"dev\") to be true")
	}
	if r.Has("qa") {
		t.Error("expected Has(\"qa\") to be false")
	}
}
