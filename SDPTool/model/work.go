package model

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func artifactDir(a Artifact) string {
	switch a.Kind {
	case "work":
		return "WORK--" + a.Name
	case "release":
		return "RELEASE--V" + a.Name
	default:
		return strings.ToUpper(a.Kind) + "--" + a.Name + "--" + a.ID[len(a.ID)-4:]
	}
}
func nameAvailable(area, kind, name string) error {
	entries, e := os.ReadDir(area)
	if e != nil {
		return e
	}
	for _, d := range entries {
		if !strings.HasPrefix(d.Name(), strings.ToUpper(kind)+"--") {
			continue
		}
		a, e := readArtifact(filepath.Join(area, d.Name()))
		if e != nil {
			return e
		}
		if strings.EqualFold(a.Name, name) {
			return fail("exists", "artifact name already exists")
		}
	}
	return nil
}
func strip(ledger []Record) []Record {
	r := append([]Record(nil), ledger...)
	for i := range r {
		r[i].Payload = ""
	}
	return r
}
func saveBaseline(dir string, a *Artifact, f Files, message string, parents []string) error {
	payload := fmt.Sprintf(".commits/#%05d", a.Sequence)
	r := newRecord(fmt.Sprintf("%s:%05d", a.ID, a.Sequence), "baseline", message, parents, f, payload)
	if e := writeFiles(filepath.Join(dir, payload, "files"), f); e != nil {
		return e
	}
	if e := writeY(filepath.Join(dir, payload, "commit.yaml"), r); e != nil {
		return e
	}
	a.Ledger = append(a.Ledger, r)
	a.Head = r.ID
	a.Digest = r.Digest
	if e := os.MkdirAll(filepath.Join(dir, ".merge"), 0700); e != nil {
		return e
	}
	return saveArtifact(dir, *a)
}
func CreateWork(area, name, from string, initial bool) (Result, error) {
	unlock, e := begin(area)
	if e != nil {
		return Result{}, e
	}
	defer unlock()
	if !validName(name) {
		return Result{}, fail("arguments", "invalid WORK name")
	}
	if e = nameAvailable(area, "work", name); e != nil {
		return Result{}, e
	}
	if initial && from != "" {
		return Result{}, fail("arguments", "initial and from are exclusive")
	}
	a := Artifact{Schema: Schema, ID: uuid(), Kind: "work", Name: name}
	f := Files{}
	parents := []string{}
	if !initial {
		if from == "" {
			from, e = latest(area)
			if e != nil {
				return Result{}, e
			}
		}
		p, e := resolve(area, from)
		if e != nil {
			return Result{}, e
		}
		source, files, e := capture(p)
		if e != nil {
			return Result{}, e
		}
		if source.Kind != "release" {
			return Result{}, fail("arguments", "new WORK requires release or --initial; use multi-source integration separately")
		}
		a.BaseRelease = source.ID
		a.Ledger = strip(source.Ledger)
		parents = []string{source.Head}
		f = files
	}
	target := artifactDir(a)
	_, e = publish(area, target, "", "", func(stage string) error {
		if e := writeFiles(stage, f); e != nil {
			return e
		}
		return saveBaseline(stage, &a, f, "WORK baseline", parents)
	})
	if e != nil {
		return Result{}, e
	}
	return Status(area, "work:"+name)
}
func reconstruct(dir string, a Artifact, id string) (Files, error) {
	records := map[string]Record{}
	for _, r := range a.Ledger {
		records[r.ID] = r
	}
	var rec func(string, int) (Files, error)
	rec = func(id string, depth int) (Files, error) {
		if depth > 256 {
			return nil, fail("limit", "reconstruction depth")
		}
		r, ok := records[id]
		if !ok || r.Payload == "" {
			return nil, fail("history", "restoration payload unavailable")
		}
		var stored Record
		b, e := os.ReadFile(filepath.Join(dir, r.Payload, "commit.yaml"))
		if e != nil {
			return nil, e
		}
		if e = strict(b, &stored); e != nil {
			return nil, e
		}
		stored.Payload = r.Payload
		if !recordEqual(stored, r) {
			return nil, fail("integrity", "commit record changed")
		}
		f := Files{}
		if r.Kind != "baseline" && r.Kind != "checkpoint" && r.Kind != "restore" {
			if len(r.Parents) == 0 {
				return nil, fail("integrity", "missing commit parent")
			}
			f, e = rec(r.Parents[0], depth+1)
			if e != nil {
				return nil, e
			}
		}
		for _, p := range r.Deleted {
			delete(f, p)
		}
		payload, e := scan(filepath.Join(dir, r.Payload, "files"), false)
		if os.IsNotExist(e) && len(r.Inventory) == 0 {
			e = nil
		}
		if e != nil {
			return nil, e
		}
		for p, b := range payload {
			want, ok := r.Inventory[p]
			if !ok || hash(b) != want {
				return nil, fail("integrity", "unexpected or corrupt after-image")
			}
			f[p] = b
		}
		if len(f) != len(r.Inventory) || Digest(f) != r.Digest {
			return nil, fail("integrity", "reconstructed digest mismatch")
		}
		return f, nil
	}
	return rec(id, 0)
}
func addCommit(stage string, a *Artifact, f Files, message, kind string, extra []string) error {
	old, e := reconstruct(stage, *a, a.Head)
	if e != nil {
		return e
	}
	a.Sequence++
	payload := fmt.Sprintf(".commits/#%05d", a.Sequence)
	r := newRecord(fmt.Sprintf("%s:%05d", a.ID, a.Sequence), kind, message, append([]string{a.Head}, extra...), f, payload)
	changed := Files{}
	for p, b := range f {
		if kind == "checkpoint" || kind == "restore" || hash(b) != hash(old[p]) || old[p] == nil {
			changed[p] = b
		}
	}
	for p := range old {
		if _, ok := f[p]; !ok {
			r.Deleted = append(r.Deleted, p)
		}
	}
	sort.Strings(r.Deleted)
	if e = os.MkdirAll(filepath.Join(stage, payload, "files"), 0700); e != nil {
		return e
	}
	if e = writeFiles(filepath.Join(stage, payload, "files"), changed); e != nil {
		return e
	}
	if e = writeY(filepath.Join(stage, payload, "commit.yaml"), r); e != nil {
		return e
	}
	a.Ledger = append(a.Ledger, r)
	a.Head = r.ID
	a.Digest = r.Digest
	return saveArtifact(stage, *a)
}
func Commit(area, ref, message string, resolved bool) (Result, error) {
	if strings.TrimSpace(message) == "" {
		return Result{}, fail("arguments", "commit message required")
	}
	unlock, e := begin(area)
	if e != nil {
		return Result{}, e
	}
	defer unlock()
	p, e := resolve(area, ref)
	if e != nil {
		return Result{}, e
	}
	a, f, e := capture(p)
	if e != nil {
		return Result{}, e
	}
	if a.Kind != "work" {
		return Result{}, fail("immutable", "only WORK can commit")
	}
	if len(a.Conflicts) > 0 && !resolved {
		return Result{}, fail("conflict", "resolve files and use commit --resolved")
	}
	if resolved {
		for _, b := range f {
			if strings.Contains(string(b), "<<<<<<<") || strings.Contains(string(b), ">>>>>>>>") {
				return Result{}, fail("conflict", "conflict markers remain")
			}
		}
		if _, e = validateSources(f); e != nil {
			return Result{}, e
		}
	}
	if Digest(f) == a.Digest && len(a.PendingParents) == 0 {
		return result("commit", p, a, f), nil
	}
	expected, e := treeDigest(p)
	if e != nil {
		return Result{}, e
	}
	_, e = publish(area, filepath.Base(p), p, expected, func(stage string) error {
		extra := a.PendingParents
		a.PendingParents = nil
		a.Conflicts = nil
		kind := "commit"
		if len(extra) > 0 {
			kind = "merge"
		}
		return addCommit(stage, &a, f, message, kind, extra)
	})
	if e != nil {
		return Result{}, e
	}
	out, e := Status(area, ref)
	out.Operation = "commit"
	return out, e
}
func replaceSources(stage string, f Files) error {
	old, e := scan(stage, true)
	if e != nil {
		return e
	}
	for p := range old {
		if e = os.Remove(filepath.Join(stage, p)); e != nil {
			return e
		}
	}
	// Remove only empty source directories, deepest first. Reserved history is untouched.
	dirs := []string{}
	e = filepath.WalkDir(stage, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p == stage {
			return nil
		}
		rel, _ := filepath.Rel(stage, p)
		if rel == ".commits" || rel == ".merge" {
			return filepath.SkipDir
		}
		if d.IsDir() {
			dirs = append(dirs, p)
		}
		return nil
	})
	if e != nil {
		return e
	}
	for i := len(dirs) - 1; i >= 0; i-- {
		if e := os.Remove(dirs[i]); e != nil {
			return e
		}
	}
	return writeFiles(stage, f)
}
func Restore(area, ref, commit string) (Result, error) {
	unlock, e := begin(area)
	if e != nil {
		return Result{}, e
	}
	defer unlock()
	p, e := resolve(area, ref)
	if e != nil {
		return Result{}, e
	}
	a, dirty, e := capture(p)
	if e != nil {
		return Result{}, e
	}
	if a.Kind != "work" {
		return Result{}, fail("immutable", "only WORK can restore")
	}
	id := commit
	if !strings.Contains(id, ":") {
		id = a.ID + ":" + strings.TrimPrefix(id, "#")
	}
	target, e := reconstruct(p, a, id)
	if e != nil {
		return Result{}, e
	}
	expected, e := treeDigest(p)
	if e != nil {
		return Result{}, e
	}
	_, e = publish(area, filepath.Base(p), p, expected, func(stage string) error {
		if Digest(dirty) != a.Digest {
			if e := addCommit(stage, &a, dirty, "preserve dirty state before restore", "checkpoint", nil); e != nil {
				return e
			}
		}
		a.Conflicts = nil
		a.PendingParents = nil
		if e := replaceSources(stage, target); e != nil {
			return e
		}
		return addCommit(stage, &a, target, "restore "+id, "restore", nil)
	})
	if e != nil {
		return Result{}, e
	}
	r, e := Status(area, ref)
	r.Operation = "restore"
	return r, e
}
