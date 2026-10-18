package governance

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestOverviewEmptyAndRegisteredWorkDoesNotWrite(t *testing.T) {
	s, sel, q := fixture(t)
	v, err := s.Overview()
	if err != nil || len(v.Runs) != 0 || v.Generation != "" || v.CheckoutID != "" {
		t.Fatal(v, err)
	}
	if _, err = os.Stat(s.state); !os.IsNotExist(err) {
		t.Fatal("inspection created state")
	}
	bind(t, s, &q)
	_, err = s.Start(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	q.OperationID, _ = identifier()
	q.RequestID, _ = identifier()
	q.Assignment.Outcome = "second task"
	if _, err = s.Start(LocalController(), q); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join(s.state, "HEAD"))
	if err != nil {
		t.Fatal(err)
	}
	v, err = s.Overview()
	if err != nil || len(v.Runs) != 2 || v.Generation != string(before) || v.Runs[0].RunID >= v.Runs[1].RunID {
		t.Fatal(v, err)
	}
	for _, r := range v.Runs {
		if r.Coverage != "known" || r.AssignedRole != "unknown" || r.LiveState != "unknown" || r.LastHostStatus != "not-observed" {
			t.Fatal(r)
		}
	}
	if err = os.WriteFile(filepath.Join(s.root, sel.Authority), []byte("Owner changes scope"), 0600); err != nil {
		t.Fatal(err)
	}
	v, err = s.Overview()
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range v.Runs {
		if r.Coverage != "stale" || len(r.Conditions) == 0 {
			t.Fatal(r)
		}
	}
	after, _ := os.ReadFile(filepath.Join(s.state, "HEAD"))
	if !bytes.Equal(before, after) {
		t.Fatal("inspection mutated operational state")
	}
	// The view does not silently drop a run when its canonical documents disappear.
	if err = os.Remove(filepath.Join(s.root, sel.Plan)); err != nil {
		t.Fatal(err)
	}
	v, err = s.Overview()
	if err != nil || len(v.Runs) != 2 {
		t.Fatal(v, err)
	}
	// A corrupt generation is a store failure, never an empty or successful view.
	if err = os.WriteFile(filepath.Join(s.state, "generations", string(before)+".json"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err = s.Overview()
	code(t, err, "reconciliation_required")
}

func TestOverviewHistoricalFactsDoNotBecomeLiveOrAccepted(t *testing.T) {
	s, _, q := fixture(t)
	bind(t, s, &q)
	result, err := s.Start(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	state, head, err := s.load()
	if err != nil {
		t.Fatal(err)
	}
	r := state.Runs[result.Run.ID]
	epoch, _ := identifier()
	newEpoch, _ := identifier()
	attempt, _ := identifier()
	p := ContextProvenance{Sources: []Source{}, AssignedRole: "Master", RoleBasis: "controller-launch", InputRevision: r.InputRevision, Host: HostContext{Status: "unavailable", InstructionStatus: "unavailable", InstructionSources: []string{}}, Skills: SkillInventory{Status: "unavailable", UseEvidence: "not-observed", Project: []SkillAvailability{}}}
	b, _ := json.Marshal(p)
	const when = "2026-01-01T00:00:00Z"
	// Frozen persisted fixture: inspection must retain the last epoch and gap, and
	// must not expose completed message contents or input/capability data.
	r.State = "verifying"
	r.Attempts = []Attempt{{ID: attempt, ThreadID: "thread", TurnID: "turn", Epoch: newEpoch, Status: "completed", Gap: true, Events: []Observation{
		{Kind: "context.observed", Text: string(b), Epoch: epoch, ReceivedAt: when},
		{Kind: "agentMessage", Text: "private completed message", ReceivedAt: when},
		{Kind: "tool.completed", ItemID: "tool-7", Status: "failed", ReceivedAt: when},
		{Kind: "turn.completed", Status: "completed", ReceivedAt: when},
	}}}
	r.Inputs = []OwnerInput{{Kind: "launch", Text: "private owner text", Payload: "private capability", AttemptID: attempt, State: "acknowledged", Skills: []SkillInput{{Name: "sdp", Path: "/project/Skills/sdp/SKILL.md", Digest: hash([]byte("skill"))}}}}
	state.Runs[r.ID] = r
	if err = s.save(state, head); err != nil {
		t.Fatal(err)
	}
	v, err := s.Overview()
	if err != nil {
		t.Fatal(err)
	}
	row := v.Runs[0]
	if row.RoutineState != "verifying" || row.LastHostStatus != "completed" || row.LiveState != "unknown" || row.OwnerAcceptance != "not-observed" || !row.ObservationGap || row.LastObservedAt != when || row.RoleEvidence != "historical-epoch" || row.AssignedRole != "Master" {
		t.Fatal(row)
	}
	if row.ToolCompletions != 1 || row.LastToolItemID != "tool-7" || row.LastToolStatus != "failed" || row.ChildMapping != "not-observed" || row.RequestedSkills != 1 || row.SkillAcceptance != "response-acknowledged" || row.SkillCompliance != "not-observed" {
		t.Fatal(row)
	}
	b, _ = json.Marshal(v)
	if bytes.Contains(b, []byte("private")) {
		t.Fatal("overview leaked input/event content")
	}
}
