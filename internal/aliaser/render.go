package aliaser

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// WriteText writes a human-readable alias report to w.
func WriteText(w io.Writer, r Result) error {
	if len(r.Groups) == 0 {
		_, err := fmt.Fprintln(w, "No alias groups detected.")
		return err
	}

	for _, g := range r.Groups {
		_, err := fmt.Fprintf(w, "[alias] %s  (value: %q)\n",
			strings.Join(g.Keys, " = "), g.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteJSON writes the result as indented JSON to w.
func WriteJSON(w io.Writer, r Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
