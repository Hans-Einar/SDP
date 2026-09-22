package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"io"
	"os"
	"path/filepath"
)

func run() int {
	flags := flag.NewFlagSet("sdui", flag.ContinueOnError)
	output := flags.String("o", "", "Output JSON file (default stdout)")
	flags.Bool("syntax-only", true, "Parse without semantic validation (G1-M1)")
	if err := flags.Parse(os.Args[1:]); err != nil {
		return 2
	}
	if flags.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "Usage: sdui [-o file] [--syntax-only] source.sdui|-")
		return 2
	}
	src := flags.Arg(0)
	var r io.Reader = os.Stdin
	if src != "-" {
		if *output != "" {
			a, _ := filepath.Abs(src)
			b, _ := filepath.Abs(*output)
			sa, ea := os.Stat(a)
			sb, eb := os.Stat(b)
			if a == b || ea == nil && eb == nil && os.SameFile(sa, sb) {
				fmt.Fprintln(os.Stderr, "Output cannot overwrite source")
				return 3
			}
		}
		f, e := os.Open(src)
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
			return 3
		}
		defer f.Close()
		r = f
	}
	data, e := io.ReadAll(io.LimitReader(r, parser.MaxBytes+1))
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		return 3
	}
	doc, e := parser.Parse(string(data))
	if e != nil {
		json.NewEncoder(os.Stderr).Encode(map[string]any{"error": e})
		return 2
	}
	out, e := json.MarshalIndent(map[string]any{"astFormat": "sdui-ast/0.2", "validation": "syntax-only", "document": parser.Data(doc)}, "", "  ")
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		return 3
	}
	out = append(out, '\n')
	if *output != "" {
		e = os.WriteFile(*output, out, 0644)
	} else {
		_, e = os.Stdout.Write(out)
	}
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		return 3
	}
	return 0
}
func main() { os.Exit(run()) }
