package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/sorter"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "ZEBRA", Kind: diff.Mismatch, ValueA: "1", ValueB: "2"},
		{Key: "ALPHA", Kind: diff.MissingInB},
		{Key: "MANGO", Kind: diff.MissingInA},
		{Key: "BETA", Kind: diff.Mismatch, ValueA: "x", ValueB: "y"},
	}
}

func TestApply_NoOptions(t *testing.T) {
	input := makeResults()
	out := sorter.Apply(input, sorter.Options{})
	if len(out) != len(input) {
		t.Fatalf("expected %d results, got %d", len(input), len(out))
	}
	// order should be unchanged
	for i, r := range out {
		if r.Key != input[i].Key {
			t.Errorf("position %d: expected key %q, got %q", i, input[i].Key, r.Key)
		}
	}
}

func TestApply_SortByKeyAsc(t *testing.T) {
	out := sorter.Apply(makeResults(), sorter.Options{By: sorter.SortByKey, Order: sorter.OrderAsc})
	expected := []string{"ALPHA", "BETA", "MANGO", "ZEBRA"}
	for i, key := range expected {
		if out[i].Key != key {
			t.Errorf("position %d: expected %q, got %q", i, key, out[i].Key)
		}
	}
}

func TestApply_SortByKeyDesc(t *testing.T) {
	out := sorter.Apply(makeResults(), sorter.Options{By: sorter.SortByKey, Order: sorter.OrderDesc})
	expected := []string{"ZEBRA", "MANGO", "BETA", "ALPHA"}
	for i, key := range expected {
		if out[i].Key != key {
			t.Errorf("position %d: expected %q, got %q", i, key, out[i].Key)
		}
	}
}

func TestApply_SortByKindAsc(t *testing.T) {
	out := sorter.Apply(makeResults(), sorter.Options{By: sorter.SortByKind, Order: sorter.OrderAsc})
	// MissingInB(0) < MissingInA(1) < Mismatch(2)
	if out[0].Kind != diff.MissingInB {
		t.Errorf("expected first kind MissingInB, got %v", out[0].Kind)
	}
	if out[1].Kind != diff.MissingInA {
		t.Errorf("expected second kind MissingInA, got %v", out[1].Kind)
	}
	// last two should be Mismatch, sorted by key
	if out[2].Key != "BETA" || out[3].Key != "ZEBRA" {
		t.Errorf("expected BETA then ZEBRA for Mismatch group, got %q %q", out[2].Key, out[3].Key)
	}
}

func TestApply_EmptySlice(t *testing.T) {
	out := sorter.Apply([]diff.Result{}, sorter.Options{By: sorter.SortByKey})
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d results", len(out))
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	input := makeResults()
	originalFirst := input[0].Key
	sorter.Apply(input, sorter.Options{By: sorter.SortByKey, Order: sorter.OrderAsc})
	if input[0].Key != originalFirst {
		t.Errorf("Apply mutated the input slice")
	}
}
