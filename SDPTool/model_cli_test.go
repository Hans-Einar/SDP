package sdptool

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCompiledModelJourney(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "sdptool")
	build := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "./cmd/sdptool")
	if b, e := build.CombinedOutput(); e != nil {
		t.Fatalf("build: %v %s", e, b)
	}
	area := t.TempDir()
	run := func(args ...string) map[string]any {
		t.Helper()
		cmd := exec.Command(binary, append([]string{area, "model", "--json"}, args...)...)
		b, e := cmd.CombinedOutput()
		if e != nil {
			t.Fatalf("%v: %v %s", args, e, b)
		}
		var r map[string]any
		if e = json.Unmarshal(b, &r); e != nil {
			t.Fatal(e, string(b))
		}
		return r
	}
	put := func(work, name, content string) {
		t.Helper()
		if e := os.WriteFile(filepath.Join(area, "WORK--"+work, name), []byte(content), 0600); e != nil {
			t.Fatal(e)
		}
	}
	run("create", "work:Initial", "--initial")
	put("Initial", "Main.design", "language design-core version 0.5.\nunit Main.\n")
	put("Initial", "Page.sdui", `sdui 0.2; Page=[<"Hello", button("OK")>];`)
	run("commit", "work:Initial", "--message", "valid sources")
	put("Initial", "Main.design", "broken intermediate")
	run("restore", "work:Initial", "to", "commit:00001")
	run("history", "work:Initial")
	run("create", "candidate:Base", "from", "work:Initial")
	run("create", "release:0.1.0", "from", "candidate:Base", "--evidence", "model-only")
	run("create", "work:A")
	run("create", "work:B")
	put("A", "notes.txt", "ours")
	put("B", "notes.txt", "theirs")
	if r := run("merge", "work:B", "into", "work:A"); r["status"] != "conflicted" {
		t.Fatal(r)
	}
	put("A", "notes.txt", "resolved")
	run("commit", "work:A", "--message", "resolve integration", "--resolved")
	run("create", "candidate:Next", "from", "work:A")
	run("create", "release:0.2.0", "from", "candidate:Next", "--evidence", "model-only")
	run("create", "work:New")
	if r := run("snapshot", "work:New"); r["preliminary"] != true {
		t.Fatal(r)
	}
	cmd := exec.Command(binary, area, "model", "status", "work:New")
	b, e := cmd.CombinedOutput()
	if e != nil || json.Valid(b) {
		t.Fatalf("human output: %v %s", e, b)
	}
	copied := t.TempDir()
	if e = os.CopyFS(copied, os.DirFS(area)); e != nil {
		t.Fatal(e)
	}
	area = copied
	run("status", "release:0.2.0")
	if _, e = os.Stat(filepath.Join(area, ".git")); !os.IsNotExist(e) {
		t.Fatal("runtime requires git")
	}
}
