package governance

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestSessionHostSelectionAndCompiledStatus(t *testing.T) {
	s, q := workFixture(t, "", "")
	binary := filepath.Join(t.TempDir(), "sdp-controller")
	if b, e := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdp-controller").CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	}
	call := func(action string, extra ...string) ([]byte, error) {
		args := []string{action, "--project", s.root, "--state", s.state, "--instance", q.InstanceID}
		return exec.Command(binary, append(args, extra...)...).CombinedOutput()
	}
	if _, e := s.SessionHostSelection(q.InstanceID); e == nil {
		t.Fatal("uncommitted assignment selected")
	}
	b, e := call("status")
	if e != nil || !json.Valid(b) {
		t.Fatal(e, string(b))
	}
	head := mustFile(t, filepath.Join(s.state, "HEAD"))
	if b, e = call("run", "--codex", "/missing/native-host", "--mcp-binary", "/missing/mcp"); e == nil || !bytes.Contains(b, []byte("no current committed assignment")) {
		t.Fatal(e, string(b))
	}
	if !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
		t.Fatal("inspection/refusal wrote state")
	}
	p, e := s.StartSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.SessionHostSelection(q.InstanceID); e == nil {
		t.Fatal("prepared start selected")
	}
	done, e := s.ResumeSession(LocalController(), workResume(q, p))
	if e != nil {
		t.Fatal(e)
	}
	selected, e := s.SessionHostSelection(q.InstanceID)
	if e != nil || selected.RunID != q.OperationID || selected.Revision != done.Revision || selected.Outcome != q.Outcome || !reflect.DeepEqual(selected.Skills, q.Skills) {
		t.Fatal(selected, e)
	}
	b, e = call("status")
	if e != nil {
		t.Fatal(e, string(b))
	}
	decoder := json.NewDecoder(bytes.NewReader(b))
	var v SessionView
	var r Result
	if decoder.Decode(&v) != nil || decoder.Decode(&r) != nil || v.InstanceID != q.InstanceID || r.Run.ID != q.OperationID {
		t.Fatal(string(b))
	}
	var extra any
	if decoder.Decode(&extra) != io.EOF {
		t.Fatal("unexpected status object")
	}
	if b, e = call("status", "--run", q.OperationID); e == nil || !bytes.Contains(b, []byte("exactly one")) {
		t.Fatal(e, string(b))
	}
	evidenceFile(t, s, "SDP/authority.md", "Changed authority; do not dispatch.\n")
	head = mustFile(t, filepath.Join(s.state, "HEAD"))
	if b, e = call("run", "--codex", "/missing/native-host", "--mcp-binary", "/missing/mcp"); e == nil || !bytes.Contains(b, []byte("stale_revision")) {
		t.Fatal(e, string(b))
	}
	if !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
		t.Fatal("stale launch changed state")
	}
}

// Opt-in setup only. Native submission is a separate explicit controller command;
// ordinary tests never submit model work or read account configuration.
func TestSessionHostNativeFixture(t *testing.T) {
	root, state := os.Getenv("PG_SESSION_HOST_LAB_ROOT"), os.Getenv("PG_SESSION_HOST_LAB_STATE")
	if root == "" {
		t.Skip("explicit native lab setup only")
	}
	if state == "" {
		t.Fatal("private state required")
	}
	source, e := filepath.Abs("../../Skills")
	if e != nil {
		t.Fatal(e)
	}
	s, q := workConfigured(t, root, state, func(t *testing.T, s *Service, b *SessionBinding) {
		evidenceFile(t, s, "AGENTS.md", "# Authorized native Session fixture\n\nThis isolated pg_lab task is authorized by its owner for PG development testing.\nUse the supplied existing Session/plan/card and pinned Master/Worker skills.\nChange only target.txt as the committed assignment requests; inspect its exact bytes\nto verify the result. This is a no-design wording change. Do not create another\nSession or plan, publish records, edit skills or configuration, spawn/delegate, use\nthe network, install packages, commit or push. The controller owns canonical\nSession/KanBan effects. Report through assignment-bound MCP and return.\nThe minimal manifest and fixture authority govern this trial; no full Toolkit\ninstallation or release claim is made. Preserve every other file.\n")
		if e := filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(source, path)
			if err != nil {
				return err
			}
			dest := filepath.Join(s.root, "Skills", rel)
			if entry.IsDir() {
				return os.MkdirAll(dest, 0700)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			return os.WriteFile(dest, data, 0600)
		}); e != nil {
			t.Fatal(e)
		}
		if e := os.MkdirAll(filepath.Join(s.root, ".agents/skills"), 0700); e != nil {
			t.Fatal(e)
		}
		for _, name := range []string{"sdp", "sdp-master", "sdp-worker"} {
			if e := os.Symlink(filepath.Join("../../Skills", name), filepath.Join(s.root, ".agents/skills", name)); e != nil {
				t.Fatal(e)
			}
		}
	}, nil)
	// The generic service fixture creates synthetic launch skills. Replace only those
	// uncommitted inputs with canonical copies before preparing this native assignment.
	for n, sk := range q.Skills {
		data := mustFile(t, filepath.Join(source, sk.Name, "SKILL.md"))
		if e = os.WriteFile(sk.Path, data, 0600); e != nil {
			t.Fatal(e)
		}
		q.Skills[n].Digest = hash(data)
	}
	q.Outcome = "Replace the entire contents of target.txt with exactly PG Session native execution verified. followed by one newline. Read the resulting file to verify exact bytes; preserve all other files. This is a wording-only correction with no design change. Report the outcome through assignment-bound MCP and return."
	q.Reason = "Start the owner-authorized native Session entry-point trial"
	_, r := startWork(t, s, q)
	selected, e := s.SessionHostSelection(q.InstanceID)
	if e != nil {
		t.Fatal(e)
	}
	if !strings.HasPrefix(selected.Outcome, "Replace") || r.Run.State != "ready" {
		t.Fatal(selected)
	}
	out, _ := json.MarshalIndent(selected, "", "  ")
	if e = os.WriteFile(filepath.Join(root, "evidence/native-launch.json"), out, 0600); e != nil {
		t.Fatal(e)
	}
}
