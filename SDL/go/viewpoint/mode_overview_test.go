package viewpoint

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestModeOverviewIsolation(t *testing.T) {
	v, e := New(canonicalOverview(t, `language design-core version 0.5.
capability Main.
capability Other.
capability Unspecified.
interface Local.
interface Remote.
interface Shared.
interface Unused.
mode Offline.
mode Online.
mode Quiet.
unit Consumer.
Consumer consumes Unused.
Consumer provides Main.
Main requires Local in mode Offline.
Main requires Shared in mode Offline.
Main requires Remote in mode Online.
Other requires Shared in mode Online.
`))
	if e != nil {
		t.Fatal(e)
	}
	index, e := v.Select(Query{Viewpoint: "VP05", Projection: "overview"})
	if e != nil || len(index.Diagrams) != 0 || len(index.Facts) != 0 || len(index.Relations) != 0 {
		t.Fatal("union index", e, index)
	}
	for _, want := range []string{"| Offline | 2 |", "| Online | 2 |", "| Quiet | 0 |"} {
		if !strings.Contains(v.ModeOverviewContext(""), want) {
			t.Fatal(want)
		}
	}
	for _, mode := range []string{"Offline", "Online", "Quiet"} {
		s, e := v.Select(Query{Viewpoint: "VP05", Projection: "overview", Mode: mode})
		if e != nil || len(s.Diagrams) != 1 {
			t.Fatal(mode, e)
		}
		d := s.Diagrams[0]
		count := 2
		if mode == "Quiet" {
			count = 0
		}
		if len(d.Edges) != count || len(s.Facts) != count || len(d.Nodes) != 3+count {
			t.Fatal(mode, d)
		}
		for _, f := range s.Facts {
			if f.S("mode") != mode || f.S("verb") != "requires" || !has(d.SourceFacts, f.S("id")) {
				t.Fatal("leak", f)
			}
		}
		if d.Nodes["n_Unspecified"] == nil || d.Nodes["n_Unused"] != nil || d.Nodes["n_Consumer"] != nil {
			t.Fatal(d.Nodes)
		}
		absent := "Remote"
		if mode == "Online" {
			absent = "Local"
		}
		if d.Nodes["n_"+absent] != nil || strings.Contains(v.ModeOverviewContext(mode), "| "+absent+" |") {
			t.Fatal("mode leak", mode)
		}
		if strings.Count(d.Mermaid(), "-.->|requires|") != count {
			t.Fatal(d.Mermaid())
		}
		if !strings.Contains(v.ModeOverviewContext(mode), "| Unspecified | No declared requirement | — |") {
			t.Fatal(mode)
		}
		again, _ := v.Select(Query{Viewpoint: "VP05", Projection: "overview", Mode: mode})
		if !reflect.DeepEqual(s, again) {
			t.Fatal("unstable")
		}
	}
	for _, d := range v.Diagrams {
		if strings.HasPrefix(d.ID, "VP05-") {
			s, e := v.Select(Query{Viewpoint: "VP05", Diagram: d.ID})
			if e != nil || len(s.Diagrams) != 1 || !reflect.DeepEqual(d, s.Diagrams[0]) {
				t.Fatal("detail changed", e)
			}
		}
	}
	for _, q := range []Query{
		{Viewpoint: "VP05", Projection: "overview", Mode: "Missing"},
		{Viewpoint: "VP05", Projection: "overview", Mode: "Online", Focus: "Main"},
		{Viewpoint: "VP05", Projection: "overview", Mode: "Online", Relations: []string{"requires"}},
		{Viewpoint: "VP05", Projection: "overview", Diagram: "VP05-Online"},
		{Viewpoint: "VP05", Projection: "overview", Depth: 3},
		{Viewpoint: "VP05", Projection: "overview", Direction: "out"},
		{Viewpoint: "VP05", Projection: "overview", Level: "A2"},
	} {
		if _, e := v.Select(q); e == nil {
			t.Fatal("partial/unknown selection", q)
		}
	}
}
func TestModeOverviewEmptyAndLarge(t *testing.T) {
	v, e := New("language design-core version 0.5.\nunit Empty.\n")
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(v.ModeOverviewContext(""), "No Modes are declared") {
		t.Fatal(v.ModeOverviewContext(""))
	}
	v, e = New("language design-core version 0.5.\nmode NormalDesktop.\n")
	if e != nil {
		t.Fatal(e)
	}
	s, e := v.Select(Query{Viewpoint: "VP05", Projection: "overview", Mode: "NormalDesktop"})
	if e != nil || len(s.Diagrams) != 0 || len(s.Facts) != 0 || !strings.Contains(v.ModeOverviewContext("NormalDesktop"), "No Capabilities are declared") {
		t.Fatal(e, s)
	}
	src := "language design-core version 0.5.\nmode Online.\ninterface Port.\n"
	for i := 0; i < 170; i++ {
		src += fmt.Sprintf("capability C%d.\n", i)
	}
	for i := 0; i < 170; i++ {
		src += fmt.Sprintf("C%d requires Port in mode Online.\n", i)
	}
	v, e = New(canonicalOverview(t, src))
	if e != nil {
		t.Fatal(e)
	}
	s, e = v.Select(Query{Viewpoint: "VP05", Projection: "overview", Mode: "Online"})
	if e != nil || len(s.Diagrams) != 1 || len(s.Diagrams[0].Nodes) != 171 || len(s.Diagrams[0].Edges) != 170 {
		t.Fatal("truncated", e)
	}
}
