package sdptool

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SDPTool/model"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	ui "github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
)

type SourceInfo struct {
	Model
	Kind         string `json:"kind"`
	ArtifactKind string `json:"artifactKind,omitempty"`
	ArtifactID   string `json:"artifactId,omitempty"`
	Preliminary  bool   `json:"preliminary,omitempty"`
	State        string `json:"state"`
	Revision     string `json:"revision,omitempty"`
	Diagnostic   string `json:"diagnostic,omitempty"`
}

// Path identity is stable across edits and distinct for equal basenames.
func sourceID(path string) string {
	base := strings.ToLower(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)))
	base = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(base, "-")
	base = strings.Trim(base, "-")
	if base == "" {
		base = "source"
	}
	if len(base) > 40 {
		base = base[:40]
	}
	return base + "-" + documents.Hash([]byte(filepath.ToSlash(path)))[:16]
}
func projectID(name string) string {
	v := regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(strings.ToLower(name), "-")
	v = strings.Trim(v, "-")
	if v == "" {
		v = "project"
	}
	if len(v) > 63 {
		v = v[:46] + "-" + documents.Hash([]byte(name))[:16]
	}
	return v
}

// Discover observes one SDP area. It does not write an index, watch the filesystem
// or consult navigation.json. A viewer owns refresh and its in-memory snapshot.
func Discover(selected string) (Project, error) { return discover(selected, true) }

func discover(selected string, withNavigation bool) (Project, error) {
	p := Project{Schema: Version, Operation: "discover", Capabilities: map[string]string{}, Installation: map[string]any{"state": "unknown"}, Sources: []SourceInfo{}, Plans: []string{}}
	bad := func(code string, e error) (Project, error) { p.Status = code; return p, failure(code, e) }
	if selected == "" {
		selected = "."
	}
	dir, e := physical(selected)
	if e != nil {
		return bad("invalid", e)
	}
	st, e := os.Stat(dir)
	if e != nil || !st.IsDir() {
		return bad("missing", fmt.Errorf("selected directory unavailable: %s", dir))
	}
	area := filepath.Join(dir, "SDP")
	if st, e = os.Lstat(area); os.IsNotExist(e) {
		if filepath.Base(dir) != "SDP" {
			return bad("missing", fmt.Errorf("no SDP directory at selection or immediate child"))
		}
		area = dir
	} else if e != nil {
		return bad("invalid", e)
	} else if !st.IsDir() || st.Mode()&os.ModeSymlink != 0 {
		return bad("invalid", fmt.Errorf("SDP child must be a directory, not a symlink"))
	}
	p.Area = area
	p.Root = filepath.Dir(area)
	pending, e := pendingInstallations(area)
	if e != nil {
		return bad("invalid", e)
	}
	if len(pending) > 0 {
		p.Status = "incomplete"
		p.Installation = map[string]any{"state": "incomplete", "operations": pending}
		return p, nil
	}
	p.Inventory = Inventory{SchemaVersion: "2.0", ProjectID: projectID(filepath.Base(p.Root)), ProcessProfile: "sdp-five-phase/0.1", Models: []Model{}, SDUI: []Model{}}
	manifest := filepath.Join(area, "SDP-project.manifest.yaml")
	if _, e = os.Lstat(manifest); e == nil {
		f, err := resolvePath(p.Root, "SDP/SDP-project.manifest.yaml")
		if err != nil {
			return bad("invalid", err)
		}
		facts, err := readYAML(f)
		if err != nil {
			return bad("invalid", err)
		}
		p.Inventory.ProjectManifest = "SDP/SDP-project.manifest.yaml"
		if project, ok := facts["project"].(map[string]any); ok {
			if name, ok := project["name"].(string); ok && name != "" {
				p.Inventory.ProjectID = projectID(name)
			}
		}
		if installed, ok := facts["installed"].(map[string]any); ok {
			rel, ok := installed["manifestPath"].(string)
			if !ok || rel == "" {
				return bad("invalid", fmt.Errorf("missing installed.manifestPath"))
			}
			dest, err := resolvePath(area, rel)
			if err != nil {
				return bad("invalid", err)
			}
			installedFacts, err := readYAML(dest, "1.0", "2.0", "3.0")
			if err != nil {
				return bad("invalid", err)
			}
			if err = validateProcessFacts(installedFacts); err != nil {
				return bad("invalid", err)
			}
			p.Installation = map[string]any{"state": "declared", "projectManifest": f, "installedManifest": dest, "facts": installedFacts, "validation": "schema version and YAML structure only; use installer validator for full conformance"}
		}
	} else if !os.IsNotExist(e) {
		return bad("invalid", e)
	}
	if st, e = os.Lstat(filepath.Join(area, "KanBan")); e == nil && st.IsDir() && st.Mode()&os.ModeSymlink == 0 {
		p.Inventory.KanBan = "SDP/KanBan"
	}
	if e = p.scanArea(); e != nil {
		return bad("limit", e)
	}
	for _, n := range p.Files {
		if n.ID == "files/Sessions" && n.Kind == "directory" {
			p.Inventory.Sessions = "SDP/Sessions"
		}
	}
	if len(p.Plans) == 1 {
		p.Inventory.ImplementationPlan = p.Plans[0]
	}
	for _, entry := range []struct {
		name    string
		present bool
	}{{"sdl", len(p.Inventory.Models) > 0}, {"sdui", len(p.Inventory.SDUI) > 0}, {"kanban", p.Inventory.KanBan != ""}, {"sessions", p.Inventory.Sessions != ""}, {"implementation-plan", len(p.Plans) > 0}} {
		state := "absent"
		if entry.present {
			state = "discovered"
		}
		p.Capabilities[entry.name] = state
	}
	p.Status = "valid"
	if !withNavigation {
		return p, nil
	}
	nav, e := Navigation(p, "")
	if e != nil {
		return bad("limit", e)
	}
	p.Navigation = &nav
	return p, nil
}

