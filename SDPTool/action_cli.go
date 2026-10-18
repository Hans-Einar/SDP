package sdptool

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool/blueprintstate"
)

func invocationCommand(ctx context.Context, args []string, in io.Reader, out, errs io.Writer) int {
	bad := func(e error) int {
		return writeInvocation(out, errs, invokeFailure(InvokeResponse{Schema: InvokeSchema}, 2, "invalid-request", e))
	}
	fs := flag.NewFlagSet("invoke", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	file := fs.String("request", "", "one request file, or - for stdin")
	actor := fs.String("as", "", "local assignment actor")
	authority := fs.String("authority", "", "local assignment authority")
	jsonMode := fs.Bool("json", true, "JSON response (required)")
	if err := fs.Parse(args); err != nil {
		return bad(err)
	}
	if fs.NArg() != 0 || *file == "" || !*jsonMode {
		return bad(fmt.Errorf("invoke requires --request FILE|- and JSON output"))
	}
	if (*actor == "") != (*authority == "") {
		return bad(fmt.Errorf("--as and --authority must be supplied together"))
	}
	if in == nil {
		in = strings.NewReader("")
	}
	if *file != "-" {
		f, err := os.Open(*file)
		if err != nil {
			return bad(err)
		}
		defer f.Close()
		in = f
	}
	return writeInvocation(out, errs, Invoke(ctx, in, ActionCaller{Principal: blueprintstate.Principal{Actor: *actor, Authority: *authority}}))
}
func writeActions(out io.Writer, r ActionCatalogue, jsonMode bool) error {
	if jsonMode {
		return json.NewEncoder(out).Encode(r)
	}
	if _, err := fmt.Fprintf(out, "SDPTool actions (%s, %s)\n", r.Producer.Version, r.Schema); err != nil {
		return err
	}
	for _, a := range r.Actions {
		if _, err := fmt.Fprintf(out, "%s — %s [%s]\n  %s\n", a.ID, a.Label, a.Availability.State, a.Availability.Message); err != nil {
			return err
		}
	}
	return nil
}
