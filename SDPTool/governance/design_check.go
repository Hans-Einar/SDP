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

const DesignCheckSchema = "sdp-governance-design-check/1"

type DesignCheckRequest struct {
	Schema              string `json:"schema"`
	AssignmentID        string `json:"assignmentId"`
	ExpectedRevision    uint64 `json:"expectedRevision"`
	ExpectedGeneration  string `json:"expectedGeneration"`
	BlueprintAssignment string `json:"blueprintAssignment"`
	BlueprintEvent      string `json:"blueprintEvent"`
}

// DesignProviderSnapshot retains observations from providers, not a new work state.
type DesignProviderSnapshot struct {
	AssignmentID string                 `json:"assignmentId"`
	EventID      string                 `json:"eventId"`
	WorkState    string                 `json:"workState"`
	Binding      blueprintstate.Binding `json:"binding"`
	Now          blueprints.Capture     `json:"now"`
	Target       blueprints.Capture     `json:"target"`
	TaskDigest   string                 `json:"taskDigest"`
	Task         blueprint.Task         `json:"task"`
	Readiness    string                 `json:"evidenceReadiness"`
	Coverage     blueprint.Coverage     `json:"coverage"`
	Unknowns     []blueprint.Unknown    `json:"unknowns"`
}

type DesignCheckResult struct {
	Schema           string                  `json:"schema"`
	Status           string                  `json:"status"`
	Generation       string                  `json:"generation"`
	Assignment       PreparedAssignment      `json:"assignment"`
	Provider         *DesignProviderSnapshot `json:"provider,omitempty"`
	ReviewRevision   string                  `json:"reviewRevision,omitempty"`
	Conditions       []string                `json:"conditions"`
	DesignAcceptance string                  `json:"designAcceptance"`
	Dispatch         string                  `json:"dispatch"`
}

// CheckAssignmentDesign builds a read-only review basis. It cannot approve or launch.
func (s *Service) CheckAssignmentDesign(q DesignCheckRequest) (DesignCheckResult, error) {
	return s.checkAssignmentDesign(q, nil)
}

func (s *Service) checkAssignmentDesign(q DesignCheckRequest, beforeRecheck func()) (DesignCheckResult, error) {
	out := DesignCheckResult{Schema: DesignCheckSchema, Status: "blocked", Conditions: []string{},
		DesignAcceptance: "not-recorded", Dispatch: "unavailable"}
	if q.Schema != DesignCheckSchema {
		return out, fail("incompatible_profile", "unsupported design check schema")
	}
	if !uuidPattern.MatchString(q.AssignmentID) || q.ExpectedRevision == 0 || !digestPattern.MatchString(q.ExpectedGeneration) ||
		!blueprintstate.ID.MatchString(q.BlueprintAssignment) || strings.TrimSpace(q.BlueprintEvent) == "" || len(q.BlueprintEvent) > 4096 {
		return out, fail("invalid_request", "pinned preparation and canonical Blueprint assignment/event required")
	}
	state, head, err := s.load()
	if err != nil {
		return out, err
	}
	p, ok := state.Prepared[q.AssignmentID]
	if !ok {
		return out, fail("missing_coverage", "prepared assignment not found")
	}
	if head != q.ExpectedGeneration || p.Revision != q.ExpectedRevision {
		return out, fail("stale_revision", "preparation or operational generation changed")
	}
	out.Generation, out.Assignment = head, p
	blocked := func(err error) (DesignCheckResult, error) {
		out.Provider, out.ReviewRevision = nil, ""
		out.Conditions = append(out.Conditions, err.Error())
		return out, nil
	}
	if p.Role != "Master" || p.State == "retired" {
		return blocked(fail("conflict", "design check requires an unretired prepared Master"))
	}
	first, err := s.designProvider(p, q)
	if err != nil {
		return blocked(err)
	}
	if beforeRecheck != nil {
		beforeRecheck()
	}
	second, err := s.designProvider(p, q)
	if err != nil {
		return blocked(err)
	}
	if !reflect.DeepEqual(first, second) {
		return blocked(fail("stale_revision", "provider observations changed during design check"))
	}
	_, finalHead, err := s.load()
	if err != nil {
		return out, err
	}
	if finalHead != head {
		return blocked(fail("stale_revision", "operational generation changed during design check"))
	}
	if err = s.checkPreparation(p.Role, p.Selection, p.InputRevision, p.Artifacts, p.Assignment); err != nil {
		return blocked(err)
	}
	// The generation guards this observation; it is not part of the review basis.
	// A future disposition write must not invalidate itself merely by advancing HEAD.
	raw, err := json.Marshal(struct {
		Schema, Project, Checkout string
		Assignment                PreparedAssignment
		Provider                  DesignProviderSnapshot
	}{DesignCheckSchema, s.project, state.CheckoutID, p, second})
	if err != nil {
		return out, err
	}
	out.Provider = &second
	out.ReviewRevision, out.Status = hash(raw), "current"
	return out, nil
}

