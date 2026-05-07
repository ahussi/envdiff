// Package grapher analyses relationships between multiple env files by
// computing pairwise key-overlap metrics.
//
// Given a set of named env maps the package produces a Result containing:
//
//   - Nodes – sorted list of env file labels
//   - Edges – every pair with shared key count, unique counts and a
//     Jaccard-style similarity score in [0, 1]
//
// A similarity of 1.0 means both files share exactly the same key set;
// 0.0 means they are completely disjoint.
//
// Typical usage:
//
//	result := grapher.Analyse(envs)
//	grapher.WriteText(os.Stdout, result)
package grapher
