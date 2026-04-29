package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpatch/internal/snapshot"
)

func makeSnap(env string) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		Environment: env,
		Path:        "secret/" + env,
		CapturedAt:  time.Now().UTC(),
		Secrets:     map[string]string{"token": "abc123"},
	}
}

func TestNewStore_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "snapshots")
	_, err := snapshot.NewStore(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestStore_SaveAndLoad(t *testing.T) {
	store, err := snapshot.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("store init failed: %v", err)
	}

	snap := makeSnap("prod")
	path, err := store.Save(snap)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := store.Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.Environment != snap.Environment {
		t.Errorf("environment mismatch: got %q want %q", loaded.Environment, snap.Environment)
	}
	if loaded.Secrets["token"] != "abc123" {
		t.Errorf("secret mismatch: got %q", loaded.Secrets["token"])
	}
}

func TestStore_LatestFor_ReturnsFile(t *testing.T) {
	store, err := snapshot.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("store init failed: %v", err)
	}

	_, err = store.Save(makeSnap("staging"))
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	path, err := store.LatestFor("staging")
	if err != nil {
		t.Fatalf("LatestFor failed: %v", err)
	}
	if path == "" {
		t.Error("expected non-empty path")
	}
}

func TestStore_LatestFor_NoSnapshots(t *testing.T) {
	store, err := snapshot.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("store init failed: %v", err)
	}

	_, err = store.LatestFor("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing environment snapshots")
	}
}

func TestStore_Load_InvalidPath(t *testing.T) {
	store, err := snapshot.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("store init failed: %v", err)
	}

	_, err = store.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
