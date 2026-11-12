package activities

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"
)

var namespaceRE = regexp.MustCompile(`^[A-Z][A-Z0-9]*$`)
var uuidRE = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
var hashRE = regexp.MustCompile(`^[a-f0-9]{64}$`)
var activityRE = regexp.MustCompile(`^ACTIVITY-([A-Z][A-Z0-9]*)-([0-9]{4,})$`)
var sessionRE = regexp.MustCompile(`^SESSION-([A-Z][A-Z0-9]*)-([0-9]{4,})$`)
var cardRE = regexp.MustCompile(`^KB-[A-Z][A-Z0-9]*-[0-9]{3,}$`)

func digest(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func portable(p string) error {
	if p == "" || !utf8.ValidString(p) || strings.ContainsAny(p, "\\\x00:\r\n<>\"|?*") || strings.HasPrefix(p, "/") || path.Clean(p) != p {
		return fmt.Errorf("invalid portable relative path %q", p)
	}
	for _, r := range p {
		if r < 32 || r == 127 {
			return fmt.Errorf("invalid control in path")
		}
	}
	for _, v := range strings.Split(p, "/") {
		if v == "." || v == ".." || strings.HasSuffix(v, ".") || strings.HasSuffix(v, " ") {
			return fmt.Errorf("invalid path component")
		}
		base := strings.ToUpper(strings.Split(v, ".")[0])
		if regexp.MustCompile(`^(CON|PRN|AUX|NUL|COM[1-9]|LPT[1-9])$`).MatchString(base) {
			return fmt.Errorf("reserved path component")
		}
	}
	return nil
}
func (s *Service) path(rel string) (string, error) {
	if e := portable(rel); e != nil {
		return "", e
	}
	if s.records != nil {
		return filepath.Join("/captured-records", filepath.FromSlash(rel)), nil
	}
	root, e := filepath.Abs(s.Area)
	if e != nil {
		return "", e
	}
	p := filepath.Join(root, filepath.FromSlash(rel))
	for q := p; ; q = filepath.Dir(q) {
		st, e := os.Lstat(q)
		if e != nil && !os.IsNotExist(e) {
			return "", e
		}
		if e == nil && st.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("symlink forbidden: %s", q)
		}
		if filepath.Dir(q) == q {
			break
		}
	}
	return p, nil
}
func (s *Service) read(rel string) ([]byte, error) {
	if s.records != nil {
		if e := portable(rel); e != nil {
			return nil, e
		}
		b, ok := s.records[rel]
		if !ok {
			return nil, os.ErrNotExist
		}
		if len(b) > 16<<20 {
			return nil, fmt.Errorf("captured record exceeds limit")
		}
		return append([]byte(nil), b...), nil
	}
	p, e := s.path(rel)
	if e != nil {
		return nil, e
	}
	st, e := os.Stat(p)
	if e != nil {
		return nil, e
	}
	if !st.Mode().IsRegular() || st.Size() > 16<<20 {
		return nil, fmt.Errorf("record must be bounded regular file")
	}
	return os.ReadFile(p)
}

