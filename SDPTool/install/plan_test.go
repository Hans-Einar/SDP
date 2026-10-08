package install

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func artifact(t *testing.T, d Descriptor) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "release.json")
	b, _ := Canonical(d)
	if e := os.WriteFile(p, b, 0600); e != nil {
		t.Fatal(e)
	}
	return p
}
func put(t *testing.T, root, p string, b []byte) {
	t.Helper()
	q := filepath.Join(root, p)
	if e := os.MkdirAll(filepath.Dir(q), 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(q, b, 0600); e != nil {
		t.Fatal(e)
	}
}
func preview(t *testing.T, root string, d Descriptor) Plan {
	t.Helper()
	p, e := Preview(Options{Root: root, Operation: "install", Artifact: artifact(t, d), AllowUnreleased: true})
	if e != nil {
		t.Fatal(e)
	}
	return p
}
func adoption(t *testing.T, root string, input Input, moves []Move, refresh []string) string {
	t.Helper()
	tree, e := Inspect(root, targetPaths(fixture()))
	if e != nil {
		t.Fatal(e)
	}
	a := Adoption{SchemaVersion: AdoptionSchema, ProjectRoot: root, Baseline: "manual", TargetDigest: input.SHA256, Inventory: tree.Snapshot, Moves: moves, RefreshManaged: refresh, AllowReferenceWarnings: true}
	b, _ := Canonical(a)
	p := filepath.Join(t.TempDir(), "adopt.yaml")
	if e = os.WriteFile(p, b, 0600); e != nil {
		t.Fatal(e)
	}
	return p
}
func TestCleanPreviewDeterministicNoWrites(t *testing.T) {
	root := t.TempDir()
	d := fixture()
	a := artifact(t, d)
	o := Options{Root: root, Operation: "install", Artifact: a, AllowUnreleased: true}
	p, e := Preview(o)
	if e != nil || !p.CanApply {
		t.Fatal(p, e)
	}
	q, e := Preview(o)
	if e != nil || q.PlanDigest != p.PlanDigest {
		t.Fatal(q, e)
	}
	entries, _ := os.ReadDir(root)
	if len(entries) != 0 {
		t.Fatal("preview wrote project")
	}
	dest := filepath.Join(t.TempDir(), "plan.json")
	if e = SavePlan(p, dest); e != nil {
		t.Fatal(e)
	}
	if e = SavePlan(p, dest); e == nil {
		t.Fatal("overwrote plan")
	}
	if e = SavePlan(p, filepath.Join(root, "plan.json")); e == nil {
		t.Fatal("plan inside inventory")
	}
	if _, e = LoadPlan(dest); e != nil {
		t.Fatal(e)
	}
}
func TestUnknownAndManualAdoption(t *testing.T) {
	root := t.TempDir()
	put(t, root, "AGENTS.md", []byte("project instructions\n"))
	put(t, root, "SDP/Agents/KanBan/board.json", []byte(`{"schemaVersion":"0.1","projectId":"XFMD","ledger":"Ledger.ndjson"}`))
	prefix := historyFixture()
	put(t, root, "SDP/Agents/KanBan/Ledger.ndjson", prefix)
	put(t, root, "SDP/Agents/KanBan/backlog/#1.md", cardFixture("[peer](%231.md)\n"))
	put(t, root, "README.md", []byte("[card](SDP/Agents/KanBan/backlog/%231.md)\n"))
	a := artifact(t, fixture())
	p, e := Preview(Options{root, "upgrade", a, "", "", true})
	if e != nil || p.CanApply {
		t.Fatal(p, e)
	}
	in, _ := LocalInput(a, true)
	m := adoption(t, root, in, []Move{{"SDP/Agents/KanBan", "SDP/KanBan"}}, []string{"AGENTS.md"})
	p, e = Preview(Options{root, "upgrade", a, "", m, true})
	if e != nil || !p.CanApply {
		t.Fatal(p.Conflicts, e)
	}
	actions := map[string]Action{}
	deleting := false
	for _, a := range p.Actions {
		if a.Action == "delete" {
			deleting = true
		} else if deleting {
			t.Fatal("write after delete")
		}
		actions[a.Path] = a
	}
	if !bytes.Equal(actions["SDP/ProjectManagement/Ledger.ndjson"].Content, prefix) {
		t.Fatal("history not preserved")
	}
	if string(actions["AGENTS-project.md"].Content) != "project instructions\n" {
		t.Fatal("instructions lost")
	}
	if string(actions["README.md"].Content) != "[card](SDP/KanBan/backlog/%231.md)\n" {
		t.Fatal("incoming link not rebased")
	}
	put(t, root, "SDP/new.md", []byte("drift"))
	p, e = Preview(Options{root, "upgrade", a, "", m, true})
	if e == nil && p.CanApply {
		t.Fatal("adoption drift accepted")
	}
}
func TestCLIConflictAndPreview(t *testing.T) {
	root := t.TempDir()
	a := artifact(t, fixture())
	var out, errs bytes.Buffer
	code := Run(context.Background(), root, "install", []string{"--artifact", a, "--allow-unreleased", "--json"}, &out, &errs)
	if code != 0 {
		t.Fatal(code, errs.String(), out.String())
	}
	var v map[string]any
	if json.Unmarshal(out.Bytes(), &v) != nil || v["schemaVersion"] != Protocol {
		t.Fatal(out.String())
	}
	out.Reset()
	code = Run(context.Background(), root, "install", []string{"--artifact", a, "--json"}, &out, &errs)
	if code != 2 {
		t.Fatal(code, out.String())
	}
}
func TestMarkdownRebaseForms(t *testing.T) {
	source := "SDP/Agents/KanBan/a.md"
	dest := "SDP/KanBan/a.md"
	moves := map[string]string{source: dest}
	b := []byte("[a](<b%20c.md> \"title\") [b](d(e).md)\n[ref]: b%20c.md\n[remote](https://example.com/x)\n")
	got := rebaseMarkdown(b, source, dest, moves, []Move{{"SDP/Agents/KanBan", "SDP/KanBan"}})
	if !bytes.Contains(got, []byte("d%28e%29.md")) {
		t.Fatal("balanced link not encoded:", string(got))
	}

	if !bytes.Contains(got, []byte("https://example.com/x")) {
		t.Fatal(string(got))
	}
}

func TestKnownUpgradeAndLocalEdits(t *testing.T) {
	root := t.TempDir()
	d := fixture()
	p := preview(t, root, d)
	for _, a := range p.Actions {
		if a.Action == "write" {
			put(t, root, a.Path, a.Content)
		}
	}
	r := Receipt{ReceiptSchema, d.Release, p.Release.SHA256, d.SourceCommit, d.ProcessProfile, d.ManagementProfile, d.Capabilities, "local-development", "install-0123456789abcdef01234567", "2026-09-27T00:00:00Z"}
	b, _ := Canonical(r)
	put(t, root, ReceiptPath, b)
	same, e := Preview(Options{root, "upgrade", p.Release.Path, "", "", true})
	if e != nil || !same.NoChange || len(same.Actions) != 0 {
		t.Fatal(same.Actions, e)
	}
	d.Release = "dev-gip-2"
	d.UpgradesFrom = []string{p.Release.SHA256}
	d.Files[0].Content = []byte("new instructions\n")
	h := Hash(d.Files[0].Content)
	d.Files[0].SHA256 = &h
	next := artifact(t, d)
	q, e := Preview(Options{root, "upgrade", next, p.Release.Path, "", true})
	if e != nil || !q.CanApply || len(q.Actions) != 1 {
		t.Fatal(q, e)
	}
	put(t, root, "AGENTS.md", []byte("local edit"))
	q, e = Preview(Options{root, "upgrade", next, p.Release.Path, "", true})
	if e != nil || q.CanApply {
		t.Fatal(q, e)
	}
}

func historyFixture() []byte {
	b, _ := Canonical(map[string]any{"schemaVersion": "1.0", "eventId": "EVT-KB-XFMD-000001", "eventType": "x-kanban:created", "occurredAt": "2026-09-27T00:00:00Z", "actor": "test", "commit": nil, "subjectId": "KB-XFMD-001", "payload": map[string]any{"schemaVersion": "0.1", "projectId": "XFMD", "previousEventId": nil, "from": nil, "to": "backlog", "fromPath": nil, "toPath": "backlog/#1.md", "reason": "Test baseline", "links": []string{}}})
	return b
}
func cardFixture(body string) []byte {
	return []byte("# Fixture\n\n| Field | Value |\n| --- | --- |\n| id | KB-XFMD-001 |\n| CardState | backlog |\n\n" + body)
}
func TestHistoryRejectsBrokenChainAndMissingCard(t *testing.T) {
	files := map[string][]byte{"SDP/KanBan/backlog/#1.md": cardFixture("body")}
	if e := validateHistory(historyFixture(), files); e != nil {
		t.Fatal(e)
	}
	bad := bytes.Replace(historyFixture(), []byte(`"previousEventId":null`), []byte(`"previousEventId":"missing"`), 1)
	if validateHistory(bad, files) == nil {
		t.Fatal("broken chain accepted")
	}
	if validateHistory(historyFixture(), map[string][]byte{}) == nil {
		t.Fatal("missing card accepted")
	}
}

func TestCleanRootInstructionsPreserved(t *testing.T) {
	root := t.TempDir()
	put(t, root, "AGENTS.md", []byte("owner instructions"))
	p := preview(t, root, fixture())
	if !p.CanApply {
		t.Fatal(p.Conflicts)
	}
	found := false
	for _, a := range p.Actions {
		if a.Path == "AGENTS-project.md" && string(a.Content) == "owner instructions" {
			found = true
		}
	}
	if !found {
		t.Fatal("instructions not preserved")
	}
	put(t, root, "AGENTS-project.md", []byte("different instructions"))
	p = preview(t, root, fixture())
	if p.CanApply {
		t.Fatal("preservation collision accepted")
	}
}
func TestAdoptionNestedMarkdownAncestorsStable(t *testing.T) {
	root := t.TempDir()
	put(t, root, "docs/deep/note.md", []byte("incoming link scope"))
	put(t, root, "SDP/README.md", []byte("manual"))
	in, _ := LocalInput(artifact(t, fixture()), true)
	m := adoption(t, root, in, []Move{}, []string{})
	p, e := Preview(Options{root, "upgrade", in.Path, "", m, true})
	if e != nil || !p.CanApply {
		t.Fatal(p.Conflicts, e)
	}
}

func TestInstallDoesNotCreateOrEditLegacyNavigation(t *testing.T) {
	for _, legacy := range []bool{false, true} {
		root := t.TempDir()
		original := []byte("historical, not valid JSON")
		if legacy {
			put(t, root, "SDP/navigation.json", original)
		}
		a := artifact(t, fixture())
		var p Plan
		var e error
		if legacy {
			in, _ := LocalInput(a, true)
			m := adoption(t, root, in, []Move{}, []string{})
			p, e = Preview(Options{root, "upgrade", a, "", m, true})
		} else {
			p, e = Preview(Options{root, "install", a, "", "", true})
		}
		if e != nil {
			t.Fatal(e)
		}
		if !p.CanApply {
			t.Fatal(p.Conflicts)
		}
		for _, action := range p.Actions {
			if action.Path == "SDP/navigation.json" {
				t.Fatal("legacy index must be untouched", action)
			}
		}
	}
}
