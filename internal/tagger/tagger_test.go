package tagger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
)

func makeResult(key, kind, a, b string) diff.Result {
	return diff.Result{Key: key, Kind: diff.Kind(kind), ValueA: a, ValueB: b}
}

func TestAnnotate_SecretKey(t *testing.T) {
	results := []diff.Result{makeResult("DB_PASSWORD", "mismatch", "old", "new")}
	annotated := Annotate(results)
	if len(annotated) != 1 {
		t.Fatalf("expected 1 result, got %d", len(annotated))
	}
	if !containsTag(annotated[0].Tags, TagSecret) {
		t.Errorf("expected secret tag for DB_PASSWORD, got %v", annotated[0].Tags)
	}
}

func TestAnnotate_URLKey(t *testing.T) {
	results := []diff.Result{makeResult("DATABASE_URL", "missing_in_b", "val", "")}
	annotated := Annotate(results)
	if !containsTag(annotated[0].Tags, TagURL) {
		t.Errorf("expected url tag for DATABASE_URL, got %v", annotated[0].Tags)
	}
}

func TestAnnotate_PortKey(t *testing.T) {
	results := []diff.Result{makeResult("APP_PORT", "mismatch", "8080", "9090")}
	annotated := Annotate(results)
	if !containsTag(annotated[0].Tags, TagPort) {
		t.Errorf("expected port tag for APP_PORT, got %v", annotated[0].Tags)
	}
}

func TestAnnotate_UnknownKey(t *testing.T) {
	results := []diff.Result{makeResult("FOOBAR", "mismatch", "x", "y")}
	annotated := Annotate(results)
	if !containsTag(annotated[0].Tags, TagUnknown) {
		t.Errorf("expected unknown tag for FOOBAR, got %v", annotated[0].Tags)
	}
}

func TestAnnotate_FeatureFlag(t *testing.T) {
	results := []diff.Result{makeResult("ENABLE_DARK_MODE", "missing_in_a", "", "true")}
	annotated := Annotate(results)
	if !containsTag(annotated[0].Tags, TagFeature) {
		t.Errorf("expected feature tag for ENABLE_DARK_MODE, got %v", annotated[0].Tags)
	}
}

func TestWriteText_NoDiffs(t *testing.T) {
	var buf bytes.Buffer
	WriteText(&buf, []Result{})
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-differences message, got: %s", buf.String())
	}
}

func TestWriteText_WithDiffs(t *testing.T) {
	results := Annotate([]diff.Result{makeResult("API_TOKEN", "mismatch", "abc", "xyz")})
	var buf bytes.Buffer
	WriteText(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "API_TOKEN") {
		t.Errorf("expected key in output, got: %s", out)
	}
	if !strings.Contains(out, "secret") {
		t.Errorf("expected secret tag in output, got: %s", out)
	}
}

func TestWriteJSON_Structure(t *testing.T) {
	results := Annotate([]diff.Result{makeResult("LOG_PATH", "missing_in_b", "/var/log", "")})
	var buf bytes.Buffer
	if err := WriteJSON(&buf, results); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "LOG_PATH") {
		t.Errorf("expected key in JSON output")
	}
	if !strings.Contains(out, "path") {
		t.Errorf("expected path tag in JSON output")
	}
}

func containsTag(tags []Tag, target Tag) bool {
	for _, t := range tags {
		if t == target {
			return true
		}
	}
	return false
}
