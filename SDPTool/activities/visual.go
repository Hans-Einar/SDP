package activities

import (
	"fmt"
	"html"
	"sort"
	"strings"
	"time"
)

var stateColors = map[string]string{"planned": "#94a3b8", "active": "#16a34a", "on-hold": "#f59e0b", "completed": "#2563eb", "canceled": "#64748b"}

// timelineAssets uses explicit causal events only. Invalid intervals are omitted
// from the proportional chart and explained in the accompanying readable history.
func (p Projection) timelineAssets(only string) map[string][]byte {
	assets := map[string][]byte{}
	rows := []Activity{}
	first, last := p.AsOf, p.AsOf
	for _, entry := range p.Entries {
		a := entry.Activity
		if only != "" && a.Descriptor.ID != only {
			continue
		}
		rows = append(rows, a)
		related, _ := p.relatedEvents(a)
		for _, e := range append(append([]Event{}, a.Events...), related...) {
			if e.Commit != nil {
				if _, when, _, err := p.commit(*e.Commit); err == nil {
					if t, err := time.Parse(time.RFC3339Nano, when); err == nil {
						if t.Before(first) {
							first = t
						}
						if t.After(last) {
							last = t
						}
					}
				}
			}
			if t, err := time.Parse(time.RFC3339Nano, e.OccurredAt); err == nil {
				if t.Before(first) {
					first = t
				}
				if t.After(last) {
					last = t
				}
			}
		}
	}
	if last.Sub(first) < time.Hour {
		last = first.Add(time.Hour)
	}
	width := 1100.0
	left := 250.0
	plot := width - left - 30
	x := func(t time.Time) float64 { return left + float64(t.Sub(first))/float64(last.Sub(first))*plot }
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1100" height="%d" viewBox="0 0 1100 %d"><rect width="100%%" height="100%%" fill="white"/><g font-family="sans-serif" font-size="12" fill="#172033">`, 120+len(rows)*64, 120+len(rows)*64)
	b.WriteString(`<text x="12" y="20" font-size="16">Activity lifecycle — calendar time, not effort (UTC)</text>`)
	for i := 0; i <= 4; i++ {
		t := first.Add(time.Duration(float64(last.Sub(first)) * float64(i) / 4))
		fmt.Fprintf(&b, `<text x="%.1f" y="43" text-anchor="middle">%s</text><line x1="%.1f" y1="50" x2="%.1f" y2="%d" stroke="#e2e8f0"/>`, x(t), t.UTC().Format("Jan 02 15:04"), x(t), x(t), 60+len(rows)*64)
	}
	for r, a := range rows {
		y := 65 + r*64
		title := []rune(a.Descriptor.Title)
		if len(title) > 32 {
			title = append(title[:31], '…')
		}
		fmt.Fprintf(&b, `<text x="12" y="%d"><title>%s</title>%s</text>`, y+17, html.EscapeString(a.Descriptor.Title), html.EscapeString(string(title)))
		transitions := []Event{}
		for _, e := range a.Events {
			if e.EventType == "x-activity:created" || e.EventType == "x-activity:transitioned" {
				transitions = append(transitions, e)
			}
		}
		for i, e := range transitions {
			start, err := time.Parse(time.RFC3339Nano, e.OccurredAt)
			if err != nil {
				continue
			}
			end := p.AsOf
			if i+1 < len(transitions) {
				end, _ = time.Parse(time.RFC3339Nano, transitions[i+1].OccurredAt)
			}
			color := stateColors[e.Payload.To]
			if e.Payload.To == "completed" || e.Payload.To == "canceled" {
				fmt.Fprintf(&b, `<circle cx="%.1f" cy="%d" r="5" fill="%s"><title>%s %s</title></circle>`, x(start), y+15, color, html.EscapeString(e.Payload.To), start.UTC().Format(time.RFC3339))
				continue
			}
			if end.IsZero() || end.Before(start) {
				continue
			}
			w := x(end) - x(start)
			if w < 1 {
				w = 1
			}
			fmt.Fprintf(&b, `<rect x="%.1f" y="%d" width="%.1f" height="24" rx="3" fill="%s"><title>%s: %s to %s</title></rect>`, x(start), y+3, w, color, html.EscapeString(e.Payload.To), start.UTC().Format(time.RFC3339), end.UTC().Format(time.RFC3339))
		}
		related, _ := p.relatedEvents(a)
		events := append(append([]Event{}, a.Events...), related...)
		markers := map[string]int{}
		seen := map[string]bool{}
		bucket := func(t time.Time) time.Time {
			t = t.UTC()
			switch p.zoom {
			case "week":
				return t.Truncate(7 * 24 * time.Hour)
			case "month":
				return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
			default:
				return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			}
		}
		for _, e := range events {
			if p.markerFilter != "none" && p.markerFilter != "commits" {
				if t, err := time.Parse(time.RFC3339Nano, e.OccurredAt); err == nil {
					markers["ledger|"+bucket(t).Format(time.RFC3339)]++
				}
			}
			if e.Commit == nil || *e.Commit == "" || seen[*e.Commit] {
				continue
			}
			seen[*e.Commit] = true
			subject, when, repo, err := p.commit(*e.Commit)
			content := "# Associated commit\n\nRepository: " + text(repo) + "\n\n" + text(subject) + "\n\nCommitter time: " + stamp(when) + "\n\nImmutable identity: `" + *e.Commit + "`\n"
			if err != nil {
				content = "# Associated commit\n\nRepository metadata unavailable; association is retained in the canonical source event.\n\nImmutable identity: `" + *e.Commit + "`\n"
			}
			// Event IDs are validated by replay; defensively reject path-like Git input.
			if hashCommit(*e.Commit) {
				assets["commits/"+*e.Commit+".md"] = []byte(content)
			}
			if err == nil && p.markerFilter != "none" && p.markerFilter != "ledger" {
				if t, err := time.Parse(time.RFC3339Nano, when); err == nil {
					markers["commit|"+bucket(t).Format(time.RFC3339)]++
				}
			}
		}
		keys := []string{}
		for key := range markers {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			v := strings.SplitN(key, "|", 2)
			t, _ := time.Parse(time.RFC3339, v[1])
			xx := x(t)
			if xx < left {
				xx = left
			}
			if xx > width-30 {
				xx = width - 30
			}
			color := "#0284c7"
			yy := y
			if v[0] == "commit" {
				color = "#9333ea"
				yy += 5
			}
			fmt.Fprintf(&b, `<line x1="%.1f" y1="%d" x2="%.1f" y2="%d" stroke="%s" stroke-width="2"><title>%d %s events in %s UTC bucket; complete sources in text</title></line>`, xx, yy, xx, yy+35, color, markers[key], v[0], v[1])
		}
	}
	y := 90 + len(rows)*64
	fmt.Fprintf(&b, `<text x="12" y="%d">Planned: grey · Active: green · On hold: amber · Completed: blue · Ledger: cyan marks · Commits: purple marks</text></g></svg>`, y)
	assets["timeline.svg"] = []byte(b.String())
	return assets
}

func hashCommit(s string) bool {
	if len(s) != 40 && len(s) != 64 {
		return false
	}
	for _, r := range s {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F') {
			return false
		}
	}
	return true
}
