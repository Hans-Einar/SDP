package blueprints

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SDPTool/model"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
)

const MaxCatalogueEntries = 256
const MaxCatalogueFiles = 10000
const MaxCatalogueBytes = 64 << 20

var hashPattern = regexp.MustCompile("^[a-f0-9]{64}$")

type Entry struct {
	RetainedRevision string
	Path             string
	Key              string
	Revision         string
	TaskID           string
	System           string
	State            string
	Diagnostic       string
	Preliminary      bool
	Files            map[string]string
}
type Catalogue struct {
	State      string
	Diagnostic string
	Entries    []Entry
}
type budget struct {
	nodes int
	files int
	bytes int64
}

func (b *budget) read(p string, limit int64) ([]byte, error) {
	if b.files >= MaxCatalogueFiles {
		return nil, fmt.Errorf("catalogue file limit")
	}
	if e := noSymlinks(p); e != nil {
		return nil, e
	}
	st, e := os.Lstat(p)
	if e != nil {
		return nil, e
	}
	if !st.Mode().IsRegular() {
		return nil, fmt.Errorf("expected regular file")
	}
	b.files++
	if st.Size() > limit || st.Size() > MaxCatalogueBytes-b.bytes {
		return nil, fmt.Errorf("catalogue byte limit")
	}
	f, e := os.Open(p)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	readLimit := min(limit, MaxCatalogueBytes-b.bytes)
	data, e := io.ReadAll(io.LimitReader(f, readLimit+1))
	b.bytes += int64(len(data))
	if int64(len(data)) > limit || b.bytes > MaxCatalogueBytes {
		return nil, fmt.Errorf("catalogue byte limit")
	}
	return data, e
}
func key(d Document) string { return documents.Hash([]byte(d.Analysis.System + "\x00" + d.TaskID)) }
func safeRelative(p string) bool {
	return p != "" && p != "." && !filepath.IsAbs(p) && filepath.ToSlash(filepath.Clean(p)) == p && p != ".." && !strings.HasPrefix(p, "../") && !strings.ContainsAny(p, "\\:\x00")
}

