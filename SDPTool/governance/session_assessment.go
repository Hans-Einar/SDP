package governance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"
)

const sessionAssessmentStateSchema = "sdp-project-governance-state/0.15"
const maxSessionAssessments = 8

// An attributed decision with explicit criteria; never an automated proof that
// prose, source code, or the host's actual skill consumption has been verified.
type SessionDesignCriterion struct {
	Domain    string `json:"domain"`
	Impact    string `json:"impact"`
	Rationale string `json:"rationale"`
}
type SessionBlueprintDecision struct {
	Requirement string `json:"requirement"`
	Rationale   string `json:"rationale"`
}
type SessionAssessRequest struct {
	Schema           string                   `json:"schema"`
	ProjectRef       string                   `json:"projectRef"`
	CheckoutID       string                   `json:"checkoutId"`
	OperationID      string                   `json:"operationId"`
	InstanceID       string                   `json:"instanceId"`
	ExpectedRevision uint64                   `json:"expectedRevision"`
	StepID           string                   `json:"stepId"`
	Reason           string                   `json:"reason"`
	Criteria         []SessionDesignCriterion `json:"criteria"`
	Blueprint        SessionBlueprintDecision `json:"blueprint"`
	Evidence         []Source                 `json:"evidence"`
	Skills           []SkillInput             `json:"skills"`
}
type SessionAssessment struct {
	Request SessionAssessRequest `json:"request"`
	Effects SessionRecordEffects `json:"effects"`
}
type SessionAssessmentView struct {
	OperationID  string                   `json:"operationId"`
	RecordPath   string                   `json:"recordPath"`
	DesignImpact string                   `json:"designImpact"`
	Blueprint    SessionBlueprintDecision `json:"blueprint"`
	Reason       string                   `json:"reason"`
	Criteria     []SessionDesignCriterion `json:"criteria"`
	Evidence     []Source                 `json:"evidence"`
	Skills       []SkillInput             `json:"skills"`
}

var sessionDesignDomains = []string{"contracts", "data-model", "responsibilities", "state-flow", "interfaces", "persistence", "deployment"}

