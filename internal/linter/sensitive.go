package linter

import "strings"

// sensitiveSubstrings is the shared list of substrings that indicate a key
// holds sensitive data. It is used by both the linter and the templater.
var sensitiveSubstrings = []string{
	"secret",
	"password",
	"passwd",
	"token",
	"apikey",
	"api_key",
	"private",
	"credential",
}

// IsSensitiveKey reports whether the given (already lower-cased) key name
// suggests the variable holds sensitive data.
func IsSensitiveKey(lowerKey string) bool {
	for _, sub := range sensitiveSubstrings {
		if strings.Contains(lowerKey, sub) {
			return true
		}
	}
	return false
}
