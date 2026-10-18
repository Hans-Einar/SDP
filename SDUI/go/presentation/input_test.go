package presentation

import (
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strings"
	"testing"
)

func TestExtendedInputStructuralDescription(t *testing.T) {
	doc, roots, err := parser.Compile(`sdui 0.3; ref: app "not-read"; Main=[note=input("Notes",value="å\n🙂",multiline=true,readOnly=false,required=true,placeholder="",callback=app.Save.@invoke) {visible=false};fallback=input("Fallback",required=false)];`)
	if err != nil {
		t.Fatal(err)
	}
	root := roots["Main"]
	for _, render := range []func() (string, error){func() (string, error) { return Dump(root, 400) }, func() (string, error) { return Markdown(root, 400) }, func() (string, error) { return Combined(root, 400, doc) }} {
		out, err := render()
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range []string{"Static input", "Main/note", `value="å\n🙂"`, "source initial text; not live state", "multiline=true", "readOnly=false", "required=true", `placeholder=""`, `placeholder="Fallback" (text fallback)`, "Commit callback=app.Save.@invoke", "not executed; binding not verified", "visible=false", "source="} {
			if !strings.Contains(out, want) {
				t.Errorf("missing %q: %s", want, out)
			}
		}
	}
}
