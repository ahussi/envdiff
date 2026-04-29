package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
	"github.com/yourorg/envdiff/internal/report"
)

func TestWriteMarkdown_NoDiffs(t *testing.T) {
	var buf bytes.Buffer
	err := report.WriteMarkdown(&buf, []diff.Result{}, "a.env", "b.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "No differences found") {
		t.Errorf("expected no-diff message, got:\n%s", out)
	}
}

func TestWriteMarkdown_WithDiffs(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", Kind: diff.Mismatch, ValueA: "localhost", ValueB: "prod.db"},
		{Key: "API_KEY", Kind: diff.MissingInB, ValueA: "abc123", ValueB: ""},
	}
	var buf bytes.Buffer
	err := report.WriteMarkdown(&buf, results, "dev.env", "prod.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()

	if !strings.Contains(out, "# envdiff report") {
		t.Errorf("expected markdown heading, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got:\n%s", out)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output, got:\n%s", out)
	}
	if !strings.Contains(out, "2 issue(s) found") {
		t.Errorf("expected issue count, got:\n%s", out)
	}
	if !strings.Contains(out, "| Key | Kind") {
		t.Errorf("expected markdown table header, got:\n%s", out)
	}
}

func TestWriteMarkdown_EmptyValues(t *testing.T) {
	results := []diff.Result{
		{Key: "EMPTY_KEY", Kind: diff.MissingInA, ValueA: "", ValueB: "somevalue"},
	}
	var buf bytes.Buffer
	_ = report.WriteMarkdown(&buf, results, "a.env", "b.env")
	out := buf.String()
	if !strings.Contains(out, "_(empty)_") {
		t.Errorf("expected empty placeholder, got:\n%s", out)
	}
}
