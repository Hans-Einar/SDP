package model

import (
	"fmt"
	ui "github.com/Hans-Einar/SDP/SDUI/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/sourcegraph"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

func validateSources(f Files) ([]string, error) {
	targets := []string{}
	covered := map[string]bool{}
	fragments := map[string]bool{}
	cache := &sdl.SyntaxCache{}
	_ = cache
	for name, b := range f {
		switch strings.ToLower(path.Ext(name)) {
		case ".sdui":
			d, e := ui.Parse(string(b))
			if e != nil {
				return nil, fail("validation", name+": "+e.Error())
			}
			if _, e = ui.Normalize(d); e != nil {
				return nil, e
			}
			targets = append(targets, name)
		case ".design":
			m, e := sdl.Parse(string(b))
			if e != nil {
				return nil, fail("validation", name+": "+e.Error())
			}
			if m.Header.Version == "0.6" {
				isRoot := false
				for _, d := range m.Declarations {
					if d.Kind == "system" {
						isRoot = true
					}
				}
				if !isRoot {
					fragments[name] = true
					continue
				}
				files := []*sdl.File{}
				base := path.Dir(name)
				for n, data := range f {
					rel, e := filepath.Rel(base, n)
					if e != nil || strings.HasPrefix(rel, "..") {
						continue
					}
					if path.Ext(n) != ".design" {
						continue
					}
					file, e := cache.Parse(filepath.ToSlash(rel), string(data))
					if e != nil {
						return nil, e
					}
					files = append(files, file)
				}
				snapshot, ds := sourcegraph.Compile(path.Base(name), files, false)
				if len(ds) > 0 {
					return nil, fail("validation", ds.Error())
				}
				for _, s := range snapshot.Sources() {
					covered[path.Join(base, s.Path)] = true
				}
			} else {
				_, ds := sdl.Check(string(b))
				if len(ds) > 0 {
					return nil, fail("validation", fmt.Sprint(ds))
				}
			}
			targets = append(targets, name)
		}
	}
	for name := range fragments {
		if !covered[name] {
			return nil, fail("validation", "orphan SDL fragment "+name)
		}
	}
	if len(targets) == 0 {
		return nil, fail("validation", "no supported SDL/SDUI entrypoints")
	}
	sort.Strings(targets)
	return targets, nil
}
func releases(area string) ([]Artifact, error) {
	ds, e := os.ReadDir(area)
	if e != nil {
		return nil, e
	}
	out := []Artifact{}
	names := map[string]string{}
	ids := map[string]bool{}
	for _, d := range ds {
		if !strings.HasPrefix(d.Name(), "RELEASE--") {
			continue
		}
		a, _, e := capture(filepath.Join(area, d.Name()))
		if e != nil {
			return nil, e
		}
		if a.Kind != "release" {
			return nil, fail("integrity", "release directory has wrong kind")
		}
		if id, ok := names[a.Name]; ok && id != a.ID {
			return nil, fail("ambiguous", "duplicate release version")
		}
		if ids[a.ID] {
			return nil, fail("ambiguous", "duplicate release identity")
		}
		names[a.Name] = a.ID
		ids[a.ID] = true
		out = append(out, a)
	}
	return out, nil
}
func latest(area string) (string, error) {
	rs, e := releases(area)
	if e != nil {
		return "", e
	}
	parents := map[string]bool{}
	for _, r := range rs {
		if r.Predecessor != "" {
			parents[r.Predecessor] = true
		}
	}
	heads := []Artifact{}
	for _, r := range rs {
		if !parents[r.ID] {
			heads = append(heads, r)
		}
	}
	if len(heads) != 1 {
		return "", fail("ambiguous", fmt.Sprintf("%d accepted release heads; use --initial or explicit source, reconcile competing heads", len(heads)))
	}
	return "release:" + heads[0].Name, nil
}
func Freeze(area, kind, name, ref, evidence string) (Result, error) {
	unlock, e := begin(area)
	if e != nil {
		return Result{}, e
	}
	defer unlock()
	if kind != "candidate" && kind != "proposal" && kind != "release" {
		return Result{}, fail("arguments", "invalid frozen kind")
	}
	if (kind == "release" && !versionRE.MatchString(name)) || (kind != "release" && !validName(name)) {
		return Result{}, fail("arguments", "invalid name/version")
	}
	if e = nameAvailable(area, kind, name); e != nil {
		return Result{}, e
	}
	p, e := resolve(area, ref)
	if e != nil {
		return Result{}, e
	}
	source, f, e := capture(p)
	if e != nil {
		return Result{}, e
	}
	if len(source.Conflicts) > 0 || len(source.PendingParents) > 0 {
		return Result{}, fail("conflict", "unresolved merge")
	}
	if kind == "release" && source.Kind != "candidate" {
		return Result{}, fail("arguments", "release requires candidate")
	}
	if kind != "release" && source.Kind != "work" {
		return Result{}, fail("arguments", "submission requires WORK")
	}
	targets, e := validateSources(f)
	if e != nil {
		return Result{}, e
	}
	a := Artifact{Schema: Schema, ID: uuid(), Kind: kind, Name: name, BaseRelease: source.BaseRelease, Ledger: strip(source.Ledger), Digest: Digest(f), Validation: targets, Evidence: evidence}
	if kind == "release" {
		if evidence != "model-only" {
			return Result{}, fail("evidence", "v1 accepts --evidence model-only; verified implementation evidence is not yet authenticated")
		}
		rs, e := releases(area)
		if e != nil {
			return Result{}, e
		}
		if len(rs) > 0 {
			last, e := latest(area)
			if e != nil {
				return Result{}, e
			}
			lp, e := resolve(area, last)
			if e != nil {
				return Result{}, e
			}
			la, _, e := capture(lp)
			if e != nil {
				return Result{}, e
			}
			if source.BaseRelease != la.ID {
				return Result{}, fail("stale", "candidate not based on current release")
			}
			a.Predecessor = la.ID
		} else if source.BaseRelease != "" {
			return Result{}, fail("reference", "candidate base release absent")
		}
	}
	r := newRecord(a.ID+":00000", kind, "create "+kind, []string{source.Head}, f, "")
	a.Head = r.ID
	a.Ledger = append(a.Ledger, r)
	expected, e := treeDigest(p)
	if e != nil {
		return Result{}, e
	}
	target := artifactDir(a)
	_, e = publish(area, target, "", "", func(stage string) error {
		if e := writeFiles(stage, f); e != nil {
			return e
		}
		if e := saveArtifact(stage, a); e != nil {
			return e
		}
		actual, e := treeDigest(p)
		if e != nil {
			return e
		}
		if actual != expected {
			return fail("stale", "promotion input changed")
		}
		return nil
	})
	if e != nil {
		return Result{}, e
	}
	out, e := Status(area, kind+":"+name)
	out.Operation = "create"
	return out, e
}
