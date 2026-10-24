package governance

import (
	"encoding/base64"
	"encoding/json"
	"sort"
	"time"
)

const CoordinationSchema = "sdp-governance-coordination/1"

type CoordinationRequest struct {
	RunID  string `json:"runId,omitempty"`
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}
type ReviewNotice struct {
	PublicationApplicable bool     `json:"publicationApplicable"`
	ID                    string   `json:"id"`
	Revision              uint64   `json:"revision"`
	State                 string   `json:"state"`
	Candidate             string   `json:"candidate"`
	ReportID              string   `json:"reportId"`
	Verdict               string   `json:"verdict"`
	Summary               string   `json:"summary"`
	FindingCount          int      `json:"findingCount"`
	Applicable            bool     `json:"applicable"`
	Qualified             bool     `json:"qualified"`
	ContextSeparation     string   `json:"contextSeparation"`
	Conditions            []string `json:"conditions"`
	PublicationID         string   `json:"publicationId"`
	PublicationStatus     string   `json:"publicationStatus"`
}
type CoordinationItem struct {
	Review    *ReviewNotice  `json:"review,omitempty"`
	ID        string         `json:"id"`
	RunID     string         `json:"runId"`
	Kind      string         `json:"kind"`
	Priority  int            `json:"priority"`
	Work      *WorkView      `json:"work,omitempty"`
	Report    *ReportMessage `json:"report,omitempty"`
	Attention *AttentionView `json:"attention,omitempty"`
}
type CoordinationView struct {
	Schema      string             `json:"schema"`
	ProjectRef  string             `json:"projectRef"`
	CheckoutID  string             `json:"checkoutId"`
	Generation  string             `json:"generation"`
	Revision    string             `json:"revision"`
	InspectedAt string             `json:"inspectedAt"`
	Scope       string             `json:"scope"`
	Items       []CoordinationItem `json:"items"`
	Total       int                `json:"total"`
	NextCursor  string             `json:"nextCursor"`
}
type coordinationCursor struct {
	Revision string `json:"revision"`
	RunID    string `json:"runId"`
	Offset   int    `json:"offset"`
}

// Coordination reads a bounded projection. Generation pins operational facts;
// source coverage is checked separately and no host liveness is inferred.
func (s *Service) Coordination(q CoordinationRequest) (CoordinationView, error) {
	var out CoordinationView
	if q.Limit == 0 {
		q.Limit = 20
	}
	if q.Limit < 1 || q.Limit > 50 || len(q.Cursor) > 512 || (q.RunID != "" && !uuidPattern.MatchString(q.RunID)) {
		return out, fail("invalid_request", "limit must be 1..50 with a valid optional run and cursor")
	}
	work, e := s.Overview()
	if e != nil {
		return out, e
	}
	reports, e := s.ReportInbox(q.RunID)
	if e != nil {
		return out, e
	}
	attention, e := s.AttentionInbox(q.RunID)
	if e != nil {
		return out, e
	}
	state, head, e := s.load()
	if e != nil {
		return out, e
	}
	if work.Generation != head || reports.Generation != head || attention.Generation != head {
		return out, fail("stale_revision", "coordination changed during inspection")
	}
	items := []CoordinationItem{}
	for _, v := range work.Runs {
		if q.RunID == "" || v.RunID == q.RunID {
			items = append(items, CoordinationItem{ID: "work:" + v.RunID, RunID: v.RunID, Kind: "work", Priority: 4, Work: &v})
		}
	}
	for _, v := range reports.Reports {
		priority := 2
		if v.Acknowledgement != nil {
			priority = 5
		}
		items = append(items, CoordinationItem{ID: "report:" + v.ID, RunID: v.RunID, Kind: "report", Priority: priority, Report: &v})
	}
	for _, v := range attention.Items {
		priority := 3
		if v.Status == "open" {
			priority = 1
			if v.Question.Blocking {
				priority = 0
			}
		}
		items = append(items, CoordinationItem{ID: "attention:" + v.Question.ID, RunID: v.RunID, Kind: "attention", Priority: priority, Attention: &v})
	}
	for _, d := range state.Reviews {
		if q.RunID != "" && q.RunID != d.RunID {
			continue
		}
		v, err := s.ReviewRead(d.ID)
		if err != nil {
			return out, err
		}
		n := ReviewNotice{ID: d.ID, Revision: d.Revision, State: d.State, Candidate: d.Spec.Candidate.Digest, ReportID: d.Spec.ReportID, Applicable: v.Applicable, Qualified: v.Qualified, ContextSeparation: v.ContextSeparation, Conditions: v.Conditions}
		if d.Outcome != nil {
			n.Verdict = d.Outcome.Verdict
			n.Summary = d.Outcome.Summary
			n.FindingCount = len(d.Outcome.Findings)
		}
		for id, p := range state.Publications {
			if p.Request.ReviewID == d.ID && p.Status != "abandoned" {
				n.PublicationID = id
				n.PublicationStatus = p.Status
				n.PublicationApplicable = s.publicationView(state, p).Applicable
			}
		}
		items = append(items, CoordinationItem{ID: "review:" + d.ID, RunID: d.RunID, Kind: "review", Priority: 2, Review: &n})
	}
	_, latest, e := s.load()
	if e != nil {
		return out, e
	}
	if latest != head {
		return out, fail("stale_revision", "reviews changed during coordination inspection")
	}
	sort.Slice(items, func(i, j int) bool {
		a, b := items[i], items[j]
		if a.Priority != b.Priority {
			return a.Priority < b.Priority
		}
		return a.ID < b.ID
	})
	revision := requestDigest(struct {
		Project, Checkout, Generation, Run string
		Items                              []CoordinationItem
	}{s.project, work.CheckoutID, head, q.RunID, items})
	offset := 0
	if q.Cursor != "" {
		b, e := base64.RawURLEncoding.DecodeString(q.Cursor)
		if e != nil {
			return out, fail("invalid_request", "invalid coordination cursor")
		}
		var c coordinationCursor
		if e = Decode(b, &c); e != nil {
			return out, e
		}
		if c.Revision != revision || c.RunID != q.RunID {
			return out, fail("stale_revision", "coordination cursor no longer applies; inspect from the beginning")
		}
		if c.Offset < 1 || c.Offset >= len(items) {
			return out, fail("invalid_request", "invalid coordination offset")
		}
		offset = c.Offset
	}
	out = CoordinationView{Schema: CoordinationSchema, ProjectRef: s.project, CheckoutID: work.CheckoutID, Generation: head, Revision: revision, InspectedAt: time.Now().UTC().Format(time.RFC3339Nano), Scope: "registered-work-reports-attention-in-selected-checkout", Items: []CoordinationItem{}, Total: len(items)}
	if len(state.Reviews) > 0 {
		out.Schema = "sdp-governance-coordination/2"
		out.Scope = "registered-work-reports-attention-reviews-in-selected-checkout"
	}
	size := 0
	end := offset
	for end < len(items) && end-offset < q.Limit {
		b, e := json.Marshal(items[end])
		if e != nil {
			return CoordinationView{}, e
		}
		if size+len(b) > 200<<10 {
			if end == offset {
				return CoordinationView{}, fail("limit", "single coordination item exceeds page capacity; use the existing per-run views")
			}
			break
		}
		size += len(b)
		out.Items = append(out.Items, items[end])
		end++
	}
	if end < len(items) {
		b, _ := json.Marshal(coordinationCursor{Revision: revision, RunID: q.RunID, Offset: end})
		out.NextCursor = base64.RawURLEncoding.EncodeToString(b)
	}
	return out, nil
}
