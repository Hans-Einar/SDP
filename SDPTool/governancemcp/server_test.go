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

func connect(t *testing.T, s *mcp.Server) *mcp.ClientSession {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	t.Cleanup(cancel)
	a, b := mcp.NewInMemoryTransports()
	ss, e := s.Connect(ctx, a, nil)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { ss.Close() })
	cs, e := mcp.NewClient(&mcp.Implementation{Name: "adapter-test", Version: "1"}, nil).Connect(ctx, b, nil)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { cs.Close() })
	if cs.InitializeResult().ProtocolVersion != ProtocolVersion {
		t.Fatal(cs.InitializeResult())
	}
	return cs
}
func call(t *testing.T, c *mcp.ClientSession, name string, args any, wantError bool) *mcp.CallToolResult {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, e := c.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: args})
	if e != nil {
		t.Fatal(e)
	}
	if r.IsError != wantError {
		t.Fatalf("%s: %#v", name, r)
	}
	return r
}
func TestMCPMatchesCoreAndRestrictsAuthority(t *testing.T) {
	s, sel, q := fixture(t)
	bind(t, s, &q)
	r, e := s.Start(governance.LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	at, epoch := governancehost.UUID(), governancehost.UUID()
	r, e = s.RecordHost(governance.LocalController(), governance.HostRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, OperationID: governancehost.UUID(), Action: "bind", AttemptID: at, HostVersion: governancehost.HostVersion, ThreadID: "fixture-thread", Observation: governance.Observation{Kind: "bound", Epoch: epoch, ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano), Status: "idle"}})
	if e != nil {
		t.Fatal(e)
	}
	server, e := New(s, r.Run.ID, at, epoch)
	if e != nil {
		t.Fatal(e)
	}
	client := connect(t, server)
	got := call(t, client, "context_read", sel, false)
	direct, e := s.Context(sel)
	if e != nil {
		t.Fatal(e)
	}
	b, _ := json.Marshal(direct)
	var want any
	json.Unmarshal(b, &want)
	if !reflect.DeepEqual(got.StructuredContent, want) {
		t.Fatal("context differs", got.StructuredContent, want)
	}
	bad := sel
	bad.Authority = "unrelated.txt"
	call(t, client, "context_read", bad, true)
	tr := governance.TransitionRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: governancehost.UUID(), RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: "start", Reason: "controller only"}
	call(t, client, "run_transition", tr, true)
	r, e = s.Transition(governance.LocalController(), tr)
	if e != nil {
		t.Fatal(e)
	}
	tr.Action = "report"
	tr.OperationID = governancehost.UUID()
	tr.ExpectedRevision = r.Run.Revision
	got = call(t, client, "run_transition", tr, false)
	again := call(t, client, "run_transition", tr, false)
	if !reflect.DeepEqual(got.StructuredContent, again.StructuredContent) {
		t.Fatal("retry changed receipt")
	}
	read := call(t, client, "run_read", map[string]any{}, false)
	data, _ := json.Marshal(read.StructuredContent)
	var result governance.Result
	json.Unmarshal(data, &result)
	if result.Run.State != "verifying" {
		t.Fatal(result)
	}
	tr.Action = "complete"
	tr.ExpectedRevision = result.Run.Revision
	tr.OperationID = governancehost.UUID()
	call(t, client, "run_transition", tr, true)
	for _, raw := range []string{`{"runId":"forged"}`, `{"role":"owner"}`, `{"role":"owner","role":"controller"}`, `{"Role":null}`} {
		call(t, client, "run_read", json.RawMessage(raw), true)
	}
	call(t, client, "context_read", json.RawMessage(`{"card":"a","card":"b"}`), true)
	if _, e = client.CallTool(context.Background(), &mcp.CallToolParams{Name: "run_start", Arguments: map[string]any{}}); e == nil {
		t.Fatal("controller operation exposed")
	}
}
func TestCompiledStdioNegotiation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	binary := filepath.Join(t.TempDir(), "sdp-mcp")
	cmd := exec.CommandContext(ctx, filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdp-mcp")
	if out, e := cmd.CombinedOutput(); e != nil {
		t.Fatal(e, string(out))
	}
	root, e := filepath.Abs("../..")
	if e != nil {
		t.Fatal(e)
	}
	state := filepath.Join(t.TempDir(), "state")
	transport := &mcp.CommandTransport{Command: exec.CommandContext(ctx, binary, "--project", root, "--state", state)}
	c, e := mcp.NewClient(&mcp.Implementation{Name: "stdio-probe", Version: "1"}, nil).Connect(ctx, transport, nil)
	if e != nil {
		t.Fatal(e)
	}
	defer c.Close()
	if c.InitializeResult().ProtocolVersion != ProtocolVersion {
		t.Fatal(c.InitializeResult())
	}
	list, e := c.ListTools(ctx, nil)
	if e != nil || len(list.Tools) != 3 {
		t.Fatal(list, e)
	}
	q := governance.RouteRequest{Schema: governance.Schema, Categories: []string{"status"}}
	r := call(t, c, "route_evaluate", q, false)
	b, _ := json.Marshal(r.StructuredContent)
	var route governance.Route
	json.Unmarshal(b, &route)
	if route.Disposition != "no-work" {
		t.Fatal(route)
	}
	if _, e = os.Stat(state); !os.IsNotExist(e) {
		t.Fatal("inspection created operational state", e)
	}
}
