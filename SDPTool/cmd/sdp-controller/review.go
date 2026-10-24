package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancecontroller"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

func reviewCommand(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("review prepare|run|inspect|supersede --project ROOT --state DIR")
	}
	action := args[0]
	if action != "prepare" && action != "run" && action != "inspect" && action != "supersede" {
		return fmt.Errorf("unknown review action")
	}
	f := flag.NewFlagSet("review "+action, flag.ContinueOnError)
	root := f.String("project", "", "explicit project root")
	state := f.String("state", "", "private operational store")
	id := f.String("review", "", "review UUID (run/inspect)")
	op := f.String("operation", "", "stable prepare UUID")
	input := f.String("input", "", "ReviewSpec JSON file (prepare)")
	binary := f.String("codex", "codex", "managed Codex executable")
	mcp := f.String("mcp-binary", "", "read-only PG MCP executable")
	skills := []string{}
	f.Func("skill", "project-relative reviewer SKILL.md; repeat up to eight", func(s string) error { skills = append(skills, s); return nil })
	timeout := f.Duration("timeout", 5*time.Minute, "controller deadline (up to 30m)")
	interval := f.Duration("interval", time.Second, "poll interval (at least 100ms)")
	if e := f.Parse(args[1:]); e != nil {
		return e
	}
	if *root == "" || *state == "" || f.NArg() != 0 || *timeout <= 0 || *timeout > 30*time.Minute || *interval < 100*time.Millisecond {
		return fmt.Errorf("explicit project/state and bounded timing required")
	}
	if action == "prepare" {
		if *op == "" || *input == "" || *id != "" || *mcp == "" || len(skills) == 0 || len(skills) > 8 {
			return fmt.Errorf("prepare requires --operation --input --mcp-binary --skill")
		}
	} else if *id == "" || *op != "" || *input != "" || len(skills) != 0 || (action == "run" && *mcp == "") {
		return fmt.Errorf("run/inspect requires --review; run also requires --mcp-binary")
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
	if action == "supersede" {
		v, e := s.ReviewRead(*id)
		if e != nil {
			return e
		}
		r, e := s.ReviewWrite(governance.LocalController(), governance.ReviewRequest{Schema: governance.ReviewSchema, ProjectRef: v.ProjectRef, CheckoutID: v.CheckoutID, OperationID: governancehost.UUID(), RunID: v.Review.RunID, ReviewID: *id, ExpectedRevision: v.Review.Revision, Action: "supersede"})
		if e != nil {
			return e
		}
		return encoder.Encode(r)
	}
	if action == "inspect" {
		v, e := s.ReviewRead(*id)
		if e != nil {
			return e
		}
		return encoder.Encode(v)
	}
	mcpPath, e := filepath.Abs(*mcp)
	if e != nil {
		return e
	}
	statePath, e := filepath.Abs(*state)
	if e != nil {
		return e
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	host, e := governancehost.Open(ctx, *binary, p.Root)
	if e != nil {
		return e
	}
	defer host.Close()
	c := governancecontroller.ReviewController{Service: s, Host: host, Root: p.Root, StateDir: statePath, MCPBinary: mcpPath, SkillPaths: skills}
	if action == "prepare" {
		file, e := os.Open(*input)
		if e != nil {
			return e
		}
		defer file.Close()
		b, e := io.ReadAll(io.LimitReader(file, governance.MaxRequest+1))
		if e != nil {
			return e
		}
		var spec governance.ReviewSpec
		if e = governance.Decode(b, &spec); e != nil {
			return e
		}
		// Host/config/skill attribution comes from this selected adapter, not the file.
		if spec.Instructions != "" || spec.HostVersion != "" || spec.MCPBinary != "" || len(spec.Skills) > 0 {
			return fmt.Errorf("host and skill fields are supplied by the controller flags")
		}
		v, e := c.Prepare(ctx, *op, spec)
		if e != nil {
			return e
		}
		return encoder.Encode(v)
	}
	revision := uint64(0)
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()
	for {
		v, e := c.Step(ctx, *id)
		if e != nil {
			return e
		}
		if v.Review.Revision != revision {
			if e = encoder.Encode(v); e != nil {
				return e
			}
			revision = v.Review.Revision
		}
		if !governance.ReviewActive(v.Review.State) {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
