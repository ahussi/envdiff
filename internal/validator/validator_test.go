package validator_test

import (
	"testing"

	"github.com/your-org/envdiff/internal/validator"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestCheck_RequiredPresent(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost")
	rules := []validator.Rule{{Key: "DB_HOST", Required: true}}
	got := validator.Check(env, rules)
	if len(got) != 0 {
		t.Fatalf("expected no violations, got %v", got)
	}
}

func TestCheck_RequiredMissing(t *testing.T) {
	env := makeEnv()
	rules := []validator.Rule{{Key: "DB_HOST", Required: true}}
	got := validator.Check(env, rules)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
	if got[0].Key != "DB_HOST" {
		t.Errorf("unexpected key %q", got[0].Key)
	}
}

func TestCheck_RequiredEmptyValue(t *testing.T) {
	env := makeEnv("API_KEY", "")
	rules := []validator.Rule{{Key: "API_KEY", Required: true}}
	got := validator.Check(env, rules)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
}

func TestCheck_PatternMatch(t *testing.T) {
	env := makeEnv("PORT", "8080")
	rules := []validator.Rule{{Key: "PORT", Pattern: `^\d+$`}}
	got := validator.Check(env, rules)
	if len(got) != 0 {
		t.Fatalf("expected no violations, got %v", got)
	}
}

func TestCheck_PatternNoMatch(t *testing.T) {
	env := makeEnv("PORT", "abc")
	rules := []validator.Rule{{Key: "PORT", Pattern: `^\d+$`}}
	got := validator.Check(env, rules)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(got))
	}
}

func TestCheck_InvalidPattern(t *testing.T) {
	env := makeEnv("FOO", "bar")
	rules := []validator.Rule{{Key: "FOO", Pattern: `[invalid`}}
	got := validator.Check(env, rules)
	if len(got) != 1 {
		t.Fatalf("expected 1 violation for bad regex, got %d", len(got))
	}
}

func TestCheck_NoRules(t *testing.T) {
	env := makeEnv("X", "1")
	got := validator.Check(env, nil)
	if len(got) != 0 {
		t.Fatalf("expected no violations, got %v", got)
	}
}

func TestCheck_MultipleViolations(t *testing.T) {
	env := makeEnv("PORT", "xyz")
	rules := []validator.Rule{
		{Key: "DB_HOST", Required: true},
		{Key: "PORT", Pattern: `^\d+$`},
	}
	got := validator.Check(env, rules)
	if len(got) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(got))
	}
}
