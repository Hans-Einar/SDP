package governance

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

func TestRunLeasesPreserveLegacyExclusion(t *testing.T) {
	s, _, q := fixture(t)
	bind(t, s, &q)
	a, e := s.Start(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	q.OperationID, _ = identifier()
	q.RequestID, _ = identifier()
	q.Assignment.WritablePaths = []string{"second.txt"}
	b, e := s.Start(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	for _, kind := range []string{"controller", "live-controller"} {
		t.Run(kind, func(t *testing.T) {
			legacy := func() (func(), error) { return lock(filepath.Join(s.state, kind)) }
			acquire := func(id string) (func(), error) { return s.runLease(LocalController(), id, kind) }
			old, e := legacy()
			if e != nil {
				t.Fatal(e)
			}
			_, e = acquire(a.Run.ID)
			code(t, e, "busy")
			old()
			one, e := acquire(a.Run.ID)
			if e != nil {
				t.Fatal(e)
			}
			_, e = legacy()
			code(t, e, "busy")
			_, e = acquire(a.Run.ID)
			code(t, e, "busy")
			two, e := acquire(b.Run.ID)
			if e != nil {
				t.Fatal(e)
			}
			one()
			_, e = legacy()
			code(t, e, "busy")
			two()
			// Failed same-run acquisition must release its shared compatibility lock.
			old, e = legacy()
			if e != nil {
				t.Fatal("leaked compatibility lock", e)
			}
			old()
			one, e = acquire(a.Run.ID)
			if e != nil {
				t.Fatal(e)
			}
			one()
		})
	}
}

func TestRunLeaseRefusalsDoNotCreateDirectories(t *testing.T) {
	s, _, q := fixture(t)
	bind(t, s, &q)
	before := inventoryDesign(t, filepath.Dir(s.root))
	_, e := s.RunControllerLease(Caller{}, q.RequestID)
	code(t, e, "unauthorized")
	_, e = s.RunLiveControllerLease(LocalController(), "../../other")
	code(t, e, "invalid_request")
	_, e = s.RunControllerLease(LocalController(), q.RequestID)
	code(t, e, "missing_coverage")
	if !reflect.DeepEqual(before, inventoryDesign(t, filepath.Dir(s.root))) {
		t.Fatal("refusal mutated state")
	}
}

func TestRunLeaseProcessExitReleasesBothLocks(t *testing.T) {
	if os.Getenv("PG_LEASE_CHILD") == "1" {
		s, e := New(os.Getenv("PG_LEASE_ROOT"), os.Getenv("PG_LEASE_AREA"), "trial", os.Getenv("PG_LEASE_STATE"))
		if e != nil {
			t.Fatal(e)
		}
		_, e = s.RunLiveControllerLease(LocalController(), os.Getenv("PG_LEASE_RUN"))
		if e != nil {
			t.Fatal(e)
		}
		os.Stdout.WriteString("locked\n")
		var b [1]byte
		os.Stdin.Read(b[:])
		os.Exit(0)
	}
	s, _, q := fixture(t)
	bind(t, s, &q)
	r, e := s.Start(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestRunLeaseProcessExitReleasesBothLocks$")
	cmd.Env = append(os.Environ(), "PG_LEASE_CHILD=1", "PG_LEASE_ROOT="+s.root, "PG_LEASE_AREA="+s.area, "PG_LEASE_STATE="+s.state, "PG_LEASE_RUN="+r.Run.ID)
	stdin, e := cmd.StdinPipe()
	if e != nil {
		t.Fatal(e)
	}
	stdout, e := cmd.StdoutPipe()
	if e != nil {
		t.Fatal(e)
	}
	if e = cmd.Start(); e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() {
		stdin.Close()
		if cmd.ProcessState == nil {
			cmd.Process.Kill()
			cmd.Wait()
		}
	})
	line, e := bufio.NewReader(stdout).ReadString('\n')
	if e != nil || line != "locked\n" {
		t.Fatal(line, e)
	}
	_, e = s.RunLiveControllerLease(LocalController(), r.Run.ID)
	code(t, e, "busy")
	_, e = s.LiveControllerLease(LocalController())
	code(t, e, "busy")
	stdin.Close()
	if e = cmd.Wait(); e != nil {
		t.Fatal(e)
	}
	release, e := s.LiveControllerLease(LocalController())
	if e != nil {
		t.Fatal(e)
	}
	release()
	release, e = s.RunLiveControllerLease(LocalController(), r.Run.ID)
	if e != nil {
		t.Fatal(e)
	}
	release()
}
