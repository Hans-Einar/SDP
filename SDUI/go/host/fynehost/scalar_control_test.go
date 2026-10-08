package fynehost

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

const scalarSource = `sdui 0.3;
page=[
 flag=checkbox("Flag") {x=fill};
 level=slider("Level", min=10, max=20, step=2, value=10) {x=fill};
 mode=select("Mode", value="one") {x=fill};
 count=number("Count", min=-10, max=20, step=2, value=0) {x=fill};
 edit=input("Text") {x=fill}
] {x=fill,y=fill};`

func scalarRequest(t *testing.T) DocumentRequest {
	t.Helper()
	doc, e := parser.Parse(scalarSource)
	if e != nil {
		t.Fatal(e)
	}
	return DocumentRequest{Document: doc, Entry: "page", SessionID: "fields", Sequence: 1, SourceRevision: "scalar", Mode: preparation.Prototype, Choices: map[string][]ui.ChoiceOption{"page/mode": {{ID: "one", Label: "Duplicate", Enabled: true}, {ID: "two", Label: "Duplicate", Enabled: true}, {ID: "off", Label: "Disabled"}}}}
}
func scalarHost(t *testing.T) (*DocumentHost, *Bundle) {
	t.Helper()
	h := documentHost(t)
	if e := h.Adopt(scalarRequest(t)); e != nil {
		t.Fatal(e)
	}
	return h, h.Current()
}
func scalarAt(b *Bundle, path string) *scalarControl {
	return b.Controls()["page/"+path].(*scalarControl)
}
func fieldAt(t *testing.T, b *Bundle, path string) ui.FieldState {
	t.Helper()
	w, ok := b.Session.Widget("page/" + path)
	if !ok {
		t.Fatal(path)
	}
	f, ok := b.Session.Field(w.Handle)
	if !ok {
		t.Fatal(path)
	}
	return f
}
func countCommits(t *testing.T, b *Bundle, path string) *int {
	t.Helper()
	calls := new(int)
	f := fieldAt(t, b, path)
	if e := b.Session.Bind(f.Target.Handle, func(e ui.Event) ([]ui.Update, error) {
		*calls++
		gen := uint64(0)
		if e.Control.Option != nil {
			gen = e.Control.Option.OptionGeneration
		}
		return []ui.Update{{Handle: e.Handle, Property: ui.AcceptedValue, Value: e.Value, AcceptDraft: true, ExpectedValueRevision: e.Control.ValueRevision, ExpectedDraftRevision: e.DraftRevision, ExpectedOptionGeneration: gen}}, nil
	}); e != nil {
		t.Fatal(e)
	}
	return calls
}
func TestScalarGestureChangeThenSingleCommit(t *testing.T) {
	_, b := scalarHost(t)
	check := scalarAt(b, "flag")
	changes := 0
	calls := countCommits(t, b, "flag")
	b.Session.ObserveChanges(check.state.Target.Handle, func(change ui.FieldChange) {
		changes++
		if *calls != 0 {
			t.Error("Commit preceded Change")
		}
		if !change.Field.Proposed.Bool {
			t.Error("missing proposal")
		}
	})
	check.check.TypedKey(&fyne.KeyEvent{Name: fyne.KeySpace})
	if changes != 1 || *calls != 1 || !fieldAt(t, b, "flag").Accepted.Bool {
		t.Fatalf("Change=%d Commit=%d field=%+v", changes, *calls, fieldAt(t, b, "flag"))
	}
	slider := scalarAt(b, "level").slider
	calls = countCommits(t, b, "level")
	key := &fyne.KeyEvent{Name: fyne.KeyRight}
	slider.KeyDown(key)
	slider.TypedKey(key)
	slider.TypedKey(key)
	f := fieldAt(t, b, "level")
	if *calls != 0 || !f.Dirty || f.Proposed.Number != 14 || f.Accepted.Number != 10 {
		t.Fatalf("held key committed: %+v calls%d", f, *calls)
	}
	if scalarAt(b, "level").slider != slider {
		t.Fatal("sync recreated native slider")
	}
	slider.KeyUp(key)
	slider.KeyUp(key)
	if *calls != 1 || fieldAt(t, b, "level").Accepted.Number != 14 {
		t.Fatal("release Commit", *calls, fieldAt(t, b, "level"))
	}
	slider.KeyDown(key)
	slider.TypedKey(key)
	slider.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEscape})
	slider.KeyUp(key)
	if *calls != 1 || fieldAt(t, b, "level").Dirty || fieldAt(t, b, "level").Accepted.Number != 14 {
		t.Fatal("Escape committed or retained draft")
	}
}
func TestScalarFailedPreparationDoesNotPaintOrCommit(t *testing.T) {
	h := documentHost(t)
	r := scalarRequest(t)
	fail := false
	r.PrepareResources = func(ui.Snapshot) error {
		if fail {
			return errors.New("field-resource-failure")
		}
		return nil
	}
	if e := h.Adopt(r); e != nil {
		t.Fatal(e)
	}
	b := h.Current()
	entry := scalarAt(b, "count").entry
	calls := countCommits(t, b, "count")
	changes := 0
	b.Session.ObserveChanges(entry.owner.state.Target.Handle, func(ui.FieldChange) { changes++ })
	before := b.presentation
	fail = true
	entry.SetText("2")
	if changes != 0 || *calls != 0 || entry.Text != "0" || b.presentation != before || fieldAt(t, b, "count").Dirty {
		t.Fatal("rejected native edit leaked", changes, *calls, entry.Text)
	}
	slider := scalarAt(b, "level").slider
	calls = countCommits(t, b, "level")
	key := &fyne.KeyEvent{Name: fyne.KeyRight}
	slider.KeyDown(key)
	slider.TypedKey(key)
	slider.KeyUp(key)
	if *calls != 0 || fieldAt(t, b, "level").Dirty {
		t.Fatal("rejected slider edit caused Commit")
	}
}
func TestScalarInvalidNumberAndReadOnlyCopy(t *testing.T) {
	_, b := scalarHost(t)
	c := scalarAt(b, "count")
	calls := countCommits(t, b, "count")
	c.entry.SetText("1") // off-grid
	f := fieldAt(t, b, "count")
	if f.RawDraft == nil || *f.RawDraft != "1" || f.Validation.Code == "" || !f.Dirty {
		t.Fatalf("invalid raw discarded: %+v", f)
	}
	c.increment.Tapped(nil)
	c.entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	if *calls != 0 || c.entry.Text != "1" {
		t.Fatal("invalid number stepped or committed")
	}
	c.entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEscape})
	c.increment.Tapped(nil)
	if *calls != 1 || fieldAt(t, b, "count").Accepted.Number != 2 {
		t.Fatal("number increment did not commit legal point", fieldAt(t, b, "count"))
	}
	f = fieldAt(t, b, "count")
	if e := b.owner.Mutate(func(s *ui.Session) error {
		return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: f.Target.Handle, Property: ui.ReadOnly, Value: ui.Bool(true), ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}})
	}); e != nil {
		t.Fatal(e)
	}
	c.entry.TypedShortcut(&fyne.ShortcutSelectAll{})
	clip := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
	c.entry.TypedShortcut(&fyne.ShortcutCopy{Clipboard: clip})
	if clip.Content() != "2" {
		t.Fatal("readonly cannot copy", clip.Content())
	}
	c.entry.TypedRune('8')
	c.entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
	c.increment.Tapped(nil)
	c.entry.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	if *calls != 1 || c.entry.Text != "2" || fieldAt(t, b, "count").Dirty {
		t.Fatal("readonly mutation")
	}
}
func TestScalarChoicesStableIDAndStaleOpening(t *testing.T) {
	_, b := scalarHost(t)
	c := scalarAt(b, "mode")
	calls := countCommits(t, b, "mode")
	c.openChoices()
	popup := c.popup
	if popup == nil {
		t.Fatal("no native choices")
	}
	if len(popup.targets) != 3 || popup.targets[0].OptionID == popup.targets[1].OptionID {
		t.Fatal("labels became identity")
	}
	popup.native.selection(func() { popup.native.popup.Items[1].(fyne.Tappable).Tapped(&fyne.PointEvent{}) })
	if *calls != 1 || fieldAt(t, b, "mode").Accepted.OptionID != "two" {
		t.Fatal("duplicate-label selected wrong ID", *calls, fieldAt(t, b, "mode"))
	}
	c.openChoices()
	old := c.popup
	f := fieldAt(t, b, "mode")
	if e := b.owner.Mutate(func(s *ui.Session) error {
		return s.ReplaceChoices(f.Target.Handle, f.Target.OptionGeneration, []ui.ChoiceOption{{ID: "two", Label: "New", Enabled: true}})
	}); e != nil {
		t.Fatal(e)
	}
	old.native.selection(func() { old.native.popup.Items[0].(fyne.Tappable).Tapped(&fyne.PointEvent{}) })
	if *calls != 1 || c.popup != nil || fieldAt(t, b, "mode").Accepted.OptionID != "two" {
		t.Fatal("stale choice replay")
	}
}
func TestScalarReentrantChangeCannotRestampAutomaticCommit(t *testing.T) {
	_, b := scalarHost(t)
	c := scalarAt(b, "flag")
	calls := countCommits(t, b, "flag")
	edit, _ := b.Session.Widget("page/edit")
	b.Session.ObserveChanges(c.state.Target.Handle, func(ui.FieldChange) {
		if e := b.Session.Draft(edit.Handle, "observer edit"); e != nil {
			t.Error(e)
		}
	})
	c.check.Tapped(nil)
	if *calls != 0 || fieldAt(t, b, "flag").Accepted.Bool || !fieldAt(t, b, "flag").Dirty {
		t.Fatal("reentrant automatic Commit recaptured", *calls)
	}
	if b.presentation.snapshot.Root == nil || b.Controls()["page/edit"].(*Input).Text != "observer edit" {
		t.Fatal("accepted reentrant state not synchronized")
	}
}
func TestScalarHiddenChoiceReadinessBeforeResources(t *testing.T) {
	h := documentHost(t)
	r := scalarRequest(t)
	doc, e := parser.Parse(strings.Replace(scalarSource, `mode=select("Mode", value="one") {x=fill}`, `mode=select("Mode", value="one") {x=fill,visible=false}`, 1))
	if e != nil {
		t.Fatal(e)
	}
	r.Document = doc
	r.Choices = map[string][]ui.ChoiceOption{}
	calls := 0
	r.PrepareResources = func(ui.Snapshot) error { calls++; return nil }
	if _, e = h.Prepare(r); e == nil || calls != 0 {
		t.Fatal("hidden unsupplied choices admitted or allocated", e, calls)
	}
}

