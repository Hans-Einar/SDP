package governancecontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

type leaderHost struct {
	*selectedHost
	service                   *governance.Service
	root                      string
	thread                    governancehost.Thread
	opens, sends              int
	lostOpen, lostSend, drift bool
	enqueueDuringSend         bool
}

func (h *leaderHost) StartWorkThread(_ context.Context, root, sandbox, instructions, approvals string, config map[string]any) (governancehost.Thread, error) {
	h.opens++
	v, e := h.service.LeaderRead()
	if e != nil {
		h.t.Fatal(e)
	}
	d := v.Leader.Dispatches[len(v.Leader.Dispatches)-1]
	if d.State != "opening" || sandbox != "read-only" || approvals != "never" || !strings.Contains(instructions, "Do not impersonate") {
		h.t.Fatal("effect before durable restricted intent", d)
	}
	raw, _ := json.Marshal(config)
	if !strings.Contains(string(raw), "--coordinator") || strings.Contains(string(raw), "run_transition") {
		h.t.Fatal(string(raw))
	}
	if h.lostOpen {
		return governancehost.Thread{}, fmt.Errorf("lost thread response")
	}
	h.thread = governancehost.Thread{ID: fmt.Sprintf("thread-%d", h.opens), Cwd: root, CLIVersion: h.Version()}
	h.thread.Status.Type = "idle"
	return h.thread, nil
}
func (h *leaderHost) MCPTools(context.Context, string, string) ([]string, string, error) {
	if h.drift {
		os.WriteFile(filepath.Join(h.root, "Skills/sdp/SKILL.md"), []byte("changed"), 0600)
	}
	return []string{"coordination_read", "report_acknowledge"}, "unsupported", nil
}
func (h *leaderHost) StartSkillTurn(_ context.Context, thread, id, text string, skills []governance.SkillInput) (governancehost.Turn, error) {
	h.sends++
	v, e := h.service.LeaderRead()
	if e != nil {
		h.t.Fatal(e)
	}
	d := v.Leader.Dispatches[len(v.Leader.Dispatches)-1]
	if d.State != "sending" || id != d.ID || text != d.Payload {
		h.t.Fatal("send before durable mapping", d)
	}
	input, _ := json.Marshal(map[string]any{"id": "input", "type": "userMessage", "clientId": id, "content": governance.TurnInputs(text, skills)})
	answer := json.RawMessage(`{"id":"answer","type":"agentMessage","phase":"final_answer","text":"Proposed next action"}`)
	private := json.RawMessage(`{"id":"private","type":"reasoning","text":"DO_NOT_RETAIN"}`)
	h.thread.Turns = []governancehost.Turn{{ID: "turn", Status: "completed", ItemsView: "full", Items: []json.RawMessage{input, private, answer}}}
	if h.enqueueDuringSend {
		if _, e := LeaderMessage(h.service, governancehost.UUID(), "owner", "fixture", "Next owner input"); e != nil {
			h.t.Fatal(e)
		}
	}
	if h.lostSend {
		return governancehost.Turn{}, fmt.Errorf("lost turn response")
	}
	return governancehost.Turn{ID: "turn", Status: "inProgress"}, nil
}
func (h *leaderHost) ReadThread(context.Context, string) (governancehost.Thread, error) {
	return h.thread, nil
}
func leaderSetup(t *testing.T) (*LeaderController, *leaderHost) {
	c, f, _ := skillFixture(t)
	h := &leaderHost{selectedHost: f, service: c.Service, root: c.Root}
	leader := &LeaderController{Service: c.Service, Host: h, Root: c.Root, StateDir: filepath.Join(filepath.Dir(c.Root), "state"), MCPBinary: "/fixture/mcp", Coordinator: "leader", SkillPaths: c.SkillPaths}
	if _, e := LeaderMessage(c.Service, governancehost.UUID(), "owner", "fixture", "Review work"); e != nil {
		t.Fatal(e)
	}
	return leader, h
}
func stepLeader(t *testing.T, c *LeaderController) governance.LeaderView {
	t.Helper()
	v, e := c.Step(context.Background())
	if e != nil {
		t.Fatal(e)
	}
	return v
}
func TestLeaderControllerLostSendRecoveryAndQueuedInput(t *testing.T) {
	c, h := leaderSetup(t)
	h.lostSend = true
	h.enqueueDuringSend = true
	stepLeader(t, c)
	stepLeader(t, c)
	if _, e := c.Step(context.Background()); e == nil {
		t.Fatal("lost response accepted")
	}
	v, e := c.Service.LeaderRead()
	if e != nil || v.Leader.Dispatches[0].State != "sending" || len(v.Leader.Messages) != 2 {
		t.Fatal(v, e)
	}
	// A new controller recovers only the mapped exact history, with no resend.
	restarted := *c
	v = stepLeader(t, &restarted)
	if v.Leader.Dispatches[0].State != "completed" || h.sends != 1 || h.opens != 1 || strings.Contains(v.Leader.Dispatches[0].Answer, "DO_NOT_RETAIN") {
		t.Fatal(v)
	}
	v = stepLeader(t, &restarted)
	if len(v.Leader.Dispatches) != 2 || len(v.Leader.Dispatches[1].MessageIDs) != 1 || v.Leader.Dispatches[1].MessageIDs[0] != v.Leader.Messages[1].ID {
		t.Fatal("queued input not next", v)
	}
}
func TestLeaderControllerRefusesUnmappedThreadAndStaleSkills(t *testing.T) {
	t.Run("opening", func(t *testing.T) {
		c, h := leaderSetup(t)
		h.lostOpen = true
		stepLeader(t, c)
		if _, e := c.Step(context.Background()); e == nil {
			t.Fatal("lost response accepted")
		}
		if _, e := c.Step(context.Background()); e == nil || h.opens != 1 || h.sends != 0 {
			t.Fatal("opening replaced", e)
		}
	})
	t.Run("skill-drift", func(t *testing.T) {
		c, h := leaderSetup(t)
		stepLeader(t, c)
		stepLeader(t, c)
		h.drift = true
		v := stepLeader(t, c)
		if v.Leader.Dispatches[0].State != "superseded" || h.sends != 0 {
			t.Fatal(v)
		}
	})
	t.Run("owner-priority", func(t *testing.T) {
		c, h := leaderSetup(t)
		stepLeader(t, c)
		if _, e := LeaderMessage(c.Service, governancehost.UUID(), "owner", "fixture", "Urgent correction"); e != nil {
			t.Fatal(e)
		}
		v := stepLeader(t, c)
		if v.Leader.Dispatches[0].State != "superseded" || h.opens != 0 {
			t.Fatal(v)
		}
		v = stepLeader(t, c)
		if len(v.Leader.Dispatches[1].MessageIDs) != 2 {
			t.Fatal(v)
		}
	})
}
func TestLeaderRecoveryExactHistoryRefusals(t *testing.T) {
	for _, kind := range []string{"missing", "duplicate", "wrong-input", "partial", "version", "cwd", "thread", "turn", "unknown-status"} {
		t.Run(kind, func(t *testing.T) {
			c, h := leaderSetup(t)
			stepLeader(t, c)
			stepLeader(t, c)
			stepLeader(t, c)
			switch kind {
			case "missing":
				h.thread.Turns = nil
			case "duplicate":
				h.thread.Turns[0].Items = append(h.thread.Turns[0].Items, h.thread.Turns[0].Items[0])
			case "wrong-input":
				h.thread.Turns[0].Items[0] = json.RawMessage(strings.Replace(string(h.thread.Turns[0].Items[0]), "Review work", "different work", 1))
			case "partial":
				h.thread.Turns[0].ItemsView = "summary"
			case "version":
				h.thread.CLIVersion = "0.165.0"
			case "cwd":
				h.thread.Cwd = "/other"
			case "thread":
				h.thread.ID = "other"
			case "turn":
				h.thread.Turns[0].ID = "other"
			case "unknown-status":
				h.thread.Turns[0].Status = "idle"
			}
			if _, e := c.Step(context.Background()); e == nil {
				t.Fatal("ambiguous history accepted")
			}
			v, e := c.Service.LeaderRead()
			if e != nil || v.Leader.Dispatches[0].State != "running" || h.sends != 1 {
				t.Fatal(v, e)
			}
		})
	}
}

func (h *leaderHost) CoordinatorTools(ctx context.Context, thread string) ([]string, string, error) {
	return h.MCPTools(ctx, "sdp_governance", thread)
}

func TestLeaderRestartSupersedesOnlyUnsentBoundThread(t *testing.T) {
	c, h := leaderSetup(t)
	stepLeader(t, c)
	first := stepLeader(t, c)
	restarted := *c
	restarted.localThreads = nil
	v := stepLeader(t, &restarted)
	if v.Leader.Dispatches[0].State != "superseded" || h.sends != 0 {
		t.Fatal(v)
	}
	v = stepLeader(t, &restarted)
	if v.Leader.Dispatches[1].MessageIDs[0] != first.Leader.Dispatches[0].MessageIDs[0] {
		t.Fatal("lost queued input")
	}
	stepLeader(t, &restarted)
	stepLeader(t, &restarted)
	restarted.localThreads = nil
	v = stepLeader(t, &restarted)
	if v.Leader.Dispatches[1].State != "completed" || h.sends != 1 {
		t.Fatal("replaced sent turn", v)
	}
}
