package viewpoint

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const encodingOverviewFixture = `language design-core version 0.5.
contract Wire.
contract Record.
dataset State.
datagram Family.
variant Body.
variant Unencoded.
field Flag.
field Number.
field Payload.
field Optional.
encoding Big.
encoding Little.
unit Unrelated.
feature Feature.
functionality Work.
mode Online.
State upholds Record.
Unrelated holds State.
Record has completeness = open.
Family from State.
Family upholds Wire.
Wire has completeness = closed.
Wire defines Body.
Wire defines Unencoded.
Wire has-field Flag.
Body has-field Number.
Body has-field Payload.
Unencoded has-field Optional.
Flag has value-type = boolean.
Flag has presence = required.
Number has value-type = signed.
Number has presence = required.
Payload has value-type = bytes.
Payload has presence = required.
Optional has value-type = text.
Optional has presence = optional.
Big encodes Body.
Big has byte-order = big-endian.
Big has bit-order = most-significant-first.
Big places Flag at 0 bits 1.
Big places Number at 1 bits 12.
Big places Payload at 13 bits 16.
Little encodes Body.
Little has byte-order = little-endian.
Little has bit-order = least-significant-first.
Little places Flag at 0 bits 1.
Little places Number at 1 bits 12.
Little places Payload at 13 bits 16.
Unrelated owns Work.
Work contributes-to Feature.
`

func TestEncodingOverviewPositionsProofCoverageAndDetails(t *testing.T) {
	v, e := New(canonicalOverview(t, encodingOverviewFixture))
	if e != nil {
		t.Fatal(e)
	}
	s, e := v.Select(Query{Viewpoint: "VP10", Projection: "overview"})
	if e != nil {
		t.Fatal(e)
	}
	if len(s.Diagrams) != 2 || len(s.MessageSets) != 0 {
		t.Fatal(s.Diagrams)
	}
	for _, d := range s.Diagrams {
		if d.Kind != "packet" || len(d.Nodes) != 3 || len(d.Edges) != 0 || len(d.Elements) != 3 {
			t.Fatal(d)
		}
		if d.Mermaid() != "packet\n    0: \"Flag\"\n    1-12: \"Number\"\n    13-28: \"Payload\"\n" {
			t.Fatal(d.Mermaid())
		}
		for _, id := range d.SourceFacts {
			if len(eq(s.Facts, "id", id)) != 1 {
				t.Fatal("missing proof", id)
			}
		}
		for _, f := range d.Elements {
			fs := eq(s.Facts, "id", f.S("fact"))
			if len(fs) != 1 || fs[0].S("field") != f.S("field") || fs[0].N("offset") != f.N("first_bit") || fs[0].N("offset")+fs[0].N("width")-1 != f.N("last_bit") {
				t.Fatal(f, fs)
			}
		}
	}
	text := s.EncodingOverviewContext()
	for _, want := range []string{"Variant Body; Contract Wire (closed); big-endian, most-significant-first; 29 bits", "little-endian, least-significant-first; 29 bits", "| Flag | Header: Wire | 0–0 / 1 | boolean / required |", "| Number | Variant: Body | 1–12 / 12 | signed / required |", "| Payload | Variant: Body | 13–28 / 16 | bytes / required |", "| Unencoded | Wire | No Encoding declared; no fixed wire layout |", "| Body | Wire | Big, Little |", "does not implement serialization"} {
		if !strings.Contains(text, want) {
			t.Fatal(want, text)
		}
	}
	if len(s.Gaps) != 1 || s.Gaps[0].S("model_id") != "Unencoded" {
		t.Fatal(s.Gaps)
	}
	for _, f := range s.Facts {
		if f.S("verb") == "owns" || f.S("verb") == "contributes-to" {
			t.Fatal("unrelated fact", f)
		}
	}
	for _, d := range v.Diagrams {
		if strings.HasPrefix(d.ID, "VP10-") {
			detail, e := v.Select(Query{Viewpoint: "VP10", Diagram: d.ID})
			if e != nil || len(detail.Diagrams) != 1 || !reflect.DeepEqual(d, detail.Diagrams[0]) {
				t.Fatal("detail changed", e)
			}
		}
	}
	again, e := v.Select(Query{Viewpoint: "VP10", Projection: "overview"})
	if e != nil || !reflect.DeepEqual(s, again) {
		t.Fatal("unstable", e)
	}
	for _, q := range []Query{{Viewpoint: "VP10", Projection: "overview", Mode: "Online"}, {Viewpoint: "VP10", Projection: "overview", Focus: "Big"}, {Viewpoint: "VP10", Projection: "overview", Relations: []string{"places"}}, {Viewpoint: "VP10", Projection: "overview", Diagram: "VP10-Big"}, {Viewpoint: "VP10", Projection: "overview", Level: "A4"}, {Viewpoint: "VP10", Projection: "overview", Depth: 3}, {Viewpoint: "VP10", Projection: "overview", Direction: "out"}} {
		if _, e := v.Select(q); e == nil {
			t.Fatal("ignored partial filter", q)
		}
	}
}

func TestEncodingOverviewEmptyLargeAndInvalidLayouts(t *testing.T) {
	v, e := New(canonicalOverview(t, "language design-core version 0.5.\nunit Solo.\n"))
	if e != nil {
		t.Fatal(e)
	}
	s, e := v.Select(Query{Viewpoint: "VP10", Projection: "overview"})
	if e != nil || len(s.Diagrams) != 0 || len(s.Facts) != 0 {
		t.Fatal(e, s)
	}
	for _, want := range []string{"No Encodings are declared", "No Variants are declared"} {
		if !strings.Contains(s.EncodingOverviewContext(), want) {
			t.Fatal(want)
		}
	}
	decls, stmts := "", ""
	for i := 0; i < 30; i++ {
		name := fmt.Sprintf("Extra%d", i)
		decls += "encoding " + name + ".\n"
		stmts += fmt.Sprintf("%s encodes Body.\n%s has byte-order = big-endian.\n%s has bit-order = most-significant-first.\n%s places Flag at 0 bits 1.\n%s places Number at 1 bits 12.\n%s places Payload at 13 bits 16.\n", name, name, name, name, name, name)
	}
	source := strings.Replace(encodingOverviewFixture, "State upholds", decls+"State upholds", 1) + stmts
	v, e = New(canonicalOverview(t, source))
	if e != nil {
		t.Fatal(e)
	}
	s, e = v.Select(Query{Viewpoint: "VP10", Projection: "overview"})
	if e != nil || len(s.Diagrams) != 32 || strings.Count(s.EncodingOverviewContext(), "### Encoding ") != 32 {
		t.Fatal("truncated", e)
	}
	for _, change := range [][2]string{{"Wire has completeness = closed", "Wire has completeness = open"}, {"Big places Number at 1 bits 12", "Big places Number at 2 bits 12"}, {"Payload has presence = required", "Payload has presence = optional"}, {"Big places Flag at 0 bits 1", "Big places Flag at 0 bits 2"}, {"Big has byte-order = big-endian.\n", ""}, {"Big places Payload at 13 bits 16", "Big places Payload at 13 bits 15"}} {
		if _, e := New(canonicalOverview(t, strings.Replace(encodingOverviewFixture, change[0], change[1], 1))); e == nil {
			t.Fatal("invalid layout accepted", change)
		}
	}
}
