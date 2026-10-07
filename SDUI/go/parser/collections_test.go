package parser

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestCollectionProfilesAndSchemas(t *testing.T) {
	for _, widget := range []string{"tree", "list"} {
		for _, args := range []string{`"Navigation"`, `label="Navigation"`, `"Navigation", callback=nav.Activate.@invoke`} {
			source := `sdui 0.3; ref: nav "never-opened.sdl"; page=[items=` + widget + `(` + args + `)];`
			d, roots, err := Compile(source)
			if err != nil {
				t.Fatal(source, err)
			}
			if d.Profile != "sdui/0.3" || roots["page"].Rows[0][0].Argument("label") != "Navigation" {
				t.Fatal("profile/label lost")
			}
			legacy := strings.Replace(source, "0.3", "0.2", 1)
			if _, err = Parse(legacy); err != nil {
				t.Fatal("legacy syntax-only changed", err)
			}
			if _, _, err = Compile(legacy); err == nil {
				t.Fatal("legacy profile accepted collection")
			}
		}
		for _, args := range []string{``, `""`, `" \t\n"`, `label=1`, `label=false`, `label=null`, `label=nav.Action.@invoke`, `"A", "B"`, `label="A", "B"`, `"A", label="B"`, `"A", callback="bad"`, `"A", items="data"`, `"A", onSelect=nav.Activate.@invoke`, `"A", callback=missing.Activate.@invoke`} {
			s := `sdui 0.3; ref: nav "unused"; page=[items=` + widget + `(` + args + `)];`
			if _, _, err := Compile(s); err == nil {
				t.Error("accepted invalid collection", s)
			}
		}
		if _, _, err := Compile(`sdui 0.3; ref: nav "unused"; page=[` + widget + `("A", callback=nav.Activate.@invoke)];`); err == nil {
			t.Fatal("anonymous callback accepted")
		}
	}
	for _, v := range []string{"0.30", "3e-1", "0.4", "0.20", "2e-1"} {
		if _, err := Parse("sdui " + v + "; page=[];"); err == nil {
			t.Error("accepted nonexact profile", v)
		}
	}
	d, _, _ := Compile(`sdui 0.2; page=[];`)
	for _, profile := range []string{"", "sdui/0.30", "sdui/0.4"} {
		d.Profile = profile
		if err := Validate(d); err == nil {
			t.Error("accepted constructed document profile", profile)
		}
		if _, err := ASTFormat(profile); err == nil {
			t.Error("accepted AST format profile", profile)
		}
	}
}

func TestCollectionReuseProfileAndSpans(t *testing.T) {
	source := "sdui 0.3;\nref: nav \"unused\";\nleaf=<items=tree(\"Æ navigation\", callback=nav.Activate.@invoke)>;\nalias=leaf;\npage=[left=alias,right=alias];"
	doc, roots, err := Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	root := roots["page"]
	if p, err := EffectiveProfile(root); err != nil || p != "sdui/0.3" {
		t.Fatal(p, err)
	}
	root.Walk(func(n *Instance) {
		if n.Profile != "sdui/0.3" {
			t.Fatalf("missing profile: %+v", n)
		}
	})
	left, right := root.Rows[0][0], root.Rows[0][1]
	if left.Path == right.Path || len(left.Uses) != 2 || left.Uses[1].Span == right.Uses[1].Span {
		t.Fatal("reuse identity/provenance lost")
	}
	item := left.Rows[0][0]
	if !strings.HasPrefix(source[item.Span.Start:item.Span.End], `items=tree("Æ navigation"`) {
		t.Fatal(item.Span)
	}
	label := item.Arguments["label"].(Literal)
	if source[label.Span.Start:label.Span.End] != `"Æ navigation"` || label.Span.Line != 3 {
		t.Fatal(label)
	}
	left.Profile = ""
	if _, err := EffectiveProfile(root); err == nil {
		t.Fatal("accepted mixed profile")
	}
	again, err := Normalize(doc)
	if err != nil || again["page"].Rows[0][0].Profile != "sdui/0.3" || right.Profile != "sdui/0.3" {
		t.Fatal("profile mutations leaked", err)
	}
}

func TestInstanceProfileSerializationAndMalformedTrees(t *testing.T) {
	for _, version := range []string{"0.2", "0.3"} {
		doc, roots, err := Compile("sdui " + version + "; page=[button(\"OK\")];")
		if err != nil {
			t.Fatal(err)
		}
		format, err := ASTFormat(doc.Profile)
		if err != nil || format != "sdui-ast/"+version {
			t.Fatal(format, err)
		}
		for _, value := range []any{roots["page"], Data(roots["page"])} {
			data, err := json.Marshal(value)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(data), `"profile"`) != (version == "0.3") {
				t.Fatal(string(data))
			}
		}
		// Data's older non-omitempty behavior is intentionally preserved.
		data := Data(roots["page"].Rows[0][0]).(map[string]any)
		if _, present := data["declaration"]; !present {
			t.Fatal("old fields newly omitted")
		}
	}
	for _, bad := range []*Instance{nil, {Profile: "sdui/0.2"}, {Profile: "sdui/0.4"}, {Rows: [][]*Instance{{nil}}}} {
		if _, err := EffectiveProfile(bad); err == nil {
			t.Fatal("accepted malformed root", bad)
		}
	}
	cycle := &Instance{}
	cycle.Rows = [][]*Instance{{cycle}}
	if _, err := EffectiveProfile(cycle); err == nil {
		t.Fatal("accepted cycle")
	}
}
