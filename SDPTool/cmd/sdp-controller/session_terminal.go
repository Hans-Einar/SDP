package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

// This adapter owns no workflow state. Every mutation uses the same guarded core
// as MCP/JSON; resume commits pending canonical effects, never starts a host.
func sessionCommand(args []string, out io.Writer, input io.Reader, prompt io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("session requires list|status|resume|select|next|rework|complete|accept|changes-requested|questions|answer")
	}
	action := args[0]
	if !slices.Contains([]string{"list", "status", "resume", "select", "next", "rework", "complete", "accept", "changes-requested", "questions", "answer"}, action) {
		return fmt.Errorf("unsupported Session command")
	}
	f := flag.NewFlagSet("session "+action, flag.ContinueOnError)
	f.SetOutput(prompt)
	root := f.String("project", "", "explicit project root")
	state := f.String("state", "", "existing private operational state directory")
	instance := f.String("instance", "", "exact registered instance UUID")
	session := f.String("session", "", "canonical Session ID when it identifies exactly one registered instance")
	step := f.String("step", "", "target step for select/next/rework")
	reason := f.String("reason", "", "required reason for workflow changes")
	basis := f.String("basis", "", "rework basis: owner, test or review")
	result := f.String("result", "", "exact failed result for rework")
	gate := f.String("gate", "", "exact displayed owner gate for accept/changes-requested")
	question := f.String("question", "", "question ID; answer selects the only current open question if omitted")
	owner := f.String("owner", "", "attributed owner for answer")
	expected := f.Uint64("expected-revision", 0, "optional inspected Session revision; refuses if changed")
	prepareOnly := f.Bool("prepare-only", false, "prepare effects without committing; resume commits later")
	jsonOut := f.Bool("json", false, "emit machine-readable views/results")
	if e := f.Parse(args[1:]); e != nil {
		return e
	}
	if *root == "" || *state == "" || f.NArg() != 0 {
		return fmt.Errorf("require --project ROOT --state DIR")
	}
	read := slices.Contains([]string{"list", "status", "questions"}, action)
	advancing := slices.Contains([]string{"next", "rework", "complete"}, action)
	deciding := action == "accept" || action == "changes-requested"
	needsStep := action == "select" || action == "next" || action == "rework"
	if (action == "list" && (*instance != "" || *session != "")) || (action != "list" && ((*instance == "") == (*session == ""))) {
		return fmt.Errorf("list needs no selector; other commands require exactly one of --session ID or --instance UUID")
	}
	if (needsStep != (*step != "")) || ((action == "rework") != (*basis != "" && *result != "")) || (action != "rework" && (*basis != "" || *result != "")) || (deciding != (*gate != "")) {
		return fmt.Errorf("step, rework basis/result and owner gate must match the selected command")
	}
	if (action == "answer") != (*owner != "") || (action != "answer" && *question != "") || (*prepareOnly && (read || action == "resume" || action == "answer")) || (action == "list" && *expected != 0) {
		return fmt.Errorf("invalid command-specific owner, question, preparation or revision flags")
	}
	if (!read && action != "answer") != (strings.TrimSpace(*reason) != "") {
		return fmt.Errorf("workflow mutations require --reason; inspection and answers do not use it")
	}
	p, e := sdptool.Discover(*root)
	if e != nil {
		return e
	}
	s, e := governance.New(p.Root, p.Area, p.Inventory.ProjectID, *state)
	if e != nil {
		return e
	}
	list, e := s.ListSessions()
	if e != nil {
		return e
	}
	if action == "list" {
		if *jsonOut {
			return json.NewEncoder(out).Encode(list)
		}
		if _, e = fmt.Fprintf(out, "Registered Sessions in this checkout: %d\n", len(list.Instances)); e != nil {
			return e
		}
		for _, v := range list.Instances {
			if _, e = fmt.Fprintf(out, "%q: %s; plan %q; instance %q; coverage %q\n", v.SessionID, v.State, v.PlanID, v.InstanceID, v.Coverage); e != nil {
				return e
			}
		}
		return nil
	}
	v, e := selectTerminalSession(list, *session, *instance)
	if e != nil {
		return e
	}
	if *expected != 0 && v.Revision != *expected {
		return fmt.Errorf("stale_revision: Session changed from inspected revision %d to %d", *expected, v.Revision)
	}
	if action == "status" {
		return emitTerminalSession(s, v, out, *jsonOut)
	}
	if action == "questions" || action == "answer" {
		return terminalSessionAttention(s, v, action, *question, *owner, *jsonOut, out, input, prompt)
	}
	if action == "resume" && v.Registration == "committed" && v.RecordEffects != "prepared" {
		if v.Coverage != "current" {
			return fmt.Errorf("stale_revision: inspect Session conditions; resume cannot rebind changed authority")
		}
		// Inspection of an already committed stage must not consume a receipt or
		// pretend it is continuing an active native attempt.
		return emitTerminalSession(s, v, out, *jsonOut)
	}
	// A selected registered instance implies an existing binding. Bind returns
	// that exact identity without creating another checkout or workflow.
	identity, e := s.Bind(governance.LocalOwner())
	if e != nil {
		return e
	}
	op := governancehost.UUID()
	if _, e = fmt.Fprintf(prompt, "Session %q, revision %d, operation %s. On interruption inspect status and use session resume; do not repeat a host send.\n", v.SessionID, v.Revision, op); e != nil {
		return e
	}
	var changed governance.SessionMutationResult
	resume := func(revision uint64, operation string) (governance.SessionMutationResult, error) {
		return s.ResumeSession(governance.LocalOwner(), governance.SessionResumeRequest{Schema: governance.SessionSchema, ProjectRef: identity.ProjectRef, CheckoutID: identity.CheckoutID, OperationID: operation, InstanceID: v.InstanceID, ExpectedRevision: revision, Reason: *reason})
	}
	switch {
	case action == "resume":
		changed, e = resume(v.Revision, op)
	case action == "select":
		changed, e = s.SelectSession(governance.LocalOwner(), governance.SessionSelectRequest{Schema: governance.SessionSchema, ProjectRef: identity.ProjectRef, CheckoutID: identity.CheckoutID, OperationID: op, InstanceID: v.InstanceID, ExpectedRevision: v.Revision, StepID: *step, Reason: *reason})
	case advancing || deciding:
		selected, err := terminalSelectedStep(v)
		if err != nil {
			return err
		}
		run, err := s.Read(selected.RunID)
		if err != nil {
			return err
		}
		if deciding {
			g := selected.OwnerGate
			if g == nil || g.Action != "request" || g.GateID != *gate {
				return fmt.Errorf("conflict: exact currently displayed waiting owner gate required")
			}
			changed, e = s.OwnerSession(governance.LocalOwner(), governance.SessionOwnerRequest{Schema: governance.SessionSchema, ProjectRef: identity.ProjectRef, CheckoutID: identity.CheckoutID, OperationID: op, InstanceID: v.InstanceID, ExpectedRevision: v.Revision, ExpectedRunRevision: run.Run.Revision, Action: action, GateID: *gate, ResultID: g.ResultID, Reason: *reason})
		} else {
			b, r := "disposition", selected.DispositionID
			if action == "rework" {
				b, r = *basis, *result
			}
			changed, e = s.AdvanceSession(governance.LocalOwner(), governance.SessionAdvanceRequest{Schema: governance.SessionSchema, ProjectRef: identity.ProjectRef, CheckoutID: identity.CheckoutID, OperationID: op, InstanceID: v.InstanceID, ExpectedRevision: v.Revision, ExpectedRunRevision: run.Run.Revision, Action: action, StepID: *step, Basis: b, ResultID: r, Reason: *reason})
		}
	}
	if e != nil {
		return e
	}
	if action != "resume" && !*prepareOnly {
		changed, e = resume(changed.Revision, governancehost.UUID())
		if e != nil {
			return fmt.Errorf("prepared operation %s; inspect and use session resume to reconcile: %w", op, e)
		}
	}
	if *jsonOut {
		return json.NewEncoder(out).Encode(changed)
	}
	v, e = s.ReadSession(v.InstanceID)
	if e != nil {
		return e
	}
	return emitTerminalSession(s, v, out, false)
}

