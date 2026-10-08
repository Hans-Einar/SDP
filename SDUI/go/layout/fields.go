package layout

import (
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// FieldMetrics contains actual native measurements local to the outer control.
// Control is the number Entry excluding its two buttons. A checkbox's native
// label can overlap Control. Zero rectangles identify absent optional parts.
// Intrinsic probes consume only Minimum; final allocation checks every part.
type FieldMetrics struct {
	Minimum                                        Size
	Label, Control, Feedback, Decrement, Increment Rect
}

// FieldMeasurer is required for WCI3-M1 scalar controls. The supplied field is a
// detached prospective snapshot projection, not permission to modify live state.
// Methods must be pure; they must not parse numeric values or publish resources.
type FieldMeasurer interface {
	MeasureField(*parser.Instance, runtime.FieldState, float64, Size) (FieldMetrics, error)
}

// FieldLayout uses the owning canvas's coordinates and effective content clips.
// ReadOnly preserves focus/readability; Enabled includes ancestor disabled state.
// Numeric values, choice identity, validation and event authority stay in runtime.
type FieldLayout struct {
	Label, Control, Feedback, Decrement, Increment                     Rect
	LabelClip, ControlClip, FeedbackClip, DecrementClip, IncrementClip Rect
	Enabled, ReadOnly                                                  bool
}

func scalarField(n *parser.Instance) bool {
	if n.Profile != "sdui/0.3" || n.Kind != "widget" {
		return false
	}
	switch n.Widget {
	case "checkbox", "slider", "select", "number":
		return true
	}
	return false
}

func (e *Engine) fieldFor(n *parser.Instance) (runtime.FieldState, error) {
	// Runtime public handle paths can differ for anonymous/reused source nodes.
	// InstancePath, not a second path resolver, joins geometry to the snapshot.
	f, ok := e.fieldState[n.Path]
	if !ok || f.InstancePath != n.Path || f.Target.Handle.Path == "" || f.Target.Handle.Kind != n.Widget {
		return runtime.FieldState{}, diag(n, "field-snapshot", "Scalar layout requires its matching typed snapshot field")
	}
	// An adapter may retain its measurement input; it must not alias snapshot data.
	if f.RawDraft != nil {
		raw := *f.RawDraft
		f.RawDraft = &raw
	}
	if f.Numeric != nil {
		numeric := *f.Numeric
		f.Numeric = &numeric
	}
	if f.Options != nil {
		f.Options = append([]runtime.ChoiceOption{}, f.Options...)
	}
	return f, nil
}

func (e *Engine) fieldMetrics(n *parser.Instance, font float64, outer Size) (FieldMetrics, error) {
	if err := e.step(n); err != nil {
		return FieldMetrics{}, err
	}
	f, err := e.fieldFor(n)
	if err != nil {
		return FieldMetrics{}, err
	}
	adapter, ok := e.Measure.(FieldMeasurer)
	if !ok {
		return FieldMetrics{}, diag(n, "field-measurement", "Scalar layout requires native field metrics")
	}
	m, err := adapter.MeasureField(n, f, font, outer)
	if err != nil {
		return FieldMetrics{}, err
	}
	if !finiteExtent(m.Minimum) || m.Minimum.W <= 0 || m.Minimum.H <= 0 {
		return FieldMetrics{}, diag(n, "field-measurement", "Field minimum must be finite, bounded and positive")
	}
	return m, nil
}

func validFieldRect(r Rect, outer Size) bool {
	return finite(r.X) && finite(r.Y) && finiteExtent(Size{r.W, r.H}) && r.X >= 0 && r.Y >= 0 && r.X+r.W <= outer.W+.01 && r.Y+r.H <= outer.H+.01
}
func positiveFieldRect(r Rect) bool { return r.W > 0 && r.H > 0 }
func fieldOverlap(a, b Rect) bool   { return positiveFieldRect(a.Intersect(b)) }

func (e *Engine) arrangeField(n *parser.Instance, b *Box) error {
	outer := Size{b.Rect.W, b.Rect.H}
	m, err := e.fieldMetrics(n, b.Font, outer)
	if err != nil {
		return err
	}
	if outer.W < m.Minimum.W-.01 || outer.H < m.Minimum.H-.01 {
		return diag(n, "native-minimum", "Assigned field size is below native minimum")
	}
	for _, r := range []Rect{m.Label, m.Control, m.Feedback, m.Decrement, m.Increment} {
		if !validFieldRect(r, outer) || r != (Rect{}) && !positiveFieldRect(r) {
			return diag(n, "field-measurement", "Field parts must be positive finite rectangles inside the assigned outer rectangle, or exactly absent")
		}
	}
	f := e.fieldState[n.Path]
	if !positiveFieldRect(m.Label) || !positiveFieldRect(m.Control) || (f.Validation.Code != "" || f.Validation.Message != "") && !positiveFieldRect(m.Feedback) {
		return diag(n, "field-measurement", "Field label, control and active validation feedback require measured regions")
	}
	if n.Widget == "number" {
		if !positiveFieldRect(m.Decrement) || !positiveFieldRect(m.Increment) || fieldOverlap(m.Decrement, m.Increment) || fieldOverlap(m.Control, m.Decrement) || fieldOverlap(m.Control, m.Increment) {
			return diag(n, "field-measurement", "Number entry and step buttons require distinct measured regions")
		}
	} else if m.Decrement != (Rect{}) || m.Increment != (Rect{}) {
		return diag(n, "field-measurement", "Only number controls have step-button regions")
	}
	if fieldOverlap(m.Feedback, m.Control) || fieldOverlap(m.Feedback, m.Label) || fieldOverlap(m.Feedback, m.Decrement) || fieldOverlap(m.Feedback, m.Increment) {
		return diag(n, "field-measurement", "Validation feedback must not overlap field label or controls")
	}
	project := func(r Rect) (Rect, Rect) {
		if r == (Rect{}) {
			return Rect{}, Rect{}
		}
		r.X += b.Rect.X
		r.Y += b.Rect.Y
		return r, r.Intersect(b.Clip)
	}
	g := FieldLayout{Enabled: b.Enabled, ReadOnly: f.ReadOnly}
	g.Label, g.LabelClip = project(m.Label)
	g.Control, g.ControlClip = project(m.Control)
	g.Feedback, g.FeedbackClip = project(m.Feedback)
	g.Decrement, g.DecrementClip = project(m.Decrement)
	g.Increment, g.IncrementClip = project(m.Increment)
	if e.fields == nil {
		e.fields = map[string]FieldLayout{}
	}
	e.fields[n.Path] = g
	return nil
}
