package governance

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func sessionFixture(t *testing.T) (*Service, SessionContext, SessionAdoptRequest) {
	t.Helper()
	return sessionFixtureAt(t, "", "")
}
func sessionFixtureAt(t *testing.T, root, state string) (*Service, SessionContext, SessionAdoptRequest) {
	t.Helper()
	s, r, _ := publicationFixtureAt(t, root, state)
	sel := r.Run.Selection
	// Existing tests exercise a complete canonical record set, including an earlier run.
	plan, e := os.ReadFile(filepath.Join(s.root, sel.Plan))
	if e != nil {
		t.Fatal(e)
	}
	evidenceFile(t, s, sel.Plan, string(plan)+"\n| Milestone | Outcome |\n| --- | --- |\n| M1 | bounded work |\n")
	path := "SDP/Sessions/execution/SESSION-TEST-0001.json"
	body, e := os.ReadFile(filepath.Join(s.root, sel.Session))
	if e != nil {
		t.Fatal(e)
	}
	evidenceFile(t, s, sel.Session, string(body)+"\n[Execution binding](execution/SESSION-TEST-0001.json)\n")
	w := SessionWorkflow()
	b := SessionBinding{Schema: SessionSchema, Session: SessionRecordRef{"SESSION-TEST-0001", sel.Session}, Plan: SessionRecordRef{"MAINT-SDP-0001", sel.Plan}, Authority: sel.Authority, Workflow: SessionWorkflowRef{w.ID, w.Version, w.Digest()}, Steps: []SessionStepBinding{{ID: "S1", Milestone: "M1", Cards: []SessionCardRef{{"KB-SDP-001", "whole"}}, DependsOn: []string{}, WritablePaths: []string{"target.txt"}, ProtectedPaths: []string{"unrelated.txt"}, Review: "independent", OwnerAcceptance: "not-required"}}}
	raw, _ := json.Marshal(b)
	evidenceFile(t, s, path, string(raw))
	c, e := s.SessionContext(SessionContextRequest{path})
	if e != nil {
		t.Fatal(e)
	}
	id, _ := identifier()
	q := SessionAdoptRequest{SessionSchema, s.project, c.CheckoutID, id, path, c.InputRevision, "Owner-selected test Session"}
	return s, c, q
}
func sessionResume(i SessionMutationResult, q SessionAdoptRequest) SessionResumeRequest {
	id, _ := identifier()
	return SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, id, i.InstanceID, i.Revision, "Recover recorded Session"}
}
func TestSessionAdoptionResumeAndIsolation(t *testing.T) {
	s, c, q := sessionFixture(t)
	before, _ := os.ReadFile(filepath.Join(s.root, c.Binding.Session.Path))
	legacy, _, _ := s.load()
	a, e := s.AdoptSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	if a.Registration != "pending" || a.Revision != 1 {
		t.Fatal(a)
	}
	unchanged, _ := os.ReadFile(filepath.Join(s.root, c.Binding.Session.Path))
	if string(before) != string(unchanged) {
		t.Fatal("adoption prematurely changed project")
	}
	repeat, e := s.AdoptSession(LocalController(), q)
	if e != nil || !reflect.DeepEqual(a, repeat) {
		t.Fatal(repeat, e)
	}
	duplicate := q
	duplicate.OperationID, _ = identifier()
	_, e = s.AdoptSession(LocalController(), duplicate)
	code(t, e, "conflict")
	changed := q
	changed.Reason = "different"
	_, e = s.AdoptSession(LocalController(), changed)
	code(t, e, "idempotency_conflict")
	v, e := s.ReadSession(a.InstanceID)
	if e != nil || v.Coverage != "current" || v.Registration != "pending" || v.Steps[0].Stage != "admission" {
		t.Fatal(v, e)
	}
	resume := sessionResume(a, q)
	done, e := s.ResumeSession(LocalController(), resume)
	if e != nil || done.Revision != 2 {
		t.Fatal(done, e)
	}
	fresh, e := New(s.root, s.area, s.project, s.state)
	if e != nil {
		t.Fatal(e)
	}
	again, e := fresh.ResumeSession(LocalController(), resume)
	if e != nil || again != done {
		t.Fatal(again, e)
	}
	view, e := fresh.ReadSession(a.InstanceID)
	if e != nil || view.Registration != "committed" || view.HostExecution != "not-dispatched" || view.Coverage != "current" {
		t.Fatal(view, e)
	}
	after, _ := os.ReadFile(filepath.Join(s.root, c.Binding.Session.Path))
	if !strings.HasPrefix(string(after), string(before)) || strings.Count(string(after), "## PG Session registration ") != 1 {
		t.Fatal("journal history changed or duplicated")
	}
	st, _, e := fresh.load()
	if e != nil || st.Schema != sessionStateSchema || !reflect.DeepEqual(st.Runs, legacy.Runs) || !reflect.DeepEqual(st.Reviews, legacy.Reviews) {
		t.Fatal("legacy state changed", e)
	}
	bad := resume
	bad.OperationID, _ = identifier()
	_, e = fresh.ResumeSession(LocalController(), bad)
	code(t, e, "stale_revision")
	b, _ := os.ReadFile(filepath.Join(s.root, "unrelated.txt"))
	if string(b) != "keep this dirty input\n" {
		t.Fatal("unrelated file changed")
	}
}
func TestSessionCrashRecovery(t *testing.T) {
	for _, point := range []string{"before-head", "after-head", "session-after-journal"} {
		t.Run(point, func(t *testing.T) {
			s, c, q := sessionFixture(t)
			if point != "session-after-journal" {
				once := true
				s.fault = func(p string) error {
					if once && p == point {
						once = false
						return errors.New("simulated crash")
					}
					return nil
				}
				_, e := s.AdoptSession(LocalController(), q)
				if e == nil {
					t.Fatal("fault not reached")
				}
				s.fault = nil
			}
			a, e := s.AdoptSession(LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			r := sessionResume(a, q)
			once := true
			s.fault = func(p string) error {
				if once && p == point {
					once = false
					return errors.New("simulated crash")
				}
				return nil
			}
			_, e = s.ResumeSession(LocalController(), r)
			if e == nil {
				t.Fatal("resume fault not reached")
			}
			s.fault = nil
			done, e := s.ResumeSession(LocalController(), r)
			if e != nil || done.Revision != 2 {
				t.Fatal(done, e)
			}
			body, _ := os.ReadFile(filepath.Join(s.root, c.Binding.Session.Path))
			if strings.Count(string(body), "## PG Session registration ") != 1 {
				t.Fatal("duplicate registration")
			}
		})
	}
}
func TestSessionDriftAndExternalEdit(t *testing.T) {
	for _, what := range []string{"binding", "plan", "session", "authority", "ledger"} {
		t.Run(what, func(t *testing.T) {
			s, c, q := sessionFixture(t)
			a, e := s.AdoptSession(LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			paths := map[string]string{"binding": q.BindingPath, "plan": c.Binding.Plan.Path, "session": c.Binding.Session.Path, "authority": c.Binding.Authority, "ledger": "SDP/ProjectManagement/Ledger.ndjson"}
			p := filepath.Join(s.root, paths[what])
			old, _ := os.ReadFile(p)
			extra := "\nExternal edit\n"
			if what == "binding" || what == "ledger" {
				extra = "\n"
			}
			os.WriteFile(p, append(old, []byte(extra)...), 0600)
			head, _ := os.ReadFile(filepath.Join(s.state, "HEAD"))
			_, e = s.ResumeSession(LocalController(), sessionResume(a, q))
			if e == nil {
				t.Fatal("drift accepted")
			}
			now, _ := os.ReadFile(filepath.Join(s.state, "HEAD"))
			if string(head) != string(now) {
				t.Fatal("failed resume changed state")
			}
			v, e := s.ReadSession(a.InstanceID)
			if e != nil || v.Coverage == "current" {
				t.Fatal(v, e)
			}
			actual, _ := os.ReadFile(p)
			if string(actual) != string(old)+extra {
				t.Fatal("external edit overwritten")
			}
		})
	}
	s, c, q := sessionFixture(t)
	a, _ := s.AdoptSession(LocalController(), q)
	r := sessionResume(a, q)
	s.fault = func(p string) error {
		if p == "session-after-journal" {
			return errors.New("crash")
		}
		return nil
	}
	_, _ = s.ResumeSession(LocalController(), r)
	s.fault = nil
	p := filepath.Join(s.root, c.Binding.Session.Path)
	b, _ := os.ReadFile(p)
	os.WriteFile(p, append(b, []byte("\nnew owner change\n")...), 0600)
	_, e := s.ResumeSession(LocalController(), r)
	code(t, e, "stale_revision")
}
func TestSessionBindingRejections(t *testing.T) {
	cases := map[string]func(*SessionBinding){
		"cycle":              func(b *SessionBinding) { b.Steps[0].DependsOn = []string{"S1"} },
		"missing dependency": func(b *SessionBinding) { b.Steps[0].DependsOn = []string{"S2"} },
		"duplicate step":     func(b *SessionBinding) { b.Steps = append(b.Steps, b.Steps[0]) },
		"foreign version":    func(b *SessionBinding) { b.Workflow.Version = "999" },
		"wrong digest":       func(b *SessionBinding) { b.Workflow.Digest = strings.Repeat("a", 64) },
		"path escape":        func(b *SessionBinding) { b.Authority = "SDP/../authority.md" },
		"scope glob":         func(b *SessionBinding) { b.Steps[0].WritablePaths = []string{"*.go"} },
		"missing review":     func(b *SessionBinding) { b.Steps[0].Review = "" },
		"missing card":       func(b *SessionBinding) { b.Steps[0].Cards[0].ID = "KB-SDP-099" },
		"missing milestone":  func(b *SessionBinding) { b.Steps[0].Milestone = "Absent" },
		"wrong Session":      func(b *SessionBinding) { b.Session.ID = "SESSION-TEST-0002" },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			s, c, q := sessionFixture(t)
			mutate(&c.Binding)
			b, _ := json.Marshal(c.Binding)
			evidenceFile(t, s, q.BindingPath, string(b))
			_, e := s.SessionContext(SessionContextRequest{q.BindingPath})
			if e == nil {
				t.Fatal("invalid binding accepted")
			}
		})
	}
	for _, input := range []string{`{"bindingPath":"a","bindingPath":"b"}`, `{"bindingPath":"a","role":"controller"}`, `{"bindingPath":null}`} {
		var q SessionContextRequest
		if e := Decode([]byte(input), &q); e == nil {
			t.Fatal("loose JSON accepted", input)
		}
	}
}
func TestSessionReadAndUnauthorizedHaveNoEffects(t *testing.T) {
	s, _, q := sessionFixture(t)
	fresh, e := New(s.root, s.area, s.project, filepath.Join(t.TempDir(), "absent"))
	if e != nil {
		t.Fatal(e)
	}
	if _, e = fresh.SessionContext(SessionContextRequest{q.BindingPath}); e != nil {
		t.Fatal(e)
	}
	list, e := fresh.ListSessions()
	if e != nil || len(list.Instances) != 0 {
		t.Fatal(list, e)
	}
	_, e = fresh.AdoptSession(Caller{}, q)
	code(t, e, "unauthorized")
	_, e = fresh.ResumeSession(Caller{}, SessionResumeRequest{})
	code(t, e, "unauthorized")
	if _, e = os.Stat(fresh.state); !os.IsNotExist(e) {
		t.Fatal("read/untrusted call created state")
	}
}
func TestSessionCorruptStateRejected(t *testing.T) {
	for _, change := range []string{"schema", "identity", "registration", "binding", "receipt", "revision-gap"} {
		t.Run(change, func(t *testing.T) {
			s, _, q := sessionFixture(t)
			a, e := s.AdoptSession(LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			st, head, _ := s.load()
			i := st.Sessions[a.InstanceID]
			switch change {
			case "schema":
				st.Schema = publicationStateSchema
			case "identity":
				i.ID = "wrong"
			case "registration":
				i.Registration.Content += "forged"
			case "binding":
				i.Context.Binding.Workflow.Version = "new"
			case "revision-gap":
				i.Revision = 3
				i.Registered = true
			case "receipt":
				delete(st.Receipts, q.OperationID)
			}
			st.Sessions[a.InstanceID] = i
			if e = s.save(st, head); e != nil {
				t.Fatal(e)
			}
			if _, _, e = s.load(); e == nil {
				t.Fatal("corrupt state accepted")
			}
		})
	}
}

func TestSessionEmptyStoreAndLaterPublication(t *testing.T) {
	s, c, q := sessionFixture(t)
	fresh, e := New(s.root, s.area, s.project, filepath.Join(t.TempDir(), "empty-state"))
	if e != nil {
		t.Fatal(e)
	}
	bound, e := fresh.Bind(LocalController())
	if e != nil {
		t.Fatal(e)
	}
	q.CheckoutID = bound.CheckoutID
	a, e := fresh.AdoptSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = fresh.ResumeSession(LocalController(), sessionResume(a, q)); e != nil {
		t.Fatal(e)
	}
	if _, e = fresh.ReadSession(a.InstanceID); e != nil {
		t.Fatal(e)
	}
	// A separate fixture proves older publication writes preserve the new profile.
	s, _, pub := publicationFixture(t)
	plan, _ := os.ReadFile(filepath.Join(s.root, "SDP/Maintenance/Plan.md"))
	evidenceFile(t, s, "SDP/Maintenance/Plan.md", string(plan)+"\n| Milestone | Outcome |\n| --- | --- |\n| M1 | Work |\n")
	body, _ := os.ReadFile(filepath.Join(s.root, c.Binding.Session.Path))
	evidenceFile(t, s, c.Binding.Session.Path, string(body)+"\n[Execution](execution/SESSION-TEST-0001.json)\n")
	raw, _ := json.Marshal(c.Binding)
	evidenceFile(t, s, q.BindingPath, string(raw))
	ctx, e := s.SessionContext(SessionContextRequest{q.BindingPath})
	if e != nil {
		t.Fatal(e)
	}
	q.CheckoutID = ctx.CheckoutID
	q.InputRevision = ctx.InputRevision
	if _, e = s.AdoptSession(LocalController(), q); e != nil {
		t.Fatal(e)
	}
	// Preparing publication can reject its stale source basis, but must never
	// downgrade or damage Session state. No result is accepted on stale inputs.
	if _, e = s.PreparePublication(LocalController(), pub); e == nil {
		t.Fatal("stale publication unexpectedly accepted")
	}
	st, _, e := s.load()
	if e != nil || st.Schema != sessionStateSchema || len(st.Sessions) != 1 {
		t.Fatal("publication changed Session profile", e)
	}
}
