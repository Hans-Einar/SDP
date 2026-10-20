package viewpoint

import (
	"fmt"
	"sort"
	"strings"
)

func (v *Views) TraceabilityOverviewContext() string {
	ref := func(name string) string {
		for _, d := range v.Model.Declarations {
			if d.Name.Name == name {
				return v.SpanReference(name, d.Span)
			}
		}
		return name
	}
	position := func(source string, line int) string {
		// Legacy single-file facts have no per-span source name. The checked
		// snapshot still identifies the one original file.
		if source == "" && len(v.Sources) == 1 {
			source = v.Sources[0].Path
		}
		return traceabilityPosition(source, line)
	}
	properties := eq(v.Facts, "verb", "has")
	var b strings.Builder
	b.WriteString("## Snapshot coverage\n\nCounts describe this checked source snapshot, not implementation coverage or verification evidence. This inspection index contains no combined architecture diagram.\n\n")
	fmt.Fprintf(&b, "| Inventory | Count |\n| --- | --- |\n| Declarations | %d |\n| Facts | %d |\n| Explicit properties | %d |\n| Reported model gaps | %d |\n\n", len(v.Model.Declarations), len(v.Facts), len(properties), len(v.Gaps))
	kinds, verbs := map[string]int{}, map[string]int{}
	for _, d := range v.Model.Declarations {
		kinds[d.Kind]++
	}
	for _, f := range v.Facts {
		verbs[f.S("verb")]++
	}
	traceabilityCounts(&b, "Declaration kinds", kinds)
	traceabilityCounts(&b, "Fact verbs", verbs)
	b.WriteString("## Explicit properties\n\nProperty values, completeness and implementation-status are source assertions, not proof of implemented behavior or owner acceptance. Missing properties are unspecified, with no inferred defaults.\n\n")
	if len(properties) == 0 {
		b.WriteString("No property assignments are declared.\n\n")
	}
	for _, f := range properties {
		fmt.Fprintf(&b, "- %s: %s = %s. Source fact: %s.\n", ref(f.S("subject")), f.S("property"), f.S("value"), v.FactReference(f.S("id")))
	}
	b.WriteString("\n## Reported model gaps\n\nThese are observations from the producer's current gap rules across viewpoints, not an exhaustive requirements audit. A missing modeled relation need not be a missing implementation or a required feature. An absence has no source statement; links identify declared context only. No reported gaps does not establish model completeness.\n\n")
	if len(v.Gaps) == 0 {
		b.WriteString("No model gaps are reported by the current rules.\n\n")
	}
	for _, spec := range Catalog {
		gaps := eq(v.Gaps, "viewpoint", spec.ID)
		if len(gaps) == 0 {
			continue
		}
		fmt.Fprintf(&b, "### %s — %s\n\n", spec.ID, spec.Title)
		for _, g := range gaps {
			context := "snapshot (no declared subject)"
			if v.Kinds[g.S("model_id")] != "" {
				context = ref(g.S("model_id"))
			}
			if g.S("mode") != "" {
				context += "; mode " + ref(g.S("mode"))
			}
			if g.S("feature") != "" {
				context += "; Feature " + ref(g.S("feature"))
			}
			fmt.Fprintf(&b, "- %s — %s: %s\n", g.S("code"), context, g.S("message"))
		}
		b.WriteString("\n")
	}
	b.WriteString("## Declaration register\n\nEvery declaration is retained, including isolated definitions. Source links identify exact spans in the original checked files.\n\n")
	if len(v.Model.Declarations) == 0 {
		b.WriteString("No declarations are present.\n\n")
	}
	for _, d := range v.Model.Declarations {
		fmt.Fprintf(&b, "- %s (%s) — %s.\n", v.SpanReference(d.Name.Name, d.Span), d.Kind, position(d.Span.Source, d.Span.Line))
	}
	b.WriteString("\n## Complete fact registry\n\nAll statements are preserved, including mode, ternary relation, sequence and placement arguments. Source traceability does not prove native code conformance.\n\n")
	if len(v.Facts) == 0 {
		b.WriteString("No facts are declared.\n\n")
	}
	for _, f := range v.Facts {
		source := ""
		if span, ok := f["span"].(map[string]any); ok {
			source, _ = span["source"].(string)
		}
		fmt.Fprintf(&b, "- %s — %s (%s).\n", v.FactReference(f.S("id")), f.S("text"), position(source, f.N("line")))
	}
	return b.String() + "\n"
}

func traceabilityCounts(b *strings.Builder, title string, counts map[string]int) {
	fmt.Fprintf(b, "### %s\n\n", title)
	if len(counts) == 0 {
		b.WriteString("None declared.\n\n")
		return
	}
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	b.WriteString("| Kind / verb | Count |\n| --- | --- |\n")
	for _, k := range keys {
		fmt.Fprintf(b, "| %s | %d |\n", k, counts[k])
	}
	b.WriteString("\n")
}

func traceabilityPosition(source string, line int) string {
	if source == "" {
		source = "inline source"
	}
	// Paths are context labels, never interpreted as links or Markdown markup.
	source = strings.NewReplacer("\\", "\\\\", "`", "\\`", "[", "\\[", "]", "\\]", "*", "\\*", "_", "\\_", "<", "&lt;", ">", "&gt;", "&", "&amp;").Replace(source)
	return fmt.Sprintf("%s, line %d", source, line)
}
