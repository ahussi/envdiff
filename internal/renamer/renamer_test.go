package renamer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/renamer"
)

func makeResult(kind diff.Kind, key, valA, valB string) diff.Result {
	return diff.Result{Kind: kind, Key: key, ValueA: valA, ValueB: valB}
}

func TestDetect_NoRenames(t *testing.T) {
	results := []diff.Result{
		makeResult(diff.MissingInB, "OLD_KEY", "alpha", ""),
		makeResult(diff.MissingInA, "NEW_KEY", "", "beta"),
	}
	got := renamer.Detect(results)
	if len(got) != 0 {
		t.Fatalf("expected 0 suggestions, got %d", len(got))
	}
}

func TestDetect_SingleRename(t *testing.T) {
	results := []diff.Result{
		makeResult(diff.MissingInB, "OLD_KEY", "shared-value", ""),
		makeResult(diff.MissingInA, "NEW_KEY", "", "shared-value"),
	}
	got := renamer.Detect(results)
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if got[0].OldKey != "OLD_KEY" {
		t.Errorf("OldKey: want OLD_KEY, got %s", got[0].OldKey)
	}
	if got[0].NewKey != "NEW_KEY" {
		t.Errorf("NewKey: want NEW_KEY, got %s", got[0].NewKey)
	}
	if got[0].Value != "shared-value" {
		t.Errorf("Value: want shared-value, got %s", got[0].Value)
	}
}

func TestDetect_IgnoresMismatch(t *testing.T) {
	results := []diff.Result{
		makeResult(diff.Mismatch, "SOME_KEY", "v1", "v2"),
	}
	got := renamer.Detect(results)
	if len(got) != 0 {
		t.Fatalf("expected 0 suggestions, got %d", len(got))
	}
}

func TestDetect_SkipsEmptyValues(t *testing.T) {
	// Empty values should not be matched as renames.
	results := []diff.Result{
		makeResult(diff.MissingInB, "OLD_KEY", "", ""),
		makeResult(diff.MissingInA, "NEW_KEY", "", ""),
	}
	got := renamer.Detect(results)
	if len(got) != 0 {
		t.Fatalf("expected 0 suggestions for empty values, got %d", len(got))
	}
}

func TestSuggestion_String(t *testing.T) {
	s := renamer.Suggestion{OldKey: "FOO", NewKey: "BAR", Value: "baz"}
	want := `FOO -> BAR (value: "baz")`
	if s.String() != want {
		t.Errorf("String(): want %q, got %q", want, s.String())
	}
}
