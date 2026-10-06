package governance

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func fixture(t *testing.T) (*Service, Selection, StartRequest) {
	t.Helper()
	base := t.TempDir()
	root := filepath.Join(base, "project")
	area := filepath.Join(root, "SDP")
	state := filepath.Join(base, "state")
	files := map[string]string{
		"SDP/ProjectManagement/Ledger.ndjson": `{"schemaVersion":"1.0","eventId":"E1","eventType":"x-kanban:created","subjectId":"KB-TEST-001","payload":{"schemaVersion":"0.2","previousEventId":null,"to":"active","toPath":"active/card.md"}}
{"schemaVersion":"1.0","eventId":"E2","eventType":"x-management:started","subjectId":"MAINT-TEST-0001","payload":{"schemaVersion":"0.2","previousEventId":null,"to":"active","toPath":"Maintenance/Plan.md"}}
`,
		"AGENTS.md":                     "Preserve unrelated files.\n",
		"SDP/SDP-project.manifest.yaml": "schemaVersion: '1.0'\nproject:\n  name: trial\n",
		"SDP/KanBan/active/card.md":     "# Work\n\n| Field | Value |\n| --- | --- |\n| id | KB-TEST-001 |\n| PlanId | MAINT-TEST-0001 |\n| CardState | in-progress |\n",
		"SDP/Maintenance/Plan.md":       "# Plan\n\n| Field | Value |\n| --- | --- |\n| id | MAINT-TEST-0001 |\n| state | active |\n| PlanType | MaintenancePlan |\n",
		"SDP/Sessions/session.md":       "# Goal\n\n| State | Step |\n| --- | --- |\n| next | S1 |\n\n| Field | Value |\n| --- | --- |\n| Session reference | SESSION-TEST-0001 |\n| Status | active |\n",
		"SDP/authority.md":              "Owner selects the bounded fixture change.\n",
		"unrelated.txt":                 "keep this dirty input\n",
	}
	for p, b := range files {
		path := filepath.Join(root, p)
		if e := os.MkdirAll(filepath.Dir(path), 0700); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(path, []byte(b), 0600); e != nil {
			t.Fatal(e)
		}
	}
	s, e := New(root, area, "trial", state)
	if e != nil {
		t.Fatal(e)
	}
	sel := Selection{"SDP/KanBan/active/card.md", "SDP/Maintenance/Plan.md", "SDP/Sessions/session.md", "SDP/authority.md", "S1"}
	c, e := s.Context(sel)
	if e != nil {
		t.Fatal(e)
	}
	if len(c.Conditions) != 0 {
		t.Fatal(c.Conditions)
	}
	op, _ := identifier()
	req, _ := identifier()
	q := StartRequest{Schema: Schema, ProjectRef: "trial", OperationID: op, RequestID: req, InputRevision: c.Revision, RoutineDigest: Maintenance().Digest(), Selection: sel, Assignment: Assignment{Outcome: "fix fixture", WritablePaths: []string{"target.txt"}, ExcludedPaths: []string{"unrelated.txt"}}}
	return s, sel, q
}
func bind(t *testing.T, s *Service, q *StartRequest) {
	t.Helper()
	b, e := s.Bind(LocalController())
	if e != nil {
		t.Fatal(e)
	}
	q.CheckoutID = b.CheckoutID
}
func code(t *testing.T, e error, want string) {
	t.Helper()
	var domain *Error
	if !errors.As(e, &domain) || domain.Code != want {
		t.Fatalf("want %s got %v", want, e)
	}
}
func transition(r Result, action string) TransitionRequest {
	op, _ := identifier()
	return TransitionRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: op, RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: action, Reason: "test observed boundary"}
}

