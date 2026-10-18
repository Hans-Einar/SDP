package parser

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestScalarSchemasLexemesAndProvenance(t *testing.T) {
	const source = `sdui 0.3;
ref: app "not-opened.sdl";
Field=number("Amount",min=-0,max=1e1,step=1.00e-1,value=0.30,callback=app.Save.@invoke);
Main=[left=Field;right=Field;flag=checkbox("Flag");choice=select("Choice",required=true);slide=slider("S",min=0,max=100,step=1,value=50)];
app.Save.setHandle(Main.left);`
	doc, roots, err := Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	root := roots["Main"]
	left, right := root.Rows[0][0], root.Rows[1][0]
	min, max, step, value, err := NumericArguments(left)
	if err != nil || min != "-0" || max != "1e1" || step != "1.00e-1" || value != "0.30" {
		t.Fatal(min, max, step, value, err)
	}
	for _, key := range []string{"min", "max", "step", "value"} {
		lit := left.Arguments[key].(Literal)
		if lit.Kind != "number-lexeme" || source[lit.Span.Start:lit.Span.End] != lit.Value {
			t.Fatal("lost exact numeric token", key, lit)
		}
		if !reflect.DeepEqual(lit, right.Arguments[key]) {
			t.Fatal("reuse changed original numeric span")
		}
	}
	if left.Path == right.Path || len(left.Uses) != 1 || left.Uses[0].Span == right.Uses[0].Span {
		t.Fatal("reuse identity")
	}
	if left.Argument("$scope") != "" || len(root.Rows[2][0].Arguments) != 1 || len(root.Rows[3][0].Arguments) != 2 {
		t.Fatal("derived/default arguments injected")
	}
	data, err := json.Marshal(Data(doc))
	if err != nil || !strings.Contains(string(data), `"kind":"number-lexeme"`) || !strings.Contains(string(data), `"value":"0.30"`) {
		t.Fatal(string(data), err)
	}
	before, _ := json.Marshal(root)
	if _, err := ResolveInteractions(root); err != nil {
		t.Fatal(err)
	}
	after, _ := json.Marshal(root)
	if string(before) != string(after) {
		t.Fatal("strict closure mutated source")
	}
	left.Arguments["value"] = Literal{Kind: "number-lexeme", Value: "0.4"}
	if right.Arguments["value"].(Literal).Value != "0.30" {
		t.Fatal("reuse aliases arguments")
	}
}
func TestScalarClosedSchemasAndTextBoundary(t *testing.T) {
	bad := []string{
		`checkbox()`, `checkbox(" ")`, `checkbox("C",value="true")`, `checkbox("C",value=null)`, `checkbox("C",readOnly=1)`,
		`checkbox("C",enabled=true)`, `checkbox("C",onChange=app.Do.@invoke)`, `checkbox("C",command="c")`, `checkbox("C",value=false,value=true)`,
		`checkbox(label="C","again")`, `checkbox("C")[]`, `select("S",options="x")`, `select("S",value=1)`, `select("S",required="yes")`,
		`slider("S",min=0,max=1,step=0.1)`, `slider("S",min="0",max=1,step=0.1,value=0)`,
		`number("N",min=0,max=1,step=0.1,value="0.3")`, `number("N",min=0,max=1,step=0.1,value=0,placeholder=false)`,
		`slider("S",min=0,max=1,step=0.1,value=0,placeholder="no")`,
		`textarea("I")`,
	}
	for _, body := range bad {
		t.Run(body, func(t *testing.T) {
			if _, _, err := Compile(`sdui 0.3; ref: app "unused"; Main=[field=` + body + `];`); err == nil {
				t.Fatal("accepted")
			}
		})
	}
	for _, body := range []string{`checkbox("C")`, `select("S")`, `number("N",min=0,max=1,step=1,value=0)`, `slider("S",min=0,max=1,step=1,value=0)`} {
		if _, _, err := Compile(`sdui 0.2; Main=[field=` + body + `];`); err == nil {
			t.Fatal("legacy accepted", body)
		}
	}
	if _, _, err := Compile(`sdui 0.3; ref: app "unused"; Main=[checkbox("C",callback=app.Do.@invoke)];`); err == nil {
		t.Fatal("anonymous callback accepted")
	}
}
func TestScalarExactNumericAdmissionAndDiagnosticSpan(t *testing.T) {
	for _, args := range []string{
		`min=0,max=1,step=0.1,value=0.3`,
		`min=-9007199254740991,max=9007199254740991,step=1,value=9007199254740991`,
		`min=1000000000000000,max=1000000000000010,step=1,value=1000000000000001`,
		`min=10000000000000000,max=10000000000000016,step=4,value=10000000000000004`,
	} {
		if _, _, err := Compile(`sdui 0.3; Main=[n=number("N",` + args + `)];`); err != nil {
			t.Fatal(args, err)
		}
	}
	for _, args := range []string{
		`min=0,max=1,step=0.1,value=0.30000000000000001`,
		`min=0,max=1,step=0.1,value=1.1`, `min=0,max=1,step=0,value=0`,
		`min=1,max=1,step=1,value=1`, `min=0,max=1,step=0.000000001,value=0`,
		`min=10000000000000000,max=10000000000000016,step=2,value=10000000000000004`,
		`min=-9007199254740991,max=9007199254740991,step=1,value=9007199254740991.1`,
	} {
		if _, _, err := Compile(`sdui 0.3; Main=[n=number("N",` + args + `)];`); err == nil {
			t.Fatal("accepted inexact/ambiguous grid", args)
		}
	}
	source := `sdui 0.3; Main=[n=number("N",min=0,max=1,step=0.1,value=0.30000000000000001)];`
	doc, err := Parse(source)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Normalize(doc)
	d, ok := err.(*Diagnostic)
	if !ok || source[d.Span.Start:d.Span.End] != "0.30000000000000001" {
		t.Fatal("error lost offending span", err)
	}
}
func TestScalarLiteralIsolationAndDialogOwnership(t *testing.T) {
	source := `sdui 0.3; Main=[s=split("horizontal",proportion=0.5)[a=[];b=[]];d=dialog("D")[n=number("N",min=0,max=1,step=0.1,value=0.3);c=checkbox("C");s=select("S");nested=dialog("Nested")[x=checkbox("X")];text=input("Text")];cmd=command("Context",context="widget",target="d/n")];`
	// Names are definition-wide; avoid colliding split/select names.
	source = strings.Replace(source, `;s=select`, `;choice=select`, 1)
	_, roots, err := Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	root := roots["Main"]
	if lit := root.Rows[0][0].Arguments["proportion"].(Literal); lit.Kind != "number" || lit.Value != float64(.5) {
		t.Fatal("split literal migrated", lit)
	}
	ids, err := ResolveInteractions(root)
	if err != nil {
		t.Fatal(err)
	}
	if ids["Main/d/n"].Dialog != "Main/d" || ids["Main/d/nested/x"].Dialog != "Main/d/nested" || ids["Main/cmd"].Target != "Main/d/n" {
		t.Fatal(ids)
	}
	if _, err := ResolveDialogField(root, "Main/d", "n"); err == nil {
		t.Fatal("typed field entered text-only mapping")
	}
	if _, err := ResolveDialogField(root, "Main/d", "text"); err != nil {
		t.Fatal(err)
	}
	for _, bad := range []any{Literal{Kind: "number", Value: .3}, Literal{Kind: "number-lexeme", Value: .3}, "0.3"} {
		node := root.Rows[1][0].Rows[0][0]
		saved := node.Arguments["value"]
		node.Arguments["value"] = bad
		if _, _, _, _, err := NumericArguments(node); err == nil {
			t.Fatal("rounded/malformed literal extracted")
		}
		if _, err := ResolveInteractions(root); err == nil {
			t.Fatal("malformed scalar passed selected closure")
		}
		node.Arguments["value"] = saved
	}
}

