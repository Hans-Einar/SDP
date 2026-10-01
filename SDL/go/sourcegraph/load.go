package sourcegraph

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Loaded owns a checked snapshot and its actual input locations.
type Loaded struct {
	Snapshot    *Snapshot
	root, entry string
	paths       []string
}

func (l *Loaded) Fresh() error {
	n, e := Load(filepath.Join(l.root, l.entry), nil, true)
	if e != nil {
		return e
	}
	if n.Snapshot.Revision() != l.Snapshot.Revision() {
		return fmt.Errorf("source changed during generation")
	}
	return nil
}
func destination(p string) (string, error) {
	a, e := filepath.Abs(p)
	if e != nil {
		return "", e
	}
	v, e := filepath.EvalSymlinks(a)
	if e == nil {
		return v, nil
	}
	if !os.IsNotExist(e) {
		return "", e
	}
	parent := filepath.Dir(a)
	if parent == a {
		return "", e
	}
	v, e = destination(parent)
	return filepath.Join(v, filepath.Base(a)), e
}
func (l *Loaded) Outside(output string) error {
	dest, e := destination(output)
	if e != nil {
		return e
	}
	for _, p := range l.paths {
		rel, e := filepath.Rel(dest, p)
		if e != nil {
			return e
		}
		if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return fmt.Errorf("source must be outside output directory: %s", p)
		}
	}
	return nil
}
func Load(entry string, cache *parser.SyntaxCache, canonical bool) (*Loaded, error) {
	abs, e := filepath.Abs(entry)
	if e != nil {
		return nil, e
	}
	root, e := filepath.EvalSymlinks(filepath.Dir(abs))
	if e != nil {
		return nil, e
	}
	if cache == nil {
		cache = &parser.SyntaxCache{}
	}
	l := &Loaded{root: root, entry: filepath.Base(abs)}
	files := []*parser.File{}
	seen := map[string]bool{}
	infos := []os.FileInfo{}
	total, tokens := 0, 0
	var visit func(string, parser.Span) error
	visit = func(name string, at parser.Span) error {
		if seen[name] {
			return nil
		}
		if len(seen) >= MaxFiles {
			return Issue{Code: "SOURCE_LIMIT", Message: "More than 128 files", Span: at}
		}
		seen[name] = true
		if name != l.entry && !ValidPath(name) {
			return Issue{Code: "SOURCE_PATH", Message: name, Span: at}
		}
		full := root
		for _, part := range strings.Split(name, "/") {
			full = filepath.Join(full, part)
			st, e := os.Lstat(full)
			if e != nil {
				return Issue{Code: "SOURCE_READ", Message: e.Error(), Span: at}
			}
			if st.Mode()&os.ModeSymlink != 0 {
				return Issue{Code: "SOURCE_PATH", Message: "Symlinked sources are not supported: " + name, Span: at}
			}
		}
		st, e := os.Stat(full)
		if e != nil {
			return e
		}
		if !st.Mode().IsRegular() {
			return Issue{Code: "SOURCE_READ", Message: "Expected regular file: " + name, Span: at}
		}
		for _, old := range infos {
			if os.SameFile(old, st) {
				return Issue{Code: "SOURCE_ALIAS", Message: name, Span: at}
			}
		}
		infos = append(infos, st)
		f, e := os.Open(full)
		if e != nil {
			return e
		}
		actual, e := f.Stat()
		if e != nil {
			f.Close()
			return e
		}
		if !actual.Mode().IsRegular() || !os.SameFile(actual, st) {
			f.Close()
			return Issue{Code: "SOURCE_CHANGED", Message: name, Span: at}
		}
		b, e := io.ReadAll(io.LimitReader(f, int64(parser.MaxBytes-total)+1))
		f.Close()
		if e != nil {
			return e
		}
		total += len(b)
		if total > parser.MaxBytes {
			return Issue{Code: "SOURCE_LIMIT", Message: "Aggregate source exceeds 2 MiB", Span: at}
		}
		ast, e := cache.Parse(name, string(b))
		if e != nil {
			return e
		}
		m := ast.Model()
		// Count grammar tokens, not whitespace-delimited words. Limit also bounds work
		// across files already parsed independently by a syntax cache.
		tokens += ast.TokenCount()
		if tokens > 250000 {
			return Issue{Code: "TOKEN_LIMIT", Message: "Aggregate token limit", Span: at}
		}
		files = append(files, ast)
		l.paths = append(l.paths, full)
		for _, dep := range Dependencies(m) {
			if !ValidPath(dep.Path) {
				return Issue{Code: "SOURCE_PATH", Message: dep.Path, Span: dep.Span}
			}
			if e := visit(dep.Path, dep.Span); e != nil {
				return fmt.Errorf("included from %s:%d: %w", name, dep.Span.Line, e)
			}
		}
		return nil
	}
	if e := visit(l.entry, parser.Span{Source: l.entry, Line: 1, Column: 1}); e != nil {
		return nil, e
	}
	snap, ds := Compile(l.entry, files, canonical)
	if len(ds) > 0 {
		return nil, ds
	}
	l.Snapshot = snap
	return l, nil
}
