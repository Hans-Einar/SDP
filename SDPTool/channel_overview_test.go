package sdptool

import (
	"context"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestChannelOverviewComposedTargetAndSourceLinks(t *testing.T) {
	root, r := projectFixture(t)
	r.Models = []Model{{"model", "Demo", "System.design", "design-core/0.6"}}
	saveInventory(t, root, r)
	for name, src := range map[string]string{
		"System.design": "language design-core version 0.6.\nsystem Demo.\nmode Online.\nmode Offline.\nincludes \"Parts.design\".\nDemo contains Shell.\n",
		"Parts.design":  "language design-core version 0.6.\ncontainer Shell.\nunit Worker.\nunit Remote.\nchannel Bus.\nchannel Idle.\ncontract Events.\ncontract Body.\nmessage Notice.\nfield Value.\nShell contains Worker.\nBody has completeness = closed.\nBody has-field Value.\nValue has presence = required.\nValue has value-type = text.\nEvents has completeness = closed.\nBus upholds Events.\nIdle upholds Events.\nEvents permits Notice.\nNotice upholds Body.\nNotice has message-kind = event.\nWorker uses Bus as sender of Notice in mode Online.\nShell uses Bus as receiver of Notice in mode Online.\nRemote uses Bus as sender of Notice in mode Offline.\nShell uses Bus as receiver of Notice in mode Offline.\n",
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
	modes := 0
	for _, n := range tree.Nodes {
		if strings.Contains(n.ID, "/VP08/overview-mode/") {
			modes++
			if n.Target == nil {
				t.Fatal(n)
			}
		}
		if n.Kind == "viewpoint" && strings.HasSuffix(n.ID, "/VP08") {
			target = n.OverviewTarget
			if target == nil || n.Target == nil || target.URI == n.Target.URI {
				t.Fatal("lost target distinction", n)
			}
		}
	}
	if target == nil || modes != 2 {
		t.Fatal("not advertised")
	}
	o := PreviewOptions{URI: target.URI + "&mode=Online", Revision: target.Revision, Output: filepath.Join(t.TempDir(), "view")}
	result, e := SelectProject(context.Background(), p, "model", o)
	if e != nil {
		t.Fatal(e)
	}
	b, e := os.ReadFile(result.Entry)
	if e != nil {
		t.Fatal(e)
	}
	text := string(b)
	for _, want := range []string{"Channel Bus — mode Online", "Bus (channel)", "Worker (unit)", "Participation register", "Selected mode: **Online**", "sdl-source://s0/", "sdl-source://s1/"} {
		if !strings.Contains(text, want) {
			t.Fatal(want, text)
		}
	}
	if strings.Count(text, "```mermaid") != 2 || strings.Count(text, "-->|sender|") != 1 || len(result.SourceLinks) != 2 {
		t.Fatal(text)
	}
	if strings.Contains(text, "OtherPort") || strings.Contains(text, "Remote (unit)") {
		t.Fatal("mode leak", text)
	}
	index, e := SelectProject(context.Background(), p, "model", PreviewOptions{URI: target.URI, Revision: target.Revision, Output: filepath.Join(t.TempDir(), "index")})
	if e != nil {
		t.Fatal(e)
	}
	ib, e := os.ReadFile(index.Entry)
	if e != nil {
		t.Fatal(e)
	}
	if strings.Contains(string(ib), "```mermaid") || !strings.Contains(string(ib), "Choose a mode") {
		t.Fatal(string(ib))
	}
	// Mutating an included source invalidates the unchanged entry/typed target.
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
