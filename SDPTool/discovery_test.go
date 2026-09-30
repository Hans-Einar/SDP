package sdptool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSourceDiscoverySnapshotAndRefresh(t *testing.T) {
	root, _ := projectFixture(t)
	put := func(rel, text string) {
		t.Helper()
		path := filepath.Join(root, "SDP", rel)
		if e := os.MkdirAll(filepath.Dir(path), 0700); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(path, []byte(text), 0600); e != nil {
			t.Fatal(e)
		}
	}
	put("SDL/Eco/One/System.design", "language design-core version 0.6.\nsystem One.\nOne contains Parts/Worker.\n")
	put("SDL/Eco/One/Parts/Worker.design", "language design-core version 0.6.\nunit Worker.\n")
	put("SDL/Eco/Two/System.design", "language design-core version 0.6.\nsystem Two.\n")
	put("SDL/Old/Desktop.design", sample)
	put("SDL/Old/Navigation.design", sample)
	put("SDL/Old/SDUI/Desktop.sdui", `sdui 0.2; Page=[<"Hello", button("OK")>];`)
	put("SDL/Old/Broken.design", "language design-core version 0.5.\ninvalid syntax")
	put("SDL/Old/Future.design", "language design-core version 99.0.\n")
	put("SDL/Old/SDUI/Broken.sdui", `sdui 0.2; bad`)
	p, e := Discover(root)
	if e != nil {
		t.Fatal(e)
	}
	if len(p.Sources) != 9 || p.Navigation == nil || p.Capabilities["sdui"] != "discovered" {
		t.Fatalf("%+v", p)
	}
	states := map[string]string{}
	ids := map[string]bool{}
	for _, s := range p.Sources {
		states[s.Source] = s.State
		if ids[s.ID] {
			t.Fatal("ID collision")
		}
		ids[s.ID] = true
	}
	for path, want := range map[string]string{"SDL/Eco/One/System.design": "validated", "SDL/Eco/One/Parts/Worker.design": "context-required", "SDL/Eco/Two/System.design": "validated", "SDL/Old/Broken.design": "invalid", "SDL/Old/Future.design": "unsupported", "SDL/Old/SDUI/Broken.sdui": "invalid"} {
		if states["SDP/"+path] != want {
			t.Fatalf("%s: %s", path, states["SDP/"+path])
		}
	}
	before := p.Navigation.InventoryRevision
	q, e := Discover(filepath.Join(root, "SDP"))
	if e != nil || before != q.Navigation.InventoryRevision {
		t.Fatalf("unstable snapshot %v", e)
	}
	if _, e := os.Stat(filepath.Join(root, "SDP/navigation.json")); !os.IsNotExist(e) {
		t.Fatal("index written")
	}
	put("SDL/Eco/One/Parts/Worker.design", "language design-core version 0.6.\nunit Extra.\nunit Worker.\n")
	q, e = Discover(root)
	if e != nil || before == q.Navigation.InventoryRevision {
		t.Fatal("edit not reflected", e)
	}
	var oldRoot, newRoot SourceInfo
	for _, s := range p.Sources {
		if strings.HasSuffix(s.Source, "One/System.design") {
			oldRoot = s
		}
	}
	for _, s := range q.Sources {
		if s.ID == oldRoot.ID {
			newRoot = s
		}
	}
	if oldRoot.Revision == newRoot.Revision || oldRoot.ID != newRoot.ID {
		t.Fatal("dependency revision or stable identity lost")
	}
	os.Rename(filepath.Join(root, "SDP/SDL/Old/Navigation.design"), filepath.Join(root, "SDP/SDL/Old/Renamed.design"))
	os.Remove(filepath.Join(root, "SDP/SDL/Old/Desktop.design"))
	put("SDL/New.design", sample)
	refreshed, e := Discover(root)
	if e != nil {
		t.Fatal(e)
	}
	b, _ := json.Marshal(refreshed)
	if strings.Contains(string(b), "SDL/Old/Navigation.design") || strings.Contains(string(b), "SDL/Old/Desktop.design") || !strings.Contains(string(b), "SDL/Old/Renamed.design") {
		t.Fatal("rename/delete not reflected")
	}
	// All navigation child/reference edges resolve even with independent Systems.
	known := map[string]bool{}
	for _, n := range refreshed.Navigation.Nodes {
		if known[n.ID] {
			t.Fatal("duplicate node", n.ID)
		}
		known[n.ID] = true
	}
	for _, n := range refreshed.Navigation.Nodes {
		for _, child := range n.Children {
			if !known[child] {
				t.Fatal("dangling child", child)
			}
		}
		if n.Reference != "" && !known[n.Reference] {
			t.Fatal("dangling reference", n.Reference)
		}
	}
}

func TestDiscoverySourceLimit(t *testing.T) {
	root, _ := projectFixture(t)
	for i := 0; i < 257; i++ {
		if e := os.WriteFile(filepath.Join(root, "SDP", fmt.Sprintf("source-%03d.design", i)), []byte("language design-core version 99.0.\n"), 0600); e != nil {
			t.Fatal(e)
		}
	}
	if _, e := Discover(root); e == nil || !strings.Contains(e.Error(), "256 source") {
		t.Fatalf("source limit: %v", e)
	}
}

func TestDiscoveryUsesOwningParserForProfile(t *testing.T) {
	root, _ := projectFixture(t)
	path := filepath.Join(root, "SDP", "Comments.sdui")
	data := []byte("sdui # profile follows\n0.2; Page=[];")
	if e := os.WriteFile(path, data, 0600); e != nil {
		t.Fatal(e)
	}
	p, e := Discover(root)
	if e != nil {
		t.Fatal(e)
	}
	if len(p.Sources) != 1 || p.Sources[0].State != "validated" || p.Sources[0].Profile != "sdui/0.2" {
		t.Fatal(p.Sources)
	}
}

func TestSelectedModelAvoidsCombinedExpansionLimit(t *testing.T) {
	root, _ := projectFixture(t)
	for i := 0; i < 255; i++ {
		name := fmt.Sprintf("M%03d.design", i)
		data := "language design-core version 0.5.\n"
		for j := 0; j < 10; j++ {
			data += fmt.Sprintf("unit U%d.\n", j)
		}
		if e := os.WriteFile(filepath.Join(root, "SDP", name), []byte(data), 0600); e != nil {
			t.Fatal(e)
		}
	}
	if _, e := Discover(root); e == nil {
		t.Fatal("expected aggregate limit")
	}
	var out, errs bytes.Buffer
	if code := Run(context.Background(), []string{root, "tree", "--model", sourceID("SDP/M000.design"), "--json"}, &out, &errs); code != 0 {
		t.Fatal(code, errs.String())
	}
	var result Tree
	if e := json.Unmarshal(out.Bytes(), &result); e != nil || result.Model != sourceID("SDP/M000.design") {
		t.Fatal(e, result.Model)
	}
}
