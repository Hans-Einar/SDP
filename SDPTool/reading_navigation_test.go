package sdptool

import (
	"context"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readingFixture(t *testing.T) string {
	root, _ := projectFixture(t)
	for name, source := range map[string]string{
		"Base.design":   "language design-core version 0.6.\nsystem Demo.\nincludes \"Shared.design\".\ncontainer Runtime.\nDemo contains Runtime.\nRuntime contains Worker.\n",
		"Mini.design":   "language design-core version 0.6.\nsystem Demo.\nincludes \"Shared.design\".\ncontainer MiniRuntime.\nDemo contains MiniRuntime.\nMiniRuntime contains Worker.\n",
		"Shared.design": "language design-core version 0.6.\nunit Worker.\ninterface SilentPort.\n",
		"Orphan.design": "language design-core version 0.6.\nunit Alone.\nAlone consumes OutsidePort.\n",
	} {
		syntax, err := parser.Parse(source)
		if err != nil {
			t.Fatal(err)
		}
		if e := os.WriteFile(filepath.Join(root, "SDP", name), []byte(parser.FormatFile(syntax)), 0600); e != nil {
			t.Fatal(e)
		}
	}
	return root
}
func TestReadingMembershipCollectionsAndSources(t *testing.T) {
	root := readingFixture(t)
	p, e := DiscoverWithOptions(root, DiscoverOptions{Navigation: "lazy-sdl"})
	if e != nil {
		t.Fatal(e)
	}
	if p.Capabilities["readingNavigation"] != "1" || p.Navigation.Presentation != readingPresentation {
		t.Fatal("missing presentation contract")
	}
	var shared, base SourceInfo
	for _, s := range p.Sources {
		if strings.HasSuffix(s.Source, "Shared.design") {
			shared = s
		}
		if strings.HasSuffix(s.Source, "Base.design") {
			base = s
		}
	}
	if len(shared.IncludedBy) != 2 {
		t.Fatal("shared context lost", shared)
	}
	for _, n := range p.Navigation.Nodes {
		if n.ID == "sdl/"+shared.ID {
			t.Fatal("included source shown as orphan")
		}
	}
	tree, e := ExpandModel(root, ModelTreeOptions{Model: base.ID, Revision: base.Revision})
	if e != nil {
		t.Fatal(e)
	}
	var rootNode Node
	collections := map[string]bool{}
	objects := map[string]bool{}
	for _, n := range tree.Nodes {
		if n.ID == "sdl" {
			rootNode = n
		}
		if n.Kind == "diagram" || n.Kind == "relationship" {
			t.Fatal("duplicate inventory route", n)
		}
		if n.Kind == "collection" {
			if n.Target == nil {
				t.Fatal("non-readable collection", n)
			}
			collections[n.ID] = true
		}
		if n.Label == "SilentPort" && strings.Contains(n.ID, "/VP04/") {
			objects[n.Label] = true
		}
	}
	if !objects["SilentPort"] {
		t.Fatal("isolated interface omitted")
	}
	if len(rootNode.Children) != 12 {
		t.Fatal(rootNode)
	}
	for i, id := range rootNode.Children[:11] {
		if !strings.HasSuffix(id, []string{"VP01", "VP02", "VP03", "VP04", "VP05", "VP06", "VP07", "VP08", "VP09", "VP10", "VP11"}[i]) {
			t.Fatal("unordered viewpoints", rootNode.Children)
		}
	}
	target := sourceTarget(p, base.Model, base.Revision, "Shared.design")
	r, e := SelectProject(context.Background(), p, base.ID, PreviewOptions{URI: target.URI, Revision: base.Revision, Output: filepath.Join(t.TempDir(), "bundle")})
	if e != nil {
		t.Fatal(e)
	}
	raw, _ := os.ReadFile(r.Entry)
	if !strings.Contains(string(raw), "SilentPort") || !strings.Contains(string(raw), "sdl-source://") {
		t.Fatal(string(raw))
	}
	if e = os.WriteFile(filepath.Join(root, "SDP", "Shared.design"), []byte("language design-core version 0.6.\nunit Worker.\ninterface ChangedPort.\n"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e = SelectProject(context.Background(), p, base.ID, PreviewOptions{URI: target.URI, Revision: base.Revision, Output: filepath.Join(t.TempDir(), "stale")}); e == nil || !strings.Contains(e.Error(), "stale") {
		t.Fatal("changed sibling accepted", e)
	}
}
func TestReadingUnassignedSyntaxNeverClaimsValidation(t *testing.T) {
	p, e := Discover(readingFixture(t))
	if e != nil {
		t.Fatal(e)
	}
	for _, s := range p.Sources {
		if !strings.HasSuffix(s.Source, "Orphan.design") {
			continue
		}
		target := sourceTarget(p, s.Model, s.Revision, "")
		r, e := SelectProject(context.Background(), p, s.ID, PreviewOptions{URI: target.URI, Revision: s.Revision, Output: filepath.Join(t.TempDir(), "orphan")})
		if e != nil {
			t.Fatal(e)
		}
		raw, _ := os.ReadFile(r.Entry)
		if !strings.Contains(string(raw), "Local syntax only") || !strings.Contains(string(raw), "OutsidePort") {
			t.Fatal(string(raw))
		}
		return
	}
	t.Fatal("missing orphan")
}
