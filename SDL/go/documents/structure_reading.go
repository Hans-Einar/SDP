package documents

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"html"
	"net/url"
	"sort"
	"strings"
)

// ReadingStructure draws channel attachments at the declared participant boundary.
// Collapsed containment is a projection only; it never creates ports or traffic.
func ReadingStructure(v *viewpoint.Views, q viewpoint.Query) (string, map[string][]byte, error) {
	assets := map[string][]byte{}
	if q.Viewpoint != "VP02" && q.Viewpoint != "VP08" {
		return "", assets, nil
	}
	if q.Register != "" || q.Collection != "" {
		return "", assets, nil
	}
	parents := map[string]string{}
	for _, f := range v.Relations["contains"] {
		if v.Kinds[f.S("subject")] == "container" || v.Kinds[f.S("subject")] == "unit" {
			parents[f.S("object")] = f.S("subject")
		}
	}
	scopes := []string{}
	if q.Focus != "" {
		if k := v.Kinds[q.Focus]; k == "container" || k == "unit" || k == "channel" {
			scopes = []string{q.Focus}
		} else {
			return "", assets, nil
		}
	} else {
		scopes = append(scopes, "")
		if q.Viewpoint == "VP02" {
			for name, kind := range v.Kinds {
				if kind == "container" || (kind == "unit" && structureHasChildren(name, parents)) {
					scopes = append(scopes, name)
				}
			}
			sort.Strings(scopes[1:])
			if len(scopes) > 25 {
				scopes = scopes[:1]
			}
		}
	}
	var text strings.Builder
	if q.Focus == "" && q.Viewpoint == "VP02" && len(scopes) == 1 {
		count := 0
		for name, kind := range v.Kinds {
			if kind == "container" || (kind == "unit" && structureHasChildren(name, parents)) {
				count++
			}
		}
		if count > 24 {
			fmt.Fprintf(&text, "\nThis model has %d containers or units with interiors. The complete topology is below; open the Containers or Units collection and an individual object for its interior drawing and chapters.\n", count)
		}
	}
	for i, scope := range scopes {
		name := fmt.Sprintf("structure-%d.svg", i)
		_, channels, participants := structureSVG(v, scope, q.Mode, parents)
		panels := structurePanels(v, scope, q.Mode, parents)
		_ = name
		title := "Structure and collaboration"
		if scope != "" {
			title = scope + " — participants and channels"
		}
		fmt.Fprintf(&text, "\n## %s\n\n", title)
		for j, panel := range panels {
			asset := fmt.Sprintf("structure-%d-%d.svg", i, j)
			assets[asset] = panel.SVG
			fmt.Fprintf(&text, "### %s\n\n![%s](%s)\n\n", panel.Title, panel.Title, asset)
		}
		text.WriteString("Containment is logical structure. Named attachments mark existing Channels at participant boundaries; they are not declared ports. Arrows describe sender/receiver roles for the stated payload and mode, not execution order.\n\n")
		if q.Mode == "" {
			text.WriteString("All explicitly declared modes are shown; no mode inheritance is inferred.\n\n")
		} else {
			fmt.Fprintf(&text, "Selected mode: **%s**.\n\n", q.Mode)
		}
		if len(channels) == 0 {
			text.WriteString("No channel participation is declared in this scope.\n\n")
		}
		if v.ReadingProject != "" {
			text.WriteString("Inspect: ")
			for j, p := range participants {
				if j > 0 {
					text.WriteString(" · ")
				}
				u := "sdl-view://" + v.ReadingProject + "/VP02?projection=reading&focus=" + url.QueryEscape(p)
				if q.Mode != "" {
					u += "&mode=" + url.QueryEscape(q.Mode)
				}
				fmt.Fprintf(&text, "[%s](%s)", p, u)
			}
			text.WriteString("\n\n")
		}
		for _, ch := range channels {
			fmt.Fprintf(&text, "### %s payloads\n\n", ch)
			text.WriteString("| Participant | Role | Message / Datagram | Dataset origin | Mode |\n| --- | --- | --- | --- | --- |\n")
			for _, f := range structureUses(v) {
				if f.S("channel") != ch || (q.Mode != "" && f.S("mode") != q.Mode) {
					continue
				}
				msg := f.S("message")
				dataset := "Not applicable"
				if v.Kinds[msg] == "datagram" {
					dataset = "Not declared"
					for _, from := range v.Relations["from"] {
						if from.S("subject") == msg {
							dataset = from.S("object")
						}
					}
				}
				fmt.Fprintf(&text, "| %s | %s | %s | %s | %s |\n", f.S("subject"), f.S("role"), msg, dataset, f.S("mode"))
			}
			text.WriteString("\n")
		}
	}
	return text.String(), assets, nil
}

