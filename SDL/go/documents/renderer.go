package documents

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Mmdr struct{ Executable, Fingerprint string }

func NewMmdr(executable string) (*Mmdr, error) {
	b, e := os.ReadFile(executable)
	if e != nil {
		return nil, e
	}
	return &Mmdr{executable, Hash(b)}, nil
}
func (m *Mmdr) Identity() string { return "mmdr:" + m.Fingerprint + ":" + SymbolProfile }

type capped struct{ bytes.Buffer }

func (b *capped) Write(p []byte) (int, error) {
	if b.Len()+len(p) > 8<<20 {
		return 0, fmt.Errorf("renderer output limit")
	}
	return b.Buffer.Write(p)
}
func (m *Mmdr) Render(ctx context.Context, source string) ([]byte, error) {
	if len(source) > 64000 {
		return nil, fmt.Errorf("diagram source limit")
	}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, m.Executable, "--input", "-")
	cmd.Stdin = strings.NewReader(source)
	var out, errs capped
	cmd.Stdout = &out
	cmd.Stderr = &errs
	if e := cmd.Run(); e != nil {
		return nil, fmt.Errorf("renderer: %w: %.1000s", e, errs.String())
	}
	if e := ValidateSVG(out.Bytes()); e != nil {
		return nil, e
	}
	return out.Bytes(), nil
}
func ValidateSVG(b []byte) error {
	if len(b) > 8<<20 {
		return fmt.Errorf("SVG limit")
	}
	d := xml.NewDecoder(bytes.NewReader(b))
	roots, depth := 0, 0
	for {
		t, e := d.Token()
		if e == io.EOF {
			break
		}
		if e != nil {
			return e
		}
		switch x := t.(type) {
		case xml.Directive:
			return fmt.Errorf("XML directives unsupported")
		case xml.StartElement:
			if depth == 0 {
				roots++
				if x.Name.Local != "svg" || x.Name.Space != "http://www.w3.org/2000/svg" {
					return fmt.Errorf("missing SVG root")
				}
			}
			depth++
			switch strings.ToLower(x.Name.Local) {
			case "script", "foreignobject", "image", "animate", "set":
				return fmt.Errorf("active SVG unsupported")
			}
			for _, a := range x.Attr {
				key := strings.ToLower(a.Name.Local)
				if strings.HasPrefix(key, "on") || key == "href" && !strings.HasPrefix(a.Value, "#") {
					return fmt.Errorf("active/external SVG attribute")
				}
			}
		case xml.EndElement:
			depth--
		}
	}
	if roots != 1 {
		return fmt.Errorf("invalid SVG roots")
	}
	return nil
}
