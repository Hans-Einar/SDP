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
			return &Diagnostic{c, n.Path, n.Span}
		}
		return nil
	}
	if profile != "sdui/0.2" {
		return &Diagnostic{Capability{Frontend, profile, 1}, root.Path, root.Span}
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
		case "markdown":
			if err := require(n, Provider, "markdown"); err != nil {
				return err
			}
		case "widget":
			id := n.Widget
			switch id {
			case "button", "input":
			case "svg":
				id = "svg-placeholder"
			default:
				return &Diagnostic{Capability{Widget, id, 1}, n.Path, n.Span}
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
			return &Diagnostic{Capability{Layout, "node:" + n.Kind, 1}, n.Path, n.Span}
		}
		for _, axis := range []string{"x", "y"} {
			if n.Layout["overflow-"+axis] == "scroll" {
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
