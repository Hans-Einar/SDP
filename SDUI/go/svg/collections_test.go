package svg

import (
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

func TestCollectionExportAndScopedNativeOmission(t *testing.T) {
	_, roots, err := parser.Compile(`sdui 0.3; page=[items=tree("Navigation"),entries=list("Entries"),ok=button("OK")];`)
	if err != nil {
		t.Fatal(err)
	}
	root := roots["page"]
	// Already supplied measured geometry cannot turn an unsupported collection
	// export into an empty or misleading successful artifact.
	box := &layout.Box{Instance: root, Path: root.Path, Rect: layout.Rect{W: 400, H: 300}, Clip: layout.Rect{W: 400, H: 300}}
	for _, opts := range []Options{{}, {SkipControls: true}, {SkipControls: true, NativeControls: map[string]string{"page/items": "tree"}}, {SkipControls: true, NativeControls: map[string]string{"page/items": "list", "page/entries": "list", "page/ok": "button"}}} {
		out, err := Render(box, opts)
		if err == nil || out != "" {
			t.Fatal("silently omitted controls", out, err)
		}
	}
	err = Check(root, Options{})
	d, ok := err.(*parser.Diagnostic)
	if !ok || d.Code != "unsupported-collection-export" || d.Span != root.Rows[0][0].Span || !strings.Contains(d.Message, "page/items") {
		t.Fatal("missing source-linked diagnostic", err)
	}
	opts := Options{Width: 400, Height: 300, SkipControls: true, NativeControls: map[string]string{"page/items": "tree", "page/entries": "list", "page/ok": "button"}}
	if out, err := Render(box, opts); err != nil || out == "" {
		t.Fatal("native inventory rejected", err)
	}
	opts.NativeControls["unused"] = "button"
	if err := Check(root, opts); err == nil {
		t.Fatal("unused native inventory accepted")
	}
	delete(opts.NativeControls, "unused")
	root.Rows[0][0].Widget = "unknown"
	opts.NativeControls["page/items"] = "unknown"
	if err := Check(root, opts); err == nil {
		t.Fatal("inventory licensed unknown kind")
	}
}

func TestHiddenCollectionExportRejected(t *testing.T) {
	_, roots, err := parser.Compile(`sdui 0.3; page=[items=list("Entries") {visible=false}];`)
	if err != nil {
		t.Fatal(err)
	}
	if err := Check(roots["page"], Options{}); err == nil {
		t.Fatal("hidden collection silently omitted")
	}
}
