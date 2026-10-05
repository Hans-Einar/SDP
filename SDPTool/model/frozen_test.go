package model

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPromotionStripsPayloadAndKeepsLineage(t *testing.T) {
	area := t.TempDir()
	p := initial(t, area, "A")
	must(t, os.WriteFile(filepath.Join(p, "Main.design"), []byte("language design-core version 0.5.\nunit Main.\n"), 0600))
	_, e := Commit(area, "work:A", "valid source", false)
	must(t, e)
	r, e := Freeze(area, "candidate", "A", "work:A", "")
	must(t, e)
	for _, record := range r.Artifact.Ledger {
		if record.Payload != "" {
			t.Fatal("payload survived")
		}
	}
	for _, name := range []string{".commits", ".merge"} {
		if _, e = os.Stat(filepath.Join(r.Path, name)); !os.IsNotExist(e) {
			t.Fatal("history copied")
		}
	}
	must(t, os.RemoveAll(p))
	_, e = Status(area, "candidate:A")
	must(t, e)
	if _, e = Commit(area, "candidate:A", "bad", false); e == nil {
		t.Fatal("frozen committed")
	}
	must(t, os.WriteFile(filepath.Join(r.Path, "Main.design"), []byte("tampered"), 0600))
	if _, e = Status(area, "candidate:A"); e == nil {
		t.Fatal("tampered accepted")
	}
}
func TestCandidateRealComposedSDL(t *testing.T) {
	area := t.TempDir()
	p := initial(t, area, "Composed")
	must(t, writeFiles(p, Files{"System.design": []byte("language design-core version 0.6.\nsystem Example.\nincludes \"Child.design\".\nExample contains Child.\n"), "Child.design": []byte("language design-core version 0.6.\nunit Child.\n")}))
	r, e := Freeze(area, "candidate", "Composed", "work:Composed", "")
	must(t, e)
	if len(r.Artifact.Validation) != 1 {
		t.Fatal(r.Artifact.Validation)
	}
	must(t, os.WriteFile(filepath.Join(p, "Orphan.design"), []byte("language design-core version 0.6.\nunit Orphan.\n"), 0600))
	if _, e = Freeze(area, "candidate", "Bad", "work:Composed", ""); e == nil {
		t.Fatal("orphan accepted")
	}
}
func TestReleaseCurrentHeadAndEvidence(t *testing.T) {
	area := t.TempDir()
	releaseFixture(t, area)
	_, e := CreateWork(area, "Next", "", false)
	must(t, e)
	_, e = Freeze(area, "candidate", "Next", "work:Next", "")
	must(t, e)
	if _, e = Freeze(area, "release", "0.2.0", "candidate:Next", ""); e == nil {
		t.Fatal("missing evidence accepted")
	}
	_, e = Freeze(area, "release", "0.2.0", "candidate:Next", "verified:"+strings.Repeat("a", 64)+":test-report-1")
	must(t, e)
	r, e := CreateWork(area, "Current", "", false)
	must(t, e)
	rp, e := resolve(area, "release:0.2.0")
	must(t, e)
	a, e := readArtifact(rp)
	must(t, e)
	if r.Artifact.BaseRelease != a.ID {
		t.Fatal("wrong default")
	}
	if _, e = Freeze(area, "release", "0.3.0", "candidate:Next", "model-only"); e == nil {
		t.Fatal("stale candidate released")
	}
}
func TestCompetingHeadsPreserved(t *testing.T) {
	base := t.TempDir()
	releaseFixture(t, base)
	clone := t.TempDir()
	must(t, os.CopyFS(clone, os.DirFS(base)))
	for _, area := range []string{base, clone} {
		_, e := CreateWork(area, "Next", "", false)
		must(t, e)
		_, e = Freeze(area, "candidate", "Next", "work:Next", "")
		must(t, e)
	}
	_, e := Freeze(base, "release", "0.2.0", "candidate:Next", "model-only")
	must(t, e)
	_, e = Freeze(clone, "release", "0.3.0", "candidate:Next", "model-only")
	must(t, e)
	incoming := filepath.Join(base, "RELEASE--V0.3.0")
	must(t, os.Mkdir(incoming, 0700))
	must(t, os.CopyFS(incoming, os.DirFS(filepath.Join(clone, "RELEASE--V0.3.0"))))
	if _, e = CreateWork(base, "Ambiguous", "", false); e == nil {
		t.Fatal("picked one of competing heads")
	}
	if _, e = Status(base, "release:0.2.0"); e != nil {
		t.Fatal(e)
	}
	if _, e = Status(base, "release:0.3.0"); e != nil {
		t.Fatal(e)
	}
}
func TestRejectInvalidModelButAllowRecoveryCommit(t *testing.T) {
	area := t.TempDir()
	p := initial(t, area, "Bad")
	must(t, os.WriteFile(filepath.Join(p, "Main.design"), []byte("not valid SDL"), 0600))
	_, e := Commit(area, "work:Bad", "intermediate", false)
	must(t, e)
	if _, e = Freeze(area, "candidate", "Bad", "work:Bad", ""); e == nil {
		t.Fatal("invalid model frozen")
	}
}

func TestTwoGitClonesTransportCompetingReleases(t *testing.T) {
	if _, e := exec.LookPath("git"); e != nil {
		t.Skip("Git transport test requires git; model runtime does not")
	}
	git := func(dir string, args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Test", "GIT_AUTHOR_EMAIL=test@example.invalid", "GIT_COMMITTER_NAME=Test", "GIT_COMMITTER_EMAIL=test@example.invalid")
		if b, e := cmd.CombinedOutput(); e != nil {
			t.Fatalf("git %v: %v %s", args, e, b)
		}
	}
	origin := t.TempDir()
	releaseFixture(t, origin)
	git(origin, "init", "-b", "main")
	git(origin, "add", "RELEASE--V0.1.0")
	git(origin, "commit", "-m", "baseline")
	a := filepath.Join(t.TempDir(), "a")
	b := filepath.Join(t.TempDir(), "b")
	git(origin, "clone", origin, a)
	git(origin, "clone", origin, b)
	for i, dir := range []string{a, b} {
		_, e := CreateWork(dir, "Next", "", false)
		must(t, e)
		_, e = Freeze(dir, "candidate", "Next", "work:Next", "")
		must(t, e)
		version := []string{"0.2.0", "0.3.0"}[i]
		_, e = Freeze(dir, "release", version, "candidate:Next", "model-only")
		must(t, e)
		git(dir, "add", "RELEASE--V"+version)
		git(dir, "commit", "-m", "independent release")
	}
	git(a, "fetch", b, "main")
	git(a, "merge", "--no-edit", "FETCH_HEAD")
	for _, ref := range []string{"release:0.2.0", "release:0.3.0"} {
		_, e := Status(a, ref)
		must(t, e)
	}
	if _, e := CreateWork(a, "Ambiguous", "", false); e == nil {
		t.Fatal("silently chose competing head")
	}
}
