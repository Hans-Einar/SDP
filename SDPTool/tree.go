package sdptool

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"net/url"
)

type Target struct {
	Entry     string `json:"entry,omitempty"`
	Operation string `json:"operation"`
	Project   string `json:"project"`
	Model     string `json:"model,omitempty"`
	URI       string `json:"uri,omitempty"`
	Path      string `json:"path,omitempty"`
	Revision  string `json:"revision,omitempty"`
}
type Node struct {
	ExpansionTarget    *Target  `json:"expansionTarget,omitempty"`
	ChildrenDeferred   bool     `json:"childrenDeferred,omitempty"`
	OverviewTarget     *Target  `json:"overviewTarget,omitempty"`
	AssignmentID       string   `json:"assignmentId,omitempty"`
	AssignmentRevision string   `json:"assignmentRevision,omitempty"`
	SourceFreshness    string   `json:"sourceFreshness,omitempty"`
	ReadinessStatus    string   `json:"readinessStatus,omitempty"`
	EvidenceStatus     string   `json:"evidenceStatus,omitempty"`
	Assignee           string   `json:"assignee,omitempty"`
	RetainedRevision   string   `json:"retainedRevision,omitempty"`
	BlueprintID        string   `json:"blueprintId,omitempty"`
	BlueprintRevision  string   `json:"blueprintRevision,omitempty"`
	AssignmentState    string   `json:"workState,omitempty"`
	ArtifactKind       string   `json:"artifactKind,omitempty"`
	ArtifactID         string   `json:"artifactId,omitempty"`
	Preliminary        bool     `json:"preliminary,omitempty"`
	Diagnostic         string   `json:"diagnostic,omitempty"`
	ID                 string   `json:"id"`
	Kind               string   `json:"kind"`
	Label              string   `json:"label"`
	State              string   `json:"state"`
	Children           []string `json:"children,omitempty"`
	Reference          string   `json:"reference,omitempty"`
	ExternalReference  string   `json:"externalReference,omitempty"`
	Target             *Target  `json:"target,omitempty"`
	WorkState          string   `json:"cardState,omitempty"`
	Sprint             string   `json:"sprintId,omitempty"`
	Scrum              string   `json:"scrumId,omitempty"`
}
type Tree struct {
	Presentation        string   `json:"presentation,omitempty"`
	InventoryRevision   string   `json:"inventoryRevision,omitempty"`
	Schema              string   `json:"schema"`
	Operation           string   `json:"operation"`
	Project             string   `json:"project"`
	Model               string   `json:"model,omitempty"`
	Revision            string   `json:"revision,omitempty"`
	Roots               []string `json:"roots"`
	Nodes               []Node   `json:"nodes"`
	ExpansionDepthLimit int      `json:"expansionDepthLimit"`
}

func selectionURI(project, vp, focus, diagram string) string {
	q := url.Values{"target": {"main"}, "consumer": {"xfmd"}}
	if focus != "" {
		q.Set("focus", focus)
	}
	if diagram != "" {
		q.Set("diagram", diagram)
	}
	return "sdl-view://" + project + "/" + vp + "?" + q.Encode()
}
func ModelTree(p Project, id string) (Tree, error) {
	m, path, e := p.model(id, false)
	if e != nil {
		return Tree{}, e
	}
	v, _, e := loadModel(path)
	if e != nil {
		return Tree{}, e
	}
	if v.Profile != m.Profile || (v.System != "" && v.System != m.System) {
		return Tree{}, failure("model", fmt.Errorf("discovered profile/System does not match source"))
	}
	return modelTree(p, m, v)
}
func modelTree(p Project, m Model, v *viewpoint.Views) (Tree, error) {
	t := Tree{Schema: Version, Operation: "tree", Presentation: readingPresentation, Project: p.Inventory.ProjectID, Model: m.ID, Revision: v.Revision, Roots: []string{"sdl"}, Nodes: []Node{}, ExpansionDepthLimit: 8}
	if len(v.Kinds) > 2000 || len(v.Facts) > 10000 {
		return Tree{}, failure("limit", fmt.Errorf("navigation model exceeds 2000 declarations or 10000 facts"))
	}
	prefix := "sdl/" + m.ID + "/"
	target := func(vp, collection, focus, register string) *Target {
		return readingTarget(p, m, v.Revision, vp, collection, focus, register)
	}
	root := Node{ID: "sdl", Kind: "tab", Label: "SDL", State: "available", Target: target("VP02", "", "", "")}
	for _, spec := range viewpoint.Catalog {
		n := Node{ID: prefix + spec.ID, Kind: "viewpoint", Label: spec.ID + " — " + spec.Title, State: "available", Target: target(spec.ID, "", "", "")}
		for _, c := range v.ReadingCollections(spec.ID) {
			if len(c.Names) == 0 {
				continue
			}
			group := Node{ID: n.ID + "/collection/" + c.ID, Kind: "collection", Label: c.Label, State: "available", Target: target(spec.ID, c.ID, "", "")}
			for _, name := range c.Names {
				object := Node{ID: group.ID + "/" + url.PathEscape(name), Kind: v.Kinds[name], Label: name, State: "available", Target: target(spec.ID, "", name, "")}
				group.Children = append(group.Children, object.ID)
				t.Nodes = append(t.Nodes, object)
			}
			n.Children = append(n.Children, group.ID)
			t.Nodes = append(t.Nodes, group)
		}
		if spec.ID == "VP11" {
			for _, r := range []struct{ id, label string }{{"facts", "Complete fact registry"}, {"declarations", "Declarations"}} {
				id := n.ID + "/register/" + r.id
				n.Children = append(n.Children, id)
				t.Nodes = append(t.Nodes, Node{ID: id, Kind: "collection", Label: r.label, State: "available", Target: target(spec.ID, "", "", r.id)})
			}
		}
		root.Children = append(root.Children, n.ID)
		t.Nodes = append(t.Nodes, n)
	}
	sources := Node{ID: prefix + "sources", Kind: "collection", Label: "Source files", State: "available", Target: sourceTarget(p, m, v.Revision, "")}
	for _, s := range v.Sources {
		id := sources.ID + "/" + url.PathEscape(s.Path)
		sources.Children = append(sources.Children, id)
		t.Nodes = append(t.Nodes, Node{ID: id, Kind: "file", Label: s.Path, State: "available", Target: sourceTarget(p, m, v.Revision, s.Path)})
	}
	root.Children = append(root.Children, sources.ID)
	t.Nodes = append(t.Nodes, sources, root)
	if len(t.Nodes) > 20000 {
		return Tree{}, failure("limit", fmt.Errorf("tree exceeds 20000 nodes"))
	}
	return t, nil
}
