// Package governancecontroller coordinates trusted owner input and observed host
// effects. It never treats a host response as SDP evidence acceptance.
package governancecontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

type Host interface {
	AccountMode(context.Context) (string, error)
	MCPTools(context.Context, string, string) ([]string, string, error)
	StartWorkThread(context.Context, string, string, string, string, map[string]any) (governancehost.Thread, error)
	StartTurn(context.Context, string, string, string) (governancehost.Turn, error)
	SteerTurn(context.Context, string, string, string, string) error
	InterruptTurn(context.Context, string, string) error
}

type Controller struct {
	SkillPaths                                                         []string
	Service                                                            *governance.Service
	Host                                                               Host
	Root, StateDir, MCPBinary, Epoch, Sandbox, Approvals, HostIdentity string
}

// Retry only bookkeeping rejected before mutation. Never retry a host effect.
func retryRecord(f func() (governance.Result, error)) (r governance.Result, e error) {
	for i := 0; i < 100; i++ {
		r, e = f()
		if e == nil {
			return r, nil
		}
		ge, ok := e.(*governance.Error)
		if !ok || (ge.Code != "busy" && ge.Code != "stale_revision") {
			return r, e
		}
		if i < 99 {
			time.Sleep(20 * time.Millisecond)
		}
	}
	return r, e
}
func (c *Controller) input(runID, action string, v governance.OwnerInput) (governance.Result, error) {
	return retryRecord(func() (governance.Result, error) {
		r, e := c.Service.Read(runID)
		if e != nil {
			return r, e
		}
		return c.Service.RecordInput(governance.LocalController(), governance.InputRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: runID, OperationID: governancehost.UUID(), ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: action, Input: v})
	})
}
func (c *Controller) transition(runID, action, reason string) (governance.Result, error) {
	r, e := c.Service.Read(runID)
	if e != nil {
		return r, e
	}
	return c.Service.Transition(governance.LocalController(), governance.TransitionRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: runID, OperationID: governancehost.UUID(), ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: action, Reason: reason})
}
func (c *Controller) record(runID, attemptID, threadID, action string, o governance.Observation) (governance.Result, error) {
	return retryRecord(func() (governance.Result, error) {
		r, e := c.Service.Read(runID)
		if e != nil {
			return r, e
		}
		return c.Service.RecordHost(governance.LocalController(), governance.HostRequest{HostIdentity: c.HostIdentity, Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: runID, OperationID: governancehost.UUID(), ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: action, AttemptID: attemptID, HostVersion: governancehost.HostVersion, ThreadID: threadID, Observation: o})
	})
}
func now() string { return time.Now().UTC().Format(time.RFC3339Nano) }

