package governance

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
)

const sessionEffectsStateSchema = "sdp-project-governance-state/0.14"
const managementLedger = "SDP/ProjectManagement/Ledger.ndjson"

type SessionSelectRequest struct {
	Schema           string `json:"schema"`
	ProjectRef       string `json:"projectRef"`
	CheckoutID       string `json:"checkoutId"`
	OperationID      string `json:"operationId"`
	InstanceID       string `json:"instanceId"`
	ExpectedRevision uint64 `json:"expectedRevision"`
	StepID           string `json:"stepId"`
	Reason           string `json:"reason"`
}
type SessionSelection struct {
	Request SessionSelectRequest `json:"request"`
	SessionRecordEffects
}
type SessionRecordEffects struct {
	Status      string              `json:"status"`
	Inputs      []Source            `json:"inputs"`
	Changes     []PublicationChange `json:"changes"`
	Coverage    []Source            `json:"coverage"`
	LedgerBytes int                 `json:"ledgerBytes"`
	CreatedAt   string              `json:"createdAt"`
}

func sessionProfile(schema string) bool {
	return schema == sessionStateSchema || sessionEffectsProfile(schema)
}
func (s *Service) canonicalTransaction() (*projecthistory.Transaction, error) {
	return projecthistory.BeginTransaction(filepath.Join(s.root, managementLedger))
}
func (s *Service) sessionEffectOwner(i SessionInstance) string {
	return s.state + "\n" + i.Selection.Request.OperationID + "\n"
}

