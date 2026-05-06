package scorer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/scorer"
)

func makeResults(kinds ...diff.Kind) []diff.Result {
	out := make([]diff.Result, len(kinds))
	for i, k := range kinds {
		out[i] = diff.Result{Kind: k, Key: fmt.Sprintf("KEY_%d", i)}
	}
	return out
}

func TestCompute_PerfectScore(t *testing.T) {
	r := scorer.Compute(nil, 10)
	if r.Score != 100 {
		t.Errorf("expected 100, got %.1f", r.Score)
	}
	if r.Grade != "A" {
		t.Errorf("expected grade A, got %s", r.Grade)
	}
}

func TestCompute_ZeroTotal(t *testing.T) {
	r := scorer.Compute(nil, 0)
	if r.Score != 100 || r.Grade != "A" {
		t.Errorf("expected perfect score for zero total, got %s", r)
	}
}

func TestCompute_AllMissing(t *testing.T) {
	results := []diff.Result{
		{Kind: diff.MissingInB, Key: "A"},
		{Kind: diff.MissingInB, Key: "B"},
	}
	r := scorer.Compute(results, 2)
	if r.Score != 0 {
		t.Errorf("expected 0, got %.1f", r.Score)
	}
	if r.Grade != "F" {
		t.Errorf("expected grade F, got %s", r.Grade)
	}
}

func TestCompute_MismatchPenalty(t *testing.T) {
	results := []diff.Result{
		{Kind: diff.Mismatch, Key: "X"},
	}
	r := scorer.Compute(results, 4)
	// penalty=1.5, maxPenalty=12 => score=100*(1-1.5/12)=87.5
	if r.Score != 87.5 {
		t.Errorf("expected 87.5, got %.1f", r.Score)
	}
	if r.Grade != "B" {
		t.Errorf("expected grade B, got %s", r.Grade)
	}
}

func TestCompute_GradeThresholds(t *testing.T) {
	tests := []struct {
		score    float64
		wantGrade string
	}{
		{95, "A"}, {89, "B"}, {74, "C"}, {59, "D"}, {39, "F"},
	}
	for _, tt := range tests {
		// drive score via penalty: score = 100*(1-p/max), so p = (1-score/100)*max
		total := 100
		maxP := float64(total) * 3.0
		penalty := (1.0 - tt.score/100.0) * maxP
		nMissing := int(penalty / 3.0)
		var res []diff.Result
		for i := 0; i < nMissing; i++ {
			res = append(res, diff.Result{Kind: diff.MissingInB, Key: "K"})
		}
		r := scorer.Compute(res, total)
		if r.Grade != tt.wantGrade {
			t.Errorf("score~%.0f: expected grade %s, got %s (actual score %.1f)",
				tt.score, tt.wantGrade, r.Grade, r.Score)
		}
	}
}

func TestResult_String(t *testing.T) {
	r := scorer.Result{Score: 75.0, Grade: "B", Penalty: 3.0, Total: 4}
	s := r.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
