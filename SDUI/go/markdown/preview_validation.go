package markdown

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

const resourceLimit = 4 << 20
const candidateLimit = 32 << 20
const svgNamespace = "http://www.w3.org/2000/svg"

// ValidateSVGResource admits only the supplied shape subset and returns a copy.
func ValidateSVGResource(r Resource) (Resource, error) { return validatePreviewResource(r, false) }

// ValidateMermaidResource admits the registered flowchart/graph output profile,
// including bounded text and local marker definitions, not arbitrary SVG/CSS.
// The actual composed backend must still accept the resource before publication.
func ValidateMermaidResource(r Resource) (Resource, error) { return validatePreviewResource(r, true) }

func validatePreviewResource(r Resource, mermaid bool) (Resource, error) {
	if len(r.SVG) == 0 || len(r.SVG) > resourceLimit || !utf8.Valid(r.SVG) {
		return Resource{}, fmt.Errorf("preview-resource: empty, invalid UTF-8 or exceeds 4 MiB")
	}
	if !finite(r.Width) || !finite(r.Height) || r.Width <= 0 || r.Height <= 0 || r.Width > 32768 || r.Height > 32768 {
		return Resource{}, fmt.Errorf("preview-resource: invalid supplied dimensions")
	}
	d := xml.NewDecoder(bytes.NewReader(r.SVG))
	var stack []string
	matrices := []matrix{unitMatrix}
	rootSeen, closed, decl := false, false, false
	for {
		token, e := d.Token()
		if e == io.EOF {
			break
		}
		if e != nil {
			return Resource{}, fmt.Errorf("preview-resource: %w", e)
		}
		switch t := token.(type) {
		case xml.Directive:
			return Resource{}, fmt.Errorf("preview-resource: XML directives forbidden")
		case xml.ProcInst:
			if t.Target != "xml" || decl || rootSeen || d.InputOffset() != int64(len(t.Inst)+len(t.Target)+5) {
				return Resource{}, fmt.Errorf("preview-resource: processing instruction forbidden")
			}
			// encoding/xml checks XML version/encoding; reject extra declaration fields.
			declaration := string(t.Inst)
			if !xmlDeclaration.MatchString(declaration) {
				return Resource{}, fmt.Errorf("preview-resource: invalid XML declaration")
			}
			decl = true
		case xml.StartElement:
			if closed || (len(stack) == 0 && t.Name.Local != "svg") || len(stack) > 0 && t.Name.Local == "svg" {
				return Resource{}, fmt.Errorf("preview-resource: expected one SVG document")
			}
			if t.Name.Space != "" && t.Name.Space != svgNamespace {
				return Resource{}, fmt.Errorf("preview-resource: unsupported namespace")
			}
			if len(stack) > 0 && (stack[len(stack)-1] == "title" || stack[len(stack)-1] == "desc") {
				return Resource{}, fmt.Errorf("preview-resource: markup in description")
			}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				container := parent == "svg" || parent == "g" || mermaid && (parent == "defs" || parent == "marker")
				textChild := mermaid && (parent == "text" || parent == "tspan") && t.Name.Local == "tspan"
				metadata := t.Name.Local == "title" || t.Name.Local == "desc"
				if !container && !textChild && !metadata {
					return Resource{}, fmt.Errorf("preview-resource: unsupported element nesting")
				}
			}
			attrs := map[string]string{}
			for _, a := range t.Attr {
				if a.Name.Space != "" {
					return Resource{}, fmt.Errorf("preview-resource: namespaced attribute forbidden")
				}
				if _, ok := attrs[a.Name.Local]; ok {
					return Resource{}, fmt.Errorf("preview-resource: duplicate attribute")
				}
				attrs[a.Name.Local] = a.Value
			}
			m := matrices[len(matrices)-1]
			if tr, ok := attrs["transform"]; ok {
				local, e := transforms(tr)
				if e != nil {
					return Resource{}, e
				}
				m, e = multiply(m, local)
				if e != nil {
					return Resource{}, e
				}
			}
			if e := validateElement(t.Name.Local, attrs, m, mermaid); e != nil {
				return Resource{}, fmt.Errorf("preview-resource: %w", e)
			}
			if !rootSeen {
				v, e := numbers(attrs["viewBox"])
				if e != nil || len(v) != 4 || v[2] != r.Width || v[3] != r.Height {
					return Resource{}, fmt.Errorf("preview-resource: viewBox dimensions must match supplied dimensions")
				}
				for _, key := range []string{"width", "height"} {
					if raw, ok := attrs[key]; ok {
						val, e := scalar(raw)
						want := r.Width
						if key == "height" {
							want = r.Height
						}
						if e != nil || val != want {
							return Resource{}, fmt.Errorf("preview-resource: inconsistent root dimension")
						}
					}
				}
				rootSeen = true
			}
			stack = append(stack, t.Name.Local)
			matrices = append(matrices, m)
		case xml.EndElement:
			if len(stack) == 0 {
				return Resource{}, fmt.Errorf("preview-resource: invalid close")
			}
			stack = stack[:len(stack)-1]
			matrices = matrices[:len(matrices)-1]
			if len(stack) == 0 {
				closed = true
			}
		case xml.CharData:
			if strings.TrimSpace(string(t)) != "" {
				if len(stack) == 0 {
					return Resource{}, fmt.Errorf("preview-resource: data outside root")
				}
				parent := stack[len(stack)-1]
				if parent != "title" && parent != "desc" && !(mermaid && (parent == "text" || parent == "tspan")) {
					return Resource{}, fmt.Errorf("preview-resource: unexpected text")
				}
			}
		}
	}
	if !rootSeen || !closed || len(stack) != 0 {
		return Resource{}, fmt.Errorf("preview-resource: incomplete SVG")
	}
	return copyResource(r), nil
}
