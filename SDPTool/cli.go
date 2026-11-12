package sdptool

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool/install"
	"github.com/Hans-Einar/SDP/SDPTool/presentation"
)

// Run retains the legacy argv/output API; stdin is explicit in RunWithInput.
func Run(ctx context.Context, args []string, out, errs io.Writer) int {
	return RunWithInput(ctx, args, strings.NewReader(""), out, errs)
}
func RunWithInput(ctx context.Context, args []string, in io.Reader, out, errs io.Writer) int {
	if len(args) > 0 {
		if command, ok := cliCommands()[args[0]]; ok && command.early != nil {
			return command.early(cliInvocation{ctx: ctx, args: args, in: in, out: out, errs: errs})
		}
	}
	args, jsonMode, formatErr := outputArguments(args)
	report := func(w io.Writer, e error) int { return reportMode(w, e, jsonMode) }
	emit := func(v any) error { return presentation.Default().Write(out, v, jsonMode) }
	if formatErr != nil {
		selected, rest := ".", args
		if len(rest) > 1 && !isCommand(rest[0]) {
			selected, rest = rest[0], rest[1:]
		}
		if len(rest) > 0 {
			if command, ok := cliCommands()[rest[0]]; ok && command.argumentError != nil {
				return command.argumentError(cliInvocation{ctx: ctx, selected: selected, args: rest, in: in, out: out, errs: errs, jsonMode: jsonMode}, formatErr)
			}
		}
		return report(errs, failure("arguments", formatErr))
	}
	if len(args) == 1 && args[0] == "--version" {
		e := emit(map[string]any{
			"schema": Version, "operation": "version", "version": BuildVersion,
			"revision": BuildRevision, "installedFactSchemas": []string{"1.0", "2.0", "3.0"}, "installationProtocol": install.Protocol,
		})
		if e != nil {
			return report(errs, e)
		}
		return 0
	}

	if len(args) == 1 && (args[0] == "--help" || args[0] == "help") {
		_, e := io.WriteString(out, "Usage: sdptool [PROJECT-OR-SDP-AREA] actions|invoke|activity|governance|model|discover|run|tree|select|view ip|sdui-preview|install|upgrade|release-log [options]\n       sdptool preview FILE --output DIRECTORY [--renderer PROGRAM]\nDefault output is human-readable; add --json for machine clients.\nSee SDPTool/Contract.md for source, model, revision and resource contracts.\n")
		if e != nil {
			return report(errs, e)
		}
		return 0
	}
	selected := "."
	explicitProject := false
	if len(args) > 1 && !isCommand(args[0]) {
		selected = args[0]
		explicitProject = true
		args = args[1:]
	}
	if len(args) > 0 {
		if command, ok := cliCommands()[args[0]]; ok {
			return command.run(cliInvocation{ctx: ctx, selected: selected, explicitProject: explicitProject, args: args, in: in, out: out, errs: errs, jsonMode: jsonMode})
		}
	}
	return report(errs, failure("arguments", fmt.Errorf("usage: sdptool preview FILE --output DIRECTORY [--viewpoint VPnn | --uri URI] [--renderer PROGRAM] [--revision HASH]")))
}

func reportMode(w io.Writer, e error, jsonMode bool) int {
	f, ok := e.(*Failure)
	if !ok {
		f = &Failure{Code: "io", Message: e.Error()}
	}
	_ = presentation.Default().Write(w, map[string]any{"schema": Version, "error": f}, jsonMode)
	return 1
}

type cliInvocation struct {
	ctx             context.Context
	selected        string
	explicitProject bool
	args            []string
	in              io.Reader
	out, errs       io.Writer
	jsonMode        bool
}

func (c cliInvocation) report(w io.Writer, e error) int { return reportMode(w, e, c.jsonMode) }
func (c cliInvocation) emit(v any) error                { return presentation.Default().Write(c.out, v, c.jsonMode) }

type cliCommand struct {
	run           func(cliInvocation) int
	early         func(cliInvocation) int
	argumentError func(cliInvocation, error) int
}

// cliCommands is the single top-level command inventory, including its exceptional
// pre-output parsing and installer error paths.
func cliCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"governance":   {run: cliGovernance},
		"activity":     {run: cliActivity},
		"invoke":       {run: cliInvoke, early: func(c cliInvocation) int { return invocationCommand(c.ctx, c.args[1:], c.in, c.out, c.errs) }},
		"actions":      {run: cliActions},
		"model":        {run: cliModel},
		"run":          {run: cliRun},
		"release-log":  {run: cliReleaseLog},
		"install":      {run: cliInstall, argumentError: cliInstallArgumentError},
		"upgrade":      {run: cliInstall, argumentError: cliInstallArgumentError},
		"discover":     {run: cliDiscover},
		"select":       {run: cliSelect},
		"sdui-preview": {run: cliSDUIPreview},
		"tree":         {run: cliTree},
		"view":         {run: cliView},
		"preview":      {run: cliPreview},
	}
}
func isCommand(s string) bool { _, ok := cliCommands()[s]; return ok }
func cliInstallArgumentError(c cliInvocation, e error) int {
	return install.ArgumentError(c.selected, c.args[0], c.jsonMode, c.out, c.errs, e)
}

