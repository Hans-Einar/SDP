// Package panes provides a bounded, real SDL action-core fixture for WCI2-M1.
// Provider barriers and failure injection are fixture controls, not product APIs.
package panes

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	sdlparser "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

const Actions = `language action-core version 0.1.
action Page.
record PageInput.
record PageOutput.
Page invokes GoPage.
Page returns PageOutput.
Page takes PageInput.
PageInput field PageId as text.
PageInput field PreviousId as text.
PageOutput field Preview as text.
`
const Source = `sdui 0.3;
ref: panes "panes.sdl";
page = [
 header="Panes interaction fixture";
 body=split(axis="horizontal",proportion=0.65,minFirst=0.25,minSecond=0.2)[
  workspace=tabs("Workspace",selected="overview",callback=panes.Page.@invoke)[
   overview=page("Overview")[
    edit=input("Overview draft",value="Overview accepted") {x=fill};
    nav=tree("Documents") {x=fill,y=fill,overflow-x=scroll,overflow-y=scroll}
   ];
   notes=page("Notes")[notesEdit=input("Notes draft",value="Notes accepted") {x=fill}]
  ];
  aside=[note=input("Side draft",value="Side accepted") {x=fill}]
 ] {x=fill,y=fill};
 footer=<preview=input("Page result",value="No page activation") {x=fill}> {x=fill}
] {x=fill,y=fill};
panes.Page.setHandle(page.footer.preview);
`
const (
	WindowTitle = "SDUI WCI2 Panes"
	TabsPath    = "page/body/workspace"
	SplitPath   = "page/body"
	TreePath    = "page/body/workspace/overview/nav"
	PreviewPath = "page/footer/preview"
)

type loadReply struct {
	data ui.CollectionData
	err  error
}
type pendingLoad struct {
	request ui.LoadRequest
	reply   chan loadReply
}
type Fixture struct {
	Engine          *sdl.Engine
	mu              sync.Mutex
	calls           int
	outcome         string
	closed          bool
	resourceFailure bool
	pending         map[string]pendingLoad
	Log             func(string, any) // Install before running providers or actions.
	Conflict        func() error      // UI-owner-only reentrant receiver-draft mutation for conflict injection.
}

