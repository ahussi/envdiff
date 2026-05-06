package profiler

import (
	"encoding/json"
	"fmt"
	"io"
)

// WriteText writes a human-readable profile report to w.
func WriteText(w io.Writer, p Profile) {
	fmt.Fprintf(w, "Profile: %s\n", p.Label)
	fmt.Fprintf(w, "  Total keys : %d\n", p.TotalKeys)
	fmt.Fprintf(w, "  Score      : %s (%d/100)\n", p.Score.Grade, p.Score.Value)
	fmt.Fprintf(w, "  Redacted   : %d key(s)\n", p.RedactedKeys)

	if len(p.LintIssues) == 0 {
		fmt.Fprintln(w, "  Lint       : no issues")
		return
	}

	fmt.Fprintf(w, "  Lint       : %d issue(s)\n", len(p.LintIssues))
	for _, iss := range p.LintIssues {
		fmt.Fprintf(w, "    [%s] %s — %s\n", iss.Severity, iss.Key, iss.Message)
	}
}

// jsonProfile is the JSON-serialisable representation of a Profile.
type jsonProfile struct {
	Label        string        `json:"label"`
	TotalKeys    int           `json:"total_keys"`
	Score        int           `json:"score"`
	Grade        string        `json:"grade"`
	RedactedKeys int           `json:"redacted_keys"`
	LintIssues   []jsonIssue   `json:"lint_issues"`
}

type jsonIssue struct {
	Key      string `json:"key"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// WriteJSON writes a JSON-encoded profile to w.
func WriteJSON(w io.Writer, p Profile) error {
	issues := make([]jsonIssue, len(p.LintIssues))
	for i, iss := range p.LintIssues {
		issues[i] = jsonIssue{
			Key:      iss.Key,
			Severity: iss.Severity,
			Message:  iss.Message,
		}
	}

	out := jsonProfile{
		Label:        p.Label,
		TotalKeys:    p.TotalKeys,
		Score:        p.Score.Value,
		Grade:        p.Score.Grade,
		RedactedKeys: p.RedactedKeys,
		LintIssues:   issues,
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