// decodeRecord rejects duplicate keys.
func decodeRecord(data []byte, out any) error {
	d := json.NewDecoder(bytes.NewReader(data))
	var walk func() error
	walk = func() error {
		t, e := d.Token()
		if e != nil {
			return e
		}
		if v, ok := t.(json.Delim); ok {
			switch v {
			case '{':
				seen := map[string]bool{}
				for d.More() {
					tok, e := d.Token()
					if e != nil {
						return e
					}
					k := tok.(string)
					if seen[k] {
						return fmt.Errorf("duplicate JSON key")
					}
					seen[k] = true
					if e = walk(); e != nil {
						return e
					}
				}
			case '[':
				for d.More() {
					if e = walk(); e != nil {
						return e
					}
				}
			default:
				return fmt.Errorf("unexpected JSON delimiter")
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
	d = json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	return d.Decode(out)
}
func verify(dir string, b *budget) (Document, *documents.Bundle, error) {
	var d Document
	var m documents.Manifest
	fail := func(e error) (Document, *documents.Bundle, error) { return d, nil, e }
	manifest, e := b.read(filepath.Join(dir, "manifest.json"), 8<<20)
	if e != nil {
		return fail(e)
	}
	if e = decodeRecord(manifest, &m); e != nil {
		return fail(e)
	}
	if m.Version != Schema || !hashPattern.MatchString(m.Revision) || len(m.Outputs) == 0 || len(m.Outputs) > MaxCatalogueFiles {
		return fail(fmt.Errorf("unsupported or invalid manifest"))
	}
	files := map[string][]byte{"manifest.json": manifest}
	names := make([]string, 0, len(m.Outputs))
	for name := range m.Outputs {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		hash := m.Outputs[name]
		if !safeRelative(name) || name == "manifest.json" || !hashPattern.MatchString(hash) {
			return fail(fmt.Errorf("unsafe manifest entry"))
		}
		data, e := b.read(filepath.Join(dir, filepath.FromSlash(name)), MaxCatalogueBytes)
		if e != nil {
			return fail(e)
		}
		if documents.Hash(data) != hash {
			return fail(fmt.Errorf("modified managed file: %s", name))
		}
		files[name] = data
	}
	for _, name := range []string{"index.md", "changes.md", "context.md", "obligations.md", "evidence.md", "blueprint.json", "assignment.json", "context.mmd", "changes.mmd"} {
		if _, ok := files[name]; !ok {
			return fail(fmt.Errorf("missing required file %s", name))
		}
	}
	if e = decodeRecord(files["blueprint.json"], &d); e != nil {
		return fail(e)
	}
	if d.Schema != Schema || d.Status != "diagnostic-preview" || d.Analysis == nil || d.TaskID == "" || d.Analysis.System == "" || d.Revision != m.Revision {
		return fail(fmt.Errorf("invalid blueprint metadata"))
	}
	revision := d.Revision
	d.Revision = ""
	encoded, _ := json.Marshal(d)
	d.Revision = revision
	if documents.Hash(encoded) != revision || documents.Hash(files["assignment.json"]) != d.TaskDigest {
		return fail(fmt.Errorf("blueprint identity mismatch"))
	}
	for _, side := range []struct {
		name     string
		expected string
	}{{"NOW", d.Now.SourceDigest}, {"TARGET", d.Target.SourceDigest}} {
		captured := model.Files{}
		for name, data := range files {
			if strings.HasPrefix(name, "sources/"+side.name+"/") {
				captured[strings.TrimPrefix(name, "sources/"+side.name+"/")] = data
			}
		}
		if model.Digest(captured) != side.expected {
			return fail(fmt.Errorf("captured source identity mismatch: %s", side.name))
		}
	}
	for _, side := range []string{"NOW", "TARGET"} {
		if !safeRelative(d.Entry) {
			return fail(fmt.Errorf("invalid entry"))
		}
		if _, ok := files["sources/"+side+"/"+d.Entry]; !ok {
			return fail(fmt.Errorf("missing captured entry"))
		}
	}
	bundle := &documents.Bundle{Files: files, Manifest: m}
	check := *bundle
	check.Files = map[string][]byte{}
	for n, data := range files {
		if strings.HasPrefix(n, "sources/") {
			check.Files[n] = nil
		} else {
			check.Files[n] = data
		}
	}
	if e = check.CheckLinks(); e != nil {
		return fail(e)
	}
	// The entire directory is immutable; extra files are not alternate authority.
	var visit func(string, int) error
	visit = func(p string, depth int) error {
		if depth > 64 {
			return fmt.Errorf("directory depth limit")
		}
		if b.nodes >= MaxCatalogueFiles {
			return fmt.Errorf("catalogue directory entry limit")
		}
		entries, err := readDirs(p, MaxCatalogueFiles-b.nodes)
		if err != nil {
			b.nodes = MaxCatalogueFiles
			return err
		}
		b.nodes += len(entries)
		for _, entry := range entries {
			child := filepath.Join(p, entry.Name())
			if entry.Type()&os.ModeSymlink != 0 {
				return fmt.Errorf("symlink in retained revision")
			}
			if entry.IsDir() {
				if err = visit(child, depth+1); err != nil {
					return err
				}
				continue
			}
			rel, _ := filepath.Rel(dir, child)
			if _, ok := files[filepath.ToSlash(rel)]; !ok {
				return fmt.Errorf("unmanaged file in retained revision")
			}
		}
		return nil
	}
	e = visit(dir, 0)
	if e != nil {
		return fail(e)
	}
	return d, bundle, nil
}
func Scan(catalogue string) Catalogue {
	c := Catalogue{State: "absent", Entries: []Entry{}}
	if _, e := os.Lstat(catalogue); os.IsNotExist(e) {
		return c
	}
	if e := noSymlinks(catalogue); e != nil {
		c.State = "unavailable"
		c.Diagnostic = e.Error()
		return c
	}
	dirs, e := readDirs(catalogue, MaxCatalogueEntries)
	if e != nil {
		c.State = "unavailable"
		c.Diagnostic = e.Error()
		return c
	}
	c.State = "empty"
	budget := &budget{}
	count := len(dirs)
	addBad := func(p, msg string) { c.Entries = append(c.Entries, Entry{Path: p, State: "invalid", Diagnostic: msg}) }
	for _, taskDir := range dirs {
		taskPath := filepath.Join(catalogue, taskDir.Name())
		if !taskDir.IsDir() || !hashPattern.MatchString(taskDir.Name()) {
			addBad(taskPath, "invalid blueprint-key directory")
			continue
		}
		revisions, e := readDirs(taskPath, MaxCatalogueEntries-count)
		if e != nil {
			addBad(taskPath, e.Error())
			c.State = "unavailable"
			c.Diagnostic = e.Error()
			return c
		}
		count += len(revisions)
		if len(revisions) == 0 {
			addBad(taskPath, "no retained revisions")
		}
		for _, revision := range revisions {
			p := filepath.Join(taskPath, revision.Name())
			if !revision.IsDir() || !hashPattern.MatchString(revision.Name()) {
				addBad(p, "invalid revision directory")
				continue
			}
			d, b, e := verify(p, budget)
			if e != nil {
				addBad(p, e.Error())
				continue
			}
			retained := model.Digest(model.Files(b.Files))
			if key(d) != taskDir.Name() || retained != revision.Name() {
				addBad(p, "directory identity mismatch or duplicate revision")
				continue
			}
			c.Entries = append(c.Entries, Entry{Path: p, Key: key(d), Revision: d.Revision, RetainedRevision: retained, TaskID: d.TaskID, System: d.Analysis.System, State: "validated", Preliminary: d.Now.Preliminary || d.Target.Preliminary, Files: b.Manifest.Outputs})
		}
	}
	if len(c.Entries) > 0 {
		c.State = "available"
	}
	return c
}
func readDirs(p string, max int) ([]os.DirEntry, error) {
	if e := noSymlinks(p); e != nil {
		return nil, e
	}
	f, e := os.Open(p)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	entries, e := f.ReadDir(max + 1)
	if e != nil && e != io.EOF {
		return nil, e
	}
	if len(entries) > max {
		return nil, fmt.Errorf("catalogue directory limit")
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	return entries, nil
}

// GenerateRetained produces a diagnostic capture and retains it atomically as a
// new immutable revision. Equal existing bytes are idempotent, never overwritten.
func GenerateRetained(ctx context.Context, o Options, catalogue string) (Result, error) {
	if o.Output != "" || catalogue == "" {
		return Result{}, fmt.Errorf("choose output or catalogue")
	}
	if e := noSymlinks(catalogue); e != nil {
		return Result{}, e
	}
	temp, e := os.MkdirTemp("", "sdp-blueprint-preview-")
	if e != nil {
		return Result{}, e
	}
	defer os.RemoveAll(temp)
	o.Output = filepath.Join(temp, "preview")
	r, e := Generate(ctx, o)
	if e != nil {
		return Result{}, e
	}
	d, b, e := verify(o.Output, &budget{})
	if e != nil {
		return Result{}, e
	}
	fresh := func() error {
		for _, input := range []struct {
			ref      string
			expected Capture
		}{{o.From, d.Now}, {o.To, d.Target}} {
			v, e := model.Snapshot(o.Area, input.ref)
			if e != nil {
				return e
			}
			if e = outputSafe(catalogue, v.Path, o.Task, filepath.Join(o.Area, ".model-operations")); e != nil {
				return e
			}
			if capture(v) != input.expected {
				return fmt.Errorf("stale retained input")
			}
		}
		data, _, e := readTask(o.Task)
		if e != nil {
			return e
		}
		if documents.Hash(data) != d.TaskDigest {
			return fmt.Errorf("stale retained task")
		}
		return ctx.Err()
	}
	if e = fresh(); e != nil {
		return Result{}, e
	}

	parent := filepath.Join(catalogue, key(d))
	retained := model.Digest(model.Files(b.Files))
	dest := filepath.Join(parent, retained)
	if e = os.MkdirAll(parent, 0700); e != nil {
		return Result{}, e
	}
	if e = noSymlinks(dest); e != nil {
		return Result{}, e
	}
	lock := dest + ".lock"
	if e = os.Mkdir(lock, 0700); e != nil {
		return Result{}, fmt.Errorf("retained revision busy: %w", e)
	}
	defer os.Remove(lock)
	if _, e = os.Lstat(dest); e == nil {
		_, old, e := verify(dest, &budget{})
		if e != nil {
			return Result{}, e
		}
		if len(old.Files) != len(b.Files) {
			return Result{}, fmt.Errorf("retained revision differs")
		}
		for p, data := range b.Files {
			if !bytes.Equal(data, old.Files[p]) {
				return Result{}, fmt.Errorf("retained revision differs")
			}
		}
	} else if !os.IsNotExist(e) {
		return Result{}, e
	} else {
		stage, e := os.MkdirTemp(parent, ".retain-")
		if e != nil {
			return Result{}, e
		}
		defer os.RemoveAll(stage)
		for p, data := range b.Files {
			if e = ctx.Err(); e != nil {
				return Result{}, e
			}
			file := filepath.Join(stage, filepath.FromSlash(p))
			if e = os.MkdirAll(filepath.Dir(file), 0700); e != nil {
				return Result{}, e
			}
			if e = os.WriteFile(file, data, 0600); e != nil {
				return Result{}, e
			}
		}
		if e = ctx.Err(); e != nil {
			return Result{}, e
		}
		if e = fresh(); e != nil {
			return Result{}, e
		}
		if e = os.Rename(stage, dest); e != nil {
			return Result{}, e
		}
	}
	r.RetainedRevision = retained
	r.Path, _ = filepath.Abs(dest)
	return r, nil
}
