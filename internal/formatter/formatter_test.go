package formatter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/formatter"
)

func TestParseStyle_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  formatter.Style
	}{
		{"plain", formatter.StylePlain},
		{"", formatter.StylePlain},
		{"color", formatter.StyleColor},
		{"markdown", formatter.StyleMarkdown},
		{"MARKDOWN", formatter.StyleMarkdown},
	}
	for _, tc := range cases {
		got, err := formatter.ParseStyle(tc.input)
		if err != nil {
			t.Errorf("ParseStyle(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseStyle(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestParseStyle_Invalid(t *testing.T) {
	_, err := formatter.ParseStyle("fancy")
	if err == nil {
		t.Error("expected error for unknown style, got nil")
	}
}

func TestKindLabel_Plain(t *testing.T) {
	got := formatter.KindLabel("mismatch", formatter.StylePlain)
	if got != "mismatch" {
		t.Errorf("expected plain label, got %q", got)
	}
}

func TestKindLabel_Markdown(t *testing.T) {
	got := formatter.KindLabel("missing_in_b", formatter.StyleMarkdown)
	if !strings.HasPrefix(got, "**") || !strings.HasSuffix(got, "**") {
		t.Errorf("expected markdown bold, got %q", got)
	}
}

func TestKindLabel_Color(t *testing.T) {
	got := formatter.KindLabel("mismatch", formatter.StyleColor)
	if !strings.Contains(got, "\033[") {
		t.Errorf("expected ANSI escape in color label, got %q", got)
	}
}

func TestKeyLabel_Markdown(t *testing.T) {
	got := formatter.KeyLabel("MY_KEY", formatter.StyleMarkdown)
	if got != "`MY_KEY`" {
		t.Errorf("expected backtick-wrapped key, got %q", got)
	}
}

func TestKeyLabel_Plain(t *testing.T) {
	got := formatter.KeyLabel("MY_KEY", formatter.StylePlain)
	if got != "MY_KEY" {
		t.Errorf("expected plain key, got %q", got)
	}
}
