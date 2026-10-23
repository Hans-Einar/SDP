package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"github.com/Hans-Einar/SDP/SDPTool"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func inboxCommand(args []string, output io.Writer) error {
	f := flag.NewFlagSet("inbox", flag.ContinueOnError)
	root := f.String("project", "", "explicit project root")
	state := f.String("state", "", "private operational state directory")
	run := f.String("run", "", "optional run UUID")
	jsonOutput := f.Bool("json", false, "emit a single JSON snapshot")
	if err := f.Parse(args); err != nil {
		return err
	}
	if *root == "" || *state == "" || f.NArg() != 0 {
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
	if _, err := fmt.Fprintf(output, "PG report inbox: %d report(s)\nProject: %q\nInspected: %s\nReports await verification; consumption and acceptance are not tracked.\n", len(view.Reports), view.ProjectRef, view.InspectedAt); err != nil {
		return err
	}
	for _, r := range view.Reports {
		if _, err := fmt.Fprintf(output, "\nReport %q / run %q\n  Task: %q\n  Submitted revision: %d; current: %d / %q\n  Coverage: %q\n  Report: %q\n", r.ID, r.RunID, r.Outcome, r.ReportedRevision, r.CurrentRevision, r.CurrentState, r.Coverage, r.Text); err != nil {
			return err
		}
		for _, c := range r.Conditions {
			if _, err := fmt.Fprintf(output, "  Condition: %q\n", c); err != nil {
				return err
			}
		}
	}
	return nil
}
