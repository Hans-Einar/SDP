package governancehost

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

// InspectSkills measures availability only. It never changes enablement, installs
// skills, reads external skill bodies or claims that an agent invoked a skill.
func (c *Client) InspectSkills(ctx context.Context, cwd string) (governance.SkillInventory, error) {
	out := governance.SkillInventory{Status: "unavailable", UseEvidence: "not-observed", Project: []governance.SkillAvailability{}}
	var response struct {
		Data []struct {
			Cwd    string      `json:"cwd"`
			Errors *[]struct{} `json:"errors"`
			Skills *[]struct {
				Name    string `json:"name"`
				Path    string `json:"path"`
				Enabled *bool  `json:"enabled"`
			} `json:"skills"`
		} `json:"data"`
	}
	if err := c.call(ctx, "skills/list", map[string]any{"cwds": []string{cwd}, "forceReload": true}, &response); err != nil {
		return out, err
	}
	if len(response.Data) != 1 || response.Data[0].Cwd != cwd {
		return out, fmt.Errorf("skill inventory cwd mismatch")
	}
	row := response.Data[0]
	if row.Skills == nil || row.Errors == nil {
		return out, fmt.Errorf("skill inventory fields unavailable")
	}
	if len(*row.Skills) > 4096 || len(*row.Errors) > 4096 {
		return out, fmt.Errorf("skill inventory limit")
	}
	root, err := os.OpenRoot(cwd)
	if err != nil {
		return out, err
	}
	defer root.Close()
	out.ErrorCount = len(*row.Errors)
	seen := map[string]bool{}
	for _, skill := range *row.Skills {
		if skill.Enabled == nil || !filepath.IsAbs(skill.Path) || len(skill.Path) > 4096 || skill.Name == "" || len(skill.Name) > 256 {
			return out, fmt.Errorf("invalid skill identity")
		}
		lexical, err := filepath.Rel(cwd, skill.Path)
		if err != nil {
			return out, err
		}
		if lexical == ".." || strings.HasPrefix(lexical, ".."+string(filepath.Separator)) {
			out.ExternalCount++
			continue
		}
		resolved, err := filepath.EvalSymlinks(skill.Path)
		if err != nil {
			return out, fmt.Errorf("skill path cannot be resolved")
		}
		rel, err := filepath.Rel(cwd, resolved)
		if err != nil {
			return out, err
		}
		if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			out.ExternalCount++
			continue
		}
		if rel == "." || seen[rel] || len(out.Project) >= 64 {
			return out, fmt.Errorf("duplicate or excessive project skill paths")
		}
		seen[rel] = true
		snapshot := governance.SkillAvailability{Name: skill.Name, Path: filepath.ToSlash(rel), Enabled: *skill.Enabled, Digest: "unavailable"}
		// Check type before opening and confine the open to the selected root.
		info, err := root.Stat(rel)
		if err == nil && info.Mode().IsRegular() && info.Size() <= 64<<10 {
			f, err := root.Open(rel)
			if err == nil {
				data, readErr := io.ReadAll(io.LimitReader(f, (64<<10)+1))
				f.Close()
				if readErr == nil && len(data) <= 64<<10 {
					sum := sha256.Sum256(data)
					snapshot.Digest = fmt.Sprintf("%x", sum)
				}
			}
		}
		out.Project = append(out.Project, snapshot)
	}
	sort.Slice(out.Project, func(i, j int) bool { return out.Project[i].Path < out.Project[j].Path })
	out.Status = "observed"
	return out, nil
}
