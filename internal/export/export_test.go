package export_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/your-org/vaultpatch/internal/export"
)

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []export.Format{export.FormatJSON, export.FormatYAML, export.FormatDotenv} {
		_, err := export.New(f)
		if err != nil {
			t.Errorf("expected no error for format %q, got %v", f, err)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := export.New("toml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestWrite_JSON(t *testing.T) {
	ex, _ := export.New(export.FormatJSON)
	secrets := map[string]string{"db_pass": "secret123", "api_key": "abc"}

	var buf bytes.Buffer
	if err := ex.Write(&buf, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]string
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if out["db_pass"] != "secret123" {
		t.Errorf("expected db_pass=secret123, got %q", out["db_pass"])
	}
}

func TestWrite_YAML(t *testing.T) {
	ex, _ := export.New(export.FormatYAML)
	secrets := map[string]string{"token": "mytoken"}

	var buf bytes.Buffer
	if err := ex.Write(&buf, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "token: mytoken") {
		t.Errorf("expected YAML to contain 'token: mytoken', got:\n%s", buf.String())
	}
}

func TestWrite_Dotenv(t *testing.T) {
	ex, _ := export.New(export.FormatDotenv)
	secrets := map[string]string{"db_host": "localhost", "db_port": "5432"}

	var buf bytes.Buffer
	if err := ex.Write(&buf, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output:\n%s", out)
	}
	if !strings.Contains(out, "DB_PORT=5432") {
		t.Errorf("expected DB_PORT=5432 in output:\n%s", out)
	}
}

func TestWrite_Dotenv_QuotesSpaces(t *testing.T) {
	ex, _ := export.New(export.FormatDotenv)
	secrets := map[string]string{"greeting": "hello world"}

	var buf bytes.Buffer
	if err := ex.Write(&buf, secrets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), `"hello world"`) {
		t.Errorf("expected quoted value in dotenv output:\n%s", buf.String())
	}
}
