package activities

// This file owns presentation classification. Hosts translate entries into their
// navigation protocol, but never infer lifecycle state from directory names.
import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Phase struct{ Key, Label, Directory string }

var Phases = []Phase{{"mandate", "Mandate", "01-Mandate"}, {"study", "Study", "02-Study"}, {"requirements", "Requirements", "03-Requirements"}, {"design", "Design", "04-Design"}, {"implementation", "Implementation", "05-Implementation"}, {"verification", "Verification", "06-Verification"}}
var States = []string{"planned", "active", "on-hold", "completed", "canceled"}

type Document struct{ Label, Path, Revision, Diagnostic string }
type Entry struct {
	Activity  Activity
	Documents map[string][]Document
}
type Page struct {
	Assets   map[string][]byte `json:"assets,omitempty"`
	Title    string            `json:"title"`
	Markdown string            `json:"markdown"`
	Revision string            `json:"revision"`
}
type Projection struct {
	Catalog      Catalog
	Entries      []Entry
	Unassigned   []Document
	Revision     string
	AsOf         time.Time
	markerFilter string
	zoom         string
	area         string
}

// Project indexes sources only. Markdown pages and timelines are rendered on selection.
func Project(area string, c Catalog, asOf time.Time) (Projection, error) {
	p := Projection{Catalog: c, AsOf: asOf.UTC(), area: area}
	assigned := map[string]bool{}
	for _, a := range c.Activities {
		e := Entry{Activity: a, Documents: map[string][]Document{}}
		seen := map[string]bool{}
		add := func(key, name string) {
			if seen[name] {
				return
			}
			seen[name] = true
			d := p.document(name)
			e.Documents[key] = append(e.Documents[key], d)
			if key == "sessions" {
				assigned[name] = true
			}
		}
		for _, ph := range Phases {
			dir := path.Join(a.Path, ph.Directory)
			if err := p.walk(dir, func(name string) { add(ph.Key, name) }); err != nil {
				return p, err
			}
		}
		if err := p.walk(path.Join(a.Path, "Sessions"), func(name string) { add("sessions", name) }); err != nil {
			return p, err
		}
		for _, r := range a.Descriptor.References {
			key := referencePhase(r.Role)
			if r.Availability == "unavailable" {
				e.Documents[key] = append(e.Documents[key], Document{Path: r.Path, Label: r.ID, Diagnostic: r.UnavailableReason})
				continue
			}
			add(key, r.Path)
		}
		for _, s := range a.Sessions {
			add("sessions", s.Path)
		}
		add("mandate", path.Join(a.Path, a.Descriptor.Mandate))
		for k := range e.Documents {
			sort.Slice(e.Documents[k], func(i, j int) bool { return e.Documents[k][i].Path < e.Documents[k][j].Path })
		}
		p.Entries = append(p.Entries, e)
	}
	sort.Slice(p.Entries, func(i, j int) bool {
		a, b := p.Entries[i].Activity.Descriptor, p.Entries[j].Activity.Descriptor
		if a.Title == b.Title {
			return a.ID < b.ID
		}
		return a.Title < b.Title
	})
	if err := p.walk("Sessions", func(name string) {
		if !assigned[name] {
			p.Unassigned = append(p.Unassigned, p.document(name))
		}
	}); err != nil {
		return p, err
	}
	b, _ := json.Marshal(struct {
		History    string
		Entries    []Entry
		Unassigned []Document
	}{c.HistoryRevision, p.Entries, p.Unassigned})
	sum := sha256.Sum256(b)
	p.Revision = hex.EncodeToString(sum[:])
	return p, nil
}
func referencePhase(role string) string {
	switch strings.ToLower(role) {
	case "session":
		return "sessions"
	case "mandate":
		return "mandate"
	case "study":
		return "study"
	case "requirement", "requirements":
		return "requirements"
	case "verification", "evidence":
		return "verification"
	case "plan", "implementation":
		return "implementation"
	default:
		return "design"
	}
}
func (p Projection) source(name string) (string, error) {
	if name == "" || path.IsAbs(name) || path.Clean(name) != name || strings.Contains(name, "\\") || name == ".." || strings.HasPrefix(name, "../") {
		return "", fmt.Errorf("invalid source path")
	}
	base, err := filepath.EvalSymlinks(p.area)
	if err != nil {
		return "", err
	}
	full, err := filepath.EvalSymlinks(filepath.Join(base, filepath.FromSlash(name)))
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(base, full)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("source escapes project area")
	}
	return full, nil
}
func (p Projection) walk(dir string, add func(string)) error {
	full, err := p.source(dir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	count := 0
	return filepath.WalkDir(full, func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("source symlink is not indexed: %s", dir)
		}
		if d.IsDir() {
			return nil
		}
		count++
		if count > 20000 {
			return fmt.Errorf("Activity source index exceeds 20000 files; select a smaller catalog")
		}
		rel, _ := filepath.Rel(full, name)
		add(path.Join(dir, filepath.ToSlash(rel)))
		return nil
	})
}
func (p Projection) document(name string) Document {
	d := Document{Path: name, Label: strings.TrimSuffix(path.Base(name), path.Ext(name))}
	full, e := p.source(name)
	if e != nil {
		d.Diagnostic = "Source unavailable: " + e.Error()
		return d
	}
	f, e := os.Open(full)
	if e != nil {
		d.Diagnostic = e.Error()
		return d
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, 32<<20+1))
	if len(b) > 32<<20 {
		d.Diagnostic = "Source exceeds 32 MiB reader limit"
		return d
	}
	if e != nil {
		d.Diagnostic = e.Error()
		return d
	}
	sum := sha256.Sum256(b)
	d.Revision = hex.EncodeToString(sum[:])
	if len(b) <= 1<<20 {
		for _, line := range strings.Split(string(b), "\n") {
			if strings.HasPrefix(line, "# ") {
				d.Label = strings.TrimSpace(strings.TrimPrefix(line, "# "))
				break
			}
		}
	}
	return d
}
func text(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "[", "\\[", "]", "\\]", "\n", " ").Replace(s)
}
func (p Projection) sourceLink(label, name string) string {
	if strings.HasPrefix(name, "git:commit/") && hashCommit(strings.TrimPrefix(name, "git:commit/")) {
		return "[" + text(label) + "](commits/" + strings.TrimPrefix(name, "git:commit/") + ".md)"
	}
	if _, e := p.source(name); e != nil {
		return text(label) + " (source unavailable)"
	}
	return "[" + text(label) + "](" + (&url.URL{Scheme: "file", Path: filepath.Join(p.area, filepath.FromSlash(name))}).String() + ")"
}
func (p Projection) docs(b *strings.Builder, ds []Document) {
	if len(ds) == 0 {
		b.WriteString("No recorded documents. Phase readiness has not been assessed.\n\n")
	}
	for _, d := range ds {
		fmt.Fprintf(b, "- %s", p.sourceLink(d.Label, d.Path))
		if d.Diagnostic != "" {
			fmt.Fprintf(b, " — %s", text(d.Diagnostic))
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
}
func (p Projection) Render(selection string) (Page, error) {
	u, err := url.Parse(selection)
	if err != nil || u.Scheme != "" || u.Host != "" || u.Fragment != "" {
		return Page{}, fmt.Errorf("invalid selection")
	}
	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return Page{}, err
	}
	for key, values := range q {
		if (key != "markers" && key != "zoom") || len(values) != 1 {
			return Page{}, fmt.Errorf("invalid timeline controls")
		}
	}
	p.markerFilter = q.Get("markers")
	p.zoom = q.Get("zoom")
	if p.markerFilter != "" && p.markerFilter != "all" && p.markerFilter != "ledger" && p.markerFilter != "commits" && p.markerFilter != "none" {
		return Page{}, fmt.Errorf("invalid marker filter")
	}
	if p.zoom != "" && p.zoom != "day" && p.zoom != "week" && p.zoom != "month" {
		return Page{}, fmt.Errorf("invalid timeline zoom")
	}
	selection = strings.Trim(u.Path, "/")
	if selection == "" {
		selection = "overview"
	}
	var b strings.Builder
	title := "Activities"
	if selection == "overview" || strings.HasPrefix(selection, "state/") {
		state := strings.TrimPrefix(selection, "state/")
		if selection != "overview" {
			valid := false
			for _, s := range States {
				valid = valid || s == state
			}
			if !valid {
				return Page{}, fmt.Errorf("unknown Activity state")
			}
			title = stateLabel(state)
		}
		fmt.Fprintf(&b, "# %s\n\n", title)
		n := 0
		for _, e := range p.Entries {
			a := e.Activity
			if selection != "overview" && a.State != state {
				continue
			}
			n++
			fmt.Fprintf(&b, "- **%s** — %s. %s\n", text(a.Descriptor.Title), stateLabel(a.State), p.sourceLink("Read mandate", path.Join(a.Path, a.Descriptor.Mandate)))
		}
		if n == 0 {
			b.WriteString("No Activities in this group.\n")
		}
		b.WriteString("\n")
		if selection == "overview" {
			p.timeline(&b, nil)
		}
	} else if selection == "unassigned" {
		title = "Unassigned sessions"
		b.WriteString("# Unassigned sessions\n\nExisting records remain in their original homes.\n\n")
		p.docs(&b, p.Unassigned)
	} else {
		parts := strings.Split(selection, "/")
		if len(parts) < 2 || len(parts) > 3 || parts[0] != "activity" {
			return Page{}, fmt.Errorf("invalid Activity selection")
		}
		var entry *Entry
		for i := range p.Entries {
			if p.Entries[i].Activity.Descriptor.ID == parts[1] {
				entry = &p.Entries[i]
				break
			}
		}
		if entry == nil {
			return Page{}, fmt.Errorf("Activity selection not found")
		}
		a := entry.Activity
		title = a.Descriptor.Title
		view := ""
		if len(parts) == 3 {
			view = parts[2]
		}
		fmt.Fprintf(&b, "# %s\n\n", text(title))
		switch view {
		case "":
			fmt.Fprintf(&b, "State: **%s**. Owner: %s.\n\n%s\n\n", stateLabel(a.State), text(a.Descriptor.Owner), p.sourceLink("Need and desired outcome", path.Join(a.Path, a.Descriptor.Mandate)))
			b.WriteString("Phase readiness is not assessed; existing documents are evidence of sources, not acceptance.\n\n")
			for _, ph := range Phases {
				fmt.Fprintf(&b, "## %s\n\n", ph.Label)
				p.docs(&b, entry.Documents[ph.Key])
			}
			b.WriteString("## Sessions\n\n")
			p.docs(&b, entry.Documents["sessions"])
			p.timeline(&b, &a)
		case "history":
			b.WriteString("History follows canonical causal order. Source events retain exact identities in the canonical ledger.\n\n")
			for _, e := range a.Events {
				fmt.Fprintf(&b, "- %s → **%s**: %s. Effective: %s; recorded: %s. %s\n", fromLabel(e), stateLabel(e.Payload.To), text(e.Payload.Reason), stamp(e.OccurredAt), stamp(e.Payload.RecordedAt), p.sourceLink("Canonical history", "ProjectManagement/Ledger.ndjson"))
				for _, link := range e.Payload.Links {
					fmt.Fprintf(&b, "  - %s\n", p.sourceLink("Supporting source", link))
				}
			}
		case "timeline":
			p.timeline(&b, &a)
		case "sessions":
			p.docs(&b, entry.Documents[view])
		default:
			found := false
			for _, ph := range Phases {
				if ph.Key == view {
					found = true
					fmt.Fprintf(&b, "## %s\n\n", ph.Label)
					p.docs(&b, entry.Documents[view])
				}
			}
			if !found {
				return Page{}, fmt.Errorf("unknown Activity phase")
			}
		}
	}
	assets := map[string][]byte{}
	if selection == "overview" || strings.HasSuffix(selection, "/timeline") || (strings.HasPrefix(selection, "activity/") && len(strings.Split(selection, "/")) == 2) {
		only := ""
		if strings.HasPrefix(selection, "activity/") {
			only = strings.Split(selection, "/")[1]
		}
		assets = p.timelineAssets(only)
		content := b.String()
		b.Reset()
		cut := strings.Index(content, "\n\n")
		if cut < 0 {
			cut = 0
		}
		b.WriteString(content[:cut])
		b.WriteString("\n\n![Activity calendar timeline](timeline.svg)\n\n")
		b.WriteString(content[cut:])
	}
	sum := sha256.Sum256([]byte(p.Revision + "\n" + b.String()))
	return Page{Title: title, Markdown: b.String(), Revision: hex.EncodeToString(sum[:]), Assets: assets}, nil
}
func stateLabel(s string) string {
	switch s {
	case "planned":
		return "Planned"
	case "active":
		return "Active"
	case "on-hold":
		return "On hold"
	case "completed":
		return "Completed"
	case "canceled":
		return "Canceled"
	}
	return "Unknown"
}
func fromLabel(e Event) string {
	if e.Payload.From == nil {
		return "Created"
	}
	return stateLabel(*e.Payload.From)
}
func stamp(s string) string {
	t, e := time.Parse(time.RFC3339Nano, s)
	if e != nil {
		return "Unknown effective time"
	}
	return t.UTC().Format(time.RFC3339Nano) + " (UTC)"
}
func (p Projection) timeline(b *strings.Builder, only *Activity) {
	b.WriteString("## Lifecycle timeline\n\nCalendar time is not effort or productive hours. Events follow causal ledger order; timestamps never reorder history. Scope starts at recorded Activity adoption.\n\n")
	for _, entry := range p.Entries {
		a := entry.Activity
		if only != nil && a.Descriptor.ID != only.Descriptor.ID {
			continue
		}
		fmt.Fprintf(b, "### %s\n\n", text(a.Descriptor.Title))
		events := []Event{}
		for _, e := range a.Events {
			if e.EventType == "x-activity:created" || e.EventType == "x-activity:transitioned" {
				events = append(events, e)
			}
		}
		for i, e := range events {
			start, err := time.Parse(time.RFC3339Nano, e.OccurredAt)
			end := p.AsOf
			endLabel := "ongoing as of " + stamp(p.AsOf.Format(time.RFC3339Nano))
			if i+1 < len(events) {
				end, _ = time.Parse(time.RFC3339Nano, events[i+1].OccurredAt)
				endLabel = stamp(events[i+1].OccurredAt)
			}
			if e.Payload.To == "completed" || e.Payload.To == "canceled" {
				fmt.Fprintf(b, "- ◆ %s at %s (terminal marker).\n", stateLabel(e.Payload.To), stamp(e.OccurredAt))
				continue
			}
			if err != nil || end.IsZero() {
				fmt.Fprintf(b, "- %s: unknown timestamp; duration unavailable.\n", stateLabel(e.Payload.To))
			} else if end.Before(start) {
				fmt.Fprintf(b, "- %s: clock anomaly (%s → %s); duration unavailable.\n", stateLabel(e.Payload.To), stamp(e.OccurredAt), endLabel)
			} else {
				fmt.Fprintf(b, "- %s: %s → %s; calendar duration %s.\n", stateLabel(e.Payload.To), stamp(e.OccurredAt), endLabel, end.Sub(start))
			}
		}
		p.markers(b, a)
		b.WriteString("\n")
	}
}
func (p Projection) markers(b *strings.Builder, a Activity) {
	b.WriteString("\nMarker legend: ◆ lifecycle transition; ◇ ledger event; │ explicitly associated commit. Commit time is unavailable unless supplied by a checked repository association; no unrelated Git traffic is inferred.\n\n")
	events, err := p.relatedEvents(a)
	if err != nil {
		fmt.Fprintf(b, "Linked management markers unavailable: %s\n\n", text(err.Error()))
	}
	all := append(append([]Event{}, a.Events...), events...)
	if len(all) > 50 {
		fmt.Fprintf(b, "%d ledger events in this Activity row. Daily marker summary (UTC):\n\n", len(all))
		counts := map[string]int{}
		for _, e := range all {
			key := "Unknown effective date"
			if t, err := time.Parse(time.RFC3339Nano, e.OccurredAt); err == nil {
				key = t.UTC().Format("2006-01-02")
			}
			counts[key]++
		}
		keys := []string{}
		for k := range counts {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(b, "- %s: %d ledger markers.\n", k, counts[k])
		}
		b.WriteString("\nComplete marker details and source links:\n\n")
	}
	seen := map[string]bool{}
	for _, e := range all {
		if p.markerFilter != "commits" && p.markerFilter != "none" {
			fmt.Fprintf(b, "- ◇ %s — %s. Recorded: %s. %s\n", stamp(e.OccurredAt), text(e.Payload.Reason), stamp(e.Payload.RecordedAt), p.sourceLink("Ledger source", "ProjectManagement/Ledger.ndjson"))
			for _, link := range e.Payload.Links {
				fmt.Fprintf(b, "  - %s\n", p.sourceLink("Attributed source", link))
			}
		}
		if p.markerFilter != "ledger" && p.markerFilter != "none" && e.Commit != nil && *e.Commit != "" && !seen[*e.Commit] {
			seen[*e.Commit] = true
			subject, when, repository, err := p.commit(*e.Commit)
			if err != nil {
				fmt.Fprintf(b, "- │ Explicitly associated commit; committer timestamp unknown (Git unavailable or association unresolved). %s\n", p.sourceLink("Immutable commit", "git:commit/"+*e.Commit))
			} else {
				fmt.Fprintf(b, "- │ %s — %s (recorded committer timestamp; repository %s). %s\n", text(subject), stamp(when), text(repository), p.sourceLink("Immutable commit", "git:commit/"+*e.Commit))
			}
		}
	}
}

// Resolve only a full immutable commit identity explicitly attached to a source
// event. Never scan a date range or infer attribution from the active interval.
func (p Projection) commit(id string) (subject, when, repository string, err error) {
	if !regexp.MustCompile(`^[0-9a-fA-F]{40}([0-9a-fA-F]{24})?$`).MatchString(id) {
		return "", "", "", fmt.Errorf("invalid immutable commit identity")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "-C", p.area, "show", "-s", "--format=%H%x00%cI%x00%s", id+"^{commit}")
	out, e := cmd.Output()
	if e != nil {
		return "", "", "", e
	}
	parts := strings.SplitN(string(out), "\x00", 3)
	if len(parts) != 3 || !strings.EqualFold(parts[0], id) {
		return "", "", "", fmt.Errorf("commit mismatch")
	}
	repo, e := exec.CommandContext(ctx, "git", "-C", p.area, "rev-parse", "--show-toplevel").Output()
	if e != nil {
		return "", "", "", e
	}
	return strings.TrimSpace(parts[2]), parts[1], strings.TrimSpace(string(repo)), nil
}

// Join management traffic by explicit identities/paths, never a timestamp window.
func (p Projection) relatedEvents(a Activity) ([]Event, error) {
	f, err := os.Open(filepath.Join(p.area, "ProjectManagement", "Ledger.ndjson"))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := io.ReadAll(io.LimitReader(f, 64<<20+1))
	if err != nil {
		return nil, err
	}
	if len(b) > 64<<20 {
		return nil, fmt.Errorf("linked marker source exceeds 64 MiB")
	}
	sum := sha256.Sum256(b)
	if p.Catalog.HistoryRevision != "" && hex.EncodeToString(sum[:]) != p.Catalog.HistoryRevision {
		return nil, fmt.Errorf("history changed; refresh Activity view")
	}
	known := map[string]bool{a.Descriptor.ID: true, a.Path: true}
	for _, id := range a.Descriptor.SourceCards {
		known[id] = true
	}
	for _, r := range a.Descriptor.References {
		if r.ID != "" {
			known[r.ID] = true
		}
		if r.Path != "" {
			known[r.Path] = true
		}
	}
	for _, s := range a.Sessions {
		known[s.ID] = true
		known[s.Path] = true
	}
	delete(known, "")
	var result []Event
	seen := map[string]bool{}
	for _, e := range a.Events {
		seen[e.EventID] = true
	}
	for _, line := range strings.Split(string(b), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var e Event
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, err
		}
		if seen[e.EventID] {
			continue
		}
		matched := known[e.SubjectID]
		var envelope struct {
			Payload map[string]json.RawMessage `json:"payload"`
		}
		if err := json.Unmarshal([]byte(line), &envelope); err != nil {
			return nil, err
		}
		for _, key := range []string{"activityId", "sessionId", "assignmentId", "planId", "cardId"} {
			var id string
			_ = json.Unmarshal(envelope.Payload[key], &id)
			matched = matched || known[id]
		}
		for _, link := range e.Payload.Links {
			matched = matched || known[link]
		}
		if matched {
			seen[e.EventID] = true
			result = append(result, e)
		}
	}
	return result, nil
}
