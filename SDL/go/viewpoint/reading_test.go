package viewpoint

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestReadingSemanticMembershipAndSelectors(t *testing.T) {
	v, e := New(canonicalOverview(t, strings.Replace(channelOverviewFixture, "unit Client.", "interface Unused.\nactor Reader.\nusecase Browse.\nfeature Readable.\ncapability Inspect.\nfunctionality Render.\ncontainer Host.\nactivity Build.\nunit Client.", 1)+"Client owns Render.\nRender realizes Inspect.\n"))
	if e != nil {
		t.Fatal(e)
	}
	v.ReadingProject = "trial"
	for _, spec := range Catalog {
		seen := map[string]bool{}
		for _, c := range v.ReadingCollections(spec.ID) {
			for _, n := range c.Names {
				if seen[n] {
					t.Fatalf("duplicate primary route %s/%s", spec.ID, n)
				}
				seen[n] = true
				if v.Kinds[n] == "" {
					t.Fatalf("invented member %s", n)
				}
			}
			q := Query{Viewpoint: spec.ID, Projection: "reading", Collection: c.ID}
			s, e := v.Select(q)
			if e != nil {
				t.Fatal(q, e)
			}
			text := v.ReadingMarkdown(q, s, false)
			for _, n := range c.Names {
				if !strings.Contains(text, "## "+n+"\n") {
					t.Fatalf("missing chapter %s/%s", spec.ID, n)
				}
			}
			if strings.Contains(text, "## Complete fact registry") {
				t.Fatal("register appended", q)
			}
		}
		q := Query{Viewpoint: spec.ID, Projection: "reading"}
		s, e := v.Select(q)
		if e != nil {
			t.Fatal(q, e)
		}
		text := v.ReadingMarkdown(q, s, false)
		if !strings.Contains(text, readingQuestions[spec.ID]) {
			t.Fatal("missing viewpoint question", spec.ID)
		}
	}
	responsibilityQuery := Query{Viewpoint: "VP03", Projection: "reading", Focus: "Client"}
	responsibility, _ := v.Select(responsibilityQuery)
	responsibilityText := v.ReadingMarkdown(responsibilityQuery, responsibility, false)
	for _, statement := range []string{"Client owns Render.", "Render realizes Inspect."} {
		if !strings.Contains(responsibilityText, statement) {
			t.Fatal("missing responsibility chain", statement)
		}
	}
	if got := v.ReadingCollections("VP04")[0].Names; !reflect.DeepEqual(got, []string{"Unused"}) {
		t.Fatal("isolated interface absent", got)
	}
	invalid := []Query{
		{Viewpoint: "VP01", Projection: "reading", Collection: "actors", Focus: "Reader"},
		{Viewpoint: "VP01", Projection: "reading", Register: "facts"},
		{Viewpoint: "VP04", Projection: "reading", Collection: "invented"},
		{Viewpoint: "VP04", Collection: "interfaces"},
		{Viewpoint: "VP04", Projection: "reading", Page: 1},
		{Viewpoint: "VP04", Projection: "reading", Collection: "interfaces", Page: math.MaxInt},
		{Viewpoint: "VP04", Projection: "reading", Mode: "Unused"},
		{Viewpoint: "VP04", Projection: "reading", Diagram: "VP04-overview"},
	}
	for _, q := range invalid {
		if _, e := v.Select(q); e == nil {
			t.Fatalf("accepted %#v", q)
		}
	}
	uri := v.ReadingTarget(Query{Viewpoint: "VP08", Projection: "reading", Collection: "contracts", Mode: "Online", Page: 2})
	parsed, e := ParseURI(uri)
	if e != nil || parsed.Project != "trial" || parsed.Query.Collection != "contracts" || parsed.Query.Mode != "Online" || parsed.Query.Page != 2 {
		t.Fatal(uri, parsed, e)
	}
}

