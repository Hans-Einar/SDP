package sdptool

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/activities"
)

func activityCatalog(p Project) (activities.Projection, error) {
	namespace, err := activityProjectID(p)
	if err != nil {
		return activities.Projection{}, err
	}
	catalog, err := activities.New(p.Area, namespace).Read()
	if err != nil {
		return activities.Projection{}, err
	}
	return activities.Project(p.Area, catalog, time.Now())
}
func activityTarget(p Project, selection, revision string) *Target {
	return &Target{Operation: "activity-view", Project: p.Inventory.ProjectID, URI: "activity-view://" + p.Inventory.ProjectID + "/" + selection, Revision: revision}
}
func activityNodes(p Project) []Node {
	projection, err := activityCatalog(p)
	root := Node{ID: "activities", Kind: "tab", Label: "Activities", State: "available"}
	if err != nil {
		root.State = "invalid"
		root.Diagnostic = err.Error()
		return []Node{root}
	}
	target := func(sel string) *Target { return activityTarget(p, sel, projection.Revision) }
	root.Target = target("overview")
	nodes := []Node{}
	groups := map[string]*Node{}
	for _, state := range activities.States {
		count := 0
		for _, e := range projection.Entries {
			if e.Activity.State == state {
				count++
			}
		}
		if state == "canceled" && count == 0 {
			continue
		}
		id := "activities/state/" + state
		g := &Node{ID: id, Kind: "group", Label: activityStateLabel(state), State: "empty", Target: target("state/" + state)}
		groups[state] = g
		root.Children = append(root.Children, id)
	}
	for _, e := range projection.Entries {
		a := e.Activity
		id := "activities/state/" + a.State + "/activity/" + a.Descriptor.ID
		selection := "activity/" + a.Descriptor.ID
		n := Node{ID: id, Kind: "activity", Label: a.Descriptor.Title, State: a.State, Target: target(selection)}
		g := groups[a.State]
		if g == nil {
			root.State = "invalid"
			root.Diagnostic = "Unknown replayed Activity state"
			continue
		}
		g.State = "available"
		g.Children = append(g.Children, id)
		sections := append([]activities.Phase{}, activities.Phases...)
		sections = append(sections, activities.Phase{Key: "sessions", Label: "Sessions"}, activities.Phase{Key: "history", Label: "History"}, activities.Phase{Key: "timeline", Label: "Timeline"})
		for _, section := range sections {
			sid := id + "/" + section.Key
			s := Node{ID: sid, Kind: "phase", Label: section.Label, State: "empty", Target: target(selection + "/" + section.Key)}
			if section.Key == "history" || section.Key == "timeline" {
				s.Kind = "document"
				s.State = "available"
			}
			n.Children = append(n.Children, sid)
			for _, d := range e.Documents[section.Key] {
				did := sid + "/document/" + url.PathEscape(d.Path)
				dn := Node{ID: did, Kind: "document", Label: d.Label, State: "available"}
				if d.Diagnostic != "" {
					dn.State = "unavailable"
					dn.Diagnostic = d.Diagnostic
				} else {
					dn.Target = &Target{Operation: "open", Project: p.Inventory.ProjectID, Path: filepath.Join(p.Area, filepath.FromSlash(d.Path)), Revision: d.Revision}
				}
				s.Children = append(s.Children, did)
				s.State = "available"
				nodes = append(nodes, dn)
			}
			nodes = append(nodes, s)
		}
		nodes = append(nodes, n)
	}
	for _, state := range activities.States {
		if g := groups[state]; g != nil {
			nodes = append(nodes, *g)
		}
	}
	if len(projection.Unassigned) > 0 {
		g := Node{ID: "activities/unassigned", Kind: "group", Label: "Unassigned sessions", State: "available", Target: target("unassigned")}
		root.Children = append(root.Children, g.ID)
		for _, d := range projection.Unassigned {
			id := g.ID + "/" + url.PathEscape(d.Path)
			n := Node{ID: id, Kind: "document", Label: d.Label, State: "available", Diagnostic: d.Diagnostic}
			if d.Diagnostic != "" {
				n.State = "unavailable"
			} else {
				n.Target = &Target{Operation: "open", Project: p.Inventory.ProjectID, Path: filepath.Join(p.Area, filepath.FromSlash(d.Path)), Revision: d.Revision}
			}
			g.Children = append(g.Children, id)
			nodes = append(nodes, n)
		}
		nodes = append(nodes, g)
	}
	return append(nodes, root)
}
func renderActivityView(p Project, selection string) (activities.Page, error) {
	if strings.HasPrefix(selection, "activity-view://") {
		u, e := url.Parse(selection)
		if e != nil || u.Scheme != "activity-view" || u.Host != p.Inventory.ProjectID || u.Fragment != "" {
			return activities.Page{}, fmt.Errorf("invalid Activity target")
		}
		selection = strings.TrimPrefix(u.Path, "/")
		if u.RawQuery != "" {
			selection += "?" + u.RawQuery
		}
	}
	projection, e := activityCatalog(p)
	if e != nil {
		return activities.Page{}, e
	}
	page, e := projection.Render(selection)
	if e != nil {
		return activities.Page{}, e
	}
	// Target revision describes source inputs. As-of time changes rendered intervals,
	// but is not a mutation of canonical history or source documents.
	page.Revision = projection.Revision
	if len(page.Assets) > 0 {
		base := strings.SplitN(selection, "?", 2)[0]
		controls := "\nMarkers: "
		for _, m := range []string{"all", "ledger", "commits", "none"} {
			controls += "[" + m + "](activity-view://" + p.Inventory.ProjectID + "/" + base + "?markers=" + m + ") "
		}
		controls += "\n\nZoom: "
		for _, z := range []string{"day", "week", "month"} {
			controls += "[" + z + "](activity-view://" + p.Inventory.ProjectID + "/" + base + "?zoom=" + z + ") "
		}
		cut := strings.Index(page.Markdown, "\n\n")
		if cut >= 0 {
			page.Markdown = page.Markdown[:cut] + "\n" + controls + page.Markdown[cut:]
		}
	}
	return page, nil
}
func activityStateLabel(state string) string {
	switch state {
	case "planned":
		return "Planned"
	case "active":
		return "Active"
	case "on-hold":
		return "On hold"
	case "completed":
		return "Completed"
	case "canceled":
		return "Canceled"
	}
	return "Unknown"
}
