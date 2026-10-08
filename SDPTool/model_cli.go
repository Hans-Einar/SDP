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
		return reportMode(errs, failure("arguments", fmt.Errorf("invalid model command; use sdptool model help")), jsonMode)
	}
	if len(args) == 1 && (args[0] == "help" || args[0] == "--help") {
		_, err := io.WriteString(out, modelHelp)
		if err != nil {
			return reportMode(errs, err, jsonMode)
		}
		return 0
	}
	if len(args) == 0 {
		return bad()
	}
	switch args[0] {
	case "snapshot":
		if len(args) != 2 {
			return bad()
		}
		view, err := model.Snapshot(area, args[1])
		if err != nil {
			return reportMode(errs, err, jsonMode)
		}
		if err = presentation.Default().Write(out, view, jsonMode); err != nil {
			return reportMode(errs, err, jsonMode)
		}
		return 0
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
			if len(args) < 4 || args[2] != "from" {
				return bad()
			}
			evidence := ""
			code := ""
			checks := ""
			if (len(args)-4)%2 != 0 {
				return bad()
			}
			seen := map[string]bool{}
			for i := 4; i < len(args); i += 2 {
				if seen[args[i]] {
					return bad()
				}
				seen[args[i]] = true
				switch args[i] {
				case "--evidence":
					evidence = args[i+1]
				case "--code-digest":
					code = args[i+1]
				case "--checks":
					checks = args[i+1]
				default:
					return bad()
				}
			}
			if evidence == "verified" {
				if code == "" || checks == "" {
					return bad()
				}
				evidence = "verified:" + code + ":" + checks
			} else if code != "" || checks != "" {
				return bad()
			}
			r, e = model.Freeze(area, target[0], target[1], args[3], evidence)
			break
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

const modelHelp = `Usage: sdptool [MODEL-AREA] model ACTION [--json]
  create work:NAME [--initial | from release:VERSION]
  create work:COMBINED from work:A work:B
  commit work:NAME --message TEXT [--resolved]
  restore work:NAME to commit:00001
  merge work:SOURCE into work:TARGET
  create candidate:NAME from work:NAME
  create proposal:NAME from work:NAME
  create release:VERSION from candidate:NAME --evidence model-only
    or --evidence verified --code-digest SHA256 --checks REFERENCE
  status|history|snapshot kind:NAME
  recover OPERATION-UUID resume|abort
MODEL-AREA is the directory containing artifact folders; no Git or SDP install required.
WORK snapshots are preliminary. Blueprint generation is not implemented.
`