// Launch consumes a new intent once while holding the controller lease. Failures
// after claim remain uncertain, even if an RPC may have succeeded. No retry sends.
func (c *Controller) Launch(ctx context.Context, runID, text string) (result governance.Result, err error) {
	unlock, e := c.Service.ControllerLease(governance.LocalController())
	if e != nil {
		return result, e
	}
	defer unlock()
	if c.MCPBinary == "" || (c.Sandbox != "read-only" && c.Sandbox != "workspace-write") {
		return result, fmt.Errorf("explicit MCP executable and sandbox required")
	}
	mode, e := c.Host.AccountMode(ctx)
	if e != nil {
		return result, e
	}
	if mode != "chatgpt" {
		return result, fmt.Errorf("managed ChatGPT mode required; no billing fallback")
	}
	r, e := c.Service.Read(runID)
	if e != nil {
		return r, e
	}
	coverage, e := c.Service.RunContext(runID)
	if e != nil {
		return r, e
	}
	inventory, selected, e := c.selectSkills(ctx)
	if e != nil {
		return r, e
	}
	instructions := "You are the fresh Master for one explicitly authorized SDP Maintenance assignment. Read the selected Session, authority, plan, card, AGENTS and relevant SDP skills before acting. Preserve writable/excluded boundaries and unrelated work. Native delegation is not granted by this launch; do not spawn children. Report candidate, tests, remaining gaps and observed outcome. Agent claims do not grant owner approval or independent review. Do not publish canonical management/Session records or mark the work complete; evidence publication is a later controller capability. Use assignment-bound MCP reads and report/fail/scope-change only."
	approvedDesign, e := c.Service.RunDesign(runID)
	if e != nil {
		return r, e
	}
	if approvedDesign != nil {
		instructions = strings.ReplaceAll(instructions, "SDP Maintenance assignment", "SDP implementation assignment")
		instructions = "You are the Master for the exact approved SDP implementation package. Use the included design basis, task constraints, Session and ImplementationPlan. Do not change NOW/TARGET, broaden scope or treat host completion as approval. " + instructions
	}
	payload, e := json.Marshal(map[string]any{"approvedDesign": approvedDesign, "ownerInput": text, "assignment": r.Run.Assignment, "context": coverage, "runId": runID, "instructions": instructions, "sandbox": c.Sandbox, "approvalMode": c.Approvals})
	if e != nil {
		return r, e
	}
	wireText := string(payload)
	if len(selected) > 0 {
		markers := ""
		for _, skill := range selected {
			markers += "$" + skill.Name + " "
		}
		wireText = markers + "\n" + wireText
	}
	v := governance.OwnerInput{Skills: selected, ID: governancehost.UUID(), Kind: "launch", Text: text, Payload: wireText, AttemptID: governancehost.UUID(), Epoch: c.Epoch, ReceivedAt: now()}
	if _, e = c.input(runID, "prepare", v); e != nil {
		return r, e
	}
	v.State = "prepared"
	if _, e = c.input(runID, "claim", v); e != nil {
		return r, e
	}
	v.State = "sending"
	defer func() {
		if err != nil {
			_, _ = c.input(runID, "uncertain", v)
		}
	}()
	config := map[string]any{"sdp_governance": map[string]any{"command": c.MCPBinary, "args": []string{"--project", c.Root, "--state", c.StateDir, "--run", runID, "--attempt", v.AttemptID, "--epoch", c.Epoch}}}
	thread, e := c.Host.StartWorkThread(ctx, c.Root, c.Sandbox, instructions, c.Approvals, config)
	if e != nil {
		return r, e
	}
	v.ThreadID = thread.ID
	r, e = c.record(runID, v.AttemptID, thread.ID, "bind", governance.Observation{Kind: "bound", Status: "idle", Epoch: c.Epoch, ReceivedAt: now()})
	if e != nil {
		return r, e
	}
	// Record exactly what the controller and host can substantiate. A discovery
	// result is not skill activation, and an assigned role is not a native role.
	provenance := governance.ContextProvenance{AssignedRole: "Master", RoleBasis: "controller-launch", InputRevision: r.Run.InputRevision, Sources: coverage.Sources, Host: thread.Context, Skills: governance.SkillInventory{Status: "unavailable", UseEvidence: "not-observed", Project: []governance.SkillAvailability{}}}
	if provenance.Host.Status == "" {
		provenance.Host = governance.HostContext{Status: "unavailable", InstructionStatus: "unavailable", InstructionSources: []string{}}
	}
	provenance.Skills = inventory
	encoded, err := json.Marshal(provenance)
	if err != nil {
		return r, err
	}
	r, e = c.record(runID, v.AttemptID, thread.ID, "observe", governance.Observation{Kind: "context.observed", ItemID: "launch", Text: string(encoded), Epoch: c.Epoch, ReceivedAt: now()})
	if e != nil {
		return r, e
	}
	names, _, e := c.Host.MCPTools(ctx, "sdp_governance", thread.ID)
	if e != nil {
		return r, e
	}
	if !slices.Equal(names, []string{"context_read", "events_read", "route_evaluate", "run_read", "run_transition"}) {
		return r, fmt.Errorf("assignment-bound MCP tools unavailable")
	}
	r, e = c.transition(runID, "start", "owner selected fresh bounded Master attempt")
	if e != nil {
		return r, e
	}
	if e = governance.CheckSkillSnapshots(c.Root, v.Skills); e != nil {
		return r, e
	}
	if e = c.Service.CheckRunReady(runID); e != nil {
		return r, e
	}
	var turn governancehost.Turn
	if len(v.Skills) > 0 {
		turn, e = c.Host.(skillHost).StartSkillTurn(ctx, thread.ID, v.ID, v.Payload, v.Skills)
	} else {
		turn, e = c.Host.StartTurn(ctx, thread.ID, v.ID, v.Payload)
	}
	if e != nil {
		return r, e
	}
	v.TurnID = turn.ID
	r, e = c.record(runID, v.AttemptID, thread.ID, "observe", governance.Observation{Kind: "turn.started", TurnID: turn.ID, Status: "running", Epoch: c.Epoch, ReceivedAt: now()})
	if e != nil {
		return r, e
	}
	return c.input(runID, "ack", v)
}

