package redactor_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/redactor"
)

func makeResults(pairs ...[]string) []diff.Result {
	results := make([]diff.Result, 0, len(pairs))
	for _, p := range pairs {
		results = append(results, diff.Result{
			Key:    p[0],
			ValueA: p[1],
			ValueB: p[2],
			Kind:   diff.KindMismatch,
		})
	}
	return results
}

func TestApply_NonSensitiveUnchanged(t *testing.T) {
	input := makeResults([]string{"APP_ENV", "production", "staging"})
	out := redactor.Apply(input, redactor.Options{})
	if out[0].ValueA != "production" || out[0].ValueB != "staging" {
		t.Errorf("expected values unchanged, got %q / %q", out[0].ValueA, out[0].ValueB)
	}
}

func TestApply_SecretKeyRedacted(t *testing.T) {
	input := makeResults([]string{"DB_PASSWORD", "hunter2", "s3cr3t"})
	out := redactor.Apply(input, redactor.Options{})
	if out[0].ValueA != "***REDACTED***" {
		t.Errorf("expected ValueA redacted, got %q", out[0].ValueA)
	}
	if out[0].ValueB != "***REDACTED***" {
		t.Errorf("expected ValueB redacted, got %q", out[0].ValueB)
	}
}

func TestApply_EmptyValueNotRedacted(t *testing.T) {
	input := makeResults([]string{"API_KEY", "abc123", ""})
	out := redactor.Apply(input, redactor.Options{})
	if out[0].ValueA != "***REDACTED***" {
		t.Errorf("expected ValueA redacted, got %q", out[0].ValueA)
	}
	if out[0].ValueB != "" {
		t.Errorf("expected empty ValueB preserved, got %q", out[0].ValueB)
	}
}

func TestApply_ExtraPatterns(t *testing.T) {
	input := makeResults([]string{"SIGNING_CERT", "cert-data", "other-cert"})
	out := redactor.Apply(input, redactor.Options{ExtraPatterns: []string{"CERT"}})
	if out[0].ValueA != "***REDACTED***" {
		t.Errorf("expected custom pattern to trigger redaction, got %q", out[0].ValueA)
	}
}

func TestApply_DisableSkipsRedaction(t *testing.T) {
	input := makeResults([]string{"DB_SECRET", "topsecret", "alsosecret"})
	out := redactor.Apply(input, redactor.Options{Disable: true})
	if out[0].ValueA != "topsecret" {
		t.Errorf("expected redaction disabled, got %q", out[0].ValueA)
	}
}

func TestApply_CaseInsensitiveKey(t *testing.T) {
	input := makeResults([]string{"db_token", "mytoken", "othertoken"})
	out := redactor.Apply(input, redactor.Options{})
	if out[0].ValueA != "***REDACTED***" {
		t.Errorf("expected lowercase key matched, got %q", out[0].ValueA)
	}
}
