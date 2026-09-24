package documents

import (
	"bytes"
	"context"
	"encoding/xml"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"io"
	"math"
	"os"
	"strings"
	"testing"
)

func TestShapeAndArrowProfile(t *testing.T) {
	v := model(t)
	var d viewpoint.Diagram
	for _, x := range v.Diagrams {
		if x.ID == "VP01-BrowseDesignViews" {
			d = x
		}
	}
	g := graphGeometry{Width: 800, Height: 500}
	for id := range d.Nodes {
		g.Nodes = append(g.Nodes, nodeGeometry{ID: id, X: 10, Y: 10, Width: 150, Height: 130, Lines: []string{"Name"}})
	}
	for _, e := range d.Edges {
		g.Edges = append(g.Edges, edgeGeometry{From: e.Source, To: e.Target, Points: [][2]float64{{0, 0}, {10, 10}}, Lines: []string{e.Relation}})
	}
	if e := validateGeometry(d, g); e != nil {
		t.Fatal(e)
	}
	svg, e := semanticSVG(d, g)
	if e != nil {
		t.Fatal(e)
	}
	for _, mark := range []string{`data-symbol="actor-head"`, `data-symbol="actor-body"`, `<ellipse data-symbol="usecase"`, `data-symbol="feature-tab"`, `data-fact="`} {
		if !strings.Contains(string(svg), mark) {
			t.Fatal(mark)
		}
	}
	g.Nodes = g.Nodes[:1]
	if validateGeometry(d, g) == nil {
		t.Fatal("backend silently lost nodes")
	}
	n := nodeGeometry{X: 0, Y: 0, Width: 200, Height: 100}
	p := ellipseBoundary(n, [2]float64{0, 0})
	if math.Abs((p[0]-100)*(p[0]-100)/10000+(p[1]-50)*(p[1]-50)/2500-1) > 1e-9 {
		t.Fatal("arrow misses ellipse")
	}
	d = viewpoint.Diagram{Title: "Interface dependency", Nodes: map[string]any{"n_A": map[string]any{"model_id": "A", "kind": "unit"}, "n_B": map[string]any{"model_id": "B", "kind": "interface"}}, Edges: []viewpoint.Edge{{Source: "n_A", Target: "n_B", Relation: "consumes", Fact: "f0"}}}
	g = graphGeometry{Width: 300, Height: 100, Nodes: []nodeGeometry{{ID: "n_A", Width: 60, Height: 50}, {ID: "n_B", X: 200, Width: 60, Height: 50}}, Edges: []edgeGeometry{{From: "n_A", To: "n_B", Points: [][2]float64{{60, 25}, {200, 25}}}}}
	svg, e = semanticSVG(d, g)
	if e != nil || !strings.Contains(string(svg), `marker-end="url(#uml-dependency)" stroke-dasharray="6 4"`) {
		t.Fatal("wrong UML dependency marker", e)
	}
	if strings.Contains(string(svg), "composition") {
		t.Fatal("invented composition")
	}
}
func TestRealSemanticBackend(t *testing.T) {
	exe := os.Getenv("MMDR")
	if exe == "" {
		t.Skip("set MMDR for real backend capability verification")
	}
	m, e := NewMmdr(exe)
	if e != nil {
		t.Fatal(e)
	}
	v := model(t)
	for _, id := range []string{"VP01-BrowseDesignViews", "VP08-SelectedViewOpened"} {
		var d viewpoint.Diagram
		for _, x := range v.Diagrams {
			if x.ID == id {
				d = x
			}
		}
		b, e := m.RenderDiagram(context.Background(), d)
		if e != nil {
			t.Fatal(e)
		}
		if d.Kind == "flowchart" {
			for _, mark := range []string{`actor-head`, `<ellipse`, `data-model-id="BrowseDesignViews"`, `data-relation="pursues"`} {
				if !strings.Contains(string(b), mark) {
					t.Fatal("actual figure missing", mark)
				}
			}
		} else {
			var text strings.Builder
			decoder := xml.NewDecoder(bytes.NewReader(b))
			for {
				tok, err := decoder.Token()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatal(err)
				}
				if c, ok := tok.(xml.CharData); ok {
					text.Write(c)
					text.WriteByte(' ')
				}
			}
			if !strings.Contains(strings.Join(strings.Fields(text.String()), " "), "Select View Request") || !strings.Contains(string(b), "path") {
				t.Fatal("sequence missing content")
			}
		}
	}
}
