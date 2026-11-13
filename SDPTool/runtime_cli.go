package sdptool

import (
	"flag"
	"fmt"
	"io"
)

func cliSetup(c cliInvocation) int {
	if c.explicitProject || len(c.args) != 1 {
		return c.report(c.errs, failure("arguments", fmt.Errorf("usage: sdptool setup [--json]; no project selection")))
	}
	r, _, err := performSetup(c.ctx, actionContext{}, noActionParameters{})
	if err != nil {
		return c.report(c.errs, err)
	}
	if err = c.emit(r); err != nil {
		return c.report(c.errs, err)
	}
	return 0
}
func cliSdui(c cliInvocation) int {
	if c.explicitProject || len(c.args) < 2 {
		return c.report(c.errs, failure("arguments", fmt.Errorf("usage: sdptool sdui preview|check|run [options]; no project selection")))
	}
	fs := flag.NewFlagSet("sdui", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	source := fs.String("source", "", "SDUI source")
	entry := fs.String("entry", "page", "Root frame")
	revision := fs.String("revision", "", "Expected source SHA256")
	output := fs.String("output", "", "Existing empty bundle directory")
	watch := fs.Bool("watch", true, "Watch source while prototype runs")
	if err := fs.Parse(c.args[2:]); err != nil {
		return c.report(c.errs, err)
	}
	var r any
	var err error
	switch c.args[1] {
	case "preview", "check":
		if fs.NArg() != 0 {
			return c.report(c.errs, fmt.Errorf("preview/check accept no positional arguments"))
		}
		r, _, err = performSource(c.ctx, sduiSourceParameters{Source: *source, Entry: *entry, Revision: *revision, Output: *output}, c.args[1] == "check")
	case "run":
		if fs.NArg() != 1 || *source != "" || *output != "" {
			return c.report(c.errs, fmt.Errorf("usage: sdptool sdui run [-entry page] [-revision SHA] [-watch=true] SOURCE"))
		}
		r, _, err = performSduiRun(c.ctx, actionContext{}, sduiRunParameters{Source: fs.Arg(0), Entry: *entry, Revision: *revision, Watch: watch})
	default:
		return c.report(c.errs, fmt.Errorf("unknown sdui operation %q", c.args[1]))
	}
	if err != nil {
		return c.report(c.errs, err)
	}
	if err = c.emit(r); err != nil {
		return c.report(c.errs, err)
	}
	return 0
}
