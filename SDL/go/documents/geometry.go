package documents

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type nodeGeometry struct {
	ID                  string
	X, Y, Width, Height float64
	LabelWidth          float64  `json:"label_width"`
	LabelHeight         float64  `json:"label_height"`
	Lines               []string `json:"label_lines"`
}
type edgeGeometry struct {
	From, To string
	Points   [][2]float64
	Anchor   *[2]float64 `json:"label_anchor"`
	Lines    []string    `json:"label_lines"`
}
type graphGeometry struct {
	Width, Height float64
	Nodes         []nodeGeometry
	Edges         []edgeGeometry
}

func renderDiagram(ctx context.Context, r Renderer, d viewpoint.Diagram) ([]byte, error) {
	if x, ok := r.(interface {
		RenderDiagram(context.Context, viewpoint.Diagram) ([]byte, error)
	}); ok {
		return x.RenderDiagram(ctx, d)
	}
	return r.Render(ctx, d.Mermaid())
}
func (m *Mmdr) RenderDiagram(ctx context.Context, d viewpoint.Diagram) ([]byte, error) {
	if d.Kind != "flowchart" {
		return m.Render(ctx, d.Mermaid())
	}
	g, e := m.geometry(ctx, d)
	if e != nil {
		return nil, e
	}
	return semanticSVG(d, g)
}
func (m *Mmdr) geometry(ctx context.Context, d viewpoint.Diagram) (graphGeometry, error) {
	var g graphGeometry
	source := d.Mermaid()
	for n, p := range d.Nodes {
		if p.(map[string]any)["kind"] == "actor" {
			old := n + "[\""
			source = strings.Replace(source, old, n+"[\"<br/><br/><br/>", 1)
		}
	}
	if len(source) > 64000 {
		return g, fmt.Errorf("diagram source limit")
	}
	dir, e := os.MkdirTemp("", "sdl-geometry-")
	if e != nil {
		return g, e
	}
	defer os.RemoveAll(dir)
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, m.Executable, "--input", "-", "--dumpLayout", filepath.Join(dir, "layout.json"))
	cmd.Stdin = strings.NewReader(source)
	var out, errors capped
	cmd.Stdout = &out
	cmd.Stderr = &errors
	if e = cmd.Run(); e != nil {
		return g, fmt.Errorf("diagram geometry: %w: %.1000s", e, errors.String())
	}
	b, e := os.ReadFile(filepath.Join(dir, "layout.json"))
	if e != nil {
		return g, e
	}
	if len(b) > 8<<20 {
		return g, fmt.Errorf("layout output limit")
	}
	if e = json.Unmarshal(b, &g); e != nil {
		return g, e
	}
	if e = validateGeometry(d, g); e != nil {
		return g, e
	}
	sort.Slice(g.Nodes, func(i, j int) bool { return g.Nodes[i].ID < g.Nodes[j].ID })
	return g, nil
}
func validateGeometry(d viewpoint.Diagram, g graphGeometry) error {
	finite := func(x float64) bool { return !math.IsNaN(x) && !math.IsInf(x, 0) && math.Abs(x) <= 32768 }
	if !(g.Width > 0 && g.Height > 0 && finite(g.Width) && finite(g.Height)) {
		return fmt.Errorf("invalid layout extent")
	}
	if len(g.Nodes) != len(d.Nodes) || len(g.Edges) != len(d.Edges) {
		return fmt.Errorf("backend lost or invented model elements")
	}
	seen := map[string]bool{}
	for _, n := range g.Nodes {
		if d.Nodes[n.ID] == nil || seen[n.ID] || !(n.Width > 0 && n.Height > 0) || !finite(n.X) || !finite(n.Y) || !finite(n.Width) || !finite(n.Height) {
			return fmt.Errorf("invalid backend node %s", n.ID)
		}
		seen[n.ID] = true
	}
	for i, e := range g.Edges {
		proof := d.Edges[i]
		if e.From != proof.Source || e.To != proof.Target || strings.Join(e.Lines, " ") != proof.Relation || len(e.Points) < 2 {
			return fmt.Errorf("backend changed edge %d", i)
		}
		for _, p := range e.Points {
			if !finite(p[0]) || !finite(p[1]) {
				return fmt.Errorf("invalid edge geometry")
			}
		}
	}
	return nil
}
