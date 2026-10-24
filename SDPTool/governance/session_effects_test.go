package governance

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
)

func selectionFixture(t *testing.T, root, state string, scopes ...string) (*Service, SessionSelectRequest) {
	t.Helper()
	s, c, q := sessionFixtureAt(t, root, state)
	// Rebuild the disposable authored fixture before adoption, preserving all history.
	old := "SDP/KanBan/active/#001--Change--fixture.md"
	back := "SDP/KanBan/backlog/#001--Change--fixture.md"
	evidenceFile(t, s, "SDP/KanBan/backlog/resource.txt", "Preserved sibling resource\n")
	card, _ := os.ReadFile(filepath.Join(s.root, old))
	card, _ = replaceMetadata(card, "CardState", "queued")
	evidenceFile(t, s, back, string(card)+"\n[Sibling resource](resource.txt)\n[Existing plan](../../Maintenance/Plan.md)\n[Self](%23001--Change--fixture.md)\n\n## Queue\n\nSelected after the fixture predecessor; explicit owner selection required.\n")
	if e := os.Remove(filepath.Join(s.root, old)); e != nil {
		t.Fatal(e)
	}
	evidenceFile(t, s, "SDP/KanBan/README.md", "# Board\n\n| KB-SDP-001 | Change | backlog | [Fixture](backlog/%23001--Change--fixture.md) |\n")
	plan, _ := os.ReadFile(filepath.Join(s.root, c.Binding.Plan.Path))
	plan, _ = replaceMetadata(plan, "state", "planned")
	evidenceFile(t, s, c.Binding.Plan.Path, string(plan)+"\n[Source card](../KanBan/backlog/%23001--Change--fixture.md)\n")
	events, _ := publicationEvents(mustFile(t, filepath.Join(s.root, managementLedger)))
	raw := []byte{}
	for _, ev := range events {
		p := ev["payload"].(map[string]any)
		if ev["subjectId"] == "KB-SDP-001" {
			p["to"] = "backlog"
			p["toPath"] = "backlog/#001--Change--fixture.md"
		} else {
			p["to"] = "planned"
		}
		b, _ := json.Marshal(ev)
		raw = append(raw, append(b, '\n')...)
	}
	evidenceFile(t, s, managementLedger, string(raw))
	if len(scopes) > 0 {
		c.Binding.Steps[0].WritablePaths = scopes
		raw, _ := json.Marshal(c.Binding)
		evidenceFile(t, s, q.BindingPath, string(raw))
	}
	c, e := s.SessionContext(SessionContextRequest{q.BindingPath})
	if e != nil {
		t.Fatal(e)
	}
	q.InputRevision = c.InputRevision
	adopted, e := s.AdoptSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	registered, e := s.ResumeSession(LocalController(), sessionResume(adopted, q))
	if e != nil {
		t.Fatal(e)
	}
	return s, SessionSelectRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newEvidenceID(t), registered.InstanceID, registered.Revision, "S1", "Owner selects the fixture step"}
}
func mustFile(t *testing.T, path string) []byte {
	t.Helper()
	b, e := os.ReadFile(path)
	if e != nil {
		t.Fatal(e)
	}
	return b
}
func selectionResume(q SessionSelectRequest, r SessionMutationResult) SessionResumeRequest {
	return SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), r.InstanceID, r.Revision, "Apply recorded selection effects"}
}
func newTestID() string { id, _ := identifier(); return id }
func TestSessionSelectionMovesAndReconciles(t *testing.T) {
	s, q := selectionFixture(t, "", "")
	history := mustFile(t, filepath.Join(s.root, managementLedger))
	prepared, e := s.SelectSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = os.Stat(filepath.Join(s.root, "SDP/KanBan/active/#001--Change--fixture.md")); !os.IsNotExist(e) {
		t.Fatal("prepare moved card")
	}
	repeat, e := s.SelectSession(LocalController(), q)
	if e != nil || repeat != prepared {
		t.Fatal("selection retry", e)
	}
	view, e := s.ReadSession(q.InstanceID)
	if e != nil || view.RecordEffects != "prepared" || view.Steps[0].Selected {
		t.Fatal(view, e)
	}
	request := selectionResume(q, prepared)
	done, e := s.ResumeSession(LocalController(), request)
	if e != nil {
		t.Fatal(e)
	}
	again, e := s.ResumeSession(LocalController(), request)
	if e != nil || again != done {
		t.Fatal(e)
	}
	if _, e = os.Stat(filepath.Join(s.root, "SDP/KanBan/backlog/#001--Change--fixture.md")); !os.IsNotExist(e) {
		t.Fatal("source still present")
	}
	card := mustFile(t, filepath.Join(s.root, "SDP/KanBan/active/#001--Change--fixture.md"))
	m, _ := metadata(card)
	if m["CardState"] != "ready" || strings.Count(string(card), "## PG step selection ") != 1 {
		t.Fatal(string(card))
	}
	if !strings.Contains(string(card), "](../backlog/resource.txt)") || !strings.Contains(string(card), "](../../Maintenance/Plan.md)") || !strings.Contains(string(card), "](%23001--Change--fixture.md)") {
		t.Fatal("moved card outgoing links changed target", string(card))
	}
	files, e := s.publicationFiles()
	if e != nil {
		t.Fatal(e)
	}
	if !strings.HasPrefix(string(files[managementLedger]), string(history)) {
		t.Fatal("history rewritten")
	}
	if strings.Contains(string(files["SDP/Maintenance/Plan.md"]), "../KanBan/backlog/") {
		t.Fatal("current link stale")
	}
	view, e = s.ReadSession(q.InstanceID)
	if e != nil || !view.Steps[0].Selected || view.Coverage != "current" || view.RecordEffects != "committed" || view.HostExecution != "not-dispatched" {
		t.Fatal(view, e)
	}
	// A canonical unrelated event is observed without accepting selected-doc edits.
	addUnrelatedCard(t, s)
	view, e = s.ReadSession(q.InstanceID)
	if e != nil || view.Coverage != "current" {
		t.Fatal("unrelated history blocked", view, e)
	}
	request.OperationID = newTestID()
	request.ExpectedRevision = done.Revision
	if _, e = s.ResumeSession(LocalController(), request); e != nil {
		t.Fatal(e)
	}
	path := filepath.Join(s.root, "SDP/Maintenance/Plan.md")
	raw := mustFile(t, path)
	os.WriteFile(path, append(raw, []byte("\nExternal change\n")...), 0600)
	view, e = s.ReadSession(q.InstanceID)
	if e != nil || view.Coverage == "current" {
		t.Fatal("selected edit accepted")
	}
}
func addUnrelatedCard(t *testing.T, s *Service) {
	t.Helper()
	path := "SDP/KanBan/backlog/#002--Change--unrelated.md"
	evidenceFile(t, s, path, "| Field | Value |\n| --- | --- |\n| id | KB-SDP-002 |\n| project | SDP |\n| type | Change |\n| created | 2026-10-24 |\n| source | fixture |\n| CardState | backlog |\n")
	index := mustFile(t, filepath.Join(s.root, "SDP/KanBan/README.md"))
	evidenceFile(t, s, "SDP/KanBan/README.md", string(index)+"| KB-SDP-002 | Change | backlog | Unrelated |\n")
	raw := mustFile(t, filepath.Join(s.root, managementLedger))
	events, _ := publicationEvents(raw)
	id, _ := nextPublicationEvent(events, "EVT-KB-SDP-")
	ev := map[string]any{"schemaVersion": "1.0", "eventId": id, "eventType": "x-kanban:created", "occurredAt": "2026-10-24T00:00:00Z", "actor": "fixture", "commit": nil, "subjectId": "KB-SDP-002", "payload": map[string]any{"schemaVersion": "0.2", "projectId": "SDP", "previousEventId": nil, "from": nil, "fromPath": nil, "to": "backlog", "toPath": "backlog/#002--Change--unrelated.md", "reason": "Unrelated new work", "links": []string{}}}
	b, _ := json.Marshal(ev)
	evidenceFile(t, s, managementLedger, appendRecord(raw, string(b)+"\n"))
}
func TestSessionSelectionEveryCrashBoundary(t *testing.T) {
	sample, q := selectionFixture(t, "", "")
	if _, e := sample.SelectSession(LocalController(), q); e != nil {
		t.Fatal(e)
	}
	st, _, _ := sample.load()
	count := len(st.Sessions[q.InstanceID].Selection.Changes)
	points := []string{"before-head", "after-head", "session-selection-after-claim"}
	for n := 0; n < count; n++ {
		points = append(points, fmt.Sprintf("session-selection-after-file:%d", n))
	}
	for _, point := range points {
		t.Run(point, func(t *testing.T) {
			s, q := selectionFixture(t, "", "")
			r, e := s.SelectSession(LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			request := selectionResume(q, r)
			once := true
			s.fault = func(p string) error {
				if once && p == point {
					once = false
					return errors.New("simulated interruption")
				}
				return nil
			}
			if _, e = s.ResumeSession(LocalController(), request); e == nil {
				t.Fatal("fault not reached")
			}
			s.fault = nil
			fresh, e := New(s.root, s.area, s.project, s.state)
			if e != nil {
				t.Fatal(e)
			}
			if _, e = fresh.ResumeSession(LocalController(), request); e != nil {
				t.Fatal(e)
			}
			if _, e = fresh.publicationFiles(); e != nil {
				t.Fatal(e)
			}
			ledger := mustFile(t, filepath.Join(s.root, managementLedger))
			events, _ := publicationEvents(ledger)
			if len(events) != 4 {
				t.Fatal("duplicate lifecycle events", len(events))
			}
			tx, e := fresh.canonicalTransaction()
			if e != nil {
				t.Fatal(e)
			}
			defer tx.Close()
			if e = tx.Check(""); e != nil {
				t.Fatal("claim not cleared", e)
			}
		})
	}
}
func TestSessionSelectionConflictPreservesEditsAndBlocksWriters(t *testing.T) {
	s, q := selectionFixture(t, "", "")
	r, e := s.SelectSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	request := selectionResume(q, r)
	s.fault = func(point string) error {
		if point == "session-selection-after-claim" {
			return errors.New("crash")
		}
		return nil
	}
	_, _ = s.ResumeSession(LocalController(), request)
	s.fault = nil
	ledger := filepath.Join(s.root, managementLedger)
	if _, e := projecthistory.Read(ledger); !errors.Is(e, projecthistory.ErrBusy) {
		t.Fatal("cooperating reader missed pending effects", e)
	}
	snapshotRevision := hash(mustFile(t, ledger))
	_, e = projecthistory.Append(ledger, snapshotRevision, []byte(`{"schemaVersion":"1.0","eventId":"EVT-X","eventType":"x-test:recorded","occurredAt":"2026-10-24T00:00:00Z","actor":"fixture","commit":null,"subjectId":"TEST","payload":{}}`))
	if !errors.Is(e, projecthistory.ErrBusy) {
		t.Fatal("canonical append bypassed pending claim", e)
	}
	source := filepath.Join(s.root, "SDP/KanBan/backlog/#001--Change--fixture.md")
	edited := append(mustFile(t, source), []byte("\nOwner edits\n")...)
	os.WriteFile(source, edited, 0600)
	head := mustFile(t, filepath.Join(s.state, "HEAD"))
	if _, e = s.ResumeSession(LocalController(), request); e == nil {
		t.Fatal("external edit accepted")
	}
	if string(mustFile(t, source)) != string(edited) || string(mustFile(t, filepath.Join(s.state, "HEAD"))) != string(head) {
		t.Fatal("external bytes or state changed")
	}
}

func TestSessionSelectionGuardedAdmission(t *testing.T) {
	for _, kind := range []string{"outside-link", "destination", "scope-claim", "stale-revision", "dependency", "unauthorized"} {
		t.Run(kind, func(t *testing.T) {
			s, q := selectionFixture(t, "", "")
			caller := LocalController()
			switch kind {
			case "outside-link":
				evidenceFile(t, s, "SDP/extra.md", "[Card](KanBan/backlog/%23001--Change--fixture.md)\n")
			case "destination":
				evidenceFile(t, s, "SDP/KanBan/active/#001--Change--fixture.md", "Do not overwrite this destination\n")
			case "scope-claim":
				raw, _ := json.Marshal(sessionScopeClaim{Schema: "sdp-session-scope-claim/0.1", OperationID: "3a49a865-45e2-4544-9099-0b4f76baebc1", SessionID: "SESSION-OTHER-0001", PlanID: "PLAN-OTHER-0001", Cards: []string{}, WritablePaths: []string{"target.txt"}, ProtectedPaths: []string{}})
				evidenceFile(t, s, "SDP/Sessions/execution/claims/3a49a865-45e2-4544-9099-0b4f76baebc1.json", string(raw))
			case "stale-revision":
				q.ExpectedRevision++
			case "dependency":
				q.StepID = "missing"
			case "unauthorized":
				caller = Caller{}
			}
			head := mustFile(t, filepath.Join(s.state, "HEAD"))
			if _, e := s.SelectSession(caller, q); e == nil {
				t.Fatal("invalid selection accepted")
			}
			if string(mustFile(t, filepath.Join(s.state, "HEAD"))) != string(head) {
				t.Fatal("refusal changed state")
			}
		})
	}
}
func TestSessionSelectionPreservesActiveContribution(t *testing.T) {
	s, c, adopt := sessionFixture(t)
	c.Binding.Steps[0].Cards[0].Coverage = "contribution"
	raw, _ := json.Marshal(c.Binding)
	evidenceFile(t, s, c.BindingPath, string(raw))
	ctx, e := s.SessionContext(SessionContextRequest{c.BindingPath})
	if e != nil {
		t.Fatal(e)
	}
	adopt.InputRevision = ctx.InputRevision
	a, e := s.AdoptSession(LocalController(), adopt)
	if e != nil {
		t.Fatal(e)
	}
	r, e := s.ResumeSession(LocalController(), sessionResume(a, adopt))
	if e != nil {
		t.Fatal(e)
	}
	q := SessionSelectRequest{SessionSchema, adopt.ProjectRef, adopt.CheckoutID, newTestID(), a.InstanceID, r.Revision, "S1", "Select a contribution without closing the broader card"}
	p, e := s.SelectSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.ResumeSession(LocalController(), selectionResume(q, p)); e != nil {
		t.Fatal(e)
	}
	raw = mustFile(t, filepath.Join(s.root, "SDP/KanBan/active/#001--Change--fixture.md"))
	m, _ := metadata(raw)
	if m["CardState"] != "in-progress" || !strings.Contains(string(raw), "Card contribution: contribution") {
		t.Fatal("active contribution changed lifecycle")
	}
}

func TestSessionSelectionCrossStoreAndCorruptCoverage(t *testing.T) {
	s, q := selectionFixture(t, "", "")
	prepared, e := s.SelectSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.ResumeSession(LocalController(), selectionResume(q, prepared)); e != nil {
		t.Fatal(e)
	}
	before := mustFile(t, filepath.Join(s.root, "SDP/Sessions/session.md"))
	other, e := New(s.root, s.area, s.project, filepath.Join(t.TempDir(), "other-state"))
	if e != nil {
		t.Fatal(e)
	}
	bound, e := other.Bind(LocalController())
	if e != nil {
		t.Fatal(e)
	}
	ctx, e := other.SessionContext(SessionContextRequest{"SDP/Sessions/execution/SESSION-TEST-0001.json"})
	if e != nil {
		t.Fatal(e)
	}
	adopt := SessionAdoptRequest{SessionSchema, s.project, bound.CheckoutID, newTestID(), ctx.BindingPath, ctx.InputRevision, "Competing store cannot register selected Session"}
	a, e := other.AdoptSession(LocalController(), adopt)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = other.ResumeSession(LocalController(), sessionResume(a, adopt)); e == nil {
		t.Fatal("second store registered a claimed Session")
	}
	if string(mustFile(t, filepath.Join(s.root, "SDP/Sessions/session.md"))) != string(before) {
		t.Fatal("competing registration changed Session")
	}
	st, head, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	i := st.Sessions[q.InstanceID]
	i.Selection.Coverage = i.Selection.Coverage[:len(i.Selection.Coverage)-1]
	st.Sessions[q.InstanceID] = i
	if e = s.save(st, head); e != nil {
		t.Fatal(e)
	}
	if _, _, e = s.load(); e == nil {
		t.Fatal("missing source pin accepted")
	}
}
