// Package governance implements the bounded PGI1 process core. It does not run
// agents, publish project records or infer owner approval from model text.
package governance

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"regexp"
)

const Schema = "sdp-project-governance/0.1"
const MaxRequest = 1 << 20
const maxState = 16 << 20

var uuidPattern = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`)
var digestPattern = regexp.MustCompile(`^[a-f0-9]{64}$`)

type Error struct {
	Code            string `json:"code"`
	Message         string `json:"message"`
	Retryable       bool   `json:"retryable"`
	CurrentRevision uint64 `json:"currentRevision,omitempty"`
}

func (e *Error) Error() string { return e.Code + ": " + e.Message }
func fail(code, message string) error {
	return &Error{Code: code, Message: message, Retryable: code == "busy"}
}
func hash(b []byte) string { return fmt.Sprintf("%x", sha256.Sum256(b)) }
func identifier() (string, error) {
	var b [16]byte
	if _, e := rand.Read(b[:]); e != nil {
		return "", e
	}
	b[6] = (b[6] & 15) | 64
	b[8] = (b[8] & 63) | 128
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

// Caller is supplied by a trusted local adapter, never decoded from a request.
// Same-user shell access is outside this API's authority boundary.
type Caller struct {
	controller bool
	runID      string
	attemptID  string
	epoch      string
}

func LocalController() Caller { return Caller{controller: true} }

// AssignedAgent is a launch-time capability supplied by a trusted local adapter.
// Never construct it from model/tool arguments. It grants only attention requests, report, fail and
// scope-change on the selected run, not controller or owner authority.
func AssignedAgent(runID, attemptID, epoch string) (Caller, error) {
	if !uuidPattern.MatchString(runID) || !uuidPattern.MatchString(attemptID) || !uuidPattern.MatchString(epoch) {
		return Caller{}, fail("invalid_request", "invalid assigned run")
	}
	return Caller{runID: runID, attemptID: attemptID, epoch: epoch}, nil
}

type Selection struct {
	Card      string `json:"card"`
	Plan      string `json:"plan"`
	Session   string `json:"session"`
	Authority string `json:"authority"`
	Step      string `json:"step"`
}
type Source struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}
type Context struct {
	Schema     string    `json:"schema"`
	Operation  string    `json:"operation"`
	ProjectRef string    `json:"projectRef"`
	Root       string    `json:"root"`
	CheckoutID string    `json:"checkoutId"`
	Revision   string    `json:"revision"`
	Selection  Selection `json:"selection"`
	Sources    []Source  `json:"sources"`
	Coverage   string    `json:"coverage"`
	Conditions []string  `json:"conditions"`
	CardID     string    `json:"cardId"`
	PlanID     string    `json:"planId"`
	SessionID  string    `json:"sessionId"`
}
type Routine struct {
	ID         string   `json:"id"`
	Version    string   `json:"version"`
	Steps      []string `json:"steps"`
	Completion string   `json:"completion"`
}

func Maintenance() Routine {
	return Routine{"maintenance", "0.1", []string{"implement", "verify", "review", "record"}, "candidate evidence, applicable independent review and checked publication required"}
}
func (r Routine) Digest() string { b, _ := json.Marshal(r); return hash(b) }

type RouteRequest struct {
	Schema         string    `json:"schema"`
	Categories     []string  `json:"categories"`
	Intent         string    `json:"intent"`
	Selection      Selection `json:"selection"`
	RoutineID      string    `json:"routineId"`
	RoutineVersion string    `json:"routineVersion"`
	InputRevision  string    `json:"inputRevision"`
}
type Route struct {
	Schema        string  `json:"schema"`
	Operation     string  `json:"operation"`
	Disposition   string  `json:"disposition"`
	Basis         string  `json:"basis"`
	Context       Context `json:"context"`
	Routine       Routine `json:"routine"`
	RoutineDigest string  `json:"routineDigest"`
}
type Assignment struct {
	Outcome       string   `json:"outcome"`
	WritablePaths []string `json:"writablePaths"`
	ExcludedPaths []string `json:"excludedPaths"`
}
type StartRequest struct {
	Schema        string     `json:"schema"`
	ProjectRef    string     `json:"projectRef"`
	CheckoutID    string     `json:"checkoutId"`
	OperationID   string     `json:"operationId"`
	RequestID     string     `json:"requestId"`
	InputRevision string     `json:"inputRevision"`
	RoutineDigest string     `json:"routineDigest"`
	Selection     Selection  `json:"selection"`
	Assignment    Assignment `json:"assignment"`
}
type TransitionRequest struct {
	Schema           string `json:"schema"`
	ProjectRef       string `json:"projectRef"`
	CheckoutID       string `json:"checkoutId"`
	OperationID      string `json:"operationId"`
	RunID            string `json:"runId"`
	ExpectedRevision uint64 `json:"expectedRevision"`
	InputRevision    string `json:"inputRevision"`
	Action           string `json:"action"`
	Reason           string `json:"reason"`
}
type Run struct {
	Attention     []AttentionQuestion `json:"attention,omitempty"`
	Inputs        []OwnerInput        `json:"inputs,omitempty"`
	Attempts      []Attempt           `json:"attempts,omitempty"`
	ID            string              `json:"id"`
	Revision      uint64              `json:"revision"`
	RequestID     string              `json:"requestId"`
	State         string              `json:"state"`
	Step          string              `json:"step"`
	InputRevision string              `json:"inputRevision"`
	Selection     Selection           `json:"selection"`
	Assignment    Assignment          `json:"assignment"`
	Routine       Routine             `json:"routine"`
	RoutineDigest string              `json:"routineDigest"`
	Reason        string              `json:"reason"`
}
type Result struct {
	Schema      string `json:"schema"`
	Operation   string `json:"operation"`
	ProjectRef  string `json:"projectRef"`
	CheckoutID  string `json:"checkoutId"`
	OperationID string `json:"operationId"`
	Run         Run    `json:"run"`
	Coverage    string `json:"coverage"`
}
type Binding struct {
	Schema     string `json:"schema"`
	Operation  string `json:"operation"`
	ProjectRef string `json:"projectRef"`
	Root       string `json:"root"`
	CheckoutID string `json:"checkoutId"`
}
