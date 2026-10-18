package layout

import (
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func inputFailure(t *testing.T, engine *Engine, snap runtime.Snapshot, code string) {
	t.Helper()
	before := snapshotJSON(t, snap)
	g, err := engine.LayoutSnapshot(snap, Size{200, 100})
	var d *parser.Diagnostic
	if g != nil || !errors.As(err, &d) || d.Code != code || !strings.Contains(d.Message, "page/n") {
		t.Fatalf("expected %s with source path: %v %v", code, g, err)
	}
	if before != snapshotJSON(t, snap) {
		t.Fatal("failed probe mutated snapshot")
	}
}
func TestInputSnapshotAndSourcePolicyAdmission(t *testing.T) {
	s := fieldSession(t, `[n=input("Notes",multiline=true,required=true,placeholder="hint") {scale=1}] {scale=1,gap=0}`)
	for _, mode := range []string{"missing", "identity", "kind", "input", "multiline", "placeholder", "required"} {
		t.Run(mode, func(t *testing.T) {
			snap := s.Snapshot()
			f := snap.Fields["page/n"]
			switch mode {
			case "identity":
				f.InstancePath = "wrong"
			case "kind":
				f.Target.Handle.Kind = "number"
			case "input":
				f.Input = nil
			case "multiline":
				f.Input.Multiline = false
			case "placeholder":
				f.Input.Placeholder = "wrong"
			case "required":
				f.Required = false
			}
			snap.Fields["page/n"] = f
			if mode == "missing" {
				delete(snap.Fields, "page/n")
			}
			m := newInputMetrics()
			inputFailure(t, &Engine{Measure: m}, snap, "field-snapshot")
			if len(m.seen) != 0 {
				t.Fatal("invalid policy reached native measurement")
			}
		})
	}
	snap := s.Snapshot()
	inputFailure(t, &Engine{Measure: fixtureMetrics{}}, snap, "field-measurement")
	if g, err := (&Engine{Measure: newInputMetrics()}).Layout(snap.Root, Size{200, 100}); g != nil || err == nil || !strings.Contains(err.Error(), "field-snapshot") {
		t.Fatal("source-only input guessed state", g, err)
	}
	snap.Root.Rows[0][0].Layout["max-x"] = 0.1
	inputFailure(t, &Engine{Measure: newInputMetrics()}, snap, "native-minimum")
	for _, axis := range []string{"x", "y"} {
		snap = s.Snapshot()
		snap.Root.Rows[0][0].Layout["overflow-"+axis] = "scroll"
		inputFailure(t, &Engine{Measure: newInputMetrics()}, snap, "unsupported-scroll")
	}
	// Layout validates manually supplied source policy even if that input is hidden.
	snap = s.Snapshot()
	n := snap.Root.Rows[0][0]
	n.Layout["visible"] = false
	n.Arguments["multiline"] = parser.Literal{Kind: "string", Value: "true", Span: n.Span}
	inputFailure(t, &Engine{Measure: newInputMetrics()}, snap, "input-argument")
}

func TestInputMalformedPartsAndEmptyLabelBoundary(t *testing.T) {
	s := fieldSession(t, `[n=input("",multiline=true,required=true) {scale=1}] {scale=1,gap=0}`)
	cases := map[string]func(*FieldMetrics){
		"zero-minimum":         func(m *FieldMetrics) { m.Minimum.H = 0 },
		"nan-minimum":          func(m *FieldMetrics) { m.Minimum.W = math.NaN() },
		"missing-entry":        func(m *FieldMetrics) { m.Control = Rect{} },
		"outside-entry":        func(m *FieldMetrics) { m.Control.X = 1 },
		"partial-zero-label":   func(m *FieldMetrics) { m.Label = Rect{W: 10} },
		"overlapping-label":    func(m *FieldMetrics) { m.Label = m.Control },
		"missing-feedback":     func(m *FieldMetrics) { m.Feedback = Rect{} },
		"overlapping-feedback": func(m *FieldMetrics) { m.Feedback = m.Control },
		"numeric-button":       func(m *FieldMetrics) { m.Increment = Rect{0, 0, 10, 10} },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			m := newInputMetrics()
			m.change = func(_ *parser.Instance, _ runtime.FieldState, _ Size, f *FieldMetrics) { change(f) }
			inputFailure(t, &Engine{Measure: m}, s.Snapshot(), "field-measurement")
		})
	}
	named := fieldSession(t, `[n=input("Name",multiline=false) {scale=1}] {scale=1,gap=0}`)
	m := newInputMetrics()
	m.change = func(_ *parser.Instance, _ runtime.FieldState, _ Size, f *FieldMetrics) { f.Label = Rect{} }
	inputFailure(t, &Engine{Measure: m}, named.Snapshot(), "field-measurement")
	// Source spaces are not normalized into a missing label.
	spaced := fieldSession(t, `[n=input(" ",multiline=false) {scale=1}] {scale=1,gap=0}`)
	inputFailure(t, &Engine{Measure: m}, spaced.Snapshot(), "field-measurement")
}

func TestInputRuntimeLabelProjectionAndMutableReadOnly(t *testing.T) {
	s := fieldSession(t, `[n=input("Before",multiline=true)] {scale=1,gap=0}`)
	w, _ := s.Widget("page/n")
	if err := s.Apply(s.Revision, s.BatchRevision+1, []runtime.Update{{Handle: w.Handle, Property: runtime.Label, Value: runtime.Text("After")}, {Handle: w.Handle, Property: runtime.ReadOnly, Value: runtime.Bool(true), ExpectedValueRevision: w.ValueRevision, ExpectedDraftRevision: w.DraftRevision}}); err != nil {
		t.Fatal(err)
	}
	m := newInputMetrics()
	g := inputGeometry(t, s.Snapshot(), Size{200, 100}, m)
	if m.seen["page/n"].Input.Placeholder != "After" || !g.Fields["page/n"].ReadOnly || !g.Fields["page/n"].Enabled {
		t.Fatal("live label/readonly rejected or source default won", m.seen, g.Fields)
	}
}
