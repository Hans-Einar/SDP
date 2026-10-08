package layout

import (
	"math"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

type fixtureMetrics struct{}

func (fixtureMetrics) Measure(*parser.Instance, float64, float64) (Size, error) {
	return Size{20, 20}, nil
}

func compileProfile(t *testing.T, profile, source string) *parser.Instance {
	t.Helper()
	_, roots, err := parser.Compile("sdui " + profile + "; page=" + source + ";")
	if err != nil {
		t.Fatal(err)
	}
	return roots["page"]
}
func snapshotGeometry(t *testing.T, source string, offsets map[string]runtime.ViewportState) *SnapshotLayout {
	t.Helper()
	root := compileProfile(t, "0.3", source)
	g, err := (&Engine{Measure: fixtureMetrics{}}).LayoutSnapshot(runtime.Snapshot{Root: root, Viewports: offsets}, Size{200, 100})
	if err != nil {
		t.Fatal(err)
	}
	return g
}

const nestedViewports = `[
 outer=[
  inner=[a=button("A") {scale=2}] {scale-x=1,scale-y=2,overflow-x=scroll,overflow-y=scroll,gap=0};
  tail=button("T") {scale=1}
 ] {scale-x=0.5,scale-y=1,overflow-x=clip,overflow-y=scroll,gap=0},
 sibling=button("S") {scale-x=0.5,scale-y=1}
] {scale=1,gap=0}`

func TestSnapshotNestedGeometryOracle(t *testing.T) {
	offsets := map[string]runtime.ViewportState{"page/outer": {X: 90, Y: 50}, "page/outer/inner": {X: 30, Y: 100}}
	g := snapshotGeometry(t, nestedViewports, offsets)
	b := boxes(g.Root)
	inner, outer := g.Viewports["page/outer/inner"], g.Viewports["page/outer"]
	if inner.Parent != "page/outer" || outer.Parent != "" {
		t.Fatalf("viewport ancestors: %+v %+v", inner, outer)
	}
	// Nested content is 400 high, but only its 200-high outer box plus tail
	// contributes to the parent's 300-high content extent.
	near(t, inner.Content.W, 200)
	near(t, inner.Content.H, 400)
	near(t, outer.Content.H, 300)
	near(t, inner.Maximum.Y, 200)
	near(t, outer.Maximum.Y, 200)
	if b["page/outer/inner"].Rect != (Rect{0, -50, 100, 200}) || b["page/outer/inner/a"].Rect != (Rect{-30, -150, 200, 400}) {
		t.Fatalf("double/missing transform: %+v %+v", b["page/outer/inner"], b["page/outer/inner/a"])
	}
	if b["page/outer/inner/a"].Clip != (Rect{0, 0, 100, 100}) {
		t.Fatalf("clip: %+v", b["page/outer/inner/a"].Clip)
	}
	if got := g.Root.Hit(10, 10); got == nil || got.Path != "page/outer/inner/a" {
		t.Fatalf("visible hit: %+v", got)
	}
	if got := g.Root.Hit(150, 10); got == nil || got.Path != "page/sibling" {
		t.Fatalf("sibling/clip hit: %+v", got)
	}
	if g.Root.Hit(10, 100) != nil || g.Root.Hit(-1, 10) != nil {
		t.Fatal("half-open clip escaped")
	}
	if outer.Offset.X != 0 || offsets["page/outer"].X != 90 {
		t.Fatal("non-scroll axis or snapshot was mutated")
	}
	out := g.EffectiveOffsets()
	out["page/outer"] = runtime.ViewportState{Y: 999}
	if g.Viewports["page/outer"].Offset.Y != 50 {
		t.Fatal("offset result aliases geometry")
	}
}

func TestScrollRoutingRemainderAndSiblingOracle(t *testing.T) {
	g := snapshotGeometry(t, nestedViewports, map[string]runtime.ViewportState{"page/outer": {Y: 50}, "page/outer/inner": {X: 30, Y: 100}})
	out, rest, err := g.RouteScroll(10, 10, 100, 150)
	if err != nil {
		t.Fatal(err)
	}
	if out["page/outer/inner"] != (runtime.ViewportState{X: 100, Y: 200}) || out["page/outer"].Y != 100 || rest != (runtime.ViewportState{X: 30}) {
		t.Fatalf("routed %+v remainder %+v", out, rest)
	}
	out, rest, err = g.RouteScroll(150, 10, 20, 40)
	if err != nil || rest != (runtime.ViewportState{X: 20, Y: 40}) || out["page/outer"].Y != 50 {
		t.Fatalf("sibling scrolled: %+v %+v %v", out, rest, err)
	}
	out, rest, err = g.RouteScroll(10, 10, -300, -300)
	if err != nil || out["page/outer/inner"] != (runtime.ViewportState{}) || out["page/outer"].Y != 0 || rest != (runtime.ViewportState{X: -270, Y: -150}) {
		t.Fatalf("reverse remainder: %+v %+v %v", out, rest, err)
	}
	if _, _, err = g.RouteScroll(0, 0, math.Inf(1), 0); err == nil {
		t.Fatal("nonfinite delta")
	}
}

func TestEnsureVisibleInnerToOuterOracle(t *testing.T) {
	g := snapshotGeometry(t, nestedViewports, map[string]runtime.ViewportState{"page/outer": {Y: 50}, "page/outer/inner": {X: 30, Y: 100}})
	out, err := g.EnsureVisible("page/outer/inner/a", Rect{150, 240, 10, 10})
	if err != nil {
		t.Fatal(err)
	}
	if out["page/outer/inner"] != (runtime.ViewportState{X: 90, Y: 200}) || out["page/outer"].Y != 100 {
		t.Fatalf("inner-to-outer reveal: %+v", out)
	}
	// An oversized target aligns its leading edge, first within the inner
	// viewport, then within the outer viewport using the translated rectangle.
	out, err = g.EnsureVisible("page/outer/inner/a", boxes(g.Root)["page/outer/inner/a"].Rect)
	if err != nil || out["page/outer/inner"] != (runtime.ViewportState{}) || out["page/outer"].Y != 0 {
		t.Fatalf("oversized reveal: %+v %v", out, err)
	}
	if _, err = g.EnsureVisible("missing", Rect{}); err == nil {
		t.Fatal("missing target")
	}
}

func TestSnapshotClampRemovalAndZeroRange(t *testing.T) {
	root := compileProfile(t, "0.3", `[a=button("A") {scale=2}] {scale=1,gap=0,overflow-x=scroll,overflow-y=scroll}`)
	e := &Engine{Measure: fixtureMetrics{}}
	snap := runtime.Snapshot{Root: root, Viewports: map[string]runtime.ViewportState{"page": {X: 900, Y: 900}, "removed": {Y: 500}}}
	g, err := e.LayoutSnapshot(snap, Size{200, 100})
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Viewports) != 1 || g.EffectiveOffsets()["page"] != (runtime.ViewportState{X: 200, Y: 100}) {
		t.Fatalf("clamp: %+v", g.Viewports)
	}
	root.Rows[0][0].Layout["scale"] = float64(1)
	g, err = e.LayoutSnapshot(snap, Size{200, 100})
	if err != nil || g.EffectiveOffsets()["page"] != (runtime.ViewportState{}) {
		t.Fatalf("shrink: %+v %v", g, err)
	}
	_, rest, err := g.RouteScroll(10, 10, 30, 40)
	if err != nil || rest != (runtime.ViewportState{X: 30, Y: 40}) {
		t.Fatalf("zero range: %+v %v", rest, err)
	}
	for _, bad := range []float64{-1, math.NaN(), math.Inf(1), math.Inf(-1)} {
		snap.Viewports["page"] = runtime.ViewportState{Y: bad}
		if _, err = e.LayoutSnapshot(snap, Size{200, 100}); err == nil {
			t.Errorf("accepted offset %v", bad)
		}
	}
}

