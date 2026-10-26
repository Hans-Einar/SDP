package governance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func successorFixture(t *testing.T, s *Service, id string) SessionContext {
	t.Helper()
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	i := st.Sessions[id]
	b := i.Context.Binding
	b.Session = SessionRecordRef{"SESSION-TEST-0002", "SDP/Sessions/session-#0002--Successor.md"}
	b.Plan = SessionRecordRef{"MAINT-SDP-0002", "SDP/Maintenance/successor.md"}
	b.Steps = append([]SessionStepBinding{}, b.Steps...)
	b.Steps[0].Cards = []SessionCardRef{{"KB-SDP-001", "whole"}}
	path := "SDP/Sessions/execution/SESSION-TEST-0002.json"
	raw, _ := json.Marshal(b)
	evidenceFile(t, s, path, string(raw))
	evidenceFile(t, s, b.Session.Path, "# Successor\n\n| Field | Value |\n| --- | --- |\n| Session reference | SESSION-TEST-0002 |\n| Status | active |\n\n## S1\n\n[Binding](execution/SESSION-TEST-0002.json)\n")
	original := mustFile(t, filepath.Join(s.root, i.Context.Binding.Plan.Path))
	body := bytes.ReplaceAll(original, []byte(i.Context.Binding.Plan.ID), []byte(b.Plan.ID))
	body, e = replaceMetadata(body, "state", "planned")
	if e != nil {
		t.Fatal(e)
	}
	evidenceFile(t, s, b.Plan.Path, string(body))
	events, e := publicationEvents(mustFile(t, filepath.Join(s.root, managementLedger)))
	if e != nil {
		t.Fatal(e)
	}
	var old map[string]any
	for _, ev := range events {
		if ev["subjectId"] == i.Context.Binding.Plan.ID {
			old = ev
		}
	}
	p := map[string]any{"schemaVersion": "0.2", "kind": "Maintenance", "planType": "MaintenancePlan", "from": nil, "fromPath": nil, "to": "planned", "toPath": "Maintenance/successor.md", "previousEventId": nil, "reason": "Register manual successor plan", "links": []string{}}
	eid, e := nextPublicationEvent(events, "EVT-PM-SDP-")
	if e != nil {
		t.Fatal(e)
	}
	ev := map[string]any{"schemaVersion": "1.0", "eventId": eid, "eventType": "x-management:created", "occurredAt": old["occurredAt"], "actor": "fixture-controller", "commit": nil, "subjectId": b.Plan.ID, "payload": p}
	raw, _ = json.Marshal(ev)
	evidenceFile(t, s, managementLedger, appendRecord(mustFile(t, filepath.Join(s.root, managementLedger)), string(raw)+"\n"))
	ctx, e := s.SessionSuccessorContext(SessionSuccessorRequest{id, path})
	if e != nil {
		t.Fatal(e)
	}
	return ctx
}
func transferFixture(t *testing.T) (*Service, SessionAdvanceRequest) {
	t.Helper()
	return transferFixtureAt(t, "", "")
}
func transferFixtureAt(t *testing.T, root, state string) (*Service, SessionAdvanceRequest) {
	t.Helper()
	s, d := dispositionFixture(t, root, state, func(_ *testing.T, _ *Service, b *SessionBinding) { b.Steps[0].Cards[0].Coverage = "contribution" })
	completeFixtureDisposition(t, s, d)
	ctx := successorFixture(t, s, d.InstanceID)
	q := advanceRequest(t, s, d.InstanceID, "transfer", "", "disposition", d.OperationID)
	q.Successor = &SessionSuccessor{ctx.BindingPath, ctx.InputRevision}
	return s, q
}
func TestSessionTransferAndFreshSuccessor(t *testing.T) {
	s, q := transferFixture(t)
	before := mustFile(t, filepath.Join(s.root, managementLedger))
	commitAdvance(t, s, q)
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	if st.Schema != sessionTransferStateSchema {
		t.Fatal(st.Schema)
	}
	i := st.Sessions[q.InstanceID]
	v, e := s.ReadSession(i.ID)
	if e != nil || v.State != "completed" || v.Coverage != "current" {
		t.Fatal(v, e)
	}
	if !bytes.HasPrefix(mustFile(t, filepath.Join(s.root, managementLedger)), before) {
		t.Fatal("rewrote history")
	}
	for _, p := range []string{i.Context.Binding.Plan.Path, i.Context.Binding.Session.Path} {
		m, _ := metadata(mustFile(t, filepath.Join(s.root, p)))
		if m["state"] != "completed" && m["Status"] != "completed" {
			t.Fatal(m)
		}
	}
	card := mustFile(t, filepath.Join(s.root, "SDP/KanBan/active/#001--Change--fixture.md"))
	m, _ := metadata(card)
	if m["CardState"] != "ready" || m["PlanId"] != "MAINT-SDP-0002" || !bytes.Contains(card, []byte("Other planned contributions remain open")) {
		t.Fatal(m, string(card))
	}
	ctx, e := s.SessionContext(SessionContextRequest{q.Successor.BindingPath})
	if e != nil {
		t.Fatal(e)
	}
	adopt := SessionAdoptRequest{SessionSchema, s.project, st.CheckoutID, newTestID(), ctx.BindingPath, ctx.InputRevision, "Adopt explicit remaining work without inherited evidence"}
	p, e := s.AdoptSession(LocalController(), adopt)
	if e != nil {
		t.Fatal(e)
	}
	done, e := s.ResumeSession(LocalController(), sessionResume(p, adopt))
	if e != nil {
		t.Fatal(e)
	}
	sel := SessionSelectRequest{SessionSchema, s.project, st.CheckoutID, newTestID(), done.InstanceID, done.Revision, "S1", "Select transferred remaining work"}
	p, e = s.SelectSession(LocalController(), sel)
	if e != nil {
		t.Fatal(e)
	}
	done, e = s.ResumeSession(LocalController(), SessionResumeRequest{SessionSchema, s.project, st.CheckoutID, newTestID(), done.InstanceID, p.Revision, "Commit successor selection"})
	if e != nil {
		t.Fatal(e)
	}
	v, e = s.ReadSession(done.InstanceID)
	if e != nil || v.Steps[0].Stage != "admission" || v.Steps[0].RunID != "" {
		t.Fatal(v, e)
	}
	st, _, e = s.load()
	if e != nil || st.Schema != sessionTransferStateSchema {
		t.Fatal(st.Schema, e)
	}
	// The terminal predecessor and all operation receipts remain readable after successor edits.
	v, e = s.ReadSession(q.InstanceID)
	if e != nil || v.State != "completed" {
		t.Fatal(v, e)
	}
	if _, e = s.AdvanceSession(LocalController(), q); e != nil {
		t.Fatal(e)
	}
}
func TestSessionTransferRefusals(t *testing.T) {
	for _, kind := range []string{"stale", "missing", "extra", "foreign-plan", "active-plan", "wrong-action", "untrusted", "downgraded", "tampered"} {
		t.Run(kind, func(t *testing.T) {
			s, q := transferFixture(t)
			before := mustFile(t, filepath.Join(s.state, "HEAD"))
			switch kind {
			case "stale":
				evidenceFile(t, s, q.Successor.BindingPath, string(mustFile(t, filepath.Join(s.root, q.Successor.BindingPath)))+"\n")
			case "missing", "extra":
				var b SessionBinding
				json.Unmarshal(mustFile(t, filepath.Join(s.root, q.Successor.BindingPath)), &b)
				b.Steps[0].Cards = []SessionCardRef{{"KB-SDP-999", "whole"}}
				raw, _ := json.Marshal(b)
				evidenceFile(t, s, q.Successor.BindingPath, string(raw))
			case "foreign-plan":
				p := "SDP/KanBan/active/#001--Change--fixture.md"
				b, _ := replaceMetadata(mustFile(t, filepath.Join(s.root, p)), "PlanId", "MAINT-SDP-9999")
				evidenceFile(t, s, p, string(b))
			case "active-plan":
				p := "SDP/Maintenance/successor.md"
				b, _ := replaceMetadata(mustFile(t, filepath.Join(s.root, p)), "state", "active")
				evidenceFile(t, s, p, string(b))
			case "wrong-action":
				q.Action = "complete"
			case "untrusted":
				_, e := s.AdvanceSession(Caller{}, q)
				if e == nil {
					t.Fatal("accepted untrusted")
				}
				return
			case "downgraded", "tampered":
				p, e := s.AdvanceSession(LocalController(), q)
				if e != nil {
					t.Fatal(e)
				}
				st, _, e := s.load()
				if e != nil {
					t.Fatal(e)
				}
				if kind == "downgraded" {
					st.Schema = sessionLifecycleStateSchema
				} else {
					i := st.Sessions[q.InstanceID]
					i.Advance.Successor.Binding.Plan.ID = "MAINT-SDP-9999"
					st.Sessions[q.InstanceID] = i
				}
				if e = validateSessionAdvances(st); e == nil {
					t.Fatal("accepted forged state", p)
				}
				return
			}
			if _, e := s.AdvanceSession(LocalController(), q); e == nil {
				t.Fatal("accepted", kind)
			}
			if !bytes.Equal(before, mustFile(t, filepath.Join(s.state, "HEAD"))) {
				t.Fatal("refusal wrote state")
			}
		})
	}
}
func TestSessionTransferEveryEffectRecovery(t *testing.T) {
	s, q := transferFixture(t)
	p, e := s.AdvanceSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	st, _, _ := s.load()
	count := len(st.Sessions[q.InstanceID].Advance.Effects.Changes)
	_ = p
	for n := 0; n < count; n++ {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			s, q := transferFixture(t)
			p, e := s.AdvanceSession(LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			fired := false
			s.fault = func(point string) error {
				if point == fmt.Sprintf("session-advance-after-file:%d", n) && !fired {
					fired = true
					return fmt.Errorf("simulated interruption")
				}
				return nil
			}
			resume := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Recover exact transfer"}
			if _, e = s.ResumeSession(LocalController(), resume); e == nil || !fired {
				t.Fatal("no fault", e)
			}
			fresh, e := New(s.root, s.area, s.project, s.state)
			if e != nil {
				t.Fatal(e)
			}
			if _, e = fresh.ResumeSession(LocalController(), resume); e != nil {
				t.Fatal(e)
			}
			v, e := fresh.ReadSession(q.InstanceID)
			if e != nil || v.State != "completed" {
				t.Fatal(v, e)
			}
			if _, e = fresh.SessionContext(SessionContextRequest{q.Successor.BindingPath}); e != nil {
				t.Fatal(e)
			}
			if _, e = os.Stat(filepath.Join(s.root, managementLedger) + ".pg-pending"); !os.IsNotExist(e) {
				t.Fatal("reservation remains", e)
			}
		})
	}
}

