package layout

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

type scalarMetrics struct {
	*paneMetrics
	seen   map[string]runtime.FieldState
	fonts  map[string]float64
	change func(*parser.Instance, runtime.FieldState, Size, *FieldMetrics)
}

func newScalarMetrics() *scalarMetrics {
	return &scalarMetrics{paneMetrics: defaultPaneMetrics(), seen: map[string]runtime.FieldState{}, fonts: map[string]float64{}}
}
func (m *scalarMetrics) MeasureField(n *parser.Instance, f runtime.FieldState, font float64, outer Size) (FieldMetrics, error) {
	m.seen[n.Path], m.fonts[n.Path] = f, font
	w := math.Max(80, outer.W)
	result := FieldMetrics{Minimum: Size{80, 60}, Label: Rect{0, 0, w, 10}, Control: Rect{0, 10, w, 30}, Feedback: Rect{0, 40, w, 20}}
	if n.Widget == "number" {
		result.Control.W -= 30
		result.Decrement = Rect{w - 30, 10, 15, 30}
		result.Increment = Rect{w - 15, 10, 15, 30}
	}
	if m.change != nil {
		m.change(n, f, outer, &result)
	}
	return result, nil
}
func fieldSession(t *testing.T, source string) *runtime.Session {
	t.Helper()
	s, err := runtime.New("scalar-layout", compileProfile(t, "0.3", source))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(s.Close)
	return s
}
func scalarGeometry(t *testing.T, s runtime.Snapshot, size Size, m *scalarMetrics) *SnapshotLayout {
	t.Helper()
	g, err := (&Engine{Measure: m}).LayoutSnapshot(s, size)
	if err != nil {
		t.Fatal(err)
	}
	return g
}
func snapshotJSON(t *testing.T, s runtime.Snapshot) string {
	t.Helper()
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestScalarTypedProjectionAndSourceRemainDistinct(t *testing.T) {
	s := fieldSession(t, `[check=checkbox("Check",value=true,readOnly=true);slide=slider("Slide",min=0,max=10,step=1,value=2);choice=select("Choice");count=number("Count",min=0,max=10,step=1,value=1.00)] {scale=1,gap=0,font=18}`)
	if err := s.BindChoices(map[string][]runtime.ChoiceOption{"page/choice": {{ID: "a", Label: "Same", Enabled: true}, {ID: "b", Label: "Same", Enabled: true}}}); err != nil {
		t.Fatal(err)
	}
	count, _ := s.Widget("page/count")
	if _, err := s.EditField(count.Handle, s.Revision, runtime.Text("-")); err != nil {
		t.Fatal(err)
	}
	snap := s.Snapshot()
	before := snapshotJSON(t, snap)
	m := newScalarMetrics()
	g := scalarGeometry(t, snap, Size{240, 300}, m)
	if len(g.Fields) != 4 {
		t.Fatal("missing scalar geometry", g.Fields)
	}
	for i, name := range []string{"check", "slide", "choice", "count"} {
		path := "page/" + name
		f := g.Fields[path]
		if f.Label != (Rect{0, float64(i * 60), 80, 10}) || f.Feedback != (Rect{0, float64(i*60 + 40), 80, 20}) || f.FeedbackClip != f.Feedback || m.fonts[path] != 18 {
			t.Fatal("label/feedback/font mismatch", path, f)
		}
	}
	f := m.seen["page/count"]
	if f.RawDraft == nil || *f.RawDraft != "-" || f.Proposed.Kind != "" || f.Accepted != runtime.Numeric(1) || f.Validation.Code == "" {
		t.Fatal("typed invalid draft lost", f)
	}
	value := boxes(g.Root)["page/count"].Instance.Arguments["value"].(parser.Literal)
	if value.Kind != "number-lexeme" || value.Value != "1.00" {
		t.Fatal("layout projected live values into lexical source", value)
	}
	if !g.Fields["page/check"].ReadOnly || !g.Fields["page/check"].Enabled || g.Root.Hit(10, 15).Path != "page/check" {
		t.Fatal("readonly incorrectly disabled readable control")
	}
	if before != snapshotJSON(t, snap) {
		t.Fatal("layout mutated source/typed snapshot")
	}
	// Retaining or modifying the adapter's measurement input cannot change snapshot state.
	*m.seen["page/count"].RawDraft = "mutated"
	m.seen["page/count"].Numeric.Min = "mutated"
	m.seen["page/choice"].Options[0].Label = "mutated"
	if before != snapshotJSON(t, snap) {
		t.Fatal("adapter field input aliases snapshot")
	}
}

func TestScalarNestedTranslationClippingAndReveal(t *testing.T) {
	source := `[outer=[inner=[count=number("Count",min=0,max=10,step=1,value=2,readOnly=true) {scale=2}] {scale-x=1,scale-y=2,overflow-x=scroll,overflow-y=scroll,gap=0};tail=button("Tail") {scale=1}] {scale-x=0.5,scale-y=1,overflow-x=clip,overflow-y=scroll,gap=0},sibling=button("S") {scale-x=0.5,scale-y=1}] {scale=1,gap=0}`
	snap := fieldSession(t, source).Snapshot()
	snap.Viewports = map[string]runtime.ViewportState{"page/outer": {Y: 5}, "page/outer/inner": {X: 30, Y: 10}}
	g := scalarGeometry(t, snap, Size{200, 100}, newScalarMetrics())
	path := "page/outer/inner/count"
	f := g.Fields[path]
	if f.Control != (Rect{-30, -5, 170, 30}) || f.ControlClip != (Rect{0, 0, 100, 25}) || f.Increment != (Rect{155, -5, 15, 30}) || f.Feedback != (Rect{-30, 25, 200, 20}) || f.FeedbackClip != (Rect{0, 25, 100, 20}) {
		t.Fatal("parts transformed or clipped incorrectly", f)
	}
	if g.Root.Hit(10, 10).Path != path || g.Root.Hit(150, 10).Path != "page/sibling" || g.Root.Hit(10, 100) != nil {
		t.Fatal("outer hit clip changed")
	}
	offsets, err := g.EnsureVisible(path, f.Increment)
	if err != nil || offsets["page/outer"] != (runtime.ViewportState{}) || offsets["page/outer/inner"] != (runtime.ViewportState{X: 100, Y: 10}) {
		t.Fatal("exact step-button reveal", offsets, err)
	}
	offsets, rest, err := g.RouteScroll(10, 10, 0, 50)
	if err != nil || rest != (runtime.ViewportState{}) || offsets["page/outer/inner"].Y != 60 || offsets["page/outer"].Y != 5 {
		t.Fatal("field became a second scroll owner", offsets, rest, err)
	}
	snap.Root.Rows[0][0].Layout["enabled"] = false
	disabled := scalarGeometry(t, snap, Size{200, 100}, newScalarMetrics())
	if disabled.Fields[path].Enabled || !disabled.Fields[path].ReadOnly || disabled.Root.Hit(10, 10) != nil {
		t.Fatal("inherited disabled field remained interactive")
	}
	if _, err := disabled.EnsureVisible(path, f.Control); err == nil {
		t.Fatal("disabled field revealed")
	}
}

func TestScalarSplitMinimaAndHiddenCanvasBranches(t *testing.T) {
	session := fieldSession(t, `[s=split(axis="horizontal",proportion=0.1)[a=checkbox("A") {x=fill};b=number("B",min=0,max=10,step=1,value=1) {x=fill}] {scale=1};hidden=checkbox("Hidden") {visible=false};d=dialog("D")[choice=select("Choice") {x=fill};nested=dialog("Nested")[bad=checkbox("Closed")]] {scale=0.5,gap=0}] {scale=1,gap=0}`)
	snap := session.Snapshot()
	m := newScalarMetrics()
	g := scalarGeometry(t, snap, Size{210, 100}, m)
	if g.Splits["page/s"].First.W != 80 || g.Splits["page/s"].Second.W != 120 || len(g.Fields) != 2 {
		t.Fatal("split ignored scalar minimum", g.Splits, g.Fields)
	}
	for _, path := range []string{"page/hidden", "page/d/choice", "page/d/nested/bad"} {
		if _, ok := m.seen[path]; ok {
			t.Fatal("inactive field measured", path)
		}
	}
	if got, err := (&Engine{Measure: m}).LayoutSnapshot(snap, Size{169, 100}); err == nil || got != nil {
		t.Fatal("insufficient split minimum accepted", err)
	}
	handle, _ := session.Surface("page/d")
	if _, err := session.OpenSurface(handle, runtime.ContextTarget{}); err != nil {
		t.Fatal(err)
	}
	snap = session.Snapshot()
	all, err := (&Engine{Measure: m}).LayoutCanvases(snap, Size{210, 100}, map[string]Size{"page/d": {160, 80}})
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Main.Fields) != 2 || len(all.Surfaces["page/d"].Fields) != 1 || all.Surfaces["page/d"].Fields["page/d/choice"].Control != (Rect{0, 10, 160, 30}) {
		t.Fatal("separate canvas field origin/branching", all.Main.Fields, all.Surfaces["page/d"].Fields)
	}
	before := snapshotJSON(t, snap)
	if got, err := (&Engine{Measure: m}).LayoutCanvases(snap, Size{210, 100}, map[string]Size{"page/d": {40, 40}}); err == nil || got != nil {
		t.Fatal("partial canvas on invalid native minimum", err)
	}
	if before != snapshotJSON(t, snap) {
		t.Fatal("failed canvas mutated source/state")
	}
}

