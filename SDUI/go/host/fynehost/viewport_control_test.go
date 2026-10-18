package fynehost

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"math"
	"testing"
)

func TestCollectionThumbGrabAndAltScroll(t *testing.T) {
	h := documentHost(t)
	p := eagerProvider()
	for i := 0; i < 100; i++ {
		p.Initial.Items = append(p.Initial.Items, ui.CollectionItem{ID: ui.ItemID(fmt.Sprint(i)), Kind: ui.Row, Label: "Long row label", ChildrenLoaded: true})
	}
	if err := h.Adopt(documentRequest(t, 1, p)); err != nil {
		t.Fatal(err)
	}
	c := h.Current().Controls()["page/tree"].(*CollectionControl)
	v := c.viewport
	grab := math.Min(8, thumbLength(v.Rect.H, v.Content.H)/2)
	y := rowHeight(c.font) + grab
	c.Dragged(&fyne.DragEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(float32(v.Rect.W+3), float32(y))}})
	if c.viewport.Offset.Y != 0 {
		t.Fatal("stationary thumb grab jumped", c.viewport.Offset)
	}
	c.Dragged(&fyne.DragEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(float32(v.Rect.W+3), float32(y+10))}, Dragged: fyne.Delta{DY: 10}})
	want := 10 * v.Maximum.Y / (v.Rect.H - thumbLength(v.Rect.H, v.Content.H))
	if math.Abs(c.viewport.Offset.Y-want) > 0.01 {
		t.Fatalf("drag=%v want %v", c.viewport.Offset.Y, want)
	}
	c.DragEnd()
	before := c.viewport.Offset.Y
	c.TypedShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyDown, Modifier: fyne.KeyModifierAlt})
	if c.viewport.Offset.Y <= before {
		t.Fatal("Alt+Down did not scroll")
	}
}

func TestFrameBackgroundAndWrappedControlWheel(t *testing.T) {
	h := documentHost(t)
	d, err := parser.Parse(`sdui 0.3; page=[outer=[a=button("A") {scale-y=2,x=fill}] {x=fill,y=fill,overflow-y=scroll}] {x=fill,y=fill};`)
	if err != nil {
		t.Fatal(err)
	}
	req := documentRequest(t, 1, eagerProvider())
	req.Document = d
	req.Providers = nil
	if err = h.Adopt(req); err != nil {
		t.Fatal(err)
	}
	b := h.Current()
	var surface *viewportSurface
	for _, obj := range b.presentation.objects {
		if v, ok := obj.(*viewportSurface); ok && v.axis == "" {
			surface = v
			break
		}
	}
	if surface == nil {
		t.Fatal("frame has no wheel target")
	}
	surface.Scrolled(&fyne.ScrollEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(2, 2)}, Scrolled: fyne.Delta{DY: -20}})
	if got := b.Session.Snapshot().Viewports["page/outer"].Y; got != 20 {
		t.Fatal("background wheel", got)
	}
	wrapper := b.view.controls["page/outer/a"].clip
	wrapper.Scrolled(&fyne.ScrollEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(2, 2)}, Scrolled: fyne.Delta{DY: -20}})
	if got := b.Session.Snapshot().Viewports["page/outer"].Y; got != 40 {
		t.Fatal("wrapped wheel swallowed", got)
	}
}

