package sdptool

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func projectFixture(t *testing.T) (string, Inventory) {
	t.Helper()
	root := t.TempDir()
	os.Mkdir(filepath.Join(root, "SDP"), 0700)
	r := Inventory{SchemaVersion: "1.0", ProjectID: "trial", ProcessProfile: "sdp-five-phase/0.1", Models: []Model{}, SDUI: []Model{}}
	saveInventory(t, root, r)
	return root, r
}

// Unit services may consume an in-memory Inventory. This helper only supplies
// project identity; source discovery never reads a saved inventory or register.
func saveInventory(t *testing.T, root string, r Inventory) {
	t.Helper()
	os.WriteFile(filepath.Join(root, "SDP/SDP-project.manifest.yaml"), []byte("schemaVersion: '1.0'\nproject:\n  name: "+r.ProjectID+"\n"), 0600)
}
func TestDiscoveryRootsAndNoParentGuess(t *testing.T) {
	root, _ := projectFixture(t)
	a, e := Discover(root)
	if e != nil || a.Status != "valid" || a.Capabilities["sdl"] != "absent" {
		t.Fatalf("%+v %v", a, e)
	}
	b, e := Discover(filepath.Join(root, "SDP"))
	if e != nil || a.Root != b.Root || a.Area != b.Area {
		t.Fatalf("roots %v", e)
	}
	nested := filepath.Join(root, "nested")
	os.Mkdir(nested, 0700)
	if _, e = Discover(nested); e == nil {
		t.Fatal("guessed parent")
	}
	var out, errs bytes.Buffer
	if Run(context.Background(), []string{root, "discover", "--json"}, &out, &errs) != 0 {
		t.Fatal(errs.String())
	}
	if !strings.Contains(out.String(), `"installation":{"state":"unknown"}`) {
		t.Fatal(out.String())
	}
}
func TestObsoleteRegistrationCannotControlDiscovery(t *testing.T) {
	root, _ := projectFixture(t)
	for _, data := range []string{`{`, `{"schemaVersion":"99","projectId":"malicious","models":[]}`, `{"command":"evil","models":null}`} {
		f := filepath.Join(root, "SDP/navigation.json")
		os.WriteFile(f, []byte(data), 0600)
		p, e := Discover(root)
		if e != nil || p.Inventory.ProjectID != "trial" {
			t.Fatalf("obsolete file used: %+v %v", p, e)
		}
		after, _ := os.ReadFile(f)
		if string(after) != data {
			t.Fatal("rewrote old project data")
		}
	}
}
func TestDiscoveryPathAndManifestSafety(t *testing.T) {
	root, _ := projectFixture(t)
	outside := t.TempDir()
	os.WriteFile(filepath.Join(outside, "secret.design"), []byte(sample), 0600)
	os.Symlink(outside, filepath.Join(root, "SDP/outside"))
	p, e := Discover(root)
	if e != nil || len(p.Inventory.Models) != 0 {
		t.Fatalf("followed symlink: %v", e)
	}
	for _, data := range []string{"schemaVersion: '99'", "schemaVersion: [", "schemaVersion: '1.0'\ninstalled:\n  manifestPath: ../secret.yaml\n"} {
		os.WriteFile(filepath.Join(root, "SDP/SDP-project.manifest.yaml"), []byte(data), 0600)
		if _, e := Discover(root); e == nil {
			t.Fatal("invalid/escaping manifest accepted")
		}
	}
}
func TestReferencedInstallationFacts(t *testing.T) {
	root, r := projectFixture(t)
	r.ProjectManifest = "SDP/SDP-project.manifest.yaml"
	saveInventory(t, root, r)
	f := filepath.Join(root, r.ProjectManifest)
	os.WriteFile(f, []byte("schemaVersion: '1.0'\ninstalled:\n  manifestPath: installed.yaml\n"), 0600)
	installed := filepath.Join(root, "SDP/installed.yaml")
	os.WriteFile(installed, []byte("schemaVersion: '1.0'\ntoolkitVersion: '0.2.0'\n"), 0600)
	p, e := Discover(root)
	if e != nil || p.Installation["state"] != "declared" {
		t.Fatalf("%v %v", p, e)
	}
	for _, bad := range []string{"schemaVersion: '2'", "schemaVersion: [broken", "schemaVersion: '1.0'\n---\nother: record"} {
		os.WriteFile(installed, []byte(bad), 0600)
		if _, e = Discover(root); e == nil {
			t.Fatal("accepted invalid installed facts")
		}
	}
	os.Remove(installed)
	if _, e = Discover(root); e == nil {
		t.Fatal("missing declared facts")
	}
}

func TestMultipleSourcesRequireExplicitSelection(t *testing.T) {
	root, _ := projectFixture(t)
	for _, name := range []string{"one", "two"} {
		os.WriteFile(filepath.Join(root, "SDP", name+".design"), []byte(sample), 0600)
	}
	p, e := Discover(root)
	if e != nil {
		t.Fatal(e)
	}
	if _, _, e = p.model("", false); e == nil {
		t.Fatal("arbitrary default chosen")
	}
	id := p.Inventory.Models[1].ID
	m, _, e := p.model(id, false)
	if e != nil || m.ID != id {
		t.Fatal(m, e)
	}
}
