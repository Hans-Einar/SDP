package governance

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/install"
)

// Publication preserves lifecycle. In addition to the install validator's replay,
// reject unknown payload fields, invalid counters, missing index/group references
// and changed typed identities before generating any new canonical events.
func validatePublicationHistory(history []byte, files map[string][]byte) error {
	if e := install.ValidateManagementHistory(history, files); e != nil {
		return e
	}
	events, e := publicationEvents(history)
	if e != nil {
		return e
	}
	latest := map[string]map[string]any{}
	numbers := map[string]uint64{}
	when := map[string]time.Time{}
	for _, ev := range events {
		typ := fmt.Sprint(ev["eventType"])
		if strings.HasPrefix(typ, "x-blueprint:") {
			continue
		}
		id := fmt.Sprint(ev["eventId"])
		parts := regexp.MustCompile(`^(EVT-(?:PM|KB)-[A-Z][A-Z0-9]*)-([0-9]{6,})$`).FindStringSubmatch(id)
		if parts == nil {
			return fmt.Errorf("invalid canonical event identity: %s", id)
		}
		n, e := strconv.ParseUint(parts[2], 10, 64)
		if e != nil || n <= numbers[parts[1]] {
			return fmt.Errorf("non-increasing event sequence: %s", id)
		}
		numbers[parts[1]] = n
		p := ev["payload"].(map[string]any)
		subject := fmt.Sprint(ev["subjectId"])
		timestamp, _ := time.Parse(time.RFC3339Nano, fmt.Sprint(ev["occurredAt"]))
		if timestamp.Before(when[subject]) {
			return fmt.Errorf("history time reversal: %s", subject)
		}
		when[subject] = timestamp
		if commit, ok := ev["commit"]; !ok || commit != nil && (fmt.Sprint(commit) == "" || !isString(commit)) {
			return fmt.Errorf("invalid event commit")
		}
		allowed := []string{"schemaVersion", "previousEventId", "from", "to", "fromPath", "toPath", "reason", "links"}
		if strings.HasPrefix(typ, "x-kanban:") {
			allowed = append(allowed, "projectId")
			if p["schemaVersion"] == "0.2" {
				allowed = append(allowed, "lineage")
			}
			ns, ok := p["projectId"].(string)
			if !ok || !regexp.MustCompile(`^[A-Z][A-Z0-9]*$`).MatchString(ns) || !strings.HasPrefix(subject, "KB-"+ns+"-") || parts[1] != "EVT-KB-"+ns {
				return fmt.Errorf("card namespace mismatch")
			}
			for _, k := range []string{"fromPath", "toPath"} {
				if p[k] != nil && !regexp.MustCompile(`^(backlog|onHold|active|completed|canceled|superseded|irrelevant)/#[0-9]{3,}--[^/\\]+\.md$`).MatchString(fmt.Sprint(p[k])) {
					return fmt.Errorf("invalid canonical card filename")
				}
			}
		} else {
			allowed = append(allowed, "kind", "members")
			if p["schemaVersion"] == "0.2" {
				allowed = append(allowed, "planType", "plans")
			}
			if p["schemaVersion"] == "0.2" {
				_, typed := p["planType"]
				needs := p["kind"] == "Plan" || p["kind"] == "Maintenance"
				if needs != typed {
					return fmt.Errorf("typed plan coverage differs")
				}
				if typed && !slices.Contains([]string{"MaintenancePlan", "RequirementPlan", "ArchitecturePlan", "DesignPlan", "ImplementationPlan", "VerificationPlan"}, fmt.Sprint(p["planType"])) {
					return fmt.Errorf("unknown plan type")
				}
				if p["kind"] == "Maintenance" && p["planType"] != "MaintenancePlan" {
					return fmt.Errorf("wrong Maintenance plan type")
				}
			}
			if p["kind"] != "Sprint" {
				if _, ok := p["members"]; ok {
					return fmt.Errorf("members outside Sprint")
				}
				if _, ok := p["plans"]; ok {
					return fmt.Errorf("plans outside Sprint")
				}
			} else {
				for _, k := range []string{"members", "plans"} {
					if k == "plans" && p["schemaVersion"] == "0.1" {
						continue
					}
					if _, e := publicationIDs(p[k]); e != nil {
						return e
					}
				}
			}
			if old := latest[subject]; old != nil && old["payload"].(map[string]any)["kind"] != p["kind"] {
				return fmt.Errorf("management kind changed")
			}
		}
		for key := range p {
			if !slices.Contains(allowed, key) {
				return fmt.Errorf("unsupported management payload field %s", key)
			}
		}
		latest[subject] = ev
	}
	metas := map[string]map[string]string{}
	cards := map[string]string{}
	for id, ev := range latest {
		p := ev["payload"].(map[string]any)
		path := "SDP/" + fmt.Sprint(p["toPath"])
		if strings.HasPrefix(id, "KB-") {
			path = "SDP/KanBan/" + fmt.Sprint(p["toPath"])
		}
		m, e := metadata(files[path])
		if e != nil {
			return e
		}
		metas[id] = m
		if strings.HasPrefix(id, "KB-") {
			if m["project"] != strings.Split(id, "-")[1] {
				return fmt.Errorf("card namespace metadata differs")
			}
			cards[id] = fmt.Sprint(p["to"])
		} else {
			links, e := publicationIDs(p["links"])
			if e != nil {
				return e
			}
			for _, link := range links {
				if latest[link] == nil {
					return fmt.Errorf("unresolved management link: %s", link)
				}
			}
		}
	}
	index := map[string]string{}
	for _, m := range regexp.MustCompile(`(?m)^\| (KB-[A-Z][A-Z0-9]*-[0-9]+) \| [^|]+ \| (\w+) \|`).FindAllStringSubmatch(string(files["SDP/KanBan/README.md"]), -1) {
		if _, ok := index[m[1]]; ok {
			return fmt.Errorf("duplicate board index")
		}
		index[m[1]] = m[2]
	}
	if requestDigest(index) != requestDigest(cards) {
		return fmt.Errorf("board index differs from canonical card history")
	}
	for id, m := range metas {
		for _, tag := range []string{"PlanId", "ScrumId", "SprintId"} {
			target := m[tag]
			if target == "" {
				continue
			}
			ev := latest[target]
			if ev == nil {
				return fmt.Errorf("unresolved %s", tag)
			}
			p := ev["payload"].(map[string]any)
			if tag == "PlanId" {
				if _, ok := p["planType"]; !ok {
					return fmt.Errorf("untyped PlanId")
				}
				if p["to"] == "active" && (m["CardState"] == "backlog" || m["CardState"] == "queued") {
					return fmt.Errorf("active plan has backlog execution card")
				}
			} else if tag == "ScrumId" {
				if p["kind"] != "Scrum" {
					return fmt.Errorf("wrong ScrumId kind")
				}
			} else {
				key := "plans"
				if strings.HasPrefix(id, "KB-") {
					key = "members"
				}
				ids, e := publicationIDs(p[key])
				if e != nil || p["kind"] != "Sprint" || !slices.Contains(ids, id) {
					return fmt.Errorf("Sprint membership differs")
				}
			}
		}
		p := latest[id]["payload"].(map[string]any)
		if p["kind"] != "Sprint" {
			continue
		}
		for key, label := range map[string]string{"members": "Members", "plans": "Plans"} {
			ids, e := publicationIDs(p[key])
			if e != nil && key == "plans" && p["schemaVersion"] == "0.1" {
				ids = []string{}
			} else if e != nil {
				return e
			}
			declared := []string{}
			for _, v := range strings.Split(m[label], ",") {
				if strings.TrimSpace(v) != "" {
					declared = append(declared, strings.TrimSpace(v))
				}
			}
			slices.Sort(ids)
			slices.Sort(declared)
			if !slices.Equal(ids, declared) {
				return fmt.Errorf("Sprint declared membership differs")
			}
			for _, member := range ids {
				mp := metas[member]
				if mp == nil || mp["SprintId"] != id {
					return fmt.Errorf("Sprint reciprocal membership differs")
				}
				if key == "plans" && mp["PlanType"] == "" {
					return fmt.Errorf("untyped Sprint plan")
				}
				if p["to"] == "active" && (mp["state"] == "planned" || mp["CardState"] == "backlog" || mp["CardState"] == "queued") {
					return fmt.Errorf("started Sprint has unstarted work")
				}
				if key == "plans" && p["to"] == "completed" && mp["state"] != "completed" && mp["state"] != "canceled" {
					return fmt.Errorf("closed Sprint has unfinished plan")
				}
			}
		}
	}
	return nil
}
func isString(v any) bool { _, ok := v.(string); return ok }
func publicationIDs(v any) ([]string, error) {
	a, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("expected canonical ID array")
	}
	out := []string{}
	for _, v := range a {
		id, ok := v.(string)
		if !ok || id == "" || slices.Contains(out, id) {
			return nil, fmt.Errorf("invalid or duplicate canonical ID")
		}
		out = append(out, id)
	}
	return out, nil
}
