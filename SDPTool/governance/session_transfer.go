package governance

import (
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

const sessionTransferStateSchema = "sdp-project-governance-state/0.21"

func sessionLifecycleProfile(schema string) bool {
	return schema == sessionLifecycleStateSchema || schema == sessionTransferStateSchema
}
func terminalSessionAdvance(action string) bool { return action == "complete" || action == "transfer" }

type SessionSuccessor struct {
	BindingPath   string `json:"bindingPath"`
	InputRevision string `json:"inputRevision"`
}
type SessionSuccessorRequest struct {
	InstanceID  string `json:"instanceId"`
	BindingPath string `json:"bindingPath"`
}

// Remaining scope comes from the last committed disposition of each card, never
// from a caller's assertion that an unfinished step has been completed.
func remainingSessionCards(i SessionInstance) (map[string]SessionCardDisposition, error) {
	for _, step := range i.Context.Binding.Steps {
		if !sessionStepDone(i, step.ID) {
			return nil, fail("missing_coverage", "all predecessor steps must be disposed before transfer")
		}
	}
	cards := map[string]SessionCardDisposition{}
	record := func(d *SessionDisposition) {
		if d == nil || d.Effects.Status != "committed" {
			return
		}
		for _, c := range d.Request.Cards {
			if c.Outcome == "completed" {
				delete(cards, c.ID)
			} else {
				cards[c.ID] = c
			}
		}
	}
	for _, p := range i.Past {
		record(p.Execution.Disposition)
	}
	record(i.Disposition)
	if len(cards) == 0 {
		return nil, fail("missing_coverage", "transfer requires explicitly remaining contribution scope")
	}
	return cards, nil
}

// Validate the successor against preimages. The predecessor still owns its
// cards at this point; only this explicit transfer may change their PlanId.
func successorContext(i SessionInstance, path string, files map[string][]byte) (SessionContext, error) {
	out := SessionContext{Schema: SessionSchema, Operation: "session-successor", ProjectRef: i.Context.ProjectRef, CheckoutID: i.Context.CheckoutID, BindingPath: path, Sources: []Source{}}
	remaining, e := remainingSessionCards(i)
	if e != nil {
		return out, e
	}
	raw, ok := files[path]
	if !ok || len(raw) > 64<<10 {
		return out, fail("missing_coverage", "bounded successor binding required")
	}
	if e = decode(raw, &out.Binding); e != nil {
		return out, e
	}
	b := out.Binding
	if e = validateSessionBinding(b); e != nil {
		return out, e
	}
	if path != "SDP/Sessions/execution/"+b.Session.ID+".json" || b.Session.ID == i.Context.Binding.Session.ID || b.Session.Path == i.Context.Binding.Session.Path || b.Plan.ID == i.Context.Binding.Plan.ID || b.Plan.Path == i.Context.Binding.Plan.Path {
		return out, fail("conflict", "distinct canonical successor Session and plan required")
	}
	sm, e := metadata(files[b.Session.Path])
	if e != nil {
		return out, e
	}
	pm, e := metadata(files[b.Plan.Path])
	if e != nil {
		return out, e
	}
	if sm["Session reference"] != b.Session.ID || sm["Status"] != "active" || pm["id"] != b.Plan.ID || pm["state"] != "planned" || !strings.Contains(string(files[b.Session.Path]), filepath.Base(path)) || strings.TrimSpace(string(files[b.Authority])) == "" {
		return out, fail("missing_coverage", "linked manual active successor Session, planned plan and authority required")
	}
	paths := map[string]bool{path: true, b.Session.Path: true, b.Plan.Path: true, b.Authority: true}
	cardPaths := map[string]string{}
	for p, raw := range files {
		if strings.HasPrefix(p, "SDP/KanBan/") {
			m, err := metadata(raw)
			if err == nil && remaining[m["id"]].ID != "" {
				if cardPaths[m["id"]] != "" {
					return out, fail("conflict", "duplicate remaining card")
				}
				cardPaths[m["id"]] = p
			}
		}
	}
	seen := map[string]bool{}
	for _, step := range b.Steps {
		if !sessionAnchor(files[b.Session.Path], step.ID) || !sessionAnchor(files[b.Plan.Path], step.Milestone) {
			return out, fail("missing_coverage", "successor anchors absent")
		}
		for _, c := range step.Cards {
			p := cardPaths[c.ID]
			if p == "" {
				return out, fail("missing_coverage", "successor must cover exactly remaining cards")
			}
			m, _ := metadata(files[p])
			if m["PlanId"] != i.Context.Binding.Plan.ID || m["CardState"] != "in-progress" {
				return out, fail("conflict", "remaining card ownership differs")
			}
			paths[p] = true
			seen[c.ID] = true
		}
	}
	if len(seen) != len(remaining) {
		return out, fail("missing_coverage", "successor omits remaining cards")
	}
	names := []string{}
	for p := range paths {
		names = append(names, p)
	}
	sort.Strings(names)
	for _, p := range names {
		raw, ok := files[p]
		if !ok {
			return out, fail("missing_coverage", "successor source absent")
		}
		out.Sources = append(out.Sources, Source{p, hash(raw)})
	}
	out.BindingDigest = hash(raw)
	out.InputRevision = sessionInputRevision(out.Sources)
	return out, nil
}
func (s *Service) readSuccessorContext(i SessionInstance, path string, files map[string][]byte) (SessionContext, error) {
	if !literalCandidatePath(path) || !strings.HasPrefix(path, "SDP/Sessions/execution/") || !strings.HasSuffix(path, ".json") {
		return SessionContext{}, fail("path", "canonical successor binding required")
	}
	raw, e := readFile(filepath.Join(s.root, path), 64<<10)
	if e != nil {
		return SessionContext{}, e
	}
	var b SessionBinding
	if e = decode(raw, &b); e != nil {
		return SessionContext{}, e
	}
	if e = validateSessionBinding(b); e != nil {
		return SessionContext{}, e
	}
	files[path] = raw
	for _, p := range []string{b.Session.Path, b.Plan.Path, b.Authority} {
		raw, e = readFile(filepath.Join(s.root, p), MaxRequest)
		if e != nil {
			return SessionContext{}, e
		}
		files[p] = raw
	}
	return successorContext(i, path, files)
}
func (s *Service) SessionSuccessorContext(q SessionSuccessorRequest) (SessionContext, error) {
	st, _, e := s.load()
	if e != nil {
		return SessionContext{}, e
	}
	i, ok := st.Sessions[q.InstanceID]
	if !ok || !i.Registered || i.Advance != nil {
		return SessionContext{}, fail("missing_coverage", "registered nonterminal predecessor required")
	}
	if e = s.sessionCurrent(i); e != nil {
		return SessionContext{}, e
	}
	files, e := s.publicationFiles()
	if e != nil {
		return SessionContext{}, e
	}
	ctx, e := s.readSuccessorContext(i, q.BindingPath, files)
	if e != nil {
		return ctx, e
	}
	e = s.transferClaims(st, i, SessionAdvance{Successor: &ctx})
	return ctx, e
}
func validateSuccessorRecord(i SessionInstance, a SessionAdvance) error {
	if a.Request.Action != "transfer" {
		if a.Successor != nil {
			return fail("reconciliation_required", "unexpected successor context")
		}
		return nil
	}
	if a.Successor == nil || a.Request.Successor == nil {
		return fail("missing_coverage", "successor context absent")
	}
	files := map[string][]byte{}
	for p, b := range a.Before {
		files[p] = []byte(b)
	}
	want, e := successorContext(i, a.Request.Successor.BindingPath, files)
	if e != nil {
		return e
	}
	if !reflect.DeepEqual(want, *a.Successor) || want.InputRevision != a.Request.Successor.InputRevision {
		return fail("reconciliation_required", "successor context differs from preimages")
	}
	return nil
}
func (s *Service) transferClaims(st snapshot, i SessionInstance, a SessionAdvance) error {
	if a.Successor == nil {
		return nil
	}
	b := a.Successor.Binding
	for _, other := range st.Sessions {
		if other.ID == i.ID {
			continue
		}
		if other.Advance == nil || other.Advance.Effects.Status != "committed" || !terminalSessionAdvance(other.Advance.Request.Action) {
			remaining, _ := remainingSessionCards(i)
			for _, step := range other.Context.Binding.Steps {
				for _, c := range step.Cards {
					if remaining[c.ID].ID != "" {
						return fail("conflict", "remaining card bound by another live Session")
					}
				}
			}
		}
		if other.Context.Binding.Session.ID == b.Session.ID || other.Context.Binding.Session.Path == b.Session.Path || other.Context.Binding.Plan.ID == b.Plan.ID {
			return fail("conflict", "successor already adopted or plan owned by another instance")
		}
	}
	// Check every successor step, including scopes which differ from the predecessor.
	for _, step := range b.Steps {
		if e := s.checkSelectionClaims(SessionInstance{}, SessionSelectRequest{}, step); e != nil {
			return e
		}
		for _, r := range st.Runs {
			if r.State != "completed" && r.State != "canceled" && successorScopeConflict(step, r.Assignment) {
				return fail("conflict", "successor overlaps active run")
			}
		}
		for _, p := range st.Prepared {
			if p.State != "retired" && successorScopeConflict(step, p.Assignment) {
				return fail("conflict", "successor overlaps prepared assignment")
			}
		}
	}
	return nil
}
func successorScopeConflict(step SessionStepBinding, a Assignment) bool {
	for _, p := range a.WritablePaths {
		if !literalCandidatePath(p) {
			return true
		}
		for _, q := range append(append([]string{}, step.WritablePaths...), step.ProtectedPaths...) {
			if overlaps(p, q) {
				return true
			}
		}
	}
	for _, p := range a.ExcludedPaths {
		if !literalCandidatePath(p) {
			return true
		}
		for _, q := range step.WritablePaths {
			if overlaps(p, q) {
				return true
			}
		}
	}
	return false
}
func mergeSessionSources(a, b []Source) []Source {
	m := map[string]string{}
	for _, src := range a {
		m[src.Path] = src.SHA256
	}
	for _, src := range b {
		m[src.Path] = src.SHA256
	}
	names := []string{}
	for p := range m {
		names = append(names, p)
	}
	sort.Strings(names)
	out := []Source{}
	for _, p := range names {
		out = append(out, Source{p, m[p]})
	}
	return out
}
