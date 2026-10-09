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

func programFixture(t *testing.T) (string, ProgramDeclaration) {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "SDP/UI"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "SDP/UI/main.sdui"), []byte(`sdui 0.2; Page=[<"Hello", button("OK")>];`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "runner"), []byte("#!/bin/sh\nprintf '%s\\n' \"$PWD\" \"$@\"\nprintf ran > ran\nexit 7\n"), 0700); err != nil {
		t.Fatal(err)
	}
	d := ProgramDeclaration{ID: "widget-lab", Label: "Widget Lab", Source: "SDP/UI/main.sdui", Entry: "Page", Command: []string{"./runner", "two words", "$(touch injected)", "--json"}}
	writePrograms(t, root, []ProgramDeclaration{d})
	return root, d
}

func writePrograms(t *testing.T, root string, declarations []ProgramDeclaration) {
	t.Helper()
	b, err := json.Marshal(map[string]any{"schemaVersion": "sdp-programs/1", "programs": declarations})
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(root, programsManifest), b, 0600); err != nil {
		t.Fatal(err)
	}
}

func TestProgramDiscoveryAndForegroundRun(t *testing.T) {
	root, _ := programFixture(t)
	p, err := Discover(root)
	if err != nil || len(p.Sources) != 1 || len(p.Programs) != 1 || p.Programs[0].State != "runnable" || p.Capabilities["sdui-programs"] != "discovered" {
		t.Fatalf("%+v %v", p, err)
	}
	if _, err := os.Stat(filepath.Join(root, "ran")); !os.IsNotExist(err) {
		t.Fatal("discovery executed command")
	}
	q, err := Discover(filepath.Join(root, "SDP"))
	if err != nil || q.Programs[0].Revision != p.Programs[0].Revision {
		t.Fatal("unstable root selection", err)
	}
	var out, diagnostic bytes.Buffer
	if code := Run(context.Background(), append([]string{root}, p.Programs[0].Run...), &out, &diagnostic); code != 7 {
		t.Fatalf("exit=%d %s", code, diagnostic.String())
	}
	if out.String() != root+"\ntwo words\n$(touch injected)\n--json\n" {
		t.Fatalf("argv/cwd: %q", out.String())
	}
	if _, err := os.Stat(filepath.Join(root, "injected")); !os.IsNotExist(err) {
		t.Fatal("shell expansion")
	}
	if _, err := os.Stat(filepath.Join(root, "ran")); err != nil {
		t.Fatal(err)
	}
}

func TestProgramBlockedDeclarationsPreserveSources(t *testing.T) {
	for _, name := range []string{"missing-command", "escape-command", "absolute-command", "missing-source", "escaping-source", "entry", "empty-command", "label", "symlink-command"} {
		t.Run(name, func(t *testing.T) {
			root, d := programFixture(t)
			switch name {
			case "missing-command":
				d.Command = []string{"./absent"}
			case "escape-command":
				d.Command = []string{"../runner"}
			case "absolute-command":
				d.Command = []string{"/bin/true"}
			case "missing-source":
				d.Source = "SDP/UI/absent.sdui"
			case "escaping-source":
				d.Source = "../main.sdui"
			case "entry":
				d.Entry = "Missing"
			case "empty-command":
				d.Command = nil
			case "label":
				d.Label = "\n"
			case "symlink-command":
				if err := os.Symlink("/bin/true", filepath.Join(root, "outside")); err != nil {
					t.Fatal(err)
				}
				d.Command = []string{"./outside"}
			}
			writePrograms(t, root, []ProgramDeclaration{d})
			p, err := Discover(root)
			if err != nil || len(p.Sources) != 1 || p.Programs[0].State != "blocked" || p.Programs[0].Diagnostic == "" || len(p.Programs[0].Run) != 0 {
				t.Fatalf("%+v %v", p, err)
			}
			var out, diagnostic bytes.Buffer
			if Run(context.Background(), []string{root, "run", "--program", d.ID}, &out, &diagnostic) == 0 {
				t.Fatal("blocked command ran")
			}
		})
	}
}

