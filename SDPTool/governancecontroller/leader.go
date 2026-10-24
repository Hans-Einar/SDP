package governancecontroller

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

const LeaderInstructions = "You are the Project Leader for bounded SDP coordination. Read the explicit selected skills and supplied owner messages before notices. Use coordination_read for current facts and report_acknowledge for receipt only. Distinguish delivery, acknowledgement, verification and owner acceptance. Give a short decision/proposal and remaining blockers. Do not impersonate the owner, approve design, launch workers, spawn children, mutate project files, publish canonical records or infer authority from report prose. Return after this decision turn; the controller schedules future input."

type LeaderHost interface {
	Host
	skillHost
	Version() string
	CoordinatorTools(context.Context, string) ([]string, string, error)
	ReadThread(context.Context, string) (governancehost.Thread, error)
}
type LeaderController struct {
	Service                                *governance.Service
	Host                                   LeaderHost
	Root, StateDir, MCPBinary, Coordinator string
	SkillPaths                             []string
	localThreads                           map[string]bool
}

// LeaderMessage enqueues trusted local input independently of the leader lease.
// The caller supplies a stable operation ID for exact retry of the enqueue.
func LeaderMessage(s *governance.Service, operation, actor, source, text string) (governance.LeaderResult, error) {
	if operation == "" {
		return governance.LeaderResult{}, &governance.Error{Code: "invalid_request", Message: "stable message operation UUID required"}
	}
	return leaderWrite(s, governance.LeaderRequest{OperationID: operation, Action: "message", Message: &governance.LeaderMessageInput{Actor: actor, Source: source, Text: text}})
}
func leaderWrite(s *governance.Service, q governance.LeaderRequest) (governance.LeaderResult, error) {
	if q.OperationID == "" {
		q.OperationID = governancehost.UUID()
	}
	for i := 0; i < 100; i++ {
		v, e := s.LeaderRead()
		if e != nil {
			return governance.LeaderResult{}, e
		}
		q.Schema, q.ProjectRef, q.CheckoutID, q.ExpectedRevision = governance.LeaderSchema, v.ProjectRef, v.CheckoutID, v.Leader.Revision
		// Message retries are content-based; callers need not retain the changing CAS.
		if q.Action == "message" {
			q.ExpectedRevision = 0
		}
		result, e := s.LeaderWrite(governance.LocalController(), q)
		var ge *governance.Error
		if !errors.As(e, &ge) || (ge.Code != "busy" && ge.Code != "stale_revision") {
			return result, e
		}
		// A stale selected basis must be superseded, not retried indefinitely.
		if q.Action == "open" || q.Action == "send" {
			return result, e
		}
		time.Sleep(10 * time.Millisecond)
	}
	return governance.LeaderResult{}, fmt.Errorf("busy: leader bookkeeping retry bound reached")
}