func TestScalarPopupFocusAndEmptyGeneration(t *testing.T) {
	h, b := scalarHost(t)
	h.canvas.SetContent(h.Container)
	c := scalarAt(b, "mode")
	c.focusCanvas()
	c.openChoices()
	popup := c.popup
	if h.canvas.Focused() != popup.native.input {
		t.Fatal("opening lacks native focus")
	}
	if e := h.Mutate(func(s *ui.Session) error { return s.SetViewports(nil) }); e != nil {
		t.Fatal(e)
	}
	if c.popup != popup || h.canvas.Focused() != popup.native.input {
		t.Fatal("ordinary sync stole popup focus")
	}
	popup.native.input.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEscape})
	if c.popup != nil || h.canvas.Focused() != c.selectBox {
		t.Fatalf("popup Escape did not restore select: popup=%p focused=%T session=%q", c.popup, h.canvas.Focused(), b.Session.Focused())
	}
	f := fieldAt(t, b, "mode")
	if e := h.Mutate(func(s *ui.Session) error { return s.ReplaceChoices(f.Target.Handle, f.Target.OptionGeneration, nil) }); e != nil {
		t.Fatal(e)
	}
	c.openChoices()
	empty := c.popup
	if empty == nil {
		t.Fatal("supplied empty not openable")
	}
	f = fieldAt(t, b, "mode")
	if e := h.Mutate(func(s *ui.Session) error {
		return s.ReplaceChoices(f.Target.Handle, f.Target.OptionGeneration, []ui.ChoiceOption{{ID: "one", Label: "One", Enabled: true}})
	}); e != nil {
		t.Fatal(e)
	}
	if c.popup != nil || !empty.native.closed {
		t.Fatal("zero-item opening lost captured generation")
	}
}

