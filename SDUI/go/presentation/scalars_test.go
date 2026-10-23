package presentation

import (
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strings"
	"testing"
)

func TestScalarStructuralDescriptions(t *testing.T) {
	source := `sdui 0.3; ref: app "not-read.sdl"; Main=[c=checkbox("Flag",value=true) {visible=false};n=number("Amount",min=-0,max=1e1,step=1.00e-1,value=0.30,readOnly=true,callback=app.Save.@invoke);s=slider("Slide",min=0,max=10,step=1,value=1);choice=select("Choice",value="same-label-ID",required=true)];`
	doc, roots, err := parser.Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	root := roots["Main"]
	for _, render := range []func() (string, error){func() (string, error) { return Dump(root, 400) }, func() (string, error) { return Markdown(root, 400) }, func() (string, error) { return Combined(root, 400, doc) }} {
		output, err := render()
		if err != nil {
			t.Fatal(err)
		}
		for _, want := range []string{"Static checkbox", "Static number", "Static slider", "Static select", "Main/c", "visible=false", "source initial values; not live state", "min=-0", "max=1e1", "step=1.00e-1", "value=0.30", "readOnly=true", "Commit callback=app.Save.@invoke", "not executed; binding not verified", "choice options not supplied; IDs are not labels", "value=same-label-ID", "source="} {
			if !strings.Contains(output, want) {
				t.Errorf("missing %q in %s", want, output)
			}
		}
	}
	root.Rows[0][0].Widget = "future-value"
	if output, err := Dump(root, 400); err == nil || output != "" {
		t.Fatal("hidden unknown widget omitted")
	}
}
