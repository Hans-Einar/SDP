// Package prototype defines the standalone local-prototype readiness contract.
package prototype

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost/admission"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	"github.com/Hans-Einar/SDP/SDUI/go/reload"
)

type Report struct {
	Schema     string `json:"schema"`
	Operation  string `json:"operation"`
	Status     string `json:"status"`
	Revision   string `json:"revision"`
	Entry      string `json:"entry"`
	Source     string `json:"source"`
	Diagnostic string `json:"diagnostic"`
	Profile    string `json:"profile,omitempty"`
}

func Check(path, entry, revision string) (reload.Candidate, Report, error) {
	c := reload.Read(path, entry, 1, nil)
	r := Report{Schema: "sdptool/0.2", Operation: "sdui-check", Entry: entry, Source: path, Revision: c.Hash}
	if c.Err != nil {
		return c, r, c.Err
	}
	if revision != "" && c.Hash != revision {
		return c, r, fmt.Errorf("stale: SDUI source changed; regenerate the preview")
	}
	// The standalone adapter has neither native panes nor application-owned
	// collection providers. Check hidden descendants before generic admission so
	// static 0.3 support cannot be confused with interactive readiness.
	if c.Document.Profile == "sdui/0.3" {
		r.Profile = c.Document.Profile
		if _, err := parser.ResolveInteractions(c.Root); err != nil {
			return c, r, err
		}
		var unavailable error
		c.Root.Walk(func(n *parser.Instance) {
			if unavailable == nil && n.Kind == "widget" && (n.Widget == "checkbox" || n.Widget == "slider" || n.Widget == "number" || n.Widget == "select") {
				capability := preparation.Capability{Dimension: preparation.Host, ID: n.Widget, Major: 1}
				message := "unsupported-value: standalone prototype has no native scalar adapter"
				if n.Widget == "select" {
					capability = preparation.Capability{Dimension: preparation.Provider, ID: "choice-options", Major: 1}
					message = "unsupported-provider: standalone prototype requires application-supplied choice options"
				}
				unavailable = fmt.Errorf("%s: %w", message, &preparation.Diagnostic{Capability: capability, Path: n.Path, Span: n.Span, Uses: append([]parser.UseSite(nil), n.Uses...)})
			}
			m2 := n.Kind == "composition" && (n.Widget == "menu" || n.Widget == "menuGroup" || n.Widget == "dialog") || n.Kind == "widget" && (n.Widget == "command" || n.Widget == "item" || n.Widget == "separator") || parser.IsCommandButton(n)
			if unavailable == nil && m2 {
				unavailable = fmt.Errorf("unsupported-interaction: standalone prototype has no command/menu/dialog adapter: %w", &preparation.Diagnostic{
					Capability: preparation.Capability{Dimension: preparation.Host, ID: n.Widget, Major: 1},
					Path:       n.Path, Span: n.Span, Uses: append([]parser.UseSite(nil), n.Uses...),
				})
			}
			if unavailable == nil && n.Kind == "composition" {
				unavailable = fmt.Errorf("unsupported-pane: standalone prototype has no native pane adapter: %w", &preparation.Diagnostic{
					Capability: preparation.Capability{Dimension: preparation.Host, ID: n.Widget, Major: 1},
					Path:       n.Path, Span: n.Span, Uses: append([]parser.UseSite(nil), n.Uses...),
				})
			}
			if unavailable == nil && n.Kind == "widget" && (n.Widget == "tree" || n.Widget == "list") {
				unavailable = fmt.Errorf("unsupported-provider: standalone prototype requires an application-supplied collection provider: %w", &preparation.Diagnostic{
					Capability: preparation.Capability{Dimension: preparation.Provider, ID: "collection-data", Major: 1},
					Path:       n.Path, Span: n.Span, Uses: append([]parser.UseSite(nil), n.Uses...),
				})
			}
		})
		if unavailable != nil {
			return c, r, unavailable
		}
	}
	prepared, err := preparation.Prepare(preparation.Request{
		Document: c.Document, Entry: entry, SessionID: "preflight", SourceRevision: c.Hash,
		Mode: preparation.Prototype, Capabilities: admission.Capabilities(),
		ValidateLayout: func(root *parser.Instance) error { return admission.Layout(root, layout.Size{W: 1280, H: 800}) },
	})
	if err != nil {
		return c, r, err
	}
	defer prepared.Close()
	unbound := prepared.Unbound
	r.Status = "prototype"
	r.Diagnostic = "Local Fyne prototype; no SDL runtime"
	if unbound > 0 {
		r.Status = "prototype-unbound"
		r.Diagnostic = fmt.Sprintf("Local Fyne prototype; %d unbound SDL callbacks (no SDL execution)", unbound)
	}
	return c, r, nil
}
