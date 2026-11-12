package activities

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestProjectionCausalTimelineAndSources(t *testing.T) {
	area := t.TempDir()
	write := func(n, b string) {
		t.Helper()
		f := filepath.Join(area, n)
		if e := os.MkdirAll(filepath.Dir(f), 0755); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(f, []byte(b), 0644); e != nil {
			t.Fatal(e)
		}
	}
	write("Activities/Activity-0001--Need/01-Mandate/Mandate.md", "# Useful need\n")
	write("Sessions/old.md", "# Earlier discussion\n")
	write("Sessions/unassigned.md", "# Unassigned discussion\n")
	write("Activities/Activity-0001--Need/04-Design/Design.md", "# Useful design\n")
	commit := "0123456789012345678901234567890123456789"
	planned := "planned"
	a := Activity{Path: "Activities/Activity-0001--Need", Descriptor: Descriptor{ID: "ACTIVITY-TEST-0001", Title: "Need", Mandate: "Activities/Activity-0001--Need/01-Mandate/Mandate.md", References: []Reference{{Role: "session", Path: "Sessions/old.md"}}}, State: "active", Events: []Event{{EventType: "x-activity:created", OccurredAt: "2026-11-12T12:00:00Z", Commit: &commit, Payload: Payload{To: "planned", Reason: "Selected need"}}, {EventType: "x-activity:transitioned", OccurredAt: "2026-11-12T11:00:00Z", Commit: &commit, Payload: Payload{From: &planned, To: "active", Reason: "Begin"}}}}
	p, e := Project(area, Catalog{HistoryRevision: "source", Activities: []Activity{a}}, time.Date(2026, 11, 12, 13, 0, 0, 0, time.UTC))
	if e != nil {
		t.Fatal(e)
	}
	if len(p.Unassigned) != 1 || p.Unassigned[0].Label != "Unassigned discussion" {
		t.Fatalf("unassigned %+v", p.Unassigned)
	}
	page, e := p.Render("activity/ACTIVITY-TEST-0001/timeline")
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(page.Markdown, "clock anomaly") || !strings.Contains(page.Markdown, "Calendar time is not effort") {
		t.Fatal(page.Markdown)
	}
	if strings.Count(page.Markdown, "Immutable commit") != 1 {
		t.Fatal("duplicate commit marker")
	}
	phase, e := p.Render("activity/ACTIVITY-TEST-0001/design")
	if e != nil || !strings.Contains(phase.Markdown, "Useful design") {
		t.Fatal(phase, e)
	}
	old := p.Revision
	write("Activities/Activity-0001--Need/04-Design/Design.md", "# Revised design\n")
	p, e = Project(area, p.Catalog, p.AsOf)
	if e != nil || old == p.Revision {
		t.Fatal("source edits must invalidate targets", e)
	}
	if _, e = p.Render("activity/ACTIVITY-TEST-0001/bogus"); e == nil {
		t.Fatal("invalid selection accepted")
	}
}
func TestProjectionUnknownTimeAndEmptyGroups(t *testing.T) {
	a := Activity{Descriptor: Descriptor{ID: "ACTIVITY-TEST-0001", Title: "Need"}, Events: []Event{{EventType: "x-activity:created", Payload: Payload{To: "planned"}}}}
	p := Projection{Entries: []Entry{{Activity: a}}, AsOf: time.Now()}
	page, e := p.Render("overview")
	if e != nil || !strings.Contains(page.Markdown, "unknown timestamp") {
		t.Fatal(page, e)
	}
	page, e = p.Render("state/completed")
	if e != nil || !strings.Contains(page.Markdown, "No Activities") {
		t.Fatal(page, e)
	}
}
func TestProjectionRejectsEscapingSource(t *testing.T) {
	area := t.TempDir()
	p := Projection{area: area}
	if _, e := p.source("../elsewhere"); e == nil {
		t.Fatal("escape accepted")
	}
	outside := t.TempDir()
	if e := os.Symlink(outside, filepath.Join(area, "Sessions")); e != nil {
		t.Fatal(e)
	}
	if e := p.walk("Sessions", func(string) {}); e == nil {
		t.Fatal("external symlink accepted")
	}
}

func TestProjectionExplicitManagementAttribution(t *testing.T) {
	area := t.TempDir()
	if e := os.MkdirAll(filepath.Join(area, "ProjectManagement"), 0755); e != nil {
		t.Fatal(e)
	}
	history := "{\"eventId\":\"linked\",\"eventType\":\"card-updated\",\"subjectId\":\"CARD-TEST-1\",\"occurredAt\":\"2026-11-12T12:00:00Z\",\"payload\":{\"reason\":\"Selected card evidence\"}}\n{\"eventId\":\"unrelated\",\"eventType\":\"card-updated\",\"subjectId\":\"CARD-TEST-2\",\"occurredAt\":\"2026-11-12T12:00:00Z\",\"payload\":{\"reason\":\"Unrelated at same timestamp\"}}\n"
	if e := os.WriteFile(filepath.Join(area, "ProjectManagement", "Ledger.ndjson"), []byte(history), 0644); e != nil {
		t.Fatal(e)
	}
	a := Activity{Descriptor: Descriptor{ID: "ACTIVITY-TEST-0001", Title: "Need", SourceCards: []string{"CARD-TEST-1"}}}
	p := Projection{area: area, Entries: []Entry{{Activity: a}}, AsOf: time.Now()}
	page, e := p.Render("activity/ACTIVITY-TEST-0001/timeline")
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(page.Markdown, "Selected card evidence") || strings.Contains(page.Markdown, "Unrelated at same timestamp") {
		t.Fatal(page.Markdown)
	}
}
