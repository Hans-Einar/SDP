package documents

import (
	"context"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
)

func ClassBundle(ctx context.Context, source string, r Renderer) (*Bundle, error) {
	m, e := parser.CheckClasses(source)
	if e != nil {
		return nil, e
	}
	d, facts, e := viewpoint.Classes(m)
	if e != nil {
		return nil, e
	}
	b := &Bundle{Files: map[string][]byte{}, Manifest: Manifest{Version: viewpoint.Version + "/class-core-0.1", Revision: Hash([]byte(source)), Facts: facts, Diagrams: []viewpoint.Diagram{d}}}
	if r != nil {
		b.Manifest.Renderer = r.Identity()
		svg, e := r.Render(ctx, d.Mermaid())
		if e != nil {
			return nil, e
		}
		b.Files["diagrams/"+d.ID+".svg"] = svg
	}
	b.Put("diagrams/"+d.ID+".mmd", d.Mermaid())
	v := &viewpoint.Views{}
	text := "# Class design — class-core 0.1\n\nExplicit classes/roles/multiplicities; structural relations are not reinterpreted.\n\n" + v.DiagramMarkdown(d, "diagrams/", r != nil) + "| ID | Statement | Line |\n| --- | --- | --- |\n"
	for _, f := range facts {
		text += fmt.Sprintf("| %s | %s | %d |\n", f.S("id"), f.S("text"), f.N("line"))
	}
	b.Put("entry.md", text)
	b.Put("index.md", "# Class design (A4)\n\n[Diagram and source facts](entry.md)\n")
	b.Seal()
	return b, b.CheckLinks()
}
