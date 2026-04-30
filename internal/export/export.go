// Package export provides functionality to export Vault secrets
// to common formats such as JSON, YAML, and dotenv.
package export

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Format represents an output format for exported secrets.
type Format string

const (
	FormatJSON   Format = "json"
	FormatYAML   Format = "yaml"
	FormatDotenv Format = "dotenv"
)

// Exporter writes secrets in a given format to a writer.
type Exporter struct {
	format Format
}

// New returns a new Exporter for the given format.
func New(format Format) (*Exporter, error) {
	switch format {
	case FormatJSON, FormatYAML, FormatDotenv:
		return &Exporter{format: format}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %q", format)
	}
}

// Write serializes the secrets map to the given writer.
func (e *Exporter) Write(w io.Writer, secrets map[string]string) error {
	switch e.format {
	case FormatJSON:
		return writeJSON(w, secrets)
	case FormatYAML:
		return writeYAML(w, secrets)
	case FormatDotenv:
		return writeDotenv(w, secrets)
	default:
		return fmt.Errorf("unsupported format: %q", e.format)
	}
}

func writeJSON(w io.Writer, secrets map[string]string) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(secrets)
}

func writeYAML(w io.Writer, secrets map[string]string) error {
	return yaml.NewEncoder(w).Encode(secrets)
}

func writeDotenv(w io.Writer, secrets map[string]string) error {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		v := secrets[k]
		// Quote values that contain spaces or special characters.
		if strings.ContainsAny(v, " \t\n\r") {
			v = fmt.Sprintf("%q", v)
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", strings.ToUpper(k), v))
	}
	_, err := io.WriteString(w, sb.String())
	return err
}
