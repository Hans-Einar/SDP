package activities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
)

const admissionType = "x-activity-admission:adopted"
const admissionDirectory = "Activities/.admission"

type ExecutionStore struct {
	StatePath string `json:"statePath"`
}
type ExecutionRef struct {
	ActivityID string `json:"activityId,omitempty"`
	StorePath  string `json:"storePath"`
	RunID      string `json:"runId"`
	Reason     string `json:"reason"`
}
type ExecutionSnapshot struct {
	Known    bool           `json:"known"`
	InFlight []ExecutionRef `json:"inFlight"`
}
type ExecutionInspector interface {
	InspectExecutions([]ExecutionStore) (ExecutionSnapshot, error)
}
type AdoptionRequest struct {
	OperationID             string           `json:"operationId"`
	ExpectedHistoryRevision string           `json:"expectedHistoryRevision"`
	Stores                  []ExecutionStore `json:"stores"`
}
type PendingLifecycle struct {
	OperationID      string `json:"operationId"`
	ActivityID       string `json:"activityId"`
	TargetActivityID string `json:"targetActivityId,omitempty"`
	Action           string `json:"action"`
}
type WorkMembership struct {
	SessionPath, PlanPath    string
	CardPaths, WritablePaths []string
}
type AdmissionBinding struct {
	ActivityID         string `json:"activityId"`
	ActivityEventID    string `json:"activityEventId"`
	DescriptorRevision string `json:"descriptorRevision"`
}
type AdmissionView struct {
	Enabled   bool              `json:"enabled"`
	ProjectID string            `json:"projectId"`
	Catalog   Catalog           `json:"catalog"`
	Stores    []ExecutionStore  `json:"stores"`
	Pending   *PendingLifecycle `json:"pending,omitempty"`
}
type admissionPayload struct {
	SchemaVersion string           `json:"schemaVersion"`
	ProjectID     string           `json:"projectId"`
	OperationID   string           `json:"operationId"`
	RequestDigest string           `json:"requestDigest"`
	Authority     Authority        `json:"authority"`
	Stores        []ExecutionStore `json:"stores"`
}
type admissionEvent struct {
	SchemaVersion string           `json:"schemaVersion"`
	EventID       string           `json:"eventId"`
	EventType     string           `json:"eventType"`
	OccurredAt    string           `json:"occurredAt"`
	Actor         string           `json:"actor"`
	Commit        *string          `json:"commit"`
	SubjectID     string           `json:"subjectId"`
	Payload       admissionPayload `json:"payload"`
}

func (s *Service) admissionRecord() (*admissionEvent, error) {
	b, err := s.read("ProjectManagement/Ledger.ndjson")
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var found *admissionEvent
	for _, line := range bytes.Split(b, []byte{'\n'}) {
		if !bytes.Contains(line, []byte("x-activity-admission")) {
			continue
		}
		var head struct {
			EventType string `json:"eventType"`
		}
		if err = json.Unmarshal(line, &head); err != nil {
			return nil, err
		}
		if !strings.HasPrefix(head.EventType, "x-activity-admission:") {
			continue
		}
		var event admissionEvent
		if err = Decode(line, &event); err != nil {
			return nil, err
		}
		p := event.Payload
		if found != nil || event.EventType != admissionType || event.SchemaVersion != "1.0" || p.SchemaVersion != "0.1" || !namespaceRE.MatchString(p.ProjectID) || !uuidRE.MatchString(p.OperationID) || !hashRE.MatchString(p.RequestDigest) || p.Authority.Role != "project-owner" || p.Authority.Source == "" || event.Actor == "" || len(p.Stores) == 0 || len(p.Stores) > 128 || event.SubjectID != "ACTIVITY-ADMISSION-"+p.ProjectID {
			return nil, fmt.Errorf("invalid or ambiguous Activity admission adoption")
		}
		if err = portable(p.Authority.Source); err != nil {
			return nil, err
		}
		seenStores := map[string]bool{}
		for _, store := range p.Stores {
			if !filepath.IsAbs(store.StatePath) || filepath.Clean(store.StatePath) != store.StatePath || seenStores[store.StatePath] {
				return nil, fmt.Errorf("invalid or duplicate adopted execution store")
			}
			seenStores[store.StatePath] = true
		}
		if s.ProjectID != "" && s.ProjectID != p.ProjectID {
			return nil, fmt.Errorf("Activity admission project mismatch")
		}
		found = &event
	}
	return found, nil
}

