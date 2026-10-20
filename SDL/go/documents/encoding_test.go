package documents

import (
	"context"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestEncodingOverviewBundleAndSourcePositions(t *testing.T) {
	source := `language design-core version 0.5.
unit Owner.
dataset State.
contract Record.
datagram Wire.
contract Shape.
variant Body.
field Value.
encoding Fixed.
Owner holds State.
State upholds Record.
Record has completeness = open.
Wire from State.
Wire upholds Shape.
Shape has completeness = closed.
Shape defines Body.
Body has-field Value.
Value has value-type = unsigned.
Value has presence = required.
Fixed encodes Body.
Fixed has byte-order = little-endian.
Fixed has bit-order = least-significant-first.
Fixed places Value at 0 bits 24.
`
	m, e := parser.Parse(source)
	if e != nil {
		t.Fatal(e)
	}
	source = parser.FormatFile(m)
	path := filepath.Join(t.TempDir(), "Packet.design")
	if e = os.WriteFile(path, []byte(source), 0600); e != nil {
		t.Fatal(e)
	}
	v, _, e := viewpoint.Load(path)
	if e != nil {
		t.Fatal(e)
	}
	v.SourceBindings = map[string]string{}
	for _, d := range v.Model.Declarations {
		v.SourceBindings[d.Span.Source] = "packet-source"
	}
	b, e := Selected(context.Background(), v, viewpoint.Query{Viewpoint: "VP10", Projection: "overview"}, nil)
	if e != nil {
		t.Fatal(e)
	}
	text := string(b.Files["entry.md"])
	for _, want := range []string{"```mermaid\npacket", "Encoding register", "Variant coverage", "Datagram family: [Wire]", "0–23 / 24", "unsigned / required", "little-endian, least-significant-first; 24 bits"} {
		if !strings.Contains(text, want) {
			t.Fatal(want, text)
		}
	}
	if len(b.Manifest.Diagrams) != 1 || strings.Count(text, "```mermaid") != 1 || strings.Contains(text, "No diagrams match") {
		t.Fatal(text)
	}
	links := regexp.MustCompile(`\[([^\]]+)\]\(sdl-source://packet-source/(\d+)/(\d+)\)`).FindAllStringSubmatch(text, -1)
	if len(links) < 15 {
		t.Fatal("missing source links", text)
	}
	sawPlacement := false
	for _, link := range links {
		start, _ := strconv.Atoi(link[2])
		end, _ := strconv.Atoi(link[3])
		if start < 0 || end <= start || end > len(source) {
			t.Fatal(link)
		}
		span := source[start:end]
		if strings.Contains(span, "Fixed places Value at 0 bits 24") {
			sawPlacement = true
		}
	}
	if !sawPlacement {
		t.Fatal("placement source inaccessible")
	}
	for _, f := range b.Manifest.Facts {
		if !strings.Contains(text, v.FactReference(f.S("id"))) {
			t.Fatal("fact has no visible proof", f)
		}
	}
	if e = b.CheckLinks(); e != nil {
		t.Fatal(e)
	}
}
