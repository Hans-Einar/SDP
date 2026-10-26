package governance

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func sessionTerminalBinary(t *testing.T) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "sdp-controller")
	if b, e := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", p, "../cmd/sdp-controller").CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	}
	return p
}
func terminalCall(binary string, s *Service, action, input string, extra ...string) ([]byte, error) {
	args := []string{"session", action, "--project", s.root, "--state", s.state}
	if action != "list" {
		args = append(args, "--session", "SESSION-TEST-0001")
	}
	cmd := exec.Command(binary, append(args, extra...)...)
	cmd.Stdin = strings.NewReader(input)
	return cmd.CombinedOutput()
}

func TestSessionTerminalClosureAndRecovery(t *testing.T) {
	binary := sessionTerminalBinary(t)
	s, d := dispositionFixture(t, "", "", nil)
	completeFixtureDisposition(t, s, d)
	head := mustFile(t, filepath.Join(s.state, "HEAD"))
	for _, action := range []string{"list", "status"} {
		b, e := terminalCall(binary, s, action, "")
		if e != nil || !bytes.Contains(b, []byte("SESSION-TEST-0001")) {
			t.Fatal(e, string(b))
		}
	}
	if !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
		t.Fatal("terminal inspection mutated state")
	}
	b, e := terminalCall(binary, s, "complete", "", "--reason", "Complete verified fixture", "--expected-revision", "1")
	if e == nil || !bytes.Contains(b, []byte("stale_revision")) || !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
		t.Fatal(e, string(b))
	}
	b, e = terminalCall(binary, s, "complete", "", "--reason", "Complete verified fixture", "--prepare-only")
	if e != nil || !bytes.Contains(b, []byte("prepared")) {
		t.Fatal(e, string(b))
	}
	v, _ := s.ReadSession(d.InstanceID)
	if v.State != "active" || v.RecordEffects != "prepared" {
		t.Fatal(v)
	}
	ledger := mustFile(t, filepath.Join(s.root, managementLedger))
	b, e = terminalCall(binary, s, "resume", "", "--reason", "Recover prepared terminal completion")
	if e != nil || !bytes.Contains(b, []byte("completed")) {
		t.Fatal(e, string(b))
	}
	if !bytes.HasPrefix(mustFile(t, filepath.Join(s.root, managementLedger)), ledger) {
		t.Fatal("rewrote old history")
	}
	v, e = s.ReadSession(d.InstanceID)
	if e != nil || v.State != "completed" {
		t.Fatal(v, e)
	}
	head = mustFile(t, filepath.Join(s.state, "HEAD"))
	if b, e = terminalCall(binary, s, "resume", "", "--reason", "Inspect already committed completion"); e != nil || !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
		t.Fatal("duplicate resume mutated state", e, string(b))
	}
	if _, e = s.publicationFiles(); e != nil {
		t.Fatal(e)
	}
}

func TestSessionTerminalExactOwnerGate(t *testing.T) {
	binary := sessionTerminalBinary(t)
	for _, action := range []string{"accept", "changes-requested"} {
		t.Run(action, func(t *testing.T) {
			s, q, _ := ownerFixture(t)
			commitOwner(t, s, LocalController(), q)
			head := mustFile(t, filepath.Join(s.state, "HEAD"))
			b, e := terminalCall(binary, s, action, "", "--gate", newTestID(), "--reason", "Synthetic terminal decision")
			if e == nil || !bytes.Contains(b, []byte("exact currently displayed")) || !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
				t.Fatal(e, string(b))
			}
			b, e = terminalCall(binary, s, action, "", "--gate", q.GateID, "--reason", "Synthetic terminal decision")
			if e != nil {
				t.Fatal(e, string(b))
			}
			st, _, e := s.load()
			g := latestOwnerGate(st.Sessions[q.InstanceID])
			if e != nil || g == nil || g.Request.Action != action || g.Effects.Status != "committed" || g.Request.ResultID != q.ResultID {
				t.Fatal(g, e)
			}
			head = mustFile(t, filepath.Join(s.state, "HEAD"))
			if b, e = terminalCall(binary, s, action, "", "--gate", q.GateID, "--reason", "Duplicate decision must not mutate"); e == nil || !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
				t.Fatal(e, string(b))
			}
			if action == "changes-requested" {
				step := selectedSessionStep(st.Sessions[q.InstanceID]).ID
				if b, e = terminalCall(binary, s, "rework", "", "--step", step, "--basis", "owner", "--result", g.Request.OperationID, "--reason", "Rework the rejected fixture"); e != nil {
					t.Fatal(e, string(b))
				}
				if b, e = terminalCall(binary, s, "select", "", "--step", step, "--reason", "Select fresh rework"); e != nil {
					t.Fatal(e, string(b))
				}
				v, e := s.ReadSession(q.InstanceID)
				if e != nil || len(v.History) != 1 || !v.Steps[0].Selected || v.Steps[0].RunID != "" || v.Steps[0].OwnerGate != nil {
					t.Fatal(v, e)
				}
			}
		})
	}
}