// AdmissionEnabled checks opt-in presence without requiring legacy histories to
// implement Activity schemas. A malformed adoption is an error, never disabled.
func (s *Service) AdmissionEnabled() (bool, error) { r, e := s.admissionRecord(); return r != nil, e }

// ReadAdmission consumes the already-held canonical transaction. It never takes
// another lock. An empty Service.ProjectID resolves the namespace pinned at adoption.
func (s *Service) ReadAdmission(tx *projecthistory.Transaction, owner string) (AdmissionView, error) {
	var out AdmissionView
	if err := tx.Check(owner); err != nil {
		return out, err
	}
	record, err := s.admissionRecord()
	if err != nil {
		return out, err
	}
	if record == nil {
		return out, nil
	}
	snap, err := tx.Read(owner)
	if err != nil {
		return out, err
	}
	out.Enabled = true
	out.ProjectID = record.Payload.ProjectID
	out.Catalog, err = Replay(s.Area, out.ProjectID, snap)
	if err != nil {
		return out, err
	}
	dir, err := s.path(admissionDirectory + "/stores")
	if err != nil {
		return out, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return out, fmt.Errorf("adopted execution registry unavailable: %w", err)
	}
	if len(entries) > 128 {
		return out, fmt.Errorf("execution registry capacity exceeded")
	}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 {
			return out, fmt.Errorf("invalid execution registry entry")
		}
		raw, err := s.read(admissionDirectory + "/stores/" + entry.Name())
		if err != nil {
			return out, err
		}
		var store ExecutionStore
		if err = Decode(raw, &store); err != nil {
			return out, err
		}
		if err = validateStore(store); err != nil {
			return out, err
		}
		if entry.Name() != digest([]byte(store.StatePath))+".json" || seen[store.StatePath] {
			return out, fmt.Errorf("ambiguous execution store")
		}
		seen[store.StatePath] = true
		out.Stores = append(out.Stores, store)
	}
	for _, store := range record.Payload.Stores {
		if !seen[store.StatePath] {
			return out, fmt.Errorf("adopted execution store missing")
		}
	}
	raw, err := s.read(admissionDirectory + "/pending.json")
	if err == nil {
		var pending PendingLifecycle
		if err = Decode(raw, &pending); err != nil {
			return out, err
		}
		if err = validatePending(pending); err != nil {
			return out, err
		}
		out.Pending = &pending
	} else if !os.IsNotExist(err) {
		return out, err
	}
	return out, nil
}