// The selection pins the full validated original ledger and accepts only an exact
// prefix plus unrelated events. Selected documents remain byte-pinned. Unknown
// Blueprint events conservatively hold until execution claims can be reconciled.
func (s *Service) selectionCoverage(i SessionInstance) error {
	sel := sessionLatestEffects(i)
	if sel.Status != "committed" {
		return fail("reconciliation_required", "Session record effects pending")
	}
	if _, e := s.publicationFiles(); e != nil {
		return e
	}
	subjects := []string{i.Context.Binding.Session.ID, i.Context.Binding.Plan.ID}
	for _, step := range i.Context.Binding.Steps {
		for _, card := range step.Cards {
			subjects = append(subjects, card.ID)
		}
	}
	for _, src := range sel.Coverage {
		raw, e := readFile(filepath.Join(s.root, src.Path), maxState)
		actual := "absent"
		if e == nil {
			actual = hash(raw)
		} else if !os.IsNotExist(e) {
			return e
		}
		if src.Path == managementLedger && actual != src.SHA256 {
			if sel.LedgerBytes <= 0 || len(raw) < sel.LedgerBytes || hash(raw[:sel.LedgerBytes]) != src.SHA256 {
				return fail("stale_revision", "canonical history prefix changed")
			}
			events, e := publicationEvents(raw[sel.LedgerBytes:])
			if e != nil {
				return e
			}
			for _, ev := range events {
				if strings.HasPrefix(fmt.Sprint(ev["eventType"]), "x-blueprint:") {
					return fail("stale_revision", "new Blueprint history requires scope reconciliation")
				}
				b, _ := json.Marshal(ev)
				for _, subject := range subjects {
					if strings.Contains(string(b), subject) {
						return fail("stale_revision", "selected history changed: "+subject)
					}
				}
			}
			continue
		}
		if actual != src.SHA256 {
			return fail("stale_revision", "selected Session coverage changed: "+src.Path)
		}
	}
	return nil
}
func replaceMetadata(raw []byte, key, value string) ([]byte, error) {
	lines := strings.Split(string(raw), "\n")
	n := 0
	for j, line := range lines {
		cells := strings.Split(line, "|")
		if len(cells) >= 4 && strings.TrimSpace(cells[1]) == key {
			lines[j] = "| " + key + " | " + value + " |"
			n++
		}
	}
	if n != 1 {
		return nil, fail("conflict", "exactly one metadata row required: "+key)
	}
	return []byte(strings.Join(lines, "\n")), nil
}
func rewriteCardLinks(raw []byte, doc, from, to string) []byte {
	// Only ordinary local Markdown links are rewritten. Immutable historical event
	// paths remain unchanged. Outgoing links are rebased separately before moves.
	old, _ := filepath.Rel(filepath.Dir(doc), from)
	next, _ := filepath.Rel(filepath.Dir(doc), to)
	escape := func(p string) string { return (&url.URL{Path: filepath.ToSlash(p)}).EscapedPath() }
	s := string(raw)
	for _, pair := range [][2]string{{filepath.ToSlash(old), filepath.ToSlash(next)}, {escape(old), escape(next)}} {
		s = strings.ReplaceAll(s, "]("+pair[0]+")", "]("+pair[1]+")")
		s = strings.ReplaceAll(s, "]("+pair[0]+"#", "]("+pair[1]+"#")
	}
	return []byte(s)
}
func (s *Service) buildSessionSelection(i SessionInstance, q SessionSelectRequest) (SessionSelection, error) {
	out := SessionSelection{Request: q, SessionRecordEffects: SessionRecordEffects{Status: "prepared", CreatedAt: time.Now().UTC().Format(time.RFC3339Nano)}}
	files, e := s.publicationFiles()
	if e != nil {
		return out, e
	}
	for _, src := range sessionSelectionSources(i) {
		if _, ok := files[src.Path]; !ok {
			b, e := readFile(filepath.Join(s.root, src.Path), maxState)
			if os.IsNotExist(e) && src.SHA256 == "absent" {
				continue
			}
			if e != nil {
				return out, e
			}
			files[src.Path] = b
		}
	}
	step := SessionStepBinding{}
	for _, v := range i.Context.Binding.Steps {
		if v.ID == q.StepID {
			step = v
		}
	}
	if step.ID == "" {
		return out, fail("missing_coverage", "selected step absent")
	}
	if !sessionDependenciesSatisfied(i, step) || sessionStepDone(i, step.ID) {
		return out, fail("conflict", "step already completed or predecessors incomplete")
	}
	if e = s.checkSelectionClaims(i, q, step); e != nil {
		return out, e
	}
	claimDir := "SDP/Sessions/execution/claims"
	before := map[string][]byte{}
	for path, b := range files {
		before[path] = b
	}
	journal := "\n## PG step selection " + q.OperationID + "\n\nSelected step " + q.StepID + " under " + i.Context.Binding.Plan.ID + ".\nWork is ready for design assessment; no host dispatch or completion.\n\nReason:\n\n" + publicationProse(q.Reason)
	events, e := publicationEvents(files[managementLedger])
	if e != nil {
		return out, e
	}
	latest := map[string]map[string]any{}
	for _, ev := range events {
		latest[fmt.Sprint(ev["subjectId"])] = ev
	}
	history := files[managementLedger]
	addEvent := func(subject, to, path, typ, reason string) error {
		old := latest[subject]
		if old == nil {
			return fail("missing_coverage", "history absent: "+subject)
		}
		prev := old["payload"].(map[string]any)
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
		p["reason"] = reason + " Selection " + q.OperationID
		p["links"] = []string{i.Context.Binding.Plan.ID}
		parts := strings.Split(subject, "-")
		prefix := "EVT-PM-" + parts[1] + "-"
		if strings.HasPrefix(subject, "KB-") {
			prefix = "EVT-KB-" + parts[1] + "-"
		}
		id, e := nextPublicationEvent(events, prefix)
		if e != nil {
			return e
		}
		ev := map[string]any{"schemaVersion": "1.0", "eventId": id, "eventType": typ, "occurredAt": out.CreatedAt, "actor": "sdp-controller", "commit": nil, "subjectId": subject, "payload": p}
		b, _ := json.Marshal(ev)
		history = []byte(appendRecord(history, string(b)+"\n"))
		events = append(events, ev)
		return nil
	}
	plan := i.Context.Binding.Plan.Path
	pm, _ := metadata(files[plan])
	planState := pm["state"]
	activationCards := append([]SessionCardRef{}, step.Cards...)
	if planState == "planned" {
		seen := map[string]bool{}
		for _, c := range activationCards {
			seen[c.ID] = true
		}
		for _, future := range i.Context.Binding.Steps {
			for _, c := range future.Cards {
				if !seen[c.ID] {
					activationCards = append(activationCards, c)
					seen[c.ID] = true
				}
			}
		}
	}
	typ := "x-management:updated"
	if planState == "planned" {
		files[plan], e = replaceMetadata(files[plan], "state", "active")
		if e != nil {
			return out, e
		}
		planState = "active"
		typ = "x-management:started"
	}
	if e = addEvent(i.Context.Binding.Plan.ID, planState, strings.TrimPrefix(plan, "SDP/"), typ, "Selected Session step; plan active."); e != nil {
		return out, e
	}
	files[plan] = []byte(appendRecord(files[plan], journal))
	session := i.Context.Binding.Session.Path
	files[session] = []byte(appendRecord(files[session], journal))
	cardPaths := map[string]string{}
	for path, b := range files {
		if strings.HasPrefix(path, "SDP/KanBan/") {
			m, e := metadata(b)
			if e == nil {
				cardPaths[m["id"]] = path
			}
		}
	}
	moved := map[string]string{}
	for _, card := range activationCards {
		old := cardPaths[card.ID]
		if old == "" {
			return out, fail("missing_coverage", "selected card absent")
		}
		m, _ := metadata(files[old])
		state := m["CardState"]
		dest := old
		evtype := "x-kanban:reviewed"
		if state == "backlog" || state == "queued" {
			dest = "SDP/KanBan/active/" + filepath.Base(old)
			if _, ok := files[dest]; ok {
				return out, fail("conflict", "card move destination exists")
			}
			files[dest], e = replaceMetadata(files[old], "CardState", "ready")
			if e != nil {
				return out, e
			}
			delete(files, old)
			moved[old] = dest
			evtype = "x-kanban:moved"
		} else if !slices.Contains([]string{"ready", "in-progress"}, state) {
			return out, fail("conflict", "card cannot be selected from current state")
		}
		if dest != old {
			files[dest], e = rebaseMovedCardLinks(files[dest], old, dest)
			if e != nil {
				return out, e
			}
		}
		cardJournal := journal + "\nCard contribution: " + card.Coverage + ". Remaining outcome is not completed by selection.\n"
		if !slices.ContainsFunc(step.Cards, func(c SessionCardRef) bool { return c.ID == card.ID }) {
			cardJournal = "\n## PG plan activation " + q.OperationID + "\n\nReady for a later bound step of " + i.Context.Binding.Plan.ID + "; no assignment started or scope completed.\n"
		}
		files[dest] = []byte(appendRecord(files[dest], cardJournal))
		if e = addEvent(card.ID, "active", strings.TrimPrefix(dest, "SDP/KanBan/"), evtype, "Selected card: "+state+" -> "+map[bool]string{true: "ready", false: state}[dest != old]+"; no execution or completion."); e != nil {
			return out, e
		}
	}
	// Reconcile maintained loaded current documents; do not rewrite ledger bytes.
	// Refuse moves with additional unsupported Markdown references instead of leaving
	// silently broken links or rewriting arbitrary historical material.
	for old, dest := range moved {
		allowed := map[string]bool{plan: true, session: true, "SDP/KanBan/README.md": true}
		for _, c := range activationCards {
			allowed[cardPaths[c.ID]] = true
			if d, ok := moved[cardPaths[c.ID]]; ok {
				allowed[d] = true
			}
		}
		if e = s.checkSelectionReferences(old, allowed); e != nil {
			return out, e
		}
		for path, b := range files {
			if strings.HasSuffix(path, ".md") {
				files[path] = rewriteCardLinks(b, path, old, dest)
			}
		}
		board := files["SDP/KanBan/README.md"]
		lines := strings.Split(string(board), "\n")
		id, _ := metadata(files[dest])
		found := 0
		for j, line := range lines {
			cells := strings.Split(line, "|")
			if len(cells) >= 6 && strings.TrimSpace(cells[1]) == id["id"] {
				cells[3] = " active "
				lines[j] = strings.Join(cells, "|")
				found++
			}
		}
		if found != 1 {
			return out, fail("conflict", "one board index row required")
		}
		files["SDP/KanBan/README.md"] = []byte(strings.Join(lines, "\n"))
	}
	files[managementLedger] = history
	if e = validatePublicationHistory(history, files); e != nil {
		return out, fail("conflict", "proposed selection invalid: "+e.Error())
	}
	claim := sessionScopeClaim{Schema: "sdp-session-scope-claim/0.1", OperationID: q.OperationID, SessionID: i.Context.Binding.Session.ID, PlanID: i.Context.Binding.Plan.ID, WritablePaths: step.WritablePaths, ProtectedPaths: step.ProtectedPaths, Cards: []string{}}
	for _, c := range step.Cards {
		claim.Cards = append(claim.Cards, c.ID)
	}
	claimPath := claimDir + "/" + q.OperationID + ".json"
	claimBytes, _ := json.MarshalIndent(claim, "", "  ")
	files[claimPath] = append(claimBytes, '\n')
	paths := map[string]bool{}
	for p := range before {
		paths[p] = true
	}
	for p := range files {
		paths[p] = true
	}
	names := []string{}
	for p := range paths {
		names = append(names, p)
	}
	sort.Strings(names)
	for _, p := range names {
		old, exists := before[p]
		next, remains := files[p]
		pre, post := "absent", "absent"
		if exists {
			pre = hash(old)
		}
		if remains {
			post = hash(next)
		}
		out.Inputs = append(out.Inputs, Source{p, pre})
		if pre != post {
			out.Changes = append(out.Changes, PublicationChange{Path: p, Before: pre, After: post, Content: string(next)})
		}
	}
	// Creates/updates first; source removals last. The reservation hides partial
	// physical states from cooperating writers until canonical validation succeeds.
	slices.SortStableFunc(out.Changes, func(a, b PublicationChange) int {
		if (a.After == "absent") != (b.After == "absent") {
			if a.After == "absent" {
				return 1
			}
			return -1
		}
		return strings.Compare(a.Path, b.Path)
	})
	coverage := map[string]bool{}
	for _, src := range sessionSelectionSources(i) {
		p := src.Path
		if dest, ok := moved[p]; ok {
			p = dest
		}
		coverage[p] = true
	}
	coverage[claimPath] = true
	names = nil
	for p := range coverage {
		names = append(names, p)
	}
	sort.Strings(names)
	for _, p := range names {
		digest := "absent"
		if b, ok := files[p]; ok {
			digest = hash(b)
		}
		out.Coverage = append(out.Coverage, Source{p, digest})
	}
	out.LedgerBytes = len(history)
	size := 0
	for _, ch := range out.Changes {
		size += len(ch.Content)
	}
	if size > 4<<20 || len(out.Changes) > 128 {
		return out, fail("limit", "bounded selection effect capacity exceeded")
	}
	return out, nil
}

