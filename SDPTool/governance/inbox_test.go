package governance

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReportInboxPersistenceAndCurrentCoverage(t *testing.T) {
	s, sel, q := fixture(t)
	empty, err := s.ReportInbox("")
	if err != nil || len(empty.Reports) != 0 || empty.Generation != "" {
		t.Fatal(empty, err)
	}
	if _, err = os.Stat(s.state); !os.IsNotExist(err) {
		t.Fatal("read created state", err)
	}
	bind(t, s, &q)
	r, err := s.Start(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	r, err = s.Transition(LocalController(), transition(r, "start"))
	if err != nil {
		t.Fatal(err)
	}
	report := transition(r, "report")
	report.Reason = "Candidate checked; review remains."
	submitted, err := s.Transition(LocalController(), report)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.Transition(LocalController(), report); err != nil {
		t.Fatal(err)
	}
	r, err = s.Transition(LocalController(), transition(submitted, "cancel"))
	if err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(filepath.Join(s.state, "HEAD"))
	reopened, err := New(s.root, s.area, s.project, s.state)
	if err != nil {
		t.Fatal(err)
	}
	view, err := reopened.ReportInbox(r.Run.ID)
	if err != nil || len(view.Reports) != 1 {
		t.Fatal(view, err)
	}
	v := view.Reports[0]
	if v.ID != report.OperationID || v.ReportedRevision != submitted.Run.Revision || v.CurrentRevision != r.Run.Revision || v.CurrentState != "canceled" || v.Text != report.Reason || v.Coverage != "known" || view.Consumption != "not-tracked" {
		t.Fatal(v)
	}
	if err = os.WriteFile(filepath.Join(s.root, sel.Authority), []byte("Changed owner mandate"), 0600); err != nil {
		t.Fatal(err)
	}
	view, err = reopened.ReportInbox("")
	if err != nil || len(view.Reports) != 1 || view.Reports[0].Coverage != "stale" || view.Reports[0].Text != report.Reason {
		t.Fatal(view, err)
	}
	after, _ := os.ReadFile(filepath.Join(s.state, "HEAD"))
	if !reflect.DeepEqual(before, after) {
		t.Fatal("inbox mutated store")
	}
	absent, _ := identifier()
	_, err = s.ReportInbox(absent)
	code(t, err, "missing_coverage")
	_, err = s.ReportInbox("invalid")
	code(t, err, "invalid_request")
}