func TestReadingContractScenarioAndEncodingMeaning(t *testing.T) {
	v, e := New(canonicalOverview(t, channelOverviewFixture))
	if e != nil {
		t.Fatal(e)
	}
	for _, tc := range []struct {
		q      Query
		want   []string
		absent string
	}{
		{Query{Viewpoint: "VP08", Projection: "reading", Focus: "Body"}, []string{"Completeness: **closed**", "Value", "text", "required", "Used by:"}, "## Complete fact registry"},
		{Query{Viewpoint: "VP08", Projection: "reading", Focus: "Bus", Mode: "Online"}, []string{"Actual participants", "Client", "Service", "sender", "receiver", "Notice", "Update"}, "Remote uses Bus"},
		{Query{Viewpoint: "VP08", Projection: "reading", Focus: "Example", Mode: "Online"}, []string{"Ordered scenario steps", "| 1 |", "| 2 |", "Only these numbered steps"}, "Away step"},
		{Query{Viewpoint: "VP08", Projection: "reading", Focus: "Away", Mode: "Online"}, []string{"no declared applicability"}, "sequenceDiagram"},
	} {
		s, e := v.Select(tc.q)
		if e != nil {
			t.Fatal(e)
		}
		text := v.ReadingMarkdown(tc.q, s, false)
		for _, want := range tc.want {
			if !strings.Contains(text, want) {
				t.Fatal(tc.q, want, text)
			}
		}
		if strings.Contains(text, tc.absent) {
			t.Fatal("scope leak", tc.q, text)
		}
	}
	q := Query{Viewpoint: "VP08", Projection: "reading", Collection: "contracts"}
	s, _ := v.Select(q)
	collection := v.ReadingMarkdown(q, s, false)
	q.Collection = ""
	q.Focus = "Body"
	s, _ = v.Select(q)
	detail := v.ReadingMarkdown(q, s, false)
	chapter := detail[strings.Index(detail, "## Body\n"):]
	if !strings.Contains(collection, chapter) {
		t.Fatal("collection/detail diverged")
	}
	enc, e := New(canonicalOverview(t, encodingOverviewFixture))
	if e != nil {
		t.Fatal(e)
	}
	q = Query{Viewpoint: "VP10", Projection: "reading", Focus: "Big"}
	s, e = enc.Select(q)
	if e != nil {
		t.Fatal(e)
	}
	text := enc.ReadingMarkdown(q, s, false)
	for _, want := range []string{"big-endian", "most-significant-first", "Explicit packet layout", "0–0", "1–12", "13–28", "required"} {
		if !strings.Contains(text, want) {
			t.Fatal(want, text)
		}
	}
}

func TestReadingCompletePaginationAndExactSourceLabels(t *testing.T) {
	source := "language design-core version 0.5.\n"
	for i := 0; i < 53; i++ {
		source += fmt.Sprintf("interface Port%03d.\n", i)
	}
	v, e := New(source)
	if e != nil {
		t.Fatal(e)
	}
	v.ReadingProject = "trial"
	seen := map[string]bool{}
	for page := 1; page <= 3; page++ {
		q := Query{Viewpoint: "VP04", Projection: "reading", Collection: "interfaces", Page: page}
		s, e := v.Select(q)
		if e != nil {
			t.Fatal(e)
		}
		text := v.ReadingMarkdown(q, s, false)
		if !strings.Contains(text, "of 53") {
			t.Fatal("lost total", text)
		}
		for _, n := range v.readingNames(q) {
			if seen[n] {
				t.Fatal("repeated", n)
			}
			seen[n] = true
		}
		if page < 3 && !strings.Contains(text, fmt.Sprintf("page=%d", page+1)) {
			t.Fatal("missing next action")
		}
	}
	if len(seen) != 53 {
		t.Fatal("omitted declarations", len(seen))
	}
	v, e = New(canonicalOverview(t, channelOverviewFixture))
	if e != nil {
		t.Fatal(e)
	}
	v.SourceBindings = map[string]string{"": "retained-source"}
	v.SourceLabels = map[string]string{"": "Example.design"}
	f := v.Facts[0]
	ref := v.FactReference(f.S("id"))
	span := f["span"].(map[string]any)
	want := fmt.Sprintf("sdl-source://retained-source/%d/%d", span["start"], span["end"])
	if !strings.Contains(ref, want) || !strings.Contains(ref, f.S("text")) || !strings.Contains(ref, "Example.design, line") || strings.Contains(ref, "["+f.S("id")+"]") {
		t.Fatal(ref, want)
	}
}