func TestScalarPointerDragEscapeAndInspection(t *testing.T) {
	_, b := scalarHost(t)
	c := scalarAt(b, "level")
	s := c.slider
	calls := countCommits(t, b, "level")
	s.Dragged(&fyne.DragEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(s.Size().Width*.6, s.Size().Height/2)}})
	if !fieldAt(t, b, "level").Dirty || *calls != 0 {
		t.Fatal("drag missing proposal or committed")
	}
	before := fieldAt(t, b, "level").Proposed.Number
	if e := b.owner.Mutate(func(session *ui.Session) error { return session.SetViewports(nil) }); e != nil {
		t.Fatal(e)
	}
	if !s.dragging || c.slider != s {
		t.Fatal("sync reset native drag")
	}
	s.DragEnd()
	s.DragEnd()
	if *calls != 1 || fieldAt(t, b, "level").Accepted.Number != before {
		t.Fatal("drag release lost exact proposal")
	}
	s.Dragged(&fyne.DragEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(s.Size().Width*.9, s.Size().Height/2)}})
	s.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEscape})
	s.Dragged(&fyne.DragEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(s.Size().Width*.7, s.Size().Height/2)}})
	s.DragEnd()
	if *calls != 1 || fieldAt(t, b, "level").Dirty || fieldAt(t, b, "level").Accepted.Number != before {
		t.Fatal("canceled drag resumed")
	}
	fields, _ := b.inspectFields()
	for _, part := range []string{"control", "slidertrack", "sliderthumb"} {
		r := fields["page/level"][part].(layout.Rect)
		if r.W <= 0 || r.H <= 0 {
			t.Fatal("missing actual slider part", part, r)
		}
	}
	for _, part := range []string{"entry", "decrement", "increment"} {
		r := fields["page/count"][part].(layout.Rect)
		if r.W <= 0 || r.H <= 0 {
			t.Fatal("missing actual number part", part, r)
		}
	}
}