type sessionScopeClaim struct {
	Schema         string   `json:"schema"`
	OperationID    string   `json:"operationId"`
	SessionID      string   `json:"sessionId"`
	PlanID         string   `json:"planId"`
	Cards          []string `json:"cards"`
	WritablePaths  []string `json:"writablePaths"`
	ProtectedPaths []string `json:"protectedPaths"`
}

func (s *Service) SelectSession(c Caller, q SessionSelectRequest) (SessionMutationResult, error) {
	if !c.controller {
		return SessionMutationResult{}, fail("unauthorized", "trusted controller required")
	}
	if q.Schema != SessionSchema || !uuidPattern.MatchString(q.InstanceID) || q.ExpectedRevision == 0 || !validLeaderText(q.Reason, 2048) || !sessionStepPattern.MatchString(q.StepID) {
		return SessionMutationResult{}, fail("invalid_request", "valid Session selection required")
	}
	dir, e := s.publicationLockDir()
	if e != nil {
		return SessionMutationResult{}, e
	}
	up, e := lock(dir)
	if e != nil {
		return SessionMutationResult{}, e
	}
	defer up()
	tx, e := s.canonicalTransaction()
	if e != nil {
		return SessionMutationResult{}, e
	}
	defer tx.Close()
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
	if e = tx.Check(""); e != nil {
		return SessionMutationResult{}, e
	}
	if _, e = readFile(filepath.Join(dir, "pending"), 8192); !os.IsNotExist(e) {
		return SessionMutationResult{}, fail("reconciliation_required", "another publication pending")
	}
	i, ok := st.Sessions[q.InstanceID]
	if !ok {
		return SessionMutationResult{}, fail("missing_coverage", "instance absent")
	}
	if i.Revision != q.ExpectedRevision {
		return SessionMutationResult{}, fail("stale_revision", "instance revision changed")
	}
	if !i.Registered || i.Selection != nil || i.Advance != nil || i.Rebind != nil {
		return SessionMutationResult{}, fail("conflict", "registered unselected instance required")
	}
	if i.Revision == ^uint64(0) || len(st.Receipts) >= 1024 {
		return SessionMutationResult{}, fail("limit", "Session capacity reached")
	}
	if e = s.sessionCurrent(i); e != nil {
		return SessionMutationResult{}, e
	}
	if len(i.Past) > 0 && i.Past[len(i.Past)-1].Transition.Request.StepID != q.StepID {
		return SessionMutationResult{}, fail("conflict", "selection differs from committed continuation")
	}
	selection, e := s.buildSessionSelection(i, q)
	if e != nil {
		return SessionMutationResult{}, e
	}
	i.Selection = &selection
	i.Revision++
	st.Sessions[i.ID] = i
	if !sessionAssessmentProfile(st.Schema) {
		st.Schema = sessionEffectsStateSchema
	}
	out := SessionMutationResult{Schema: SessionSchema, Operation: "session-select", OperationID: q.OperationID, InstanceID: i.ID, Revision: i.Revision, Registration: "committed", RecordEffects: "prepared"}
	st.Receipts[q.OperationID] = receipt{requestDigest(q), Result{Schema: Schema, Operation: "session.select", OperationID: q.OperationID, ProjectRef: s.project, CheckoutID: st.CheckoutID, Session: &out}}
	if e = s.save(st, head); e != nil {
		return SessionMutationResult{}, e
	}
	return out, nil
}
func (s *Service) applySessionSelection(i *SessionInstance, tx *projecthistory.Transaction) error {
	p := i.Selection
	step := SessionStepBinding{}
	for _, v := range i.Context.Binding.Steps {
		if v.ID == p.Request.StepID {
			step = v
		}
	}
	if e := s.checkSelectionClaims(*i, p.Request, step); e != nil {
		return e
	}
	return s.applySessionRecordEffects(&p.SessionRecordEffects, tx, s.sessionEffectOwner(*i), "session-selection")
}

