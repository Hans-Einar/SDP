package sdptool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
)

func deferredFixture(t *testing.T) string {
	t.Helper()
	root, _ := projectFixture(t)
	put := func(name, text string) {
		t.Helper()
		if e := os.WriteFile(filepath.Join(root, "SDP", name), []byte(text), 0600); e != nil {
			t.Fatal(e)
		}
	}
	put("Valid.design", sample)
	put("Invalid.design", "language design-core version 0.5.\nnot a declaration.\n")
	put("Unsupported.design", "language design-core version 9.9.\nunit Future.\n")
	put("Child.design", "language design-core version 0.6.\nunit Child.\n")
	return root
}
func compact(t *testing.T, root string) Project {
	t.Helper()
	p, e := DiscoverWithOptions(root, DiscoverOptions{Navigation: "lazy-sdl"})
	if e != nil {
		t.Fatal(e)
	}
	return p
}
func TestDeferredDiscoveryPreservesSourcesAndOtherTabs(t *testing.T) {
	root := deferredFixture(t)
	before, e := Discover(root)
	if e != nil {
		t.Fatal(e)
	}
	p := compact(t, root)
	if p.NavigationMode != "lazy-sdl" || !reflect.DeepEqual(p.Sources, before.Sources) || !reflect.DeepEqual(p.Inventory, before.Inventory) {
		t.Fatal("source set changed")
	}
	for _, n := range p.Navigation.Nodes {
		if n.Kind == "viewpoint" {
			t.Fatal("model graph eagerly included")
		}
		if n.Kind == "source" {
			var source *SourceInfo
			for i := range p.Sources {
				if p.Sources[i].ID == strings.TrimPrefix(n.ID, "sdl/") {
					source = &p.Sources[i]
					break
				}
			}
			if source == nil {
				t.Fatal(n)
			}
			if source.State == "validated" {
				if !n.ChildrenDeferred || len(n.Children) != 0 || n.Target == nil || n.Target.Revision != source.Revision || n.Target.Operation != "tree" {
					t.Fatal(n)
				}
				got, e := ExpandModel(root, ModelTreeOptions{Model: source.ID, Revision: source.Revision})
				if e != nil {
					t.Fatal(e)
				}
				want, e := ModelTree(before, source.ID)
				if e != nil {
					t.Fatal(e)
				}
				if !reflect.DeepEqual(got, want) {
					t.Fatal("expansion differs from full model tree")
				}
			} else if n.ChildrenDeferred || n.Target != nil || n.State != source.State || n.Diagnostic != source.Diagnostic {
				t.Fatal("invalid source hidden/actionable", n)
			}
		}
	}
	for _, n := range before.Navigation.Nodes {
		if n.ID == "sdl" || strings.HasPrefix(n.ID, "sdl/") {
			continue
		}
		found := false
		for _, x := range p.Navigation.Nodes {
			if x.ID == n.ID {
				found = reflect.DeepEqual(x, n)
			}
		}
		if !found {
			t.Fatal("other tab changed", n.ID)
		}
	}
	again, e := DiscoverWithOptions(root, DiscoverOptions{Navigation: "eager"})
	if e != nil {
		t.Fatal(e)
	}
	a, _ := json.Marshal(before)
	b, _ := json.Marshal(again)
	if !bytes.Equal(a, b) {
		t.Fatal("legacy payload changed")
	}
}

func TestDeferredProjectExceedsOldCombinedNodeLimit(t *testing.T) {
	root, _ := projectFixture(t)
	var text strings.Builder
	text.WriteString("language design-core version 0.5.\n")
	for i := 0; i < 300; i++ {
		fmt.Fprintf(&text, "unit U%03d.\nfunctionality F%03d.\n", i, i)
	}
	for i := 0; i < 300; i++ {
		fmt.Fprintf(&text, "U%03d owns F%03d.\n", i, i)
	}
	model, e := parser.Parse(text.String())
	if e != nil {
		t.Fatal(e)
	}
	source := []byte(parser.FormatFile(model))
	for i := 0; i < 12; i++ {
		if e = os.WriteFile(filepath.Join(root, "SDP", fmt.Sprintf("Model%02d.design", i)), source, 0600); e != nil {
			t.Fatal(e)
		}
	}
	p := compact(t, root)
	count := 0
	for _, m := range p.Sources {
		tree, e := ExpandModel(root, ModelTreeOptions{Model: m.ID, Revision: m.Revision})
		if e != nil {
			t.Fatal(e)
		}
		count += len(tree.Nodes)
	}
	if count <= 20000 {
		t.Fatal("fixture does not exceed old cap", count)
	}
	if _, e = Discover(root); e == nil || !strings.Contains(e.Error(), "combined inventory") {
		t.Fatal("legacy repro lost", e)
	}
	if len(p.Sources) != 12 || len(p.Navigation.Nodes) >= 20000 {
		t.Fatal("compact source loss or eager inventory")
	}
	t.Logf("all %d models reachable: %d expanded nodes, %d initial nodes", len(p.Sources), count, len(p.Navigation.Nodes))
}

func TestDeferredComposedFreshnessAndMissingSelectors(t *testing.T) {
	root, _ := projectFixture(t)
	base := filepath.Join(root, "SDP")
	header := "language design-core version 0.6.\n"
	for name, text := range map[string]string{"System.design": header + "system Demo.\nincludes \"Part.design\".\nDemo contains Child.\n", "Part.design": header + "unit Child.\n"} {
		if e := os.WriteFile(filepath.Join(base, name), []byte(text), 0600); e != nil {
			t.Fatal(e)
		}
	}
	p := compact(t, root)
	var o ModelTreeOptions
	for _, s := range p.Sources {
		if s.State == "validated" {
			o = ModelTreeOptions{Model: s.ID, Revision: s.Revision}
		}
	}
	if o.Model == "" {
		t.Fatal(p.Sources)
	}
	if _, e := ExpandModel(root, o); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(filepath.Join(base, "Part.design"), []byte(header+"unit Child.\nunit Extra.\n"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e := ExpandModel(root, o); e == nil || !strings.Contains(e.Error(), "stale") {
		t.Fatal("accepted changed include", e)
	}
	for _, bad := range []ModelTreeOptions{{}, {Model: o.Model}, {Revision: o.Revision}, {Model: "missing", Revision: o.Revision}} {
		if _, e := ExpandModel(root, bad); e == nil {
			t.Fatal("accepted missing selector", bad)
		}
	}
}