func New() (*Fixture, error) {
	f := &Fixture{pending: map[string]pendingLoad{}}
	p, err := sdlparser.CompileActions(Actions)
	if err != nil {
		return nil, err
	}
	f.Engine, err = sdl.New(p, sdl.Registry{"GoPage": {
		Input:  sdlparser.RecordType{"PageId": sdlparser.TextType, "PreviousId": sdlparser.TextType},
		Output: sdlparser.RecordType{"Preview": sdlparser.TextType},
		Call: func(ctx context.Context, r sdl.Record) (sdl.Record, error) {
			f.mu.Lock()
			f.calls++
			calls, mode := f.calls, f.outcome
			f.outcome = ""
			f.mu.Unlock()
			invocation, _ := sdl.CurrentInvocation(ctx)
			if f.Log != nil {
				f.Log("action", map[string]any{"page": r["PageId"].Text, "previous": r["PreviousId"].Text, "calls": calls, "sequence": invocation.Sequence, "outcome": mode})
			}
			switch mode {
			case "error":
				return nil, fmt.Errorf("injected Page failure")
			case "invalid":
				return sdl.Record{"Preview": sdl.Integer(1)}, nil
			case "draft-conflict":
				if f.Conflict == nil {
					return nil, fmt.Errorf("conflict hook missing")
				}
				if err := f.Conflict(); err != nil {
					return nil, err
				}
			}
			return sdl.Record{"Preview": sdl.Text(r["PreviousId"].Text + " -> " + r["PageId"].Text)}, nil
		},
	}})
	return f, err
}
func (f *Fixture) Calls() int { f.mu.Lock(); defer f.mu.Unlock(); return f.calls }
func (f *Fixture) NextAction(outcome string) error {
	switch outcome {
	case "success", "error", "invalid", "draft-conflict":
	default:
		return fmt.Errorf("unknown action outcome %q", outcome)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return fmt.Errorf("fixture closed")
	}
	f.outcome = outcome
	return nil
}
func (f *Fixture) load(ctx context.Context, r ui.LoadRequest) (ui.CollectionData, error) {
	key := fmt.Sprintf("%s:%d:%d", r.Target.Handle.Path, r.Target.ModelRevision, r.RequestID)
	ch := make(chan loadReply, 1)
	f.mu.Lock()
	if f.closed {
		f.mu.Unlock()
		return ui.CollectionData{}, fmt.Errorf("fixture closed")
	}
	f.pending[key] = pendingLoad{r, ch}
	f.mu.Unlock()
	if f.Log != nil {
		f.Log("load", map[string]any{"key": key, "request": r})
	}
	// Cancellation deliberately does not release this barrier: late completion
	// must be rejected by runtime identity, not hidden by provider timing.
	reply := <-ch
	if f.Log != nil {
		f.Log("load-return", map[string]any{"key": key, "canceled": ctx.Err() != nil})
	}
	return reply.data, reply.err
}
func (f *Fixture) Pending() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	keys := make([]string, 0, len(f.pending))
	for key := range f.pending {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
func (f *Fixture) Complete(key, outcome string) error {
	switch outcome {
	case "success", "error", "empty", "invalid":
	default:
		return fmt.Errorf("unknown completion outcome %q", outcome)
	}
	f.mu.Lock()
	p, ok := f.pending[key]
	if ok {
		delete(f.pending, key)
	}
	f.mu.Unlock()
	if !ok {
		return fmt.Errorf("unknown pending key %q", key)
	}
	r := loadReply{}
	switch outcome {
	case "error":
		r.err = fmt.Errorf("injected pane provider failure")
	case "invalid":
		r.data.Items = []ui.CollectionItem{{Kind: ui.Row, Label: "Missing ID", ChildrenLoaded: true}}
	case "success":
		for i := 0; i < 24; i++ {
			r.data.Items = append(r.data.Items, ui.CollectionItem{ID: ui.ItemID(fmt.Sprintf("%s/child-%02d", p.request.Target.ItemID, i)), Parent: p.request.Target.ItemID, Kind: ui.Row, Label: fmt.Sprintf("Loaded document %02d - long horizontal scrolling label", i), ChildrenLoaded: true})
		}
	}
	p.reply <- r
	return nil
}
func (f *Fixture) Providers() map[string]ui.CollectionProvider {
	data := ui.CollectionData{Items: []ui.CollectionItem{{ID: "lazy", Kind: ui.Group, Label: "Lazy documents", HasChildren: true}}}
	for i := 0; i < 64; i++ {
		data.Items = append(data.Items, ui.CollectionItem{ID: ui.ItemID(fmt.Sprintf("doc-%02d", i)), Kind: ui.Row, Label: fmt.Sprintf("Document %02d - long stable label for retained horizontal and vertical scroll", i), ChildrenLoaded: true})
	}
	return map[string]ui.CollectionProvider{TreePath: {ID: "panes-tree", Epoch: 1, Initial: data, RootLoaded: true, Load: f.load}}
}
func Plans() map[string]bridge.Plan {
	return map[string]bridge.Plan{"panes.Page": {Inputs: map[string]bridge.Source{"PageId": {EventField: bridge.TabPageID}, "PreviousId": {EventField: bridge.TabPreviousPageID}}, OutputField: "Preview"}}
}
func (f *Fixture) Request(source string, sequence uint64) (fynehost.DocumentRequest, error) {
	d, err := parser.Parse(source)
	if err != nil {
		return fynehost.DocumentRequest{}, err
	}
	revision := f.Engine.Revision()
	return fynehost.DocumentRequest{Document: d, Entry: "page", SessionID: "panes", SourceRevision: fmt.Sprintf("panes-%d", sequence), Sequence: sequence, Mode: preparation.Connected, Providers: f.Providers(), PrepareResources: f.prepareResources, Guard: func() error {
		if f.Engine.Revision() != revision {
			return fmt.Errorf("stale SDL revision")
		}
		return nil
	}, Bind: func(s *ui.Session, doc *parser.Document) error {
		_, err := bridge.Bind(context.Background(), s, doc, map[string]*sdl.Engine{"panes": f.Engine}, Plans(), nil)
		return err
	}}, nil
}

// FailResource arms one final resource-preparation failure, without calling a user action.
func (f *Fixture) FailResource() { f.mu.Lock(); f.resourceFailure = true; f.mu.Unlock() }
func (f *Fixture) prepareResources(ui.Snapshot) error {
	f.mu.Lock()
	fail := f.resourceFailure
	f.resourceFailure = false
	f.mu.Unlock()
	if fail {
		return fmt.Errorf("injected final resource preparation failure")
	}
	return nil
}
func (f *Fixture) Close() {
	f.mu.Lock()
	if f.closed {
		f.mu.Unlock()
		return
	}
	f.closed = true
	for k, p := range f.pending {
		p.reply <- loadReply{err: fmt.Errorf("fixture closed")}
		delete(f.pending, k)
	}
	f.mu.Unlock()
	f.Engine.Close()
}
