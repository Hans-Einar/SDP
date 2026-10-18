package parser

import (
	"os"
	"testing"
)

func TestAllFamilySelectedSourceInventory(t *testing.T) {
	source, err := os.ReadFile("testdata/wci4-all-families.sdui")
	if err != nil {
		t.Fatal(err)
	}
	doc, roots, err := Compile(string(source))
	if err != nil {
		t.Fatal(err)
	}
	root := roots["Main"]
	if _, err := ResolveInteractions(root); err != nil {
		t.Fatal(err)
	}
	profile, err := EffectiveProfile(root)
	format, formatErr := ASTFormat(doc.Profile)
	if err != nil || formatErr != nil || profile != "sdui/0.3" || format != "sdui-ast/0.3" {
		t.Fatal(profile, err)
	}
	counts := map[string]int{}
	explicit := 0
	root.Walk(func(n *Instance) {
		kind := n.Widget
		if kind == "" {
			kind = n.Kind
		}
		counts[kind]++
		if n.Kind == "markdown" || n.Widget == "svg" {
			p, err := PreviewOptions(n)
			if err != nil {
				t.Fatal(err)
			}
			if p.Explicit {
				explicit++
			}
		}
	})
	for _, kind := range []string{"frame", "group", "button", "input", "tree", "list", "tabs", "page", "split", "command", "menu", "menuGroup", "item", "separator", "dialog", "checkbox", "slider", "select", "number", "svg", "markdown"} {
		if counts[kind] == 0 {
			t.Error("missing family", kind)
		}
	}
	if explicit != 3 || counts["svg"] != 2 || counts["dialog"] != 2 {
		t.Fatal(counts, explicit)
	}
}