// Decode strictly decodes one request, rejecting duplicate keys, unknown fields and trailing input.
func Decode(data []byte, out any) error {
	if !utf8.Valid(data) {
		return fmt.Errorf("invalid UTF-8")
	}
	d := json.NewDecoder(bytes.NewReader(data))
	if e := uniqueJSON(d); e != nil {
		return e
	}
	if e := canonicalJSON(data, reflect.TypeOf(out)); e != nil {
		return e
	}
	d = json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if e := d.Decode(out); e != nil {
		return e
	}
	if e := d.Decode(new(any)); e != io.EOF {
		return fmt.Errorf("trailing JSON input")
	}
	return nil
}
func uniqueJSON(d *json.Decoder) error {
	v, e := d.Token()
	if e != nil {
		return e
	}
	if v == json.Delim('{') {
		seen := map[string]bool{}
		for d.More() {
			k, e := d.Token()
			if e != nil {
				return e
			}
			s := k.(string)
			if seen[s] {
				return fmt.Errorf("duplicate JSON key %s", s)
			}
			seen[s] = true
			if e = uniqueJSON(d); e != nil {
				return e
			}
		}
		_, e = d.Token()
		return e
	}
	if v == json.Delim('[') {
		for d.More() {
			if e = uniqueJSON(d); e != nil {
				return e
			}
		}
		_, e = d.Token()
		return e
	}
	return nil
}
func identityInRecord(data []byte, id string) bool {
	for _, line := range strings.Split(string(data), "\n") {
		cells := strings.Split(line, "|")
		if len(cells) < 4 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(cells[1]))
		if (key == "id" || key == "session reference") && strings.TrimSpace(cells[2]) == id {
			return true
		}
	}
	var record map[string]any
	if yaml.Unmarshal(data, &record) == nil {
		for _, key := range []string{"id", "identity", "sessionId", "planId", "blueprintId", "assignmentId"} {
			if actual, ok := record[key].(string); ok {
				return actual == id
			}
		}
	}
	return false
}
func validateRefs(s *Service, refs []Reference) error {
	seen := map[string]bool{}
	for _, r := range refs {
		if strings.TrimSpace(r.Role) == "" || strings.TrimSpace(r.ID) == "" || seen[r.Role+"\x00"+r.ID] {
			return fmt.Errorf("invalid or duplicate typed reference")
		}
		seen[r.Role+"\x00"+r.ID] = true
		if strings.Contains(strings.ToLower(r.Role), "depend") {
			return fmt.Errorf("dependency references are not implemented")
		}
		if e := validateReferenceShape(r); e != nil {
			return e
		}
		projectID := s.ProjectID
		if r.ProjectID != "" {
			projectID = r.ProjectID
		}
		switch r.Role {
		case "session", "prior-discussion":
			if !strings.HasPrefix(r.ID, "SESSION-"+projectID+"-") {
				return fmt.Errorf("reference type/project mismatch")
			}
		case "plan", "implementation-plan":
			if !strings.HasPrefix(r.ID, "PLAN-"+projectID+"-") && !strings.HasPrefix(r.ID, "MAINT-"+projectID+"-") {
				return fmt.Errorf("reference type/project mismatch")
			}
		case "card", "source-card":
			if !strings.HasPrefix(r.ID, "KB-"+projectID+"-") {
				return fmt.Errorf("reference type/project mismatch")
			}
		}
		if r.Availability == "unavailable" {
			if r.ProjectID == s.ProjectID {
				return fmt.Errorf("local reference cannot be marked external unavailable")
			}
			continue
		}
		if r.ProjectID != "" && r.ProjectID != s.ProjectID {
			return fmt.Errorf("external references require explicit unavailable disposition")
		}
		b, e := s.read(r.Path)
		if e != nil {
			return fmt.Errorf("reference %s: %w", r.ID, e)
		}
		if r.Revision != "" && digest(b) != r.Revision {
			return fmt.Errorf("reviewed reference revision changed")
		}
		expected := ""
		switch r.Role {
		case "session", "prior-discussion":
			expected = "SESSION-"
		case "implementation-plan", "plan":
			expected = "PLAN-"
			if strings.HasPrefix(r.ID, "MAINT-"+s.ProjectID+"-") {
				expected = "MAINT-"
			}
		case "source-card", "card":
			expected = "KB-"
		}
		if expected != "" {
			if !strings.HasPrefix(r.ID, expected+s.ProjectID+"-") || !identityInRecord(b, r.ID) {
				return fmt.Errorf("reference identity/type mismatch: %s", r.ID)
			}
		} else if !identityInRecord(b, r.ID) {
			return fmt.Errorf("reference record has no matching identity: %s", r.ID)
		}
	}
	return nil
}
func (s *Service) validateCards(cards []string) error {
	for _, id := range cards {
		if !cardRE.MatchString(id) || !strings.HasPrefix(id, "KB-"+s.ProjectID+"-") {
			return fmt.Errorf("source card must be an existing project card")
		}
		base, e := s.path("KanBan")
		if e != nil {
			return e
		}
		matches := 0
		e = filepath.WalkDir(base, func(p string, d os.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if d.Type()&os.ModeSymlink != 0 {
				return fmt.Errorf("card symlink")
			}
			if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
				return nil
			}
			b, e := os.ReadFile(p)
			if e != nil {
				return e
			}
			if identityInRecord(b, id) {
				matches++
			}
			return nil
		})
		if e != nil {
			return e
		}
		if matches != 1 {
			return fmt.Errorf("source card %s has %d matching records", id, matches)
		}
	}
	return nil
}
func parseDescriptor(b []byte) (Descriptor, error) {
	var out Descriptor
	d := yaml.NewDecoder(bytes.NewReader(b))
	d.KnownFields(true)
	if e := d.Decode(&out); e != nil {
		return out, e
	}
	if e := d.Decode(new(any)); e != io.EOF {
		return out, fmt.Errorf("multiple YAML documents")
	}
	if out.Schema != "sdp-activity/0.1" && out.Schema != "sdp-activity/0.2" {
		return out, fmt.Errorf("unsupported descriptor schema")
	}
	if !activityRE.MatchString(out.ID) || !namespaceRE.MatchString(out.ProjectID) || strings.TrimSpace(out.Title) == "" || strings.TrimSpace(out.Owner) == "" || out.References == nil || out.SourceCards == nil {
		return out, fmt.Errorf("invalid descriptor")
	}
	if out.Schema == "sdp-activity/0.1" && (len(out.SourceCards) == 0 || len(out.Scope) > 0 || len(out.Dependencies) > 0) {
		return out, fmt.Errorf("invalid legacy descriptor")
	}
	seen := map[string]bool{}
	for _, c := range out.SourceCards {
		if !cardRE.MatchString(c) || seen[c] {
			return out, fmt.Errorf("invalid source card")
		}
		seen[c] = true
	}
	for _, list := range [][]string{out.Systems, out.Scope} {
		seen := map[string]bool{}
		for _, v := range list {
			if strings.TrimSpace(v) == "" || seen[v] {
				return out, fmt.Errorf("invalid or duplicate descriptor list item")
			}
			seen[v] = true
		}
	}
	if e := portable(out.Mandate); e != nil {
		return out, e
	}
	for _, r := range out.References {
		if out.Schema == "sdp-activity/0.1" && (r.ProjectID != "" || r.Membership != "" || r.Revision != "" || r.Availability != "" || r.UnavailableReason != "") {
			return out, fmt.Errorf("legacy reference extension")
		}
		if e := validateReferenceShape(r); e != nil {
			return out, e
		}
		if r.Role == "" || r.ID == "" {
			return out, fmt.Errorf("invalid reference")
		}
		if e := portable(r.Path); e != nil {
			return out, e
		}
	}
	for _, p := range out.Scope {
		if e := portable(p); e != nil {
			return out, e
		}
	}
	return out, nil
}
func (s *Service) resolve(ref DescriptorReference) (Descriptor, error) {
	if !hashRE.MatchString(ref.SHA256) || !strings.HasPrefix(ref.Path, "Activities/") || path.Base(ref.Path) != "activity.yaml" {
		return Descriptor{}, fmt.Errorf("invalid descriptor reference")
	}
	if _, e := s.path(ref.Path); e != nil {
		return Descriptor{}, e
	}
	b, e := s.read(ref.Path)
	if e == nil && digest(b) == ref.SHA256 {
		return parseDescriptor(b)
	}
	b, e = s.read("Activities/.descriptors/" + ref.SHA256 + ".yaml")
	if e != nil {
		return Descriptor{}, fmt.Errorf("unavailable descriptor %s: %w", ref.SHA256, e)
	}
	if digest(b) != ref.SHA256 {
		return Descriptor{}, fmt.Errorf("descriptor snapshot hash mismatch")
	}
	return parseDescriptor(b)
}

