package sdptool

import (
	"context"
	"flag"
	"fmt"
	"io"
)

func runProgramCommand(ctx context.Context, selected string, args []string, out, errs io.Writer, jsonMode bool) int {
	fail := func(code string, err error) int { return reportMode(errs, failure(code, err), jsonMode) }
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	id := fs.String("program", "", "discovered program ID")
	revision := fs.String("revision", "", "expected declaration and SDUI source revision")
	if err := fs.Parse(args); err != nil {
		return fail("arguments", err)
	}
	if *id == "" || fs.NArg() != 0 || jsonMode {
		return fail("arguments", fmt.Errorf("run requires --program ID, accepts optional --revision HASH, and streams application output (no --json)"))
	}
	code, err := runProgram(ctx, selected, programRunRequest{Program: *id, Revision: *revision}, out, errs)
	if err != nil {
		return reportMode(errs, err, false)
	}
	return code
}
