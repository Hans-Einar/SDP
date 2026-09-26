package sdptool

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// This consumer uses the compiled executable and JSON only, not facade API calls
// for the journeys under test. Build is test setup, never viewer startup behavior.
func TestExecutableConsumerJourney(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "sdptool")
	build := exec.Command(filepath.Join(runtime.GOROOT(), "bin/go"), "build", "-o", binary, "./cmd/sdptool")
	if b, e := build.CombinedOutput(); e != nil {
		t.Fatalf("build: %s %v", b, e)
	}
	run := func(args ...string) (map[string]any, string, error) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		c := exec.CommandContext(ctx, binary, args...)
		var out, errs bytes.Buffer
		c.Stdout = &out
		c.Stderr = &errs
		e := c.Run()
		var v map[string]any
		if e == nil {
			if err := json.Unmarshal(out.Bytes(), &v); err != nil {
				t.Fatalf("invalid response: %s %v", out.String(), err)
			}
		}
		return v, errs.String(), e
	}
	root, r := projectFixture(t)
	r.Models = []Model{{"model", "System", "model with spaces.design", "design-core/0.5"}}
	r.ImplementationPlan = "SDP/plan.md"
	saveRegistration(t, root, r)
	source := filepath.Join(root, r.Models[0].Source)
	original, e := os.ReadFile("testdata/consumer.design")
	if e != nil {
		t.Fatal(e)
	}
	os.WriteFile(source, original, 0600)
	os.WriteFile(filepath.Join(root, r.ImplementationPlan), []byte("authored plan"), 0600)
	var expectations struct {
		Schema     string   `json:"producerSchema"`
		Roots      []string `json:"roots"`
		Viewpoints int      `json:"viewpoints"`
		Diagram    string   `json:"firstDiagram"`
		Files      []string `json:"requiredBundleFiles"`
		Stale      string   `json:"staleError"`
	}
	f, _ := os.ReadFile("testdata/consumer-expectations.json")
	if e = json.Unmarshal(f, &expectations); e != nil {
		t.Fatal(e)
	}
	discovered, stderr, e := run(root, "discover")
	if e != nil || discovered["schema"] != expectations.Schema || discovered["status"] != "valid" {
		t.Fatalf("discovery %v %s", e, stderr)
	}
	tree, stderr, e := run(filepath.Join(root, "SDP"), "tree", "--model", "model")
	if e != nil {
		t.Fatalf("tree %v %s", e, stderr)
	}
	revision := tree["revision"].(string)
	roots := tree["roots"].([]any)
	for i, want := range expectations.Roots {
		if roots[i] != want {
			t.Fatal(roots)
		}
	}
	nodes := tree["nodes"].([]any)
	vp := 0
	uri := ""
	ids := map[string]bool{}
	for _, raw := range nodes {
		n := raw.(map[string]any)
		ids[n["id"].(string)] = true
		if n["kind"] == "viewpoint" {
			vp++
		}
		if n["kind"] == "diagram" && strings.HasSuffix(n["id"].(string), expectations.Diagram) {
			uri = n["target"].(map[string]any)["uri"].(string)
		}
	}
	if vp != expectations.Viewpoints || uri == "" {
		t.Fatalf("viewpoints=%d uri=%s", vp, uri)
	}
	for _, raw := range nodes {
		n := raw.(map[string]any)
		if refs, ok := n["children"].([]any); ok {
			for _, ref := range refs {
				if !ids[ref.(string)] {
					t.Fatal("unresolved child")
				}
			}
		}
		if ref, ok := n["reference"].(string); ok && !ids[ref] {
			t.Fatal("unresolved reference")
		}
	}
	output := filepath.Join(t.TempDir(), "request")
	selected, stderr, e := run(root, "select", "--model", "model", "--uri", uri, "--revision", revision, "--output", output)
	if e != nil {
		t.Fatalf("select %v %s", e, stderr)
	}
	if selected["revision"] != revision || selected["operation"] != "select" {
		t.Fatal(selected)
	}
	for _, file := range expectations.Files {
		if _, e = os.Stat(filepath.Join(output, file)); e != nil {
			t.Fatal(e)
		}
	}
	previous, _ := os.ReadFile(filepath.Join(output, "entry.md"))
	os.WriteFile(source, []byte(strings.ReplaceAll(string(original), "Child", "Edited")), 0600)
	_, stderr, e = run(root, "select", "--model", "model", "--uri", uri, "--revision", revision, "--output", output)
	if e == nil || !strings.Contains(stderr, `"code":"`+expectations.Stale+`"`) {
		t.Fatalf("stale: %v %s", e, stderr)
	}
	after, _ := os.ReadFile(filepath.Join(output, "entry.md"))
	if !bytes.Equal(previous, after) {
		t.Fatal("stale response replaced visible document")
	}
	preview, stderr, e := run("preview", source, "--output", output)
	if e != nil || preview["revision"] == revision {
		t.Fatalf("regenerate %v %s", e, stderr)
	}
	viewer := executable(t, `if [ "$1" = '--help' ]; then printf '%s\n' '--navigator --sdl-tool'; exit 0; fi
[ -f "$1" ] && [ -f "$3" ]
`)
	sdl := executable(t, "exit 0\n")
	if _, stderr, e = run(root, "view", "ip", "--model", "model", "--viewer", viewer, "--sdl-tool", sdl); e != nil {
		t.Fatalf("view %v %s", e, stderr)
	}
	for _, args := range [][]string{{root, "generate", "ip"}, {root, "tree", "--model", "missing"}, {"preview", source, "--output", output, "--viewpoint", "VP99"}} {
		if _, stderr, e = run(args...); e == nil || !json.Valid([]byte(stderr)) {
			t.Fatalf("error response %v %s", e, stderr)
		}
	}
}
func TestReviewStrictOptionalMetadata(t *testing.T) {
	root, _ := projectFixture(t)
	file := filepath.Join(root, "SDP/navigation.json")
	base, _ := os.ReadFile(file)
	for _, value := range []string{"null", `""`, "false"} {
		v := strings.Replace(string(base), `"models"`, `"implementationPlan":`+value+`,"models"`, 1)
		os.WriteFile(file, []byte(v), 0600)
		if _, e := Discover(root); e == nil {
			t.Fatal("accepted invalid optional path", value)
		}
	}
}
