package previews

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	sp "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

type DomainState struct {
	Commits       map[string]int    `json:"commits"`
	Values        map[string]string `json:"values"`
	Marked        bool              `json:"marked"`
	DetailCommits int               `json:"detailCommits"`
	DetailNote    string            `json:"detailNote"`
}
type RequestIdentity struct {
	SourceSHA256 string            `json:"sourceSHA256"`
	Revision     uint64            `json:"revision"`
	Resources    map[string]string `json:"resources"`
	Renderers    map[string]string `json:"renderers"`
}
type Fixture struct {
	Engine                                                           *sdl.Engine
	Log                                                              func(string, any)
	PrepareResources                                                 func(ui.Snapshot) error
	mu                                                               sync.Mutex
	nonmodal                                                         bool
	revision                                                         uint64
	resources, rendererModes, rendererRevisions, policies, documents map[string]string
	faults                                                           map[string]string
	loans, returns                                                   map[string][]byte
	requests                                                         map[uint64]RequestIdentity
	calls, changes, rendererCalls                                    map[string]int
	nextActions                                                      map[string]string
	domain                                                           DomainState
	closed                                                           bool
}

func New(nonmodal bool) (*Fixture, error) {
	f := &Fixture{nonmodal: nonmodal, revision: 1,
		resources:     map[string]string{"Hero": "good", "ReportSVG": "good", "DetailSVG": "good"},
		rendererModes: map[string]string{"ReportA": "good", "ReportB": "good"}, rendererRevisions: map[string]string{"ReportA": "r1", "ReportB": "r1"},
		policies: map[string]string{}, documents: map[string]string{}, faults: map[string]string{}, loans: map[string][]byte{}, returns: map[string][]byte{}, requests: map[uint64]RequestIdentity{},
		calls: map[string]int{}, changes: map[string]int{}, rendererCalls: map[string]int{Targets["ReportA"]: 0, Targets["ReportB"]: 0}, nextActions: map[string]string{},
		domain: DomainState{Commits: map[string]int{}, Values: map[string]string{}}}
	for _, alias := range []string{"Note", "DetailNote"} {
		f.changes[Targets[alias]] = 0
	}
	program, err := sp.CompileActions(Actions)
	if err != nil {
		return nil, err
	}
	registry := sdl.Registry{}
	for name, a := range program.Actions {
		f.calls[name] = 0
		f.domain.Commits[name] = 0
		registry[a.GoSymbol] = sdl.Binding{Input: program.Records[a.Input], Output: program.Records[a.Output], Call: f.handler(name)}
	}
	f.Engine, err = sdl.New(program, registry)
	return f, err
}
func (f *Fixture) emit(kind string, v any) {
	if f.Log != nil {
		f.Log(kind, v)
	}
}
func (f *Fixture) Source() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return sourceFor(f.nonmodal, f.policies, f.documents)
}
func (f *Fixture) handler(name string) sdl.Handler {
	return func(ctx context.Context, r sdl.Record) (sdl.Record, error) {
		f.mu.Lock()
		f.calls[name]++
		count := f.calls[name]
		mode := f.nextActions[name]
		delete(f.nextActions, name)
		f.mu.Unlock()
		if mode == "" {
			mode = "success"
		}
		invocation, _ := sdl.CurrentInvocation(ctx)
		f.emit("action", map[string]any{"action": name, "input": r, "sequence": invocation.Sequence, "calls": count, "outcome": mode})
		if mode == "error" {
			return nil, fmt.Errorf("injected %s error", name)
		}
		if name == "AcceptDetail" && mode == "false" {
			return sdl.Record{"Accepted": sdl.Boolean(false), "Message": sdl.Text("Detail rejected")}, nil
		}
		value := r["Value"].Text
		f.mu.Lock()
		switch name {
		case "Load":
			value = "Loaded preview metadata"
		case "Mark":
			f.domain.Marked = r["Checked"].Boolean
			value = fmt.Sprintf("Marked: %t", f.domain.Marked)
		case "AcceptDetail":
			f.domain.DetailCommits++
			f.domain.DetailNote = r["Note"].Text
		}
		if name != "Load" {
			f.domain.Commits[name]++
			f.domain.Values[name] = value
		}
		f.mu.Unlock()
		if name == "AcceptDetail" {
			return sdl.Record{"Accepted": sdl.Boolean(true), "Message": sdl.Text("Detail saved")}, nil
		}
		return sdl.Record{"Value": sdl.Text(value)}, nil
	}
}
func Plans() map[string]bridge.Plan {
	token := sdl.Text("preview-metadata")
	return map[string]bridge.Plan{
		"previews.Load":         {Inputs: map[string]bridge.Source{"Value": {Literal: &token}}, OutputField: "Value"},
		"previews.Mark":         {Inputs: map[string]bridge.Source{"Checked": {EventField: bridge.CommandChecked}}, OutputField: "Value"},
		"previews.SaveNote":     {Inputs: map[string]bridge.Source{"Value": {EventField: bridge.ControlText}}, OutputField: "Value"},
		"previews.SaveDetail":   {Inputs: map[string]bridge.Source{"Value": {EventField: bridge.ControlText}}, OutputField: "Value"},
		"previews.AcceptDetail": {ResultMode: bridge.DialogAcceptResult, Inputs: map[string]bridge.Source{"Note": {EventField: bridge.DialogFieldValue, FieldPath: "detailNote"}}, AcceptField: "Accepted", MessageField: "Message"},
	}
}

