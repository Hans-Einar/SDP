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

func projectFixture(t *testing.T) (string, Registration) {
	t.Helper()
	root := t.TempDir()
	os.Mkdir(filepath.Join(root, "SDP"), 0700)
	r := Registration{SchemaVersion: "1.0", ProjectID: "trial", ProcessProfile: "sdp-five-phase/0.1", Models: []Model{}, SDUI: []Model{}}
	saveRegistration(t, root, r)
	return root, r
}
func saveRegistration(t *testing.T, root string, r Registration) {
	t.Helper()
	b, e := json.Marshal(r)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(root, "SDP/navigation.json"), b, 0600); e != nil {
		t.Fatal(e)
	}
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
	if Run(context.Background(), []string{root, "discover"}, &out, &errs) != 0 {
		t.Fatal(errs.String())
	}
	if !strings.Contains(out.String(), `"installation":{"state":"unknown"}`) {
		t.Fatal(out.String())
	}
}
func TestDiscoveryInvalidAndUnsupported(t *testing.T) {
	for _, change := range []string{"unknown", "duplicate", "trailing", "null", "escape", "symlink", "schema", "profile", "id"} {
		t.Run(change, func(t *testing.T) {
			root, r := projectFixture(t)
			file := filepath.Join(root, "SDP/navigation.json")
			switch change {
			case "unknown":
				b, _ := os.ReadFile(file)
				os.WriteFile(file, append([]byte(`{"command":"evil",`), b[1:]...), 0600)
			case "duplicate":
				b, _ := os.ReadFile(file)
				os.WriteFile(file, append([]byte(`{"projectId":"other",`), b[1:]...), 0600)
			case "trailing":
				b, _ := os.ReadFile(file)
				os.WriteFile(file, append(b, []byte(" {}")...), 0600)
			case "null":
				r.Models = nil
				saveRegistration(t, root, r)
			case "escape":
				r.ImplementationPlan = "../secret"
				saveRegistration(t, root, r)
			case "symlink":
				os.Symlink(t.TempDir(), filepath.Join(root, "outside"))
				r.KanBan = "outside/board"
				saveRegistration(t, root, r)
			case "schema":
				r.SchemaVersion = "2"
				saveRegistration(t, root, r)
			case "profile":
				r.ProcessProfile = "future"
				saveRegistration(t, root, r)
			case "id":
				r.ProjectID = "../bad"
				saveRegistration(t, root, r)
			}
			before, _ := os.ReadFile(file)
			p, e := Discover(root)
			if e == nil {
				t.Fatal("accepted invalid")
			}
			want := "invalid"
			if change == "schema" || change == "profile" {
				want = "unsupported"
			}
			if p.Status != want {
				t.Fatalf("%s %v", p.Status, e)
			}
			after, _ := os.ReadFile(file)
			if !bytes.Equal(before, after) {
				t.Fatal("modified metadata")
			}
		})
	}
}
func TestReferencedInstallationFacts(t *testing.T) {
	root, r := projectFixture(t)
	r.ProjectManifest = "SDP/SDP-project.manifest.yaml"
	saveRegistration(t, root, r)
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
