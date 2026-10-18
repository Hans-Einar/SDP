package blueprints

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Hans-Einar/SDP/SDPTool/blueprintstate"
	"github.com/Hans-Einar/SDP/SDPTool/model"
	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"path/filepath"
	"reflect"
	"sort"
)

var errAssignmentStale = errors.New("assignment source references are stale")

type AssignmentResult struct {
	Schema          string               `json:"schema"`
	Operation       string               `json:"operation"`
	HistoryRevision string               `json:"historyRevision"`
	Appended        bool                 `json:"appended"`
	Assignment      blueprintstate.State `json:"assignment"`
}
type AssignmentView struct {
	OpenRevision string `json:"openRevision,omitempty"`
	blueprintstate.State
	Validation     string `json:"validation"`
	Freshness      string `json:"sourceFreshness"`
	Readiness      string `json:"readinessStatus"`
	EvidenceStatus string `json:"evidenceStatus"`
	Diagnostic     string `json:"diagnostic,omitempty"`
	OpenPath       string `json:"openPath,omitempty"`
}

func assignmentPath(area, rel string) (string, error) {
	if !blueprintstate.Relative(rel) {
		return "", fmt.Errorf("invalid project-relative path")
	}
	p := filepath.Join(area, filepath.FromSlash(rel))
	return p, noSymlinks(p)
}
func boundDocument(area string, b blueprintstate.Binding) (Document, error) {
	d, _, e := boundCapture(area, b, &budget{})
	return d, e
}
func boundCapture(area string, b blueprintstate.Binding, reads *budget) (Document, *documents.Bundle, error) {
	p, err := assignmentPath(area, b.Bundle)
	if err != nil {
		return Document{}, nil, err
	}
	d, bundle, err := verify(p, reads)
	if err != nil {
		return d, nil, err
	}
	retained := model.Digest(model.Files(bundle.Files))
	if d.Revision != b.Revision || retained != b.Retained || d.TaskID != b.Task || d.Analysis.System != b.System || filepath.Base(p) != retained || filepath.Base(filepath.Dir(p)) != key(d) {
		return d, nil, fmt.Errorf("assignment bundle identity mismatch")
	}
	return d, bundle, nil
}
func freshSources(area string, b blueprintstate.Binding, d Document) error {
	root := area
	if b.ModelArea != "." {
		var err error
		root, err = assignmentPath(area, b.ModelArea)
		if err != nil {
			return err
		}
	}
	now, err := model.Snapshot(root, b.From)
	if err != nil {
		return err
	}
	target, err := model.Snapshot(root, b.To)
	if err != nil {
		return err
	}
	if capture(now) != d.Now || capture(target) != d.Target {
		return errAssignmentStale
	}
	return nil
}
func planPresent(raw []byte, id string) error {
	found := false
	for _, line := range bytes.Split(raw, []byte{'\n'}) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var e struct {
			Subject string `json:"subjectId"`
			Type    string `json:"eventType"`
			Payload struct {
				Kind string `json:"kind"`
				To   string `json:"to"`
			} `json:"payload"`
		}
		if err := json.Unmarshal(line, &e); err != nil {
			return err
		}
		if e.Subject == id && e.Type != "" && e.Payload.Kind == "Plan" {
			found = e.Payload.To == "active"
		}
	}
	if !found {
		return fmt.Errorf("assignment requires an active canonical plan: %s", id)
	}
	return nil
}

// ImplementationReceipt is attributed scoped evidence, not tool-executed proof.
type ImplementationReceipt struct {
	Schema         string   `json:"schema"`
	Assignment     string   `json:"assignmentId"`
	Retained       string   `json:"retained"`
	EvidenceDigest string   `json:"evidenceDigest"`
	CodeDigest     string   `json:"codeDigest"`
	Checks         []string `json:"checks"`
}

