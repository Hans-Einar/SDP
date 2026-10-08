package markdown

import (
	"testing"
)

func TestPreviewActualFiniteShapeGeometry(t *testing.T) {
	cases := []string{
		`<g transform="scale(1e308)"><line x1="0" y1="0" x2="1" y2="0" stroke="red" stroke-width="4"/></g>`,
		`<g transform="scale(1e308)"><path d="M1 0 A2 2 0 1 1 -1 0"/></g>`,
		`<g transform="matrix(1e308 0 -1e308 1 0 0)"><line x1="1" y1="-1" x2="0" y2="0"/></g>`,
		`<rect x="1e308" y="0" width="1e308" height="1"/>`,
		`<rect x="0" y="1e308" width="1" height="1e308"/>`,
		`<circle cx="1e308" cy="0" r="1e308"/>`,
		`<ellipse cx="0" cy="-1e308" rx="1" ry="1e308"/>`,
		`<g transform="matrix(1e308 0 -1e308 1 0 0)"><rect x="1" y="-1" width="1" height="1"/></g>`,
		`<path d="M0 0 Q-1e308 0 1e308 0 T0 0"/>`,
		`<path d="M0 0 C0 0 -1e308 0 1e308 0 S0 0 0 0"/>`,
		`<path d="M0 0 q-1e308 0 1e308 0 t-1e308 0"/>`,
		`<g transform="matrix(1e308 0 -1e308 1 0 0)"><path d="M0 0 Q0 0 1 0 T0 0"/></g>`,
	}
	for i, body := range cases {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			r := shapeResource(body)
			for _, validate := range []func(Resource) (Resource, error){ValidateSVGResource, ValidateMermaidResource} {
				if _, e := validate(r); e == nil {
					t.Errorf("admitted nonfinite derived geometry: %s", body)
				}
			}
			for _, policy := range []string{"label", "reject"} {
				root := previewRoot(t, svgNodeSource("a", policy))
				calls := 0
				backend := fullBackend()
				backend.SVG = func(Resource) error { calls++; return nil }
				p, e := PreparePreviews(root, bindSVG(root, r), nil, backend)
				if calls != 0 {
					t.Error("invalid geometry reached backend")
				}
				if policy == "reject" {
					if e == nil {
						t.Error("reject policy admitted invalid geometry")
					}
				} else {
					if e != nil {
						t.Fatal(e)
					}
					o, _ := p.Outcome("Main/a")
					if o.Status != "label" || o.Resource != nil {
						t.Error("invalid geometry became rendered")
					}
				}
			}
		})
	}
}
func TestPreviewFiniteGeometryPreservesValidControls(t *testing.T) {
	for _, body := range []string{
		`<g transform="matrix(2 1 -1 2 3 4)"><rect x="-3" y="-2" width="7" height="4"/><circle cx="1" cy="-2" r="3"/><ellipse cx="-1" cy="2" rx="3" ry="4"/><line x1="1" y1="-1" x2="3" y2="2"/></g>`,
		`<path d="M0 0 Q1 2 3 4 T5 6 7 8 C1 2 3 4 5 6 S7 8 9 10 11 12 13 14 z"/>`,
		`<path d="M0 0 Q-1e308 0 1e308 0 L0 0 T1 1"/>`,
		`<path d="M0 0 C0 0 -1e308 0 1e308 0 M0 0 S1 1 2 2"/>`,
		`<path d="M1e308 0 Q1e308 0 1e308 0 T1e308 0"/>`,
		`<path d="M0 0 A20 10 45 0 1 30 10 a5 8 90 1 0 10 5"/>`,
		`<path d="M0 0 A0 3 0 0 1 5 5 A2 3 0 1 1 5 5"/>`,
		`<g transform="matrix(2 1 -1 2 3 4)"><path d="M0 0 A5 10 30 0 1 20 20"/></g>`,
		// Pair validation must not turn x into a fictitious y that overflows.
		`<g transform="matrix(1 0 1e308 1 0 0)"><line x1="2" y1="0" x2="3" y2="0"/></g>`,
	} {
		if _, e := ValidateSVGResource(shapeResource(body)); e != nil {
			t.Errorf("valid geometry rejected: %s: %v", body, e)
		}
	}
}
