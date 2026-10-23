package sdptool

import (
	"context"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDataOverviewComposedTargetAndSourceLinks(t *testing.T) {
	root, r := projectFixture(t)
	r.Models = []Model{{"model", "Demo", "System.design", "design-core/0.6"}}
	saveInventory(t, root, r)
	for name, src := range map[string]string{
		"System.design": "language design-core version 0.6.\nsystem Demo.\ncontract Unused.\nincludes \"Parts.design\".\nDemo contains Owner.\nUnused has completeness = open.\n",
		"Parts.design":  "language design-core version 0.6.\nunit Owner.\ndatabase Store.\ndatabase Idle.\ndataset State.\ndataset Temporary.\ndatagram Notice.\ncontract Record.\ncontract Wire.\nvariant Change.\nfield Value.\nfield Payload.\nfunctionality Project.\nfeature Feature.\nmode Online.\nOwner owns Store.\nOwner owns Idle.\nOwner owns Project.\nOwner holds Temporary.\nStore holds State.\nState upholds Record.\nTemporary upholds Record.\nNotice from State.\nNotice upholds Wire.\nRecord has completeness = closed.\nRecord has-field Value.\nValue has presence = required.\nValue has value-type = text.\nWire has completeness = closed.\nWire defines Change.\nChange has-field Payload.\nPayload has presence = optional.\nPayload has value-type = bytes.\nProject projects State into Notice.\nProject contributes-to Feature.\n",
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
		if n.Kind == "viewpoint" && strings.HasSuffix(n.ID, "/VP09") {
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
	for _, want := range []string{"Data definitions", "Database-backed persistence", "Unit-held; persistence unspecified", "Contract register", "Variant register", "Field register", "Projection responsibilities", "sdl-source://s0/", "sdl-source://s1/"} {
		if !strings.Contains(text, want) {
			t.Fatal(want, text)
		}
	}
	if strings.Count(text, "```mermaid") != 2 || strings.Count(text, "-->|from|") != 1 || strings.Count(text, "-->|holds|") != 2 || len(result.SourceLinks) != 2 {
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
