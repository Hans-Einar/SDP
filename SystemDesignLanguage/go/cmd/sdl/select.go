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

func selectView(args []string, out, errs io.Writer) int {
	fs := flag.NewFlagSet("view", flag.ContinueOnError)
	fs.SetOutput(errs)
	uri := fs.String("uri", "", "Typed sdl-view URI")
	output := fs.String("output", "", "Complete selected bundle directory")
	renderer := fs.String("renderer", "", "Registered mmdr program")
	if len(args) == 0 {
		return 2
	}
	source := args[0]
	if fs.Parse(args[1:]) != nil {
		return 2
	}
	fail := func(e error) int { fmt.Fprintln(errs, e); return 2 }
	if *output == "" {
		return fail(fmt.Errorf("--output is required"))
	}
	abs, _ := filepath.Abs(source)
	target, _ := filepath.Abs(*output)
	rel, _ := filepath.Rel(target, abs)
	if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fail(fmt.Errorf("keep source outside output"))
	}
	s, e := viewpoint.ParseURI(*uri)
	if e != nil {
		return fail(e)
	}
	f, e := os.Open(source)
	if e != nil {
		return fail(e)
	}
	defer f.Close()
	bytes, e := io.ReadAll(io.LimitReader(f, parser.MaxBytes+1))
	if e != nil {
		return fail(e)
	}
	v, e := viewpoint.New(string(bytes))
	if e != nil {
		return fail(e)
	}
	var r documents.Renderer
	if *renderer != "" {
		r, e = documents.NewMmdr(*renderer)
		if e != nil {
			return fail(e)
		}
	}
	b, e := documents.Selected(context.Background(), v, s.Query, r)
	if e != nil {
		return fail(e)
	}
	if e = b.Publish(*output); e != nil {
		return fail(e)
	}
	json.NewEncoder(out).Encode(map[string]any{"revision": v.Revision, "entry": *output + "/entry.md"})
	return 0
}