func TestScrollProfilesAndPolicies(t *testing.T) {
	cases := []struct{ profile, source, code string }{
		{"0.2", `[button("A")] {scale=1,overflow-y=scroll}`, "unsupported-scroll"},
		{"0.2", `button("A") {scale=1,overflow-y=scroll}`, "unsupported-scroll"},
		{"0.3", `[button("A")] {x=fill,overflow-y=scroll}`, "scroll-layout"},
		{"0.3", `[hidden=[] {overflow-y=scroll,visible=false}] {scale=1}`, "scroll-layout"},
		{"0.3", `[button("A")] {scale=1,overflow-y=scroll,justify=end}`, "scroll-layout"},
		{"0.3", `[button("A") {scale=1,overflow-y=scroll}] {scale=1}`, "unsupported-scroll"},
		{"0.3", `[button("A") {overflow-y=scroll,visible=false}] {scale=1}`, "unsupported-scroll"},
		{"0.3", `[button("A") {scale=2}] {scale=1,gap=0}`, "overflow-x"},
	}
	for _, c := range cases {
		t.Run(c.profile+c.code+c.source, func(t *testing.T) {
			root := compileProfile(t, c.profile, c.source)
			_, err := (&Engine{Measure: fixtureMetrics{}}).Layout(root, Size{200, 100})
			d, ok := err.(*parser.Diagnostic)
			if !ok || d.Code != c.code {
				t.Fatalf("got %v want %s", err, c.code)
			}
		})
	}
	g := snapshotGeometry(t, `[button("A") {scale=2}] {scale=1,gap=0,overflow-x=clip,overflow-y=clip}`, nil)
	if len(g.Viewports) != 0 || g.Root.Hit(199, 99) == nil || g.Root.Hit(200, 99) != nil {
		t.Fatal("clip is not zero-offset clipping")
	}
}

