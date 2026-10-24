package governancecontroller

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

func TestCompiledEvidenceCaptureAttachInspectAndRework(t *testing.T) {
	s, sel, q := fixture(t)
	bind(t, s, &q)
	r, e := s.Start(governance.LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	for _, action := range []string{"start", "report"} {
		r, e = s.Transition(governance.LocalController(), governance.TransitionRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID,
			OperationID: governancehost.UUID(), RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: action, Reason: "compiled evidence fixture"})
		if e != nil {
			t.Fatal(e)
		}
	}
	c, e := s.Context(sel)
	if e != nil {
		t.Fatal(e)
	}
	state := filepath.Join(filepath.Dir(c.Root), "state")
	artifact := []byte("external verification found a defect\n")
	if e = os.WriteFile(filepath.Join(c.Root, "verification.log"), artifact, 0600); e != nil {
		t.Fatal(e)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	binary := filepath.Join(t.TempDir(), "sdp-controller")
	if b, err := exec.CommandContext(ctx, filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdp-controller").CombinedOutput(); err != nil {
		t.Fatal(err, string(b))
	}
	invoke := func(args ...string) ([]byte, error) {
		args = append(append([]string{"evidence"}, args...), "--project", c.Root, "--state", state)
		return exec.CommandContext(ctx, binary, args...).CombinedOutput()
	}
	raw, e := invoke("capture", "--run", r.Run.ID)
	if e != nil {
		t.Fatal(e, string(raw))
	}
	var candidate governance.Candidate
	if e = json.Unmarshal(raw, &candidate); e != nil {
		t.Fatal(e)
	}
	request := governance.EvidenceRequest{Schema: governance.EvidenceSchema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: governancehost.UUID(), RunID: r.Run.ID,
		ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, ReportID: r.OperationID, Candidate: candidate, Kind: "test", Verdict: "failed",
		Issuer: "compiled fixture", Note: "external test result", Artifact: governance.Source{Path: "verification.log", SHA256: fmt.Sprintf("%x", sha256.Sum256(artifact))}}
	input := filepath.Join(t.TempDir(), "input.json")
	write := func(v any) {
		t.Helper()
		b, _ := json.Marshal(v)
		if err := os.WriteFile(input, b, 0600); err != nil {
			t.Fatal(err)
		}
	}
	write(request)
	attached, e := invoke("attach", "--input", input)
	if e != nil {
		t.Fatal(e, string(attached))
	}
	repeated, e := invoke("attach", "--input", input)
	if e != nil || !bytes.Equal(attached, repeated) {
		t.Fatal("retry did not preserve receipt", e)
	}
	raw, e = invoke("inspect", "--run", r.Run.ID, "--json")
	if e != nil {
		t.Fatal(e, string(raw))
	}
	var view governance.EvidenceView
	if e = json.Unmarshal(raw, &view); e != nil || len(view.Items) != 1 || !view.Items[0].Applicable {
		t.Fatal(view, e)
	}
	human, e := invoke("inspect", "--run", r.Run.ID)
	if e != nil || !strings.Contains(string(human), "review independence unverified") {
		t.Fatal(e, string(human))
	}
	if old := os.Getenv("PG_EVIDENCE_OLD_CONTROLLER"); old != "" {
		b, err := exec.CommandContext(ctx, old, "inbox", "--project", c.Root, "--state", state, "--json").CombinedOutput()
		if err == nil || (!strings.Contains(string(b), "incompatible_profile") && !strings.Contains(string(b), "reconciliation_required")) {
			t.Fatal("old reader accepted new state", err, string(b))
		}
	}
	rework := governance.EvidenceReworkRequest{Schema: governance.EvidenceSchema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID,
		OperationID: governancehost.UUID(), ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, EvidenceID: request.OperationID, Reason: "fix the failed case"}
	write(rework)
	raw, e = invoke("rework", "--input", input)
	if e != nil {
		t.Fatal(e, string(raw))
	}
	var result governance.Result
	if e = json.Unmarshal(raw, &result); e != nil || result.Run.State != "ready" || result.Run.Selection != r.Run.Selection {
		t.Fatal(result, e)
	}
	final, e := s.Evidence(r.Run.ID)
	if e != nil || final.Items[0].Applicable {
		t.Fatal(final, e)
	}
	if path := os.Getenv("PG_EVIDENCE_CLI_WIRE"); path != "" {
		cases := []map[string]any{}
		var av any
		if e = json.Unmarshal(attached, &av); e != nil {
			t.Fatal(e)
		}
		for _, v := range []struct {
			name  string
			value any
		}{{"Candidate", candidate}, {"EvidenceRequest", request}, {"Result", av}, {"EvidenceView", view}, {"EvidenceReworkRequest", rework}, {"Result", result}, {"EvidenceView", final}} {
			cases = append(cases, map[string]any{"definition": v.name, "value": v.value})
		}
		b, _ := json.MarshalIndent(cases, "", "  ")
		if e = os.WriteFile(path, b, 0600); e != nil {
			t.Fatal(e)
		}
	}
}
