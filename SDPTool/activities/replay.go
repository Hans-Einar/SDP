package activities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
	"os"
	"path"
	"strings"
	"time"
)

func Replay(area, projectID string, snapshot projecthistory.Snapshot) (Catalog, error) {
	c, e := replayService(New(area, projectID), snapshot)
	if e != nil {
		return Catalog{}, e
	}
	return c, nil
}
func replayService(s *Service, snapshot projecthistory.Snapshot) (Catalog, error) {
	projectID := s.ProjectID
	c := Catalog{ProjectID: projectID, HistoryRevision: snapshot.Revision, Activities: []Activity{}}
	if !namespaceRE.MatchString(projectID) || digest(snapshot.Bytes) != snapshot.Revision {
		return c, fmt.Errorf("invalid project or snapshot revision")
	}
	if e := projecthistory.Validate(snapshot.Bytes); e != nil {
		return c, e
	}
	if err := validateBatches(snapshot.Bytes); err != nil {
		return c, err
	}
	indices := map[string]int{}
	ops := map[string]bool{}
	sessions := map[string]bool{}
	for _, line := range bytes.Split(snapshot.Bytes, []byte{'\n'}) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var head struct {
			EventType string `json:"eventType"`
		}
		if e := json.Unmarshal(line, &head); e != nil {
			return c, e
		}
		if !strings.HasPrefix(head.EventType, "x-activity:") {
			continue
		}
		var e Event
		if err := Decode(line, &e); err != nil {
			return c, err
		}
		var fields map[string]json.RawMessage
		_ = json.Unmarshal(line, &fields)
		var pf map[string]json.RawMessage
		_ = json.Unmarshal(fields["payload"], &pf)
		for _, k := range []string{"schemaVersion", "projectId", "operationId", "recordedAt", "authority", "previousEventId", "from", "to", "reason", "descriptor", "links"} {
			if _, ok := pf[k]; !ok {
				return c, fmt.Errorf("missing Activity payload %s", k)
			}
		}
		p := e.Payload
		if p.ProjectID != projectID {
			return c, fmt.Errorf("cross-project Activity event")
		}
		if p.SchemaVersion != "0.1" && p.SchemaVersion != "0.2" {
			return c, fmt.Errorf("unsupported Activity payload")
		}
		if strings.TrimSpace(p.OperationID) == "" || (p.SchemaVersion == "0.2" && !uuidRE.MatchString(p.OperationID)) || (ops[p.OperationID] && p.Batch == nil) || p.Authority.Role != "project-owner" || strings.TrimSpace(p.Authority.Source) == "" || p.Reason == "" || p.Links == nil {
			return c, fmt.Errorf("invalid Activity payload")
		}
		ops[p.OperationID] = true
		seenLinks := map[string]bool{}
		for _, link := range p.Links {
			if strings.TrimSpace(link) == "" || seenLinks[link] {
				return c, fmt.Errorf("invalid or duplicate Activity link")
			}
			seenLinks[link] = true
		}
		if _, err := time.Parse(time.RFC3339Nano, p.RecordedAt); err != nil {
			return c, err
		}
		if err := portable(p.Authority.Source); err != nil {
			return c, err
		}
		if p.SchemaVersion == "0.1" && (p.RequestDigest != "" || p.Session != nil || p.Batch != nil || p.Outcome != "" || p.Evidence != nil) {
			return c, fmt.Errorf("legacy payload extension")
		}
		if p.SchemaVersion == "0.2" && !hashRE.MatchString(p.RequestDigest) {
			return c, fmt.Errorf("invalid request digest")
		}
		for _, ref := range p.Evidence {
			if strings.TrimSpace(ref.Role) == "" || strings.TrimSpace(ref.ID) == "" {
				return c, fmt.Errorf("invalid evidence reference")
			}
			if err := portable(ref.Path); err != nil {
				return c, err
			}
		}
		d, err := s.resolve(p.Descriptor)
		if err != nil {
			return c, err
		}
		if d.ID != e.SubjectID || d.ProjectID != projectID || !strings.HasPrefix(d.ID, "ACTIVITY-"+projectID+"-") {
			return c, fmt.Errorf("descriptor subject/project mismatch")
		}
		m := activityRE.FindStringSubmatch(d.ID)
		folder := strings.Split(path.Base(path.Dir(p.Descriptor.Path)), "--")[0]
		if folder != "Activity-"+m[2] || path.Dir(path.Dir(p.Descriptor.Path)) != "Activities" {
			return c, fmt.Errorf("descriptor path/identity mismatch")
		}
		if err := validateDependencies(c, d); err != nil {
			return c, err
		}
		if p.To == "active" {
			if err := checkDependencySatisfaction(c, Activity{Descriptor: d}); err != nil {
				return c, err
			}
		}
		i, exists := indices[d.ID]
		if exists && path.Dir(p.Descriptor.Path) != c.Activities[i].Path {
			return c, fmt.Errorf("Activity home changed")
		}
		switch e.EventType {
		case "x-activity:created":
			if exists || p.PreviousEventID != nil || p.From != nil || p.To != "planned" || p.Session != nil {
				return c, fmt.Errorf("invalid creation")
			}
			i = len(c.Activities)
			indices[d.ID] = i
			c.Activities = append(c.Activities, Activity{Path: path.Dir(p.Descriptor.Path), Sessions: []Session{}})
		case "x-activity:transitioned", "x-activity:session-created", "x-activity:updated":
			if !exists || p.PreviousEventID == nil || *p.PreviousEventID != c.Activities[i].LastEventID || p.From == nil || *p.From != c.Activities[i].State {
				return c, fmt.Errorf("broken Activity predecessor")
			}
			if e.EventType == "x-activity:transitioned" {
				if p.SchemaVersion == "0.2" && p.To == "on-hold" && strings.TrimSpace(p.ResumeCondition) == "" {
					return c, fmt.Errorf("missing resume condition")
				}
				if p.SchemaVersion == "0.2" && (p.To == "completed" || p.To == "canceled") {
					action := "complete"
					if p.To == "canceled" {
						action = "cancel"
					}
					if err := validateContributionDispositions(d, action, p.ContributionDispositions); err != nil {
						return c, err
					}
				}
				if p.Session != nil || !validTransition(*p.From, p.To) || (p.SchemaVersion == "0.2" && p.To == "completed" && (strings.TrimSpace(p.Outcome) == "" || len(p.Evidence) == 0)) {
					return c, fmt.Errorf("invalid state transition")
				}
			} else if e.EventType == "x-activity:updated" {
				if p.SchemaVersion != "0.2" || p.Session != nil || p.To != *p.From || p.Batch != nil {
					return c, fmt.Errorf("invalid metadata update")
				}
			} else {
				if p.SchemaVersion != "0.2" || p.Session == nil || p.To != *p.From || !sessionRE.MatchString(p.Session.ID) || !strings.HasPrefix(p.Session.ID, "SESSION-"+projectID+"-") || sessions[p.Session.ID] || !hashRE.MatchString(p.Session.SHA256) || !strings.HasPrefix(p.Session.Path, c.Activities[i].Path+"/Sessions/") {
					return c, fmt.Errorf("invalid Session event")
				}
				if _, err := s.path(p.Session.Path); err != nil {
					return c, err
				}
				b, err := s.read("Activities/.sessions/" + p.Session.SHA256 + ".md")
				if err != nil || digest(b) != p.Session.SHA256 {
					return c, fmt.Errorf("Session content hash mismatch")
				}
				sessions[p.Session.ID] = true
				c.Activities[i].Sessions = append(c.Activities[i].Sessions, *p.Session)
			}
		default:
			return c, fmt.Errorf("unsupported Activity event %s", e.EventType)
		}
		a := &c.Activities[i]
		a.Descriptor = d
		a.DescriptorRef = p.Descriptor
		a.State = p.To
		a.LastEventID = e.EventID
		a.Events = append(a.Events, e)
		active := 0
		for _, a := range c.Activities {
			if a.State == "active" {
				active++
			}
		}
		if active > 1 {
			return c, fmt.Errorf("multiple active Activities")
		}
	}
	for _, a := range c.Activities {
		if a.State == "active" {
			if e := checkDependencySatisfaction(c, a); e != nil {
				return c, e
			}
		}
		if e := validateDependencies(c, a.Descriptor); e != nil {
			return c, e
		}
	}
	for i := range c.Activities {
		a := &c.Activities[i]
		seen := map[string]bool{}
		for _, v := range a.Sessions {
			seen[v.ID] = true
		}
		for _, r := range a.Descriptor.References {
			if r.Role == "session" && r.Availability != "unavailable" && (r.ProjectID == "" || r.ProjectID == s.ProjectID) && !seen[r.ID] {
				if !sessionRE.MatchString(r.ID) {
					return c, fmt.Errorf("invalid Session reference")
				}
				if _, e := s.path(r.Path); e != nil {
					return c, e
				}
				a.Sessions = append(a.Sessions, Session{ID: r.ID, Path: r.Path})
				seen[r.ID] = true
			}
		}
	}
	if err := s.validateHomes(c); err != nil {
		return c, err
	}
	return c, nil
}
func validTransition(from, to string) bool {
	switch from {
	case "planned":
		return to == "active" || to == "canceled"
	case "active":
		return to == "on-hold" || to == "completed" || to == "canceled"
	case "on-hold":
		return to == "active" || to == "canceled"
	case "completed", "canceled":
		return to == "planned"
	}
	return false
}

