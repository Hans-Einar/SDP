package governancecontroller

import (
	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
	"os"
	"path/filepath"
	"testing"
)

func fixture(t *testing.T) (*governance.Service, governance.Selection, governance.StartRequest) {
	t.Helper()
	base := t.TempDir()
	root := filepath.Join(base, "project")
	area := filepath.Join(root, "SDP")
	state := filepath.Join(base, "state")
	files := map[string]string{
		"SDP/ProjectManagement/Ledger.ndjson": `{"schemaVersion":"1.0","eventId":"E1","eventType":"x-kanban:created","subjectId":"KB-TEST-001","payload":{"schemaVersion":"0.2","previousEventId":null,"to":"active","toPath":"active/card.md"}}
{"schemaVersion":"1.0","eventId":"E2","eventType":"x-management:started","subjectId":"MAINT-TEST-0001","payload":{"schemaVersion":"0.2","previousEventId":null,"to":"active","toPath":"Maintenance/Plan.md"}}
`,
		"AGENTS.md":                     "Preserve unrelated files.\n",
		"SDP/SDP-project.manifest.yaml": "schemaVersion: '1.0'\nproject:\n  name: trial\n",
		"SDP/KanBan/active/card.md":     "# Work\n\n| Field | Value |\n| --- | --- |\n| id | KB-TEST-001 |\n| PlanId | MAINT-TEST-0001 |\n| CardState | in-progress |\n",
		"SDP/Maintenance/Plan.md":       "# Plan\n\n| Field | Value |\n| --- | --- |\n| id | MAINT-TEST-0001 |\n| state | active |\n| PlanType | MaintenancePlan |\n",
		"SDP/Sessions/session.md":       "# Goal\n\n| State | Step |\n| --- | --- |\n| next | S1 |\n\n| Field | Value |\n| --- | --- |\n| Session reference | SESSION-TEST-0001 |\n| Status | active |\n",
		"SDP/authority.md":              "Owner selects the bounded fixture change.\n",
		"unrelated.txt":                 "keep this dirty input\n",
	}
	for p, b := range files {
		path := filepath.Join(root, p)
		if e := os.MkdirAll(filepath.Dir(path), 0700); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(path, []byte(b), 0600); e != nil {
			t.Fatal(e)
		}
	}
	s, e := governance.New(root, area, "trial", state)
	if e != nil {
		t.Fatal(e)
	}
	sel := governance.Selection{Card: "SDP/KanBan/active/card.md", Plan: "SDP/Maintenance/Plan.md", Session: "SDP/Sessions/session.md", Authority: "SDP/authority.md", Step: "S1"}
	c, e := s.Context(sel)
	if e != nil {
		t.Fatal(e)
	}
	if len(c.Conditions) != 0 {
		t.Fatal(c.Conditions)
	}
	op := governancehost.UUID()
	req := governancehost.UUID()
	q := governance.StartRequest{Schema: governance.Schema, ProjectRef: "trial", OperationID: op, RequestID: req, InputRevision: c.Revision, RoutineDigest: governance.Maintenance().Digest(), Selection: sel, Assignment: governance.Assignment{Outcome: "fix fixture", WritablePaths: []string{"target.txt"}, ExcludedPaths: []string{"unrelated.txt"}}}
	return s, sel, q
}
func bind(t *testing.T, s *governance.Service, q *governance.StartRequest) {
	t.Helper()
	b, e := s.Bind(governance.LocalController())
	if e != nil {
		t.Fatal(e)
	}
	q.CheckoutID = b.CheckoutID
}
