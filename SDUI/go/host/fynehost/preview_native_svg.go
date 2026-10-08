package fynehost

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
	"github.com/fyne-io/oksvg"
)

// nativePreviewSVG receives an already validated provider resource. Fyne's
// pinned SVG decoder scales nonzero viewBox origins incorrectly. Derive a
// zero-origin native representation during preparation, retaining the provider
// bytes and digest as the source identity. SVG 2 places the root transform
// outside the viewBox mapping: displayScale * rootTransform * originTranslation.
func nativePreviewSVG(r markdown.Resource) ([]byte, error) {
	source, changed, err := nativeUniformScale(r.SVG)
	if err != nil {
		return nil, err
	}
	icon, err := oksvg.ReadIconStream(bytes.NewReader(source), oksvg.StrictErrorMode)
	if err != nil {
		return nil, err
	}
	if icon.ViewBox.X == 0 && icon.ViewBox.Y == 0 {
		if !changed {
			return r.SVG, nil
		}
		return validateNativePreviewSVG(source, r)
	}
	d := xml.NewDecoder(bytes.NewReader(source))
	var root xml.StartElement
	var start, end int64
	depth := 0
	for {
		before := d.InputOffset()
		token, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch token := token.(type) {
		case xml.StartElement:
			if depth == 0 {
				root, start = token, d.InputOffset()
			}
			depth++
		case xml.EndElement:
			depth--
			if depth == 0 {
				end = before
			}
		}
	}
	// A self-closing root has no content (the decoder synthesizes its EndElement).
	if end < start || root.Name.Local != "svg" {
		return nil, fmt.Errorf("native preview: invalid SVG root")
	}
	var out bytes.Buffer
	number := func(v float64) string { return strconv.FormatFloat(v, 'g', -1, 64) }
	fmt.Fprintf(&out, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %s %s"><g`, number(r.Width), number(r.Height))
	for _, a := range root.Attr {
		switch a.Name.Local {
		case "xmlns", "viewBox", "width", "height":
			continue
		}
		// The provider's closed subset permits only unqualified presentation /
		// transform attributes here. Preserve all of them on the outer group.
		out.WriteByte(' ')
		out.WriteString(a.Name.Local)
		out.WriteString(`="`)
		if err := xml.EscapeText(&out, []byte(a.Value)); err != nil {
			return nil, err
		}
		out.WriteByte('"')
	}
	fmt.Fprintf(&out, `><g transform="translate(%s %s)">`, number(-icon.ViewBox.X), number(-icon.ViewBox.Y))
	out.Write(source[start:end])
	out.WriteString(`</g></g></svg>`)
	return validateNativePreviewSVG(out.Bytes(), r)
}

func validateNativePreviewSVG(derived []byte, r markdown.Resource) ([]byte, error) {
	// Recheck the actual composed transform/geometry and the actual native
	// decoder before PreparePreviews freezes its outcome or capabilities.
	if _, err := markdown.ValidateSVGResource(markdown.Resource{SVG: derived, Width: r.Width, Height: r.Height}); err != nil {
		return nil, fmt.Errorf("native preview derivative: %w", err)
	}
	if err := checkPreviewSVG(derived); err != nil {
		return nil, err
	}
	return derived, nil
}

// The pinned oksvg decoder interprets scale(s) as scale(s, 0). Expand only
// this spelling in actual transform attributes, including descendant shapes.
// Already validated source excludes other attribute namespaces / CSS transforms.
var nativeScale = regexp.MustCompile(`scale[^\(]*\([^\)]*\)`)

func nativeUniformScale(source []byte) ([]byte, bool, error) {
	d := xml.NewDecoder(bytes.NewReader(source))
	var out bytes.Buffer
	var copied int64
	changed := false
	for {
		before := d.InputOffset()
		token, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, false, err
		}
		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}
		modified := false
		for i, a := range start.Attr {
			if a.Name.Local != "transform" {
				continue
			}
			value := nativeScale.ReplaceAllStringFunc(a.Value, func(call string) string {
				open := strings.IndexByte(call, '(')
				operand := strings.TrimSpace(call[open+1 : len(call)-1])
				// Match exactly one source number, including exponent/sign and
				// the whitespace accepted by the provider's transform grammar.
				if _, err := strconv.ParseFloat(operand, 64); err != nil {
					return call
				}
				return "scale(" + operand + " " + operand + ")"
			})
			if value != a.Value {
				start.Attr[i].Value = value
				modified = true
			}
		}
		if !modified {
			continue
		}
		out.Write(source[copied:before])
		out.WriteByte('<')
		out.WriteString(start.Name.Local)
		for _, a := range start.Attr {
			out.WriteByte(' ')
			out.WriteString(a.Name.Local)
			out.WriteString(`="`)
			if err := xml.EscapeText(&out, []byte(a.Value)); err != nil {
				return nil, false, err
			}
			out.WriteByte('"')
		}
		copied = d.InputOffset()
		if bytes.HasSuffix(source[before:copied], []byte("/>")) {
			out.WriteByte('/')
		}
		out.WriteByte('>')
		changed = true
	}
	if !changed {
		return source, false, nil
	}
	out.Write(source[copied:])
	return out.Bytes(), true, nil
}
