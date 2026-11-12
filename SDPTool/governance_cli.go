package sdptool

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governanceapp"
	"github.com/Hans-Einar/SDP/SDPTool/presentation"
)

const governanceHelp = `Usage: sdptool [PROJECT] governance context|route|bind|start|read|transition [options]
  --input FILE   bounded JSON Selection, RouteRequest, StartRequest or TransitionRequest
  --run UUID     run to read (read only)
  --instance UUID  Session instance for session-read
Session operations: session-definition, session-context, session-successor, session-adopt, session-read,
session-list, session-design-context, session-rebind-context, session-rebind, session-select, session-assess, session-start, session-progress, session-verify, session-review, session-dispose, session-owner, session-advance, session-resume. Selection prepares checked ready-card effects; resume commits them.
Record effects commit through session-resume; host launch remains separate.
session-owner records an actual owner decision or requests a concrete owner review.
  --state DIR    private operational state outside project (defaults to XDG state home)
  --json         machine output; human-readable output is the default
context and route only read; bind is explicit. start/transition are trusted local
controller operations, not agent/owner approval endpoints. PGI1 cannot complete
work, publish records or launch agents. See SDPTool/governance/README.md.
`

func governanceCommand(selected string, args []string, out, errs io.Writer, jsonMode bool) int {
	return governanceCommandAs(selected, args, out, errs, jsonMode, governance.LocalOwner())
}
func governanceCommandAs(selected string, args []string, out, errs io.Writer, jsonMode bool, caller governance.Caller) int {
	report := func(e error) int {
		var domain *governance.Error
		if !errors.As(e, &domain) {
			domain = &governance.Error{Code: "unavailable", Message: e.Error()}
		}
		v := map[string]any{"schema": governance.Schema, "operation": "error", "status": "error", "error": domain}
		_ = presentation.Default().Write(errs, v, jsonMode)
		return 2
	}
	if len(args) == 1 && (args[0] == "help" || args[0] == "--help") {
		_, e := io.WriteString(out, governanceHelp)
		if e != nil {
			return report(e)
		}
		return 0
	}
	if len(args) == 0 {
		return report(fmt.Errorf("governance command required"))
	}
	op := args[0]
	fs := flag.NewFlagSet("governance", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	input := fs.String("input", "", "JSON request")
	run := fs.String("run", "", "run identity")
	instance := fs.String("instance", "", "Session instance identity")
	state := fs.String("state", "", "state directory")
	if e := fs.Parse(args[1:]); e != nil {
		return report(e)
	}
	if fs.NArg() != 0 {
		return report(fmt.Errorf("unexpected argument"))
	}
	operation, ok := governanceapp.Find(op)
	if !ok {
		if *instance != "" {
			return report(fmt.Errorf("--instance is only valid for session-read"))
		}
		return report(fmt.Errorf("unknown governance command"))
	}
	if e := operation.Validate(*input, *run, *instance); e != nil {
		return report(e)
	}
	service, e := governanceService(selected, *state)
	if e != nil {
		return report(e)
	}
	var b []byte
	if *input != "" {
		b, e = boundedFile(*input)
		if e != nil {
			return report(e)
		}
	}
	identity := *run
	if operation.Input == governanceapp.InstanceInput {
		identity = *instance
	}
	result, e := operation.Execute(context.Background(), service, caller, b, identity)
	if e != nil {
		return report(e)
	}
	if e = presentation.Default().Write(out, result, jsonMode); e != nil {
		return report(e)
	}
	return 0
}