func TestReadingModeFiltersGeometryBeforeRestoringProof(t *testing.T) {
	v, e := New(canonicalOverview(t, channelOverviewFixture))
	if e != nil {
		t.Fatal(e)
	}
	bindings := eq(v.Facts, "verb", "uses")
	d := Diagram{ID: "VP08-Bus", Kind: "flowchart", Title: "Mixed static participation", Nodes: v.nodes([]string{"Bus", "Client", "Service", "Remote"}), SourceFacts: values(bindings, "id")}
	for _, f := range bindings {
		d.Edges = append(d.Edges, Edge{"n_" + f.S("subject"), "uses", "n_Bus", f.S("id")})
	}
	d.Elements = []Fact{{"mode": "Online", "proof": values(eq(bindings, "mode", "Online"), "id")}, {"mode": "Offline", "proof": values(eq(bindings, "mode", "Offline"), "id")}}
	v.Diagrams = append(v.Diagrams, d)
	q := Query{Viewpoint: "VP08", Projection: "reading", Focus: "Bus", Mode: "Online"}
	selected, e := v.Select(q)
	if e != nil {
		t.Fatal(e)
	}
	if len(selected.Diagrams) != 1 || len(selected.Diagrams[0].Edges) != 4 || selected.Diagrams[0].Nodes["n_Remote"] != nil || len(selected.Diagrams[0].Elements) != 1 || selected.Diagrams[0].Elements[0].S("mode") != "Online" {
		t.Fatal("mode geometry leak", selected.Diagrams)
	}
	for _, f := range selected.Facts {
		if f.S("mode") != "" && f.S("mode") != "Online" {
			t.Fatal("mode proof leak", f)
		}
	}
	for _, id := range selected.Diagrams[0].SourceFacts {
		for _, f := range eq(v.Facts, "id", id) {
			if f.S("mode") != "" && f.S("mode") != "Online" {
				t.Fatal("mode diagram proof leak", f)
			}
		}
	}
}

func TestReadingRegistersPaginateCompleteScope(t *testing.T) {
	source := "language design-core version 0.5.\n"
	for i := 0; i < 205; i++ {
		source += fmt.Sprintf("activity Work%03d.\n", i)
	}
	for i := 0; i < 205; i++ {
		source += fmt.Sprintf("Work%03d has implementation-status = planned.\n", i)
	}
	v, e := New(canonicalOverview(t, source))
	if e != nil {
		t.Fatal(e)
	}
	for _, reg := range []string{"facts", "declarations"} {
		total := len(v.Facts)
		if reg == "declarations" {
			total = len(v.Model.Declarations)
		}
		covered := 0
		for page := 1; page <= (total+ReadingRegisterPageSize-1)/ReadingRegisterPageSize; page++ {
			q := Query{Viewpoint: "VP11", Projection: "reading", Register: reg, Page: page}
			selected, e := v.Select(q)
			if e != nil {
				t.Fatal(e)
			}
			text := v.ReadingMarkdown(q, selected, false)
			a, z, scope, e := v.readingPage(q)
			if e != nil || scope != total {
				t.Fatal(scope, total, e)
			}
			for i := a; i < z; i++ {
				want := ""
				if reg == "facts" {
					want = v.Facts[i].S("text")
				} else {
					want = "| " + v.Model.Declarations[i].Name.Name + " |"
				}
				if !strings.Contains(text, want) {
					t.Fatal("omitted register row", reg, i, want)
				}
			}
			covered += z - a
		}
		if covered != total {
			t.Fatal("incomplete register", reg, covered, total)
		}
	}
}
