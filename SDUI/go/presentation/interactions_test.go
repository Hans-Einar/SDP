package presentation

import (
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strings"
	"testing"
)

func TestInteractionDescriptionsAndStrictEntry(t *testing.T) {
	source := `sdui 0.3; ref: app "not-opened.sdl"; Main=[c=command("Change",toggle=true,checked=true,exclusive="g",key="Primary+K",callback=app.Do.@invoke) {visible=false};button(command="c",tooltip="Tip");menu("Actions")[menuGroup("Group")[item(command="c");separator()]];d=dialog("Edit",modal=false,callback=app.Save.@invoke)[field=input("Text");button("Cancel",effect="cancel")]];`
	doc, roots, err := parser.Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	for _, render := range []func() (string, error){func() (string, error) { return Dump(roots["Main"], 200) }, func() (string, error) { return Markdown(roots["Main"], 200) }, func() (string, error) { return Combined(roots["Main"], 200, doc) }} {
		out, err := render()
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range []string{"Static command", "Main/c", "checked=true", "exclusive=g", "Primary+K", "not live state", "not executed", "visible=false", "Static menu", "Static menuGroup", "Static item", "Static separator", "Static dialog", "initially closed", "Accept/Cancel/Close not simulated", "modal=false", "effect=cancel", "source="} {
			if !strings.Contains(out, want) {
				t.Errorf("missing %q in %s", want, out)
			}
		}
	}
	roots["Main"].Rows[1][0].Arguments["command"] = parser.Literal{Kind: "string", Value: "missing"}
	if out, err := Dump(roots["Main"], 200); err == nil || out != "" {
		t.Fatal("invalid reference produced artifact")
	}
}
