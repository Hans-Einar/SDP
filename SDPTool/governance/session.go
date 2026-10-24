package governance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
)

const SessionSchema = "sdp-session-execution/0.1"
const sessionStateSchema = "sdp-project-governance-state/0.13"

// The bounded definition is code-owned and pinned. Stages are not scripts.
func SessionWorkflow() Routine {
	return Routine{"session-lifecycle", "0.1", []string{"admission", "design", "implementation", "verification", "review", "model-release", "disposition", "done"}, "checked records, candidate evidence, required review and model release; no host-turn implied completion"}
}

type SessionRecordRef struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}
type SessionWorkflowRef struct {
	ID      string `json:"id"`
	Version string `json:"version"`
	Digest  string `json:"digest"`
}
type SessionCardRef struct {
	ID       string `json:"id"`
	Coverage string `json:"coverage"`
}
type SessionStepBinding struct {
	ID              string           `json:"id"`
	Milestone       string           `json:"milestone"`
	Cards           []SessionCardRef `json:"cards"`
	DependsOn       []string         `json:"dependsOn"`
	WritablePaths   []string         `json:"writablePaths"`
	ProtectedPaths  []string         `json:"protectedPaths"`
	Review          string           `json:"review"`
	OwnerAcceptance string           `json:"ownerAcceptance"`
}
type SessionBinding struct {
	Schema    string               `json:"schema"`
	Session   SessionRecordRef     `json:"session"`
	Plan      SessionRecordRef     `json:"plan"`
	Authority string               `json:"authority"`
	Workflow  SessionWorkflowRef   `json:"workflow"`
	Steps     []SessionStepBinding `json:"steps"`
}
type SessionContextRequest struct {
	BindingPath string `json:"bindingPath"`
}
type SessionContext struct {
	Schema        string         `json:"schema"`
	Operation     string         `json:"operation"`
	ProjectRef    string         `json:"projectRef"`
	CheckoutID    string         `json:"checkoutId"`
	BindingPath   string         `json:"bindingPath"`
	BindingDigest string         `json:"bindingDigest"`
	Binding       SessionBinding `json:"binding"`
	InputRevision string         `json:"inputRevision"`
	Sources       []Source       `json:"sources"`
}
type SessionAdoptRequest struct {
	Schema        string `json:"schema"`
	ProjectRef    string `json:"projectRef"`
	CheckoutID    string `json:"checkoutId"`
	OperationID   string `json:"operationId"`
	BindingPath   string `json:"bindingPath"`
	InputRevision string `json:"inputRevision"`
	Reason        string `json:"reason"`
}
type SessionResumeRequest struct {
	Schema           string `json:"schema"`
	ProjectRef       string `json:"projectRef"`
	CheckoutID       string `json:"checkoutId"`
	OperationID      string `json:"operationId"`
	InstanceID       string `json:"instanceId"`
	ExpectedRevision uint64 `json:"expectedRevision"`
	Reason           string `json:"reason"`
}
type SessionInstance struct {
	Disposition  *SessionDisposition `json:"disposition,omitempty"`
	Work         []SessionWork       `json:"work,omitempty"`
	Assessments  []SessionAssessment `json:"assessments,omitempty"`
	Selection    *SessionSelection   `json:"selection,omitempty"`
	ID           string              `json:"id"`
	Revision     uint64              `json:"revision"`
	Context      SessionContext      `json:"context"`
	Reason       string              `json:"reason"`
	Registration PublicationChange   `json:"registration"`
	Registered   bool                `json:"registered"`
}
type SessionMutationResult struct {
	RecordEffects string `json:"recordEffects,omitempty"`
	Schema        string `json:"schema"`
	Operation     string `json:"operation"`
	OperationID   string `json:"operationId"`
	InstanceID    string `json:"instanceId"`
	Revision      uint64 `json:"revision"`
	Registration  string `json:"registration"`
}
type SessionStepView struct {
	DispositionID        string                   `json:"dispositionId,omitempty"`
	CardDispositions     []SessionCardDisposition `json:"cardDispositions,omitempty"`
	VerificationEvidence string                   `json:"verificationEvidence,omitempty"`
	ReviewID             string                   `json:"reviewId,omitempty"`
	RunID                string                   `json:"runId,omitempty"`
	RunState             string                   `json:"runState,omitempty"`
	ProgressEvidence     string                   `json:"progressEvidence,omitempty"`
	Assessment           *SessionAssessmentView   `json:"assessment,omitempty"`
	Selected             bool                     `json:"selected,omitempty"`
	ID                   string                   `json:"id"`
	Stage                string                   `json:"stage"`
	Holds                []string                 `json:"holds"`
}
type SessionView struct {
	Binding       *SessionBinding    `json:"binding,omitempty"`
	RecordEffects string             `json:"recordEffects,omitempty"`
	Schema        string             `json:"schema"`
	Operation     string             `json:"operation"`
	InstanceID    string             `json:"instanceId"`
	Revision      uint64             `json:"revision"`
	SessionID     string             `json:"sessionId"`
	PlanID        string             `json:"planId"`
	Workflow      SessionWorkflowRef `json:"workflow"`
	Registration  string             `json:"registration"`
	Coverage      string             `json:"coverage"`
	Conditions    []string           `json:"conditions"`
	Steps         []SessionStepView  `json:"steps"`
	NextActions   []string           `json:"nextActions"`
	HostExecution string             `json:"hostExecution"`
}
type SessionList struct {
	Schema     string        `json:"schema"`
	Operation  string        `json:"operation"`
	Generation string        `json:"generation"`
	Instances  []SessionView `json:"instances"`
}

