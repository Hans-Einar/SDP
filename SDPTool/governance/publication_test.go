package governance

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func publicationFixture(t *testing.T) (*Service, Result, PublicationRequest) {
	return publicationFixtureAt(t, "", "")
}
func publicationFixtureAt(t *testing.T, root, state string) (*Service, Result, PublicationRequest) {
	t.Helper()
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	s, sel, q := fixture(t)
	if root != "" {
		if _, e := os.Stat(root); !os.IsNotExist(e) {
			t.Fatal("lab root must be absent", root)
		}
		e := filepath.WalkDir(s.root, func(path string, d os.DirEntry, e error) error {
			if e != nil {
				return e
			}
			rel, _ := filepath.Rel(s.root, path)
			dest := filepath.Join(root, rel)
			if d.IsDir() {
				return os.MkdirAll(dest, 0700)
			}
			b, e := os.ReadFile(path)
			if e != nil {
				return e
			}
			return os.WriteFile(dest, b, 0600)
		})
		if e != nil {
			t.Fatal(e)
		}
		s, e = New(root, filepath.Join(root, "SDP"), "trial", state)
		if e != nil {
			t.Fatal(e)
		}
	}
	evidenceFile(t, s, "SDP/KanBan/README.md", "# Publication fixture board\n\n| KB-SDP-001 | Change | active | Fixture |\n")
	oldCard := sel.Card
	sel.Card = "SDP/KanBan/active/#001--Change--fixture.md"
	q.Selection = sel
	evidenceFile(t, s, sel.Card, "# Work\n\n| Field | Value |\n| --- | --- |\n| id | KB-SDP-001 |\n| project | SDP |\n| type | Change |\n| created | 2026-10-24 |\n| source | authorized disposable fixture |\n| PlanId | MAINT-SDP-0001 |\n| CardState | in-progress |\n")
	if e := os.Remove(filepath.Join(s.root, oldCard)); e != nil {
		t.Fatal(e)
	}
	evidenceFile(t, s, sel.Plan, "# Plan\n\n| Field | Value |\n| --- | --- |\n| id | MAINT-SDP-0001 |\n| project | SDP |\n| state | active |\n| PlanType | MaintenancePlan |\n| BranchPolicy | current |\n| CommitPolicy | milestone |\n\nPreserve original authored plan.\n")
	event := func(id, subject, typ, path string, payload map[string]any) string {
		payload["schemaVersion"] = "0.2"
		payload["previousEventId"] = nil
		payload["from"] = nil
		payload["to"] = "active"
		payload["fromPath"] = nil
		payload["toPath"] = path
		payload["reason"] = "Authorized canonical publication fixture"
		payload["links"] = []string{}
		b, _ := json.Marshal(map[string]any{"schemaVersion": "1.0", "eventId": id, "eventType": typ, "occurredAt": "2026-10-24T00:00:00Z", "actor": "fixture", "commit": nil, "subjectId": subject, "payload": payload})
		return string(b) + "\n"
	}
	ledger := event("EVT-KB-SDP-000001", "KB-SDP-001", "x-kanban:created", strings.TrimPrefix(sel.Card, "SDP/KanBan/"), map[string]any{"projectId": "SDP"}) + event("EVT-PM-SDP-000001", "MAINT-SDP-0001", "x-management:created", "Maintenance/Plan.md", map[string]any{"kind": "Maintenance", "planType": "MaintenancePlan"})
	evidenceFile(t, s, "SDP/ProjectManagement/Ledger.ndjson", ledger)
	ctx, e := s.Context(sel)
	if e != nil {
		t.Fatal(e)
	}
	q.InputRevision = ctx.Revision
	bind(t, s, &q)
	r, e := s.Start(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.Transition(LocalController(), transition(r, "start"))
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.Transition(LocalController(), transition(r, "report"))
	if e != nil {
		t.Fatal(e)
	}
	evidenceFile(t, s, "target.txt", "verified fixture candidate\n")
	evidenceFile(t, s, "evidence/test.log", "Fixture test result: passed. This is protocol-test evidence, not native reviewer execution.\n")
	c, e := s.CaptureCandidate(r.Run.ID, nil)
	if e != nil {
		t.Fatal(e)
	}
	artifact, _ := os.ReadFile(filepath.Join(s.root, "evidence/test.log"))
	ev := EvidenceRequest{Schema: EvidenceSchema, ProjectRef: s.project, CheckoutID: r.CheckoutID, OperationID: newEvidenceID(t), RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, ReportID: r.OperationID, Candidate: c, Kind: "test", Verdict: "passed", Issuer: "protocol fixture", Note: "Synthetic adapter evidence; no native execution claim", Artifact: Source{Path: "evidence/test.log", SHA256: hash(artifact)}}
	if _, e = s.AttachEvidence(LocalController(), ev); e != nil {
		t.Fatal(e)
	}
	evidenceFile(t, s, "Skills/sdp-reviewer/SKILL.md", "Fixture review procedure")
	id := newEvidenceID(t)
	rq := ReviewRequest{Schema: ReviewSchema, ProjectRef: s.project, CheckoutID: r.CheckoutID, OperationID: id, RunID: r.Run.ID, ReviewID: id, ExpectedRunRevision: r.Run.Revision, Action: "prepare", Spec: &ReviewSpec{Instructions: "Read-only fixture review", Candidate: c, ReportID: r.OperationID, EvidenceIDs: []string{ev.OperationID}, Focus: "Check fixture", HostVersion: "0.164.0", MCPBinary: "/fixture/sdp-mcp", Skills: []SkillInput{{Name: "sdp-reviewer", Path: filepath.Join(s.root, "Skills/sdp-reviewer/SKILL.md"), Digest: hash([]byte("Fixture review procedure"))}}}}
	reviewApply(t, s, rq)
	reviewFinish(t, s, id, `{"verdict":"passed","summary":"Protocol fixture review passed; not a native model observation","findings":[]}`)
	return s, r, PublicationRequest{Schema: PublicationSchema, ProjectRef: s.project, CheckoutID: r.CheckoutID, OperationID: newEvidenceID(t), RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, ReviewID: id, System: "TEST", Summary: "Publish the tested fixture result", Remaining: "Owner acceptance and broader Session completion remain open"}
}
func mustPreparePublication(t *testing.T, s *Service, q PublicationRequest) PublicationView {
	t.Helper()
	v, e := s.PreparePublication(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	return v
}
func TestPublicationPreservesHistoryAndRetries(t *testing.T) {
	s, r, q := publicationFixture(t)
	old, _ := os.ReadFile(filepath.Join(s.root, "SDP/ProjectManagement/Ledger.ndjson"))
	session, _ := os.ReadFile(filepath.Join(s.root, r.Run.Selection.Session))
	prepared := mustPreparePublication(t, s, q)
	if !prepared.Applicable || prepared.Publication.Status != "prepared" || len(prepared.Publication.Changes) != 6 {
		t.Fatal(prepared)
	}
	if _, e := os.Stat(filepath.Join(s.root, prepared.Publication.EvidencePath)); !os.IsNotExist(e) {
		t.Fatal("prepare wrote project")
	}
	_, e := s.ApplyPublication(Caller{}, q.OperationID)
	code(t, e, "unauthorized")
	done, e := s.ApplyPublication(LocalController(), q.OperationID)
	if e != nil || !done.Applicable || done.Publication.Status != "published" {
		t.Fatal(done, e)
	}
	repeat, e := s.ApplyPublication(LocalController(), q.OperationID)
	if e != nil || !reflect.DeepEqual(done, repeat) {
		t.Fatal("repeat differs", e)
	}
	reprepare, e := s.PreparePublication(LocalController(), q)
	if e != nil || !reflect.DeepEqual(done, reprepare) {
		t.Fatal("prepare retry differs", e)
	}
	q.Summary = "different"
	_, e = s.PreparePublication(LocalController(), q)
	code(t, e, "idempotency_conflict")
	b, _ := os.ReadFile(filepath.Join(s.root, "SDP/ProjectManagement/Ledger.ndjson"))
	if !bytes.HasPrefix(b, old) || bytes.Count(b, []byte(done.Publication.EventIDs[0])) != 1 {
		t.Fatal("ledger history lost/duplicated")
	}
	b, _ = os.ReadFile(filepath.Join(s.root, r.Run.Selection.Session))
	if !bytes.HasPrefix(b, session) {
		t.Fatal("authored session rewritten")
	}
	view, e := s.Read(r.Run.ID)
	if e != nil || view.Run.State != "verifying" || view.Run.Revision != r.Run.Revision {
		t.Fatal("publication closed run", e)
	}
	evidenceFile(t, s, "target.txt", "later change")
	historical, e := s.PublicationRead(q.OperationID)
	if e != nil || historical.Applicable || historical.Publication.Status != "published" {
		t.Fatal("lost historical publication", e)
	}
	if file := os.Getenv("PG_PUBLICATION_WIRE"); file != "" {
		cases := []map[string]any{{"definition": "PublicationRequest", "value": done.Publication.Request}, {"definition": "PublicationView", "value": prepared}, {"definition": "PublicationView", "value": done}, {"definition": "PublicationView", "value": historical}}
		b, _ := json.MarshalIndent(cases, "", "  ")
		if e = os.WriteFile(file, b, 0600); e != nil {
			t.Fatal(e)
		}
	}
}
func TestPublicationCrashRecovery(t *testing.T) {
	for _, point := range []string{"publication-after-intent", "publication-after-file:0", "publication-after-file:1", "publication-after-file:2", "publication-after-file:3", "publication-after-file:4", "publication-after-file:5", "final-before-head", "final-after-head"} {
		t.Run(point, func(t *testing.T) {
			s, r, q := publicationFixture(t)
			v := mustPreparePublication(t, s, q)
			boom := errors.New("simulated process loss")
			fileEffects := false
			s.fault = func(stage string) error {
				if stage == "publication-after-file:5" {
					fileEffects = true
				}
				if stage == point || fileEffects && stage == strings.TrimPrefix(point, "final-") {
					return boom
				}
				return nil
			}
			_, e := s.ApplyPublication(LocalController(), q.OperationID)
			if !errors.Is(e, boom) {
				t.Fatal("missing fault", e)
			}
			restarted, e := New(s.root, s.area, s.project, s.state)
			if e != nil {
				t.Fatal(e)
			}
			if point != "final-after-head" {
				_, e = restarted.Transition(LocalController(), transition(r, "cancel"))
				code(t, e, "reconciliation_required")
			}
			done, e := restarted.ApplyPublication(LocalController(), q.OperationID)
			if e != nil || done.Publication.Status != "published" || !done.Applicable {
				t.Fatal(done, e)
			}
			for _, c := range v.Publication.Changes {
				b, e := os.ReadFile(filepath.Join(s.root, c.Path))
				if e != nil || hash(b) != c.After {
					t.Fatal("not exact", c.Path, e)
				}
			}
		})
	}
}
func TestPublicationConflictAndUnsafeInputs(t *testing.T) {
	for _, mode := range []string{"source", "test", "authority", "ledger", "partial-edit", "postimage-removed", "symlink", "active-review", "failed-test", "incomplete-history"} {
		t.Run(mode, func(t *testing.T) {
			s, _, q := publicationFixture(t)
			if mode == "incomplete-history" {
				evidenceFile(t, s, "SDP/ProjectManagement/Ledger.ndjson", `{"eventId":"fixture"}`+"\n")
				_, e := s.PreparePublication(LocalController(), q)
				if e == nil {
					t.Fatal("accepted invalid history")
				}
				return
			}
			if mode == "failed-test" {
				st, _, _ := s.load()
				review := st.Reviews[q.ReviewID]
				ev := *st.Receipts[review.Spec.EvidenceIDs[0]].Result.Evidence
				ev.Request.OperationID = newEvidenceID(t)
				ev.Request.Verdict = "failed"
				if _, e := s.AttachEvidence(LocalController(), ev.Request); e != nil {
					t.Fatal(e)
				}
				_, e := s.PreparePublication(LocalController(), q)
				code(t, e, "evidence_missing")
				return
			}
			if mode == "active-review" {
				st, _, _ := s.load()
				d := st.Reviews[q.ReviewID]
				rq := ReviewRequest{Schema: ReviewSchema, ProjectRef: s.project, CheckoutID: q.CheckoutID, OperationID: newEvidenceID(t), RunID: q.RunID, ExpectedRunRevision: q.ExpectedRevision, Action: "prepare", Spec: &d.Spec}
				rq.ReviewID = rq.OperationID
				reviewApply(t, s, rq)
				_, e := s.PreparePublication(LocalController(), q)
				code(t, e, "conflict")
				return
			}
			p := mustPreparePublication(t, s, q).Publication
			if mode == "partial-edit" || mode == "postimage-removed" {
				s.fault = func(stage string) error {
					if stage == "publication-after-file:1" {
						return errors.New("stop")
					}
					return nil
				}
				s.ApplyPublication(LocalController(), q.OperationID)
				s.fault = nil
			}
			switch mode {
			case "source":
				evidenceFile(t, s, "target.txt", "drift")
			case "test":
				evidenceFile(t, s, "evidence/test.log", "drift")
			case "authority":
				evidenceFile(t, s, "SDP/authority.md", "revoked")
			case "ledger":
				evidenceFile(t, s, "SDP/ProjectManagement/Ledger.ndjson", "changed\n")
			case "partial-edit":
				evidenceFile(t, s, p.Changes[1].Path, "external authored change")
			case "postimage-removed":
				os.Remove(filepath.Join(s.root, p.Changes[0].Path)) // absent is original preimage; recovery must restore once
			case "symlink":
				path := filepath.Join(s.root, p.Changes[1].Path)
				os.Remove(path)
				if e := os.Symlink(filepath.Join(s.root, "unrelated.txt"), path); e != nil {
					t.Fatal(e)
				}
			}
			v, e := s.ApplyPublication(LocalController(), q.OperationID)
			if mode == "postimage-removed" {
				if e != nil || !v.Applicable {
					t.Fatal(v, e)
				}
				return
			}
			if e == nil {
				t.Fatal("unsafe publication succeeded")
			}
			if mode == "partial-edit" {
				b, _ := os.ReadFile(filepath.Join(s.root, p.Changes[1].Path))
				if string(b) != "external authored change" {
					t.Fatal("overwrote external change")
				}
			}
		})
	}
}
func TestPublicationPreparedAbandonAndSharedProjectClaim(t *testing.T) {
	s, _, q := publicationFixture(t)
	v := mustPreparePublication(t, s, q)
	_, e := s.AbandonPublication(Caller{}, q.OperationID)
	code(t, e, "unauthorized")
	a, e := s.AbandonPublication(LocalController(), q.OperationID)
	if e != nil || a.Applicable || a.Publication.Status != "abandoned" {
		t.Fatal(a, e)
	}
	_, e = s.ApplyPublication(LocalController(), q.OperationID)
	code(t, e, "conflict")
	q.OperationID = newEvidenceID(t)
	v = mustPreparePublication(t, s, q)
	dir, _ := s.publicationLockDir()
	if e = mkdirDurable(dir); e != nil {
		t.Fatal(e)
	}
	if e = writeNew(filepath.Join(dir, "pending"), []byte("another-store\noperation\n")); e != nil {
		t.Fatal(e)
	}
	_, e = s.ApplyPublication(LocalController(), q.OperationID)
	code(t, e, "reconciliation_required")
	for _, ch := range v.Publication.Changes {
		b, e := os.ReadFile(filepath.Join(s.root, ch.Path))
		if ch.Before == "absent" {
			if !os.IsNotExist(e) {
				t.Fatal("cross-store write")
			}
		} else if hash(b) != ch.Before {
			t.Fatal("cross-store changed file")
		}
	}
}
func TestReviewCoordinationAndLeaderDelivery(t *testing.T) {
	s, _, q := publicationFixture(t)
	v, e := s.Coordination(CoordinationRequest{})
	if e != nil || v.Schema != "sdp-governance-coordination/2" {
		t.Fatal(v, e)
	}
	var notice *ReviewNotice
	for _, i := range v.Items {
		if i.Review != nil {
			notice = i.Review
		}
	}
	if notice == nil || !notice.Qualified || notice.Verdict != "passed" {
		t.Fatal(notice)
	}
	config := leaderConfigFixture(t, s)
	req := leaderRequest(t, s, "prepare", "")
	req.Config = &config
	dispatch := leaderApply(t, s, req).Dispatch
	if !strings.Contains(dispatch.Payload, q.ReviewID) {
		t.Fatal("review omitted from leader context")
	}
	// A delivered terminal notice is not selected again solely because generation changed.
	notices, e := s.leaderNotices()
	if e != nil {
		t.Fatal(e)
	}
	delivered := LeaderState{Dispatches: []LeaderDispatch{*dispatch}}
	if _, e = s.prepareLeader(delivered, newEvidenceID(t), config); e == nil {
		t.Fatal("duplicate notice delivery")
	}
	evidenceFile(t, s, "target.txt", "changed")
	changed, e := s.leaderNotices()
	if e != nil || reflect.DeepEqual(changed, notices) {
		t.Fatal("applicability change lost", e)
	}
	if _, e = s.prepareLeader(delivered, newEvidenceID(t), config); e != nil {
		t.Fatal("changed review not delivered", e)
	}
	if path := os.Getenv("PG_PUBLICATION_COORDINATION_WIRE"); path != "" {
		b, _ := json.MarshalIndent(v, "", "  ")
		if e = os.WriteFile(path, b, 0600); e != nil {
			t.Fatal(e)
		}
	}
}

// Optional compiled-CLI exercise retains a disposable canonical lab and synthetic
// protocol review. It intentionally does not claim a fresh native model review.
func TestPublicationCompiledLab(t *testing.T) {
	root := os.Getenv("PG_PUBLICATION_LAB_ROOT")
	if root == "" {
		t.Skip("explicit lab root not selected")
	}
	state := os.Getenv("PG_PUBLICATION_LAB_STATE")
	binary := os.Getenv("PG_PUBLICATION_CONTROLLER")
	if state == "" || binary == "" {
		t.Fatal("explicit lab state and controller required")
	}
	s, _, q := publicationFixtureAt(t, root, state)
	input := filepath.Join(t.TempDir(), "request.json")
	b, _ := json.Marshal(q)
	os.WriteFile(input, b, 0600)
	outputs := []json.RawMessage{}
	for _, args := range [][]string{{"prepare", "--input", input}, {"inspect", "--publication", q.OperationID}, {"apply", "--publication", q.OperationID}, {"apply", "--publication", q.OperationID}, {"inspect", "--publication", q.OperationID}} {
		argv := append([]string{"publication", args[0], "--project", root, "--state", state}, args[1:]...)
		cmd := exec.Command(binary, argv...)
		b, e := cmd.CombinedOutput()
		if e != nil {
			t.Fatal(string(b), e)
		}
		var v PublicationView
		if e = Decode(b, &v); e != nil {
			t.Fatal(e)
		}
		outputs = append(outputs, json.RawMessage(b))
	}
	v, e := s.PublicationRead(q.OperationID)
	if e != nil || v.Publication.Status != "published" || !v.Applicable {
		t.Fatal(v, e)
	}
	if file := os.Getenv("PG_PUBLICATION_LAB_OUTPUT"); file != "" {
		b, _ = json.MarshalIndent(outputs, "", "  ")
		os.WriteFile(file, b, 0600)
	}
	fmt.Println("PASS compiled publication prepare/inspect/apply/retry; synthetic review provenance:", q.OperationID)
}

func TestPublicationCanonicalRecordValidation(t *testing.T) {
	for _, mode := range []string{"unknown-payload", "index", "unresolved-link", "typed-plan", "namespace", "counter"} {
		t.Run(mode, func(t *testing.T) {
			s, _, _ := publicationFixture(t)
			files, e := s.publicationFiles()
			if e != nil {
				t.Fatal(e)
			}
			events, e := publicationEvents(files["SDP/ProjectManagement/Ledger.ndjson"])
			if e != nil {
				t.Fatal(e)
			}
			p := events[1]["payload"].(map[string]any)
			switch mode {
			case "unknown-payload":
				p["ownerApproval"] = true
			case "index":
				files["SDP/KanBan/README.md"] = []byte("missing card")
			case "unresolved-link":
				p["links"] = []string{"MAINT-SDP-9999"}
			case "typed-plan":
				p["planType"] = "AnythingPlan"
			case "namespace":
				files["SDP/KanBan/active/#001--Change--fixture.md"] = bytes.ReplaceAll(files["SDP/KanBan/active/#001--Change--fixture.md"], []byte("| project | SDP |"), []byte("| project | WRONG |"))
			case "counter":
				events[1]["eventId"] = "EVT-PM-SDP-malformed"
			}
			data := []byte{}
			for _, ev := range events {
				b, _ := json.Marshal(ev)
				data = append(data, append(b, '\n')...)
			}
			if e = validatePublicationHistory(data, files); e == nil {
				t.Fatal("accepted invalid canonical record set")
			}
		})
	}
	if root := os.Getenv("PG_PUBLICATION_INSPECT_ROOT"); root != "" {
		s, e := New(root, filepath.Join(root, "SDP"), "SDP", filepath.Join(t.TempDir(), "state"))
		if e != nil {
			t.Fatal(e)
		}
		files, e := s.publicationFiles()
		if e != nil {
			t.Fatal("actual repository record set", e)
		}
		t.Logf("actual repository canonical files validated: %d", len(files))
	}
}

func TestPublicationIntentSaveFailureCanRecoverOrAbandon(t *testing.T) {
	for _, abandon := range []bool{false, true} {
		t.Run(fmt.Sprint(abandon), func(t *testing.T) {
			s, _, q := publicationFixture(t)
			p := mustPreparePublication(t, s, q).Publication
			s.fault = func(stage string) error {
				if stage == "before-head" {
					return errors.New("intent persistence lost")
				}
				return nil
			}
			if _, e := s.ApplyPublication(LocalController(), q.OperationID); e == nil {
				t.Fatal("expected intent fault")
			}
			s.fault = nil
			for _, ch := range p.Changes {
				b, e := os.ReadFile(filepath.Join(s.root, ch.Path))
				if ch.Before == "absent" {
					if !os.IsNotExist(e) {
						t.Fatal("effect before intent")
					}
				} else if hash(b) != ch.Before {
					t.Fatal("effect before intent")
				}
			}
			if abandon {
				v, e := s.AbandonPublication(LocalController(), q.OperationID)
				if e != nil || v.Publication.Status != "abandoned" {
					t.Fatal(v, e)
				}
			} else {
				v, e := s.ApplyPublication(LocalController(), q.OperationID)
				if e != nil || v.Publication.Status != "published" {
					t.Fatal(v, e)
				}
			}
		})
	}
}
