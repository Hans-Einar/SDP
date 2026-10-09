package governance

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// SkillInput records an explicit request plus a local snapshot, not consumption.
type SkillInput struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Digest string `json:"digest"`
}

var skillNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,127}$`)

func ValidSkillInputs(skills []SkillInput) bool {
	if len(skills) > 8 {
		return false
	}
	names, paths := map[string]bool{}, map[string]bool{}
	for _, s := range skills {
		if !skillNamePattern.MatchString(s.Name) || !filepath.IsAbs(s.Path) || filepath.Clean(s.Path) != s.Path || len(s.Path) > 4096 || strings.ContainsAny(s.Path, "\x00\r\n") || !digestPattern.MatchString(s.Digest) || names[s.Name] || paths[s.Path] {
			return false
		}
		names[s.Name] = true
		paths[s.Path] = true
	}
	return true
}

// CheckSkillSnapshots confines reads to the project. It does not assert which
// bytes the host later consumes, and is never applied retroactively to history.
func CheckSkillSnapshots(root string, skills []SkillInput) error {
	if !ValidSkillInputs(skills) {
		return fail("invalid_request", "invalid selected skills")
	}
	if len(skills) == 0 {
		return nil
	}
	r, err := os.OpenRoot(root)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, s := range skills {
		rel, err := filepath.Rel(root, s.Path)
		if err != nil || !relative(filepath.ToSlash(rel)) {
			return fail("invalid_request", "selected skill must be project-local")
		}
		info, err := r.Stat(rel)
		if err != nil || !info.Mode().IsRegular() || info.Size() > 64<<10 {
			return fail("stale_revision", "selected skill snapshot unavailable")
		}
		f, err := r.Open(rel)
		if err != nil {
			return fail("stale_revision", "selected skill cannot be opened within project")
		}
		b, err := io.ReadAll(io.LimitReader(f, (64<<10)+1))
		f.Close()
		if err != nil || len(b) > 64<<10 || hash(b) != s.Digest {
			return fail("stale_revision", "selected skill bytes changed")
		}
	}
	return nil
}

// TurnInputs is the exact allowlisted outbound array. Snapshot digests are local
// provenance and are not part of the host's skill input protocol.
func TurnInputs(text string, skills []SkillInput) []map[string]string {
	out := []map[string]string{{"type": "text", "text": text}}
	for _, s := range skills {
		out = append(out, map[string]string{"type": "skill", "name": s.Name, "path": s.Path})
	}
	return out
}
func (v OwnerInput) WireDigest() string { return requestDigest(TurnInputs(v.Payload, v.Skills)) }

// ObservedInputDigest accepts only the measured lossless protocol subset.
// Unknown content/fields yield no correlation proof; ordinary capture can still
// retain a visible message without inventing exact outbound acceptance.
func ObservedInputDigest(raw json.RawMessage) string {
	if CheckWire(raw) != nil {
		return ""
	}
	var parts []map[string]json.RawMessage
	if json.Unmarshal(raw, &parts) != nil || len(parts) < 1 || len(parts) > 9 {
		return ""
	}
	out := []map[string]string{}
	for i, p := range parts {
		var kind string
		if json.Unmarshal(p["type"], &kind) != nil {
			return ""
		}
		keys := []string{"type", "name", "path"}
		if i == 0 && kind == "text" {
			keys = []string{"type", "text"}
			if elements, ok := p["text_elements"]; ok {
				var values []json.RawMessage
				if string(elements) == "null" || json.Unmarshal(elements, &values) != nil || len(values) != 0 {
					return ""
				}
				delete(p, "text_elements")
			}
		} else if i == 0 || kind != "skill" {
			return ""
		}
		if len(p) != len(keys) {
			return ""
		}
		part := map[string]string{}
		for _, key := range keys {
			var partValue string
			v, ok := p[key]
			if !ok || string(v) == "null" || json.Unmarshal(v, &partValue) != nil {
				return ""
			}
			part[key] = partValue
		}
		out = append(out, part)
	}
	return requestDigest(out)
}

// Require a unique consistent client message. Duplicate receipts of the same item
// are removed by RecordHost; different items with the same client ID are ambiguous.
func inputObserved(a Attempt, v OwnerInput, turnID string) bool {
	matches := 0
	for _, o := range a.Events {
		if o.Kind != "userMessage" || o.ClientID != v.ID {
			continue
		}
		if o.TurnID != turnID || o.Text != v.Payload {
			return false
		}
		if o.InputDigest != "" {
			if o.InputDigest != v.WireDigest() {
				return false
			}
		} else if len(v.Skills) > 0 || len(o.Attachments) > 0 {
			return false
		}
		matches++
	}
	return matches == 1
}
