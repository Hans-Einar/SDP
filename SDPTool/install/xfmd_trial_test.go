package install

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestXFMDDisposableAdoption never writes to the selected source. Product commands
// operate exclusively on a temporary copy of the observed installation/link scope.
func TestXFMDDisposableAdoption(t *testing.T) {
	source := os.Getenv("SDP_TEST_XFMD")
	binary := os.Getenv("SDP_TEST_BINARY")
	client := os.Getenv("GH_SDP_BINARY")
	artifactPath := os.Getenv("SDP_TEST_ARTIFACT")
	if source == "" || binary == "" || client == "" || artifactPath == "" {
		t.Skip("set SDP_TEST_XFMD, SDP_TEST_BINARY, GH_SDP_BINARY and SDP_TEST_ARTIFACT")
	}
	data, e := Read(artifactPath, MetadataLimit)
	if e != nil {
		t.Fatal(e)
	}
	var d Descriptor
	if e = Decode(data, MetadataLimit, &d); e != nil {
		t.Fatal(e)
	}
	if e = ValidateDescriptor(d); e != nil {
		t.Fatal(e)
	}
	before, e := Inspect(source, targetPaths(d))
	if e != nil {
		t.Fatal(e)
	}
	git := func(args ...string) []byte {
		cmd := exec.Command("git", append([]string{"-C", source}, args...)...)
		out, e := cmd.Output()
		if e != nil {
			t.Fatal(e)
		}
		return out
	}
	status := git("status", "--porcelain=v1", "-z")
	commit := strings.TrimSpace(string(git("rev-parse", "HEAD")))
	work := t.TempDir()
	root := filepath.Join(work, "project")
	descriptorPath := filepath.Join(work, "release.json")
	keyPath := filepath.Join(work, "test.pub")
	toolPath := filepath.Join(work, "sdptool")
	bb, e := Read(binary, PayloadLimit)
	if e != nil {
		t.Fatal(e)
	}
	if len(d.Binaries) != 1 || d.Binaries[0].SHA256 != Hash(bb) {
		t.Fatal("descriptor does not pin actual candidate")
	}
	if e = os.WriteFile(toolPath, bb, 0700); e != nil {
		t.Fatal(e)
	}
	os.WriteFile(descriptorPath, data, 0600)
	public, private, _ := ed25519.GenerateKey(nil)
	signature, _ := json.Marshal(map[string]any{"keyId": Hash(public), "signature": ed25519.Sign(private, data)})
	os.WriteFile(descriptorPath+".sig", signature, 0600)
	os.WriteFile(keyPath, []byte(base64.StdEncoding.EncodeToString(public)), 0600)
	t.Setenv("SDP_RELEASE", descriptorPath)
	t.Setenv("SDP_TEST_KEY", keyPath)
	t.Setenv("SDP_CACHE_DIR", filepath.Join(work, "cache"))
	t.Setenv("SDP_OFFLINE", "false")
	t.Setenv("GH_CONFIG_DIR", filepath.Join(work, "gh-config"))
	t.Setenv("XDG_DATA_HOME", filepath.Join(work, "gh-data"))
	extension := filepath.Join(work, "gh-sdp")
	os.Mkdir(extension, 0700)
	cb, e := Read(client, PayloadLimit)
	if e != nil {
		t.Fatal(e)
	}
	os.WriteFile(filepath.Join(extension, "gh-sdp"), cb, 0700)
	cmd := exec.Command("gh", "extension", "install", ".")
	cmd.Dir = extension
	if out, e := cmd.CombinedOutput(); e != nil {
		t.Fatal(e, string(out))
	}
	restore := func() {
		t.Helper()
		os.RemoveAll(root)
		os.Mkdir(root, 0700)
		for p, o := range before.Snapshot {
			if o.Type == "directory" {
				os.MkdirAll(filepath.Join(root, p), 0700)
			} else if o.Type == "file" {
				b, e := Read(filepath.Join(source, p), FileLimit)
				if e != nil || Hash(b) != *o.SHA256 {
					t.Fatal("source changed during copy", p, e)
				}
				put(t, root, p, b)
			}
		}
	}
	restore()
	ad := Adoption{SchemaVersion: AdoptionSchema, ProjectRoot: root, Baseline: "manual", ObservedCommit: &commit, TargetDigest: Hash(data), Inventory: before.Snapshot, Moves: []Move{{"SDP/Agents/KanBan", "SDP/KanBan"}}, RefreshManaged: []string{"AGENTS.md"}, AllowReferenceWarnings: true}
	manifest := filepath.Join(work, "xfmd-upgrade.yaml")
	ab, _ := Canonical(ad)
	os.WriteFile(manifest, ab, 0600)
	run := func(viaGH bool, args ...string) []byte {
		t.Helper()
		name := binary
		if viaGH {
			name = "gh"
			args = append([]string{"sdp"}, args...)
		}
		cmd := exec.Command(name, args...)
		cmd.Dir = root
		out, e := cmd.CombinedOutput()
		if e != nil {
			t.Fatalf("%s %v: %v\n%s", name, args, e, out)
		}
		return out
	}
	directPlan := filepath.Join(work, "direct-plan.json")
	wrapperPlan := filepath.Join(work, "wrapper-plan.json")
	run(false, "upgrade", "--manifest", manifest, "--plan-output", directPlan, "--json")
	run(true, "upgrade", "--manifest", manifest, "--plan-output", wrapperPlan, "--json")
	p, e := LoadPlan(directPlan)
	if e != nil {
		t.Fatal(e)
	}
	q, e := LoadPlan(wrapperPlan)
	if e != nil || q.PlanDigest != p.PlanDigest {
		t.Fatal("direct/client plans differ", e)
	}
	baselineCopy, _ := Inspect(root, targetPaths(d))
	if !SameSnapshot(before.Snapshot, baselineCopy.Snapshot) {
		t.Fatal("preview changed project")
	}
	historySource := "SDP/Agents/KanBan/Ledger.ndjson"
	history := before.Files[historySource]
	verify := func() {
		t.Helper()
		got, e := Read(filepath.Join(root, "AGENTS-project.md"), FileLimit)
		if e != nil || !bytes.Equal(got, bytes.ReplaceAll(before.Files["AGENTS.md"], []byte("[KanBan board](SDP/Agents/KanBan/README.md)"), []byte("[KanBan board](SDP/KanBan/README.md)"))) {
			t.Fatal("root instructions not preserved")
		}
		// The only accepted instruction edit is the independently inspected board
		// link relocation above. The original bytes must also survive in the backup.
		foundBackup := false
		for i, a := range p.Actions {
			if a.Path == "AGENTS.md" {
				original, e := Read(filepath.Join(root, Operations, "install-"+p.PlanDigest[:24], "backups", fmt.Sprint(i)), FileLimit)
				if e != nil || !bytes.Equal(original, before.Files["AGENTS.md"]) {
					t.Fatal("original instruction backup lost", e)
				}
				foundBackup = true
			}
		}
		if !foundBackup {
			t.Fatal("missing root-instruction replacement")
		}
		ledger, e := Read(filepath.Join(root, "SDP/ProjectManagement/Ledger.ndjson"), FileLimit)
		if e != nil || !bytes.HasPrefix(ledger, history) {
			t.Fatal("history prefix changed")
		}
		if _, e = os.Stat(filepath.Join(root, historySource)); !os.IsNotExist(e) {
			t.Fatal("old writable history remains")
		}
		for _, a := range p.Actions {
			if a.Action == "write" && a.Path != "SDP/ProjectManagement/Ledger.ndjson" {
				b, e := Read(filepath.Join(root, a.Path), FileLimit)
				if e != nil || Hash(b) != *a.After {
					t.Fatal("planned output mismatch", a.Path, e)
				}
			}
		}
		for _, name := range p.Preserved {
			b, e := Read(filepath.Join(root, name), FileLimit)
			if e != nil || Hash(b) != *p.Snapshot[name].SHA256 {
				t.Fatal("preserved path changed", name, e)
			}
		}
		b, e := Read(filepath.Join(root, ReceiptPath), MetadataLimit)
		if e != nil {
			t.Fatal(e)
		}
		var receipt Receipt
		if e = Decode(b, MetadataLimit, &receipt); e != nil || receipt.Provenance != "test-signed" {
			t.Fatal("receipt provenance", e)
		}
		pending, e := Pending(root)
		if e != nil || len(pending) != 0 {
			t.Fatal(pending, e)
		}
	}
	var result Result
	for _, viaGH := range []bool{false, true} {
		if viaGH {
			restore()
		}
		out := run(viaGH, "upgrade", "--apply", directPlan, "--json")
		var envelope struct {
			Result Result `json:"result"`
		}
		if e = json.Unmarshal(out, &envelope); e != nil {
			t.Fatal(e)
		}
		result = envelope.Result
		verify()
		snap, _ := Inspect(root, nil)
		repeat := filepath.Join(work, fmt.Sprintf("repeat-%t.json", viaGH))
		run(viaGH, "upgrade", "--plan-output", repeat, "--json")
		r, e := LoadPlan(repeat)
		if e != nil || !r.NoChange {
			t.Fatal("repeat not no-change", e)
		}
		run(viaGH, "upgrade", "--apply", repeat, "--json")
		after, _ := Inspect(root, nil)
		if !SameSnapshot(snap.Snapshot, after.Snapshot) {
			t.Fatal("repeat changed bytes")
		}
	}
	// Exercise real process-exit recovery on the full copied project, then resume
	// using the actual gh sdp entrypoint. Fault injection exists only in test code.
	restore()
	cmd = exec.Command(os.Args[0], "-test.run=^TestCrashHelper$")
	cmd.Env = append(os.Environ(), "SDP_TEST_CRASH_CHILD=1", "SDP_TEST_PLAN="+directPlan, "SDP_TEST_BOUNDARY=write", "SDP_TEST_INDEX=5")
	out, err := cmd.CombinedOutput()
	ee, ok := err.(*exec.ExitError)
	if !ok || ee.ExitCode() != 73 {
		t.Fatal(err, string(out))
	}
	pending, e := Pending(root)
	if e != nil || len(pending) != 1 {
		t.Fatal(pending, e)
	}
	run(true, "upgrade", "--resume", pending[0], "--json")
	verify()
	after, e := Inspect(source, targetPaths(d))
	if e != nil || !SameSnapshot(before.Snapshot, after.Snapshot) || !bytes.Equal(status, git("status", "--porcelain=v1", "-z")) {
		t.Fatal("live source observations/status changed", e)
	}
	if dest := os.Getenv("SDP_TEST_EVIDENCE"); dest != "" {
		proof := map[string]any{"schemaVersion": "gip-xfmd-trial/1", "sourceRoot": source, "sourceCommit": commit, "sourceStatus": string(status), "liveSnapshot": before.Snapshot, "liveUnchanged": true, "engineSHA256": Hash(bb), "clientSHA256": Hash(cb), "descriptorSHA256": Hash(data), "sourceEngineCommit": d.SourceCommit, "adoptionManifest": ad, "planDigest": p.PlanDigest, "actions": len(p.Actions), "preserved": len(p.Preserved), "historyPrefixSHA256": Hash(history), "result": result, "directAndGhPlansEqual": true, "directAndGhApplyPassed": true, "processExitAndGhResumePassed": true, "repeatNoChange": true, "limitations": []string{"Linux native evidence only", "ephemeral test signing key, no production release", "non-Markdown references and excluded build/dependency trees require separate disposition", "no live mutation"}}
		actionProof := []map[string]any{}
		for _, a := range p.Actions {
			record := map[string]any{"action": a.Action, "path": a.Path, "before": a.Before, "after": a.After}
			if a.Path == "SDP/KanBan/board.json" || a.Path == "SDP/navigation.json" {
				var value any
				if e = json.Unmarshal(a.Content, &value); e != nil {
					t.Fatal(e)
				}
				record["jsonValue"] = value
			}
			actionProof = append(actionProof, record)
		}
		proof["plannedActions"] = actionProof
		proof["instructionPreservation"] = "Only the KanBan board link is rebased in AGENTS-project.md; original AGENTS.md bytes verified in the operation backup."
		b, _ := json.MarshalIndent(proof, "", "  ")
		if e = os.WriteFile(dest, append(b, '\n'), 0600); e != nil {
			t.Fatal(e)
		}
	}
}