// Request owns fresh loan buffers. Mutating those buffers is deliberately separate
// from changing application-authoritative identities, which invalidates Guard.
func (f *Fixture) Request(sequence uint64) (fynehost.DocumentRequest, error) {
	f.mu.Lock()
	source := sourceFor(f.nonmodal, f.policies, f.documents)
	version := f.revision
	resources := map[string]markdown.PreparedSVG{}
	renderers := map[string]markdown.MarkdownRenderer{}
	id := RequestIdentity{SourceSHA256: digest([]byte(source)), Revision: version, Resources: map[string]string{}, Renderers: map[string]string{}}
	f.loans = map[string][]byte{}
	for _, alias := range []string{"Hero", "ReportSVG", "DetailSVG"} {
		slot := f.resources[alias]
		if slot == "missing" {
			continue
		}
		r, err := supplied(slot)
		if err != nil {
			f.mu.Unlock()
			return fynehost.DocumentRequest{}, err
		}
		path := Targets[alias]
		resources[path] = r
		f.loans[path] = r.Resource.SVG
		id.Resources[path] = r.ProviderID + ":" + r.SHA256
	}
	for _, alias := range []string{"ReportA", "ReportB"} {
		mode := f.rendererModes[alias]
		doc, overridden := f.documents[alias]
		if mode == "missing" || (overridden && !strings.Contains(doc, "```mermaid")) {
			continue
		}
		path := Targets[alias]
		r := markdown.MarkdownRenderer{ProviderID: "fixture-diagram-" + alias + "/1", Revision: f.rendererRevisions[alias], Renderer: countedRenderer{f, path, mode}}
		renderers[path] = r
		id.Renderers[path] = r.ProviderID + ":" + r.Revision
	}
	for alias, fault := range f.faults {
		injectFault(resources, renderers, alias, fault)
	}
	f.faults = map[string]string{}
	f.requests[sequence] = id
	f.mu.Unlock()
	doc, err := parser.Parse(source)
	if err != nil {
		return fynehost.DocumentRequest{}, err
	}
	engineRevision := f.Engine.Revision()
	return fynehost.DocumentRequest{Document: doc, Entry: "page", SessionID: "previews", SourceRevision: fmt.Sprintf("previews-%d-%s", version, id.SourceSHA256), Sequence: sequence, Mode: preparation.Connected, SVGResources: resources, MarkdownRenderers: renderers,
		Guard: func() error {
			f.mu.Lock()
			stale := f.revision != version
			f.mu.Unlock()
			if stale || f.Engine.Revision() != engineRevision {
				return fmt.Errorf("stale preview source/resource/renderer identity")
			}
			return nil
		},
		PrepareResources: func(s ui.Snapshot) error {
			if f.PrepareResources != nil {
				return f.PrepareResources(s)
			}
			return nil
		},
		Bind: func(s *ui.Session, d *parser.Document) error {
			if _, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"previews": f.Engine}, Plans(), nil); err != nil {
				return err
			}
			for _, alias := range []string{"Note", "DetailNote"} {
				w, ok := s.Widget(Targets[alias])
				if !ok {
					return fmt.Errorf("missing %s", alias)
				}
				if err := s.ObserveChanges(w.Handle, f.changed); err != nil {
					return err
				}
			}
			return nil
		}}, nil
}
func injectFault(svg map[string]markdown.PreparedSVG, md map[string]markdown.MarkdownRenderer, alias, fault string) {
	path := Targets[alias]
	if r, ok := svg[path]; ok {
		switch fault {
		case "digest":
			r.SHA256 = strings.Repeat("0", 64)
		case "source":
			r.Source.Object = "Wrong"
		case "provider":
			r.ProviderID = ""
		case "extra":
			svg["page/missing"] = r
		case "wrong-kind":
			delete(svg, path)
			svg[Targets["Note"]] = r
			return
		}
		svg[path] = r
	} else {
		r := md[path]
		switch fault {
		case "provider":
			r.ProviderID = ""
		case "revision":
			r.Revision = ""
		case "nil":
			r.Renderer = nil
		case "extra":
			md["page/missing"] = r
		case "wrong-kind":
			delete(md, path)
			md[Targets["Note"]] = r
			return
		case "unused":
			delete(md, path)
			md[Targets["Prose"]] = r
			return
		}
		md[path] = r
	}
}
func (f *Fixture) changed(change ui.FieldChange) {
	path := change.Field.InstancePath
	f.mu.Lock()
	f.changes[path]++
	count := f.changes[path]
	f.mu.Unlock()
	f.emit("change", map[string]any{"path": path, "count": count, "field": change.Field})
}
func copyCounts(in map[string]int) map[string]int {
	out := map[string]int{}
	for k, v := range in {
		out[k] = v
	}
	return out
}
func (f *Fixture) Counts() (map[string]int, map[string]int, map[string]int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return copyCounts(f.calls), copyCounts(f.changes), copyCounts(f.rendererCalls)
}
func (f *Fixture) Domain() DomainState {
	f.mu.Lock()
	defer f.mu.Unlock()
	d := f.domain
	d.Commits = copyCounts(d.Commits)
	d.Values = map[string]string{}
	for k, v := range f.domain.Values {
		d.Values[k] = v
	}
	return d
}
func (f *Fixture) NextAction(name, mode string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.calls[name]; !ok {
		return fmt.Errorf("unknown action %q", name)
	}
	if mode != "success" && mode != "error" && !(name == "AcceptDetail" && (mode == "true" || mode == "false")) {
		return fmt.Errorf("invalid action outcome")
	}
	f.nextActions[name] = mode
	return nil
}
func (f *Fixture) Close() {
	f.mu.Lock()
	if f.closed {
		f.mu.Unlock()
		return
	}
	f.closed = true
	f.mu.Unlock()
	f.Engine.Close()
}
