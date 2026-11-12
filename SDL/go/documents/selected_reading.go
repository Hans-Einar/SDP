package documents

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
)

func selectedReading(ctx context.Context, v *viewpoint.Views, q viewpoint.Query, r Renderer) (*Bundle, error) {
	selected, err := v.Select(q)
	if err != nil {
		return nil, err
	}
	structure, assets, err := ReadingStructure(v, q)
	if err != nil {
		return nil, err
	}
	if structure != "" {
		selected.Diagrams = nil
	}
	b := &Bundle{Files: map[string][]byte{}, Manifest: Manifest{Version: viewpoint.Version, Revision: v.Revision, Facts: selected.Facts, Gaps: selected.Gaps, Diagrams: selected.Diagrams}}
	for name, data := range assets {
		b.Files[name] = data
	}
	b.sourceProvenance(v)
	if r != nil {
		b.Manifest.Renderer = r.Identity()
	}
	for _, d := range selected.Diagrams {
		if err = ctx.Err(); err != nil {
			return nil, err
		}
		b.Put("diagrams/"+d.ID+".mmd", d.Mermaid())
		if r != nil {
			svg, e := renderDiagram(ctx, r, d)
			if e != nil {
				return nil, fmt.Errorf("%s: %w", d.ID, e)
			}
			b.Files["diagrams/"+d.ID+".svg"] = svg
		}
	}
	b.Put("entry.md", v.ReadingMarkdown(q, selected, r != nil, structure))
	j, _ := json.Marshal(q.Canonical())
	b.Files["selection.json"] = append(j, '\n')
	b.Seal()
	return b, b.CheckLinks()
}
