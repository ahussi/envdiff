// Package linter provides heuristic checks on parsed .env files,
// flagging suspicious or potentially problematic entries.
package linter

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/parser"
)

// Severity indicates how serious a lint finding is.
type Severity string

const (
	Warn  Severity = "warn"
	Error Severity = "error"
)

// Finding represents a single lint issue found in an env file.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

// Check runs all lint rules against the provided env entries and returns
// any findings. The label is used in messages to identify the source file.
func Check(entries []parser.Entry) []Finding {
	var findings []Finding

	for _, e := range entries {
		if f, ok := checkEmptyValue(e); ok {
			findings = append(findings, f)
		}
		if f, ok := checkKeyCase(e); ok {
			findings = append(findings, f)
		}
		if f, ok := checkWhitespaceValue(e); ok {
			findings = append(findings, f)
		}
		if f, ok := checkSensitivePlaintext(e); ok {
			findings = append(findings, f)
		}
	}

	return findings
}

// checkEmptyValue warns when a key is present but has no value.
func checkEmptyValue(e parser.Entry) (Finding, bool) {
	if e.Value == "" {
		return Finding{
			Key:      e.Key,
			Message:  "key has an empty value",
			Severity: Warn,
		}, true
	}
	return Finding{}, false
}

// checkKeyCase warns when a key contains lowercase letters (convention is ALL_CAPS).
func checkKeyCase(e parser.Entry) (Finding, bool) {
	if e.Key != strings.ToUpper(e.Key) {
		return Finding{
			Key:      e.Key,
			Message:  fmt.Sprintf("key %q is not upper-case", e.Key),
			Severity: Warn,
		}, true
	}
	return Finding{}, false
}

// checkWhitespaceValue warns when a value is only whitespace.
func checkWhitespaceValue(e parser.Entry) (Finding, bool) {
	if e.Value != "" && strings.TrimSpace(e.Value) == "" {
		return Finding{
			Key:      e.Key,
			Message:  "value contains only whitespace",
			Severity: Warn,
		}, true
	}
	return Finding{}, false
}

// checkSensitivePlaintext errors when a key looks sensitive but its value
// appears to be a plaintext secret (short, no variable references).
func checkSensitivePlaintext(e parser.Entry) (Finding, bool) {
	sensitivePatterns := []string{"SECRET", "PASSWORD", "TOKEN", "API_KEY", "PRIVATE_KEY"}
	for _, p := range sensitivePatterns {
		if strings.Contains(strings.ToUpper(e.Key), p) {
			if e.Value != "" && !strings.HasPrefix(e.Value, "${") && !strings.HasPrefix(e.Value, "$(") {
				return Finding{
					Key:      e.Key,
					Message:  fmt.Sprintf("key %q may contain a plaintext secret", e.Key),
					Severity: Error,
				}, true
			}
		}
	}
	return Finding{}, false
}
