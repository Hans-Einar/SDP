package governance

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type receipt struct {
	Digest string `json:"digest"`
	Result Result `json:"result"`
}
type snapshot struct {
	Sessions               map[string]SessionInstance       `json:"sessions,omitempty"`
	Publications           map[string]Publication           `json:"publications,omitempty"`
	Reviews                map[string]Review                `json:"reviews,omitempty"`
	Leader                 *LeaderState                     `json:"leader,omitempty"`
	LeaderReceipts         map[string]leaderReceipt         `json:"leaderReceipts,omitempty"`
	AssignmentRuns         map[string]AssignmentRun         `json:"assignmentRuns,omitempty"`
	DesignDecisions        map[string]string                `json:"designDecisions,omitempty"`
	DesignDecisionReceipts map[string]designDecisionReceipt `json:"designDecisionReceipts,omitempty"`
	Prepared               map[string]PreparedAssignment    `json:"prepared,omitempty"`
	PreparationReceipts    map[string]preparationReceipt    `json:"preparationReceipts,omitempty"`
	Schema                 string                           `json:"schema"`
	Root                   string                           `json:"root"`
	Project                string                           `json:"project"`
	CheckoutID             string                           `json:"checkoutId"`
	Previous               string                           `json:"previous"`
	Runs                   map[string]Run                   `json:"runs"`
	Receipts               map[string]receipt               `json:"receipts"`
}

func (s *Service) load() (snapshot, string, error) {
	empty := snapshot{Schema: Schema, Root: s.root, Project: s.project, Runs: map[string]Run{}, Receipts: map[string]receipt{}}
	b, e := readFile(filepath.Join(s.state, "HEAD"), 64)
	if os.IsNotExist(e) {
		return empty, "", nil
	}
	if e != nil {
		return empty, "", e
	}
	head := string(b)
	if !digestPattern.MatchString(head) {
		return empty, "", fail("reconciliation_required", "invalid state head")
	}
	b, e = readFile(filepath.Join(s.state, "generations", head+".json"), maxState)
	if e != nil {
		return empty, "", fail("reconciliation_required", "state generation unavailable: "+e.Error())
	}
	if hash(b) != head {
		return empty, "", fail("reconciliation_required", "state generation digest mismatch")
	}
	var state snapshot
	if e = decode(b, &state); e != nil {
		return empty, "", fail("reconciliation_required", e.Error())
	}
	if state.Schema != Schema && state.Schema != preparationStateSchema && state.Schema != preparationLifecycleStateSchema && state.Schema != designDecisionStateSchema && state.Schema != activationStateSchema && !attentionProfile(state.Schema) {
		return empty, "", fail("incompatible_profile", "unsupported persisted state")
	}
	if (state.Schema == Schema && (state.Prepared != nil || state.PreparationReceipts != nil)) || (state.Schema != Schema && !attentionProfile(state.Schema) && (state.Prepared == nil || state.PreparationReceipts == nil)) || ((state.Prepared == nil) != (state.PreparationReceipts == nil)) {
		return empty, "", fail("incompatible_profile", "assignment preparation requires operational profile 0.2 through 0.13")
	}
	for id, p := range state.Prepared {
		if p.ID != id || !uuidPattern.MatchString(id) || p.Revision == 0 ||
			(p.State != "" && p.State != "prepared" && p.State != "retired") ||
			(p.Revision == 1 && (p.State != "" || p.ChangeReason != "")) ||
			(p.Revision > 1 && ((state.Schema != preparationLifecycleStateSchema && state.Schema != designDecisionStateSchema && state.Schema != activationStateSchema && !attentionProfile(state.Schema)) || p.State == "" || !validPreparationReason(p.ChangeReason))) {
			return empty, "", fail("reconciliation_required", "invalid persisted preparation lifecycle")
		}
	}
	if state.Root != s.root || state.Project != s.project {
		return empty, "", fail("unauthorized", "state belongs to another checkout")
	}
	if !uuidPattern.MatchString(state.CheckoutID) || state.Runs == nil || state.Receipts == nil {
		return empty, "", fail("reconciliation_required", "invalid persisted state")
	}
	if err := validateDesignDecisionState(state); err != nil {
		return empty, "", err
	}
	if err := validateAssignmentRuns(state); err != nil {
		return empty, "", err
	}
	if err := validateAttentionState(state); err != nil {
		return empty, "", err
	}
	if err := validateAttentionContinuations(state); err != nil {
		return empty, "", err
	}
	if err := validateReportAcknowledgements(state); err != nil {
		return empty, "", err
	}
	if e := validateLeaderState(state); e != nil {
		return empty, "", e
	}
	if e := validateEvidenceState(state); e != nil {
		return empty, "", e
	}
	if e := validateReviewState(state); e != nil {
		return empty, "", e
	}
	if e := validatePublicationState(state); e != nil {
		return empty, "", e
	}
	if e := validateSessionState(state); e != nil {
		return empty, "", e
	}
	return state, head, nil
}
func syncDir(path string) error {
	f, e := os.Open(path)
	if e != nil {
		return e
	}
	defer f.Close()
	return f.Sync()
}