func TestScrollFiniteReferencesRegionsAndAlignment(t *testing.T) {
	// Scale, padding and gap use the finite explicit ancestor, including when
	// fixed content overflows. Footer remains viewport-relative, then scrolls
	// with all other content; it is not a sticky overlay.
	g := snapshotGeometry(t, `[header=button("H"); body=[a=button("A") {scale-x=2,scale-y=2}] {scale=1,gap=0,overflow-x=clip,overflow-y=clip}; footer=button("F")] {scale=1,gap=0,padding=0.1,overflow-y=scroll}`, map[string]runtime.ViewportState{"page": {Y: 5}})
	b := boxes(g.Root)
	near(t, b["page/header"].Rect.X, 20)
	near(t, b["page/header"].Rect.Y, 10)
	near(t, b["page/body"].Rect.H, 40)
	near(t, b["page/footer"].Rect.Y, 70)
	near(t, b["page/body/a"].Rect.W, 320)
	// A center/end-aligned overflowing child in a scroll owner must not move
	// to a negative, unreachable content coordinate.
	g = snapshotGeometry(t, `[a=button("A") {scale-x=1,scale-y=2,align-y=end}] {scale=1,gap=0,overflow-y=scroll}`, nil)
	near(t, boxes(g.Root)["page/a"].Rect.Y, 0)
	// A ratio-resolved dimension is definite, while a cyclic content ancestor
	// remains an error rather than being solved against an infinite plane.
	g = snapshotGeometry(t, `[button("A") {scale-y=2}] {2:1,x=fill,overflow-y=scroll,gap=0}`, nil)
	near(t, g.Viewports["page"].Rect.H, 100)
	root := compileProfile(t, "0.3", `[<button("A") {scale-y=2}> {x=fill}] {scale=1,overflow-y=scroll}`)
	if _, err := (&Engine{Measure: fixtureMetrics{}}).Layout(root, Size{200, 100}); err == nil || !strings.Contains(err.Error(), "layout-dependency") {
		t.Fatalf("indefinite child accepted: %v", err)
	}
}

func TestScrollFooterMovesAndNeverStartsNegative(t *testing.T) {
	g := snapshotGeometry(t, `[header=button("H") {scale-y=0.2}; body=button("B") {scale-y=2}; footer=button("F") {scale-y=0.2}] {scale=1,gap=0,overflow-y=scroll}`, map[string]runtime.ViewportState{"page": {Y: 30}})
	b := boxes(g.Root)
	near(t, b["page/header"].Rect.Y, -30)
	near(t, b["page/footer"].Rect.Y, 50)
	near(t, b["page/body"].Rect.Y, -10)
	g = snapshotGeometry(t, `[footer=button("F") {scale-y=2}] {scale=1,gap=0,overflow-y=scroll}`, nil)
	near(t, boxes(g.Root)["page/footer"].Rect.Y, 0)
	near(t, g.Viewports["page"].Content.H, 200)
}

