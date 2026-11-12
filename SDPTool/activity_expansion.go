package sdptool

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const activityPageSize = 100

func activityPresent(p Project) bool {
	_, err := os.Lstat(filepath.Join(p.Area, "Activities"))
	return !os.IsNotExist(err)
}

func initialActivityTree(p Project) ([]Node, error) {
	full := activityNodes(p)
	byID := map[string]Node{}
	for _, n := range full {
		byID[n.ID] = n
	}
	root, ok := byID["activities"]
	if !ok {
		return nil, fmt.Errorf("missing Activity root")
	}
	nodes := []Node{root}
	for _, id := range root.Children {
		n := byID[id]
		nodes = append(nodes, deferActivityNode(p, n, n.ID))
	}
	return nodes, nil
}

func deferActivityNode(p Project, n Node, placedID string) Node {
	original := n.ID
	n.ID = placedID
	if len(n.Children) > 0 {
		n.Children = nil
		n.ChildrenDeferred = true
		n.ExpansionTarget = activityExpansionTarget(p, original, placedID, 0, n.Target.Revision)
	}
	return n
}

func activityExpansionTarget(p Project, original, placed string, offset int, revision string) *Target {
	q := url.Values{"node": {original}, "parent": {placed}, "offset": {strconv.Itoa(offset)}}
	return &Target{Operation: "activity-tree", Project: p.Inventory.ProjectID, URI: "activity-tree://" + p.Inventory.ProjectID + "/children?" + q.Encode(), Revision: revision}
}

var activityPageSegment = regexp.MustCompile(`/more/[1-9][0-9]*`)

// ExpandActivity publishes one source-bound level. More nodes preserve all children
// without exceeding the native protocol's response size or constructing pages.
func ExpandActivity(p Project, uri, revision string) (Tree, error) {
	bad := func(e string) (Tree, error) { return Tree{}, failure("activity-tree", fmt.Errorf("%s", e)) }
	u, e := url.Parse(uri)
	if e != nil || u.Scheme != "activity-tree" || u.Host != p.Inventory.ProjectID || u.Path != "/children" || u.Fragment != "" || u.User != nil {
		return bad("invalid expansion target")
	}
	q, e := url.ParseQuery(u.RawQuery)
	if e != nil || len(q) != 3 {
		return bad("invalid expansion parameters")
	}
	for _, k := range []string{"node", "parent", "offset"} {
		if len(q[k]) != 1 {
			return bad("ambiguous expansion parameters")
		}
	}
	base, placed := q.Get("node"), q.Get("parent")
	offset, e := strconv.Atoi(q.Get("offset"))
	if e != nil || offset < 0 || offset%activityPageSize != 0 || len(placed) > 8192 || strings.Contains(placed, "..") || activityPageSegment.ReplaceAllString(placed, "") != base {
		return bad("invalid expansion placement")
	}
	if offset > 0 && !strings.HasSuffix(placed, "/more/"+strconv.Itoa(offset)) {
		return bad("invalid continuation")
	}
	full := activityNodes(p)
	index := map[string]Node{}
	for _, n := range full {
		index[n.ID] = n
	}
	parent, ok := index[base]
	if !ok || parent.Target == nil || parent.Target.Operation != "activity-view" {
		return bad("unknown expansion parent")
	}
	if revision == "" || revision != parent.Target.Revision {
		return Tree{}, failure("stale", fmt.Errorf("Activity sources changed; refresh inventory"))
	}
	if offset >= len(parent.Children) {
		return bad("expansion offset outside source children")
	}
	result := Tree{Presentation: readingPresentation, Schema: Version, Operation: "activity-tree", Project: p.Inventory.ProjectID, Revision: revision, Roots: []string{placed}, ExpansionDepthLimit: 64}
	children := parent.Children
	parent.ID, parent.Children, parent.ChildrenDeferred, parent.ExpansionTarget = placed, nil, false, nil
	end := offset + activityPageSize
	if end > len(children) {
		end = len(children)
	}
	for _, child := range children[offset:end] {
		n := index[child]
		id := placed + strings.TrimPrefix(child, base)
		if !strings.HasPrefix(id, placed+"/") {
			return bad("invalid source hierarchy")
		}
		n = deferActivityNode(p, n, id)
		parent.Children = append(parent.Children, id)
		result.Nodes = append(result.Nodes, n)
	}
	if end < len(children) {
		id := placed + "/more/" + strconv.Itoa(end)
		more := Node{ID: id, Kind: "group", Label: fmt.Sprintf("More (%d remaining)", len(children)-end), State: "available", Target: parent.Target, ChildrenDeferred: true, ExpansionTarget: activityExpansionTarget(p, base, id, end, revision)}
		parent.Children = append(parent.Children, id)
		result.Nodes = append(result.Nodes, more)
	}
	result.Nodes = append([]Node{parent}, result.Nodes...)
	return result, nil
}
