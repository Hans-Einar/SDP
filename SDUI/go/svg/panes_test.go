package svg

import (
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

func TestPaneSVGRejectionAndExactNativeInventory(t *testing.T) {
	_, roots, err := parser.Compile(`sdui 0.3; Main=[s=split("horizontal")[left=tabs("T")[a=page("A")[ok=button("SECRET BUTTON")];b=page("B")[edit=input("SECRET INPUT")] {visible=false}];right=["Visible background"]]];`)
	if err != nil {
		t.Fatal(err)
	}
	root := roots["Main"]
	var box func(*parser.Instance) *layout.Box
	box = func(n *parser.Instance) *layout.Box {
		b := &layout.Box{Instance: n, Path: n.Path, Rect: layout.Rect{W: 400, H: 300}, Clip: layout.Rect{W: 400, H: 300}, Font: 14, Enabled: true}
		for _, row := range n.Rows {
			for _, child := range row {
				b.Children = append(b.Children, box(child))
			}
		}
		return b
	}
	measured := box(root)
	if out, err := Render(measured, Options{}); err == nil || out != "" || !strings.Contains(err.Error(), "unsupported-pane-export") || !strings.Contains(err.Error(), "Main/s") {
		t.Fatal(out, err)
	}
	root.Rows[0][0].Layout["visible"] = false
	if err := Check(root, Options{}); err == nil || !strings.Contains(err.Error(), "unsupported-pane-export") {
		t.Fatal("hidden pane export accepted", err)
	}
	opts := Options{Width: 400, Height: 300, SkipControls: true, NativeControls: map[string]string{
		"Main/s": "split", "Main/s/left": "tabs", "Main/s/left/a/ok": "button", "Main/s/left/b/edit": "input",
	}}
	if out, err := Render(measured, opts); err != nil || !strings.Contains(out, "Visible background") || strings.Contains(out, "SECRET") {
		t.Fatal("native background content/omission", out, err)
	}
	for path, kind := range opts.NativeControls {
		delete(opts.NativeControls, path)
		if err := Check(root, opts); err == nil {
			t.Error("missing prepared control accepted", path)
		}
		opts.NativeControls[path] = kind
	}
	opts.NativeControls["Main/s/left/a"] = "page"
	if err := Check(root, opts); err == nil {
		t.Fatal("extraneous page inventory accepted")
	}
	delete(opts.NativeControls, "Main/s/left/a")
	opts.NativeControls["Main/s"] = "tabs"
	if err := Check(root, opts); err == nil {
		t.Fatal("wrong pane inventory kind accepted")
	}
	opts.NativeControls["Main/s"] = "split"
	root.Rows[0][0].Rows[0][0].Widget = "dialog"
	if err := Check(root, opts); err == nil {
		t.Fatal("inventory licensed M2 composition")
	}
}
