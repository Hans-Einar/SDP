package fynehost

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/numeric"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"math"
	"strconv"
	"strings"
)

func isScalar(kind string) bool {
	return kind == "checkbox" || kind == "slider" || kind == "select" || kind == "number"
}

// scalarControl is one persistent native composition. Runtime owns every value;
// the native widgets hold only its last accepted presentation and gesture input.
type scalarControl struct {
	widget.BaseWidget
	bundle               *Bundle
	path, kind           string
	state                ui.FieldState
	label, feedback      *widget.Label
	check                *scalarCheck
	slider               *scalarSlider
	selectBox            *scalarSelect
	entry                *numberEntry
	decrement, increment *scalarStep
	body                 fyne.CanvasObject
	parts                layout.FieldMetrics
	disabled             bool
	popup                *choicePopup
}

func newScalarControl(b *Bundle, n *parser.Instance) *scalarControl {
	c := &scalarControl{bundle: b, path: n.Path, kind: n.Widget}
	c.ExtendBaseWidget(c)
	c.label = widget.NewLabel(n.Argument("label"))
	c.feedback = widget.NewLabel("")
	c.feedback.Truncation = fyne.TextTruncateEllipsis
	switch c.kind {
	case "checkbox":
		c.check = &scalarCheck{owner: c}
		c.check.Text = n.Argument("label")
		c.check.ExtendBaseWidget(c.check)
		c.body = c.check
	case "slider":
		c.label.Truncation = fyne.TextTruncateEllipsis
		c.slider = &scalarSlider{owner: c}
		c.slider.Min = 0
		c.slider.Max = 1
		c.slider.Step = 0
		c.slider.ExtendBaseWidget(c.slider)
		c.body = c.slider
		c.slider.OnChanged = func(v float64) { c.slider.propose(v) }
	case "select":
		c.selectBox = &scalarSelect{owner: c}
		c.selectBox.PlaceHolder = "Select…"
		c.selectBox.ExtendBaseWidget(c.selectBox)
		c.body = c.selectBox
	case "number":
		c.entry = &numberEntry{owner: c}
		c.entry.ExtendBaseWidget(c.entry)
		c.entry.SetPlaceHolder(n.Argument("placeholder"))
		c.body = c.entry
		c.entry.OnChanged = func(v string) {
			if c.live() {
				c.edit(ui.Text(v), false)
			}
		}
		c.entry.OnSubmitted = func(string) { c.commitCurrent() }
		c.decrement = newScalarStep(c, -1)
		c.increment = newScalarStep(c, 1)
	}
	return c
}
func (c *scalarControl) CreateRenderer() fyne.WidgetRenderer {
	objects := []fyne.CanvasObject{c.label, c.body, c.feedback}
	if c.kind == "checkbox" {
		c.label.Hide()
	}
	if c.entry != nil {
		objects = append(objects, c.decrement, c.increment)
	}
	return widget.NewSimpleRenderer(container.New(scalarPartsLayout{c}, objects...))
}
func (c *scalarControl) MinSize() fyne.Size {
	return fyne.NewSize(float32(c.parts.Minimum.W), float32(c.parts.Minimum.H))
}

type scalarPartsLayout struct{ c *scalarControl }

func (p scalarPartsLayout) MinSize([]fyne.CanvasObject) fyne.Size { return p.c.MinSize() }
func (p scalarPartsLayout) Layout(_ []fyne.CanvasObject, _ fyne.Size) {
	c := p.c
	place := func(o fyne.CanvasObject, r layout.Rect) {
		if o == nil {
			return
		}
		o.Move(fyne.NewPos(float32(r.X), float32(r.Y)))
		o.Resize(fyne.NewSize(float32(r.W), float32(r.H)))
	}
	place(c.label, c.parts.Label)
	place(c.body, c.parts.Control)
	place(c.feedback, c.parts.Feedback)
	if c.entry != nil {
		place(c.decrement, c.parts.Decrement)
		place(c.increment, c.parts.Increment)
	}
}

