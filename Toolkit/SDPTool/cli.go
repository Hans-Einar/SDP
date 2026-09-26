package sdptool

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

func Run(ctx context.Context, args []string, out, errs io.Writer) int {
	if len(args) == 1 && args[0] == "--version" {
		e := json.NewEncoder(out).Encode(map[string]any{
			"schema": Version, "operation": "version", "version": BuildVersion,
			"revision": BuildRevision, "installedFactSchemas": []string{"1.0", "2.0"},
		})
		if e != nil {
			return report(errs, e)
		}
		return 0
	}

	if len(args) == 1 && (args[0] == "--help" || args[0] == "help") {
		_, e := io.WriteString(out, "Usage: sdptool [PROJECT-OR-SDP-AREA] discover|tree|select|view ip|sdui-preview [options]\n       sdptool preview FILE --output DIRECTORY [--renderer PROGRAM]\nSee Toolkit/SDPTool/Contract.md for source, model, revision and resource contracts.\n")
		if e != nil {
			return report(errs, e)
		}
		return 0
	}
	selected := "."
	if len(args) > 1 && !isCommand(args[0]) {
		selected = args[0]
		args = args[1:]
	}
	if len(args) > 0 && args[0] == "discover" {
		if len(args) != 1 {
			return report(errs, failure("arguments", fmt.Errorf("discover takes no positional arguments")))
		}
		p, e := Discover(selected)
		if e != nil {
			return report(errs, e)
		}
		if e = json.NewEncoder(out).Encode(p); e != nil {
			return report(errs, e)
		}
		return 0
	}

	if len(args) > 0 && args[0] == "select" {
		fs := flag.NewFlagSet("select", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		var o PreviewOptions
		var model string
		fs.StringVar(&model, "model", "", "model ID")
		fs.StringVar(&o.Output, "output", "", "bundle directory")
		fs.StringVar(&o.URI, "uri", "", "selection URI")
		fs.StringVar(&o.Revision, "revision", "", "expected revision")
		fs.StringVar(&o.Renderer, "renderer", os.Getenv("SDP_MMDR"), "renderer")
		if e := fs.Parse(args[1:]); e != nil {
			return report(errs, failure("arguments", e))
		}
		if fs.NArg() != 0 {
			return report(errs, failure("arguments", fmt.Errorf("unexpected arguments")))
		}
		p, e := Discover(selected)
		if e != nil {
			return report(errs, e)
		}
		r, e := SelectProject(ctx, p, model, o)
		if e != nil {
			return report(errs, e)
		}
		if e = json.NewEncoder(out).Encode(r); e != nil {
			return report(errs, e)
		}
		return 0
	}

	if len(args) > 0 && args[0] == "sdui-preview" {
		fs := flag.NewFlagSet("sdui-preview", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		model := fs.String("model", "", "registered source")
		entry := fs.String("entry", "", "root frame")
		output := fs.String("output", "", "bundle directory")
		revision := fs.String("revision", "", "expected revision")
		if e := fs.Parse(args[1:]); e != nil {
			return report(errs, failure("arguments", e))
		}
		if fs.NArg() != 0 {
			return report(errs, failure("arguments", fmt.Errorf("unexpected arguments")))
		}
		p, e := Discover(selected)
		if e != nil {
			return report(errs, e)
		}
		r, e := UIPreview(ctx, p, *model, *entry, *output, *revision)
		if e != nil {
			return report(errs, e)
		}
		if e = json.NewEncoder(out).Encode(r); e != nil {
			return report(errs, e)
		}
		return 0
	}

	if len(args) > 0 && args[0] == "tree" {
		fs := flag.NewFlagSet("tree", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		model := fs.String("model", "", "registered model")
		if e := fs.Parse(args[1:]); e != nil {
			return report(errs, failure("arguments", e))
		}
		if fs.NArg() != 0 {
			return report(errs, failure("arguments", fmt.Errorf("unexpected arguments")))
		}
		p, e := Discover(selected)
		if e != nil {
			return report(errs, e)
		}
		t, e := Navigation(p, *model)
		if e != nil {
			return report(errs, e)
		}
		if e = json.NewEncoder(out).Encode(t); e != nil {
			return report(errs, e)
		}
		return 0
	}

	if len(args) > 0 && args[0] == "view" {
		if len(args) < 2 || (args[1] != "ip" && args[1] != "implementation-plan") {
			return report(errs, failure("arguments", fmt.Errorf("view requires ip or implementation-plan")))
		}
		fs := flag.NewFlagSet("view", flag.ContinueOnError)
		fs.SetOutput(io.Discard)
		h := Host{}
		var model string
		fs.StringVar(&model, "model", "", "registered model")
		fs.StringVar(&h.Viewer, "viewer", os.Getenv("SDP_XFMD"), "prebuilt viewer")
		fs.StringVar(&h.SDLTool, "sdl-tool", os.Getenv("SDP_SDL_TOOL"), "prebuilt SDL tool")
		fs.StringVar(&h.Renderer, "renderer", os.Getenv("SDP_MMDR"), "prebuilt renderer")
		if e := fs.Parse(args[2:]); e != nil {
			return report(errs, failure("arguments", e))
		}
		if fs.NArg() != 0 {
			return report(errs, failure("arguments", fmt.Errorf("unexpected positional arguments")))
		}
		p, e := Discover(selected)
		if e != nil {
			return report(errs, e)
		}
		if e = ViewPlan(ctx, p, model, h); e != nil {
			return report(errs, e)
		}
		if e = json.NewEncoder(out).Encode(map[string]any{"schema": Version, "operation": "view", "status": "viewer-exited"}); e != nil {
			return report(errs, e)
		}
		return 0
	}

	if len(args) < 2 || args[0] != "preview" {
		return report(errs, failure("arguments", fmt.Errorf("usage: sdptool preview FILE --output DIRECTORY [--viewpoint VPnn | --uri URI] [--renderer PROGRAM] [--revision HASH]")))
	}
	fs := flag.NewFlagSet("preview", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	o := PreviewOptions{Source: args[1]}
	fs.StringVar(&o.Output, "output", "", "bundle directory")
	fs.StringVar(&o.Renderer, "renderer", "", "prebuilt renderer")
	fs.StringVar(&o.Viewpoint, "viewpoint", "", "viewpoint")
	fs.StringVar(&o.URI, "uri", "", "typed selection")
	fs.StringVar(&o.Revision, "revision", "", "expected source revision")
	if e := fs.Parse(args[2:]); e != nil {
		return report(errs, failure("arguments", e))
	}
	if fs.NArg() != 0 {
		return report(errs, failure("arguments", fmt.Errorf("unexpected positional arguments")))
	}
	r, e := Preview(ctx, o)
	if e != nil {
		return report(errs, e)
	}
	if e = json.NewEncoder(out).Encode(r); e != nil {
		return report(errs, e)
	}
	return 0
}
func report(w io.Writer, e error) int {
	f, ok := e.(*Failure)
	if !ok {
		f = &Failure{Code: "io", Message: e.Error()}
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"schema": Version, "error": f})
	return 1
}

func isCommand(s string) bool {
	switch s {
	case "preview", "discover", "view", "tree", "select", "sdui-preview":
		return true
	}
	return false
}
