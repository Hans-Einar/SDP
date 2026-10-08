package probe

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func must(t *testing.T, e error) {
	t.Helper()
	if e != nil {
		t.Fatal(e)
	}
}
func newStore(t *testing.T, id string, f Files) Store {
	t.Helper()
	s, e := New(filepath.Join(t.TempDir(), id), id, f)
	must(t, e)
	return s
}
func equal(t *testing.T, a, b Files) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("got %#v want %#v", a, b)
	}
}
func TestCommitDeleteRenameRestoreAndDirtyPreservation(t *testing.T) {
	base := Files{"A.design": "old", "nested/B.design": "keep", "empty": ""}
	s := newStore(t, "A", base)
	changed := Files{"Renamed.design": "new", "nested/B.design": "keep", "empty": ""}
	n, e := s.Commit(changed, "rename/edit A", nil, false)
	must(t, e)
	if n != 1 {
		t.Fatal(n)
	}
	if _, e = os.Stat(filepath.Join(s.dir(1), "files/nested/B.design")); !os.IsNotExist(e) {
		t.Fatal("unchanged file copied")
	}
	f, r, e := s.Read(1)
	must(t, e)
	equal(t, f, changed)
	if !reflect.DeepEqual(r[1].Deleted, []string{"A.design"}) {
		t.Fatal(r[1])
	}
	dirty := Files{"Renamed.design": "unsaved", "nested/B.design": "keep", "empty": ""}
	n, e = s.Restore(0, dirty)
	must(t, e)
	if n != 3 {
		t.Fatal("restore reused a number", n)
	}
	f, _, e = s.Read(n)
	must(t, e)
	equal(t, f, base)
	f, _, e = s.Read(2)
	must(t, e)
	equal(t, f, dirty)
}
func TestCorruptPayloadCannotRestore(t *testing.T) {
	s := newStore(t, "A", Files{"A": "original"})
	must(t, os.WriteFile(filepath.Join(s.dir(0), "files/A"), []byte("tampered"), 0600))
	if _, _, e := s.Read(0); e == nil {
		t.Fatal("tampering accepted")
	}
}
func TestEmptyInitialAndNoop(t *testing.T) {
	s := newStore(t, "I", Files{})
	n, e := s.Commit(Files{}, "nothing", nil, false)
	must(t, e)
	if n != 0 {
		t.Fatal(n)
	}
	if _, e = New(s.Root, "I", Files{}); e == nil {
		t.Fatal("overwrote existing work")
	}
}
func TestNoGitCopyReopensHistory(t *testing.T) {
	s := newStore(t, "A", Files{"file": "base"})
	_, e := s.Commit(Files{"file": "next"}, "next", nil, false)
	must(t, e)
	copyPath := filepath.Join(t.TempDir(), "copied")
	must(t, os.Mkdir(copyPath, 0700))
	must(t, os.CopyFS(copyPath, os.DirFS(s.Root)))
	must(t, os.RemoveAll(s.Root))
	f, _, e := (Store{copyPath}).Read(1)
	must(t, e)
	equal(t, f, Files{"file": "next"})
	if _, e = os.Stat(filepath.Join(copyPath, ".git")); !os.IsNotExist(e) {
		t.Fatal("unexpected Git")
	}
}
func TestProcessInterruptionAndResume(t *testing.T) {
	if root := os.Getenv("MG_CRASH_ROOT"); root != "" {
		s := Store{root}
		_, e := s.Commit(Files{"A": "next"}, "interrupted", nil, true)
		if e == nil {
			os.Exit(22)
		}
		os.Exit(23)
	}
	s := newStore(t, "A", Files{"A": "base"})
	cmd := exec.Command(os.Args[0], "-test.run=^TestProcessInterruptionAndResume$")
	cmd.Env = append(os.Environ(), "MG_CRASH_ROOT="+s.Root)
	e := cmd.Run()
	if ee, ok := e.(*exec.ExitError); !ok || ee.ExitCode() != 23 {
		t.Fatalf("wrong child failure: %v", e)
	}
	reopened := Store{s.Root}
	h, e := reopened.head()
	must(t, e)
	if h.Seq != 0 {
		t.Fatal("premature head")
	}
	f, _, e := reopened.Read(0)
	must(t, e)
	equal(t, f, Files{"A": "base"})
	if _, e = reopened.Commit(Files{"A": "other"}, "blocked", nil, false); e == nil {
		t.Fatal("ignored pending transaction")
	}
	must(t, reopened.Resume())
	f, _, e = reopened.Read(1)
	must(t, e)
	equal(t, f, Files{"A": "next"})
}
func TestInterruptedCorruptCommitNotPublished(t *testing.T) {
	s := newStore(t, "A", Files{"A": "base"})
	_, e := s.Commit(Files{"A": "next"}, "bad pending", nil, true)
	if e == nil {
		t.Fatal("missing fault")
	}
	must(t, os.WriteFile(filepath.Join(s.dir(1), "files/A"), []byte("bad"), 0600))
	if e = s.Resume(); e == nil {
		t.Fatal("published corrupt state")
	}
	h, e := s.head()
	must(t, e)
	if h.Seq != 0 {
		t.Fatal(h)
	}
}
func TestMergeConflictMatrix(t *testing.T) {
	for _, tc := range []struct {
		name             string
		base, a, b, want Files
		conflicts        []string
	}{
		{"disjoint", Files{"a": "0", "b": "0"}, Files{"a": "1", "b": "0"}, Files{"a": "0", "b": "2"}, Files{"a": "1", "b": "2"}, []string{}},
		{"identical", Files{"a": "0"}, Files{"a": "1"}, Files{"a": "1"}, Files{"a": "1"}, []string{}},
		{"delete-modify", Files{"a": "0"}, Files{}, Files{"a": "1"}, Files{}, []string{"a"}},
		{"same-path-add", Files{}, Files{"a": "1"}, Files{"a": "2"}, Files{}, []string{"a"}},
		{"overlap", Files{"a": "0"}, Files{"a": "1"}, Files{"a": "2"}, Files{}, []string{"a"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, c := Merge(tc.base, tc.a, tc.b)
			equal(t, f, tc.want)
			if !reflect.DeepEqual(c, tc.conflicts) {
				t.Fatal(c)
			}
		})
	}
}
func TestMergeArchivePromotionAndWholeStateRollback(t *testing.T) {
	base := Files{"a": "0", "b": "0"}
	a := newStore(t, "A", base)
	b := newStore(t, "B", base)
	av := Files{"a": "1", "b": "0"}
	bv := Files{"a": "0", "b": "2"}
	_, e := a.Commit(av, "A edit", nil, false)
	must(t, e)
	_, e = b.Commit(bv, "B edit", nil, false)
	must(t, e)
	// Archive both pre-merge inputs outside target traversal, then place archives.
	staging := t.TempDir()
	for _, s := range []Store{a, b} {
		h, e := s.head()
		must(t, e)
		dest := filepath.Join(staging, h.ID)
		must(t, os.Mkdir(dest, 0700))
		must(t, os.CopyFS(dest, os.DirFS(s.Root)))
	}
	for _, name := range []string{"A", "B"} {
		must(t, os.Rename(filepath.Join(staging, name), filepath.Join(a.Root, ".merge", name)))
	}
	merged, c := Merge(base, av, bv)
	if len(c) != 0 {
		t.Fatal(c)
	}
	_, br, e := b.Read(1)
	must(t, e)
	n, e := a.Commit(merged, "merge B", []string{br[1].ID}, false)
	must(t, e)
	_, ar, e := a.Read(n)
	must(t, e)
	dest := filepath.Join(t.TempDir(), "CANDIDATE--combined")
	must(t, Promote(dest, "C", merged, append(ar, br...)))
	var meta Frozen
	must(t, readYAML(filepath.Join(dest, "artifact.yaml"), &meta))
	if meta.RestorePayloadAvailable || len(meta.Lineage) != 5 {
		t.Fatal(meta)
	}
	for _, name := range []string{".commits", ".merge"} {
		if _, e = os.Stat(filepath.Join(dest, name)); !os.IsNotExist(e) {
			t.Fatal("history payload leaked")
		}
	}
	restored, e := a.Restore(1, merged)
	must(t, e)
	f, _, e := a.Read(restored)
	must(t, e)
	equal(t, f, av)
	// Promotion remains browsable after both source WORK directories disappear.
	must(t, os.RemoveAll(a.Root))
	must(t, os.RemoveAll(b.Root))
	must(t, readYAML(filepath.Join(dest, "artifact.yaml"), &meta))
	if meta.Digest != digest(merged) {
		t.Fatal("wrong promoted digest")
	}
}
func TestLineageDeduplicatesAndRejectsIdentityConflict(t *testing.T) {
	s := newStore(t, "A", Files{"a": "1"})
	f, r, e := s.Read(0)
	must(t, e)
	must(t, Promote(filepath.Join(t.TempDir(), "ok"), "C", f, append(r, r...)))
	bad := r[0]
	bad.Message = "different"
	if e = Promote(filepath.Join(t.TempDir(), "bad"), "D", f, append(r, bad)); e == nil {
		t.Fatal("identity conflict accepted")
	}
}
func TestYAMLRejectsUnknownAndDuplicateFields(t *testing.T) {
	p := filepath.Join(t.TempDir(), "bad.yaml")
	for _, b := range []string{"schema: x\nid: a\nseq: 0\nunknown: true\n", "id: a\nid: b\n"} {
		must(t, os.WriteFile(p, []byte(b), 0600))
		var h Head
		if e := readYAML(p, &h); e == nil {
			t.Fatal("invalid metadata accepted")
		}
	}
}
func TestSameFileThreeWayPrimitive(t *testing.T) {
	// External Git executable for algorithm comparison only: no repository or Git commit.
	git, e := exec.LookPath("git")
	if e != nil {
		t.Fatal("probe requires git merge-file for comparison")
	}
	dir := t.TempDir()
	base := "one\ntwo\nthree\nfour\nfive\nsix\nseven\n"
	a := "ONE\ntwo\nthree\nfour\nfive\nsix\nseven\n"
	b := "one\ntwo\nthree\nfour\nfive\nsix\nSEVEN\n"
	must(t, writeFiles(dir, Files{"base": base, "a": a, "b": b}))
	cmd := exec.Command(git, "merge-file", "-p", "a", "base", "b")
	cmd.Dir = dir
	out, e := cmd.Output()
	must(t, e)
	if string(out) != "ONE\ntwo\nthree\nfour\nfive\nsix\nSEVEN\n" {
		t.Fatal(string(out))
	}
	must(t, os.WriteFile(filepath.Join(dir, "b"), []byte("OTHER\ntwo\nthree\nfour\nfive\nsix\nseven\n"), 0600))
	cmd = exec.Command(git, "merge-file", "-p", "a", "base", "b")
	cmd.Dir = dir
	if e = cmd.Run(); e == nil {
		t.Fatal("overlap reported clean")
	}
}
func TestMissingEmptyFileIsNotAnExistingEmptyFile(t *testing.T) {
	s := newStore(t, "A", Files{"empty": ""})
	must(t, os.Remove(filepath.Join(s.dir(0), "files/empty")))
	if _, _, e := s.Read(0); e == nil {
		t.Fatal("missing empty file accepted")
	}
}
func TestRepeatedFileMergeIsContentIdempotent(t *testing.T) {
	base := Files{"a": "0", "b": "0"}
	a := Files{"a": "1", "b": "0"}
	b := Files{"a": "0", "b": "2"}
	first, c := Merge(base, a, b)
	if len(c) != 0 {
		t.Fatal(c)
	}
	second, c := Merge(base, first, b)
	if len(c) != 0 {
		t.Fatal(c)
	}
	equal(t, first, second)
}