func validateBatches(data []byte) error {
	lines := bytes.Split(bytes.TrimSpace(data), []byte{'\n'})
	seen := map[string]bool{}
	for i := 0; i < len(lines); i++ {
		if len(bytes.TrimSpace(lines[i])) == 0 {
			continue
		}
		var head struct {
			EventType string `json:"eventType"`
		}
		if e := json.Unmarshal(lines[i], &head); e != nil {
			return e
		}
		if !strings.HasPrefix(head.EventType, "x-activity:") {
			continue
		}
		var first Event
		if e := json.Unmarshal(lines[i], &first); e != nil {
			return e
		}
		b := first.Payload.Batch
		if !strings.HasPrefix(first.EventType, "x-activity:") {
			continue
		}
		if seen[first.Payload.OperationID] {
			return fmt.Errorf("reused Activity operation")
		}
		seen[first.Payload.OperationID] = true
		if b == nil {
			continue
		}
		if first.EventType != "x-activity:transitioned" || first.Payload.SchemaVersion != "0.2" || b.Index != 0 || b.Count != 2 || i+1 >= len(lines) || first.Payload.From == nil || *first.Payload.From != "active" || first.Payload.To != "on-hold" {
			return fmt.Errorf("invalid or incomplete switch batch")
		}
		var second Event
		if e := Decode(lines[i+1], &second); e != nil {
			return e
		}
		p := second.Payload
		if p.Batch == nil || p.Batch.Index != 1 || p.Batch.Count != 2 || p.OperationID != first.Payload.OperationID || p.RequestDigest != first.Payload.RequestDigest || second.Actor != first.Actor || p.Authority != first.Payload.Authority || p.Reason != first.Payload.Reason || p.ProjectID != first.Payload.ProjectID || second.SubjectID == first.SubjectID || second.EventType != "x-activity:transitioned" || p.From == nil || (*p.From != "planned" && *p.From != "on-hold") || p.To != "active" {
			return fmt.Errorf("conflicting switch batch")
		}
		i++
	}
	return nil
}

