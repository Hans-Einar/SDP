package model

import (
	"os"
	"path/filepath"
	"testing"
)

func releaseFixture(t *testing.T, area string) {
	t.Helper()
	p := initial(t, area, "Initial")
	must(t, os.WriteFile(filepath.Join(p, "Main.design"), []byte("language design-core version 0.5.\nunit Main.\n"), 0600))
	_, e := Freeze(area, "candidate", "Initial", "work:Initial", "")
	must(t, e)
	_, e = Freeze(area, "release", "0.1.0", "candidate:Initial", "model-only")
	must(t, e)
}
func TestThreeWayText(t *testing.T) {
	base := []byte("a\nb\nc\nd\ne\nf\ng\n")
	a := []byte("A\nb\nc\nd\ne\nf\ng\n")
	b := []byte("a\nb\nc\nd\ne\nf\nG\n")
	r, ok := mergeText(base, a, b)
	if !ok || string(r) != "A\nb\nc\nd\ne\nf\nG\n" {
		t.Fatalf("%v %q", ok, r)
	}
	if _, ok = mergeText(base, a, []byte("OTHER\nb\nc\nd\ne\nf\ng\n")); ok {
		t.Fatal("overlap accepted")
	}
}
func TestMergeLifecycle(t *testing.T) {
	area := t.TempDir()
	releaseFixture(t, area)
	ar, e := CreateWork(area, "A", "", false)
	must(t, e)
	br, e := CreateWork(area, "B", "", false)
	must(t, e)
	must(t, os.WriteFile(filepath.Join(ar.Path, "a.txt"), []byte("A"), 0600))
	must(t, os.WriteFile(filepath.Join(br.Path, "b.txt"), []byte("B"), 0600))
	_, e = Commit(area, "work:A", "a", false)
	must(t, e)
	_, e = Commit(area, "work:B", "b", false)
	must(t, e)
	r, e := Merge(area, "work:B", "work:A", "")
	must(t, e)
	head := r.Artifact.Head
	f, e := scan(ar.Path, true)
	must(t, e)
	if string(f["a.txt"]) != "A" || string(f["b.txt"]) != "B" {
		t.Fatal(f)
	}
	r, e = Merge(area, "work:B", "work:A", "")
	must(t, e)
	if r.Artifact.Head != head {
		t.Fatal("repeat created event")
	}
	_, e = Restore(area, "work:A", "00001")
	must(t, e)
	if _, e = os.Stat(filepath.Join(ar.Path, "b.txt")); !os.IsNotExist(e) {
		t.Fatal("restore kept merged file")
	}
}
func TestCombinedWorkAndConflictResolution(t *testing.T) {
	area := t.TempDir()
	releaseFixture(t, area)
	a, e := CreateWork(area, "A", "", false)
	must(t, e)
	b, e := CreateWork(area, "B", "", false)
	must(t, e)
	must(t, os.WriteFile(filepath.Join(a.Path, "note.txt"), []byte("one"), 0600))
	must(t, os.WriteFile(filepath.Join(b.Path, "note.txt"), []byte("two"), 0600))
	r, e := Merge(area, "work:B", "work:A", "Combined")
	must(t, e)
	if r.Status != "conflicted" {
		t.Fatal(r)
	}
	if _, e = Freeze(area, "candidate", "bad", "work:Combined", ""); e == nil {
		t.Fatal("conflicted candidate")
	}
	must(t, os.WriteFile(filepath.Join(r.Path, "note.txt"), []byte("resolved"), 0600))
	r, e = Commit(area, "work:Combined", "resolved both", true)
	must(t, e)
	if len(r.Artifact.Conflicts) != 0 {
		t.Fatal("conflicts remain")
	}
	_, e = Freeze(area, "candidate", "good", "work:Combined", "")
	must(t, e)
	content, e := os.ReadFile(filepath.Join(a.Path, "note.txt"))
	must(t, e)
	if string(content) != "one" {
		t.Fatal("source changed")
	}
}
func TestMergeNoCommonBase(t *testing.T) {
	area := t.TempDir()
	initial(t, area, "A")
	initial(t, area, "B")
	if _, e := Merge(area, "work:A", "work:B", ""); e == nil {
		t.Fatal("unrelated roots merged")
	}
}
func TestMergeThreeGenerations(t *testing.T) {
	area := t.TempDir()
	releaseFixture(t, area)
	for _, name := range []string{"A", "B", "C", "D"} {
		r, e := CreateWork(area, name, "", false)
		must(t, e)
		must(t, os.WriteFile(filepath.Join(r.Path, name+".txt"), []byte(name), 0600))
		_, e = Commit(area, "work:"+name, name, false)
		must(t, e)
	}
	for _, pair := range [][2]string{{"A", "B"}, {"B", "C"}, {"C", "D"}} {
		_, e := Merge(area, "work:"+pair[0], "work:"+pair[1], "")
		must(t, e)
	}
	r, e := Status(area, "work:D")
	must(t, e)
	for _, name := range []string{"A", "B", "C", "D"} {
		b, e := os.ReadFile(filepath.Join(r.Path, name+".txt"))
		must(t, e)
		if string(b) != name {
			t.Fatal(name)
		}
	}
}

func TestDirtyCaptureDoesNotAliasFutureSourceCommit(t *testing.T) {
	area := t.TempDir()
	releaseFixture(t, area)
	a, e := CreateWork(area, "A", "", false)
	must(t, e)
	b, e := CreateWork(area, "B", "", false)
	must(t, e)
	must(t, os.WriteFile(filepath.Join(b.Path, "b.txt"), []byte("B"), 0600))
	first, e := Merge(area, "work:B", "work:A", "")
	must(t, e)
	again, e := Merge(area, "work:B", "work:A", "")
	must(t, e)
	if again.Artifact.Head != first.Artifact.Head {
		t.Fatal("duplicate dirty integration")
	}
	must(t, os.WriteFile(filepath.Join(b.Path, "b.txt"), []byte("NEW"), 0600))
	_, e = Commit(area, "work:B", "real source commit", false)
	must(t, e)
	r, e := Merge(area, "work:B", "work:A", "")
	must(t, e)
	data, e := os.ReadFile(filepath.Join(a.Path, "b.txt"))
	must(t, e)
	if r.Status != "conflicted" && string(data) != "NEW" {
		t.Fatal("silently skipped new source commit")
	}
}

func TestBinaryConflictRemainsReadable(t *testing.T) {
	area := t.TempDir()
	releaseFixture(t, area)
	a, e := CreateWork(area, "A", "", false)
	must(t, e)
	b, e := CreateWork(area, "B", "", false)
	must(t, e)
	must(t, os.WriteFile(filepath.Join(a.Path, "image.bin"), []byte{0xff}, 0600))
	must(t, os.WriteFile(filepath.Join(b.Path, "image.bin"), []byte{0xfe}, 0600))
	r, e := Merge(area, "work:B", "work:A", "")
	must(t, e)
	if r.Status != "conflicted" {
		t.Fatal(r)
	}
	r, e = Status(area, "work:A")
	must(t, e)
	if r.Artifact.Conflicts[0].Theirs == nil || *r.Artifact.Conflicts[0].Theirs != "/g==" {
		t.Fatal("wrong binary payload")
	}
	_, e = Restore(area, "work:A", "00000")
	must(t, e)
}