func (s *Service) designProvider(p PreparedAssignment, q DesignCheckRequest) (DesignProviderSnapshot, error) {
	var out DesignProviderSnapshot
	if err := s.checkPreparation(p.Role, p.Selection, p.InputRevision, p.Artifacts, p.Assignment); err != nil {
		return out, err
	}
	c, err := s.contextFor(p.Selection, []string{"ImplementationPlan"}, "PLAN-")
	if err != nil {
		return out, err
	}
	views, err := blueprints.AssignmentViews(s.area)
	if err != nil {
		return out, fail("provider_unavailable", err.Error())
	}
	var selected *blueprints.AssignmentView
	for i := range views {
		if views[i].ID == q.BlueprintAssignment {
			selected = &views[i]
			break
		}
	}
	if selected == nil {
		return out, fail("missing_coverage", "canonical Blueprint assignment not found")
	}
	v := *selected
	if v.EventID != q.BlueprintEvent || v.Binding.Plan != c.PlanID {
		return out, fail("stale_revision", "Blueprint event or governing plan does not match")
	}
	if !slices.Contains([]string{"draft", "ready", "assigned"}, v.State.State) {
		return out, fail("conflict", "Blueprint work state is outside pre-implementation review")
	}
	if v.Validation != "validated" || v.Freshness != "current" {
		return out, fail("stale_revision", "Blueprint integrity/live-source check: "+v.Validation+"/"+v.Freshness+"; "+v.Diagnostic)
	}
	pins := map[string]ArtifactPin{}
	for _, pin := range p.Artifacts {
		pins[pin.Kind] = pin
	}
	relativePath := func(abs string) string {
		rel, err := filepath.Rel(s.root, abs)
		if err != nil || !relative(filepath.ToSlash(rel)) {
			return ""
		}
		return filepath.ToSlash(rel)
	}
	bundleDocument := filepath.Join(s.area, filepath.FromSlash(v.Binding.Bundle), "blueprint.json")
	if pins["blueprint"].Path != relativePath(bundleDocument) {
		return out, fail("conflict", "prepared blueprint pin does not name the canonical retained bundle")
	}
	raw, err := readFile(bundleDocument, MaxRequest)
	if err != nil || hash(raw) != pins["blueprint"].SHA256 {
		return out, fail("stale_revision", "pinned blueprint document unavailable or changed")
	}
	var doc blueprints.Document
	if err = json.Unmarshal(raw, &doc); err != nil || doc.Analysis == nil || doc.Revision != v.Binding.Revision {
		return out, fail("conflict", "verified blueprint document binding changed")
	}
	if !doc.Analysis.ConstraintsPass {
		return out, fail("conflict", "Blueprint structural constraints failed")
	}
	for _, side := range []struct {
		kind, ref string
		captured  blueprints.Capture
	}{
		{"now", v.Binding.From, doc.Now}, {"target", v.Binding.To, doc.Target},
	} {
		source, err := model.Snapshot(filepath.Join(s.area, filepath.FromSlash(v.Binding.ModelArea)), side.ref)
		if err != nil {
			return out, fail("provider_unavailable", err.Error())
		}
		if source.Artifact == nil || source.Status == "conflicted" {
			return out, fail("conflict", "model source is absent or conflicted")
		}
		capture := blueprints.Capture{ID: source.Artifact.ID, Kind: source.Artifact.Kind, Head: source.Artifact.Head,
			MetadataDigest: source.Artifact.MetadataDigest, SourceDigest: source.LiveDigest, Preliminary: source.Preliminary}
		if capture != side.captured {
			return out, fail("stale_revision", side.kind+" live model capture changed")
		}
		metadata := filepath.Join(source.Path, filepath.Base(source.Path)+".yaml")
		if pins[side.kind].Path != relativePath(metadata) {
			return out, fail("conflict", side.kind+" pin must name resolved model metadata")
		}
	}
	taskBytes, err := readFile(filepath.Join(s.root, filepath.FromSlash(pins["task"].Path)), MaxRequest)
	if err != nil || hash(taskBytes) != doc.TaskDigest || doc.TaskDigest != pins["task"].SHA256 {
		return out, fail("stale_revision", "prepared task does not match retained blueprint task")
	}
	var task blueprint.Task
	if err = json.Unmarshal(taskBytes, &task); err != nil || task.ID != v.Binding.Task {
		return out, fail("conflict", "retained task identity mismatch")
	}
	for _, path := range p.Assignment.WritablePaths {
		if !slices.Contains(task.AllowedCodePaths, path) {
			return out, fail("conflict", "prepared writable path is not an exact task permission: "+path)
		}
	}
	out = DesignProviderSnapshot{AssignmentID: v.ID, EventID: v.EventID, WorkState: v.State.State,
		Binding: v.Binding, Now: doc.Now, Target: doc.Target, TaskDigest: doc.TaskDigest, Task: task,
		Readiness: v.Readiness, Coverage: doc.Analysis.Coverage, Unknowns: doc.Analysis.Unknowns}
	return out, nil
}
