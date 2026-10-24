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
	"time"

	"github.com/Hans-Einar/SDP/SDPTool"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

type coordinationCommand struct {
	Action          string                                   `json:"action"`
	Acknowledgement *governance.ReportAcknowledgementRequest `json:"acknowledgement,omitempty"`
	Disposition     *governance.AttentionDisposition         `json:"disposition,omitempty"`
	Reply           *coordinationReply                       `json:"reply,omitempty"`
}
type coordinationReply struct {
	Text     string `json:"text"`
	OptionID string `json:"optionId,omitempty"`
}
type coordinationInput struct {
	Command coordinationCommand
	Err     error
}
type coordinationEvent struct {
	Schema string                       `json:"schema"`
	Kind   string                       `json:"kind"`
	View   *governance.CoordinationView `json:"view,omitempty"`
	Result *governance.Result           `json:"result,omitempty"`
	Error  *governance.Error            `json:"error,omitempty"`
}

func coordinationCommandMain(args []string, input io.Reader, out io.Writer) error {
	f := flag.NewFlagSet("coordinate", flag.ContinueOnError)
	root := f.String("project", "", "explicit project root")
	state := f.String("state", "", "private operational state directory")
	run := f.String("run", "", "optional run filter")
	cursor := f.String("cursor", "", "snapshot page cursor; not used with watch")
	limit := f.Int("limit", 20, "page item limit (1..50)")
	watch := f.Bool("watch", false, "refresh changed snapshots and accept bounded JSON input")
	interval := f.Duration("interval", time.Second, "watch polling interval (at least 100ms)")
	lifetime := f.Duration("timeout", 30*time.Minute, "maximum watch lifetime")
	coordinator := f.String("coordinator", "", "local attribution allowed for report acknowledgements")
	owner := f.String("owner", "", "local owner attribution allowed for attention replies")
	jsonOut := f.Bool("json", false, "emit JSON snapshots or watch event lines")
	if e := f.Parse(args); e != nil {
		return e
	}
	if *root == "" || *state == "" || f.NArg() != 0 || *limit < 1 || *limit > 50 || *interval < 100*time.Millisecond || *lifetime <= 0 || (*watch && *cursor != "") || (!*watch && (*coordinator != "" || *owner != "")) {
		return fmt.Errorf("coordinate requires --project ROOT --state DIR; --watch permits coordinator/owner input but no cursor")
	}
	p, e := sdptool.Discover(*root)
	if e != nil {
		return e
	}
	s, e := governance.New(p.Root, p.Area, p.Inventory.ProjectID, *state)
	if e != nil {
		return e
	}
	q := governance.CoordinationRequest{RunID: *run, Cursor: *cursor, Limit: *limit}
	if !*watch {
		v, e := s.Coordination(q)
		if e != nil {
			return e
		}
		if *jsonOut {
			return json.NewEncoder(out).Encode(v)
		}
		return renderCoordination(out, v)
	}
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	ctx, cancel := context.WithTimeout(signalCtx, *lifetime)
	defer cancel()
	commands := readCoordinationInput(ctx, input)
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()
	return watchCoordination(ctx, s, q, *coordinator, *owner, commands, ticker.C, out, *jsonOut)
}
func readCoordinationInput(ctx context.Context, input io.Reader) <-chan coordinationInput {
	commands := make(chan coordinationInput, 1)
	go func() {
		defer close(commands)
		scanner := bufio.NewScanner(input)
		scanner.Buffer(make([]byte, 4096), governance.MaxRequest+1)
		for scanner.Scan() {
			var v coordinationCommand
			e := governance.Decode(scanner.Bytes(), &v)
			select {
			case commands <- coordinationInput{v, e}:
			case <-ctx.Done():
				return
			}
		}
		if e := scanner.Err(); e != nil {
			select {
			case commands <- coordinationInput{Err: e}:
			case <-ctx.Done():
			}
		}
	}()
	return commands
}
func applyCoordinationCommand(s *governance.Service, q governance.CoordinationRequest, coordinator, owner string, v coordinationCommand) (*governance.Result, error) {
	invalid := func() (*governance.Result, error) {
		return nil, &governance.Error{Code: "invalid_request", Message: "use refresh, quit, exact acknowledge, or exact reply input"}
	}
	switch v.Action {
	case "refresh", "quit":
		if v.Acknowledgement != nil || v.Disposition != nil || v.Reply != nil {
			return invalid()
		}
		return nil, nil
	case "acknowledge":
		if v.Acknowledgement == nil || v.Disposition != nil || v.Reply != nil {
			return invalid()
		}
		a := v.Acknowledgement
		if coordinator == "" || a.Coordinator != coordinator || (q.RunID != "" && q.RunID != a.RunID) {
			return nil, &governance.Error{Code: "unauthorized", Message: "acknowledgement is outside console launch scope"}
		}
		r, e := s.AcknowledgeReport(governance.LocalController(), *a)
		return &r, e
	case "reply":
		if v.Disposition == nil || v.Reply == nil || v.Acknowledgement != nil {
			return invalid()
		}
		if owner == "" || (q.RunID != "" && q.RunID != v.Disposition.RunID) {
			return nil, &governance.Error{Code: "unauthorized", Message: "reply is outside console launch scope"}
		}
		r, e := s.ReplyAttention(governance.LocalController(), *v.Disposition, governance.OwnerAttentionInput{Actor: owner, Source: "sdp-controller coordinate stdin", Text: v.Reply.Text, OptionID: v.Reply.OptionID})
		return &r, e
	default:
		return invalid()
	}
}
func watchCoordination(ctx context.Context, s *governance.Service, q governance.CoordinationRequest, coordinator, owner string, commands <-chan coordinationInput, ticks <-chan time.Time, out io.Writer, jsonOut bool) error {
	emit := func(v coordinationEvent) error {
		v.Schema = "sdp-governance-coordination-event/1"
		if jsonOut {
			return json.NewEncoder(out).Encode(v)
		}
		if v.View != nil {
			return renderCoordination(out, *v.View)
		}
		if v.Error != nil {
			_, e := fmt.Fprintf(out, "Coordination error %q: %q\n", v.Error.Code, v.Error.Message)
			return e
		}
		_, e := fmt.Fprintf(out, "Recorded %q operation %q; run %q remains %q\n", v.Result.Operation, v.Result.OperationID, v.Result.Run.ID, v.Result.Run.State)
		return e
	}
	reportError := func(e error) error {
		var detail *governance.Error
		if !errors.As(e, &detail) {
			detail = &governance.Error{Code: "invalid_request", Message: e.Error()}
		}
		return emit(coordinationEvent{Kind: "error", Error: detail})
	}
	revision := ""
	refresh := func(force bool) error {
		v, e := s.Coordination(q)
		if e != nil {
			return reportError(e)
		}
		if force || v.Revision != revision {
			if e = emit(coordinationEvent{Kind: "snapshot", View: &v}); e != nil {
				return e
			}
			revision = v.Revision
		}
		return nil
	}
	if e := refresh(true); e != nil {
		return e
	}
	process := func(v coordinationInput) (bool, error) {
		if v.Err != nil {
			return false, reportError(v.Err)
		}
		result, e := applyCoordinationCommand(s, q, coordinator, owner, v.Command)
		if e != nil {
			return false, reportError(e)
		}
		if v.Command.Action == "quit" {
			return true, nil
		}
		if result != nil {
			if e = emit(coordinationEvent{Kind: "result", Result: result}); e != nil {
				return false, e
			}
		}
		return false, refresh(v.Command.Action == "refresh")
	}
	for {
		if ctx.Err() != nil {
			return nil
		}
		// Consume queued owner input before taking another polling turn.
		select {
		case v, ok := <-commands:
			if !ok {
				commands = nil
			} else {
				done, e := process(v)
				if done || e != nil {
					return e
				}
				continue
			}
		default:
		}
		select {
		case <-ctx.Done():
			return nil
		case v, ok := <-commands:
			if !ok {
				commands = nil
				continue
			}
			done, e := process(v)
			if done || e != nil {
				return e
			}
		case _, ok := <-ticks:
			if !ok {
				ticks = nil
				continue
			}
			if e := refresh(false); e != nil {
				return e
			}
		}
	}
}
func renderCoordination(out io.Writer, v governance.CoordinationView) error {
	if _, e := fmt.Fprintf(out, "Coordination for %q: %d of %d item(s)\nInspected: %s; host liveness remains unknown.\n", v.ProjectRef, len(v.Items), v.Total, v.InspectedAt); e != nil {
		return e
	}
	for _, item := range v.Items {
		var text, status string
		switch {
		case item.Attention != nil:
			text = item.Attention.Question.Prompt
			status = item.Attention.Status + " / " + item.Attention.Applicability
		case item.Report != nil:
			text = item.Report.Text
			status = "acknowledgement pending / " + item.Report.Coverage
			if item.Report.Acknowledgement != nil {
				status = "acknowledged / " + item.Report.Coverage
			}
		case item.Work != nil:
			text = item.Work.Outcome
			status = item.Work.RoutineState + " / " + item.Work.Coverage
		}
		if _, e := fmt.Fprintf(out, "  %q run %q: %q — %q\n", item.ID, item.RunID, status, text); e != nil {
			return e
		}
	}
	if v.NextCursor != "" {
		_, e := fmt.Fprintf(out, "More items: inspect with --cursor %q or narrow --run.\n", v.NextCursor)
		return e
	}
	return nil
}