func TestSessionTerminalNextAndBroaderScopeHold(t *testing.T) {
	binary := sessionTerminalBinary(t)
	s, d := dispositionFixture(t, "", "", func(t *testing.T, s *Service, b *SessionBinding) {
		multiCardImplementation(t, s, b)
		next := b.Steps[0]
		next.ID = "second"
		next.Cards = []SessionCardRef{{"KB-SDP-001", "whole"}}
		next.DependsOn = []string{b.Steps[0].ID}
		b.Steps = append(b.Steps, next)
		evidenceFile(t, s, b.Session.Path, string(mustFile(t, filepath.Join(s.root, b.Session.Path)))+"\n## second\n\nComplete the remaining card contribution.\n")
	})
	completeFixtureDisposition(t, s, d)
	head := mustFile(t, filepath.Join(s.state, "HEAD"))
	if b, e := terminalCall(binary, s, "complete", "", "--reason", "Must retain broader scope"); e == nil || !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
		t.Fatal(e, string(b))
	}
	if b, e := terminalCall(binary, s, "next", "", "--step", "second", "--reason", "Continue the existing selected plan"); e != nil {
		t.Fatal(e, string(b))
	}
	if b, e := terminalCall(binary, s, "select", "", "--step", "second", "--reason", "Select the dependency-ready step"); e != nil {
		t.Fatal(e, string(b))
	}
	v, e := s.ReadSession(d.InstanceID)
	if e != nil || len(v.History) != 1 || !v.Steps[1].Selected || v.Steps[0].Stage != "done" || v.Steps[1].RunID != "" {
		t.Fatal(v, e)
	}
	evidenceFile(t, s, "SDP/authority.md", "Changed authority must require explicit reconciliation.\n")
	head = mustFile(t, filepath.Join(s.state, "HEAD"))
	if b, e := terminalCall(binary, s, "resume", "", "--reason", "Cannot silently refresh authority"); e == nil || !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
		t.Fatal(e, string(b))
	}
}

func TestSessionTerminalCurrentQuestionReply(t *testing.T) {
	binary := sessionTerminalBinary(t)
	s, q := workFixture(t, "", "")
	_, r := startWork(t, s, q)
	r, e := s.RecordHost(LocalController(), HostRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: newTestID(), RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: "bind", AttemptID: newTestID(), HostVersion: "0.164.0", ThreadID: "terminal-question-fixture", Observation: Observation{Kind: "bound", Epoch: newTestID(), ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano), Status: "idle"}})
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.Transition(LocalController(), transition(r, "start"))
	if e != nil {
		t.Fatal(e)
	}
	// Controller-authored question is enough to verify terminal capture. This is
	// a service/protocol fixture, not a native agent observation.
	r, e = s.RequestAttention(LocalController(), questionRequest(r))
	if e != nil {
		t.Fatal(e)
	}
	head := mustFile(t, filepath.Join(s.state, "HEAD"))
	for _, input := range []string{"", "short\n"} {
		b, e := terminalCall(binary, s, "answer", input, "--owner", "fixture-owner")
		if e == nil || !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
			t.Fatal(e, string(b))
		}
	}
	b, e := terminalCall(binary, s, "answer", "short\nUse the short label.\n", "--owner", "fixture-owner")
	if e != nil || !bytes.Contains(b, []byte("Which label")) || !bytes.Contains(b, []byte("Reply recorded")) {
		t.Fatal(e, string(b))
	}
	inbox, e := s.AttentionInbox(r.Run.ID)
	if e != nil || len(inbox.Items) != 1 || inbox.Items[0].Question.Reply == nil || inbox.Items[0].Question.Reply.Text != "Use the short label." {
		t.Fatal(inbox, e)
	}
	if inbox.Items[0].Question.Reply.Source != "local-controller-stdin" || inbox.Items[0].Question.Reply.Actor != "fixture-owner" {
		t.Fatal("lost actual terminal attribution")
	}
	head = mustFile(t, filepath.Join(s.state, "HEAD"))
	if b, e = terminalCall(binary, s, "answer", "short\nDuplicate\n", "--owner", "fixture-owner", "--question", inbox.Items[0].Question.ID); e == nil || !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
		t.Fatal(e, string(b))
	}
	if b, e = terminalCall(binary, s, "status", "", "--expected-revision", fmt.Sprint(q.ExpectedRevision)); e == nil {
		t.Fatal("accepted obsolete inspected revision", string(b))
	}
}
