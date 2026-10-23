package sdptool

import (
	"context"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestActivityOverviewComposedTargetAndSourceLinks(t *testing.T) {
	root, r := projectFixture(t)
	r.Models = []Model{{"model", "Demo", "System.design", "design-core/0.6"}}
	saveInventory(t, root, r)
	for name, src := range map[string]string{
		"System.design": "language design-core version 0.6.\nsystem Demo.\nactivity Plan.\nincludes \"Parts.design\".\nDemo contains Shell.\n",
		"Parts.design":  "language design-core version 0.6.\ncontainer Shell.\nunit Worker.\nunit Isolated.\nactivity Build.\nactivity Check.\nfunctionality Work.\nfeature Product.\nBuild refines Plan.\nBuild addresses Work.\nBuild delivers Product.\nBuild has implementation-status = implemented.\nCheck depends-on Build.\nWorker owns Work.\nShell contains Worker.\n",
	} {
		m, e := parser.Parse(src)
		if e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(filepath.Join(root, name), []byte(parser.FormatFile(m)), 0600); e != nil {
			t.Fatal(e)
		}
	}
	p, e := Discover(root)
	if e != nil {
		t.Fatal(e)
	}
	p.Inventory = r
	tree, e := ModelTree(p, "model")
	if e != nil {
		t.Fatal(e)
	}
	var target *Target
	for _, n := range tree.Nodes {
		if n.Kind == "viewpoint" && strings.HasSuffix(n.ID, "/VP06") {
			target = n.OverviewTarget
			if target == nil || n.Target == nil || target.URI == n.Target.URI {
				t.Fatal("lost target distinction", n)
			}
		}
	}
	if target == nil {
		t.Fatal("not advertised")
	}
	o := PreviewOptions{URI: target.URI, Revision: target.Revision, Output: filepath.Join(t.TempDir(), "view")}
	result, e := SelectProject(context.Background(), p, "model", o)
	if e != nil {
		t.Fatal(e)
	}
	b, e := os.ReadFile(result.Entry)
	if e != nil {
		t.Fatal(e)
	}
	text := string(b)
	for _, want := range []string{"All activities and explicit refinement", "Plan (activity)", "Build (activity)", "Check (activity)", "Explicit activity prerequisites", "Source-claimed status", "Addressed responsibilities", "implemented", "sdl-source://s0/", "sdl-source://s1/"} {
		if !strings.Contains(text, want) {
			t.Fatal(want, text)
		}
	}
	if strings.Count(text, "```mermaid") != 2 || strings.Count(text, "-->|refines|") != 1 || strings.Count(text, "-.->|depends-on|") != 1 || len(result.SourceLinks) != 2 {
		t.Fatal(text)
	}
	// Mutating a dependency invalidates the unchanged entry/typed target.
	child := filepath.Join(root, "Parts.design")
	saved, _ := os.ReadFile(child)
	changed, e := parser.Parse(strings.Replace(string(saved), "version 0.6.\n", "version 0.6.\nunit NewUnit.\n", 1))
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(child, []byte(parser.FormatFile(changed)), 0600); e != nil {
		t.Fatal(e)
	}
	o.Output = filepath.Join(t.TempDir(), "stale")
	if _, e = SelectProject(context.Background(), p, "model", o); e == nil || !strings.Contains(e.Error(), "stale") {
		t.Fatal(e)
	}
	if _, e = os.Stat(o.Output); !os.IsNotExist(e) {
		t.Fatal("stale selection published output", e)
	}
}
