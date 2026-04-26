package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return p
}

func TestParse_BasicKeyValue(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=envdiff\nPORT=8080\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Entries["APP_NAME"].Value != "envdiff" {
		t.Errorf("expected APP_NAME=envdiff, got %q", env.Entries["APP_NAME"].Value)
	}
	if env.Entries["PORT"].Value != "8080" {
		t.Errorf("expected PORT=8080, got %q", env.Entries["PORT"].Value)
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `DB_URL="postgres://localhost/mydb"` + "\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Entries["DB_URL"].Value != "postgres://localhost/mydb" {
		t.Errorf("unexpected value: %q", env.Entries["DB_URL"].Value)
	}
}

func TestParse_Comments(t *testing.T) {
	path := writeTempEnv(t, "# database url\nDB_URL=localhost\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.Entries["DB_URL"].Comment != "database url" {
		t.Errorf("expected comment 'database url', got %q", env.Entries["DB_URL"].Comment)
	}
}

func TestParse_SkipsBlankLines(t *testing.T) {
	path := writeTempEnv(t, "\nKEY=value\n\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(env.Entries))
	}
}

func TestParse_FileNotFound(t *testing.T) {
	_, err := Parse("/nonexistent/.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
