package sdptool

import (
	"net/url"
	"strings"
)

// Existing overview URI behavior remains supported, while the new reading tree
// deliberately no longer advertises the old inventory's duplicate mode rows.
// Adapt the old overview regression fixtures to explicit retained URI requests;
// reading_navigation_test exercises the actual new tree contract independently.
func legacyOverviewFixture(tree Tree) Tree {
	out := tree
	out.Nodes = append([]Node(nil), tree.Nodes...)
	for i, n := range out.Nodes {
		if n.Kind == "viewpoint" && n.Target != nil {
			legacy := *n.Target
			legacy.URI = strings.Replace(legacy.URI, "projection=reading", "projection=overview", 1)
			out.Nodes[i].OverviewTarget = &legacy
		}
	}
	for _, n := range tree.Nodes {
		if n.Kind != "mode" || n.Target == nil {
			continue
		}
		u, e := url.Parse(n.Target.URI)
		if e != nil {
			continue
		}
		vp := strings.TrimPrefix(u.Path, "/")
		if vp != "VP05" && vp != "VP07" && vp != "VP08" {
			continue
		}
		q := u.Query()
		q.Del("focus")
		q.Set("projection", "overview")
		q.Set("mode", n.Label)
		u.RawQuery = q.Encode()
		target := *n.Target
		target.URI = u.String()
		out.Nodes = append(out.Nodes, Node{ID: "sdl/" + tree.Model + "/" + vp + "/overview-mode/" + url.PathEscape(n.Label), Kind: "diagram", Label: n.Label, State: "available", Target: &target})
	}
	return out
}
