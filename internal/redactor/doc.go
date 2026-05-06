// Package redactor provides value-masking for sensitive environment variables.
//
// It inspects the key name of each diff.Result and replaces the ValueA / ValueB
// fields with a placeholder string when the key matches a known sensitive
// pattern (e.g. SECRET, TOKEN, PASSWORD).
//
// Usage:
//
//	import "github.com/user/envdiff/internal/redactor"
//
//	results = redactor.Apply(results, redactor.Options{
//		ExtraPatterns: []string{"CERT", "SIGNING"},
//	})
//
// Redaction is purely cosmetic — it does not alter the underlying parsed data.
package redactor