func TestScalarFeedbackFitsMinimumAtMultipleFonts(t *testing.T) {
	_, b := scalarHost(t)
	snapshot := b.Session.Snapshot()
	measure := &collectionMeasure{snapshot: snapshot}
	for _, font := range []float64{14, 20, 28} {
		snapshot.Root.Walk(func(n *parser.Instance) {
			if !isScalar(n.Widget) {
				return
			}
			f := snapshot.Fields[n.Path]
			m, err := measure.MeasureField(n, f, font, layout.Size{})
			if err != nil {
				t.Fatal(err)
			}
			f.Validation = ui.FieldValidation{Code: "invalid", Message: strings.Repeat("very long diagnostic ", 100)}
			with, err := measure.MeasureField(n, f, font, m.Minimum)
			if err != nil {
				t.Fatal(err)
			}
			if with.Minimum != m.Minimum || with.Feedback.W > m.Minimum.W || with.Feedback.Y+with.Feedback.H > m.Minimum.H+.001 {
				t.Fatal("feedback changed admitted minimum", font, n.Widget, m, with)
			}
		})
	}
}

func TestScalarNumberNativeStepHitAndChrome(t *testing.T) {
	h, b := scalarHost(t)
	h.canvas.SetContent(h.Container)
	h.Container.Resize(fyne.NewSize(800, 500))
	c := scalarAt(b, "count")
	calls := countCommits(t, b, "count")
	for _, step := range []*scalarStep{c.increment, c.decrement} {
		if step.button.Size() != step.Size() || step.Size().Width <= 0 || step.Size().Height <= 0 {
			t.Fatal("native button owner lost allocated size", step.button.Size(), step.Size())
		}
		r := inspectionRect(step)
		test.TapCanvas(h.canvas, fyne.NewPos(float32(r.X+r.W/2), float32(r.Y+r.H/2)))
	}
	if *calls != 2 || fieldAt(t, b, "count").Accepted.Number != 0 {
		t.Fatal("native step hit failed", *calls, fieldAt(t, b, "count"))
	}
	// Entry is the only Tab stop for the number composition.
	h.canvas.Focus(c.entry)
	h.canvas.FocusNext()
	if h.canvas.Focused() != b.Controls()["page/edit"].(*Input) {
		t.Fatalf("number step became extra Tab stop: %T", h.canvas.Focused())
	}
}

