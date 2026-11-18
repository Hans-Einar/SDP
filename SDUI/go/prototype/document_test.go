package prototype

import (
	"encoding/json"
	"errors"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func documentFile(t *testing.T, path, contents string) string {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}
func documentFixture(t *testing.T) (source, contextPath, resource string) {
	t.Helper()
	dir := t.TempDir()
	source = documentFile(t, filepath.Join(dir, "page.sdui"), `sdui 0.3; ref: art "never-opened.sdl"; Main=[t=tabs("Tasks")[p=page("Hidden fixtures")[choice=select("Choice",value="one"); image=svg(art.X.@resource,description="Figure",fallback="reject")] {visible=false}]];`)
	contextPath = documentFile(t, filepath.Join(dir, "context.json"), `{"schema":"sdui-prototype-context/1","viewport":{"width":1680,"height":1150},"choices":{"Main/t/p/choice":[{"id":"one","label":"One","enabled":true}]},"svg":{"Main/t/p/image":{"file":"image.svg","width":640,"height":500}}}`)
	resource = documentFile(t, filepath.Join(dir, "image.svg"), `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 500"><rect x="0" y="0" width="640" height="500" fill="#ffffff"/></svg>`)
	return
}
func TestDocumentContextNativePreflightAndSnapshot(t *testing.T) {
	source, contextPath, resource := documentFixture(t)
	input, err := LoadDocument(source, "Main", "", contextPath)
	if err != nil {
		t.Fatal(err)
	}
	if input.Size.W != 1680 || input.Size.H != 1150 || len(input.ContextRevision) != 64 || len(input.Choices["Main/t/p/choice"]) != 1 || len(input.SVGResources) != 1 {
		t.Fatal(input)
	}
	report, err := CheckDocument(source, "Main", input.Candidate.Hash, contextPath)
	if err != nil {
		t.Fatal(err)
	}
	if report.Status != "preflight" || report.Host != "document" || report.NativeLayout != "pending" || report.ContextRevision != input.ContextRevision || !strings.Contains(report.Diagnostic, "pending") {
		t.Fatal(report)
	}
	documentFile(t, resource, `<svg viewBox="0 0 640 500"><circle cx="100" cy="100" r="50"/></svg>`)
	if err = input.Guard(); err == nil || !strings.Contains(err.Error(), "stale") {
		t.Fatal(err)
	}
	newer, err := LoadDocument(source, "Main", "", contextPath)
	if err != nil {
		t.Fatal(err)
	}
	if newer.ContextRevision == input.ContextRevision || newer.Candidate.Hash != input.Candidate.Hash {
		t.Fatal("resource change must change aggregate context identity only")
	}
	documentFile(t, source, `sdui 0.3; Main=[];`)
	if err = newer.Guard(); err == nil {
		t.Fatal("source guard accepted replaced model")
	}
	if _, err = LoadDocument(source, "Main", input.Candidate.Hash, contextPath); err == nil || !strings.Contains(err.Error(), "stale") {
		t.Fatal(err)
	}
}
func TestDocumentContextStrictSchema(t *testing.T) {
	cases := []string{
		`null`, `[]`, `{"schema":"sdui-prototype-context/1"} {}`, `{"schema":"sdui-prototype-context/1","schema":"sdui-prototype-context/1"}`,
		`{"Schema":"sdui-prototype-context/1"}`, `{"schema":"sdui-prototype-context/1","unknown":true}`, `{"schema":"wrong"}`,
		`{"schema":"sdui-prototype-context/1","viewport":null}`, `{"schema":"sdui-prototype-context/1","viewport":{"width":1680}}`,
		`{"schema":"sdui-prototype-context/1","viewport":{"width":-1,"height":800}}`, `{"schema":"sdui-prototype-context/1","viewport":{"width":1e999,"height":800}}`,
		`{"schema":"sdui-prototype-context/1","choices":null}`, `{"schema":"sdui-prototype-context/1","choices":{"P/f":null}}`,
		`{"schema":"sdui-prototype-context/1","choices":{"P/f":[{"id":"a","label":"A"}]}}`,
		`{"schema":"sdui-prototype-context/1","choices":{"P/f":[{"id":"a","label":"A","Enabled":true}]}}`,
		`{"schema":"sdui-prototype-context/1","choices":{"P/f":[],"P/f":[]}}`,
		`{"schema":"sdui-prototype-context/1","svg":{"P/f":{"file":"x.svg","width":1,"height":null}}}`,
	}
	for _, data := range cases {
		t.Run(data, func(t *testing.T) {
			if _, err := parseDocumentContext([]byte(data)); err == nil {
				t.Fatal("accepted malformed context", data)
			}
		})
	}
}
func TestDocumentContextExactPathsAndHiddenChoiceInventory(t *testing.T) {
	source, contextPath, _ := documentFixture(t)
	for _, change := range []func(map[string]any){
		func(c map[string]any) { delete(c, "choices") },
		func(c map[string]any) { c["choices"] = map[string]any{"choice": []any{}} },
		func(c map[string]any) { c["choices"] = map[string]any{"Main/t/p/image": []any{}} },
		func(c map[string]any) {
			c["choices"] = map[string]any{"Main/t/p/choice": []any{map[string]any{"id": "one", "label": "One", "enabled": false}}}
		},
		func(c map[string]any) {
			c["svg"] = map[string]any{"Main/t/p/choice": map[string]any{"file": "image.svg", "width": 640, "height": 500}}
		},
		func(c map[string]any) {
			c["svg"] = map[string]any{"Main/missing": map[string]any{"file": "image.svg", "width": 640, "height": 500}}
		},
	} {
		original, _ := os.ReadFile(contextPath)
		var c map[string]any
		if err := json.Unmarshal(original, &c); err != nil {
			t.Fatal(err)
		}
		change(c)
		encoded, _ := json.Marshal(c)
		documentFile(t, contextPath, string(encoded))
		if _, err := LoadDocument(source, "Main", "", contextPath); err == nil {
			t.Fatal("accepted wrong/missing fixture", string(encoded))
		}
		documentFile(t, contextPath, string(original))
	}
	// No choices.json discovery even when a plausible sibling fixture exists.
	documentFile(t, filepath.Join(filepath.Dir(source), "choices.json"), `{"choice":[["one","One"]]}`)
	if _, err := LoadDocument(source, "Main", "", ""); err == nil || !strings.Contains(err.Error(), "missing Main/t/p/choice") {
		t.Fatal(err)
	}
}
func TestDocumentContextResourceContainmentAndLimits(t *testing.T) {
	source, contextPath, _ := documentFixture(t)
	outside := documentFile(t, filepath.Join(t.TempDir(), "outside.svg"), `<svg viewBox="0 0 1 1"/>`)
	symlink := filepath.Join(filepath.Dir(contextPath), "escape.svg")
	if err := os.Symlink(outside, symlink); err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{outside, "../outside.svg", "sub/../../outside.svg", "escape.svg", "missing.svg"} {
		original, _ := os.ReadFile(contextPath)
		context := strings.Replace(string(original), "image.svg", file, 1)
		documentFile(t, contextPath, context)
		if _, err := LoadDocument(source, "Main", "", contextPath); err == nil {
			t.Fatal("accepted escaped/missing resource", file)
		}
		documentFile(t, contextPath, string(original))
	}
	original, _ := os.ReadFile(contextPath)
	documentFile(t, contextPath, strings.Repeat(" ", contextLimit+1))
	if _, err := LoadDocument(source, "Main", "", contextPath); err == nil {
		t.Fatal("accepted oversized context")
	}
	documentFile(t, contextPath, string(original))
	documentFile(t, filepath.Join(filepath.Dir(contextPath), "image.svg"), strings.Repeat("x", svgLimit+1))
	if _, err := LoadDocument(source, "Main", "", contextPath); err == nil {
		t.Fatal("accepted oversized SVG")
	}
}
func TestDocumentPreflightNoInventedCollectionProviders(t *testing.T) {
	source := documentFile(t, filepath.Join(t.TempDir(), "page.sdui"), `sdui 0.3; Main=[items=list("Items") {visible=false}];`)
	report, err := CheckDocument(source, "Main", "", "")
	var diagnostic *preparation.Diagnostic
	if !errors.As(err, &diagnostic) || diagnostic.Capability != (preparation.Capability{Dimension: preparation.Provider, ID: "collection-data", Major: 1}) || diagnostic.Path != "Main/items" || diagnostic.Span.Line == 0 || report.Status != "" {
		t.Fatal(report, err)
	}
	if _, err := LoadDocument(source, "Main", "", ""); !errors.As(err, &diagnostic) {
		t.Fatal("native loader must reject absent collection provider", err)
	}
}
func TestDocumentDefaultViewportAndLegacyCheckPolicy(t *testing.T) {
	source := documentFile(t, filepath.Join(t.TempDir(), "page.sdui"), `sdui 0.3; Main=[t=tabs("T")[p=page("P")[]]];`)
	input, err := LoadDocument(source, "Main", "", "")
	if err != nil {
		t.Fatal(err)
	}
	if input.Size.W != 1280 || input.Size.H != 800 || input.ContextRevision != "" {
		t.Fatal(input)
	}
	if _, err := CheckDocument(source, "Main", "", ""); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Check(source, "Main", ""); err == nil || !strings.Contains(err.Error(), "unsupported-pane") {
		t.Fatal("legacy policy changed", err)
	}
}
