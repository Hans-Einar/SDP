package sdptool

import (
	"context"
	"fmt"
)

// DiscoverOptions changes projection only, never which ordinary sources exist.
type DiscoverOptions struct {
	Navigation string `json:"navigation,omitempty"`
}

// DiscoverWithOptions preserves the default discovery payload for old clients.
func DiscoverWithOptions(selected string, options DiscoverOptions) (Project, error) {
	switch options.Navigation {
	case "", "eager":
		return Discover(selected)
	case "lazy-sdl":
	default:
		return Project{}, failure("arguments", fmt.Errorf("navigation must be eager or lazy-sdl"))
	}
	p, err := discover(selected, false)
	if err != nil || p.Status != "valid" {
		return p, err
	}
	nav, err := navigation(p, "", true)
	if err != nil {
		p.Status = "limit"
		return p, err
	}
	p.Navigation = &nav
	p.NavigationMode = "lazy-sdl"
	return p, nil
}

// ModelTreeOptions binds a deferred expansion to the discovered source revision.
type ModelTreeOptions struct {
	Model    string `json:"model" nonempty:"true"`
	Revision string `json:"revision" nonempty:"true"`
}

// ExpandModel returns one SDL graph without other project tabs or model graphs.
func ExpandModel(selected string, options ModelTreeOptions) (Tree, error) {
	if options.Model == "" || options.Revision == "" {
		return Tree{}, failure("arguments", fmt.Errorf("model and revision are required for model-only expansion"))
	}
	p, err := discover(selected, false)
	if err != nil {
		return Tree{}, err
	}
	if p.Status != "valid" {
		return Tree{}, failure("incomplete", fmt.Errorf("project is not available for model expansion"))
	}
	t, err := ModelTree(p, options.Model)
	if err != nil {
		return Tree{}, err
	}
	if t.Revision != options.Revision {
		return Tree{}, failure("stale", fmt.Errorf("model revision changed; refresh discovery before expanding"))
	}
	return t, nil
}

func discoveryAction() actionSpec {
	spec := defineAction(ActionDescriptor{ID: "project.discover", Label: "Discover project", Description: "Read project sources, navigation and available content without writing an index.", Group: "Project", IconKey: "project-discover", Mutation: "none", OutputSchema: Version}, performDiscover)
	spec.descriptor.InputSchema.Properties["navigation"].Enum = []string{"eager", "lazy-sdl"}
	return spec
}

func modelTreeAction() actionSpec {
	return defineAction(ActionDescriptor{ID: "project.model-tree", Label: "Expand SDL model", Description: "Read one complete SDL model tree at its discovered revision.", Group: "Project", IconKey: "project-discover", Mutation: "none", OutputSchema: Version}, performModelTree)
}
func performModelTree(ctx context.Context, c actionContext, o ModelTreeOptions) (any, int, error) {
	if err := ctx.Err(); err != nil {
		return nil, 130, err
	}
	t, err := ExpandModel(c.Root, o)
	return t, 0, err
}
