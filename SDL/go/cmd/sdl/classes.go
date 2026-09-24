package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func classCommand(args []string, out, errs io.Writer) int {
	if len(args) < 2 {
		return 2
	}
	command, source := args[0], args[1]
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(errs)
	output := fs.String("output", "", "Class document bundle")
	renderer := fs.String("renderer", "", "Explicit mmdr executable")
	if fs.Parse(args[2:]) != nil {
		return 2
	}
	fail := func(e error) int { fmt.Fprintln(errs, e); return 2 }
	f, e := os.Open(source)
	if e != nil {
		return fail(e)
	}
	defer f.Close()
	data, e := io.ReadAll(io.LimitReader(f, parser.MaxBytes+1))
	if e != nil {
		return fail(e)
	}
	m, e := parser.CheckClasses(string(data))
	if e != nil {
		return fail(e)
	}
	if command == "class-check" {
		json.NewEncoder(out).Encode(map[string]any{"valid": true, "profile": "class-core/0.1", "ast": m})
		return 0
	}
	if *output == "" {
		return fail(fmt.Errorf("--output required"))
	}
	abs, _ := filepath.Abs(source)
	target, _ := filepath.Abs(*output)
	rel, _ := filepath.Rel(target, abs)
	if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fail(fmt.Errorf("keep source outside output"))
	}
	var r documents.Renderer
	if *renderer != "" {
		r, e = documents.NewMmdr(*renderer)
		if e != nil {
			return fail(e)
		}
	}
	b, e := documents.ClassBundle(context.Background(), string(data), r)
	if e != nil {
		return fail(e)
	}
	if e = b.Publish(*output); e != nil {
		return fail(e)
	}
	json.NewEncoder(out).Encode(map[string]any{"entry": *output + "/entry.md", "revision": b.Manifest.Revision})
	return 0
}
