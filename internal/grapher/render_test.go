package grapher

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func makeResult() Result {
	return Result{
		Nodes: []string{"dev", "prod"},
		Edges: []Edge{
			{
				From:       "dev",
				To:         "prod",
				SharedKeys: 3,
				UniqueToA:  1,
				UniqueToB:  2,
				Similarity: 0.5,
			},
		},
	}
}

func TestWriteText_ContainsNodes(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteText(&buf, makeResult()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "dev") || !strings.Contains(buf.String(), "prod") {
		t.Errorf("expected node labels in output")
	}
}

func TestWriteText_ContainsSimilarity(t *testing.T) {
	var buf bytes.Buffer
	WriteText(&buf, makeResult())
	if !strings.Contains(buf.String(), "0.50") {
		t.Errorf("expected similarity score in output")
	}
}

func TestWriteText_EmptyResult(t *testing.T) {
	var buf bytes.Buffer
	WriteText(&buf, Result{})
	if !strings.Contains(buf.String(), "No edges") {
		t.Errorf("expected no-edges message")
	}
}

func TestWriteJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteJSON(&buf, makeResult()); err != nil {
		t.Fatal(err)
	}
	var out struct {
		Nodes []string `json:"nodes"`
		Edges []struct {
			From       string  `json:"from"`
			Similarity float64 `json:"similarity"`
		} `json:"edges"`
	}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out.Nodes) != 2 {
		t.Errorf("nodes=%d want 2", len(out.Nodes))
	}
	if out.Edges[0].Similarity != 0.5 {
		t.Errorf("similarity=%.2f want 0.50", out.Edges[0].Similarity)
	}
}
