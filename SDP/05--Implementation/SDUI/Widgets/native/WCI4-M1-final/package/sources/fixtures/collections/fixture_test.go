package collections

import (
	"fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"testing"
)

func TestRealConnectedFixtureAcrossReloads(t *testing.T) {
	a := test.NewTempApp(t)
	w := a.NewWindow("connected")
	defer w.Close()
	f, e := New()
	if e != nil {
		t.Fatal(e)
	}
	defer f.Close()
	h := fynehost.NewDocumentHost(w.Canvas(), layout.Size{W: 1000, H: 650})
	defer h.Close()
	var old *ui.Session
	for seq := uint64(1); seq <= 3; seq++ {
		r, e := f.Request(Source, seq, false)
		if e != nil {
			t.Fatal(e)
		}
		next, e := h.Prepare(r)
		if e != nil {
			t.Fatal(e)
		}
		if f.Calls() != int(seq-1) {
			t.Fatal("preparation invoked action")
		}
		if e = h.Commit(next); e != nil {
			t.Fatal(e)
		}
		if old != nil && !old.Closed() {
			t.Fatal("old session alive")
		}
		old = next.Session
		w, _ := next.Session.Widget("page/body/entries")
		target, e := next.Session.Target(w.Handle, "entry-001")
		if e != nil {
			t.Fatal(e)
		}
		if e = h.Mutate(func(s *ui.Session) error {
			return s.Dispatch(ui.Event{Handle: w.Handle, ModelRevision: s.Revision, Sequence: s.Sequence() + 1, Kind: ui.Select, Collection: &target})
		}); e != nil {
			t.Fatal(e)
		}
		if f.Calls() != int(seq-1) {
			t.Fatal("selection invoked action")
		}
		if e = h.Mutate(func(s *ui.Session) error {
			return s.Dispatch(ui.Event{Handle: w.Handle, ModelRevision: s.Revision, Sequence: s.Sequence() + 1, Kind: ui.Activate, Collection: &target})
		}); e != nil {
			t.Fatal(e)
		}
		preview, _ := next.Session.Widget("page/footer/preview")
		if preview.Value != "Preview: Entry 001 [entry-001]" || f.Calls() != int(seq) {
			t.Fatal("missing actual SDL preview", preview.Value, f.Calls())
		}
	}
}
