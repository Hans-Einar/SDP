package viewpoint

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestInterfaceOverviewComplete(t *testing.T) {
	v, e := New(canonicalOverview(t, `language design-core version 0.5.
container Shell.
unit Client.
unit Unrelated.
interface Alpha.
interface Beta.
interface Isolated.
capability Ability.
mode Online.
Shell contains Client.
Shell consumes Alpha.
Client consumes Alpha.
Client consumes Beta.
Client provides Ability.
Ability requires Isolated in mode Online.
`))
	if e != nil {
		t.Fatal(e)
	}
	s, e := v.Select(Query{Viewpoint: "VP04", Projection: "overview"})
	if e != nil {
		t.Fatal(e)
	}
	if len(s.Diagrams) != 1 {
		t.Fatal(s.Diagrams)
	}
	d := s.Diagrams[0]
	if d.ID != "VP04-overview" || len(d.Nodes) != 5 || len(d.Edges) != 3 || len(s.Facts) != 3 {
		t.Fatal(d, s.Facts)
	}
	for _, name := range []string{"Shell", "Client", "Alpha", "Beta", "Isolated"} {
		if d.Nodes["n_"+name] == nil {
			t.Fatal("lost", name)
		}
	}
	if d.Nodes["n_Unrelated"] != nil || d.Nodes["n_Ability"] != nil || d.Nodes["n_Online"] != nil {
		t.Fatal("inferred unrelated node", d.Nodes)
	}
	if !reflect.DeepEqual(d.SourceFacts, values(v.Relations["consumes"], "id")) {
		t.Fatal(d.SourceFacts)
	}
	for _, e := range d.Edges {
		f := eq(v.Facts, "id", e.Fact)[0]
		if f.S("verb") != "consumes" || e.Source != "n_"+f.S("subject") || e.Target != "n_"+f.S("object") {
			t.Fatal(e)
		}
	}
	syntax := d.Mermaid()
	if strings.Count(syntax, "-.->|consumes|") != 3 || !strings.Contains(syntax, `n_Isolated("Isolated (interface)")`) {
		t.Fatal(syntax)
	}
	text := s.InterfaceOverviewContext()
	for _, want := range []string{"| Alpha | Client, Shell |", "| Beta | Client |", "| Isolated | No declared consumer | — |", "No provider, Channel, transport"} {
		if !strings.Contains(text, want) {
			t.Fatal(want, text)
		}
	}
	// Existing ordinary/table selection keeps its behavior. The overview retains
	// that exact consumes table; it does not acquire other relation semantics.
	ordinary, e := v.Select(Query{Viewpoint: "VP04"})
	if e != nil {
		t.Fatal(e)
	}
	if len(ordinary.Diagrams) != 0 || ordinary.Tables("VP04") != s.Tables("VP04") {
		t.Fatal("legacy table changed")
	}
	again, e := v.Select(Query{Viewpoint: "VP04", Projection: "overview"})
	if e != nil || !reflect.DeepEqual(s, again) {
		t.Fatal("nondeterministic", e)
	}
	for _, q := range []Query{
		{Viewpoint: "VP04", Projection: "overview", Focus: "Client"},
		{Viewpoint: "VP04", Projection: "overview", Mode: "Online"},
		{Viewpoint: "VP04", Projection: "overview", Relations: []string{"consumes"}},
		{Viewpoint: "VP04", Projection: "overview", Diagram: "VP04-overview"},
		{Viewpoint: "VP04", Projection: "overview", Depth: 3},
		{Viewpoint: "VP04", Projection: "overview", Direction: "out"},
		{Viewpoint: "VP04", Projection: "overview", Level: "A2"},
	} {
		if _, e := v.Select(q); e == nil {
			t.Fatal("partial overview accepted", q)
		}
	}
}

func TestInterfaceOverviewEmptyAndLarge(t *testing.T) {
	v, e := New("language design-core version 0.5.\nunit Client.\n")
	if e != nil {
		t.Fatal(e)
	}
	s, e := v.Select(Query{Viewpoint: "VP04", Projection: "overview"})
	if e != nil || len(s.Diagrams) != 0 || !strings.Contains(s.InterfaceOverviewContext(), "No Interfaces are declared") {
		t.Fatal(e, s)
	}
	v, e = New("language design-core version 0.5.\ninterface Isolated.\n")
	if e != nil {
		t.Fatal(e)
	}
	s, e = v.Select(Query{Viewpoint: "VP04", Projection: "overview"})
	if e != nil || len(s.Diagrams) != 1 || len(s.Diagrams[0].Nodes) != 1 || len(s.Diagrams[0].Edges) != 0 || !strings.Contains(s.InterfaceOverviewContext(), "No interface consumption is declared") {
		t.Fatal(e, s)
	}
	src := "language design-core version 0.5.\nunit Client.\n"
	for i := 0; i < 170; i++ {
		src += fmt.Sprintf("interface I%d.\n", i)
	}
	for i := 0; i < 169; i++ {
		src += fmt.Sprintf("Client consumes I%d.\n", i)
	}
	v, e = New(canonicalOverview(t, src))
	if e != nil {
		t.Fatal(e)
	}
	s, e = v.Select(Query{Viewpoint: "VP04", Projection: "overview"})
	if e != nil || len(s.Diagrams) != 1 || len(s.Diagrams[0].Nodes) != 171 || len(s.Diagrams[0].Edges) != 169 || len(s.Facts) != 169 {
		t.Fatal("truncated", e)
	}
	if strings.Count(s.InterfaceOverviewContext(), "\n| I") != 171 {
		t.Fatal("register truncated")
	}
}
