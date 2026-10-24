package governance

import (
	"sort"
	"time"
)

// ReportInbox is a read-only projection of immutable report receipts. It does
// not claim semantic consumption, acceptance or automatic continuation.
type ReportInbox struct {
	Schema      string          `json:"schema"`
	ProjectRef  string          `json:"projectRef"`
	CheckoutID  string          `json:"checkoutId"`
	Generation  string          `json:"generation"`
	InspectedAt string          `json:"inspectedAt"`
	Scope       string          `json:"scope"`
	Consumption string          `json:"consumption"`
	Reports     []ReportMessage `json:"reports"`
}

type ReportMessage struct {
	Acknowledgement   *ReportAcknowledgement `json:"acknowledgement,omitempty"`
	AcknowledgementID string                 `json:"acknowledgementId,omitempty"`
	ID                string                 `json:"id"`
	RunID             string                 `json:"runId"`
	ReportedRevision  uint64                 `json:"reportedRevision"`
	InputRevision     string                 `json:"inputRevision"`
	Selection         Selection              `json:"selection"`
	Outcome           string                 `json:"outcome"`
	Text              string                 `json:"text"`
	CurrentRevision   uint64                 `json:"currentRevision"`
	CurrentState      string                 `json:"currentState"`
	Coverage          string                 `json:"coverage"`
	Conditions        []string               `json:"conditions"`
}

func (s *Service) ReportInbox(runID string) (ReportInbox, error) {
	out := ReportInbox{Schema: "sdp-governance-inbox/1", ProjectRef: s.project,
		InspectedAt: time.Now().UTC().Format(time.RFC3339Nano), Scope: "report-receipts-in-selected-checkout",
		Consumption: "not-tracked", Reports: []ReportMessage{}}
	if runID != "" && !uuidPattern.MatchString(runID) {
		return out, fail("invalid_request", "invalid runId")
	}
	state, head, err := s.load()
	if err != nil {
		return out, err
	}
	if runID != "" {
		if _, ok := state.Runs[runID]; !ok {
			return out, fail("missing_coverage", "run not found")
		}
	}
	out.CheckoutID, out.Generation = state.CheckoutID, head
	acknowledgements := map[string]Result{}
	for _, receipt := range state.Receipts {
		if a := receipt.Result.ReportAcknowledgement; a != nil {
			acknowledgements[a.Request.ReportID] = receipt.Result
		}
	}
	if len(acknowledgements) > 0 {
		out.Schema = "sdp-governance-inbox/2"
		out.Consumption = "acknowledgement-recorded-separately"
	}
	views := map[string]ReportMessage{}
	for _, receipt := range state.Receipts {
		result := receipt.Result
		r := result.Run
		if result.Operation != "transition" || r.State != "verifying" || (runID != "" && r.ID != runID) {
			continue
		}
		v, ok := views[r.ID]
		if !ok {
			current, exists := state.Runs[r.ID]
			if !exists {
				return out, fail("reconciliation_required", "report has no current run")
			}
			v = ReportMessage{RunID: r.ID, CurrentRevision: current.Revision, CurrentState: current.State, Conditions: []string{}}
			c, e := s.runContext(current)
			switch {
			case e != nil:
				v.Coverage = "unavailable"
				v.Conditions = append(v.Conditions, e.Error())
			case c.Revision != current.InputRevision || len(c.Conditions) > 0:
				v.Coverage = "stale"
				v.Conditions = append(v.Conditions, c.Conditions...)
				if c.Revision != current.InputRevision {
					v.Conditions = append(v.Conditions, "governing sources differ from pinned input revision")
				}
			default:
				v.Coverage = "known"
			}
			views[r.ID] = v
		}
		v.ID, v.ReportedRevision, v.InputRevision = result.OperationID, r.Revision, r.InputRevision
		v.Selection, v.Outcome, v.Text = r.Selection, r.Assignment.Outcome, r.Reason
		ack := acknowledgements[v.ID]
		v.Acknowledgement, v.AcknowledgementID = ack.ReportAcknowledgement, ack.OperationID
		out.Reports = append(out.Reports, v)
	}
	sort.Slice(out.Reports, func(i, j int) bool {
		a, b := out.Reports[i], out.Reports[j]
		if a.RunID != b.RunID {
			return a.RunID < b.RunID
		}
		if a.ReportedRevision != b.ReportedRevision {
			return a.ReportedRevision < b.ReportedRevision
		}
		return a.ID < b.ID
	})
	// Coverage is a separate filesystem check, not an atomic filesystem snapshot.
	_, latest, err := s.load()
	if err != nil {
		return out, err
	}
	if latest != head {
		return out, fail("stale_revision", "inbox changed during inspection")
	}
	return out, nil
}
