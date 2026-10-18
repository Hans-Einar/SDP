package viewpoint

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"sort"
	"strings"
)

// SourceBindings is supplied only by a filesystem-aware host after it has bound
// a complete checked source snapshot to original paths and retained source bytes.
func (v *Views) SpanReference(label string, span parser.Span) string {
	id := v.SourceBindings[span.Source]
	if id == "" || span.Start < 0 || span.End <= span.Start {
		return label
	}
	return fmt.Sprintf("[%s](sdl-source://%s/%d/%d)", label, id, span.Start, span.End)
}
func (v *Views) FactReference(id string) string {
	for _, f := range v.Facts {
		if f.S("id") == id {
			if data, ok := f["span"].(map[string]any); ok {
				start, _ := data["start"].(int)
				end, _ := data["end"].(int)
				source, _ := data["source"].(string)
				return v.SpanReference(id, parser.Span{Source: source, Start: start, End: end})
			}
		}
	}
	return id
}
func (v *Views) FactReferences(ids []string) string {
	labels := make([]string, 0, len(ids))
	for _, id := range ids {
		labels = append(labels, v.FactReference(id))
	}
	return strings.Join(labels, ", ")
}
func (v *Views) DeclarationReferences(d Diagram) string {
	if v.Model == nil {
		return ""
	}
	names := map[string]bool{}
	for _, value := range d.Nodes {
		if m, ok := value.(map[string]any); ok {
			if n, ok := m["model_id"].(string); ok {
				names[n] = true
			}
		}
	}
	labels := []string{}
	for _, decl := range v.Model.Declarations {
		if names[decl.Name.Name] {
			labels = append(labels, v.SpanReference(decl.Name.Name, decl.Span))
		}
	}
	sort.Strings(labels)
	return strings.Join(labels, ", ")
}
