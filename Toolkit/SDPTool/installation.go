package sdptool

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// These are producer build facts, not a consuming project's installed version.
var BuildVersion = "0.1.0-dev"
var BuildRevision = "unknown"

func validateProcessFacts(facts map[string]any) error {
	if facts["schemaVersion"] != "2.0" {
		return nil
	} // v1's declared-facts contract remains supported.
	fields := []string{"schemaVersion", "toolkitVersion", "frameworkVersion", "agentsContractVersion", "installerVersion", "toolkitInstalledAt", "sourceCommit", "skills", "capabilities", "processProfile", "managementProfile", "configurationDigest"}
	if len(facts) != len(fields) {
		return fmt.Errorf("unknown/missing installed process fields")
	}
	for _, key := range fields {
		if _, ok := facts[key]; !ok {
			return fmt.Errorf("missing installed field %s", key)
		}
	}
	if facts["processProfile"] != "sdp-five-phase/0.1" || facts["managementProfile"] != "sdp-project-management/0.1" {
		return fmt.Errorf("unsupported installed process profile")
	}
	digest, ok := facts["configurationDigest"].(string)
	if !ok || !regexp.MustCompile("^[a-f0-9]{64}$").MatchString(digest) {
		return fmt.Errorf("invalid configuration digest")
	}
	for _, key := range []string{"toolkitVersion", "frameworkVersion", "agentsContractVersion", "installerVersion", "toolkitInstalledAt"} {
		s, ok := facts[key].(string)
		if !ok || s == "" {
			return fmt.Errorf("invalid installed field %s", key)
		}
	}
	if _, ok := facts["skills"].(map[string]any); !ok {
		return fmt.Errorf("invalid installed skills")
	}
	if _, ok := facts["capabilities"].([]any); !ok {
		return fmt.Errorf("invalid installed capabilities")
	}
	return nil
}

func pendingInstallations(area string) ([]string, error) {
	root, err := resolvePath(area, ".sdp-operations")
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(root)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	pending := []string{}
	id := regexp.MustCompile("^install-[a-f0-9]{24}$")
	for _, entry := range entries {
		if !id.MatchString(entry.Name()) {
			if entry.IsDir() || (entry.Name() != "install.lock" && !regexp.MustCompile("^temp-[a-f0-9]{64}$").MatchString(entry.Name())) {
				return nil, fmt.Errorf("unknown installation operation entry %s", entry.Name())
			}
			continue
		}
		if !entry.IsDir() || entry.Type()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("invalid installation operation directory")
		}
		path, err := resolvePath(area, filepath.Join(".sdp-operations", entry.Name(), "journal.json"))
		if err != nil {
			return nil, err
		}
		b, err := boundedFileLimit(path, 64<<20)
		if os.IsNotExist(err) {
			// Preparation may leave an empty directory before publishing its journal.
			contents, e := os.ReadDir(filepath.Dir(path))
			if e == nil && len(contents) == 0 {
				continue
			}
		}
		if err != nil {
			return nil, err
		}
		var journal struct {
			Schema string `json:"schemaVersion"`
			ID     string `json:"operationId"`
			Status string `json:"status"`
		}
		if err = json.Unmarshal(b, &journal); err != nil {
			return nil, err
		}
		if journal.Schema != "2.0" || journal.ID != entry.Name() {
			return nil, fmt.Errorf("invalid installation journal identity")
		}
		switch journal.Status {
		case "active", "failed":
			pending = append(pending, journal.ID)
		case "completed":
		default:
			return nil, fmt.Errorf("unsupported installation status")
		}
	}
	return pending, nil
}