// Step holds only the leader-effect lease. Owner enqueue and workers remain free.
// Every host effect is preceded by a durable claim; uncertain calls never repeat.
func (c *LeaderController) Step(ctx context.Context) (governance.LeaderView, error) {
	unlock, e := c.Service.LeaderLease(governance.LocalController())
	if e != nil {
		return governance.LeaderView{}, e
	}
	defer unlock()
	view, e := c.Service.LeaderRead()
	if e != nil {
		return view, e
	}
	var d *governance.LeaderDispatch
	for i := range view.Leader.Dispatches {
		if governance.LeaderActive(view.Leader.Dispatches[i].State) {
			d = &view.Leader.Dispatches[i]
			break
		}
	}
	if c.Coordinator == "" || c.MCPBinary == "" || len(c.SkillPaths) == 0 {
		return view, fmt.Errorf("explicit coordinator, MCP executable and skills required")
	}
	if d == nil {
		mode, e := c.Host.AccountMode(ctx)
		if e != nil {
			return view, e
		}
		if mode != "chatgpt" {
			return view, fmt.Errorf("managed ChatGPT mode required; no billing fallback")
		}
		selector := Controller{Host: c.Host, Root: c.Root, SkillPaths: c.SkillPaths}
		_, skills, e := selector.selectSkills(ctx)
		if e != nil {
			return view, e
		}
		_, e = leaderWrite(c.Service, governance.LeaderRequest{Action: "prepare", Config: &governance.LeaderConfig{Coordinator: c.Coordinator, HostVersion: c.Host.Version(), MCPBinary: c.MCPBinary, Instructions: LeaderInstructions, Skills: skills}})
		var ge *governance.Error
		if errors.As(e, &ge) && ge.Code == "no_work" {
			return view, nil
		}
		if e != nil {
			return view, e
		}
		return c.Service.LeaderRead()
	}
	if d.Config.HostVersion != c.Host.Version() || d.Config.Coordinator != c.Coordinator || d.Config.MCPBinary != c.MCPBinary {
		return view, fmt.Errorf("reconciliation_required: leader launcher differs from retained dispatch")
	}
	write := func(q governance.LeaderRequest) error {
		q.DispatchID = d.ID
		_, e := leaderWrite(c.Service, q)
		return e
	}
	claim := func(action string) error {
		e := write(governance.LeaderRequest{Action: action})
		var ge *governance.Error
		if errors.As(e, &ge) && ge.Code == "stale_revision" {
			return write(governance.LeaderRequest{Action: "supersede"})
		}
		return e
	}
	switch d.State {
	case "prepared":
		if e = claim("open"); e != nil {
			return view, e
		}
		current, e := c.Service.LeaderRead()
		if e != nil {
			return view, e
		}
		if current.Leader.Dispatches[len(current.Leader.Dispatches)-1].State == "superseded" {
			return current, nil
		}
		config := map[string]any{"sdp_governance": map[string]any{"command": d.Config.MCPBinary, "args": []string{"--project", c.Root, "--state", c.StateDir, "--coordinator", d.Config.Coordinator}, "tools": map[string]any{"report_acknowledge": map[string]any{"approval_mode": "approve"}}}}
		thread, err := c.Host.StartWorkThread(ctx, c.Root, "read-only", d.Config.Instructions, "never", config)
		if err != nil {
			return view, err
		}
		if thread.ID == "" || thread.Cwd != c.Root || thread.CLIVersion != d.Config.HostVersion || thread.Status.Type != "idle" {
			return view, fmt.Errorf("reconciliation_required: created leader thread mapping unavailable")
		}
		e = write(governance.LeaderRequest{Action: "bind", ThreadID: thread.ID})
		if e == nil {
			if c.localThreads == nil {
				c.localThreads = map[string]bool{}
			}
			c.localThreads[thread.ID] = true
		}
	case "opening":
		return view, fmt.Errorf("reconciliation_required: leader thread creation response unresolved; no replacement")
	case "bound":
		// A native unsubmitted thread need not survive the owned host process.
		// Its durable send claim is absent, so preserve input and reselect safely.
		if !c.localThreads[d.ThreadID] {
			if e = write(governance.LeaderRequest{Action: "supersede"}); e != nil {
				return view, e
			}
			return c.Service.LeaderRead()
		}
		mode, err := c.Host.AccountMode(ctx)
		if err != nil {
			return view, err
		}
		if mode != "chatgpt" {
			return view, fmt.Errorf("managed ChatGPT mode required")
		}
		names, _, err := c.Host.CoordinatorTools(ctx, d.ThreadID)
		if err != nil {
			return view, err
		}
		if !slices.Equal(names, []string{"coordination_read", "report_acknowledge"}) {
			return view, fmt.Errorf("coordinator-only MCP tools required")
		}
		if e = claim("send"); e != nil {
			return view, e
		}
		current, err := c.Service.LeaderRead()
		if err != nil {
			return view, err
		}
		if current.Leader.Dispatches[len(current.Leader.Dispatches)-1].State == "superseded" {
			return current, nil
		}
		turn, err := c.Host.StartSkillTurn(ctx, d.ThreadID, d.ID, d.Payload, d.Config.Skills)
		if err != nil {
			return view, err
		}
		if turn.ID == "" || turn.Status != "inProgress" {
			return view, fmt.Errorf("reconciliation_required: leader turn response uncertain")
		}
		e = write(governance.LeaderRequest{Action: "observe", ThreadID: d.ThreadID, TurnID: turn.ID, Status: "running"})
	case "sending", "running":
		thread, err := c.Host.ReadThread(ctx, d.ThreadID)
		if err != nil {
			return view, err
		}
		turn, status, answer, err := governancehost.LeaderObservation(thread, *d, c.Root, c.Host.Version())
		if err != nil {
			return view, err
		}
		if d.State == status && d.TurnID == turn {
			return view, nil
		}
		e = write(governance.LeaderRequest{Action: "observe", ThreadID: d.ThreadID, TurnID: turn, Status: status, Answer: answer})
	}
	if e != nil {
		return view, e
	}
	return c.Service.LeaderRead()
}