func (p *Project) scanArea() error {
	p.Files = []Node{{ID: "files", Kind: "directory", Label: "SDP", State: "available"}}
	positions := map[string]int{".": 0}
	sources, total := 0, int64(0)
	artifacts := map[string]*model.Artifact{}
	return filepath.WalkDir(p.Area, func(path string, d fs.DirEntry, walkErr error) error {
		if path == p.Area {
			return walkErr
		}
		rel, e := filepath.Rel(p.Area, path)
		if e != nil {
			return e
		}
		rel = filepath.ToSlash(rel)
		if d != nil && d.IsDir() && (d.Name() == ".git" || d.Name() == ".sdp-operations" || d.Name() == ".sdp-backups") {
			return filepath.SkipDir
		}
		if d != nil && d.IsDir() {
			if model.OperationArea(path) {
				return filepath.SkipDir
			}
			if artifacts[filepath.Dir(path)] != nil && (d.Name() == ".commits" || d.Name() == ".merge") {
				return filepath.SkipDir
			}
		}
		// WalkDir reports an unreadable directory again after its initial visit.
		// Amend the canonical node rather than publishing a duplicate identity.
		if walkErr != nil {
			if i, exists := positions[rel]; exists {
				p.Files[i].State = "unavailable"
				p.Files[i].Diagnostic = walkErr.Error()
				p.Files[i].Target = nil
				return filepath.SkipDir
			}
		}
		if len(p.Files) >= 10000 {
			return fmt.Errorf("SDP discovery exceeds 10000 filesystem entries")
		}
		parent := filepath.ToSlash(filepath.Dir(rel))
		pos, ok := positions[parent]
		if !ok {
			return fmt.Errorf("missing discovery parent")
		}
		id := "files/" + rel
		n := Node{ID: id, Kind: "file", Label: filepath.Base(path), State: "available", Target: &Target{Operation: "open", Project: p.Inventory.ProjectID, Path: path}}
		if walkErr != nil {
			n.State = "unavailable"
			n.Diagnostic = walkErr.Error()
		}
		if d != nil && d.IsDir() {
			n.Kind = "directory"
			positions[rel] = len(p.Files)
			artifact, err := model.Inspect(path)
			if err != nil {
				n.State = "invalid"
				n.Diagnostic = err.Error()
				n.Target = nil
				p.Files[pos].Children = append(p.Files[pos].Children, id)
				p.Files = append(p.Files, n)
				return filepath.SkipDir
			}
			if artifact != nil {
				artifacts[path] = artifact
				n.ArtifactKind = artifact.Kind
				n.ArtifactID = artifact.ID
				n.Preliminary = artifact.Kind == "work"
			}
		}
		if d != nil && d.Type()&os.ModeSymlink != 0 {
			n.Kind = "symlink"
			n.State = "not-followed"
			n.Target = nil
		}
		if d != nil && !d.IsDir() && d.Type()&os.ModeType != 0 && d.Type()&os.ModeSymlink == 0 {
			n.State = "unsupported"
			n.Target = nil
		}
		if info, err := dInfo(d); err == nil {
			n.Target = copyRevision(n.Target, fmt.Sprintf("%d:%d", info.Size(), info.ModTime().UnixNano()))
		}
		if walkErr == nil && d != nil && d.Type().IsRegular() && !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".design" || ext == ".sdui" {
				sources++
				if sources > 256 {
					return fmt.Errorf("SDP discovery exceeds 256 source files")
				}
				b, err := readSource(path)
				total += int64(len(b))
				if total > 64<<20 {
					return fmt.Errorf("SDP discovery exceeds 64 MiB source input")
				}
				sourceRel := "SDP/" + rel
				info := inspectSource(sourceRel, path, b, err)
				for parent := filepath.Dir(path); parent != p.Area && parent != filepath.Dir(parent); parent = filepath.Dir(parent) {
					if a := artifacts[parent]; a != nil {
						info.ArtifactKind = a.Kind
						info.ArtifactID = a.ID
						info.Preliminary = a.Kind == "work"
						n.ArtifactKind = a.Kind
						n.ArtifactID = a.ID
						n.Preliminary = info.Preliminary
						break
					}
				}
				p.Sources = append(p.Sources, info)
				n.State = info.State
				n.Diagnostic = info.Diagnostic
				n.Target = copyRevision(n.Target, info.Revision)
				if info.Kind == "sdui" {
					p.Inventory.SDUI = append(p.Inventory.SDUI, info.Model)
				} else {
					p.Inventory.Models = append(p.Inventory.Models, info.Model)
				}
			} else if ext == ".md" && strings.HasPrefix(rel, "Sessions/") {
				b, err := boundedFile(path)
				total += int64(len(b))
				if total > 64<<20 {
					return fmt.Errorf("SDP discovery exceeds 64 MiB source and Session input")
				}
				if err != nil {
					n.State, n.Diagnostic, n.Target = "unavailable", err.Error(), nil
				} else {
					n.Target = copyRevision(n.Target, documents.Hash(b))
				}
			} else if ext == ".md" && strings.HasPrefix(rel, "05--Implementation/") {
				if b, err := boundedFile(path); err == nil && strings.Contains(string(b), "| PlanType | ImplementationPlan |") {
					p.Plans = append(p.Plans, "SDP/"+rel)
				}
			}
		}
		p.Files[pos].Children = append(p.Files[pos].Children, id)
		p.Files = append(p.Files, n)
		if strings.Count(rel, "/") >= 64 && d != nil && d.IsDir() {
			p.Files[len(p.Files)-1].State = "unavailable"
			p.Files[len(p.Files)-1].Diagnostic = "directory depth limit"
			return filepath.SkipDir
		}
		if walkErr != nil && d != nil && d.IsDir() {
			return filepath.SkipDir
		}
		return nil
	})
}
func dInfo(d fs.DirEntry) (fs.FileInfo, error) {
	if d == nil {
		return nil, fmt.Errorf("unavailable")
	}
	return d.Info()
}
func copyRevision(t *Target, r string) *Target {
	if t != nil {
		t.Revision = r
	}
	return t
}

