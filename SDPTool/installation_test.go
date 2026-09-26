package sdptool

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessInstallationFacts(t *testing.T) {
	root, r := projectFixture(t)
	r.ProjectManifest = "SDP/SDP-project.manifest.yaml"
	saveRegistration(t, root, r)
	os.WriteFile(filepath.Join(root, r.ProjectManifest), []byte("schemaVersion: '1.0'\ninstalled:\n  manifestPath: installed.yaml\n"), 0600)
	facts := `schemaVersion: "2.0"
toolkitVersion: "0.2.0"
frameworkVersion: "1.0.0"
agentsContractVersion: "2.0.0"
installerVersion: "0.2.0"
toolkitInstalledAt: "2026-09-25T10:00:00Z"
sourceCommit: null
processProfile: "sdp-five-phase/0.1"
managementProfile: "sdp-project-management/0.1"
configurationDigest: "` + strings.Repeat("a", 64) + `"
skills:
  sdp: "1.0.0"
capabilities:
  - sdp.install.v2
`
	dest := filepath.Join(root, "SDP/installed.yaml")
	os.WriteFile(dest, []byte(facts), 0600)
	p, e := Discover(root)
	if e != nil || p.Installation["state"] != "declared" {
		t.Fatalf("%v %v", p, e)
	}
	for _, profile := range []string{"sdp-project-management/0.1", "sdp-project-management/0.2"} {
		os.WriteFile(dest, []byte(strings.ReplaceAll(facts, "sdp-project-management/0.1", profile)), 0600)
		if got, err := Discover(root); err != nil || got.Installation["state"] != "declared" {
			t.Fatalf("profile %s: %v %v", profile, got, err)
		}
	}
	os.WriteFile(dest, []byte(strings.ReplaceAll(facts, "sdp-project-management/0.1", "sdp-project-management/9.0")), 0600)
	if _, err := Discover(root); err == nil {
		t.Fatal("unknown management profile accepted")
	}
	for _, bad := range []string{facts + "unknown: true\n", strings.Replace(facts, strings.Repeat("a", 64), "broken", 1), strings.Replace(facts, "schemaVersion: \"2.0\"", "schemaVersion: \"3.0\"", 1)} {
		os.WriteFile(dest, []byte(bad), 0600)
		if _, e = Discover(root); e == nil {
			t.Fatal("invalid facts accepted")
		}
	}
	os.WriteFile(dest, []byte(facts), 0600)
	id := "install-" + strings.Repeat("a", 24)
	directory := filepath.Join(root, "SDP/.sdp-operations", id)
	os.MkdirAll(directory, 0700)
	file := filepath.Join(directory, "journal.json")
	for _, status := range []string{"active", "failed", "completed"} {
		b, _ := json.Marshal(map[string]any{"schemaVersion": "2.0", "operationId": id, "status": status})
		os.WriteFile(file, b, 0600)
		p, e = Discover(root)
		if e != nil {
			t.Fatal(e)
		}
		expected := "incomplete"
		if status == "completed" {
			expected = "declared"
		}
		if p.Installation["state"] != expected {
			t.Fatal(p.Installation)
		}
	}
	os.WriteFile(file, []byte(`{"schemaVersion":"2.0","operationId":"`+id+`","status":"active"}`), 0600)
	os.Remove(dest)
	p, e = Discover(root)
	if e != nil || p.Installation["state"] != "incomplete" {
		t.Fatalf("%v %v", p, e)
	}
}
func TestBuildVersionProtocol(t *testing.T) {
	var out, errs bytes.Buffer
	if Run(context.Background(), []string{"--version"}, &out, &errs) != 0 {
		t.Fatal(errs.String())
	}
	var result map[string]any
	if e := json.Unmarshal(out.Bytes(), &result); e != nil {
		t.Fatal(e)
	}
	if result["operation"] != "version" || result["version"] != BuildVersion {
		t.Fatal(result)
	}
	if len(result["installedFactSchemas"].([]any)) != 2 {
		t.Fatal(result)
	}
}

func TestEarlyIncompleteInstallation(t *testing.T) {
	root := t.TempDir()
	id := "install-" + strings.Repeat("b", 24)
	area := filepath.Join(root, "SDP")
	op := filepath.Join(area, ".sdp-operations", id)
	os.MkdirAll(op, 0700)
	os.WriteFile(filepath.Join(op, "journal.json"), []byte(`{"schemaVersion":"2.0","operationId":"`+id+`","status":"active"}`), 0600)
	for _, selected := range []string{root, area} {
		p, e := Discover(selected)
		if e != nil || p.Status != "incomplete" || p.Installation["state"] != "incomplete" {
			t.Fatalf("%v %v", p, e)
		}
		if _, err := Navigation(p, ""); err == nil {
			t.Fatal("early incomplete project presented as empty navigation")
		}
		if p.Registration.ProjectID != "" || len(p.Registration.Models) != 0 {
			t.Fatal("invented registration")
		}
	}
	os.Rename(op, filepath.Join(area, ".sdp-operations", "install-bad"))
	if _, e := Discover(root); e == nil {
		t.Fatal("unknown operation directory accepted")
	}
}
