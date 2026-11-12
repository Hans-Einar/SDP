package sdptool

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/Hans-Einar/SDP/SDPTool/activities"
	"github.com/Hans-Einar/SDP/SDPTool/presentation"
)

func activityCommand(ctx context.Context, selected string, args []string, out, errs io.Writer, jsonMode bool) int {
	bad := func(e error) int { return reportMode(errs, failure("activity", e), jsonMode) }
	if len(args) > 1 && args[0] == "session" && args[1] == "create" {
		args = append([]string{"session-create"}, args[2:]...)
	}
	if len(args) == 0 {
		return bad(fmt.Errorf("activity requires an operation; use activity list --json"))
	}
	op, ok := activities.FindOperation(args[0])
	if !ok {
		return bad(fmt.Errorf("unknown Activity operation %q", args[0]))
	}
	fs := flag.NewFlagSet("activity "+op.Name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	request := fs.String("request", "", "JSON request file")
	id := fs.String("activity", "", "Activity identity")
	recovery := fs.String("operation", "", "retained operation identity")
	uri := fs.String("uri", "", "producer selection URI")
	revision := fs.String("revision", "", "expected source revision")
	output := fs.String("output", "", "caller-owned view bundle")
	source := fs.String("authority-source", "", "existing owner decision document, relative to SDP")
	if e := fs.Parse(args[1:]); e != nil {
		return bad(e)
	}
	if fs.NArg() != 0 {
		return bad(fmt.Errorf("unexpected Activity arguments"))
	}
	if op.Write && *source == "" {
		return bad(fmt.Errorf("Activity writes require --authority-source PATH to the existing owner decision document under SDP"))
	}
	if (op.Input == activities.RequestInput) != (*request != "") || (op.Input == activities.ActivityInput) != (*id != "") || (op.Input == activities.RecoveryInput) != (*recovery != "") {
		return bad(fmt.Errorf("%s expects input kind %s", op.Name, op.Input))
	}
	s, e := activityService(selected)
	if e != nil {
		return bad(e)
	}
	if op.Input == activities.ViewInput || op.Input == activities.TreeInput {
		if *uri == "" {
			return bad(fmt.Errorf("selection URI required"))
		}
		if op.Input == activities.TreeInput && (*output != "" || *revision == "") {
			return bad(fmt.Errorf("tree requires revision and no output"))
		}
	} else if *uri != "" || *revision != "" || *output != "" {
		return bad(fmt.Errorf("unexpected view/tree flags"))
	}
	in := activities.Input{ActivityID: *id, OperationID: *recovery, URI: *uri, Revision: *revision, Output: *output}
	if *request != "" {
		in.Request, e = boundedFile(*request)
		if e != nil {
			return bad(e)
		}
	}
	r, e := op.Execute(ctx, s, activities.Caller{Actor: "local-cli", Authority: activities.Authority{Role: "project-owner", Source: *source}}, in, activityOperationContext(ctx, selected, s))
	if e != nil {
		return bad(e)
	}
	var result any = r
	if op.Input == activities.ViewInput || op.Input == activities.TreeInput {
		result = r.Data
	}
	if e = presentation.Default().Write(out, result, jsonMode); e != nil {
		return bad(e)
	}
	return 0
}
