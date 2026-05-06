package merger_test

import (
	"testing"

	"github.com/user/envdiff/internal/merger"
	"github.com/user/envdiff/internal/parser"
)

func makeNamed(label string, pairs ...string) merger.NamedEnv {
	env := make(parser.EnvFile)
	for i := 0; i+1 < len(pairs); i += 2 {
		env[pairs[i]] = pairs[i+1]
	}
	return merger.NamedEnv{Label: label, Env: env}
}

func TestMerge_StrategyFirst(t *testing.T) {
	a := makeNamed("a.env", "KEY", "from-a", "ONLY_A", "yes")
	b := makeNamed("b.env", "KEY", "from-b", "ONLY_B", "yes")

	res, err := merger.Merge([]merger.NamedEnv{a, b}, merger.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := res.Env["KEY"]; got != "from-a" {
		t.Errorf("KEY = %q, want %q", got, "from-a")
	}
	if res.Sources["KEY"] != "a.env" {
		t.Errorf("source for KEY = %q, want %q", res.Sources["KEY"], "a.env")
	}
	if _, ok := res.Env["ONLY_B"]; !ok {
		t.Error("expected ONLY_B to be present")
	}
}

func TestMerge_StrategyLast(t *testing.T) {
	a := makeNamed("a.env", "KEY", "from-a")
	b := makeNamed("b.env", "KEY", "from-b")

	res, err := merger.Merge([]merger.NamedEnv{a, b}, merger.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := res.Env["KEY"]; got != "from-b" {
		t.Errorf("KEY = %q, want %q", got, "from-b")
	}
	if res.Sources["KEY"] != "b.env" {
		t.Errorf("source for KEY = %q, want %q", res.Sources["KEY"], "b.env")
	}
}

func TestMerge_StrategyError_Conflict(t *testing.T) {
	a := makeNamed("a.env", "KEY", "from-a")
	b := makeNamed("b.env", "KEY", "from-b")

	_, err := merger.Merge([]merger.NamedEnv{a, b}, merger.StrategyError)
	if err == nil {
		t.Fatal("expected error for conflicting keys, got nil")
	}
}

func TestMerge_StrategyError_NoConflict(t *testing.T) {
	a := makeNamed("a.env", "ALPHA", "1")
	b := makeNamed("b.env", "BETA", "2")

	res, err := merger.Merge([]merger.NamedEnv{a, b}, merger.StrategyError)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Env))
	}
}

func TestParseStrategy_Valid(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  merger.Strategy
	}{
		{"first", merger.StrategyFirst},
		{"last", merger.StrategyLast},
		{"error", merger.StrategyError},
	} {
		got, err := merger.ParseStrategy(tc.input)
		if err != nil {
			t.Errorf("ParseStrategy(%q) error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseStrategy(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestParseStrategy_Invalid(t *testing.T) {
	_, err := merger.ParseStrategy("unknown")
	if err == nil {
		t.Error("expected error for unknown strategy, got nil")
	}
}
