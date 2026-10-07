package governancemcp

import (
	"os"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

func TestMCPContentionRetainsExactOperation(t *testing.T) {
	s, sel, start := fixture(t)
	bind(t, s, &start)
	r, e := s.Start(governance.LocalController(), start)
	if e != nil {
		t.Fatal(e)
	}
	at, epoch := governancehost.UUID(), governancehost.UUID()
	r, e = s.RecordHost(governance.LocalController(), governance.HostRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, OperationID: governancehost.UUID(), Action: "bind", AttemptID: at, HostVersion: governancehost.HostVersion, ThreadID: "thread", Observation: governance.Observation{Kind: "bound", Epoch: epoch, ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano), Status: "idle"}})
	if e != nil {
		t.Fatal(e)
	}
	q := governance.TransitionRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, OperationID: governancehost.UUID(), ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: "start", Reason: "fixture"}
	r, e = s.Transition(governance.LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	q.Action = "report"
	q.ExpectedRevision = r.Run.Revision
	q.OperationID = governancehost.UUID()
	server, e := New(s, r.Run.ID, at, epoch)
	if e != nil {
		t.Fatal(e)
	}
	client := connect(t, server)
	coverage, _ := s.Context(sel)
	file, e := os.OpenFile(filepath.Join(filepath.Dir(coverage.Root), "state", "lock"), os.O_RDWR, 0600)
	if e != nil {
		t.Fatal(e)
	}
	defer file.Close()
	if e = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); e != nil {
		t.Fatal(e)
	}
	released := make(chan struct{})
	go func() {
		time.Sleep(100 * time.Millisecond)
		syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		close(released)
	}()
	first := call(t, client, "run_transition", q, false)
	<-released
	again := call(t, client, "run_transition", q, false)
	if !reflect.DeepEqual(first.StructuredContent, again.StructuredContent) {
		t.Fatal("receipt changed on retry")
	}
	after, e := s.Read(r.Run.ID)
	if e != nil || after.Run.State != "verifying" || after.Run.Revision != r.Run.Revision+1 {
		t.Fatal(after, e)
	}
}