// The composite's children participate in native Tab order exactly once. Only
// the entry, not its pointer step affordances, is a number's keyboard focus stop.
func controlFocusable(o fyne.CanvasObject) (fyne.Focusable, bool) {
	if c, ok := o.(*textControl); ok {
		return c.entry, true
	}
	if c, ok := o.(*scalarControl); ok {
		f, ok := c.body.(fyne.Focusable)
		return f, ok
	}
	f, ok := o.(fyne.Focusable)
	return f, ok
}
func (c *scalarControl) Disable() {
	c.disabled = true
	if d, ok := c.body.(fyne.Disableable); ok {
		d.Disable()
	}
	c.Refresh()
}
func (c *scalarControl) Enable() {
	c.disabled = false
	if d, ok := c.body.(fyne.Disableable); ok {
		d.Enable()
	}
	c.Refresh()
}
func (c *scalarControl) Disabled() bool { return c.disabled }
func (c *scalarControl) live() bool     { return c.bundle != nil && c.bundle.livePane() }
func (c *scalarControl) editable() bool { return c.live() && !c.disabled && !c.state.ReadOnly }
func (c *scalarControl) focus() {
	if !c.live() || c.disabled {
		return
	}
	b := c.bundle
	if err := b.owner.Mutate(func(s *ui.Session) error { return s.Focus(c.state.Target.Handle) }); err == nil {
		b.owner.ensureWidget(c.path)
	}
}
func (c *scalarControl) focusCanvas() {
	if c.live() && !c.disabled {
		f, _ := controlFocusable(c)
		c.bundle.canvasFor(c.path).Focus(f)
	}
}
func (c *scalarControl) shortcut(s fyne.Shortcut) {
	if c.live() {
		c.bundle.routeShortcut(c.bundle.canvasFor(c.path), s)
	}
}
func (c *scalarControl) context(e *fyne.PointEvent) {
	if c.live() {
		c.bundle.openContext(c.path, "", e.AbsolutePosition)
	}
}
func (c *scalarControl) escape() {
	if !c.live() {
		return
	}
	if c.popup != nil {
		c.popup.native.close()
		return
	}
	if c.slider != nil {
		c.slider.cancelGesture()
	}
	f, ok := c.bundle.Session.Field(c.state.Target.Handle)
	if ok && f.Dirty {
		c.bundle.owner.Mutate(func(s *ui.Session) error { return s.RevertField(f.Target.Handle) })
		return
	}
	c.bundle.escapeSurface(c.path)
}
func (c *scalarControl) edit(v ui.Value, commit bool) {
	if !c.editable() {
		return
	}
	b := c.bundle
	b.owner.Mutate(func(s *ui.Session) error {
		change, err := s.EditField(c.state.Target.Handle, c.state.Target.ModelRevision, v)
		if err != nil {
			return err
		}
		if commit {
			return c.dispatch(s, change.Field.Target)
		}
		return nil
	})
	if c.entry != nil && b.owner.current == b && !b.closed {
		if f, ok := b.Session.Field(c.state.Target.Handle); ok && f.RawDraft != nil && c.entry.Text != *f.RawDraft {
			muted := b.muted
			b.muted = true
			c.entry.SetText(*f.RawDraft)
			b.muted = muted
		}
	}
}
func (c *scalarControl) dispatch(s *ui.Session, target ui.FieldTarget) error {
	event, err := s.CaptureCommit(target)
	if err != nil {
		return err
	}
	return s.Dispatch(event)
}
func (c *scalarControl) commitCurrent() {
	if !c.editable() {
		return
	}
	c.bundle.owner.Mutate(func(s *ui.Session) error {
		f, ok := s.Field(c.state.Target.Handle)
		if !ok {
			return fmt.Errorf("stale-field")
		}
		return c.dispatch(s, f.Target)
	})
}
func (c *scalarControl) grid() (*numeric.Grid, error) {
	f := c.state
	if f.Numeric == nil {
		return nil, fmt.Errorf("numeric constraints unavailable")
	}
	return numeric.NewGrid(f.Numeric.Min, f.Numeric.Max, f.Numeric.Step)
}
func (c *scalarControl) step(direction int, commit bool) bool {
	if !c.editable() {
		return false
	}
	g, err := c.grid()
	if err != nil {
		c.bundle.owner.status(err)
		return false
	}
	value := c.state.Proposed.Number
	if c.entry != nil {
		if c.state.RawDraft == nil {
			return false
		}
		value, err = g.Parse(*c.state.RawDraft)
		if err != nil {
			c.bundle.owner.status(err)
			return false
		}
	}
	tick, err := g.Tick(value)
	if err != nil {
		c.bundle.owner.status(err)
		return false
	}
	if direction < 0 {
		if tick == 0 {
			return false
		}
		tick--
	} else {
		if tick == g.LastTick() {
			return false
		}
		tick++
	}
	return c.editTick(tick, commit)
}
func (c *scalarControl) editTick(tick uint64, commit bool) bool {
	if !c.editable() {
		return false
	}
	changed := false
	c.bundle.owner.Mutate(func(s *ui.Session) error {
		change, err := s.EditTick(c.state.Target.Handle, c.state.Target.ModelRevision, tick)
		if err != nil {
			return err
		}
		changed = true
		if commit {
			return c.dispatch(s, change.Field.Target)
		}
		return nil
	})
	return changed
}
func (c *scalarControl) sync(state ui.FieldState, parts layout.FieldLayout, box *layout.Box) {
	c.state = state
	local := func(r layout.Rect) layout.Rect {
		if r == (layout.Rect{}) {
			return r
		}
		r.X -= box.Rect.X
		r.Y -= box.Rect.Y
		return r
	}
	c.parts = layout.FieldMetrics{Label: local(parts.Label), Control: local(parts.Control), Feedback: local(parts.Feedback), Decrement: local(parts.Decrement), Increment: local(parts.Increment)}
	label := box.Instance.Argument("label")
	if c.slider != nil {
		label = sliderLabel(label, state)
	}
	c.label.SetText(label)
	text := state.Validation.Message
	if text == "" {
		text = state.Validation.Code
	}
	c.feedback.SetText(text)
	if c.check != nil {
		c.check.Text = box.Instance.Argument("label")
		c.check.Checked = state.Proposed.Bool
		c.check.Refresh()
	}
	if c.slider != nil {
		g, e := c.grid()
		if e == nil {
			tick, e := g.Tick(state.Proposed.Number)
			if e == nil {
				v := 0.0
				if g.LastTick() > 0 {
					v = float64(tick) / float64(g.LastTick())
				}
				c.slider.Value = v
				c.slider.Refresh()
			}
		}
	}
	if c.selectBox != nil {
		label := ""
		labels := make([]string, 0, len(state.Options))
		for _, o := range state.Options {
			labels = append(labels, o.Label)
			if o.ID == state.Proposed.OptionID {
				label = o.Label
			}
		}
		c.selectBox.Options = labels
		c.selectBox.Selected = label
		c.selectBox.Refresh()
	}
	if c.entry != nil && state.RawDraft != nil && c.entry.Text != *state.RawDraft {
		c.entry.SetText(*state.RawDraft)
	}
	if c.entry != nil {
		for _, step := range []*scalarStep{c.decrement, c.increment} {
			if state.ReadOnly || !box.Enabled {
				step.button.Disable()
			} else {
				step.button.Enable()
			}
		}
	}
	c.Refresh()
}

