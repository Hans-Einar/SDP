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

func hostCommandFixture(t *testing.T) (*Service, SessionHostRequest) {
	t.Helper()
	s, work := workFixture(t, "", "")
	done, run := startWork(t, s, work)
	path := "evidence/host-leader/SKILL.md"
	evidenceFile(t, s, path, "---\nname: sdp-project-leader\ndescription: Test leader.\n---\n\nPreserve Session scope and owner gates.\n")
	skill, err := s.SessionLeaderSkill(work.InstanceID, path)
	if err != nil {
		t.Fatal(err)
	}
	return s, SessionHostRequest{Schema: SessionHostSchema, ProjectRef: s.project, CheckoutID: run.CheckoutID, OperationID: newTestID(), InstanceID: work.InstanceID, ExpectedRevision: done.Revision, RunID: run.Run.ID, ExpectedRunRevision: run.Run.Revision, Action: "launch", Reason: "Execute the committed test assignment", LeaderSkill: skill}
}

func TestSessionHostRequestRecoveryAndIdempotency(t *testing.T) {
	s, q := hostCommandFixture(t)
	head := mustFile(t, filepath.Join(s.state, "HEAD"))
	if _, err := s.RequestSessionHost(Caller{}, q); err == nil {
		t.Fatal("untrusted request accepted")
	}
	if !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
		t.Fatal("refusal mutated state")
	}
	command, err := s.RequestSessionHost(LocalController(), q)
	if err != nil || command.State != "queued" {
		t.Fatal(command, err)
	}
	again, err := s.RequestSessionHost(LocalController(), q)
	if err != nil || !reflect.DeepEqual(again, command) {
		t.Fatal("retry differs", again, err)
	}
	changed := q
	changed.Reason = "Different intent"
	if _, err = s.RequestSessionHost(LocalController(), changed); err == nil {
		t.Fatal("identity reused")
	}
	changed = q
	changed.OperationID = newTestID()
	if _, err = s.RequestSessionHost(LocalController(), changed); err == nil {
		t.Fatal("parallel command accepted")
	}
	command, err = s.ClaimSessionHost(LocalController(), q.OperationID)
	if err != nil || command.State != "claimed" {
		t.Fatal(command, err)
	}
	if _, err = s.ClaimSessionHost(LocalController(), q.OperationID); err == nil {
		t.Fatal("claim replayed")
	}
	command, err = s.SettleSessionHost(LocalController(), q.OperationID)
	if err != nil || command.State != "uncertain" {
		t.Fatal(command, err)
	}
	// A process lost before preparing an input can be explicitly canceled under lease.
	command, err = s.CancelSessionHost(LocalController(), q.OperationID)
	if err != nil || command.State != "canceled" {
		t.Fatal(command, err)
	}
	command, err = s.RequestSessionHost(LocalController(), q)
	if err != nil || command.State != "canceled" {
		t.Fatal("retry revived canceled request", command, err)
	}
	if _, err = s.RequestSessionHost(LocalController(), changed); err != nil {
		t.Fatal(err)
	}
	restarted, err := New(s.root, s.area, s.project, s.state)
	if err != nil {
		t.Fatal(err)
	}
	commands, err := restarted.ReadSessionHostCommands(q.InstanceID)
	if err != nil || len(commands.Commands) != 2 {
		t.Fatal(commands, err)
	}
	st, _, err := s.load()
	if err != nil || st.Schema != sessionHostStateSchema {
		t.Fatal(st.Schema, err)
	}
	if err = sessionHostHold(st, q.InstanceID); err == nil {
		t.Fatal("pending request does not block Session close/rebind")
	}
}

func TestSessionHostRequestExactBasis(t *testing.T) {
	for _, change := range []string{"session-revision", "run-revision", "skill", "action", "identity"} {
		t.Run(change, func(t *testing.T) {
			s, q := hostCommandFixture(t)
			switch change {
			case "session-revision":
				q.ExpectedRevision++
			case "run-revision":
				q.ExpectedRunRevision++
			case "skill":
				if err := os.WriteFile(q.LeaderSkill.Path, []byte("changed"), 0600); err != nil {
					t.Fatal(err)
				}
			case "action":
				q.Action = "interrupt"
			case "identity":
				q.OperationID = q.RunID
			}
			head := mustFile(t, filepath.Join(s.state, "HEAD"))
			if _, err := s.RequestSessionHost(LocalController(), q); err == nil {
				t.Fatal("invalid basis accepted")
			}
			if !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
				t.Fatal("refusal mutated state")
			}
		})
	}
}

