// Package model manages bounded local model artifacts, independently of Git.
package model

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

const Schema = "sdp-model/0.1"
const maxBytes = 128 << 20
const maxMetadata = 16 << 20
const maxFiles = 10000
const operations = ".model-operations"

type Files map[string][]byte
type Record struct {
	ID        string            `yaml:"id" json:"id"`
	Parents   []string          `yaml:"parents" json:"parents"`
	Kind      string            `yaml:"kind" json:"kind"`
	Message   string            `yaml:"message" json:"message"`
	At        string            `yaml:"recordedAt" json:"recordedAt"`
	Digest    string            `yaml:"sourceDigest" json:"sourceDigest"`
	Inventory map[string]string `yaml:"inventory" json:"inventory"`
	Deleted   []string          `yaml:"deleted" json:"deleted"`
	Payload   string            `yaml:"payload,omitempty" json:"payload,omitempty"`
}
type Conflict struct {
	Path   string  `yaml:"path" json:"path"`
	Base   *string `yaml:"base" json:"base"`
	Ours   *string `yaml:"ours" json:"ours"`
	Theirs *string `yaml:"theirs" json:"theirs"`
}
type Artifact struct {
	Schema         string     `yaml:"schema" json:"schema"`
	ID             string     `yaml:"id" json:"id"`
	Kind           string     `yaml:"kind" json:"kind"`
	Name           string     `yaml:"name" json:"name"`
	Head           string     `yaml:"head" json:"head"`
	Sequence       int        `yaml:"sequence" json:"sequence"`
	BaseRelease    string     `yaml:"baseRelease,omitempty" json:"baseRelease,omitempty"`
	Predecessor    string     `yaml:"predecessor,omitempty" json:"predecessor,omitempty"`
	Digest         string     `yaml:"sourceDigest" json:"sourceDigest"`
	MetadataDigest string     `yaml:"metadataDigest" json:"metadataDigest"`
	Ledger         []Record   `yaml:"ledger" json:"ledger"`
	Conflicts      []Conflict `yaml:"conflicts,omitempty" json:"conflicts,omitempty"`
	PendingParents []string   `yaml:"pendingParents,omitempty" json:"pendingParents,omitempty"`
	Validation     []string   `yaml:"validationTargets,omitempty" json:"validationTargets,omitempty"`
	Evidence       string     `yaml:"evidence,omitempty" json:"evidence,omitempty"`
}
type Result struct {
	Schema      string    `json:"schema"`
	Operation   string    `json:"operation"`
	Status      string    `json:"status"`
	Path        string    `json:"path,omitempty"`
	Artifact    *Artifact `json:"artifact,omitempty"`
	LiveDigest  string    `json:"liveDigest,omitempty"`
	Preliminary bool      `json:"preliminary,omitempty"`
	Recovery    string    `json:"recovery,omitempty"`
}
type Error struct {
	Code   string
	Detail string
}