func structureSVG(v *viewpoint.Views, scope, mode string, parents map[string]string, forcedChannels ...string) ([]byte, []string, []string) {
	nodes, inside, channels := map[string]bool{}, map[string]bool{}, map[string]bool{}
	for _, ch := range forcedChannels {
		channels[ch] = true
	}
	nested := structureHasChildren(scope, parents)
	if v.Kinds[scope] == "channel" {
		channels[scope] = true
	}
	for name, kind := range v.Kinds {
		if scope == "" {
			if kind == "channel" {
				channels[name] = true
			}
			if (kind == "container" || kind == "unit") && parents[name] == "" {
				nodes[name] = true
			}
		}
		if scope != "" && (parents[name] == scope || (name == scope && !nested && (kind == "unit" || kind == "container"))) {
			nodes[name] = true
			inside[name] = true
		}
	}
	for _, f := range structureUses(v) {
		if mode != "" && f.S("mode") != mode {
			continue
		}
		if scope == "" || structureDescendant(f.S("subject"), scope, parents) || f.S("channel") == scope {
			channels[f.S("channel")] = true
		}
	}
	// One attachment per actual participant/role/mode. Payload rows remain complete
	// in the adjacent channel chapter instead of repeating parallel edges per message.
	groups := map[string]viewpoint.Fact{}
	payloads := map[string]map[string]bool{}
	for _, f := range structureUses(v) {
		if !channels[f.S("channel")] || (mode != "" && f.S("mode") != mode) {
			continue
		}
		p := structureVisible(f.S("subject"), scope, parents)
		if p != scope || !nested {
			nodes[p] = true
		}
		key := strings.Join([]string{f.S("subject"), f.S("channel"), f.S("role"), f.S("mode")}, "/")
		groups[key] = f
		if payloads[key] == nil {
			payloads[key] = map[string]bool{}
		}
		payloads[key][f.S("message")] = true
	}
	names := []string{}
	for n := range nodes {
		names = append(names, n)
	}
	sort.Slice(names, func(i, j int) bool {
		if inside[names[i]] != inside[names[j]] {
			return inside[names[i]]
		}
		return names[i] < names[j]
	})
	chs := []string{}
	for n := range channels {
		chs = append(chs, n)
	}
	sort.Strings(chs)
	keys := []string{}
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	boxH, laneH := 80, 60
	for _, n := range names {
		boxH = max(boxH, 28+16*svgLineCount(structureKind(v.Kinds[n], parents[n]), 29)+18*svgLineCount(n, 23))
	}
	for _, n := range chs {
		boxH = max(boxH, 45+18*svgLineCount(n, 17))
	}
	for _, k := range keys {
		f := groups[k]
		laneH = max(laneH, 20+16*(svgLineCount(f.S("subject")+" · "+f.S("role"), 34)+svgLineCount(f.S("mode"), 34)+svgLineCount(structurePayloadSummary(payloads[k]), 34)))
	}
	height := 150 + max(len(names)*(boxH+28), max(len(chs)*(boxH+40), len(keys)*laneH))
	if height < 250 {
		height = 250
	}
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="800" height="%d" viewBox="0 0 800 %d"><defs><marker id="arrow" markerWidth="8" markerHeight="8" refX="7" refY="4" orient="auto-start-reverse"><path d="M0,0 L8,4 L0,8 Z" fill="#42627d"/></marker></defs><rect width="800" height="100%%" fill="#fff"/>`, height, height)
	step := (height - 150) / max(1, len(names))
	top := 70
	if scope != "" && nested {
		top = max(top, 35+(svgLineCount(scope, 24)+1)*18)
	}
	if scope != "" && nested {
		h := top - 20 + lenInside(inside)*step
		fmt.Fprintf(&b, `<rect x="12" y="12" width="250" height="%d" rx="12" fill="#f2f7fb" stroke="#8aa3b7"/>`, h)
		svgText(&b, 24, 34, 24, 16, scope, "#253d52")
	}
	ys := map[string]int{}
	if nested && scope != "" {
		ys[scope] = 45
	}
	for i, n := range names {
		y := top + i*step
		ys[n] = y + boxH/2
		fill := "#eef4fa"
		if scope != "" && !inside[n] {
			fill = "#faf2e7"
		}
		fmt.Fprintf(&b, `<rect x="20" y="%d" width="230" height="%d" rx="5" fill="%s" stroke="#56758e"/>`, y, boxH, fill)
		end := svgText(&b, 32, y+19, 29, 14, structureKind(v.Kinds[n], parents[n]), "#5c6d7a")
		svgText(&b, 32, end+20, 23, 16, n, "#142d43")
	}
	cy := map[string]int{}
	for i, ch := range chs {
		y := top + i*(height-150)/max(1, len(chs))
		cy[ch] = y + boxH/2
		fmt.Fprintf(&b, `<rect x="600" y="%d" width="180" height="%d" rx="7" fill="#eaf6ef" stroke="#477d62"/>`, y, boxH)
		svgText(&b, 610, y+19, 20, 14, "Channel", "#476956")
		svgText(&b, 610, y+41, 17, 16, ch, "#143c29")
	}
	for i, k := range keys {
		f := groups[k]
		p := structureVisible(f.S("subject"), scope, parents)
		y, end := ys[p], cy[f.S("channel")]
		startX := 250
		if nested && p == scope {
			startX = 262
		}
		lane := top + laneH/2 + i*(height-150)/max(1, len(keys))
		arrow := `marker-end="url(#arrow)"`
		if f.S("role") == "receiver" {
			arrow = `marker-start="url(#arrow)"`
		}
		if f.S("role") == "uses (aggregated)" {
			arrow = ""
		}
		fmt.Fprintf(&b, `<path d="M%d,%d L265,%d L280,%d L580,%d L590,%d L600,%d" fill="none" stroke="#42627d" stroke-width="1.5" %s/><rect x="%d" y="%d" width="14" height="14" rx="2" fill="#fff" stroke="#42627d"/><rect x="280" y="%d" width="310" height="%d" rx="3" fill="#fff"/>`, startX, y, y, lane, lane, end, end, arrow, startX-7, y-7, lane-laneH/2, laneH)
		labelEnd := svgText(&b, 288, lane-laneH/2+15, 34, 14, f.S("subject")+" · "+f.S("role"), "#203c52")
		labelEnd = svgText(&b, 288, labelEnd+16, 34, 14, f.S("mode"), "#526b7d")
		summary := structurePayloadSummary(payloads[k])
		svgText(&b, 288, labelEnd+16, 34, 14, summary, "#526b7d")
	}
	if len(names) == 0 {
		svgText(&b, 30, 65, 70, 16, "No container or unit participation is declared in this scope.", "#253d52")
	}
	b.WriteString("</svg>")
	return []byte(b.String()), chs, names
}
func structureDescendant(name, scope string, parents map[string]string) bool {
	seen := map[string]bool{}
	for name != "" && !seen[name] {
		if name == scope {
			return true
		}
		seen[name] = true
		name = parents[name]
	}
	return false
}
func structureVisible(name, scope string, parents map[string]string) string {
	if name == scope {
		return name
	}
	seen := map[string]bool{}
	for parents[name] != "" && !seen[name] {
		if parents[name] == scope {
			return name
		}
		if scope != "" && !structureDescendant(name, scope, parents) {
			return name
		}
		seen[name] = true
		name = parents[name]
	}
	return name
}
func structureHasChildren(name string, parents map[string]string) bool {
	if name == "" {
		return false
	}
	for _, parent := range parents {
		if parent == name {
			return true
		}
	}
	return false
}
func lenInside(m map[string]bool) int {
	n := 0
	for _, b := range m {
		if b {
			n++
		}
	}
	return n
}
func structureUses(v *viewpoint.Views) []viewpoint.Fact {
	out := []viewpoint.Fact{}
	for _, f := range v.Facts {
		if f.S("verb") == "uses" {
			out = append(out, f)
		}
	}
	return out
}
func structureKind(kind, parent string) string {
	if parent != "" {
		return kind + " — contained by " + parent
	}
	return kind
}
func svgLineCount(s string, width int) int { return max(1, (len([]rune(s))+width-1)/width) }
func svgText(b *strings.Builder, x, y, width, size int, s, color string) int {
	r := []rune(s)
	if len(r) == 0 {
		return y
	}
	for len(r) > 0 {
		n := min(width, len(r))
		fmt.Fprintf(b, `<text x="%d" y="%d" font-family="sans-serif" font-size="%d" fill="%s">%s</text>`, x, y, size, color, html.EscapeString(string(r[:n])))
		r = r[n:]
		if len(r) > 0 {
			y += size + 2
		}
	}
	return y
}

func structurePayloadSummary(p map[string]bool) string {
	if len(p) == 1 {
		for n := range p {
			return n
		}
	}
	return fmt.Sprintf("%d declared payloads — see table", len(p))
}
