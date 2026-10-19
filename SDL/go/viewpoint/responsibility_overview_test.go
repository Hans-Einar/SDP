package viewpoint

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestResponsibilityOverviewSemantics(t *testing.T) {
	v, e := New(canonicalOverview(t, `language design-core version 0.5.
container Shell.
unit Owner.
unit Provider.
unit Isolated.
functionality Work.
functionality Unmapped.
capability Both.
capability ContributionOnly.
capability OfferedOnly.
capability Unconnected.
feature Feature.
mode Online.
Shell contains Owner.
Owner owns Work.
Provider owns Unmapped.
Work realizes Both.
Work realizes ContributionOnly.
Provider provides Both.
Provider provides OfferedOnly.
Work contributes-to Feature.
Work allocated-to Shell in mode Online.
`))
	if e != nil {
		t.Fatal(e)
	}
	s, e := v.Select(Query{Viewpoint: "VP03", Projection: "overview"})
	if e != nil {
		t.Fatal(e)
	}
	if len(s.Diagrams) != 1 {
		t.Fatal(s.Diagrams)
	}
	d := s.Diagrams[0]
	if len(d.Nodes) != 5 || len(d.Edges) != 1 || len(s.Facts) != 7 || len(d.Elements) != 6 {
		t.Fatal(d, s.Facts)
	}
	for _, name := range []string{"Shell", "Owner", "Provider", "Isolated", "Unconnected"} {
		if d.Nodes["n_"+name] == nil {
			t.Fatal(name)
		}
	}
	for _, f := range s.Facts {
		if !has([]string{"contains", "owns", "realizes", "provides"}, f.S("verb")) || !has(d.SourceFacts, f.S("id")) {
			t.Fatal(f)
		}
	}
	if !reflect.DeepEqual(d.SourceFacts, values(s.Facts, "id")) {
		t.Fatal(d.SourceFacts)
	}
	for _, a := range d.Elements {
		proof := a["source_facts"].([]string)
		switch a.S("role") {
		case "contributes":
			if a.S("node") != "n_Owner" || len(proof) != 2 {
				t.Fatal(a)
			}
			if eq(s.Facts, "id", proof[0])[0].S("verb") != "owns" || eq(s.Facts, "id", proof[1])[0].S("verb") != "realizes" {
				t.Fatal(a)
			}
		case "provides":
			if a.S("node") != "n_Provider" || len(proof) != 1 || eq(s.Facts, "id", proof[0])[0].S("verb") != "provides" {
				t.Fatal(a)
			}
		case "owns":
			if len(proof) != 1 || eq(s.Facts, "id", proof[0])[0].S("verb") != "owns" {
				t.Fatal(a)
			}
		default:
			t.Fatal(a)
		}
	}
	text := s.ResponsibilityOverviewContext()
	for _, want := range []string{"R01 | Unmapped | Provider | Not specified", "R02 | Work | Owner | Both, ContributionOnly", "C01 | Both | Provider | Work", "C02 | ContributionOnly | Not specified | Work", "C03 | OfferedOnly | Provider | Not specified", "C04 | Unconnected | Not specified | Not specified"} {
		if !strings.Contains(text, want) {
			t.Fatal(want, text)
		}
	}
	for _, absent := range []string{"allocated-to", "Feature (feature)"} {
		if strings.Contains(d.Mermaid(), absent) {
			t.Fatal(d.Mermaid())
		}
	}
	for _, original := range v.Diagrams {
		if strings.HasPrefix(original.ID, "VP03-") {
			detail, e := v.Select(Query{Viewpoint: "VP03", Diagram: original.ID})
			if e != nil || len(detail.Diagrams) != 1 || !reflect.DeepEqual(original, detail.Diagrams[0]) {
				t.Fatal(original.ID, e)
			}
		}
	}
	for _, q := range []Query{
		{Viewpoint: "VP03", Projection: "overview", Focus: "Owner"},
		{Viewpoint: "VP03", Projection: "overview", Mode: "Online"},
		{Viewpoint: "VP03", Projection: "overview", Relations: []string{"owns"}},
		{Viewpoint: "VP03", Projection: "overview", Diagram: "VP03-Both"},
		{Viewpoint: "VP03", Projection: "overview", Depth: 3},
		{Viewpoint: "VP03", Projection: "overview", Direction: "out"},
	} {
		if _, e := v.Select(q); e == nil {
			t.Fatal("partial overview accepted", q)
		}
	}
}
func TestResponsibilityOverviewEmptyAndLarge(t *testing.T) {
	v, e := New("language design-core version 0.5.\nactor Reader.\n")
	if e != nil {
		t.Fatal(e)
	}
	s, e := v.Select(Query{Viewpoint: "VP03", Projection: "overview"})
	if e != nil || len(s.Diagrams) != 0 {
		t.Fatal(e, s)
	}
	for _, want := range []string{"No Functionality responsibilities", "No Capabilities are declared"} {
		if !strings.Contains(s.ResponsibilityOverviewContext(), want) {
			t.Fatal(want)
		}
	}
	src := "language design-core version 0.5.\n"
	for i := 0; i < 170; i++ {
		src += fmt.Sprintf("unit U%d.\nfunctionality F%d.\ncapability C%d.\n", i, i, i)
	}
	for i := 0; i < 170; i++ {
		src += fmt.Sprintf("U%d owns F%d.\n", i, i)
	}
	v, e = New(canonicalOverview(t, src))
	if e != nil {
		t.Fatal(e)
	}
	s, e = v.Select(Query{Viewpoint: "VP03", Projection: "overview"})
	if e != nil || len(s.Diagrams) != 1 || len(s.Diagrams[0].Nodes) != 340 || len(s.Diagrams[0].Elements) != 170 || len(s.Facts) != 170 {
		t.Fatal("truncated", e)
	}
	if strings.Count(s.ResponsibilityOverviewContext(), "\n| R") != 170 || strings.Count(s.ResponsibilityOverviewContext(), "\n| C") != 170 {
		t.Fatal("register truncated")
	}
}
