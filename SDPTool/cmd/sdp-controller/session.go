package main

import (
	"bufio"
	"context"
	"encoding/json"
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
	run := f.String("run", "", "existing authorized Maintenance run UUID")
	binary := f.String("codex", "codex", "pinned Codex executable")
	mcp := f.String("mcp-binary", "", "sdp-mcp executable")
	message := f.String("message", "", "file with exact owner input for launch")
	approvals := f.String("approvals", "never", "never or ask for per-request SDP MCP approval")
	sandbox := f.String("sandbox", "read-only", "read-only or explicitly authorized workspace-write")
	timeout := f.Duration("timeout", 30*time.Minute, "maximum controller lifetime")
	if e := f.Parse(args[1:]); e != nil {
		return e
	}
	if *root == "" || *state == "" || *run == "" || f.NArg() != 0 || *timeout <= 0 {
		return fmt.Errorf("require --project ROOT --state DIR --run UUID")
	}
	if *approvals != "never" && *approvals != "ask" {
		return fmt.Errorf("unsupported approval mode")
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
	var text string
	if args[0] == "run" {
		if *message == "" || *mcp == "" {
			return fmt.Errorf("run requires --message FILE --mcp-binary BINARY")
		}
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
	unlock, e := service.LiveControllerLease(governance.LocalController())
	if e != nil {
		return e
	}
	defer unlock()
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ctx, cancel := context.WithTimeout(signalCtx, *timeout)
	defer cancel()
	client, e := governancehost.Open(ctx, *binary, p.Root)
	if e != nil {
		return e
	}
	defer client.Close()
	ctl := &governancecontroller.Controller{Service: service, Host: client, Root: p.Root, StateDir: dir, MCPBinary: *mcp, Epoch: client.Epoch, Sandbox: *sandbox, Approvals: *approvals}
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
	if args[0] != "run" {
		return fmt.Errorf("unknown owner command")
	}
	r, e := ctl.Launch(ctx, *run, text)
	if e != nil {
		_, _ = ctl.Disconnect(*run)
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
	for {
		select {
		case <-ctx.Done():
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
			if in.value.Action == "status" {
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
				_, _ = ctl.Disconnect(*run)
				return fmt.Errorf("host transport ended; reconciliation required: %v", client.Err())
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
				emit(r)
				return nil
			}
		}
	}
}
