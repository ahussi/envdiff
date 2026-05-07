// Package grapher builds a dependency graph between env files based on
// shared and divergent keys, producing a weighted adjacency representation
// useful for visualising env file relationships.
package grapher

import "sort"

// Edge represents a relationship between two env file labels.
type Edge struct {
	From        string
	To          string
	SharedKeys  int
	UniqueToA   int
	UniqueToB   int
	Similarity  float64 // 0.0 – 1.0
}

// Result holds the full graph output.
type Result struct {
	Nodes []string
	Edges []Edge
}

// Analyse computes pairwise similarity between named env maps.
// envs maps a label (e.g. ".env.production") to its key→value pairs.
func Analyse(envs map[string]map[string]string) Result {
	labels := make([]string, 0, len(envs))
	for l := range envs {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	var edges []Edge
	for i := 0; i < len(labels); i++ {
		for j := i + 1; j < len(labels); j++ {
			a, b := labels[i], labels[j]
			e := buildEdge(a, b, envs[a], envs[b])
			edges = append(edges, e)
		}
	}

	return Result{Nodes: labels, Edges: edges}
}

func buildEdge(a, b string, ma, mb map[string]string) Edge {
	shared := 0
	for k := range ma {
		if _, ok := mb[k]; ok {
			shared++
		}
	}
	uniqueA := len(ma) - shared
	uniqueB := len(mb) - shared
	total := len(ma) + len(mb) - shared
	var sim float64
	if total > 0 {
		sim = float64(shared) / float64(total)
	}
	return Edge{
		From:       a,
		To:         b,
		SharedKeys: shared,
		UniqueToA:  uniqueA,
		UniqueToB:  uniqueB,
		Similarity: sim,
	}
}
