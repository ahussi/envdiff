// Package scorer provides a simple health-scoring mechanism for env-file
// comparison results.
//
// A score of 100 means the two files are identical in structure; a score of 0
// means every key differs or is absent.  The score is computed by applying
// weighted penalties per diff kind:
//
//   - MissingInA / MissingInB  → 3.0 penalty points each
//   - Mismatch                 → 1.5 penalty points each
//   - Other                    → 1.0 penalty point  each
//
// The raw penalty is normalised against the theoretical maximum (all keys
// missing) and mapped to the 0–100 range, then converted to a letter grade:
//
//	A ≥ 90   B ≥ 75   C ≥ 60   D ≥ 40   F < 40
//
// Usage:
//
//	results := diff.Compare(a, b)
//	total   := len(a) + len(b) // approximate unique key count
//	r       := scorer.Compute(results, total)
//	fmt.Println(r) // score=82.5 grade=B penalty=10.5 total=20
package scorer
