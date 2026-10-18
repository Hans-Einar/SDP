package layout

import (
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func TestScalarMalformedNativeMetricsRejectPurely(t *testing.T) {
	s := fieldSession(t, `[n=number("N",min=0,max=10,step=1,value=1) {scale=1}] {scale=1,gap=0}`)
	w, _ := s.Widget("page/n")
	if _, err := s.EditField(w.Handle, s.Revision, runtime.Text("-")); err != nil {
		t.Fatal(err)
	}
	snap := s.Snapshot()
	before := snapshotJSON(t, snap)
	cases := map[string]func(*FieldMetrics){
		"zero-minimum":         func(m *FieldMetrics) { m.Minimum.W = 0 },
		"negative-minimum":     func(m *FieldMetrics) { m.Minimum.H = -1 },
		"nan-minimum":          func(m *FieldMetrics) { m.Minimum.W = math.NaN() },
		"infinite-minimum":     func(m *FieldMetrics) { m.Minimum.H = math.Inf(1) },
		"unbounded-minimum":    func(m *FieldMetrics) { m.Minimum.W = 1e8 },
		"missing-label":        func(m *FieldMetrics) { m.Label = Rect{} },
		"partial-zero":         func(m *FieldMetrics) { m.Control.W = 0 },
		"negative-part":        func(m *FieldMetrics) { m.Control.X = -1 },
		"nan-part":             func(m *FieldMetrics) { m.Control.Y = math.NaN() },
		"infinite-part":        func(m *FieldMetrics) { m.Label.W = math.Inf(1) },
		"outside-part":         func(m *FieldMetrics) { m.Label.X = 1 },
		"missing-feedback":     func(m *FieldMetrics) { m.Feedback = Rect{} },
		"overlapping-feedback": func(m *FieldMetrics) { m.Feedback = m.Control },
		"missing-increment":    func(m *FieldMetrics) { m.Increment = Rect{} },
		"overlapping-buttons":  func(m *FieldMetrics) { m.Increment = m.Decrement },
		"overlapping-entry":    func(m *FieldMetrics) { m.Increment = m.Control },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			m := newScalarMetrics()
			m.change = func(_ *parser.Instance, _ runtime.FieldState, _ Size, out *FieldMetrics) { change(out) }
			g, err := (&Engine{Measure: m}).LayoutSnapshot(snap, Size{160, 80})
			var diagnostic *parser.Diagnostic
			if g != nil || !errors.As(err, &diagnostic) || diagnostic.Code != "field-measurement" || !strings.Contains(diagnostic.Message, "page/n") {
				t.Fatal("bad metrics lost typed source diagnostic", g, err)
			}
			if before != snapshotJSON(t, snap) {
				t.Fatal("failed geometry mutated snapshot")
			}
		})
	}
}

func TestScalarSnapshotAndFinalMinimumAdmission(t *testing.T) {
	s := fieldSession(t, `[n=checkbox("N") {x=fill,max-x=0.1}] {scale=1,gap=0}`)
	snap := s.Snapshot()
	m := newScalarMetrics()
	if g, err := (&Engine{Measure: m}).LayoutSnapshot(snap, Size{200, 100}); g != nil || err == nil || !strings.Contains(err.Error(), "native-minimum") {
		t.Fatal("source max bypassed native minimum", g, err)
	}
	delete(snap.Root.Rows[0][0].Layout, "max-x")
	if _, err := (&Engine{Measure: fixtureMetrics{}}).LayoutSnapshot(snap, Size{200, 100}); err == nil || !strings.Contains(err.Error(), "field-measurement") {
		t.Fatal("missing native adapter accepted", err)
	}
	if _, err := (&Engine{Measure: m}).Layout(snap.Root, Size{200, 100}); err == nil || !strings.Contains(err.Error(), "field-snapshot") {
		t.Fatal("source-only layout guessed typed state", err)
	}
	original := snap.Fields["page/n"]
	for _, mode := range []string{"missing", "path", "kind", "handle"} {
		f := original
		switch mode {
		case "missing":
			delete(snap.Fields, "page/n")
		case "path":
			f.InstancePath = "wrong"
			snap.Fields["page/n"] = f
		case "kind":
			f.Target.Handle.Kind = "input"
			snap.Fields["page/n"] = f
		case "handle":
			f.Target.Handle.Path = ""
			snap.Fields["page/n"] = f
		}
		if g, err := (&Engine{Measure: m}).LayoutSnapshot(snap, Size{200, 100}); g != nil || err == nil || !strings.Contains(err.Error(), "field-snapshot") {
			t.Fatal("mismatched snapshot accepted", mode, err)
		}
		snap.Fields["page/n"] = original
	}
	m.change = func(_ *parser.Instance, _ runtime.FieldState, _ Size, f *FieldMetrics) {
		f.Increment = Rect{0, 0, 10, 10}
	}
	if _, err := (&Engine{Measure: m}).LayoutSnapshot(snap, Size{200, 100}); err == nil {
		t.Fatal("checkbox acquired numeric subcontrols")
	}
	// The native Check label belongs inside the combined control's hit rectangle.
	m.change = func(_ *parser.Instance, _ runtime.FieldState, _ Size, f *FieldMetrics) { f.Label = f.Control }
	scalarGeometry(t, snap, Size{200, 100}, m)
}

func TestScalarIntrinsicProbeUsesOnlyMinimumAndInactivePageSkipsMetrics(t *testing.T) {
	s := fieldSession(t, `[s=split(axis="horizontal")[left=[a=checkbox("A") {x=fill}] {scale=1,gap=0};right=tabs("Tabs")[one=page("One")[b=checkbox("B") {x=fill}];two=page("Two")[hidden=checkbox("Hidden")]] {scale=1}] {scale=1,gap=0}] {scale=1,gap=0}`)
	snap := s.Snapshot()
	m := newScalarMetrics()
	probes := 0
	m.change = func(n *parser.Instance, _ runtime.FieldState, size Size, f *FieldMetrics) {
		if n.Argument("label") == "Hidden" {
			t.Fatal("inactive field measured")
		}
		if size.W < 80 || size.H < 60 {
			probes++
			f.Label.X = math.NaN()
		}
	}
	g := scalarGeometry(t, snap, Size{210, 150}, m)
	if probes == 0 || len(g.Fields) != 2 || g.Fields["page/s/right/one/b"].Control.Y != 34 {
		t.Fatal("intrinsic metrics/active pane composition", probes, g.Fields)
	}
}
