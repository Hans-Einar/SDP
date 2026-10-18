// Package collections supplies bounded in-memory data and a real action-core
// fixture. It never opens collection IDs as paths or accesses an XFMD workspace.
package collections

import (
	"context"
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	sdlparser "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
	"sort"
	"sync"
)

const Actions = `language action-core version 0.1.
action Activate.
record ActivateInput.
record ActivateOutput.
Activate invokes GoActivate.
Activate returns ActivateOutput.
Activate takes ActivateInput.
ActivateInput field ItemId as text.
ActivateOutput field Preview as text.
`
const Source = `sdui 0.3;
ref: nav "collections.sdl";
page = [
 header="Collection navigation fixture";
 body=<
  nav=tree("Navigation",callback=nav.Activate.@invoke) {x=1fr,y=fill,overflow-x=scroll,overflow-y=scroll},
  entries=list("Entries",callback=nav.Activate.@invoke) {x=1fr,y=fill,overflow-x=scroll,overflow-y=scroll}
 > {x=fill,y=fill};
 footer=<preview=input("Preview",value="No activation") {x=fill}> {x=fill}
] {x=fill,y=fill};
nav.Activate.setHandle(page.footer.preview);
`

type reply struct {
	data ui.CollectionData
	err  error
}
type pending struct {
	request ui.LoadRequest
	reply   chan reply
}
type Fixture struct {
	Engine  *sdl.Engine
	mu      sync.Mutex
	pending map[string]pending
	known   map[string]string
	calls   int
	closed  bool
	Log     func(string, any)
}

func New() (*Fixture, error) {
	f := &Fixture{pending: map[string]pending{}, known: map[string]string{"folder-b": "Folder B", "file": "Duplicate label"}}
	for i := 0; i < 100; i++ {
		if i%20 != 0 && i%20 != 19 {
			f.known[fmt.Sprintf("entry-%03d", i)] = fmt.Sprintf("Entry %03d", i)
		}
	}
	p, err := sdlparser.CompileActions(Actions)
	if err != nil {
		return nil, err
	}
	f.Engine, err = sdl.New(p, sdl.Registry{"GoActivate": {Input: sdlparser.RecordType{"ItemId": sdlparser.TextType}, Output: sdlparser.RecordType{"Preview": sdlparser.TextType}, Call: func(_ context.Context, r sdl.Record) (sdl.Record, error) {
		id := r["ItemId"].Text
		f.mu.Lock()
		f.calls++
		count := f.calls
		preview, known := f.known[id]
		f.mu.Unlock()
		if f.Log != nil {
			f.Log("action", map[string]any{"item": id, "calls": count})
		}
		if !known {
			return nil, fmt.Errorf("unknown item")
		}
		return sdl.Record{"Preview": sdl.Text("Preview: " + preview + " [" + id + "]")}, nil
	}}})
	return f, err
}
func (f *Fixture) Calls() int { f.mu.Lock(); defer f.mu.Unlock(); return f.calls }
func (f *Fixture) load(ctx context.Context, r ui.LoadRequest) (ui.CollectionData, error) {
	key := fmt.Sprintf("%s:%d:%d", r.Target.Handle.Path, r.Target.ModelRevision, r.RequestID)
	ch := make(chan reply, 1)
	f.mu.Lock()
	if f.closed {
		f.mu.Unlock()
		return ui.CollectionData{}, fmt.Errorf("fixture closed")
	}
	f.pending[key] = pending{r, ch}
	f.mu.Unlock()
	if f.Log != nil {
		f.Log("load", map[string]any{"key": key, "request": r})
	}
	// Deliberately controllable even after cancellation: tests can deliver late
	// success/error and establish that the host revoked publication rights.
	result := <-ch
	if f.Log != nil {
		f.Log("load-return", map[string]any{"key": key, "canceled": ctx.Err() != nil})
	}
	return result.data, result.err
}

