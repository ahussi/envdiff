// Package profiler provides a single-file health analysis pipeline for env
// files used by envdiff.
//
// # Overview
//
// Analyse parses a .env file, runs the built-in linter, counts redacted
// (sensitive) keys, and computes a numeric health score.  The result is
// returned as a [Profile] value that callers can render in any format.
//
// # Usage
//
//	p, err := profiler.Analyse(".env.production", "production", nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Score: %s (%d/100)\n", p.Score.Grade, p.Score.Value)
//
// # Extra redaction patterns
//
// Pass additional glob patterns via the extraPatterns argument to treat
// project-specific keys as sensitive (e.g. "*_TOKEN", "*_SECRET").
package profiler
