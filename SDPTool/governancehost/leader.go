package governancehost

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

// LeaderObservation requires the original thread, exact unique input and full
// history. Empty/idle history is never permission to resubmit an uncertain send.
func LeaderObservation(t Thread, d governance.LeaderDispatch, root, version string) (turnID, status, answer string, err error) {
	if version != d.Config.HostVersion || t.CLIVersion != version || t.ID != d.ThreadID || !filepath.IsAbs(t.Cwd) || filepath.Clean(t.Cwd) != filepath.Clean(root) {
		return "", "", "", fmt.Errorf("reconciliation_required: leader host/thread/project mapping differs")
	}
	turnID, _, err = correlatedInputs(t, d.TurnID, []governance.OwnerInput{{ID: d.ID, Payload: d.Payload, Skills: d.Config.Skills}})
	if err != nil {
		return
	}
	matches := 0
	for _, turn := range t.Turns {
		if turn.ID == turnID {
			matches++
		}
	}
	if matches != 1 {
		return "", "", "", fmt.Errorf("reconciliation_required: ambiguous leader turn history")
	}
	events, status, err := History(t, turnID, "leader-inspection", time.Now().UTC())
	if err != nil {
		return "", "", "", err
	}
	parts := []string{}
	for _, event := range events {
		if event.Observation.Kind == "agentMessage" {
			parts = append(parts, event.Observation.Text)
		}
	}
	answer = strings.Join(parts, "\n\n")
	if len(answer) > 64<<10 {
		return "", "", "", fmt.Errorf("limit: visible leader answer exceeds 64 KiB")
	}
	return turnID, status, answer, nil
}
