package documents

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
)

func Selected(ctx context.Context, v *viewpoint.Views, q viewpoint.Query, r Renderer) (*Bundle, error) {
	s, e := v.Select(q)
	if e != nil {
		return nil, e
	}
	spec, _ := viewpoint.SpecFor(q.Viewpoint)
	b := &Bundle{Files: map[string][]byte{}, Manifest: Manifest{Version: viewpoint.Version, Revision: v.Revision, Facts: s.Facts, Gaps: s.Gaps, Diagrams: s.Diagrams}}
	if r != nil {
		b.Manifest.Renderer = r.Identity()
	}
	text := "# " + q.Viewpoint + " — " + spec.Title + "\n\nRevisjon: `" + v.Revision + "`.\n\n" + spec.Note + "\n\n"
	j, _ := json.Marshal(q.Canonical())
	text += "Utvalg: `" + string(j) + "`.\n\n"
	b.Files["selection.json"] = append(j, '\n')
	if len(s.Diagrams) == 0 {
		text += "Ingen diagramtreff i dette utvalget. Det betyr ikke at begrepet er umulig eller at systemet mangler det.\n\n"
	}
	for _, d := range s.Diagrams {
		b.Put("diagrams/"+d.ID+".mmd", d.Mermaid())
		if r != nil {
			svg, e := r.Render(ctx, d.Mermaid())
			if e != nil {
				return nil, fmt.Errorf("%s: %w", d.ID, e)
			}
			b.Files["diagrams/"+d.ID+".svg"] = svg
		}
		text += s.DiagramMarkdown(d, "diagrams/", r != nil)
	}
	text += s.Tables(q.Viewpoint)
	b.Put("entry.md", text)
	b.Seal()
	return b, b.CheckLinks()
}
