package governance

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"slices"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool/blueprints"
	"github.com/Hans-Einar/SDP/SDPTool/blueprintstate"
	"github.com/Hans-Einar/SDP/SDPTool/model"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/blueprint"
)

// SessionDesignRequest selects sources to inspect, never approval or authority.
// ModelArea is relative to the SDP area; TaskPath is project-relative.
type SessionDesignRequest struct {
	Schema              string `json:"schema"`
	InstanceID          string `json:"instanceId"`
	ExpectedRevision    uint64 `json:"expectedRevision"`
	AssessmentID        string `json:"assessmentId"`
	ModelArea           string `json:"modelArea"`
	WorkRef             string `json:"workRef"`
	Entry               string `json:"entry"`
	TaskPath            string `json:"taskPath"`
	BlueprintAssignment string `json:"blueprintAssignment,omitempty"`
	BlueprintEvent      string `json:"blueprintEvent,omitempty"`
}

type SessionDesignBasis struct {
	Request       SessionDesignRequest    `json:"request"`
	BindingDigest string                  `json:"bindingDigest"`
	Assessment    SessionAssessmentView   `json:"assessment"`
	Scope         SessionStepBinding      `json:"scope"`
	Work          model.Validation        `json:"work"`
	TaskDigest    string                  `json:"taskDigest"`
	Task          blueprint.Task          `json:"task"`
	Blueprint     *DesignProviderSnapshot `json:"blueprint,omitempty"`
}

type SessionDesignContext struct {
	Schema           string              `json:"schema"`
	Operation        string              `json:"operation"`
	InstanceID       string              `json:"instanceId"`
	Revision         uint64              `json:"revision"`
	Status           string              `json:"status"`
	Basis            *SessionDesignBasis `json:"basis,omitempty"`
	ReviewRevision   string              `json:"reviewRevision,omitempty"`
	Holds            []string            `json:"holds"`
	DesignAcceptance string              `json:"designAcceptance"`
	Dispatch         string              `json:"dispatch"`
}

// InspectSessionDesign builds an exact, read-only review basis. A current result
// does not authorize implementation, accept design or waive model release.
func (s *Service) InspectSessionDesign(q SessionDesignRequest) (SessionDesignContext, error) {
	return s.inspectSessionDesign(q, nil)
}

func (s *Service) inspectSessionDesign(q SessionDesignRequest, beforeRecheck func()) (SessionDesignContext, error) {
	out := SessionDesignContext{Schema: SessionSchema, Operation: "session-design-context", InstanceID: q.InstanceID,
		Revision: q.ExpectedRevision, Status: "blocked", Holds: []string{}, DesignAcceptance: "not-recorded", Dispatch: "unavailable"}
	if q.Schema != SessionSchema || !uuidPattern.MatchString(q.InstanceID) || q.ExpectedRevision == 0 || !uuidPattern.MatchString(q.AssessmentID) ||
		!relative(q.ModelArea) || !relative(q.Entry) || !strings.HasSuffix(q.Entry, ".design") || !relative(q.TaskPath) || !strings.HasPrefix(q.WorkRef, "work:") || len(q.WorkRef) > 256 {
		return out, fail("invalid_request", "exact Session/assessment and contained SDL WORK, entry and task required")
	}
	if (q.BlueprintAssignment == "") != (q.BlueprintEvent == "") || (q.BlueprintAssignment != "" && (!blueprintstate.ID.MatchString(q.BlueprintAssignment) || !validLeaderText(q.BlueprintEvent, 4096))) {
		return out, fail("invalid_request", "canonical Blueprint assignment and event must be supplied together")
	}
	st, head, err := s.load()
	if err != nil {
		return out, err
	}
	i, ok := st.Sessions[q.InstanceID]
	if !ok {
		return out, fail("missing_coverage", "Session absent")
	}
	if i.Revision != q.ExpectedRevision {
		return out, fail("stale_revision", "Session revision changed")
	}
	blocked := func(e error) (SessionDesignContext, error) { out.Holds = append(out.Holds, e.Error()); return out, nil }
	first, err := s.sessionDesignBasis(st, i, q)
	if err != nil {
		return blocked(err)
	}
	if beforeRecheck != nil {
		beforeRecheck()
	}
	second, err := s.sessionDesignBasis(st, i, q)
	if err != nil {
		return blocked(err)
	}
	if !reflect.DeepEqual(first, second) {
		return blocked(fail("stale_revision", "design sources changed during inspection"))
	}
	_, last, err := s.load()
	if err != nil {
		return out, err
	}
	if last != head {
		return blocked(fail("stale_revision", "operational state changed during inspection"))
	}
	raw, err := json.Marshal(struct {
		Project, Checkout string
		Basis             SessionDesignBasis
	}{s.project, st.CheckoutID, second})
	if err != nil {
		return out, err
	}
	out.Status, out.Basis, out.ReviewRevision = "current", &second, hash(raw)
	return out, nil
}