func TestScalarExactLargeOriginSteppingAndProgrammaticMute(t *testing.T) {
	h := documentHost(t)
	r := scalarRequest(t)
	doc, err := parser.Parse(strings.Replace(scalarSource, `min=-10, max=20, step=2, value=0`, `min=9007199254740987, max=9007199254740991, step=1, value=9007199254740988`, 1))
	if err != nil {
		t.Fatal(err)
	}
	r.Document = doc
	if err = h.Adopt(r); err != nil {
		t.Fatal(err)
	}
	b := h.Current()
	c := scalarAt(b, "count")
	calls := countCommits(t, b, "count")
	changes := 0
	b.Session.ObserveChanges(c.state.Target.Handle, func(ui.FieldChange) { changes++ })
	c.increment.Tapped(nil)
	if fieldAt(t, b, "count").Accepted.Number != 9007199254740989 || c.entry.Text != "9007199254740989" || changes != 1 || *calls != 1 {
		t.Fatal("native step rounded from zero-origin", fieldAt(t, b, "count"), changes, *calls)
	}
	f := fieldAt(t, b, "count")
	if err = h.Mutate(func(s *ui.Session) error {
		return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: f.Target.Handle, Property: ui.AcceptedValue, Value: ui.Numeric(9007199254740991), ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}})
	}); err != nil {
		t.Fatal(err)
	}
	if c.entry.Text != "9007199254740991" || changes != 1 || *calls != 1 {
		t.Fatal("programmatic projection emitted user gesture", c.entry.Text, changes, *calls)
	}
	c.increment.Tapped(nil)
	if changes != 1 || *calls != 1 {
		t.Fatal("maximum step repeated Commit")
	}
}

func TestScalarSliderTapKeepsOriginalChangeTarget(t *testing.T) {
	_, b := scalarHost(t)
	c := scalarAt(b, "level")
	calls := countCommits(t, b, "level")
	edit, _ := b.Session.Widget("page/edit")
	if err := b.Session.ObserveChanges(c.state.Target.Handle, func(ui.FieldChange) {
		if err := b.Session.Draft(edit.Handle, "observer changed another field"); err != nil {
			t.Error(err)
		}
	}); err != nil {
		t.Fatal(err)
	}
	s := c.slider
	s.Tapped(&fyne.PointEvent{Position: fyne.NewPos(s.Size().Width*.8, s.Size().Height/2)})
	f := fieldAt(t, b, "level")
	if *calls != 0 || f.Accepted.Number != 10 || f.Proposed.Number != 18 || !f.Dirty {
		t.Fatalf("automatic tap restamped observer state: calls=%d field=%+v", *calls, f)
	}
	if b.Controls()["page/edit"].(*Input).Text != "observer changed another field" {
		t.Fatal("accepted reentrant edit lost from native presentation")
	}
	// A later independent tap works after removing the reentrant observer.
	if err := b.Session.ObserveChanges(c.state.Target.Handle, nil); err != nil {
		t.Fatal(err)
	}
	s.Tapped(&fyne.PointEvent{Position: fyne.NewPos(s.Size().Width*.4, s.Size().Height/2)})
	if *calls != 1 || fieldAt(t, b, "level").Dirty {
		t.Fatal("fresh tap failed or committed twice", *calls)
	}
}

