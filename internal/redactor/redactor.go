// Package redactor masks sensitive values in env entries before display or export.
package redactor

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// defaultSensitivePatterns are substrings that indicate a key holds a secret.
var defaultSensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY",
	"PRIVATE", "CREDENTIAL", "AUTH", "PWD",
}

const masked = "***REDACTED***"

// Options controls redaction behaviour.
type Options struct {
	// ExtraPatterns are additional substrings (case-insensitive) to treat as sensitive.
	ExtraPatterns []string
	// Disable turns off all redaction.
	Disable bool
}

// Apply returns a copy of results with sensitive values masked.
func Apply(results []diff.Result, opts Options) []diff.Result {
	if opts.Disable {
		return results
	}

	patterns := buildPatterns(opts.ExtraPatterns)
	out := make([]diff.Result, len(results))

	for i, r := range results {
		if isSensitive(r.Key, patterns) {
			r.ValueA = redactNonEmpty(r.ValueA)
			r.ValueB = redactNonEmpty(r.ValueB)
		}
		out[i] = r
	}

	return out
}

func buildPatterns(extra []string) []string {
	patterns := make([]string, len(defaultSensitivePatterns))
	copy(patterns, defaultSensitivePatterns)
	for _, p := range extra {
		patterns = append(patterns, strings.ToUpper(p))
	}
	return patterns
}

func isSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

func redactNonEmpty(v string) string {
	if v == "" {
		return v
	}
	return masked
}
