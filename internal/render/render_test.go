package render_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vaultpatch/vaultpatch/internal/diff"
	"github.com/vaultpatch/vaultpatch/internal/render"
)

func TestRenderDiff_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	r := render.New(&buf, false)
	r.RenderDiff(nil)
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected no-changes message, got: %q", buf.String())
	}
}

func TestRenderDiff_Added(t *testing.T) {
	var buf bytes.Buffer
	r := render.New(&buf, false)
	r.RenderDiff([]diff.Diff{
		{Op: diff.OpAdd, Path: "secret/app", Key: "token", NewValue: "abc123"},
	})
	out := buf.String()
	if !strings.Contains(out, "+ [secret/app] token = abc123") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRenderDiff_Removed(t *testing.T) {
	var buf bytes.Buffer
	r := render.New(&buf, false)
	r.RenderDiff([]diff.Diff{
		{Op: diff.OpRemove, Path: "secret/app", Key: "old_key", OldValue: "stale"},
	})
	out := buf.String()
	if !strings.Contains(out, "- [secret/app] old_key = stale") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRenderDiff_Modified(t *testing.T) {
	var buf bytes.Buffer
	r := render.New(&buf, false)
	r.RenderDiff([]diff.Diff{
		{Op: diff.OpModify, Path: "secret/db", Key: "password", OldValue: "old", NewValue: "new"},
	})
	out := buf.String()
	if !strings.Contains(out, "~ [secret/db] password: old -> new") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRenderDiff_ColorEnabled(t *testing.T) {
	var buf bytes.Buffer
	r := render.New(&buf, true)
	r.RenderDiff([]diff.Diff{
		{Op: diff.OpAdd, Path: "secret/app", Key: "k", NewValue: "v"},
	})
	out := buf.String()
	if !strings.Contains(out, "\033[") {
		t.Errorf("expected ANSI codes in output, got: %q", out)
	}
}

func TestRenderSummary_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	r := render.New(&buf, false)
	r.RenderSummary(0, 0, 0)
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected no-changes summary, got: %q", buf.String())
	}
}

func TestRenderSummary_Mixed(t *testing.T) {
	var buf bytes.Buffer
	r := render.New(&buf, false)
	r.RenderSummary(2, 1, 3)
	out := buf.String()
	for _, want := range []string{"2 added", "1 removed", "3 modified"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in summary: %q", want, out)
		}
	}
}
