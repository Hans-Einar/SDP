package install

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
)

func Run(ctx context.Context, root, op string, args []string, out, errs io.Writer) int {
	fs := flag.NewFlagSet(op, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	o := Options{Root: root, Operation: op}
	fs.StringVar(&o.Artifact, "artifact", "", "explicit local development descriptor")
	fs.StringVar(&o.PreviousArtifact, "previous-artifact", "", "original development descriptor")
	fs.StringVar(&o.Manifest, "manifest", "", "adoption JSON/YAML")
	fs.BoolVar(&o.AllowUnreleased, "allow-unreleased", false, "allow development artifacts")
	dest := fs.String("plan-output", "", "new plan file outside project")
	jsonMode := fs.Bool("json", false, "machine-readable result")
	e := fs.Parse(args)
	if e == nil && fs.NArg() != 0 {
		e = fmt.Errorf("unexpected positional arguments")
	}
	if e == nil && o.Artifact == "" {
		e = fmt.Errorf("select --artifact and --allow-unreleased; signed releases are not implemented yet")
	}
	if e != nil {
		return output(out, errs, *jsonMode, root, op, nil, fail("arguments", 2, "%v", e))
	}
	select {
	case <-ctx.Done():
		return output(out, errs, *jsonMode, root, op, nil, fail("canceled", 2, "%v", ctx.Err()))
	default:
	}
	p, e := Preview(o)
	if e == nil && *dest != "" {
		e = SavePlan(p, *dest)
	}
	if e == nil && !p.CanApply {
		e = fail("conflict", 3, "preview contains %d conflicts", len(p.Conflicts))
	}
	return output(out, errs, *jsonMode, p.ProjectRoot, op, &p, e)
}
func output(out, errs io.Writer, jsonMode bool, root, op string, result any, e error) int {
	status := "ok"
	exit := 0
	var problem *Error
	if e != nil {
		status = "error"
		var ok bool
		problem, ok = e.(*Error)
		if !ok {
			problem = &Error{"io", e.Error(), 4}
		}
		exit = problem.Exit
	}
	if jsonMode {
		envelope := map[string]any{"schemaVersion": Protocol, "operation": op, "projectRoot": root, "status": status}
		if result != nil {
			envelope["result"] = result
		}
		if problem != nil {
			envelope["error"] = problem
		}
		if err := json.NewEncoder(out).Encode(envelope); err != nil {
			fmt.Fprintln(errs, err)
			return 4
		}
	} else {
		if p, ok := result.(*Plan); ok {
			fmt.Fprintf(out, "%s preview: %d changes, %d preserved, %d conflicts; plan %s\n", op, len(p.Actions), len(p.Preserved), len(p.Conflicts), p.PlanDigest)
			for _, c := range p.Conflicts {
				fmt.Fprintln(out, "Conflict:", c)
			}
			for _, w := range p.Warnings {
				fmt.Fprintln(out, "Warning:", w)
			}
		}
		if problem != nil {
			fmt.Fprintln(errs, problem.Error())
		}
	}
	return exit
}
