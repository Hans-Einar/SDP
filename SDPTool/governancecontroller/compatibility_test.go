package governancecontroller

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

// The executable below is this test's recorded protocol host, not Codex/model execution.
func TestFutureHostPersistsActualVersionAndRejectsCrossVersionRecovery(t *testing.T) {
	t.Setenv("SDP_PGI3_FAKE_HOST", "1")
	t.Setenv("SDP_PGC1_FAKE_HOST_VERSION", "0.190.3")
	s, sel, q := fixture(t)
	bind(t, s, &q)
	r, err := s.Start(governance.LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	coverage, err := s.Context(sel)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	client, err := governancehost.Open(ctx, exe, coverage.Root)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	c := &Controller{Service: s, Host: client, Root: coverage.Root, StateDir: filepath.Join(filepath.Dir(coverage.Root), "state"), MCPBinary: "/fixture/mcp", Epoch: client.Epoch, Sandbox: "read-only"}
	r, err = c.Launch(ctx, r.Run.ID, "Fix the selected fixture")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Run.Attempts) != 1 || r.Run.Attempts[0].HostVersion != "0.190.3" {
		t.Fatal("actual host version lost", r)
	}
	// A new supported version must not silently rebind an existing attempt.
	t.Setenv("SDP_PGC1_FAKE_HOST_VERSION", "0.191.0")
	other, err := governancehost.Open(ctx, exe, coverage.Root)
	if err != nil {
		t.Fatal(err)
	}
	defer other.Close()
	_, err = other.PrepareRecovery(ctx, s, r.Run.ID, r.Run.Attempts[0].ID)
	if err == nil || !strings.Contains(err.Error(), "recorded attempt host version") {
		t.Fatal(err)
	}
	after, err := s.Read(r.Run.ID)
	if err != nil || after.Run.Revision != r.Run.Revision {
		t.Fatal("refusal mutated the run", after, err)
	}
}
