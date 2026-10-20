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
	b.sourceProvenance(v)
	if r != nil {
		b.Manifest.Renderer = r.Identity()
	}
	text := "# " + q.Viewpoint + " — " + spec.Title + "\n\nRevision: `" + v.Revision + "`.\n\n" + spec.Note + "\n\n"
	j, _ := json.Marshal(q.Canonical())
	if q.Projection != "overview" {
		text += "Selection: `" + string(j) + "`.\n\n"
	}
	b.Files["selection.json"] = append(j, '\n')
	if len(s.Diagrams) == 0 && !(q.Projection == "overview" && (q.Viewpoint == "VP05" || q.Viewpoint == "VP06" || q.Viewpoint == "VP08" || q.Viewpoint == "VP09" || q.Viewpoint == "VP10")) {
		text += "No diagrams match this selection. This does not imply the concept is impossible or absent from the system.\n\n"
	}
	for _, d := range s.Diagrams {
		b.Put("diagrams/"+d.ID+".mmd", d.Mermaid())
		if r != nil {
			svg, e := renderDiagram(ctx, r, d)
			if e != nil {
				return nil, fmt.Errorf("%s: %w", d.ID, e)
			}
			b.Files["diagrams/"+d.ID+".svg"] = svg
		}
		text += s.DiagramMarkdown(d, "diagrams/", r != nil)
	}
	if r != nil {
		text += "\nSymbol profile: SDL 1. Actor/ellipse use UML shapes; labeled SDL relations retain their language semantics. consumes is a dashed dependency; realizes is a contribution.\n\n"
	}
	if q.Projection == "overview" {
		if q.Viewpoint == "VP02" {
			text += s.ArchitectureOverviewContext()
		} else if q.Viewpoint == "VP03" {
			text += s.ResponsibilityOverviewContext()
		} else if q.Viewpoint == "VP04" {
			text += s.InterfaceOverviewContext()
		} else if q.Viewpoint == "VP05" {
			text += v.ModeOverviewContext(q.Mode)
		} else if q.Viewpoint == "VP06" {
			text += s.ActivityOverviewContext()
		} else if q.Viewpoint == "VP10" {
			text += s.EncodingOverviewContext()
		} else if q.Viewpoint == "VP09" {
			text += s.DataOverviewContext()
		} else if q.Viewpoint == "VP08" {
			text += v.ChannelOverviewContext(q.Mode)
		} else if q.Viewpoint == "VP07" {
			text += s.FeatureOverviewContext(q.Mode)
		} else {
			text += s.OverviewContext()
		}
	}
	if q.Viewpoint == "VP01" && q.Projection == "" {
		text += s.UseCaseContext()
	}
	if !(q.Projection == "overview" && (q.Viewpoint == "VP09" || q.Viewpoint == "VP10")) {
		text += s.Tables(q.Viewpoint)
	}
	b.Put("entry.md", text)
	b.Seal()
	return b, b.CheckLinks()
}
