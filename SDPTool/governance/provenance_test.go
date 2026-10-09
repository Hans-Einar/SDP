package governance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProvenanceIsBoundEvidenceNotRoleAuthority(t *testing.T) {
	s, _, q := fixture(t)
	bind(t, s, &q)
	r, err := s.Start(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	view, err := s.Provenance(r.Run.ID)
	if err != nil || view.Status != "unavailable" || view.Context != nil {
		t.Fatal(view, err)
	}
	context, err := s.Context(r.Run.Selection)
	if err != nil {
		t.Fatal(err)
	}
	epoch, _ := identifier()
	attempt, _ := identifier()
	op, _ := identifier()
	h := HostRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: op, RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: "bind", AttemptID: attempt, HostVersion: "0.161.0", ThreadID: "thread", Observation: Observation{Kind: "bound", Epoch: epoch, ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano), Status: "idle"}}
	r, err = s.RecordHost(LocalController(), h)
	if err != nil {
		t.Fatal(err)
	}
	p := ContextProvenance{AssignedRole: "Master", RoleBasis: "controller-launch", InputRevision: r.Run.InputRevision, Sources: context.Sources, Host: HostContext{Status: "observed", SessionID: "session", NativeRole: "claimed-host-role", InstructionStatus: "unavailable", InstructionSources: []string{}}, Skills: SkillInventory{Status: "observed", UseEvidence: "not-observed", Project: []SkillAvailability{{Name: "sdp", Path: "Skills/sdp/SKILL.md", Enabled: true, Digest: "unavailable"}}}}
	h.Action = "observe"
	h.ExpectedRevision = r.Run.Revision
	h.OperationID, _ = identifier()
	h.Observation.Kind = "context.observed"
	h.Observation.Status = ""
	h.Observation.ItemID = "launch"
	b, _ := json.Marshal(p)
	h.Observation.Text = string(b)
	agent, _ := AssignedAgent(r.Run.ID, attempt, epoch)
	_, err = s.RecordHost(agent, h)
	code(t, err, "unauthorized")
	bad := p
	bad.AssignedRole = "owner"
	b, _ = json.Marshal(bad)
	h.Observation.Text = string(b)
	_, err = s.RecordHost(LocalController(), h)
	code(t, err, "invalid_request")
	b, _ = json.Marshal(p)
	h.Observation.Text = string(b)
	r, err = s.RecordHost(LocalController(), h)
	if err != nil {
		t.Fatal(err)
	}
	view, err = s.Provenance(r.Run.ID)
	if err != nil || view.Status != "observed" || view.Context.Host.NativeRole != "claimed-host-role" || view.Context.Skills.UseEvidence != "not-observed" {
		t.Fatal(view, err)
	}
	again, err := s.RecordHost(LocalController(), h)
	if err != nil || again.Run.Revision != r.Run.Revision {
		t.Fatal(again, err)
	}
	_, err = s.Transition(agent, transition(r, "complete"))
	code(t, err, "unauthorized")
	// A new epoch does not relabel the captured context as a fresh measurement.
	h.Action = "reconcile"
	h.OperationID, _ = identifier()
	h.ExpectedRevision = r.Run.Revision
	next, _ := identifier()
	h.Observation = Observation{Kind: "reconciled", Epoch: next, Status: "idle", ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	r, err = s.RecordHost(LocalController(), h)
	if err != nil {
		t.Fatal(err)
	}
	view, err = s.Provenance(r.Run.ID)
	if err != nil || view.Status != "historical-epoch" || view.CaptureEpoch != epoch || view.CurrentEpoch != next {
		t.Fatal(view, err)
	}
	after, err := s.Read(r.Run.ID)
	if err != nil || after.Run.Revision != r.Run.Revision {
		t.Fatal("read mutated run", after, err)
	}
	if dir := os.Getenv("PGI_SCHEMA_EVIDENCE"); dir != "" {
		b, _ := json.Marshal([]map[string]any{{"definition": "ProvenanceView", "value": view}, {"definition": "ContextProvenance", "value": p}})
		if err = os.WriteFile(filepath.Join(dir, "provenance-wire-cases.json"), b, 0600); err != nil {
			t.Fatal(err)
		}
	}
}
func TestProvenanceRejectsFalseActivationAndInvalidRecords(t *testing.T) {
	base := ContextProvenance{AssignedRole: "Master", RoleBasis: "controller-launch", InputRevision: hash([]byte("source")), Host: HostContext{Status: "unavailable", InstructionStatus: "unavailable"}, Skills: SkillInventory{Status: "unavailable", UseEvidence: "not-observed"}}
	if !base.Valid() {
		t.Fatal("unknown metadata must be representable")
	}
	for _, edit := range []func(*ContextProvenance){func(p *ContextProvenance) { p.Skills.UseEvidence = "compliant" }, func(p *ContextProvenance) { p.Host.InstructionSources = []string{"relative/path"} }, func(p *ContextProvenance) { p.Host.InstructionSources = make([]string, 65) }, func(p *ContextProvenance) {
		p.Skills.Project = []SkillAvailability{{Name: "sdp", Path: "../outside", Digest: "unavailable"}}
	}, func(p *ContextProvenance) { p.Skills.ExternalCount = -1 }} {
		p := base
		edit(&p)
		if p.Valid() {
			t.Fatal("invalid provenance accepted", p)
		}
	}
}
