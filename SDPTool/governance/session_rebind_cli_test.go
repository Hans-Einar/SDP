package governance

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// Compiled adapters share exact preparation, commitment and replay. The pinned
// review is synthetic controller attribution, not native independent review.
func TestSessionRebindCompiledCLIAndMCP(t *testing.T) {
	root, state := os.Getenv("PG_REBIND_LAB_ROOT"), os.Getenv("PG_REBIND_LAB_STATE")
	if root != "" && state == "" {
		t.Fatal("explicit private state required")
	}
	s, id := rebindFixtureAt(t, root, state)
	editRebindPlan(t, s, id, func(b *SessionBinding) { b.Steps[0].WritablePaths = append(b.Steps[0].WritablePaths, "extra.txt") })
	q := rebindRequest(t, s, id)
	binary := os.Getenv("PG_REBIND_BINARY")
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
	cq := SessionRebindContextRequest{id}
	cases = append(cases, map[string]any{"definition": "rebindContextRequest", "value": cq})
	ctxRaw := call("session-rebind-context", cq, "rebindContext")
	requestBytes, _ := json.Marshal(cq)
	requestPath := filepath.Join(t.TempDir(), "context.json")
	os.WriteFile(requestPath, requestBytes, 0600)
	invocation, _ := json.Marshal(map[string]any{"schema": "sdptool.invoke/1", "requestId": "rebind-preview", "action": "governance.session.rebind-context", "context": map[string]string{"project": s.root}, "parameters": map[string]string{"state": s.state, "requestFile": requestPath}})
	cmd := exec.Command(binary, "invoke", "--request", "-")
	cmd.Stdin = bytes.NewReader(invocation)
	invoked, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err, string(invoked))
	}
	var envelope struct {
		Result SessionRebindContext `json:"result"`
	}
	if err = json.Unmarshal(invoked, &envelope); err != nil || envelope.Result.Next.InputRevision != q.InputRevision {
		t.Fatal(err, string(invoked))
	}
	cs := sessionMCPConnect(t, sessionMCPBinary(t), s, id)
	var cliCtx, mcpCtx SessionRebindContext
	json.Unmarshal(ctxRaw, &cliCtx)
	json.Unmarshal(sessionMCPCall(t, cs, "session_rebind_context", cq, false), &mcpCtx)
	if !reflect.DeepEqual(cliCtx, mcpCtx) {
		t.Fatal("adapter context mismatch")
	}
	sessionMCPCall(t, cs, "session_rebind_context", SessionRebindContextRequest{newTestID()}, true)
	wrong := q
	wrong.InstanceID = newTestID()
	sessionMCPCall(t, cs, "session_rebind", wrong, true)
	forged, _ := json.Marshal(q)
	var fields map[string]any
	json.Unmarshal(forged, &fields)
	fields["authority"] = "owner"
	sessionMCPCall(t, cs, "session_rebind", fields, true)
	cases = append(cases, map[string]any{"definition": "rebindRequest", "value": q}, map[string]any{"definition": "rebindReview", "value": json.RawMessage(mustFile(t, filepath.Join(s.root, q.Review.Path)))})
	var p SessionMutationResult
	json.Unmarshal(call("session-rebind", q, "mutation"), &p)
	retry := sessionMCPMutation(t, cs, "session_rebind", q)
	if retry != p {
		t.Fatal("cross-adapter replay differs")
	}
	resume := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), id, p.Revision, "Commit exact reviewed binding revision"}
	done := sessionMCPMutation(t, cs, "session_resume", resume)
	cases = append(cases, map[string]any{"definition": "mutation", "value": done})
	oldHead := mustFile(t, filepath.Join(s.state, "HEAD"))
	call("session-resume", resume, "mutation")
	if !bytes.Equal(oldHead, mustFile(t, filepath.Join(s.state, "HEAD"))) {
		t.Fatal("replay changed HEAD")
	}
	v, e := s.ReadSession(id)
	if e != nil || v.Coverage != "current" || len(v.BindingRevisions) != 1 {
		t.Fatal(v, e)
	}
	cases = append(cases, map[string]any{"definition": "view", "value": v})
	b, _ := json.MarshalIndent(cases, "", "  ")
	path := filepath.Join(t.TempDir(), "wire.json")
	os.WriteFile(path, b, 0600)
	if out, e := exec.Command("python3", "check_session_schema.py", path).CombinedOutput(); e != nil {
		t.Fatal(e, string(out))
	}
	if root != "" {
		evidenceFile(t, s, "evidence/rebind-wire.json", string(b)+"\n")
		summary, _ := json.MarshalIndent(map[string]any{"schema": "pg-rebind-lab/1", "instanceId": id, "bindingRevisions": v.BindingRevisions, "coverage": v.Coverage, "canonicalEffects": "committed", "crossAdapterRetry": "unchanged HEAD", "limitations": "Synthetic attributed review; compiled CLI/JSON/MCP binding revision only. Fresh selection/assessment required; no native implementation or owner acceptance claimed."}, "", "  ")
		evidenceFile(t, s, "evidence/rebind-result.json", string(summary)+"\n")
	}
}
