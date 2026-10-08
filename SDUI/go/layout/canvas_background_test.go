package layout_test

import (
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SDUI/go/svg"
)

func TestSurfaceBackgroundUsesExactSelectedSnapshotAndScopedInventory(t *testing.T) {
	const source = `sdui 0.3;Main=[cmd=command("Toggle",toggle=true);d=dialog("Dialog")[control=button(command="cmd");"Visible body";nested=dialog("Nested")["Never paint"]] {scale=1,gap=0}] {scale=1,gap=0};`
	_, roots, err := parser.Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	root := roots["Main"]
	snapshot := runtime.Snapshot{Root: root, Surfaces: map[string]runtime.SurfaceState{"Main/d": {InstancePath: "Main/d", Open: true}, "Main/d/nested": {InstancePath: "Main/d/nested"}}}
	g, err := (&layout.Engine{}).LayoutCanvases(snapshot, layout.Size{W: 500, H: 300}, map[string]layout.Size{"Main/d": {W: 300, H: 200}})
	if err != nil {
		t.Fatal(err)
	}
	body := g.Surfaces["Main/d"].Root
	opts := svg.Options{Width: 300, Height: 200, SkipControls: true, InteractionRoot: root, NativeControls: map[string]string{"Main/d": "dialog", "Main/d/control": "button", "Main/d/nested": "dialog"}}
	output, err := svg.Render(body, opts)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output, "Visible body") || strings.Contains(output, "Never paint") {
		t.Fatal("wrong surface background content")
	}
	// The button's command is outside this dialog subtree; it needs the actual
	// full selected snapshot, not a subtree treated as an independent entry.
	opts.InteractionRoot = nil
	if out, err := svg.Render(body, opts); err == nil || out != "" {
		t.Fatal("missing selected-root context produced partial artifact")
	}
	_, other, err := parser.Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	opts.InteractionRoot = other["Main"]
	if out, err := svg.Render(body, opts); err == nil || out != "" {
		t.Fatal("different snapshot pointer accepted")
	}
	opts.InteractionRoot = root
	delete(opts.NativeControls, "Main/d/nested")
	if out, err := svg.Render(body, opts); err == nil || out != "" {
		t.Fatal("closed nested surface skipped adapter preflight")
	}
}
