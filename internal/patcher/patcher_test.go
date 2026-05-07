package patcher_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/patcher"
)

func makeResult(kind diff.Kind, key, a, b string) diff.Result {
	return diff.Result{Kind: kind, Key: key, ValueA: a, ValueB: b}
}

func TestParseFormat_Valid(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  patcher.Format
	}{
		{"env", patcher.FormatEnv},
		{"ENV", patcher.FormatEnv},
		{"export", patcher.FormatExport},
		{"Export", patcher.FormatExport},
	} {
		got, err := patcher.ParseFormat(tc.input)
		if err != nil {
			t.Fatalf("ParseFormat(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseFormat(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := patcher.ParseFormat("yaml")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestWrite_EnvFormat_MissingInB(t *testing.T) {
	results := []diff.Result{
		makeResult(diff.MissingInB, "DB_HOST", "localhost", ""),
	}
	var buf strings.Builder
	if err := patcher.Write(&buf, results, patcher.FormatEnv); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if got != "DB_HOST=localhost" {
		t.Errorf("got %q, want %q", got, "DB_HOST=localhost")
	}
}

func TestWrite_ExportFormat_Mismatch(t *testing.T) {
	results := []diff.Result{
		makeResult(diff.Mismatch, "APP_ENV", "production", "staging"),
	}
	var buf strings.Builder
	if err := patcher.Write(&buf, results, patcher.FormatExport); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if got != "export APP_ENV=production" {
		t.Errorf("got %q, want %q", got, "export APP_ENV=production")
	}
}

func TestWrite_SkipsMissingInA(t *testing.T) {
	results := []diff.Result{
		makeResult(diff.MissingInA, "ONLY_IN_B", "", "value"),
	}
	var buf strings.Builder
	if err := patcher.Write(&buf, results, patcher.FormatEnv); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}

func TestWrite_QuotesValueWithSpaces(t *testing.T) {
	results := []diff.Result{
		makeResult(diff.MissingInB, "GREETING", "hello world", ""),
	}
	var buf strings.Builder
	if err := patcher.Write(&buf, results, patcher.FormatEnv); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	want := `GREETING="hello world"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
