package sdptool

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
)

const programsManifest = "SDP/programs.json"

type ProgramDeclaration struct {
	ID      string   `json:"id"`
	Label   string   `json:"label"`
	Source  string   `json:"source"`
	Entry   string   `json:"entry"`
	Command []string `json:"command"`
}

type ProgramInfo struct {
	ProgramDeclaration
	Model      string   `json:"model,omitempty"`
	State      string   `json:"state"`
	Diagnostic string   `json:"diagnostic,omitempty"`
	Revision   string   `json:"revision"`
	Directory  string   `json:"directory"`
	Run        []string `json:"run,omitempty"`
}

type ProgramDiscovery struct {
	Manifest   string `json:"manifest"`
	State      string `json:"state"`
	Revision   string `json:"revision,omitempty"`
	Diagnostic string `json:"diagnostic,omitempty"`
}

// Program discovery observes declarations and executable availability, never runs
// build/probe commands. Source discovery remains available for an invalid catalog.
func (p *Project) discoverPrograms() {
	p.Programs = []ProgramInfo{}
	p.ProgramDiscovery = ProgramDiscovery{Manifest: programsManifest, State: "absent"}
	p.Capabilities["sdui-programs"] = "absent"
	if _, err := os.Lstat(filepath.Join(p.Root, programsManifest)); os.IsNotExist(err) {
		return
	}
	p.ProgramDiscovery.State = "invalid"
	p.Capabilities["sdui-programs"] = "invalid"
	path, err := resolvePath(p.Root, programsManifest)
	var data []byte
	if err == nil {
		data, err = boundedFileLimit(path, 64<<10)
	}
	var catalog struct {
		SchemaVersion string               `json:"schemaVersion"`
		Programs      []ProgramDeclaration `json:"programs"`
	}
	if err == nil {
		err = strictJSON(data, &catalog)
		p.ProgramDiscovery.Revision = documents.Hash(data)
	}
	if err == nil && (catalog.SchemaVersion != "sdp-programs/1" || catalog.Programs == nil || len(catalog.Programs) > 64) {
		err = fmt.Errorf("expected sdp-programs/1 with a programs array of at most 64 declarations")
	}
	seen := map[string]bool{}
	if err == nil {
		for _, d := range catalog.Programs {
			if !identifier.MatchString(d.ID) || seen[d.ID] {
				err = fmt.Errorf("invalid or duplicate program ID %q", d.ID)
				break
			}
			seen[d.ID] = true
		}
	}
	if err != nil {
		p.ProgramDiscovery.Diagnostic = err.Error()
		return
	}
	p.ProgramDiscovery.State = "valid"
	if len(catalog.Programs) > 0 {
		p.Capabilities["sdui-programs"] = "discovered"
	}
	for _, d := range catalog.Programs {
		x := ProgramInfo{ProgramDeclaration: d, State: "blocked", Directory: p.Root}
		b, _ := json.Marshal(d)
		x.Revision = documents.Hash(b)
		if err = p.checkProgram(&x); err != nil {
			x.Diagnostic = err.Error()
		} else {
			x.State = "runnable"
			x.Run = []string{"run", "--program", d.ID, "--revision", x.Revision}
		}
		p.Programs = append(p.Programs, x)
	}
}

func (p Project) checkProgram(x *ProgramInfo) error {
	if strings.TrimSpace(x.Label) == "" || len(x.Label) > 200 || strings.IndexFunc(x.Label, unicode.IsControl) >= 0 {
		return fmt.Errorf("program label must be nonempty text, at most 200 bytes")
	}
	var source *SourceInfo
	for i := range p.Sources {
		if p.Sources[i].Source == x.Source && p.Sources[i].Kind == "sdui" {
			source = &p.Sources[i]
			break
		}
	}
	if source == nil || source.State != "validated" {
		return fmt.Errorf("program source must name a validated discovered SDUI source")
	}
	x.Model = source.ID
	path, err := resolvePath(p.Root, x.Source)
	if err != nil {
		return err
	}
	roots, b, err := uiRoots(path, source.Profile)
	if err != nil {
		return err
	}
	root, found := roots[x.Entry]
	if !found || root.Kind != "frame" {
		return fmt.Errorf("entry %q is not a root frame in %s", x.Entry, x.Source)
	}
	x.Revision = documents.Hash([]byte(x.Revision + "\x00" + documents.Hash(b)))
	if len(x.Command) == 0 || len(x.Command) > 64 {
		return fmt.Errorf("command must contain 1 to 64 arguments")
	}
	for _, arg := range x.Command {
		if len(arg) > 8192 || strings.IndexByte(arg, 0) >= 0 {
			return fmt.Errorf("invalid command argument")
		}
	}
	_, err = programExecutable(p.Root, x.Command[0])
	return err
}

func programExecutable(root, name string) (string, error) {
	if name == "" || filepath.IsAbs(name) {
		return "", fmt.Errorf("command executable must be a PATH name or project-relative path")
	}
	if strings.ContainsAny(name, "/\\") {
		path, err := resolvePath(root, strings.TrimPrefix(name, "./"))
		if err != nil {
			return "", err
		}
		info, err := os.Stat(path)
		if err != nil || !info.Mode().IsRegular() || info.Mode()&0111 == 0 {
			return "", fmt.Errorf("program executable is missing or not executable: %s", name)
		}
		return path, nil
	}
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("program executable unavailable: %s: %w", name, err)
	}
	return filepath.Abs(path)
}
