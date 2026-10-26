package governance

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// Cross-process CLI preparation and MCP commitment share the same receipt and
// canonical publisher. Reports/review in this fixture are explicitly synthetic.
func TestSessionTransferCompiledCLIAndMCP(t *testing.T) {
	root, state := os.Getenv("PG_TRANSFER_LAB_ROOT"), os.Getenv("PG_TRANSFER_LAB_STATE")
	if root != "" && state == "" {
		t.Fatal("explicit private state required")
	}
	s, q := transferFixtureAt(t, root, state)
	binary := os.Getenv("PG_ADVANCE_BINARY")
	if binary == "" {
		binary = filepath.Join(t.TempDir(), "sdptool")
		if b, e := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdptool").CombinedOutput(); e != nil {
			t.Fatal(e, string(b))
		}
	}
	cases := []map[string]any{}
	call := func(op string, q any, def string) json.RawMessage {
		t.Helper()
		b, _ := json.Marshal(q)
		p := filepath.Join(t.TempDir(), "request.json")
		os.WriteFile(p, b, 0600)
		raw, e := exec.Command(binary, s.root, "governance", op, "--state", s.state, "--input", p, "--json").CombinedOutput()
		if e != nil || !json.Valid(raw) {
			t.Fatal(op, e, string(raw))
		}
		cases = append(cases, map[string]any{"definition": def, "value": json.RawMessage(raw)})
		return raw
	}
	ctxQ := SessionSuccessorRequest{q.InstanceID, q.Successor.BindingPath}
	cases = append(cases, map[string]any{"definition": "successorRequest", "value": ctxQ})
	ctxRaw := call("session-successor", ctxQ, "successorContext")
	requestBytes, _ := json.Marshal(ctxQ)
	requestPath := filepath.Join(t.TempDir(), "successor.json")
	os.WriteFile(requestPath, requestBytes, 0600)
	invocation, _ := json.Marshal(map[string]any{"schema": "sdptool.invoke/1", "requestId": "successor", "action": "governance.session.successor", "context": map[string]string{"project": s.root}, "parameters": map[string]string{"state": s.state, "requestFile": requestPath}})
	cmd := exec.Command(binary, "invoke", "--request", "-")
	cmd.Stdin = bytes.NewReader(invocation)
	invoked, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err, string(invoked))
	}
	var envelope struct {
		Result SessionContext `json:"result"`
	}
	if err = json.Unmarshal(invoked, &envelope); err != nil || envelope.Result.InputRevision != q.Successor.InputRevision {
		t.Fatal(err, string(invoked))
	}

	cs := sessionMCPConnect(t, sessionMCPBinary(t), s, q.InstanceID)
	var cliCtx, mcpCtx SessionContext
	json.Unmarshal(ctxRaw, &cliCtx)
	json.Unmarshal(sessionMCPCall(t, cs, "session_successor", ctxQ, false), &mcpCtx)
	if cliCtx.InputRevision != mcpCtx.InputRevision {
		t.Fatal("adapter context mismatch")
	}
	wrong := ctxQ
	wrong.InstanceID = newTestID()
	sessionMCPCall(t, cs, "session_successor", wrong, true)
	cases = append(cases, map[string]any{"definition": "advanceRequest", "value": q})
	raw := call("session-advance", q, "mutation")
	var p SessionMutationResult
	json.Unmarshal(raw, &p)
	retry := sessionMCPMutation(t, cs, "session_advance", q)
	if retry != p {
		t.Fatal("cross-adapter replay differs")
	}
	resume := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Commit explicit remaining-scope transfer"}
	done := sessionMCPMutation(t, cs, "session_resume", resume)
	cases = append(cases, map[string]any{"definition": "mutation", "value": done})
	oldHead := mustFile(t, filepath.Join(s.state, "HEAD"))
	call("session-resume", resume, "mutation")
	if !bytes.Equal(oldHead, mustFile(t, filepath.Join(s.state, "HEAD"))) {
		t.Fatal("retry changed state")
	}
	view, e := s.ReadSession(q.InstanceID)
	if e != nil || view.State != "completed" {
		t.Fatal(view, e)
	}
	cases = append(cases, map[string]any{"definition": "view", "value": view})
	ctx, e := s.SessionContext(SessionContextRequest{q.Successor.BindingPath})
	if e != nil {
		t.Fatal(e)
	}
	cases = append(cases, map[string]any{"definition": "context", "value": ctx})
	b, _ := json.MarshalIndent(cases, "", "  ")
	path := filepath.Join(t.TempDir(), "wire.json")
	os.WriteFile(path, b, 0600)
	if out, e := exec.Command("python3", "check_session_schema.py", path).CombinedOutput(); e != nil {
		t.Fatal(e, string(out))
	}
	if root != "" {
		evidenceFile(t, s, "evidence/transfer-wire.json", string(b)+"\n")
		summary, _ := json.MarshalIndent(map[string]any{"schema": "pg-transfer-lab/1", "predecessorInstance": q.InstanceID, "state": view.State, "successorSession": ctx.Binding.Session.ID, "successorPlan": ctx.Binding.Plan.ID, "successorInputRevision": ctx.InputRevision, "canonicalEffects": "committed", "history": "preserved", "crossAdapterRetry": "unchanged HEAD", "limitations": "Synthetic report/test/review foundation; compiled CLI and MCP transfer only; successor remains unadopted and requires fresh assessment."}, "", "  ")
		evidenceFile(t, s, "evidence/transfer-result.json", string(summary)+"\n")
	}
}
