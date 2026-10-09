package governancecontroller

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

type skillHost interface {
	StartSkillTurn(context.Context, string, string, string, []governance.SkillInput) (governancehost.Turn, error)
}

func (c *Controller) selectSkills(ctx context.Context) (governance.SkillInventory, []governance.SkillInput, error) {
	inventory := governance.SkillInventory{Status: "unavailable", UseEvidence: "not-observed", Project: []governance.SkillAvailability{}}
	if len(c.SkillPaths) > 8 {
		return inventory, nil, fmt.Errorf("at most eight explicit skills")
	}
	inspector, ok := c.Host.(interface {
		InspectSkills(context.Context, string) (governance.SkillInventory, error)
	})
	if ok {
		observed, err := inspector.InspectSkills(ctx, c.Root)
		if err == nil {
			inventory = observed
		} else if len(c.SkillPaths) > 0 {
			return inventory, nil, fmt.Errorf("selected skill inventory unavailable: %w", err)
		}
	}
	if len(c.SkillPaths) == 0 {
		return inventory, nil, nil
	}
	if inventory.Status != "observed" || inventory.ErrorCount != 0 {
		return inventory, nil, fmt.Errorf("selected skill inventory unavailable or incomplete")
	}
	if _, ok := c.Host.(skillHost); !ok {
		return inventory, nil, fmt.Errorf("host cannot submit explicit skill input")
	}
	selected := []governance.SkillInput{}
	for _, path := range c.SkillPaths {
		if filepath.IsAbs(path) || filepath.Clean(path) != path || path == "." || strings.HasPrefix(path, "..") || filepath.Base(path) != "SKILL.md" {
			return inventory, nil, fmt.Errorf("skill requires a project-relative SKILL.md path")
		}
		resolved, err := filepath.EvalSymlinks(filepath.Join(c.Root, path))
		if err != nil {
			return inventory, nil, fmt.Errorf("selected skill path unavailable")
		}
		rel, err := filepath.Rel(c.Root, resolved)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return inventory, nil, fmt.Errorf("selected skill must remain in project")
		}
		matches := []governance.SkillAvailability{}
		for _, s := range inventory.Project {
			if s.Path == filepath.ToSlash(rel) {
				matches = append(matches, s)
			}
		}
		if len(matches) != 1 || !matches[0].Enabled || matches[0].Digest == "unavailable" {
			return inventory, nil, fmt.Errorf("selected skill missing, disabled or unreadable")
		}
		skill := matches[0]
		names := 0
		for _, s := range inventory.Project {
			if s.Name == skill.Name {
				names++
			}
		}
		if names != 1 {
			return inventory, nil, fmt.Errorf("selected skill name ambiguous")
		}
		selected = append(selected, governance.SkillInput{Name: skill.Name, Path: resolved, Digest: skill.Digest})
	}
	if err := governance.CheckSkillSnapshots(c.Root, selected); err != nil {
		return inventory, nil, err
	}
	return inventory, selected, nil
}
