// Package export serialises Vault secret maps into portable output formats.
//
// Supported formats:
//
//   - json    — pretty-printed JSON object
//   - yaml    — YAML mapping
//   - dotenv  — KEY=VALUE pairs suitable for use with direnv or Docker
//
// Usage:
//
//	ex, err := export.New(export.FormatDotenv)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := ex.Write(os.Stdout, secrets); err != nil {
//		log.Fatal(err)
//	}
package export
