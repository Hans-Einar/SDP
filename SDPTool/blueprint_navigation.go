package sdptool

import (
	"github.com/Hans-Einar/SDP/SDPTool/blueprints"
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
	nodes = append(nodes, tab)
	return nodes
}
