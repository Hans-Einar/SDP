package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancecontroller"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

func ownerCommand(args []string) error {
	signal.Ignore(syscall.SIGPIPE) // A closed output observer must not terminate work.
	f := flag.NewFlagSet(args[0], flag.ContinueOnError)
	root := f.String("project", "", "explicit project root")
	state := f.String("state", "", "private operational state directory")
	run := f.String("run", "", "existing authorized Maintenance or implementation run UUID")
	instance := f.String("instance", "", "registered Session instance instead of --run; resolves its current committed assignment")
	socket := f.String("host-socket", "", "private local socket of separately running pinned host")
	binary := f.String("codex", "codex", "pinned Codex executable")
	mcp := f.String("mcp-binary", "", "sdp-mcp executable")
	message := f.String("message", "", "file with exact owner input for launch")
	approvals := f.String("approvals", "never", "never or ask for per-request SDP MCP approval")
	reporting := f.String("reporting", "manual", "manual or automatic assignment-bound reporting (run/continue only)")
	continuations := f.Int("attention-continuations", 0, "maximum automatic fresh attempts after current attention replies (0..8; run/continue only)")
	sandbox := f.String("sandbox", "read-only", "read-only or explicitly authorized workspace-write")
	var skillPaths []string
	f.Func("skill", "project-relative enabled SKILL.md to invoke (repeatable; run only)", func(v string) error { skillPaths = append(skillPaths, v); return nil })
	timeout := f.Duration("timeout", 30*time.Minute, "maximum controller lifetime")
	if e := f.Parse(args[1:]); e != nil {
		return e
	}
	if *root == "" || *state == "" || ((*run == "") == (*instance == "")) || f.NArg() != 0 || *timeout <= 0 {
		return fmt.Errorf("require --project ROOT --state DIR and exactly one of --run UUID or --instance UUID")
	}
	if *continuations < 0 || *continuations > 8 || (*continuations > 0 && args[0] != "run" && args[0] != "continue") {
		return fmt.Errorf("--attention-continuations requires run/continue and 0..8")
	}
	if args[0] == "continue" && (*message != "" || len(skillPaths) > 0 || *mcp == "") {
		return fmt.Errorf("continue requires --mcp-binary and reuses original input/skills")
	}
	if len(skillPaths) > 0 && args[0] != "run" {
		return fmt.Errorf("--skill is only supported by run")
	}
	if *approvals != "never" && *approvals != "ask" {
		return fmt.Errorf("unsupported approval mode")
	}
	if (*reporting != "manual" && *reporting != "automatic") || (*reporting != "manual" && args[0] != "run" && args[0] != "continue") {
		return fmt.Errorf("--reporting requires run/continue and manual or automatic")
	}
	p, e := sdptool.Discover(*root)
	if e != nil {
		return e
	}
	dir, e := filepath.Abs(*state)
	if e != nil {
		return e
	}
	service, e := governance.New(p.Root, p.Area, p.Inventory.ProjectID, dir)
	if e != nil {
		return e
	}
	emit := func(v any) { _ = json.NewEncoder(os.Stdout).Encode(v) } // Observers do not own execution.
	var sessionInput string
	if *instance != "" {
		view, err := service.ReadSession(*instance)
		if err != nil {
			return err
		}
		if args[0] == "status" {
			emit(view)
		}
		selected, err := service.SessionHostSelection(*instance)
		if err != nil {
			var missing *governance.Error
			if args[0] == "status" && errors.As(err, &missing) && missing.Code == "missing_coverage" {
				return nil
			}
			return err
		}
		*run = selected.RunID
		if args[0] == "run" {
			current, err := service.Read(*run)
			if err != nil {
				return err
			}
			if current.Run.State != "ready" {
				return &governance.Error{Code: "conflict", Message: "current Session run is not ready; inspect or use the explicit continuation/recovery operation"}
			}
			if err := service.CheckRunReady(*run); err != nil {
				return err
			}
			if len(skillPaths) == 0 {
				for _, skill := range selected.Skills {
					rel, err := filepath.Rel(p.Root, skill.Path)
					if err != nil {
						return err
					}
					skillPaths = append(skillPaths, rel)
				}
			}
			sessionInput = "Execute the committed Session assignment: " + selected.Outcome
		}
	}

	if args[0] == "provenance" {
		p, err := service.Provenance(*run)
		if err == nil {
			emit(p)
		}
		return err
	}
	if args[0] == "status" {
		r, e := service.Read(*run)
		if e == nil {
			emit(r)
			prompts, err := service.PendingApprovals(*run)
			if err != nil {
				return err
			}
			if len(prompts) > 0 {
				emit(map[string]any{"pendingApprovals": prompts})
			}
		}
		return e
	}
	text := sessionInput
	if args[0] == "run" {
		if (*message == "" && *instance == "") || *mcp == "" {
			return fmt.Errorf("run requires --mcp-binary BINARY and --message FILE unless a committed --instance supplies its outcome")
		}
		if *message != "" {
			file, e := os.Open(*message)
			if e != nil {
				return e
			}
			defer file.Close()
			b, e := io.ReadAll(io.LimitReader(file, (64<<10)+1))
			if e != nil {
				return e
			}
			if len(b) > 64<<10 {
				return fmt.Errorf("owner input exceeds 64 KiB")
			}
			text = string(b)
		}
	}
	unlock, e := service.RunLiveControllerLease(governance.LocalController(), *run)
	if e != nil {
		return e
	}
	defer unlock()
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ctx, cancel := context.WithTimeout(signalCtx, *timeout)
	defer cancel()
	var client *governancehost.Client
	if args[0] == "attach" && *socket == "" {
		return fmt.Errorf("attach requires --host-socket")
	}
	if *socket != "" {
		client, e = governancehost.OpenSocket(ctx, *socket)
	} else {
		client, e = governancehost.Open(ctx, *binary, p.Root)
	}
	if e != nil {
		return e
	}
	defer client.Close()
	ctl := &governancecontroller.Controller{SkillPaths: skillPaths, Service: service, Host: client, Root: p.Root, StateDir: dir, MCPBinary: *mcp, HostIdentity: client.HostIdentity, Epoch: client.Epoch, Sandbox: *sandbox, Approvals: *approvals, Reporting: *reporting}
	if args[0] == "recover" {
		r, e := service.Read(*run)
		if e != nil {
			return e
		}
		if len(r.Run.Attempts) == 0 {
			return fmt.Errorf("no mapped attempt")
		}
		recovery, e := client.PrepareRecovery(ctx, service, *run, r.Run.Attempts[len(r.Run.Attempts)-1].ID)
		if e != nil {
			return e
		}
		r, e = service.RecordHost(governance.LocalController(), recovery.ReconcileRequest())
		if e != nil {
			return e
		}
		for _, event := range recovery.Events {
			r, e = ctl.Observe(*run, event)
			if e != nil {
				return e
			}
		}
		for _, input := range recovery.Inputs {
			r, e = service.RecordInput(governance.LocalController(), governance.InputRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, OperationID: governancehost.UUID(), ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: "resolve", Input: input})
			if e != nil {
				return e
			}
		}
		emit(r)
		return nil // History reconciliation only; never resume/resend a model turn.
	}
	if args[0] != "run" && args[0] != "attach" && args[0] != "continue" {
		return fmt.Errorf("unknown owner command")
	}
	var r governance.Result
	if args[0] == "attach" {
		recovery, err := client.PrepareAttachment(ctx, service, *run, p.Root, *approvals)
		if err != nil {
			return err
		}
		r, e = service.RecordHost(governance.LocalController(), recovery.ReattachRequest(client.HostIdentity))
	} else if args[0] == "continue" {
		r, e = ctl.ContinueAttention(ctx, *run)
	} else {
		r, e = ctl.Launch(ctx, *run, text)
	}
	if e != nil {
		// A refused continuation must preserve the completed predecessor. Launch
		// already records post-claim delivery uncertainty before returning.
		if args[0] != "continue" {
			_, _ = ctl.Disconnect(*run)
		}
		return e
	}
	emit(r)
	type command struct {
		Action   string                    `json:"action"`
		Text     string                    `json:"text"`
		Approval *governance.ApprovalReply `json:"approval,omitempty"`
	}
	type input struct {
		value command
		err   error
		exact string
	}
	commands := make(chan input, 1)
	go func() {
		defer close(commands)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Buffer(make([]byte, 4096), governance.MaxRequest+1)
		for scanner.Scan() {
			var v command
			e := governance.Decode(scanner.Bytes(), &v)
			select {
			case commands <- input{value: v, err: e, exact: string(scanner.Bytes())}:
			case <-ctx.Done():
				return
			}
		}
		if e := scanner.Err(); e != nil {
			select {
			case commands <- input{err: e}:
			case <-ctx.Done():
			}
		}
	}()
	waitingForReply := false
	var ticks <-chan time.Time
	if *continuations > 0 {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		ticks = ticker.C
	}
	for {
		select {
		case <-ticks:
			if !waitingForReply {
				continue
			}
			waiting, ready, err := ctl.AttentionPending(*run)
			if err != nil {
				var changed *governance.Error
				if errors.As(err, &changed) && changed.Code == "stale_revision" {
					continue
				} // Read-only generation changed; no host effect occurred.
				return err
			}
			if !waiting && !ready {
				return nil
			}
			if ready {
				r, e = ctl.ContinueAttention(ctx, *run)
				if e != nil {
					return e
				}
				*continuations--
				waitingForReply = false
				emit(r)
			}
		case <-ctx.Done():
			if waitingForReply {
				return fmt.Errorf("attention wait ended; completed attempt and durable question retained: %w", ctx.Err())
			}
			_, _ = ctl.Disconnect(*run)
			return fmt.Errorf("controller ended; host outcome requires reconciliation: %w", ctx.Err())
		case in, ok := <-commands:
			if !ok {
				commands = nil
				continue
			} // Input EOF is observation-only, not cancellation.
			if in.err != nil {
				emit(map[string]any{"error": in.err.Error()})
				continue
			}
			if in.value.Action == "detach" {
				if client.HostIdentity == "" {
					emit(map[string]any{"error": "detach requires a separately running socket host"})
					continue
				}
				r, e = ctl.Disconnect(*run)
				if e != nil {
					return e
				}
				emit(r)
				return nil
			} else if in.value.Action == "status" {
				r, e = service.Read(*run)
			} else if in.value.Action == "approval" && in.value.Approval != nil {
				r, e = ctl.Approve(*run, *in.value.Approval, in.exact)
			} else {
				r, e = ctl.Control(ctx, *run, in.value.Action, in.value.Text)
			}
			if e != nil {
				emit(map[string]any{"error": e.Error()})
			} else {
				emit(r)
			}
		case event, ok := <-client.Notifications:
			if !ok {
				if waitingForReply {
					return fmt.Errorf("host ended while waiting; completed attempt and question retained: %v", client.Err())
				}
				_, _ = ctl.Disconnect(*run)
				return fmt.Errorf("host transport ended; reconciliation required: %v", client.Err())
			}
			// Events from completed predecessor threads cannot control the fresh attempt.
			if len(r.Run.Attempts) > 1 && event.ThreadID != r.Run.Attempts[len(r.Run.Attempts)-1].ThreadID {
				historical := false
				for _, a := range r.Run.Attempts[:len(r.Run.Attempts)-1] {
					if a.ThreadID == event.ThreadID && a.Status == "completed" {
						historical = true
					}
				}
				if historical {
					emit(map[string]any{"ignoredCompletedPredecessorEvent": event})
					continue
				}
			}
			if event.Observation.Kind == "approval.resolved" {
				r, e = ctl.ResolveApproval(*run, event)
			} else {
				r, e = ctl.Observe(*run, event)
			}
			if e != nil {
				_, _ = ctl.Disconnect(*run)
				return e
			}
			emit(event)
			if event.Observation.Kind == "approval.requested" {
				var prompt governance.ApprovalPrompt
				if e := json.Unmarshal([]byte(event.Observation.Text), &prompt); e != nil {
					return e
				}
				emit(map[string]any{"approval": prompt})
			}
			if event.Observation.Kind == "turn.completed" {
				// A reconnect gap remains explicit. Fetch only terminal history to recover
				// completed messages omitted by the conservative active snapshot.
				if args[0] == "attach" {
					recovery, err := client.PrepareRecovery(ctx, service, *run, r.Run.Attempts[len(r.Run.Attempts)-1].ID)
					if err != nil {
						return err
					}
					for _, replay := range recovery.Events {
						r, e = ctl.Observe(*run, replay)
						if e != nil {
							return e
						}
					}
				}
				emit(r)
				if *continuations > 0 {
					waitingForReply = true
					emit(map[string]any{"attentionChecking": true, "remainingContinuations": *continuations})
					continue
				}
				return nil
			}
		}
	}
}
