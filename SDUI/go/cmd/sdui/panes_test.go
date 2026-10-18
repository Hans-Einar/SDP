package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPaneCLIProfileAndAtomicUnsupportedExport(t *testing.T) {
	const source = `sdui 0.3; Main=[t=tabs("Workspace")[a=page("Overview")[];b=page("Notes")[]]];`
	var out, diagnostic bytes.Buffer
	if code := execute([]string{"-"}, strings.NewReader(source), &out, &diagnostic); code != 0 {
		t.Fatal(code, diagnostic.String())
	}
	var envelope map[string]any
	if err := json.Unmarshal(out.Bytes(), &envelope); err != nil || len(envelope) != 3 || envelope["astFormat"] != "sdui-ast/0.3" || !strings.Contains(out.String(), `"composition"`) {
		t.Fatal(out.String(), err)
	}
	for _, format := range []string{"dump", "markdown"} {
		out.Reset()
		diagnostic.Reset()
		if code := execute([]string{"-", "--format", format}, strings.NewReader(source), &out, &diagnostic); code != 0 || !strings.Contains(out.String(), "Static tabs") {
			t.Fatal(code, out.String(), diagnostic.String())
		}
	}
	file := filepath.Join(t.TempDir(), "existing.svg")
	if err := os.WriteFile(file, []byte("previous artifact"), 0600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	diagnostic.Reset()
	if code := execute([]string{"-", "--format", "svg", "-o", file}, strings.NewReader(source), &out, &diagnostic); code != 2 || out.Len() != 0 || !strings.Contains(diagnostic.String(), "unsupported-pane-export") || !strings.Contains(diagnostic.String(), "Main/t") {
		t.Fatal(code, out.String(), diagnostic.String())
	}
	if after, err := os.ReadFile(file); err != nil || string(after) != "previous artifact" {
		t.Fatal("failed export modified artifact", err)
	}
	out.Reset()
	diagnostic.Reset()
	if code := execute([]string{"-"}, strings.NewReader(strings.Replace(source, "0.3", "0.2", 1)), &out, &diagnostic); code == 0 || out.Len() != 0 {
		t.Fatal("pane syntax admitted as frozen0.2")
	}
}
