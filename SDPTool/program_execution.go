package sdptool

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// programRunRequest contains the already parsed identity and optional revision.
type programRunRequest struct{ Program, Revision string }

// runProgram executes a discovered declaration. Its subprocess is an intentional
// application boundary; it retains the process group and terminal lifecycle.
func runProgram(ctx context.Context, selected string, request programRunRequest, out, errs io.Writer) (int, error) {
	failed := func(code string, err error) (int, error) { return 1, failure(code, err) }
	p, err := discover(selected, false)
	if err != nil {
		return 1, err
	}
	if p.Status != "valid" || p.ProgramDiscovery.State != "valid" {
		return failed("program", fmt.Errorf("program discovery %s: %s", p.ProgramDiscovery.State, p.ProgramDiscovery.Diagnostic))
	}
	for _, program := range p.Programs {
		if program.ID != request.Program {
			continue
		}
		if program.State != "runnable" {
			return failed("program", fmt.Errorf("%s: %s", program.ID, program.Diagnostic))
		}
		if request.Revision != "" && request.Revision != program.Revision {
			return failed("stale", fmt.Errorf("program declaration or SDUI source changed; discover again"))
		}
		binary, err := programExecutable(p.Root, program.Command[0])
		if err != nil {
			return failed("program", err)
		}
		cmd := exec.CommandContext(ctx, binary, program.Command[1:]...)
		cmd.Dir, cmd.Stdin, cmd.Stdout, cmd.Stderr = p.Root, os.Stdin, out, errs
		restore, err := configureProgramProcess(cmd)
		if err != nil {
			return failed("program", err)
		}
		err = cmd.Run()
		if restoreErr := restore(); restoreErr != nil {
			return failed("program", errors.Join(err, restoreErr))
		}
		if err != nil {
			var child *exec.ExitError
			if errors.As(err, &child) && child.ExitCode() >= 0 {
				return child.ExitCode(), nil
			}
			return failed("program", err)
		}
		return 0, nil
	}
	return failed("selection", fmt.Errorf("unknown program %q; use discover", request.Program))
}
