package fynehost

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func TestExtendedTextEmptyMeasurementMatchesNativeContent(t *testing.T) {
	test.NewTempApp(t)
	for _, multi := range []bool{false, true} {
		_, roots, err := parser.Compile(fmt.Sprintf(`sdui 0.3;page=[edit=input("Label",multiline=%t,placeholder="Placeholder")];`, multi))
		if err != nil {
			t.Fatal(err)
		}
		var n *parser.Instance
		roots["page"].Walk(func(x *parser.Instance) {
			if x.Widget == "input" {
				n = x
			}
		})
		f := ui.FieldState{Input: &ui.InputState{Multiline: multi, Placeholder: "Placeholder"}}
		for _, font := range []float64{14, 20, 28} {
			base, err := measureTextField(n, f, font, layout.Size{})
			if err != nil {
				t.Fatal(err)
			}
			for _, value := range []string{"", "short 世界", strings.Repeat("wide世界", 600), strings.Repeat("long wrapped words 世界\r\n", 180)} {
				if !multi && strings.ContainsAny(value, "\r\n") {
					continue
				}
				native := newTextControlValue(nil, n, value)
				body := container.NewThemeOverride(native.entry, componentTheme{float32(font)}).MinSize()
				if math.Abs(float64(body.Height)-base.Control.H) > .001 || float64(body.Width) > base.Control.W+.001 {
					t.Fatalf("empty probe differs mode=%v font=%v bytes=%d native=%v parts=%+v", multi, font, len(value), body, base)
				}
				empty := newTextControlValue(nil, n, "")
				emptyMin := container.NewThemeOverride(empty.entry, componentTheme{float32(font)}).MinSize()
				if body != emptyMin {
					t.Fatalf("content changes native minimum mode=%v font=%v bytes=%d %v/%v", multi, font, len(value), body, emptyMin)
				}
			}
		}
	}
}
