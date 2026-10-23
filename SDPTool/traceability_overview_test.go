package sdptool

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestTraceabilityComposedOverviewExactSourcesAndFreshness(t *testing.T) {
	root, r := projectFixture(t)
	r.Models = []Model{{"model", "Demo", "System.design", "design-core/0.6"}}
	saveInventory(t, root, r)
	for name, src := range map[string]string{
		"System.design":        "language design-core version 0.6.\nsystem Demo.\ncontract Unused.\nincludes \"Parts_[draft].design\".\nDemo contains Owner.\nUnused has completeness = open.\n",
		"Parts_[draft].design": "language design-core version 0.6.\nunit Owner.\nactor Guest.\nusecase Orphan.\nfunctionality Work.\nfeature Feature.\nactivity Build.\nOwner owns Work.\nWork contributes-to Feature.\nWork has repeatability = deterministic.\nWork has state-retention = stateful.\nBuild addresses Work.\nBuild has implementation-status = verified.\n",
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
		if n.Kind == "viewpoint" && strings.HasSuffix(n.ID, "/VP11") {
			target = n.OverviewTarget
			if target == nil || n.Target == nil || target.URI == n.Target.URI {
				t.Fatal("lost target distinction", n)
			}
		}
	}
	if target == nil {
		t.Fatal("overview not advertised")
	}
	o := PreviewOptions{URI: target.URI, Revision: target.Revision, Output: filepath.Join(t.TempDir(), "view")}
	result, e := SelectProject(context.Background(), p, "model", o)
	if e != nil {
		t.Fatal(e)
	}
	data, e := os.ReadFile(result.Entry)
	if e != nil {
		t.Fatal(e)
	}
	text := string(data)
	for _, want := range []string{"Snapshot coverage", "Explicit properties", "Reported model gaps", "Declaration register", "Complete fact registry", "sdl-source://s0/", "sdl-source://s1/", "Parts\\_\\[draft\\].design", "UNSPECIFIED_ALLOCATION", "not proof of implemented behavior"} {
		if !strings.Contains(text, want) {
			t.Fatal(want, text)
		}
	}
	if strings.Contains(text, "```mermaid") || strings.Contains(text, "No diagrams match") || strings.Count(text, "## Complete fact registry") != 1 || len(result.SourceLinks) != 2 {
		t.Fatal(text)
	}
	v, _, e := loadModel(filepath.Join(root, "System.design"))
	if e != nil {
		t.Fatal(e)
	}
	var manifest documents.Manifest
	b, e := os.ReadFile(filepath.Join(result.Directory, "manifest.json"))
	if e != nil {
		t.Fatal(e)
	}
	if e = json.Unmarshal(b, &manifest); e != nil || len(manifest.Facts) != len(v.Facts) || len(manifest.Gaps) != len(v.Gaps) {
		t.Fatal("incomplete manifest", e)
	}
	sources := map[string][]byte{}
	for _, l := range result.SourceLinks {
		original, e := os.ReadFile(l.Path)
		if e != nil {
			t.Fatal(e)
		}
		saved, e := os.ReadFile(filepath.Join(result.Directory, l.Snapshot))
		if e != nil || string(saved) != string(original) || documents.Hash(saved) != l.Revision {
			t.Fatal("source mismatch", l, e)
		}
		sources[l.ID] = original
	}
	covered := map[string]bool{}
	for _, m := range regexp.MustCompile(`\[([^]]+)\]\(sdl-source://(s[0-9]+)/([0-9]+)/([0-9]+)\)`).FindAllStringSubmatch(text, -1) {
		b, ok := sources[m[2]]
		if !ok {
			t.Fatal(m)
		}
		a, _ := strconv.Atoi(m[3])
		z, _ := strconv.Atoi(m[4])
		if a < 0 || z <= a || z > len(b) {
			t.Fatal(m)
		}
		covered[string(b[a:z])] = true
	}
	for _, f := range v.Facts {
		if !covered[f.S("text")] || !strings.Contains(text, f.S("id")) {
			t.Fatal("fact lost source", f)
		}
	}
	for _, d := range v.Model.Declarations {
		if !covered[fmt.Sprintf("%s %s.", d.Kind, d.Name.Name)] {
			t.Fatal("declaration lost source", d)
		}
	}
	child := filepath.Join(root, "Parts_[draft].design")
	saved, _ := os.ReadFile(child)
	changed, e := parser.Parse(strings.Replace(string(saved), "version 0.6.\n", "version 0.6.\nunit Later.\n", 1))
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(child, []byte(parser.FormatFile(changed)), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e = SelectProject(context.Background(), p, "model", o); e == nil || !strings.Contains(e.Error(), "stale") {
		t.Fatal("stale target accepted", e)
	}
	after, e := os.ReadFile(result.Entry)
	if e != nil || string(after) != text {
		t.Fatal("stale selection replaced view", e)
	}
}
