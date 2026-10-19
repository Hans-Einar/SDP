package viewpoint

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"strings"
	"testing"
)

func TestVP01OverviewCompleteAndSeparate(t *testing.T) {
	source := "language design-core version 0.5.\nactor User.\nactor IsolatedActor.\nfeature Feature.\nfunctionality Function.\nunit Owner.\n"
	for i := 0; i < 30; i++ {
		source += fmt.Sprintf("usecase Goal%d.\n", i)
	}
	source += "Owner owns Function.\nUser pursues Goal0.\nFeature supports Goal0.\nFunction contributes-to Feature.\n"
	v, e := New(canonicalOverview(t, source))
	if e != nil {
		t.Fatal(e)
	}
	s, e := v.Select(Query{Viewpoint: "VP01", Projection: "overview"})
	if e != nil {
		t.Fatal(e)
	}
	if len(s.Diagrams) != 1 || len(s.Diagrams[0].Nodes) != 32 || len(s.Diagrams[0].Edges) != 1 {
		t.Fatalf("incomplete: %#v", s.Diagrams)
	}
	seen := map[string]bool{}
	for _, data := range s.Diagrams[0].Nodes {
		seen[data.(map[string]any)["model_id"].(string)] = true
	}
	for name, kind := range v.Kinds {
		if (kind == "actor" || kind == "usecase") && !seen[name] {
			t.Fatal("missing", name)
		}
	}
	syntax := s.Diagrams[0].Mermaid()
	if !strings.HasPrefix(syntax, "usecase-beta") || strings.Contains(syntax, "Feature") || strings.Contains(syntax, "include") || strings.Contains(syntax, "systemBoundary") {
		t.Fatal(syntax)
	}
	context := s.OverviewContext()
	for _, text := range []string{"Feature | supports | Goal0", "Function | contributes-to | Feature", "not inferred"} {
		if !strings.Contains(context, text) {
			t.Fatal(context)
		}
	}
	for _, d := range v.Diagrams {
		if d.Kind == "usecase" {
			t.Fatal("overview generated eagerly")
		}
	}
	for _, q := range []Query{{Viewpoint: "VP06", Projection: "overview"}, {Viewpoint: "VP01", Projection: "bad"}, {Viewpoint: "VP01", Projection: "overview", Focus: "Goal0"}, {Viewpoint: "VP01", Projection: "overview", Relations: []string{"pursues"}}, {Viewpoint: "VP01", Projection: "overview", Depth: 1}, {Viewpoint: "VP01", Projection: "overview", Direction: "out"}, {Viewpoint: "VP01", Projection: "overview", Diagram: "VP01-Goal0"}} {
		if _, e := v.Select(q); e == nil {
			t.Fatal("ignored overview filter", q)
		}
	}
	q, e := ParseURI("sdl-view://trial/VP01?projection=overview&target=main&consumer=xfmd")
	if e != nil {
		t.Fatal(e)
	}
	if _, e = v.Select(q.Query); e != nil {
		t.Fatal(e)
	}
	// Existing diagram identity is not captured by the new projection query.
	existing, e := v.Select(Query{Viewpoint: "VP01", Diagram: "VP01-Goal0"})
	if e != nil || len(existing.Diagrams) != 1 || existing.Diagrams[0].Kind != "usecase" {
		t.Fatal(e)
	}
}
func TestVP01OverviewEmptyAndLimit(t *testing.T) {
	v, e := New("language design-core version 0.5.\nunit Internal.\n")
	if e != nil {
		t.Fatal(e)
	}
	s, e := v.Select(Query{Viewpoint: "VP01", Projection: "overview"})
	if e != nil || len(s.Diagrams) != 0 || len(s.Gaps) == 0 {
		t.Fatal(e, s)
	}
	source := "language design-core version 0.5.\n"
	for i := 0; i < 129; i++ {
		source += fmt.Sprintf("usecase Goal%d.\n", i)
	}
	v, e = New(canonicalOverview(t, source))
	if e != nil {
		t.Fatal(e)
	}
	if _, e = v.Select(Query{Viewpoint: "VP01", Projection: "overview"}); e == nil || !strings.Contains(e.Error(), "no partial") {
		t.Fatal(e)
	}
}

func canonicalOverview(t *testing.T, source string) string {
	t.Helper()
	m, e := parser.Parse(source)
	if e != nil {
		t.Fatal(e)
	}
	return parser.FormatFile(m)
}

func TestOverviewURIRejectsExplicitTraversal(t *testing.T) {
	if _, e := ParseURI("sdl-view://trial/VP01?projection=overview&depth=0"); e == nil {
		t.Fatal("explicit depth silently ignored")
	}
}
