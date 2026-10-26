package governance

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// SessionLeaderSkill pins a launcher-selected project-local skill. It records
// availability, not proof that an agent loaded or followed the procedure.
func (s *Service) SessionLeaderSkill(id, path string) (SkillInput, error) {
	skill, e := s.readSessionLeaderSkill(path)
	if e != nil {
		return skill, e
	}
	return skill, s.CheckSessionLeader(id, skill)
}

func (s *Service) readSessionLeaderSkill(path string) (SkillInput, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(s.root, path)
	}
	skill := SkillInput{Name: "sdp-project-leader", Path: filepath.Clean(path)}
	rel, e := filepath.Rel(s.root, skill.Path)
	if e != nil || !relative(filepath.ToSlash(rel)) {
		return skill, fail("unauthorized", "Project Leader skill must be project-local")
	}
	root, e := os.OpenRoot(s.root)
	if e != nil {
		return skill, e
	}
	defer root.Close()
	info, e := root.Stat(rel)
	if e != nil || !info.Mode().IsRegular() || info.Size() > 64<<10 {
		return skill, fail("invalid_request", "bounded regular Project Leader skill required")
	}
	f, e := root.Open(rel)
	if e != nil {
		return skill, e
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, (64<<10)+1))
	if e != nil {
		return skill, e
	}
	if len(b) > 64<<10 {
		return skill, fail("limit", "Project Leader skill too large")
	}
	text := strings.ReplaceAll(string(b), "\r\n", "\n")
	if filepath.Base(skill.Path) != "SKILL.md" || !strings.HasPrefix(text, "---\n") {
		return skill, fail("invalid_request", "named Project Leader SKILL.md frontmatter required")
	}
	end := strings.Index(text[4:], "\n---\n")
	if end < 0 {
		return skill, fail("invalid_request", "closed skill frontmatter required")
	}
	var metadata struct {
		Name string `yaml:"name"`
	}
	if e := yaml.Unmarshal([]byte(text[4:4+end]), &metadata); e != nil || metadata.Name != skill.Name {
		return skill, fail("unauthorized", "selected skill is not sdp-project-leader")
	}
	skill.Digest = hash(b)
	return skill, nil
}

// CheckSessionLeader is an adapter precondition. Core mutations independently
// check canonical scope, revisions and authority under their own locks.
func (s *Service) CheckSessionLeader(id string, skill SkillInput) error {
	if skill.Name != "sdp-project-leader" {
		return fail("unauthorized", "mandatory Project Leader skill required")
	}
	if e := CheckSkillSnapshots(s.root, []SkillInput{skill}); e != nil {
		return e
	}
	actual, e := s.readSessionLeaderSkill(skill.Path)
	if e != nil {
		return e
	}
	if actual != skill {
		return fail("stale_revision", "launcher skill identity changed")
	}
	v, e := s.ReadSession(id)
	if e != nil {
		return e
	}
	if v.Registration != "committed" {
		return fail("missing_coverage", "registered Session required before granting leader capability")
	}
	return nil
}
