package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
)

func execute(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) > 0 && (args[0] == "class-check" || args[0] == "class-view") {
		return classCommand(args, stdout, stderr)
	}
	if len(args) > 0 && args[0] == "view" {
		return selectView(args[1:], stdout, stderr)
	}
	if len(args) > 0 && args[0] == "viewpoints" {
		return exportViews(args[1:], stdout, stderr)
	}
	if len(args) != 2 || (args[0] != "check" && args[0] != "ast" && args[0] != "format" && args[0] != "action-check") {
		fmt.Fprintln(stderr, "Usage: sdl check|ast|format|action-check file|-")
		return 2
	}
	input := stdin
	if args[1] != "-" {
		f, e := os.Open(args[1])
		if e != nil {
			fmt.Fprintln(stderr, e)
			return 2
		}
		defer f.Close()
		input = f
	}
	source, e := io.ReadAll(io.LimitReader(input, parser.MaxBytes+1))
	if e != nil {
		fmt.Fprintln(stderr, e)
		return 2
	}
	if args[0] == "action-check" {
		_, err := parser.CompileActions(string(source))
		report := map[string]any{"valid": err == nil, "profile": "action-core/0.1"}
		if err != nil {
			report["diagnostic"] = parser.Data(err.(parser.Diagnostic))
		}
		if json.NewEncoder(stdout).Encode(report) != nil {
			return 2
		}
		if err != nil {
			return 1
		}
		return 0
	}
	if args[0] == "format" {
		m, e := parser.Parse(string(source))
		if e != nil {
			_ = json.NewEncoder(stderr).Encode(parser.Data(e.(parser.Diagnostic)))
			return 1
		}
		canonical, d := parser.Canonical(m)
		if len(d) > 0 {
			_ = json.NewEncoder(stderr).Encode(parser.Data(d))
			return 1
		}
		_, e = io.WriteString(stdout, canonical)
		if e != nil {
			return 2
		}
		return 0
	}
	m, d := parser.Check(string(source))
	report := map[string]any{"valid": len(d) == 0, "diagnostics": parser.Data(d)}
	if args[0] == "ast" {
		report["ast"] = parser.Data(m)
		symbols := map[string]any{}
		if m != nil {
			for name, decl := range parser.Symbols(m) {
				symbols[name] = parser.Data(decl)
			}
		}
		report["symbols"] = symbols
	}
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	if e = enc.Encode(report); e != nil {
		return 2
	}
	if len(d) > 0 {
		return 1
	}
	return 0
}
func main() { os.Exit(execute(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)) }