func (s *Service) applySessionRecordEffects(p *SessionRecordEffects, tx *projecthistory.Transaction, owner, faultPrefix string, check ...func() error) error {
	intent := Publication{Inputs: p.Inputs, Changes: p.Changes}
	for _, f := range check {
		if e := f(); e != nil {
			return e
		}
	}
	if e := tx.Check(owner); e != nil {
		return e
	}
	if e := s.publicationInputs(intent, true); e != nil {
		return e
	}
	if e := s.validateSessionEffectRecords(*p); e != nil {
		return e
	}
	if e := tx.Reserve(owner); e != nil {
		return e
	}
	if s.fault != nil {
		if e := s.fault(faultPrefix + "-after-claim"); e != nil {
			return e
		}
	}
	for j, ch := range p.Changes {
		for _, f := range check {
			if e := f(); e != nil {
				return e
			}
		}
		if e := s.publicationInputs(intent, true); e != nil {
			return e
		}
		if e := s.applyPublicationFile(ch); e != nil {
			return e
		}
		if s.fault != nil {
			if e := s.fault(fmt.Sprintf("%s-after-file:%d", faultPrefix, j)); e != nil {
				return e
			}
		}
	}
	if e := s.publicationInputs(intent, true); e != nil {
		return e
	}
	for _, ch := range p.Changes {
		b, e := readFile(filepath.Join(s.root, ch.Path), maxState)
		actual := "absent"
		if e == nil {
			actual = hash(b)
		} else if !os.IsNotExist(e) {
			return e
		}
		if actual != ch.After {
			return fail("reconciliation_required", "selection postimage differs")
		}
	}
	if _, e := s.publicationFiles(); e != nil {
		return e
	}
	for _, f := range check {
		if e := f(); e != nil {
			return e
		}
	}
	p.Status = "committed"
	return nil
}

