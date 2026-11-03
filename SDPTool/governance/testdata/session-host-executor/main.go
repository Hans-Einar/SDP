// Synthetic host bridge for the governance package's integration tests. This
// executable never opens Codex or submits native/model work.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancecontroller"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

type host struct {
	service            *governance.Service
	runID, root        string
	skills             []governance.SkillInput
	starts, interrupts int
	lost               bool
}

func (h *host) Version() string                             { return "0.165.0" }
func (h *host) AccountMode(context.Context) (string, error) { return "chatgpt", nil }
func (h *host) MCPTools(context.Context, string, string) ([]string, string, error) {
	return []string{"attention_read", "attention_request", "context_read", "events_read", "route_evaluate", "run_read", "run_transition"}, "fixture", nil
}
func (h *host) InspectSkills(context.Context, string) (governance.SkillInventory, error) {
	v := governance.SkillInventory{Status: "observed", UseEvidence: "not-observed", Project: []governance.SkillAvailability{}}
	for _, s := range h.skills {
		p, _ := filepath.Rel(h.root, s.Path)
		v.Project = append(v.Project, governance.SkillAvailability{Name: s.Name, Path: filepath.ToSlash(p), Digest: s.Digest, Enabled: true})
	}
	return v, nil
}
func (h *host) StartWorkThread(context.Context, string, string, string, string, map[string]any) (governancehost.Thread, error) {
	r, e := h.service.Read(h.runID)
	if e != nil {
		return governancehost.Thread{}, e
	}
	if len(r.Run.Inputs) != 1 || r.Run.Inputs[0].State != "sending" {
		return governancehost.Thread{}, fmt.Errorf("effect before durable claim")
	}
	t := governancehost.Thread{ID: "synthetic-thread"}
	t.Status.Type = "idle"
	return t, nil
}
func (h *host) StartTurn(context.Context, string, string, string) (governancehost.Turn, error) {
	return governancehost.Turn{}, fmt.Errorf("explicit skill input required")
}
func (h *host) StartSkillTurn(_ context.Context, thread, id, text string, skills []governance.SkillInput) (governancehost.Turn, error) {
	h.starts++
	if h.lost {
		return governancehost.Turn{}, fmt.Errorf("lost response after possible send")
	}
	v, e := h.service.OwnerInput(governance.LocalController(), h.runID, id)
	if e != nil || v.State != "sending" || len(skills) != 2 {
		return governancehost.Turn{}, fmt.Errorf("missing exact input/skills: %v", e)
	}
	return governancehost.Turn{ID: "synthetic-turn", Status: "inProgress"}, nil
}
func (h *host) SteerTurn(context.Context, string, string, string, string) error {
	return fmt.Errorf("not allowed")
}
func (h *host) InterruptTurn(context.Context, string, string) error { h.interrupts++; return nil }
func main() {
	s, e := governance.New(os.Args[1], filepath.Join(os.Args[1], "SDP"), os.Args[3], os.Args[2])
	check(e)
	command, e := s.ReadSessionHostCommand(os.Args[4])
	check(e)
	selected, e := s.SessionHostSelection(command.Request.InstanceID)
	check(e)
	h := &host{service: s, runID: command.Request.RunID, root: os.Args[1], skills: selected.Skills, lost: os.Args[5] == "lost"}
	c := &governancecontroller.Controller{Service: s, Host: h, Root: os.Args[1], StateDir: os.Args[2], MCPBinary: "/fixture/mcp", Epoch: governancehost.UUID(), Sandbox: "read-only"}
	result, first := c.ExecuteSessionHost(context.Background(), command.Request.OperationID)
	again, second := c.ExecuteSessionHost(context.Background(), command.Request.OperationID)
	if h.starts != 1 {
		panic(fmt.Sprintf("duplicate/missing submission: %d; %v", h.starts, first))
	}
	if h.lost {
		if first == nil || second == nil || result.State != "uncertain" || again.State != "uncertain" {
			panic("lost response accepted")
		}
		if _, e = s.CancelSessionHost(governance.LocalController(), command.Request.OperationID); e == nil {
			panic("uncertain input canceled")
		}
	} else {
		check(first)
		check(second)
		if result.State != "acknowledged" || again.State != "acknowledged" {
			panic("delivery not acknowledged")
		}
		run, e := s.Read(h.runID)
		check(e)
		q := command.Request
		q.OperationID = governancehost.UUID()
		q.ExpectedRunRevision = run.Run.Revision
		q.Action = "scope-change"
		q.Reason = "New owner idea requires a separate Session"
		stop, e := s.RequestSessionHost(governance.LocalController(), q)
		check(e)
		stop, e = c.ExecuteSessionHost(context.Background(), q.OperationID)
		check(e)
		_, e = c.ExecuteSessionHost(context.Background(), q.OperationID)
		check(e)
		run, e = s.Read(h.runID)
		check(e)
		if h.interrupts != 1 || stop.State != "acknowledged" || run.Run.State != "replan-required" {
			panic("scope control failed")
		}
	}
	check(json.NewEncoder(os.Stdout).Encode(map[string]any{"starts": h.starts, "interrupts": h.interrupts, "state": result.State}))
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
