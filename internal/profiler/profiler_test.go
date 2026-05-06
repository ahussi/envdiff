package profiler_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envdiff/internal/profiler"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestAnalyse_BasicFile(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=envdiff\nDEBUG=true\n")
	p, err := profiler.Analyse(path, "test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.TotalKeys != 2 {
		t.Errorf("TotalKeys = %d, want 2", p.TotalKeys)
	}
	if p.Label != "test" {
		t.Errorf("Label = %q, want \"test\"", p.Label)
	}
}

func TestAnalyse_LintIssuesDetected(t *testing.T) {
	// lowercase key triggers linter
	path := writeTempEnv(t, "app_name=envdiff\n")
	p, err := profiler.Analyse(path, "lint-test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.LintIssues) == 0 {
		t.Error("expected at least one lint issue for lowercase key")
	}
}

func TestAnalyse_RedactedKeysCount(t *testing.T) {
	path := writeTempEnv(t, "DB_PASSWORD=secret123\nAPP_NAME=envdiff\n")
	p, err := profiler.Analyse(path, "redact-test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.RedactedKeys == 0 {
		t.Error("expected DB_PASSWORD to be counted as redacted")
	}
}

func TestAnalyse_FileNotFound(t *testing.T) {
	_, err := profiler.Analyse("/nonexistent/.env", "missing", nil)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestAnalyse_EmptyFile(t *testing.T) {
	path := writeTempEnv(t, "")
	p, err := profiler.Analyse(path, "empty", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.TotalKeys != 0 {
		t.Errorf("TotalKeys = %d, want 0", p.TotalKeys)
	}
}
