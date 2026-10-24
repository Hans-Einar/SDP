package governancecontroller

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

type consoleEvent struct {
	Raw    json.RawMessage              `json:"-"`
	Schema string                       `json:"schema"`
	Kind   string                       `json:"kind"`
	View   *governance.CoordinationView `json:"view"`
	Result *governance.Result           `json:"result"`
	Error  *governance.Error            `json:"error"`
}

func TestCompiledCoordinationWatchWithIndependentWork(t *testing.T) {
	s, sel, q := fixture(t)
	bind(t, s, &q)
	first, e := s.Start(governance.LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	change := func(r governance.Result, action string) governance.Result {
		t.Helper()
		next, e := s.Transition(governance.LocalController(), governance.TransitionRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: governancehost.UUID(), RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: action, Reason: "independent fixture work"})
		if e != nil {
			t.Fatal(e)
		}
		return next
	}
	first = change(first, "start")
	q.OperationID = governancehost.UUID()
	q.RequestID = governancehost.UUID()
	q.Assignment.WritablePaths = []string{"second.txt"}
	second, e := s.Start(governance.LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	second = change(second, "start")
	second, e = s.RequestAttention(governance.LocalController(), governance.AttentionRequest{Schema: governance.AttentionSchema, ProjectRef: second.ProjectRef, CheckoutID: second.CheckoutID, OperationID: governancehost.UUID(), RunID: second.Run.ID, ExpectedRevision: second.Run.Revision, InputRevision: second.Run.InputRevision, Prompt: "Which fixture label?", Blocking: true})
	if e != nil {
		t.Fatal(e)
	}
	coverage, e := s.Context(sel)
	if e != nil {
		t.Fatal(e)
	}
	state := filepath.Join(filepath.Dir(coverage.Root), "state")
	untouched, e := os.ReadFile(filepath.Join(coverage.Root, "unrelated.txt"))
	if e != nil {
		t.Fatal(e)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	binary := filepath.Join(t.TempDir(), "sdp-controller")
	if b, e := exec.CommandContext(ctx, filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdp-controller").CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	}
	cmd := exec.CommandContext(ctx, binary, "coordinate", "--project", coverage.Root, "--state", state, "--watch", "--interval", "100ms", "--timeout", "20s", "--json", "--coordinator", "fixture-leader", "--owner", "fixture-owner")
	input, e := cmd.StdinPipe()
	if e != nil {
		t.Fatal(e)
	}
	output, e := cmd.StdoutPipe()
	if e != nil {
		t.Fatal(e)
	}
	cmd.Stderr = os.Stderr
	if e = cmd.Start(); e != nil {
		t.Fatal(e)
	}
	defer func() {
		input.Close()
		if cmd.ProcessState == nil {
			cmd.Process.Kill()
			cmd.Wait()
		}
	}()
	events := make(chan consoleEvent, 32)
	scanErrors := make(chan error, 1)
	go func() {
		defer close(events)
		scan := bufio.NewScanner(output)
		scan.Buffer(make([]byte, 4096), 1<<20)
		for scan.Scan() {
			var v consoleEvent
			if e := json.Unmarshal(scan.Bytes(), &v); e != nil {
				scanErrors <- e
				return
			}
			v.Raw = append(json.RawMessage(nil), scan.Bytes()...)
			select {
			case events <- v:
			case <-ctx.Done():
				return
			}
		}
		scanErrors <- scan.Err()
	}()
	captured := []json.RawMessage{}
	receive := func(kind string) consoleEvent {
		t.Helper()
		for {
			select {
			case v, ok := <-events:
				if !ok {
					t.Fatal("console ended", <-scanErrors)
				}
				captured = append(captured, v.Raw)
				if v.Kind == kind {
					return v
				}
				if v.Kind == "error" {
					t.Fatal(v.Error)
				}
			case <-ctx.Done():
				t.Fatal("console timeout")
			}
		}
	}
	initial := receive("snapshot")
	if initial.View.Total != 3 || initial.View.Items[0].Kind != "attention" {
		t.Fatal(initial)
	}
	// Bad input produces an error without killing the observer.
	if _, e = io.WriteString(input, "{broken}\n"); e != nil {
		t.Fatal(e)
	}
	if v := receive("error"); v.Error == nil {
		t.Fatal(v)
	}
	question := second.Run.Attention[0]
	disposition := governance.AttentionDisposition{Schema: governance.AttentionSchema, ProjectRef: second.ProjectRef, CheckoutID: second.CheckoutID, OperationID: governancehost.UUID(), RunID: second.Run.ID, QuestionID: question.ID, Basis: question.Basis}
	send := func(v any) {
		t.Helper()
		if e := json.NewEncoder(input).Encode(v); e != nil {
			t.Fatal(e)
		}
	}
	send(map[string]any{"action": "reply", "disposition": disposition, "reply": map[string]any{"text": "Use short label"}})
	replied := receive("result")
	if replied.Result.Run.Attention[0].Reply.Actor != "fixture-owner" || replied.Result.Run.State != "running" {
		t.Fatal(replied)
	}
	receive("snapshot")
	// Worker mutation succeeds while the console waits on input; no global lease.
	first = change(first, "report")
	view := receive("snapshot").View
	var report *governance.ReportMessage
	for _, item := range view.Items {
		if item.Report != nil && item.RunID == first.Run.ID {
			report = item.Report
		}
	}
	if report == nil {
		t.Fatal("report lost during owner dialogue")
	}
	ack := governance.ReportAcknowledgementRequest{Schema: governance.ReportAcknowledgementSchema, ProjectRef: view.ProjectRef, CheckoutID: view.CheckoutID, OperationID: governancehost.UUID(), RunID: report.RunID, ReportID: report.ID, ReportedRevision: report.ReportedRevision, InputRevision: report.InputRevision, ExpectedRevision: report.CurrentRevision, ExpectedGeneration: view.Generation, Coordinator: "fixture-leader", Note: "Fixture report processed, awaiting verification."}
	send(map[string]any{"action": "acknowledge", "acknowledgement": ack})
	received := receive("result")
	if received.Result.Run.State != "verifying" || received.Result.Run.Revision != first.Run.Revision {
		t.Fatal(received)
	}
	receive("snapshot")
	send(map[string]any{"action": "quit"})
	if e = cmd.Wait(); e != nil {
		t.Fatal(e)
	}
	// A separate process recovers the durable report/reply view after observer exit.
	b, e := exec.CommandContext(ctx, binary, "coordinate", "--project", coverage.Root, "--state", state, "--json").CombinedOutput()
	if e != nil {
		t.Fatal(e, string(b))
	}
	var reopened governance.CoordinationView
	if e = json.Unmarshal(b, &reopened); e != nil {
		t.Fatal(e)
	}
	acknowledged, answered := false, false
	for _, item := range reopened.Items {
		if item.Report != nil {
			acknowledged = item.Report.Acknowledgement != nil
		}
		if item.Attention != nil {
			answered = item.Attention.Status == "answered"
		}
	}
	if !acknowledged || !answered {
		t.Fatal(reopened)
	}
	now, e := os.ReadFile(filepath.Join(coverage.Root, "unrelated.txt"))
	if e != nil || string(now) != string(untouched) {
		t.Fatal("unrelated source changed", e)
	}
	if path := os.Getenv("SDP_PGL3_C6_CONSOLE_WIRE"); path != "" {
		b, e = json.MarshalIndent(captured, "", "  ")
		if e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(path, b, 0600); e != nil {
			t.Fatal(e)
		}
	}
}
