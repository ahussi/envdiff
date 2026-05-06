// Package validator provides rule-based validation for environment variable
// maps loaded by envdiff.
//
// A Rule binds a key name to one or more constraints:
//
//   - Required: the key must exist in the env map with a non-empty value.
//   - Pattern:  the value must fully match the supplied regular expression.
//
// Usage:
//
//	rules := []validator.Rule{
//		{Key: "DATABASE_URL", Required: true},
//		{Key: "PORT", Required: true, Pattern: `^\d+$`},
//		{Key: "LOG_LEVEL", Pattern: `^(debug|info|warn|error)$`},
//	}
//
//	violations := validator.Check(env, rules)
//	for _, v := range violations {
//		fmt.Printf("[%s] %s\n", v.Key, v.Message)
//	}
//
// Violations are returned in the same order as the rules slice, making
// output deterministic and easy to test.
package validator
