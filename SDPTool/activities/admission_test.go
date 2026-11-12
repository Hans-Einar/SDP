package activities

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
)

type admissionInspector struct {
	snapshot ExecutionSnapshot
	err      error
}

func (i admissionInspector) InspectExecutions([]ExecutionStore) (ExecutionSnapshot, error) {
	return i.snapshot, i.err
}
func knownEmpty() ExecutionInspector {
	return admissionInspector{snapshot: ExecutionSnapshot{Known: true}}
}
func adoptedAdmission(t *testing.T) (*Service, Caller, AdoptionRequest) {
	t.Helper()
	s, caller := fixture(t)
	cat, err := s.Read()
	if err != nil {
		t.Fatal(err)
	}
	q := AdoptionRequest{OperationID: op1, ExpectedHistoryRevision: cat.HistoryRevision, Stores: []ExecutionStore{{StatePath: t.TempDir()}}}
	view, err := s.AdoptAdmission(q, caller, knownEmpty())
	if err != nil || !view.Enabled || view.ProjectID != "TEST" {
		t.Fatal(view, err)
	}
	return s, caller, q
}
func TestAdmissionExplicitAdoptionRetryAndMissingRegistry(t *testing.T) {
	s, caller, q := adoptedAdmission(t)
	if _, err := s.AdoptAdmission(q, caller, knownEmpty()); err != nil {
		t.Fatal(err)
	}
	changed := q
	changed.OperationID = op2
	if _, err := s.AdoptAdmission(changed, caller, knownEmpty()); err == nil {
		t.Fatal("conflicting adoption accepted")
	}
	// Governance uses a source slug separately: an unspecified namespace resolves
	// the canonical namespace pinned in policy, without uppercasing that slug.
	reader := New(s.Area, "")
	tx, err := projecthistory.BeginTransaction(filepath.Join(s.Area, "ProjectManagement/Ledger.ndjson"))
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Close()
	view, err := reader.ReadAdmission(tx, "")
	if err != nil || view.ProjectID != "TEST" {
		t.Fatal(view, err)
	}
	p := filepath.Join(s.Area, admissionDirectory, "stores", digest([]byte(q.Stores[0].StatePath))+".json")
	if err = os.Remove(p); err != nil {
		t.Fatal(err)
	}
	if _, err = reader.ReadAdmission(tx, ""); err == nil {
		t.Fatal("missing adopted inventory was disabled")
	}
}
func TestAdmissionUnknownWorkNeverMeansEmpty(t *testing.T) {
	for _, snapshot := range []ExecutionSnapshot{{Known: false}, {Known: true, InFlight: []ExecutionRef{{RunID: "running"}}}} {
		s, caller := fixture(t)
		c, err := s.Read()
		if err != nil {
			t.Fatal(err)
		}
		q := AdoptionRequest{OperationID: op1, ExpectedHistoryRevision: c.HistoryRevision, Stores: []ExecutionStore{{StatePath: t.TempDir()}}}
		if _, err = s.AdoptAdmission(q, caller, admissionInspector{snapshot: snapshot}); err == nil {
			t.Fatal("uncertain adoption accepted")
		}
		enabled, err := s.AdmissionEnabled()
		if err != nil || enabled {
			t.Fatal(enabled, err)
		}
	}
}
func TestAdmissionPendingUsesCanonicalGateWithoutReservation(t *testing.T) {
	s, _, _ := adoptedAdmission(t)
	ledger := filepath.Join(s.Area, "ProjectManagement/Ledger.ndjson")
	tx, err := projecthistory.BeginTransaction(ledger)
	if err != nil {
		t.Fatal(err)
	}
	pending := PendingLifecycle{OperationID: op2, ActivityID: "ACTIVITY-TEST-0001", Action: "hold"}
	if err = s.SetPendingLifecycle(tx, "", pending); err != nil {
		t.Fatal(err)
	}
	if err = s.SetPendingLifecycle(tx, "", pending); err != nil {
		t.Fatal(err)
	}
	other := pending
	other.OperationID = op1
	if err = s.SetPendingLifecycle(tx, "", other); err == nil {
		t.Fatal("pending operation replaced")
	}
	if err = s.RegisterExecutionStore(tx, "", ExecutionStore{StatePath: t.TempDir()}); err == nil {
		t.Fatal("pending admitted another store")
	}
	tx.Close()
	// Pending lifecycle does not reserve/block canonical history throughout a run.
	if _, err = projecthistory.Read(ledger); err != nil {
		t.Fatal(err)
	}
	tx, err = projecthistory.BeginTransaction(ledger)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Close()
	if err = s.ClearPendingLifecycle(tx, "", op1); err == nil {
		t.Fatal("foreign operation cleared gate")
	}
	if err = tx.Reserve("recovery"); err != nil {
		t.Fatal(err)
	}
	if _, err = s.ReadAdmission(tx, "recovery"); err != nil {
		t.Fatal(err)
	}
	if err = s.ClearPendingLifecycle(tx, "recovery", op2); err != nil {
		t.Fatal(err)
	}
	if err = tx.Release("recovery"); err != nil {
		t.Fatal(err)
	}
}
func TestAdmissionMembershipScopeAndDisabledCompatibility(t *testing.T) {
	work := WorkMembership{SessionPath: "SDP/Sessions/session.md", WritablePaths: []string{"Source/file.go"}}
	if _, err := CheckAdmission(AdmissionView{}, work); err != nil {
		t.Fatal(err)
	}
	view := AdmissionView{Enabled: true, Catalog: Catalog{Activities: []Activity{{State: "active", LastEventID: "event", DescriptorRef: DescriptorReference{SHA256: "revision"}, Descriptor: Descriptor{ID: "ACTIVITY-TEST-0001", References: []Reference{{Role: "session", Path: "Sessions/session.md"}}, Scope: []string{"Source"}}}}}}
	if _, err := CheckAdmission(view, work); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"SourceElse/file.go", "../Source/file.go", "Source/*.go"} {
		bad := work
		bad.WritablePaths = []string{path}
		if _, err := CheckAdmission(view, bad); err == nil {
			t.Fatal("unsafe scope accepted", path)
		}
	}
	bad := work
	bad.SessionPath = "SDP/Sessions/unrelated.md"
	if _, err := CheckAdmission(view, bad); err == nil {
		t.Fatal("unrelated Session accepted")
	}
	view.Pending = &PendingLifecycle{OperationID: op1}
	if _, err := CheckAdmission(view, work); err == nil {
		t.Fatal("pending admitted work")
	}
}