func validateSessionSelections(st snapshot) error {
	for _, receipt := range st.Receipts {
		if r := receipt.Result.Session; r != nil && r.Operation == "session-select" {
			found := false
			for _, i := range sessionExecutions(st) {
				if i.ID == r.InstanceID && i.Selection != nil && i.Selection.Request.OperationID == r.OperationID {
					found = true
				}
			}
			if !found {
				return fail("reconciliation_required", "orphan selection receipt")
			}
		}
	}
	for _, i := range sessionExecutions(st) {
		p := i.Selection
		if p == nil {
			continue
		}
		if !sessionEffectsProfile(st.Schema) || !i.Registered || p.Request.Schema != SessionSchema || p.Request.InstanceID != i.ID || p.Request.ProjectRef != st.Project || p.Request.CheckoutID != st.CheckoutID || !uuidPattern.MatchString(p.Request.OperationID) || !slices.Contains([]string{"prepared", "committed"}, p.Status) || len(p.Changes) == 0 || len(p.Changes) > 128 || p.LedgerBytes <= 0 {
			return fail("reconciliation_required", "invalid Session selection")
		}
		r, ok := st.Receipts[p.Request.OperationID]
		if !ok || r.Digest != requestDigest(p.Request) || r.Result.Session == nil || r.Result.Session.Operation != "session-select" || r.Result.Session.Revision != p.Request.ExpectedRevision+1 {
			return fail("reconciliation_required", "selection receipt mismatch")
		}

		stepFound := false
		for _, step := range i.Context.Binding.Steps {
			if step.ID == p.Request.StepID && sessionDependenciesSatisfied(i, step) {
				stepFound = true
			}
		}
		if !stepFound || !validLeaderText(p.Request.Reason, 2048) || p.Request.ExpectedRevision == 0 || len(p.Inputs) > 2048 || len(p.Coverage) > 512 {
			return fail("reconciliation_required", "invalid selection binding/limits")
		}
		if _, e := time.Parse(time.RFC3339Nano, p.CreatedAt); e != nil {
			return fail("reconciliation_required", "invalid selection time")
		}
		inputs := map[string]string{}
		for _, src := range p.Inputs {
			if !literalCandidatePath(src.Path) || inputs[src.Path] != "" || (!digestPattern.MatchString(src.SHA256) && src.SHA256 != "absent") {
				return fail("reconciliation_required", "invalid selection preimage")
			}
			inputs[src.Path] = src.SHA256
		}
		changed := map[string]bool{}
		for _, ch := range p.Changes {
			if changed[ch.Path] || inputs[ch.Path] != ch.Before || (!digestPattern.MatchString(ch.After) && ch.After != "absent") || (ch.After != "absent" && hash([]byte(ch.Content)) != ch.After) || (ch.After == "absent" && (ch.Content != "" || !strings.HasPrefix(ch.Path, "SDP/KanBan/backlog/"))) {
				return fail("reconciliation_required", "invalid selection effect")
			}
			changed[ch.Path] = true
		}
		pins := map[string]bool{}
		for _, src := range p.Coverage {
			if !literalCandidatePath(src.Path) || pins[src.Path] || (!digestPattern.MatchString(src.SHA256) && src.SHA256 != "absent") {
				return fail("reconciliation_required", "invalid selection coverage")
			}
			pins[src.Path] = true
		}

		post := map[string]string{}
		for _, ch := range p.Changes {
			post[ch.Path] = ch.After
		}
		required := map[string]bool{}
		for _, src := range sessionSelectionSources(i) {
			path := src.Path
			if post[path] == "absent" {
				path = "SDP/KanBan/active/" + filepath.Base(path)
			}
			required[path] = true
		}
		required["SDP/Sessions/execution/claims/"+p.Request.OperationID+".json"] = true
		if len(required) != len(p.Coverage) {
			return fail("reconciliation_required", "selection coverage set differs")
		}
		for _, src := range p.Coverage {
			if !required[src.Path] {
				return fail("reconciliation_required", "unexpected selection coverage")
			}
			if src.Path != managementLedger {
				want := inputs[src.Path]
				if h, ok := post[src.Path]; ok {
					want = h
				}
				if want == "" {
					want = "absent"
				}
				if src.SHA256 != want {
					return fail("reconciliation_required", "selection coverage digest differs")
				}
			}
		}
		if !pins[managementLedger] || !pins[i.Context.BindingPath] || !pins[i.Context.Binding.Session.Path] || !pins[i.Context.Binding.Plan.Path] || !pins[i.Context.Binding.Authority] {
			return fail("reconciliation_required", "missing selected coverage")
		}
	}
	return nil
}

