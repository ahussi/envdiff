package grapher

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteText writes a human-readable graph summary to w.
func WriteText(w io.Writer, r Result) error {
	if len(r.Edges) == 0 {
		_, err := fmt.Fprintln(w, "No edges: fewer than two env files provided.")
		return err
	}
	fmt.Fprintf(w, "Nodes (%d): %v\n\n", len(r.Nodes), r.Nodes)
	for _, e := range r.Edges {
		fmt.Fprintf(w, "%s <-> %s\n", e.From, e.To)
		fmt.Fprintf(w, "  shared=%d  unique_a=%d  unique_b=%d  similarity=%.2f\n",
			e.SharedKeys, e.UniqueToA, e.UniqueToB, e.Similarity)
	}
	return nil
}

// jsonEdge is the serialisable form of Edge.
type jsonEdge struct {
	From       string  `json:"from"`
	To         string  `json:"to"`
	Shared     int     `json:"shared_keys"`
	UniqueA    int     `json:"unique_to_a"`
	UniqueB    int     `json:"unique_to_b"`
	Similarity float64 `json:"similarity"`
}

// WriteJSON writes the graph as a JSON object to w.
func WriteJSON(w io.Writer, r Result) error {
	type payload struct {
		Nodes []string    `json:"nodes"`
		Edges []jsonEdge  `json:"edges"`
	}
	p := payload{Nodes: r.Nodes}
	for _, e := range r.Edges {
		p.Edges = append(p.Edges, jsonEdge{
			From:       e.From,
			To:         e.To,
			Shared:     e.SharedKeys,
			UniqueA:    e.UniqueToA,
			UniqueB:    e.UniqueToB,
			Similarity: e.Similarity,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p)
}
