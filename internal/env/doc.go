// Package env provides environment-aware path resolution for Vault secrets.
//
// An Environment pairs a logical name (e.g. "dev", "staging", "prod") with
// a Vault path prefix (e.g. "secret/dev"). A Resolver built from a set of
// environments can translate a short secret key into a fully-qualified Vault
// path, making it easy to operate on the same logical secret across multiple
// environments without hard-coding paths throughout the codebase.
//
// Example:
//
//	r := env.NewResolver([]env.Environment{
//		{Name: "dev",  Prefix: "secret/dev"},
//		{Name: "prod", Prefix: "secret/prod"},
//	})
//
//	path, err := r.Resolve("prod", "database/password")
//	// path == "secret/prod/database/password"
package env
