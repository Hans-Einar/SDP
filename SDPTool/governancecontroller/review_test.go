package governancecontroller

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

type reviewHost struct {
	*selectedHost
	service                    *governance.Service
	reviewID                   string
	thread                     governancehost.Thread
	opens, sends               int
	lostOpen, lostSend, unsafe bool
}

func (h *reviewHost) StartWorkThread(_ context.Context, root, sandbox, instructions, approval string, cfg map[string]any) (governancehost.Thread, error) {
	h.opens++
	v, e := h.service.ReviewRead(h.reviewID)
	if e != nil || v.Review.State != "opening" || sandbox != "read-only" || approval != "never" || !strings.Contains(instructions, "independently assigned Reviewer") {
		h.t.Fatal(v, e)
	}
	b, _ := json.Marshal(cfg)
	if strings.Contains(string(b), "--coordinator") || !strings.Contains(string(b), "--run") || strings.Contains(string(b), "--attempt") || strings.Contains(string(b), "--epoch") {
		h.t.Fatal("review received mutation capability")
	}
	if h.lostOpen {
		return governancehost.Thread{}, fmt.Errorf("lost opening")
	}
	h.thread = governancehost.Thread{ID: "isolated-review", Cwd: root, CLIVersion: h.Version()}
	h.thread.Status.Type = "idle"
	return h.thread, nil
}
func (h *reviewHost) CoordinatorTools(context.Context, string) ([]string, string, error) {
	if h.unsafe {
		return []string{"run_transition"}, "", nil
	}
	return []string{"context_read", "route_evaluate", "run_read"}, "unsupported", nil
}
func (h *reviewHost) ReadThread(context.Context, string) (governancehost.Thread, error) {
	return h.thread, nil
}
func (h *reviewHost) StartSkillTurn(_ context.Context, thread, id, text string, skills []governance.SkillInput) (governancehost.Turn, error) {
	h.sends++
	v, e := h.service.ReviewRead(id)
	if e != nil || v.Review.State != "sending" || text != v.Review.Payload || len(skills) != 1 || skills[0].Name != "sdp-reviewer" {
		h.t.Fatal(v, e)
	}
	in, _ := json.Marshal(map[string]any{"id": "input", "type": "userMessage", "clientId": id, "content": governance.TurnInputs(text, skills)})
	answer, _ := json.Marshal(map[string]any{"id": "answer", "type": "agentMessage", "phase": "final_answer", "text": `{"verdict":"passed","summary":"checked exact candidate","findings":[]}`})
	h.thread.Turns = []governancehost.Turn{{ID: "turn", Status: "completed", ItemsView: "full", Items: []json.RawMessage{in, answer}}}
	if h.lostSend {
		return governancehost.Turn{}, fmt.Errorf("lost send")
	}
	return governancehost.Turn{ID: "turn", Status: "inProgress"}, nil
}
func reviewControllerFixture(t *testing.T) (*ReviewController, *reviewHost, string) {
	t.Helper()
	c, base, id := skillFixture(t)
	s := c.Service
	r, e := s.Read(id)
	if e != nil {
		t.Fatal(e)
	}
	for _, action := range []string{"start", "report"} {
		r, e = s.Transition(governance.LocalController(), governance.TransitionRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: governancehost.UUID(), RunID: id, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: action, Reason: "fixture"})
		if e != nil {
			t.Fatal(e)
		}
	}
	os.WriteFile(filepath.Join(c.Root, "test.log"), []byte("observed external test pass"), 0600)
	candidate, e := s.CaptureCandidate(id, nil)
	if e != nil {
		t.Fatal(e)
	}
	ev := governance.EvidenceRequest{Schema: governance.EvidenceSchema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: governancehost.UUID(), RunID: id, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, ReportID: r.OperationID, Candidate: candidate, Kind: "test", Verdict: "passed", Issuer: "fixture", Note: "external test", Artifact: governance.Source{Path: "test.log", SHA256: fmt.Sprintf("%x", sha256.Sum256([]byte("observed external test pass")))}}
	if _, e = s.AttachEvidence(governance.LocalController(), ev); e != nil {
		t.Fatal(e)
	}
	base.inventory.Project[0].Name = "sdp-reviewer"
	host := &reviewHost{selectedHost: base, service: s, reviewID: governancehost.UUID()}
	rc := &ReviewController{Service: s, Host: host, Root: c.Root, StateDir: c.StateDir, MCPBinary: "/fixture/mcp", SkillPaths: c.SkillPaths}
	_, e = rc.Prepare(context.Background(), host.reviewID, governance.ReviewSpec{Candidate: candidate, ReportID: r.OperationID, EvidenceIDs: []string{ev.OperationID}, Focus: "review candidate independently"})
	if e != nil {
		t.Fatal(e)
	}
	return rc, host, host.reviewID
}
func TestReviewControllerExactRecoveryAndNoResend(t *testing.T) {
	c, h, id := reviewControllerFixture(t)
	ctx := context.Background()
	if _, e := c.Step(ctx, id); e != nil {
		t.Fatal(e)
	}
	h.lostSend = true
	if _, e := c.Step(ctx, id); e == nil {
		t.Fatal("uncertain send accepted")
	}
	recovered := *c
	recovered.localThreads = nil
	v, e := recovered.Step(ctx, id)
	if e != nil || !v.Qualified || h.opens != 1 || h.sends != 1 {
		t.Fatal(v, e, h.opens, h.sends)
	}
	if _, e = recovered.Step(ctx, id); e != nil || h.sends != 1 {
		t.Fatal(e)
	}
	if path := os.Getenv("PG_REVIEW_CONTROLLER_WIRE"); path != "" {
		b, _ := json.MarshalIndent([]map[string]any{{"definition": "ReviewView", "value": v}}, "", "  ")
		if e = os.WriteFile(path, b, 0600); e != nil {
			t.Fatal(e)
		}
	}
}
func TestReviewControllerUnsentAndUnknownOpening(t *testing.T) {
	for _, lost := range []bool{false, true} {
		c, h, id := reviewControllerFixture(t)
		h.lostOpen = lost
		_, e := c.Step(context.Background(), id)
		if lost && e == nil {
			t.Fatal("lost open accepted")
		}
		restarted := *c
		restarted.localThreads = nil
		v, e := restarted.Step(context.Background(), id)
		if lost {
			if e == nil {
				t.Fatal("unknown opening replaced")
			}
		} else if e != nil || v.Review.State != "superseded" {
			t.Fatal(v, e)
		}
		if h.opens != 1 || h.sends != 0 {
			t.Fatal(h.opens, h.sends)
		}
	}
}
func TestReviewControllerRejectsMutationToolsAndSourceDrift(t *testing.T) {
	for _, unsafe := range []bool{true, false} {
		c, h, id := reviewControllerFixture(t)
		if _, e := c.Step(context.Background(), id); e != nil {
			t.Fatal(e)
		}
		if unsafe {
			h.unsafe = true
		} else {
			r, _ := c.Service.Read(h.runID)
			path := r.Run.Assignment.WritablePaths[0]
			os.WriteFile(filepath.Join(c.Root, path), []byte("changed"), 0600)
		}
		if _, e := c.Step(context.Background(), id); e == nil {
			t.Fatal("unsafe review submitted")
		}
		if h.sends != 0 {
			t.Fatal("host send happened")
		}
	}
}
