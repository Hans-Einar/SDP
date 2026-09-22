package parser

import (
	"os"
	"strings"
	"testing"
)

func TestClassProfile(t *testing.T) {
	b, e := os.ReadFile("../examples/runtime-classes.sdl")
	if e != nil {
		t.Fatal(e)
	}
	src := string(b)
	m, e := CheckClasses(src)
	if e != nil {
		t.Fatal(e)
	}
	if CanonicalClasses(m) != src || len(m.Declarations) != 5 || len(m.Statements) != 8 {
		t.Fatal("roundtrip")
	}
	for _, s := range m.Statements {
		if src[s.Span.Start:s.Span.End] != s.Sentence() {
			t.Fatal("source map", s)
		}
	}
	for _, c := range []struct{ old, new string }{
		{"as text", "as Missing"}, {"1 to 1 with Widget", "1 to many with Widget"}, {"0 to many with DomainModule", "4 to 2 with DomainModule"}, {"role Widgets", "role Owner"}, {"ownership = aggregation", "ownership = contains"}, {"RuntimeModules has ownership = aggregation.\n", ""}, {"RuntimeWidgets links Runtime", "RuntimeWidgets links Widget"}, {"Widget attribute Value as text.", "Widget attribute Value as text.\nWidget attribute Value as text."},
	} {
		if _, e := CheckClasses(strings.Replace(src, c.old, c.new, 1)); e == nil {
			t.Errorf("accepted %s", c.new)
		}
	}
	if _, e := CheckClasses(strings.Replace(src, "Runtime operation Reload returns boolean.", "Runtime operation Reload returns boolean.\nRuntime operation Run invokes Go.", 1)); e == nil {
		t.Fatal("invented execution")
	}
}