func sessionEffectsProfile(schema string) bool {
	return schema == sessionEffectsStateSchema || sessionAssessmentProfile(schema)
}
func sessionLatestEffects(i SessionInstance) *SessionRecordEffects {
	if i.Advance != nil {
		return &i.Advance.Effects
	}
	if i.Disposition != nil {
		return &i.Disposition.Effects
	}
	if g := latestOwnerGate(i); g != nil {
		return &g.Effects
	}
	if len(i.Work) > 0 {
		return &i.Work[len(i.Work)-1].Effects
	}
	if len(i.Assessments) > 0 {
		return &i.Assessments[len(i.Assessments)-1].Effects
	}
	if i.Selection != nil {
		return &i.Selection.SessionRecordEffects
	}
	if len(i.Past) > 0 {
		return &i.Past[len(i.Past)-1].Transition.Effects
	}
	return nil
}
func (s *Service) sessionPendingOwner(i SessionInstance) string {
	if i.Advance != nil {
		return s.state + "\n" + i.Advance.Request.OperationID + "\n"
	}
	if i.Selection == nil && len(i.Past) > 0 {
		return s.state + "\n" + i.Past[len(i.Past)-1].Transition.Request.OperationID + "\n"
	}
	if i.Disposition != nil {
		return s.state + "\n" + i.Disposition.Request.OperationID + "\n"
	}
	if g := latestOwnerGate(i); g != nil {
		return s.state + "\n" + g.Request.OperationID + "\n"
	}
	if len(i.Work) > 0 {
		return s.state + "\n" + i.Work[len(i.Work)-1].Request.OperationID + "\n"
	}
	if len(i.Assessments) > 0 {
		return s.state + "\n" + i.Assessments[len(i.Assessments)-1].Request.OperationID + "\n"
	}
	return s.sessionEffectOwner(i)
}
func assessmentPath(q SessionAssessRequest) string {
	return "SDP/Sessions/execution/assessments/" + q.OperationID + ".json"
}
func selectedSessionStep(i SessionInstance) SessionStepBinding {
	if i.Selection != nil {
		for _, step := range i.Context.Binding.Steps {
			if step.ID == i.Selection.Request.StepID {
				return step
			}
		}
	}
	return SessionStepBinding{}
}
func assessmentImpact(q SessionAssessRequest, step SessionStepBinding) string {
	impact := "unchanged"
	for _, c := range q.Criteria {
		if c.Impact == "changed" {
			return "changed"
		}
		if c.Impact == "unknown" {
			impact = "unknown"
		}
	}
	// A negative prose decision cannot override explicitly writable model sources.
	for _, p := range step.WritablePaths {
		if strings.HasSuffix(strings.ToLower(p), ".design") || strings.HasSuffix(strings.ToLower(p), ".sdui") || p == "SDP/SDL" || strings.HasPrefix(p, "SDP/SDL/") {
			return "changed"
		}
	}
	return impact
}
func assessmentStage(q SessionAssessRequest, step SessionStepBinding) (string, []string) {
	switch assessmentImpact(q, step) {
	case "changed":
		return "design", []string{"SDL WORK and accepted design handoff required; design dispatch is not implemented"}
	case "unknown":
		return "admission", []string{"design impact is unknown; record a resolved assessment before execution"}
	}
	if q.Blueprint.Requirement == "unknown" {
		return "admission", []string{"Blueprint applicability is unknown; record a resolved assessment before execution"}
	}
	if q.Blueprint.Requirement == "required" {
		return "implementation", []string{"required Blueprint handoff is not yet bound", "Session start/progress integration is not implemented"}
	}
	return "implementation", []string{"explicit no-design disposition recorded; Session start/progress integration is not implemented"}
}
func assessmentView(a SessionAssessment) *SessionAssessmentView {
	q := a.Request
	return &SessionAssessmentView{OperationID: q.OperationID, RecordPath: assessmentPath(q), Blueprint: q.Blueprint, Reason: q.Reason, Criteria: q.Criteria, Evidence: q.Evidence, Skills: q.Skills}
}
func validateSessionAssessRequest(q SessionAssessRequest) error {
	if q.Schema != SessionSchema || !uuidPattern.MatchString(q.InstanceID) || !uuidPattern.MatchString(q.OperationID) || q.ExpectedRevision == 0 || !sessionStepPattern.MatchString(q.StepID) || !validLeaderText(q.Reason, 2048) || len(q.Criteria) != len(sessionDesignDomains) || len(q.Evidence) < 1 || len(q.Evidence) > 8 || !ValidSkillInputs(q.Skills) {
		return fail("invalid_request", "complete bounded Session design assessment required")
	}
	seen := map[string]bool{}
	for _, c := range q.Criteria {
		if !slices.Contains(sessionDesignDomains, c.Domain) || seen[c.Domain] || !slices.Contains([]string{"changed", "unchanged", "unknown"}, c.Impact) || !validLeaderText(c.Rationale, 1024) {
			return fail("invalid_request", "each design domain needs one impact and rationale")
		}
		seen[c.Domain] = true
	}
	if !slices.Contains([]string{"required", "not-required", "unknown"}, q.Blueprint.Requirement) || !validLeaderText(q.Blueprint.Rationale, 1024) {
		return fail("invalid_request", "separate Blueprint applicability and rationale required")
	}
	seen = map[string]bool{}
	for _, e := range q.Evidence {
		if !literalCandidatePath(e.Path) || !digestPattern.MatchString(e.SHA256) || seen[e.Path] {
			return fail("invalid_request", "distinct project-local evidence pins required")
		}
		seen[e.Path] = true
	}
	leader, analysis := false, false
	for _, skill := range q.Skills {
		if filepath.Base(skill.Path) != "SKILL.md" || filepath.Base(filepath.Dir(skill.Path)) != skill.Name {
			return fail("invalid_request", "skill name must identify its selected SKILL.md")
		}
		leader = leader || skill.Name == "sdp-project-leader"
		analysis = analysis || skill.Name == "sdp-change-analysis" || skill.Name == "sdp-architect"
	}
	if !leader || !analysis {
		return fail("missing_coverage", "Project Leader and Change Analysis or Architect skill snapshots required")
	}
	return nil
}
func assessmentArtifact(a SessionAssessment, i SessionInstance) []byte {
	b, _ := json.MarshalIndent(struct {
		Schema        string               `json:"schema"`
		Attribution   string               `json:"attribution"`
		CreatedAt     string               `json:"createdAt"`
		BindingDigest string               `json:"bindingDigest"`
		DesignImpact  string               `json:"designImpact"`
		Request       SessionAssessRequest `json:"request"`
	}{"sdp-session-assessment/0.1", "trusted-controller", a.Effects.CreatedAt, i.Context.BindingDigest, assessmentImpact(a.Request, selectedSessionStep(i)), a.Request}, "", "  ")
	return append(b, '\n')
}
func assessmentJournal(a SessionAssessment, i SessionInstance, path string) string {
	q := a.Request
	rel, _ := filepath.Rel(filepath.Dir(path), assessmentPath(q))
	return "\n## PG design assessment " + q.OperationID + "\n\nStep " + q.StepID + "; design impact: " + assessmentImpact(q, selectedSessionStep(i)) + "; Blueprint: " + q.Blueprint.Requirement + ".\n[Attributed assessment and pinned inputs](" + filepath.ToSlash(rel) + ").\nNo host dispatch, implementation result, independent approval or completion.\n\nReason:\n\n" + publicationProse(q.Reason)
}
func (s *Service) buildSessionAssessment(i SessionInstance, q SessionAssessRequest) (SessionAssessment, error) {
	a := SessionAssessment{Request: q, Effects: SessionRecordEffects{Status: "prepared", CreatedAt: time.Now().UTC().Format(time.RFC3339Nano)}}
	files, e := s.publicationFiles()
	if e != nil {
		return a, e
	}
	pins := append([]Source{}, sessionLatestEffects(i).Coverage...)
	for _, pin := range q.Evidence {
		raw, err := readFile(filepath.Join(s.root, pin.Path), maxState)
		if err != nil || hash(raw) != pin.SHA256 {
			return a, fail("stale_revision", "assessment evidence changed: "+pin.Path)
		}
	}
	pins = append(pins, q.Evidence...)
	if e = CheckSkillSnapshots(s.root, q.Skills); e != nil {
		return a, e
	}
	for _, skill := range q.Skills {
		p, _ := filepath.Rel(s.root, skill.Path)
		pins = append(pins, Source{filepath.ToSlash(p), skill.Digest})
	}
	for _, pin := range pins {
		raw, err := readFile(filepath.Join(s.root, pin.Path), maxState)
		if os.IsNotExist(err) && pin.SHA256 == "absent" {
			continue
		}
		if err != nil {
			return a, err
		}
		// History may have acquired an unrelated validated suffix since last observation.
		if pin.Path != managementLedger && hash(raw) != pin.SHA256 {
			return a, fail("stale_revision", "assessment input changed: "+pin.Path)
		}
		files[pin.Path] = raw
	}
	before := map[string][]byte{}
	for p, b := range files {
		before[p] = b
	}
	path := assessmentPath(q)
	if _, err := os.Lstat(filepath.Join(s.root, path)); !os.IsNotExist(err) {
		return a, fail("conflict", "assessment record destination exists or unavailable")
	}
	files[path] = assessmentArtifact(a, i)
	subjects := map[string]string{i.Context.Binding.Plan.ID: i.Context.Binding.Plan.Path}
	for _, card := range selectedSessionStep(i).Cards {
		for p, b := range files {
			if strings.HasPrefix(p, "SDP/KanBan/") {
				m, err := metadata(b)
				if err == nil && m["id"] == card.ID {
					subjects[card.ID] = p
				}
			}
		}
		if subjects[card.ID] == "" {
			return a, fail("missing_coverage", "assessment card absent")
		}
	}
	events, e := publicationEvents(files[managementLedger])
	if e != nil {
		return a, e
	}
	latest := map[string]map[string]any{}
	for _, ev := range events {
		latest[fmt.Sprint(ev["subjectId"])] = ev
	}
	ids := []string{}
	for id := range subjects {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		p := subjects[id]
		files[p] = []byte(appendRecord(files[p], assessmentJournal(a, i, p)))
		old := latest[id]
		if old == nil {
			return a, fail("missing_coverage", "assessment history absent")
		}
		payload := map[string]any{}
		for k, v := range old["payload"].(map[string]any) {
			payload[k] = v
		}
		delete(payload, "lineage")
		payload["previousEventId"] = old["eventId"]
		payload["from"] = payload["to"]
		payload["fromPath"] = payload["toPath"]
		payload["reason"] = "Recorded design assessment " + q.OperationID + "; records remain ready/in-progress; no dispatch or completion."
		prefix, typ := "EVT-PM-"+strings.Split(id, "-")[1]+"-", "x-management:updated"
		if strings.HasPrefix(id, "KB-") {
			prefix = "EVT-KB-" + strings.Split(id, "-")[1] + "-"
			typ = "x-kanban:reviewed"
		}
		eventID, err := nextPublicationEvent(events, prefix)
		if err != nil {
			return a, err
		}
		ev := map[string]any{"schemaVersion": "1.0", "eventId": eventID, "eventType": typ, "occurredAt": a.Effects.CreatedAt, "actor": "sdp-controller", "commit": nil, "subjectId": id, "payload": payload}
		b, _ := json.Marshal(ev)
		files[managementLedger] = []byte(appendRecord(files[managementLedger], string(b)+"\n"))
		events = append(events, ev)
	}
	sp := i.Context.Binding.Session.Path
	files[sp] = []byte(appendRecord(files[sp], assessmentJournal(a, i, sp)))
	if e = validatePublicationHistory(files[managementLedger], files); e != nil {
		return a, e
	}
	names := []string{}
	for p := range files {
		names = append(names, p)
	}
	for _, pin := range pins {
		if _, ok := files[pin.Path]; !ok {
			names = append(names, pin.Path)
		}
	}
	sort.Strings(names)
	names = slices.Compact(names)
	for _, p := range names {
		pre, post := "absent", "absent"
		if b, ok := before[p]; ok {
			pre = hash(b)
		}
		if b, ok := files[p]; ok {
			post = hash(b)
		}
		a.Effects.Inputs = append(a.Effects.Inputs, Source{p, pre})
		if pre != post {
			a.Effects.Changes = append(a.Effects.Changes, PublicationChange{p, pre, post, string(files[p])})
		}
	}
	covered := map[string]bool{path: true}
	for _, pin := range pins {
		covered[pin.Path] = true
	}
	for _, p := range names {
		if covered[p] {
			h := "absent"
			if b, ok := files[p]; ok {
				h = hash(b)
			}
			a.Effects.Coverage = append(a.Effects.Coverage, Source{p, h})
		}
	}
	a.Effects.LedgerBytes = len(files[managementLedger])
	if len(a.Effects.Inputs) > 2048 {
		return a, fail("limit", "assessment input capacity exceeded")
	}
	size := 0
	for _, ch := range a.Effects.Changes {
		size += len(ch.Content)
	}
	if size > 4<<20 {
		return a, fail("limit", "assessment effect capacity exceeded")
	}
	return a, nil
}
func (s *Service) AssessSession(c Caller, q SessionAssessRequest) (SessionMutationResult, error) {
	if !c.controller {
		return SessionMutationResult{}, fail("unauthorized", "trusted controller required")
	}
	if e := validateSessionAssessRequest(q); e != nil {
		return SessionMutationResult{}, e
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
		return r, e
	}
	if e = tx.Check(""); e != nil {
		return SessionMutationResult{}, e
	}
	if _, e = readFile(filepath.Join(dir, "pending"), 8192); !os.IsNotExist(e) {
		return SessionMutationResult{}, fail("reconciliation_required", "canonical publication pending")
	}
	i, ok := st.Sessions[q.InstanceID]
	if !ok {
		return SessionMutationResult{}, fail("missing_coverage", "Session instance absent")
	}
	if i.Revision != q.ExpectedRevision {
		return SessionMutationResult{}, fail("stale_revision", "Session revision changed")
	}
	if len(i.Work) > 0 || !i.Registered || i.Selection == nil || i.Selection.Status != "committed" || selectedSessionStep(i).ID != q.StepID {
		return SessionMutationResult{}, fail("conflict", "committed selected step required")
	}
	if len(i.Assessments) >= maxSessionAssessments || i.Revision == ^uint64(0) || len(st.Receipts) >= 1024 {
		return SessionMutationResult{}, fail("limit", "Session assessment capacity reached")
	}
	if e = s.sessionCurrent(i); e != nil {
		return SessionMutationResult{}, e
	}
	if e = s.checkSelectionClaims(i, i.Selection.Request, selectedSessionStep(i)); e != nil {
		return SessionMutationResult{}, e
	}
	a, e := s.buildSessionAssessment(i, q)
	if e != nil {
		return SessionMutationResult{}, e
	}
	i.Assessments = append(i.Assessments, a)
	i.Revision++
	st.Sessions[i.ID] = i
	if !sessionWorkProfile(st.Schema) {
		st.Schema = sessionAssessmentStateSchema
	}
	out := SessionMutationResult{Schema: SessionSchema, Operation: "session-assess", OperationID: q.OperationID, InstanceID: i.ID, Revision: i.Revision, Registration: "committed", RecordEffects: "prepared"}
	st.Receipts[q.OperationID] = receipt{requestDigest(q), Result{Schema: Schema, Operation: "session.assess", OperationID: q.OperationID, ProjectRef: s.project, CheckoutID: st.CheckoutID, Session: &out}}
	if e = s.save(st, head); e != nil {
		return SessionMutationResult{}, e
	}
	return out, nil
}

