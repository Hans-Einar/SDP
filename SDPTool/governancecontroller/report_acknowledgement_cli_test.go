package governancecontroller

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

func TestCompiledReportAcknowledgement(t *testing.T) {
	s, sel, q := fixture(t)
	bind(t, s, &q)
	r, e := s.Start(governance.LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	change := func(action string) {
		t.Helper()
		r, e = s.Transition(governance.LocalController(), governance.TransitionRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, InputRevision: r.Run.InputRevision, OperationID: governancehost.UUID(), ExpectedRevision: r.Run.Revision, Action: action, Reason: "compiled report fixture"})
		if e != nil {
			t.Fatal(e)
		}
	}
	change("start")
	change("report")
	c, e := s.Context(sel)
	if e != nil {
		t.Fatal(e)
	}
	state := filepath.Join(filepath.Dir(c.Root), "state")
	before := map[string][]byte{}
	if e = filepath.WalkDir(c.Root, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if !d.IsDir() {
			b, e := os.ReadFile(p)
			if e != nil {
				return e
			}
			before[p] = b
		}
		return nil
	}); e != nil {
		t.Fatal(e)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	binary := filepath.Join(t.TempDir(), "sdp-controller")
	if b, e := exec.CommandContext(ctx, filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdp-controller").CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	}
	invoke := func(args ...string) ([]byte, error) {
		args = append(args, "--project", c.Root, "--state", state)
		return exec.CommandContext(ctx, binary, args...).CombinedOutput()
	}
	raw, e := invoke("inbox", "--json")
	if e != nil {
		t.Fatal(e, string(raw))
	}
	var view governance.ReportInbox
	if e = json.Unmarshal(raw, &view); e != nil || len(view.Reports) != 1 {
		t.Fatal(view, e)
	}
	report := view.Reports[0]
	ack := governance.ReportAcknowledgementRequest{Schema: governance.ReportAcknowledgementSchema, ProjectRef: view.ProjectRef, CheckoutID: view.CheckoutID, OperationID: governancehost.UUID(), RunID: report.RunID, ReportID: report.ID, ReportedRevision: report.ReportedRevision, InputRevision: report.InputRevision, ExpectedRevision: report.CurrentRevision, ExpectedGeneration: view.Generation, Coordinator: "fixture-leader", Note: "Read and queued for verification; not accepted."}
	path := filepath.Join(t.TempDir(), "ack.json")
	data, _ := json.Marshal(ack)
	if e = os.WriteFile(path, data, 0600); e != nil {
		t.Fatal(e)
	}
	result, e := invoke("inbox", "acknowledge", "--input", path)
	if e != nil {
		t.Fatal(e, string(result))
	}
	var out governance.Result
	if e = json.Unmarshal(result, &out); e != nil || !reflect.DeepEqual(out.Run, r.Run) || out.ReportAcknowledgement == nil {
		t.Fatal(out, e)
	}
	head, e := os.ReadFile(filepath.Join(state, "HEAD"))
	if e != nil {
		t.Fatal(e)
	}
	retry, e := invoke("inbox", "acknowledge", "--input", path)
	if e != nil || !bytes.Equal(retry, result) {
		t.Fatal("lost-response retry", string(retry), e)
	}
	final, e := invoke("inbox", "--json")
	if e != nil {
		t.Fatal(e, string(final))
	}
	if e = json.Unmarshal(final, &view); e != nil || view.Schema != "sdp-governance-inbox/2" || view.Reports[0].AcknowledgementID != ack.OperationID {
		t.Fatal(view, e)
	}
	human, e := invoke("inbox")
	if e != nil || !strings.Contains(string(human), "fixture-leader") || !strings.Contains(string(human), "verifying") {
		t.Fatal(e, string(human))
	}
	// Actual old binary, when supplied by the evidence runner, must refuse new state.
	if old := os.Getenv("SDP_PGL3_C5_OLD_READER"); old != "" {
		b, e := exec.CommandContext(ctx, old, "inbox", "--project", c.Root, "--state", state, "--json").CombinedOutput()
		if e == nil || !(strings.Contains(string(b), "reconciliation_required") || strings.Contains(string(b), "incompatible_profile")) {
			t.Fatal("old reader accepted new state", e, string(b))
		}
	}
	after, _ := os.ReadFile(filepath.Join(state, "HEAD"))
	if !bytes.Equal(after, head) {
		t.Fatal("reads/retry changed generation")
	}
	for p, b := range before {
		now, e := os.ReadFile(p)
		if e != nil || !bytes.Equal(now, b) {
			t.Fatal("project source changed", p, e)
		}
	}
	if path := os.Getenv("SDP_PGL3_C5_WIRE"); path != "" {
		wire := map[string]json.RawMessage{"request": data, "result": result, "inboxV1": raw, "inboxV2": final}
		b, e := json.MarshalIndent(wire, "", "  ")
		if e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(path, b, 0600); e != nil {
			t.Fatal(e)
		}
	}
}