func validateStore(store ExecutionStore) error {
	if !filepath.IsAbs(store.StatePath) || filepath.Clean(store.StatePath) != store.StatePath {
		return fmt.Errorf("canonical absolute execution store required")
	}
	for p := store.StatePath; ; p = filepath.Dir(p) {
		st, err := os.Lstat(p)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		if err == nil && st.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("execution store symlink")
		}
		if filepath.Dir(p) == p {
			break
		}
	}
	return nil
}
func (s *Service) writeStore(store ExecutionStore) error {
	if err := validateStore(store); err != nil {
		return err
	}
	b, _ := json.Marshal(store)
	rel := admissionDirectory + "/stores/" + digest([]byte(store.StatePath)) + ".json"
	old, err := s.read(rel)
	if err == nil {
		if !bytes.Equal(old, b) {
			return fmt.Errorf("execution registry conflict")
		}
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	return s.write(rel, b)
}

// RegisterExecutionStore must precede every new admission into that store. Store
// metadata locates PG state; it never substitutes a cached execution status.
func (s *Service) RegisterExecutionStore(tx *projecthistory.Transaction, owner string, store ExecutionStore) error {
	view, err := s.ReadAdmission(tx, owner)
	if err != nil {
		return err
	}
	if !view.Enabled {
		return nil
	}
	if view.Pending != nil {
		return fmt.Errorf("Activity lifecycle quiescence pending")
	}
	return s.writeStore(store)
}

// AdoptAdmission records explicit inventory only. Existing execution must first
// have a known quiescent disposition; arbitrary unlisted legacy stores are not discovered.
func (s *Service) AdoptAdmission(q AdoptionRequest, caller Caller, inspector ExecutionInspector) (AdmissionView, error) {
	var out AdmissionView
	if !namespaceRE.MatchString(s.ProjectID) || !uuidRE.MatchString(q.OperationID) || !hashRE.MatchString(q.ExpectedHistoryRevision) || len(q.Stores) == 0 || len(q.Stores) > 128 || caller.Actor == "" || caller.Authority.Role != "project-owner" || caller.Authority.Source == "" || inspector == nil {
		return out, fmt.Errorf("trusted adoption, explicit store inventory and inspector required")
	}
	if _, err := s.read(caller.Authority.Source); err != nil {
		return out, err
	}
	seen := map[string]bool{}
	for _, store := range q.Stores {
		if err := validateStore(store); err != nil {
			return out, err
		}
		if seen[store.StatePath] {
			return out, fmt.Errorf("duplicate execution store")
		}
		seen[store.StatePath] = true
	}
	sort.Slice(q.Stores, func(i, j int) bool { return q.Stores[i].StatePath < q.Stores[j].StatePath })
	request, _ := json.Marshal(struct {
		Request AdoptionRequest
		Caller  Caller
	}{q, caller})
	requestDigest := digest(request)
	ledger, err := s.path("ProjectManagement/Ledger.ndjson")
	if err != nil {
		return out, err
	}
	tx, err := projecthistory.BeginTransaction(ledger)
	if err != nil {
		return out, err
	}
	defer tx.Close()
	old, err := s.admissionRecord()
	if err != nil {
		return out, err
	}
	if old != nil {
		if old.Payload.OperationID != q.OperationID || old.Payload.RequestDigest != requestDigest {
			return out, fmt.Errorf("admission adoption conflict")
		}
		return s.ReadAdmission(tx, "")
	}
	snap, err := tx.Read("")
	if err != nil {
		return out, err
	}
	if snap.Revision != q.ExpectedHistoryRevision {
		return out, projecthistory.ErrStale
	}
	if _, err = Replay(s.Area, s.ProjectID, snap); err != nil {
		return out, err
	}
	observed, err := inspector.InspectExecutions(q.Stores)
	if err != nil {
		return out, err
	}
	if !observed.Known || len(observed.InFlight) > 0 {
		return out, fmt.Errorf("execution inventory is not known quiescent")
	}
	for _, store := range q.Stores {
		if err = s.writeStore(store); err != nil {
			return out, err
		}
	}
	event := admissionEvent{SchemaVersion: "1.0", EventID: "EVT-ACTIVITY-ADMISSION-" + strings.ToUpper(q.OperationID), EventType: admissionType, OccurredAt: s.Now().UTC().Format(time.RFC3339Nano), Actor: caller.Actor, SubjectID: "ACTIVITY-ADMISSION-" + s.ProjectID, Payload: admissionPayload{SchemaVersion: "0.1", ProjectID: s.ProjectID, OperationID: q.OperationID, RequestDigest: requestDigest, Authority: caller.Authority, Stores: q.Stores}}
	raw, _ := json.Marshal(event)
	if _, err = tx.AppendBatch("", snap.Revision, [][]byte{raw}); err != nil {
		return out, err
	}
	return s.ReadAdmission(tx, "")
}

func validatePending(p PendingLifecycle) error {
	if !uuidRE.MatchString(p.OperationID) || !activityRE.MatchString(p.ActivityID) || p.Action == "" || p.TargetActivityID != "" && !activityRE.MatchString(p.TargetActivityID) {
		return fmt.Errorf("invalid pending lifecycle")
	}
	return nil
}
func (s *Service) SetPendingLifecycle(tx *projecthistory.Transaction, owner string, p PendingLifecycle) error {
	if err := tx.Check(owner); err != nil {
		return err
	}
	if err := validatePending(p); err != nil {
		return err
	}
	view, err := s.ReadAdmission(tx, owner)
	if err != nil {
		return err
	}
	if !view.Enabled {
		return fmt.Errorf("Activity admission adoption required")
	}
	raw, _ := json.Marshal(p)
	if view.Pending != nil {
		old, _ := json.Marshal(view.Pending)
		if !bytes.Equal(old, raw) {
			return fmt.Errorf("another lifecycle request pending")
		}
		return nil
	}
	return s.write(admissionDirectory+"/pending.json", raw)
}
func (s *Service) ClearPendingLifecycle(tx *projecthistory.Transaction, owner, operationID string) error {
	if err := tx.Check(owner); err != nil {
		return err
	}
	view, err := s.ReadAdmission(tx, owner)
	if err != nil {
		return err
	}
	if view.Pending == nil {
		return nil
	}
	if view.Pending.OperationID != operationID {
		return fmt.Errorf("another lifecycle request pending")
	}
	p, err := s.path(admissionDirectory + "/pending.json")
	if err != nil {
		return err
	}
	if err = os.Remove(p); err != nil {
		return err
	}
	dir, err := os.Open(filepath.Dir(p))
	if err != nil {
		return err
	}
	defer dir.Close()
	return dir.Sync()
}

// CheckAdmission requires explicit Session or plan membership and every assigned
// write within an owner's literal project-relative prefix. It never grants scope
// from a caller-asserted Activity ID or from the currently open directory.
func CheckAdmission(view AdmissionView, work WorkMembership) (AdmissionBinding, error) {
	var out AdmissionBinding
	if !view.Enabled {
		return out, nil
	}
	if view.Pending != nil {
		return out, fmt.Errorf("Activity lifecycle quiescence pending")
	}
	var selected *Activity
	for i := range view.Catalog.Activities {
		a := &view.Catalog.Activities[i]
		if a.State == "active" {
			if selected != nil {
				return out, fmt.Errorf("multiple active Activities")
			}
			selected = a
		}
	}
	if selected == nil {
		return out, fmt.Errorf("no active Activity")
	}
	// A Session has one primary Activity. A prior-discussion/external link is
	// context rather than a second membership.
	owners := map[string]bool{}
	if work.SessionPath != "" {
		for _, a := range view.Catalog.Activities {
			for _, linked := range a.Sessions {
				if sameSDPPath(work.SessionPath, linked.Path) {
					owners[a.Descriptor.ID] = true
				}
			}
			for _, ref := range a.Descriptor.References {
				if ref.Role == "session" && ref.Availability != "unavailable" && (ref.ProjectID == "" || ref.ProjectID == view.ProjectID) && sameSDPPath(work.SessionPath, ref.Path) {
					owners[a.Descriptor.ID] = true
				}
			}
		}
		if len(owners) != 1 || !owners[selected.Descriptor.ID] {
			return out, fmt.Errorf("Session requires one primary active Activity membership")
		}
	}
	member := false
	for _, ref := range selected.Descriptor.References {
		if ref.Availability == "unavailable" || ref.ProjectID != "" && ref.ProjectID != view.ProjectID {
			continue
		}
		if (ref.Role == "session" && sameSDPPath(work.SessionPath, ref.Path)) || ((ref.Role == "plan" || strings.HasSuffix(ref.Role, "-plan")) && sameSDPPath(work.PlanPath, ref.Path)) {
			member = true
		}
	}
	for _, session := range selected.Sessions {
		if sameSDPPath(work.SessionPath, session.Path) {
			member = true
		}
	}
	if !member {
		return out, fmt.Errorf("selected work has no explicit active Activity membership")
	}
	for _, write := range work.WritablePaths {
		if err := portable(write); err != nil {
			return out, err
		}
		if strings.ContainsAny(write, "*?[]") {
			return out, fmt.Errorf("Activity scope requires literal assigned paths")
		}
		allowed := false
		for _, scope := range selected.Descriptor.Scope {
			if write == scope || strings.HasPrefix(write, scope+"/") {
				allowed = true
			}
		}
		if !allowed {
			return out, fmt.Errorf("write %s is outside active Activity scope", write)
		}
	}
	return AdmissionBinding{ActivityID: selected.Descriptor.ID, ActivityEventID: selected.LastEventID, DescriptorRevision: selected.DescriptorRef.SHA256}, nil
}
func sameSDPPath(projectPath, areaPath string) bool {
	return projectPath != "" && (projectPath == "SDP/"+areaPath || projectPath == areaPath)
}

// ResolveProjectID recovers the canonical namespace from the board and Activity
// history, rejecting ambiguity. A source-project slug is never uppercased as identity.
func ResolveProjectID(area, fallback string) (string, error) {
	service := New(area, "")
	found := ""
	take := func(id string) error {
		if !namespaceRE.MatchString(id) {
			return fmt.Errorf("invalid Activity project namespace")
		}
		if found != "" && found != id {
			return fmt.Errorf("ambiguous Activity project namespace")
		}
		found = id
		return nil
	}
	board, err := service.read("KanBan/board.json")
	if err == nil {
		var fields struct {
			ProjectID string `json:"projectId"`
		}
		if err = json.Unmarshal(board, &fields); err != nil {
			return "", err
		}
		if err = take(fields.ProjectID); err != nil {
			return "", err
		}
	} else if !os.IsNotExist(err) {
		return "", err
	}
	history, err := service.read("ProjectManagement/Ledger.ndjson")
	if err == nil {
		for _, line := range bytes.Split(history, []byte{'\n'}) {
			if len(bytes.TrimSpace(line)) == 0 {
				continue
			}
			var fields struct {
				EventType string `json:"eventType"`
				Payload   struct {
					ProjectID string `json:"projectId"`
				} `json:"payload"`
			}
			if err = json.Unmarshal(line, &fields); err != nil {
				return "", err
			}
			if strings.HasPrefix(fields.EventType, "x-activity:") || strings.HasPrefix(fields.EventType, "x-activity-admission:") {
				if err = take(fields.Payload.ProjectID); err != nil {
					return "", err
				}
			}
		}
	} else if !os.IsNotExist(err) {
		return "", err
	}
	if found == "" && fallback != "" {
		if err = take(fallback); err != nil {
			return "", err
		}
	}
	if found == "" {
		return "", fmt.Errorf("canonical Activity project namespace unavailable")
	}
	return found, nil
}

// ValidateAdmissionRecords checks a captured installation/publication inventory
// without opening external state stores. It validates registry structure, never
// claims execution quiescence or grants another canonical writer.
func ValidateAdmissionRecords(history []byte, records map[string][]byte) error {
	if err := projecthistory.Validate(history); err != nil {
		return err
	}
	service := New("/captured-records", "")
	service.records = make(map[string][]byte, len(records)+1)
	for p, b := range records {
		if err := portable(p); err != nil {
			return err
		}
		service.records[p] = b
	}
	service.records["ProjectManagement/Ledger.ndjson"] = history
	record, err := service.admissionRecord()
	if err != nil || record == nil {
		return err
	}
	seen := map[string]bool{}
	for p, b := range service.records {
		if !strings.HasPrefix(p, admissionDirectory+"/stores/") {
			continue
		}
		var store ExecutionStore
		if err = Decode(b, &store); err != nil {
			return err
		}
		if !filepath.IsAbs(store.StatePath) || filepath.Clean(store.StatePath) != store.StatePath || p != admissionDirectory+"/stores/"+digest([]byte(store.StatePath))+".json" || seen[store.StatePath] {
			return fmt.Errorf("invalid captured execution registry")
		}
		seen[store.StatePath] = true
	}
	if len(seen) > 128 {
		return fmt.Errorf("execution registry capacity exceeded")
	}
	for _, store := range record.Payload.Stores {
		if !seen[store.StatePath] {
			return fmt.Errorf("adopted execution store missing")
		}
	}
	if raw, ok := service.records[admissionDirectory+"/pending.json"]; ok {
		var pending PendingLifecycle
		if err = Decode(raw, &pending); err != nil {
			return err
		}
		return validatePending(pending)
	}
	return nil
}
