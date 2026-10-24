package governancehost

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

// ReviewObservation admits one exact fresh-context turn, not arbitrary review prose.
func ReviewObservation(t Thread, d governance.Review, root, version string) (turnID, status, answer string, err error) {
	if version != d.Spec.HostVersion || t.CLIVersion != version || t.ID != d.ThreadID || !filepath.IsAbs(t.Cwd) || filepath.Clean(t.Cwd) != filepath.Clean(root) || t.ParentThreadID != "" || len(t.Turns) != 1 {
		return "", "", "", fmt.Errorf("reconciliation_required: review host/thread/project or isolated history differs")
	}
	turnID, _, err = correlatedInputs(t, d.TurnID, []governance.OwnerInput{{ID: d.ID, Payload: d.Payload, Skills: d.Spec.Skills}})
	if err != nil {
		return
	}
	events, state, e := History(t, turnID, "review-inspection", time.Now().UTC())
	if e != nil {
		return "", "", "", e
	}
	status = state
	parts := []string{}
	for _, event := range events {
		if event.Observation.Kind == "agentMessage" && event.Observation.Phase == "final_answer" {
			parts = append(parts, event.Observation.Text)
		}
	}
	answer = strings.Join(parts, "\n\n")
	if len(answer) > 64<<10 {
		return "", "", "", fmt.Errorf("limit: final review exceeds 64 KiB")
	}
	if status == "running" {
		answer = ""
	}
	return
}
