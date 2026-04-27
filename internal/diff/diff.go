package diff

// DiffKind describes the type of difference found.
type DiffKind int

const (
	MissingInB DiffKind = iota
	MissingInA
	Mismatch
)

// Result holds a single diff between two env files.
type Result struct {
	Kind   DiffKind
	Key    string
	ValueA string
	ValueB string
}

// Compare returns the list of differences between two parsed env maps.
// envA and envB map key -> value.
func Compare(envA, envB map[string]string) []Result {
	var results []Result

	for k, vA := range envA {
		vB, ok := envB[k]
		if !ok {
			results = append(results, Result{Kind: MissingInB, Key: k, ValueA: vA})
			continue
		}
		if vA != vB {
			results = append(results, Result{Kind: Mismatch, Key: k, ValueA: vA, ValueB: vB})
		}
	}

	for k, vB := range envB {
		if _, ok := envA[k]; !ok {
			results = append(results, Result{Kind: MissingInA, Key: k, ValueB: vB})
		}
	}

	return results
}