func canonicalJSON(data []byte, t reflect.Type) error {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Struct:
		var fields map[string]json.RawMessage
		if e := json.Unmarshal(data, &fields); e != nil {
			return e
		}
		allowed := map[string]reflect.Type{}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			name := strings.Split(f.Tag.Get("json"), ",")[0]
			if name == "" {
				name = f.Name
			}
			if name != "-" {
				allowed[name] = f.Type
			}
		}
		for k, v := range fields {
			ft, ok := allowed[k]
			if !ok {
				return fmt.Errorf("unknown or noncanonical field %s", k)
			}
			if e := canonicalJSON(v, ft); e != nil {
				return e
			}
		}
	case reflect.Slice, reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			return nil
		}
		var items []json.RawMessage
		if e := json.Unmarshal(data, &items); e != nil {
			return e
		}
		for _, v := range items {
			if e := canonicalJSON(v, t.Elem()); e != nil {
				return e
			}
		}
	}
	return nil
}

func validateDependencies(c Catalog, d Descriptor) error {
	graph := map[string][]Dependency{}
	for _, a := range c.Activities {
		graph[a.Descriptor.ID] = a.Descriptor.Dependencies
	}
	graph[d.ID] = d.Dependencies
	for id, deps := range graph {
		seen := map[string]bool{}
		for _, dep := range deps {
			if !activityRE.MatchString(dep.ActivityID) || seen[dep.ActivityID] || dep.ActivityID == id {
				return fmt.Errorf("invalid, duplicate or self Activity dependency")
			}
			seen[dep.ActivityID] = true
			if _, ok := graph[dep.ActivityID]; !ok {
				return fmt.Errorf("unknown Activity dependency %s", dep.ActivityID)
			}
			if dep.SatisfiedByEventID != "" {
				target, ok := c.Find(dep.ActivityID)
				if !ok {
					return fmt.Errorf("unknown dependency")
				}
				found := false
				for _, ev := range target.Events {
					if ev.EventID == dep.SatisfiedByEventID && ev.EventType == "x-activity:transitioned" && ev.Payload.To == "completed" {
						found = true
					}
				}
				if !found {
					return fmt.Errorf("dependency satisfaction event mismatch")
				}
			}
		}
	}
	visiting := map[string]bool{}
	done := map[string]bool{}
	var visit func(string) error
	visit = func(id string) error {
		if visiting[id] {
			return fmt.Errorf("Activity dependency cycle")
		}
		if done[id] {
			return nil
		}
		visiting[id] = true
		for _, dep := range graph[id] {
			if e := visit(dep.ActivityID); e != nil {
				return e
			}
		}
		visiting[id] = false
		done[id] = true
		return nil
	}
	for id := range graph {
		if e := visit(id); e != nil {
			return e
		}
	}
	return nil
}
func checkDependencySatisfaction(c Catalog, a Activity) error {
	for _, dep := range a.Descriptor.Dependencies {
		target, ok := c.Find(dep.ActivityID)
		if !ok || target.State != "completed" || dep.SatisfiedByEventID == "" {
			return fmt.Errorf("Activity dependency %s is not explicitly satisfied", dep.ActivityID)
		}
		latest := ""
		for _, ev := range target.Events {
			if ev.EventType == "x-activity:transitioned" && ev.Payload.To == "completed" {
				latest = ev.EventID
			}
		}
		if latest != dep.SatisfiedByEventID {
			return fmt.Errorf("Activity dependency satisfaction is stale")
		}
	}
	return nil
}

func validateReferenceShape(r Reference) error {
	if strings.TrimSpace(r.Role) == "" || strings.TrimSpace(r.ID) == "" {
		return fmt.Errorf("reference role and identity required")
	}
	if e := portable(r.Path); e != nil {
		return e
	}
	if r.ProjectID != "" && !namespaceRE.MatchString(r.ProjectID) {
		return fmt.Errorf("invalid reference project")
	}
	if r.Membership != "" && r.Membership != "contribution" && r.Membership != "whole-scope" {
		return fmt.Errorf("invalid reference membership")
	}
	if r.Revision != "" && !hashRE.MatchString(r.Revision) {
		return fmt.Errorf("invalid reviewed reference revision")
	}
	if r.Availability != "" && r.Availability != "unavailable" {
		return fmt.Errorf("invalid reference availability")
	}
	if r.Availability == "unavailable" && (r.ProjectID == "" || strings.TrimSpace(r.UnavailableReason) == "") {
		return fmt.Errorf("external unavailable reference requires project and reason")
	}
	if r.Availability == "" && r.UnavailableReason != "" {
		return fmt.Errorf("unavailability reason requires disposition")
	}
	return nil
}
