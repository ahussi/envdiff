package diff

import "github.com/user/envdiff/internal/parser"

// Result holds the comparison outcome between two env files.
type Result struct {
	MissingInB  []parser.Entry // keys present in A but not in B
	MissingInA  []parser.Entry // keys present in B but not in A
	Mismatched  []Mismatch     // keys present in both but with different values
}

// Mismatch describes a key whose value differs between two env files.
type Mismatch struct {
	Key    string
	ValueA string
	ValueB string
}

// Compare compares two EnvFile instances and returns a Result.
func Compare(a, b *parser.EnvFile) Result {
	result := Result{}

	for key, entryA := range a.Entries {
		entryB, exists := b.Entries[key]
		if !exists {
			result.MissingInB = append(result.MissingInB, entryA)
			continue
		}
		if entryA.Value != entryB.Value {
			result.Mismatched = append(result.Mismatched, Mismatch{
				Key:    key,
				ValueA: entryA.Value,
				ValueB: entryB.Value,
			})
		}
	}

	for key, entryB := range b.Entries {
		if _, exists := a.Entries[key]; !exists {
			result.MissingInA = append(result.MissingInA, entryB)
		}
	}

	return result
}

// HasDifferences returns true if the result contains any differences.
func (r Result) HasDifferences() bool {
	return len(r.MissingInA) > 0 || len(r.MissingInB) > 0 || len(r.Mismatched) > 0
}