func (s *Service) sessionDesignBasis(st snapshot, i SessionInstance, q SessionDesignRequest) (SessionDesignBasis, error) {
	var out SessionDesignBasis
	if !i.Registered || i.Selection == nil || i.Selection.Status != "committed" || len(i.Assessments) == 0 || len(i.Work) > 0 || i.Rebind != nil || i.Advance != nil || i.Disposition != nil {
		return out, fail("conflict", "design inspection requires a selected, assessed Session before implementation")
	}
	a := i.Assessments[len(i.Assessments)-1]
	step := selectedSessionStep(i)
	if a.Request.OperationID != q.AssessmentID || a.Effects.Status != "committed" {
		return out, fail("stale_revision", "latest committed assessment required")
	}
	impact := assessmentImpact(a.Request, step)
	if impact == "unknown" || a.Request.Blueprint.Requirement == "unknown" {
		return out, fail("missing_coverage", "resolve design and Blueprint assessment first")
	}
	if impact != "changed" && a.Request.Blueprint.Requirement != "required" {
		return out, fail("conflict", "no-design Session does not require this handoff")
	}
	if a.Request.Blueprint.Requirement == "required" && q.BlueprintAssignment == "" {
		return out, fail("missing_coverage", "assessment requires a retained Blueprint")
	}
	if err := s.sessionCurrent(i); err != nil {
		return out, err
	}
	if err := s.sessionWorkClaims(st, i); err != nil {
		return out, err
	}
	area := filepath.Join(s.area, filepath.FromSlash(q.ModelArea))
	if err := safeAbsolute(area); err != nil {
		return out, err
	}
	work, err := model.Validate(area, q.WorkRef)
	if err != nil {
		return out, fail("provider_unavailable", err.Error())
	}
	if work.Source.Artifact == nil || work.Source.Artifact.Kind != "work" || work.Source.Status != "clean" || work.Source.Artifact.Head == "" || !slices.Contains(work.Targets, q.Entry) {
		return out, fail("conflict", "committed, conflict-free SDL WORK and validated entry required")
	}
	raw, err := readFile(filepath.Join(s.root, filepath.FromSlash(q.TaskPath)), MaxRequest)
	if err != nil {
		return out, err
	}
	var task blueprint.Task
	// Task uses the Blueprint producer's existing untagged JSON contract.
	if err = decode(raw, &task); err != nil {
		return out, err
	}
	if !validLeaderText(task.ID, 256) || !validLeaderText(task.Intent, 4096) {
		return out, fail("invalid_request", "task identity and intent required")
	}
	paths := slices.Clone(task.AllowedCodePaths)
	slices.Sort(paths)
	writable := slices.Clone(step.WritablePaths)
	slices.Sort(writable)
	if !slices.Equal(paths, writable) {
		return out, fail("conflict", "task code permissions must exactly match selected Session scope")
	}
	av := assessmentView(a)
	av.DesignImpact = impact
	out = SessionDesignBasis{Request: q, BindingDigest: i.Context.BindingDigest, Assessment: *av, Scope: step, Work: work, TaskDigest: hash(raw), Task: task}
	if q.BlueprintAssignment != "" {
		provider, err := s.sessionDesignBlueprint(i, out)
		if err != nil {
			return SessionDesignBasis{}, err
		}
		out.Blueprint = &provider
	}
	return out, nil
}

func (s *Service) sessionDesignBlueprint(i SessionInstance, basis SessionDesignBasis) (DesignProviderSnapshot, error) {
	var out DesignProviderSnapshot
	q := basis.Request
	views, err := blueprints.AssignmentViews(s.area)
	if err != nil {
		return out, fail("provider_unavailable", err.Error())
	}
	var selected *blueprints.AssignmentView
	for n := range views {
		if views[n].ID == q.BlueprintAssignment {
			selected = &views[n]
			break
		}
	}
	if selected == nil {
		return out, fail("missing_coverage", "canonical Blueprint assignment absent")
	}
	v := *selected
	if v.EventID != q.BlueprintEvent || v.Binding.Plan != i.Context.Binding.Plan.ID || v.Binding.Milestone != basis.Scope.Milestone || v.Binding.ModelArea != q.ModelArea || v.Binding.To != q.WorkRef {
		return out, fail("stale_revision", "Blueprint does not bind the exact Session plan, milestone and WORK")
	}
	if !slices.Contains([]string{"draft", "ready", "assigned"}, v.State.State) || v.Validation != "validated" || v.Freshness != "current" {
		return out, fail("conflict", "Blueprint is outside current pre-implementation review: "+v.Validation+"/"+v.Freshness)
	}
	raw, err := readFile(filepath.Join(s.area, filepath.FromSlash(v.Binding.Bundle), "blueprint.json"), MaxRequest)
	if err != nil {
		return out, err
	}
	var doc blueprints.Document
	if err = json.Unmarshal(raw, &doc); err != nil {
		return out, err
	}
	if doc.Analysis == nil || !doc.Analysis.ConstraintsPass || doc.Revision != v.Binding.Revision || doc.Entry != q.Entry || doc.TaskID != basis.Task.ID || doc.TaskDigest != basis.TaskDigest {
		return out, fail("conflict", "Blueprint constraints, entry or exact task differ")
	}
	source := basis.Work.Source
	capture := blueprints.Capture{ID: source.Artifact.ID, Kind: source.Artifact.Kind, Head: source.Artifact.Head, MetadataDigest: source.Artifact.MetadataDigest, SourceDigest: source.LiveDigest, Preliminary: source.Preliminary}
	if capture != doc.Target {
		return out, fail("stale_revision", "Blueprint target differs from validated WORK")
	}
	out = DesignProviderSnapshot{AssignmentID: v.ID, EventID: v.EventID, WorkState: v.State.State, Binding: v.Binding, Now: doc.Now, Target: doc.Target, TaskDigest: doc.TaskDigest, Task: basis.Task, Readiness: v.Readiness, Coverage: doc.Analysis.Coverage, Unknowns: doc.Analysis.Unknowns}
	return out, nil
}