// Value-first projection preserves numeric feedback when a long source label
// is ellipsized. The marker identifies a proposal; neither string changes state.
func sliderLabel(label string, state ui.FieldState) string {
	value := strconv.FormatFloat(state.Proposed.Number, 'g', -1, 64)
	if state.Dirty {
		value += " *"
	}
	return value + " | " + label
}

// Shortest round-trip binary64 formatting needs at most 24 ASCII characters.
// Reserve the widest possible numeric glyph repeated to that bound, plus draft
// marker, separator and ellipsis. This is independent of value/constraints, so
// full numeric feedback fits even when the source label must be truncated.
func sliderLabelMinimum(font float64) float64 {
	width := float64(0)
	for _, glyph := range "0123456789-+.e" {
		reserved := widget.NewLabel(strings.Repeat(string(glyph), 24) + " * | …")
		measured := container.NewThemeOverride(reserved, componentTheme{float32(font)}).MinSize().Width
		width = math.Max(width, float64(measured))
	}
	return width
}

// Use native themed minima and one fixed feedback row. Diagnostic width never
// changes admission or the grid; full validation remains in Snapshot.Fields.
func (m *collectionMeasure) MeasureField(n *parser.Instance, f ui.FieldState, font float64, outer layout.Size) (layout.FieldMetrics, error) {
	if extendedInput(n) {
		return measureTextField(n, f, font, outer)
	}
	c := newScalarControl(nil, n)
	// Measure the original label, never current proposal/dirty text. Truncation
	// is a rendering policy, not permission to shrink source-label admission.
	c.label.Truncation = fyne.TextTruncateOff
	label := container.NewThemeOverride(c.label, componentTheme{float32(font)}).MinSize()
	body := container.NewThemeOverride(c.body, componentTheme{float32(font)}).MinSize()
	feedback := widget.NewLabel("Invalid")
	feedback.Truncation = fyne.TextTruncateEllipsis
	ft := container.NewThemeOverride(feedback, componentTheme{float32(font)})
	fh := float64(ft.MinSize().Height)
	w, h := math.Max(float64(label.Width), float64(body.Width)), float64(label.Height)+float64(body.Height)+fh
	if n.Widget == "slider" {
		w = math.Max(w, sliderLabelMinimum(font))
	}

	step := float64(0)
	if n.Widget == "number" {
		step = math.Max(float64(body.Height), float64(c.increment.MinSize().Width))
		w = math.Max(w, float64(body.Width)+2*step)
	}
	if n.Widget == "checkbox" {
		h = float64(body.Height) + fh
		w = float64(body.Width)
	}
	r := layout.FieldMetrics{Minimum: layout.Size{W: w, H: h}}
	width := math.Max(w, outer.W)
	height := math.Max(h, outer.H)
	lh := float64(label.Height)
	bh := height - lh - fh
	if n.Widget == "checkbox" {
		lh = 0
		bh = height - fh
	}
	r.Label = layout.Rect{W: width, H: float64(label.Height)}
	r.Control = layout.Rect{Y: lh, W: width, H: bh}
	if n.Widget == "checkbox" {
		r.Label = r.Control
	}
	if n.Widget == "number" {
		r.Decrement = layout.Rect{Y: lh, W: step, H: bh}
		r.Increment = layout.Rect{X: width - step, Y: lh, W: step, H: bh}
		r.Control.X = step
		r.Control.W = width - 2*step
	}
	r.Feedback = layout.Rect{Y: height - fh, W: width, H: fh}
	return r, nil
}

