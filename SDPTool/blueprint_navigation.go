package sdptool

import (
	"github.com/Hans-Einar/SDP/SDPTool/blueprints"
	"github.com/Hans-Einar/SDP/SDPTool/blueprintstate"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"path/filepath"
	"sort"
)

func blueprintNodes(p Project) []Node {
	catalogue := blueprints.Scan(filepath.Join(p.Area, "Blueprints"))
	tab := Node{ID: "blueprints", Kind: "tab", Label: "Blueprints", State: catalogue.State, Diagnostic: catalogue.Diagnostic}
	nodes := []Node{}
	tasks := map[string]*Node{}
	for _, entry := range catalogue.Entries {
		if entry.State != "validated" {
			id := "blueprints/invalid/" + documents.Hash([]byte(entry.Path))
			tab.Children = append(tab.Children, id)
			nodes = append(nodes, Node{ID: id, Kind: "blueprint-revision", Label: filepath.Base(entry.Path), State: "invalid", Diagnostic: entry.Diagnostic})
			continue
		}
		taskID := "blueprints/" + entry.Key
		task := tasks[taskID]
		if task == nil {
			task = &Node{ID: taskID, Kind: "blueprint", Label: entry.System + " / " + entry.TaskID, State: "available", BlueprintID: entry.Key, AssignmentState: "unknown"}
			tasks[taskID] = task
			tab.Children = append(tab.Children, taskID)
		}
		id := taskID + "/" + entry.RetainedRevision
		task.Children = append(task.Children, id)
		n := Node{ID: id, Kind: "blueprint-revision", Label: entry.RetainedRevision[:12], State: "validated", BlueprintID: entry.Key, BlueprintRevision: entry.Revision, RetainedRevision: entry.RetainedRevision, AssignmentState: "unknown", Preliminary: entry.Preliminary, Diagnostic: "Assignment progress and live-source freshness are not evaluated"}
		for _, name := range []string{"index.md", "changes.md", "context.md", "obligations.md", "evidence.md", "changes.mmd", "context.mmd"} {
			child := id + "/" + name
			n.Children = append(n.Children, child)
			target := &Target{Operation: "open", Project: p.Inventory.ProjectID, Path: filepath.Join(entry.Path, name), Revision: entry.Files[name]}
			nodes = append(nodes, Node{ID: child, Kind: "document", Label: name, State: "available", Target: target})
			if name == "index.md" {
				n.Target = target
			}
		}
		nodes = append(nodes, n)
	}
	ids := []string{}
	for id := range tasks {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		nodes = append(nodes, *tasks[id])
	}
	views, err := blueprints.AssignmentViews(p.Area)
	if err != nil {
		tab.State = "unavailable"
		tab.Diagnostic = err.Error()
		id := "blueprints/assignments-invalid"
		tab.Children = append(tab.Children, id)
		nodes = append(nodes, Node{ID: id, Kind: "diagnostic", Label: "Assignment history unavailable", State: "invalid", AssignmentState: "unknown", Diagnostic: err.Error()})
	} else {
		groups := map[string]*Node{}
		for _, v := range views {
			group := groups[v.State.State]
			if group == nil {
				group = &Node{ID: "blueprints/assignments/" + v.State.State, Kind: "group", Label: v.State.State, State: "available"}
				groups[v.State.State] = group
			}
			id := "blueprints/assignment/" + v.ID
			group.Children = append(group.Children, id)
			n := Node{ID: id, Kind: "blueprint-assignment", Label: v.ID + " / " + v.Binding.System + " / " + v.Binding.Task, State: v.Validation, AssignmentID: v.ID, AssignmentRevision: v.EventID, AssignmentState: v.State.State, RetainedRevision: v.Binding.Retained, BlueprintRevision: v.Binding.Revision, SourceFreshness: v.Freshness, ReadinessStatus: v.Readiness, EvidenceStatus: v.EvidenceStatus, Assignee: v.Assignee, Diagnostic: v.Diagnostic}
			if v.OpenPath != "" {
				n.Target = &Target{Operation: "open", Project: p.Inventory.ProjectID, Path: v.OpenPath, Revision: v.OpenRevision}
			}
			nodes = append(nodes, n)
		}
		for _, state := range blueprintstate.States {
			if g := groups[state]; g != nil {
				tab.Children = append(tab.Children, g.ID)
				nodes = append(nodes, *g)
			}
		}
		if len(views) > 0 && (tab.State == "absent" || tab.State == "empty") {
			tab.State = "available"
		}
	}
	nodes = append(nodes, tab)
	return nodes
}
