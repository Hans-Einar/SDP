package main

import (
	"flag"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/codegen"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	actions := flag.String("actions", "", "SDL action-core source")
	source := flag.String("ui", "", "SDUI source")
	entry := flag.String("entry", "page", "UI root")
	pkg := flag.String("package", "model", "Generated Go package")
	output := flag.String("output", "", "Output directory")
	flag.Parse()
	if *actions == "" || *source == "" || *output == "" || flag.NArg() != 0 {
		fatal(fmt.Errorf("require -actions, -ui, -output"))
	}
	target, e := filepath.Abs(*output)
	check(e)
	read := func(name string) string {
		path, e := filepath.Abs(name)
		check(e)
		rel, e := filepath.Rel(target, path)
		check(e)
		if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			fatal(fmt.Errorf("source must be outside output"))
		}
		f, e := os.Open(path)
		check(e)
		defer f.Close()
		b, e := io.ReadAll(io.LimitReader(f, parser.MaxBytes+1))
		check(e)
		return string(b)
	}
	b, e := codegen.Bundle(read(*actions), read(*source), *entry, *pkg)
	check(e)
	check(b.Publish(target))
	fmt.Println(filepath.Join(target, "manifest.json"))
}
func fatal(e error) { fmt.Fprintln(os.Stderr, e); os.Exit(2) }
func check(e error) {
	if e != nil {
		fatal(e)
	}
}
