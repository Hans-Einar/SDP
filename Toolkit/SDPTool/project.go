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
type Registration struct {
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
	Registration Registration      `json:"registration"`
	Capabilities map[string]string `json:"capabilities"`
	Installation map[string]any    `json:"installation"`
}

var identifier = regexp.MustCompile(`^[a-z][a-z0-9-]{0,63}$`)

func boundedFile(path string) ([]byte, error) {
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
	b, e := io.ReadAll(io.LimitReader(f, (1<<20)+1))
	if e == nil && len(b) > 1<<20 {
		e = fmt.Errorf("metadata exceeds 1 MiB")
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
	if p == "" || filepath.IsAbs(p) || filepath.Clean(p) != p || strings.ContainsAny(p, "\\:") || strings.IndexFunc(p, unicode.IsControl) >= 0 {
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
func readYAML(path string) (map[string]any, error) {
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
	if v["schemaVersion"] != "1.0" {
		return nil, failure("unsupported", fmt.Errorf("unsupported manifest schema in %s", path))
	}
	return v, nil
}
func Discover(selected string) (Project, error) {
	p := Project{Schema: Version, Operation: "discover", Capabilities: map[string]string{}, Installation: map[string]any{"state": "unknown"}}
	bad := func(code string, e error) (Project, error) { p.Status = code; return p, failure(code, e) }
	if selected == "" {
		selected = "."
	}
	dir, e := physical(selected)
	if e != nil {
		return bad("invalid", e)
	}
	st, e := os.Stat(dir)
	if e != nil || !st.IsDir() {
		return bad("missing", fmt.Errorf("selected directory unavailable: %s", dir))
	}
	area := dir
	file := filepath.Join(area, "navigation.json")
	if _, e = os.Lstat(file); os.IsNotExist(e) {
		area = filepath.Join(dir, "SDP")
		file = filepath.Join(area, "navigation.json")
	}
	p.Area = area
	p.Root = filepath.Dir(area)
	b, e := boundedFile(file)
	if os.IsNotExist(e) {
		return bad("missing", fmt.Errorf("no navigation.json in selected area or SDP child; register navigation without inferring installation"))
	}
	if e != nil {
		return bad("invalid", e)
	}
	var raw map[string]any
	if e = json.Unmarshal(b, &raw); e != nil {
		return bad("invalid", e)
	}
	if raw["schemaVersion"] != "1.0" || raw["processProfile"] != "sdp-five-phase/0.1" {
		return bad("unsupported", fmt.Errorf("unsupported navigation schema or process profile"))
	}
	if e = strictJSON(b, &p.Registration); e != nil {
		return bad("invalid", e)
	}
	r := p.Registration
	if !identifier.MatchString(r.ProjectID) || r.Models == nil || r.SDUI == nil || len(r.Models) > 32 || len(r.SDUI) > 32 {
		return bad("invalid", fmt.Errorf("invalid project ID or source lists"))
	}
	ids := map[string]bool{}
	for _, m := range append(append([]Model{}, r.Models...), r.SDUI...) {
		if !identifier.MatchString(m.ID) || ids[m.ID] || m.System == "" || m.Profile == "" {
			return bad("invalid", fmt.Errorf("invalid/duplicate model registration %s", m.ID))
		}
		ids[m.ID] = true
		if _, e = resolvePath(p.Root, m.Source); e != nil {
			return bad("invalid", e)
		}
	}
	for _, v := range []string{r.ImplementationPlan, r.KanBan, r.ProjectManifest} {
		if v != "" {
			if _, e = resolvePath(p.Root, v); e != nil {
				return bad("invalid", e)
			}
		}
	}
	for name, value := range map[string]string{"implementation-plan": r.ImplementationPlan, "kanban": r.KanBan} {
		state := "absent"
		if value != "" {
			state = "declared"
		}
		p.Capabilities[name] = state
	}
	for name, models := range map[string][]Model{"sdl": r.Models, "sdui": r.SDUI} {
		state := "absent"
		if len(models) > 0 {
			state = "declared"
		}
		p.Capabilities[name] = state
	}
	if r.ProjectManifest != "" {
		f, _ := resolvePath(p.Root, r.ProjectManifest)
		manifest, err := readYAML(f)
		if err != nil {
			if x, ok := err.(*Failure); ok {
				return bad(x.Code, err)
			}
			return bad("invalid", err)
		}
		installed, ok := manifest["installed"].(map[string]any)
		if !ok {
			return bad("invalid", fmt.Errorf("project manifest lacks installed mapping"))
		}
		rel, ok := installed["manifestPath"].(string)
		if !ok {
			return bad("invalid", fmt.Errorf("missing installed.manifestPath"))
		}
		dest, err := resolvePath(filepath.Dir(f), rel)
		if err != nil {
			return bad("invalid", err)
		}
		facts, err := readYAML(dest)
		if err != nil {
			if x, ok := err.(*Failure); ok {
				return bad(x.Code, err)
			}
			return bad("invalid", err)
		}
		p.Installation = map[string]any{"state": "declared", "projectManifest": f, "installedManifest": dest, "facts": facts, "validation": "schema version and YAML structure only; use installer validator for full conformance"}
	}
	p.Status = "valid"
	return p, nil
}
func (p Project) model(id string, sdui bool) (Model, string, error) {
	list := p.Registration.Models
	if sdui {
		list = p.Registration.SDUI
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
			if m.Profile != expected {
				return m, "", failure("unsupported", fmt.Errorf("profile %s", m.Profile))
			}
			path, e := resolvePath(p.Root, m.Source)
			return m, path, e
		}
	}
	return Model{}, "", failure("selection", fmt.Errorf("select a registered model ID (required when several exist)"))
}