func assignmentProof(area string, b blueprintstate.Binding, evidence, trace, id string) (*blueprintstate.Proof, error) {
	return assignmentProofBudget(area, b, evidence, trace, id, &budget{})
}
func assignmentProofBudget(area string, b blueprintstate.Binding, evidence, trace, id string, reads *budget) (*blueprintstate.Proof, error) {
	ep, err := assignmentPath(area, evidence)
	if err != nil {
		return nil, err
	}
	bp, err := assignmentPath(area, b.Bundle)
	if err != nil {
		return nil, err
	}
	a, err := assessBudget(bp, ep, nil, reads)
	if err != nil {
		return nil, err
	}
	if a.Status != "ready" {
		return nil, fmt.Errorf("assignment evidence blocked: %v", a.Blockers)
	}
	proof := &blueprintstate.Proof{Assessment: a.Identity, EvidenceDigest: a.EvidenceDigest}
	if trace == "" {
		return proof, nil
	}
	raw, err := reads.read(ep, 1<<20)
	if err != nil {
		return nil, err
	}
	if documents.Hash(raw) != a.EvidenceDigest {
		return nil, fmt.Errorf("evidence changed")
	}
	var evidenceRecord Evidence
	if err = decodeRecord(raw, &evidenceRecord); err != nil {
		return nil, err
	}
	code, ok := evidenceRecord.Code["TARGET"]
	if !ok || len(code.Files) == 0 {
		return nil, fmt.Errorf("submission requires scoped target code")
	}
	for _, u := range a.Unknowns {
		if u.ID == "code/TARGET" {
			return nil, fmt.Errorf("stale target code")
		}
	}
	checks := []string{}
	for _, c := range a.Checks {
		if c.Side == "TARGET" {
			if c.Effective != "pass" {
				return nil, fmt.Errorf("submission requires passing target checks")
			}
			checks = append(checks, c.ID)
		}
	}
	if len(checks) == 0 {
		return nil, fmt.Errorf("submission requires target verification")
	}
	for _, r := range evidenceRecord.Receipts {
		if r.Side == "TARGET" && r.Traceability != trace {
			return nil, fmt.Errorf("receipt traceability mismatch")
		}
	}
	sort.Strings(checks)
	traceBytes, err := reads.read(filepath.Join(area, "Traceability", "Ledger.ndjson"), projecthistory.MaxBytes)
	if err != nil {
		return nil, err
	}
	if err = projecthistory.Validate(traceBytes); err != nil {
		return nil, err
	}
	for _, line := range bytes.Split(traceBytes, []byte{'\n'}) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var e struct {
			ID      string          `json:"eventId"`
			Type    string          `json:"eventType"`
			Subject string          `json:"subjectId"`
			Payload json.RawMessage `json:"payload"`
		}
		if err = json.Unmarshal(line, &e); err != nil {
			return nil, err
		}
		if e.ID != trace {
			continue
		}
		var receipt ImplementationReceipt
		if err = blueprintstate.Strict(e.Payload, &receipt); err != nil {
			return nil, err
		}
		if e.Type != "x-verification:recorded" || e.Subject != id || receipt.Schema != "sdp-blueprint-implementation/1" || receipt.Assignment != id || receipt.Retained != b.Retained || receipt.EvidenceDigest != a.EvidenceDigest || receipt.CodeDigest != code.Digest {
			return nil, fmt.Errorf("traceability binding mismatch")
		}
		sort.Strings(receipt.Checks)
		if !reflect.DeepEqual(receipt.Checks, checks) {
			return nil, fmt.Errorf("traceability checks mismatch")
		}
		proof.CodeDigest = code.Digest
		proof.TraceEvent = trace
		proof.TraceDigest = documents.Hash(line)
		return proof, nil
	}
	return nil, fmt.Errorf("missing implementation trace event")
}

