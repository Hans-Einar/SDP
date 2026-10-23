package sdptool

import (
	"context"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOverviewTargetGenerationAndFreshness(t *testing.T) {
	root, r := projectFixture(t)
	r.Models = []Model{{"model", "System", "model.design", "design-core/0.5"}}
	saveInventory(t, root, r)
	source := "language design-core version 0.5.\nactor User.\nusecase Browse.\nusecase Isolated.\nUser pursues Browse.\n"
	m, err := parser.Parse(source)
	if err != nil {
		t.Fatal(err)
	}
	source = parser.FormatFile(m)
	path := filepath.Join(root, "model.design")
	os.WriteFile(path, []byte(source), 0600)
	p, _ := Discover(root)
	p.Inventory = r
	tree, e := ModelTree(p, "model")
	if e != nil {
		t.Fatal(e)
	}
	var overview *Target
	for _, n := range tree.Nodes {
		if n.OverviewTarget != nil {
			if (!strings.HasSuffix(n.ID, "/VP01") && !strings.HasSuffix(n.ID, "/VP07")) || n.Target.URI == n.OverviewTarget.URI {
				t.Fatal(n)
			}
			if strings.HasSuffix(n.ID, "/VP01") {
				overview = n.OverviewTarget
			}
		}
	}
	if overview == nil {
		t.Fatal("missing capability")
	}
	o := PreviewOptions{URI: overview.URI, Revision: overview.Revision, Output: filepath.Join(t.TempDir(), "selected")}
	result, e := SelectProject(context.Background(), p, "model", o)
	if e != nil {
		t.Fatal(e)
	}
	b, e := os.ReadFile(result.Entry)
	if e != nil {
		t.Fatal(e)
	}
	for _, s := range []string{"usecase-beta", "Isolated", "Functionality realization", "pursues"} {
		if !strings.Contains(string(b), s) {
			t.Fatal(s, string(b))
		}
	}
	if strings.Count(string(b), "```mermaid") != 1 {
		t.Fatal("not one overview")
	}
	o.Output = filepath.Join(t.TempDir(), "foreign")
	o.URI = strings.Replace(overview.URI, "sdl-view://trial/", "sdl-view://foreign/", 1)
	if _, e = SelectProject(context.Background(), p, "model", o); e == nil {
		t.Fatal("foreign accepted")
	}
	o.URI = overview.URI
	m, err = parser.Parse(strings.Replace(source, "version 0.5.\n", "version 0.5.\nactor NewActor.\n", 1))
	if err != nil {
		t.Fatal(err)
	}
	os.WriteFile(path, []byte(parser.FormatFile(m)), 0600)
	if _, e = SelectProject(context.Background(), p, "model", o); e == nil || !strings.Contains(e.Error(), "stale") {
		t.Fatal(e)
	}
	if _, e = os.Stat(o.Output); !os.IsNotExist(e) {
		t.Fatal("invalid request published output", e)
	}
}

func TestFeatureOverviewTargetsGenerateOnlySelectedMode(t *testing.T) {
	root, r := projectFixture(t)
	r.Models = []Model{{"model", "System", "model.design", "design-core/0.5"}}
	saveInventory(t, root, r)
	m, e := parser.Parse(`language design-core version 0.5.
unit Owner.
container Local.
container Remote.
functionality Work.
feature Feature.
mode Online.
mode Offline.
Owner owns Work.
Work contributes-to Feature.
Work allocated-to Local in mode Offline.
Work allocated-to Remote in mode Online.
`)
	if e != nil {
		t.Fatal(e)
	}
	os.WriteFile(filepath.Join(root, "model.design"), []byte(parser.FormatFile(m)), 0600)
	p, _ := Discover(root)
	p.Inventory = r
	tree, e := ModelTree(p, "model")
	if e != nil {
		t.Fatal(e)
	}
	modes := 0
	logical := false
	for _, n := range tree.Nodes {
		var action *Target
		if strings.HasSuffix(n.ID, "/VP07") {
			action = n.OverviewTarget
			logical = action != nil
		}
		if strings.Contains(n.ID, "/overview-mode/") {
			action = n.Target
			modes++
		}
		if action == nil {
			continue
		}
		result, e := SelectProject(context.Background(), p, "model", PreviewOptions{URI: action.URI, Revision: action.Revision, Output: filepath.Join(t.TempDir(), "view")})
		if e != nil {
			t.Fatal(e)
		}
		b, e := os.ReadFile(result.Entry)
		if e != nil {
			t.Fatal(e)
		}
		text := string(b)
		if strings.Count(text, "```mermaid") != 1 || !strings.Contains(text, "Feature contributions") {
			t.Fatal(text)
		}
		if strings.Contains(action.URI, "mode=Offline") && (!strings.Contains(text, "n_Local[\"Local (container) | allocated-to: F01") || strings.Contains(text, "n_Remote[\"Remote (container) | allocated-to: F01")) {
			t.Fatal("mixed mode", text)
		}
		if strings.Contains(action.URI, "mode=Online") && (!strings.Contains(text, "n_Remote[\"Remote (container) | allocated-to: F01") || strings.Contains(text, "n_Local[\"Local (container) | allocated-to: F01")) {
			t.Fatal("mixed mode", text)
		}
	}
	if !logical || modes != 2 {
		t.Fatal("missing lazy targets", logical, modes)
	}
}
