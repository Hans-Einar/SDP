package viewpoint

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestActivityOverviewComplete(t *testing.T) {
	v, e := New(canonicalOverview(t, `language design-core version 0.5.
activity Plan.
activity Build.
activity Check.
activity Isolated.
unit Owner.
functionality Work.
functionality Uncovered.
feature Feature.
scenario Example.
mode Online.
unit Handler.
channel Bus.
contract Events.
contract Body.
message Notice.
field Value.
Build refines Plan.
Check refines Plan.
Check depends-on Build.
Build addresses Work.
Build delivers Feature.
Example has completeness = open.
Example illustrates Build.
Example runs-in Online.
Example step 1 sends Notice from Owner to Handler via Bus.
Body has completeness = closed.
Body has-field Value.
Value has presence = required.
Value has value-type = text.
Events has completeness = closed.
Bus upholds Events.
Events permits Notice.
Notice upholds Body.
Notice has message-kind = event.
Owner uses Bus as sender of Notice in mode Online.
Handler uses Bus as receiver of Notice in mode Online.
Owner owns Work.
Owner owns Uncovered.
Work contributes-to Feature.
Plan has implementation-status = planned.
Build has implementation-status = implemented.
Check has implementation-status = verified.
`))
	if e != nil {
		t.Fatal(e)
	}
	before := v.Implementation()
	s, e := v.Select(Query{Viewpoint: "VP06", Projection: "overview"})
	if e != nil {
		t.Fatal(e)
	}
	if len(s.Diagrams) != 2 || len(s.Facts) != 10 {
		t.Fatal(s.Diagrams, s.Facts)
	}
	for i, d := range s.Diagrams {
		expected := 2
		verb := "refines"
		if i == 1 {
			expected = 1
			verb = "depends-on"
		}
		if len(d.Nodes) != 4 || len(d.Edges) != expected || d.Nodes["n_Isolated"] == nil {
			t.Fatal(d)
		}
		if !reflect.DeepEqual(d.SourceFacts, values(v.Relations[verb], "id")) {
			t.Fatal(d.SourceFacts)
		}
		for _, edge := range d.Edges {
			f := eq(v.Facts, "id", edge.Fact)[0]
			if edge.Source != "n_"+f.S("subject") || edge.Target != "n_"+f.S("object") || edge.Relation != verb {
				t.Fatal(edge)
			}
		}
	}
	if strings.Count(s.Diagrams[1].Mermaid(), "-.->|depends-on|") != 1 {
		t.Fatal(s.Diagrams[1])
	}
	text := s.ActivityOverviewContext()
	for _, want := range []string{"| Build | implemented | Plan |", "| Check | verified | Plan |", "| Plan | planned | No declared parent |", "| Isolated | Unspecified | No declared parent |", "| Build | Work | Owner |", "| Build | Feature |", "| Build | Example |", "- Uncovered", "not test evidence or owner acceptance"} {
		if !strings.Contains(text, want) {
			t.Fatal(want, text)
		}
	}
	for _, f := range s.Facts {
		if f.S("verb") == "contributes-to" || f.S("verb") == "owns" && f.S("object") == "Uncovered" {
			t.Fatal("unrelated fact", f)
		}
	}
	if before != v.Implementation() {
		t.Fatal("implementation report changed")
	}
	for _, d := range v.Diagrams {
		if strings.HasPrefix(d.ID, "VP06-") {
			detail, e := v.Select(Query{Viewpoint: "VP06", Diagram: d.ID})
			if e != nil || len(detail.Diagrams) != 1 || !reflect.DeepEqual(d, detail.Diagrams[0]) {
				t.Fatal("detail changed", d.ID, e)
			}
		}
	}
	again, e := v.Select(Query{Viewpoint: "VP06", Projection: "overview"})
	if e != nil || !reflect.DeepEqual(s, again) {
		t.Fatal("unstable", e)
	}
	for _, q := range []Query{
		{Viewpoint: "VP06", Projection: "overview", Mode: "Online"},
		{Viewpoint: "VP06", Projection: "overview", Focus: "Build"},
		{Viewpoint: "VP06", Projection: "overview", Relations: []string{"addresses"}},
		{Viewpoint: "VP06", Projection: "overview", Diagram: "VP06-roots"},
		{Viewpoint: "VP06", Projection: "overview", Depth: 3},
		{Viewpoint: "VP06", Projection: "overview", Direction: "out"},
	} {
		if _, e := v.Select(q); e == nil {
			t.Fatal("partial selection", q)
		}
	}
}
func TestActivityOverviewEmptyAndLarge(t *testing.T) {
	v, e := New(canonicalOverview(t, "language design-core version 0.5.\nunit Owner.\nfunctionality Work.\nOwner owns Work.\n"))
	if e != nil {
		t.Fatal(e)
	}
	s, e := v.Select(Query{Viewpoint: "VP06", Projection: "overview"})
	if e != nil || len(s.Diagrams) != 0 || len(s.Facts) != 0 {
		t.Fatal(e, s)
	}
	for _, want := range []string{"No Activities are declared", "- Work", "not unimplemented native code", "No activity prerequisites"} {
		if !strings.Contains(s.ActivityOverviewContext(), want) {
			t.Fatal(want)
		}
	}
	src := "language design-core version 0.5.\n"
	for i := 0; i < 170; i++ {
		src += fmt.Sprintf("activity A%d.\n", i)
	}
	v, e = New(canonicalOverview(t, src))
	if e != nil {
		t.Fatal(e)
	}
	s, e = v.Select(Query{Viewpoint: "VP06", Projection: "overview"})
	if e != nil || len(s.Diagrams) != 1 || len(s.Diagrams[0].Nodes) != 170 || len(s.Diagrams[0].Edges) != 0 {
		t.Fatal("truncated", e)
	}
	if strings.Count(s.ActivityOverviewContext(), "| Unspecified |") != 170 {
		t.Fatal("register truncated")
	}
}
