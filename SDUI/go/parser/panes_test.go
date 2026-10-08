package parser

import (
	"encoding/json"
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestPaneBodiesReuseAndStableIdentity(t *testing.T) {
	const source = `sdui 0.3;
ref: actions "never-opened.sdl";
reusable=page("Résumé",icon="document")[header=input("Draft",value="λ")];
set=tabs("Workspace",selected="second")[first=reusable;second=page("Résumé")[]];
alias=set;
Main=[
 pair=split(axis="horizontal")[left=alias;right=alias];
 live=tabs("Live",callback=actions.Page.@invoke)[overview=page("Overview")[]];
 preview=input("Preview",value="")
];
actions.Page.setHandle(Main.preview);`
	doc, roots, err := Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	split := roots["Main"].Rows[0][0]
	if split.Kind != "composition" || split.Widget != "split" || split.Arguments["proportion"].(Literal).Value != .5 || split.Arguments["collapsible"].(Literal).Value != true {
		t.Fatal(split)
	}
	children, err := PaneChildren(split)
	if err != nil || children[0].ID != "left" || children[1].ID != "right" {
		t.Fatal(children, err)
	}
	left, right := children[0].Node, children[1].Node
	pages, err := PaneChildren(left)
	if err != nil || pages[0].ID != "first" || pages[1].ID != "second" || pages[0].Node.Argument("label") != pages[1].Node.Argument("label") {
		t.Fatal(pages, err)
	}
	if len(left.Uses) != 2 || left.Path == right.Path || left.Uses[1].Span == right.Uses[1].Span {
		t.Fatal("reuse identity/provenance lost")
	}
	page := pages[0].Node
	if len(page.Regions) != 0 || page.Rows[0][0].Path != "Main/pair/left/first/header" || source[page.Span.Start:page.Span.End] != `page("Résumé",icon="document")[header=input("Draft",value="λ")]` {
		t.Fatal("body name became region or UTF-8 span lost", page)
	}
	if rootProfile, err := EffectiveProfile(roots["Main"]); err != nil || rootProfile != "sdui/0.3" {
		t.Fatal(rootProfile, err)
	}
	encoded, err := json.Marshal(Data(doc))
	if err != nil || !strings.Contains(string(encoded), `"kind":"composition"`) {
		t.Fatal(string(encoded), err)
	}
	before, _ := json.Marshal(Data(doc))
	page.Arguments["label"] = Literal{Kind: "string", Value: "changed"}
	after, _ := json.Marshal(Data(doc))
	if string(before) != string(after) || right.Rows[0][0].Argument("label") != "Résumé" {
		t.Fatal("normalized trees share AST/argument state")
	}
}

func TestPaneSchemaAndPlacementRejection(t *testing.T) {
	for _, body := range []string{
		`tabs("Tabs")`, `tabs("Tabs")[]`, `tabs("")[p=page("P")[]]`,
		`tabs("T")[page("P")[]]`, `tabs("T")[p=page("P")[],q=page("Q")[]]`,
		`tabs("T")[p=input("P")]`, `tabs("T",selected="none")[p=page("P")[]]`,
		`tabs("T",selected="p")[p=page("P")[] {visible=false}]`,
		`tabs("T",selected="p")[p=page("P")[] {enabled=false}]`,
		`tabs("T",selected="")[p=page("P")[]]`,
		`tabs("T",selected=3)[p=page("P")[]]`,
		`tabs("T",label="T")[p=page("P")[]]`,
		`tabs(label="T","extra")[p=page("P")[]]`,
		`tabs("T",command="future")[p=page("P")[]]`,
		`tabs("T")[p=page("P",icon=false)[]]`, `tabs("T")[p=page("P",icon=" ")[]]`,
		`page("P")[]`, `split("horizontal")[left=page("P")[];right=[]]`,
		`split("horizontal")[]`, `split("horizontal")[left=[]]`,
		`split("horizontal")[left=[];right=[];extra=[]]`,
		`split("horizontal")[left=[],right=[]]`, `split("horizontal")[[];right=[]]`,
		`split("diagonal")[left=[];right=[]]`, `split()[left=[];right=[]]`,
		`split("horizontal",proportion="0.5")[left=[];right=[]]`,
		`split("horizontal",proportion=1.1)[left=[];right=[]]`,
		`split("horizontal",minFirst=-0.1)[left=[];right=[]]`,
		`split("horizontal",minFirst=0.5,minSecond=0.5)[left=[];right=[]]`,
		`split("horizontal",minFirst=0.6)[left=[];right=[]]`,
		`split("horizontal",collapsible=1)[left=[];right=[]]`,
		`button("OK")[]`, `tree("Navigation")[]`, `menu("File")`,
		`dialog("Dialog")[]`, `command("Future")`,
	} {
		if _, _, err := Compile(`sdui 0.3; Main=[` + body + `];`); err == nil {
			t.Errorf("accepted %s", body)
		}
	}
	for _, source := range []string{
		`sdui 0.2; Main=[tabs("T")[p=page("P")[]]];`,
		`sdui 0.2; Main=[tabs("T")];`,
		`sdui 0.3; reusable=page("P")[]; Main=[bad=reusable];`,
		`sdui 0.3; reusable=input("P"); Main=[t=tabs("T")[bad=reusable]];`,
		`sdui 0.3; Main=[tabs("T")[p=page("P")[];p=page("P")[]]];`,
		`sdui 0.3; ref: app "x.sdl"; Main=[tabs("T",callback=app.Page.@invoke)[p=page("P")[]]];`,
		`sdui 0.3; Main=[t=tabs("T",callback=missing.Page.@invoke)[p=page("P")[]]];`,
	} {
		if _, _, err := Compile(source); err == nil {
			t.Errorf("accepted %s", source)
		}
	}
}

func TestPaneOptionalSelectionAndAliases(t *testing.T) {
	for _, source := range []string{
		`sdui 0.3; P=page("P")[]; Main=[t=tabs("T")[P]];`,
		`sdui 0.3; Main=[t=tabs("T")[p=page("P")[] {visible=false}]];`,
		`sdui 0.3; Main=[t=tabs("T",selected="p")[p=page("P")[]] {visible=false}];`,
		`sdui 0.3; Main=[s=split("vertical",proportion=0,minSecond=0.2,collapsible=false)[left=[];right=[]]];`,
	} {
		if _, _, err := Compile(source); err != nil {
			t.Fatal(source, err)
		}
	}
}

func TestConstructedPaneNumericValidationAndLimits(t *testing.T) {
	for _, invalid := range []any{math.Inf(1), math.NaN(), "not a float"} {
		doc, err := Parse(`sdui 0.3; Main=[s=split("horizontal",proportion=0.5)[left=[];right=[]]];`)
		if err != nil {
			t.Fatal(err)
		}
		n := doc.Definitions[0].Root.Rows[0].Items[0]
		n.Arguments[1].Value = Literal{Kind: "number", Value: invalid, Span: n.Span}
		if err := Validate(doc); err == nil {
			t.Fatal("constructed nonfinite/wrong payload accepted", invalid)
		}
	}
	deep := `"leaf"`
	for i := 0; i < 65; i++ {
		deep = `page("P")[` + deep + `]`
	}
	if _, err := Parse(`sdui 0.3; Main=` + deep + `;`); err == nil {
		t.Fatal("composition body bypassed source depth limit")
	}
	if children, err := PaneChildren(nil); err == nil || children != nil {
		t.Fatal("nil pane accepted")
	}
	_, roots, err := Compile(`sdui 0.3; Main=[s=split("horizontal")[left=[];right=[]]];`)
	if err != nil {
		t.Fatal(err)
	}
	n := roots["Main"].Rows[0][0]
	for _, path := range []string{"Main/s/$anonymous", "Main/s/left/nested", "Main/s/..", "different/left"} {
		n.Rows[0][0].Path = path
		if _, err := PaneChildren(n); err == nil {
			t.Fatal("invalid direct identity accepted", path)
		}
	}
	if !reflect.DeepEqual(n.Arguments["minFirst"].(Literal).Value, float64(0)) {
		t.Fatal("split numeric defaults changed")
	}
}
