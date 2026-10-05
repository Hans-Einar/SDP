package sdptool

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SDPTool/model"
	"github.com/Hans-Einar/SDP/SDPTool/presentation"
	"io"
	"strings"
)

func modelCommand(area string, args []string, out, errs io.Writer, jsonMode bool) int {
	var r model.Result
	var e error
	bad := func() int {
		return reportMode(errs, failure("arguments", fmt.Errorf("model: create work:NAME [--initial|from release:VERSION], status kind:NAME")), jsonMode)
	}
	if len(args) == 0 {
		return bad()
	}
	switch args[0] {
	case "merge":
		if len(args) != 4 || args[2] != "into" {
			return bad()
		}
		r, e = model.Merge(area, args[1], args[3], "")
	case "commit":
		if len(args) != 4 && len(args) != 5 {
			return bad()
		}
		if args[2] != "--message" {
			return bad()
		}
		resolved := len(args) == 5 && args[4] == "--resolved"
		if len(args) == 5 && !resolved {
			return bad()
		}
		r, e = model.Commit(area, args[1], args[3], resolved)
	case "restore":
		if len(args) != 4 || args[2] != "to" || !strings.HasPrefix(args[3], "commit:") {
			return bad()
		}
		r, e = model.Restore(area, args[1], strings.TrimPrefix(args[3], "commit:"))
	case "recover":
		if len(args) != 3 {
			return bad()
		}
		r, e = model.Recover(area, args[1], args[2])
	case "status", "history":
		if len(args) != 2 {
			return bad()
		}
		r, e = model.Status(area, args[1])
		r.Operation = args[0]
	case "create":
		if len(args) < 2 {
			return bad()
		}
		target := strings.SplitN(args[1], ":", 2)
		if len(target) != 2 {
			return bad()
		}
		if target[0] != "work" {
			return bad()
		}
		initial := false
		from := ""
		switch {
		case len(args) == 2:
		case len(args) == 3 && args[2] == "--initial":
			initial = true
		case len(args) == 5 && args[2] == "from":
			r, e = model.Merge(area, args[4], args[3], target[1])
			break
		case len(args) == 4 && args[2] == "from":
			from = args[3]
		default:
			return bad()
		}
		if len(args) != 5 {
			r, e = model.CreateWork(area, target[1], from, initial)
		}
	default:
		return bad()
	}
	if e != nil {
		if me, ok := e.(*model.Error); ok {
			return reportMode(errs, failure(me.Code, fmt.Errorf("%s", me.Detail)), jsonMode)
		}
		return reportMode(errs, e, jsonMode)
	}
	if e = presentation.Default().Write(out, r, jsonMode); e != nil {
		return reportMode(errs, e, jsonMode)
	}
	return 0
}