func cliGovernance(c cliInvocation) int {
	return governanceCommand(c.selected, c.args[1:], c.out, c.errs, c.jsonMode)
}

func cliActivity(c cliInvocation) int {
	return activityCommand(c.ctx, c.selected, c.args[1:], c.out, c.errs, c.jsonMode)
}

func cliInvoke(c cliInvocation) int {
	return writeInvocation(c.out, c.errs, invokeFailure(InvokeResponse{Schema: InvokeSchema}, 2, "invalid-request", fmt.Errorf("invoke must be the first argument; project belongs in request context")))
}

func cliActions(c cliInvocation) int {
	if len(c.args) != 1 {
		return c.report(c.errs, failure("arguments", fmt.Errorf("actions takes only a project selection and --json")))
	}
	catalogueContext := ""
	if c.explicitProject {
		catalogueContext = c.selected
	}
	result, e := Actions(c.ctx, catalogueContext)
	if e != nil {
		return c.report(c.errs, e)
	}
	if e = writeActions(c.out, result, c.jsonMode); e != nil {
		return c.report(c.errs, e)
	}
	if len(result.Diagnostics) > 0 {
		return 1
	}
	return 0
}

func cliModel(c cliInvocation) int {
	return modelCommand(c.ctx, c.selected, c.args[1:], c.out, c.errs, c.jsonMode)
}

func cliRun(c cliInvocation) int {
	return runProgramCommand(c.ctx, c.selected, c.args[1:], c.out, c.errs, c.jsonMode)
}

func cliReleaseLog(c cliInvocation) int {
	return releaseLog(c.selected, c.args[1:], c.out, c.errs, c.jsonMode)
}

