package sdptool

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func governanceService(selected, state string) (*governance.Service, error) {
	p, e := Discover(selected)
	if e != nil {
		return nil, e
	}
	if state == "" {
		base := os.Getenv("XDG_STATE_HOME")
		if base == "" {
			home, e := os.UserHomeDir()
			if e != nil {
				return nil, e
			}
			base = filepath.Join(home, ".local", "state")
		}
		if !filepath.IsAbs(base) {
			return nil, fmt.Errorf("XDG_STATE_HOME must be absolute")
		}
		state = filepath.Join(base, "sdptool", "project-governance", fmt.Sprintf("%x", sha256.Sum256([]byte(p.Root))))
	}
	service, e := governance.New(p.Root, p.Area, p.Inventory.ProjectID, state)
	if e != nil {
		return nil, e
	}
	return service, nil
}
