package differ_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envdiff/internal/differ"
)

func makeMap(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestAnalyse_IdenticalMaps(t *testing.T) {
	a := makeMap("FOO", "1", "BAR", "2")
	b := makeMap("FOO", "1", "BAR", "2")
	r := differ.Analyse(a, b)
	if r.Shared != 2 || r.TotalUnique != 2 {
		t.Fatalf("expected 2 shared / 2 total, got %d/%d", r.Shared, r.TotalUnique)
	}
	if r.SimilarityPct != 100.0 {
		t.Errorf("expected 100%% similarity, got %.1f", r.SimilarityPct)
	}
	if len(r.OnlyInA) != 0 || len(r.OnlyInB) != 0 {
		t.Errorf("expected no exclusive keys")
	}
}

func TestAnalyse_DisjointMaps(t *testing.T) {
	a := makeMap("FOO", "1")
	b := makeMap("BAR", "2")
	r := differ.Analyse(a, b)
	if r.Shared != 0 {
		t.Errorf("expected 0 shared, got %d", r.Shared)
	}
	if r.TotalUnique != 2 {
		t.Errorf("expected 2 total unique, got %d", r.TotalUnique)
	}
	if r.SimilarityPct != 0.0 {
		t.Errorf("expected 0%% similarity, got %.1f", r.SimilarityPct)
	}
}

func TestAnalyse_PartialOverlap(t *testing.T) {
	a := makeMap("FOO", "1", "BAR", "2", "BAZ", "3")
	b := makeMap("FOO", "x", "QUX", "4")
	r := differ.Analyse(a, b)
	// union: FOO BAR BAZ QUX = 4
	if r.TotalUnique != 4 {
		t.Errorf("expected 4 total unique, got %d", r.TotalUnique)
	}
	if r.Shared != 1 {
		t.Errorf("expected 1 shared key, got %d", r.Shared)
	}
	expectedPct := 25.0
	if r.SimilarityPct != expectedPct {
		t.Errorf("expected %.1f%% similarity, got %.1f", expectedPct, r.SimilarityPct)
	}
}

func TestAnalyse_EmptyMaps(t *testing.T) {
	r := differ.Analyse(map[string]string{}, map[string]string{})
	if r.TotalUnique != 0 || r.SimilarityPct != 0 {
		t.Errorf("expected zero result for empty maps")
	}
}

func TestSummary_ContainsPercentage(t *testing.T) {
	a := makeMap("FOO", "1", "BAR", "2")
	b := makeMap("FOO", "1")
	r := differ.Analyse(a, b)
	s := differ.Summary(r)
	if !strings.Contains(s, "similarity:") {
		t.Errorf("summary missing 'similarity:' prefix: %q", s)
	}
	if !strings.Contains(s, "only-in-A") {
		t.Errorf("summary missing only-in-A info: %q", s)
	}
}