func TestSessionTransferGuardsAndPendingAdoption(t *testing.T) {
	t.Run("whole-card-needs-no-transfer", func(t *testing.T) {
		s, d := dispositionFixture(t, "", "", nil)
		completeFixtureDisposition(t, s, d)
		st, _, _ := s.load()
		if _, e := remainingSessionCards(st.Sessions[d.InstanceID]); e == nil {
			t.Fatal("invented remaining scope")
		}
	})
	t.Run("unfinished-step", func(t *testing.T) {
		s, d := dispositionFixture(t, "", "", func(t *testing.T, s *Service, b *SessionBinding) {
			b.Steps[0].Cards[0].Coverage = "contribution"
			next := b.Steps[0]
			next.ID = "second"
			next.DependsOn = []string{"S1"}
			b.Steps = append(b.Steps, next)
			evidenceFile(t, s, b.Session.Path, string(mustFile(t, filepath.Join(s.root, b.Session.Path)))+"\n## second\n")
		})
		completeFixtureDisposition(t, s, d)
		st, _, _ := s.load()
		if _, e := remainingSessionCards(st.Sessions[d.InstanceID]); e == nil {
			t.Fatal("abandoned unfinished step")
		}
	})
	for _, point := range []string{"session-advance-after-claim", "before-head", "after-head"} {
		t.Run(point, func(t *testing.T) {
			s, q := transferFixture(t)
			p, e := s.AdvanceSession(LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			fired := false
			s.fault = func(at string) error {
				if !fired && at == point {
					fired = true
					return fmt.Errorf("interruption")
				}
				return nil
			}
			resume := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Recover transfer commit boundary"}
			if _, e = s.ResumeSession(LocalController(), resume); e == nil || !fired {
				t.Fatal(e)
			}
			if point == "before-head" {
				// All canonical postimages are present, but ownership remains reserved until
				// the predecessor's operational commit. Adoption cannot race that boundary.
				ctx, e := s.SessionContext(SessionContextRequest{q.Successor.BindingPath})
				if e != nil {
					t.Fatal(e)
				}
				adopt := SessionAdoptRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), ctx.BindingPath, ctx.InputRevision, "Must wait for canonical transfer commitment"}
				if _, e = s.AdoptSession(LocalController(), adopt); e == nil {
					t.Fatal("adopted during reserved transfer")
				}
			}
			fresh, e := New(s.root, s.area, s.project, s.state)
			if e != nil {
				t.Fatal(e)
			}
			if _, e = fresh.ResumeSession(LocalController(), resume); e != nil {
				t.Fatal(e)
			}
		})
	}
	t.Run("changed-prepared-target", func(t *testing.T) {
		s, q := transferFixture(t)
		p, e := s.AdvanceSession(LocalController(), q)
		if e != nil {
			t.Fatal(e)
		}
		old := mustFile(t, filepath.Join(s.root, "SDP/Maintenance/successor.md"))
		evidenceFile(t, s, "SDP/Maintenance/successor.md", string(old)+"\nChanged goal\n")
		_, e = s.ResumeSession(LocalController(), SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Refuse changed successor"})
		if e == nil {
			t.Fatal("ignored edited successor")
		}
		if !bytes.Contains(mustFile(t, filepath.Join(s.root, "SDP/Maintenance/successor.md")), []byte("Changed goal")) {
			t.Fatal("overwrote edit")
		}
	})
}
