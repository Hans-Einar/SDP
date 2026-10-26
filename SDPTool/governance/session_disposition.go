package governance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
)

const sessionDispositionStateSchema = "sdp-project-governance-state/0.18"
const dispositionTrace = "SDP/Traceability/Ledger.ndjson"
const dispositionRunReason = "Selected Session step disposed; canonical results committed and claim released"

type SessionCardDisposition struct {
	ID        string `json:"id"`
	Outcome   string `json:"outcome"`
	Summary   string `json:"summary"`
	Remaining string `json:"remaining"`
}
type SessionDisposeRequest struct {
	Schema              string                   `json:"schema"`
	ProjectRef          string                   `json:"projectRef"`
	CheckoutID          string                   `json:"checkoutId"`
	OperationID         string                   `json:"operationId"`
	InstanceID          string                   `json:"instanceId"`
	ExpectedRevision    uint64                   `json:"expectedRevision"`
	ExpectedRunRevision uint64                   `json:"expectedRunRevision"`
	ResultID            string                   `json:"resultId"`
	System              string                   `json:"system"`
	Reason              string                   `json:"reason"`
	Cards               []SessionCardDisposition `json:"cards"`
}
type SessionDisposition struct {
	Request SessionDisposeRequest `json:"request"`
	Effects SessionRecordEffects  `json:"effects"`
	Before  map[string]string     `json:"before"`
}

