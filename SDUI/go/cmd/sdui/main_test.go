package main

import (
	"bytes"
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