func TestNestedNativeThumbsUseDisjointSharedGutters(t *testing.T) {
	h := documentHost(t)
	d, err := parser.Parse(`sdui 0.3; page=[outer=[tree=tree("Tree") {x=fill,scale-y=2,overflow-y=scroll}] {x=fill,y=fill,overflow-y=scroll}] {x=fill,y=fill};`)
	if err != nil {
		t.Fatal(err)
	}
	p := eagerProvider()
	for i := 0; i < 100; i++ {
		p.Initial.Items = append(p.Initial.Items, ui.CollectionItem{ID: ui.ItemID(fmt.Sprint(i)), Kind: ui.Row, Label: "Row", ChildrenLoaded: true})
	}
	req := documentRequest(t, 1, p)
	req.Document = d
	req.Providers = map[string]ui.CollectionProvider{"page/outer/tree": p}
	if err = h.Adopt(req); err != nil {
		t.Fatal(err)
	}
	b := h.Current()
	outer := b.Geometry().Viewports["page/outer"]
	c := b.Controls()["page/outer/tree"].(*CollectionControl)
	inner := c.viewport
	if outer.VerticalGutter.W != collectionGutter || outer.VerticalGutter.X != outer.Rect.X+outer.Rect.W {
		t.Fatal("native measured gutter missing", outer)
	}
	if inner.Rect.X+inner.Rect.W+collectionGutter > outer.VerticalGutter.X {
		t.Fatal("inner track intrudes into parent gutter", inner, outer)
	}
	var thumb *viewportSurface
	for _, o := range b.presentation.objects {
		if v, ok := o.(*viewportSurface); ok && v.axis == "y" && v.viewport.Path == "page/outer" {
			thumb = v
			break
		}
	}
	if thumb == nil {
		t.Fatal("missing outer native thumb")
	}
	if float64(thumb.Position().X) != outer.VerticalGutter.X {
		t.Fatal("thumb ignores measured strip")
	}
	// Each adapter updates its own runtime offset; geometry keeps their native
	// hit regions disjoint. Actual OS drag remains a separate acceptance check.
	thumb.Dragged(&fyne.DragEvent{Dragged: fyne.Delta{DY: 10}})
	offsets := b.Session.Snapshot().Viewports
	if offsets["page/outer"].Y <= 0 || offsets["page/outer/tree"].Y != 0 {
		t.Fatal("outer drag routed to child", offsets)
	}
	c.Dragged(&fyne.DragEvent{PointEvent: fyne.PointEvent{Position: fyne.NewPos(float32(c.viewport.Rect.W+3), float32(rowHeight(c.font)+15))}, Dragged: fyne.Delta{DY: 10}})
	offsets = b.Session.Snapshot().Viewports
	if offsets["page/outer/tree"].Y <= 0 {
		t.Fatal("inner drag not independent", offsets)
	}
}

func TestNativeViewportInsetsDeclaredAxesOnly(t *testing.T) {
	m := &collectionMeasure{}
	for _, kind := range []string{"frame", "group", "widget"} {
		for _, profile := range []string{"", "sdui/0.3"} {
			for _, axis := range []string{"x", "y"} {
				n := &parser.Instance{Kind: kind, Profile: profile, Layout: map[string]any{"overflow-" + axis: "scroll"}}
				got, err := m.MeasureViewport(n, 14)
				if err != nil {
					t.Fatal(err)
				}
				right, bottom := 0.0, 0.0
				if profile == "sdui/0.3" && kind != "widget" {
					if axis == "y" {
						right = collectionGutter
					} else {
						bottom = collectionGutter
					}
				}
				if got.Right != right || got.Bottom != bottom {
					t.Fatal(kind, profile, axis, got)
				}
			}
		}
	}
}

func TestCollectionPageKeysUseAncestorClippedHeight(t *testing.T) {
	h := documentHost(t)
	d, err := parser.Parse(`sdui 0.3; page=[outer=[tree=tree("Tree") {x=fill,scale-y=2,overflow-y=scroll}] {x=fill,y=fill,overflow-y=scroll}] {x=fill,y=fill};`)
	if err != nil {
		t.Fatal(err)
	}
	p := eagerProvider()
	for i := 0; i < 100; i++ {
		p.Initial.Items = append(p.Initial.Items, ui.CollectionItem{ID: ui.ItemID(fmt.Sprint(i)), Kind: ui.Row, Label: "Row", ChildrenLoaded: true})
	}
	req := documentRequest(t, 1, p)
	req.Document = d
	req.Providers = map[string]ui.CollectionProvider{"page/outer/tree": p}
	if err = h.Adopt(req); err != nil {
		t.Fatal(err)
	}
	b := h.Current()
	c := b.Controls()["page/outer/tree"].(*CollectionControl)
	if c.viewport.Clip.H >= c.viewport.Rect.H {
		t.Fatal("fixture must clip the collection through its ancestor")
	}
	want := c.viewport.Clip.H * .9
	c.TypedKey(&fyne.KeyEvent{Name: fyne.KeyPageDown})
	if got := c.viewport.Offset.Y; math.Abs(got-want) > 0.001 {
		t.Fatalf("PageDown moved %v, visible-page step is %v", got, want)
	}
	if b.Session.Snapshot().Viewports["page/outer"].Y != 0 {
		t.Fatal("PageDown unexpectedly scrolled ancestor before inner limit")
	}
	c.TypedKey(&fyne.KeyEvent{Name: fyne.KeyPageUp})
	if got := c.viewport.Offset.Y; math.Abs(got) > 0.001 {
		t.Fatalf("PageUp did not reverse visible-page step: %v", got)
	}
}
