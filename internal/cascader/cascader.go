// Package cascader merges multiple .env files in priority order,
// simulating how tools like docker-compose or dotenv-flow resolve
// environment variables across base, environment, and local overrides.
package cascader

import (
	"fmt"

	"github.com/your/envdiff/internal/loader"
)

// Layer represents a single .env file in the cascade chain.
type Layer struct {
	Label string
	Path  string
}

// Result holds the resolved environment map and the provenance of each key.
type Result struct {
	// Resolved is the final merged map of key -> value.
	Resolved map[string]string
	// Provenance maps each key to the label of the layer that provided it.
	Provenance map[string]string
	// Layers is the ordered list of layers used (lowest to highest priority).
	Layers []Layer
}

// Cascade resolves environment variables by loading each layer in order.
// Later layers override earlier ones (last-write-wins). Empty paths are skipped.
func Cascade(layers []Layer) (*Result, error) {
	if len(layers) == 0 {
		return nil, fmt.Errorf("cascader: at least one layer is required")
	}

	resolved := make(map[string]string)
	provenance := make(map[string]string)

	for _, layer := range layers {
		if layer.Path == "" {
			continue
		}
		env, err := loader.Load(layer.Path, layer.Label)
		if err != nil {
			return nil, fmt.Errorf("cascader: loading layer %q: %w", layer.Label, err)
		}
		for _, entry := range env.Entries {
			resolved[entry.Key] = entry.Value
			provenance[entry.Key] = layer.Label
		}
	}

	return &Result{
		Resolved:   resolved,
		Provenance: provenance,
		Layers:     layers,
	}, nil
}
