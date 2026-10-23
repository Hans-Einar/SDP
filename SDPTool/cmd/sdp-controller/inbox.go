package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Hans-Einar/SDP/SDPTool"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func inboxCommand(args []string, output io.Writer) error {
	acknowledge := len(args) > 0 && args[0] == "acknowledge"
	if acknowledge {
		args = args[1:]
	}
	f := flag.NewFlagSet("inbox", flag.ContinueOnError)
	request := f.String("input", "", "pinned acknowledgement JSON request (acknowledge only)")
	root := f.String("project", "", "explicit project root")
	state := f.String("state", "", "private operational state directory")
	run := f.String("run", "", "optional run UUID")
	jsonOutput := f.Bool("json", false, "emit a single JSON snapshot")
	if err := f.Parse(args); err != nil {
		return err
	}
	if *root == "" || *state == "" || f.NArg() != 0 || (acknowledge && (*request == "" || *run != "")) || (!acknowledge && *request != "") {
		return fmt.Errorf("inbox requires --project ROOT --state DIR [--run UUID] [--json]")
	}
	p, err := sdptool.Discover(*root)
	if err != nil {
		return err
	}
	s, err := governance.New(p.Root, p.Area, p.Inventory.ProjectID, *state)
	if err != nil {
		return err
	}
	if acknowledge {
		file, e := os.Open(*request)
		if e != nil {
			return e
		}
		defer file.Close()
		data, e := io.ReadAll(io.LimitReader(file, governance.MaxRequest+1))
		if e != nil {
			return e
		}
		var q governance.ReportAcknowledgementRequest
		if e = governance.Decode(data, &q); e != nil {
			return e
		}
		result, e := s.AcknowledgeReport(governance.LocalController(), q)
		if e != nil {
			return e
		}
		return json.NewEncoder(output).Encode(result)
	}
	view, err := s.ReportInbox(*run)
	if err != nil {
		return err
	}
	if *jsonOutput {
		return json.NewEncoder(output).Encode(view)
	}
	return renderInbox(output, view)
}

func renderInbox(output io.Writer, view governance.ReportInbox) error {
	if _, err := fmt.Fprintf(output, "PG report inbox: %d report(s)\nProject: %q\nInspected: %s\nAcknowledgement records processing only; verification and acceptance remain separate.\n", len(view.Reports), view.ProjectRef, view.InspectedAt); err != nil {
		return err
	}
	for _, r := range view.Reports {
		if _, err := fmt.Fprintf(output, "\nReport %q / run %q\n  Task: %q\n  Submitted revision: %d; current: %d / %q\n  Coverage: %q\n  Report: %q\n", r.ID, r.RunID, r.Outcome, r.ReportedRevision, r.CurrentRevision, r.CurrentState, r.Coverage, r.Text); err != nil {
			return err
		}
		if a := r.Acknowledgement; a != nil {
			if _, err := fmt.Fprintf(output, "  Acknowledged by %q at %s: %q\n", a.Request.Coordinator, a.RecordedAt, a.Request.Note); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintln(output, "  Acknowledgement: pending"); err != nil {
				return err
			}
		}
		for _, c := range r.Conditions {
			if _, err := fmt.Fprintf(output, "  Condition: %q\n", c); err != nil {
				return err
			}
		}
	}
	return nil
}
