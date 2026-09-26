package fynehost

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/reload"
	uiruntime "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func candidate(t *testing.T, source string, seq uint64) reload.Candidate {
	t.Helper()
	d, roots, e := parser.Compile("sdui 0.2; " + source)
	if e != nil {
		t.Fatal(e)
	}
	return reload.Candidate{Sequence: seq, Root: roots["page"], Document: d}
}
func TestRuntimeReloadNativeDraftFocusAndNoReplay(t *testing.T) {
	a := test.NewTempApp(t)
	first := candidate(t, `page=[edit=input("Edit",value="1"),ok=button("OK")] {scale=1};`, 1)
	s, e := uiruntime.New("test", first.Root)
	if e != nil {
		t.Fatal(e)
	}
	w := a.NewWindow("runtime")
	defer w.Close()
	v, e := NewRuntime(s, w.Canvas())
	if e != nil {
		t.Fatal(e)
	}
	w.SetPadded(false)
	w.SetContent(v.Container)
	w.Resize(fyne.NewSize(640, 400))
	w.Show()
	v.Container.Resize(fyne.NewSize(640, 400))
	edit, _ := s.Widget("page/edit")
	calls := 0
	s.Bind(edit.Handle, func(e uiruntime.Event) ([]uiruntime.Update, error) {
		calls++
		current, _ := s.Widget(e.Handle.Path)
		return []uiruntime.Update{{Handle: e.Handle, Property: uiruntime.AcceptedValue, Value: e.Value, ExpectedValueRevision: current.ValueRevision, AcceptDraft: true}}, nil
	})
	old := v.View.Controls[edit.InstancePath].(*Input)
	w.Canvas().Focus(old)
	test.Type(old, "draft")
	state, _ := s.Widget("page/edit")
	draft := state.Draft
	if !state.Dirty || calls != 0 {
		t.Fatal("typing dispatched domain action")
	}
	next := candidate(t, `page=["Added heading";<edit=input("Renamed",value="default"),ok=button("OK")>] {scale=1};`, 2)
	if e = v.Adopt(next); e != nil {
		t.Fatal(e)
	}
	state, _ = s.Widget("page/edit")
	current := v.View.Controls[state.InstancePath].(*Input)
	if current.Text != draft || w.Canvas().Focused() != current || calls != 0 {
		t.Fatal("reload lost state/focus or replayed")
	}
	old.OnSubmitted(draft)
	if calls != 0 {
		t.Fatal("retired native callback dispatched")
	}
	current.OnSubmitted(draft)
	state, _ = s.Widget("page/edit")
	if calls != 1 || state.Dirty || state.Value != draft {
		t.Fatal(state, calls)
	}
	test.Type(current, "more")
	current.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEscape})
	if current.Text != state.Value || calls != 1 {
		t.Fatal("Escape dispatched callback or lost accepted value")
	}
	revision := s.Revision
	bad := candidate(t, `page=[button("Impossible") {scale-x=10}] {scale=1};`, 3)
	if v.Adopt(bad) == nil || s.Revision != revision || v.View.Controls[state.InstancePath] != current {
		t.Fatal("invalid measured reload replaced last good UI")
	}
	v.Close()
	current.OnSubmitted(current.Text)
	if calls != 1 {
		t.Fatal("callback after closed session")
	}
}