func cliDiscover(c cliInvocation) int {
	fs := flag.NewFlagSet("discover", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	navigation := fs.String("navigation", "", "eager or lazy-sdl navigation")
	if e := fs.Parse(c.args[1:]); e != nil {
		return c.report(c.errs, failure("arguments", e))
	}
	if fs.NArg() != 0 {
		return c.report(c.errs, failure("arguments", fmt.Errorf("discover takes no positional arguments")))
	}
	value, _, e := performDiscover(c.ctx, actionContext{Root: c.selected}, DiscoverOptions{Navigation: *navigation})
	p, _ := value.(Project)
	if e != nil {
		return c.report(c.errs, e)
	}
	if e = c.emit(p); e != nil {
		return c.report(c.errs, e)
	}
	return 0
}

func cliSelect(c cliInvocation) int {
	fs := flag.NewFlagSet("select", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var o PreviewOptions
	var model string
	fs.StringVar(&model, "model", "", "model ID")
	fs.StringVar(&o.Output, "output", "", "bundle directory")
	fs.StringVar(&o.URI, "uri", "", "selection URI")
	fs.StringVar(&o.Revision, "revision", "", "expected revision")
	fs.StringVar(&o.Renderer, "renderer", os.Getenv("SDP_MMDR"), "renderer")
	if e := fs.Parse(c.args[1:]); e != nil {
		return c.report(c.errs, failure("arguments", e))
	}
	if fs.NArg() != 0 {
		return c.report(c.errs, failure("arguments", fmt.Errorf("unexpected arguments")))
	}
	p, e := discover(c.selected, false)
	if e != nil {
		return c.report(c.errs, e)
	}
	r, e := SelectProject(c.ctx, p, model, o)
	if e != nil {
		return c.report(c.errs, e)
	}
	if e = c.emit(r); e != nil {
		return c.report(c.errs, e)
	}
	return 0
}

func cliSDUIPreview(c cliInvocation) int {
	fs := flag.NewFlagSet("sdui-preview", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	model := fs.String("model", "", "discovered source")
	entry := fs.String("entry", "", "root frame")
	output := fs.String("output", "", "bundle directory")
	revision := fs.String("revision", "", "expected revision")
	if e := fs.Parse(c.args[1:]); e != nil {
		return c.report(c.errs, failure("arguments", e))
	}
	if fs.NArg() != 0 {
		return c.report(c.errs, failure("arguments", fmt.Errorf("unexpected arguments")))
	}
	p, e := discover(c.selected, false)
	if e != nil {
		return c.report(c.errs, e)
	}
	r, e := UIPreview(c.ctx, p, *model, *entry, *output, *revision)
	if e != nil {
		return c.report(c.errs, e)
	}
	if e = c.emit(r); e != nil {
		return c.report(c.errs, e)
	}
	return 0
}

func cliTree(c cliInvocation) int {
	fs := flag.NewFlagSet("tree", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	model := fs.String("model", "", "discovered model")
	revision := fs.String("revision", "", "expected model revision")
	modelOnly := fs.Bool("model-only", false, "return only the revision-bound SDL model tree")
	if e := fs.Parse(c.args[1:]); e != nil {
		return c.report(c.errs, failure("arguments", e))
	}
	if fs.NArg() != 0 {
		return c.report(c.errs, failure("arguments", fmt.Errorf("unexpected arguments")))
	}
	if *modelOnly {
		t, e := ExpandModel(c.selected, ModelTreeOptions{Model: *model, Revision: *revision})
		if e != nil {
			return c.report(c.errs, e)
		}
		if e = c.emit(t); e != nil {
			return c.report(c.errs, e)
		}
		return 0
	}
	if *revision != "" {
		return c.report(c.errs, failure("arguments", fmt.Errorf("revision requires model-only")))
	}
	p, e := discover(c.selected, false)
	if e != nil {
		return c.report(c.errs, e)
	}
	t, e := Navigation(p, *model)
	if e != nil {
		return c.report(c.errs, e)
	}
	if e = c.emit(t); e != nil {
		return c.report(c.errs, e)
	}
	return 0
}

func cliView(c cliInvocation) int {
	if len(c.args) < 2 || (c.args[1] != "ip" && c.args[1] != "implementation-plan") {
		return c.report(c.errs, failure("arguments", fmt.Errorf("view requires ip or implementation-plan")))
	}
	fs := flag.NewFlagSet("view", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	h := Host{}
	var model string
	plan := fs.String("plan", "", "project-relative implementation plan")
	fs.StringVar(&model, "model", "", "discovered model")
	fs.StringVar(&h.Viewer, "viewer", os.Getenv("SDP_XFMD"), "prebuilt viewer")
	fs.StringVar(&h.SDLTool, "sdl-tool", os.Getenv("SDP_SDL_TOOL"), "prebuilt SDL tool")
	fs.StringVar(&h.Renderer, "renderer", os.Getenv("SDP_MMDR"), "prebuilt renderer")
	if e := fs.Parse(c.args[2:]); e != nil {
		return c.report(c.errs, failure("arguments", e))
	}
	if fs.NArg() != 0 {
		return c.report(c.errs, failure("arguments", fmt.Errorf("unexpected positional arguments")))
	}
	p, e := discover(c.selected, false)
	if e != nil {
		return c.report(c.errs, e)
	}
	if *plan != "" {
		p.Inventory.ImplementationPlan = *plan
	}
	if e = ViewPlan(c.ctx, p, model, h); e != nil {
		return c.report(c.errs, e)
	}
	if e = c.emit(map[string]any{"schema": Version, "operation": "view", "status": "viewer-exited"}); e != nil {
		return c.report(c.errs, e)
	}
	return 0
}

func cliInstall(c cliInvocation) int {
	installArgs := c.args[1:]
	if c.jsonMode {
		installArgs = append([]string{"--json"}, installArgs...)
	}
	return install.Run(c.ctx, c.selected, c.args[0], installArgs, c.out, c.errs)
}

func cliPreview(c cliInvocation) int {
	if len(c.args) < 2 {
		return c.report(c.errs, failure("arguments", fmt.Errorf("usage: sdptool preview FILE --output DIRECTORY [--viewpoint VPnn | --uri URI] [--renderer PROGRAM] [--revision HASH]")))
	}
	fs := flag.NewFlagSet("preview", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	o := PreviewOptions{Source: c.args[1]}
	fs.StringVar(&o.Output, "output", "", "bundle directory")
	fs.StringVar(&o.Renderer, "renderer", "", "prebuilt renderer")
	fs.StringVar(&o.Viewpoint, "viewpoint", "", "viewpoint")
	fs.StringVar(&o.URI, "uri", "", "typed selection")
	fs.StringVar(&o.Revision, "revision", "", "expected source revision")
	if e := fs.Parse(c.args[2:]); e != nil {
		return c.report(c.errs, failure("arguments", e))
	}
	if fs.NArg() != 0 {
		return c.report(c.errs, failure("arguments", fmt.Errorf("unexpected positional arguments")))
	}
	r, e := Preview(c.ctx, o)
	if e != nil {
		return c.report(c.errs, e)
	}
	if e = c.emit(r); e != nil {
		return c.report(c.errs, e)
	}
	return 0
}