func TestScalarBoundedNumericLexemeAndConstructedAST(t *testing.T) {
	raw := "0." + strings.Repeat("0", 32768) + "1"
	source := `sdui 0.3; Main=[n=number("N",min=0,max=1,step=0.1,value=` + raw + `)];`
	doc, err := Parse(source)
	if err != nil {
		t.Fatal("finite source-number syntax changed", err)
	}
	if _, err := Normalize(doc); err == nil {
		t.Fatal("oversized mantissa admitted")
	} else if d, ok := err.(*Diagnostic); !ok || source[d.Span.Start:d.Span.End] != raw {
		t.Fatal("numeric limit span", err)
	}
	doc, err = Parse(`sdui 0.3; Main=[n=number("N",min=0,max=1,step=0.1,value=0.3)];`)
	if err != nil {
		t.Fatal(err)
	}
	node := doc.Definitions[0].Root.Rows[0].Items[0]
	for i, arg := range node.Arguments {
		if arg.Name != nil && *arg.Name == "value" {
			node.Arguments[i].Value = Literal{Kind: "number-lexeme", Value: float64(.3), Span: arg.Span}
		}
	}
	if _, err := Normalize(doc); err == nil {
		t.Fatal("constructed AST lost original token but passed")
	}
	if _, _, _, _, err := NumericArguments(nil); err == nil {
		t.Fatal("nil numeric source")
	}
}