func TestProgramInvalidCatalogAndRemoval(t *testing.T) {
	root, d := programFixture(t)
	writePrograms(t, root, []ProgramDeclaration{d, d})
	p, err := Discover(root)
	if err != nil || p.ProgramDiscovery.State != "invalid" || len(p.Programs) != 0 || len(p.Sources) != 1 {
		t.Fatalf("%+v %v", p, err)
	}
	for _, data := range []string{`{"schemaVersion":"future","programs":[]}`, `{"schemaVersion":"sdp-programs/1","programs":[],"programs":[]}`, `{"schemaVersion":"sdp-programs/1","programs":[],"extra":1}`} {
		if err := os.WriteFile(filepath.Join(root, programsManifest), []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
		p, err = Discover(root)
		if err != nil || p.ProgramDiscovery.State != "invalid" || len(p.Sources) != 1 {
			t.Fatalf("%+v %v", p, err)
		}
	}
	if err := os.Remove(filepath.Join(root, programsManifest)); err != nil {
		t.Fatal(err)
	}
	p, err = Discover(root)
	if err != nil || p.ProgramDiscovery.State != "absent" || len(p.Programs) != 0 || len(p.Sources) != 1 {
		t.Fatalf("%+v %v", p, err)
	}
}

func TestProgramRefreshAndStaleSelection(t *testing.T) {
	for _, change := range []string{"source", "declaration", "remove"} {
		t.Run(change, func(t *testing.T) {
			root, d := programFixture(t)
			p, err := Discover(root)
			if err != nil {
				t.Fatal(err)
			}
			switch change {
			case "source":
				err = os.WriteFile(filepath.Join(root, d.Source), []byte(`sdui 0.2; Page=[<"Changed", button("OK")>];`), 0600)
			case "declaration":
				d.Command = []string{"./runner", "changed"}
				writePrograms(t, root, []ProgramDeclaration{d})
			case "remove":
				err = os.Remove(filepath.Join(root, programsManifest))
			}
			if err != nil {
				t.Fatal(err)
			}
			var out, diagnostic bytes.Buffer
			if Run(context.Background(), append([]string{root}, p.Programs[0].Run...), &out, &diagnostic) == 0 {
				t.Fatal("stale selection ran")
			}
			if _, err := os.Stat(filepath.Join(root, "ran")); !os.IsNotExist(err) {
				t.Fatal("stale command executed")
			}
		})
	}
}

func TestProgramCLISelectionAndPresentation(t *testing.T) {
	cleaned, jsonMode, err := outputArguments([]string{"run", "--program", "--json"})
	if err != nil || jsonMode || len(cleaned) != 3 {
		t.Fatal("program value treated as global option", cleaned, jsonMode, err)
	}
	root, _ := programFixture(t)
	for _, args := range [][]string{{"run"}, {"run", "--program", "missing"}, {"run", "--program", "widget-lab", "--json"}, {"run", "--program", "widget-lab", "extra"}} {
		var out, diagnostic bytes.Buffer
		if Run(context.Background(), append([]string{root}, args...), &out, &diagnostic) == 0 {
			t.Fatal(args)
		}
		if _, err := os.Stat(filepath.Join(root, "ran")); !os.IsNotExist(err) {
			t.Fatal("invalid invocation executed")
		}
	}
	var out, diagnostic bytes.Buffer
	if code := Run(context.Background(), []string{root, "discover"}, &out, &diagnostic); code != 0 || !strings.Contains(out.String(), "widget-lab — Widget Lab [runnable]") {
		t.Fatalf("%d %s %s", code, out.String(), diagnostic.String())
	}
	out.Reset()
	if code := Run(context.Background(), []string{root, "discover", "--json"}, &out, &diagnostic); code != 0 {
		t.Fatal(code, diagnostic.String())
	}
	var p Project
	if err := json.Unmarshal(out.Bytes(), &p); err != nil || len(p.Programs) != 1 {
		t.Fatal(err, out.String())
	}
}
