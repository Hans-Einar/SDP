// Package preparation validates detached SDUI candidates without I/O, SDL or GUI dependencies.
package preparation

import (
	"fmt"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

type Dimension string

const (
	Frontend Dimension = "frontend"
	Layout   Dimension = "layout"
	Widget   Dimension = "widget"
	Viewport Dimension = "viewport"
	Provider Dimension = "provider"
	Host     Dimension = "host"
)

// Capability versions match exactly. The source profile is part of the frontend
// identifier (sdui/0.2); its admission contract has major version 1.
type Capability struct {
	Dimension Dimension
	ID        string
	Major     int
}
type Capabilities []Capability
type Diagnostic struct {
	Capability Capability
	Path       string
	Span       parser.Span
	Uses       []parser.UseSite
}

func (d *Diagnostic) Error() string {
	return fmt.Sprintf("capability: %s %s/%d at %s (%d:%d)", d.Capability.Dimension, d.Capability.ID, d.Capability.Major, d.Path, d.Span.Line, d.Span.Column)
}

// Check inspects hidden descendants too. Unknown normalized kinds cannot be
// enabled just by advertising a capability for an unknown rendering contract.
func Check(profile string, root *parser.Instance, supported Capabilities) error {
	if root == nil {
		return fmt.Errorf("preparation: missing root")
	}
	has := map[Capability]bool{}
	for _, c := range supported {
		has[c] = true
	}
	require := func(n *parser.Instance, dimension Dimension, id string) error {
		c := Capability{dimension, id, 1}
		if !has[c] {
			return &Diagnostic{Capability: c, Path: n.Path, Span: n.Span, Uses: append([]parser.UseSite(nil), n.Uses...)}
		}
		return nil
	}
	effective, err := parser.EffectiveProfile(root)
	if err != nil {
		return err
	}
	if profile != effective {
		return fmt.Errorf("profile-mismatch: document %s, root %s", profile, effective)
	}
	if err := require(root, Frontend, profile); err != nil {
		return err
	}
	if err := require(root, Layout, "relative"); err != nil {
		return err
	}
	seen := map[*parser.Instance]bool{}
	var walk func(*parser.Instance) error
	walk = func(n *parser.Instance) error {
		if n == nil || seen[n] {
			return fmt.Errorf("preparation: nil or repeated normalized node")
		}
		seen[n] = true
		switch n.Kind {
		case "frame", "group":
		case "composition":
			if profile != "sdui/0.3" {
				return &Diagnostic{Capability: Capability{Widget, n.Widget, 1}, Path: n.Path, Span: n.Span, Uses: n.Uses}
			}
			switch n.Widget {
			case "tabs", "split":
				if err := require(n, Widget, n.Widget); err != nil {
					return err
				}
				if err := require(n, Host, n.Widget); err != nil {
					return err
				}
				if n.Widget == "tabs" {
					if err := require(n, Host, "tab-activate"); err != nil {
						return err
					}
				}
			case "page":
				if n.Argument("icon") != "" {
					if err := require(n, Provider, "icon"); err != nil {
						return err
					}
				}
			default:
				return &Diagnostic{Capability: Capability{Widget, n.Widget, 1}, Path: n.Path, Span: n.Span, Uses: n.Uses}
			}
		case "markdown":
			if err := require(n, Provider, "markdown"); err != nil {
				return err
			}
		case "widget":
			id := n.Widget
			switch id {
			case "button", "input":
			case "tree", "list":
				if profile != "sdui/0.3" {
					return &Diagnostic{Capability: Capability{Widget, id, 1}, Path: n.Path, Span: n.Span, Uses: n.Uses}
				}
				if err := require(n, Layout, "collections"); err != nil {
					return err
				}
				if err := require(n, Provider, "collection-data"); err != nil {
					return err
				}
			case "svg":
				id = "svg-placeholder"
			default:
				return &Diagnostic{Capability: Capability{Widget, id, 1}, Path: n.Path, Span: n.Span, Uses: n.Uses}
			}
			if err := require(n, Widget, id); err != nil {
				return err
			}
			dimension := Host
			if n.Widget == "svg" {
				dimension = Provider
			}
			if err := require(n, dimension, id); err != nil {
				return err
			}
		default:
			return &Diagnostic{Capability: Capability{Layout, "node:" + n.Kind, 1}, Path: n.Path, Span: n.Span, Uses: n.Uses}
		}
		for _, axis := range []string{"x", "y"} {
			if n.Layout["overflow-"+axis] == "scroll" {
				if profile != "sdui/0.3" || !(n.Kind == "frame" || n.Kind == "group" || n.Kind == "widget" && (n.Widget == "tree" || n.Widget == "list")) {
					return &Diagnostic{Capability: Capability{Viewport, "scroll-" + axis, 1}, Path: n.Path, Span: n.Span, Uses: n.Uses}
				}
				if err := require(n, Host, "viewport"); err != nil {
					return err
				}
				if err := require(n, Viewport, "scroll-"+axis); err != nil {
					return err
				}
			}
		}
		for _, r := range n.Regions {
			if err := walk(r.Node); err != nil {
				return err
			}
		}
		for _, row := range n.Rows {
			for _, child := range row {
				if err := walk(child); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return walk(root)
}
