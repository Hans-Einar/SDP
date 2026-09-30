// Standalone study probe. Run from the archived SDUI Go module, not as product code.
package main

import (
	"encoding/json"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
	"github.com/Hans-Einar/SDP/SDUI/go/svg"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/sfnt"
	"os"
	"strings"
)

func main() {
	f, err := sfnt.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
	var buf sfnt.Buffer
	glyphs := map[string]any{}
	for _, r := range []rune{'A', '↻', '↺', '→', '✓', 'ø', '界'} {
		index, e := f.GlyphIndex(&buf, r)
		var out strings.Builder
		renderErr := svg.Text(&out, string(r), 0, 20, 14, "black")
		glyphs[string(r)] = map[string]any{"index": index, "lookupError": e, "renderError": renderErr, "svgBytes": out.Len()}
	}
	d, err := markdown.Parse("- outer\n  - inner\n    - deepest\n")
	if err != nil {
		panic(err)
	}
	lines := layout.Lines("alpha beta gamma", layout.TextWidth("alpha b", 14), 14)
	if err = json.NewEncoder(os.Stdout).Encode(map[string]any{"glyphs": glyphs, "wrappedLines": lines, "nestedListBlocks": d.Blocks}); err != nil {
		panic(err)
	}
}