// Principal is supplied by a trusted local adapter, never by request JSON.
func ApplyAssignment(area string, principal blueprintstate.Principal, r blueprintstate.Request) (AssignmentResult, error) {
	result := AssignmentResult{Schema: blueprintstate.Schema, Operation: "assignment"}
	if err := blueprintstate.ValidateRequest(r); err != nil {
		return result, err
	}
	if principal.Actor == "" {
		return result, fmt.Errorf("missing trusted principal")
	}
	historyPath := filepath.Join(area, "ProjectManagement", "Ledger.ndjson")
	history, err := projecthistory.Read(historyPath)
	if err != nil {
		return result, err
	}
	states, events, err := blueprintstate.Replay(history.Bytes)
	if err != nil {
		return result, err
	}
	if prior, ok := events[r.EventID]; ok {
		if prior.Actor != principal.Actor || prior.Payload.Authority != principal.Authority || !reflect.DeepEqual(prior.Payload.Request, r) {
			return result, projecthistory.ErrConflict
		}
		result.HistoryRevision = history.Revision
		result.Assignment = states[r.ID]
		return result, nil
	}
	old := states[r.ID]
	e := blueprintstate.Event{SchemaVersion: "1.0", EventID: r.EventID, EventType: "x-blueprint:" + r.Action, OccurredAt: r.OccurredAt, Actor: principal.Actor, SubjectID: r.ID, Payload: blueprintstate.Payload{Schema: blueprintstate.Schema, Authority: principal.Authority, Request: r}}
	b := old.Binding
	if r.Binding != nil {
		b = *r.Binding
	}
	if r.Action == "create" || r.Action == "adopt-readiness" || r.Action == "assign" {
		if err = planPresent(history.Bytes, b.Plan); err != nil {
			return result, err
		}
	}
	if r.Action == "supersede" {
		successor, ok := states[r.Successor]
		if !ok || successor.Binding.System != b.System {
			return result, fmt.Errorf("invalid successor assignment")
		}
	}
	gated := r.Action != "cancel" && r.Action != "supersede" && r.Action != "hold" && r.Action != "reject"
	check := func() error {
		if !gated {
			return nil
		}
		d, err := boundDocument(area, b)
		if err != nil {
			return err
		}
		if err = freshSources(area, b, d); err != nil {
			return err
		}
		switch r.Action {
		case "adopt-readiness", "submit":
			p, err := assignmentProof(area, b, r.Evidence, r.Trace, r.ID)
			if err != nil {
				return err
			}
			e.Payload.Proof = p
		case "assign", "start":
			p, err := assignmentProof(area, b, old.ReadyEvidence, "", r.ID)
			if err != nil {
				return err
			}
			if !reflect.DeepEqual(p, old.Ready) {
				return fmt.Errorf("adopted readiness changed")
			}
		case "accept-review", "complete":
			if old.Submission == nil {
				return fmt.Errorf("missing submitted evidence")
			}
			p, err := assignmentProof(area, b, old.SubmissionEvidence, old.Submission.TraceEvent, r.ID)
			if err != nil {
				return err
			}
			if !reflect.DeepEqual(p, old.Submission) {
				return fmt.Errorf("submitted evidence changed")
			}
			e.Payload.Proof = p
		case "resume":
			if old.Paused == "review" {
				if old.Submission == nil {
					return fmt.Errorf("missing submission")
				}
				p, err := assignmentProof(area, b, old.SubmissionEvidence, old.Submission.TraceEvent, r.ID)
				if err != nil {
					return err
				}
				if !reflect.DeepEqual(p, old.Submission) {
					return fmt.Errorf("submission changed")
				}
			} else if old.Paused != "draft" && old.Paused != "in-progress" {
				p, err := assignmentProof(area, b, old.ReadyEvidence, "", r.ID)
				if err != nil {
					return err
				}
				if !reflect.DeepEqual(p, old.Ready) {
					return fmt.Errorf("readiness changed")
				}
			}
		}
		return nil
	}
	if err = check(); err != nil {
		return result, err
	}
	n, err := blueprintstate.Next(old, e)
	if err != nil {
		return result, err
	}
	firstProof := e.Payload.Proof
	if err = check(); err != nil {
		return result, err
	}
	if !reflect.DeepEqual(firstProof, e.Payload.Proof) {
		return result, fmt.Errorf("evidence changed during transition")
	}
	raw, _ := json.Marshal(e)
	candidate := append(append(append([]byte{}, history.Bytes...), raw...), '\n')
	if _, _, err = blueprintstate.Replay(candidate); err != nil {
		return result, err
	}
	appended, err := projecthistory.Append(historyPath, history.Revision, raw)
	result.HistoryRevision = appended.Revision
	result.Appended = appended.Appended
	result.Assignment = n
	return result, err
}
func AssignmentViews(area string) ([]AssignmentView, error) {
	h, err := projecthistory.Read(filepath.Join(area, "ProjectManagement", "Ledger.ndjson"))
	if err != nil {
		return nil, err
	}
	states, _, err := blueprintstate.Replay(h.Bytes)
	if err != nil {
		return nil, err
	}
	if len(states) > 256 {
		return nil, fmt.Errorf("assignment projection limit")
	}
	ids := []string{}
	for id := range states {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := []AssignmentView{}
	reads := &budget{}
	freshness := map[blueprintstate.Binding]error{}
	sourceComparisons := 0
	for _, id := range ids {
		s := states[id]
		v := AssignmentView{State: s, Validation: "unknown", Freshness: "unknown", Readiness: "not-adopted", EvidenceStatus: "not-submitted"}
		if s.Ready != nil {
			v.Readiness = "changed-or-unavailable"
		}
		if s.Submission != nil {
			v.EvidenceStatus = "changed-or-unavailable"
		}
		d, bundle, err := boundCapture(area, s.Binding, reads)
		if err != nil {
			v.Diagnostic = err.Error()
		} else {
			v.Validation = "validated"
			v.OpenPath = filepath.Join(area, s.Binding.Bundle, "index.md")
			v.OpenRevision = documents.Hash(bundle.Files["index.md"])
			if previous, ok := freshness[s.Binding]; ok {
				err = previous
			} else if sourceComparisons < 2 {
				err = freshSources(area, s.Binding, d)
				freshness[s.Binding] = err
				sourceComparisons++
			} else {
				err = fmt.Errorf("live-source comparison budget exhausted")
			}
			if err != nil {
				v.Freshness = "unknown"
				if errors.Is(err, errAssignmentStale) {
					v.Freshness = "stale"
				}
				v.Diagnostic = err.Error()
			} else {
				v.Freshness = "current"
			}
			if s.Ready != nil {
				v.Readiness = "changed-or-unavailable"
				p, e := assignmentProofBudget(area, s.Binding, s.ReadyEvidence, "", id, reads)
				if e == nil && reflect.DeepEqual(p, s.Ready) {
					v.Readiness = "ready"
				} else {
					v.Diagnostic += fmt.Sprintf("; readiness: %v", e)
					if e == nil {
						v.Diagnostic += " pinned assessment changed"
					}
				}
			}
			if s.Submission != nil {
				v.EvidenceStatus = "changed-or-unavailable"
				p, e := assignmentProofBudget(area, s.Binding, s.SubmissionEvidence, s.Submission.TraceEvent, id, reads)
				if e == nil && reflect.DeepEqual(p, s.Submission) {
					v.EvidenceStatus = "verified-capture"
				} else {
					v.Diagnostic += fmt.Sprintf("; submission: %v", e)
					if e == nil {
						v.Diagnostic += " pinned assessment changed"
					}
				}
			}
		}
		out = append(out, v)
	}
	return out, nil
}

func LoadAssignmentRequest(path string) (blueprintstate.Request, error) {
	var r blueprintstate.Request
	raw, err := (&budget{}).read(path, 1<<20)
	if err != nil {
		return r, err
	}
	if err = blueprintstate.Strict(raw, &r); err != nil {
		return r, err
	}
	return r, blueprintstate.ValidateRequest(r)
}
