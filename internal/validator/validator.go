// Package validator checks .env values against simple rules such as
// required keys and allowed value patterns.
package validator

import (
	"fmt"
	"regexp"
)

// Rule describes a validation constraint for a single key.
type Rule struct {
	// Key is the env key this rule applies to.
	Key string
	// Required asserts the key must be present and non-empty.
	Required bool
	// Pattern, when non-empty, is a regex the value must fully match.
	Pattern string
}

// Violation records a failed validation.
type Violation struct {
	Key     string
	Message string
}

// Check validates the provided env map against the given rules and returns
// any violations found. An empty slice means all rules passed.
func Check(env map[string]string, rules []Rule) []Violation {
	var violations []Violation

	for _, r := range rules {
		val, exists := env[r.Key]

		if r.Required && (!exists || val == "") {
			violations = append(violations, Violation{
				Key:     r.Key,
				Message: "required key is missing or empty",
			})
			continue
		}

		if r.Pattern != "" && exists && val != "" {
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				violations = append(violations, Violation{
					Key:     r.Key,
					Message: fmt.Sprintf("invalid pattern %q: %v", r.Pattern, err),
				})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, Violation{
					Key:     r.Key,
					Message: fmt.Sprintf("value %q does not match pattern %q", val, r.Pattern),
				})
			}
		}
	}

	return violations
}
