package previews

import (
	"encoding/xml"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
)

// Resource/condition validation only: no app, host, window or rasterizer.
func TestOriginResourceSlots(t *testing.T) {
	cases := []struct {
		slot, viewBox, transform string
		x, y, w, h               float64
	}{
		{"negative-origin", "-40 -20 400 200", "", -40, -20, 400, 200},
		{"positive-origin", "40 20 400 200", "", 40, 20, 400, 200},
		{"mixed-origin", "-40 20 400 200", "", -40, 20, 400, 200},
		{"root-transform", "40 20 400 200", "scale(0.5)", 120, 60, 160, 80},
	}
	for _, tc := range cases {
		t.Run(tc.slot, func(t *testing.T) {
			r, err := supplied(tc.slot)
			if err != nil {
				t.Fatal(err)
			}
			if r.Resource.Width != 400 || r.Resource.Height != 200 || r.SHA256 != digest(r.Resource.SVG) {
				t.Fatal("incorrect dimensions/digest", r)
			}
			if r.Source.Module != "art" || r.Source.Object != "Chart" || r.Source.Member != "resource" || r.ProviderID != "fixture-svg/1" {
				t.Fatal("binding changed", r)
			}
			if _, err = markdown.ValidateSVGResource(r.Resource); err != nil {
				t.Fatal("closed subset rejected", err)
			}
			var doc struct {
				XMLName   xml.Name
				ViewBox   string `xml:"viewBox,attr"`
				Transform string `xml:"transform,attr"`
				Rects     []struct {
					X    float64 `xml:"x,attr"`
					Y    float64 `xml:"y,attr"`
					W    float64 `xml:"width,attr"`
					H    float64 `xml:"height,attr"`
					Fill string  `xml:"fill,attr"`
				} `xml:"rect"`
			}
			if err = xml.Unmarshal(r.Resource.SVG, &doc); err != nil {
				t.Fatal(err)
			}
			if doc.XMLName.Local != "svg" || doc.ViewBox != tc.viewBox || doc.Transform != tc.transform || len(doc.Rects) != 1 {
				t.Fatalf("bounds memo/source mismatch: %+v", doc)
			}
			rect := doc.Rects[0]
			if rect.X != tc.x || rect.Y != tc.y || rect.W != tc.w || rect.H != tc.h || rect.Fill != "#238248" {
				t.Fatalf("pixel oracle source mismatch: %+v", rect)
			}
			// Exercise the existing stdin condition path; no SDL engine or host is needed.
			f := &Fixture{resources: map[string]string{}}
			c := &Controller{Fixture: f}
			for _, alias := range []string{"Hero", "ReportSVG", "DetailSVG"} {
				if err = c.Run("resource " + alias + " " + tc.slot); err != nil {
					t.Fatal(alias, err)
				}
				if f.resources[alias] != tc.slot {
					t.Fatal("condition did not stage slot", alias)
				}
			}
			if f.revision != 3 {
				t.Fatal("condition revisions", f.revision)
			}
		})
	}
}
