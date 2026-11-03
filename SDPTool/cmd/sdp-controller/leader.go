package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancecontroller"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

func leaderCommand(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("leader inspect|message|run --project ROOT --state DIR")
	}
	action, args := args[0], args[1:]
	if action != "inspect" && action != "message" && action != "run" {
		return fmt.Errorf("unknown leader action")
	}
	f := flag.NewFlagSet("leader "+action, flag.ContinueOnError)
	root := f.String("project", "", "explicit project root")
	state := f.String("state", "", "private operational state directory")
	var owner, operation, input, coordinator, mcp, binary string
	var skills, sessions []string
	var executeSessions bool
	var workerSandbox string
	var maximum int
	var timeout, interval time.Duration
	if action == "message" {
		f.StringVar(&owner, "owner", "", "trusted local owner attribution")
		f.StringVar(&operation, "operation", "", "stable enqueue UUID; reuse for retry")
		f.StringVar(&input, "input", "", "UTF-8 message file (up to 4096 bytes)")
	}
	if action == "run" {
		f.StringVar(&coordinator, "coordinator", "", "explicit coordinator identity")
		f.StringVar(&mcp, "mcp-binary", "", "sdp-mcp executable")
		f.StringVar(&binary, "codex", "codex", "Codex executable (managed account)")
		f.Func("skill", "explicit project-relative SKILL.md; repeat up to eight", func(v string) error { skills = append(skills, v); return nil })
		f.BoolVar(&executeSessions, "execute-session-requests", false, "consume granted Session launch/continue requests with up to two separate worker controllers")
		f.StringVar(&workerSandbox, "worker-sandbox", "read-only", "sandbox for explicitly enabled Session execution: read-only or workspace-write")
		f.Func("session", "registered Session UUID to grant workflow and host-request tools; repeat up to eight", func(v string) error { sessions = append(sessions, v); return nil })
		f.IntVar(&maximum, "max-turns", 1, "terminal leader turns before exit (1..16)")
		f.DurationVar(&timeout, "timeout", 5*time.Minute, "maximum controller lifetime (up to 30m)")
		f.DurationVar(&interval, "interval", time.Second, "poll interval (at least 100ms)")
	}
	if e := f.Parse(args); e != nil {
		return e
	}
	if *root == "" || *state == "" || f.NArg() != 0 {
		return fmt.Errorf("explicit project and state required")
	}
	if action == "message" && (owner == "" || operation == "" || input == "") {
		return fmt.Errorf("message requires --owner --operation --input")
	}
	if action == "run" && (coordinator == "" || mcp == "" || len(skills) == 0 || len(skills) > 8 || len(sessions) > 8 || maximum < 1 || maximum > 16 || timeout <= 0 || timeout > 30*time.Minute || interval < 100*time.Millisecond) {
		return fmt.Errorf("run requires coordinator, MCP executable, explicit skills and bounded turns/time")
	}
	if action == "run" && ((executeSessions && len(sessions) == 0) || (workerSandbox != "read-only" && workerSandbox != "workspace-write") || (!executeSessions && workerSandbox != "read-only")) {
		return fmt.Errorf("Session execution requires explicit Session grants and supported worker sandbox")
	}
	p, e := sdptool.Discover(*root)
	if e != nil {
		return e
	}
	s, e := governance.New(p.Root, p.Area, p.Inventory.ProjectID, *state)
	if e != nil {
		return e
	}
	encoder := json.NewEncoder(out)
	if action == "inspect" {
		v, e := s.LeaderRead()
		if e != nil {
			return e
		}
		return encoder.Encode(v)
	}
	if action == "message" {
		file, e := os.Open(input)
		if e != nil {
			return e
		}
		defer file.Close()
		b, e := io.ReadAll(io.LimitReader(file, 4097))
		if e != nil {
			return e
		}
		// Explicit message enqueue may establish the checkout binding, never a run.
		if _, e = s.Bind(governance.LocalController()); e != nil {
			return e
		}
		result, e := governancecontroller.LeaderMessage(s, operation, owner, "sdp-controller leader message", string(b))
		if e != nil {
			return e
		}
		return encoder.Encode(result)
	}
	mcp, e = filepath.Abs(mcp)
	if e != nil {
		return e
	}
	statePath, e := filepath.Abs(*state)
	if e != nil {
		return e
	}
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ctx, cancel := context.WithTimeout(signalCtx, timeout)
	defer cancel()
	host, e := governancehost.Open(ctx, binary, p.Root)
	if e != nil {
		return e
	}
	defer host.Close()
	c := governancecontroller.LeaderController{SessionIDs: sessions, Service: s, Host: host, Root: p.Root, StateDir: statePath, MCPBinary: mcp, Coordinator: coordinator, SkillPaths: skills}
	v, e := s.LeaderRead()
	if e != nil {
		return e
	}
	baseline := leaderTerminalCount(v)
	dispatcher := sessionDispatcher{service: s, sessions: sessions, root: p.Root, state: statePath, codex: binary, mcp: mcp, sandbox: workerSandbox, attempted: map[string]bool{}, active: map[string]sessionProcess{}}
	revision := ^uint64(0)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		if leaderTerminalCount(v)-baseline < maximum {
			v, e = c.Step(ctx)
			if e != nil {
				return e
			}
		}
		if executeSessions {
			if e := dispatcher.tick(ctx, encoder); e != nil {
				return e
			}
		}
		if v.Leader.Revision != revision {
			if e = encoder.Encode(v); e != nil {
				return e
			}
			revision = v.Leader.Revision
		}
		if leaderTerminalCount(v)-baseline >= maximum && len(dispatcher.active) == 0 {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
func leaderTerminalCount(v governance.LeaderView) int {
	n := 0
	for _, d := range v.Leader.Dispatches {
		if d.State == "completed" || d.State == "failed" || d.State == "interrupted" {
			n++
		}
	}
	return n
}
