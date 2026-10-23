package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScalarCLILexemesAndAtomicUnsupportedSVG(t *testing.T) {
	source := `sdui 0.3; Main=[n=number("Amount",min=-0,max=1e1,step=1.00e-1,value=0.30)];`
	var out, errors bytes.Buffer
	if code := execute([]string{"-"}, strings.NewReader(source), &out, &errors); code != 0 || !strings.Contains(out.String(), `"kind": "number-lexeme"`) || !strings.Contains(out.String(), `"value": "0.30"`) {
		t.Fatal(code, out.String(), errors.String())
	}
	for _, format := range []string{"dump", "markdown"} {
		out.Reset()
		errors.Reset()
		if code := execute([]string{"-", "--format", format, "--entry", "Main"}, strings.NewReader(source), &out, &errors); code != 0 || !strings.Contains(out.String(), "value=0.30") {
			t.Fatal(code, out.String(), errors.String())
		}
	}
	output := filepath.Join(t.TempDir(), "prior.svg")
	if err := os.WriteFile(output, []byte("old artifact"), 0600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errors.Reset()
	if code := execute([]string{"-", "--format", "svg", "--entry", "Main", "-o", output}, strings.NewReader(source), &out, &errors); code == 0 || out.Len() != 0 || !strings.Contains(errors.String(), "unsupported-value-export") {
		t.Fatal(code, out.String(), errors.String())
	}
	bytes, _ := os.ReadFile(output)
	if string(bytes) != "old artifact" {
		t.Fatal("unsupported export overwrote old output")
	}
}
