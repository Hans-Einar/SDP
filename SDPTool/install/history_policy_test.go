package install

import (
	"encoding/json"
	"testing"
)

func TestCanonicalRecordedTypingDoesNotRelaxInstallationPolicy(t *testing.T) {
	first := map[string]any{"schemaVersion": "1.0", "eventId": "EVT-PM-TEST-000001", "eventType": "x-management:created", "subjectId": "MAINT-TEST-0001", "occurredAt": "2026-10-24T00:00:00Z", "actor": "fixture", "commit": nil, "payload": map[string]any{"schemaVersion": "0.1", "kind": "Maintenance", "previousEventId": nil, "from": nil, "to": "active", "fromPath": nil, "toPath": "Maintenance/Plan.md", "reason": "initial untyped record", "links": []string{}}}
	second := map[string]any{"schemaVersion": "1.0", "eventId": "EVT-PM-TEST-000002", "eventType": "x-management:updated", "subjectId": "MAINT-TEST-0001", "occurredAt": "2026-10-24T01:00:00Z", "actor": "fixture", "commit": nil, "payload": map[string]any{"schemaVersion": "0.2", "kind": "Maintenance", "planType": "MaintenancePlan", "previousEventId": "EVT-PM-TEST-000001", "from": "active", "to": "active", "fromPath": "Maintenance/Plan.md", "toPath": "Maintenance/Plan.md", "reason": "recorded typed-plan adoption", "links": []string{}}}
	a, _ := json.Marshal(first)
	b, _ := json.Marshal(second)
	history := append(append(append(a, '\n'), b...), '\n')
	files := map[string][]byte{"SDP/Maintenance/Plan.md": []byte("# Plan\n\n| Field | Value |\n| --- | --- |\n| id | MAINT-TEST-0001 |\n| state | active |\n| PlanType | MaintenancePlan |\n| BranchPolicy | current |\n| CommitPolicy | milestone |\n")}
	if e := ValidateManagementHistory(history, files); e != nil {
		t.Fatal(e)
	}
	if e := validateHistory(history, files); e == nil {
		t.Fatal("installation policy was silently relaxed")
	}
}
