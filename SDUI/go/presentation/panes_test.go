package presentation

import (
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

func TestPaneDescriptionsIncludeInactiveDeclarations(t *testing.T) {
	const source = `sdui 0.3; ref: app "absent.sdl";
Main=[s=split("vertical",minFirst=0.2)[
left=tabs("Workspace",callback=app.Page.@invoke)[
a=page("First")["FIRST BODY"];b=page("Second",icon="notes")["SECOND BODY"] {visible=false}];
right=[]]];`
	doc, roots, err := parser.Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	for _, render := range []func() (string, error){
		func() (string, error) { return Dump(roots["Main"], 160) },
		func() (string, error) { return Markdown(roots["Main"], 160) },
		func() (string, error) { return Combined(roots["Main"], 160, doc) },
	} {
		out, err := render()
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range []string{"Static split", "axis=vertical", "minFirst=0.2", "Static tabs", "initial selection=first eligible", "not live state", "Static page", "Main/s/left/b", "Second", "symbolic icon=notes", "visible=false", "app.Page.@invoke", "not executed", "source="} {
			if !strings.Contains(out, want) {
				t.Errorf("missing %q in %s", want, out)
			}
		}
	}
	// Unsupported future kinds must reject even inside hidden page content.
	roots["Main"].Rows[0][0].Rows[0][0].Rows[1][0].Rows[0][0].Kind = "composition"
	roots["Main"].Rows[0][0].Rows[0][0].Rows[1][0].Rows[0][0].Widget = "futurePane"
	if out, err := Combined(roots["Main"], 160, doc); err == nil || out != "" {
		t.Fatal("hidden future composition silently accepted", err)
	}
}
