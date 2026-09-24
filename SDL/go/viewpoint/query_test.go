package viewpoint

import (
	"os"
	"reflect"
	"testing"
)

func full(t *testing.T) *Views {
	b, e := os.ReadFile("../../../SDUI/design/architecture.design")
	if e != nil {
		t.Fatal(e)
	}
	v, e := New(string(b))
	if e != nil {
		t.Fatal(e)
	}
	return v
}
func TestSelectionParityAndLimits(t *testing.T) {
	v := full(t)
	for _, d := range v.Diagrams {
		q := Query{Viewpoint: d.ID[:4], Diagram: d.ID}
		s, e := v.Select(q)
		if e != nil {
			t.Fatal(d.ID, e)
		}
		if len(s.Diagrams) != 1 || !reflect.DeepEqual(s.Diagrams[0], d) {
			t.Fatal("full/subset mismatch", d.ID)
		}
	}
	for _, q := range []Query{{Viewpoint: "bogus"}, {Viewpoint: "VP02", Focus: "Missing"}, {Viewpoint: "VP02", Relations: []string{"composition"}}, {Viewpoint: "VP02", Level: "A0"}, {Viewpoint: "VP02", Depth: 9}, {Viewpoint: "VP02", Direction: "sideways"}, {Viewpoint: "VP02", Mode: "Missing"}, {Viewpoint: "VP02", Diagram: "../no"}} {
		if _, e := v.Select(q); e == nil {
			t.Fatalf("accepted %#v", q)
		}
	}
	for _, direction := range []string{"in", "out", "both"} {
		s, e := v.Select(Query{Viewpoint: "VP02", Focus: "SduiFrontend", Depth: 1, Direction: direction, Relations: []string{"contains"}})
		if e != nil {
			t.Fatal(e)
		}
		for _, d := range s.Diagrams {
			for _, edge := range d.Edges {
				if edge.Relation != "contains" {
					t.Fatal(edge)
				}
				if direction == "in" && edge.Target != "n_SduiFrontend" {
					t.Fatal(edge)
				}
				if direction == "out" && edge.Source != "n_SduiFrontend" {
					t.Fatal(edge)
				}
			}
		}
	}
}
func TestURIRejectsUnregisteredCommands(t *testing.T) {
	for _, uri := range []string{"https://x/VP02", "sdl-view://x/VP02?command=rm", "sdl-view://x/VP02?depth=1&depth=2", "sdl-view://x/VP02?target=other", "sdl-view://x/VP02?consumer=shell", "sdl-view://x/VP02?depth=%zz"} {
		if _, e := ParseURI(uri); e == nil {
			t.Fatal(uri)
		}
	}
	s, e := ParseURI("sdl-view://sdui-design/VP02?focus=Host&relations=contains&depth=1&direction=out&level=A2")
	if e != nil || s.Query.Depth != 1 || s.Query.Focus != "Host" {
		t.Fatal(s, e)
	}
}
