package filter_test

import (
	"testing"

	"github.com/yourusername/envdiff/internal/diff"
	"github.com/yourusername/envdiff/internal/filter"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Kind: diff.MissingInB},
		{Key: "DB_PORT", Kind: diff.Mismatch, ValA: "5432", ValB: "3306"},
		{Key: "APP_SECRET", Kind: diff.MissingInA},
		{Key: "APP_DEBUG", Kind: diff.Mismatch, ValA: "true", ValB: "false"},
	}
}

func TestApply_NoOptions(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{})
	if len(got) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(got))
	}
}

func TestApply_FilterByKind(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{OnlyKinds: []string{"mismatch"}})
	if len(got) != 2 {
		t.Fatalf("expected 2 mismatch results, got %d", len(got))
	}
	for _, r := range got {
		if r.Kind != diff.Mismatch {
			t.Errorf("expected Mismatch kind, got %v", r.Kind)
		}
	}
}

func TestApply_FilterByPrefix(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{KeyPrefix: "APP_"})
	if len(got) != 2 {
		t.Fatalf("expected 2 APP_ results, got %d", len(got))
	}
	for _, r := range got {
		if !startsWith(r.Key, "APP_") {
			t.Errorf("unexpected key %q", r.Key)
		}
	}
}

func TestApply_FilterByKindAndPrefix(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{
		OnlyKinds: []string{"mismatch"},
		KeyPrefix:  "APP_",
	})
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].Key != "APP_DEBUG" {
		t.Errorf("expected APP_DEBUG, got %q", got[0].Key)
	}
}

func TestApply_NoMatches(t *testing.T) {
	results := makeResults()
	got := filter.Apply(results, filter.Options{KeyPrefix: "REDIS_"})
	if len(got) != 0 {
		t.Errorf("expected 0 results, got %d", len(got))
	}
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