func TestSessionHostClaimRechecksAndLease(t *testing.T) {
	s, q := hostCommandFixture(t)
	if _, err := s.RequestSessionHost(LocalController(), q); err != nil {
		t.Fatal(err)
	}
	release, err := s.SessionHostCommandLease(LocalController(), q.OperationID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.CancelSessionHost(LocalController(), q.OperationID); err == nil {
		t.Fatal("canceled leased command")
	}
	release()
	evidenceFile(t, s, "SDP/authority.md", "Changed authority after queueing.\n")
	if _, err = s.ClaimSessionHost(LocalController(), q.OperationID); err == nil {
		t.Fatal("stale queued basis dispatched")
	}
	if _, err = s.CancelSessionHost(LocalController(), q.OperationID); err != nil {
		t.Fatal(err)
	}
}

func TestSessionHostExecutorDeliveryAndLostResponse(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "executor")
	if b, e := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "./testdata/session-host-executor").CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	}
	for _, mode := range []string{"success", "lost"} {
		t.Run(mode, func(t *testing.T) {
			s, q := hostCommandFixture(t)
			if _, err := s.RequestSessionHost(LocalController(), q); err != nil {
				t.Fatal(err)
			}
			b, e := exec.Command(binary, s.root, s.state, s.project, q.OperationID, mode).CombinedOutput()
			if e != nil {
				t.Fatal(e, string(b))
			}
			if !bytes.Contains(b, []byte(`"starts":1`)) {
				t.Fatal(string(b))
			}
		})
	}
}

func TestSessionHostCompiledCLIAndInvocation(t *testing.T) {
	s, q := hostCommandFixture(t)
	binary := filepath.Join(t.TempDir(), "sdptool")
	if b, e := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdptool").CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	}
	requestPath := filepath.Join(t.TempDir(), "request.json")
	request, _ := json.Marshal(q)
	if e := os.WriteFile(requestPath, request, 0600); e != nil {
		t.Fatal(e)
	}
	raw, e := exec.Command(binary, s.root, "governance", "session-host-request", "--state", s.state, "--input", requestPath, "--json").CombinedOutput()
	if e != nil {
		t.Fatal(e, string(raw))
	}
	var cli struct {
		Schema  string             `json:"schema"`
		Command SessionHostCommand `json:"command"`
	}
	if json.Unmarshal(raw, &cli) != nil || cli.Schema != SessionHostSchema || cli.Command.State != "queued" {
		t.Fatal(string(raw))
	}
	invoke, _ := json.Marshal(map[string]any{"schema": "sdptool.invoke/1", "requestId": "host-request", "action": "governance.session.host-request", "context": map[string]string{"project": s.root}, "parameters": map[string]string{"state": s.state, "requestFile": requestPath}})
	command := exec.Command(binary, "invoke", "--request", "-", "--as", "test-controller", "--authority", "controller")
	command.Stdin = bytes.NewReader(invoke)
	raw, e = command.CombinedOutput()
	if e != nil {
		t.Fatal(e, string(raw))
	}
	var out struct {
		Result struct {
			Command SessionHostCommand `json:"command"`
		} `json:"result"`
	}
	if json.Unmarshal(raw, &out) != nil || !reflect.DeepEqual(cli.Command, out.Result.Command) {
		t.Fatal(string(raw))
	}
	raw, e = exec.Command(binary, s.root, "governance", "session-host-read", "--state", s.state, "--instance", q.InstanceID, "--json").CombinedOutput()
	if e != nil {
		t.Fatal(e, string(raw))
	}
	var view SessionHostCommands
	if json.Unmarshal(raw, &view) != nil || view.Current == nil || view.Current.RunRevision != q.ExpectedRunRevision || len(view.Commands) != 1 {
		t.Fatal(string(raw))
	}
	cancel, _ := json.Marshal(map[string]string{"operationId": q.OperationID})
	os.WriteFile(requestPath, cancel, 0600)
	raw, e = exec.Command(binary, s.root, "governance", "session-host-cancel", "--state", s.state, "--input", requestPath, "--json").CombinedOutput()
	if e != nil || !bytes.Contains(raw, []byte(`"canceled"`)) {
		t.Fatal(e, string(raw))
	}
}