func TestScalarAnonymousAndReusedRuntimeIdentity(t *testing.T) {
	_, roots, err := parser.Compile(`sdui 0.3;Part=[checkbox("Anonymous");named=number("Named",min=0,max=10,step=1,value=2)];page=[left=Part;right=Part] {scale=1,gap=0};`)
	if err != nil {
		t.Fatal(err)
	}
	s, err := runtime.New("reuse-fields", roots["page"])
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
		t.Fatal("fixture failed to exercise normalized/public identity distinction")
	}
	g := scalarGeometry(t, snap, Size{240, 300}, newScalarMetrics())
	if len(g.Fields) != 4 {
		t.Fatal("anonymous/reused fields lost", g.Fields)
	}
}

func TestScalarAdjunctLeavesExact02GeometryUnchanged(t *testing.T) {
	for _, profile := range []string{"0.2", "0.3"} {
		root := compileProfile(t, profile, `[a=button("A"),b=input("B")] {scale=1,gap=0}`)
		m := newScalarMetrics()
		got, err := (&Engine{Measure: m}).LayoutSnapshot(runtime.Snapshot{Root: root}, Size{200, 100})
		if err != nil {
			t.Fatal(err)
		}
		want, err := (&Engine{Measure: defaultPaneMetrics()}).LayoutSnapshot(runtime.Snapshot{Root: root}, Size{200, 100})
		if err != nil {
			t.Fatal(err)
		}
		a, _ := json.Marshal(got)
		b, _ := json.Marshal(want)
		if string(a) != string(b) || strings.Contains(string(a), `"Fields"`) || len(m.seen) != 0 || len(got.Fields) != 0 {
			t.Fatal("legacy projection changed", string(a), string(b))
		}
		if !reflect.DeepEqual(got.Root, want.Root) {
			t.Fatal("legacy box changed")
		}
	}
	for _, kind := range []string{"checkbox", "slider", "select", "number"} {
		if _, _, err := parser.Compile(fmt.Sprintf(`sdui 0.2;page=[%s("Field")];`, kind)); err == nil {
			t.Fatal("new kind admitted under .2", kind)
		}
	}
	root := compileProfile(t, "0.2", `[button("A")] {scale=1,overflow-y=scroll}`)
	m := newScalarMetrics()
	if _, err := (&Engine{Measure: m}).LayoutSnapshot(runtime.Snapshot{Root: root}, Size{200, 100}); err == nil || !strings.Contains(err.Error(), "unsupported-scroll") || len(m.seen) != 0 {
		t.Fatal("legacy unsupported scroll changed", err)
	}
}
