package activities

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
	"os"
	"path/filepath"
	"testing"
)

func fixture(t *testing.T) (*Service, Caller) {
	t.Helper()
	area := t.TempDir()
	if e := os.MkdirAll(filepath.Join(area, "ProjectManagement"), 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(filepath.Join(area, "authority.md"), []byte("Owner authority"), 0600); e != nil {
		t.Fatal(e)
	}
	return New(area, "TEST"), Caller{Actor: "test/host", Authority: Authority{Role: "project-owner", Source: "authority.md"}}
}
func request(t *testing.T, s *Service, id string) CreateRequest {
	t.Helper()
	c, e := s.Read()
	if e != nil {
		t.Fatal(e)
	}
	return CreateRequest{OperationID: id, ExpectedHistoryRevision: c.HistoryRevision, Title: "Disposable activity", Owner: "project-owner", Mandate: "Need", Reason: "Owner selects need", Scope: []string{"Source"}}
}

const op1 = "11111111-1111-4111-8111-111111111111"
const op2 = "22222222-2222-4222-8222-222222222222"

func TestCreateRetryStaleAndChangedContent(t *testing.T) {
	s, caller := fixture(t)
	q := request(t, s, op1)
	r, e := s.Create(q, caller)
	if e != nil {
		t.Fatal(e)
	}
	c, e := s.Read()
	if e != nil || len(c.Activities) != 1 || c.Activities[0].State != "planned" {
		t.Fatalf("%+v %v", c, e)
	}
	other := request(t, s, op2)
	if _, e = s.Create(other, caller); e != nil {
		t.Fatal(e)
	}
	again, e := s.Create(q, caller)
	if e != nil || again != r {
		t.Fatalf("retry changed result: %+v %+v %v", r, again, e)
	}
	q.Title = "Changed"
	if _, e = s.Create(q, caller); !errors.Is(e, projecthistory.ErrConflict) {
		t.Fatal(e)
	}
	q.OperationID = "33333333-3333-4333-8333-333333333333"
	if _, e = s.Create(q, caller); !errors.Is(e, projecthistory.ErrStale) {
		t.Fatal(e)
	}
}
func TestCreateCrashBoundaries(t *testing.T) {
	for _, point := range []string{"after-intent", "after-reserve", "after-file-0", "after-file-1", "after-file-2", "after-file-3", "after-file-4", "after-file-5", "after-file-6", "after-file-7", "after-file-8", "after-file-9", "before-history", "after-history", "after-release"} {
		t.Run(point, func(t *testing.T) {
			s, caller := fixture(t)
			q := request(t, s, op1)
			s.fault = func(p string) error {
				if p == point {
					return errors.New("crash")
				}
				return nil
			}
			if _, e := s.Create(q, caller); e == nil {
				t.Fatal("fault not reached")
			}
			s.fault = nil
			r, e := s.Recover(op1, caller)
			if e != nil {
				t.Fatal(e)
			}
			c, e := s.Read()
			if e != nil || len(c.Activities) != 1 || c.Activities[0].Descriptor.ID != r.ActivityID {
				t.Fatalf("%+v %v", c, e)
			}
			if _, e = s.Create(q, caller); e != nil {
				t.Fatal(e)
			}
		})
	}
}
func sessionRequest(t *testing.T, s *Service) CreateSessionRequest {
	t.Helper()
	c, e := s.Read()
	if e != nil {
		t.Fatal(e)
	}
	a := c.Activities[0]
	return CreateSessionRequest{OperationID: op2, ExpectedHistoryRevision: c.HistoryRevision, ActivityID: a.Descriptor.ID, ExpectedActivityEventID: a.LastEventID, ExpectedDescriptorRevision: a.DescriptorRef.SHA256, Title: "Discuss", Objective: "Study first", Phase: "04-Design", Reason: "Owner discussion"}
}
func TestSessionAllocationMutableDocumentAndCrashRecovery(t *testing.T) {
	for _, point := range []string{"after-intent", "after-reserve", "after-file-1", "after-file-2", "before-history", "after-history", "after-release"} {
		t.Run(point, func(t *testing.T) {
			s, caller := fixture(t)
			if _, e := s.Create(request(t, s, op1), caller); e != nil {
				t.Fatal(e)
			}
			if e := os.MkdirAll(filepath.Join(s.Area, "Sessions"), 0700); e != nil {
				t.Fatal(e)
			}
			if e := os.WriteFile(filepath.Join(s.Area, "Sessions/session-#0019--Legacy.md"), []byte("old"), 0600); e != nil {
				t.Fatal(e)
			}
			q := sessionRequest(t, s)
			s.fault = func(p string) error {
				if p == point {
					return errors.New("crash")
				}
				return nil
			}
			if _, e := s.CreateSession(q, caller); e == nil {
				t.Fatal("fault not reached")
			}
			s.fault = nil
			r, e := s.Recover(op2, caller)
			if e != nil || r.SessionID != "SESSION-TEST-0020" {
				t.Fatalf("%+v %v", r, e)
			}
			if e := os.WriteFile(filepath.Join(s.Area, filepath.FromSlash(r.Path)), []byte("working log edited"), 0600); e != nil {
				t.Fatal(e)
			}
			if _, e = s.Read(); e != nil {
				t.Fatal(e)
			}
			if _, e = s.CreateSession(q, caller); e != nil {
				t.Fatal(e)
			}
		})
	}
}
func TestRecoveryRefusesExternalConflictAndBlocksOtherWriters(t *testing.T) {
	s, caller := fixture(t)
	q := request(t, s, op1)
	s.fault = func(p string) error {
		if p == "after-file-0" {
			return errors.New("crash")
		}
		return nil
	}
	if _, e := s.Create(q, caller); e == nil {
		t.Fatal("fault")
	}
	s.fault = nil
	in, e := s.load(op1)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(s.Area, filepath.FromSlash(in.Files[0].Path)), []byte("external conflict"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Recover(op1, caller); e == nil {
		t.Fatal("overwrote external change")
	}
	if _, e = projecthistory.Read(filepath.Join(s.Area, "ProjectManagement/Ledger.ndjson")); !errors.Is(e, projecthistory.ErrBusy) {
		t.Fatal(e)
	}
}
func TestPathsAndStrictDecoding(t *testing.T) {
	for _, p := range []string{"../escape", "/absolute", "a\\b", "a/../b", "C:/bad", "a//b", "a/CON", "a/", "a/name?", "a/tab\tfile"} {
		if portable(p) == nil {
			t.Fatalf("accepted %q", p)
		}
	}
	var q CreateRequest
	for _, b := range []string{`{"title":"a","title":"b"}`, `{"unknown":1}`, `{} {}`} {
		if Decode([]byte(b), &q) == nil {
			t.Fatal(b)
		}
	}
	s, caller := fixture(t)
	q = request(t, s, op1)
	if e := os.Symlink(t.TempDir(), filepath.Join(s.Area, "Activities")); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Create(q, caller); e == nil {
		t.Fatal("accepted symlink")
	}
}
func TestRealManualReplayBytePreservationAndHistoricalSnapshots(t *testing.T) {
	area := filepath.Join("..", "..", "SDP")
	ledger, e := os.ReadFile(filepath.Join(area, "ProjectManagement/Ledger.ndjson"))
	if e != nil {
		t.Fatal(e)
	}
	prefix := manualPilotPrefix(t, ledger)
	manualArea := t.TempDir()
	for _, home := range []string{"Activity-0001--SDPTool-Restructuring", "Activity-0002--Activity-Management"} {
		rel := filepath.Join("Activities", home, "activity.yaml")
		raw, err := os.ReadFile(filepath.Join(area, rel))
		if err != nil {
			t.Fatal(err)
		}
		target := filepath.Join(manualArea, rel)
		if err = os.MkdirAll(filepath.Dir(target), 0700); err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(target, raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	c, e := Replay(manualArea, "SDP", projecthistory.Snapshot{Bytes: prefix, Revision: digest(prefix)})
	if e != nil {
		t.Fatal(e)
	}
	if len(c.Activities) != 2 || c.Activities[0].State != "completed" || c.Activities[1].State != "active" {
		t.Fatalf("%+v", c)
	}
	after, _ := os.ReadFile(filepath.Join(area, "ProjectManagement/Ledger.ndjson"))
	if !bytes.Equal(ledger, after) {
		t.Fatal("manual bytes changed")
	}
	s, caller := fixture(t)
	r, e := s.Create(request(t, s, op1), caller)
	if e != nil {
		t.Fatal(e)
	}
	cat, _ := s.Read()
	p := cat.Activities[0].DescriptorRef.Path
	b, _ := s.read(p)
	if e = os.WriteFile(filepath.Join(s.Area, filepath.FromSlash(p)), append(b, []byte("# later descriptor bytes\n")...), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Read(); e != nil {
		t.Fatal(e)
	}
	pin := cat.Activities[0].DescriptorRef
	os.WriteFile(filepath.Join(s.Area, "Activities/.descriptors/"+pin.SHA256+".yaml"), []byte("forged"), 0600)
	if _, e = s.Read(); e == nil {
		t.Fatal("accepted forged snapshot")
	}
	_ = r
}
func TestReplayRejectsCrossProjectAndMissingNullableField(t *testing.T) {
	s, caller := fixture(t)
	if _, e := s.Create(request(t, s, op1), caller); e != nil {
		t.Fatal(e)
	}
	b, _ := os.ReadFile(filepath.Join(s.Area, "ProjectManagement/Ledger.ndjson"))
	for _, kind := range []string{"cross", "missing"} {
		var v map[string]any
		json.Unmarshal(bytes.TrimSpace(b), &v)
		p := v["payload"].(map[string]any)
		if kind == "cross" {
			p["projectId"] = "OTHER"
		} else {
			delete(p, "previousEventId")
		}
		raw, _ := json.Marshal(v)
		raw = append(raw, '\n')
		if _, e := Replay(s.Area, s.ProjectID, projecthistory.Snapshot{Bytes: raw, Revision: digest(raw)}); e == nil {
			t.Fatal(kind)
		}
	}
}

func captureRecords(t *testing.T, s *Service) map[string][]byte {
	t.Helper()
	records := map[string][]byte{}
	e := filepath.WalkDir(s.Area, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() {
			return nil
		}
		rel, e := filepath.Rel(s.Area, p)
		if e != nil {
			return e
		}
		if filepath.Ext(p) == ".lock" {
			return nil
		}
		b, e := os.ReadFile(p)
		if e != nil {
			return e
		}
		records[filepath.ToSlash(rel)] = b
		return nil
	})
	if e != nil {
		t.Fatal(e)
	}
	return records
}
func TestCapturedReplayMatchesLiveAndNoPartialCatalog(t *testing.T) {
	s, c := fixture(t)
	if _, e := s.Create(request(t, s, op1), c); e != nil {
		t.Fatal(e)
	}
	q := sessionRequest(t, s)
	if _, e := s.CreateSession(q, c); e != nil {
		t.Fatal(e)
	}
	live, e := s.Read()
	if e != nil {
		t.Fatal(e)
	}
	snap, e := projecthistory.Read(filepath.Join(s.Area, "ProjectManagement/Ledger.ndjson"))
	if e != nil {
		t.Fatal(e)
	}
	records := captureRecords(t, s)
	captured, e := ReplayRecords(s.ProjectID, snap, records)
	if e != nil {
		t.Fatal(e)
	}
	a, _ := json.Marshal(live)
	b, _ := json.Marshal(captured)
	if !bytes.Equal(a, b) {
		t.Fatal("captured/live reducer diverged")
	}
	for p := range records {
		if filepath.Ext(p) == ".md" && bytes.Contains([]byte(p), []byte(".sessions")) {
			delete(records, p)
		}
	}
	invalid, e := ReplayRecords(s.ProjectID, snap, records)
	if e == nil || len(invalid.Activities) != 0 {
		t.Fatal("invalid replay exposed partial catalog", e)
	}
}
func TestMissingDuplicateAndUnregisteredDescriptorHomes(t *testing.T) {
	for _, kind := range []string{"missing", "duplicate", "unregistered"} {
		t.Run(kind, func(t *testing.T) {
			s, c := fixture(t)
			r, e := s.Create(request(t, s, op1), c)
			if e != nil {
				t.Fatal(e)
			}
			cat, _ := s.Read()
			b, _ := s.read(cat.Activities[0].DescriptorRef.Path)
			switch kind {
			case "missing":
				os.Remove(filepath.Join(s.Area, filepath.FromSlash(cat.Activities[0].DescriptorRef.Path)))
			case "duplicate":
				s.write("Activities/Activity-0001--Duplicate/activity.yaml", b)
			case "unregistered":
				b = bytes.ReplaceAll(b, []byte("ACTIVITY-TEST-0001"), []byte("ACTIVITY-TEST-0099"))
				s.write("Activities/Activity-0099--Unknown/activity.yaml", b)
			}
			_ = r
			if invalid, e := s.Read(); e == nil || len(invalid.Activities) != 0 {
				t.Fatal("invalid homes accepted", e)
			}
		})
	}
}

func TestSessionPhaseContributionsMilestonesAndHandoff(t *testing.T) {
	s, c := fixture(t)
	if _, e := s.Create(request(t, s, op1), c); e != nil {
		t.Fatal(e)
	}
	if e := s.write("Plan.md", []byte("# Plan\n\n| Field | Value |\n| --- | --- |\n| id | PLAN-TEST-0001 |\n\n| Milestone | Outcome |\n| --- | --- |\n| A1 Study | Accepted study |\n")); e != nil {
		t.Fatal(e)
	}
	if e := s.write("Archive/Prior.md", []byte("# Prior Session\n\n| Field | Value |\n| --- | --- |\n| Session reference | SESSION-TEST-0012 |\n")); e != nil {
		t.Fatal(e)
	}
	if e := s.write("Handoff.md", []byte("# Handoff\n\n| Field | Value |\n| --- | --- |\n| id | HANDOFF-TEST-0001 |\n")); e != nil {
		t.Fatal(e)
	}
	q := sessionRequest(t, s)
	q.Phase = ""
	q.PhaseContributions = []string{"02-Study", "04-Design"}
	q.Plan = &Reference{Role: "plan", ID: "PLAN-TEST-0001", Path: "Plan.md"}
	q.MilestoneRefs = []Reference{{Role: "milestone", ID: "PLAN-TEST-0001:A1", Path: "Plan.md"}}
	q.Predecessor = &Reference{Role: "predecessor", ID: "SESSION-TEST-0012", Path: "Archive/Prior.md"}
	q.Handoff = &Reference{Role: "handoff", ID: "HANDOFF-TEST-0001", Path: "Handoff.md"}
	r, e := s.CreateSession(q, c)
	if e != nil || r.SessionID != "SESSION-TEST-0013" {
		t.Fatal(r, e)
	}
	b, e := s.read(r.Path)
	if e != nil {
		t.Fatal(e)
	}
	for _, text := range []string{"| Phase contributions | 02-Study, 04-Design |", "| Milestones | PLAN-TEST-0001:A1 |", "| Predecessor | SESSION-TEST-0012 |", "| Handoff | HANDOFF-TEST-0001 |", "../../../Plan.md", "../../../Archive/Prior.md"} {
		if !bytes.Contains(b, []byte(text)) {
			t.Fatal("missing real linked Session context", text, string(b))
		}
	}
	cat, e := s.Read()
	if e != nil {
		t.Fatal(e)
	}
	last := cat.Activities[0].Events[len(cat.Activities[0].Events)-1]
	if last.Payload.Session == nil || last.Payload.Session.SHA256 != digest(b) {
		t.Fatal("Session event did not pin complete authored context")
	}
}
func TestSessionRefusesInvalidPhaseAndSelectedRecordLinks(t *testing.T) {
	for _, kind := range []string{"missing-phase", "duplicate-phase", "invalid-phase", "conflicting-phase", "wrong-plan", "missing-milestone", "wrong-predecessor", "unavailable-handoff"} {
		t.Run(kind, func(t *testing.T) {
			s, c := fixture(t)
			if _, e := s.Create(request(t, s, op1), c); e != nil {
				t.Fatal(e)
			}
			s.write("Plan.md", []byte("| id | PLAN-TEST-0001 |\n| A1 Study | Result |\n"))
			q := sessionRequest(t, s)
			switch kind {
			case "missing-phase":
				q.Phase = ""
			case "duplicate-phase":
				q.PhaseContributions = []string{"04-Design", "04-Design"}
			case "invalid-phase":
				q.PhaseContributions = []string{"03-Architecture"}
			case "conflicting-phase":
				q.PhaseContributions = []string{"02-Study"}
			case "wrong-plan":
				q.Plan = &Reference{Role: "plan", ID: "PLAN-TEST-9999", Path: "Plan.md"}
			case "missing-milestone":
				q.MilestoneRefs = []Reference{{Role: "milestone", ID: "PLAN-TEST-0001:A9", Path: "Plan.md"}}
			case "wrong-predecessor":
				q.Predecessor = &Reference{Role: "session", ID: "SESSION-TEST-9999", Path: "Plan.md"}
			case "unavailable-handoff":
				q.Handoff = &Reference{Role: "handoff", ID: "HANDOFF-OTHER-0001", Path: "missing.md", ProjectID: "OTHER", Availability: "unavailable", UnavailableReason: "not installed"}
			}
			if _, e := s.CreateSession(q, c); e == nil {
				t.Fatal("invalid Session context accepted", kind)
			}
			cat, e := s.Read()
			if e != nil || len(cat.Activities[0].Sessions) != 0 {
				t.Fatal("invalid Session context produced records", e)
			}
		})
	}
}