func TestScalarSliderFocusLossReleasesHeldKey(t *testing.T) {
	_, b := scalarHost(t)
	s := scalarAt(b, "level").slider
	calls := countCommits(t, b, "level")
	right, left := &fyne.KeyEvent{Name: fyne.KeyRight}, &fyne.KeyEvent{Name: fyne.KeyLeft}
	s.KeyDown(right)
	s.TypedKey(right)
	s.FocusLost()
	if *calls != 0 || s.keyHeld || s.key != "" || !fieldAt(t, b, "level").Dirty {
		t.Fatal("focus loss committed/reverted or retained old key")
	}
	s.KeyUp(right) // late release of the abandoned gesture
	s.FocusGained()
	s.KeyDown(left)
	s.TypedKey(left)
	s.KeyUp(left)
	s.KeyUp(right)
	if *calls != 1 || s.keyHeld || fieldAt(t, b, "level").Accepted.Number != 10 {
		t.Fatal("old key stranded/replayed after new gesture", *calls, s.keyHeld, fieldAt(t, b, "level"))
	}
}

func TestScalarSliderFocusLossReleasesHeldPointer(t *testing.T) {
	_, b := scalarHost(t)
	s := scalarAt(b, "level").slider
	calls := countCommits(t, b, "level")
	drag := func(ratio float32) {
		s.Dragged(&fyne.DragEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(s.Size().Width*ratio, s.Size().Height/2)}})
	}
	drag(.6)
	s.FocusLost()
	if s.dragging || *calls != 0 || !fieldAt(t, b, "level").Dirty {
		t.Fatal("focus loss retained pointer or committed/reverted")
	}
	s.DragEnd() // obsolete release cannot accept the proposal
	if *calls != 0 {
		t.Fatal("late pointer release committed")
	}
	s.FocusGained()
	drag(.8)
	s.DragEnd()
	s.DragEnd()
	if *calls != 1 || fieldAt(t, b, "level").Accepted.Number != 18 || fieldAt(t, b, "level").Dirty {
		t.Fatal("new pointer gesture stranded or duplicated", *calls, fieldAt(t, b, "level"))
	}
}

func TestScalarSliderVisibleProposalAndFixedGeometry(t *testing.T) {
	h, b := scalarHost(t)
	c := scalarAt(b, "level")
	slider := c.slider
	calls := countCommits(t, b, "level")
	changes := 0
	if err := b.Session.ObserveChanges(c.state.Target.Handle, func(ui.FieldChange) { changes++ }); err != nil {
		t.Fatal(err)
	}
	geometry := b.Geometry().Fields[c.path]
	check := func(text string) {
		t.Helper()
		if c.label.Text != text {
			t.Fatalf("native value label=%q, want %q", c.label.Text, text)
		}
		fields, _ := b.inspectFields()
		if fields[c.path]["labelText"] != c.label.Text {
			t.Fatal("inspector did not expose native label.Text")
		}
		got := b.Geometry().Fields[c.path]
		got.ReadOnly, got.Enabled = geometry.ReadOnly, geometry.Enabled // policy is not geometry
		if got != geometry {
			t.Fatal("value/feedback changed admitted field geometry", geometry, got)
		}
	}
	check("10 | Level")
	if err := h.Mutate(func(s *ui.Session) error {
		return s.ValidateFieldWith(c.state.Target.Handle, func(ui.FieldState) ui.FieldValidation {
			return ui.FieldValidation{Code: "local-validation", Message: "A separate validation message"}
		})
	}); err != nil {
		t.Fatal(err)
	}
	key := &fyne.KeyEvent{Name: fyne.KeyRight}
	slider.KeyDown(key)
	slider.TypedKey(key)
	check("12 * | Level")
	if *calls != 0 || changes != 1 || c.feedback.Text != "A separate validation message" || geometry.Feedback.H <= 0 || geometry.Label.Intersect(geometry.Feedback).H > 0 {
		t.Fatal("held value/independent validation/event policy", *calls, changes, c.feedback.Text)
	}
	slider.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEscape})
	slider.KeyUp(key)
	check("10 | Level")
	if *calls != 0 || changes != 1 {
		t.Fatal("revert emitted user callback")
	}
	if err := h.Mutate(func(s *ui.Session) error { return s.ValidateFieldWith(c.state.Target.Handle, nil) }); err != nil {
		t.Fatal(err)
	}
	f := fieldAt(t, b, "level")
	if err := h.Mutate(func(s *ui.Session) error {
		return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: f.Target.Handle, Property: ui.ReadOnly, Value: ui.Bool(true), ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}})
	}); err != nil {
		t.Fatal(err)
	}
	f = fieldAt(t, b, "level")
	if err := h.Mutate(func(s *ui.Session) error {
		return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: f.Target.Handle, Property: ui.AcceptedValue, Value: ui.Numeric(18), ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}})
	}); err != nil {
		t.Fatal(err)
	}
	check("18 | Level")
	slider.KeyDown(key)
	slider.TypedKey(key)
	slider.KeyUp(key)
	check("18 | Level")
	if *calls != 0 || changes != 1 || c.feedback.Text != "" {
		t.Fatal("programmatic/readonly sync emitted event or reused validation row")
	}
}

