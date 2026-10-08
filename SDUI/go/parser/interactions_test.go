package parser

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func interactionsRoot(t *testing.T, source string) *Instance {
	t.Helper()
	_, roots, err := Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	return roots["Main"]
}
func TestInteractionReuseIdentityAndGeneratedScope(t *testing.T) {
	source := `sdui 0.3;
Part=[local=command("Local",toggle=true,checked=true,exclusive="choice"); <button(command="local"),button(command="/shared")>];
Alias=Part;
Main=[shared=command("Shared",toggle=true);left=Alias;right=Part;
dialogue=dialog("Edit",modal=false)[<name=input("Name")>;nested=dialog("Nested")[inside=input("Inside")];ok=button("OK",effect="accept")]];`
	root := interactionsRoot(t, source)
	before, _ := json.Marshal(root)
	identities, err := ResolveInteractions(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, side := range []string{"left", "right"} {
		base := "Main/" + side
		if got := identities[base+"/$r1c0/$r0c0"]; got.Command != base+"/local" || got.Scope != base {
			t.Fatalf("relative reuse: %+v", got)
		}
		if got := identities[base+"/$r1c0/$r0c1"]; got.Command != "Main/shared" {
			t.Fatalf("absolute reuse: %+v", got)
		}
	}
	if identities["Main/dialogue/ok"].Dialog != "Main/dialogue" || identities["Main/dialogue/nested"].Dialog != "Main/dialogue" || identities["Main/dialogue/nested/inside"].Dialog != "Main/dialogue/nested" {
		t.Fatal(identities)
	}
	field, err := ResolveDialogField(root, "Main/dialogue", "name")
	if err != nil || field.Path != "Main/dialogue/$r0c0/name" {
		t.Fatal(field, err)
	}
	for _, path := range []string{"nested/inside", "/name", "../name", "$r0c0/name", "name[0]", "missing"} {
		if _, err := ResolveDialogField(root, "Main/dialogue", path); err == nil {
			t.Fatal("accepted field", path)
		}
	}
	after, _ := json.Marshal(root)
	if string(before) != string(after) {
		t.Fatal("resolver mutated tree")
	}
	if len(root.Rows[1][0].Uses) != 2 {
		t.Fatal("alias provenance lost")
	}
	_, roots, err := Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ResolveInteractions(roots["Part"]); err == nil {
		t.Fatal("unresolved standalone template claimed valid")
	}
}
func TestInteractionOptInPreservesLegacy(t *testing.T) {
	root := interactionsRoot(t, `sdui 0.3; ref: app "unopened.sdl"; Main=[basic=button("B",callback=app.Do.@invoke);decorated=button("D",icon="symbol",tooltip="Tip");explicit=button("E",toggle=false,callback=app.Do.@invoke)];`)
	ids, err := ResolveInteractions(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range root.Rows[:2] {
		if IsCommandButton(row[0]) || row[0].Argument("$scope") != "" {
			t.Fatal("legacy promoted")
		}
		if _, ok := ids[row[0].Path]; ok {
			t.Fatal("legacy command identity")
		}
	}
	if !IsCommandButton(root.Rows[2][0]) || ids["Main/explicit"].Command != "Main/explicit" {
		t.Fatal("explicit opt-in missing")
	}
	_, roots, err := Compile(`sdui 0.2; Main=[b=button("B")];`)
	if err != nil {
		t.Fatal(err)
	}
	before, _ := json.Marshal(roots)
	ids, err = ResolveInteractions(roots["Main"])
	after, _ := json.Marshal(roots)
	if err != nil || len(ids) != 0 || string(before) != string(after) {
		t.Fatal("legacy resolution changed model")
	}
}
func TestInteractionSchemaRejections(t *testing.T) {
	for _, body := range []string{
		`command("Unnamed")`, `dialog("Unnamed")[]`, `menu("M")`, `separator(label="x")`, `item("x")`,
		`c=command("C",checked=false)`, `c=command("C",exclusive="g")`, `c=command("C",toggle=true,exclusive=" ")`,
		`c=command("C",context="other")`, `c=command("C",context="item")`, `c=command("C",target="x")`,
		`c=command("C",effect="unknown")`, `c=command("C",effect="accept",toggle=false)`, `c=command("C",effect="open")`,
		`button(command="c",callback=x.Do.@invoke)`, `button(command="c",key="F1")`, `button(command="c",toggle=false)`,
		`button()`, `button(" ")`, `c=command("C") {x=fill}`, `menu("C",mode="context",target="x")[] {padding=1}`,
		`menu("M",mode="bogus")[]`, `menu("M",mode="context")[]`, `menu("M",target="x")[]`,
		`menu("M")[button("Wrong")]`, `menu("M")[separator(),separator()]`, `menu("M")[menu("Wrong")[]]`,
		`menu("Wrong",mode="submenu")[]`, `menuGroup("Wrong")[]`, `item(command="c")`, `separator()`,
		`menu("M")[separator() {visible=false}]`, `s=split("horizontal")[c=command("C");r=[]]`,
		`s=split("horizontal")[d=dialog("D")[];r=[]]`,
	} {
		t.Run(body, func(t *testing.T) {
			if _, _, err := Compile(`sdui 0.3; ref: x "no.sdl"; Main=[` + body + `];`); err == nil {
				t.Fatal("accepted")
			}
		})
	}
	for _, key := range []string{"A", "Primary+Ctrl+A", "Shift+Alt+A", "Alt+Alt+A", "Meta+A", "F13", "Primary+a", "+F1"} {
		if _, _, err := Compile(`sdui 0.3; Main=[c=command("C",key="` + key + `")];`); err == nil {
			t.Fatal("accepted key", key)
		}
	}
	for _, body := range []string{`c=command("C")`, `m=menu("M")[]`, `d=dialog("D")[]`, `button("B",toggle=true)`} {
		if _, _, err := Compile(`sdui 0.2; Main=[` + body + `];`); err == nil {
			t.Fatal("0.2 accepted", body)
		}
	}
}
func TestInteractionSelectedRootRejections(t *testing.T) {
	for _, body := range []string{
		`button(command="missing") {visible=false}`, `x=input("X");button(command="x")`, `button(command="../x")`,
		`c=command("C",context="item",target="x");x=input("X")`,
		`c=command("C",effect="open",target="x");x=input("X")`,
		`button("Accept",effect="accept")`,
		`a=command("A",toggle=true,checked=true,exclusive="g");b=command("B",toggle=true,checked=true,exclusive="g")`,
		`a=command("A",key="F1");b=command("B",key="F1")`,
		`a=tree("A");b=list("B");c=command("C",context="item",target="a");menu("M",mode="context",target="b")[item(command="c")]`,
	} {
		root := interactionsRoot(t, `sdui 0.3; Main=[`+body+`];`)
		if _, err := ResolveInteractions(root); err == nil {
			t.Error("accepted", body)
		}
	}
}
func TestInteractionMenuStructureAndFields(t *testing.T) {
	root := interactionsRoot(t, `sdui 0.3; Main=[rows=tree("Rows");c=command("Action",context="item",target="rows");menu("Menu",mode="context",target="rows")[menuGroup("Group")[item(command="c");separator();menu("More",mode="submenu")[item(command="c")]]];d=dialog("D")[ok=button("OK",effect="accept")]];`)
	ids, err := ResolveInteractions(root)
	if err != nil {
		t.Fatal(err)
	}
	if ids["Main/c"].Target != "Main/rows" {
		t.Fatal(ids)
	}
	root.Walk(func(n *Instance) {
		if n.Widget == "command" || n.Widget == "dialog" || n.Widget == "menuGroup" || n.Widget == "item" || n.Widget == "separator" {
			if !IsAuxiliary(n) {
				t.Error("in flow", n.Path)
			}
		}
	})
	// An ambiguous normalized named path must never pick an arbitrary widget.
	clone := *root.Rows[0][0]
	clone.Path = "Main/$extra/rows"
	root.Rows = append(root.Rows, []*Instance{&clone})
	if _, err := ResolveInteractions(root); err == nil || !strings.Contains(err.Error(), "ambiguous") {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(root.Rows[0][0].Arguments, clone.Arguments) {
		t.Fatal("test setup")
	}
}

func TestInteractionSerializedArgumentsRejectWithoutPanic(t *testing.T) {
	for _, mutation := range []func(*Instance){
		func(n *Instance) { n.Arguments["toggle"] = Literal{Kind: "boolean", Value: "wrong"} },
		func(n *Instance) { n.Arguments["label"] = Literal{Kind: "string", Value: 42} },
		func(n *Instance) { n.Arguments["future"] = Literal{Kind: "string", Value: "unsupported"} },
		func(n *Instance) { n.Arguments["$scope"] = Literal{Kind: "string", Value: "other"} },
		func(n *Instance) { delete(n.Arguments, "$scope") },
		func(n *Instance) { n.Arguments["callback"] = Reference{Module: "x", Object: "Do", Member: "other"} },
	} {
		root := interactionsRoot(t, `sdui 0.3; Main=[c=command("C",toggle=true)];`)
		mutation(root.Rows[0][0])
		if _, err := ResolveInteractions(root); err == nil {
			t.Fatal("serialized invalid argument accepted")
		}
	}
}
func TestInteractionDefinitionRootsRetainNamesAndLimits(t *testing.T) {
	source := `sdui 0.3; ref: app "absent.sdl"; Action=command("A",toggle=true,callback=app.Do.@invoke);Form=dialog("F",callback=app.Save.@invoke)[];Main=[act=Action;form=Form;button(command="act")];`
	root := interactionsRoot(t, source)
	ids, err := ResolveInteractions(root)
	if err != nil {
		t.Fatal(err)
	}
	if ids["Main/act"].Scope != "Main/act" || ids["Main/form"].Scope != "Main/form" {
		t.Fatal(ids)
	}
	for _, source := range []string{`sdui 0.3; ref: app "x"; Main=[c=command("C",callback=app.Do.@other)];`, `sdui 0.3; Main=[c=command("C",icon="https://example.invalid/icon")];`} {
		if _, _, err := Compile(source); err == nil {
			t.Fatal("invalid M2 symbolic declaration accepted")
		}
	}
	root.Rows = append(root.Rows, []*Instance{root})
	if _, err := ResolveInteractions(root); err == nil {
		t.Fatal("cyclic tree accepted")
	}
}
