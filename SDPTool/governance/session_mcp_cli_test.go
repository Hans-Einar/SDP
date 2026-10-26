package governance

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func sessionMCPBinary(t *testing.T) string {
	t.Helper()
	if b := os.Getenv("PG_SESSION_MCP_BINARY"); b != "" {
		return b
	}
	b := filepath.Join(t.TempDir(), "sdp-mcp")
	if out, e := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", b, "../cmd/sdp-mcp").CombinedOutput(); e != nil {
		t.Fatal(e, string(out))
	}
	return b
}
func sessionMCPConnect(t *testing.T, binary string, s *Service, id string) *mcp.ClientSession {
	t.Helper()
	// A named launcher skill is separate from synthetic assessment evidence;
	// this fixture makes no native enabled-skill inventory/consumption claim.
	skillPath := filepath.Join(s.root, "evidence/launcher-leader/SKILL.md")
	if _, e := os.Stat(skillPath); os.IsNotExist(e) {
		evidenceFile(t, s, "evidence/launcher-leader/SKILL.md", "---\nname: sdp-project-leader\ndescription: Synthetic launcher procedure.\n---\n\nUse only the assigned Session and preserve owner gates.\n")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	t.Cleanup(cancel)
	cmd := exec.CommandContext(ctx, binary, "--project", s.root, "--state", s.state, "--session-controller", id, "--leader-skill", "evidence/launcher-leader/SKILL.md")
	cs, e := mcp.NewClient(&mcp.Implementation{Name: "session-protocol-fixture", Version: "1"}, nil).Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { cs.Close() })
	return cs
}
func sessionMCPCall(t *testing.T, c *mcp.ClientSession, name string, q any, wantError bool) json.RawMessage {
	t.Helper()
	r, e := c.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: q})
	if wantError {
		if e == nil && !r.IsError {
			t.Fatal("unexpected privilege", name, r)
		}
		return nil
	}
	if e != nil {
		t.Fatal(name, e)
	}
	b, _ := json.Marshal(r.StructuredContent)
	if r.IsError {
		t.Fatal(name, string(b))
	}
	return b
}
func sessionMCPMutation(t *testing.T, c *mcp.ClientSession, name string, q any) SessionMutationResult {
	t.Helper()
	var out SessionMutationResult
	if e := json.Unmarshal(sessionMCPCall(t, c, name, q, false), &out); e != nil {
		t.Fatal(e)
	}
	return out
}
func TestSessionMCPCompiledLifecycle(t *testing.T) {
	root, state := os.Getenv("PG_SESSION_MCP_LAB_ROOT"), os.Getenv("PG_SESSION_MCP_LAB_STATE")
	if root != "" && state == "" {
		t.Fatal("explicit private state required")
	}
	s, q := workFixture(t, root, state)
	binary := sessionMCPBinary(t)
	c := sessionMCPConnect(t, binary, s, q.InstanceID)
	list, e := c.ListTools(context.Background(), nil)
	if e != nil || len(list.Tools) != 15 {
		t.Fatal(list, e)
	}
	for _, tool := range list.Tools {
		if tool.Annotations.ReadOnlyHint != (tool.Name == "session_read" || tool.Name == "session_definition" || tool.Name == "session_successor" || tool.Name == "session_rebind_context") {
			t.Fatal(tool.Name)
		}
	}
	sessionMCPCall(t, c, "session_definition", map[string]any{}, false)
	head := mustFile(t, filepath.Join(s.state, "HEAD"))
	bad := q
	bad.InstanceID = newTestID()
	sessionMCPCall(t, c, "session_start", bad, true)
	forged, _ := json.Marshal(q)
	var fields map[string]any
	json.Unmarshal(forged, &fields)
	fields["authority"] = "owner"
	sessionMCPCall(t, c, "session_start", fields, true)
	if !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
		t.Fatal("refused adapter changed state")
	}
	p := sessionMCPMutation(t, c, "session_start", q)
	direct, e := s.StartSession(LocalController(), q)
	if e != nil || direct != p {
		t.Fatal("core/MCP receipt differs", direct, e)
	}
	r := workResume(q, p)
	done := sessionMCPMutation(t, c, "session_resume", r)
	c.Close()
	c = sessionMCPConnect(t, binary, s, q.InstanceID)
	if got := sessionMCPMutation(t, c, "session_resume", r); got != done {
		t.Fatal("reconnect replay differs")
	}
	run, e := s.Read(q.OperationID)
	if e != nil {
		t.Fatal(e)
	}
	run = reportWork(t, s, run)
	pq := progressRequest(t, s, q, done, run)
	p = sessionMCPMutation(t, c, "session_progress", pq)
	done = sessionMCPMutation(t, c, "session_resume", workResume(pq, p))
	pq.OperationID = newTestID()
	pq.ExpectedRevision = done.Revision
	p = sessionMCPMutation(t, c, "session_verify", pq)
	done = sessionMCPMutation(t, c, "session_resume", workResume(pq, p))
	rq := passingSessionReview(t, s, pq, done)
	p = sessionMCPMutation(t, c, "session_review", rq)
	done = sessionMCPMutation(t, c, "session_resume", workResume(rq, p))
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	d := SessionDisposeRequest{Schema: SessionSchema, ProjectRef: s.project, CheckoutID: st.CheckoutID, OperationID: newTestID(), InstanceID: q.InstanceID, ExpectedRevision: done.Revision, ExpectedRunRevision: st.Runs[q.OperationID].Revision, ResultID: rq.OperationID, System: "SDPTOOL", Reason: "Dispose the checked compiled MCP fixture", Cards: []SessionCardDisposition{{ID: "KB-SDP-001", Outcome: "completed", Summary: "Exact synthetic fixture verified"}}}
	p = sessionMCPMutation(t, c, "session_dispose", d)
	sessionMCPMutation(t, c, "session_resume", disposeResume(d, p))
	aq := advanceRequest(t, s, q.InstanceID, "complete", "", "disposition", d.OperationID)
	p = sessionMCPMutation(t, c, "session_advance", aq)
	sessionMCPMutation(t, c, "session_resume", SessionResumeRequest{SessionSchema, aq.ProjectRef, aq.CheckoutID, newTestID(), aq.InstanceID, p.Revision, "Commit exact Session closure through shared core"})
	raw := sessionMCPCall(t, c, "session_read", map[string]any{}, false)
	var v SessionView
	json.Unmarshal(raw, &v)
	if v.State != "completed" || v.Coverage != "current" {
		t.Fatal(string(raw))
	}
	directView, e := s.ReadSession(q.InstanceID)
	if e != nil || !reflect.DeepEqual(v, directView) {
		t.Fatal("read adapter differs", e)
	}
	if _, e = s.publicationFiles(); e != nil {
		t.Fatal(e)
	}
	// An already running launcher must not silently consume changed procedure bytes.
	skillPath := filepath.Join(s.root, "evidence/launcher-leader/SKILL.md")
	before := mustFile(t, skillPath)
	os.WriteFile(skillPath, append(before, []byte("\nChanged procedure.\n")...), 0600)
	sessionMCPCall(t, c, "session_read", map[string]any{}, true)
	os.WriteFile(skillPath, before, 0600)
	sessionMCPCall(t, c, "session_read", map[string]any{}, false)
	if root != "" {
		os.WriteFile(filepath.Join(root, "evidence/mcp-session-view.json"), raw, 0600)
	}
}
func TestSessionMCPRejectsOwnerAndMixedLaunch(t *testing.T) {
	s, q, _ := ownerFixture(t)
	binary := sessionMCPBinary(t)
	c := sessionMCPConnect(t, binary, s, q.InstanceID)
	p := sessionMCPMutation(t, c, "session_owner_request", q)
	done := sessionMCPMutation(t, c, "session_resume", SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Publish concrete owner request"})
	head := mustFile(t, filepath.Join(s.state, "HEAD"))
	for _, action := range []string{"accept", "changes-requested"} {
		sessionMCPCall(t, c, "session_owner_request", ownerDecision(q, done.Revision, action), true)
	}
	for _, name := range []string{"session_adopt", "run_start", "design_decision", "owner_accept"} {
		if _, e := c.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: map[string]any{}}); e == nil {
			t.Fatal("ungranted tool", name)
		}
	}
	if !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
		t.Fatal("owner refusal wrote state")
	}
	for _, extra := range [][]string{{"--coordinator", "other"}, {"--run", newTestID()}, {"--attempt", newTestID()}, {"--epoch", newTestID()}} {
		args := []string{"--project", s.root, "--state", s.state, "--session-controller", q.InstanceID, "--leader-skill", "evidence/launcher-leader/SKILL.md"}
		args = append(args, extra...)
		if b, e := exec.Command(binary, args...).CombinedOutput(); e == nil {
			t.Fatal("mixed launch", string(b))
		}
	}
	if _, e := s.SessionLeaderSkill(q.InstanceID, "Skills/sdp-worker/SKILL.md"); e == nil {
		t.Fatal("worker relabeled as leader")
	}
	if _, e := s.SessionLeaderSkill(q.InstanceID, "../outside.md"); e == nil {
		t.Fatal("outside skill")
	}
}