var declaredSDL = regexp.MustCompile(`(?m)^\s*language\s+([A-Za-z][A-Za-z0-9-]*)\s+version\s+([0-9]+\.[0-9]+)\.`)
var declaredSDUI = regexp.MustCompile(`(?m)^\s*sdui\s+([0-9]+\.[0-9]+)\s*;`)

func inspectSource(rel, path string, b []byte, readErr error) SourceInfo {
	x := SourceInfo{Model: Model{ID: sourceID(rel), Source: rel, System: filepath.Base(filepath.Dir(path))}, Kind: "sdl", State: "invalid", Revision: documents.Hash(b)}
	if strings.EqualFold(filepath.Ext(path), ".sdui") {
		x.Kind = "sdui"
	}
	if readErr != nil {
		x.Diagnostic = readErr.Error()
		return x
	}
	text := string(b)
	if x.Kind == "sdui" {
		if h := declaredSDUI.FindStringSubmatch(text); h != nil {
			x.Profile = "sdui/" + h[1]
		}
		doc, err := ui.Parse(text)
		if err == nil {
			x.Profile = doc.Profile
			_, err = ui.Normalize(doc)
		}
		if err != nil {
			if x.Profile != "" && x.Profile != "sdui/0.2" {
				x.State = "unsupported"
			}
			x.Diagnostic = err.Error()
			return x
		}
		x.State = "validated"
		return x
	}
	if h := declaredSDL.FindStringSubmatch(text); h != nil {
		x.Profile = h[1] + "/" + h[2]
	}
	model, err := parser.Parse(text)
	if err != nil {
		if x.Profile != "" && x.Profile != "design-core/0.5" && x.Profile != "design-core/0.6" {
			x.State = "unsupported"
		}
		x.Diagnostic = err.Error()
		return x
	}
	x.Profile = "design-core/" + model.Header.Version
	if model.Header.Version == "0.6" {
		systems := []string{}
		for _, d := range model.Declarations {
			if d.Kind == "system" {
				systems = append(systems, d.Name.Name)
			}
		}
		if len(systems) == 0 {
			x.State = "context-required"
			return x
		}
		x.System = systems[0]
	}
	v, _, err := loadModel(path)
	if err != nil {
		x.Diagnostic = err.Error()
		return x
	}
	x.State = "validated"
	x.Revision = v.Revision
	if v.System != "" {
		x.System = v.System
	}
	return x
}
