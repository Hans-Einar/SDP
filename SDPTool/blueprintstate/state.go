// Package blueprintstate replays canonical assignment events without filesystem I/O.
package blueprintstate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Hans-Einar/SDP/SDPTool/bootstrap"
	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
)

const Schema = "sdp-blueprint-assignment/1"

var Timestamp = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(?:\.[0-9]{1,9})?(?:Z|[+-][0-9]{2}:[0-9]{2})$`)
var Hash = regexp.MustCompile(`^[a-f0-9]{64}$`)
var ID = regexp.MustCompile(`^BPA-[A-Za-z0-9][A-Za-z0-9.-]{0,95}$`)
var PlanID = regexp.MustCompile(`^PLAN-[A-Z][A-Z0-9]*-[0-9]{4,}$`)
var States = []string{"draft", "ready", "assigned", "in-progress", "review", "completed", "on-hold", "canceled", "superseded"}

type Principal struct{ Actor, Authority string }
type Binding struct {
	Bundle    string `json:"bundle"`
	Task      string `json:"task"`
	System    string `json:"system"`
	Revision  string `json:"revision"`
	Retained  string `json:"retained"`
	Plan      string `json:"plan"`
	Milestone string `json:"milestone"`
	ModelArea string `json:"modelArea"`
	From      string `json:"from"`
	To        string `json:"to"`
}
type Request struct {
	Schema     string   `json:"schema"`
	EventID    string   `json:"eventId"`
	OccurredAt string   `json:"occurredAt"`
	ID         string   `json:"assignmentId"`
	Action     string   `json:"action"`
	Expected   string   `json:"expectedEvent"`
	Binding    *Binding `json:"binding,omitempty"`
	Evidence   string   `json:"evidence,omitempty"`
	Trace      string   `json:"traceEvent,omitempty"`
	Assignee   string   `json:"assignee,omitempty"`
	Reason     string   `json:"reason"`
	Successor  string   `json:"successor,omitempty"`
}
type Proof struct {
	Assessment     string `json:"assessment"`
	EvidenceDigest string `json:"evidenceDigest"`
	CodeDigest     string `json:"codeDigest,omitempty"`
	TraceEvent     string `json:"traceEvent,omitempty"`
	TraceDigest    string `json:"traceDigest,omitempty"`
}
type Payload struct {
	Schema    string  `json:"schema"`
	Authority string  `json:"authority"`
	Request   Request `json:"request"`
	Proof     *Proof  `json:"proof,omitempty"`
}
type Event struct {
	SchemaVersion string  `json:"schemaVersion"`
	EventID       string  `json:"eventId"`
	EventType     string  `json:"eventType"`
	OccurredAt    string  `json:"occurredAt"`
	Actor         string  `json:"actor"`
	Commit        *string `json:"commit"`
	SubjectID     string  `json:"subjectId"`
	Payload       Payload `json:"payload"`
}
type State struct {
	ID                 string  `json:"assignmentId"`
	State              string  `json:"workState"`
	EventID            string  `json:"eventId"`
	OccurredAt         string  `json:"occurredAt"`
	Binding            Binding `json:"binding"`
	Assignee           string  `json:"assignee,omitempty"`
	Paused             string  `json:"pausedState,omitempty"`
	Ready              *Proof  `json:"readiness,omitempty"`
	ReadyEvidence      string  `json:"readinessEvidence,omitempty"`
	Submission         *Proof  `json:"submission,omitempty"`
	SubmissionEvidence string  `json:"submissionEvidence,omitempty"`
	Reviewer           string  `json:"reviewer,omitempty"`
	AcceptedEvent      string  `json:"acceptedEvent,omitempty"`
	Successor          string  `json:"successor,omitempty"`
}

func Relative(p string) bool {
	return len(p) <= 4096 && p != "" && p != "." && path.Clean(p) == p && !strings.HasPrefix(p, "/") && p != ".." && !strings.HasPrefix(p, "../") && !strings.ContainsAny(p, "\\:\x00")
}
func nonempty(s string) bool { return strings.TrimSpace(s) != "" && len(s) <= 4096 }
func Strict(data []byte, out any) error {
	if !utf8.Valid(data) {
		return fmt.Errorf("invalid UTF-8")
	}
	if err := bootstrap.StrictJSON(data, out); err != nil {
		return err
	}
	return exact(data, reflect.TypeOf(out).Elem())
}
func exact(raw []byte, t reflect.Type) error {
	if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		if t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.String {
			return nil
		}
		return fmt.Errorf("null field is not allowed")
	}
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	names := map[string]reflect.Type{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := strings.Split(f.Tag.Get("json"), ",")[0]
		names[name] = f.Type
		if !strings.Contains(f.Tag.Get("json"), ",omitempty") {
			if _, ok := fields[name]; !ok {
				return fmt.Errorf("missing required field %s", name)
			}
		}
	}
	for k, v := range fields {
		ft, ok := names[k]
		if !ok {
			return fmt.Errorf("unknown or noncanonical field %s", k)
		}
		if err := exact(v, ft); err != nil {
			return err
		}
	}
	return nil
}
func ValidateRequest(r Request) error {
	if len(r.EventID) > 4096 || len(r.Expected) > 4096 || len(r.Trace) > 4096 {
		return fmt.Errorf("assignment field limit")
	}
	if r.Schema != Schema || !ID.MatchString(r.ID) || !nonempty(r.Reason) {
		return fmt.Errorf("invalid assignment request")
	}
	if !Timestamp.MatchString(r.OccurredAt) {
		return fmt.Errorf("invalid assignment timestamp")
	}
	if _, err := time.Parse(time.RFC3339Nano, r.OccurredAt); err != nil {
		return err
	}
	if r.Binding != nil {
		b := r.Binding
		if !Relative(b.Bundle) || !Hash.MatchString(b.Revision) || !Hash.MatchString(b.Retained) || !PlanID.MatchString(b.Plan) || !nonempty(b.Milestone) || !nonempty(b.Task) || !nonempty(b.System) || !(b.ModelArea == "." || Relative(b.ModelArea)) || !nonempty(b.From) || !nonempty(b.To) {
			return fmt.Errorf("invalid assignment binding")
		}
	}
	if (r.Action == "create") != (r.Binding != nil) || (r.Action == "assign") != (r.Assignee != "") || (r.Action == "supersede") != (r.Successor != "") {
		return fmt.Errorf("unexpected action fields")
	}
	wantsEvidence := r.Action == "adopt-readiness" || r.Action == "submit"
	if wantsEvidence != (r.Evidence != "") || (r.Evidence != "" && !Relative(r.Evidence)) || (r.Action == "submit") != (r.Trace != "") {
		return fmt.Errorf("unexpected evidence fields")
	}
	if r.Successor != "" && (!ID.MatchString(r.Successor) || r.Successor == r.ID) {
		return fmt.Errorf("invalid successor")
	}
	return nil
}
func Next(old State, e Event) (State, error) {
	r := e.Payload.Request
	bad := func() (State, error) {
		return State{}, fmt.Errorf("invalid %s transition for %s from %s", r.Action, r.ID, old.State)
	}
	if err := ValidateRequest(r); err != nil {
		return State{}, err
	}
	if e.SchemaVersion != "1.0" || e.EventType != "x-blueprint:"+r.Action || e.SubjectID != r.ID || e.EventID != r.EventID || e.OccurredAt != r.OccurredAt || !nonempty(e.Actor) || e.Payload.Schema != Schema {
		return bad()
	}
	if old.EventID != r.Expected || (old.ID != "" && old.ID != r.ID) {
		return bad()
	}
	if old.OccurredAt != "" {
		a, _ := time.Parse(time.RFC3339Nano, old.OccurredAt)
		b, _ := time.Parse(time.RFC3339Nano, e.OccurredAt)
		if b.Before(a) {
			return bad()
		}
	}
	if old.State == "completed" || old.State == "canceled" || old.State == "superseded" {
		return bad()
	}
	authority := e.Payload.Authority
	required := "controller"
	switch r.Action {
	case "start", "submit":
		required = "assignee"
	case "reject", "accept-review":
		required = "reviewer"
	}
	if authority != required || (required == "assignee" && (old.Assignee == "" || old.Assignee != e.Actor)) || (required == "reviewer" && (old.Assignee == "" || old.Assignee == e.Actor)) {
		return bad()
	}
	proof := e.Payload.Proof
	needsProof := r.Action == "adopt-readiness" || r.Action == "submit" || r.Action == "accept-review" || r.Action == "complete"
	if needsProof != (proof != nil) {
		return bad()
	}
	if proof != nil {
		if !Hash.MatchString(proof.Assessment) || !Hash.MatchString(proof.EvidenceDigest) {
			return bad()
		}
		if r.Action != "adopt-readiness" && (!Hash.MatchString(proof.CodeDigest) || !Hash.MatchString(proof.TraceDigest) || proof.TraceEvent == "") {
			return bad()
		}
		if r.Action == "adopt-readiness" && (proof.CodeDigest != "" || proof.TraceEvent != "" || proof.TraceDigest != "") {
			return bad()
		}
	}
	n := old
	n.ID = r.ID
	n.EventID = e.EventID
	n.OccurredAt = e.OccurredAt
	switch r.Action {
	case "create":
		if old.ID != "" || r.Expected != "" {
			return bad()
		}
		n.Binding = *r.Binding
		n.State = "draft"
	case "adopt-readiness":
		if old.State != "draft" {
			return bad()
		}
		n.Ready = proof
		n.ReadyEvidence = r.Evidence
		n.State = "ready"
	case "assign":
		if old.State != "ready" || !nonempty(r.Assignee) {
			return bad()
		}
		n.Assignee = r.Assignee
		n.State = "assigned"
	case "start":
		if old.State != "assigned" {
			return bad()
		}
		n.State = "in-progress"
	case "submit":
		if old.State != "in-progress" || proof.TraceEvent != r.Trace {
			return bad()
		}
		n.Submission = proof
		n.SubmissionEvidence = r.Evidence
		n.AcceptedEvent = ""
		n.Reviewer = ""
		n.State = "review"
	case "reject":
		if old.State != "review" {
			return bad()
		}
		n.Submission = nil
		n.SubmissionEvidence = ""
		n.AcceptedEvent = ""
		n.Reviewer = ""
		n.State = "in-progress"
	case "accept-review":
		if old.State != "review" || old.AcceptedEvent != "" || !reflect.DeepEqual(proof, old.Submission) {
			return bad()
		}
		n.AcceptedEvent = e.EventID
		n.Reviewer = e.Actor
	case "complete":
		if old.State != "review" || old.AcceptedEvent == "" || !reflect.DeepEqual(proof, old.Submission) {
			return bad()
		}
		n.State = "completed"
	case "hold":
		if old.State == "" || old.State == "on-hold" {
			return bad()
		}
		n.Paused = old.State
		n.State = "on-hold"
	case "resume":
		if old.State != "on-hold" || old.Paused == "" {
			return bad()
		}
		n.State = old.Paused
		n.Paused = ""
	case "cancel", "supersede":
		if old.State == "" {
			return bad()
		}
		n.State = map[string]string{"cancel": "canceled", "supersede": "superseded"}[r.Action]
		n.Successor = r.Successor
	default:
		return bad()
	}
	return n, nil
}
func Replay(raw []byte) (map[string]State, map[string]Event, error) {
	states := map[string]State{}
	events := map[string]Event{}
	if err := projecthistory.Validate(raw); err != nil {
		return nil, nil, err
	}
	for _, line := range bytes.Split(raw, []byte{'\n'}) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var head struct {
			Type string `json:"eventType"`
		}
		if err := json.Unmarshal(line, &head); err != nil {
			return nil, nil, err
		}
		if !strings.HasPrefix(head.Type, "x-blueprint:") {
			continue
		}
		var e Event
		if err := Strict(line, &e); err != nil {
			return nil, nil, err
		}
		if e.Payload.Request.Action == "supersede" {
			if successor, ok := states[e.Payload.Request.Successor]; !ok || successor.Binding.System != states[e.SubjectID].Binding.System {
				return nil, nil, fmt.Errorf("missing successor assignment")
			}
		}
		n, err := Next(states[e.SubjectID], e)
		if err != nil {
			return nil, nil, err
		}
		states[n.ID] = n
		events[e.EventID] = e
	}
	return states, events, nil
}
