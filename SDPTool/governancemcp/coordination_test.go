package governancemcp

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestCompiledCoordinatorMCPAuthority(t *testing.T) {
	s, sel, q := fixture(t)
	bind(t, s, &q)
	r, e := s.Start(governance.LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	for _, action := range []string{"start", "report"} {
		r, e = s.Transition(governance.LocalController(), governance.TransitionRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: governancehost.UUID(), RunID: r.Run.ID, InputRevision: r.Run.InputRevision, ExpectedRevision: r.Run.Revision, Action: action, Reason: "protocol fixture"})
		if e != nil {
			t.Fatal(e)
		}
	}
	c, e := s.Context(sel)
	if e != nil {
		t.Fatal(e)
	}
	state := filepath.Join(filepath.Dir(c.Root), "state")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	binary := filepath.Join(t.TempDir(), "sdp-mcp")
	if b, e := exec.CommandContext(ctx, filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdp-mcp").CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	}
	connectStdio := func() *mcp.ClientSession {
		t.Helper()
		transport := &mcp.CommandTransport{Command: exec.CommandContext(ctx, binary, "--project", c.Root, "--state", state, "--coordinator", "fixture-leader")}
		client, e := mcp.NewClient(&mcp.Implementation{Name: "coordinator-test", Version: "1"}, nil).Connect(ctx, transport, nil)
		if e != nil {
			t.Fatal(e)
		}
		t.Cleanup(func() { client.Close() })
		return client
	}
	client := connectStdio()
	list, e := client.ListTools(ctx, nil)
	if e != nil || len(list.Tools) != 2 {
		t.Fatal(list, e)
	}
	tools := map[string]bool{}
	for _, tool := range list.Tools {
		tools[tool.Name] = true
		if (tool.Name == "coordination_read") != tool.Annotations.ReadOnlyHint {
			t.Fatal(tool)
		}
	}
	if !tools["coordination_read"] || !tools["report_acknowledge"] {
		t.Fatal(tools)
	}
	read := call(t, client, "coordination_read", governance.CoordinationRequest{}, false)
	b, _ := json.Marshal(read.StructuredContent)
	var view governance.CoordinationView
	if e = json.Unmarshal(b, &view); e != nil || view.Total != 2 {
		t.Fatal(view, e)
	}
	report := view.Items[0].Report
	ack := governance.ReportAcknowledgementRequest{Schema: governance.ReportAcknowledgementSchema, ProjectRef: view.ProjectRef, CheckoutID: view.CheckoutID, OperationID: governancehost.UUID(), RunID: report.RunID, ReportID: report.ID, ReportedRevision: report.ReportedRevision, InputRevision: report.InputRevision, ExpectedRevision: report.CurrentRevision, ExpectedGeneration: view.Generation, Coordinator: "forged-owner", Note: "Read; evidence verification remains."}
	call(t, client, "report_acknowledge", ack, true)
	ack.Coordinator = "fixture-leader"
	result := call(t, client, "report_acknowledge", ack, false)
	client.Close()
	client = connectStdio()
	retry := call(t, client, "report_acknowledge", ack, false)
	if !reflect.DeepEqual(result.StructuredContent, retry.StructuredContent) {
		t.Fatal("restart lost exact receipt")
	}
	for _, name := range []string{"attention_reply", "run_start", "run_transition", "design_decision"} {
		if _, e = client.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: map[string]any{}}); e == nil {
			t.Fatal("privilege exposed", name)
		}
	}
	call(t, client, "coordination_read", map[string]any{"role": "owner"}, true)
	current, e := s.Read(r.Run.ID)
	if e != nil || !reflect.DeepEqual(current.Run, r.Run) {
		t.Fatal("ack advanced work", current, e)
	}
	if b, e = exec.CommandContext(ctx, binary, "--project", c.Root, "--state", state, "--coordinator", "fixture-leader", "--run", r.Run.ID).CombinedOutput(); e == nil {
		t.Fatal("mixed launcher scope accepted", string(b))
	}
	if _, e = NewCoordinator(s, " "); e == nil {
		t.Fatal("empty identity accepted")
	}
	ordinary, e := New(s, "", "", "")
	if e != nil {
		t.Fatal(e)
	}
	defaultClient := connect(t, ordinary)
	if _, e = defaultClient.CallTool(ctx, &mcp.CallToolParams{Name: "report_acknowledge", Arguments: ack}); e == nil {
		t.Fatal("default adapter gained authority")
	}
	if path := os.Getenv("SDP_PGL3_C6_MCP_WIRE"); path != "" {
		wire := map[string]any{"view": view, "acknowledgement": ack, "result": result.StructuredContent, "tools": list}
		b, e = json.MarshalIndent(wire, "", "  ")
		if e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(path, b, 0600); e != nil {
			t.Fatal(e)
		}
	}
}
