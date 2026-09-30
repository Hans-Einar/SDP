package sdptool

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const releaseNotes = "# Releases\n\n## [Unreleased]\n\nRelease-Date: unreleased\n\n### Added\n\n- Not released.\n\n## [1.1.0] - 2026-09-30\n\n### Added\n\n- Sessions.\n\n```text\n## [9.0.0] - 2026-09-30\n```\n\n## [1.0.0] - 2026-09-28\n\n### Fixed\n\n- Prior.\n"

func TestReleaseLogProjection(t *testing.T) {
	logs, e := ReleaseLogs(releaseNotes)
	if e != nil || len(logs) != 2 {
		t.Fatal(logs, e)
	}
	if bytes.Contains(logs["1.1.0"], []byte("Not released")) || !bytes.Contains(logs["1.1.0"], []byte("9.0.0")) {
		t.Fatal("wrong sections", logs)
	}
	for _, text := range []string{releaseNotes + "\n## [1.1.0] - 2026-09-30\n### Fixed\n- Duplicate", strings.Replace(releaseNotes, "2026-09-30", "2026-02-30", 1), "## [1.0.0] - 2026-09-30\n"} {
		if _, e := ReleaseLogs(text); e == nil {
			t.Fatal("invalid notes accepted")
		}
	}
}
func TestReleaseLogCLIImmutableAndCheck(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "RELEASE-NOTES.md"), []byte(releaseNotes), 0600)
	run := func(args ...string) int {
		var out, err bytes.Buffer
		return Run(context.Background(), append([]string{root, "release-log"}, args...), &out, &err)
	}
	args := []string{"--all", "--output", "Releases"}
	if run(append(args, "--check")...) == 0 {
		t.Fatal("check accepted absent output")
	}
	if run(args...) != 0 || run(append(args, "--check")...) != 0 || run(args...) != 0 {
		t.Fatal("generation/check/repeat failed")
	}
	p := filepath.Join(root, "Releases", "1.0.0.md")
	os.WriteFile(p, []byte("owner edit"), 0600)
	if run(args...) == 0 {
		t.Fatal("overwrote existing log")
	}
	b, _ := os.ReadFile(p)
	if string(b) != "owner edit" {
		t.Fatal("lost bytes")
	}
	if run("--version", "1.9.0") == 0 {
		t.Fatal("invented release")
	}
}
