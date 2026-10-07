package presentation

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strings"
	"testing"
)

func TestCombinedProvenanceAndText(t *testing.T) {
	source := `sdui 0.2;
ref: model "unopened.design";
Buttons=<ok=button("Å save", callback=model.save.@callback)>;
page=[header="Hello"; a=Buttons, b=Buttons {visible=false}];`
	d, roots, e := parser.Compile(source)
	if e != nil {
		t.Fatal(e)
	}
	root := roots["page"]
	out, e := Combined(root, 160, d)
	if e != nil {
		t.Fatal(e)
	}
	again, _ := Combined(root, 160, d)
	if again != out {
		t.Fatal("nondeterministic")
	}
	text, e := Markdown(root, 160)
	if e != nil || !strings.HasSuffix(out, text) {
		t.Fatal("text exporter changed")
	}
	for _, want := range []string{"```mermaid", "mindmap", "Row 1", "hidden", "reuse Buttons", "page/b/ok", "module model", "unopened.design", "sdui-source://"} {
		if !strings.Contains(out, want) {
			t.Fatal(want, out)
		}
	}
	a := root.Rows[0][0]
	b := root.Rows[0][1]
	if a.Span != b.Span || len(a.Uses) != 1 || a.Uses[0].Span == b.Uses[0].Span {
		t.Fatal("reuse provenance lost")
	}
	span := a.Uses[0].Span
	if source[span.Start:span.End] != "a=Buttons" {
		t.Fatal(source[span.Start:span.End])
	}
	widget := a.Rows[0][0]
	if !strings.Contains(out, fmt.Sprintf("sdui-source://%d/%d", widget.Span.Start, widget.Span.End)) {
		t.Fatal("widget span absent")
	}
}
