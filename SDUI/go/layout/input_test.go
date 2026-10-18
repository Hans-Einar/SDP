package layout

import (
	"encoding/json"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// Synthetic row sizes exercise the shared seam, not actual native font/wrap behavior.
type inputMetrics struct{ *scalarMetrics }

func newInputMetrics() *inputMetrics { return &inputMetrics{newScalarMetrics()} }
func (m *inputMetrics) MeasureField(n *parser.Instance, f runtime.FieldState, font float64, outer Size) (FieldMetrics, error) {
	if n.Widget != "input" {
		return m.scalarMetrics.MeasureField(n, f, font, outer)
	}
	m.seen[n.Path], m.fonts[n.Path] = f, font
	label, body := 10., 30.
	if n.Argument("text") == "" {
		label = 0
	}
	if f.Input != nil && f.Input.Multiline {
		body = 60
	}
	w, h := math.Max(80, outer.W), math.Max(label+body+20, outer.H)
	out := FieldMetrics{Minimum: Size{80, label + body + 20}, Control: Rect{0, label, w, h - label - 20}, Feedback: Rect{0, h - 20, w, 20}}
	if label > 0 {
		out.Label = Rect{0, 0, w, label}
	}
	if m.change != nil {
		m.change(n, f, outer, &out)
	}
	return out, nil
}
func (*inputMetrics) MeasureViewport(n *parser.Instance, font float64) (ViewportInsets, error) {
	return (&gutterMetrics{}).MeasureViewport(n, font)
}
func inputGeometry(t *testing.T, s runtime.Snapshot, size Size, m *inputMetrics) *SnapshotLayout {
	t.Helper()
	g, err := (&Engine{Measure: m}).LayoutSnapshot(s, size)
	if err != nil {
		t.Fatal(err)
	}
	return g
}

func TestInputExplicitOptInAndLegacyGeometry(t *testing.T) {
	for _, profile := range []string{"0.2", "0.3"} {
		s, err := runtime.New("legacy-input", compileProfile(t, profile, `[n=input("Old",value="original")] {scale=1,gap=0}`))
		if err != nil {
			t.Fatal(err)
		}
		w, _ := s.Widget("page/n")
		if err = s.Draft(w.Handle, "retained legacy draft"); err != nil {
			t.Fatal(err)
		}
		snap := s.Snapshot()
		before := snapshotJSON(t, snap)
		m := newInputMetrics()
		got := inputGeometry(t, snap, Size{200, 100}, m)
		want, err := (&Engine{Measure: defaultPaneMetrics()}).LayoutSnapshot(snap, Size{200, 100})
		if err != nil {
			t.Fatal(err)
		}
		a, _ := json.Marshal(got)
		b, _ := json.Marshal(want)
		if string(a) != string(b) || strings.Contains(string(a), `"Fields"`) || len(m.seen) != 0 || before != snapshotJSON(t, snap) {
			t.Fatal("legacy geometry/projection changed", profile, string(a), string(b))
		}
		f := snap.Fields["page/n"]
		f.Input = &runtime.InputState{Multiline: true}
		snap.Fields["page/n"] = f
		if g := inputGeometry(t, snap, Size{200, 100}, m); len(g.Fields) != 0 {
			t.Fatal("field metadata promoted legacy input")
		}
		s.Close()
	}
	for _, arg := range []string{`multiline=false`, `readOnly=false`, `required=false`, `placeholder=""`} {
		t.Run(arg, func(t *testing.T) {
			s := fieldSession(t, `[n=input("Name",`+arg+`)] {scale=1,gap=0}`)
			m := newInputMetrics()
			g := inputGeometry(t, s.Snapshot(), Size{200, 100}, m)
			if len(g.Fields) != 1 || m.seen["page/n"].Input == nil || g.Fields["page/n"].Control != (Rect{0, 10, 80, 30}) {
				t.Fatal("explicit default did not select measured input", g.Fields)
			}
			if _, _, err := parser.Compile(`sdui 0.2;page=[n=input("Name",` + arg + `)];`); err == nil {
				t.Fatal("extended argument admitted under .2")
			}
		})
	}
}

func TestInputLabelsFeedbackAndDetachedProjection(t *testing.T) {
	s := fieldSession(t, `[blank=input("",multiline=true,required=true,placeholder="");named=input("Notes",multiline=true,readOnly=true)] {scale=1,gap=0,font=18}`)
	snap := s.Snapshot()
	before := snapshotJSON(t, snap)
	m := newInputMetrics()
	g := inputGeometry(t, snap, Size{200, 200}, m)
	blank, named := g.Fields["page/blank"], g.Fields["page/named"]
	if blank.Label != (Rect{}) || blank.LabelClip != (Rect{}) || blank.Control != (Rect{0, 0, 80, 60}) || blank.Feedback != (Rect{0, 60, 80, 20}) {
		t.Fatal("empty label not absent", blank)
	}
	if named.Label != (Rect{0, 80, 80, 10}) || named.Control != (Rect{0, 90, 80, 60}) || !named.ReadOnly || !named.Enabled {
		t.Fatal("named readonly input geometry", named)
	}
	if m.seen["page/blank"].Validation.Code == "" || m.seen["page/blank"].Input.Placeholder != "" || m.seen["page/named"].Input.Placeholder != "Notes" || m.fonts["page/named"] != 18 {
		t.Fatal("policy/feedback/font projection lost", m.seen)
	}
	if g.Root.Hit(10, 100).Path != "page/named" || len(g.Viewports) != 0 {
		t.Fatal("readonly input lost hit or owns text viewport")
	}
	m.seen["page/blank"].Input.Placeholder = "mutated"
	m.seen["page/blank"].Input.Multiline = false
	*m.seen["page/blank"].RawDraft = "mutated"
	if before != snapshotJSON(t, snap) {
		t.Fatal("adapter aliases snapshot input metadata/text")
	}
	snap.Root.Rows[1][0].Layout["enabled"] = false
	disabled := inputGeometry(t, snap, Size{200, 200}, newInputMetrics())
	if disabled.Fields["page/named"].Enabled || !disabled.Fields["page/named"].ReadOnly || disabled.Root.Hit(10, 100) != nil {
		t.Fatal("disabled input policy/hit")
	}
}

func TestInputNestedGuttersExtentsAndReveal(t *testing.T) {
	s := fieldSession(t, `[inner=<n=input("Notes",multiline=true) {scale-y=2,x=fill}> {x=fill,scale-y=2,gap=0,overflow-y=scroll}] {scale=1,gap=0,overflow-y=scroll}`)
	snap := s.Snapshot()
	snap.Viewports = map[string]runtime.ViewportState{"page": {Y: 30}, "page/inner": {Y: 50}}
	m := newInputMetrics()
	g := inputGeometry(t, snap, Size{200, 100}, m)
	path := "page/inner/n"
	f := g.Fields[path]
	if f.Control != (Rect{0, -70, 176, 370}) || f.ControlClip != (Rect{0, 0, 176, 100}) || f.Feedback != (Rect{0, 300, 176, 20}) {
		t.Fatal("entry transformed/clipped twice", f)
	}
	if len(g.Viewports) != 2 || g.Viewports["page"].Maximum.Y != 100 || g.Viewports["page/inner"].Maximum.Y != 200 {
		t.Fatal("entry created shared viewport/extent", g.Viewports)
	}
	if g.Root.Hit(175, 10).Path != path || g.Root.Hit(180, 10) != nil || g.Root.Hit(194, 10) != nil {
		t.Fatal("entry hits leaked into ancestor gutter")
	}
	next, rest, err := g.RouteScroll(180, 20, 0, 10)
	if err != nil || rest.Y != 0 || next["page/inner"].Y != 60 || next["page"].Y != 30 {
		t.Fatal("inner gutter unreachable", next, rest, err)
	}
	next, rest, err = g.RouteScroll(194, 20, 0, 10)
	if err != nil || rest.Y != 0 || next["page"].Y != 40 || next["page/inner"].Y != 50 {
		t.Fatal("outer gutter unreachable", next, rest, err)
	}
	next, err = g.EnsureVisible(path, Rect{0, -80, 176, 10})
	if err != nil || next["page"].Y != 0 || next["page/inner"].Y != 0 {
		t.Fatal("outer reveal failed", next, err)
	}
	fstate := snap.Fields[path]
	fstate.Proposed = runtime.Text(strings.Repeat("long Unicode 文\r\n", 1000))
	fstate.RawDraft = &fstate.Proposed.Text
	snap.Fields[path] = fstate
	long := inputGeometry(t, snap, Size{200, 100}, newInputMetrics())
	if !reflect.DeepEqual(long.Fields, g.Fields) || !reflect.DeepEqual(long.Viewports, g.Viewports) {
		t.Fatal("draft content became outer scroll extent")
	}
	// RouteScroll intentionally remains an outer-only API. The native Entry must
	// consume its wheel before invoking it; this synthetic test cannot prove that.
}

func TestInputRelativeMinimumAndReusedIdentity(t *testing.T) {
	_, roots, err := parser.Compile(`sdui 0.3;Part=[input("Anonymous",multiline=true) {x=fill}] {scale=1,padding=0.1,gap=0};page=[s=split(axis="horizontal")[left=Part;right=Part] {scale=1}] {scale=1,gap=0};`)
	if err != nil {
		t.Fatal(err)
	}
	s, err := runtime.New("reused-input", roots["page"])
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	snap := s.Snapshot()
	different := false
	for path, f := range snap.Fields {
		different = different || path != f.Target.Handle.Path
	}
	if !different {
		t.Fatal("fixture did not exercise distinct normalized/runtime identity")
	}
	m := newInputMetrics()
	probes := 0
	m.change = func(_ *parser.Instance, _ runtime.FieldState, size Size, f *FieldMetrics) {
		if size.W < 80 || size.H < 90 {
			probes++
			f.Label.X = math.NaN()
		}
	}
	g := inputGeometry(t, snap, Size{330, 150}, m)
	if len(g.Fields) != 2 || probes == 0 || g.Splits["page/s"].First != (Rect{0, 0, 160, 150}) {
		t.Fatal("input intrinsic probe/relative split", probes, g.Fields, g.Splits)
	}
	for _, f := range g.Fields {
		if f.Control.W != 96 || f.Control.H != 60 || f.ControlClip != f.Control {
			t.Fatal("relative finite inner allocation", f)
		}
	}
	before := snapshotJSON(t, snap)
	if got, err := (&Engine{Measure: newInputMetrics()}).LayoutSnapshot(snap, Size{200, 150}); got != nil || err == nil {
		t.Fatal("relative native minima ignored", got, err)
	}
	if before != snapshotJSON(t, snap) {
		t.Fatal("unfit relative minimum mutated snapshot")
	}
}
