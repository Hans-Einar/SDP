package sdptool

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRuntimeRegistrationConflictsAndShell(t *testing.T) {
	ctx := context.Background()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("PATH", t.TempDir())
	engine := filepath.Join(home, "gh-sdp")
	os.WriteFile(engine, []byte("fixture"), 0700)
	setup := func() (SetupResult, error) {
		return registerRuntime(ctx, func() (string, error) { return engine, nil })
	}
	r, err := setup()
	if err != nil || r.Status != "registered" || len(r.Warnings) != 1 {
		t.Fatalf("%+v %v", r, err)
	}
	again, err := setup()
	if err != nil || again.Status != "already-registered" {
		t.Fatal(again, err)
	}
	t.Setenv("PATH", filepath.Dir(r.Entrypoint))
	r, err = setup()
	if err != nil || len(r.Warnings) != 0 {
		t.Fatal(r, err)
	}
	if err = os.Remove(r.Entrypoint); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(r.Entrypoint, []byte("owned elsewhere"), 0700); err != nil {
		t.Fatal(err)
	}
	if _, err = setup(); err == nil {
		t.Fatal("overwrote foreign regular file")
	}
	b, _ := os.ReadFile(r.Entrypoint)
	if string(b) != "owned elsewhere" {
		t.Fatal("modified foreign file")
	}
	os.Remove(r.Entrypoint)
	os.Symlink("/missing/foreign", r.Entrypoint)
	if _, err = setup(); err == nil {
		t.Fatal("overwrote foreign dangling symlink")
	}
	target, _ := os.Readlink(r.Entrypoint)
	if target != "/missing/foreign" {
		t.Fatal(target)
	}
}
func TestRuntimeActionsKeepContextAndSourceRevision(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	ctx := context.Background()
	invoke := func(project, action string, params any) InvokeResponse {
		return Invoke(ctx, bytes.NewReader(requestBytes(t, project, action, params)), ActionCaller{})
	}
	// The test executable is transient and intentionally cannot register itself.
	if r := invoke("", "runtime.setup", map[string]any{}); r.Error == nil || !strings.Contains(r.Error.Message, "transient") {
		t.Fatalf("%+v", r)
	}
	if r := invoke(t.TempDir(), "runtime.setup", map[string]any{}); r.ExitCode == 0 {
		t.Fatal("host action accepted project")
	}
	if r := invoke(t.TempDir(), "project.discover", map[string]any{}); r.ExitCode == 0 {
		t.Fatal("project action bypassed context")
	}
	source := filepath.Join(t.TempDir(), "page.sdui")
	os.WriteFile(source, []byte(`sdui 0.2; page=[button("OK")];`), 0600)
	p := map[string]any{"source": source, "entry": "page"}
	r := invoke("", "sdui.source.check", p)
	if r.ExitCode != 0 {
		t.Fatalf("%+v", r)
	}
	var out, errs bytes.Buffer
	if code := Run(ctx, []string{"sdui", "check", "-source", source, "-entry", "page", "--json"}, &out, &errs); code != 0 {
		t.Fatal(errs.String())
	}
	var cli any
	json.Unmarshal(out.Bytes(), &cli)
	a, _ := json.Marshal(r.Result)
	b, _ := json.Marshal(cli)
	if string(a) != string(b) {
		t.Fatalf("CLI/action diverged: %s %s", a, b)
	}
	p["revision"] = strings.Repeat("0", 64)
	if r := invoke("", "sdui.source.check", p); r.ExitCode == 0 {
		t.Fatal("stale source accepted")
	}
}
func TestEmbeddedNativeCacheIntegrity(t *testing.T) {
	old := bundledNativeHost
	t.Cleanup(func() { bundledNativeHost = old })
	var data bytes.Buffer
	z := gzip.NewWriter(&data)
	z.Write([]byte("native fixture bytes"))
	z.Close()
	bundledNativeHost = data.Bytes()
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	p, err := nativeHost()
	if err != nil {
		t.Fatal(err)
	}
	if p2, e := nativeHost(); e != nil || p2 != p {
		t.Fatal(p2, e)
	}
	if err = os.WriteFile(p, []byte("corrupted"), 0700); err != nil {
		t.Fatal(err)
	}
	if _, err = nativeHost(); err == nil {
		t.Fatal("cached helper corruption accepted")
	}
	os.Remove(p)
	os.Symlink("/bin/true", p)
	if _, err = nativeHost(); err == nil {
		t.Fatal("cached helper symlink accepted")
	}
}

func TestSduiGlobalOutputArguments(t *testing.T) {
	for _, args := range [][]string{
		{"sdui", "preview", "--json", "-source", "Page.sdui"},
		{"--json", "sdui", "preview", "-source", "Page.sdui"},
		{"sdui", "preview", "-source", "Page.sdui", "--json"},
	} {
		got, mode, err := outputArguments(args)
		if err != nil || !mode || strings.Join(got, " ") != "sdui preview -source Page.sdui" {
			t.Fatal(got, mode, err)
		}
	}
	got, mode, err := outputArguments([]string{"sdui", "check", "-source", "--json"})
	if err != nil || mode || got[3] != "--json" {
		t.Fatal(got, mode, err)
	}
	got, mode, err = outputArguments([]string{"preview", "--json"})
	if err != nil || mode || got[1] != "--json" {
		t.Fatal("existing FILE positional changed", got, mode, err)
	}
}
