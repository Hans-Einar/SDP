package model

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCommitRestoreDirtyDelete(t *testing.T) {
	area := t.TempDir()
	p := initial(t, area, "A")
	must(t, os.WriteFile(filepath.Join(p, "old.design"), []byte("old"), 0600))
	_, e := Commit(area, "work:A", "first", false)
	must(t, e)
	must(t, os.Remove(filepath.Join(p, "old.design")))
	must(t, os.WriteFile(filepath.Join(p, "new.design"), []byte("new"), 0600))
	_, e = Commit(area, "work:A", "rename", false)
	must(t, e)
	must(t, os.WriteFile(filepath.Join(p, "new.design"), []byte("dirty"), 0600))
	r, e := Restore(area, "work:A", "00001")
	must(t, e)
	if r.Artifact.Sequence != 4 {
		t.Fatal(r.Artifact.Sequence)
	}
	b, e := os.ReadFile(filepath.Join(p, "old.design"))
	must(t, e)
	if string(b) != "old" {
		t.Fatal(string(b))
	}
	if _, e = os.Stat(filepath.Join(p, "new.design")); !os.IsNotExist(e) {
		t.Fatal("deleted file retained")
	}
	f, e := reconstruct(p, *r.Artifact, r.Artifact.ID+":00003")
	must(t, e)
	if string(f["new.design"]) != "dirty" {
		t.Fatal("dirty content lost")
	}
}
func TestInterruptedProcesses(t *testing.T) {
	if root := os.Getenv("MGI_TEST_AREA"); root != "" {
		phase := os.Getenv("MGI_TEST_PHASE")
		faultHook = func(s string) {
			if s == phase {
				os.Exit(71)
			}
		}
		var e error
		if os.Getenv("MGI_TEST_RESTORE") == "1" {
			_, e = Restore(root, "work:A", "00000")
		} else {
			_, e = Commit(root, "work:A", "child", false)
		}
		if e != nil {
			os.Exit(72)
		}
		os.Exit(73)
	}
	for _, phase := range []string{"building", "prepared", "backup", "installed"} {
		t.Run(phase, func(t *testing.T) {
			area := t.TempDir()
			p := initial(t, area, "A")
			must(t, os.WriteFile(filepath.Join(p, "change"), []byte("new"), 0600))
			cmd := exec.Command(os.Args[0], "-test.run=^TestInterruptedProcesses$")
			cmd.Env = append(os.Environ(), "MGI_TEST_AREA="+area, "MGI_TEST_PHASE="+phase)
			e := cmd.Run()
			if x, ok := e.(*exec.ExitError); !ok || x.ExitCode() != 71 {
				t.Fatalf("bad child result %v", e)
			}
			ts, e := journals(area)
			must(t, e)
			if len(ts) != 1 {
				t.Fatal(ts)
			}
			if _, e = CreateWork(area, "Blocked", "", true); e == nil {
				t.Fatal("pending ignored")
			}
			action := "resume"
			if phase == "building" {
				action = "abort"
			}
			_, e = Recover(area, ts[0].ID, action)
			must(t, e)
			r, e := Status(area, "work:A")
			must(t, e)
			if phase != "building" && (r.Artifact.Sequence != 1 || r.Status != "clean") {
				t.Fatal(r)
			}
			b, e := os.ReadFile(filepath.Join(p, "change"))
			must(t, e)
			if string(b) != "new" {
				t.Fatal("lost source")
			}
		})
	}
}
func TestRecoveryRefusesEditedTarget(t *testing.T) {
	area := t.TempDir()
	p := initial(t, area, "A")
	must(t, os.WriteFile(filepath.Join(p, "x"), []byte("one"), 0600))
	faultHook = func(s string) {
		if s == "prepared" {
			panic("stop")
		}
	}
	func() {
		defer func() { _ = recover(); faultHook = nil }()
		_, _ = Commit(area, "work:A", "stop", false)
	}()
	must(t, os.WriteFile(filepath.Join(p, "x"), []byte("external"), 0600))
	ts, e := journals(area)
	must(t, e)
	if _, e = Recover(area, ts[0].ID, "resume"); e == nil {
		t.Fatal("overwrote external edits")
	}
	b, e := os.ReadFile(filepath.Join(p, "x"))
	must(t, e)
	if string(b) != "external" {
		t.Fatal("lost external edits")
	}
}
func TestLocalWriterExclusion(t *testing.T) {
	area := t.TempDir()
	unlock, e := lockArea(area)
	must(t, e)
	defer unlock()
	if _, e = CreateWork(area, "A", "", true); e == nil {
		t.Fatal("writer contention ignored")
	}
}
func TestRestoreMissingPayload(t *testing.T) {
	area := t.TempDir()
	p := initial(t, area, "A")
	must(t, os.WriteFile(filepath.Join(p, "empty"), nil, 0600))
	r, e := Commit(area, "work:A", "empty", false)
	must(t, e)
	must(t, os.Remove(filepath.Join(p, ".commits/#00001/files/empty")))
	if _, e = reconstruct(p, *r.Artifact, r.Artifact.Head); e == nil {
		t.Fatal("missing empty file accepted")
	}
}

func TestRestoreFileDirectoryTransitions(t *testing.T) {
	area := t.TempDir()
	p := initial(t, area, "A")
	must(t, os.WriteFile(filepath.Join(p, "x"), []byte("file"), 0600))
	_, e := Commit(area, "work:A", "file", false)
	must(t, e)
	must(t, os.Remove(filepath.Join(p, "x")))
	must(t, os.Mkdir(filepath.Join(p, "x"), 0700))
	must(t, os.WriteFile(filepath.Join(p, "x/y"), []byte("child"), 0600))
	_, e = Commit(area, "work:A", "directory", false)
	must(t, e)
	_, e = Restore(area, "work:A", "00001")
	must(t, e)
	_, e = Restore(area, "work:A", "00002")
	must(t, e)
	b, e := os.ReadFile(filepath.Join(p, "x/y"))
	must(t, e)
	if string(b) != "child" {
		t.Fatal("wrong restore")
	}
}

func TestCrossProcessWriter(t *testing.T) {
	if root := os.Getenv("MGI_LOCK_TEST"); root != "" {
		_, e := CreateWork(root, "A", "", true)
		if e == nil {
			os.Exit(9)
		}
		os.Exit(0)
	}
	area := t.TempDir()
	unlock, e := lockArea(area)
	must(t, e)
	cmd := exec.Command(os.Args[0], "-test.run=^TestCrossProcessWriter$")
	cmd.Env = append(os.Environ(), "MGI_LOCK_TEST="+area)
	must(t, cmd.Run())
	unlock()
	initial(t, area, "A")
}
