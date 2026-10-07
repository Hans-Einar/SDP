package blueprint

import (
	"fmt"
	"html"
	"net/url"
	"strings"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
)

// Text escapes authored text as literal Markdown prose, including HTML/control
// characters. Source bytes themselves are preserved separately.
func Text(s string) string {
	s = html.EscapeString(s)
	var b strings.Builder
	for _, r := range s {
		if strings.ContainsRune("\\[]()*_`#|!~", r) || r < 32 || r == 127 {
			fmt.Fprintf(&b, "&#%d;", r)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
func sourceLink(side string, sp parser.Span) string {
	u := url.URL{Path: "sources/" + side + "/" + sp.Source}
	return fmt.Sprintf("[%s line %d](%s)", Text(side+" "+sp.Source), sp.Line, u.String())
}

// Documents renders only validated analysis facts. It does not grant readiness.
// Hyperedges use explicit fact boxes to retain n-ary relationships.
func Documents(a *Analysis, task Task) map[string][]byte {
	files := map[string][]byte{}
	put := func(name, text string) { files[name] = []byte(text) }
	index := "# Blueprint — " + Text(task.ID) + "\n\n" + Text(task.Intent) + "\n\nDiagnostic preview; not an executable assignment.\n\n" +
		"[Changes](changes.md) · [Context](context.md) · [Obligations](obligations.md) · [Evidence](evidence.md)\n\n" +
		"System: " + Text(a.System) + ".\n\nAnalysis: " + a.Identity + ".\n"
	put("index.md", index)
	changes := "# Semantic changes\n\n[Overview](index.md)\n\n"
	for _, n := range a.Nodes {
		if n.Change != "unchanged" {
			changes += "- " + Text(n.Change+" "+n.Key) + "\n"
		}
	}
	d := viewpoint.Diagram{ID: "blueprint-context", Title: "Modeled affected context", Kind: "flowchart", Nodes: map[string]any{}, Edges: []viewpoint.Edge{}}
	preserved := map[string]bool{}
	for _, o := range a.Obligations {
		if strings.HasPrefix(o.ID, "protect/") {
			for _, protection := range task.Protect {
				prefix := "protect/" + protection.ID + "/"
				if strings.HasPrefix(o.ID, prefix) {
					preserved[strings.TrimPrefix(o.ID, prefix)] = true
				}
			}
		}
	}
	for _, rule := range task.Rules {
		if rule.Now && rule.Target {
			for _, f := range a.Facts {
				if f.Statement == rule.Fact {
					preserved[f.Key] = true
				}
			}
		}
	}
	label := func(key, change string) string {
		if preserved[key] {
			return "PRESERVE " + change
		}
		if change == "unchanged" {
			return "CONTEXT"
		}
		return "CHANGE " + change
	}
	ids := map[string]string{}
	for i, n := range a.Nodes {
		id := fmt.Sprintf("n%d", i)
		ids[n.Key] = id
		d.Nodes[id] = map[string]any{"model_id": n.Name, "kind": n.Kind + " " + label(n.Key, n.Change)}
	}
	context := "# Modeled context\n\n[Overview](index.md)\n\nCHANGE is added/removed; unchanged nodes are CONTEXT, not permission to edit.\nPRESERVE obligations are listed separately. Runtime and code coverage remain UNKNOWN.\n\n"
	for i, f := range a.Facts {
		id := fmt.Sprintf("f%d", i)
		d.Nodes[id] = map[string]any{"model_id": f.Statement, "kind": label(f.Key, f.Change)}
		for _, k := range f.Endpoints {
			d.Edges = append(d.Edges, viewpoint.Edge{Source: ids[k], Target: id, Relation: "operand"})
		}
		row := "- " + Text(f.Change+" "+f.Statement)
		for _, side := range []struct {
			name  string
			spans []parser.Span
		}{{"NOW", f.Origins.Now}, {"TARGET", f.Origins.Target}} {
			for _, sp := range side.spans {
				row += " " + sourceLink(side.name, sp)
			}
		}
		context += row + "\n"
		if f.Change != "unchanged" {
			changes += row + "\n"
		}
	}
	// Keep the complete context available, with a bounded delta diagram first.
	focus := viewpoint.Diagram{ID: "blueprint-changes", Title: "Semantic delta", Kind: "flowchart", Nodes: map[string]any{}, Edges: []viewpoint.Edge{}}
	for i, f := range a.Facts {
		if f.Change != "unchanged" {
			id := fmt.Sprintf("f%d", i)
			focus.Nodes[id] = d.Nodes[id]
			for _, k := range f.Endpoints {
				focus.Nodes[ids[k]] = d.Nodes[ids[k]]
				focus.Edges = append(focus.Edges, viewpoint.Edge{Source: ids[k], Target: id, Relation: "operand"})
			}
		}
	}
	for _, n := range a.Nodes {
		if n.Change != "unchanged" {
			focus.Nodes[ids[n.Key]] = d.Nodes[ids[n.Key]]
		}
	}
	delta := focus.Mermaid()
	changes += "\n```mermaid\n" + delta + "```\n"
	put("changes.mmd", delta)
	mermaid := d.Mermaid()
	context += "\n```mermaid\n" + mermaid + "```\n\n## Cuts and exclusions\n\n"
	for _, b := range append(append([]Boundary{}, a.Frontier...), a.Excluded...) {
		context += "- " + Text(b.Key+": "+b.Reason) + "\n"
	}
	put("context.md", context)
	put("changes.md", changes)
	put("context.mmd", mermaid)
	obligations := "# Structural obligations\n\n[Overview](index.md)\n\n"
	for _, o := range a.Obligations {
		obligations += fmt.Sprintf("- %s: **%s** (NOW present=%t, TARGET present=%t)\n", Text(o.ID), o.Result, o.NowPresent, o.TargetPresent)
	}
	put("obligations.md", obligations)
	evidence := "# Evidence and unknowns\n\n[Overview](index.md)\n\nSources passed structural validation. Behavioral checks: not run.\nNo implementation or independent assignment acceptance is inferred.\nSDUI semantic analysis is unsupported; captured SDUI bytes remain available.\n\n"
	for _, u := range a.Unknowns {
		evidence += "- " + Text(u.Side+" "+u.Code+" "+u.Subject+": "+u.Detail) + "\n"
	}
	put("evidence.md", evidence)
	return files
}