func dispositionPath(d SessionDisposition) string {
	return "SDP/Sessions/execution/dispositions/" + d.Request.OperationID + ".json"
}
func dispositionClaim(i SessionInstance) string {
	return "SDP/Sessions/execution/claims/" + i.Selection.Request.OperationID + ".json"
}
func dispositionResult(i SessionInstance) SessionWork { return i.Work[len(i.Work)-1] }
func validSessionDispose(q SessionDisposeRequest) error {
	if q.Schema != SessionSchema || !uuidPattern.MatchString(q.OperationID) || !uuidPattern.MatchString(q.InstanceID) || !uuidPattern.MatchString(q.ResultID) || q.ExpectedRevision == 0 || q.ExpectedRunRevision == 0 || !regexp.MustCompile(`^[A-Z][A-Z0-9]{0,31}$`).MatchString(q.System) || !validLeaderText(q.Reason, 2048) || len(q.Cards) == 0 || len(q.Cards) > 16 {
		return fail("invalid_request", "exact bounded disposition identity, evidence and cards required")
	}
	seen := map[string]bool{}
	for _, c := range q.Cards {
		if seen[c.ID] || !sessionCardPattern.MatchString(c.ID) || !validLeaderText(c.Summary, 4096) || !slices.Contains([]string{"completed", "contribution-completed"}, c.Outcome) || (c.Outcome == "completed" && c.Remaining != "") || (c.Outcome == "contribution-completed" && !validLeaderText(c.Remaining, 4096)) {
			return fail("invalid_request", "unique explicit card dispositions and remaining contribution scope required")
		}
		seen[c.ID] = true
	}
	return nil
}
func dispositionEligibility(i SessionInstance, q SessionDisposeRequest) error {
	if i.Selection == nil || len(i.Work) < 3 || len(i.Work) > 4 || i.Work[2].Evidence == nil {
		return fail("missing_coverage", "committed Session verification/review required")
	}
	step := selectedSessionStep(i)
	w := dispositionResult(i)
	if step.OwnerAcceptance != "not-required" && !sessionOwnerAccepted(i) {
		return fail("missing_coverage", "actual committed owner acceptance required")
	}
	if w.Effects.Status != "committed" || w.Request.OperationID != q.ResultID || (step.Review == "independent" && w.Kind != "review") || (w.Kind != "verify" && w.Kind != "review") {
		return fail("missing_coverage", "exact committed result required")
	}
	if len(q.Cards) != len(step.Cards) {
		return fail("missing_coverage", "all selected cards require disposition")
	}
	for n, c := range step.Cards {
		d := q.Cards[n]
		if d.ID != c.ID || (c.Coverage == "whole" && d.Outcome != "completed") || (c.Coverage == "contribution" && d.Outcome != "contribution-completed") {
			return fail("conflict", "card outcome must match ordered selected coverage")
		}
		if c.Coverage == "whole" {
			for _, other := range i.Context.Binding.Steps {
				if other.ID == step.ID || sessionStepDone(i, other.ID) {
					continue
				}
				for _, ref := range other.Cards {
					if ref.ID == c.ID {
						return fail("conflict", "another bound step retains card scope")
					}
				}
			}
		}
	}
	return nil
}
func (s *Service) dispositionBasis(st snapshot, i SessionInstance, d SessionDisposition) error {
	if e := dispositionEligibility(i, d.Request); e != nil {
		return e
	}
	r := st.Runs[i.Work[0].Request.OperationID]
	if r.Revision != d.Request.ExpectedRunRevision || r.Revision == ^uint64(0) {
		return fail("stale_revision", "disposition run revision changed or exhausted")
	}
	return s.sessionResultSnapshot(st, i, dispositionResult(i))
}
func dispositionArtifact(i SessionInstance, d SessionDisposition) []byte {
	w := dispositionResult(i)
	ownerLimit := "owner acceptance explicitly not required by selected binding"
	if sessionOwnerAccepted(i) {
		ownerLimit = "actual owner acceptance recorded at " + ownerGatePath(latestOwnerGate(i).Request)
	}
	b, _ := json.MarshalIndent(struct {
		Schema        string                `json:"schema"`
		Request       SessionDisposeRequest `json:"request"`
		CreatedAt     string                `json:"createdAt"`
		Binding       SessionBinding        `json:"binding"`
		BindingDigest string                `json:"bindingDigest"`
		RunID         string                `json:"runId"`
		ResultPath    string                `json:"resultPath"`
		Evidence      *EvidenceRecord       `json:"evidence,omitempty"`
		Review        *Review               `json:"review,omitempty"`
		Tests         []EvidenceRecord      `json:"tests,omitempty"`
		Limits        []string              `json:"limitations"`
	}{"sdp-session-disposition/0.1", d.Request, d.Effects.CreatedAt, i.Context.Binding, i.Context.BindingDigest, i.Work[0].Request.OperationID, workPath(w), w.Evidence, w.Review, w.Tests, []string{"external tests and coverage judgment are controller-attributed", ownerLimit, "selected step only; plan/Session closure is separate"}}, "", "  ")
	return append(b, '\n')
}
func dispositionJournal(i SessionInstance, d SessionDisposition, path string, c *SessionCardDisposition) string {
	rel, _ := filepath.Rel(filepath.Dir(path), dispositionPath(d))
	out := "\n## PG Session disposition " + d.Request.OperationID + "\n\nSelected step " + selectedSessionStep(i).ID + " completed; plan/Session closure is separate.\n[Exact result and disposition](" + filepath.ToSlash(rel) + ").\n\nReason:\n\n" + publicationProse(d.Request.Reason)
	if c != nil {
		out += "\nCard outcome: " + c.Outcome + ".\n\nWork summary:\n\n" + publicationProse(c.Summary) + "\nRemaining card work:\n\n"
		if c.Remaining == "" {
			out += "None within the selected whole-card outcome.\n"
		} else {
			out += publicationProse(c.Remaining)
		}
	}
	return out
}
func dispositionInputs(i SessionInstance) ([]string, map[string]string, error) {
	paths := []string{i.Context.Binding.Session.Path, i.Context.Binding.Plan.Path, managementLedger, "SDP/KanBan/README.md", dispositionTrace, dispositionClaim(i)}
	cards := map[string]string{}
	for _, src := range dispositionPredecessor(i).Coverage {
		if !strings.HasPrefix(src.Path, "SDP/KanBan/active/") {
			continue
		}
		for _, ch := range dispositionResult(i).Effects.Changes {
			if ch.Path != src.Path {
				continue
			}
			m, e := metadata([]byte(ch.Content))
			if e == nil {
				cards[m["id"]] = src.Path
			}
		}
	}
	for _, c := range selectedSessionStep(i).Cards {
		p := cards[c.ID]
		if p == "" {
			return nil, nil, fail("missing_coverage", "selected active card absent from result")
		}
		paths = append(paths, p)
	}
	return paths, cards, nil
}

