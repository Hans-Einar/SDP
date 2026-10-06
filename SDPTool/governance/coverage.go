package governance

import (
	"bytes"
	"encoding/json"
	"strings"
)

// ledgerCoverage checks the selected subjects using the canonical management
// envelope and per-subject predecessor convention. Full repository validation
// remains with ProjectManagement/validate.py; this never writes a projection.
func ledgerCoverage(b []byte, card, plan, cardPath, planPath string) error {
	type event struct {
		Schema  string `json:"schemaVersion"`
		ID      string `json:"eventId"`
		Type    string `json:"eventType"`
		Subject string `json:"subjectId"`
		Payload struct {
			Schema   string  `json:"schemaVersion"`
			Previous *string `json:"previousEventId"`
			To       string  `json:"to"`
			Path     string  `json:"toPath"`
		} `json:"payload"`
	}
	latest := map[string]event{}
	seen := map[string]bool{}
	for _, line := range bytes.Split(b, []byte("\n")) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		// Historical extension fields are allowed, duplicate keys are not.
		var raw map[string]any
		if e := decode(line, &raw); e != nil {
			return fail("conflict", "invalid management ledger JSON")
		}
		var ev event
		if e := json.Unmarshal(line, &ev); e != nil {
			return fail("conflict", "invalid management event")
		}
		if ev.Subject != card && ev.Subject != plan {
			continue
		}
		prefix := "x-management:"
		if ev.Subject == card {
			prefix = "x-kanban:"
		}
		if !strings.HasPrefix(ev.Type, prefix) {
			continue
		}
		if ev.ID == "" || seen[ev.ID] || ev.Schema != "1.0" || (ev.Payload.Schema != "0.1" && ev.Payload.Schema != "0.2") {
			return fail("conflict", "invalid selected management event identity/profile")
		}
		previous := ""
		if ev.Payload.Previous != nil {
			previous = *ev.Payload.Previous
		}
		if previous != latest[ev.Subject].ID {
			return fail("conflict", "broken selected management chain")
		}
		latest[ev.Subject] = ev
		seen[ev.ID] = true
	}
	for subject, path := range map[string]string{card: cardPath, plan: planPath} {
		ev, ok := latest[subject]
		if !ok {
			return fail("missing_coverage", "selected subject missing from canonical management ledger")
		}
		if ev.Payload.Path != path || ev.Payload.To != "active" {
			return fail("conflict", "selected record and canonical management state disagree")
		}
	}
	return nil
}