// The native slider's public renderer supplies exact painted track/thumb objects
// for diagnostics. No private fields or duplicated geometry calculations.
type scalarSlider struct {
	widget.Slider
	owner                                                *scalarControl
	renderer                                             fyne.WidgetRenderer
	dragging, canceled, keyHeld, changed, shift, tapping bool
	key                                                  fyne.KeyName
}

func (s *scalarSlider) CreateRenderer() fyne.WidgetRenderer {
	s.renderer = s.Slider.CreateRenderer()
	return s.renderer
}
func (s *scalarSlider) FocusGained() { s.Slider.FocusGained(); s.owner.focus() }
func (s *scalarSlider) FocusLost() {
	// Losing focus ends ownership of a held gesture, without committing or
	// reverting its proposal. A late release must not finalize it; the next
	// physical gesture starts with fresh pointer/key bookkeeping.
	s.shift = false
	s.cancelGesture()
	s.dragging = false
	s.Slider.FocusLost()
}
func (s *scalarSlider) TypedRune(rune)                     {}
func (s *scalarSlider) TypedShortcut(e fyne.Shortcut)      { s.owner.shortcut(e) }
func (s *scalarSlider) TappedSecondary(e *fyne.PointEvent) { s.owner.context(e) }
func (s *scalarSlider) propose(ratio float64) {
	if !s.owner.editable() || s.canceled {
		return
	}
	g, err := s.owner.grid()
	if err != nil {
		s.owner.bundle.owner.status(err)
		return
	}
	ratio = math.Max(0, math.Min(1, ratio))
	tick := uint64(math.Round(ratio * float64(g.LastTick())))
	if tick > g.LastTick() {
		tick = g.LastTick()
	}
	// A tap is one automatic gesture: editTick commits the original Change
	// target inside this call. Only a later, explicit drag/key release recaptures.
	s.changed = s.owner.editTick(tick, s.tapping) || s.changed
	// Native Slider updates first. Restore accepted projection on a rejected gate.
	if f, ok := s.owner.bundle.Session.Field(s.owner.state.Target.Handle); ok {
		if k, e := g.Tick(f.Proposed.Number); e == nil {
			s.Value = 0
			if g.LastTick() > 0 {
				s.Value = float64(k) / float64(g.LastTick())
			}
			s.Refresh()
		}
	}
}
func (s *scalarSlider) Dragged(e *fyne.DragEvent) {
	if !s.owner.editable() {
		return
	}
	if !s.dragging {
		s.dragging = true
		s.canceled = false
		s.changed = false
		s.owner.focusCanvas()
	}
	if !s.canceled {
		s.Slider.Dragged(e)
	}
}
func (s *scalarSlider) DragEnd() {
	if !s.dragging {
		return
	}
	s.dragging = false
	if !s.canceled && s.changed {
		s.owner.commitCurrent()
	}
	s.changed = false
	s.canceled = false
}
func (s *scalarSlider) Tapped(e *fyne.PointEvent) {
	if !s.owner.editable() {
		s.owner.focusCanvas()
		return
	}
	s.canceled = false
	s.changed = false
	s.tapping = true
	defer func() { s.tapping = false; s.changed = false }()
	s.Slider.Tapped(e)
}
func (s *scalarSlider) cancelGesture() {
	s.canceled = true
	s.changed = false
	s.keyHeld = false
	s.key = ""
}
func (s *scalarSlider) KeyDown(e *fyne.KeyEvent) {
	if shiftKey(e) {
		s.shift = true
	}
	switch e.Name {
	case fyne.KeyLeft, fyne.KeyRight, fyne.KeyUp, fyne.KeyDown, fyne.KeyHome, fyne.KeyEnd:
		if !s.keyHeld {
			s.keyHeld = true
			s.key = e.Name
			s.changed = false
			s.canceled = false
		}
	}
}
func (s *scalarSlider) KeyUp(e *fyne.KeyEvent) {
	if shiftKey(e) {
		s.shift = false
	}
	if s.keyHeld && e.Name == s.key {
		s.keyHeld = false
		if !s.canceled && s.changed {
			s.owner.commitCurrent()
		}
		s.changed = false
	}
}
func (s *scalarSlider) TypedKey(e *fyne.KeyEvent) {
	if commandKeyEvent(e, s.shift, false, s.owner.shortcut) {
		return
	}
	if e.Name == fyne.KeyEscape {
		s.owner.escape()
		return
	}
	if !s.owner.editable() {
		return
	}
	switch e.Name {
	case fyne.KeyLeft, fyne.KeyDown:
		s.changed = s.owner.step(-1, false) || s.changed
	case fyne.KeyRight, fyne.KeyUp:
		s.changed = s.owner.step(1, false) || s.changed
	case fyne.KeyHome:
		s.changed = s.owner.editTick(0, false) || s.changed
	case fyne.KeyEnd:
		if g, err := s.owner.grid(); err == nil {
			s.changed = s.owner.editTick(g.LastTick(), false) || s.changed
		}
	}
}

