package governance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func leaderConfigFixture(t *testing.T, s *Service) LeaderConfig {
	t.Helper()
	path := filepath.Join(s.root, "Skills/leader/SKILL.md")
	b := []byte("bounded leader fixture")
	if e := os.MkdirAll(filepath.Dir(path), 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(path, b, 0600); e != nil {
		t.Fatal(e)
	}
	return LeaderConfig{Coordinator: "leader", HostVersion: "0.164.0", MCPBinary: "/fixture/sdp-mcp", Instructions: "bounded coordination", Skills: []SkillInput{{Name: "leader", Path: path, Digest: hash(b)}}}
}
func leaderRequest(t *testing.T, s *Service, action, id string) LeaderRequest {
	t.Helper()
	v, e := s.LeaderRead()
	if e != nil {
		t.Fatal(e)
	}
	op, _ := identifier()
	q := LeaderRequest{Schema: LeaderSchema, ProjectRef: v.ProjectRef, CheckoutID: v.CheckoutID, OperationID: op, ExpectedRevision: v.Leader.Revision, Action: action, DispatchID: id}
	if action == "message" {
		q.ExpectedRevision = 0
		q.Message = &LeaderMessageInput{Actor: "owner", Source: "fixture", Text: "Review current work"}
	}
	return q
}
func leaderApply(t *testing.T, s *Service, q LeaderRequest) LeaderResult {
	t.Helper()
	r, e := s.LeaderWrite(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	return r
}
func completeLeader(t *testing.T, s *Service, d LeaderDispatch) {
	t.Helper()
	leaderApply(t, s, leaderRequest(t, s, "open", d.ID))
	q := leaderRequest(t, s, "bind", d.ID)
	q.ThreadID = "thread"
	leaderApply(t, s, q)
	leaderApply(t, s, leaderRequest(t, s, "send", d.ID))
	q = leaderRequest(t, s, "observe", d.ID)
	q.ThreadID = "thread"
	q.TurnID = "turn"
	q.Status = "completed"
	q.Answer = "Proposed next step, no approval"
	leaderApply(t, s, q)
}
func TestLeaderQueueSerializationSupersessionAndRestart(t *testing.T) {
	s, _, start := fixture(t)
	v, e := s.LeaderRead()
	if e != nil || v.Leader.Revision != 0 {
		t.Fatal(v, e)
	}
	if _, e = os.Stat(s.state); !os.IsNotExist(e) {
		t.Fatal("inspection mutated state")
	}
	bind(t, s, &start)
	config := leaderConfigFixture(t, s)
	message := leaderRequest(t, s, "message", "")
	first := leaderApply(t, s, message)
	q := leaderRequest(t, s, "prepare", "")
	q.Config = &config
	d := *leaderApply(t, s, q).Dispatch
	other := leaderRequest(t, s, "prepare", "")
	other.Config = &config
	_, e = s.LeaderWrite(LocalController(), other)
	code(t, e, "conflict")
	leaderApply(t, s, leaderRequest(t, s, "message", ""))
	_, e = s.LeaderWrite(LocalController(), leaderRequest(t, s, "open", d.ID))
	code(t, e, "stale_revision")
	leaderApply(t, s, leaderRequest(t, s, "supersede", d.ID))
	q = leaderRequest(t, s, "prepare", "")
	q.Config = &config
	d = *leaderApply(t, s, q).Dispatch
	if len(d.MessageIDs) != 2 {
		t.Fatal("supersession dropped owner input", d)
	}
	completeLeader(t, s, d)
	head := stateHead(t, s)
	retry := leaderApply(t, s, message)
	if !reflect.DeepEqual(first, retry) || head != stateHead(t, s) {
		t.Fatal("enqueue retry changed receipt")
	}
	s, e = New(s.root, s.area, s.project, s.state)
	if e != nil {
		t.Fatal(e)
	}
	q = leaderRequest(t, s, "prepare", "")
	q.Config = &config
	_, e = s.LeaderWrite(LocalController(), q)
	code(t, e, "no_work")
	message.Message.Text = "altered"
	_, e = s.LeaderWrite(LocalController(), message)
	code(t, e, "idempotency_conflict")
	state, _, e := s.load()
	if e != nil || state.Schema != leaderStateSchema {
		t.Fatal(e)
	}
	if len(state.Runs) != 0 {
		t.Fatal("leader queue created an assignment")
	}
	// Leader operations share the global operation namespace with worker services.
	start.OperationID = first.OperationID
	_, e = s.Start(LocalController(), start)
	code(t, e, "idempotency_conflict")
}
func TestLeaderReportKeysCoexistenceAndNoAcknowledgement(t *testing.T) {
	s, _, r, _ := reportFixture(t)
	config := leaderConfigFixture(t, s)
	q := leaderRequest(t, s, "prepare", "")
	q.Config = &config
	d := *leaderApply(t, s, q).Dispatch
	if len(d.NoticeKeys) != 1 || len(d.MessageIDs) != 0 {
		t.Fatal(d)
	}
	var payload struct {
		Notices []leaderNotice `json:"notices"`
	}
	if e := json.Unmarshal([]byte(d.Payload), &payload); e != nil {
		t.Fatal(e)
	}
	if payload.Notices[0].Item.Report.CurrentRevision == 0 {
		t.Fatal("normalization damaged actionable notice")
	}
	completeLeader(t, s, d)
	q = leaderRequest(t, s, "prepare", "")
	q.Config = &config
	_, e := s.LeaderWrite(LocalController(), q)
	code(t, e, "no_work")
	inbox, e := s.ReportInbox("")
	if e != nil || inbox.Reports[0].Acknowledgement != nil {
		t.Fatal("delivery acknowledged work", e)
	}
	ack := acknowledgementRequest(t, s, r)
	if _, e = s.AcknowledgeReport(LocalController(), ack); e != nil {
		t.Fatal(e)
	}
	state, _, e := s.load()
	if e != nil || state.Schema != leaderStateSchema || state.Runs[r.Run.ID].State != "verifying" {
		t.Fatal("profile downgraded or run advanced", e)
	}
}
func TestLeaderLostBoundariesCannotBeSuperseded(t *testing.T) {
	for _, boundary := range []string{"opening", "sending", "running"} {
		t.Run(boundary, func(t *testing.T) {
			s, _, start := fixture(t)
			bind(t, s, &start)
			config := leaderConfigFixture(t, s)
			leaderApply(t, s, leaderRequest(t, s, "message", ""))
			q := leaderRequest(t, s, "prepare", "")
			q.Config = &config
			d := *leaderApply(t, s, q).Dispatch
			leaderApply(t, s, leaderRequest(t, s, "open", d.ID))
			if boundary != "opening" {
				q = leaderRequest(t, s, "bind", d.ID)
				q.ThreadID = "thread"
				leaderApply(t, s, q)
				leaderApply(t, s, leaderRequest(t, s, "send", d.ID))
			}
			if boundary == "running" {
				q = leaderRequest(t, s, "observe", d.ID)
				q.ThreadID = "thread"
				q.TurnID = "turn"
				q.Status = "running"
				leaderApply(t, s, q)
			}
			leaderApply(t, s, leaderRequest(t, s, "message", ""))
			_, e := s.LeaderWrite(LocalController(), leaderRequest(t, s, "supersede", d.ID))
			code(t, e, "reconciliation_required")
			q = leaderRequest(t, s, "prepare", "")
			q.Config = &config
			_, e = s.LeaderWrite(LocalController(), q)
			code(t, e, "conflict")
		})
	}
}
func TestLeaderLeaseAndSnapshotValidation(t *testing.T) {
	s, _, start := fixture(t)
	bind(t, s, &start)
	config := leaderConfigFixture(t, s)
	unlock, e := s.LeaderLease(LocalController())
	if e != nil {
		t.Fatal(e)
	}
	defer unlock()
	_, e = s.LeaderLease(LocalController())
	code(t, e, "busy")
	leaderApply(t, s, leaderRequest(t, s, "message", "")) // Queue does not take leader lease.
	q := leaderRequest(t, s, "prepare", "")
	q.Config = &config
	d := *leaderApply(t, s, q).Dispatch
	for _, change := range []string{"payload", "duplicate-active", "version", "receipt", "profile", "missing-message"} {
		t.Run(change, func(t *testing.T) {
			st, _, e := s.load()
			if e != nil {
				t.Fatal(e)
			}
			switch change {
			case "payload":
				st.Leader.Dispatches[0].Payload = "changed"
			case "duplicate-active":
				st.Leader.Dispatches = append(st.Leader.Dispatches, d)
			case "version":
				st.Leader.Dispatches[0].Config.HostVersion = "0.163.0"
			case "receipt":
				delete(st.LeaderReceipts, d.ID)
			case "profile":
				st.Schema = reportAcknowledgementStateSchema
			case "missing-message":
				st.Leader.Messages = nil
			}
			if e := validateLeaderState(st); e == nil {
				t.Fatal("corrupt leader accepted")
			}
		})
	}
	q = leaderRequest(t, s, "message", "")
	_, e = s.LeaderWrite(Caller{}, q)
	code(t, e, "unauthorized")
}

func TestLeaderMessageBoundsAndEncodedBatch(t *testing.T) {
	s, _, start := fixture(t)
	bind(t, s, &start)
	config := leaderConfigFixture(t, s)
	for i := 0; i < 64; i++ {
		q := leaderRequest(t, s, "message", "")
		q.Message.Text = strings.Repeat("\x01", 4096)
		leaderApply(t, s, q)
	}
	_, e := s.LeaderWrite(LocalController(), leaderRequest(t, s, "message", ""))
	code(t, e, "limit")
	q := leaderRequest(t, s, "prepare", "")
	q.Config = &config
	d := *leaderApply(t, s, q).Dispatch
	if len(d.Payload) > 64<<10 || len(d.MessageIDs) < 1 || len(d.MessageIDs) >= 8 {
		t.Fatal("encoded owner batch not bounded", len(d.Payload), len(d.MessageIDs))
	}
	for _, text := range []string{string([]byte{0xff}), "a\x00b", strings.Repeat("x", 4097)} {
		q := leaderRequest(t, s, "message", "")
		q.Message.Text = text
		_, e = s.LeaderWrite(LocalController(), q)
		code(t, e, "invalid_request")
	}
}
