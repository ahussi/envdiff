// Package linter provides heuristic checks against parsed .env file entries.
//
// It applies a set of built-in rules to surface common mistakes:
//
//   - Empty values: a key is defined but has no value.
//   - Non-upper-case keys: convention expects ALL_CAPS keys.
//   - Whitespace-only values: value exists but is invisible.
//   - Plaintext secrets: keys that look sensitive (e.g. *_SECRET, *_TOKEN)
//     whose values appear to be literal strings rather than variable references.
//
// Usage:
//
//	entries, _ := parser.ParseFile(".env")
//	findings := linter.Check(entries)
//	for _, f := range findings {
//		fmt.Printf("[%s] %s: %s\n", f.Severity, f.Key, f.Message)
//	}
//
// Findings carry a Severity of either Warn or Error. Callers decide how to
// surface or gate on these findings (e.g. exit non-zero on any Error).
package linter
