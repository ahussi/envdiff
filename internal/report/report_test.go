package report_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/report"
)

var sampleResults = []diff.Result{
	{Kind: diff.MissingInB, Key: "DB_HOST", ValueA: "localhost", ValueB: ""},
	{Kind: diff.MissingInA, Key: "API_KEY", ValueA: "", ValueB: "secret"},
	{Kind: diff.Mismatch, Key: "PORT", ValueA: "8080", ValueB: "9090"},
}

func TestWriteText_WithDiffs(t *testing.T) {
	var buf bytes.Buffer
	report.WriteText(&buf, sampleResults, ".env.dev", ".env.prod", false)
	out := buf.String()

	if !strings.Contains(out, "MISSING") {
		t.Error("expected MISSING label in output")
	}
	if !strings.Contains(out, "MISMATCH") {
		t.Error("expected MISMATCH label in output")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected key DB_HOST in output")
	}
}

func TestWriteText_NoDiffs(t *testing.T) {
	var buf bytes.Buffer
	report.WriteText(&buf, nil, ".env.dev", ".env.prod", false)
	out := buf.String()

	if !strings.Contains(out, "No differences") {
		t.Errorf("expected no-diff message, got: %s", out)
	}
}

func TestWriteJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	if err := report.WriteJSON(&buf, sampleResults, ".env.dev", ".env.prod"); err != nil {
		t.Fatalf("WriteJSON returned error: %v", err)
	}

	var out report.JSONReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}

	if out.Count != 3 {
		t.Errorf("expected 3 diffs, got %d", out.Count)
	}
	if out.FileA != ".env.dev" {
		t.Errorf("unexpected file_a: %s", out.FileA)
	}
	if out.Entries[0].Kind != "missing_in_b" {
		t.Errorf("unexpected kind: %s", out.Entries[0].Kind)
	}
}
