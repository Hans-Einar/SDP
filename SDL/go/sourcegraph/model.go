// Package sourcegraph composes source-owned SDL. File syntax remains I/O-free;
// loading and publication freshness live at the filesystem boundary.
package sourcegraph

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"path"
	"sort"
	"strings"
	"unicode"
)

const Version = "sdl-source-graph/1"
const MaxFiles = 128
const MaxDiagnostics = 100
const MaxRelated = 8

type Issue struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Span    parser.Span   `json:"span"`
	Related []parser.Span `json:"related,omitempty"`
}

func (i Issue) Error() string {
	return fmt.Sprintf("%s at %s:%d:%d: %s", i.Code, i.Span.Source, i.Span.Line, i.Span.Column, i.Message)
}

type Issues []Issue

func (i Issues) Error() string {
	if len(i) == 0 {
		return ""
	}
	return i[0].Error()
}

type Source struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}
type Edge struct {
	From string      `json:"from"`
	To   string      `json:"to"`
	Span parser.Span `json:"span"`
}
type Snapshot struct {
	entry, revision, system string
	model                   *parser.Model
	files                   []*parser.File
	sources                 []Source
	edges                   []Edge
	warnings                []Issue
}

func (s *Snapshot) Model() *parser.Model { return parser.CopyModel(s.model) }
func (s *Snapshot) Revision() string     { return s.revision }
func (s *Snapshot) System() string       { return s.system }
func (s *Snapshot) Profile() string      { return "design-core/" + s.model.Header.Version }
func (s *Snapshot) Sources() []Source    { return append([]Source{}, s.sources...) }
func (s *Snapshot) Edges() []Edge        { return append([]Edge{}, s.edges...) }
func (s *Snapshot) Warnings() []Issue    { return append([]Issue{}, s.warnings...) }
func (s *Snapshot) FileCount() int       { return len(s.files) }
func (s *Snapshot) Format() map[string]string {
	out := map[string]string{}
	for _, f := range s.files {
		out[f.Name()] = parser.FormatFile(f.Model())
	}
	return out
}
func (s *Snapshot) AST() map[string]any {
	files := []any{}
	for _, f := range s.files {
		files = append(files, map[string]any{"path": f.Name(), "ast": parser.Data(f.Model())})
	}
	symbols := []any{}
	for _, d := range s.model.Declarations {
		symbols = append(symbols, map[string]any{"system": s.system, "name": d.Name.Name, "kind": d.Kind, "span": parser.Data(d.Span)})
	}
	return map[string]any{"schema": Version, "valid": true, "profile": s.Profile(), "system": s.system, "revision": s.revision, "files": files, "symbols": symbols, "sources": s.Sources(), "edges": s.Edges(), "warnings": s.Warnings()}
}
func hash(text string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(text))) }
func ValidPath(p string) bool {
	if len(p) == 0 || len(p) > 1024 || path.IsAbs(p) || path.Clean(p) != p || !strings.HasSuffix(p, ".design") || strings.ContainsAny(p, "\\:") || strings.IndexFunc(p, unicode.IsControl) >= 0 {
		return false
	}
	for _, part := range strings.Split(p, "/") {
		if part == "." || part == ".." || part == "" {
			return false
		}
	}
	return true
}
func Dependencies(m *parser.Model) []parser.Include {
	out := append([]parser.Include{}, m.Includes...)
	for _, s := range m.Statements {
		if s.Path != "" {
			out = append(out, parser.Include{Path: s.Path + ".design", Span: s.PathSpan})
		}
	}
	return out
}

