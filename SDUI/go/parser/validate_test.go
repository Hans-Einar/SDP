package parser

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func instanceSummary(roots map[string]*Instance) map[string]any {
	out := map[string]any{}
	for _, root := range roots {
		root.Walk(func(n *Instance) {
			args := map[string]any{}
			for k, v := range n.Arguments {
				args[k] = Data(v)
			}
			rows := [][]string{}
			for _, row := range n.Rows {
				names := []string{}
				for _, child := range row {
					names = append(names, child.Path)
				}
				rows = append(rows, names)
			}
			regions := map[string]string{}
			for _, reg := range n.Regions {
				regions[reg.Role] = reg.Node.Path
			}
			out[n.Path] = map[string]any{"kind": n.Kind, "path": n.Path, "widget": n.Widget, "variant": n.Variant, "text": n.Text, "layout": n.Layout, "arguments": args, "rows": rows, "regions": regions}
		})
	}
	return out
}
func TestPythonSemanticPort(t *testing.T) {
	raw, e := os.ReadFile("testdata/python-port-cases.json")
	if e != nil {
		t.Fatal(e)
	}
	var cases []struct {
		Source, Code string
		Instances    any
	}
	if e = json.Unmarshal(raw, &cases); e != nil {
		t.Fatal(e)
	}
	for index, c := range cases {
		_, roots, err := Compile(c.Source)
		if c.Code != "" {
			d, ok := err.(*Diagnostic)
			if !ok || d.Code != c.Code {
				t.Errorf("case %d expected %s got %v", index, c.Code, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("case %d: %v", index, err)
			continue
		}
		encoded, _ := json.Marshal(instanceSummary(roots))
		var actual any
		json.Unmarshal(encoded, &actual)
		if !reflect.DeepEqual(actual, c.Instances) {
			t.Errorf("case %d normalized tree differs", index)
		}
	}
}
func TestExpansionAndIndependence(t *testing.T) {
	src := `sdui 0.2; P=[left=G {font=12},right=G]; G=<b=button("OK")>;`
	doc, roots, e := Compile(src)
	if e != nil {
		t.Fatal(e)
	}
	left, right := roots["P"].Rows[0][0], roots["P"].Rows[0][1]
	left.Layout["font"] = 99.0
	if _, ok := right.Layout["font"]; ok {
		t.Fatal("Shared instance layout")
	}
	again, e := Normalize(doc)
	if e != nil {
		t.Fatal(e)
	}
	if again["P"].Rows[0][0].Layout["font"] != 12.0 {
		t.Fatal("AST was mutated")
	}
	many := "sdui 0.2; P=[<" + strings.Repeat(`button("x"),`, 2048) + `"x">];`
	_, _, e = Compile(many)
	if d, ok := e.(*Diagnostic); !ok || d.Code != "node-limit" {
		t.Fatal(e)
	}
}

func TestRatioDoesNotAliasAST(t *testing.T) {
	d, roots, err := Compile(`sdui 0.2; P=[] {16:9,<->};`)
	if err != nil {
		t.Fatal(err)
	}
	roots["P"].Layout["ratio"].([]float64)[0] = 1
	again, err := Normalize(d)
	if err != nil || again["P"].Layout["ratio"].([]float64)[0] != 16 {
		t.Fatal("Layout aliases AST", err)
	}
}
