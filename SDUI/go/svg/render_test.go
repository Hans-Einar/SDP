package svg

import (
	"encoding/xml"
	"io"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

func TestExportGeometryAndLiterals(t *testing.T) {
	_, roots, e := parser.Compile(`sdui 0.2; page=[header="# <script>",<ok=button("<&>"),edit=input("Value",value="42")>]*b {scale=1};`)
	if e != nil {
		t.Fatal(e)
	}
	b, e := (&layout.Engine{}).Layout(roots["page"], layout.Size{W: 500, H: 300})
	if e != nil {
		t.Fatal(e)
	}
	s, e := Render(b, Options{Width: 500, Height: 300})
	if e != nil {
		t.Fatal(e)
	}
	if strings.Contains(s, "<script>") || !strings.Contains(s, "&lt;&amp;&gt;") {
		t.Fatal("literal escaped incorrectly")
	}
	d := xml.NewDecoder(strings.NewReader(s))
	rects, paths := 0, 0
	for {
		token, e := d.Token()
		if e == io.EOF {
			break
		}
		if e != nil {
			t.Fatal(e)
		}
		if x, ok := token.(xml.StartElement); ok {
			switch x.Name.Local {
			case "rect":
				rects++
			case "path":
				paths++
			}
		}
	}
	if rects < 7 || paths != 3 {
		t.Fatalf("missing geometry %d %d", rects, paths)
	}
	again, e := Render(b, Options{Width: 500, Height: 300})
	if e != nil || again != s {
		t.Fatal("non-deterministic export")
	}
}