func selectTerminalSession(list governance.SessionList, session, instance string) (governance.SessionView, error) {
	var matches []governance.SessionView
	for _, v := range list.Instances {
		if (instance != "" && v.InstanceID == instance) || (session != "" && v.SessionID == session) {
			matches = append(matches, v)
		}
	}
	if len(matches) != 1 {
		return governance.SessionView{}, fmt.Errorf("Session selector matched %d instances; use session list and an exact --instance", len(matches))
	}
	return matches[0], nil
}

func terminalSelectedStep(v governance.SessionView) (governance.SessionStepView, error) {
	for _, step := range v.Steps {
		if step.Selected {
			return step, nil
		}
	}
	return governance.SessionStepView{}, fmt.Errorf("no selected current Session step")
}

func emitTerminalSession(s *governance.Service, v governance.SessionView, out io.Writer, jsonOut bool) error {
	if jsonOut {
		return json.NewEncoder(out).Encode(v)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Session %q: %s\nPlan: %q\nRegistration: %q; record effects: %q; coverage: %q\n", v.SessionID, v.State, v.PlanID, v.Registration, v.RecordEffects, v.Coverage)
	fmt.Fprintln(&b, "Snapshot of recorded work; host completion and owner acceptance are separate.")
	for _, c := range v.Conditions {
		fmt.Fprintf(&b, "Condition: %q\n", c)
	}
	for _, step := range v.Steps {
		fmt.Fprintf(&b, "\nStep %q: %q; selected: %t\n", step.ID, step.Stage, step.Selected)
		if step.RunID != "" {
			fmt.Fprintf(&b, "  Run: %q; state: %q\n", step.RunID, step.RunState)
		}
		for _, hold := range step.Holds {
			fmt.Fprintf(&b, "  Hold: %q\n", hold)
		}
		if g := step.OwnerGate; g != nil {
			fmt.Fprintf(&b, "  Owner gate: %q; decision: %q; exact result: %q\n", g.GateID, g.Action, g.ResultID)
		}
		if step.DispositionID != "" {
			fmt.Fprintf(&b, "  Disposition: %q\n", step.DispositionID)
		}
	}
	fmt.Fprintf(&b, "\nAvailable core operations: %q\nInstance: %q; revision: %d; retained prior executions: %d\n", v.NextActions, v.InstanceID, v.Revision, len(v.History))
	if _, e := io.WriteString(out, b.String()); e != nil {
		return e
	}
	step, e := terminalSelectedStep(v)
	if e != nil || step.RunID == "" {
		return nil
	}
	view, e := s.Overview()
	if e != nil {
		return e
	}
	for _, r := range view.Runs {
		if r.RunID == step.RunID {
			_, e = fmt.Fprintf(out, "Role: %q; last host status: %q at %q; observation gap: %t; observed tool completions: %d; skill acceptance: %q; compliance: %q\n", r.AssignedRole, r.LastHostStatus, r.LastObservedAt, r.ObservationGap, r.ToolCompletions, r.SkillAcceptance, r.SkillCompliance)
			return e
		}
	}
	return nil
}

func terminalSessionAttention(s *governance.Service, v governance.SessionView, action, question, owner string, jsonOut bool, out io.Writer, input io.Reader, prompt io.Writer) error {
	selected, e := s.SessionHostSelection(v.InstanceID)
	if e != nil {
		return e
	}
	inbox, e := s.AttentionInbox(selected.RunID)
	if e != nil {
		return e
	}
	if action == "questions" {
		if jsonOut {
			return json.NewEncoder(out).Encode(inbox)
		}
		return renderAttention(out, inbox)
	}
	var matches []governance.AttentionView
	for _, item := range inbox.Items {
		if (question != "" && item.Question.ID == question) || (question == "" && item.Status == "open" && item.Applicability == "current") {
			matches = append(matches, item)
		}
	}
	if len(matches) != 1 {
		return fmt.Errorf("answer matched %d questions; inspect session questions and supply --question ID", len(matches))
	}
	q := matches[0].Question
	r, e := captureAttentionReply(s, governance.AttentionDisposition{Schema: governance.AttentionSchema, ProjectRef: inbox.ProjectRef, CheckoutID: inbox.CheckoutID, OperationID: governancehost.UUID(), RunID: selected.RunID, QuestionID: q.ID, Basis: q.Basis}, owner, input, prompt)
	if e != nil {
		return e
	}
	if jsonOut {
		return json.NewEncoder(out).Encode(r)
	}
	_, e = fmt.Fprintf(out, "Reply recorded for question %q in Session %q. The current controller may continue within existing authority.\n", q.ID, v.SessionID)
	return e
}
