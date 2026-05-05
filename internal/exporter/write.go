package exporter

import (
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/report"
)

// WriteOptions controls how the export is rendered.
type WriteOptions struct {
	Format      Format
	FileA       string
	FileB       string
	ShowSummary bool
}

// Write renders diff results to w using the given options.
func Write(w io.Writer, results []diff.Result, opts WriteOptions) error {
	switch opts.Format {
	case FormatJSON:
		if err := report.WriteJSON(w, results, opts.FileA, opts.FileB); err != nil {
			return fmt.Errorf("exporter: write json: %w", err)
		}
		if opts.ShowSummary {
			if err := report.WriteSummaryJSON(w, results); err != nil {
				return fmt.Errorf("exporter: write summary json: %w", err)
			}
		}
	case FormatMarkdown:
		if err := report.WriteMarkdown(w, results, opts.FileA, opts.FileB); err != nil {
			return fmt.Errorf("exporter: write markdown: %w", err)
		}
	default:
		if err := report.WriteText(w, results, opts.FileA, opts.FileB); err != nil {
			return fmt.Errorf("exporter: write text: %w", err)
		}
		if opts.ShowSummary {
			if err := report.WriteSummaryText(w, results); err != nil {
				return fmt.Errorf("exporter: write summary text: %w", err)
			}
		}
	}
	return nil
}

// WriteToFile renders diff results to the file at path.
// The format is inferred from the file extension unless overridden.
func WriteToFile(path string, results []diff.Result, opts WriteOptions) error {
	if opts.Format == "" {
		fmt, err := InferFormat(path)
		if err != nil {
			return err
		}
		opts.Format = fmt
	}
	w, err := OpenOutput(path)
	if err != nil {
		return err
	}
	defer w.Close()
	return Write(w, results, opts)
}