// Compile links only the entry's reachable files. Unrelated cached syntax cannot
// satisfy an unresolved name. It never opens a file or mutates an input AST.
func Compile(entry string, files []*parser.File, canonical bool) (*Snapshot, Issues) {
	byName := map[string]*parser.File{}
	var issues Issues
	for _, f := range files {
		if byName[f.Name()] != nil {
			return nil, Issues{{Code: "DUPLICATE_SOURCE", Message: f.Name()}}
		}
		byName[f.Name()] = f
	}
	snap := &Snapshot{entry: entry, model: &parser.Model{}, sources: []Source{}, edges: []Edge{}}
	seen := map[string]bool{}
	total, tokens := 0, 0
	var visit func(string, parser.Span)
	visit = func(name string, at parser.Span) {
		if seen[name] {
			return
		}
		seen[name] = true
		f := byName[name]
		if f == nil {
			issues = append(issues, Issue{Code: "MISSING_SOURCE", Message: name, Span: at})
			return
		}
		if len(seen) > MaxFiles {
			issues = append(issues, Issue{Code: "SOURCE_LIMIT", Message: "More than 128 files", Span: at})
			return
		}
		total += len(f.Text())
		if total > parser.MaxBytes {
			issues = append(issues, Issue{Code: "SOURCE_LIMIT", Message: "Aggregate source exceeds 2 MiB", Span: at})
			return
		}
		tokens += f.TokenCount()
		if tokens > 250000 {
			issues = append(issues, Issue{Code: "TOKEN_LIMIT", Message: "Aggregate token limit", Span: at})
			return
		}
		m := f.Model()
		snap.files = append(snap.files, f)
		if name == entry {
			snap.model.Header = m.Header
		}
		for _, dep := range Dependencies(m) {
			if !ValidPath(dep.Path) {
				issues = append(issues, Issue{Code: "SOURCE_PATH", Message: dep.Path, Span: dep.Span})
				continue
			}
			snap.edges = append(snap.edges, Edge{name, dep.Path, dep.Span})
			if seen[dep.Path] && len(snap.warnings) < MaxDiagnostics {
				snap.warnings = append(snap.warnings, Issue{Code: "REUSED_SOURCE", Message: dep.Path, Span: dep.Span})
			}
			visit(dep.Path, dep.Span)
		}
	}
	visit(entry, parser.Span{Source: entry, Line: 1, Column: 1})
	if len(issues) > 0 {
		return nil, ordered(issues)
	}
	sort.Slice(snap.files, func(i, j int) bool { return snap.files[i].Name() < snap.files[j].Name() })
	version := snap.model.Header.Version
	systems := []parser.Declaration{}
	declared := map[string]parser.Declaration{}
	for _, f := range snap.files {
		m := f.Model()
		snap.sources = append(snap.sources, Source{f.Name(), hash(f.Text())})
		if m.Header.Version != version {
			issues = append(issues, Issue{Code: "PROFILE_MISMATCH", Message: "All files must use the entry profile", Span: m.Header.Span})
		}
		for _, d := range m.Declarations {
			if old, ok := declared[d.Name.Name]; ok {
				issues = append(issues, Issue{Code: "DUPLICATE_DECLARATION", Message: d.Name.Name, Span: d.Name.Span, Related: []parser.Span{old.Name.Span}})
			} else {
				declared[d.Name.Name] = d
			}
			if d.Kind == "system" {
				systems = append(systems, d)
				if f.Name() != entry {
					issues = append(issues, Issue{Code: "SYSTEM_ENTRY", Message: "System must be declared in the selected entry", Span: d.Span})
				}
			}
		}
		snap.model.Declarations = append(snap.model.Declarations, m.Declarations...)
		for _, s := range m.Statements {
			if s.Path != "" {
				target := byName[s.Path+".design"]
				found := false
				if target != nil {
					for _, d := range target.Model().Declarations {
						if d.Name.Name == s.Object.Name {
							found = true
						}
					}
				}
				if !found {
					issues = append(issues, Issue{Code: "PATH_TARGET", Message: s.Path + " must declare " + s.Object.Name, Span: s.PathSpan})
				}
				s.Path = "" // Keep the authored path in FileAST; resolved facts use object identity.
			}
			snap.model.Statements = append(snap.model.Statements, s)
		}
		if canonical && f.Text() != parser.FormatFile(m) {
			issues = append(issues, Issue{Code: "NONCANONICAL_FORM", Message: "Format this file in its complete source context", Span: m.Header.Span})
		}
	}
	if version == "0.6" && len(systems) != 1 {
		sp := snap.model.Header.Span
		rel := []parser.Span{}
		for _, d := range systems {
			rel = append(rel, d.Span)
		}
		issues = append(issues, Issue{Code: "SYSTEM_CARDINALITY", Message: "Exactly one System is required; fragment parsing is available separately", Span: sp, Related: rel})
	}
	if len(issues) > 0 {
		return nil, ordered(issues)
	}
	if len(systems) == 1 {
		snap.system = systems[0].Name.Name
	}
	// Preserve the historical single-file span and ordering for profile 0.5.
	if version == "0.5" {
		snap.model = byName[entry].Model()
	} else {
		sort.Slice(snap.model.Declarations, func(i, j int) bool {
			return snap.model.Declarations[i].Name.Name < snap.model.Declarations[j].Name.Name
		})
		sort.Slice(snap.model.Statements, func(i, j int) bool { return snap.model.Statements[i].Sentence() < snap.model.Statements[j].Sentence() })
	}
	// Bound diagnostic enrichment independently of malformed input size. Index both
	// endpoints, including reverse edges in cross-file cycles. Related evidence is
	// a bounded neighbourhood, not a claim to enumerate every contributing fact.
	bySpan := map[parser.Span]parser.Statement{}
	byEndpoint := map[string]map[string][]parser.Span{}
	for _, fact := range snap.model.Statements {
		bySpan[fact.Span] = fact
		for _, name := range []string{fact.Subject.Name, fact.Object.Name} {
			if name == "" {
				continue
			}
			if byEndpoint[name] == nil {
				byEndpoint[name] = map[string][]parser.Span{}
			}
			group := byEndpoint[name][fact.Span.Source]
			if len(group) < MaxRelated {
				byEndpoint[name][fact.Span.Source] = append(group, fact.Span)
			}
		}
	}
	for _, d := range parser.Validate(snap.model) {
		if len(issues) == MaxDiagnostics {
			break
		}
		related := []parser.Span{}
		used := map[parser.Span]bool{}
		fact := bySpan[d.Span]
		for _, file := range snap.files {
			if file.Name() == d.Span.Source {
				continue
			}
			for _, name := range []string{fact.Subject.Name, fact.Object.Name} {
				for _, span := range byEndpoint[name][file.Name()] {
					if !used[span] && len(related) < MaxRelated {
						related = append(related, span)
						used[span] = true
					}
				}
			}
		}
		issues = append(issues, Issue{d.Code, d.Message, d.Span, related})
	}
	if len(issues) > 0 {
		return nil, ordered(issues)
	}
	if version == "0.5" {
		snap.revision = hash(byName[entry].Text())
	} else {
		h := sha256.New()
		put := func(s string) {
			var n [8]byte
			binary.BigEndian.PutUint64(n[:], uint64(len(s)))
			h.Write(n[:])
			h.Write([]byte(s))
		}
		put("SDL-source-graph-revision/1")
		put(snap.Profile())
		put(entry)
		for _, s := range snap.sources {
			put(s.Path)
			put(s.SHA256)
		}
		snap.revision = fmt.Sprintf("%x", h.Sum(nil))
	}
	return snap, nil
}
func ordered(ds Issues) Issues {
	sort.SliceStable(ds, func(i, j int) bool {
		a, b := ds[i], ds[j]
		if a.Span.Source != b.Span.Source {
			return a.Span.Source < b.Span.Source
		}
		if a.Span.Start != b.Span.Start {
			return a.Span.Start < b.Span.Start
		}
		return a.Code+a.Message < b.Code+b.Message
	})
	if len(ds) > MaxDiagnostics {
		ds = ds[:MaxDiagnostics]
	}
	for i := range ds {
		if len(ds[i].Related) > MaxRelated {
			ds[i].Related = ds[i].Related[:MaxRelated]
		}
	}
	return ds
}