// Control records exact owner input before steering or interrupting. A scope
// change stops the routine before attempting a host interrupt; it never expands
// the existing assignment. Uncertain requests remain visible and are not resent.
func (c *Controller) Control(ctx context.Context, runID, kind, text string) (result governance.Result, err error) {
	unlock, e := c.Service.ControllerLease(governance.LocalController())
	if e != nil {
		return result, e
	}
	defer unlock()
	if kind != "steer" && kind != "interrupt" && kind != "scope-change" {
		return result, fmt.Errorf("unsupported owner command")
	}
	r, e := c.Service.Read(runID)
	if e != nil {
		return r, e
	}
	if len(r.Run.Attempts) == 0 {
		return r, fmt.Errorf("no mapped attempt")
	}
	a := r.Run.Attempts[len(r.Run.Attempts)-1]
	if a.Epoch != c.Epoch {
		return r, fmt.Errorf("reconciliation_required: controller epoch changed")
	}
	v := governance.OwnerInput{ID: governancehost.UUID(), Kind: kind, Text: text, Payload: text, AttemptID: a.ID, Epoch: c.Epoch, ThreadID: a.ThreadID, TurnID: a.TurnID, ReceivedAt: now()}
	if _, e = c.input(runID, "prepare", v); e != nil {
		return r, e
	}
	v.State = "prepared"
	if _, e = c.input(runID, "claim", v); e != nil {
		return r, e
	}
	v.State = "sending"
	defer func() {
		if err != nil {
			_, _ = c.input(runID, "uncertain", v)
		}
	}()
	if kind == "steer" {
		e = c.Host.SteerTurn(ctx, a.ThreadID, a.TurnID, v.ID, text)
	} else {
		action := "interrupt"
		if kind == "scope-change" {
			action = kind
		}
		if _, e = c.transition(runID, action, text); e != nil {
			return r, e
		}
		e = c.Host.InterruptTurn(ctx, a.ThreadID, a.TurnID)
	}
	if e != nil {
		return r, e
	}
	return c.input(runID, "ack", v)
}
func (c *Controller) Observe(runID string, event governancehost.Event) (governance.Result, error) {
	r, e := c.Service.Read(runID)
	if e != nil {
		return r, e
	}
	if len(r.Run.Attempts) == 0 {
		return r, fmt.Errorf("no mapped attempt")
	}
	a := r.Run.Attempts[len(r.Run.Attempts)-1]
	if event.ThreadID != a.ThreadID || event.Observation.Epoch != c.Epoch || a.Epoch != c.Epoch {
		return r, fmt.Errorf("observation outside current attempt")
	}
	return c.record(runID, a.ID, a.ThreadID, "observe", event.Observation)
}
func (c *Controller) Disconnect(runID string) (governance.Result, error) {
	r, e := c.Service.Read(runID)
	if e != nil {
		return r, e
	}
	if len(r.Run.Attempts) == 0 {
		return r, nil
	}
	a := r.Run.Attempts[len(r.Run.Attempts)-1]
	return c.record(runID, a.ID, a.ThreadID, "disconnect", governance.Observation{Kind: "disconnected", Status: "unknown", Epoch: c.Epoch, ReceivedAt: now()})
}
