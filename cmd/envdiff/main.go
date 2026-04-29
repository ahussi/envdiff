package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envdiff/internal/diff"
	"github.com/yourorg/envdiff/internal/filter"
	"github.com/yourorg/envdiff/internal/formatter"
	"github.com/yourorg/envdiff/internal/parser"
	"github.com/yourorg/envdiff/internal/report"
	"github.com/yourorg/envdiff/internal/sorter"
)

func main() {
	format := flag.String("format", "text", "output format: text, json, markdown")
	style := flag.String("style", "plain", "output style: plain, color, markdown")
	kindFlag := flag.String("kind", "", "filter by kind: missing_in_a, missing_in_b, mismatch")
	prefixFlag := flag.String("prefix", "", "filter keys by prefix")
	sortFlag := flag.String("sort", "", "sort results: key_asc, key_desc, kind_asc, kind_desc")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: envdiff [flags] <file-a> <file-b>")
		os.Exit(1)
	}

	_, err := formatter.ParseStyle(*style)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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
	results = filter.Apply(results, filter.Options{
		Kind:   *kindFlag,
		Prefix: *prefixFlag,
	})
	results = sorter.Apply(results, sorter.Options{Sort: *sortFlag})

	switch *format {
	case "json":
		if err := report.WriteJSON(os.Stdout, results); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "markdown":
		if err := report.WriteMarkdown(os.Stdout, results, args[0], args[1]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		if err := report.WriteText(os.Stdout, results); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if len(results) > 0 {
		os.Exit(2)
	}
}