func (e *Error) Error() string       { return e.Code + ": " + e.Detail }
func fail(code, detail string) error { return &Error{code, detail} }
func uuid() string {
	b := make([]byte, 16)
	if _, e := rand.Read(b); e != nil {
		panic(e)
	}
	b[6] = (b[6] & 15) | 64
	b[8] = (b[8] & 63) | 128
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

var uuidRE = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
var hashRE = regexp.MustCompile(`^[0-9a-f]{64}$`)
var nameRE = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,79}$`)
var versionRE = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)

func validName(s string) bool {
	if !nameRE.MatchString(s) {
		return false
	}
	u := strings.ToUpper(s)
	for _, n := range []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"} {
		if u == n {
			return false
		}
	}
	return true
}
func safePath(s string) bool {
	if s == "" || filepath.IsAbs(s) || strings.ContainsAny(s, "\\:\x00") || filepath.ToSlash(filepath.Clean(s)) != s {
		return false
	}
	for _, part := range strings.Split(s, "/") {
		if part == ".." || part == "." || strings.HasSuffix(part, ".") || strings.HasSuffix(part, " ") {
			return false
		}
		for _, c := range part {
			if c < 32 || c > 126 || strings.ContainsRune(`<>"|?*`, c) {
				return false
			}
		}
		base := strings.Split(part, ".")[0]
		if regexp.MustCompile(`(?i)^(CON|PRN|AUX|NUL|COM[1-9]|LPT[1-9])$`).MatchString(base) {
			return false
		}
	}
	return true
}
func hash(b []byte) string { v := sha256.Sum256(b); return hex.EncodeToString(v[:]) }
func inventory(f Files) map[string]string {
	m := map[string]string{}
	for p, b := range f {
		m[p] = hash(b)
	}
	return m
}
func Digest(f Files) string {
	h := sha256.New()
	h.Write([]byte("SDP-model-sources/1\x00"))
	keys := make([]string, 0, len(f))
	for p := range f {
		keys = append(keys, p)
	}
	sort.Strings(keys)
	for _, p := range keys {
		var n [8]byte
		binary.BigEndian.PutUint64(n[:], uint64(len(p)))
		h.Write(n[:])
		h.Write([]byte(p))
		b := sha256.Sum256(f[p])
		h.Write(b[:])
	}
	return hex.EncodeToString(h.Sum(nil))
}
func metadataHash(a Artifact) string {
	a.MetadataDigest = ""
	b, _ := yaml.Marshal(a)
	var v any
	_ = yaml.Unmarshal(b, &v)
	b, _ = json.Marshal(v)
	return hash(b)
}
func strict(data []byte, v any) error {
	if len(data) > maxMetadata {
		return fail("limit", "metadata too large")
	}
	var node yaml.Node
	dec := yaml.NewDecoder(bytes.NewReader(data))
	if e := dec.Decode(&node); e != nil {
		return e
	}
	var extra yaml.Node
	if e := dec.Decode(&extra); e != io.EOF {
		return fail("schema", "one YAML document required")
	}
	var walk func(*yaml.Node) error
	walk = func(n *yaml.Node) error {
		if n.Kind == yaml.AliasNode || n.Anchor != "" {
			return fail("schema", "YAML aliases/anchors forbidden")
		}
		switch n.Tag {
		case "", "!!map", "!!seq", "!!str", "!!int", "!!bool", "!!null":
		default:
			return fail("schema", "unsupported YAML tag")
		}
		for _, c := range n.Content {
			if e := walk(c); e != nil {
				return e
			}
		}
		return nil
	}
	if e := walk(&node); e != nil {
		return e
	}
	d := yaml.NewDecoder(bytes.NewReader(data))
	d.KnownFields(true)
	return d.Decode(v)
}
func recordEqual(a, b Record) bool {
	x, _ := yaml.Marshal(a)
	y, _ := yaml.Marshal(b)
	return bytes.Equal(x, y)
}
func validate(a Artifact) error {
	if a.Schema != Schema || !uuidRE.MatchString(a.ID) || !hashRE.MatchString(a.Digest) || a.MetadataDigest != metadataHash(a) {
		return fail("integrity", "schema/identity/hash mismatch")
	}
	if a.Kind != "work" && a.Kind != "proposal" && a.Kind != "candidate" && a.Kind != "release" {
		return fail("schema", "unknown artifact kind")
	}
	if (a.Kind == "release" && !versionRE.MatchString(a.Name)) || (a.Kind != "release" && !validName(a.Name)) {
		return fail("schema", "invalid artifact name")
	}
	if len(a.Ledger) == 0 || len(a.Ledger) > maxFiles {
		return fail("limit", "invalid ledger size")
	}
	records := map[string]Record{}
	for _, r := range a.Ledger {
		if _, ok := records[r.ID]; ok {
			return fail("integrity", "duplicate ledger ID")
		}
		if !hashRE.MatchString(r.Digest) || r.ID == "" {
			return fail("integrity", "invalid commit")
		}
		if _, e := time.Parse(time.RFC3339Nano, r.At); e != nil {
			return fail("schema", "invalid time")
		}
		for p, h := range r.Inventory {
			if !safePath(p) || !hashRE.MatchString(h) {
				return fail("schema", "invalid inventory")
			}
		}
		for _, p := range r.Deleted {
			if !safePath(p) {
				return fail("schema", "invalid deleted path")
			}
		}
		if r.Payload != "" && (!safePath(r.Payload) || !strings.HasPrefix(r.Payload, ".commits/")) {
			return fail("schema", "invalid payload path")
		}
		if a.Kind != "work" && r.Payload != "" {
			return fail("schema", "frozen payload reference")
		}
		records[r.ID] = r
	}
	state := map[string]int{}
	var visit func(string, int) error
	visit = func(id string, depth int) error {
		if depth > 256 {
			return fail("limit", "lineage depth exceeds 256")
		}
		r, ok := records[id]
		if !ok {
			return fail("integrity", "missing parent")
		}
		if state[id] == 1 {
			return fail("integrity", "lineage cycle")
		}
		if state[id] == 2 {
			return nil
		}
		state[id] = 1
		for _, p := range r.Parents {
			if e := visit(p, depth+1); e != nil {
				return e
			}
		}
		state[id] = 2
		return nil
	}
	for id := range records {
		if e := visit(id, 0); e != nil {
			return e
		}
	}
	r, ok := records[a.Head]
	if !ok || r.Digest != a.Digest {
		return fail("integrity", "head mismatch")
	}
	return nil
}
func metadataName(dir string) string { return filepath.Base(dir) + ".yaml" }
func readArtifact(dir string) (Artifact, error) {
	var a Artifact
	st, e := os.Lstat(dir)
	if e != nil {
		return a, e
	}
	if !st.IsDir() || st.Mode()&os.ModeSymlink != 0 {
		return a, fail("path", "artifact must be a real directory")
	}
	b, e := os.ReadFile(filepath.Join(dir, metadataName(dir)))
	if e != nil {
		return a, e
	}
	if e = strict(b, &a); e != nil {
		return a, e
	}
	return a, validate(a)
}
func saveArtifact(dir string, a Artifact) error {
	a.MetadataDigest = metadataHash(a)
	if e := validate(a); e != nil {
		return e
	}
	b, e := yaml.Marshal(a)
	if e != nil {
		return e
	}
	return os.WriteFile(filepath.Join(dir, metadataName(dir)), b, 0600)
}
func scan(dir string, exclude bool) (Files, error) {
	f := Files{}
	seen := map[string]bool{}
	total := 0
	e := filepath.WalkDir(dir, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if p == dir {
			return nil
		}
		rel, _ := filepath.Rel(dir, p)
		rel = filepath.ToSlash(rel)
		if !safePath(rel) {
			return fail("path", "unsafe path: "+rel)
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fail("path", "symlink: "+rel)
		}
		if exclude && (rel == ".commits" || rel == ".merge") {
			if !d.IsDir() {
				return fail("path", "history must be directory")
			}
			return filepath.SkipDir
		}
		if exclude && rel == metadataName(dir) {
			return nil
		}
		key := strings.ToLower(rel)
		if seen[key] {
			return fail("path", "case collision")
		}
		seen[key] = true
		if d.IsDir() {
			return nil
		}
		info, e := d.Info()
		if e != nil {
			return e
		}
		if !info.Mode().IsRegular() {
			return fail("path", "nonregular file")
		}
		if info.Size() > maxBytes {
			return fail("limit", "file too large")
		}
		b, e := os.ReadFile(p)
		if e != nil {
			return e
		}
		total += len(b)
		if total > maxBytes || len(f) >= maxFiles {
			return fail("limit", "snapshot too large")
		}
		f[rel] = b
		return nil
	})
	return f, e
}
func writeFiles(dir string, f Files) error {
	for p, b := range f {
		if !safePath(p) {
			return fail("path", p)
		}
		dest := filepath.Join(dir, filepath.FromSlash(p))
		if e := os.MkdirAll(filepath.Dir(dest), 0700); e != nil {
			return e
		}
		if e := os.WriteFile(dest, b, 0600); e != nil {
			return e
		}
	}
	return nil
}
func capture(dir string) (Artifact, Files, error) {
	a, e := readArtifact(dir)
	if e != nil {
		return a, nil, e
	}
	f, e := scan(dir, true)
	if e != nil {
		return a, nil, e
	}
	g, e := scan(dir, true)
	if e != nil {
		return a, nil, e
	}
	again, e := readArtifact(dir)
	if e != nil {
		return a, nil, e
	}
	if Digest(f) != Digest(g) || a.MetadataDigest != again.MetadataDigest {
		return a, nil, fail("stale", "source changed during capture")
	}
	if a.Kind != "work" && Digest(f) != a.Digest {
		return a, nil, fail("integrity", "frozen content changed")
	}
	return a, f, nil
}
func resolve(area, ref string) (string, error) {
	parts := strings.SplitN(ref, ":", 2)
	if len(parts) != 2 {
		return "", fail("arguments", "expected kind:NAME")
	}
	entries, e := os.ReadDir(area)
	if e != nil {
		return "", e
	}
	matches := []string{}
	for _, d := range entries {
		if !d.IsDir() || !strings.HasPrefix(d.Name(), strings.ToUpper(parts[0])+"--") {
			continue
		}
		p := filepath.Join(area, d.Name())
		a, e := readArtifact(p)
		if e != nil {
			return "", e
		}
		if a.Kind == parts[0] && strings.EqualFold(a.Name, parts[1]) {
			matches = append(matches, p)
		}
	}
	if len(matches) != 1 {
		return "", fail("reference", fmt.Sprintf("%s resolves to %d artifacts", ref, len(matches)))
	}
	return matches[0], nil
}
func result(op, path string, a Artifact, f Files) Result {
	state := "clean"
	if a.Kind == "work" && Digest(f) != a.Digest {
		state = "dirty"
	}
	if len(a.Conflicts) > 0 {
		state = "conflicted"
	}
	return Result{Schema: Schema, Operation: op, Status: state, Path: path, Artifact: &a, LiveDigest: Digest(f), Preliminary: a.Kind == "work"}
}
func Status(area, ref string) (Result, error) {
	p, e := resolve(area, ref)
	if e != nil {
		return Result{}, e
	}
	a, f, e := capture(p)
	if e != nil {
		return Result{}, e
	}
	return result("status", p, a, f), nil
}
func now() string { return time.Now().UTC().Format(time.RFC3339Nano) }
func newRecord(id, kind, message string, parents []string, f Files, payload string) Record {
	return Record{ID: id, Kind: kind, Message: message, Parents: parents, At: now(), Digest: Digest(f), Inventory: inventory(f), Payload: payload}
}
func textual(b []byte) bool { return utf8.Valid(b) && !bytes.ContainsRune(b, 0) }

var _ = errors.New
