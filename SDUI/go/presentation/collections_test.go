package presentation

import (
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

func TestCollectionsAreTruthfulStructuralDescriptions(t *testing.T) {
	doc, roots, err := parser.Compile(`sdui 0.3; ref: nav "missing.sdl";
leaf=<items=tree("Folders", callback=nav.Activate.@invoke) {overflow-y=scroll}>;
page=[left=leaf,right=leaf; files=list("Files"); hidden=list("Hidden") {visible=false}];`)
	if err != nil {
		t.Fatal(err)
	}
	for _, render := range []func(*parser.Instance) (string, error){
		func(n *parser.Instance) (string, error) { return Dump(n, 160) },
		func(n *parser.Instance) (string, error) { return Markdown(n, 160) },
		func(n *parser.Instance) (string, error) { return Combined(n, 160, doc) },
	} {
		out, err := render(roots["page"])
		if err != nil {
			t.Fatal(err)
		}
		for _, part := range []string{"tree", "list", "page/left/items", "page/right/items", "provider data not supplied"} {
			if !strings.Contains(out, part) {
				t.Fatal("missing "+part, out)
			}
		}
		if strings.Contains(out, "SVG placeholder") || strings.Contains(out, "SVG plassholder") {
			t.Fatal("collection mislabelled SVG")
		}
	}
	out, err := Combined(roots["page"], 160, doc)
	if err != nil || !strings.Contains(out, "sdui-source://") || !strings.Contains(out, "page/hidden") || !strings.Contains(out, "overflow-y") || !strings.Contains(out, "Activate") {
		t.Fatal("source/binding/hidden metadata missing", err, out)
	}
	// Even hidden unsupported content must reject instead of silently disappearing.
	roots["page"].Rows[2][0].Widget = "unknown"
	for _, render := range []func(*parser.Instance, int) (string, error){Dump, Markdown} {
		if out, err := render(roots["page"], 160); err == nil || out != "" {
			t.Fatal("unknown hidden kind accepted", out, err)
		}
	}
	if out, err := Combined(roots["page"], 160, doc); err == nil || out != "" {
		t.Fatal(out, err)
	}
}