var sessionIDPattern = regexp.MustCompile(`^SESSION-[A-Z][A-Z0-9]*-[0-9]{4,}$`)
var sessionPlanPattern = regexp.MustCompile(`^(PLAN|MAINT)-[A-Z][A-Z0-9]*-[0-9]{4,}$`)
var sessionCardPattern = regexp.MustCompile(`^KB-[A-Z][A-Z0-9]*-[0-9]{3,}$`)
var sessionStepPattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]{0,63}$`)

func validateSessionBinding(b SessionBinding) error {
	w := SessionWorkflow()
	if b.Schema != SessionSchema || b.Workflow.ID != w.ID || b.Workflow.Version != w.Version || b.Workflow.Digest != w.Digest() {
		return fail("incompatible_profile", "unsupported Session binding or workflow definition")
	}
	if !sessionIDPattern.MatchString(b.Session.ID) || !sessionPlanPattern.MatchString(b.Plan.ID) || !strings.HasPrefix(b.Session.Path, "SDP/Sessions/") {
		return fail("invalid_request", "canonical Session and plan identity required")
	}
	for _, p := range []string{b.Session.Path, b.Plan.Path, b.Authority} {
		if !literalCandidatePath(p) || !strings.HasPrefix(p, "SDP/") || !strings.HasSuffix(p, ".md") {
			return fail("path", "Session references must be literal SDP Markdown paths")
		}
	}
	if b.Session.Path == b.Plan.Path || b.Authority == b.Session.Path || b.Authority == b.Plan.Path {
		return fail("invalid_request", "distinct Session, plan and authority documents required")
	}
	if len(b.Steps) == 0 || len(b.Steps) > 32 {
		return fail("limit", "one to 32 selected Session steps required")
	}
	steps := map[string]SessionStepBinding{}
	cards := map[string]bool{}
	for _, step := range b.Steps {
		if _, ok := steps[step.ID]; ok {
			return fail("conflict", "duplicate Session step")
		}
		if !sessionStepPattern.MatchString(step.ID) || !sessionStepPattern.MatchString(step.Milestone) || !slices.Contains([]string{"independent", "not-required"}, step.Review) || !slices.Contains([]string{"required", "not-required"}, step.OwnerAcceptance) {
			return fail("invalid_request", "step, milestone and explicit review/acceptance policy required")
		}
		if step.DependsOn == nil || len(step.DependsOn) > 32 || len(step.Cards) == 0 || len(step.Cards) > 16 || step.ProtectedPaths == nil {
			return fail("invalid_request", "bounded explicit step collections required")
		}
		if e := validateAssignment(Assignment{Outcome: step.ID, WritablePaths: step.WritablePaths, ExcludedPaths: step.ProtectedPaths}); e != nil {
			return e
		}
		seenPaths := map[string]bool{}
		for _, p := range append(append([]string{}, step.WritablePaths...), step.ProtectedPaths...) {
			if !literalCandidatePath(p) || seenPaths[p] {
				return fail("path", "unique literal scope paths required")
			}
			seenPaths[p] = true
		}
		seenCards := map[string]bool{}
		for _, c := range step.Cards {
			if !sessionCardPattern.MatchString(c.ID) || !slices.Contains([]string{"whole", "contribution"}, c.Coverage) || seenCards[c.ID] {
				return fail("invalid_request", "unique valid card coverage required")
			}
			seenCards[c.ID] = true
			cards[c.ID] = true
		}
		steps[step.ID] = step
	}
	if len(cards) > 64 {
		return fail("limit", "Session card count exceeds 64")
	}
	visiting, done := map[string]bool{}, map[string]bool{}
	var visit func(string) error
	visit = func(id string) error {
		if visiting[id] {
			return fail("conflict", "cyclic Session dependencies")
		}
		if done[id] {
			return nil
		}
		step, ok := steps[id]
		if !ok {
			return fail("missing_coverage", "unknown Session dependency")
		}
		visiting[id] = true
		seen := map[string]bool{}
		for _, dep := range step.DependsOn {
			if seen[dep] {
				return fail("conflict", "duplicate Session dependency")
			}
			seen[dep] = true
			if e := visit(dep); e != nil {
				return e
			}
		}
		visiting[id] = false
		done[id] = true
		return nil
	}
	for _, step := range b.Steps {
		if e := visit(step.ID); e != nil {
			return e
		}
	}
	return nil
}

// Match a roadmap/table cell or heading, not an incidental prose substring.
func sessionAnchor(body []byte, id string) bool {
	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSpace(line)
		candidates := []string{}
		if strings.HasPrefix(line, "|") {
			candidates = strings.Split(line, "|")
		}
		if strings.HasPrefix(line, "#") {
			candidates = append(candidates, strings.TrimLeft(line, "# "))
		}
		for _, c := range candidates {
			words := strings.Fields(c)
			if len(words) > 0 && strings.Trim(words[0], "`*") == id {
				return true
			}
		}
	}
	return false
}
func (s *Service) sessionContext(path string) (SessionContext, error) {
	out := SessionContext{Schema: SessionSchema, Operation: "session-context", ProjectRef: s.project, BindingPath: path, Sources: []Source{}}
	if !literalCandidatePath(path) || !strings.HasPrefix(path, "SDP/Sessions/execution/") || !strings.HasSuffix(path, ".json") {
		return out, fail("path", "Session binding must be under SDP/Sessions/execution")
	}
	raw, e := readFile(filepath.Join(s.root, path), 64<<10)
	if e != nil {
		return out, e
	}
	if e = decode(raw, &out.Binding); e != nil {
		return out, e
	}
	if e = validateSessionBinding(out.Binding); e != nil {
		return out, e
	}
	if path != "SDP/Sessions/execution/"+out.Binding.Session.ID+".json" {
		return out, fail("conflict", "binding filename must match Session identity")
	}
	files, e := s.publicationFiles()
	if e != nil {
		return out, e
	}
	b := out.Binding
	paths := map[string]bool{path: true, b.Session.Path: true, b.Plan.Path: true, b.Authority: true, "AGENTS.md": true, "SDP/SDP-project.manifest.yaml": true, "SDP/ProjectManagement/Ledger.ndjson": true}
	for _, p := range []string{b.Session.Path, b.Authority} {
		v, e := readFile(filepath.Join(s.root, p), MaxRequest)
		if e != nil {
			return out, e
		}
		files[p] = v
	}
	sm, e := metadata(files[b.Session.Path])
	if e != nil {
		return out, e
	}
	pm, e := metadata(files[b.Plan.Path])
	if e != nil {
		return out, e
	}
	if sm["Session reference"] != b.Session.ID || sm["Status"] != "active" || pm["id"] != b.Plan.ID || !slices.Contains([]string{"planned", "active"}, pm["state"]) {
		return out, fail("missing_coverage", "active Session and selected live plan must match binding")
	}
	if !strings.Contains(string(files[b.Session.Path]), filepath.Base(path)) {
		return out, fail("missing_coverage", "Session must link its execution binding")
	}
	if strings.TrimSpace(string(files[b.Authority])) == "" {
		return out, fail("missing_coverage", "authority document is empty")
	}
	cards := map[string]string{}
	for p, v := range files {
		if strings.HasPrefix(p, "SDP/KanBan/") {
			m, err := metadata(v)
			if err == nil && sessionCardPattern.MatchString(m["id"]) {
				cards[m["id"]] = p
			}
		}
	}
	for _, step := range b.Steps {
		if !sessionAnchor(files[b.Session.Path], step.ID) || !sessionAnchor(files[b.Plan.Path], step.Milestone) {
			return out, fail("missing_coverage", "selected Session step or plan milestone absent")
		}
		for _, c := range step.Cards {
			p, ok := cards[c.ID]
			if !ok {
				return out, fail("missing_coverage", "selected card absent")
			}
			m, _ := metadata(files[p])
			if m["PlanId"] != b.Plan.ID || !slices.Contains([]string{"backlog", "queued", "ready", "in-progress", "gate-review"}, m["CardState"]) {
				return out, fail("missing_coverage", "selected live card must reference plan")
			}
			paths[p] = true
		}
	}
	files[path] = raw
	names := []string{}
	for p := range paths {
		names = append(names, p)
	}
	sort.Strings(names)
	for _, p := range names {
		v, ok := files[p]
		if !ok {
			v, e = readFile(filepath.Join(s.root, p), MaxRequest)
			if os.IsNotExist(e) && (p == "AGENTS.md" || p == "SDP/SDP-project.manifest.yaml") {
				out.Sources = append(out.Sources, Source{p, "absent"})
				continue
			}
			if e != nil {
				return out, e
			}
		}
		out.Sources = append(out.Sources, Source{p, hash(v)})
	}
	out.BindingDigest = hash(raw)
	out.InputRevision = sessionInputRevision(out.Sources)
	return out, nil
}
func sessionInputRevision(sources []Source) string { b, _ := json.Marshal(sources); return hash(b) }
func (s *Service) SessionContext(q SessionContextRequest) (SessionContext, error) {
	st, _, e := s.load()
	if e != nil {
		return SessionContext{}, e
	}
	out, e := s.sessionContext(q.BindingPath)
	out.CheckoutID = st.CheckoutID
	return out, e
}
func sessionJournal(i SessionInstance) string {
	return fmt.Sprintf("\n## PG Session registration %s\n\nBound %s to workflow %s/%s at binding digest %s.\nThis records adoption only; no card activation, agent dispatch or completion.\n\nReason:\n\n%s", i.ID, i.Context.Binding.Session.ID, i.Context.Binding.Workflow.ID, i.Context.Binding.Workflow.Version, i.Context.BindingDigest, publicationProse(i.Reason))
}
func sessionRegistration(i SessionInstance) string {
	if i.Registered {
		return "committed"
	}
	return "pending"
}
func (s *Service) sessionCurrent(i SessionInstance) error {
	if i.Selection != nil {
		return s.selectionCoverage(i)
	}
	if _, e := s.publicationFiles(); e != nil {
		return e
	}
	for _, src := range i.Context.Sources {
		raw, e := readFile(filepath.Join(s.root, src.Path), maxState)
		actual := "absent"
		if e == nil {
			actual = hash(raw)
		} else if !os.IsNotExist(e) {
			return e
		}
		want := src.SHA256
		if src.Path == i.Registration.Path {
			if i.Registered {
				want = i.Registration.After
			} else if actual == i.Registration.After {
				continue
			}
		}
		if actual != want {
			return fail("stale_revision", "Session coverage changed: "+src.Path)
		}
	}
	return nil
}
func sessionRetry(st snapshot, id, digest string) (SessionMutationResult, bool, error) {
	r, ok, e := retry(st, id, digest)
	if e != nil {
		return SessionMutationResult{}, true, e
	}
	if !ok {
		return SessionMutationResult{}, false, nil
	}
	if r.Session == nil {
		return SessionMutationResult{}, true, fail("idempotency_conflict", "operation belongs to another workflow")
	}
	return *r.Session, true, nil
}
func (s *Service) AdoptSession(c Caller, q SessionAdoptRequest) (SessionMutationResult, error) {
	if !c.controller {
		return SessionMutationResult{}, fail("unauthorized", "trusted controller required")
	}
	if q.Schema != SessionSchema {
		return SessionMutationResult{}, fail("incompatible_profile", "unsupported Session operation")
	}
	if !digestPattern.MatchString(q.InputRevision) || !validLeaderText(q.Reason, 2048) {
		return SessionMutationResult{}, fail("invalid_request", "coverage revision and bounded reason required")
	}
	unlock, e := lock(s.state)
	if e != nil {
		return SessionMutationResult{}, e
	}
	defer unlock()
	st, head, e := s.load()
	if e != nil {
		return SessionMutationResult{}, e
	}
	if e = s.authorize(c, Schema, q.ProjectRef, q.CheckoutID, q.OperationID, st); e != nil {
		return SessionMutationResult{}, e
	}
	if r, ok, e := sessionRetry(st, q.OperationID, requestDigest(q)); ok {
		return r, e
	}
	ctx, e := s.sessionContext(q.BindingPath)
	if e != nil {
		return SessionMutationResult{}, e
	}
	ctx.CheckoutID = st.CheckoutID
	if ctx.InputRevision != q.InputRevision {
		return SessionMutationResult{}, fail("stale_revision", "Session adoption coverage changed")
	}
	for _, i := range st.Sessions {
		if i.Context.Binding.Session.ID == ctx.Binding.Session.ID || i.Context.Binding.Session.Path == ctx.Binding.Session.Path {
			return SessionMutationResult{}, fail("conflict", "Session already adopted; inspect/resume the existing instance")
		}
	}
	if len(st.Sessions) >= 32 || len(st.Receipts) >= 1024 {
		return SessionMutationResult{}, fail("limit", "Session or receipt capacity reached")
	}
	raw, e := readFile(filepath.Join(s.root, ctx.Binding.Session.Path), MaxRequest)
	if e != nil {
		return SessionMutationResult{}, e
	}
	i := SessionInstance{ID: q.OperationID, Revision: 1, Context: ctx, Reason: q.Reason}
	content := string(raw) + sessionJournal(i)
	if len(content) > MaxRequest {
		return SessionMutationResult{}, fail("limit", "Session journal exceeds limit")
	}
	i.Registration = PublicationChange{Path: ctx.Binding.Session.Path, Before: hash(raw), After: hash([]byte(content)), Content: content}
	if e = s.sessionCurrent(i); e != nil {
		return SessionMutationResult{}, e
	}
	if st.Sessions == nil {
		st.Sessions = map[string]SessionInstance{}
	}
	st.Sessions[i.ID] = i
	if !sessionEffectsProfile(st.Schema) {
		st.Schema = sessionStateSchema
	}
	out := SessionMutationResult{Schema: SessionSchema, Operation: "session-adopt", OperationID: q.OperationID, InstanceID: i.ID, Revision: i.Revision, Registration: "pending"}
	st.Receipts[q.OperationID] = receipt{requestDigest(q), Result{Schema: Schema, Operation: "session.adopt", OperationID: q.OperationID, ProjectRef: s.project, CheckoutID: st.CheckoutID, Session: &out}}
	if e = s.save(st, head); e != nil {
		return SessionMutationResult{}, e
	}
	return out, nil
}
func (s *Service) ResumeSession(c Caller, q SessionResumeRequest) (SessionMutationResult, error) {
	if !c.controller {
		return SessionMutationResult{}, fail("unauthorized", "trusted controller required")
	}
	if q.Schema != SessionSchema {
		return SessionMutationResult{}, fail("incompatible_profile", "unsupported Session operation")
	}
	if !uuidPattern.MatchString(q.InstanceID) || q.ExpectedRevision == 0 || !validLeaderText(q.Reason, 2048) {
		return SessionMutationResult{}, fail("invalid_request", "instance, revision and reason required")
	}
	dir, tx, unlock, e := s.lockRecordWriter()
	if e != nil {
		return SessionMutationResult{}, e
	}
	defer unlock()
	st, head, e := s.load()
	if e != nil {
		return SessionMutationResult{}, e
	}
	if e = s.authorize(c, Schema, q.ProjectRef, q.CheckoutID, q.OperationID, st); e != nil {
		return SessionMutationResult{}, e
	}
	if r, ok, e := sessionRetry(st, q.OperationID, requestDigest(q)); ok {
		if e == nil {
			if instance, ok := st.Sessions[q.InstanceID]; ok && instance.Selection != nil && sessionLatestEffects(instance).Status == "committed" {
				if err := tx.Release(s.sessionPendingOwner(instance)); err != nil {
					return SessionMutationResult{}, err
				}
			}
		}
		return r, e
	}
	i, ok := st.Sessions[q.InstanceID]
	if !ok {
		return SessionMutationResult{}, fail("missing_coverage", "Session instance absent")
	}
	if i.Revision != q.ExpectedRevision {
		return SessionMutationResult{}, &Error{Code: "stale_revision", Message: "Session revision changed", CurrentRevision: i.Revision}
	}
	if i.Revision == ^uint64(0) || len(st.Receipts) >= 1024 {
		return SessionMutationResult{}, fail("limit", "Session revision or receipts exhausted")
	}
	if _, e = readFile(filepath.Join(dir, "pending"), 8192); !os.IsNotExist(e) {
		if e != nil {
			return SessionMutationResult{}, e
		}
		return SessionMutationResult{}, fail("reconciliation_required", "another canonical publication is pending")
	}
	owner := ""
	if i.Selection != nil {
		owner = s.sessionPendingOwner(i)
	}
	if e = tx.Check(owner); e != nil {
		return SessionMutationResult{}, e
	}
	encoded, _ := json.Marshal(st)
	if len(encoded) > maxState-8192 {
		return SessionMutationResult{}, fail("limit", "insufficient receipt capacity before effects")
	}
	if i.Disposition != nil && i.Disposition.Effects.Status == "prepared" {
		if e = s.applySessionDisposition(&st, &i, tx); e != nil {
			return SessionMutationResult{}, e
		}
	} else if len(i.Work) > 0 && sessionLatestEffects(i).Status == "prepared" {
		if e = s.applySessionWork(&st, &i, tx); e != nil {
			return SessionMutationResult{}, e
		}
	} else if len(i.Assessments) > 0 && sessionLatestEffects(i).Status == "prepared" {
		a := &i.Assessments[len(i.Assessments)-1]
		if e = s.checkSelectionClaims(i, i.Selection.Request, selectedSessionStep(i)); e != nil {
			return SessionMutationResult{}, e
		}
		if e = s.applySessionRecordEffects(&a.Effects, tx, s.sessionPendingOwner(i), "session-assessment"); e != nil {
			return SessionMutationResult{}, e
		}
	} else if i.Selection != nil && i.Selection.Status == "prepared" {
		if e = s.applySessionSelection(&i, tx); e != nil {
			return SessionMutationResult{}, e
		}
	} else if e = s.sessionCurrent(i); e != nil {
		return SessionMutationResult{}, e
	}
	if !i.Registered {
		if e = s.registrationUnclaimed(i); e != nil {
			return SessionMutationResult{}, e
		}
		if e = s.applyPublicationFile(i.Registration); e != nil {
			return SessionMutationResult{}, e
		}
		if s.fault != nil {
			if e = s.fault("session-after-journal"); e != nil {
				return SessionMutationResult{}, e
			}
		}
		i.Registered = true
	}
	if e = s.sessionCurrent(i); e != nil {
		return SessionMutationResult{}, e
	}
	if i.Selection != nil {
		ledger, e := readFile(filepath.Join(s.root, managementLedger), maxState)
		if e != nil {
			return SessionMutationResult{}, e
		}
		for n, src := range sessionLatestEffects(i).Coverage {
			if src.Path == managementLedger {
				sessionLatestEffects(i).Coverage[n].SHA256 = hash(ledger)
			}
		}
		sessionLatestEffects(i).LedgerBytes = len(ledger)
	}
	i.Revision++
	st.Sessions[i.ID] = i
	out := SessionMutationResult{Schema: SessionSchema, Operation: "session-resume", OperationID: q.OperationID, InstanceID: i.ID, Revision: i.Revision, Registration: "committed"}
	if i.Selection != nil {
		out.RecordEffects = sessionLatestEffects(i).Status
	}
	st.Receipts[q.OperationID] = receipt{requestDigest(q), Result{Schema: Schema, Operation: "session.resume", OperationID: q.OperationID, ProjectRef: s.project, CheckoutID: st.CheckoutID, Session: &out}}
	if e = s.save(st, head); e != nil {
		return SessionMutationResult{}, e
	}
	if i.Selection != nil {
		if e = tx.Release(s.sessionPendingOwner(i)); e != nil {
			return SessionMutationResult{}, e
		}
	}
	return out, nil
}
func (s *Service) sessionView(i SessionInstance, st snapshot) SessionView {
	v := SessionView{Schema: SessionSchema, Operation: "session-read", InstanceID: i.ID, Revision: i.Revision, SessionID: i.Context.Binding.Session.ID, PlanID: i.Context.Binding.Plan.ID, Workflow: i.Context.Binding.Workflow, Registration: sessionRegistration(i), Coverage: "current", Conditions: []string{}, Steps: []SessionStepView{}, NextActions: []string{"session-read", "session-resume"}, HostExecution: "not-dispatched"}
	if e := s.sessionCurrent(i); e != nil {
		v.Coverage = "stale-or-unavailable"
		v.Conditions = append(v.Conditions, e.Error())
		v.NextActions = []string{"session-read"}
	}
	if i.Selection != nil {
		v.RecordEffects = sessionLatestEffects(i).Status
		if v.RecordEffects == "prepared" {
			v.NextActions = []string{"session-read", "session-resume"}
		}
	} else if i.Registered && v.Coverage == "current" {
		v.NextActions = append(v.NextActions, "session-select")
	}
	for _, step := range i.Context.Binding.Steps {
		holds := []string{"step activation is not implemented in this increment"}
		if !i.Registered {
			holds = append(holds, "Session registration pending")
		}
		if len(step.DependsOn) > 0 {
			holds = append(holds, "predecessor steps not completed")
		}
		holds = append(holds, v.Conditions...)
		selected := i.Selection != nil && i.Selection.Status == "committed" && i.Selection.Request.StepID == step.ID
		if selected {
			holds[0] = "design assessment is required before execution"
		}
		sv := SessionStepView{ID: step.ID, Stage: "admission", Holds: holds, Selected: selected}
		if selected {
			if v.Coverage == "current" && len(i.Assessments) < maxSessionAssessments && len(i.Work) == 0 {
				v.NextActions = append(v.NextActions, "session-assess")
			}
			for n := len(i.Assessments) - 1; n >= 0; n-- {
				a := i.Assessments[n]
				if a.Effects.Status != "committed" {
					continue
				}
				sv.Assessment = assessmentView(a)
				sv.Assessment.DesignImpact = assessmentImpact(a.Request, step)
				sv.Stage, sv.Holds = assessmentStage(a.Request, step)
				sv.Holds = append(sv.Holds, v.Conditions...)
				break
			}
		}
		s.projectSessionWork(i, st, &v, &sv)
		s.projectSessionDisposition(i, st, &v, &sv)
		v.Steps = append(v.Steps, sv)
	}
	return v
}
func (s *Service) ReadSession(id string) (SessionView, error) {
	if !uuidPattern.MatchString(id) {
		return SessionView{}, fail("invalid_request", "Session instance UUID required")
	}
	st, _, e := s.load()
	if e != nil {
		return SessionView{}, e
	}
	i, ok := st.Sessions[id]
	if !ok {
		return SessionView{}, fail("missing_coverage", "Session instance absent")
	}
	return s.sessionView(i, st), nil
}
func (s *Service) ListSessions() (SessionList, error) {
	st, head, e := s.load()
	out := SessionList{SessionSchema, "session-list", head, []SessionView{}}
	if e != nil {
		return out, e
	}
	ids := []string{}
	for id := range st.Sessions {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		out.Instances = append(out.Instances, s.sessionView(st.Sessions[id], st))
	}
	return out, nil
}
func validateSessionState(st snapshot) error {
	if !sessionProfile(st.Schema) {
		if st.Sessions != nil {
			return fail("incompatible_profile", "Sessions require profile 0.13")
		}
		for _, r := range st.Receipts {
			if r.Result.Session != nil {
				return fail("incompatible_profile", "Session receipt outside profile")
			}
		}
		return nil
	}
	if len(st.Sessions) == 0 || len(st.Sessions) > 32 {
		return fail("reconciliation_required", "invalid Session collection")
	}
	seen := map[string]bool{}
	paths := map[string]bool{}
	revisions := map[string]map[uint64]bool{}
	for id, i := range st.Sessions {
		c := i.Context
		if !uuidPattern.MatchString(id) || id != i.ID || i.Revision == 0 || validateSessionBinding(c.Binding) != nil || seen[c.Binding.Session.ID] || paths[c.Binding.Session.Path] || c.Operation != "session-context" || c.Schema != SessionSchema || c.ProjectRef != st.Project || c.CheckoutID != st.CheckoutID || c.BindingPath != "SDP/Sessions/execution/"+c.Binding.Session.ID+".json" || !digestPattern.MatchString(c.BindingDigest) || !validLeaderText(i.Reason, 2048) || len(c.Sources) < 5 || len(c.Sources) > 72 || c.InputRevision != sessionInputRevision(c.Sources) {
			return fail("reconciliation_required", "invalid stored Session identity/binding")
		}
		seen[c.Binding.Session.ID] = true
		paths[c.Binding.Session.Path] = true
		revisions[id] = map[uint64]bool{}
		reg := i.Registration
		prefix, ok := strings.CutSuffix(reg.Content, sessionJournal(i))
		if !ok || reg.Path != c.Binding.Session.Path || hash([]byte(prefix)) != reg.Before || hash([]byte(reg.Content)) != reg.After || len(reg.Content) > MaxRequest || (i.Registered && i.Revision < 2) {
			return fail("reconciliation_required", "invalid Session registration intent")
		}
		last := ""
		pins := map[string]string{}
		for _, p := range c.Sources {
			if !literalCandidatePath(p.Path) || p.Path <= last || (!digestPattern.MatchString(p.SHA256) && p.SHA256 != "absent") {
				return fail("reconciliation_required", "invalid Session coverage pins")
			}
			last = p.Path
			pins[p.Path] = p.SHA256
		}
		if pins[c.Binding.Session.Path] != reg.Before || pins[c.BindingPath] != c.BindingDigest || !digestPattern.MatchString(pins[c.Binding.Plan.Path]) || !digestPattern.MatchString(pins[c.Binding.Authority]) || !digestPattern.MatchString(pins["SDP/ProjectManagement/Ledger.ndjson"]) {
			return fail("reconciliation_required", "missing Session coverage pins")
		}
		r, ok := st.Receipts[id]
		if !ok || r.Result.Session == nil || r.Result.Session.Operation != "session-adopt" || r.Result.Session.InstanceID != id || r.Result.Session.Revision != 1 {
			return fail("reconciliation_required", "missing Session adoption receipt")
		}
	}
	for id, r := range st.Receipts {
		v := r.Result.Session
		if v == nil {
			if strings.HasPrefix(r.Result.Operation, "session.") {
				return fail("reconciliation_required", "missing Session receipt result")
			}
			continue
		}
		i, ok := st.Sessions[v.InstanceID]
		if !ok || !uuidPattern.MatchString(id) || id != v.OperationID || v.Schema != SessionSchema || v.Revision == 0 || v.Revision > i.Revision || !digestPattern.MatchString(r.Digest) || r.Result.ProjectRef != st.Project || r.Result.CheckoutID != st.CheckoutID || r.Result.OperationID != id || r.Result.Schema != Schema || !slices.Contains([]string{"session-adopt", "session-resume", "session-select", "session-assess", "session-start", "session-progress", "session-verify", "session-review", "session-dispose"}, v.Operation) || r.Result.Operation != strings.Replace(v.Operation, "-", ".", 1) || (v.Operation == "session-adopt" && (v.InstanceID != id || v.Registration != "pending")) || (v.Operation == "session-resume" && (v.Revision < 2 || v.Registration != "committed")) {
			return fail("reconciliation_required", "invalid Session receipt")
		}
		if revisions[v.InstanceID][v.Revision] {
			return fail("reconciliation_required", "duplicate Session revision receipt")
		}
		revisions[v.InstanceID][v.Revision] = true
	}
	for id, i := range st.Sessions {
		if uint64(len(revisions[id])) != i.Revision || (!i.Registered && i.Revision != 1) {
			return fail("reconciliation_required", "incomplete Session revision history")
		}
	}
	return nil
}