func TestHiddenDisabledAndObscuredBranchesDoNotRoute(t *testing.T) {
	g := snapshotGeometry(t, nestedViewports, map[string]runtime.ViewportState{"page/outer": {Y: 50}})
	// Disabled hits do not tunnel into a sibling. Outer paint/hit order is
	// also the route order when an overlay obscures a scrolling branch.
	sibling := boxes(g.Root)["page/sibling"]
	sibling.Rect = Rect{0, 0, 200, 100}
	sibling.Clip = sibling.Rect
	sibling.Enabled = false
	out, rest, err := g.RouteScroll(10, 10, 0, 50)
	if err != nil || rest.Y != 50 || out["page/outer"].Y != 50 {
		t.Fatalf("obscured route: %+v %+v %v", out, rest, err)
	}
	if _, err = g.EnsureVisible(sibling.Path, sibling.Rect); err == nil {
		t.Fatal("disabled focus target accepted")
	}
	root := compileProfile(t, "0.3", `[hidden=[button("A"){scale=2}]{scale=1,overflow-x=scroll,overflow-y=scroll,visible=false}] {scale=1}`)
	snap, err := (&Engine{Measure: fixtureMetrics{}}).LayoutSnapshot(runtime.Snapshot{Root: root, Viewports: map[string]runtime.ViewportState{"page/hidden": {Y: 10}}}, Size{200, 100})
	if err != nil || len(snap.Viewports) != 0 {
		t.Fatalf("hidden viewport retained: %+v %v", snap, err)
	}
	if _, err = snap.EnsureVisible("page/hidden", Rect{}); err == nil {
		t.Fatal("hidden focus target accepted")
	}
}

func TestProfileMismatchAndEngineReuse(t *testing.T) {
	e := &Engine{Measure: fixtureMetrics{}}
	root := compileProfile(t, "0.3", `[button("A"){scale=2}] {scale=1,overflow-x=scroll,overflow-y=scroll}`)
	g, err := e.LayoutSnapshot(runtime.Snapshot{Root: root, Viewports: map[string]runtime.ViewportState{"page": {X: 40, Y: 40}}}, Size{200, 100})
	if err != nil || g.Root.Children[0].Rect.X != -40 {
		t.Fatalf("snapshot run: %+v %v", g, err)
	}
	b, err := e.Layout(root, Size{200, 100})
	if err != nil || b.Children[0].Rect.X != 0 {
		t.Fatalf("offset leaked to root-only call: %+v %v", b, err)
	}
	root.Rows[0][0].Profile = ""
	if _, err = e.Layout(root, Size{200, 100}); err == nil {
		t.Fatal("mixed tree accepted")
	}
	root.Profile = "sdui/0.9"
	if _, err = e.Layout(root, Size{200, 100}); err == nil {
		t.Fatal("unknown profile accepted")
	}
}

func TestScrollGroupAndFiniteTrackAllocation(t *testing.T) {
	g := snapshotGeometry(t, `[g=<a=button("A") {scale=2}> {scale=1,gap=0,overflow-x=scroll,overflow-y=scroll}] {scale=1,gap=0}`, map[string]runtime.ViewportState{"page/g": {X: 50, Y: 25}})
	if g.Viewports["page/g"].Maximum != (runtime.ViewportState{X: 200, Y: 100}) || boxes(g.Root)["page/g/a"].Rect != (Rect{-50, -25, 400, 200}) {
		t.Fatalf("group scroll: %+v", g)
	}
	g = snapshotGeometry(t, `[a=button("A") {scale-x=0.75},b=button("B") {x=1fr,min-x=0.5}] {scale=1,gap=0,overflow-x=scroll}`, nil)
	near(t, boxes(g.Root)["page/a"].Rect.W, 150)
	near(t, boxes(g.Root)["page/b"].Rect.W, 100)
	near(t, g.Viewports["page"].Content.W, 250)
	near(t, g.Viewports["page"].Maximum.X, 50)
	g = snapshotGeometry(t, `[a=button("A") {scale-y=0.75};b=button("B") {y=1fr,min-y=0.5}] {scale=1,gap=0,overflow-y=scroll}`, nil)
	near(t, boxes(g.Root)["page/a"].Rect.H, 75)
	near(t, boxes(g.Root)["page/b"].Rect.H, 50)
	near(t, g.Viewports["page"].Content.H, 125)
	near(t, g.Viewports["page"].Maximum.Y, 25)
}
