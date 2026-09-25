package sdptool

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
)

func Run(ctx context.Context, args []string, out, errs io.Writer) int {
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
