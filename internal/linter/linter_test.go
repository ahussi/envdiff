package linter_test

import (
	"testing"

	"github.com/user/envdiff/internal/linter"
	"github.com/user/envdiff/internal/parser"
)

func makeEntries(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestCheck_EmptyValue(t *testing.T) {
	entries := makeEntries("DB_HOST", "")
	findings := linter.Check(entries)
	if len(findings) == 0 {
		t.Fatal("expected a finding for empty value")
	}
	if findings[0].Severity != linter.Warn {
		t.Errorf("expected Warn, got %s", findings[0].Severity)
	}
}

func TestCheck_LowercaseKey(t *testing.T) {
	entries := makeEntries("db_host", "localhost")
	findings := linter.Check(entries)
	if len(findings) == 0 {
		t.Fatal("expected a finding for lowercase key")
	}
	if findings[0].Severity != linter.Warn {
		t.Errorf("expected Warn, got %s", findings[0].Severity)
	}
}

func TestCheck_WhitespaceValue(t *testing.T) {
	entries := makeEntries("DB_HOST", "   ")
	findings := linter.Check(entries)
	if len(findings) == 0 {
		t.Fatal("expected a finding for whitespace-only value")
	}
}

func TestCheck_SensitivePlaintext(t *testing.T) {
	entries := makeEntries("API_SECRET", "hunter2")
	findings := linter.Check(entries)
	if len(findings) == 0 {
		t.Fatal("expected a finding for plaintext secret")
	}
	if findings[0].Severity != linter.Error {
		t.Errorf("expected Error severity, got %s", findings[0].Severity)
	}
}

func TestCheck_SensitiveVariable_NoFinding(t *testing.T) {
	// Value is a variable reference — should not flag.
	entries := makeEntries("API_SECRET", "${SECRET_FROM_VAULT}")
	findings := linter.Check(entries)
	for _, f := range findings {
		if f.Key == "API_SECRET" && f.Severity == linter.Error {
			t.Errorf("did not expect an error finding for variable reference value")
		}
	}
}

func TestCheck_CleanEntry_NoFindings(t *testing.T) {
	entries := makeEntries("DB_HOST", "localhost")
	findings := linter.Check(entries)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d", len(findings))
	}
}

func TestCheck_MultipleIssues(t *testing.T) {
	entries := makeEntries(
		"DB_HOST", "",
		"api_key", "myplainkey",
	)
	findings := linter.Check(entries)
	// Expect at least: empty value warning + lowercase key warning + plaintext secret error
	if len(findings) < 3 {
		t.Errorf("expected at least 3 findings, got %d", len(findings))
	}
}
