// Package install owns SDP installation policy, independently of navigation.
package install

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const (
	Protocol       = "sdp-install-command/1"
	ReleaseSchema  = "sdp-release-descriptor/1"
	AdoptionSchema = "sdp-adoption/1"
	PlanSchema     = "sdp-install-plan/1"
	JournalSchema  = "sdp-install-journal/1"
	ReceiptSchema  = "3.0"
	ReceiptPath    = "SDP/Framework/installed-toolkit.manifest.yaml"
	Operations     = "SDP/.sdp-operations"
	MetadataLimit  = 16 << 20
	RecordLimit    = 64 << 20
	FileLimit      = 64 << 20
	PayloadLimit   = 512 << 20
	EntryLimit     = 100000
)

var digestRE = regexp.MustCompile(`^[a-f0-9]{64}$`)
var identityRE = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/-]{0,127}$`)
var operationRE = regexp.MustCompile(`^install-[a-f0-9]{24}$`)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Exit    int    `json:"exit"`
}

func (e *Error) Error() string { return e.Code + ": " + e.Message }
func fail(code string, exit int, format string, args ...any) error {
	return &Error{code, fmt.Sprintf(format, args...), exit}
}
func Hash(b []byte) string { v := sha256.Sum256(b); return hex.EncodeToString(v[:]) }

type File struct {
	Path      string  `json:"path"`
	Type      string  `json:"type"`
	Ownership string  `json:"ownership"`
	SHA256    *string `json:"sha256"`
	Content   []byte  `json:"content"`
}
type Move struct {
	From string `json:"from"`
	To   string `json:"to"`
}
type Asset struct {
	Platform string `json:"platform"`
	Path     string `json:"path"`
	SHA256   string `json:"sha256"`
	Size     int64  `json:"size"`
}
type Descriptor struct {
	SchemaVersion     string   `json:"schemaVersion"`
	Release           string   `json:"release"`
	SourceCommit      string   `json:"sourceCommit"`
	Protocol          string   `json:"protocol"`
	ProcessProfile    string   `json:"processProfile"`
	ManagementProfile string   `json:"managementProfile"`
	Capabilities      []string `json:"capabilities"`
	Files             []File   `json:"files"`
	Retired           []string `json:"retired"`
	UpgradesFrom      []string `json:"upgradesFrom"`
	Binaries          []Asset  `json:"binaries"`
}
type Observation struct {
	Type   string  `json:"type"`
	SHA256 *string `json:"sha256"`
}
type Adoption struct {
	SchemaVersion          string                 `json:"schemaVersion"`
	ProjectRoot            string                 `json:"projectRoot"`
	Baseline               string                 `json:"baseline"`
	ObservedCommit         *string                `json:"observedCommit"`
	TargetDigest           string                 `json:"targetDigest"`
	Inventory              map[string]Observation `json:"inventory"`
	Moves                  []Move                 `json:"moves"`
	RefreshManaged         []string               `json:"refreshManaged"`
	AllowReferenceWarnings bool                   `json:"allowReferenceWarnings"`
}
type Input struct {
	Path       string `json:"path"`
	SHA256     string `json:"sha256"`
	Bytes      []byte `json:"bytes"`
	Provenance string `json:"provenance"`
	Signature  []byte `json:"signature"`
	KeyID      string `json:"keyId"`
}
type Receipt struct {
	SchemaVersion     string   `json:"schemaVersion"`
	Release           string   `json:"release"`
	DescriptorDigest  string   `json:"descriptorDigest"`
	SourceCommit      string   `json:"sourceCommit"`
	ProcessProfile    string   `json:"processProfile"`
	ManagementProfile string   `json:"managementProfile"`
	Capabilities      []string `json:"capabilities"`
	Provenance        string   `json:"provenance"`
	OperationID       string   `json:"operationId"`
	InstalledAt       string   `json:"installedAt"`
}
type Action struct {
	Action  string  `json:"action"`
	Path    string  `json:"path"`
	Before  *string `json:"before"`
	After   *string `json:"after"`
	Content []byte  `json:"content"`
}
type Plan struct {
	SchemaVersion string                 `json:"schemaVersion"`
	PlanDigest    string                 `json:"planDigest"`
	Operation     string                 `json:"operation"`
	ProjectRoot   string                 `json:"projectRoot"`
	Release       Input                  `json:"release"`
	Previous      *Input                 `json:"previous"`
	Adoption      *Input                 `json:"adoption"`
	OldReceipt    *Receipt               `json:"oldReceipt"`
	Baseline      string                 `json:"baseline"`
	ProjectID     string                 `json:"projectId"`
	Policy        string                 `json:"policy"`
	Snapshot      map[string]Observation `json:"snapshot"`
	Actions       []Action               `json:"actions"`
	Preserved     []string               `json:"preserved"`
	Conflicts     []string               `json:"conflicts"`
	Warnings      []string               `json:"warnings"`
	CanApply      bool                   `json:"canApply"`
	NoChange      bool                   `json:"noChange"`
}
type Journal struct {
	SchemaVersion string   `json:"schemaVersion"`
	OperationID   string   `json:"operationId"`
	Plan          Plan     `json:"plan"`
	CreatedAt     string   `json:"createdAt"`
	MaintenanceID string   `json:"maintenanceId"`
	Steps         []Action `json:"steps"`
	Next          int      `json:"next"`
	Status        string   `json:"status"`
	Integrity     string   `json:"integrity"`
}

// Canonical uses sorted object keys, integer numbers, UTF-8 and one final LF.
func Canonical(v any) ([]byte, error) {
	b, e := json.Marshal(v)
	if e != nil {
		return nil, e
	}
	var x any
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	if e = d.Decode(&x); e != nil {
		return nil, e
	}
	var out bytes.Buffer
	enc := json.NewEncoder(&out)
	enc.SetEscapeHTML(false)
	e = enc.Encode(x)
	return out.Bytes(), e
}
func planHash(p Plan) string {
	b, _ := Canonical(p)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	delete(m, "planDigest")
	b, _ = Canonical(m)
	return Hash(b)
}

// Decode accepts a strict JSON/YAML subset. Every declared field is required;
// explicit null is allowed only for nullable fields, never as an omitted value.
func Decode(b []byte, limit int, v any) error {
	if len(b) > limit {
		return fail("limit", 2, "record exceeds %d bytes", limit)
	}
	var n yaml.Node
	d := yaml.NewDecoder(bytes.NewReader(b))
	if e := d.Decode(&n); e != nil {
		return fail("schema", 2, "%v", e)
	}
	var extra yaml.Node
	if e := d.Decode(&extra); e != io.EOF {
		return fail("schema", 2, "expected one document")
	}
	budget := 1000000
	var convert func(*yaml.Node, int) (any, error)
	convert = func(n *yaml.Node, depth int) (any, error) {
		budget--
		if depth > 64 || budget < 0 {
			return nil, fmt.Errorf("record nesting/entry limit")
		}
		if n.Anchor != "" || n.Kind == yaml.AliasNode || n.Style&yaml.TaggedStyle != 0 {
			return nil, fmt.Errorf("aliases, anchors and explicit tags forbidden")
		}
		switch n.Kind {
		case yaml.DocumentNode:
			if len(n.Content) != 1 {
				return nil, fmt.Errorf("empty document")
			}
			return convert(n.Content[0], depth+1)
		case yaml.MappingNode:
			m := map[string]any{}
			seen := map[string]bool{}
			for i := 0; i < len(n.Content); i += 2 {
				k := n.Content[i]
				if k.Tag != "!!str" || k.Kind != yaml.ScalarNode {
					return nil, fmt.Errorf("string mapping key required")
				}
				fold := strings.ToLower(k.Value)
				if seen[fold] {
					return nil, fmt.Errorf("duplicate/case-colliding key %s", k.Value)
				}
				seen[fold] = true
				x, e := convert(n.Content[i+1], depth+1)
				if e != nil {
					return nil, e
				}
				m[k.Value] = x
			}
			return m, nil
		case yaml.SequenceNode:
			a := []any{}
			for _, c := range n.Content {
				x, e := convert(c, depth+1)
				if e != nil {
					return nil, e
				}
				a = append(a, x)
			}
			return a, nil
		case yaml.ScalarNode:
			switch n.Tag {
			case "!!str":
				return n.Value, nil
			case "!!null":
				return nil, nil
			case "!!bool":
				return strconv.ParseBool(n.Value)
			case "!!int":
				x, e := strconv.ParseInt(n.Value, 10, 64)
				return x, e
			}
			return nil, fmt.Errorf("unsupported scalar %s", n.Tag)
		}
		return nil, fmt.Errorf("unsupported node")
	}
	x, e := convert(&n, 0)
	if e != nil {
		return fail("schema", 2, "%v", e)
	}
	if e = shape(x, reflect.TypeOf(v).Elem()); e != nil {
		return fail("schema", 2, "%v", e)
	}
	b, e = json.Marshal(x)
	if e != nil {
		return e
	}
	if e = json.Unmarshal(b, v); e != nil {
		return fail("schema", 2, "%v", e)
	}
	return nil
}
func shape(v any, t reflect.Type) error {
	if t.Kind() == reflect.Pointer {
		if v == nil {
			return nil
		}
		return shape(v, t.Elem())
	}
	if v == nil {
		if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
			return nil
		}
		return fmt.Errorf("null for %s", t)
	}
	switch t.Kind() {
	case reflect.Struct:
		m, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("expected object")
		}
		if len(m) != t.NumField() {
			return fmt.Errorf("unknown or missing fields in %s", t)
		}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			k := f.Tag.Get("json")
			x, ok := m[k]
			if !ok {
				return fmt.Errorf("missing %s", k)
			}
			if e := shape(x, f.Type); e != nil {
				return fmt.Errorf("%s: %w", k, e)
			}
		}
	case reflect.Map:
		m, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("expected map")
		}
		for _, x := range m {
			if e := shape(x, t.Elem()); e != nil {
				return e
			}
		}
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			if _, ok := v.(string); !ok {
				return fmt.Errorf("expected base64")
			}
			return nil
		}
		a, ok := v.([]any)
		if !ok {
			return fmt.Errorf("expected array")
		}
		for _, x := range a {
			if e := shape(x, t.Elem()); e != nil {
				return e
			}
		}
	case reflect.String:
		if _, ok := v.(string); !ok {
			return fmt.Errorf("expected string")
		}
	case reflect.Bool:
		if _, ok := v.(bool); !ok {
			return fmt.Errorf("expected boolean")
		}
	case reflect.Int, reflect.Int64:
		if _, ok := v.(int64); !ok {
			return fmt.Errorf("expected integer")
		}
	default:
		return fmt.Errorf("unsupported schema type")
	}
	return nil
}
func Read(path string, limit int) ([]byte, error) {
	if e := SafeAbsolute(path); e != nil {
		return nil, e
	}
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	s, e := f.Stat()
	if e != nil {
		return nil, e
	}
	if !s.Mode().IsRegular() || s.Size() > int64(limit) {
		return nil, fail("limit", 2, "not a bounded regular file: %s", path)
	}
	b, e := io.ReadAll(io.LimitReader(f, int64(limit)+1))
	if len(b) > limit {
		return nil, fail("limit", 2, "file grew beyond limit")
	}
	return b, e
}
func SafeAbsolute(p string) error {
	p, e := filepath.Abs(p)
	if e != nil {
		return e
	}
	for q := p; ; q = filepath.Dir(q) {
		s, e := os.Lstat(q)
		if e == nil && s.Mode()&os.ModeSymlink != 0 {
			return fail("path", 2, "symlink: %s", q)
		}
		if e != nil && !os.IsNotExist(e) {
			return e
		}
		if filepath.Dir(q) == q {
			break
		}
	}
	return nil
}
func Root(selected string) (string, error) {
	p, e := filepath.Abs(selected)
	if e != nil {
		return "", e
	}
	p = filepath.Clean(p)
	if e = SafeAbsolute(p); e != nil {
		return "", e
	}
	s, e := os.Stat(p)
	if e != nil || !s.IsDir() {
		return "", fail("root", 2, "existing project directory required")
	}
	if filepath.Base(p) == "SDP" {
		p = filepath.Dir(p)
	}
	return p, nil
}
func Relative(p string) error {
	if p == "" || p == "." || path.Clean(p) != p || strings.HasPrefix(p, "/") || strings.ContainsAny(p, "\\:\x00") || strings.IndexFunc(p, unicode.IsControl) >= 0 {
		return fail("path", 2, "invalid relative path %q", p)
	}
	for _, s := range strings.Split(p, "/") {
		base := strings.ToUpper(strings.Split(s, ".")[0])
		if s == ".." || strings.HasSuffix(s, ".") || strings.HasSuffix(s, " ") || strings.ContainsAny(s, `<>"|?*`) || regexp.MustCompile(`^(CON|PRN|AUX|NUL|COM[0-9]|LPT[0-9])$`).MatchString(base) {
			return fail("path", 2, "nonportable path %q", p)
		}
	}
	return nil
}
func Reserved(p string) bool {
	for _, v := range []string{Operations, "SDP/.sdp-backups", ReceiptPath, "SDP/ProjectManagement/Ledger.ndjson", "SDP/navigation.json", "SDP/Traceability/Ledger.ndjson", "SDP/KanBan/board.json"} {
		if overlap(strings.ToLower(p), strings.ToLower(v)) {
			return true
		}
	}
	return false
}
func overlap(a, b string) bool {
	return a == b || strings.HasPrefix(a, b+"/") || strings.HasPrefix(b, a+"/")
}
func managed(p string) bool {
	return p == "AGENTS.md" || strings.HasPrefix(p, ".codex/skills/") || strings.HasPrefix(p, "SDP/Framework/")
}
func ValidateDescriptor(d Descriptor) error {
	if d.SchemaVersion != ReleaseSchema || d.Protocol != Protocol || !identityRE.MatchString(d.Release) || !regexp.MustCompile(`^[a-f0-9]{40}$`).MatchString(d.SourceCommit) || d.ProcessProfile != "sdp-five-phase/0.1" || d.ManagementProfile != "sdp-project-management/0.2" {
		return fail("unsupported", 2, "release identity/protocol/profile")
	}
	if len(d.Files) > EntryLimit {
		return fail("limit", 2, "too many entries")
	}
	seen := map[string]string{}
	total := 0
	for _, f := range d.Files {
		if e := Relative(f.Path); e != nil {
			return e
		}
		if Reserved(f.Path) || !(managed(f.Path) || strings.HasPrefix(f.Path, "SDP/")) {
			return fail("scope", 2, "invalid release destination %s", f.Path)
		}
		if f.Ownership != "managed" && f.Ownership != "initialize-if-missing" || f.Ownership == "managed" && !managed(f.Path) {
			return fail("ownership", 2, "%s", f.Path)
		}
		for q := f.Path; q != "."; q = path.Dir(q) {
			fold := strings.ToLower(q)
			if old, ok := seen[fold]; ok && old != q {
				return fail("path", 2, "case collision %s", q)
			}
			seen[fold] = q
		}
		if f.Type == "directory" {
			if f.SHA256 != nil || f.Content != nil {
				return fail("schema", 2, "directory content")
			}
		} else if f.Type == "file" {
			if f.SHA256 == nil || *f.SHA256 != Hash(f.Content) || len(f.Content) > FileLimit {
				return fail("digest", 4, "payload %s", f.Path)
			}
			total += len(f.Content)
		} else {
			return fail("schema", 2, "unknown file type")
		}
	}
	if total > PayloadLimit {
		return fail("limit", 2, "payload too large")
	}
	for i, a := range d.Files {
		for _, b := range d.Files[i+1:] {
			if strings.EqualFold(a.Path, b.Path) || a.Type == "file" && strings.HasPrefix(b.Path, a.Path+"/") || b.Type == "file" && strings.HasPrefix(a.Path, b.Path+"/") {
				return fail("path", 2, "duplicate or overlapping files")
			}
		}
	}
	for _, p := range d.Retired {
		if e := Relative(p); e != nil {
			return e
		}
		if !managed(p) || Reserved(p) {
			return fail("scope", 2, "invalid retirement")
		}
	}
	for _, a := range d.Binaries {
		if e := Relative(a.Path); e != nil {
			return e
		}
		if !digestRE.MatchString(a.SHA256) || a.Size < 1 || a.Size > PayloadLimit || !identityRE.MatchString(a.Platform) {
			return fail("schema", 2, "invalid binary")
		}
	}
	return nil
}
func ValidateReceipt(r Receipt) error {
	if r.SchemaVersion != ReceiptSchema || !identityRE.MatchString(r.Release) || !digestRE.MatchString(r.DescriptorDigest) || !operationRE.MatchString(r.OperationID) || r.ProcessProfile != "sdp-five-phase/0.1" || r.ManagementProfile != "sdp-project-management/0.2" || r.InstalledAt == "" || (r.Provenance != "local-development" && r.Provenance != "signed") {
		return fail("receipt", 2, "invalid installed receipt")
	}
	return nil
}
