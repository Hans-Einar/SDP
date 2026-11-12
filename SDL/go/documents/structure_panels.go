package documents

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"sort"
)

type structurePanel struct {
	Title string
	SVG   []byte
	Facts []viewpoint.Fact
}

// Semantically partition dense figures before rendering. Every binding appears
// in exactly one Channel+Mode participant panel; declarations have a separate
// complete containment overview. This is not a sample or truncation.
func structurePanels(v *viewpoint.Views, scope, mode string, parents map[string]string) []structurePanel {
	svg, channels, names := structureSVG(v, scope, mode, parents)
	uses := []viewpoint.Fact{}
	channelSet := map[string]bool{}
	for _, ch := range channels {
		channelSet[ch] = true
	}
	for _, f := range structureUses(v) {
		if channelSet[f.S("channel")] && (mode == "" || mode == f.S("mode")) {
			uses = append(uses, f)
		}
	}
	if len(uses) <= 8 && len(names) <= 6 {
		return []structurePanel{{Title: "Declared structure and participation", SVG: svg, Facts: uses}}
	}
	panels := []structurePanel{}
	declarations := []string{}
	for n, k := range v.Kinds {
		if (k == "unit" || k == "container") && ((scope == "" && parents[n] == "") || (scope != "" && parents[n] == scope) || n == scope && !structureHasChildren(scope, parents)) {
			declarations = append(declarations, n)
		}
	}
	sort.Strings(declarations)
	for start := 0; start < len(declarations); start += 6 {
		end := min(start+6, len(declarations))
		copy := *v
		copy.Kinds = map[string]string{}
		copy.Facts = nil
		for _, n := range declarations[start:end] {
			copy.Kinds[n] = v.Kinds[n]
		}
		if scope != "" {
			copy.Kinds[scope] = v.Kinds[scope]
		}
		image, _, _ := structureSVG(&copy, scope, mode, parents)
		panels = append(panels, structurePanel{Title: fmt.Sprintf("Containment — %d–%d of %d objects", start+1, end, len(declarations)), SVG: image})
	}
	for _, ch := range channels {
		modes := map[string]bool{}
		for _, f := range uses {
			if f.S("channel") == ch {
				modes[f.S("mode")] = true
			}
		}
		if len(modes) == 0 {
			modes[mode] = true
		}
		sortedModes := []string{}
		for m := range modes {
			sortedModes = append(sortedModes, m)
		}
		sort.Strings(sortedModes)
		for _, m := range sortedModes {
			people := map[string]bool{}
			for _, f := range uses {
				if f.S("channel") == ch && f.S("mode") == m {
					people[f.S("subject")] = true
				}
			}
			sorted := []string{}
			for n := range people {
				sorted = append(sorted, n)
			}
			sort.Strings(sorted)
			total := max(1, (len(sorted)+3)/4)
			for page := 0; page < total; page++ {
				copy := *v
				copy.Kinds = map[string]string{ch: "channel"}
				copy.Facts = nil
				selected := map[string]bool{}
				for _, n := range sorted[min(page*4, len(sorted)):min((page+1)*4, len(sorted))] {
					selected[n] = true
					seen := map[string]bool{}
					for n != "" && !seen[n] {
						seen[n] = true
						copy.Kinds[n] = v.Kinds[n]
						n = parents[n]
					}
				}
				if scope != "" {
					copy.Kinds[scope] = v.Kinds[scope]
				}
				for _, f := range uses {
					if f.S("channel") == ch && f.S("mode") == m && selected[f.S("subject")] {
						copy.Facts = append(copy.Facts, f)
					}
				}
				image, _, _ := structureSVG(&copy, scope, m, parents, ch)
				label := ch + " — mode " + m
				if m == "" {
					label = ch + " — no declared participation"
				}
				if total > 1 {
					label += fmt.Sprintf(" — participants %d–%d of %d (panel %d/%d)", page*4+1, min((page+1)*4, len(sorted)), len(sorted), page+1, total)
				}
				panels = append(panels, structurePanel{Title: label, SVG: image, Facts: copy.Facts})
			}
		}
	}
	return panels
}
