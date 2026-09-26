package sdptool

import (
	"context"
	"encoding/json"
	"fmt"
	ui "github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/presentation"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"path/filepath"
	"sort"
)

func uiRoots(source string) (map[string]*ui.Instance, []byte, error) {
	b, e := readSource(source)
	if e != nil {
		return nil, nil, e
	}
	d, e := ui.Parse(string(b))
	if e != nil {
		return nil, nil, e
	}
	roots, e := ui.Normalize(d)
	return roots, b, e
}

// UIPreview delegates a static Markdown prototype to SDUI. No widget runtime or
// browser/native UI implementation is added by this facade.
func UIPreview(ctx context.Context, p Project, id, entry, output, revision string) (Result, error) {
	_, source, e := p.model(id, true)
	if e != nil {
		return Result{}, e
	}
	if output == "" {
		return Result{}, failure("arguments", fmt.Errorf("--output required"))
	}
	if e = outside(output, source); e != nil {
		return Result{}, failure("output", e)
	}
	roots, b, e := uiRoots(source)
	if e != nil {
		return Result{}, failure("model", e)
	}
	hash := documents.Hash(b)
	if revision != "" && hash != revision {
		return Result{}, failure("stale", fmt.Errorf("SDUI source changed"))
	}
	if entry == "" {
		names := []string{}
		for name, n := range roots {
			if n.Kind == "frame" {
				names = append(names, name)
			}
		}
		sort.Strings(names)
		if len(names) != 1 {
			return Result{}, failure("selection", fmt.Errorf("choose --entry for multiple/no root frames"))
		}
		entry = names[0]
	}
	root := roots[entry]
	if root == nil || root.Kind != "frame" {
		return Result{}, failure("selection", fmt.Errorf("entry must be a defined frame"))
	}
	text, e := presentation.Markdown(root, 160)
	if e != nil {
		return Result{}, failure("render", e)
	}
	if len(text) > 32<<20 {
		return Result{}, failure("limit", fmt.Errorf("SDUI output limit"))
	}
	if e = ctx.Err(); e != nil {
		return Result{}, failure("canceled", e)
	}
	current, e := readSource(source)
	if e != nil {
		return Result{}, e
	}
	if documents.Hash(current) != hash {
		return Result{}, failure("stale", fmt.Errorf("SDUI source changed during generation"))
	}
	out, e := filepath.Abs(output)
	if e != nil {
		return Result{}, e
	}
	result := Result{Version, "sdui-preview", source, "sdui/0.2", hash, filepath.Join(out, "entry.md"), out, "caller-owned; release after consumer"}
	bundle := &documents.Bundle{Files: map[string][]byte{}, Manifest: documents.Manifest{Version: "sdptool-sdui-markdown/0.1", Revision: hash}}
	bundle.Put("entry.md", text)
	j, _ := json.MarshalIndent(result, "", "  ")
	bundle.Files["sdptool.json"] = append(j, '\n')
	bundle.Seal()
	if e = bundle.Publish(out); e != nil {
		return Result{}, failure("output", e)
	}
	return result, nil
}
func UINodes(p Project) ([]Node, string, error) {
	tab := Node{ID: "sdui", Kind: "tab", Label: "SDUI", State: "absent"}
	nodes := []Node{}
	all := []byte{}
	for _, m := range p.Registration.SDUI {
		key := "sdui/" + m.ID
		tab.Children = append(tab.Children, key)
		tab.State = "available"
		n := Node{ID: key, Kind: "source", Label: m.ID, State: "unsupported"}
		if m.Profile == "sdui/0.2" {
			_, source, e := p.model(m.ID, true)
			if e != nil {
				return nil, "", e
			}
			roots, b, e := uiRoots(source)
			if e != nil {
				return nil, "", e
			}
			n.State = "validated"
			n.Target = &Target{Operation: "sdui-preview", Project: p.Registration.ProjectID, Model: m.ID, Path: source, Revision: documents.Hash(b)}
			names := []string{}
			for name, root := range roots {
				if root.Kind == "frame" {
					names = append(names, name)
				}
			}
			sort.Strings(names)
			for _, name := range names {
				child := key + "/frame/" + name
				n.Children = append(n.Children, child)
				nodes = append(nodes, Node{ID: child, Kind: "frame", Label: name, State: "validated", Target: &Target{Operation: "sdui-preview", Project: p.Registration.ProjectID, Model: m.ID, Entry: name, Path: source, Revision: documents.Hash(b)}})
			}
			all = append(all, []byte(m.ID+"\x00")...)
			all = append(all, b...)
		}
		nodes = append(nodes, n)
	}
	nodes = append(nodes, tab)
	return nodes, documents.Hash(all), nil
}
