package viewpoint

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const channelOverviewFixture = `language design-core version 0.5.
unit Client.
unit Service.
unit Remote.
channel Bus.
channel Idle.
contract Events.
contract Body.
message Notice.
message Update.
field Value.
mode Online.
mode Offline.
mode Quiet.
scenario Example.
scenario Away.
Body has completeness = closed.
Body has-field Value.
Value has presence = required.
Value has value-type = text.
Events has completeness = closed.
Bus upholds Events.
Idle upholds Events.
Events permits Notice.
Events permits Update.
Notice upholds Body.
Notice has message-kind = event.
Update upholds Body.
Update has message-kind = event.
Client uses Bus as sender of Notice in mode Online.
Client uses Bus as sender of Update in mode Online.
Service uses Bus as receiver of Notice in mode Online.
Service uses Bus as receiver of Update in mode Online.
Remote uses Bus as sender of Notice in mode Offline.
Client uses Bus as receiver of Notice in mode Offline.
Example has completeness = open.
Example runs-in Online.
Example step 1 sends Notice from Client to Service via Bus.
Example step 2 sends Update from Client to Service via Bus.
Away has completeness = open.
Away runs-in Offline.
Away step 1 sends Notice from Remote to Client via Bus.
`

func TestChannelOverviewIsolationProofAndDetails(t *testing.T) {
	v, e := New(canonicalOverview(t, channelOverviewFixture))
	if e != nil {
		t.Fatal(e)
	}
	index, e := v.Select(Query{Viewpoint: "VP08", Projection: "overview"})
	if e != nil {
		t.Fatal(e)
	}
	if len(index.Diagrams) != 0 || len(index.MessageSets) != 0 {
		t.Fatal("eager/merged topology")
	}
	for _, f := range index.Facts {
		if f.S("verb") == "uses" {
			t.Fatal("index mode leak", f)
		}
	}
	text := v.ChannelOverviewContext("")
	for _, want := range []string{"| Offline | 2 | 1 |", "| Online | 4 | 1 |", "| Quiet | 0 | 0 |", "| Example | Online | open | 2 |", "| Away | Offline | open | 1 |", "| Idle | Events | Notice, Update |"} {
		if !strings.Contains(text, want) {
			t.Fatal(want, text)
		}
	}
	for _, mode := range []string{"Online", "Offline", "Quiet"} {
		s, e := v.Select(Query{Viewpoint: "VP08", Projection: "overview", Mode: mode})
		if e != nil {
			t.Fatal(e)
		}
		if len(s.Diagrams) != 2 {
			t.Fatal(s.Diagrams)
		}
		d := s.Diagrams[0]
		nodes, edges := 3, 2
		if mode == "Quiet" {
			nodes, edges = 1, 0
		}
		if len(d.Nodes) != nodes || len(d.Edges) != edges || d.Nodes["n_Bus"] == nil {
			t.Fatal(mode, d)
		}
		if len(s.Diagrams[1].Nodes) != 1 || s.Diagrams[1].Nodes["n_Idle"] == nil || len(s.Diagrams[1].Edges) != 0 {
			t.Fatal("lost isolated channel")
		}
		if !reflect.DeepEqual(d.SourceFacts, values(v.channelBindings(mode), "id")) {
			t.Fatal(d.SourceFacts)
		}
		proof := []string{}
		for i, x := range d.Elements {
			ids := x["proof"].([]string)
			proof = append(proof, ids...)
			if mode == "Online" && len(ids) != 2 {
				t.Fatal("lost grouped message proof", x)
			}
			edge := d.Edges[i]
			a, b := "n_"+x.S("participant"), "n_"+x.S("channel")
			if x.S("role") == "receiver" {
				a, b = b, a
			}
			if edge.Source != a || edge.Target != b || edge.Fact != ids[0] {
				t.Fatal(edge, x)
			}
			for _, id := range ids {
				f := eq(s.Facts, "id", id)
				if len(f) != 1 || f[0].S("mode") != mode || f[0].S("subject") != x.S("participant") {
					t.Fatal(f, x)
				}
			}
		}
		if !reflect.DeepEqual(unique(proof), unique(d.SourceFacts)) {
			t.Fatal("incomplete edge proof")
		}
		for _, m := range s.MessageSets {
			if m.S("mode") != mode {
				t.Fatal("messageset mode leak")
			}
		}
		for _, f := range s.Facts {
			if f.S("verb") == "uses" && f.S("mode") != mode || f.S("verb") == "runs-in" && f.S("object") != mode {
				t.Fatal("fact mode leak", f)
			}
		}
		txt := v.ChannelOverviewContext(mode)
		if mode == "Online" && (strings.Contains(txt, "Away") || strings.Contains(txt, "Remote")) {
			t.Fatal("context mode leak", txt)
		}
		if mode == "Quiet" && (!strings.Contains(txt, "No participation is declared") || !strings.Contains(txt, "No Scenarios are declared")) {
			t.Fatal(txt)
		}
		again, e := v.Select(Query{Viewpoint: "VP08", Projection: "overview", Mode: mode})
		if e != nil || !reflect.DeepEqual(s, again) {
			t.Fatal("nondeterministic", e)
		}
	}
	for _, d := range v.Diagrams {
		if strings.HasPrefix(d.ID, "VP08-") {
			detail, e := v.Select(Query{Viewpoint: "VP08", Diagram: d.ID})
			if e != nil || len(detail.Diagrams) != 1 || !reflect.DeepEqual(d, detail.Diagrams[0]) {
				t.Fatal("sequence changed", e)
			}
		}
	}
	for _, q := range []Query{{Viewpoint: "VP08", Projection: "overview", Mode: "Missing"}, {Viewpoint: "VP08", Projection: "overview", Focus: "Bus"}, {Viewpoint: "VP08", Projection: "overview", Relations: []string{"uses"}}, {Viewpoint: "VP08", Projection: "overview", Diagram: "VP08-Example"}, {Viewpoint: "VP08", Projection: "overview", Depth: 3}, {Viewpoint: "VP08", Projection: "overview", Direction: "out"}} {
		if _, e := v.Select(q); e == nil {
			t.Fatal("ignored filter", q)
		}
	}
}
func TestChannelOverviewEmptyAndLarge(t *testing.T) {
	v, e := New(canonicalOverview(t, "language design-core version 0.5.\nunit Solo.\n"))
	if e != nil {
		t.Fatal(e)
	}
	s, e := v.Select(Query{Viewpoint: "VP08", Projection: "overview"})
	if e != nil || len(s.Diagrams) != 0 {
		t.Fatal(e)
	}
	for _, want := range []string{"No Channels are declared", "No Modes are declared", "No Scenarios are declared"} {
		if !strings.Contains(v.ChannelOverviewContext(""), want) {
			t.Fatal(want)
		}
	}
	src := channelOverviewFixture
	for i := 0; i < 170; i++ {
		src = strings.Replace(src, "Body has completeness", fmt.Sprintf("channel Idle%d.\nBody has completeness", i), 1)
		src += fmt.Sprintf("Idle%d upholds Events.\n", i)
	}
	v, e = New(canonicalOverview(t, src))
	if e != nil {
		t.Fatal(e)
	}
	s, e = v.Select(Query{Viewpoint: "VP08", Projection: "overview", Mode: "Quiet"})
	if e != nil || len(s.Diagrams) != 172 {
		t.Fatal("inventory truncated", e)
	}
	for _, d := range s.Diagrams {
		if len(d.Nodes) != 1 || len(d.Edges) != 0 {
			t.Fatal(d)
		}
	}
}
