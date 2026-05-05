// Package exporter provides utilities for rendering and persisting diff
// results to various output formats (text, JSON, Markdown).
//
// Usage:
//
//	// Write to an io.Writer with an explicit format:
//	 err := exporter.Write(os.Stdout, results, exporter.WriteOptions{
//	     Format: exporter.FormatMarkdown,
//	     FileA:  "staging.env",
//	     FileB:  "production.env",
//	 })
//
//	// Write to a file, inferring the format from the extension:
//	 err := exporter.WriteToFile("report.md", results, exporter.WriteOptions{
//	     FileA: "staging.env",
//	     FileB: "production.env",
//	 })
//
// Supported formats: text (.txt), json (.json), markdown (.md / .markdown).
package exporter
