package fynehost

import (
	"errors"
	"fyne.io/fyne/v2"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"strings"
	"testing"
)

func TestDynamicSurfaceParentOwnsModalCanvas(t *testing.T) {
	h, w := lifecycleHost(t)
	var statuses []error
	h.OnStatus = func(e error) { statuses = append(statuses, e) }
	doc, e := parser.Parse(`sdui 0.3;page=[open=button("Parent",effect="open",target="longparent");longparent=dialog("Parent",modal=false)[openChild=button("Child",effect="open",target="x")] {scale-x=0.6,scale-y=0.6};x=dialog("Child")[field=input("Field")] {scale-x=0.3,scale-y=0.3}] {x=fill,y=fill};`)
	if e != nil {
		t.Fatal(e)
	}
	e = h.Adopt(DocumentRequest{Document: doc, Entry: "page", SessionID: "parent", SourceRevision: "parent", Sequence: 1, Mode: preparation.Prototype})
	if e != nil {
		t.Fatal(e)
	}
	w.SetContent(h.Container)
	w.Resize(fyne.NewSize(800, 500))
	w.Show()
	b := h.Current()
	b.Controls()["page/open"].(*commandButton).Tapped(&fyne.PointEvent{})
	p := b.surfaces["page/longparent"]
	if p == nil {
		t.Fatal("parent not opened")
	}
	b.Controls()["page/longparent/openChild"].(*commandButton).Tapped(&fyne.PointEvent{})
	c := b.surfaces["page/x"]
	if c == nil {
		t.Fatalf("child not opened: %v", statuses)
	}
	if c.canvas != p.canvas {
		t.Errorf("modal child attached to main=%v instead of actual nonmodal parent canvas", c.canvas == w.Canvas())
	}
	if c.size.W != p.size.W*0.3 || c.size.H != p.size.H*0.3 {
		t.Fatalf("child size %+v not 30 percent actual parent %+v", c.size, p.size)
	}
	if w.Canvas().Overlays().Top() != nil || p.canvas.Overlays().Top() == nil {
		t.Fatal("modal overlay mounted on wrong window")
	}
	// Prospective malformed parent state must reject; it cannot silently use main.
	snapshot := b.Session.Snapshot()
	child := snapshot.Surfaces["page/x"]
	stale := *child.ParentSurface
	stale.OpenGeneration++
	child.ParentSurface = &stale
	snapshot.Surfaces["page/x"] = child
	if _, err := b.canvasSizes(snapshot, &layout.Engine{}); err == nil || !strings.Contains(err.Error(), "missing exact parent") {
		t.Fatal("missing parent accepted", err)
	}
	snapshot = b.Session.Snapshot()
	parent := snapshot.Surfaces["page/longparent"]
	target := snapshot.Surfaces["page/x"].Target
	parent.ParentSurface = &target
	snapshot.Surfaces["page/longparent"] = parent
	if _, err := orderedSurfacePaths(snapshot); err == nil || !strings.Contains(err.Error(), "cyclic") {
		t.Fatal("parent cycle accepted", err)
	}
}

