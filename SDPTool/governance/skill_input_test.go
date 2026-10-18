package governance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSkillInputExactCorrelation(t *testing.T) {
	v := OwnerInput{ID: "client", Payload: "exact", Skills: []SkillInput{{Name: "sdp", Path: "/project/Skills/sdp/SKILL.md", Digest: hash([]byte("skill"))}}}
	b, _ := json.Marshal(TurnInputs(v.Payload, v.Skills))
	if ObservedInputDigest(b) != v.WireDigest() {
		t.Fatal("exact input failed")
	}
	for _, raw := range []string{
		`[{"type":"text","text":"exact"}]`,
		`[{"type":"text","text":"exact"},{"type":"skill","name":"other","path":"/project/Skills/sdp/SKILL.md"}]`,
		`[{"type":"text","text":"exact"},{"type":"skill","name":"sdp","path":"/elsewhere"}]`,
		`[{"type":"skill","name":"sdp","path":"/project/Skills/sdp/SKILL.md"},{"type":"text","text":"exact"}]`,
		`[{"type":"text","text":"exact","text_elements":[{}]},{"type":"skill","name":"sdp","path":"/project/Skills/sdp/SKILL.md"}]`,
		`[{"type":"text","text":"exact","extra":true},{"type":"skill","name":"sdp","path":"/project/Skills/sdp/SKILL.md"}]`,
		`[{"type":"text","text":"exact","text_elements":null}]`,
		`[{"type":"text","text":"wrong","text":"exact"}]`,
	} {
		if ObservedInputDigest([]byte(raw)) == v.WireDigest() {
			t.Fatal("false match", raw)
		}
	}
	normalized := strings.Replace(string(b), `"text":"exact"`, `"text":"exact","text_elements":[]`, 1)
	if ObservedInputDigest([]byte(normalized)) != v.WireDigest() {
		t.Fatal("measured empty metadata refused")
	}
	o := Observation{Kind: "userMessage", ClientID: v.ID, TurnID: "turn", Text: v.Payload, InputDigest: v.WireDigest()}
	if !inputObserved(Attempt{Events: []Observation{o}}, v, "turn") {
		t.Fatal("missing exact proof")
	}
	if inputObserved(Attempt{Events: []Observation{o, o}}, v, "turn") {
		t.Fatal("duplicate proof accepted")
	}
	o.InputDigest = ""
	if inputObserved(Attempt{Events: []Observation{o}}, v, "turn") {
		t.Fatal("legacy text proves skill")
	}
	o.InputDigest = "unavailable"
	v.Skills = nil
	if inputObserved(Attempt{Events: []Observation{o}}, v, "turn") {
		t.Fatal("unsupported content became legacy text proof")
	}
}

func TestSelectedSkillPersistenceDriftAndRecovery(t *testing.T) {
	s, _, q := fixture(t)
	bind(t, s, &q)
	r, e := s.Start(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	path := filepath.Join(s.root, "Skills/sdp/SKILL.md")
	os.MkdirAll(filepath.Dir(path), 0700)
	os.WriteFile(path, []byte("skill bytes"), 0600)
	id, _ := identifier()
	at, _ := identifier()
	epoch, _ := identifier()
	v := OwnerInput{ID: id, Kind: "launch", Text: "owner exact", Payload: "$sdp\nassignment", Skills: []SkillInput{{Name: "sdp", Path: path, Digest: hash([]byte("skill bytes"))}}, AttemptID: at, Epoch: epoch, ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	input := func(action string) (Result, error) {
		op, _ := identifier()
		return s.RecordInput(LocalController(), InputRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, OperationID: op, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: action, Input: v})
	}
	r, e = input("prepare")
	if e != nil {
		t.Fatal(e)
	}
	v.State = "prepared"
	os.WriteFile(path, []byte("changed"), 0600)
	_, e = input("claim")
	code(t, e, "stale_revision")
	os.WriteFile(path, []byte("skill bytes"), 0600)
	old := v.Skills[0]
	v.Skills[0].Name = "different"
	_, e = input("claim")
	code(t, e, "conflict")
	v.Skills[0] = old
	r, e = input("claim")
	if e != nil {
		t.Fatal(e)
	}
	v.State = "sending"
	host := func(action string, o Observation) (Result, error) {
		op, _ := identifier()
		o.Epoch = epoch
		o.ReceivedAt = v.ReceivedAt
		return s.RecordHost(LocalController(), HostRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, OperationID: op, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: action, AttemptID: at, HostVersion: "0.161.0", ThreadID: "thread", Observation: o})
	}
	r, e = host("bind", Observation{Kind: "bound", Status: "idle"})
	if e != nil {
		t.Fatal(e)
	}
	v.ThreadID = "thread"
	r, e = input("uncertain")
	if e != nil {
		t.Fatal(e)
	}
	v.State = "uncertain"
	v.TurnID = "turn"
	r, e = host("reconcile", Observation{Kind: "reconciled", TurnID: "turn", Status: "completed"})
	if e != nil {
		t.Fatal(e)
	}
	_, e = input("resolve")
	code(t, e, "reconciliation_required")
	o := Observation{Kind: "userMessage", ItemID: "item", ClientID: id, TurnID: "turn", Text: v.Payload, InputDigest: v.WireDigest(), Attachments: []Attachment{{Kind: "skill", Reference: path, Availability: "unknown"}}}
	r, e = host("observe", o)
	if e != nil {
		t.Fatal(e)
	}
	changed := o
	changed.InputDigest = hash([]byte("wrong"))
	_, e = host("observe", changed)
	code(t, e, "observation_gap")
	// Historical acceptance survives later local drift; it cannot claim current bytes.
	os.WriteFile(path, []byte("later"), 0600)
	r, e = input("resolve")
	if e != nil {
		t.Fatal(e)
	}
	stored, e := s.OwnerInput(LocalController(), r.Run.ID, id)
	if e != nil || stored.State != "acknowledged" || len(stored.Skills) != 1 {
		t.Fatal(stored, e)
	}
	view, e := s.Provenance(r.Run.ID)
	if e != nil || view.Invocation == nil || view.Invocation.Acceptance != "host-message-observed" || view.Invocation.Compliance != "not-observed" {
		t.Fatal(view, e)
	}
	if len(r.Run.Inputs[0].Skills) != 0 {
		t.Fatal("unbounded summary exposed content")
	}
	o.Epoch = epoch
	o.ReceivedAt = v.ReceivedAt
	if dir := os.Getenv("PGI_SCHEMA_EVIDENCE"); dir != "" {
		b, _ := json.Marshal([]map[string]any{{"definition": "OwnerInput", "value": stored}, {"definition": "ProvenanceView", "value": view}, {"definition": "Observation", "value": o}})
		os.WriteFile(filepath.Join(dir, "skill-wire-cases.json"), b, 0600)
	}
}
func TestSkillSnapshotConfinement(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "SKILL.md")
	os.WriteFile(outside, []byte("outside"), 0600)
	os.Symlink(outside, filepath.Join(root, "SKILL.md"))
	s := []SkillInput{{Name: "sdp", Path: filepath.Join(root, "SKILL.md"), Digest: hash([]byte("outside"))}}
	if CheckSkillSnapshots(root, s) == nil {
		t.Fatal("external symlink accepted")
	}
	s[0].Path = outside
	if CheckSkillSnapshots(root, s) == nil {
		t.Fatal("external path accepted")
	}
	s[0].Name = "sdp\nowner"
	if ValidSkillInputs(s) {
		t.Fatal("marker injection")
	}
}