func TestReadRoutingAndNoMutation(t *testing.T) {
	s, sel, q := fixture(t)
	r, e := s.Route(RouteRequest{Schema: Schema, Categories: []string{"status"}})
	if e != nil || r.Disposition != "no-work" {
		t.Fatal(r, e)
	}
	route := RouteRequest{Schema: Schema, Categories: []string{"maintenance"}, Selection: sel, RoutineID: "maintenance", RoutineVersion: "0.1"}
	r, e = s.Route(route)
	if e != nil || r.Disposition != "matched" {
		t.Fatal(r, e)
	}
	if _, e = os.Stat(s.state); !os.IsNotExist(e) {
		t.Fatal("read created state")
	}
	_, e = s.Start(Caller{}, q)
	code(t, e, "unauthorized")
	if _, e = os.Stat(s.state); !os.IsNotExist(e) {
		t.Fatal("untrusted call created state")
	}
	for _, tc := range []struct {
		category          []string
		id, version, want string
	}{
		{[]string{"release"}, "maintenance", "0.1", "missing"},
		{[]string{"maintenance", "status"}, "maintenance", "0.1", "ambiguous"},
		{[]string{"maintenance"}, "other", "0.1", "missing"},
		{[]string{"maintenance"}, "maintenance", "99", "incompatible"},
	} {
		route.Categories = tc.category
		route.RoutineID = tc.id
		route.RoutineVersion = tc.version
		r, e = s.Route(route)
		if e != nil || r.Disposition != tc.want {
			t.Fatal(tc, r, e)
		}
	}
	route.Categories = []string{"maintenance"}
	route.RoutineID = "maintenance"
	route.RoutineVersion = "0.1"
	os.Remove(filepath.Join(s.root, sel.Plan))
	r, e = s.Route(route)
	if e != nil || r.Disposition != "unavailable" {
		t.Fatal(r, e)
	}
}
func TestDurableRunRetryGuardsAndResume(t *testing.T) {
	s, _, q := fixture(t)
	bind(t, s, &q)
	r, e := s.Start(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	other, e := New(s.root, s.area, s.project, s.state)
	if e != nil {
		t.Fatal(e)
	}
	repeated, e := other.Start(LocalController(), q)
	if e != nil || repeated.Run.ID != r.Run.ID || repeated.Run.Revision != 1 {
		t.Fatal(repeated, e)
	}
	changed := q
	changed.Assignment.Outcome = "new outcome"
	_, e = s.Start(LocalController(), changed)
	code(t, e, "idempotency_conflict")
	changed = q
	changed.OperationID, _ = identifier()
	_, e = s.Start(LocalController(), changed)
	code(t, e, "conflict")
	start := transition(r, "start")
	r, e = other.Transition(LocalController(), start)
	if e != nil || r.Run.State != "running" {
		t.Fatal(r, e)
	}
	retry, e := s.Transition(LocalController(), start)
	if e != nil || retry.Run.Revision != 2 {
		t.Fatal(retry, e)
	}
	stale := start
	stale.OperationID, _ = identifier()
	_, e = s.Transition(LocalController(), stale)
	code(t, e, "stale_revision")
	for _, action := range []string{"interrupt", "resume", "start", "report"} {
		r, e = s.Transition(LocalController(), transition(r, action))
		if e != nil {
			t.Fatal(action, e)
		}
	}
	_, e = s.Transition(LocalController(), transition(r, "complete"))
	code(t, e, "evidence_missing")
	read, e := other.Read(r.Run.ID)
	if e != nil || read.Run.State != "verifying" {
		t.Fatal(read, e)
	}
	if read.Run.RoutineDigest != Maintenance().Digest() {
		t.Fatal("routine changed")
	}
	unchanged, _ := os.ReadFile(filepath.Join(s.root, "unrelated.txt"))
	if string(unchanged) != "keep this dirty input\n" {
		t.Fatal("modified unrelated input")
	}
}
func TestCrashBeforeAndAfterCommitPoint(t *testing.T) {
	for _, point := range []string{"before-head", "after-head"} {
		t.Run(point, func(t *testing.T) {
			s, _, q := fixture(t)
			bind(t, s, &q)
			r, e := s.Start(LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			tr := transition(r, "start")
			s.fault = func(p string) error {
				if p == point {
					return fmt.Errorf("simulated crash")
				}
				return nil
			}
			if _, e = s.Transition(LocalController(), tr); e == nil {
				t.Fatal("fault not injected")
			}
			s.fault = nil
			current, e := s.Read(r.Run.ID)
			if e != nil {
				t.Fatal(e)
			}
			want := uint64(1)
			if point == "after-head" {
				want = 2
			}
			if current.Run.Revision != want {
				t.Fatal(current)
			}
			retried, e := s.Transition(LocalController(), tr)
			if e != nil || retried.Run.Revision != 2 {
				t.Fatal(retried, e)
			}
			retried, e = s.Transition(LocalController(), tr)
			if e != nil || retried.Run.Revision != 2 {
				t.Fatal(retried, e)
			}
		})
	}
}
func TestWorkInputDriftAndStop(t *testing.T) {
	s, sel, q := fixture(t)
	bind(t, s, &q)
	r, e := s.Start(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.Transition(LocalController(), transition(r, "start"))
	if e != nil {
		t.Fatal(e)
	}
	os.WriteFile(filepath.Join(s.root, "unrelated.txt"), []byte("another unrelated edit"), 0600)
	view, e := s.Read(r.Run.ID)
	if e != nil || view.Coverage != "known" {
		t.Fatal(view, e)
	}
	f, _ := os.OpenFile(filepath.Join(s.root, sel.Authority), os.O_APPEND|os.O_WRONLY, 0600)
	f.WriteString("New scope constraint\n")
	f.Close()
	view, e = s.Read(r.Run.ID)
	if e != nil || view.Coverage != "stale" {
		t.Fatal(view, e)
	}
	_, e = s.Transition(LocalController(), transition(r, "report"))
	code(t, e, "stale_revision")
	r, e = s.Transition(LocalController(), transition(r, "interrupt"))
	if e != nil || r.Coverage != "stale" {
		t.Fatal(r, e)
	}
	_, e = s.Transition(LocalController(), transition(r, "resume"))
	code(t, e, "stale_revision")
	_, e = s.Transition(LocalController(), transition(r, "cancel"))
	if e != nil {
		t.Fatal(e)
	}
}
func TestPathsBindingsAndCorruptState(t *testing.T) {
	s, sel, q := fixture(t)
	bad := sel
	bad.Authority = "../secret.md"
	_, e := s.Context(bad)
	code(t, e, "path")
	os.Remove(filepath.Join(s.root, sel.Authority))
	os.Symlink(filepath.Join(s.root, "unrelated.txt"), filepath.Join(s.root, sel.Authority))
	_, e = s.Context(sel)
	code(t, e, "path")
	os.Remove(filepath.Join(s.root, sel.Authority))
	os.WriteFile(filepath.Join(s.root, sel.Authority), []byte("owner"), 0600)
	bind(t, s, &q)
	other := *s
	other.project = "other"
	_, e = other.Bind(LocalController())
	code(t, e, "unauthorized")
	head, _ := os.ReadFile(filepath.Join(s.state, "HEAD"))
	os.WriteFile(filepath.Join(s.state, "generations", string(head)+".json"), []byte("{}"), 0600)
	_, e = s.Context(sel)
	code(t, e, "reconciliation_required")
}
func TestConcurrentRevisionAndLock(t *testing.T) {
	s, _, q := fixture(t)
	bind(t, s, &q)
	r, e := s.Start(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	release, e := lock(s.state)
	if e != nil {
		t.Fatal(e)
	}
	_, e = s.Transition(LocalController(), transition(r, "start"))
	code(t, e, "busy")
	release()
	var wg sync.WaitGroup
	results := make(chan error, 2)
	for i := 0; i < 2; i++ {
		tr := transition(r, "start")
		wg.Add(1)
		go func() { defer wg.Done(); _, e := s.Transition(LocalController(), tr); results <- e }()
	}
	wg.Wait()
	close(results)
	successes := 0
	for e := range results {
		if e == nil {
			successes++
		} else {
			var d *Error
			if !errors.As(e, &d) || (d.Code != "busy" && d.Code != "stale_revision") {
				t.Fatal(e)
			}
		}
	}
	if successes != 1 {
		t.Fatal(successes)
	}
}
func TestStrictInputs(t *testing.T) {
	var aliases RouteRequest
	code(t, Decode([]byte(`{"schema":"sdp-project-governance/0.1","Schema":"sdp-project-governance/0.1","categories":["status"]}`), &aliases), "invalid_request")

	for _, b := range []string{`{"schema":"x","schema":"y"}`, `{"role":"owner"}`, `{"expectedRevision":1.5}`, `{"schema":null}`, `{} {}`, `{"selection":{"card":"a","card":"b"}}`} {
		var q TransitionRequest
		if Decode([]byte(b), &q) == nil {
			t.Fatal("accepted", b)
		}
	}
	var q RouteRequest
	if e := Decode([]byte(`{"schema":"sdp-project-governance/0.1","categories":["status"]}`), &q); e != nil {
		t.Fatal(e)
	}
	if e := Decode(make([]byte, MaxRequest+1), &q); e == nil {
		t.Fatal("unbounded input")
	}
	// Field order cannot change the semantic idempotency digest.
	a := map[string]any{"x": 1, "y": 2}
	b := map[string]any{"y": 2, "x": 1}
	if requestDigest(a) != requestDigest(b) {
		t.Fatal("unstable normalized digest")
	}
	raw, _ := json.Marshal(Maintenance())
	if !json.Valid(raw) {
		t.Fatal("invalid pinned routine")
	}
}

func TestCanonicalCoverageConflict(t *testing.T) {
	s, sel, q := fixture(t)
	bind(t, s, &q)
	p := filepath.Join(s.area, "ProjectManagement/Ledger.ndjson")
	b, _ := os.ReadFile(p)
	b = append(b, []byte(`{"schemaVersion":"1.0","eventId":"E3","eventType":"x-management:completed","subjectId":"MAINT-TEST-0001","payload":{"schemaVersion":"0.2","previousEventId":"E2","to":"completed","toPath":"Maintenance/Plan.md"}}`)...)
	os.WriteFile(p, b, 0600)
	c, e := s.Context(sel)
	if e != nil || c.Coverage != "conflict" {
		t.Fatal(c, e)
	}
	q.InputRevision = c.Revision
	_, e = s.Start(LocalController(), q)
	code(t, e, "missing_coverage")
}
