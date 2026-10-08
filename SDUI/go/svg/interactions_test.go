package svg

import (
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strings"
	"testing"
)

func TestInteractionExportRequiresExactPreparedAdapters(t *testing.T) {
	_, roots, err := parser.Compile(`sdui 0.3; Main=[c=command("C",toggle=true);b=button(command="c");m=menu("M")[menuGroup("G")[item(command="c");separator();sub=menu("Sub",mode="submenu")[item(command="c")]]];d=dialog("D")[button(command="/c");field=input("Field")]];`)
	if err != nil {
		t.Fatal(err)
	}
	root := roots["Main"]
	for _, hidden := range []bool{false, true} {
		root.Rows[0][0].Layout["visible"] = !hidden
		box := &layout.Box{Instance: root, Path: root.Path}
		if out, err := Render(box, Options{}); err == nil || out != "" || !strings.Contains(err.Error(), "unsupported-interaction-export") {
			t.Fatal(out, err)
		}
	}
	controls := map[string]string{"Main/b": "button", "Main/m": "menu", "Main/m/$r0c0/sub": "menu", "Main/d": "dialog", "Main/d/$r0c0": "button", "Main/d/field": "input"}
	options := Options{SkipControls: true, NativeControls: controls}
	if err := Check(root, options); err != nil {
		t.Fatal(err)
	}
	for path, kind := range controls {
		delete(controls, path)
		if err := Check(root, options); err == nil {
			t.Error("missing adapter", path)
		}
		controls[path] = kind
	}
	controls["Main/c"] = "command"
	if err := Check(root, options); err == nil {
		t.Fatal("fabricated native command accepted")
	}
	delete(controls, "Main/c")
	controls["Main/m/$r0c0"] = "menuGroup"
	if err := Check(root, options); err == nil {
		t.Fatal("extra native group accepted")
	}
	delete(controls, "Main/m/$r0c0")
	dialog := root.Rows[3][0]
	scoped := Options{SkipControls: true, InteractionRoot: root, NativeControls: map[string]string{"Main/d": "dialog", "Main/d/$r0c0": "button", "Main/d/field": "input"}}
	if err := Check(dialog, scoped); err != nil {
		t.Fatal("full entry context", err)
	}
	scoped.InteractionRoot = nil
	if err := Check(dialog, scoped); err == nil {
		t.Fatal("standalone subtree resolved external reference")
	}
	scoped.InteractionRoot = root
	copy := *dialog
	if err := Check(&copy, scoped); err == nil {
		t.Fatal("wrong snapshot pointer accepted")
	}
	scoped.SkipControls = false
	if err := Check(dialog, scoped); err == nil {
		t.Fatal("public context allowance")
	}
}

func TestDecoratedLegacyButtonExportIsExplicit(t *testing.T) {
	_, roots, err := parser.Compile(`sdui 0.3; Main=[b=button("B",tooltip="Tip")];`)
	if err != nil {
		t.Fatal(err)
	}
	if parser.IsCommandButton(roots["Main"].Rows[0][0]) {
		t.Fatal("decoration promoted legacy button")
	}
	if err := Check(roots["Main"], Options{}); err == nil || !strings.Contains(err.Error(), "unsupported-interaction-export") {
		t.Fatal(err)
	}
	if err := Check(roots["Main"], Options{SkipControls: true, NativeControls: map[string]string{"Main/b": "button"}}); err != nil {
		t.Fatal(err)
	}
}
