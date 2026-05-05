package ignorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/ignorer"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "DATABASE_URL", Kind: diff.MissingInB},
		{Key: "AWS_SECRET_KEY", Kind: diff.Mismatch},
		{Key: "APP_PORT", Kind: diff.MissingInA},
		{Key: "AWS_ACCESS_KEY_ID", Kind: diff.Mismatch},
		{Key: "LOG_LEVEL", Kind: diff.MissingInB},
	}
}

func TestApply_NoPatterns(t *testing.T) {
	results := makeResults()
	out, err := ignorer.Apply(results, ignorer.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(out))
	}
}

func TestApply_ExactMatch(t *testing.T) {
	out, err := ignorer.Apply(makeResults(), ignorer.Options{
		Patterns: []string{"DATABASE_URL"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range out {
		if r.Key == "DATABASE_URL" {
			t.Errorf("DATABASE_URL should have been ignored")
		}
	}
	if len(out) != 4 {
		t.Errorf("expected 4 results, got %d", len(out))
	}
}

func TestApply_WildcardPrefix(t *testing.T) {
	out, err := ignorer.Apply(makeResults(), ignorer.Options{
		Patterns: []string{"AWS_*"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range out {
		if r.Key == "AWS_SECRET_KEY" || r.Key == "AWS_ACCESS_KEY_ID" {
			t.Errorf("key %q should have been ignored", r.Key)
		}
	}
	if len(out) != 3 {
		t.Errorf("expected 3 results, got %d", len(out))
	}
}

func TestApply_MultiplePatterns(t *testing.T) {
	out, err := ignorer.Apply(makeResults(), ignorer.Options{
		Patterns: []string{"AWS_*", "LOG_LEVEL"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 results, got %d", len(out))
	}
}

func TestApply_CaseInsensitive(t *testing.T) {
	out, err := ignorer.Apply(makeResults(), ignorer.Options{
		Patterns: []string{"aws_*"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range out {
		if r.Key == "AWS_SECRET_KEY" || r.Key == "AWS_ACCESS_KEY_ID" {
			t.Errorf("key %q should have been ignored via case-insensitive match", r.Key)
		}
	}
}

func TestApply_InvalidPattern(t *testing.T) {
	_, err := ignorer.Apply(makeResults(), ignorer.Options{
		Patterns: []string{"[invalid"},
	})
	if err == nil {
		t.Error("expected error for invalid glob pattern, got nil")
	}
}
