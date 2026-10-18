package governance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestReattachPreservesOnlyBoundCapability(t *testing.T) {
	s, _, q := fixture(t)
	bind(t, s, &q)
	r, err := s.Start(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	r, err = s.Transition(LocalController(), transition(r, "start"))
	if err != nil {
		t.Fatal(err)
	}
	attempt, _ := identifier()
	original, _ := identifier()
	observer, _ := identifier()
	host := strings.Repeat("a", 64)
	agent, _ := AssignedAgent(r.Run.ID, attempt, original)
	request := func(action, kind, epoch, status string) HostRequest {
		op, _ := identifier()
		return HostRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, OperationID: op, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, AttemptID: attempt, HostVersion: "0.161.0", HostIdentity: host, ThreadID: "thread", Action: action, Observation: Observation{Kind: kind, Epoch: epoch, Status: status, TurnID: "turn", ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}}
	}
	h := request("bind", "bound", original, "idle")
	h.Observation.TurnID = ""
	r, err = s.RecordHost(LocalController(), h)
	if err != nil {
		t.Fatal(err)
	}
	r, err = s.RecordHost(LocalController(), request("observe", "turn.started", original, "running"))
	if err != nil {
		t.Fatal(err)
	}
	r, err = s.RecordHost(LocalController(), request("disconnect", "disconnected", original, "unknown"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Transition(agent, transition(r, "report"))
	code(t, err, "unauthorized")
	h = request("reattach", "reconciled", observer, "running")
	h.HostIdentity = strings.Repeat("b", 64)
	_, err = s.RecordHost(LocalController(), h)
	code(t, err, "unauthorized")
	h = request("reattach", "reconciled", observer, "running")
	h.Observation.TurnID = "other"
	_, err = s.RecordHost(LocalController(), h)
	code(t, err, "unauthorized")
	h = request("reattach", "reconciled", observer, "running")
	_, err = s.RecordHost(agent, h)
	code(t, err, "unauthorized")
	r, err = s.RecordHost(LocalController(), h)
	if err != nil {
		t.Fatal(err)
	}
	if dir := os.Getenv("PGI_SCHEMA_EVIDENCE"); dir != "" {
		cases, _ := json.Marshal([]map[string]any{{"definition": "HostRequest", "value": h}, {"definition": "Result", "value": r}})
		if err = os.WriteFile(filepath.Join(dir, "attach-wire-cases.json"), cases, 0600); err != nil {
			t.Fatal(err)
		}
	}
	a := r.Run.Attempts[0]
	if !a.Gap || a.Epoch != observer || a.CapabilityEpoch() != original || len(r.Run.Attempts) != 1 {
		t.Fatal(a)
	}
	replay, err := s.RecordHost(LocalController(), h)
	if err != nil || replay.Run.Revision != r.Run.Revision {
		t.Fatal(replay, err)
	}
	_, err = s.RecordHost(LocalController(), request("observe", "turn.completed", original, "completed"))
	code(t, err, "reconciliation_required")
	forged, _ := AssignedAgent(r.Run.ID, attempt, observer)
	_, err = s.Transition(forged, transition(r, "report"))
	code(t, err, "unauthorized")
	r, err = s.Transition(agent, transition(r, "report"))
	if err != nil || r.Run.State != "verifying" {
		t.Fatal(r, err)
	}
}

func TestReattachRefusesMissingAuthorityAndUnresolvedWork(t *testing.T) {
	for _, scenario := range []string{"legacy", "approval", "input", "source"} {
		t.Run(scenario, func(t *testing.T) {
			s, _, q := fixture(t)
			bind(t, s, &q)
			r, err := s.Start(LocalController(), q)
			if err != nil {
				t.Fatal(err)
			}
			r, err = s.Transition(LocalController(), transition(r, "start"))
			if err != nil {
				t.Fatal(err)
			}
			attempt, _ := identifier()
			epoch, _ := identifier()
			op, _ := identifier()
			host := strings.Repeat("a", 64)
			h := HostRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, OperationID: op, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, AttemptID: attempt, HostVersion: "0.161.0", HostIdentity: host, ThreadID: "thread", Action: "bind", Observation: Observation{Kind: "bound", Epoch: epoch, Status: "idle", ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}}
			record := func() {
				t.Helper()
				h.OperationID, _ = identifier()
				h.ExpectedRevision = r.Run.Revision
				r, err = s.RecordHost(LocalController(), h)
				if err != nil {
					t.Fatal(err)
				}
			}
			if scenario == "legacy" {
				h.HostIdentity = ""
			}
			record()
			h.Action = "observe"
			h.Observation.Kind = "turn.started"
			h.Observation.TurnID = "turn"
			h.Observation.Status = "running"
			record()
			if scenario == "approval" {
				prompt := ApprovalPrompt{RequestID: "0", Epoch: epoch, ThreadID: "thread", TurnID: "turn", Message: `Allow the sdp_governance MCP server to run tool "run_transition"?`, Request: transition(r, "report")}
				b, _ := json.Marshal(prompt)
				h.Observation.Kind = "approval.requested"
				h.Observation.ItemID = "0"
				h.Observation.Status = ""
				h.Observation.Text = string(b)
				record()
			}
			if scenario == "input" {
				id, _ := identifier()
				op, _ := identifier()
				r, err = s.RecordInput(LocalController(), InputRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, OperationID: op, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: "prepare", Input: OwnerInput{ID: id, Kind: "steer", Text: "exact", Payload: "exact", AttemptID: attempt, Epoch: epoch, ThreadID: "thread", TurnID: "turn", ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}})
				if err != nil {
					t.Fatal(err)
				}
			}
			h.Action = "disconnect"
			h.Observation = Observation{Kind: "disconnected", Epoch: epoch, Status: "unknown", TurnID: "turn", ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}
			record()
			if scenario == "source" {
				if err = os.WriteFile(filepath.Join(s.root, q.Selection.Authority), []byte("changed authority"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			h.Action = "reattach"
			h.HostIdentity = host
			h.ExpectedRevision = r.Run.Revision
			h.OperationID, _ = identifier()
			h.Observation.Kind = "reconciled"
			h.Observation.Epoch, _ = identifier()
			h.Observation.Status = "running"
			_, err = s.RecordHost(LocalController(), h)
			want := "reconciliation_required"
			if scenario == "legacy" {
				want = "unauthorized"
			}
			if scenario == "source" {
				want = "stale_revision"
			}
			code(t, err, want)
		})
	}
}
