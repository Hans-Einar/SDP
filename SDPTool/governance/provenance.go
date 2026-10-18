package governance

import (
	"path/filepath"
	"strings"
)

// ContextProvenance is evidence of reported/snapshotted context, never a grant of
// role authority or proof of skill invocation, comprehension or compliance.
type ContextProvenance struct {
	AssignedRole  string         `json:"assignedRole"`
	RoleBasis     string         `json:"roleBasis"`
	InputRevision string         `json:"inputRevision"`
	Sources       []Source       `json:"sources"`
	Host          HostContext    `json:"host"`
	Skills        SkillInventory `json:"skills"`
}
type HostContext struct {
	Status             string   `json:"status"`
	SessionID          string   `json:"sessionId"`
	ParentThreadID     string   `json:"parentThreadId"`
	NativeRole         string   `json:"nativeRole"`
	Source             string   `json:"source"`
	InstructionStatus  string   `json:"instructionStatus"`
	InstructionSources []string `json:"instructionSources"`
}
type SkillInventory struct {
	Status        string              `json:"status"`
	UseEvidence   string              `json:"useEvidence"`
	ExternalCount int                 `json:"externalCount"`
	ErrorCount    int                 `json:"errorCount"`
	Project       []SkillAvailability `json:"project"`
}
type SkillAvailability struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Enabled bool   `json:"enabled"`
	Digest  string `json:"digest"`
}

func (p ContextProvenance) Valid() bool {
	if p.AssignedRole != "Master" || p.RoleBasis != "controller-launch" || !digestPattern.MatchString(p.InputRevision) || len(p.Sources) > 64 || len(p.Host.InstructionSources) > 64 || len(p.Skills.Project) > 64 || p.Skills.ExternalCount < 0 || p.Skills.ExternalCount > 4096 || p.Skills.ErrorCount < 0 || p.Skills.ErrorCount > 4096 || p.Skills.UseEvidence != "not-observed" {
		return false
	}
	status := func(s string) bool { return s == "observed" || s == "unavailable" }
	if !status(p.Host.Status) || !status(p.Host.InstructionStatus) || !status(p.Skills.Status) {
		return false
	}
	if p.Host.Status == "unavailable" && (p.Host.SessionID != "" || p.Host.ParentThreadID != "" || p.Host.NativeRole != "" || p.Host.Source != "" || p.Host.InstructionStatus != "unavailable") {
		return false
	}
	if p.Host.InstructionStatus == "unavailable" && len(p.Host.InstructionSources) > 0 {
		return false
	}
	if p.Skills.Status == "unavailable" && (len(p.Skills.Project) > 0 || p.Skills.ExternalCount != 0 || p.Skills.ErrorCount != 0) {
		return false
	}
	for _, s := range []string{p.Host.SessionID, p.Host.ParentThreadID, p.Host.NativeRole, p.Host.Source} {
		if len(s) > 1024 || strings.ContainsAny(s, "\x00\r\n") {
			return false
		}
	}
	for _, s := range p.Sources {
		if !relative(s.Path) || (!digestPattern.MatchString(s.SHA256) && s.SHA256 != "absent") {
			return false
		}
	}
	for _, path := range p.Host.InstructionSources {
		if !filepath.IsAbs(path) || len(path) > 4096 || strings.ContainsAny(path, "\x00\r\n") {
			return false
		}
	}
	seen := map[string]bool{}
	for _, skill := range p.Skills.Project {
		if !opaque(skill.Name) || !relative(skill.Path) || seen[skill.Path] || (!digestPattern.MatchString(skill.Digest) && skill.Digest != "unavailable") {
			return false
		}
		seen[skill.Path] = true
	}
	return true
}

type SkillInvocationView struct {
	AttemptID  string       `json:"attemptId"`
	InputID    string       `json:"inputId"`
	InputState string       `json:"inputState"`
	InputEpoch string       `json:"inputEpoch"`
	TurnID     string       `json:"turnId"`
	Acceptance string       `json:"acceptance"`
	Compliance string       `json:"compliance"`
	Skills     []SkillInput `json:"skills"`
}

type ProvenanceView struct {
	Invocation   *SkillInvocationView `json:"invocation,omitempty"`
	ThreadID     string               `json:"threadId"`
	HostVersion  string               `json:"hostVersion"`
	Schema       string               `json:"schema"`
	Operation    string               `json:"operation"`
	RunID        string               `json:"runId"`
	AttemptID    string               `json:"attemptId"`
	CaptureEpoch string               `json:"captureEpoch"`
	CurrentEpoch string               `json:"currentEpoch"`
	ReceivedAt   string               `json:"receivedAt"`
	Status       string               `json:"status"`
	Context      *ContextProvenance   `json:"context,omitempty"`
}

func (s *Service) Provenance(runID string) (ProvenanceView, error) {
	out := ProvenanceView{Schema: Schema, Operation: "provenance", RunID: runID, Status: "unavailable"}
	if !uuidPattern.MatchString(runID) {
		return out, fail("invalid_request", "invalid runId")
	}
	state, _, err := s.load()
	if err != nil {
		return out, err
	}
	r, ok := state.Runs[runID]
	if !ok {
		return out, fail("missing_coverage", "run not found")
	}
	return runProvenance(r)
}

// Project one already loaded run, keeping multi-run views on one generation.
func runProvenance(r Run) (ProvenanceView, error) {
	out := ProvenanceView{Schema: Schema, Operation: "provenance", RunID: r.ID, Status: "unavailable"}
	for i := len(r.Inputs) - 1; i >= 0; i-- {
		v := r.Inputs[i]
		if v.Kind != "launch" {
			continue
		}
		if len(v.Skills) > 0 {
			acceptance := "unconfirmed"
			if v.State == "acknowledged" {
				acceptance = "response-acknowledged"
			}
			for _, a := range r.Attempts {
				if a.ID == v.AttemptID && inputObserved(a, v, a.TurnID) {
					acceptance = "host-message-observed"
				}
			}
			out.Invocation = &SkillInvocationView{AttemptID: v.AttemptID, InputID: v.ID, InputState: v.State, InputEpoch: v.Epoch, TurnID: v.TurnID, Acceptance: acceptance, Compliance: "not-observed", Skills: v.Skills}
		}
		break
	}
	if len(r.Attempts) == 0 {
		return out, nil
	}
	a := r.Attempts[len(r.Attempts)-1]
	out.ThreadID = a.ThreadID
	out.HostVersion = a.HostVersion
	out.AttemptID = a.ID
	out.CurrentEpoch = a.Epoch
	for _, o := range a.Events {
		if o.Kind != "context.observed" {
			continue
		}
		var p ContextProvenance
		if err := Decode([]byte(o.Text), &p); err != nil || !p.Valid() {
			return out, fail("conflict", "invalid stored provenance")
		}
		out.Context = &p
		out.ReceivedAt = o.ReceivedAt
		out.CaptureEpoch = o.Epoch
		out.Status = "observed"
		if o.Epoch != a.Epoch {
			out.Status = "historical-epoch"
		}
	}
	return out, nil
}
