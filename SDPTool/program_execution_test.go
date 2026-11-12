package sdptool

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestTypedProgramExecutionRevisionAndExit(t *testing.T) {
	root, declaration := programFixture(t)
	var out, errs bytes.Buffer
	code, err := runProgram(context.Background(), root, programRunRequest{Program: declaration.ID, Revision: "stale"}, &out, &errs)
	problem, ok := err.(*Failure)
	if code != 1 || !ok || problem.Code != "stale" {
		t.Fatalf("stale: %d %v", code, err)
	}
	if _, err := os.Stat(filepath.Join(root, "ran")); !os.IsNotExist(err) {
		t.Fatal("stale request executed program")
	}
	code, err = runProgram(context.Background(), root, programRunRequest{Program: declaration.ID}, &out, &errs)
	if code != 7 || err != nil {
		t.Fatalf("child exit: %d %v %s", code, err, &errs)
	}
	if out.String() != root+"\ntwo words\n$(touch injected)\n--json\n" {
		t.Fatalf("argv/cwd: %q", out.String())
	}
}
