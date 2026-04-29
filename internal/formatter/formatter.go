package formatter

import (
	"fmt"
	"strings"
)

// Style defines the output formatting style for diff results.
type Style int

const (
	StylePlain Style = iota
	StyleColor
	StyleMarkdown
)

// ParseStyle converts a string to a Style value.
func ParseStyle(s string) (Style, error) {
	switch strings.ToLower(s) {
	case "plain", "":
		return StylePlain, nil
	case "color":
		return StyleColor, nil
	case "markdown":
		return StyleMarkdown, nil
	default:
		return StylePlain, fmt.Errorf("unknown style %q: must be plain, color, or markdown", s)
	}
}

// String returns the ANSI-colored or plain representation of a diff kind label.
func KindLabel(kind string, style Style) string {
	switch style {
	case StyleColor:
		return colorize(kind)
	case StyleMarkdown:
		return fmt.Sprintf("**%s**", kind)
	default:
		return kind
	}
}

// KeyLabel formats a key according to the style.
func KeyLabel(key string, style Style) string {
	switch style {
	case StyleColor:
		return fmt.Sprintf("\033[1m%s\033[0m", key)
	case StyleMarkdown:
		return fmt.Sprintf("`%s`", key)
	default:
		return key
	}
}

func colorize(kind string) string {
	const (
		red    = "\033[31m"
		yellow = "\033[33m"
		cyan   = "\033[36m"
		reset  = "\033[0m"
	)
	switch strings.ToLower(kind) {
	case "missing_in_b", "missing":
		return red + kind + reset
	case "missing_in_a":
		return cyan + kind + reset
	case "mismatch":
		return yellow + kind + reset
	default:
		return kind
	}
}
