package governance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestAssignmentAuthorityAndHostReplay(t *testing.T) {
	s, _, q := fixture(t)
	bind(t, s, &q)
	r, e := s.Start(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	at, _ := identifier()
	epoch, _ := identifier()
	agent, _ := AssignedAgent(r.Run.ID, at, epoch)
	_, e = s.Transition(agent, transition(r, "start"))
	code(t, e, "unauthorized")
	r, e = s.Transition(LocalController(), transition(r, "start"))
	if e != nil {
		t.Fatal(e)
	}
	other, _ := identifier()
	forged := transition(r, "report")
	forged.RunID = other
	_, e = s.Transition(agent, forged)
	code(t, e, "unauthorized")
	now := time.Now().UTC().Format(time.RFC3339Nano)
	request := func(action, kind string) HostRequest {
		op, _ := identifier()
		return HostRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, OperationID: op, Action: action, AttemptID: at, HostVersion: "0.160.1", ThreadID: "observed-thread", Observation: Observation{Kind: kind, Epoch: epoch, ReceivedAt: now, Status: "idle"}}
	}
	h := request("bind", "bound")
	_, e = s.RecordHost(agent, h)
	code(t, e, "unauthorized")
	r, e = s.RecordHost(LocalController(), h)
	if e != nil {
		t.Fatal(e)
	}
	hostResult := r
	again, e := s.RecordHost(LocalController(), h)
	if e != nil || again.Run.Revision != r.Run.Revision {
		t.Fatal(e, again)
	}
	h = request("observe", "agentMessage")
	h.Observation.Status = ""
	h.Observation.Phase = "commentary"
	h.Observation.Text = "Working"
	h.Observation.TurnID = "turn"
	h.Observation.ItemID = "item"
	r, e = s.RecordHost(LocalController(), h)
	if e != nil {
		t.Fatal(e)
	}
	h.OperationID, _ = identifier()
	h.ExpectedRevision = r.Run.Revision
	h.Observation.ReceivedAt = time.Now().UTC().Format(time.RFC3339Nano)
	r, e = s.RecordHost(LocalController(), h)
	if e != nil || r.Run.Attempts[0].EventCount != 2 {
		t.Fatal("replay duplicated item", e)
	}
	h.OperationID, _ = identifier()
	h.ExpectedRevision = r.Run.Revision
	h.Observation.Text = "conflicting replay"
	_, e = s.RecordHost(LocalController(), h)
	code(t, e, "observation_gap")
	page, e := s.Events(r.Run.ID, at, "", 1)
	if e != nil || page.NextCursor == "" {
		t.Fatal(page, e)
	}
	h = request("disconnect", "disconnected")
	h.Observation.Status = "unknown"
	r, e = s.RecordHost(LocalController(), h)
	if e != nil {
		t.Fatal(e)
	}
	_, e = s.Events(r.Run.ID, at, page.NextCursor, 1)
	code(t, e, "stale_revision")
	s2, e := New(s.root, s.area, s.project, s.state)
	if e != nil {
		t.Fatal(e)
	}
	r, e = s2.Read(r.Run.ID)
	if e != nil || r.Run.Attempts[0].Status != "unknown" || !r.Run.Attempts[0].Gap {
		t.Fatal(r, e)
	}
	h = request("observe", "turn.completed")
	h.Observation.TurnID = "turn"
	h.Observation.Status = "completed"
	_, e = s2.RecordHost(LocalController(), h)
	code(t, e, "reconciliation_required")
	epoch, _ = identifier()
	h = request("observe", "turn.completed")
	h.Observation.TurnID = "turn"
	h.Observation.Status = "completed"
	_, e = s2.RecordHost(LocalController(), h)
	code(t, e, "reconciliation_required")
	h = request("reconcile", "reconciled")
	h.Observation.TurnID = "turn"
	h.Observation.Status = "completed"
	r, e = s2.RecordHost(LocalController(), h)
	if e != nil {
		t.Fatal(e)
	}
	if r.Run.State != "running" || r.Run.Attempts[0].Status != "completed" {
		t.Fatal("host completion changed routine state", r)
	}
	_, e = s2.Transition(agent, transition(r, "complete"))
	code(t, e, "unauthorized")
	_, e = s2.Transition(agent, transition(r, "report"))
	code(t, e, "unauthorized")
	agent, _ = AssignedAgent(r.Run.ID, at, epoch)
	r, e = s2.Transition(agent, transition(r, "report"))
	if e != nil || r.Run.State != "verifying" {
		t.Fatal(r, e)
	}
	_, e = s2.Transition(LocalController(), transition(r, "complete"))
	code(t, e, "evidence_missing")
	_, e = s2.Events(r.Run.ID, at, "", 257)
	code(t, e, "invalid_request")
	page, e = s2.Events(r.Run.ID, at, "", 256)
	if e != nil {
		t.Fatal(e)
	}
	if dir := os.Getenv("PGI_SCHEMA_EVIDENCE"); dir != "" {
		cases := []map[string]any{{"definition": "Result", "value": hostResult}, {"definition": "HostRequest", "value": h}, {"definition": "Result", "value": r}, {"definition": "EventPage", "value": page}}
		b, e := json.MarshalIndent(cases, "", "  ")
		if e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(filepath.Join(dir, "host-wire-cases.json"), b, 0600); e != nil {
			t.Fatal(e)
		}
	}
}

func TestObservationPagesBoundBytesWithoutLosingContent(t *testing.T) {
	s, _, q := fixture(t)
	bind(t, s, &q)
	r, err := s.Start(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	at, _ := identifier()
	epoch, _ := identifier()
	record := func(kind, item, content string) {
		t.Helper()
		op, _ := identifier()
		h := HostRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, OperationID: op, Action: "observe", AttemptID: at, HostVersion: "0.160.1", ThreadID: "page-thread", Observation: Observation{Kind: kind, Epoch: epoch, ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano), ItemID: item, Text: content, TurnID: "turn", Phase: "final_answer"}}
		if kind == "bound" {
			h.Action = "bind"
			h.Observation.Status = "idle"
			h.Observation.TurnID = ""
			h.Observation.Phase = ""
		}
		r, err = s.RecordHost(LocalController(), h)
		if err != nil {
			t.Fatal(err)
		}
	}
	record("bound", "", "")
	content := strings.Repeat("x", 64<<10)
	for i := 0; i < 5; i++ {
		record("agentMessage", strconv.Itoa(i), content)
	}
	read, err := s.Read(r.Run.ID)
	if err != nil || read.Run.Attempts[0].EventCount != 6 || len(read.Run.Attempts[0].Events) != 0 {
		t.Fatal(read, err)
	}
	cursor := ""
	count, pages := 0, 0
	for {
		page, err := s.Events(r.Run.ID, at, cursor, 256)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := json.Marshal(page)
		if len(b) > 256<<10 {
			t.Fatal("oversized page", len(b))
		}
		for _, event := range page.Events {
			if event.Kind == "agentMessage" {
				if event.Text != content {
					t.Fatal("truncated text")
				}
				count++
			}
		}
		pages++
		cursor = page.NextCursor
		if cursor == "" {
			break
		}
	}
	if count != 5 || pages != 2 {
		t.Fatal(count, pages)
	}
	o := Observation{Kind: "agentMessage", TurnID: "turn", ItemID: "huge", Phase: "commentary", Text: strings.Repeat("\x00", 64<<10), Epoch: epoch, ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	if validObservation(o) {
		t.Fatal("encoded observation can exceed readable page")
	}
}
