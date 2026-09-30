package sdptool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

type Model struct {
	ID      string `json:"id"`
	System  string `json:"system"`
	Source  string `json:"source"`
	Profile string `json:"profile"`
}
type Inventory struct {
	SchemaVersion      string  `json:"schemaVersion"`
	ProjectID          string  `json:"projectId"`
	ProcessProfile     string  `json:"processProfile"`
	ProjectManifest    string  `json:"projectManifest,omitempty"`
	ImplementationPlan string  `json:"implementationPlan,omitempty"`
	KanBan             string  `json:"kanban,omitempty"`
	Models             []Model `json:"models"`
	SDUI               []Model `json:"sdui"`
}
type Project struct {
	Schema       string            `json:"schema"`
	Operation    string            `json:"operation"`
	Status       string            `json:"status"`
	Root         string            `json:"root"`
	Area         string            `json:"area"`
	Inventory    Inventory         `json:"inventory"`
	Capabilities map[string]string `json:"capabilities"`
	Installation map[string]any    `json:"installation"`
	Sources      []SourceInfo      `json:"sources"`
	Files        []Node            `json:"-"`
	Navigation   *Tree             `json:"navigation,omitempty"`
	Plans        []string          `json:"plans"`
}

var identifier = regexp.MustCompile(`^[a-z][a-z0-9-]{0,63}$`)

func boundedFile(path string) ([]byte, error) {
	return boundedFileLimit(path, 1<<20)
}
func boundedFileLimit(path string, limit int64) ([]byte, error) {
	before, e := os.Stat(path)
	if e != nil {
		return nil, e
	}
	if !before.Mode().IsRegular() {
		return nil, fmt.Errorf("expected regular file: %s", path)
	}

	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	st, e := f.Stat()
	if e != nil {
		return nil, e
	}
	if !st.Mode().IsRegular() {
		return nil, fmt.Errorf("expected regular file: %s", path)
	}
	b, e := io.ReadAll(io.LimitReader(f, limit+1))
	if e == nil && int64(len(b)) > limit {
		e = fmt.Errorf("metadata exceeds %d bytes", limit)
	}
	return b, e
}
func strictJSON(b []byte, v any) error {
	// encoding/json normally accepts duplicate keys; configuration must not be ambiguous.
	d := json.NewDecoder(bytes.NewReader(b))
	var walk func() error
	walk = func() error {
		t, e := d.Token()
		if e != nil {
			return e
		}
		if delim, ok := t.(json.Delim); ok {
			if delim == '{' {
				keys := map[string]bool{}
				for d.More() {
					k, e := d.Token()
					if e != nil {
						return e
					}
					s := k.(string)
					if keys[s] {
						return fmt.Errorf("duplicate key %s", s)
					}
					keys[s] = true
					if e = walk(); e != nil {
						return e
					}
				}
			} else if delim == '[' {
				for d.More() {
					if e = walk(); e != nil {
						return e
					}
				}
			} else {
				return fmt.Errorf("unexpected delimiter")
			}
			_, e = d.Token()
			return e
		}
		return nil
	}
	if e := walk(); e != nil {
		return e
	}
	if _, e := d.Token(); e != io.EOF {
		return fmt.Errorf("trailing JSON")
	}
	d = json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	return d.Decode(v)
}
func resolvePath(root, p string) (string, error) {
	if p == "" || filepath.IsAbs(p) || filepath.ToSlash(filepath.Clean(p)) != p || strings.ContainsAny(p, "\\:") || strings.IndexFunc(p, unicode.IsControl) >= 0 {
		return "", fmt.Errorf("invalid relative path %q", p)
	}
	for _, part := range strings.Split(p, "/") {
		if part == ".." || part == "." {
			return "", fmt.Errorf("escaping path %q", p)
		}
	}
	r, e := physical(root)
	if e != nil {
		return "", e
	}
	q, e := physical(filepath.Join(r, p))
	if e != nil {
		return "", e
	}
	rel, e := filepath.Rel(r, q)
	if e != nil || rel == ".." || strings.HasPrefix(rel, "../") {
		return "", fmt.Errorf("path escapes project: %s", p)
	}
	return q, nil
}
func readYAML(path string, schemas ...string) (map[string]any, error) {
	b, e := boundedFile(path)
	if e != nil {
		return nil, e
	}
	d := yaml.NewDecoder(bytes.NewReader(b))
	var v map[string]any
	if e = d.Decode(&v); e != nil {
		return nil, e
	}
	var extra any
	if e = d.Decode(&extra); e != io.EOF {
		return nil, fmt.Errorf("expected one YAML document")
	}
	if v == nil {
		return nil, fmt.Errorf("expected YAML mapping")
	}
	if len(schemas) == 0 {
		schemas = []string{"1.0"}
	}
	supported := false
	for _, schema := range schemas {
		if v["schemaVersion"] == schema {
			supported = true
		}
	}
	if !supported {
		return nil, failure("unsupported", fmt.Errorf("unsupported manifest schema in %s", path))
	}
	return v, nil
}
func (p Project) model(id string, sdui bool) (Model, string, error) {
	list := p.Inventory.Models
	if sdui {
		list = p.Inventory.SDUI
	}
	if id == "" && len(list) == 1 {
		id = list[0].ID
	}
	for _, m := range list {
		if m.ID == id {
			expected := "design-core/0.5"
			if sdui {
				expected = "sdui/0.2"
			}
			if m.Profile != expected && (sdui || m.Profile != "design-core/0.6") {
				return m, "", failure("unsupported", fmt.Errorf("profile %s", m.Profile))
			}
			path, e := resolvePath(p.Root, m.Source)
			return m, path, e
		}
	}
	return Model{}, "", failure("selection", fmt.Errorf("select a discovered model ID (required when several exist)"))
}