func TestScalarSliderLabelMinimumDoesNotFollowValue(t *testing.T) {
	_, b := scalarHost(t)
	var source *parser.Instance
	b.Session.Snapshot().Root.Walk(func(n *parser.Instance) {
		if n.Path == "page/level" {
			source = n
		}
	})
	f := fieldAt(t, b, "level")
	measure := &collectionMeasure{}
	tinyDoc, err := parser.Parse(strings.Replace(scalarSource, `slider("Level"`, `slider("X"`, 1))
	if err != nil {
		t.Fatal(err)
	}
	roots, err := parser.Normalize(tinyDoc)
	if err != nil {
		t.Fatal(err)
	}
	var tiny *parser.Instance
	roots["page"].Walk(func(n *parser.Instance) {
		if n.Path == "page/level" {
			tiny = n
		}
	})
	for _, source := range []*parser.Instance{source, tiny} {
		for _, font := range []float64{14, 20, 28} {
			baseline, err := measure.MeasureField(source, f, font, layout.Size{})
			if err != nil {
				t.Fatal(err)
			}
			if baseline.Minimum.W < sliderLabelMinimum(font) {
				t.Fatal("minimum can hide complete value", font, baseline.Minimum)
			}
			for _, value := range []float64{-1.7976931348623157e308, 0, 0.0000000000000001, 9007199254740991} {
				// Measurement consumes state without numeric parsing or range authority.
				f.Proposed = ui.Numeric(value)
				f.Dirty = true
				f.Validation = ui.FieldValidation{Code: "invalid", Message: strings.Repeat("feedback ", 100)}
				candidate, err := measure.MeasureField(source, f, font, baseline.Minimum)
				if err != nil {
					t.Fatal(err)
				}
				if candidate.Minimum != baseline.Minimum {
					t.Fatal("numeric readout changed minimum", font, value)
				}
				native := newScalarControl(nil, source)
				if native.label.Truncation != fyne.TextTruncateEllipsis {
					t.Fatal("unbounded native value label")
				}
				text := sliderLabel(source.Argument("label"), f)
				reserved := widget.NewLabel(strconv.FormatFloat(value, 'g', -1, 64) + " * | …")
				needed := container.NewThemeOverride(reserved, componentTheme{float32(font)}).MinSize()
				if float64(needed.Width) > candidate.Label.W {
					t.Fatal("full signed/exponent value does not fit admitted label", font, value, needed, candidate.Label)
				}
				if !strings.HasPrefix(text, strconv.FormatFloat(value, 'g', -1, 64)+" *") {
					t.Fatal("label hides proposal behind source text", text)
				}
			}
		}
	}
}
