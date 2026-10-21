package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func attentionCommand(args []string, out io.Writer, input io.Reader, prompt io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("attention requires list|request|reply|withdraw")
	}
	action := args[0]
	if action != "list" && action != "request" && action != "reply" && action != "withdraw" {
		return fmt.Errorf("unknown attention operation")
	}
	f := flag.NewFlagSet("attention "+action, flag.ContinueOnError)
	root := f.String("project", "", "explicit project root")
	state := f.String("state", "", "private operational state directory")
	request := f.String("input", "", "bounded JSON request file")
	run := f.String("run", "", "optional run filter (list only)")
	owner := f.String("owner", "", "attributed owner (reply only)")
	reason := f.String("reason", "", "required withdrawal reason")
	jsonOut := f.Bool("json", false, "machine-readable list; mutations always return JSON")
	if e := f.Parse(args[1:]); e != nil {
		return e
	}
	if *root == "" || *state == "" || f.NArg() != 0 {
		return fmt.Errorf("explicit --project and --state required")
	}
	if (action == "list" && *request != "") || (action != "list" && (*request == "" || *run != "")) || (action == "reply" && *owner == "") || (action != "reply" && *owner != "") || (action == "withdraw" && *reason == "") || (action != "withdraw" && *reason != "") {
		return fmt.Errorf("invalid operation-specific attention flags")
	}
	p, e := sdptool.Discover(*root)
	if e != nil {
		return e
	}
	s, e := governance.New(p.Root, p.Area, p.Inventory.ProjectID, *state)
	if e != nil {
		return e
	}
	if action == "list" {
		v, e := s.AttentionInbox(*run)
		if e != nil {
			return e
		}
		if *jsonOut {
			return json.NewEncoder(out).Encode(v)
		}
		return renderAttention(out, v)
	}
	file, e := os.Open(*request)
	if e != nil {
		return e
	}
	defer file.Close()
	data, e := io.ReadAll(io.LimitReader(file, governance.MaxRequest+1))
	if e != nil {
		return e
	}
	var result governance.Result
	if action == "request" {
		var q governance.AttentionRequest
		if e = governance.Decode(data, &q); e == nil {
			result, e = s.RequestAttention(governance.LocalController(), q)
		}
	} else {
		var q governance.AttentionDisposition
		if e = governance.Decode(data, &q); e == nil {
			if action == "withdraw" {
				result, e = s.WithdrawAttention(governance.LocalController(), q, *reason)
			} else {
				result, e = captureAttentionReply(s, q, *owner, input, prompt)
			}
		}
	}
	if e != nil {
		return e
	}
	return json.NewEncoder(out).Encode(result)
}
func renderAttention(out io.Writer, v governance.AttentionInbox) error {
	if _, e := fmt.Fprintf(out, "PG attention: %d question(s)\nProject: %q\nReplies are recorded input. Dispatch acknowledgement does not prove agent consumption or acceptance.\nAn authorized waiting controller may continue after a current reply.\n", len(v.Items), v.ProjectRef); e != nil {
		return e
	}
	for _, item := range v.Items {
		q := item.Question
		if _, e := fmt.Fprintf(out, "\nQuestion %q / run %q\n  Status: %q; applicability: %q; coverage: %q; blocking: %t; delivery: %q\n  Prompt: %q\n", q.ID, item.RunID, item.Status, item.Applicability, item.Coverage, q.Blocking, item.Consumption, q.Prompt); e != nil {
			return e
		}
		for _, option := range q.Options {
			if _, e := fmt.Fprintf(out, "  Option %q: %q\n", option.ID, option.Label); e != nil {
				return e
			}
		}
		if q.Reply != nil {
			if _, e := fmt.Fprintf(out, "  Reply by %q: %q (option %q)\n", q.Reply.Actor, q.Reply.Text, q.Reply.OptionID); e != nil {
				return e
			}
		}
		if q.Withdrawal != nil {
			if _, e := fmt.Fprintf(out, "  Withdrawn: %q\n", q.Withdrawal.Reason); e != nil {
				return e
			}
		}
	}
	return nil
}
func captureAttentionReply(s *governance.Service, q governance.AttentionDisposition, owner string, input io.Reader, prompt io.Writer) (governance.Result, error) {
	var empty governance.Result
	inbox, e := s.AttentionInbox(q.RunID)
	if e != nil {
		return empty, e
	}
	for _, item := range inbox.Items {
		question := item.Question
		if question.ID != q.QuestionID {
			continue
		}
		if old := question.Reply; old != nil && old.OperationID == q.OperationID {
			if old.Actor != owner || old.Source != "local-controller-stdin" {
				return empty, fmt.Errorf("idempotency_conflict: reply attribution differs")
			}
			return s.ReplyAttention(governance.LocalController(), q, governance.OwnerAttentionInput{Actor: owner, Source: old.Source, OptionID: old.OptionID, Text: old.Text})
		}
		if item.Status != "open" || item.Applicability != "current" || question.Basis != q.Basis {
			return empty, fmt.Errorf("stale_revision: displayed question is not open/current")
		}
		if e = renderAttention(prompt, governance.AttentionInbox{ProjectRef: inbox.ProjectRef, Items: []governance.AttentionView{item}}); e != nil {
			return empty, e
		}
		if _, e = fmt.Fprintf(prompt, "Reply attributed to %q. Enter an option ID, or leave blank for free text (no default):\n", owner); e != nil {
			return empty, e
		}
		scanner := bufio.NewScanner(input)
		scanner.Buffer(make([]byte, 256), 4098)
		if !scanner.Scan() {
			return empty, fmt.Errorf("no reply recorded: input ended: %v", scanner.Err())
		}
		option := strings.TrimSpace(scanner.Text())
		if _, e = fmt.Fprintln(prompt, "Enter your answer/comment (required; an authorized waiting controller may continue):"); e != nil {
			return empty, e
		}
		if !scanner.Scan() {
			return empty, fmt.Errorf("no reply recorded: answer ended: %v", scanner.Err())
		}
		return s.ReplyAttention(governance.LocalController(), q, governance.OwnerAttentionInput{Actor: owner, Source: "local-controller-stdin", OptionID: option, Text: scanner.Text()})
	}
	return empty, fmt.Errorf("question absent from selected run")
}
