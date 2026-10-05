package model

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Bounded line-based three-way merge. LCS is deliberately capped; large/ambiguous
// edits become explicit conflicts rather than unbounded allocation or guessed edits.
type edit struct {
	lo, hi int
	lines  []string
}

func edits(base, other []string) ([]edit, bool) {
	if len(base)*len(other) > 2000000 {
		return nil, false
	}
	w := len(other) + 1
	dp := make([]int, (len(base)+1)*w)
	for i := len(base) - 1; i >= 0; i-- {
		for j := len(other) - 1; j >= 0; j-- {
			if base[i] == other[j] {
				dp[i*w+j] = dp[(i+1)*w+j+1] + 1
			} else {
				a, b := dp[(i+1)*w+j], dp[i*w+j+1]
				if a > b {
					dp[i*w+j] = a
				} else {
					dp[i*w+j] = b
				}
			}
		}
	}
	out := []edit{}
	i, j := 0, 0
	for i < len(base) || j < len(other) {
		if i < len(base) && j < len(other) && base[i] == other[j] {
			i++
			j++
			continue
		}
		e := edit{lo: i}
		for i < len(base) || j < len(other) {
			if i < len(base) && j < len(other) && base[i] == other[j] {
				break
			}
			if j < len(other) && (i == len(base) || dp[i*w+j+1] >= dp[(i+1)*w+j]) {
				e.lines = append(e.lines, other[j])
				j++
			} else {
				i++
			}
		}
		e.hi = i
		out = append(out, e)
	}
	return out, true
}
func mergeText(base, a, b []byte) ([]byte, bool) {
	if !textual(base) || !textual(a) || !textual(b) {
		return nil, false
	}
	lines := strings.SplitAfter(string(base), "\n")
	ea, ok := edits(lines, strings.SplitAfter(string(a), "\n"))
	if !ok {
		return nil, false
	}
	eb, ok := edits(lines, strings.SplitAfter(string(b), "\n"))
	if !ok {
		return nil, false
	}
	all := append(ea, eb...)
	sort.SliceStable(all, func(i, j int) bool { return all[i].lo < all[j].lo })
	merged := []edit{}
	for _, e := range all {
		if len(merged) > 0 {
			prev := merged[len(merged)-1]
			if e.lo <= prev.hi {
				if e.lo == prev.lo && e.hi == prev.hi && strings.Join(e.lines, "") == strings.Join(prev.lines, "") {
					continue
				}
				return nil, false
			}
		}
		merged = append(merged, e)
	}
	out := []string{}
	at := 0
	for _, e := range merged {
		out = append(out, lines[at:e.lo]...)
		out = append(out, e.lines...)
		at = e.hi
	}
	out = append(out, lines[at:]...)
	return []byte(strings.Join(out, "")), true
}
func pointer(b []byte, exists bool) *string {
	if !exists {
		return nil
	}
	s := string(b)
	return &s
}
func mergeFiles(base, a, b Files) (Files, []Conflict) {
	out := Files{}
	conflicts := []Conflict{}
	keys := map[string]bool{}
	for _, f := range []Files{base, a, b} {
		for p := range f {
			keys[p] = true
		}
	}
	for p := range keys {
		v0, o0 := base[p]
		va, oa := a[p]
		vb, ob := b[p]
		switch {
		case oa == ob && bytes.Equal(va, vb):
			if oa {
				out[p] = va
			}
		case oa == o0 && bytes.Equal(va, v0):
			if ob {
				out[p] = vb
			}
		case ob == o0 && bytes.Equal(vb, v0):
			if oa {
				out[p] = va
			}
		default:
			if o0 && oa && ob {
				if text, ok := mergeText(v0, va, vb); ok {
					out[p] = text
					continue
				}
			}
			if oa {
				out[p] = va
			}
			conflicts = append(conflicts, Conflict{p, pointer(v0, o0), pointer(va, oa), pointer(vb, ob)})
		}
	}
	sort.Slice(conflicts, func(i, j int) bool { return conflicts[i].Path < conflicts[j].Path })
	return out, conflicts
}
func ancestors(a Artifact, head string) map[string]bool {
	records := map[string]Record{}
	for _, r := range a.Ledger {
		records[r.ID] = r
	}
	seen := map[string]bool{}
	var walk func(string)
	walk = func(id string) {
		if seen[id] {
			return
		}
		seen[id] = true
		for _, p := range records[id].Parents {
			walk(p)
		}
	}
	walk(head)
	return seen
}
func mergeBase(ap string, a Artifact, bp string, b Artifact) (Files, error) {
	aa := ancestors(a, a.Head)
	bb := ancestors(b, b.Head)
	candidates := []Record{}
	for _, r := range a.Ledger {
		if aa[r.ID] && bb[r.ID] {
			candidates = append(candidates, r)
		}
	}
	best := []Record{}
	for _, r := range candidates {
		older := false
		for _, s := range candidates {
			if s.ID != r.ID && ancestors(a, s.ID)[r.ID] {
				older = true
				break
			}
		}
		if !older {
			best = append(best, r)
		}
	}
	if len(best) != 1 {
		return nil, fail("merge-base", fmt.Sprintf("expected one common ancestor, got %d", len(best)))
	}
	for i, art := range []Artifact{a, b} {
		dir := []string{ap, bp}[i]
		for _, r := range art.Ledger {
			if r.Digest == best[0].Digest && r.Payload != "" {
				f, e := reconstruct(dir, art, r.ID)
				if e == nil {
					return f, nil
				}
			}
		}
	}
	return nil, fail("history", "common ancestor payload unavailable")
}
func Merge(area, source, target, newName string) (Result, error) {
	unlock, e := begin(area)
	if e != nil {
		return Result{}, e
	}
	defer unlock()
	sp, e := resolve(area, source)
	if e != nil {
		return Result{}, e
	}
	tp, e := resolve(area, target)
	if e != nil {
		return Result{}, e
	}
	if sp == tp {
		return Result{}, fail("arguments", "cannot merge artifact into itself")
	}
	sa, sf, e := capture(sp)
	if e != nil {
		return Result{}, e
	}
	ta, tf, e := capture(tp)
	if e != nil {
		return Result{}, e
	}
	if ta.Kind != "work" || sa.Kind != "work" {
		return Result{}, fail("arguments", "v1 integration requires WORK inputs")
	}
	if len(ta.Conflicts) > 0 || len(sa.Conflicts) > 0 {
		return Result{}, fail("conflict", "resolve existing conflicts first")
	}
	// Validate shared identities before any ancestry-based fast path.
	knownInput := map[string]Record{}
	for _, r := range ta.Ledger {
		r.Payload = ""
		knownInput[r.ID] = r
	}
	for _, r := range sa.Ledger {
		r.Payload = ""
		if old, ok := knownInput[r.ID]; ok && !recordEqual(old, r) {
			return Result{}, fail("integrity", "shared identity changed")
		}
	}
	if ancestors(ta, ta.Head)[sa.Head] && Digest(sf) == sa.Digest && newName == "" {
		return result("merge", tp, ta, tf), nil
	}
	if newName == "" && Digest(sf) != sa.Digest {
		reachable := ancestors(ta, ta.Head)
		for _, r := range ta.Ledger {
			if r.Origin == sa.ID && r.Digest == Digest(sf) && reachable[r.ID] {
				return result("merge", tp, ta, tf), nil
			}
		}
	}
	base, e := mergeBase(tp, ta, sp, sa)
	if e != nil {
		return Result{}, e
	}
	merged, conflicts := mergeFiles(base, tf, sf)
	expected, e := treeDigest(tp)
	if e != nil {
		return Result{}, e
	}
	sourceExpected, e := treeDigest(sp)
	if e != nil {
		return Result{}, e
	}
	dest := filepath.Base(tp)
	if newName != "" {
		if !validName(newName) {
			return Result{}, fail("arguments", "invalid combined WORK name")
		}
		if e = nameAvailable(area, "work", newName); e != nil {
			return Result{}, e
		}
		dest = "WORK--" + newName
	}
	_, e = publish(area, dest, "", func() string {
		if newName != "" {
			return ""
		}
		return expected
	}(), func(stage string) error {
		// Preserve target under its original name before changing the root identity.
		if e := os.CopyFS(stage, os.DirFS(tp)); e != nil {
			return e
		}
		if dest != filepath.Base(tp) {
			if e := os.Remove(filepath.Join(stage, metadataName(tp))); e != nil {
				return e
			}
		}
		// Snapshot target before integration; never traverse the staging destination.
		targetArchive := filepath.Join(stage, ".merge", ta.ID+"-"+expected[:12], filepath.Base(tp))
		if e := os.MkdirAll(targetArchive, 0700); e != nil {
			return e
		}
		if e := os.CopyFS(targetArchive, os.DirFS(tp)); e != nil {
			return e
		}
		originalTarget := ta
		if newName != "" {
			ta = Artifact{Schema: Schema, ID: uuid(), Kind: "work", Name: newName, BaseRelease: ta.BaseRelease, Ledger: strip(ta.Ledger)}
			if e := os.RemoveAll(filepath.Join(stage, ".commits")); e != nil {
				return e
			}
			if e := saveBaseline(stage, &ta, tf, "combined WORK baseline", []string{originalTarget.Head}); e != nil {
				return e
			}
		}
		archive := filepath.Join(".merge", sa.ID+"-"+sourceExpected[:12], filepath.Base(sp))
		ad := filepath.Join(stage, archive)
		if e := os.MkdirAll(ad, 0700); e != nil {
			return e
		}
		if e := os.CopyFS(ad, os.DirFS(sp)); e != nil {
			return e
		}
		if Digest(sf) != sa.Digest {
			captureID := uuid()
			payload := ".commits/capture-" + captureID
			r := newRecord(captureID+":00000", "checkpoint", "captured merge input", []string{sa.Head}, sf, payload)
			r.Origin = sa.ID
			if e := os.MkdirAll(filepath.Join(ad, payload, "files"), 0700); e != nil {
				return e
			}
			if e := writeFiles(filepath.Join(ad, payload, "files"), sf); e != nil {
				return e
			}
			if e := writeY(filepath.Join(ad, payload, "commit.yaml"), r); e != nil {
				return e
			}
			sa.Ledger = append(sa.Ledger, r)
			sa.Head = r.ID
			sa.Digest = r.Digest
		}
		known := map[string]Record{}
		for _, r := range ta.Ledger {
			known[r.ID] = r
		}
		for _, r := range sa.Ledger {
			if old, ok := known[r.ID]; ok {
				old.Payload = ""
				logical := r
				logical.Payload = ""
				if !recordEqual(old, logical) {
					return fail("integrity", "same commit ID has different metadata")
				}
				continue
			}
			if r.Payload != "" {
				r.Payload = filepath.ToSlash(filepath.Join(archive, r.Payload))
			}
			ta.Ledger = append(ta.Ledger, r)
		}
		if Digest(tf) != ta.Digest {
			if e := addCommit(stage, &ta, tf, "pre-merge dirty checkpoint", "checkpoint", nil); e != nil {
				return e
			}
		}
		if e := replaceSources(stage, merged); e != nil {
			return e
		}
		if len(conflicts) > 0 {
			ta.Conflicts = conflicts
			ta.PendingParents = []string{sa.Head}
			if e := saveArtifact(stage, ta); e != nil {
				return e
			}
		} else {
			if e := addCommit(stage, &ta, merged, "merge "+source, "merge", []string{sa.Head}); e != nil {
				return e
			}
		}
		actual, e := treeDigest(sp)
		if e != nil {
			return e
		}
		if actual != sourceExpected {
			return fail("stale", "merge source changed")
		}
		actual, e = treeDigest(tp)
		if e != nil {
			return e
		}
		if actual != expected {
			return fail("stale", "merge target changed")
		}
		return nil
	})
	if e != nil {
		return Result{}, e
	}
	ref := target
	if newName != "" {
		ref = "work:" + newName
	}
	r, e := Status(area, ref)
	r.Operation = "merge"
	return r, e
}
