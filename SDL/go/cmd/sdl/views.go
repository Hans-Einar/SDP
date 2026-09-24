package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func exportViews(args []string, out, errs io.Writer) int {
	fs := flag.NewFlagSet("viewpoints", flag.ContinueOnError)
	fs.SetOutput(errs)
	output := fs.String("output", "", "Output directory")
	form := fs.String("format", "navigator", "navigator or static")
	project := fs.String("project", "local", "Registered project ID")
	selection := fs.String("viewpoint", "", "Comma-separated viewpoint IDs")
	renderer := fs.String("renderer", "", "Explicit mmdr executable")
	mono := fs.Bool("monolithic", false, "Also produce one complete Markdown report")
	// Source precedes flags for compatibility with the toolkit's original CLI.
	if len(args) < 1 {
		return 2
	}
	source := args[0]
	if e := fs.Parse(args[1:]); e != nil {
		return 2
	}
	fail := func(e error) int { fmt.Fprintln(errs, e); return 2 }
	if *output == "" || (*form != "navigator" && *form != "static") {
		return fail(fmt.Errorf("select --output and --format navigator|static"))
	}
	abs, _ := filepath.Abs(source)
	target, _ := filepath.Abs(*output)
	rel, _ := filepath.Rel(target, abs)
	if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fail(fmt.Errorf("keep source outside output"))
	}
	f, e := os.Open(source)
	if e != nil {
		return fail(e)
	}
	defer f.Close()
	data, e := io.ReadAll(io.LimitReader(f, parser.MaxBytes+1))
	if e != nil {
		return fail(e)
	}
	v, e := viewpoint.New(string(data))
	if e != nil {
		return fail(e)
	}
	o := documents.Options{Navigator: *form == "navigator", Monolithic: *mono, Project: *project}
	if *selection != "" {
		o.Viewpoints = strings.Split(*selection, ",")
	}
	if *renderer != "" {
		o.Renderer, e = documents.NewMmdr(*renderer)
		if e != nil {
			return fail(e)
		}
	}
	b, e := documents.Build(context.Background(), v, o)
	if e != nil {
		return fail(e)
	}
	if e = b.Publish(*output); e != nil {
		return fail(e)
	}
	json.NewEncoder(out).Encode(map[string]any{"revision": v.Revision, "diagrams": len(b.Manifest.Diagrams), "files": len(b.Files), "format": *form})
	return 0
}
