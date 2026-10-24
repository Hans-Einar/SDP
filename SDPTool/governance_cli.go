package sdptool

import (
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/presentation"
)

const governanceHelp = `Usage: sdptool [PROJECT] governance context|route|bind|start|read|transition [options]
  --input FILE   bounded JSON Selection, RouteRequest, StartRequest or TransitionRequest
  --run UUID     run to read (read only)
  --instance UUID  Session instance for session-read
Session operations: session-definition, session-context, session-adopt, session-read,
session-list, session-select, session-assess, session-start, session-progress, session-verify, session-review, session-resume. Selection prepares checked ready-card effects; resume commits them.
These operations do not launch agents or complete work.
  --state DIR    private operational state outside project (defaults to XDG state home)
  --json         machine output; human-readable output is the default
context and route only read; bind is explicit. start/transition are trusted local
controller operations, not agent/owner approval endpoints. PGI1 cannot complete
work, publish records or launch agents. See SDPTool/governance/README.md.
`

func governanceCommand(selected string, args []string, out, errs io.Writer, jsonMode bool) int {
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
	if *instance != "" && op != "session-read" {
		return report(fmt.Errorf("--instance is only valid for session-read"))
	}
	switch op {
	case "session-definition", "session-list":
		if *input != "" || *run != "" {
			return report(fmt.Errorf("operation takes no request/run"))
		}
	case "session-read":
		if *input != "" || *run != "" || *instance == "" {
			return report(fmt.Errorf("session-read requires --instance only"))
		}
	case "session-context", "session-adopt", "session-resume", "session-select", "session-assess", "session-start", "session-progress", "session-verify", "session-review":
		if *input == "" || *run != "" {
			return report(fmt.Errorf("operation requires --input only"))
		}
	case "bind":
		if *input != "" || *run != "" {
			return report(fmt.Errorf("bind takes no request/run"))
		}
	case "read":
		if *input != "" || *run == "" {
			return report(fmt.Errorf("read requires --run only"))
		}
	case "context", "route", "start", "transition":
		if *input == "" || *run != "" {
			return report(fmt.Errorf("operation requires --input only"))
		}
	default:
		return report(fmt.Errorf("unknown governance command"))
	}
	p, e := Discover(selected)
	if e != nil {
		return report(e)
	}
	if *state == "" {
		base := os.Getenv("XDG_STATE_HOME")
		if base == "" {
			home, e := os.UserHomeDir()
			if e != nil {
				return report(e)
			}
			base = filepath.Join(home, ".local", "state")
		}
		if !filepath.IsAbs(base) {
			return report(fmt.Errorf("XDG_STATE_HOME must be absolute"))
		}
		*state = filepath.Join(base, "sdptool", "project-governance", fmt.Sprintf("%x", sha256.Sum256([]byte(p.Root))))
	}
	service, e := governance.New(p.Root, p.Area, p.Inventory.ProjectID, *state)
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
	var result any
	switch op {
	case "session-definition":
		w := governance.SessionWorkflow()
		result = map[string]any{"schema": governance.SessionSchema, "operation": op, "definition": w, "digest": w.Digest()}
	case "session-context":
		var q governance.SessionContextRequest
		if e = governance.Decode(b, &q); e == nil {
			result, e = service.SessionContext(q)
		}
	case "session-adopt":
		var q governance.SessionAdoptRequest
		if e = governance.Decode(b, &q); e == nil {
			result, e = service.AdoptSession(governance.LocalController(), q)
		}
	case "session-resume":
		var q governance.SessionResumeRequest
		if e = governance.Decode(b, &q); e == nil {
			result, e = service.ResumeSession(governance.LocalController(), q)
		}
	case "session-select":
		var q governance.SessionSelectRequest
		if e = governance.Decode(b, &q); e == nil {
			result, e = service.SelectSession(governance.LocalController(), q)
		}
	case "session-assess":
		var q governance.SessionAssessRequest
		if e = governance.Decode(b, &q); e == nil {
			result, e = service.AssessSession(governance.LocalController(), q)
		}
	case "session-start", "session-progress", "session-verify", "session-review":
		kind := strings.TrimPrefix(op, "session-")
		var q governance.SessionWorkRequest
		if q, e = governance.DecodeSessionWork(b, kind); e == nil {
			if kind == "start" {
				result, e = service.StartSession(governance.LocalController(), q)
			} else if kind == "progress" {
				result, e = service.ProgressSession(governance.LocalController(), q)
			} else if kind == "verify" {
				result, e = service.VerifySession(governance.LocalController(), q)
			} else {
				result, e = service.ReviewSession(governance.LocalController(), q)
			}
		}
	case "session-read":
		result, e = service.ReadSession(*instance)
	case "session-list":
		result, e = service.ListSessions()
	case "context":
		var q governance.Selection
		if e = governance.Decode(b, &q); e == nil {
			result, e = service.Context(q)
		}
	case "route":
		var q governance.RouteRequest
		if e = governance.Decode(b, &q); e == nil {
			result, e = service.Route(q)
		}
	case "bind":
		result, e = service.Bind(governance.LocalController())
	case "start":
		var q governance.StartRequest
		if e = governance.Decode(b, &q); e == nil {
			result, e = service.Start(governance.LocalController(), q)
		}
	case "read":
		result, e = service.Read(*run)
	case "transition":
		var q governance.TransitionRequest
		if e = governance.Decode(b, &q); e == nil {
			result, e = service.Transition(governance.LocalController(), q)
		}
	}
	if e != nil {
		return report(e)
	}
	if e = presentation.Default().Write(out, result, jsonMode); e != nil {
		return report(e)
	}
	return 0
}
