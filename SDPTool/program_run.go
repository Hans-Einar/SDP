package sdptool

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
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
	p, err := discover(selected, false)
	if err != nil {
		return reportMode(errs, err, false)
	}
	if p.Status != "valid" || p.ProgramDiscovery.State != "valid" {
		return fail("program", fmt.Errorf("program discovery %s: %s", p.ProgramDiscovery.State, p.ProgramDiscovery.Diagnostic))
	}
	for _, program := range p.Programs {
		if program.ID != *id {
			continue
		}
		if program.State != "runnable" {
			return fail("program", fmt.Errorf("%s: %s", program.ID, program.Diagnostic))
		}
		if *revision != "" && *revision != program.Revision {
			return fail("stale", fmt.Errorf("program declaration or SDUI source changed; discover again"))
		}
		binary, err := programExecutable(p.Root, program.Command[0])
		if err != nil {
			return fail("program", err)
		}
		cmd := exec.CommandContext(ctx, binary, program.Command[1:]...)
		cmd.Dir, cmd.Stdin, cmd.Stdout, cmd.Stderr = p.Root, os.Stdin, out, errs
		restore, err := configureProgramProcess(cmd)
		if err != nil {
			return fail("program", err)
		}
		err = cmd.Run()
		if restoreErr := restore(); restoreErr != nil {
			return fail("program", errors.Join(err, restoreErr))
		}
		if err != nil {
			var child *exec.ExitError
			if errors.As(err, &child) && child.ExitCode() >= 0 {
				return child.ExitCode()
			}
			return fail("program", err)
		}
		return 0
	}
	return fail("selection", fmt.Errorf("unknown program %q; use discover", *id))
}
