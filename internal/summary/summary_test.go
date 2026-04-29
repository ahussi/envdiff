package summary_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/summary"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "A", Kind: diff.MissingInB},
		{Key: "B", Kind: diff.MissingInA},
		{Key: "C", Kind: diff.ValueMismatch, ValA: "x", ValB: "y"},
		{Key: "D", Kind: diff.MissingInB},
	}
}

func TestCompute_Counts(t *testing.T) {
	s := summary.Compute(makeResults())
	if s.Total != 4 {
		t.Errorf("Total: want 4, got %d", s.Total)
	}
	if s.Missing != 2 {
		t.Errorf("Missing: want 2, got %d", s.Missing)
	}
	if s.Extra != 1 {
		t.Errorf("Extra: want 1, got %d", s.Extra)
	}
	if s.Mismatch != 1 {
		t.Errorf("Mismatch: want 1, got %d", s.Mismatch)
	}
}

func TestCompute_Empty(t *testing.T) {
	s := summary.Compute(nil)
	if s.Total != 0 || s.Missing != 0 || s.Extra != 0 || s.Mismatch != 0 {
		t.Errorf("expected all zeros, got %+v", s)
	}
}

func TestWriteText_NoDiffs(t *testing.T) {
	var buf bytes.Buffer
	summary.WriteText(&buf, summary.Stats{})
	if !strings.Contains(buf.String(), "no differences") {
		t.Errorf("expected 'no differences' in output, got: %q", buf.String())
	}
}

func TestWriteText_WithDiffs(t *testing.T) {
	var buf bytes.Buffer
	s := summary.Compute(makeResults())
	summary.WriteText(&buf, s)
	out := buf.String()
	for _, want := range []string{"4", "2", "1", "missing", "extra", "mismatched"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output: %q", want, out)
		}
	}
}

func TestWriteJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	s := summary.Compute(makeResults())
	summary.WriteJSON(&buf, s)
	out := buf.String()
	for _, want := range []string{`"total":4`, `"missing":2`, `"extra":1`, `"mismatch":1`} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in JSON output: %q", want, out)
		}
	}
}