// Check and Select override the native mutation hooks before any optimistic
// value paint. The inherited Fyne renderer, accessibility and focus remain native.
type scalarCheck struct {
	widget.Check
	owner *scalarControl
	shift bool
}

func (c *scalarCheck) Tapped(*fyne.PointEvent) {
	c.owner.focusCanvas()
	if c.owner.editable() {
		c.owner.edit(ui.Bool(!c.owner.state.Proposed.Bool), true)
	}
}
func (c *scalarCheck) TypedKey(e *fyne.KeyEvent) {
	if commandKeyEvent(e, c.shift, false, c.owner.shortcut) {
		return
	}
	switch e.Name {
	case fyne.KeySpace:
		c.Tapped(nil)
	case fyne.KeyEscape:
		c.owner.escape()
	}
}
func (c *scalarCheck) FocusGained() { c.Check.FocusGained(); c.owner.focus() }
func (c *scalarCheck) FocusLost()   { c.shift = false; c.Check.FocusLost() }
func (c *scalarCheck) KeyDown(e *fyne.KeyEvent) {
	if shiftKey(e) {
		c.shift = true
	}
}
func (c *scalarCheck) KeyUp(e *fyne.KeyEvent) {
	if shiftKey(e) {
		c.shift = false
	}
}
func (c *scalarCheck) TypedShortcut(s fyne.Shortcut)      { c.owner.shortcut(s) }
func (c *scalarCheck) TappedSecondary(e *fyne.PointEvent) { c.owner.context(e) }

type scalarSelect struct {
	widget.Select
	owner *scalarControl
	shift bool
}