// All canonical PG mutations take project coordination, canonical history, then
// operational state locks. Every lock is nonblocking; none spans a host wait.
func (s *Service) lockRecordWriter() (string, *projecthistory.Transaction, func(), error) {
	dir, e := s.publicationLockDir()
	if e != nil {
		return "", nil, nil, e
	}
	project, e := lock(dir)
	if e != nil {
		return "", nil, nil, e
	}
	tx, e := s.canonicalTransaction()
	if e != nil {
		project()
		return "", nil, nil, e
	}
	state, e := lock(s.state)
	if e != nil {
		tx.Close()
		project()
		return "", nil, nil, e
	}
	return dir, tx, func() { state(); tx.Close(); project() }, nil
}

func (s *Service) checkSelectionClaims(i SessionInstance, q SessionSelectRequest, step SessionStepBinding) error {
	// Persistent selection claims are shared between stores. They are not host
	// leases and are not released by process exit. Later disposition must release.
	claimDir := "SDP/Sessions/execution/claims"
	entries, e := os.ReadDir(filepath.Join(s.root, claimDir))
	if e != nil && !os.IsNotExist(e) {
		return e
	}
	if len(entries) > 128 {
		return fail("limit", "selection claim capacity reached")
	}
	for _, entry := range entries {
		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 {
			return fail("path", "unexpected selection claim entry")
		}
		raw, e := readFile(filepath.Join(s.root, claimDir, entry.Name()), MaxRequest)
		if e != nil {
			return e
		}
		var claim sessionScopeClaim
		if e = decode(raw, &claim); e != nil {
			return e
		}
		if claim.OperationID == q.OperationID {
			continue
		}
		if claim.Schema != "sdp-session-scope-claim/0.1" || !uuidPattern.MatchString(claim.OperationID) || len(claim.WritablePaths) == 0 || entry.Name() != claim.OperationID+".json" || !sessionIDPattern.MatchString(claim.SessionID) || !sessionPlanPattern.MatchString(claim.PlanID) || validateAssignment(Assignment{Outcome: claim.SessionID, WritablePaths: claim.WritablePaths, ExcludedPaths: claim.ProtectedPaths}) != nil {
			return fail("reconciliation_required", "invalid selection scope claim")
		}
		if claim.SessionID == i.Context.Binding.Session.ID || claim.PlanID == i.Context.Binding.Plan.ID {
			return fail("conflict", "Session/plan already selected by another claim")
		}
		for _, card := range step.Cards {
			if slices.Contains(claim.Cards, card.ID) {
				return fail("conflict", "card already selected")
			}
		}
		for _, a := range step.WritablePaths {
			for _, b := range append(append([]string{}, claim.WritablePaths...), claim.ProtectedPaths...) {
				if within(a, b) || within(b, a) {
					return fail("conflict", "selection scope overlaps another claim")
				}
			}
		}
		for _, a := range step.ProtectedPaths {
			for _, b := range claim.WritablePaths {
				if within(a, b) || within(b, a) {
					return fail("conflict", "selection reader overlaps another writer")
				}
			}
		}
	}
	return nil
}

