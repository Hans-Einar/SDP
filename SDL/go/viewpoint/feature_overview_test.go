package viewpoint

import (
	"strings"
	"testing"
)

func TestVP07OverviewArchitectureAndModeIsolation(t *testing.T) {
	source := `language design-core version 0.5.
unit Root.
unit Owner.
unit Unrelated.
container Local.
container Remote.
functionality Work.
functionality Missing.
feature Feature.
feature Isolated.
mode Online.
mode Offline.
Root contains Owner.
Owner owns Work.
Owner owns Missing.
Work contributes-to Feature.
Missing contributes-to Feature.
Work allocated-to Local in mode Offline.
Work allocated-to Remote in mode Online.
`
	v, e := New(canonicalOverview(t, source))
	if e != nil {
		t.Fatal(e)
	}
	for _, mode := range []string{"", "Offline", "Online"} {
		s, e := v.Select(Query{Viewpoint: "VP07", Projection: "overview", Mode: mode})
		if e != nil {
			t.Fatal(e)
		}
		if len(s.Diagrams) != 1 || len(s.Diagrams[0].Nodes) != 6 {
			t.Fatalf("incomplete architecture: %#v", s.Diagrams)
		}
		d := s.Diagrams[0]
		syntax := d.Mermaid()
		for _, part := range []string{"Unrelated", "Isolated", "n_Isolated(\"", "owns: F01", "-->|contains|"} {
			if !strings.Contains(syntax, part) {
				t.Fatal(part, syntax)
			}
		}
		alloc := s.Relations["allocated-to"]
		if mode == "" {
			if len(alloc) != 0 || strings.Contains(syntax, "allocated-to") {
				t.Fatal("logical view mixed allocations")
			}
		} else {
			if len(alloc) != 1 || alloc[0].S("mode") != mode {
				t.Fatal("cross-mode allocation", alloc)
			}
			found := false
			for _, g := range s.Gaps {
				if g.S("code") == "UNSPECIFIED_ALLOCATION" {
					if g.S("model_id") != "Missing" || g.S("mode") != mode {
						t.Fatal(g)
					}
					found = true
				}
			}
			if !found {
				t.Fatal("missing unallocated contribution gap")
			}
		}
		for _, annotation := range d.Elements {
			if annotation.S("role") == "allocated-to" {
				expected := "n_Local"
				if mode == "Online" {
					expected = "n_Remote"
				}
				if mode == "" || annotation.S("node") != expected {
					t.Fatal("cross-mode annotation", annotation)
				}
			}
		}
		proof := map[string]bool{}
		for _, f := range s.Facts {
			proof[f.S("id")] = true
		}
		for _, annotation := range d.Elements {
			for _, id := range annotation["source_facts"].([]string) {
				if !proof[id] {
					t.Fatal("unproven annotation", annotation)
				}
			}
		}
		for _, edge := range d.Edges {
			if !proof[edge.Fact] {
				t.Fatal("invented edge", edge)
			}
		}
		if !strings.Contains(s.FeatureOverviewContext(mode), "Isolated | Not specified") {
			t.Fatal("isolated feature lost")
		}
	}
	if _, e = v.Select(Query{Viewpoint: "VP07", Projection: "overview", Mode: "NoSuchMode"}); e == nil {
		t.Fatal("unknown mode accepted")
	}
	if _, e = v.Select(Query{Viewpoint: "VP01", Projection: "overview", Mode: "Online"}); e == nil {
		t.Fatal("VP01 mode ignored")
	}
	for _, d := range v.Diagrams {
		if d.ID == "VP07-overview" {
			t.Fatal("overview generated eagerly")
		}
	}
}

func TestVP07OverviewEmptyAndNoModes(t *testing.T) {
	for _, source := range []string{"language design-core version 0.5.\nactor Alone.\n", "language design-core version 0.5.\nfeature Alone.\n"} {
		v, e := New(source)
		if e != nil {
			t.Fatal(e)
		}
		s, e := v.Select(Query{Viewpoint: "VP07", Projection: "overview"})
		if e != nil {
			t.Fatal(e)
		}
		if len(s.Gaps) == 0 {
			t.Fatal("missing truthful empty-model gap")
		}
		if v.Kinds["Alone"] == "actor" && len(s.Diagrams) != 0 {
			t.Fatal("invented architecture")
		}
		if v.Kinds["Alone"] == "feature" && (len(s.Diagrams) != 1 || len(s.Diagrams[0].Nodes) != 1) {
			t.Fatal("isolated feature omitted")
		}
	}
}
