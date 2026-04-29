package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/envdiff/internal/diff"
	"github.com/yourusername/envdiff/internal/filter"
	"github.com/yourusername/envdiff/internal/parser"
	"github.com/yourusername/envdiff/internal/report"
)

func main() {
	format := flag.String("format", "text", "Output format: text or json")
	onlyKinds := flag.String("only", "", "Comma-separated kinds to show: missing_in_a,missing_in_b,mismatch")
	keyPrefix := flag.String("prefix", "", "Filter results to keys starting with this prefix")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: envdiff [flags] <file-a> <file-b>")
		os.Exit(1)
	}

	envA, err := parser.Parse(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", args[0], err)
		os.Exit(1)
	}

	envB, err := parser.Parse(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", args[1], err)
		os.Exit(1)
	}

	results := diff.Compare(envA, envB)

	var kinds []string
	if *onlyKinds != "" {
		for _, k := range strings.Split(*onlyKinds, ",") {
			if k = strings.TrimSpace(k); k != "" {
				kinds = append(kinds, k)
			}
		}
	}

	results = filter.Apply(results, filter.Options{
		OnlyKinds: kinds,
		KeyPrefix:  *keyPrefix,
	})

	switch *format {
	case "json":
		if err := report.WriteJSON(os.Stdout, results, args[0], args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
			os.Exit(1)
		}
	default:
		if err := report.WriteText(os.Stdout, results, args[0], args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
			os.Exit(1)
		}
	}

	if len(results) > 0 {
		os.Exit(2)
	}
}