func validateSessionAssessments(st snapshot) error {
	seen := map[string]bool{}
	for _, i := range sessionExecutions(st) {
		if len(i.Assessments) == 0 {
			continue
		}
		if !sessionAssessmentProfile(st.Schema) || !i.Registered || i.Selection == nil || i.Selection.Status != "committed" || len(i.Assessments) > maxSessionAssessments {
			return fail("reconciliation_required", "invalid assessment owner/profile")
		}
		previous := i.Selection.SessionRecordEffects
		previousRevision := i.Selection.Request.ExpectedRevision + 1
		for n, a := range i.Assessments {
			q, p := a.Request, a.Effects
			if validateSessionAssessRequest(q) != nil || q.ProjectRef != st.Project || q.CheckoutID != st.CheckoutID || q.InstanceID != i.ID || q.StepID != selectedSessionStep(i).ID || q.ExpectedRevision <= previousRevision || q.ExpectedRevision >= i.Revision || seen[q.OperationID] {
				return fail("reconciliation_required", "invalid assessment identity or request")
			}
			seen[q.OperationID] = true
			previousRevision = q.ExpectedRevision + 1
			r, ok := st.Receipts[q.OperationID]
			if !ok || r.Digest != requestDigest(q) || r.Result.Session == nil || r.Result.Session.Operation != "session-assess" || r.Result.Session.RecordEffects != "prepared" || r.Result.Session.Revision != q.ExpectedRevision+1 {
				return fail("reconciliation_required", "assessment receipt mismatch")
			}
			if !slices.Contains([]string{"prepared", "committed"}, p.Status) || (p.Status == "prepared" && (n != len(i.Assessments)-1 || i.Revision != q.ExpectedRevision+1)) || (p.Status == "committed" && i.Revision <= q.ExpectedRevision+1) || len(p.Inputs) > 2048 || len(p.Changes) < 4 || len(p.Changes) > 20 || len(p.Coverage) > 218 || p.LedgerBytes <= 0 {
				return fail("reconciliation_required", "invalid assessment effect state")
			}
			if _, e := time.Parse(time.RFC3339Nano, p.CreatedAt); e != nil {
				return fail("reconciliation_required", "invalid assessment time")
			}
			inputs := map[string]string{}
			for _, pin := range p.Inputs {
				if !literalCandidatePath(pin.Path) || inputs[pin.Path] != "" || (!digestPattern.MatchString(pin.SHA256) && pin.SHA256 != "absent") {
					return fail("reconciliation_required", "invalid assessment preimage")
				}
				inputs[pin.Path] = pin.SHA256
			}
			required := map[string]string{}
			for _, pin := range previous.Coverage {
				required[pin.Path] = pin.SHA256
			}
			for _, pin := range q.Evidence {
				if h, ok := required[pin.Path]; ok && pin.Path != managementLedger && h != pin.SHA256 {
					return fail("reconciliation_required", "assessment evidence contradicts existing pin")
				}
				required[pin.Path] = pin.SHA256
			}
			for _, skill := range q.Skills {
				rel, e := filepath.Rel(st.Root, skill.Path)
				rel = filepath.ToSlash(rel)
				if e != nil || !literalCandidatePath(rel) {
					return fail("reconciliation_required", "assessment skill outside project")
				}
				if h, ok := required[rel]; ok && h != skill.Digest {
					return fail("reconciliation_required", "assessment skill contradicts existing pin")
				}
				required[rel] = skill.Digest
			}
			for path, digest := range required {
				if path != managementLedger && inputs[path] != digest {
					return fail("reconciliation_required", "assessment prior coverage differs")
				}
			}
			changed := map[string]bool{}
			cards := map[string]bool{}
			for _, card := range selectedSessionStep(i).Cards {
				cards[card.ID] = true
			}
			size := 0
			for _, ch := range p.Changes {
				size += len(ch.Content)
				if changed[ch.Path] || inputs[ch.Path] != ch.Before || !digestPattern.MatchString(ch.After) || hash([]byte(ch.Content)) != ch.After {
					return fail("reconciliation_required", "invalid assessment effect")
				}
				changed[ch.Path] = true
				switch ch.Path {
				case assessmentPath(q):
					if ch.Before != "absent" || ch.Content != string(assessmentArtifact(a, i)) {
						return fail("reconciliation_required", "assessment artifact differs")
					}
				case managementLedger: // The appended records are validated under the canonical writer lock before application.
				default:
					allowed := ch.Path == i.Context.Binding.Session.Path || ch.Path == i.Context.Binding.Plan.Path
					if strings.HasPrefix(ch.Path, "SDP/KanBan/active/") && required[ch.Path] != "" {
						m, e := metadata([]byte(ch.Content))
						allowed = e == nil && cards[m["id"]]
						if allowed {
							delete(cards, m["id"])
						}
					}
					prefix, ok := strings.CutSuffix(ch.Content, assessmentJournal(a, i, ch.Path))
					if !allowed || !ok || hash([]byte(prefix)) != ch.Before {
						return fail("reconciliation_required", "assessment journal differs or writes unselected path")
					}
				}
				required[ch.Path] = ch.After
			}
			if size > 4<<20 || len(cards) > 0 || !changed[assessmentPath(q)] || !changed[managementLedger] || !changed[i.Context.Binding.Session.Path] || !changed[i.Context.Binding.Plan.Path] || len(required) != len(p.Coverage) {
				return fail("reconciliation_required", "incomplete assessment record effects")
			}
			pins := map[string]bool{}
			for _, pin := range p.Coverage {
				if pins[pin.Path] || required[pin.Path] == "" || (!digestPattern.MatchString(pin.SHA256) && pin.SHA256 != "absent") || (pin.Path != managementLedger && required[pin.Path] != pin.SHA256) {
					return fail("reconciliation_required", "assessment coverage differs")
				}
				pins[pin.Path] = true
			}
			previous = p
		}
	}
	for _, r := range st.Receipts {
		if v := r.Result.Session; v != nil && v.Operation == "session-assess" && !seen[v.OperationID] {
			return fail("reconciliation_required", "orphan assessment receipt")
		}
	}
	return nil
}
