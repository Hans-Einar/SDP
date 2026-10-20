package governancecontroller

import (
	"context"
	"errors"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
	"testing"
)

type fakeHost struct {
	t                          *testing.T
	c                          *Controller
	runID, mode                string
	starts, steers, interrupts int
	lost                       bool
}

func (f *fakeHost) AccountMode(context.Context) (string, error) { return f.mode, nil }
func (f *fakeHost) StartWorkThread(context.Context, string, string, string, string, map[string]any) (governancehost.Thread, error) {
	r, e := f.c.Service.Read(f.runID)
	if e != nil {
		f.t.Fatal(e)
	}
	if r.Run.Inputs[len(r.Run.Inputs)-1].State != "sending" {
		f.t.Fatal("host effect preceded durable intent")
	}
	v, e := f.c.Service.OwnerInput(governance.LocalController(), f.runID, r.Run.Inputs[len(r.Run.Inputs)-1].ID)
	if e != nil || v.Text != "Fix the selected fixture" || v.Payload == "" {
		f.t.Fatal(v, e)
	}
	thread := governancehost.Thread{ID: "thread"}
	thread.Status.Type = "idle"
	return thread, nil
}
func (f *fakeHost) StartTurn(context.Context, string, string, string) (governancehost.Turn, error) {
	f.starts++
	if f.lost {
		return governancehost.Turn{}, errors.New("response lost after possible execution")
	}
	return governancehost.Turn{ID: "turn", Status: "inProgress"}, nil
}
func (f *fakeHost) SteerTurn(context.Context, string, string, string, string) error {
	f.steers++
	return nil
}
func (f *fakeHost) InterruptTurn(context.Context, string, string) error { f.interrupts++; return nil }
func setup(t *testing.T) (*Controller, *fakeHost, string) {
	s, _, q := fixture(t)
	bind(t, s, &q)
	r, e := s.Start(governance.LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	c := &Controller{Service: s, Root: "fixture", StateDir: "private", MCPBinary: "mcp", Epoch: governancehost.UUID(), Sandbox: "read-only"}
	f := &fakeHost{t: t, c: c, runID: r.Run.ID, mode: "chatgpt"}
	c.Host = f
	return c, f, r.Run.ID
}
func TestOwnerJourneyPreservesIdentityAndStopsScope(t *testing.T) {
	c, f, id := setup(t)
	r, e := c.Launch(context.Background(), id, "Fix the selected fixture")
	if e != nil {
		t.Fatal(e)
	}
	session := r.Run.Selection.Session
	if f.starts != 1 || r.Run.State != "running" || r.Run.Inputs[0].State != "acknowledged" || r.Run.Inputs[0].Text != "" {
		t.Fatal(r)
	}
	r, e = c.Control(context.Background(), id, "steer", "Keep the original scope; add a negative test")
	if e != nil {
		t.Fatal(e)
	}
	if f.steers != 1 || len(r.Run.Inputs) != 2 {
		t.Fatal(r)
	}
	r, e = c.Control(context.Background(), id, "scope-change", "Owner proposes a separate broader feature")
	if e != nil {
		t.Fatal(e)
	}
	if f.interrupts != 1 || r.Run.State != "replan-required" || r.Run.Selection.Session != session {
		t.Fatal(r)
	}
	if _, e = c.Control(context.Background(), id, "steer", "continue without replanning"); e == nil {
		t.Fatal("scope expanded without new coverage")
	}
	if _, e = c.Observe(id, governancehost.Event{ThreadID: "other", Observation: governance.Observation{Epoch: c.Epoch}}); e == nil {
		t.Fatal("cross-thread event accepted")
	}
}
func TestUncertainSendIsNotResubmitted(t *testing.T) {
	c, f, id := setup(t)
	f.lost = true
	if _, e := c.Launch(context.Background(), id, "Fix the selected fixture"); e == nil {
		t.Fatal("lost response accepted")
	}
	r, e := c.Service.Read(id)
	if e != nil || r.Run.Inputs[0].State != "uncertain" {
		t.Fatal(r, e)
	}
	for i := 0; i < 2; i++ {
		if _, e = c.Launch(context.Background(), id, "Fix the selected fixture"); e == nil {
			t.Fatal("uncertain launch repeated")
		}
	}
	if f.starts != 1 {
		t.Fatal("duplicate host submission", f.starts)
	}
	if _, e = c.Service.OwnerInput(governance.Caller{}, id, r.Run.Inputs[0].ID); e == nil {
		t.Fatal("untrusted input read")
	}
}
func TestManagedAccountAndControllerLease(t *testing.T) {
	c, f, id := setup(t)
	f.mode = "apiKey"
	if _, e := c.Launch(context.Background(), id, "Fix the selected fixture"); e == nil {
		t.Fatal("billing fallback")
	}
	r, _ := c.Service.Read(id)
	if len(r.Run.Inputs) != 0 || f.starts != 0 {
		t.Fatal("unmanaged request dispatched")
	}
	unlock, e := c.Service.ControllerLease(governance.LocalController())
	if e != nil {
		t.Fatal(e)
	}
	defer unlock()
	f.mode = "chatgpt"
	if _, e = c.Launch(context.Background(), id, "Fix the selected fixture"); e == nil {
		t.Fatal("competing controller launched")
	}
}

func (f *fakeHost) MCPTools(context.Context, string, string) ([]string, string, error) {
	return []string{"attention_read", "attention_request", "context_read", "events_read", "route_evaluate", "run_read", "run_transition"}, "unsupported", nil
}
