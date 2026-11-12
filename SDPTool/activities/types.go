// Package activities owns Activity records, canonical replay and recoverable record effects.
// It deliberately has no dependency on host dispatch or governance execution state.
package activities

import (
	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
	"time"
)

type Reference struct {
	ProjectID         string `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	Membership        string `json:"membership,omitempty" yaml:"membership,omitempty"`
	Revision          string `json:"revision,omitempty" yaml:"revision,omitempty"`
	Availability      string `json:"availability,omitempty" yaml:"availability,omitempty"`
	UnavailableReason string `json:"unavailableReason,omitempty" yaml:"unavailableReason,omitempty"`
	Role              string `json:"role" yaml:"role"`
	ID                string `json:"id" yaml:"id"`
	Path              string `json:"path" yaml:"path"`
}
type Descriptor struct {
	Dependencies []Dependency `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Schema       string       `json:"schema" yaml:"schema"`
	ID           string       `json:"id" yaml:"id"`
	ProjectID    string       `json:"projectId" yaml:"projectId"`
	Title        string       `json:"title" yaml:"title"`
	Owner        string       `json:"owner" yaml:"owner"`
	Mandate      string       `json:"mandate" yaml:"mandate"`
	SourceCards  []string     `json:"sourceCards" yaml:"sourceCards"`
	Systems      []string     `json:"systems,omitempty" yaml:"systems,omitempty"`
	References   []Reference  `json:"references" yaml:"references"`
	Scope        []string     `json:"scope,omitempty" yaml:"scope,omitempty"`
}
type DescriptorReference struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}
type Authority struct {
	Role   string `json:"role"`
	Source string `json:"source"`
}
type Caller struct {
	Actor     string
	Authority Authority
}
type Session struct {
	ID     string `json:"id"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}
type BatchReference struct {
	Index int `json:"index"`
	Count int `json:"count"`
}
type Payload struct {
	ResumeCondition          string                    `json:"resumeCondition,omitempty"`
	ContributionDispositions []ContributionDisposition `json:"contributionDispositions,omitempty"`
	SchemaVersion            string                    `json:"schemaVersion"`
	ProjectID                string                    `json:"projectId"`
	OperationID              string                    `json:"operationId"`
	RecordedAt               string                    `json:"recordedAt"`
	Authority                Authority                 `json:"authority"`
	PreviousEventID          *string                   `json:"previousEventId"`
	From                     *string                   `json:"from"`
	To                       string                    `json:"to"`
	Reason                   string                    `json:"reason"`
	Descriptor               DescriptorReference       `json:"descriptor"`
	Links                    []string                  `json:"links"`
	RequestDigest            string                    `json:"requestDigest,omitempty"`
	Session                  *Session                  `json:"session,omitempty"`
	Batch                    *BatchReference           `json:"batch,omitempty"`
	Outcome                  string                    `json:"outcome,omitempty"`
	Evidence                 []Reference               `json:"evidence,omitempty"`
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
type Activity struct {
	Path          string              `json:"path"`
	Descriptor    Descriptor          `json:"descriptor"`
	DescriptorRef DescriptorReference `json:"descriptorRef"`
	State         string              `json:"state"`
	LastEventID   string              `json:"lastEventId"`
	Events        []Event             `json:"events"`
	Sessions      []Session           `json:"sessions"`
}
type Catalog struct {
	ProjectID       string     `json:"projectId"`
	HistoryRevision string     `json:"historyRevision"`
	Activities      []Activity `json:"activities"`
}

func (c Catalog) Find(id string) (Activity, bool) {
	for _, a := range c.Activities {
		if a.Descriptor.ID == id {
			return a, true
		}
	}
	return Activity{}, false
}

type CreateRequest struct {
	Dependencies            []Dependency `json:"dependencies,omitempty"`
	OperationID             string       `json:"operationId"`
	ExpectedHistoryRevision string       `json:"expectedHistoryRevision"`
	Title                   string       `json:"title"`
	Owner                   string       `json:"owner"`
	Mandate                 string       `json:"mandate"`
	SourceCards             []string     `json:"sourceCards,omitempty"`
	Systems                 []string     `json:"systems,omitempty"`
	References              []Reference  `json:"references,omitempty"`
	Scope                   []string     `json:"scope,omitempty"`
	Reason                  string       `json:"reason"`
}
type CreateSessionRequest struct {
	PhaseContributions         []string    `json:"phaseContributions,omitempty"`
	MilestoneRefs              []Reference `json:"milestoneRefs,omitempty"`
	Predecessor                *Reference  `json:"predecessor,omitempty"`
	Handoff                    *Reference  `json:"handoff,omitempty"`
	OperationID                string      `json:"operationId"`
	ExpectedHistoryRevision    string      `json:"expectedHistoryRevision"`
	ActivityID                 string      `json:"activityId"`
	ExpectedActivityEventID    string      `json:"expectedActivityEventId"`
	ExpectedDescriptorRevision string      `json:"expectedDescriptorRevision"`
	Title                      string      `json:"title"`
	Objective                  string      `json:"objective"`
	Phase                      string      `json:"phase,omitempty"`
	Plan                       *Reference  `json:"plan,omitempty"`
	Reason                     string      `json:"reason"`
}
type Result struct {
	OperationID        string `json:"operationId"`
	ActivityID         string `json:"activityId"`
	SessionID          string `json:"sessionId,omitempty"`
	Path               string `json:"path"`
	EventID            string `json:"eventId"`
	HistoryRevision    string `json:"historyRevision"`
	DescriptorRevision string `json:"descriptorRevision"`
	Pending            bool   `json:"pending,omitempty"`
	Recovered          bool   `json:"recovered,omitempty"`
}
type Service struct {
	Area, ProjectID string
	Now             func() time.Time
	fault           func(string) error
	records         map[string][]byte
}

func New(area, projectID string) *Service {
	return &Service{Area: area, ProjectID: projectID, Now: time.Now}
}
func (s *Service) Read() (Catalog, error) {
	p, e := s.path("ProjectManagement/Ledger.ndjson")
	if e != nil {
		return Catalog{}, e
	}
	snap, e := projecthistory.Read(p)
	if e != nil {
		return Catalog{}, e
	}
	return Replay(s.Area, s.ProjectID, snap)
}

type TransitionRequest struct {
	ResumeCondition                  string                    `json:"resumeCondition,omitempty"`
	ContributionDispositions         []ContributionDisposition `json:"contributionDispositions,omitempty"`
	OperationID                      string                    `json:"operationId"`
	ExpectedHistoryRevision          string                    `json:"expectedHistoryRevision"`
	ActivityID                       string                    `json:"activityId"`
	ExpectedActivityEventID          string                    `json:"expectedActivityEventId"`
	ExpectedDescriptorRevision       string                    `json:"expectedDescriptorRevision"`
	Action                           string                    `json:"action"`
	Reason                           string                    `json:"reason"`
	TargetActivityID                 string                    `json:"targetActivityId,omitempty"`
	ExpectedTargetEventID            string                    `json:"expectedTargetEventId,omitempty"`
	ExpectedTargetDescriptorRevision string                    `json:"expectedTargetDescriptorRevision,omitempty"`
	Outcome                          string                    `json:"outcome,omitempty"`
	Evidence                         []Reference               `json:"evidence,omitempty"`
}
type UpdateRequest struct {
	Dependencies               []Dependency `json:"dependencies"`
	OperationID                string       `json:"operationId"`
	ExpectedHistoryRevision    string       `json:"expectedHistoryRevision"`
	ActivityID                 string       `json:"activityId"`
	ExpectedActivityEventID    string       `json:"expectedActivityEventId"`
	ExpectedDescriptorRevision string       `json:"expectedDescriptorRevision"`
	Reason                     string       `json:"reason"`
	Title                      string       `json:"title"`
	Owner                      string       `json:"owner"`
	Scope                      []string     `json:"scope"`
	References                 []Reference  `json:"references"`
	SourceCards                []string     `json:"sourceCards"`
	Systems                    []string     `json:"systems"`
}

type Dependency struct {
	ActivityID         string `json:"activityId" yaml:"activityId"`
	SatisfiedByEventID string `json:"satisfiedByEventId,omitempty" yaml:"satisfiedByEventId,omitempty"`
}

type ContributionDisposition struct {
	ReferenceID string `json:"referenceId"`
	Disposition string `json:"disposition"`
	Outcome     string `json:"outcome"`
}
