package governancecontroller

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/activities"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
)

type pendingAtThreadHost struct {
	*fakeHost
	pending func()
}

func (h *pendingAtThreadHost) StartWorkThread(ctx context.Context, a, b, c, d string, config map[string]any) (governancehost.Thread, error) {
	thread, err := h.fakeHost.StartWorkThread(ctx, a, b, c, d, config)
	if err == nil {
		h.pending()
	}
	return thread, err
}

func TestActivityHoldDuringClaimedLaunchUsesExistingSendingRecovery(t *testing.T) {
	service, sel, q := fixture(t)
	coverage, err := service.Context(sel)
	if err != nil {
		t.Fatal(err)
	}
	root := coverage.Root
	area := filepath.Join(root, "SDP")
	state := filepath.Join(filepath.Dir(root), "state")
	ledger := filepath.Join(area, "ProjectManagement/Ledger.ndjson")
	raw, err := os.ReadFile(ledger)
	if err != nil {
		t.Fatal(err)
	}
	var valid []byte
	for _, line := range bytes.Split(raw, []byte{'\n'}) {
		if len(line) == 0 {
			continue
		}
		var event map[string]any
		if err = json.Unmarshal(line, &event); err != nil {
			t.Fatal(err)
		}
		event["eventId"] = "EVT-" + event["eventId"].(string)
		event["actor"] = "fixture"
		event["occurredAt"] = "2026-11-12T00:00:00Z"
		event["commit"] = nil
		b, _ := json.Marshal(event)
		valid = append(valid, b...)
		valid = append(valid, '\n')
	}
	if err = os.WriteFile(ledger, valid, 0600); err != nil {
		t.Fatal(err)
	}
	bind(t, service, &q)
	domain := activities.New(area, "TEST")
	owner := activities.Caller{Actor: "owner", Authority: activities.Authority{Role: "project-owner", Source: "authority.md"}}
	snapshot, err := projecthistory.Read(ledger)
	if err != nil {
		t.Fatal(err)
	}
	created, err := domain.Create(activities.CreateRequest{OperationID: governancehost.UUID(), ExpectedHistoryRevision: snapshot.Revision, Title: "Need", Owner: "owner", Mandate: "Need", Reason: "Owner selected", Scope: []string{"target.txt"}, References: []activities.Reference{{Role: "session", ID: "SESSION-TEST-0001", Path: "Sessions/session.md"}}}, owner)
	if err != nil {
		t.Fatal(err)
	}
	prev := created.EventID
	from := "planned"
	activation := activities.Event{SchemaVersion: "1.0", EventID: "EVT-ACTIVITY-TEST-900001", EventType: "x-activity:transitioned", OccurredAt: "2026-11-12T00:00:00Z", Actor: "owner", SubjectID: created.ActivityID, Payload: activities.Payload{SchemaVersion: "0.2", ProjectID: "TEST", OperationID: governancehost.UUID(), RecordedAt: "2026-11-12T00:00:00Z", Authority: owner.Authority, PreviousEventID: &prev, From: &from, To: "active", Reason: "Owner selects", Descriptor: activities.DescriptorReference{Path: created.Path + "/activity.yaml", SHA256: created.DescriptorRevision}, RequestDigest: created.DescriptorRevision, Links: []string{}}}
	b, _ := json.Marshal(activation)
	appended, err := projecthistory.Append(ledger, created.HistoryRevision, b)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = domain.AdoptAdmission(activities.AdoptionRequest{OperationID: governancehost.UUID(), ExpectedHistoryRevision: appended.Revision, Stores: []activities.ExecutionStore{{StatePath: state}}}, owner, governance.ActivityExecutionInspector(root, area, "trial")); err != nil {
		t.Fatal(err)
	}
	coverage, err = service.Context(sel)
	if err != nil {
		t.Fatal(err)
	}
	q.InputRevision = coverage.Revision
	run, err := service.Start(governance.LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	controller := &Controller{Service: service, Root: root, StateDir: state, MCPBinary: "mcp", Epoch: governancehost.UUID(), Sandbox: "read-only"}
	fake := &fakeHost{t: t, c: controller, runID: run.Run.ID, mode: "chatgpt"}
	controller.Host = &pendingAtThreadHost{fakeHost: fake, pending: func() {
		// This host callback occurs after durable sending and outside any ledger lock.
		tx, err := projecthistory.BeginTransaction(ledger)
		if err != nil {
			t.Fatal("host RPC retained ledger lock", err)
		}
		defer tx.Close()
		if err = domain.SetPendingLifecycle(tx, "", activities.PendingLifecycle{OperationID: governancehost.UUID(), ActivityID: created.ActivityID, Action: "hold"}); err != nil {
			t.Fatal(err)
		}
		inspected, err := governance.ActivityExecutionInspector(root, area, "trial").InspectExecutions([]activities.ExecutionStore{{StatePath: state}})
		if err != nil || len(inspected.InFlight) != 1 {
			t.Fatal(inspected, err)
		}
	}}
	result, err := controller.Launch(context.Background(), run.Run.ID, "Fix the selected fixture")
	if err != nil || fake.starts != 1 || result.Run.State != "running" {
		t.Fatal(result, err, fake.starts)
	}
	inspected, err := governance.ActivityExecutionInspector(root, area, "trial").InspectExecutions([]activities.ExecutionStore{{StatePath: state}})
	if err != nil || len(inspected.InFlight) != 1 {
		t.Fatal("running work disappeared behind hold", inspected, err)
	}
	if _, err = controller.Control(context.Background(), run.Run.ID, "interrupt", "Owner requests quiescence"); err != nil {
		t.Fatal("pending hold blocked interrupt", err)
	}
	inspected, err = governance.ActivityExecutionInspector(root, area, "trial").InspectExecutions([]activities.ExecutionStore{{StatePath: state}})
	if err != nil || len(inspected.InFlight) != 1 {
		t.Fatal("interrupt request was mistaken for observed stop", inspected, err)
	}
	_, err = controller.Observe(run.Run.ID, governancehost.Event{ThreadID: "thread", Observation: governance.Observation{Kind: "turn.completed", TurnID: "turn", Status: "interrupted", Epoch: controller.Epoch, ReceivedAt: "2026-11-12T01:00:00Z"}})
	if err != nil {
		t.Fatal("pending gate blocked terminal observation", err)
	}
	inspected, err = governance.ActivityExecutionInspector(root, area, "trial").InspectExecutions([]activities.ExecutionStore{{StatePath: state}})
	if err != nil || len(inspected.InFlight) != 0 {
		t.Fatal("observed interruption did not establish quiescence", inspected, err)
	}
}