func (s *Service) validateHomes(c Catalog) error {
	if s.records != nil {
		seen := map[string]string{}
		registered := map[string]string{}
		for _, a := range c.Activities {
			registered[a.Descriptor.ID] = a.Path
		}
		for p, b := range s.records {
			parts := strings.Split(p, "/")
			if len(parts) != 3 || parts[0] != "Activities" || strings.HasPrefix(parts[1], ".") || parts[2] != "activity.yaml" {
				continue
			}
			if e := portable(p); e != nil {
				return e
			}
			d, e := parseDescriptor(b)
			if e != nil {
				return e
			}
			home := path.Dir(p)
			if old, ok := seen[d.ID]; ok && old != home {
				return fmt.Errorf("duplicate Activity identity homes")
			}
			seen[d.ID] = home
			if registered[d.ID] != home {
				return fmt.Errorf("unregistered or conflicting Activity descriptor %s", d.ID)
			}
		}
		for id, p := range registered {
			if seen[id] != p {
				return fmt.Errorf("registered Activity descriptor missing: %s", id)
			}
		}
		return nil
	}

	root, e := s.path("Activities")
	if e != nil {
		return e
	}
	entries, e := os.ReadDir(root)
	if os.IsNotExist(e) {
		if len(c.Activities) == 0 {
			return nil
		}
		return fmt.Errorf("registered Activity homes missing")
	}
	if e != nil {
		return e
	}
	seen := map[string]string{}
	registered := map[string]string{}
	for _, a := range c.Activities {
		registered[a.Descriptor.ID] = a.Path
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("Activity home symlink")
		}
		if !entry.IsDir() {
			continue
		}
		rel := "Activities/" + entry.Name()
		b, e := s.read(rel + "/activity.yaml")
		if os.IsNotExist(e) {
			continue
		}
		if e != nil {
			return e
		}
		d, e := parseDescriptor(b)
		if e != nil {
			return e
		}
		if old, ok := seen[d.ID]; ok && old != rel {
			return fmt.Errorf("duplicate Activity identity homes")
		}
		seen[d.ID] = rel
		if registered[d.ID] != rel {
			return fmt.Errorf("unregistered or conflicting Activity descriptor %s", d.ID)
		}
	}
	for id, rel := range registered {
		if seen[id] != rel {
			return fmt.Errorf("registered Activity descriptor missing: %s", id)
		}
	}
	return nil
}

// ReplayRecords validates a captured inventory using the same reducer as live
// reads. Keys are portable paths relative to the SDP document area; no filesystem
// state is consulted. Callers supply current descriptors and retained snapshots.
func ReplayRecords(projectID string, snapshot projecthistory.Snapshot, records map[string][]byte) (Catalog, error) {
	captured := map[string][]byte{}
	for p, b := range records {
		if e := portable(p); e != nil {
			return Catalog{}, e
		}
		if len(b) > projecthistory.MaxBytes {
			return Catalog{}, fmt.Errorf("captured record exceeds limit")
		}
		captured[p] = append([]byte(nil), b...)
	}
	s := New("/captured-records", projectID)
	s.records = captured
	c, e := replayService(s, snapshot)
	if e != nil {
		return Catalog{}, e
	}
	return c, nil
}
