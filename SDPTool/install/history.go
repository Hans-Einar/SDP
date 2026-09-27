package install

import (
	"bytes"
	"fmt"
	"github.com/Hans-Einar/SDP/SDPTool/bootstrap"
	"regexp"
	"strings"
	"time"
)

// Validate history before planning any relocation. History bytes remain immutable;
// validation checks both predecessor chains and the documents they identify.
func validateHistory(history []byte, files map[string][]byte) error {
	latest := map[string]map[string]any{}
	ids := map[string]bool{}
	for _, line := range bytes.Split(history, []byte{'\n'}) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var event map[string]any
		if e := bootstrap.StrictJSON(line, &event); e != nil {
			return e
		}
		id, _ := event["eventId"].(string)
		subject, _ := event["subjectId"].(string)
		typ, _ := event["eventType"].(string)
		when, _ := event["occurredAt"].(string)
		actor, _ := event["actor"].(string)
		payload, ok := event["payload"].(map[string]any)
		if !ok || event["schemaVersion"] != "1.0" || id == "" || subject == "" || actor == "" || ids[id] {
			return fmt.Errorf("invalid/duplicate management event")
		}
		if _, e := time.Parse(time.RFC3339, when); e != nil {
			return e
		}
		ids[id] = true
		if payload["schemaVersion"] != "0.1" && payload["schemaVersion"] != "0.2" {
			return fmt.Errorf("unsupported history payload")
		}
		for _, k := range []string{"previousEventId", "from", "to", "fromPath", "toPath", "reason", "links"} {
			if _, ok := payload[k]; !ok {
				return fmt.Errorf("missing payload %s", k)
			}
		}
		reason, ok := payload["reason"].(string)
		if !ok || reason == "" {
			return fmt.Errorf("missing history reason")
		}
		links, ok := payload["links"].([]any)
		if !ok {
			return fmt.Errorf("invalid history links")
		}
		seen := map[string]bool{}
		for _, x := range links {
			s, ok := x.(string)
			if !ok || s == "" || seen[s] {
				return fmt.Errorf("invalid/duplicate history link")
			}
			seen[s] = true
		}
		dest, ok := payload["toPath"].(string)
		if !ok {
			return fmt.Errorf("missing history path")
		}
		if e := Relative(dest); e != nil {
			return e
		}
		prior := latest[subject]
		if prior == nil {
			if payload["previousEventId"] != nil || payload["from"] != nil || payload["fromPath"] != nil || !strings.HasSuffix(typ, ":created") {
				return fmt.Errorf("broken initial history chain")
			}
		} else {
			pp := prior["payload"].(map[string]any)
			if payload["previousEventId"] != prior["eventId"] || payload["from"] != pp["to"] || payload["fromPath"] != pp["toPath"] {
				return fmt.Errorf("broken history chain: %s", id)
			}
			if pp["planType"] != payload["planType"] {
				return fmt.Errorf("plan type changed")
			}
		}
		if strings.HasPrefix(typ, "x-kanban:") {
			if !regexp.MustCompile(`^KB-[A-Z][A-Z0-9]*-[0-9]{3,}$`).MatchString(subject) || !contains([]string{"x-kanban:created", "x-kanban:moved", "x-kanban:reviewed"}, typ) {
				return fmt.Errorf("invalid card event")
			}
			if !strings.HasPrefix(dest, fmt.Sprint(payload["to"])+"/") {
				return fmt.Errorf("card path/state mismatch")
			}
		} else if strings.HasPrefix(typ, "x-management:") {
			kind, _ := payload["kind"].(string)
			prefix := map[string]string{"Plan": "PLAN", "Maintenance": "MAINT", "Sprint": "SPR", "Scrum": "SCRUM", "CodeReview": "REVIEW", "Refactor": "REFACTOR"}[kind]
			if prefix == "" || !regexp.MustCompile(`^`+prefix+`-[A-Z][A-Z0-9]*-[0-9]{4,}$`).MatchString(subject) {
				return fmt.Errorf("management kind/identity mismatch")
			}
			to, _ := payload["to"].(string)
			if !contains([]string{"planned", "active", "completed", "canceled"}, to) {
				return fmt.Errorf("unknown management state")
			}
			action := strings.TrimPrefix(typ, "x-management:")
			switch action {
			case "created":
				if prior != nil || to != "planned" && to != "active" {
					return fmt.Errorf("invalid creation")
				}
			case "started":
				if payload["from"] != "planned" || to != "active" {
					return fmt.Errorf("invalid start")
				}
			case "updated":
				if prior == nil || payload["from"] != to || payload["fromPath"] != dest || to == "completed" || to == "canceled" {
					return fmt.Errorf("invalid update")
				}
			case "completed", "canceled":
				if prior == nil || to != action || len(links) == 0 {
					return fmt.Errorf("invalid closure")
				}
			default:
				return fmt.Errorf("unsupported management event")
			}
		} else {
			return fmt.Errorf("unsupported management event type")
		}
		latest[subject] = event
	}
	cards := map[string]bool{}
	states := map[string][]string{"backlog": {"backlog", "queued"}, "active": {"ready", "in-progress", "gate-review"}, "onHold": {"onHold"}, "completed": {"completed"}, "canceled": {"canceled"}, "superseded": {"superseded"}, "irrelevant": {"irrelevant"}}
	for p, b := range files {
		if !strings.HasPrefix(p, "SDP/KanBan/") {
			continue
		}
		parts := strings.Split(strings.TrimPrefix(p, "SDP/KanBan/"), "/")
		if len(parts) != 2 || !strings.HasPrefix(parts[1], "#") || !strings.HasSuffix(parts[1], ".md") {
			continue
		}
		allowed, ok := states[parts[0]]
		if !ok {
			continue
		}
		m, e := metadataTable(b)
		if e != nil {
			return e
		}
		id := m["id"]
		ev := latest[id]
		if id == "" || cards[id] || ev == nil || !contains(allowed, m["CardState"]) {
			return fmt.Errorf("invalid card state/history: %s", p)
		}
		payload := ev["payload"].(map[string]any)
		if payload["to"] != parts[0] || payload["toPath"] != strings.TrimPrefix(p, "SDP/KanBan/") {
			return fmt.Errorf("card/history location mismatch: %s", p)
		}
		cards[id] = true
	}
	for id, ev := range latest {
		payload := ev["payload"].(map[string]any)
		if strings.HasPrefix(id, "KB-") {
			if !cards[id] {
				return fmt.Errorf("missing card %s", id)
			}
			continue
		}
		if payload["schemaVersion"] != "0.2" {
			continue
		}
		dest := payload["toPath"].(string)
		m, e := metadataTable(files["SDP/"+dest])
		if e != nil {
			return e
		}
		if m["id"] != id || m["state"] != payload["to"] {
			return fmt.Errorf("management document/state mismatch: %s", id)
		}
		if pt, ok := payload["planType"]; ok {
			if m["PlanType"] != pt || !contains([]string{"phase", "current"}, m["BranchPolicy"]) || !contains([]string{"phase", "milestone"}, m["CommitPolicy"]) {
				return fmt.Errorf("typed plan metadata mismatch: %s", id)
			}
		}
	}
	return nil
}
func metadataTable(b []byte) (map[string]string, error) {
	m := map[string]string{}
	started := false
	for _, line := range strings.Split(string(b), "\n") {
		if strings.TrimSpace(line) == "| Field | Value |" {
			started = true
			continue
		}
		if !started {
			continue
		}
		if !strings.HasPrefix(line, "|") {
			break
		}
		parts := strings.Split(line, "|")
		if len(parts) != 4 {
			return nil, fmt.Errorf("invalid metadata row")
		}
		k, v := strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2])
		if strings.HasPrefix(k, "---") {
			continue
		}
		if _, ok := m[k]; ok {
			return nil, fmt.Errorf("duplicate metadata row")
		}
		m[k] = v
	}
	return m, nil
}
