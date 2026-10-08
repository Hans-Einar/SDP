package svg

import (
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strings"
	"testing"
)

func TestScalarSVGExplicitRejectionAndNativeInventory(t *testing.T) {
	for _, body := range []string{`checkbox("C")`, `slider("S",min=0,max=10,step=1,value=1)`, `number("N",min=0,max=1,step=0.1,value=0.3)`, `select("S")`} {
		_, roots, err := parser.Compile(`sdui 0.3; Main=[field=` + body + ` {visible=false}];`)
		if err != nil {
			t.Fatal(err)
		}
		root := roots["Main"]
		field := root.Rows[0][0]
		box := &layout.Box{Instance: root, Path: root.Path}
		if out, err := Render(box, Options{}); err == nil || out != "" || !strings.Contains(err.Error(), "unsupported-value-export") || !strings.Contains(err.Error(), field.Path) {
			t.Fatal(out, err)
		}
		opts := Options{SkipControls: true, NativeControls: map[string]string{field.Path: field.Widget}}
		if err := Check(root, opts); err != nil {
			t.Fatal(err)
		}
		opts.NativeControls[field.Path] = "button"
		if err := Check(root, opts); err == nil {
			t.Fatal("wrong native kind accepted")
		}
		delete(opts.NativeControls, field.Path)
		if err := Check(root, opts); err == nil {
			t.Fatal("silent scalar omission")
		}
		opts.NativeControls[field.Path] = field.Widget
		opts.NativeControls["Main/extra"] = "checkbox"
		if err := Check(root, opts); err == nil {
			t.Fatal("extra inventory accepted")
		}
	}
}
func TestScalarNativeDialogContextValidatesLexemes(t *testing.T) {
	_, roots, err := parser.Compile(`sdui 0.3; Main=[d=dialog("D")[n=number("N",min=0,max=1,step=0.1,value=0.3)]];`)
	if err != nil {
		t.Fatal(err)
	}
	root := roots["Main"]
	dialog := root.Rows[0][0]
	opts := Options{SkipControls: true, InteractionRoot: root, NativeControls: map[string]string{"Main/d": "dialog", "Main/d/n": "number"}}
	if err := Check(dialog, opts); err != nil {
		t.Fatal(err)
	}
	dialog.Rows[0][0].Arguments["value"] = parser.Literal{Kind: "number", Value: .3}
	if err := Check(dialog, opts); err == nil {
		t.Fatal("native inventory licensed lost numeric provenance")
	}
}
