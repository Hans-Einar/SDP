// Protocol test host for the real controller. No Codex/model execution occurs.
package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancecontroller"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

type host struct {
	s         *governance.Service
	run, root string
	starts    int
	lost      bool
}

func (h *host) AccountMode(context.Context) (string, error) { return "chatgpt", nil }
func (h *host) Version() string                             { return "0.164.0" }
func (h *host) MCPTools(context.Context, string, string) ([]string, string, error) {
	return []string{"attention_read", "attention_request", "context_read", "events_read", "route_evaluate", "run_read", "run_transition"}, "fixture", nil
}
func (h *host) StartWorkThread(_ context.Context, root, sandbox, instructions, approval string, config map[string]any) (governancehost.Thread, error) {
	r, e := h.s.Read(h.run)
	if e != nil {
		return governancehost.Thread{}, e
	}
	if r.Run.State != "ready" || len(r.Run.Inputs) != 1 || r.Run.Inputs[0].State != "sending" {
		return governancehost.Thread{}, errors.New("host preceded committed run/intent")
	}
	b, e := os.ReadFile(filepath.Join(h.root, r.Run.Selection.Card))
	if e != nil || !strings.Contains(string(b), "| CardState | in-progress |") {
		return governancehost.Thread{}, errors.New("host preceded in-progress card")
	}
	if !strings.Contains(instructions, "Session assignment") {
		return governancehost.Thread{}, errors.New("wrong assignment instructions")
	}
	thread := governancehost.Thread{ID: "session-controller-thread"}
	thread.Status.Type = "idle"
	return thread, nil
}
func (h *host) StartTurn(context.Context, string, string, string) (governancehost.Turn, error) {
	return governancehost.Turn{}, errors.New("explicit skills required")
}
func (h *host) StartSkillTurn(_ context.Context, thread, id, text string, skills []governance.SkillInput) (governancehost.Turn, error) {
	h.starts++
	if len(skills) != 2 {
		return governancehost.Turn{}, errors.New("missing skills")
	}
	v, e := h.s.OwnerInput(governance.LocalController(), h.run, id)
	if e != nil || v.State != "sending" || len(v.Skills) != 2 {
		return governancehost.Turn{}, errors.New("input not durably bound")
	}
	_, body, ok := strings.Cut(text, "\n")
	if !ok {
		return governancehost.Turn{}, errors.New("skill markers absent")
	}
	var p struct {
		SessionWork *governance.SessionView `json:"sessionWork"`
	}
	if json.Unmarshal([]byte(body), &p) != nil || p.SessionWork == nil || p.SessionWork.Binding == nil || p.SessionWork.Steps[0].RunID != h.run || p.SessionWork.Steps[0].Assessment == nil {
		return governancehost.Turn{}, errors.New("Session context missing from payload")
	}
	if h.lost {
		return governancehost.Turn{}, errors.New("simulated lost start reply")
	}
	return governancehost.Turn{ID: "session-controller-turn", Status: "inProgress"}, nil
}
func (h *host) SteerTurn(context.Context, string, string, string, string) error { return nil }
func (h *host) InterruptTurn(context.Context, string, string) error             { return nil }
func (h *host) InspectSkills(context.Context, string) (governance.SkillInventory, error) {
	v := governance.SkillInventory{Status: "observed", UseEvidence: "not-observed", Project: []governance.SkillAvailability{}}
	for _, name := range []string{"sdp-master", "sdp-worker"} {
		p := "Skills/" + name + "/SKILL.md"
		b, e := os.ReadFile(filepath.Join(h.root, p))
		if e != nil {
			return v, e
		}
		v.Project = append(v.Project, governance.SkillAvailability{Name: name, Path: p, Enabled: true, Digest: fmt.Sprintf("%x", sha256.Sum256(b))})
	}
	return v, nil
}
func main() {
	if e := run(); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
func run() error {
	if len(os.Args) != 5 {
		return errors.New("root state run lost|success required")
	}
	root, state, id := os.Args[1], os.Args[2], os.Args[3]
	s, e := governance.New(root, filepath.Join(root, "SDP"), "trial", state)
	if e != nil {
		return e
	}
	h := &host{s: s, run: id, root: root, lost: os.Args[4] == "lost"}
	c := governancecontroller.Controller{Service: s, Host: h, Root: root, StateDir: state, MCPBinary: "/fixture/sdp-mcp", Epoch: governancehost.UUID(), Sandbox: "workspace-write", SkillPaths: []string{"Skills/sdp-master/SKILL.md", "Skills/sdp-worker/SKILL.md"}}
	r, e := c.Launch(context.Background(), id, "Execute the exact selected Session task")
	if h.lost {
		if e == nil || h.starts != 1 {
			return errors.New("lost reply was not retained")
		}
		if _, e = c.Launch(context.Background(), id, "Retry"); e == nil || h.starts != 1 {
			return errors.New("uncertain launch resent")
		}
		v, e := s.RunSession(id)
		if e != nil || v.HostExecution != "unknown" {
			return fmt.Errorf("uncertain Session observation: %+v %v", v, e)
		}
		return json.NewEncoder(os.Stdout).Encode(v)
	}
	if e != nil {
		return e
	}
	if _, e = c.Launch(context.Background(), id, "Retry"); e == nil || h.starts != 1 {
		return errors.New("launch duplicated")
	}
	at := r.Run.Attempts[len(r.Run.Attempts)-1]
	r, e = s.RecordHost(governance.LocalController(), governance.HostRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: governancehost.UUID(), RunID: id, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: "observe", AttemptID: at.ID, HostVersion: at.HostVersion, ThreadID: at.ThreadID, Observation: governance.Observation{Kind: "turn.completed", Status: "completed", Epoch: at.Epoch, TurnID: at.TurnID, ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}})
	if e != nil {
		return e
	}
	caller, e := governance.AssignedAgent(id, at.ID, at.CapabilityEpoch())
	if e != nil {
		return e
	}
	r, e = s.Transition(caller, governance.TransitionRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: governancehost.UUID(), RunID: id, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: "report", Reason: "Protocol fixture report; no native agent/model execution"})
	if e != nil {
		return e
	}
	return json.NewEncoder(os.Stdout).Encode(r)
}
