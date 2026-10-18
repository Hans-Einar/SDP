package svg

import (
	"errors"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strings"
	"testing"
)

func TestExtendedInputSVGRejectsEvenFalseEmptyAndHidden(t *testing.T) {
	for _, arg := range []string{`multiline=false`, `readOnly=false`, `required=false`, `placeholder=""`, `multiline=true,value="å\n🙂"`} {
		_, roots, err := parser.Compile(`sdui 0.3; Main=[field=input("Text",` + arg + `) {visible=false}];`)
		if err != nil {
			t.Fatal(err)
		}
		root := roots["Main"]
		field := root.Rows[0][0]
		out, err := Render(&layout.Box{Instance: root, Path: root.Path}, Options{})
		var d *parser.Diagnostic
		if out != "" || !errors.As(err, &d) || d.Code != "unsupported-text-export" || d.Span != field.Span || !strings.Contains(d.Message, field.Path) {
			t.Fatal(out, err)
		}
		for _, inventory := range []map[string]string{nil, {field.Path: "button"}, {field.Path: "input", "extra": "input"}} {
			if err := Check(root, Options{SkipControls: true, NativeControls: inventory}); err == nil {
				t.Fatal("invalid prepared inventory", inventory)
			}
		}
		if err := Check(root, Options{SkipControls: true, NativeControls: map[string]string{field.Path: "input"}}); err != nil {
			t.Fatal(err)
		}
	}
	for _, profile := range []string{"0.2", "0.3"} {
		_, roots, err := parser.Compile(`sdui ` + profile + `; Main=[field=input("Text",value="old")];`)
		if err != nil {
			t.Fatal(err)
		}
		if err := Check(roots["Main"], Options{}); err != nil {
			t.Fatal("legacy export changed", err)
		}
	}
}
