package governance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCoordinationPagingFreshnessAndReadOnly(t *testing.T) {
	s, sel, start := fixture(t)
	empty, e := s.Coordination(CoordinationRequest{})
	if e != nil || empty.Total != 0 || empty.Generation != "" {
		t.Fatal(empty, e)
	}
	if _, e = os.Stat(s.state); !os.IsNotExist(e) {
		t.Fatal("read created store", e)
	}
	bind(t, s, &start)
	r, e := s.Start(LocalController(), start)
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.Transition(LocalController(), transition(r, "start"))
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.RequestAttention(LocalController(), questionRequest(r))
	if e != nil {
		t.Fatal(e)
	}
	_, _, other := fixture(t)
	other.CheckoutID = r.CheckoutID
	other.InputRevision = r.Run.InputRevision
	other.Assignment.WritablePaths = []string{"other.txt"}
	second, e := s.Start(LocalController(), other)
	if e != nil {
		t.Fatal(e)
	}
	second, e = s.Transition(LocalController(), transition(second, "start"))
	if e != nil {
		t.Fatal(e)
	}
	second, e = s.Transition(LocalController(), transition(second, "report"))
	if e != nil {
		t.Fatal(e)
	}
	head := stateHead(t, s)
	first, e := s.Coordination(CoordinationRequest{Limit: 1})
	if e != nil || first.Total != 4 || len(first.Items) != 1 || first.Items[0].Kind != "attention" || first.NextCursor == "" {
		t.Fatal(first, e)
	}
	page, e := s.Coordination(CoordinationRequest{Limit: 1, Cursor: first.NextCursor})
	if e != nil || page.Items[0].Kind != "report" {
		t.Fatal(page, e)
	}
	seen := map[string]bool{first.Items[0].ID: true, page.Items[0].ID: true}
	for page.NextCursor != "" {
		page, e = s.Coordination(CoordinationRequest{Limit: 1, Cursor: page.NextCursor})
		if e != nil {
			t.Fatal(e)
		}
		for _, i := range page.Items {
			if seen[i.ID] {
				t.Fatal("duplicate page item")
			}
			seen[i.ID] = true
		}
	}
	if len(seen) != 4 || stateHead(t, s) != head {
		t.Fatal("lost items or read mutation")
	}
	again, e := s.Coordination(CoordinationRequest{Limit: 50})
	if e != nil || again.Revision != first.Revision {
		t.Fatal("inspection timestamp changed revision", e)
	}
	filtered, e := s.Coordination(CoordinationRequest{RunID: second.Run.ID})
	if e != nil || filtered.Total != 2 {
		t.Fatal(filtered, e)
	}
	_, e = s.Coordination(CoordinationRequest{RunID: second.Run.ID, Cursor: first.NextCursor})
	code(t, e, "stale_revision")
	if e = os.WriteFile(filepath.Join(s.root, sel.Authority), []byte("changed owner basis"), 0600); e != nil {
		t.Fatal(e)
	}
	_, e = s.Coordination(CoordinationRequest{Cursor: first.NextCursor})
	code(t, e, "stale_revision")
	changed, e := s.Coordination(CoordinationRequest{})
	if e != nil || changed.Generation != head || changed.Revision == first.Revision || changed.Items[0].Attention.Applicability != "stale" {
		t.Fatal(changed, e)
	}
}
func TestCoordinationAcknowledgementInvalidatesCursor(t *testing.T) {
	s, _, _, q := reportFixture(t)
	first, e := s.Coordination(CoordinationRequest{Limit: 1})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.AcknowledgeReport(LocalController(), q); e != nil {
		t.Fatal(e)
	}
	_, e = s.Coordination(CoordinationRequest{Cursor: first.NextCursor})
	code(t, e, "stale_revision")
	fresh, e := s.Coordination(CoordinationRequest{})
	if e != nil || fresh.Items[0].Kind != "work" || fresh.Items[1].Report.Acknowledgement == nil {
		t.Fatal(fresh, e)
	}
}
func TestCoordinationRejectsInvalidRequests(t *testing.T) {
	s, _, _ := fixture(t)
	for _, q := range []CoordinationRequest{{Limit: -1}, {Limit: 51}, {RunID: "invalid"}, {Cursor: "?"}} {
		if _, e := s.Coordination(q); e == nil {
			t.Fatal(q)
		}
	}
}

func TestCoordinationPageByteBound(t *testing.T) {
	s, _, q := fixture(t)
	bind(t, s, &q)
	q.Assignment.Outcome = strings.Repeat("x", 60000)
	for i := 0; i < 4; i++ {
		q.OperationID, _ = identifier()
		q.RequestID, _ = identifier()
		if _, e := s.Start(LocalController(), q); e != nil {
			t.Fatal(e)
		}
	}
	view, e := s.Coordination(CoordinationRequest{Limit: 50})
	if e != nil {
		t.Fatal(e)
	}
	b, _ := json.Marshal(view)
	if len(b) > 256<<10 || view.NextCursor == "" || len(view.Items) >= 4 {
		t.Fatal("page not bounded", len(b), len(view.Items))
	}
	next, e := s.Coordination(CoordinationRequest{Limit: 50, Cursor: view.NextCursor})
	if e != nil || len(next.Items)+len(view.Items) != 4 {
		t.Fatal(next, e)
	}
}
