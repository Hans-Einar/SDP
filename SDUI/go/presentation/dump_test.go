package presentation

import (
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"os"
	"strings"
	"testing"
)

func TestConcept1Snapshots(t *testing.T) {
	source, e := os.ReadFile("../../examples/concept1-bucking.sdui")
	if e != nil {
		t.Fatal(e)
	}
	_, roots, e := parser.Compile(string(source))
	if e != nil {
		t.Fatal(e)
	}
	for _, format := range []string{"txt", "md"} {
		var text string
		if format == "txt" {
			text, e = Dump(roots["bucking"], 160)
		} else {
			text, e = Markdown(roots["bucking"], 160)
		}
		if e != nil {
			t.Fatal(e)
		}
		want, e := os.ReadFile("../../examples/concept1-bucking.dump." + format)
		if e != nil {
			t.Fatal(e)
		}
		if text != string(want) {
			t.Fatalf("Concept1 %s differs from Python snapshot", format)
		}
	}
}
func TestBoundedRawMarkdown(t *testing.T) {
	_, roots, e := parser.Compile(`sdui 0.2; P=["\u001b[31m红色 é"; "\u202eabc"; "` + "```mermaid\\nSECRET-->B\\n```" + `"];`)
	if e != nil {
		t.Fatal(e)
	}
	s, e := Dump(roots["P"], 80)
	if e != nil {
		t.Fatal(e)
	}
	if strings.ContainsAny(s, "\x1b\u202e") || strings.Contains(s, "SECRET") {
		t.Fatal(s)
	}
	for _, line := range strings.Split(s, "\n")[1:] {
		if CellWidth(line) > 80 {
			t.Fatal(line)
		}
	}
	if _, e = Dump(roots["P"], 0); e == nil {
		t.Fatal("accepted zero width")
	}
	if codeSpan("`x``y`") != "``` `x``y` ```" {
		t.Fatal("Uncontained code span")
	}
	if strings.Join(contentLines("```python\nx=1"), "\n") != "```python\nx=1\n```" {
		t.Fatal("Unclosed fence")
	}
}
