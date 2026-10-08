package parser

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestInputOptionsPresenceAndLegacyDefaults(t *testing.T) {
	for _, profile := range []string{"0.2", "0.3"} {
		_, roots, err := Compile(`sdui ` + profile + `; Main=[field=input("Label",value="old")];`)
		if err != nil {
			t.Fatal(err)
		}
		n := roots["Main"].Rows[0][0]
		before, _ := json.Marshal(n)
		p, err := InputOptions(n)
		if err != nil || p != (InputPolicy{Placeholder: "Label"}) {
			t.Fatal(p, err)
		}
		after, _ := json.Marshal(n)
		if string(before) != string(after) || len(n.Arguments) != 2 {
			t.Fatal("accessor mutated legacy source")
		}
		if profile == "0.2" && n.Profile != "" {
			t.Fatal("lost legacy sentinel")
		}
	}
	for _, arg := range []string{`multiline=false`, `readOnly=false`, `required=false`, `placeholder=""`} {
		_, roots, err := Compile(`sdui 0.3; Main=[field=input("Label",` + arg + `)];`)
		if err != nil {
			t.Fatal(err)
		}
		n := roots["Main"].Rows[0][0]
		p, err := InputOptions(n)
		if err != nil || !p.Extended || p.Multiline || p.ReadOnly || p.Required || len(n.Arguments) != 2 {
			t.Fatal(arg, p, err)
		}
		if arg == `placeholder=""` {
			if !p.PlaceholderSet || p.Placeholder != "" {
				t.Fatal(p)
			}
		} else if p.PlaceholderSet || p.Placeholder != "Label" {
			t.Fatal(p)
		}
		if _, _, err := Compile(`sdui 0.2; Main=[field=input("Label",` + arg + `)];`); err == nil {
			t.Fatal("legacy admitted", arg)
		}
	}
}

func TestExtendedInputSchemasAndReuse(t *testing.T) {
	const source = `sdui 0.3; ref: app "not-opened.sdl";
Field=input("",value="å🙂\nnext",multiline=true,readOnly=false,required=true,placeholder="",callback=app.Save.@invoke);
Main=[left=Field;right=Field;form=dialog("Form")[note=input("Notes",required=false)]];
app.Save.setHandle(Main.left);`
	doc, roots, err := Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	root := roots["Main"]
	left, right := root.Rows[0][0], root.Rows[1][0]
	p, err := InputOptions(left)
	if err != nil || p != (InputPolicy{Extended: true, Multiline: true, Required: true, PlaceholderSet: true}) {
		t.Fatal(p, err)
	}
	if left.Argument("value") != "å🙂\nnext" || left.Path == right.Path || len(left.Uses) != 1 || left.Uses[0].Span == right.Uses[0].Span {
		t.Fatal("text/reuse identity lost")
	}
	if !reflect.DeepEqual(left.Arguments, right.Arguments) || len(left.Arguments) != 7 {
		t.Fatal("default/derived fields or span drift")
	}
	lit := left.Arguments["placeholder"].(Literal)
	if source[lit.Span.Start:lit.Span.End] != `""` {
		t.Fatal("original span lost", lit)
	}
	ids, err := ResolveInteractions(root)
	if err != nil || ids["Main/form/note"].Dialog != "Main/form" {
		t.Fatal(ids, err)
	}
	if _, err := ResolveDialogField(root, "Main/form", "note"); err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(Data(doc))
	if !strings.Contains(string(data), `"kind":"boolean"`) || strings.Contains(string(data), `InputPolicy`) {
		t.Fatal(string(data))
	}
	left.Arguments["placeholder"] = Literal{Kind: "string", Value: "changed"}
	if right.Argument("placeholder") != "" {
		t.Fatal("reuse aliases arguments")
	}
}

func TestExtendedInputClosedSchemaAndConstructedModels(t *testing.T) {
	for _, body := range []string{
		`input(multiline=true)`, `input("I",multiline="true")`, `input("I",readOnly=0)`,
		`input("I",required=null)`, `input("I",placeholder=false)`, `input("I",value=true)`,
		`input("I",multiline=true,multiline=false)`, `input(text="I","other")`,
		`input("I",onChange=app.Save.@invoke)`, `input("I",wrap=true)`, `input("I",enabled=true)`,
		`input("I",multiline=true)[]`, `input("I",callback="Save")`,
	} {
		if _, _, err := Compile(`sdui 0.3; ref: app "unused"; Main=[field=` + body + `];`); err == nil {
			t.Fatal("admitted", body)
		}
	}
	if _, _, err := Compile(`sdui 0.3; ref: app "unused"; Main=[input("I",multiline=true,callback=app.Save.@invoke)];`); err == nil {
		t.Fatal("anonymous callback")
	}
	// Empty required text is runtime validation, not a parser grammar failure.
	_, roots, err := Compile(`sdui 0.3; Main=[field=input("",required=true,value="")];`)
	if err != nil {
		t.Fatal(err)
	}
	n := roots["Main"].Rows[0][0]
	span := n.Arguments["required"].(Literal).Span
	for _, bad := range []any{Literal{Kind: "boolean", Value: "true", Span: span}, Literal{Kind: "string", Value: true, Span: span}, Reference{Span: span}, true} {
		n.Arguments["required"] = bad
		if _, err := InputOptions(n); err == nil {
			t.Fatal("malformed normalized argument accepted", bad)
		}
		if _, err := ResolveInteractions(roots["Main"]); err == nil {
			t.Fatal("strict selected-root validation bypass")
		}
	}
	if _, err := InputOptions(nil); err == nil {
		t.Fatal("nil accepted")
	}
	n.Arguments["required"] = Literal{Kind: "boolean", Value: false}
	n.Profile = "sdui/0.2"
	if _, err := InputOptions(n); err == nil {
		t.Fatal("explicit legacy profile instead of sentinel")
	}
	n.Profile = ""
	if _, err := InputOptions(n); err == nil {
		t.Fatal("constructed legacy new args")
	}
	n.Profile = "sdui/0.3"
	n.Arguments["unknown"] = Literal{Kind: "string", Value: "x"}
	if _, err := InputOptions(n); err == nil {
		t.Fatal("unknown key")
	}
	delete(n.Arguments, "unknown")
	n.Rows = [][]*Instance{{{Profile: "sdui/0.3", Kind: "frame", Path: n.Path + "/child"}}}
	if _, err := InputOptions(n); err == nil {
		t.Fatal("body accepted")
	}
}

func TestExtendedInputConstructedSourceAST(t *testing.T) {
	doc, err := Parse(`sdui 0.3; Main=[field=input("",required=false)];`)
	if err != nil {
		t.Fatal(err)
	}
	n := doc.Definitions[0].Root.Rows[0].Items[0]
	original := n.Arguments[1].Value
	span := original.(Literal).Span
	n.Arguments[1].Value = Literal{Kind: "boolean", Value: "false", Span: span}
	if _, err := Normalize(doc); err == nil {
		t.Fatal("malformed typed AST admitted")
	} else if d, ok := err.(*Diagnostic); !ok || d.Span != span {
		t.Fatal("lost bad literal span", err)
	}
	n.Arguments[1].Value = original
	n.Rows = []Row{{Items: []*Node{{Kind: "frame"}}}}
	if _, err := Normalize(doc); err == nil {
		t.Fatal("constructed input body admitted")
	}
}