func (s *Service) checkSelectionReferences(old string, allowed map[string]bool) error {
	count := 0
	return filepath.WalkDir(s.area, func(path string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() {
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fail("path", "symlink in Session reference coverage")
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		count++
		if count > 4096 {
			return fail("limit", "reference scan exceeds 4096 documents")
		}
		rel, _ := filepath.Rel(s.root, path)
		rel = filepath.ToSlash(rel)
		b, e := readFile(path, maxState)
		if e != nil {
			return e
		}
		target, _ := filepath.Rel(filepath.Dir(rel), old)
		escaped := (&url.URL{Path: filepath.ToSlash(target)}).EscapedPath()
		references := strings.Contains(string(b), filepath.ToSlash(target)) || strings.Contains(string(b), escaped)
		if references && !allowed[rel] {
			return fail("missing_coverage", "card reference requires explicit maintenance scope: "+rel)
		}
		if references && strings.Contains(string(rewriteCardLinks(b, rel, old, "SDP/__pg_reference_marker__")), escaped) {
			return fail("incompatible_profile", "unsupported card link syntax: "+rel)
		}
		return nil
	})
}

func (s *Service) validateSessionEffectRecords(p SessionRecordEffects) error {
	files := map[string][]byte{}
	changes := map[string]PublicationChange{}
	for _, ch := range p.Changes {
		changes[ch.Path] = ch
	}
	for _, src := range p.Inputs {
		if ch, ok := changes[src.Path]; ok {
			if ch.After != "absent" {
				files[src.Path] = []byte(ch.Content)
			}
			continue
		}
		if src.SHA256 == "absent" {
			continue
		}
		raw, e := readFile(filepath.Join(s.root, src.Path), maxState)
		if e != nil || hash(raw) != src.SHA256 {
			return fail("reconciliation_required", "selection validation input changed")
		}
		files[src.Path] = raw
	}
	if e := validatePublicationHistory(files[managementLedger], files); e != nil {
		return fail("reconciliation_required", "invalid proposed selection records: "+e.Error())
	}
	return nil
}

func (s *Service) registrationUnclaimed(i SessionInstance) error {
	dir := filepath.Join(s.root, "SDP/Sessions/execution/claims")
	entries, e := os.ReadDir(dir)
	if e != nil && !os.IsNotExist(e) {
		return e
	}
	if len(entries) > 128 {
		return fail("limit", "selection claims exceed capacity")
	}
	for _, entry := range entries {
		raw, e := readFile(filepath.Join(dir, entry.Name()), MaxRequest)
		if e != nil {
			return e
		}
		var claim sessionScopeClaim
		if e = decode(raw, &claim); e != nil {
			return e
		}
		if claim.SessionID == i.Context.Binding.Session.ID {
			return fail("conflict", "Session already selected in another instance")
		}
	}
	return nil
}

// Bounded inline Markdown link profile. Refuse unsupported reference/HTML/title
// syntax instead of silently changing destinations. Code lines with link-like
// text require explicit document maintenance before the card can move.
func rebaseMovedCardLinks(raw []byte, from, to string) ([]byte, error) {
	text := string(raw)
	if regexp.MustCompile(`(?m)^\s{0,3}\[[^\]]+\]:`).MatchString(text) || regexp.MustCompile(`(?i)(href|src)\s*=`).MatchString(text) {
		return nil, fail("incompatible_profile", "reference/HTML links require explicit card maintenance")
	}
	pattern := regexp.MustCompile(`\]\(([^)]*)\)`)
	lines := strings.Split(text, "\n")
	fenced := false
	for n, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			fenced = !fenced
		}
		if !strings.Contains(line, "](") {
			continue
		}
		if fenced || strings.Contains(line, "`") || strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
			return nil, fail("incompatible_profile", "link-like code requires explicit card maintenance")
		}
		var problem error
		lines[n] = pattern.ReplaceAllStringFunc(line, func(match string) string {
			value := pattern.FindStringSubmatch(match)[1]
			if value == "" || strings.ContainsAny(value, " \t\r\n<>") {
				problem = fail("incompatible_profile", "unsupported inline link syntax")
				return match
			}
			u, e := url.Parse(value)
			if e != nil {
				problem = e
				return match
			}
			if u.IsAbs() || u.Host != "" || strings.HasPrefix(u.Path, "/") || u.Path == "" {
				return match
			}
			target := filepath.Clean(filepath.Join(filepath.Dir(from), u.Path))
			if target == from {
				target = to
			}
			next, e := filepath.Rel(filepath.Dir(to), target)
			if e != nil {
				problem = e
				return match
			}
			next = filepath.ToSlash(next)
			if strings.HasSuffix(u.Path, "/") {
				next += "/"
			}
			if next == u.Path {
				return match
			}
			u.Path = next
			u.RawPath = ""
			return "](" + u.String() + ")"
		})
		if problem != nil {
			return nil, problem
		}
	}
	return []byte(strings.Join(lines, "\n")), nil
}
