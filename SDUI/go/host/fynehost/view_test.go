package fynehost

import (
	"math"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

func TestSharedGeometryNativeInteraction(t *testing.T) {
	a := test.NewTempApp(t)
	_, roots, e := parser.Compile(`sdui 0.2; page=[header="# Form",ok=button("OK"),edit=input("Value",value="12"),off=button("Disabled") {enabled=false}]*b {scale=1};`)
	if e != nil {
		t.Fatal(e)
	}
	v := New(roots["page"])
	w := a.NewWindow("test")
	defer w.Close()
	w.SetPadded(false)
	w.SetContent(v.Container)
	w.Resize(fyne.NewSize(640, 400))
	w.Show()
	var status error
	v.OnStatus = func(e error) { status = e }
	v.Container.Resize(fyne.NewSize(640, 400))
	if v.Geometry == nil {
		t.Fatal("no layout")
	}
	calls := 0
	draft := ""
	v.Actions["page/ok"] = func(_, value string) error { calls++; return nil }
	v.Actions["page/edit"] = func(_, value string) error { draft = value; return nil }
	ok := v.Controls["page/ok"].(*widget.Button)
	edit := v.Controls["page/edit"].(*Input)
	off := v.Controls["page/off"].(*widget.Button)
	test.Tap(ok)
	if calls != 1 || status != nil {
		t.Fatalf("click %d %v", calls, status)
	}
	test.Type(edit, "3")
	if draft == "" || draft != edit.Text {
		t.Fatalf("edit callback %q %q", draft, edit.Text)
	}
	if !off.Disabled() {
		t.Fatal("disabled inheritance")
	}
	w.Canvas().Focus(ok)
	w.Canvas().FocusNext()
	if w.Canvas().Focused() != edit {
		t.Fatalf("tab focus: %T", w.Canvas().Focused())
	}
	w.Canvas().FocusNext()
	if w.Canvas().Focused() != ok {
		t.Fatalf("disabled focus not skipped: %T", w.Canvas().Focused())
	}
	v.Geometry.Walk(func(b *layout.Box) {
		if c := v.controls[b.Path]; c != nil {
			if math.Abs(float64(c.clip.Position().X)-b.Clip.X) > .01 || math.Abs(float64(c.clip.Size().Width)-b.Clip.W) > .01 {
				t.Fatalf("native geometry diverged %s", b.Path)
			}
		}
	})
	v.Container.Resize(fyne.NewSize(800, 500))
	if calls != 1 {
		t.Fatal("resize executed callback")
	}
	if edit.Text != draft {
		t.Fatal("resize lost input")
	}
	delete(v.Actions, "page/ok")
	test.Tap(ok)
	if status == nil {
		t.Fatal("missing unbound status")
	}
	v.Close()
	test.Tap(ok)
	if calls != 1 {
		t.Fatal("callback after teardown")
	}
}

func TestNativeClippedHitTesting(t *testing.T) {
	a := test.NewTempApp(t)
	_, roots, e := parser.Compile(`sdui 0.2; page=[ok=button("OK") {scale-x=2}] {scale-x=0.3333333333,scale-y=1,overflow-x=clip};`)
	if e != nil {
		t.Fatal(e)
	}
	v := New(roots["page"])
	w := a.NewWindow("clip")
	defer w.Close()
	w.SetPadded(false)
	w.SetContent(v.Container)
	w.Resize(fyne.NewSize(300, 100))
	w.Show()
	v.Container.Resize(fyne.NewSize(300, 100))
	calls := 0
	v.Actions["page/ok"] = func(_, _ string) error { calls++; return nil }
	test.TapCanvas(w.Canvas(), fyne.NewPos(50, 15))
	test.TapCanvas(w.Canvas(), fyne.NewPos(150, 15))
	if calls != 1 {
		t.Fatalf("native clip delivered %d events", calls)
	}
}
