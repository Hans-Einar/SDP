package install

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
)

func TestCrashHelper(t *testing.T) {
	if os.Getenv("SDP_TEST_CRASH_CHILD") != "1" {
		return
	}
	p, e := LoadPlan(os.Getenv("SDP_TEST_PLAN"))
	if e != nil {
		t.Fatal(e)
	}
	index, _ := strconv.Atoi(os.Getenv("SDP_TEST_INDEX"))
	x := Executor{Fault: func(name string, i int) error {
		if name == os.Getenv("SDP_TEST_BOUNDARY") && i == index {
			os.Exit(73)
		}
		return nil
	}}
	_, e = x.Apply(context.Background(), p.ProjectRoot, p)
	if e != nil {
		t.Fatal(e)
	}
	t.Fatal("fault boundary not reached")
}
func recoveryFixture(t *testing.T, kind string) (string, Plan) {
	t.Helper()
	root := t.TempDir()
	d := fixture()
	if kind == "known" {
		first := preview(t, root, d)
		if _, e := (Executor{}).Apply(context.Background(), root, first); e != nil {
			t.Fatal(e)
		}
		d.Release = "dev-next"
		d.UpgradesFrom = []string{first.Release.SHA256}
		d.Files[0].Content = []byte("next version")
		h := Hash(d.Files[0].Content)
		d.Files[0].SHA256 = &h
		p, e := Preview(Options{root, "upgrade", artifact(t, d), first.Release.Path, "", true})
		if e != nil {
			t.Fatal(e)
		}
		return root, p
	}
	if kind == "manual" {
		put(t, root, "AGENTS.md", []byte("owner instructions\n"))
		put(t, root, "SDP/Agents/KanBan/board.json", []byte(`{"schemaVersion":"0.1","projectId":"XFMD","ledger":"Ledger.ndjson"}`))
		put(t, root, "SDP/Agents/KanBan/Ledger.ndjson", historyFixture())
		put(t, root, "SDP/Agents/KanBan/backlog/#1.md", cardFixture("project card\n"))
		a := artifact(t, d)
		in, _ := LocalInput(a, true)
		m := adoption(t, root, in, []Move{{"SDP/Agents/KanBan", "SDP/KanBan"}}, []string{"AGENTS.md"})
		p, e := Preview(Options{root, "upgrade", a, "", m, true})
		if e != nil || !p.CanApply {
			t.Fatal(p, e)
		}
		return root, p
	}
	return root, preview(t, root, d)
}
func TestProcessExitRecoveryMatrix(t *testing.T) {
	for _, kind := range []string{"clean", "known", "manual"} {
		t.Run(kind, func(t *testing.T) {
			_, initial := recoveryFixture(t, kind)
			journal, e := prepare(initial)
			if e != nil {
				t.Fatal(e)
			}
			cases := []struct {
				name  string
				index int
			}{{"prepared", 0}, {"complete", len(journal.Steps)}}
			for i := range journal.Steps {
				for _, name := range []string{"backup", "write", "journal"} {
					cases = append(cases, struct {
						name  string
						index int
					}{name, i})
				}
			}
			for _, c := range cases {
				t.Run(fmt.Sprintf("%s-%d", c.name, c.index), func(t *testing.T) {
					root, p := recoveryFixture(t, kind)
					plan := filepath.Join(t.TempDir(), "plan.json")
					if e = SavePlan(p, plan); e != nil {
						t.Fatal(e)
					}
					cmd := exec.Command(os.Args[0], "-test.run=^TestCrashHelper$")
					cmd.Env = append(os.Environ(), "SDP_TEST_CRASH_CHILD=1", "SDP_TEST_PLAN="+plan, "SDP_TEST_BOUNDARY="+c.name, "SDP_TEST_INDEX="+strconv.Itoa(c.index))
					out, err := cmd.CombinedOutput()
					ee, ok := err.(*exec.ExitError)
					if !ok || ee.ExitCode() != 73 {
						t.Fatalf("child: %v %s", err, out)
					}
					pending, e := Pending(root)
					if e != nil || len(pending) != 1 {
						t.Fatal(pending, e)
					}
					jb, e := Read(filepath.Join(root, Operations, pending[0], "journal.json"), RecordLimit)
					if e != nil {
						t.Fatal(e)
					}
					var j Journal
					if e = Decode(jb, RecordLimit, &j); e != nil {
						t.Fatal(e)
					}
					steps, _ := Canonical(j.Steps)
					// Recovery relies on reserved, embedded inputs; no network or mutable artifact.
					os.Remove(p.Release.Path)
					r, e := (Executor{}).Resume(context.Background(), root, pending[0])
					if e != nil || r.Status != "completed" {
						t.Fatal(r, e)
					}
					b, _ := Read(filepath.Join(root, Operations, pending[0], "journal.json"), RecordLimit)
					var finished Journal
					if e = Decode(b, RecordLimit, &finished); e != nil {
						t.Fatal(e)
					}
					actual, _ := Canonical(finished.Steps)
					if !bytes.Equal(steps, actual) || finished.MaintenanceID != j.MaintenanceID || finished.CreatedAt != j.CreatedAt {
						t.Fatal("reserved finalization changed")
					}
					before, _ := Inspect(root, nil)
					if _, e = (Executor{}).Resume(context.Background(), root, pending[0]); e != nil {
						t.Fatal(e)
					}
					after, _ := Inspect(root, nil)
					if !SameSnapshot(before.Snapshot, after.Snapshot) {
						t.Fatal("repeat resume changed files")
					}
					if kind == "manual" {
						history, _ := os.ReadFile(filepath.Join(root, "SDP/ProjectManagement/Ledger.ndjson"))
						if !bytes.HasPrefix(history, historyFixture()) {
							t.Fatal("history prefix altered")
						}
					}
				})
			}
		})
	}
}
func TestRecoveryRejectsPostFailureEdits(t *testing.T) {
	root, p := recoveryFixture(t, "clean")
	x := Executor{Fault: func(name string, i int) error {
		if name == "write" && i == 0 {
			return fmt.Errorf("stop")
		}
		return nil
	}}
	r, e := x.Apply(context.Background(), root, p)
	if e == nil {
		t.Fatal("no interruption")
	}
	put(t, root, p.Actions[0].Path, []byte("outside edit"))
	if _, e = (Executor{}).Resume(context.Background(), root, r.OperationID); e == nil {
		t.Fatal("post-failure edit accepted")
	}
}