// Deterministic record transformation, shared by preparation and retained-state
// validation. Only the named selected records may change; original event bytes
// and the exact reviewed context remain historical.
func planDisposition(i SessionInstance, d SessionDisposition) (map[string]string, error) {
	paths, cards, e := dispositionInputs(i)
	if e != nil {
		return nil, e
	}
	files := map[string]string{}
	for p, b := range d.Before {
		files[p] = b
	}
	for _, p := range paths {
		if _, ok := files[p]; !ok && p != dispositionTrace {
			return nil, fail("missing_coverage", "disposition preimage absent: "+p)
		}
	}
	events, e := publicationEvents([]byte(files[managementLedger]))
	if e != nil {
		return nil, e
	}
	latest := map[string]map[string]any{}
	for _, ev := range events {
		latest[fmt.Sprint(ev["subjectId"])] = ev
	}
	addEvent := func(id, to, path, typ, reason string) error {
		old := latest[id]
		if old == nil {
			return fail("missing_coverage", "selected history absent")
		}
		prev, ok := old["payload"].(map[string]any)
		if !ok {
			return fail("reconciliation_required", "selected history payload malformed")
		}
		p := map[string]any{}
		for k, v := range prev {
			p[k] = v
		}
		delete(p, "lineage")
		p["previousEventId"] = old["eventId"]
		p["from"] = prev["to"]
		p["fromPath"] = prev["toPath"]
		p["to"] = to
		p["toPath"] = path
		p["reason"] = reason + " Disposition " + d.Request.OperationID
		prefix := "EVT-PM-" + strings.Split(id, "-")[1] + "-"
		if strings.HasPrefix(id, "KB-") {
			prefix = "EVT-KB-" + strings.Split(id, "-")[1] + "-"
		}
		eid, e := nextPublicationEvent(events, prefix)
		if e != nil {
			return e
		}
		ev := map[string]any{"schemaVersion": "1.0", "eventId": eid, "eventType": typ, "occurredAt": d.Effects.CreatedAt, "actor": "sdp-controller", "commit": nil, "subjectId": id, "payload": p}
		b, _ := json.Marshal(ev)
		files[managementLedger] = appendRecord([]byte(files[managementLedger]), string(b)+"\n")
		events = append(events, ev)
		return nil
	}
	moved := map[string]string{}
	for _, c := range d.Request.Cards {
		old := cards[c.ID]
		m, e := metadata([]byte(files[old]))
		if e != nil || m["CardState"] != "in-progress" {
			return nil, fail("conflict", "in-progress selected card required")
		}
		dest, state, typ := old, "active", "x-kanban:reviewed"
		if c.Outcome == "completed" {
			dest = "SDP/KanBan/completed/" + filepath.Base(old)
			state = "completed"
			typ = "x-kanban:moved"
			b, e := replaceMetadata([]byte(files[old]), "CardState", "completed")
			if e != nil {
				return nil, e
			}
			b, e = rebaseMovedCardLinks(b, old, dest)
			if e != nil {
				return nil, e
			}
			files[dest] = string(b)
			delete(files, old)
			moved[old] = dest
		}
		files[dest] = appendRecord([]byte(files[dest]), dispositionJournal(i, d, dest, &c))
		if e = addEvent(c.ID, state, strings.TrimPrefix(dest, "SDP/KanBan/"), typ, "Recorded "+c.Outcome+"; remaining scope is explicit."); e != nil {
			return nil, e
		}
	}
	plan := i.Context.Binding.Plan.Path
	files[plan] = appendRecord([]byte(files[plan]), dispositionJournal(i, d, plan, nil))
	if e = addEvent(i.Context.Binding.Plan.ID, "active", strings.TrimPrefix(plan, "SDP/"), "x-management:updated", "Selected step completed; overall plan remains active."); e != nil {
		return nil, e
	}
	pmEvent := fmt.Sprint(events[len(events)-1]["eventId"])
	sp := i.Context.Binding.Session.Path
	files[sp] = appendRecord([]byte(files[sp]), dispositionJournal(i, d, sp, nil))
	for old, dest := range moved {
		for p, b := range files {
			if strings.HasSuffix(p, ".md") {
				files[p] = string(rewriteCardLinks([]byte(b), p, old, dest))
			}
		}
		lines := strings.Split(files["SDP/KanBan/README.md"], "\n")
		m, _ := metadata([]byte(files[dest]))
		found := 0
		for j, line := range lines {
			cells := strings.Split(line, "|")
			if len(cells) >= 6 && strings.TrimSpace(cells[1]) == m["id"] {
				cells[3] = " completed "
				lines[j] = strings.Join(cells, "|")
				found++
			}
		}
		if found != 1 {
			return nil, fail("conflict", "one board index row required")
		}
		files["SDP/KanBan/README.md"] = strings.Join(lines, "\n")
	}
	delete(files, dispositionClaim(i))
	files[dispositionPath(d)] = string(dispositionArtifact(i, d))
	trace, e := publicationEvents([]byte(files[dispositionTrace]))
	if e != nil {
		return nil, e
	}
	eid := "EVT-" + d.Request.System + "-PG-" + strings.ToUpper(d.Request.OperationID)
	for _, ev := range trace {
		if ev["eventId"] == eid {
			return nil, fail("conflict", "trace event identity collision")
		}
	}
	w := dispositionResult(i)
	candidate := i.Work[2].Evidence.Request.Candidate.Digest
	ev := map[string]any{"schemaVersion": "1.0", "eventId": eid, "eventType": "x-verification:recorded", "occurredAt": d.Effects.CreatedAt, "actor": "sdp-controller", "commit": nil, "subjectId": d.Request.System + "-VER-PG-" + strings.ToUpper(d.Request.OperationID), "payload": map[string]any{"plan": i.Context.Binding.Plan.ID, "cards": d.Request.Cards, "session": i.Context.Binding.Session.ID, "managementEvent": pmEvent, "evidence": strings.TrimPrefix(dispositionPath(d), "SDP/"), "candidate": candidate, "result": w.Request.OperationID, "note": "Published scoped candidate result with required review disposition; no owner acceptance inferred."}}
	raw, _ := json.Marshal(ev)
	files[dispositionTrace] = appendRecord([]byte(files[dispositionTrace]), string(raw)+"\n")
	return files, nil
}