// Pending returns an owned, sorted list of started provider barriers. Complete
// must name one of these keys; no timing sleep chooses which request receives data.
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
		return fmt.Errorf("unknown pending request %s", key)
	}
	result := reply{}
	switch outcome {
	case "error":
		result.err = fmt.Errorf("controlled provider failure")
	case "empty":
	case "invalid":
		result.data.Items = []ui.CollectionItem{{ID: "bad", Kind: ui.Row, Label: "bad"}}
	default:
		for i := 0; i < 24; i++ {
			result.data.Items = append(result.data.Items, ui.CollectionItem{ID: ui.ItemID(fmt.Sprintf("%s/file-%02d", p.request.Target.ItemID, i)), Parent: p.request.Target.ItemID, Kind: ui.Row, Label: fmt.Sprintf("Document %02d - a long label for horizontal scrolling and identity inspection", i), ChildrenLoaded: true})
		}
	}
	f.mu.Lock()
	for _, item := range result.data.Items {
		f.known[string(item.ID)] = item.Label
	}
	f.mu.Unlock()
	p.reply <- result
	return nil
}
func (f *Fixture) Providers(empty bool) map[string]ui.CollectionProvider {
	tree := ui.CollectionData{Items: []ui.CollectionItem{{ID: "folder-a", Kind: ui.Group, Label: "Folder A (lazy group)", HasChildren: true}, {ID: "folder-b", Kind: ui.Row, Label: "Folder B (lazy selectable)", HasChildren: true}, {ID: "file", Kind: ui.Row, Label: "Duplicate label", ChildrenLoaded: true}}}
	list := ui.CollectionData{}
	for i := 0; i < 100; i++ {
		item := ui.CollectionItem{ID: ui.ItemID(fmt.Sprintf("entry-%03d", i)), Kind: ui.Row, Label: fmt.Sprintf("Entry %03d - long text for horizontal scrolling", i), ChildrenLoaded: true}
		if i%20 == 0 {
			item.Kind = ui.Group
			item.Label = "Group"
		}
		if i%20 == 19 {
			item.Kind = ui.Separator
			item.Label = ""
		}
		list.Items = append(list.Items, item)
	}
	if empty {
		tree = ui.CollectionData{}
	}
	return map[string]ui.CollectionProvider{"page/body/nav": {ID: "fixture-tree", Epoch: 1, Initial: tree, RootLoaded: !empty, Load: f.load}, "page/body/entries": {ID: "fixture-list", Epoch: 1, Initial: list, RootLoaded: true}}
}
func (f *Fixture) Request(source string, sequence uint64, empty bool) (fynehost.DocumentRequest, error) {
	d, err := parser.Parse(source)
	if err != nil {
		return fynehost.DocumentRequest{}, err
	}
	revision := f.Engine.Revision()
	return fynehost.DocumentRequest{Document: d, Entry: "page", SessionID: "collections", SourceRevision: fmt.Sprintf("fixture-%d", sequence), Sequence: sequence, Mode: preparation.Connected, Providers: f.Providers(empty), Guard: func() error {
		if f.Engine.Revision() != revision {
			return fmt.Errorf("stale SDL revision")
		}
		return nil
	}, Bind: func(s *ui.Session, doc *parser.Document) error {
		_, err := bridge.Bind(context.Background(), s, doc, map[string]*sdl.Engine{"nav": f.Engine}, map[string]bridge.Plan{"nav.Activate": {Inputs: map[string]bridge.Source{"ItemId": {EventField: bridge.CollectionItemID}}, OutputField: "Preview"}}, nil)
		return err
	}}, nil
}
func (f *Fixture) Close() {
	f.mu.Lock()
	if f.closed {
		f.mu.Unlock()
		return
	}
	f.closed = true
	for k, p := range f.pending {
		p.reply <- reply{err: fmt.Errorf("fixture closed")}
		delete(f.pending, k)
	}
	f.mu.Unlock()
	f.Engine.Close()
}
