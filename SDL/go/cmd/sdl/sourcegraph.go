package main

import (
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/sourcegraph"
	"io"
	"os"
	"path/filepath"
)

func graphCommand(args []string, in io.Reader, out, errs io.Writer) int {
	encode := func(v any) int {
		e := json.NewEncoder(out)
		e.SetIndent("", "  ")
		if e.Encode(v) != nil {
			return 2
		}
		return 0
	}
	fail := func(e error) int {
		var detail any = e.Error()
		switch d := e.(type) {
		case sourcegraph.Issues:
			detail = d
		case sourcegraph.Issue:
			detail = d
		case parser.Diagnostic:
			detail = parser.Data(d)
		}
		encode(map[string]any{"valid": false, "diagnostics": detail})
		return 1
	}
	if len(args) < 2 || len(args) > 3 {
		return 2
	}
	name := args[1]
	cache := &parser.SyntaxCache{}
	var snap *sourcegraph.Snapshot
	if args[0] == "fragment" || name == "-" {
		if name != "-" {
			f, e := os.Open(name)
			if e != nil {
				return fail(e)
			}
			defer f.Close()
			st, e := f.Stat()
			if e != nil {
				return fail(e)
			}
			if !st.Mode().IsRegular() {
				return fail(fmt.Errorf("expected regular source"))
			}
			in = f
		}
		b, e := io.ReadAll(io.LimitReader(in, parser.MaxBytes+1))
		if e != nil {
			return fail(e)
		}
		file, e := cache.Parse(filepath.Base(name), string(b))
		if e != nil {
			return fail(e)
		}
		if args[0] == "fragment" {
			return encode(file.Inspect())
		}
		var ds sourcegraph.Issues
		snap, ds = sourcegraph.Compile(file.Name(), []*parser.File{file}, args[0] != "format")
		if len(ds) > 0 {
			return fail(ds)
		}
	} else {
		loaded, e := sourcegraph.Load(name, cache, args[0] != "format")
		if e != nil {
			return fail(e)
		}
		snap = loaded.Snapshot
	}
	if args[0] == "format" {
		if len(args) == 3 && args[2] == "--file-map" {
			return encode(map[string]any{"schema": "sdl-formatted-files/1", "files": snap.Format()})
		}
		if snap.FileCount() != 1 {
			return fail(fmt.Errorf("multi-file format requires --file-map"))
		}
		for _, s := range snap.Format() {
			_, e := io.WriteString(out, s)
			if e != nil {
				return 2
			}
		}
		return 0
	}
	if args[0] == "ast" {
		return encode(snap.AST())
	}
	return encode(map[string]any{"valid": true, "profile": snap.Profile(), "system": snap.System(), "revision": snap.Revision(), "sources": snap.Sources(), "warnings": snap.Warnings()})
}
