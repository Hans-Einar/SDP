package parser

import (
	"math"
	"strings"
)

var panes = map[string]widgetSchema{
	"tabs":  {"label", map[string]string{"label": "string", "selected": "string", "callback": "reference"}, "label"},
	"page":  {"label", map[string]string{"label": "string", "icon": "string"}, "label"},
	"split": {"axis", map[string]string{"axis": "string", "proportion": "number", "minFirst": "number", "minSecond": "number", "collapsible": "boolean"}, "axis"},
}

func validatePaneSource(n *Node, profile string) {
	_, pane := panes[val(n.Widget)]
	if n.Kind == "composition" && (profile != "sdui/0.3" || !pane) {
		fail("composition-kind", "Only 0.3 tabs/page/split accept bodies", n.Span)
	}
	if !pane {
		return
	}
	if n.Kind != "composition" {
		fail("pane-body", "Pane requires a bracketed body", n.Span)
	}
	a := widgetArguments(n, profile)
	for _, key := range []string{"label", "selected", "icon"} {
		if v, ok := a[key].(Literal); ok {
			s, ok := v.Value.(string)
			if !ok || strings.TrimSpace(s) == "" {
				fail("pane-argument", "Nonempty "+key+" required", v.Span)
			}
		}
	}
	if val(n.Widget) == "split" {
		axis := a["axis"].(Literal)
		axisValue, ok := axis.Value.(string)
		if !ok || !member(axisValue, "horizontal vertical") {
			fail("split-axis", "Expected horizontal or vertical axis", axis.Span)
		}
		if _, ok := a["collapsible"].(Literal).Value.(bool); !ok {
			fail("split-collapsible", "Boolean collapsible required", n.Span)
		}
		values := map[string]float64{}
		for _, key := range []string{"proportion", "minFirst", "minSecond"} {
			v := a[key].(Literal)
			x, ok := v.Value.(float64)
			if !ok || math.IsNaN(x) || math.IsInf(x, 0) {
				fail("split-number", "Finite numeric "+key+" required", v.Span)
			}
			values[key] = x
		}
		p, first, second := values["proportion"], values["minFirst"], values["minSecond"]
		if first < 0 || second < 0 || first+second >= 1 || p < first || p > 1-second {
			fail("split-range", "Minima must be nonnegative with sum < 1; proportion must satisfy both", n.Span)
		}
	}
	for _, row := range n.Rows {
		for _, child := range row.Items {
			if child.Role != nil {
				fail("pane-region", "Body assignments are names, not frame regions", child.Span)
			}
		}
	}
}

// PaneChild is a direct named child in declaration order. ID is local to its owner;
// Node.Path is the full normalized identity, including every reuse instance.
type PaneChild struct {
	ID   string
	Node *Instance
}

// PaneChildren validates the bounded direct shape of a normalized tabs or split.
// It does not choose a live page or resolve M2 command reference paths.
func PaneChildren(owner *Instance) ([]PaneChild, error) {
	span := Span{}
	if owner != nil {
		span = owner.Span
	}
	bad := func(message string) ([]PaneChild, error) {
		return nil, &Diagnostic{Code: "pane-children", Message: message, Span: span}
	}
	if owner == nil || owner.Profile != "sdui/0.3" || owner.Kind != "composition" || !member(owner.Widget, "tabs split") {
		return bad("Expected a 0.3 tabs or split composition")
	}
	if len(owner.Regions) != 0 || len(owner.Rows) == 0 || len(owner.Rows) > 8192 || owner.Widget == "split" && len(owner.Rows) != 2 {
		return bad(owner.Path + ": tabs needs pages; split needs exactly two children; no regions")
	}
	children := make([]PaneChild, 0, len(owner.Rows))
	seen := map[string]bool{}
	for _, row := range owner.Rows {
		if len(row) != 1 || row[0] == nil {
			return bad(owner.Path + ": exactly one named child per row required")
		}
		child := row[0]
		id := strings.TrimPrefix(child.Path, owner.Path+"/")
		if id == child.Path || id == "" || id == "." || id == ".." || strings.ContainsAny(id, "/$") || seen[id] || child.Profile != owner.Profile {
			return bad(owner.Path + ": children need distinct direct named paths")
		}
		if owner.Widget == "tabs" && (child.Kind != "composition" || child.Widget != "page") {
			return bad(child.Path + ": tabs accepts only page children")
		}
		seen[id] = true
		children = append(children, PaneChild{ID: id, Node: child})
	}
	return children, nil
}

func validatePanePlacement(n, parent *Instance) {
	if n.Kind == "composition" {
		if n.Widget == "page" && parent != nil && (parent.Kind != "composition" || parent.Widget != "tabs") {
			fail("page-placement", n.Path+": page must be a direct child of tabs", n.Span)
		}
		if n.Widget == "tabs" || n.Widget == "split" {
			children, err := PaneChildren(n)
			if err != nil {
				panic(err)
			}
			if selected, exists := n.Arguments["selected"]; n.Widget == "tabs" && exists {
				id := selected.(Literal).Value.(string)
				valid := false
				for _, child := range children {
					if child.ID == id && child.Node.Layout["enabled"] != false && child.Node.Layout["visible"] != false {
						valid = true
					}
				}
				if !valid {
					fail("tab-selection", n.Path+": selected page must name an eligible direct page", selected.(Literal).Span)
				}
			}
		}
	}
	for _, row := range n.Rows {
		for _, child := range row {
			validatePanePlacement(child, n)
		}
	}
	for _, region := range n.Regions {
		validatePanePlacement(region.Node, n)
	}
}
