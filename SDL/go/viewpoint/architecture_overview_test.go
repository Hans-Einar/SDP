package viewpoint

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestArchitectureOverviewComplete(t *testing.T) {
	v, e := New(canonicalOverview(t, `language design-core version 0.5.
container Shell.
container Nested.
unit Worker.
unit Isolated.
functionality Work.
mode Online.
feature Feature.
Shell contains Nested.
Nested contains Worker.
Worker owns Work.
Work allocated-to Shell in mode Online.
Work contributes-to Feature.
`))
	if e != nil {
		t.Fatal(e)
	}
	s, e := v.Select(Query{Viewpoint: "VP02", Projection: "overview"})
	if e != nil {
		t.Fatal(e)
	}
	if len(s.Diagrams) != 1 {
		t.Fatal(s.Diagrams)
	}
	d := s.Diagrams[0]
	if d.ID != "VP02-overview" || len(d.Nodes) != 4 || len(d.Edges) != 2 || len(s.Facts) != 2 {
		t.Fatal(d, s.Facts)
	}
	for _, n := range []string{"Shell", "Nested", "Worker", "Isolated"} {
		if d.Nodes["n_"+n] == nil {
			t.Fatal("lost", n)
		}
	}
	for _, f := range s.Facts {
		if f.S("verb") != "contains" || !has(d.SourceFacts, f.S("id")) {
			t.Fatal(f)
		}
	}
	if strings.Contains(d.Mermaid(), "allocated-to") || strings.Contains(d.Mermaid(), "Feature") {
		t.Fatal(d.Mermaid())
	}
	if !strings.Contains(s.ArchitectureOverviewContext(), "| Isolated | unit | Root (no declared parent)") {
		t.Fatal(s.ArchitectureOverviewContext())
	}
	// The existing roots/parent views are independent, unchanged detail slices.
	for _, original := range v.Diagrams {
		if strings.HasPrefix(original.ID, "VP02-") {
			detail, e := v.Select(Query{Viewpoint: "VP02", Diagram: original.ID})
			if e != nil || len(detail.Diagrams) != 1 || !reflect.DeepEqual(original, detail.Diagrams[0]) {
				t.Fatal(original.ID, e)
			}
		}
	}
	for _, q := range []Query{
		{Viewpoint: "VP02", Projection: "overview", Focus: "Worker"},
		{Viewpoint: "VP02", Projection: "overview", Mode: "Online"},
		{Viewpoint: "VP02", Projection: "overview", Relations: []string{"contains"}},
		{Viewpoint: "VP02", Projection: "overview", Diagram: "VP02-Shell"},
		{Viewpoint: "VP02", Projection: "overview", Depth: 3},
	} {
		if _, e := v.Select(q); e == nil {
			t.Fatal("partial overview accepted", q)
		}
	}
}
func TestArchitectureOverviewEmptyAndLarge(t *testing.T) {
	v, e := New("language design-core version 0.5.\nactor Reader.\n")
	if e != nil {
		t.Fatal(e)
	}
	s, e := v.Select(Query{Viewpoint: "VP02", Projection: "overview"})
	if e != nil || len(s.Diagrams) != 0 || !strings.Contains(s.ArchitectureOverviewContext(), "No System, Container or Unit") {
		t.Fatal(e, s)
	}
	src := "language design-core version 0.5.\n"
	for i := 0; i < 170; i++ {
		src += fmt.Sprintf("unit Item%d.\n", i)
	}
	v, e = New(canonicalOverview(t, src))
	if e != nil {
		t.Fatal(e)
	}
	s, e = v.Select(Query{Viewpoint: "VP02", Projection: "overview"})
	if e != nil || len(s.Diagrams) != 1 || len(s.Diagrams[0].Nodes) != 170 {
		t.Fatal("truncated", e)
	}
}