// Sync each newly created directory entry before acknowledging any state write.
func mkdirDurable(dir string) error {
	if e := safeAbsolute(dir); e != nil {
		return e
	}
	if st, e := os.Stat(dir); e == nil {
		if !st.IsDir() {
			return fail("path", "state path is not a directory")
		}
		return nil
	} else if !os.IsNotExist(e) {
		return e
	}
	parent := filepath.Dir(dir)
	if e := mkdirDurable(parent); e != nil {
		return e
	}
	if e := os.Mkdir(dir, 0700); e != nil && !os.IsExist(e) {
		return e
	}
	return syncDir(parent)
}

func writeNew(path string, b []byte) error {
	if e := safeAbsolute(path); e != nil {
		return e
	}
	f, e := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		return e
	}
	if _, e = f.Write(b); e != nil {
		f.Close()
		return e
	}
	if e = f.Sync(); e != nil {
		f.Close()
		return e
	}
	return f.Close()
}
func (s *Service) save(state snapshot, head string) error {
	state.Previous = head
	b, e := json.Marshal(state)
	if e != nil {
		return e
	}
	if len(b) > maxState {
		return fail("limit", "operational snapshot exceeds 16 MiB")
	}
	dir := filepath.Join(s.state, "generations")
	if e = safeAbsolute(dir); e != nil {
		return e
	}
	if e = mkdirDurable(dir); e != nil {
		return e
	}
	next := hash(b)
	path := filepath.Join(dir, next+".json")
	if e = writeNew(path, b); os.IsExist(e) {
		var old []byte
		old, e = readFile(path, maxState)
		if e == nil && hash(old) != next {
			e = fail("reconciliation_required", "generation collision")
		}
	}
	if e != nil {
		return e
	}
	if e = syncDir(dir); e != nil {
		return e
	}
	if s.fault != nil {
		if e = s.fault("before-head"); e != nil {
			return e
		}
	}
	temp, e := os.CreateTemp(s.state, "head-")
	if e != nil {
		return e
	}
	name := temp.Name()
	defer os.Remove(name)
	if _, e = temp.WriteString(next); e != nil {
		temp.Close()
		return e
	}
	if e = temp.Sync(); e != nil {
		temp.Close()
		return e
	}
	if e = temp.Close(); e != nil {
		return e
	}
	if e = os.Rename(name, filepath.Join(s.state, "HEAD")); e != nil {
		return e
	}
	if e = syncDir(s.state); e != nil {
		return e
	}
	if s.fault != nil {
		return s.fault("after-head")
	}
	return nil
}
func (s *Service) Bind(caller Caller) (Binding, error) {
	b := Binding{Schema: Schema, Operation: "bind", ProjectRef: s.project, Root: s.root}
	if !caller.controller {
		return b, fail("unauthorized", "controller binding required")
	}
	unlock, e := lock(s.state)
	if e != nil {
		return b, e
	}
	defer unlock()
	state, head, e := s.load()
	if e != nil {
		return b, e
	}
	if state.CheckoutID == "" {
		state.CheckoutID, e = identifier()
		if e != nil {
			return b, e
		}
		if e = s.save(state, head); e != nil {
			return b, e
		}
	}
	b.CheckoutID = state.CheckoutID
	return b, nil
}
