// Package snapshot exports selected UI state through the production geometry path.
package snapshot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/presentation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SDUI/go/svg"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"runtime"
	"sort"
)

const Version = "sdl-sdui-document/1"

type WidgetState struct {
	Value   *string `json:"value,omitempty"`
	Label   *string `json:"label,omitempty"`
	Enabled *bool   `json:"enabled,omitempty"`
	Visible *bool   `json:"visible,omitempty"`
}
type Options struct {
	Entry         string
	Width, Height float64
	State         map[string]WidgetState
	Design        string
	Renderer      markdown.Renderer
	RendererID    string
}

func Build(ctx context.Context, source string, o Options) (*documents.Bundle, error) {
	_, roots, e := parser.Compile(source)
	if e != nil {
		return nil, e
	}
	root := roots[o.Entry]
	if root == nil {
		return nil, fmt.Errorf("unknown UI entry %q", o.Entry)
	}
	session, e := ui.New("document", root)
	if e != nil {
		return nil, e
	}
	defer session.Close()
	paths := []string{}
	for p := range o.State {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	updates := []ui.Update{}
	for _, p := range paths {
		w, ok := session.Widget(p)
		if !ok {
			return nil, fmt.Errorf("unknown state widget %s", p)
		}
		v := o.State[p]
		add := func(prop ui.Property, value ui.Value) {
			updates = append(updates, ui.Update{Handle: w.Handle, Property: prop, Value: value, ExpectedValueRevision: w.ValueRevision})
		}
		if v.Value != nil {
			add(ui.AcceptedValue, ui.Text(*v.Value))
		}
		if v.Label != nil {
			add(ui.Label, ui.Text(*v.Label))
		}
		if v.Enabled != nil {
			add(ui.Enabled, ui.Bool(*v.Enabled))
		}
		if v.Visible != nil {
			add(ui.Visible, ui.Bool(*v.Visible))
		}
	}
	if e = session.Apply(session.Revision, 1, updates); e != nil {
		return nil, e
	}
	root = session.SnapshotRoot()
	provider, e := markdown.Prepare(root, o.Renderer)
	if e != nil {
		return nil, e
	}
	geometry, e := (&layout.Engine{Measure: provider}).Layout(root, layout.Size{W: o.Width, H: o.Height})
	if e != nil {
		return nil, e
	}
	image, e := svg.Render(geometry, svg.Options{Width: o.Width, Height: o.Height, Content: provider})
	if e != nil {
		return nil, e
	}
	md, e := presentation.Markdown(root, 120)
	if e != nil {
		return nil, e
	}
	state, _ := json.MarshalIndent(o.State, "", "  ")
	b := &documents.Bundle{Files: map[string][]byte{}, Manifest: documents.Manifest{Version: Version, Revision: documents.Hash([]byte(source)), Renderer: o.RendererID}}
	b.Put("ui.svg", image)
	b.Put("structure.md", md)
	b.Put("inputs/ui.sdui", source)
	b.Files["inputs/state.json"] = append(state, '\n')
	entry := "# SDUI — valgt statisk tilstand\n\n![UI fra felles layout](ui.svg)\n\n[Struktur og rå Markdown](structure.md) · [Proveniens](provenance.json)\n\nStatisk eksport: ingen callbacks er kjørt. UI-verdier er eksplisitte dokumentasjonsdata.\n"
	if o.Design != "" {
		v, e := viewpoint.New(o.Design)
		if e != nil {
			return nil, e
		}
		nav, e := documents.Build(ctx, v, documents.Options{Navigator: true, Project: "design"})
		if e != nil {
			return nil, e
		}
		for n, f := range nav.Files {
			b.Files["design/"+n] = f
		}
		b.Put("inputs/architecture.design", o.Design)
		entry += "\n[SDL-arkitektur og viewpoints](design/index.md) · [Navigator](design/navigator.md)\n"
	}
	provenance := map[string]any{"tool": Version, "toolchain": runtime.Version(), "sdui_profile": "0.2", "layout_profile": "relative-layout/1", "font": "Go Regular logical DIP", "ui_sha256": documents.Hash([]byte(source)), "state_sha256": documents.Hash(state), "design_sha256": documents.Hash([]byte(o.Design)), "projector": viewpoint.Version, "entry": o.Entry, "width": o.Width, "height": o.Height, "renderer": o.RendererID, "widgets": session.Widgets()}
	p, e := json.MarshalIndent(provenance, "", "  ")
	if e != nil {
		return nil, e
	}
	b.Files["provenance.json"] = append(p, '\n')
	b.Put("entry.md", entry)
	b.Seal()
	return b, b.CheckLinks()
}
