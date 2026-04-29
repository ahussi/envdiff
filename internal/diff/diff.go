package diff

// Kind describes the type of difference found between two env files.
type Kind int

const (
	MissingInB Kind = iota // key present in A but absent in B
	MissingInA             // key present in B but absent in A
	Mismatch               // key present in both but values differ
)

// Result holds a single comparison finding.
type Result struct {
	Key    string
	Kind   Kind
	ValueA string // value from file A (empty when MissingInA)
	ValueB string // value from file B (empty when MissingInB)
}

// EnvFile is a map of environment variable key → value.
type EnvFile map[string]string

// Compare compares two EnvFile maps and returns all differences.
// Results are returned in deterministic order (sorted by key within each kind).
func Compare(a, b EnvFile) []Result {
	var results []Result

	for key, valA := range a {
		if valB, ok := b[key]; !ok {
			results = append(results, Result{Key: key, Kind: MissingInB, ValueA: valA})
		} else if valA != valB {
			results = append(results, Result{Key: key, Kind: Mismatch, ValueA: valA, ValueB: valB})
		}
	}

	for key, valB := range b {
		if _, ok := a[key]; !ok {
			results = append(results, Result{Key: key, Kind: MissingInA, ValueB: valB})
		}
	}

	return results
}
