package main

import (
	"flag"
	"fmt"
	"os"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
	"envdiff/internal/report"
)

func main() {
	var (
		fileA   = flag.String("a", "", "Path to the first .env file (required)")
		fileB   = flag.String("b", "", "Path to the second .env file (required)")
		format  = flag.String("format", "text", "Output format: text or json")
		labelA  = flag.String("label-a", "", "Label for the first file (defaults to filename)")
		labelB  = flag.String("label-b", "", "Label for the second file (defaults to filename)")
	)
	flag.Parse()

	if *fileA == "" || *fileB == "" {
		fmt.Fprintln(os.Stderr, "error: both -a and -b flags are required")
		flag.Usage()
		os.Exit(1)
	}

	envA, err := parser.Parse(*fileA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", *fileA, err)
		os.Exit(1)
	}

	envB, err := parser.Parse(*fileB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", *fileB, err)
		os.Exit(1)
	}

	if *labelA == "" {
		*labelA = *fileA
	}
	if *labelB == "" {
		*labelB = *fileB
	}

	diffs := diff.Compare(envA, envB)

	switch *format {
	case "json":
		if err := report.WriteJSON(os.Stdout, diffs, *labelA, *labelB); err != nil {
			fmt.Fprintf(os.Stderr, "error writing JSON report: %v\n", err)
			os.Exit(1)
		}
	case "text":
		report.WriteText(os.Stdout, diffs, *labelA, *labelB)
	default:
		fmt.Fprintf(os.Stderr, "error: unknown format %q, expected 'text' or 'json'\n", *format)
		os.Exit(1)
	}

	if len(diffs) > 0 {
		os.Exit(2)
	}
}