func TestNativeOwnerCloseRevokesModalChildWithoutGateOrAcceptance(t *testing.T) {
	h, w := lifecycleHost(t)
	doc, err := parser.Parse(`sdui 0.3;page=[open=button("Parent",effect="open",target="longparent");longparent=dialog("Parent",modal=false)[parentInput=input("Parent input",value="parent saved");openChild=button("Child",effect="open",target="x")] {scale-x=0.6,scale-y=0.6};x=dialog("Child")[field=input("Field",value="child saved")] {scale-x=0.3,scale-y=0.3};other=dialog("Unrelated",modal=false)[otherInput=input("Other")] {scale-x=0.3,scale-y=0.3}] {x=fill,y=fill};`)
	if err != nil {
		t.Fatal(err)
	}
	if err = h.Adopt(DocumentRequest{Document: doc, Entry: "page", SessionID: "wm-parent", SourceRevision: "parent", Sequence: 1, Mode: preparation.Prototype}); err != nil {
		t.Fatal(err)
	}
	w.SetContent(h.Container)
	w.Resize(fyne.NewSize(800, 500))
	w.Show()
	b := h.Current()
	b.Controls()["page/open"].(*commandButton).Tapped(&fyne.PointEvent{})
	parent := b.surfaces["page/longparent"]
	calls := 0
	handle, _ := b.Session.Surface("page/longparent")
	if err = h.Mutate(func(s *ui.Session) error {
		return s.BindInteraction(handle, func(ui.Event) (ui.InteractionReply, error) {
			calls++
			return ui.InteractionReply{Domain: ui.DomainUnknown}, errors.New("unknown remote outcome")
		})
	}); err != nil {
		t.Fatal(err)
	}
	b.dispatchDialog(parent.target, ui.Accept)
	b.Controls()["page/longparent/parentInput"].(*Input).SetText("parent draft")
	b.Controls()["page/longparent/openChild"].(*commandButton).Tapped(&fyne.PointEvent{})
	child := b.surfaces["page/x"]
	if child == nil {
		t.Fatal("child not opened")
	}
	b.Controls()["page/x/field"].(*Input).SetText("child draft")
	if _, err = b.Session.CaptureDialog(parent.target, ui.Close); err == nil {
		t.Fatal("source Close bypassed modal child guard")
	}
	var results []ui.DialogResult
	h.OnDialogResult = func(r ui.DialogResult) { results = append(results, r) }
	gateCalls := 0
	b.request.PrepareResources = func(ui.Snapshot) error { gateCalls++; return errors.New("resources unavailable") }
	parent.nativeCloseRequested() // exact callback installed in SetCloseIntercept
	if gateCalls != 0 || calls != 1 {
		t.Fatal("WM close invoked gate or acceptance", gateCalls, calls)
	}
	if len(results) != 2 || len(b.surfaces) != 0 || !parent.retiring || !child.retiring {
		t.Fatal("native parent cascade incomplete", results)
	}
	for _, r := range results {
		if r.Kind != "close" || r.Reason != "parent-closed" || r.Sequence != 0 {
			t.Fatal("incorrect native lifetime receipt", r)
		}
		if r.Surface == parent.target && (r.Domain != ui.DomainUnknown || r.AcceptSequence == 0) {
			t.Fatal("lost prior domain outcome", r)
		}
	}
	for _, path := range []string{"page/longparent/parentInput", "page/x/field"} {
		value, _ := b.Session.Widget(path)
		if value.Dirty || value.Draft != value.Value {
			t.Fatal("native close retained draft", value)
		}
	}
	parent.nativeCloseRequested()
	parent.nativeClosed()
	if len(results) != 2 {
		t.Fatal("duplicate callback emitted another receipt")
	}
	// Once a replacement opens, obsolete native callbacks may not revoke it.
	b.request.PrepareResources = nil
	b.Controls()["page/open"].(*commandButton).Tapped(&fyne.PointEvent{})
	replacement := b.surfaces["page/longparent"]
	if replacement == nil || replacement.target == parent.target {
		t.Fatal("fresh opening missing")
	}
	parent.nativeCloseRequested()
	parent.nativeClosed()
	if b.surfaces["page/longparent"] != replacement || replacement.retiring || len(results) != 2 {
		t.Fatal("obsolete WM callback revoked replacement")
	}

	// Closing an unfocused owner cannot change another root nonmodal's focus.
	otherHandle, _ := b.Session.Surface("page/other")
	mainOpener, _ := b.Session.Widget("page/open")
	if err = h.Mutate(func(s *ui.Session) error {
		_, e := s.OpenSurfaceFrom(otherHandle, mainOpener.Handle, ui.ContextTarget{})
		return e
	}); err != nil {
		t.Fatal(err)
	}
	other := b.surfaces["page/other"]
	if other == nil {
		t.Fatal("unrelated surface missing")
	}
	beforeFocus := b.Session.Focused()
	beforeNative := other.canvas.Focused()
	replacement.nativeCloseRequested()
	if active := b.Session.Snapshot().ActiveSurface; active == nil || *active != other.target || b.Session.Focused() != beforeFocus || other.canvas.Focused() != beforeNative {
		t.Fatal("background owner close stole unrelated focus")
	}
}
