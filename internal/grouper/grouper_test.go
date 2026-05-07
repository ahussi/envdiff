package grouper_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/grouper"
)

func makeResults(keys ...string) []diff.Result {
	out := make([]diff.Result, len(keys))
	for i, k := range keys {
		out[i] = diff.Result{Key: k, Kind: diff.KindMissingInB}
	}
	return out
}

func TestAnalyse_GroupsByPrefix(t *testing.T) {
	results := makeResults("DB_HOST", "DB_PORT", "AWS_KEY", "PORT")
	groups := grouper.Analyse(results)

	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
	if groups[0].Prefix != "AWS" {
		t.Errorf("expected first group AWS, got %s", groups[0].Prefix)
	}
	if groups[1].Prefix != "DB" {
		t.Errorf("expected second group DB, got %s", groups[1].Prefix)
	}
	if groups[2].Prefix != "_ungrouped" {
		t.Errorf("expected last group _ungrouped, got %s", groups[2].Prefix)
	}
}

func TestAnalyse_UngroupedLast(t *testing.T) {
	results := makeResults("PLAIN", "AWS_SECRET", "GCP_PROJECT")
	groups := grouper.Analyse(results)

	last := groups[len(groups)-1]
	if last.Prefix != "_ungrouped" {
		t.Errorf("expected _ungrouped last, got %s", last.Prefix)
	}
	if len(last.Results) != 1 || last.Results[0].Key != "PLAIN" {
		t.Errorf("unexpected ungrouped results: %+v", last.Results)
	}
}

func TestAnalyse_EmptyInput(t *testing.T) {
	groups := grouper.Analyse(nil)
	if len(groups) != 0 {
		t.Errorf("expected 0 groups for empty input, got %d", len(groups))
	}
}

func TestAnalyse_AllSamePrefix(t *testing.T) {
	results := makeResults("DB_HOST", "DB_PORT", "DB_NAME")
	groups := grouper.Analyse(results)

	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Prefix != "DB" {
		t.Errorf("expected prefix DB, got %s", groups[0].Prefix)
	}
	if len(groups[0].Results) != 3 {
		t.Errorf("expected 3 results, got %d", len(groups[0].Results))
	}
}

func TestAnalyse_LeadingUnderscoreUngrouped(t *testing.T) {
	// A key starting with "_" has no prefix before the first underscore.
	results := makeResults("_INTERNAL")
	groups := grouper.Analyse(results)

	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Prefix != "_ungrouped" {
		t.Errorf("expected _ungrouped, got %s", groups[0].Prefix)
	}
}
