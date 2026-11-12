package activities

import (
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type pendingRequest struct {
	Request TransitionRequest `json:"request"`
	Caller  Caller            `json:"caller"`
	Digest  string            `json:"digest"`
}

func pendingRequestPath(id string) string { return "Activities/.requests/" + id + ".json" }
func (s *Service) checkActivity(c Catalog, id, event, descriptor string) (Activity, error) {
	a, ok := c.Find(id)
	if !ok {
		return a, fmt.Errorf("unknown Activity")
	}
	if a.LastEventID != event || a.DescriptorRef.SHA256 != descriptor {
		return a, projecthistory.ErrStale
	}
	return a, nil
}
func (s *Service) descriptorImage(a Activity) ([]image, error) {
	b, e := s.read(a.DescriptorRef.Path)
	if e != nil || digest(b) != a.DescriptorRef.SHA256 {
		return nil, fmt.Errorf("live descriptor differs from current pin")
	}
	return []image{{Path: "Activities/.descriptors/" + a.DescriptorRef.SHA256 + ".yaml", After: b}}, nil
}
func (s *Service) Transition(q TransitionRequest, caller Caller, inspector ExecutionInspector) (Result, error) {
	if e := s.validateCaller(caller); e != nil {
		return Result{}, e
	}
	if !uuidRE.MatchString(q.OperationID) || !hashRE.MatchString(q.ExpectedHistoryRevision) || strings.TrimSpace(q.Reason) == "" {
		return Result{}, fmt.Errorf("operation, revision and reason required")
	}
	raw, _ := json.Marshal(struct {
		Request any
		Caller  Caller
	}{q, caller})
	dg := digest(raw)
	ledger, e := s.path("ProjectManagement/Ledger.ndjson")
	if e != nil {
		return Result{}, e
	}
	tx, e := projecthistory.BeginTransaction(ledger)
	if e != nil {
		return Result{}, e
	}
	defer tx.Close()
	owner := "activity:" + q.OperationID
	old, e := s.load(q.OperationID)
	if e == nil {
		if old.Digest != dg {
			return Result{}, projecthistory.ErrConflict
		}
		return s.apply(tx, old)
	}
	if !os.IsNotExist(e) {
		return Result{}, e
	}
	snap, e := tx.Read(owner)
	if e != nil {
		return Result{}, e
	}
	view, e := s.ReadAdmission(tx, owner)
	if e != nil {
		return Result{}, e
	}
	waiting := false
	if view.Pending != nil {
		if view.Pending.OperationID != q.OperationID {
			return Result{}, fmt.Errorf("another lifecycle request pending")
		}
		b, e := s.read(pendingRequestPath(q.OperationID))
		if e != nil {
			return Result{}, e
		}
		var p pendingRequest
		if e = Decode(b, &p); e != nil {
			return Result{}, e
		}
		if p.Digest != dg {
			return Result{}, projecthistory.ErrConflict
		}
		waiting = true
	}
	if !waiting && snap.Revision != q.ExpectedHistoryRevision {
		return Result{}, projecthistory.ErrStale
	}
	c, e := Replay(s.Area, s.ProjectID, snap)
	if e != nil {
		return Result{}, e
	}
	a, e := s.checkActivity(c, q.ActivityID, q.ExpectedActivityEventID, q.ExpectedDescriptorRevision)
	if e != nil {
		return Result{}, e
	}
	target := ""
	switch q.Action {
	case "activate":
		if a.State != "planned" {
			return Result{}, fmt.Errorf("activate requires planned")
		}
		target = "active"
	case "resume":
		if a.State != "on-hold" {
			return Result{}, fmt.Errorf("resume requires on-hold")
		}
		target = "active"
	case "hold":
		if strings.TrimSpace(q.ResumeCondition) == "" {
			return Result{}, fmt.Errorf("hold requires explicit resume condition")
		}
		if a.State != "active" {
			return Result{}, fmt.Errorf("hold requires active")
		}
		target = "on-hold"
	case "complete":
		if a.State != "active" || strings.TrimSpace(q.Outcome) == "" || len(q.Evidence) == 0 {
			return Result{}, fmt.Errorf("complete requires active Activity, explicit outcome and evidence")
		}
		for _, ref := range q.Evidence {
			if ref.Availability == "unavailable" {
				return Result{}, fmt.Errorf("unavailable reference is not completion evidence")
			}
		}
		if e = validateRefs(s, q.Evidence); e != nil {
			return Result{}, e
		}
		if e = validateContributionDispositions(a.Descriptor, q.Action, q.ContributionDispositions); e != nil {
			return Result{}, e
		}
		if e = s.checkWholeScopeCompletion(a.Descriptor, q.ContributionDispositions); e != nil {
			return Result{}, e
		}
		target = "completed"
	case "cancel":
		if a.State != "planned" && a.State != "active" && a.State != "on-hold" {
			return Result{}, fmt.Errorf("cannot cancel terminal Activity")
		}
		if e = validateContributionDispositions(a.Descriptor, q.Action, q.ContributionDispositions); e != nil {
			return Result{}, e
		}
		target = "canceled"
	case "reopen":
		if a.State != "completed" && a.State != "canceled" {
			return Result{}, fmt.Errorf("reopen requires terminal Activity")
		}
		for _, other := range c.Activities {
			if other.State == "active" {
				for _, dep := range other.Descriptor.Dependencies {
					if dep.ActivityID == a.Descriptor.ID {
						return Result{}, fmt.Errorf("hold dependent active Activity before reopening its dependency")
					}
				}
			}
		}
		target = "planned"
	case "switch":
		if strings.TrimSpace(q.ResumeCondition) == "" {
			return Result{}, fmt.Errorf("switch requires explicit former Activity resume condition")
		}
		if a.State != "active" || q.TargetActivityID == a.Descriptor.ID {
			return Result{}, fmt.Errorf("switch requires distinct active source")
		}
		target = "on-hold"
	default:
		return Result{}, fmt.Errorf("unknown lifecycle action")
	}
	var selected Activity
	if q.Action == "switch" {
		selected, e = s.checkActivity(c, q.TargetActivityID, q.ExpectedTargetEventID, q.ExpectedTargetDescriptorRevision)
		if e != nil {
			return Result{}, e
		}
		if e = checkDependencySatisfaction(c, selected); e != nil {
			return Result{}, e
		}
		if selected.State != "planned" && selected.State != "on-hold" {
			return Result{}, fmt.Errorf("switch target must be planned or held")
		}
	} else if q.TargetActivityID != "" || q.ExpectedTargetEventID != "" || q.ExpectedTargetDescriptorRevision != "" {
		return Result{}, fmt.Errorf("target fields only valid for switch")
	}
	if target == "active" {
		if e = checkDependencySatisfaction(c, a); e != nil {
			return Result{}, e
		}
		for _, other := range c.Activities {
			if other.State == "active" && other.Descriptor.ID != a.Descriptor.ID {
				return Result{}, fmt.Errorf("another Activity is active; use switch")
			}
		}
	}
	if a.State == "active" && target != "active" {
		if !view.Enabled || inspector == nil || len(view.Stores) == 0 {
			return Result{}, fmt.Errorf("active exit requires explicit adopted execution inventory and inspector")
		}
		observed, e := inspector.InspectExecutions(view.Stores)
		if e != nil {
			return Result{}, e
		}
		if !observed.Known {
			return Result{}, fmt.Errorf("execution disposition unknown")
		}
		if len(observed.InFlight) > 0 {
			p := pendingRequest{Request: q, Caller: caller, Digest: dg}
			b, _ := json.Marshal(p)
			old, e := s.read(pendingRequestPath(q.OperationID))
			if e == nil {
				var prior pendingRequest
				if e = Decode(old, &prior); e != nil {
					return Result{}, e
				}
				if prior.Digest != dg {
					return Result{}, projecthistory.ErrConflict
				}
			} else if os.IsNotExist(e) {
				if e = s.write(pendingRequestPath(q.OperationID), b); e != nil {
					return Result{}, e
				}
			} else {
				return Result{}, e
			}
			if e = s.SetPendingLifecycle(tx, owner, PendingLifecycle{OperationID: q.OperationID, ActivityID: a.Descriptor.ID, TargetActivityID: q.TargetActivityID, Action: q.Action}); e != nil {
				return Result{}, e
			}
			return Result{OperationID: q.OperationID, ActivityID: a.Descriptor.ID, Path: a.Path, HistoryRevision: snap.Revision, DescriptorRevision: a.DescriptorRef.SHA256, Pending: true}, nil
		}
	}
	files, e := s.descriptorImage(a)
	if e != nil {
		return Result{}, e
	}
	ev := s.event(q.OperationID, "x-activity:transitioned", a.Descriptor.ID, q.Reason, caller, dg, a.DescriptorRef, &a.LastEventID, &a.State, target, nil)
	ev.Payload.ResumeCondition = q.ResumeCondition
	ev.Payload.ContributionDispositions = q.ContributionDispositions
	ev.Payload.Outcome = q.Outcome
	ev.Payload.Evidence = q.Evidence
	in := makeIntent(s.ProjectID, q.OperationID, dg, snap.Revision, files, ev, Result{OperationID: q.OperationID, ActivityID: a.Descriptor.ID, Path: a.Path, DescriptorRevision: a.DescriptorRef.SHA256})
	in.Pending = waiting
	if q.Action == "switch" {
		otherFiles, e := s.descriptorImage(selected)
		if e != nil {
			return Result{}, e
		}
		in.Files = append(in.Files, otherFiles...)
		ev.Payload.Batch = &BatchReference{Index: 0, Count: 2}
		second := s.event(q.OperationID, "x-activity:transitioned", selected.Descriptor.ID, q.Reason, caller, dg, selected.DescriptorRef, &selected.LastEventID, &selected.State, "active", nil)
		second.EventID += "-TARGET"
		second.Payload.Batch = &BatchReference{Index: 1, Count: 2}
		firstBytes, _ := json.Marshal(ev)
		secondBytes, _ := json.Marshal(second)
		in.Events = [][]byte{firstBytes, secondBytes}
		in.Result.ActivityID = selected.Descriptor.ID
		in.Result.Path = selected.Path
		in.Result.EventID = second.EventID
		in.Result.DescriptorRevision = selected.DescriptorRef.SHA256
	}
	return s.publish(tx, in)
}
func (s *Service) Update(q UpdateRequest, caller Caller) (Result, error) {
	tx, dg, prior, e := s.begin(q.OperationID, q.ExpectedHistoryRevision, q, caller)
	if e != nil {
		return Result{}, e
	}
	if prior != nil {
		return *prior, nil
	}
	defer tx.Close()
	view, e := s.ReadAdmission(tx, "activity:"+q.OperationID)
	if e != nil {
		return Result{}, e
	}
	if view.Pending != nil {
		return Result{}, fmt.Errorf("lifecycle request pending")
	}
	if strings.TrimSpace(q.Reason) == "" || strings.TrimSpace(q.Title) == "" || strings.TrimSpace(q.Owner) == "" || q.Scope == nil || q.References == nil || q.SourceCards == nil {
		return Result{}, fmt.Errorf("explicit complete descriptor update required")
	}
	if e = validateRefs(s, q.References); e != nil {
		return Result{}, e
	}
	if e = s.validateCards(q.SourceCards); e != nil {
		return Result{}, e
	}
	for _, p := range q.Scope {
		if _, e = s.path(p); e != nil {
			return Result{}, e
		}
	}
	snap, e := tx.Read("activity:" + q.OperationID)
	if e != nil {
		return Result{}, e
	}
	c, e := Replay(s.Area, s.ProjectID, snap)
	if e != nil {
		return Result{}, e
	}
	a, e := s.checkActivity(c, q.ActivityID, q.ExpectedActivityEventID, q.ExpectedDescriptorRevision)
	if e != nil {
		return Result{}, e
	}
	files, e := s.descriptorImage(a)
	if e != nil {
		return Result{}, e
	}
	d := a.Descriptor
	d.Schema = "sdp-activity/0.2"
	d.Title = q.Title
	d.Owner = q.Owner
	d.Scope = q.Scope
	d.References = q.References
	d.SourceCards = q.SourceCards
	d.Systems = q.Systems
	d.Dependencies = q.Dependencies
	if a.State == "active" {
		if e = checkDependencySatisfaction(c, Activity{Descriptor: d}); e != nil {
			return Result{}, e
		}
	}
	if e = validateDependencies(c, d); e != nil {
		return Result{}, e
	}
	b, e := yaml.Marshal(d)
	if e != nil {
		return Result{}, e
	}
	if _, e = parseDescriptor(b); e != nil {
		return Result{}, e
	}
	old, e := s.read(a.DescriptorRef.Path)
	if e != nil {
		return Result{}, e
	}
	ref := DescriptorReference{Path: a.DescriptorRef.Path, SHA256: digest(b)}
	files = append(files, image{Path: "Activities/.descriptors/" + ref.SHA256 + ".yaml", After: b}, image{Path: ref.Path, Before: old, After: b})
	ev := s.event(q.OperationID, "x-activity:updated", a.Descriptor.ID, q.Reason, caller, dg, ref, &a.LastEventID, &a.State, a.State, nil)
	return s.publish(tx, makeIntent(s.ProjectID, q.OperationID, dg, snap.Revision, files, ev, Result{OperationID: q.OperationID, ActivityID: a.Descriptor.ID, Path: a.Path, DescriptorRevision: ref.SHA256}))
}

// RecoverLifecycle can observe actual quiescence for a retained pending request;
// ordinary Recover returns its pending disposition without claiming work stopped.
func (s *Service) RecoverLifecycle(id string, caller Caller, inspector ExecutionInspector) (Result, error) {
	if !uuidRE.MatchString(id) {
		return Result{}, fmt.Errorf("invalid operation UUID")
	}
	b, e := s.read(pendingRequestPath(id))
	if os.IsNotExist(e) {
		return s.Recover(id, caller)
	}
	if e != nil {
		return Result{}, e
	}
	var p pendingRequest
	if e = Decode(b, &p); e != nil {
		return Result{}, e
	}
	return s.Transition(p.Request, caller, inspector)
}

func validateContributionDispositions(d Descriptor, action string, dispositions []ContributionDisposition) error {
	selected := map[string]bool{}
	for _, ref := range d.References {
		if ref.Membership != "" {
			selected[ref.ID] = true
		}
	}
	seen := map[string]bool{}
	for _, v := range dispositions {
		if !selected[v.ReferenceID] || seen[v.ReferenceID] || strings.TrimSpace(v.Outcome) == "" {
			return fmt.Errorf("invalid contribution disposition")
		}
		seen[v.ReferenceID] = true
		switch v.Disposition {
		case "completed", "not-applicable":
		case "returned-to-backlog", "transferred":
			if action != "cancel" {
				return fmt.Errorf("unmet selected contribution requires explicit scope change before completion")
			}
		default:
			return fmt.Errorf("unknown contribution disposition")
		}
	}
	for id := range selected {
		if !seen[id] {
			return fmt.Errorf("selected contribution %s requires explicit disposition", id)
		}
	}
	return nil
}

func (s *Service) checkWholeScopeCompletion(d Descriptor, dispositions []ContributionDisposition) error {
	byID := map[string]string{}
	for _, v := range dispositions {
		byID[v.ReferenceID] = v.Disposition
	}
	for _, ref := range d.References {
		if ref.Membership != "whole-scope" || byID[ref.ID] != "completed" {
			continue
		}
		if ref.Availability == "unavailable" {
			return fmt.Errorf("unavailable whole-scope input cannot prove completion")
		}
		b, e := s.read(ref.Path)
		if e != nil {
			return e
		}
		satisfied := false
		for _, line := range strings.Split(string(b), "\n") {
			parts := strings.Split(line, "|")
			if len(parts) < 4 {
				continue
			}
			key := strings.ToLower(strings.TrimSpace(parts[1]))
			value := strings.ToLower(strings.TrimSpace(parts[2]))
			if key == "state" || key == "status" || key == "cardstate" {
				switch value {
				case "completed", "closed", "done", "delivered", "accepted":
					satisfied = true
				}
			}
		}
		if !satisfied {
			return fmt.Errorf("whole-scope reference %s has no authoritative completed disposition", ref.ID)
		}
	}
	return nil
}
