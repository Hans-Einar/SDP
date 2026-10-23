package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInteractionCLIAtomicExportAndTemplateBoundary(t *testing.T) {
	source := `sdui 0.3; Part=<button(command="/c")>; Main=[c=command("C",toggle=true);buttons=Part;d=dialog("D")[]];`
	var out, errors bytes.Buffer
	if code := execute([]string{"-"}, strings.NewReader(source), &out, &errors); code != 0 || !strings.Contains(out.String(), `"sdui-ast/0.3"`) {
		t.Fatal(code, out.String(), errors.String())
	}
	for _, format := range []string{"dump", "markdown"} {
		out.Reset()
		errors.Reset()
		if code := execute([]string{"-", "--format", format, "--entry", "Main"}, strings.NewReader(source), &out, &errors); code != 0 || !strings.Contains(out.String(), "Static command") {
			t.Fatal(code, out.String(), errors.String())
		}
	}
	output := filepath.Join(t.TempDir(), "out.svg")
	if err := os.WriteFile(output, []byte("keep previous"), 0600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errors.Reset()
	if code := execute([]string{"-", "--format", "svg", "--entry", "Main", "-o", output}, strings.NewReader(source), &out, &errors); code == 0 || !strings.Contains(errors.String(), "unsupported-interaction-export") {
		t.Fatal(code, errors.String())
	}
	data, _ := os.ReadFile(output)
	if string(data) != "keep previous" {
		t.Fatal("failed export replaced prior artifact")
	}
	out.Reset()
	errors.Reset()
	bad := `sdui 0.3; Main=[button(command="missing")];`
	if code := execute([]string{"-", "--format", "dump", "--entry", "Main"}, strings.NewReader(bad), &out, &errors); code == 0 || out.Len() != 0 || !strings.Contains(errors.String(), "interaction-reference") {
		t.Fatal(code, out.String(), errors.String())
	}
}
