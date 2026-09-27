package sdptool

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func executable(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "program with spaces")
	if e := os.WriteFile(p, []byte("#!/bin/sh\n"+content), 0700); e != nil {
		t.Fatal(e)
	}
	return p
}
func TestViewPlanArgumentsAndCleanup(t *testing.T) {
	root, r := projectFixture(t)
	r.ImplementationPlan = "SDP/plan with spaces.md"
	r.Models = []Model{{"model", "SDPTOOL", "model.design", "design-core/0.5"}}
	saveRegistration(t, root, r)
	os.WriteFile(filepath.Join(root, r.ImplementationPlan), []byte("authored plan"), 0600)
	os.WriteFile(filepath.Join(root, "model.design"), []byte(sample), 0600)
	log := filepath.Join(t.TempDir(), "args")
	t.Setenv("ARG_LOG", log)
	runtime := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", runtime)
	viewer := executable(t, `if [ "$1" = '--help' ]; then printf '%s\n' '--navigator --sdl-tool'; exit 0; fi
printf '%s\n' "$@" > "$ARG_LOG"
[ -f "$1" ] && [ "$2" = '--navigator' ] && [ -f "$3" ]
`)
	tool := executable(t, "exit 0\n")
	p, e := Discover(root)
	if e != nil {
		t.Fatal(e)
	}
	if e = ViewPlan(context.Background(), p, "model", Host{Viewer: viewer, SDLTool: tool}); e != nil {
		t.Fatal(e)
	}
	args, _ := os.ReadFile(log)
	if !strings.Contains(string(args), "--sdl-source\n"+root+"/model.design\n") || !strings.HasPrefix(string(args), root+"/SDP/plan with spaces.md\n") {
		t.Fatal(string(args))
	}
	entries, _ := os.ReadDir(runtime)
	if len(entries) != 0 {
		t.Fatal("session leaked")
	}
	b, _ := os.ReadFile(filepath.Join(root, r.ImplementationPlan))
	if string(b) != "authored plan" {
		t.Fatal("plan replaced")
	}
	viewer = executable(t, "exit 4\n")
	if e = ViewPlan(context.Background(), p, "model", Host{Viewer: viewer, SDLTool: tool}); e == nil {
		t.Fatal("failed probe accepted")
	}
}
func TestPlanOnlyCancellationAndMissingTools(t *testing.T) {
	root, r := projectFixture(t)
	r.ImplementationPlan = "SDP/plan.md"
	saveRegistration(t, root, r)
	os.WriteFile(filepath.Join(root, r.ImplementationPlan), []byte("plan"), 0600)
	p, _ := Discover(root)
	if e := ViewPlan(context.Background(), p, "", Host{Viewer: "/missing/viewer"}); e == nil {
		t.Fatal("missing viewer accepted")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	viewer := executable(t, "exec sleep 10\n")
	if e := ViewPlan(ctx, p, "", Host{Viewer: viewer}); e == nil || !strings.Contains(e.Error(), "canceled") {
		t.Fatalf("%v", e)
	}
}
