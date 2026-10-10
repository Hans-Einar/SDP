package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func overviewCommand(args []string, output io.Writer) error {
	f := flag.NewFlagSet("overview", flag.ContinueOnError)
	root := f.String("project", "", "explicit project root")
	state := f.String("state", "", "existing private operational state directory")
	jsonOutput := f.Bool("json", false, "emit a single JSON snapshot")
	if err := f.Parse(args); err != nil {
		return err
	}
	if *root == "" || *state == "" || f.NArg() != 0 {
		return fmt.Errorf("overview requires --project ROOT --state DIR [--json]")
	}
	p, err := sdptool.Discover(*root)
	if err != nil {
		return err
	}
	s, err := governance.New(p.Root, p.Area, p.Inventory.ProjectID, *state)
	if err != nil {
		return err
	}
	view, err := s.Overview()
	if err != nil {
		return err
	}
	if *jsonOutput {
		return json.NewEncoder(output).Encode(view)
	}
	return renderOverview(output, view)
}

func renderOverview(output io.Writer, view governance.Overview) error {
	var b strings.Builder
	fmt.Fprintf(&b, "Project %q\nRoot: %q\nInspected: %s\n", view.ProjectRef, view.Root, view.InspectedAt)
	fmt.Fprintf(&b, "Registered work: %d run(s) in this checkout store\n", len(view.Runs))
	fmt.Fprintln(&b, "Snapshot only; live state, child mapping and owner acceptance are not observed.")
	if len(view.Runs) == 0 {
		fmt.Fprintln(&b, "No registered runs here. Other stores and unregistered project work are outside this view.")
	}
	for _, r := range view.Runs {
		fmt.Fprintf(&b, "\nRun %q (revision %d)\n  Task: %q\n", r.RunID, r.Revision, r.Outcome)
		fmt.Fprintf(&b, "  Role: %q; evidence: %q\n", r.AssignedRole, r.RoleEvidence)
		fmt.Fprintf(&b, "  Session: %q; step: %q\n  Plan: %q\n  Card: %q\n", r.Selection.Session, r.Selection.Step, r.Selection.Plan, r.Selection.Card)
		fmt.Fprintf(&b, "  Routine: %q / %q; coverage: %q\n", r.RoutineState, r.RoutineStep, r.Coverage)
		for _, c := range r.Conditions {
			fmt.Fprintf(&b, "  Coverage condition: %q\n", c)
		}
		fmt.Fprintf(&b, "  Last host status: %q at %q; observation gap: %t\n", r.LastHostStatus, r.LastObservedAt, r.ObservationGap)
		fmt.Fprintf(&b, "  Attempt: %q; thread: %q; turn: %q\n", r.AttemptID, r.ThreadID, r.TurnID)
		fmt.Fprintf(&b, "  Tool completions recorded in latest attempt: %d; last item: %q (%q)\n", r.ToolCompletions, r.LastToolItemID, r.LastToolStatus)
		fmt.Fprintf(&b, "  Explicit skills (latest launch input): %d; acceptance: %q; compliance: %q\n", r.RequestedSkills, r.SkillAcceptance, r.SkillCompliance)
		if r.SkillAttemptID != "" {
			fmt.Fprintf(&b, "  Skill input: %q; attempt: %q\n", r.SkillInputID, r.SkillAttemptID)
		}
	}
	_, err := io.WriteString(output, b.String())
	return err
}
