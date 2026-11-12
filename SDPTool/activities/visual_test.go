package activities

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

func TestVisualCalendarEscapingAndFilters(t *testing.T) {
	start := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	a := Activity{Descriptor: Descriptor{ID: "ACTIVITY-TEST-0001", Title: "Need <script>&"}, Events: []Event{{EventType: "x-activity:created", OccurredAt: start.Format(time.RFC3339), Payload: Payload{To: "planned"}}}}
	p := Projection{Entries: []Entry{{Activity: a}}, AsOf: start.Add(time.Hour), area: t.TempDir()}
	page, e := p.Render("overview?markers=ledger&zoom=week")
	if e != nil {
		t.Fatal(e)
	}
	svg := page.Assets["timeline.svg"]
	if len(svg) == 0 || strings.Contains(string(svg), "<script>") {
		t.Fatal("unsafe or missing visual")
	}
	var doc any
	if e = xml.Unmarshal(svg, &doc); e != nil {
		t.Fatal("invalid SVG", e)
	}
	if !strings.Contains(string(svg), "ledger events") || !strings.Contains(string(svg), "calendar time") {
		t.Fatal(string(svg))
	}
	no, e := p.Render("overview?markers=none")
	if e != nil {
		t.Fatal(e)
	}
	if strings.Contains(string(no.Assets["timeline.svg"]), "ledger events in") {
		t.Fatal("filter ignored")
	}
	if _, e = p.Render("overview?markers=all&markers=none"); e == nil {
		t.Fatal("ambiguous controls accepted")
	}
}
