package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	for _, args := range [][]string{{"-"}, {"-", "--format", "dump"}, {"--format", "markdown", "-"}, {"-", "--format", "svg"}} {
		var out, err bytes.Buffer
		code := execute(args, strings.NewReader(`sdui 0.2; P=[button("OK")];`), &out, &err)
		if code != 0 || out.Len() == 0 {
			t.Fatal(args, code, err.String())
		}
	}
	for _, args := range [][]string{{"-", "--format", "prototype-svg"}, {"-", "--format", "prototype-html"}, {"-", "--format", "unknown"}, {"-", "--format", "dump", "--syntax-only"}, {"-", "--entry", "missing", "--format", "dump"}} {
		var out, err bytes.Buffer
		if execute(args, strings.NewReader(`sdui 0.2; P=[];`), &out, &err) != 2 || out.Len() != 0 {
			t.Fatal(args, out.String(), err.String())
		}
	}
	dir := t.TempDir()
	source := filepath.Join(dir, "ui.sdui")
	os.WriteFile(source, []byte(`sdui 0.2; P=[];`), 0600)
	var out, err bytes.Buffer
	if execute([]string{source, "-o", source}, nil, &out, &err) != 3 {
		t.Fatal("overwrote source")
	}
}

func TestProfileEnvelopesAndCollectionExport(t *testing.T) {
	for _, version := range []string{"0.2", "0.3"} {
		for _, syntax := range []bool{false, true} {
			args := []string{"-"}
			if syntax {
				args = append(args, "--syntax-only")
			}
			var out, diagnostic bytes.Buffer
			if code := execute(args, strings.NewReader("sdui "+version+"; page=[];"), &out, &diagnostic); code != 0 {
				t.Fatal(code, diagnostic.String())
			}
			var envelope map[string]any
			if err := json.Unmarshal(out.Bytes(), &envelope); err != nil {
				t.Fatal(err)
			}
			validation := "local-profile"
			if syntax {
				validation = "syntax-only"
			}
			if len(envelope) != 3 || envelope["astFormat"] != "sdui-ast/"+version || envelope["validation"] != validation || envelope["document"].(map[string]any)["profile"] != "sdui/"+version {
				t.Fatal(envelope)
			}
		}
	}
	const source = `sdui 0.3; ref: nav "not-present.sdl"; page=[items=tree("Navigation",callback=nav.Activate.@invoke) {overflow-y=scroll}];`
	for _, format := range []string{"ast", "dump", "markdown"} {
		var out, diagnostic bytes.Buffer
		if code := execute([]string{"-", "--format", format}, strings.NewReader(source), &out, &diagnostic); code != 0 || out.Len() == 0 {
			t.Fatal(format, code, diagnostic.String())
		}
	}
	file := filepath.Join(t.TempDir(), "existing.svg")
	if err := os.WriteFile(file, []byte("retain previous artifact"), 0600); err != nil {
		t.Fatal(err)
	}
	var out, diagnostic bytes.Buffer
	if code := execute([]string{"-", "--format", "svg", "-o", file}, strings.NewReader(source), &out, &diagnostic); code != 2 || out.Len() != 0 || !strings.Contains(diagnostic.String(), "unsupported-collection-export") || !strings.Contains(diagnostic.String(), "page/items") {
		t.Fatal(code, out.String(), diagnostic.String())
	}
	retained, err := os.ReadFile(file)
	if err != nil || string(retained) != "retain previous artifact" {
		t.Fatal("failed export changed previous artifact", err)
	}
}
