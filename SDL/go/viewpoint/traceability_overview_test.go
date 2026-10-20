package viewpoint

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func traceabilityFixture() string {
	declarations := `actor Guest.
usecase Orphan.
activity Build.
functionality OtherWork.
container Host.
mode Offline.
unit Target.
channel Bus.
channel Idle.
contract Protocol.
contract MessageBody.
field Content.
message Notice.
scenario Example.
`
	statements := `Work has state-retention = stateful.
Work has repeatability = deterministic.
Work projects State into Family.
Work allocated-to Host in mode Online.
Unrelated owns OtherWork.
OtherWork contributes-to Feature.
Build addresses Work.
Build has implementation-status = verified.
Bus upholds Protocol.
Idle upholds Protocol.
Protocol has completeness = closed.
Protocol permits Notice.
MessageBody has completeness = closed.
MessageBody has-field Content.
Content has presence = required.
Content has value-type = text.
Notice has message-kind = event.
Notice upholds MessageBody.
Unrelated uses Bus as sender of Notice in mode Online.
Target uses Bus as receiver of Notice in mode Online.
Unrelated uses Bus as sender of Notice in mode Offline.
Example has completeness = open.
Example runs-in Online.
Example step 1 sends Notice from Unrelated to Target via Bus.
`
	return strings.Replace(encodingOverviewFixture, "State upholds", declarations+"State upholds", 1) + statements
}

func TestTraceabilityOverviewCompleteClaimsAndGapContext(t *testing.T) {
	v, e := New(canonicalOverview(t, traceabilityFixture()))
	if e != nil {
		t.Fatal(e)
	}
	before := v.Tables("VP11")
	s, e := v.Select(Query{Viewpoint: "VP11", Projection: "overview"})
	if e != nil {
		t.Fatal(e)
	}
	if len(s.Diagrams) != 0 || len(s.MessageSets) != 0 || !reflect.DeepEqual(s.Facts, v.Facts) || !reflect.DeepEqual(s.Gaps, v.Gaps) || !reflect.DeepEqual(s.Model.Declarations, v.Model.Declarations) {
		t.Fatal("incomplete inventory")
	}
	text := s.TraceabilityOverviewContext()
	for _, f := range v.Facts {
		if !strings.Contains(text, f.S("text")) {
			t.Fatal("missing statement", f)
		}
	}
	for _, d := range v.Model.Declarations {
		if !strings.Contains(text, fmt.Sprintf("- %s (%s)", d.Name.Name, d.Kind)) {
			t.Fatal("missing declaration", d)
		}
	}
	for _, g := range v.Gaps {
		if !strings.Contains(text, g.S("code")) || !strings.Contains(text, g.S("message")) {
			t.Fatal("missing gap", g)
		}
	}
	for _, want := range []string{fmt.Sprintf("| Declarations | %d |", len(v.Model.Declarations)), fmt.Sprintf("| Facts | %d |", len(v.Facts)), fmt.Sprintf("| Explicit properties | %d |", len(eq(v.Facts, "verb", "has"))), fmt.Sprintf("| Reported model gaps | %d |", len(v.Gaps)), "Work: state-retention = stateful", "Work: repeatability = deterministic", "Build: implementation-status = verified", "Notice: message-kind = event", "UNSPECIFIED_ALLOCATION — OtherWork; mode Online; Feature Feature", "INCOMPLETE_PARTICIPATION — Bus; mode Offline", "NO_ENCODING — Unencoded", "NO_ACTOR — Orphan", "NO_GOAL — Guest", "not proof of implemented behavior", "Missing properties are unspecified", "not an exhaustive requirements audit", "absence has no source statement", "No reported gaps does not establish model completeness"} {
		if !strings.Contains(text, want) {
			t.Fatal(want, text)
		}
	}
	legacy, e := v.Select(Query{Viewpoint: "VP11"})
	if e != nil || legacy.Tables("VP11") != before {
		t.Fatal("ordinary registry changed", e)
	}
	again, e := v.Select(Query{Viewpoint: "VP11", Projection: "overview"})
	if e != nil || again.TraceabilityOverviewContext() != text || !reflect.DeepEqual(s, again) {
		t.Fatal("unstable", e)
	}
	for _, q := range []Query{{Viewpoint: "VP11", Projection: "overview", Mode: "Online"}, {Viewpoint: "VP11", Projection: "overview", Focus: "Work"}, {Viewpoint: "VP11", Projection: "overview", Relations: []string{"has"}}, {Viewpoint: "VP11", Projection: "overview", Diagram: "VP10-Big"}, {Viewpoint: "VP11", Projection: "overview", Level: "A4"}, {Viewpoint: "VP11", Projection: "overview", Depth: 3}, {Viewpoint: "VP11", Projection: "overview", Direction: "out"}} {
		if _, e := v.Select(q); e == nil {
			t.Fatal("ignored partial filter", q)
		}
	}
}

func TestTraceabilityOverviewEmptyLargeAndPathLabels(t *testing.T) {
	v, e := New(canonicalOverview(t, "language design-core version 0.5.\n"))
	if e != nil {
		t.Fatal(e)
	}
	s, e := v.Select(Query{Viewpoint: "VP11", Projection: "overview"})
	if e != nil {
		t.Fatal(e)
	}
	for _, want := range []string{"| Declarations | 0 |", "| Facts | 0 |", "No property assignments are declared", "No declarations are present", "No facts are declared", "snapshot (no declared subject)"} {
		if !strings.Contains(s.TraceabilityOverviewContext(), want) {
			t.Fatal(want)
		}
	}
	// Exercise the explicit zero-reported-gap state without claiming the producer's
	// current detectors establish completeness for any particular language model.
	zero := *s
	zero.Gaps = []Fact{}
	if !strings.Contains(zero.TraceabilityOverviewContext(), "No model gaps are reported by the current rules") {
		t.Fatal("missing explicit state")
	}
	decls, stmts := "language design-core version 0.5.\nunit Owner.\n", ""
	for i := 0; i < 170; i++ {
		decls += fmt.Sprintf("functionality Work%d.\n", i)
		stmts += fmt.Sprintf("Owner owns Work%d.\nWork%d has repeatability = deterministic.\n", i, i)
	}
	v, e = New(canonicalOverview(t, decls+stmts))
	if e != nil {
		t.Fatal(e)
	}
	s, e = v.Select(Query{Viewpoint: "VP11", Projection: "overview"})
	if e != nil || len(s.Facts) != 340 || len(s.Model.Declarations) != 171 || strings.Count(s.TraceabilityOverviewContext(), ": repeatability = deterministic.") != 170 {
		t.Fatal("truncated", e)
	}
	for _, f := range s.Facts {
		if !strings.Contains(s.TraceabilityOverviewContext(), f.S("text")) {
			t.Fatal("missing large fact", f)
		}
	}
	got := traceabilityPosition("Parts_[draft]&<x>.design", 3)
	if got != "Parts\\_\\[draft\\]&amp;&lt;x&gt;.design, line 3" {
		t.Fatal("path markup", got)
	}
}