func TestSessionMCPFinalLabReplay(t *testing.T) {
	root, state := os.Getenv("PG_SESSION_MCP_REPLAY_ROOT"), os.Getenv("PG_SESSION_MCP_REPLAY_STATE")
	if root == "" {
		t.Skip("explicit durable lab replay only")
	}
	if state == "" {
		t.Fatal("private state required")
	}
	s, e := New(root, filepath.Join(root, "SDP"), "trial", state)
	if e != nil {
		t.Fatal(e)
	}
	st, _, e := s.load()
	if e != nil || len(st.Sessions) != 1 {
		t.Fatal(e)
	}
	var instance SessionInstance
	for _, i := range st.Sessions {
		instance = i
	}
	if instance.Advance == nil || instance.Advance.Request.Action != "complete" || instance.Advance.Effects.Status != "committed" {
		t.Fatal("closed fixture required")
	}
	q := instance.Advance.Request
	var resumeID string
	for id, r := range st.Receipts {
		if v := r.Result.Session; v != nil && v.Operation == "session-resume" && v.InstanceID == instance.ID && v.Revision == q.ExpectedRevision+2 {
			resumeID = id
		}
	}
	if resumeID == "" {
		t.Fatal("exact original resume receipt missing")
	}
	head := mustFile(t, filepath.Join(state, "HEAD"))
	c := sessionMCPConnect(t, sessionMCPBinary(t), s, instance.ID)
	p := sessionMCPMutation(t, c, "session_advance", q)
	rr := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, resumeID, q.InstanceID, p.Revision, "Commit exact Session closure through shared core"}
	done := sessionMCPMutation(t, c, "session_resume", rr)
	raw := sessionMCPCall(t, c, "session_read", map[string]any{}, false)
	if done != *st.Receipts[resumeID].Result.Session || !bytes.Equal(head, mustFile(t, filepath.Join(state, "HEAD"))) {
		t.Fatal("historical replay mutated state or receipt")
	}
	out, _ := json.MarshalIndent(map[string]any{"headUnchanged": true, "advance": p, "resume": done, "view": json.RawMessage(raw), "limitation": "actual compiled MCP replay; synthetic earlier agent/test/review observations"}, "", "  ")
	if e = os.WriteFile(filepath.Join(root, "evidence/mcp-final-replay.json"), out, 0600); e != nil {
		t.Fatal(e)
	}
}
