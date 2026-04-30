// Package rollback implements point-in-time restoration of Vault secrets
// using previously captured snapshots.
//
// # Overview
//
// Given a [snapshot.Snapshot] and the current state of secrets in Vault,
// Apply computes the delta and writes or deletes secrets so that Vault
// matches the snapshot exactly.
//
// # Dry-run mode
//
// Pass Options{DryRun: true} to preview which paths would be restored or
// deleted without performing any actual Vault writes.
//
// # Error handling
//
// Individual write/delete failures are collected in Result.Errors rather
// than aborting the entire rollback, allowing partial progress to be
// inspected and retried.
package rollback