func dispositionChanges(before, after map[string]string) []PublicationChange {
	names := map[string]bool{}
	for p := range before {
		names[p] = true
	}
	for p := range after {
		names[p] = true
	}
	changes := []PublicationChange{}
	for p := range names {
		a, b := "absent", "absent"
		if v, ok := before[p]; ok {
			a = hash([]byte(v))
		}
		if v, ok := after[p]; ok {
			b = hash([]byte(v))
		}
		if a != b {
			changes = append(changes, PublicationChange{p, a, b, after[p]})
		}
	}
	slices.SortFunc(changes, func(a, b PublicationChange) int {
		if (a.After == "absent") != (b.After == "absent") {
			if a.After == "absent" {
				return 1
			}
			return -1
		}
		return strings.Compare(a.Path, b.Path)
	})
	return changes
}
func dispositionCoverage(i SessionInstance, d SessionDisposition) []Source {
	pins := map[string]string{}
	for _, src := range dispositionPredecessor(i).Coverage {
		pins[src.Path] = src.SHA256
	}
	for _, ch := range d.Effects.Changes {
		pins[ch.Path] = ch.After
	}
	names := []string{}
	for p := range pins {
		names = append(names, p)
	}
	sort.Strings(names)
	out := []Source{}
	for _, p := range names {
		out = append(out, Source{p, pins[p]})
	}
	return out
}
func (s *Service) buildSessionDisposition(i SessionInstance, d SessionDisposition) (SessionDisposition, error) {
	files, e := s.publicationFiles()
	if e != nil {
		return d, e
	}
	for _, src := range dispositionPredecessor(i).Coverage {
		if _, ok := files[src.Path]; ok {
			continue
		}
		b, err := readFile(filepath.Join(s.root, src.Path), maxState)
		if os.IsNotExist(err) && src.SHA256 == "absent" {
			continue
		}
		if err != nil {
			return d, err
		}
		files[src.Path] = b
	}
	paths, cards, e := dispositionInputs(i)
	if e != nil {
		return d, e
	}
	d.Before = map[string]string{}
	for _, p := range paths {
		b, ok := files[p]
		if !ok {
			b, e = readFile(filepath.Join(s.root, p), maxState)
			if os.IsNotExist(e) && p == dispositionTrace {
				continue
			}
			if e != nil {
				return d, e
			}
			files[p] = b
		}
		d.Before[p] = string(b)
	}
	after, e := planDisposition(i, d)
	if e != nil {
		return d, e
	}
	d.Effects.Changes = dispositionChanges(d.Before, after)
	allowed := map[string]bool{i.Context.Binding.Plan.Path: true, i.Context.Binding.Session.Path: true, "SDP/KanBan/README.md": true}
	for _, p := range cards {
		allowed[p] = true
		allowed["SDP/KanBan/completed/"+filepath.Base(p)] = true
	}
	for _, c := range d.Request.Cards {
		if c.Outcome == "completed" {
			if e = s.checkSelectionReferences(cards[c.ID], allowed); e != nil {
				return d, e
			}
		}
	}
	candidate := i.Work[2].Evidence.Request.Candidate
	for _, ch := range d.Effects.Changes {
		if ch.Before == "absent" {
			if _, err := os.Lstat(filepath.Join(s.root, ch.Path)); !os.IsNotExist(err) {
				return d, fail("conflict", "disposition destination exists or unavailable: "+ch.Path)
			}
		}
		for _, p := range append(append([]string{}, candidate.WritablePaths...), candidate.Dependencies...) {
			if overlaps(p, ch.Path) {
				return d, fail("conflict", "disposition overlaps candidate scope")
			}
		}
		if ch.Path == i.Context.Binding.Authority {
			return d, fail("conflict", "disposition cannot rewrite authority")
		}
	}
	for p := range d.Before {
		delete(files, p)
	}
	for p, b := range after {
		files[p] = []byte(b)
	}
	if e = validatePublicationHistory(files[managementLedger], files); e != nil {
		return d, e
	}
	pre := map[string]string{}
	for p, b := range files {
		pre[p] = hash(b)
	}
	for p, b := range d.Before {
		pre[p] = hash([]byte(b))
	}
	for _, ch := range d.Effects.Changes {
		pre[ch.Path] = ch.Before
	}
	for _, src := range dispositionPredecessor(i).Coverage {
		if _, ok := pre[src.Path]; !ok {
			pre[src.Path] = src.SHA256
		}
	}
	names := []string{}
	for p := range pre {
		names = append(names, p)
	}
	sort.Strings(names)
	for _, p := range names {
		d.Effects.Inputs = append(d.Effects.Inputs, Source{p, pre[p]})
	}
	d.Effects.Coverage = dispositionCoverage(i, d)
	d.Effects.LedgerBytes = len(after[managementLedger])
	size := 0
	for _, b := range d.Before {
		size += len(b)
	}
	for _, ch := range d.Effects.Changes {
		size += len(ch.Content)
	}
	if len(d.Effects.Inputs) > 2048 || len(d.Effects.Changes) > 128 || size > 8<<20 {
		return d, fail("limit", "disposition pre/postimage capacity exceeded")
	}
	return d, nil
}
func (s *Service) DisposeSession(c Caller, q SessionDisposeRequest) (SessionMutationResult, error) {
	if !c.controller {
		return SessionMutationResult{}, fail("unauthorized", "trusted controller required")
	}
	if e := validSessionDispose(q); e != nil {
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
	if out, ok, e := sessionRetry(st, q.OperationID, requestDigest(q)); ok {
		return out, e
	}
	if e = tx.Check(""); e != nil {
		return SessionMutationResult{}, e
	}
	if _, e = readFile(filepath.Join(dir, "pending"), 8192); !os.IsNotExist(e) {
		return SessionMutationResult{}, fail("reconciliation_required", "canonical publication pending")
	}
	i, ok := st.Sessions[q.InstanceID]
	if !ok {
		return SessionMutationResult{}, fail("missing_coverage", "Session absent")
	}
	if i.Revision != q.ExpectedRevision {
		return SessionMutationResult{}, fail("stale_revision", "Session revision changed")
	}
	if i.Disposition != nil {
		return SessionMutationResult{}, fail("conflict", "disposition already prepared")
	}
	if i.Revision == ^uint64(0) || len(st.Receipts) >= 1024 {
		return SessionMutationResult{}, fail("limit", "Session capacity exhausted")
	}
	if e = s.sessionCurrent(i); e != nil {
		return SessionMutationResult{}, e
	}
	d := SessionDisposition{Request: q, Effects: SessionRecordEffects{Status: "prepared", CreatedAt: time.Now().UTC().Format(time.RFC3339Nano)}}
	if e = s.dispositionBasis(st, i, d); e != nil {
		return SessionMutationResult{}, e
	}
	if e = s.sessionWorkClaims(st, i); e != nil {
		return SessionMutationResult{}, e
	}
	d, e = s.buildSessionDisposition(i, d)
	if e != nil {
		return SessionMutationResult{}, e
	}
	i.Disposition = &d
	i.Revision++
	st.Sessions[i.ID] = i
	if st.Schema != sessionOwnerStateSchema && !sessionLifecycleProfile(st.Schema) {
		st.Schema = sessionDispositionStateSchema
	}
	out := SessionMutationResult{Schema: SessionSchema, Operation: "session-dispose", OperationID: q.OperationID, InstanceID: i.ID, Revision: i.Revision, Registration: "committed", RecordEffects: "prepared"}
	st.Receipts[q.OperationID] = receipt{requestDigest(q), Result{Schema: Schema, Operation: "session.dispose", OperationID: q.OperationID, ProjectRef: s.project, CheckoutID: st.CheckoutID, Session: &out}}
	if e = validateSessionDispositions(st); e != nil {
		return SessionMutationResult{}, e
	}
	if e = s.save(st, head); e != nil {
		return SessionMutationResult{}, e
	}
	return out, nil
}
func (s *Service) applySessionDisposition(st *snapshot, i *SessionInstance, tx *projecthistory.Transaction) error {
	if e := s.sessionWorkClaims(*st, *i); e != nil {
		return e
	}
	d := i.Disposition
	check := func() error { return s.dispositionBasis(*st, *i, *d) }
	if e := s.applySessionRecordEffects(&d.Effects, tx, s.sessionPendingOwner(*i), "session-disposition", check); e != nil {
		return e
	}
	r := st.Runs[i.Work[0].Request.OperationID]
	r.State = "completed"
	r.Revision++
	r.Reason = dispositionRunReason
	st.Runs[r.ID] = r
	return nil
}
func validateSessionDispositions(st snapshot) error {
	seen := map[string]bool{}
	for _, i := range sessionExecutions(st) {
		d := i.Disposition
		if d == nil {
			if len(i.Work) > 0 && st.Runs[i.Work[0].Request.OperationID].State == "completed" {
				return fail("reconciliation_required", "Session run completed without disposition")
			}
			continue
		}
		q, p := d.Request, d.Effects
		if !sessionDispositionProfile(st.Schema) || validSessionDispose(q) != nil || dispositionEligibility(i, q) != nil || q.ProjectRef != st.Project || q.CheckoutID != st.CheckoutID || q.InstanceID != i.ID || q.ExpectedRevision >= i.Revision || q.ExpectedRevision <= dispositionResult(i).Request.ExpectedRevision || len(p.Inputs) > 2048 || len(p.Changes) > 128 || len(p.Coverage) > 512 || p.LedgerBytes <= 0 {
			return fail("reconciliation_required", "invalid Session disposition identity")
		}
		if _, e := time.Parse(time.RFC3339Nano, p.CreatedAt); e != nil {
			return fail("reconciliation_required", "invalid disposition time")
		}
		r := st.Runs[i.Work[0].Request.OperationID]
		if (p.Status == "prepared" && (i.Revision != q.ExpectedRevision+1 || r.State == "completed")) || (p.Status == "committed" && (i.Revision <= q.ExpectedRevision+1 || r.State != "completed" || r.Revision <= q.ExpectedRunRevision)) || !slices.Contains([]string{"prepared", "committed"}, p.Status) {
			return fail("reconciliation_required", "invalid disposition/run commitment")
		}
		receipt, ok := st.Receipts[q.OperationID]
		if !ok || receipt.Digest != requestDigest(q) || receipt.Result.Session == nil || receipt.Result.Session.Operation != "session-dispose" || receipt.Result.Session.InstanceID != i.ID || receipt.Result.Session.Revision != q.ExpectedRevision+1 || receipt.Result.Session.RecordEffects != "prepared" {
			return fail("reconciliation_required", "disposition receipt differs")
		}
		seen[q.OperationID] = true
		paths, _, e := dispositionInputs(i)
		if e != nil {
			return e
		}
		for path := range d.Before {
			if !slices.Contains(paths, path) {
				return fail("reconciliation_required", "unselected disposition preimage")
			}
		}
		after, e := planDisposition(i, *d)
		if e != nil {
			return e
		}
		expected := dispositionChanges(d.Before, after)
		if !reflect.DeepEqual(expected, p.Changes) {
			return fail("reconciliation_required", "disposition postimages differ from selected transformation")
		}
		inputs := map[string]string{}
		for _, src := range p.Inputs {
			if !literalCandidatePath(src.Path) || inputs[src.Path] != "" || (!digestPattern.MatchString(src.SHA256) && src.SHA256 != "absent") {
				return fail("reconciliation_required", "invalid disposition inputs")
			}
			inputs[src.Path] = src.SHA256
		}
		for _, src := range dispositionPredecessor(i).Coverage {
			if src.Path != managementLedger && inputs[src.Path] != src.SHA256 {
				return fail("reconciliation_required", "disposition lost result coverage")
			}
		}
		size := 0
		for path, b := range d.Before {
			size += len(b)
			if inputs[path] != hash([]byte(b)) {
				return fail("reconciliation_required", "disposition preimage digest differs")
			}
		}
		for _, ch := range p.Changes {
			size += len(ch.Content)
			if inputs[ch.Path] != ch.Before {
				return fail("reconciliation_required", "disposition effect preimage differs")
			}
		}
		if size > 8<<20 {
			return fail("reconciliation_required", "disposition capacity exceeded")
		}
		want := dispositionCoverage(i, *d)
		if len(want) != len(p.Coverage) {
			return fail("reconciliation_required", "disposition coverage differs")
		}
		for n, src := range want {
			got := p.Coverage[n]
			if src.Path != got.Path || (src.Path != managementLedger && src.SHA256 != got.SHA256) || !digestPattern.MatchString(got.SHA256) && got.SHA256 != "absent" {
				return fail("reconciliation_required", "disposition coverage differs")
			}
		}
	}
	for _, r := range st.Receipts {
		if v := r.Result.Session; v != nil && v.Operation == "session-dispose" && !seen[v.OperationID] {
			return fail("reconciliation_required", "orphan disposition receipt")
		}
	}
	return nil
}
func (s *Service) projectSessionDisposition(i SessionInstance, st snapshot, v *SessionView, step *SessionStepView) {
	if !step.Selected {
		return
	}
	if i.Disposition == nil {
		if (step.Stage == "disposition") && selectedSessionStep(i).OwnerAcceptance == "not-required" && v.Coverage == "current" {
			v.NextActions = append(v.NextActions, "session-dispose")
		}
		return
	}
	d := i.Disposition
	step.DispositionID = d.Request.OperationID
	step.Holds = append([]string{}, v.Conditions...)
	v.NextActions = []string{"session-read", "session-resume"}
	if d.Effects.Status != "committed" {
		step.Holds = append(step.Holds, "disposition effects pending; run and claim not yet committed")
		return
	}
	step.Stage = "done"
	step.CardDispositions = d.Request.Cards
	step.Holds = append(step.Holds, "selected step completed; plan/Session closure and other steps remain separate")
}
