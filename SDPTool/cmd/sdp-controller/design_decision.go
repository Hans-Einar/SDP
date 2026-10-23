package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func captureDesignDecision(service *governance.Service, q governance.DesignDecisionRequest, owner string, input io.Reader, prompt io.Writer) (governance.DesignDecisionResult, error) {
	var empty governance.DesignDecisionResult
	old, err := service.DesignDecisionReceipt(governance.LocalController(), q)
	if err != nil {
		return empty, err
	}
	if old != nil {
		if old.Record.Input.Actor != owner || old.Record.Input.Source != "local-controller-stdin" {
			return empty, fmt.Errorf("idempotency_conflict: receipt belongs to different input attribution")
		}
		fmt.Fprintln(prompt, "Returning historical decision receipt; this does not revalidate current acceptance.")
		return *old, nil
	}
	review, err := service.CheckAssignmentDesign(q.Review)
	if err != nil {
		return empty, err
	}
	if review.Status != "current" || review.ReviewRevision != q.ReviewRevision {
		return empty, fmt.Errorf("stale_revision: reviewed design is no longer current")
	}
	p := review.Assignment
	// Quote every project-controlled string, including terminal escape characters.
	fmt.Fprintf(prompt, "Design decision for owner %q\nOutcome: %q\nTask: %q\nNOW: %q\nTARGET: %q\n", owner, p.Assignment.Outcome, review.Provider.Task.Intent, review.Provider.Now.ID, review.Provider.Target.ID)
	for _, v := range p.Assignment.WritablePaths {
		fmt.Fprintf(prompt, "May change: %q\n", v)
	}
	for _, v := range p.Assignment.ExcludedPaths {
		fmt.Fprintf(prompt, "Excluded: %q\n", v)
	}
	for _, v := range p.ReturnCriteria {
		fmt.Fprintf(prompt, "Return criterion: %q\n", v)
	}
	for _, v := range p.Artifacts {
		fmt.Fprintf(prompt, "Review source (%s): %q\n", v.Kind, v.Path)
	}
	fmt.Fprintf(prompt, "Preliminary models: NOW=%t TARGET=%t. Blueprint evidence readiness: %q.\n", review.Provider.Now.Preliminary, review.Provider.Target.Preliminary, review.Provider.Readiness)
	fmt.Fprintf(prompt, "Model coverage: %q\nUnresolved model observations: %q\n", fmt.Sprint(review.Provider.Coverage), fmt.Sprint(review.Provider.Unknowns))
	fmt.Fprintln(prompt, "Review the proposal and sources before deciding. This records your decision; it does not launch work.")
	fmt.Fprintln(prompt, "Enter approved, changes-requested, rejected, or cancel (no default):")
	lines := bufio.NewScanner(input)
	lines.Buffer(make([]byte, 256), 4096)
	if !lines.Scan() {
		return empty, fmt.Errorf("no decision recorded: input ended: %v", lines.Err())
	}
	decision := strings.TrimSpace(lines.Text())
	if decision == "cancel" || decision == "" {
		return empty, fmt.Errorf("no decision recorded: canceled")
	}
	if decision != "approved" && decision != "changes-requested" && decision != "rejected" {
		return empty, fmt.Errorf("no decision recorded: invalid choice")
	}
	fmt.Fprintln(prompt, "Enter a comment describing your decision (required):")
	if !lines.Scan() {
		return empty, fmt.Errorf("no decision recorded: comment input ended: %v", lines.Err())
	}
	return service.RecordDesignDecision(governance.LocalController(), q, governance.OwnerDesignInput{Actor: owner, Source: "local-controller-stdin", Decision: decision, Comment: lines.Text()})
}
