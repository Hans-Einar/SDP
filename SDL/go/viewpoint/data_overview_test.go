package viewpoint

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const dataOverviewFixture = `language design-core version 0.5.
unit Owner.
database Store.
database Idle.
dataset State.
dataset Temporary.
datagram Notice.
contract Record.
contract Wire.
contract Unused.
variant Change.
field Value.
field Payload.
functionality Project.
feature Feature.
mode Online.
Owner owns Store.
Owner owns Idle.
Owner owns Project.
Owner holds Temporary.
Store holds State.
State upholds Record.
Temporary upholds Record.
Notice from State.
Notice upholds Wire.
Record has completeness = closed.
Record has-field Value.
Value has presence = required.
Value has value-type = text.
Wire has completeness = closed.
Wire defines Change.
Change has-field Payload.
Payload has presence = optional.
Payload has value-type = bytes.
Unused has completeness = open.
Project projects State into Notice.
Project contributes-to Feature.
`

func TestDataOverviewInventoryProofAndDetails(t *testing.T) {
	v, e := New(canonicalOverview(t, dataOverviewFixture))
	if e != nil {
		t.Fatal(e)
	}
	s, e := v.Select(Query{Viewpoint: "VP09", Projection: "overview"})
	if e != nil {
		t.Fatal(e)
	}
	if len(s.Diagrams) != 2 || len(s.MessageSets) != 0 {
		t.Fatal(s.Diagrams)
	}
	for i, d := range s.Diagrams {
		nodes, edges := 6, 5
		if i == 1 {
			nodes, edges = 3, 2
		}
		if len(d.Nodes) != nodes || len(d.Edges) != edges {
			t.Fatal(d)
		}
		for _, edge := range d.Edges {
			fs := eq(s.Facts, "id", edge.Fact)
			if len(fs) != 1 || edge.Source != "n_"+fs[0].S("subject") || edge.Target != "n_"+fs[0].S("object") {
				t.Fatal(edge, fs)
			}
		}
	}
	text := s.DataOverviewContext()
	for _, want := range []string{"State (Dataset): holder Store; contract Record; Database-backed persistence", "Temporary (Dataset): holder Owner; contract Record; Unit-held; persistence unspecified", "Notice (Datagram): logical source Dataset State; contract Wire", "| Idle | Owner | No Datasets declared |", "| Unused | open | 0 / 0 / 0 |", "| Change | Wire |", "| Payload | Change | bytes | optional |", "| Project | State | Notice | Owner |", "| Wire | defines | Change |", "| Record | has-field | Value |", "not SQL tables", "not imply null"} {
		if !strings.Contains(text, want) {
			t.Fatal(want, text)
		}
	}
	for _, f := range s.Facts {
		if f.S("verb") == "contributes-to" {
			t.Fatal("unrelated feature fact", f)
		}
	}
	for _, d := range v.Diagrams {
		if strings.HasPrefix(d.ID, "VP09-") {
			detail, e := v.Select(Query{Viewpoint: "VP09", Diagram: d.ID})
			if e != nil || len(detail.Diagrams) != 1 || !reflect.DeepEqual(d, detail.Diagrams[0]) {
				t.Fatal("detail changed", e)
			}
		}
	}
	again, e := v.Select(Query{Viewpoint: "VP09", Projection: "overview"})
	if e != nil || !reflect.DeepEqual(s, again) {
		t.Fatal("unstable", e)
	}
	for _, q := range []Query{{Viewpoint: "VP09", Projection: "overview", Mode: "Online"}, {Viewpoint: "VP09", Projection: "overview", Focus: "State"}, {Viewpoint: "VP09", Projection: "overview", Relations: []string{"holds"}}, {Viewpoint: "VP09", Projection: "overview", Diagram: "VP09-data-State"}, {Viewpoint: "VP09", Projection: "overview", Depth: 3}, {Viewpoint: "VP09", Projection: "overview", Direction: "out"}} {
		if _, e := v.Select(q); e == nil {
			t.Fatal("ignored partial filter", q)
		}
	}
}
func TestDataOverviewEmptyAndLarge(t *testing.T) {
	v, e := New(canonicalOverview(t, "language design-core version 0.5.\nunit Solo.\n"))
	if e != nil {
		t.Fatal(e)
	}
	s, e := v.Select(Query{Viewpoint: "VP09", Projection: "overview"})
	if e != nil || len(s.Diagrams) != 0 || len(s.Facts) != 0 {
		t.Fatal(e, s)
	}
	for _, want := range []string{"No dataset declarations", "No datagram declarations", "No Databases", "No Contracts", "No Variants", "No Fields", "No projections"} {
		if !strings.Contains(s.DataOverviewContext(), want) {
			t.Fatal(want)
		}
	}
	decls, stmts := "language design-core version 0.5.\n", ""
	for i := 0; i < 170; i++ {
		decls += fmt.Sprintf("contract Unused%d.\n", i)
		stmts += fmt.Sprintf("Unused%d has completeness = open.\n", i)
	}
	v, e = New(canonicalOverview(t, decls+stmts))
	if e != nil {
		t.Fatal(e)
	}
	s, e = v.Select(Query{Viewpoint: "VP09", Projection: "overview"})
	if e != nil || len(s.Diagrams) != 0 || strings.Count(s.DataOverviewContext(), "| open | 0 / 0 / 0 |") != 170 {
		t.Fatal("truncated inventory", e)
	}
}
