package presentation

import (
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strings"
	"testing"
)

func TestPreviewStructuralDescriptionIncludesHiddenIntent(t *testing.T) {
	doc, roots, err := parser.Compile(`sdui 0.3; ref: art "not-read"; M=markdown("<img src='not-fetched'>",description="Reading",fallback="label"); Main=[a=M;b=M {visible=false};image=svg(art.Image.@resource,label="Caption",description="Figure",fallback="reject") {visible=false}];`)
	if err != nil {
		t.Fatal(err)
	}
	root := roots["Main"]
	for _, render := range []func() (string, error){func() (string, error) { return Dump(root, 400) }, func() (string, error) { return Markdown(root, 400) }, func() (string, error) { return Combined(root, 400, doc) }} {
		out, err := render()
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range []string{"Static Markdown preview", "Static SVG preview", "Main/a", "Main/b", "Main/image", "source=art.Image.@resource", `description="Reading"`, `description="Figure"`, `fallback="label"`, `fallback="reject"`, `label="Caption"`, "resources not supplied; source intent only", "visible=false", "source="} {
			if !strings.Contains(out, want) {
				t.Errorf("missing %q: %s", want, out)
			}
		}
	}
}
