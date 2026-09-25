package sdptool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func metadata(b []byte) (map[string]string, error) {
	out := map[string]string{}
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
			return nil, fmt.Errorf("malformed metadata table")
		}
		k, v := strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2])
		if strings.HasPrefix(k, "---") {
			continue
		}
		if _, ok := out[k]; ok {
			return nil, fmt.Errorf("duplicate metadata %s", k)
		}
		out[k] = v
	}
	if out["id"] == "" || out["CardState"] == "" {
		return nil, fmt.Errorf("missing card ID or CardState")
	}
	return out, nil
}
func BoardNodes(p Project) ([]Node, string, error) {
	root, e := resolvePath(p.Root, p.Registration.KanBan)
	if e != nil {
		return nil, "", e
	}
	b, e := boundedFile(filepath.Join(root, "board.json"))
	if e != nil {
		return nil, "", e
	}
	var descriptor struct {
		Schema     string   `json:"schemaVersion"`
		Project    string   `json:"projectId"`
		Namespaces []string `json:"namespaces"`
		Ledger     string   `json:"ledger"`
		Profile    string   `json:"profile"`
	}
	if e = strictJSON(b, &descriptor); e != nil {
		return nil, "", e
	}
	if descriptor.Schema != "0.2" || descriptor.Profile != "sdp-project-management/0.1" {
		return nil, "", failure("unsupported", fmt.Errorf("unsupported KanBan board/profile"))
	}
	rel, e := filepath.Rel(p.Root, filepath.Join(root, descriptor.Ledger))
	if e != nil {
		return nil, "", e
	}
	ledger, e := resolvePath(p.Root, rel)
	if e != nil {
		return nil, "", e
	}
	history, e := boundedFile(ledger)
	if e != nil {
		return nil, "", e
	}
	latest := map[string]string{}
	previous := map[string]string{}
	for _, line := range bytes.Split(history, []byte("\n")) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var ev struct {
			ID      string `json:"eventId"`
			Type    string `json:"eventType"`
			Subject string `json:"subjectId"`
			Payload struct {
				Schema   string  `json:"schemaVersion"`
				Previous *string `json:"previousEventId"`
				Path     string  `json:"toPath"`
			} `json:"payload"`
		}
		if e = json.Unmarshal(line, &ev); e != nil {
			return nil, "", e
		}
		if !strings.HasPrefix(ev.Type, "x-kanban:") {
			continue
		}
		if ev.Payload.Schema != "0.1" && ev.Payload.Schema != "0.2" {
			return nil, "", failure("unsupported", fmt.Errorf("unsupported card event payload"))
		}
		if ev.Payload.Previous == nil {
			if previous[ev.Subject] != "" {
				return nil, "", fmt.Errorf("broken card chain")
			}
		} else if previous[ev.Subject] != *ev.Payload.Previous {
			return nil, "", fmt.Errorf("broken card chain")
		}
		previous[ev.Subject] = ev.ID
		latest[ev.Subject] = ev.Payload.Path
	}
	states := map[string][]string{"backlog": {"backlog", "queued"}, "active": {"ready", "in-progress", "gate-review"}, "onHold": {"onHold"}, "completed": {"completed"}, "canceled": {"canceled"}, "superseded": {"superseded"}, "irrelevant": {"irrelevant"}}
	folders := []string{"backlog", "active", "onHold", "completed", "canceled", "superseded", "irrelevant"}
	tab := Node{ID: "kanban", Kind: "tab", Label: "KanBan", State: "validated"}
	nodes := []Node{}
	ids := map[string]bool{}
	refs := map[string]string{}
	digest := append(append([]byte{}, b...), history...)
	for _, folder := range folders {
		status := Node{ID: "kanban/" + folder, Kind: "status", Label: folder, State: "empty"}
		tab.Children = append(tab.Children, status.ID)
		entries, err := os.ReadDir(filepath.Join(root, folder))
		if err != nil {
			return nil, "", err
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasPrefix(entry.Name(), "#") || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			relPath := filepath.Join(folder, entry.Name())
			file := filepath.Join(root, relPath)
			q, err := filepath.Rel(p.Root, file)
			if err != nil {
				return nil, "", err
			}
			file, err = resolvePath(p.Root, q)
			if err != nil {
				return nil, "", err
			}
			data, err := boundedFile(file)
			if err != nil {
				return nil, "", err
			}
			m, err := metadata(data)
			if err != nil {
				return nil, "", err
			}
			id := m["id"]
			if ids[id] || latest[id] != relPath {
				return nil, "", fmt.Errorf("card/ledger placement mismatch: %s", id)
			}
			ids[id] = true
			valid := false
			for _, s := range states[folder] {
				if s == m["CardState"] {
					valid = true
				}
			}
			if !valid {
				return nil, "", fmt.Errorf("invalid CardState for %s", id)
			}
			key := "kanban/card/" + id
			status.Children = append(status.Children, key)
			status.State = "available"
			n := Node{ID: key, Kind: "card", Label: entry.Name(), State: "available", WorkState: m["CardState"], Sprint: m["SprintId"], Scrum: m["ScrumId"], Target: &Target{Operation: "open", Project: p.Registration.ProjectID, Path: file, Revision: documents.Hash(data)}}
			if primary := m["primary"]; primary != "" {
				refs[key] = primary
			}
			nodes = append(nodes, n)
			digest = append(digest, []byte(relPath+"\x00")...)
			digest = append(digest, data...)
		}
		nodes = append(nodes, status)
	}
	if len(ids) != len(latest) {
		return nil, "", fmt.Errorf("ledger contains missing cards")
	}
	for i := range nodes {
		if primary := refs[nodes[i].ID]; primary != "" {
			if !ids[primary] {
				return nil, "", fmt.Errorf("unresolved Ref %s", primary)
			}
			nodes[i].Reference = "kanban/card/" + primary
		}
	}
	nodes = append(nodes, tab)
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID < nodes[j].ID })
	return nodes, documents.Hash(digest), nil
}
