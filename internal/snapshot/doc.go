// Package snapshot provides point-in-time capture and persistence of
// HashiCorp Vault secret state for a named environment.
//
// A Snapshot records the full set of key/value secrets at a given Vault
// path along with the environment name and capture timestamp. Snapshots
// can be serialised to JSON and stored on disk via a Store.
//
// Typical usage:
//
//	// Capture current state
//	snap, err := snapshot.Capture("prod", "secret/prod", vaultClient)
//
//	// Persist to disk
//	store, _ := snapshot.NewStore(".vaultpatch/snapshots")
//	path, _ := store.Save(snap)
//
//	// Reload later for diffing
//	old, _ := store.Load(path)
package snapshot
