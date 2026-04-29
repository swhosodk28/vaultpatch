package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Store handles persisting and loading snapshots to/from the filesystem.
type Store struct {
	dir string
}

// NewStore creates a Store that saves snapshots under the given directory.
// The directory is created if it does not exist.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, fmt.Errorf("snapshot store: cannot create directory %q: %w", dir, err)
	}
	return &Store{dir: dir}, nil
}

// Save writes the snapshot to a timestamped file in the store directory.
// Returns the path of the written file.
func (s *Store) Save(snap *Snapshot) (string, error) {
	data, err := snap.Marshal()
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s-%s.json",
		snap.Environment,
		snap.CapturedAt.Format("20060102T150405Z"),
	)
	path := filepath.Join(s.dir, filename)

	if err := os.WriteFile(path, data, 0o640); err != nil {
		return "", fmt.Errorf("snapshot store: write failed: %w", err)
	}
	return path, nil
}

// Load reads and deserialises a snapshot from the given file path.
func (s *Store) Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot store: read failed: %w", err)
	}
	return Unmarshal(data)
}

// LatestFor returns the most recently saved snapshot file path for the
// given environment, or an error if none is found.
func (s *Store) LatestFor(env string) (string, error) {
	pattern := filepath.Join(s.dir, env+"-*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("snapshot store: glob error: %w", err)
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("snapshot store: no snapshots found for environment %q", env)
	}

	// Filenames are timestamped; lexicographic sort gives us the latest.
	latest := matches[0]
	var latestTime time.Time
	for _, m := range matches {
		info, err := os.Stat(m)
		if err != nil {
			continue
		}
		if info.ModTime().After(latestTime) {
			latestTime = info.ModTime()
			latest = m
		}
	}
	return latest, nil
}
