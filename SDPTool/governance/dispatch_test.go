package governance

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDispatchCrashAndImmutableIntent(t *testing.T) {
	for _, point := range []string{"before-head", "after-head"} {
		t.Run(point, func(t *testing.T) {
			s, _, start := fixture(t)
			bind(t, s, &start)
			r, e := s.Start(LocalController(), start)
			if e != nil {
				t.Fatal(e)
			}
			id, _ := identifier()
			at, _ := identifier()
			epoch, _ := identifier()
			op, _ := identifier()
			q := InputRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: op, RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: "prepare", Input: OwnerInput{ID: id, Kind: "launch", Text: "owner input\nexact", Payload: "bounded assignment", AttemptID: at, Epoch: epoch, ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}}
			_, e = s.RecordInput(Caller{}, q)
			code(t, e, "unauthorized")
			r, e = s.RecordInput(LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			q.Action = "claim"
			q.OperationID, _ = identifier()
			q.ExpectedRevision = r.Run.Revision
			q.Input.State = "prepared"
			s.fault = func(p string) error {
				if p == point {
					return errors.New("injected crash")
				}
				return nil
			}
			if _, e = s.RecordInput(LocalController(), q); e == nil {
				t.Fatal("fault absent")
			}
			reopened, e := New(s.root, s.area, s.project, s.state)
			if e != nil {
				t.Fatal(e)
			}
			v, e := reopened.OwnerInput(LocalController(), r.Run.ID, id)
			if e != nil {
				t.Fatal(e)
			}
			want := "prepared"
			if point == "after-head" {
				want = "sending"
			}
			if v.State != want || v.Text != "owner input\nexact" {
				t.Fatal(v)
			}
			// Retrying the state operation resolves its receipt; it does not send anything.
			retry, e := reopened.RecordInput(LocalController(), q)
			if e != nil || retry.Run.Inputs[0].State != "sending" {
				t.Fatal(retry, e)
			}
			q.OperationID, _ = identifier()
			q.ExpectedRevision = retry.Run.Revision
			_, e = reopened.RecordInput(LocalController(), q)
			if e == nil {
				t.Fatal("intent claimed twice")
			}
			q.Input.Text = "different"
			_, e = reopened.RecordInput(LocalController(), q)
			if e == nil {
				t.Fatal("input rewritten")
			}
			if dir := os.Getenv("PGI_SCHEMA_EVIDENCE"); dir != "" && point == "after-head" {
				b, _ := json.MarshalIndent([]map[string]any{{"definition": "InputRequest", "value": q}, {"definition": "OwnerInput", "value": v}, {"definition": "Result", "value": retry}}, "", "  ")
				if e = os.WriteFile(filepath.Join(dir, "dispatch-wire-cases.json"), b, 0600); e != nil {
					t.Fatal(e)
				}
			}
		})
	}
}
func TestPreparedInputCannotAdvanceAfterDrift(t *testing.T) {
	s, _, start := fixture(t)
	bind(t, s, &start)
	r, e := s.Start(LocalController(), start)
	if e != nil {
		t.Fatal(e)
	}
	id, _ := identifier()
	at, _ := identifier()
	epoch, _ := identifier()
	op, _ := identifier()
	q := InputRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: op, RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: "prepare", Input: OwnerInput{ID: id, Kind: "launch", Text: "input", Payload: "payload", AttemptID: at, Epoch: epoch, ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}}
	r, e = s.RecordInput(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(s.root, start.Selection.Authority), []byte("Changed authority"), 0600); e != nil {
		t.Fatal(e)
	}
	q.Action = "claim"
	q.OperationID, _ = identifier()
	q.ExpectedRevision = r.Run.Revision
	q.Input.State = "prepared"
	_, e = s.RecordInput(LocalController(), q)
	code(t, e, "stale_revision")
	v, _ := s.OwnerInput(LocalController(), r.Run.ID, id)
	if v.State != "prepared" {
		t.Fatal(v)
	}
}