func (s *scalarSelect) Tapped(*fyne.PointEvent) { s.owner.focusCanvas(); s.owner.openChoices() }
func (s *scalarSelect) TypedKey(e *fyne.KeyEvent) {
	if commandKeyEvent(e, s.shift, false, s.owner.shortcut) {
		return
	}
	switch e.Name {
	case fyne.KeySpace, fyne.KeyEnter, fyne.KeyReturn, fyne.KeyUp, fyne.KeyDown:
		s.owner.openChoices()
	case fyne.KeyEscape:
		s.owner.escape()
	}
}
func (s *scalarSelect) FocusGained() { s.Select.FocusGained(); s.owner.focus() }
func (s *scalarSelect) FocusLost()   { s.shift = false; s.Select.FocusLost() }
func (s *scalarSelect) KeyDown(e *fyne.KeyEvent) {
	if shiftKey(e) {
		s.shift = true
	}
}
func (s *scalarSelect) KeyUp(e *fyne.KeyEvent) {
	if shiftKey(e) {
		s.shift = false
	}
}
func (s *scalarSelect) TypedShortcut(e fyne.Shortcut)      { s.owner.shortcut(e) }
func (s *scalarSelect) TappedSecondary(e *fyne.PointEvent) { s.owner.context(e) }

// Read-only number entry retains native selection and copy. Entry edits and
// mutations from paste/cut/undo are intercepted, not reverted after publication.
type numberEntry struct {
	widget.Entry
	owner *scalarControl
	shift bool
}

func (e *numberEntry) FocusGained() { e.Entry.FocusGained(); e.owner.focus() }
func (e *numberEntry) FocusLost()   { e.shift = false; e.Entry.FocusLost() }
func (e *numberEntry) KeyDown(k *fyne.KeyEvent) {
	e.Entry.KeyDown(k)
	if shiftKey(k) {
		e.shift = true
	}
}
func (e *numberEntry) KeyUp(k *fyne.KeyEvent) {
	e.Entry.KeyUp(k)
	if shiftKey(k) {
		e.shift = false
	}
}
func (e *numberEntry) TypedRune(r rune) {
	if e.owner.editable() {
		e.Entry.TypedRune(r)
	}
}
func (e *numberEntry) TypedKey(k *fyne.KeyEvent) {
	if commandKeyEvent(k, e.shift, true, e.owner.shortcut) {
		return
	}
	switch k.Name {
	case fyne.KeyEscape:
		e.owner.escape()
		return
	case fyne.KeyUp:
		e.owner.step(1, true)
		return
	case fyne.KeyDown:
		e.owner.step(-1, true)
		return
	}
	if e.owner.state.ReadOnly {
		switch k.Name {
		case fyne.KeyBackspace, fyne.KeyDelete, fyne.KeyEnter, fyne.KeyReturn:
			return
		}
	}
	e.Entry.TypedKey(k)
}
func (e *numberEntry) TypedShortcut(s fyne.Shortcut) {
	switch s.(type) {
	case *fyne.ShortcutCopy, *fyne.ShortcutSelectAll:
		e.Entry.TypedShortcut(s)
		return
	}
	if nativeEditingShortcut(s) {
		if e.owner.editable() {
			e.Entry.TypedShortcut(s)
		}
		return
	}
	e.owner.shortcut(s)
}
func (e *numberEntry) TappedSecondary(p *fyne.PointEvent) {
	if e.owner.bundle != nil && e.owner.bundle.hasContext(e.owner.path) {
		e.owner.context(p)
	} else { /* suppress Entry's unchecked cut/paste popup; keyboard copy remains native */
	}
}

// Pointer-only native step affordance: no extra Tab stops and no hidden Button
// focus mutation. The public Button renderer supplies its themed chrome.
type scalarStep struct {
	widget.BaseWidget
	button    *widget.Button
	owner     *scalarControl
	direction int
}

func newScalarStep(c *scalarControl, d int) *scalarStep {
	label := "+"
	if d < 0 {
		label = "-"
	}
	s := &scalarStep{owner: c, direction: d, button: widget.NewButton(label, nil)}
	s.ExtendBaseWidget(s)
	return s
}
func (s *scalarStep) CreateRenderer() fyne.WidgetRenderer { return s.button.CreateRenderer() }
func (s *scalarStep) MinSize() fyne.Size                  { return s.button.MinSize() }
func (s *scalarStep) Resize(size fyne.Size) {
	// Button.Refresh uses its public Size; its adapter alone is mounted/focusable.
	s.button.Resize(size)
	s.BaseWidget.Resize(size)
}
func (s *scalarStep) Tapped(*fyne.PointEvent) { s.owner.focusCanvas(); s.owner.step(s.direction, true) }

var _ fyne.Layout = scalarPartsLayout{}
