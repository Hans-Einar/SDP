package commands

import (
	"context"
	"fmt"
	"sort"

	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

type loadReply struct {
	data ui.CollectionData
	err  error
}
type pendingLoad struct {
	request ui.LoadRequest
	reply   chan loadReply
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
	// A revoked request remains releasable to establish rejection of late replies.
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
func (f *Fixture) Complete(key, mode string) error {
	switch mode {
	case "success", "error", "empty", "invalid":
	default:
		return fmt.Errorf("unknown completion outcome %q", mode)
	}
	f.mu.Lock()
	p, ok := f.pending[key]
	if ok {
		delete(f.pending, key)
	}
	f.mu.Unlock()
	if !ok {
		return fmt.Errorf("unknown pending provider key %q", key)
	}
	result := loadReply{}
	switch mode {
	case "error":
		result.err = fmt.Errorf("injected commands provider failure")
	case "invalid":
		result.data.Items = []ui.CollectionItem{{Kind: ui.Row, Label: "No identity", ChildrenLoaded: true}}
	case "success":
		for i := 0; i < 12; i++ {
			result.data.Items = append(result.data.Items, ui.CollectionItem{ID: ui.ItemID(fmt.Sprintf("%s/item-%02d", p.request.Target.ItemID, i)), Parent: p.request.Target.ItemID, Kind: ui.Row, Label: fmt.Sprintf("Loaded item %02d", i), ChildrenLoaded: true})
		}
	}
	p.reply <- result
	return nil
}
func (f *Fixture) Providers() map[string]ui.CollectionProvider {
	return map[string]ui.CollectionProvider{ItemsPath: {ID: "commands-items", Epoch: 1, RootLoaded: true, Load: f.load, Initial: ui.CollectionData{Items: []ui.CollectionItem{
		{ID: "alpha", Kind: ui.Row, Label: "Alpha item", ChildrenLoaded: true},
		{ID: "beta", Kind: ui.Row, Label: "Beta item", ChildrenLoaded: true},
		{ID: "group", Kind: ui.Group, Label: "Lazy group", HasChildren: true},
	}}}}
}
