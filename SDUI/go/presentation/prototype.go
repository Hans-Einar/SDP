package presentation

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"html"
	"strings"
)

type PrototypeWidget struct {
	Path, Kind, Label, Value string
	Enabled                  bool
}

func PrototypeWidgets(root *parser.Instance) []PrototypeWidget {
	out := []PrototypeWidget{}
	var visit func(*parser.Instance, bool)
	visit = func(n *parser.Instance, enabled bool) {
		if !visible(n) {
			return
		}
		enabled = enabled && n.Layout["enabled"] != false
		if n.Widget == "button" || n.Widget == "input" {
			key := "label"
			if n.Widget == "input" {
				key = "text"
			}
			out = append(out, PrototypeWidget{n.Path, n.Widget, n.Argument(key), n.Argument("value"), enabled})
		}
		for _, r := range n.Regions {
			visit(r.Node, enabled)
		}
		for _, row := range n.Rows {
			for _, c := range row {
				visit(c, enabled)
			}
		}
	}
	visit(root, true)
	return out
}

// PrototypeSVG is an explicit control gallery. General SDUI geometry belongs to G2.
func PrototypeSVG(root *parser.Instance) (string, error) {
	items := PrototypeWidgets(root)
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="640" height="%d" viewBox="0 0 640 %d"><title>SDUI prototype control gallery</title><rect width="100%%" height="100%%" fill="#edf1e8"/>`, max(1, len(items))*86+40, max(1, len(items))*86+40)
	b.WriteString(`<text x="16" y="24" font-family="sans-serif" font-size="14">Static control gallery — no callbacks</text>`)
	for index, w := range items {
		color := "#20302a"
		if !w.Enabled {
			color = "#7e8880"
		}
		y := index*86 + 36
		fmt.Fprintf(&b, `<svg x="16" y="%d" width="608" height="80" data-path="%s" data-widget="%s"><title>%s</title>`, y, html.EscapeString(w.Path), w.Kind, html.EscapeString(w.Label))
		if w.Kind == "button" {
			fmt.Fprintf(&b, `<rect x="1" y="2" width="605" height="36" rx="5" fill="#f3f6f0" stroke="#9aaa9d"/><text x="303" y="26" text-anchor="middle" font-family="sans-serif" font-size="14" fill="%s">%s</text>`, color, html.EscapeString(Safe(w.Label)))
		} else {
			fmt.Fprintf(&b, `<text x="0" y="16" font-family="sans-serif" font-size="14" fill="%s">%s</text><rect x="1" y="24" width="605" height="36" rx="4" fill="white" stroke="#9aaa9d"/><text x="10" y="48" font-family="sans-serif" font-size="14" fill="%s">%s</text>`, color, html.EscapeString(Safe(w.Label)), color, html.EscapeString(Safe(w.Value)))
		}
		b.WriteString(`</svg>`)
	}
	b.WriteString("</svg>\n")
	if b.Len() > MaxCells {
		return "", diagnostic("prototype-limit", "Prototype exceeds output budget", root)
	}
	return b.String(), nil
}
func PrototypeHTML(root *parser.Instance) (string, error) {
	var b strings.Builder
	b.WriteString(`<!doctype html><html lang="nb"><meta charset="utf-8"><title>SDUI prototype</title><style>body{font:16px sans-serif;max-width:50em;margin:2em auto;background:#edf1e8;color:#20302a}section{padding:1em;background:white;margin:1em 0}button,input{font:inherit;padding:.5em;margin:.3em}button:active{transform:translateY(2px);background:#d6e4d5}input:focus,button:focus{outline:2px solid #3975b0}label{display:block}.print-value{display:none}@media print{input{display:none}.print-value{display:inline}}</style><h1>SDUI-kontroller</h1><p>Lokal prototype. Ingen SDL-callbacks eller innsending.</p><section>`)
	for i, w := range PrototypeWidgets(root) {
		disabled := ""
		if !w.Enabled {
			disabled = " disabled"
		}
		if w.Kind == "button" {
			fmt.Fprintf(&b, `<button type="button"%s>%s</button>`, disabled, html.EscapeString(w.Label))
		} else {
			fmt.Fprintf(&b, `<label for="widget-%d">%s<input id="widget-%d" value="%s"%s><span class="print-value">%s</span></label>`, i, html.EscapeString(w.Label), i, html.EscapeString(w.Value), disabled, html.EscapeString(w.Value))
		}
	}
	b.WriteString(`</section><p id="status" aria-live="polite"></p><script>document.querySelectorAll('button').forEach(b=>b.addEventListener('click',()=>document.querySelector('#status').textContent='Trykket: '+b.textContent));document.querySelectorAll('input').forEach(i=>i.addEventListener('input',()=>i.nextElementSibling.textContent=i.value));</script></html>`)
	if b.Len() > MaxCells {
		return "", diagnostic("prototype-limit", "Prototype exceeds output budget", root)
	}
	return b.String(), nil
}
