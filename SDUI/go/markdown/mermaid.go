package markdown

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Mmdr is an explicit trusted program registration. Source text is stdin, never
// a shell command. The first supported capability is flowchart/graph only.
type Mmdr struct{ Executable string }

func (m Mmdr) Render(source string) (Resource, error) {
	first := strings.Fields(source)
	if len(first) == 0 || (first[0] != "flowchart" && first[0] != "graph") {
		return Resource{}, fmt.Errorf("mermaid-profile: only flowchart/graph is verified")
	}
	if len(source) > 12000 || strings.Contains(source, "%%{") {
		return Resource{}, fmt.Errorf("mermaid-profile: source limit or unregistered configuration")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, m.Executable, "--input", "-")
	cmd.Stdin = strings.NewReader(source)
	var out capped
	var errors capped
	cmd.Stdout = &out
	cmd.Stderr = &errors
	if err := cmd.Run(); err != nil {
		return Resource{}, fmt.Errorf("mermaid-render: %w: %.1000s", err, errors.String())
	}
	return validateResource(out.Bytes())
}

type capped struct{ bytes.Buffer }

func (b *capped) Write(data []byte) (int, error) {
	if b.Len()+len(data) > 4<<20 {
		return 0, fmt.Errorf("renderer output exceeds 4 MiB")
	}
	return b.Buffer.Write(data)
}
func validateResource(data []byte) (Resource, error) {
	r := Resource{SVG: data}
	d := xml.NewDecoder(bytes.NewReader(data))
	first := true
	for {
		token, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Resource{}, err
		}
		switch x := token.(type) {
		case xml.Directive:
			return Resource{}, fmt.Errorf("mermaid-resource: XML directives unsupported")
		case xml.StartElement:
			if first {
				if x.Name.Local != "svg" {
					return Resource{}, fmt.Errorf("mermaid-resource: missing SVG root")
				}
				first = false
				for _, a := range x.Attr {
					if a.Name.Local == "viewBox" {
						v := strings.Fields(a.Value)
						if len(v) == 4 {
							r.Width, _ = strconv.ParseFloat(v[2], 64)
							r.Height, _ = strconv.ParseFloat(v[3], 64)
						}
					}
				}
			}
			switch strings.ToLower(x.Name.Local) {
			case "script", "foreignobject", "image", "animate", "set":
				return Resource{}, fmt.Errorf("mermaid-resource: active or external SVG element")
			}
			for _, a := range x.Attr {
				key := strings.ToLower(a.Name.Local)
				if strings.HasPrefix(key, "on") || key == "href" && !strings.HasPrefix(a.Value, "#") || externalCSS(a.Value) {
					return Resource{}, fmt.Errorf("mermaid-resource: external or active attribute")
				}
			}
		case xml.CharData:
			t := strings.ToLower(string(x))
			if externalCSS(t) {
				return Resource{}, fmt.Errorf("mermaid-resource: external CSS")
			}
		}
	}
	if !(r.Width > 0 && r.Width <= 32768 && r.Height > 0 && r.Height <= 32768) {
		return Resource{}, fmt.Errorf("mermaid-resource: invalid dimensions")
	}
	return r, nil
}

// WriteResources exports immutable content-addressed SVGs separately from their
// embedded preview. Consumers can publish the directory with their own lease policy.
func (p *Provider) WriteResources(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	for id, r := range p.Resources {
		path := filepath.Join(dir, id+".svg")
		old, err := os.ReadFile(path)
		if err == nil {
			if !bytes.Equal(old, r.SVG) {
				return fmt.Errorf("resource collision: %s", id)
			}
			continue
		}
		if !os.IsNotExist(err) {
			return err
		}
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			return err
		}
		_, err = f.Write(r.SVG)
		closeErr := f.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
	}
	return nil
}

var cssURLs = regexp.MustCompile(`(?i)url\(\s*['"]?([^)'"\s]+)`)

func externalCSS(s string) bool {
	if strings.Contains(strings.ToLower(s), "@import") {
		return true
	}
	for _, m := range cssURLs.FindAllStringSubmatch(s, -1) {
		if !strings.HasPrefix(m[1], "#") {
			return true
		}
	}
	return false
}
