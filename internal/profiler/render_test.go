package profiler_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/linter"
	"github.com/yourorg/envdiff/internal/profiler"
	"github.com/yourorg/envdiff/internal/scorer"
)

func makeProfile() profiler.Profile {
	return profiler.Profile{
		Label:        "staging",
		TotalKeys:    5,
		RedactedKeys: 2,
		Score:        scorer.Score{Value: 80, Grade: "B"},
		LintIssues: []linter.Issue{
			{Key: "db_host", Severity: "warn", Message: "key should be uppercase"},
		},
	}
}

func TestWriteText_ContainsLabel(t *testing.T) {
	var buf bytes.Buffer
	profiler.WriteText(&buf, makeProfile())
	if !strings.Contains(buf.String(), "staging") {
		t.Error("expected label 'staging' in output")
	}
}

func TestWriteText_ContainsScore(t *testing.T) {
	var buf bytes.Buffer
	profiler.WriteText(&buf, makeProfile())
	out := buf.String()
	if !strings.Contains(out, "80") || !strings.Contains(out, "B") {
		t.Errorf("expected score in output, got: %s", out)
	}
}

func TestWriteText_LintIssues(t *testing.T) {
	var buf bytes.Buffer
	profiler.WriteText(&buf, makeProfile())
	if !strings.Contains(buf.String(), "db_host") {
		t.Error("expected lint issue key 'db_host' in output")
	}
}

func TestWriteText_NoIssues(t *testing.T) {
	p := makeProfile()
	p.LintIssues = nil
	var buf bytes.Buffer
	profiler.WriteText(&buf, p)
	if !strings.Contains(buf.String(), "no issues") {
		t.Error("expected 'no issues' when lint slice is empty")
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := profiler.WriteJSON(&buf, makeProfile()); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["label"] != "staging" {
		t.Errorf("label = %v, want staging", out["label"])
	}
	if out["grade"] != "B" {
		t.Errorf("grade = %v, want B", out["grade"])
	}
}
