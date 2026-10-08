package fynehost

import (
	"context"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	fynetest "fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"reflect"
	"testing"
	"time"
)

const paneSource = `sdui 0.3; page=[panes=split(axis="horizontal",proportion=0.6,minFirst=0.15,minSecond=0.2)[tabs=tabs("Workspace")[one=page("One")[edit=input("First",value="first") {x=fill}];two=page("Two")[second=input("Second",value="second") {x=fill}]];aside=[ok=button("OK")]] {x=fill,y=fill};preview=input("Preview") {x=fill}] {x=fill,y=fill};`

func paneRequest(t *testing.T) DocumentRequest {
	t.Helper()
	d, err := parser.Parse(paneSource)
	if err != nil {
		t.Fatal(err)
	}
	return DocumentRequest{Document: d, Entry: "page", SessionID: "panes", SourceRevision: "panes-source", Sequence: 1, Mode: preparation.Prototype}
}
func TestPaneNativeKeyboardSelectionAndSplitProjection(t *testing.T) {
	h := documentHost(t)
	if err := h.Adopt(paneRequest(t)); err != nil {
		t.Fatal(err)
	}
	b := h.Current()
	header := b.Controls()["page/panes/tabs"].(*paneHeader)
	header.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	state := b.Session.Snapshot().Tabs["page/panes/tabs"]
	if state.Selected != "two" || header.selected != "two" {
		t.Fatal("header/runtime mismatch", state, header.selected)
	}
	header.TypedKey(&fyne.KeyEvent{Name: fyne.KeyHome})
	if b.Session.Snapshot().Tabs["page/panes/tabs"].Selected != "one" {
		t.Fatal("Home did not select first")
	}
	divider := b.Controls()["page/panes"].(*paneDivider)
	divider.TypedShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyHome, Modifier: fyne.KeyModifierControl})
	if s := b.Session.Snapshot().Splits["page/panes"]; s.Collapsed != ui.SplitFirst || b.Geometry().Splits["page/panes"].First.W != 0 {
		t.Fatal("collapse not projected", s, b.Geometry().Splits["page/panes"])
	}
	divider.TypedKey(&fyne.KeyEvent{Name: fyne.KeySpace})
	if b.Session.Snapshot().Splits["page/panes"].Collapsed != ui.SplitNone {
		t.Fatal("restore failed")
	}
}
func TestFailedPaneCallbackNeverPromotesProbe(t *testing.T) {
	h := documentHost(t)
	if err := h.Adopt(paneRequest(t)); err != nil {
		t.Fatal(err)
	}
	b := h.Current()
	tabs := b.Session.Snapshot().Tabs["page/panes/tabs"]
	calls := 0
	if err := b.Session.BindInteraction(tabs.Handle, func(ui.Event) (ui.InteractionReply, error) {
		calls++
		return ui.InteractionReply{}, errors.New("rejected")
	}); err != nil {
		t.Fatal(err)
	}
	before := b.presentation
	sequence := b.Session.Sequence()
	b.Controls()["page/panes/tabs"].(*paneHeader).TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	if calls != 1 || b.Session.Sequence() != sequence+1 {
		t.Fatal("callback/sequence", calls, b.Session.Sequence())
	}
	if b.Session.Snapshot().Tabs["page/panes/tabs"].Selected != "one" || b.presentation != before || b.pending != nil {
		t.Fatal("rejected probe became native presentation")
	}
}
func TestReentrantAcceptedDraftSurvivesPaneConflict(t *testing.T) {
	h := documentHost(t)
	if err := h.Adopt(paneRequest(t)); err != nil {
		t.Fatal(err)
	}
	b := h.Current()
	tabs := b.Session.Snapshot().Tabs["page/panes/tabs"]
	preview, _ := b.Session.Widget("page/preview")
	if err := b.Session.BindInteraction(tabs.Handle, func(ui.Event) (ui.InteractionReply, error) {
		if err := b.Session.Draft(preview.Handle, "accepted reentrant draft"); err != nil {
			return ui.InteractionReply{}, err
		}
		return ui.InteractionReply{Domain: ui.DomainSucceeded}, nil
	}); err != nil {
		t.Fatal(err)
	}
	b.Controls()["page/panes/tabs"].(*paneHeader).TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	if b.Session.Snapshot().Tabs["page/panes/tabs"].Selected != "one" {
		t.Fatal("conflicting selection published")
	}
	if b.Controls()["page/preview"].(*Input).Text != "accepted reentrant draft" {
		t.Fatal("accepted inner ticket lost")
	}
}

func TestPaneHeaderNativePointerFocusAndInspect(t *testing.T) {
	h, w := lifecycleHost(t)
	if err := h.Adopt(paneRequest(t)); err != nil {
		t.Fatal(err)
	}
	w.SetContent(h.Container)
	w.Resize(fyne.NewSize(800, 500))
	w.Show()
	b := h.Current()
	header := b.Controls()["page/panes/tabs"].(*paneHeader)
	buttons := header.tabs.buttons()
	if len(buttons) != 2 {
		t.Fatal("native headers unavailable", len(buttons))
	}
	buttons[1].(fyne.Tappable).Tapped(&fyne.PointEvent{})
	if b.Session.Snapshot().Tabs["page/panes/tabs"].Selected != "two" || w.Canvas().Focused() != header {
		t.Fatal("native header pointer did not select/focus", w.Canvas().Focused())
	}
	header.TypedKey(&fyne.KeyEvent{Name: fyne.KeyTab})
	second := b.Controls()["page/panes/tabs/two/second"].(*Input)
	if w.Canvas().Focused() != second {
		t.Fatal("Tab did not enter selected page", w.Canvas().Focused())
	}
	second.SetText("retained draft")
	header.TypedKey(&fyne.KeyEvent{Name: fyne.KeyHome})
	header.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnd})
	header.TypedKey(&fyne.KeyEvent{Name: fyne.KeyTab})
	if w.Canvas().Focused() != second || second.Text != "retained draft" {
		t.Fatal("page lost native remembered focus/draft")
	}
	controls := b.Inspect()["controls"].(map[string]map[string]any)
	input := controls["page/panes/tabs/two/second"]
	if input["visible"] != true || input["clip"].(layout.Rect).H <= 0 || controls["page/panes/tabs/one/edit"]["visible"] != false {
		t.Fatal("control inspection lost active page geometry", controls)
	}
	inspected := b.Inspect()["tabs"].(map[string]map[string]any)["page/panes/tabs"]
	pages := inspected["pages"].([]map[string]any)
	if len(pages) != 2 {
		t.Fatal("missing header inspection")
	}
	for _, p := range pages {
		r := p["clip"].(layout.Rect)
		if r.W <= 0 || r.H <= 0 {
			t.Fatal("header not native-hit visible", p)
		}
	}
}

func TestPaneSelectedHeaderPointerFocusWithoutAction(t *testing.T) {
	h, w := lifecycleHost(t)
	if err := h.Adopt(paneRequest(t)); err != nil {
		t.Fatal(err)
	}
	w.SetContent(h.Container)
	w.Resize(fyne.NewSize(800, 500))
	w.Show()
	b := h.Current()
	header := b.Controls()["page/panes/tabs"].(*paneHeader)
	edit := b.Controls()["page/panes/tabs/one/edit"].(*Input)
	tabs := b.Session.Snapshot().Tabs["page/panes/tabs"]
	calls := 0
	if err := b.Session.BindInteraction(tabs.Handle, func(ui.Event) (ui.InteractionReply, error) {
		calls++
		return ui.InteractionReply{Domain: ui.DomainRejected}, errors.New("keep old page and focus")
	}); err != nil {
		t.Fatal(err)
	}
	w.Canvas().Focus(edit)
	edit.SetText("uncommitted raw draft")
	sequence := b.Session.Sequence()
	tapHeader := func(index int) {
		button := header.tabs.buttons()[index]
		pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(button)
		size := button.Size()
		fynetest.TapCanvas(w.Canvas(), pos.Add(fyne.NewPos(size.Width/2, size.Height/2)))
	}
	tapHeader(0)
	if w.Canvas().Focused() != header || b.Session.Focused() != tabs.Handle.Path {
		t.Fatal("selected header click did not focus tabs", w.Canvas().Focused(), b.Session.Focused())
	}
	if calls != 0 || b.Session.Sequence() != sequence || edit.Text != "uncommitted raw draft" {
		t.Fatal("selected header click dispatched or changed draft", calls, b.Session.Sequence(), edit.Text)
	}
	w.Canvas().Focus(edit)
	focused := b.Session.Focused()
	tapHeader(1)
	if calls != 1 || b.Session.Snapshot().Tabs["page/panes/tabs"].Selected != "one" || w.Canvas().Focused() != edit || b.Session.Focused() != focused || edit.Text != "uncommitted raw draft" {
		t.Fatal("rejected changed-page click changed accepted focus/page/draft", calls, w.Canvas().Focused(), b.Session.Focused())
	}
}

func TestPaneHeaderKeepsNativeTargetsAcrossUnrelatedState(t *testing.T) {
	h, w := lifecycleHost(t)
	if err := h.Adopt(paneRequest(t)); err != nil {
		t.Fatal(err)
	}
	w.SetContent(h.Container)
	w.Resize(fyne.NewSize(800, 500))
	w.Show()
	b := h.Current()
	header := b.Controls()["page/panes/tabs"].(*paneHeader)
	items := append([]*container.TabItem(nil), header.tabs.Items...)
	buttons := append([]fyne.CanvasObject(nil), header.tabs.buttons()...)
	if len(buttons) != 2 {
		t.Fatal("native headers unavailable")
	}
	positions := []fyne.Position{buttons[0].Position(), buttons[1].Position()}
	sizes := []fyne.Size{buttons[0].Size(), buttons[1].Size()}
	check := func() {
		t.Helper()
		for i, button := range header.tabs.buttons() {
			if header.tabs.Items[i] != items[i] || button != buttons[i] {
				t.Fatal("unchanged header recreated native target", i)
			}
			if button.Position() != positions[i] || button.Size() != sizes[i] {
				t.Fatal("unchanged header lost laid-out geometry", i)
			}
		}
	}
	b.Controls()["page/preview"].(*Input).SetText("unrelated draft")
	check()
	header.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	check()
	if header.tabs.Selected() != items[1] {
		t.Fatal("stable items did not project accepted selection")
	}
	pages := append([]headerPage(nil), header.pages...)
	pages[1].label = "Renamed"
	pages[1].enabled = false
	header.sync(pages, "one")
	if header.tabs.Items[1] == items[1] || header.tabs.Items[1].Text != "Renamed" || !header.tabs.Items[1].Disabled() {
		t.Fatal("changed header metadata did not rebuild native target")
	}
}

func TestPaneFinalResourceFailureKeepsAcceptedPresentation(t *testing.T) {
	for _, outcome := range []ui.DomainOutcome{ui.DomainSucceeded, ui.DomainUnknown} {
		t.Run(string(outcome), func(t *testing.T) {
			h := documentHost(t)
			req := paneRequest(t)
			fail := false
			cause := &ui.Fault{Code: "resource-failure", Message: "native resource refusal"}
			req.PrepareResources = func(s ui.Snapshot) error {
				if fail && s.Tabs["page/panes/tabs"].Selected == "two" {
					return cause
				}
				return nil
			}
			if err := h.Adopt(req); err != nil {
				t.Fatal(err)
			}
			b := h.Current()
			tabs := b.Session.Snapshot().Tabs["page/panes/tabs"]
			calls := 0
			if err := b.Session.BindInteraction(tabs.Handle, func(ui.Event) (ui.InteractionReply, error) {
				calls++
				return ui.InteractionReply{Domain: outcome}, nil
			}); err != nil {
				t.Fatal(err)
			}
			var reported []error
			h.OnStatus = func(err error) { reported = append(reported, err) }
			before := b.presentation
			sequence := b.Session.Sequence() + 1
			fail = true
			b.Controls()["page/panes/tabs"].(*paneHeader).TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
			if calls != 1 || b.Session.Snapshot().Tabs["page/panes/tabs"].Selected != "one" || b.presentation != before || b.pending != nil {
				t.Fatal("final preparation failure published speculative state", calls)
			}
			if len(reported) != 1 {
				t.Fatal("expected one native outcome report", reported)
			}
			var interaction *InteractionError
			if !errors.As(reported[0], &interaction) || interaction.Result != (ui.InteractionResult{Sequence: sequence, Status: "ui-conflict", Domain: outcome}) {
				t.Fatal("native error lost interaction outcome", reported[0])
			}
			var fault *ui.Fault
			if !errors.Is(reported[0], cause) || !errors.As(reported[0], &fault) || fault != cause {
				t.Fatal("native error hid original cause", reported[0])
			}
			want := fmt.Sprintf("interaction status=ui-conflict domain=%s sequence=%d: %s", outcome, sequence, cause)
			if reported[0].Error() != want {
				t.Fatalf("native text outcome: got %q, want %q", reported[0], want)
			}
		})
	}
}

func TestPaneHandlerErrorReportsUnknownDomain(t *testing.T) {
	h := documentHost(t)
	if err := h.Adopt(paneRequest(t)); err != nil {
		t.Fatal(err)
	}
	b := h.Current()
	cause := errors.New("domain outcome unavailable")
	tabs := b.Session.Snapshot().Tabs["page/panes/tabs"]
	if err := b.Session.BindInteraction(tabs.Handle, func(ui.Event) (ui.InteractionReply, error) {
		return ui.InteractionReply{}, cause
	}); err != nil {
		t.Fatal(err)
	}
	var reported error
	h.OnStatus = func(err error) { reported = err }
	sequence := b.Session.Sequence() + 1
	b.Controls()["page/panes/tabs"].(*paneHeader).TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	var interaction *InteractionError
	if !errors.As(reported, &interaction) || interaction.Result != (ui.InteractionResult{Sequence: sequence, Status: "ui-conflict", Domain: ui.DomainUnknown}) || !errors.Is(reported, cause) {
		t.Fatal("handler error lost unknown domain outcome/cause", reported)
	}
}

func TestPaneRetiredNativeControlsCannotDispatch(t *testing.T) {
	h := documentHost(t)
	r := paneRequest(t)
	if err := h.Adopt(r); err != nil {
		t.Fatal(err)
	}
	old := h.Current()
	header := old.Controls()["page/panes/tabs"].(*paneHeader)
	divider := old.Controls()["page/panes"].(*paneDivider)
	r.Sequence = 2
	r.SourceRevision = "next"
	if err := h.Adopt(r); err != nil {
		t.Fatal(err)
	}
	b := h.Current()
	before := b.Session.Snapshot()
	header.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	divider.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnd})
	if !reflect.DeepEqual(before, b.Session.Snapshot()) {
		t.Fatal("retired native input reached successor")
	}
}

func TestPaneDisabledTabsKeepNativeSelection(t *testing.T) {
	h, w := lifecycleHost(t)
	if err := h.Adopt(paneRequest(t)); err != nil {
		t.Fatal(err)
	}
	w.SetContent(h.Container)
	w.Resize(fyne.NewSize(800, 500))
	w.Show()
	b := h.Current()
	tabs := b.Session.Snapshot().Tabs["page/panes/tabs"]
	calls := 0
	if err := b.Session.BindInteraction(tabs.Handle, func(ui.Event) (ui.InteractionReply, error) {
		calls++
		return ui.InteractionReply{}, nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := h.Mutate(func(s *ui.Session) error {
		return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: tabs.Handle, Property: ui.Enabled, Value: ui.Bool(false)}})
	}); err != nil {
		t.Fatal(err)
	}
	header := b.Controls()["page/panes/tabs"].(*paneHeader)
	before := b.Session.Snapshot()
	nativeBefore := header.tabs.Selected()
	button := header.tabs.buttons()[1]
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(button)
	size := button.Size()
	fynetest.TapCanvas(w.Canvas(), pos.Add(fyne.NewPos(size.Width/2, size.Height/2)))
	// Also guard AppTabs' eager selection if its native callback is reached
	// independently of the pointer surface (for example toolkit dispatch).
	button.(fyne.Tappable).Tapped(&fyne.PointEvent{})
	if calls != 0 || !reflect.DeepEqual(before, b.Session.Snapshot()) || header.tabs.Selected() != nativeBefore {
		t.Fatal("disabled native header diverged or dispatched", calls, header.tabs.SelectedIndex())
	}
	inspected := b.Inspect()["tabs"].(map[string]map[string]any)["page/panes/tabs"]
	if inspected["nativeSelected"] != inspected["selected"] {
		t.Fatal("disabled native and accepted selection disagree", inspected)
	}
	setEnabled := func(handle ui.Handle, enabled bool) {
		t.Helper()
		if err := h.Mutate(func(s *ui.Session) error {
			return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: handle, Property: ui.Enabled, Value: ui.Bool(enabled)}})
		}); err != nil {
			t.Fatal(err)
		}
	}
	setEnabled(tabs.Pages[1].Handle, false)
	setEnabled(tabs.Handle, true)
	header.tabs.buttons()[1].(fyne.Tappable).Tapped(&fyne.PointEvent{})
	if header.Disabled() || !header.tabs.Items[1].Disabled() || calls != 0 || header.tabs.Selected() != header.tabs.Items[0] {
		t.Fatal("re-enabling owner widened individual page eligibility", calls)
	}
	setEnabled(tabs.Pages[1].Handle, true)
	header.tabs.buttons()[1].(fyne.Tappable).Tapped(&fyne.PointEvent{})
	if calls != 1 || header.tabs.Selected() != header.tabs.Items[1] || b.Session.Snapshot().Tabs["page/panes/tabs"].Selected != "two" {
		t.Fatal("re-enabled eligible page remained unusable", calls)
	}
}

func TestPaneDisabledPageSkippedAndEmptySelection(t *testing.T) {
	h := documentHost(t)
	if err := h.Adopt(paneRequest(t)); err != nil {
		t.Fatal(err)
	}
	b := h.Current()
	tabs := b.Session.Snapshot().Tabs["page/panes/tabs"]
	disable := func(handle ui.Handle) error {
		return h.Mutate(func(s *ui.Session) error {
			return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: handle, Property: ui.Enabled, Value: ui.Bool(false)}})
		})
	}
	if err := disable(tabs.Pages[1].Handle); err != nil {
		t.Fatal(err)
	}
	before := b.Session.Sequence()
	header := b.Controls()["page/panes/tabs"].(*paneHeader)
	header.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnd})
	if b.Session.Sequence() != before || header.selected != "one" {
		t.Fatal("disabled page activated")
	}
	if err := disable(tabs.Pages[0].Handle); err != nil {
		t.Fatal(err)
	}
	if b.Session.Snapshot().Tabs["page/panes/tabs"].Selected != "" {
		t.Fatal("zero eligible pages retained active selection")
	}
	if b.view.controls["page/panes/tabs/one/edit"].clip.Visible() || b.view.controls["page/panes/tabs/two/second"].clip.Visible() {
		t.Fatal("ineligible page controls remained visible")
	}
}

func TestPaneAllHiddenPagesKeepHeaderAndEmptyBody(t *testing.T) {
	h := documentHost(t)
	if err := h.Adopt(paneRequest(t)); err != nil {
		t.Fatal(err)
	}
	b := h.Current()
	tabs := b.Session.Snapshot().Tabs["page/panes/tabs"]
	updates := []ui.Update{}
	for _, p := range tabs.Pages {
		updates = append(updates, ui.Update{Handle: p.Handle, Property: ui.Visible, Value: ui.Bool(false)})
	}
	if err := h.Mutate(func(s *ui.Session) error { return s.Apply(s.Revision, s.BatchRevision+1, updates) }); err != nil {
		t.Fatal(err)
	}
	if b.Session.Snapshot().Tabs["page/panes/tabs"].Selected != "" || b.Geometry().Tabs["page/panes/tabs"].Header.H <= 0 {
		t.Fatal("empty tabs lost header geometry")
	}
}

func TestPaneFailedActivationPreservesLoadAcceptedHideCancels(t *testing.T) {
	h := documentHost(t)
	d, err := parser.Parse(`sdui 0.3;page=[tabs=tabs("Workspace")[one=page("One")[tree=tree("Tree") {x=fill,y=fill,overflow-y=scroll}];two=page("Two")[]] {x=fill,y=fill}] {x=fill,y=fill};`)
	if err != nil {
		t.Fatal(err)
	}
	started := make(chan context.Context, 2)
	replies := make(chan ui.CollectionData, 2)
	posts := make(chan func(), 2)
	h.Post = func(f func()) { posts <- f }
	p := ui.CollectionProvider{ID: "lazy", Epoch: 1, Load: func(ctx context.Context, _ ui.LoadRequest) (ui.CollectionData, error) {
		started <- ctx
		return <-replies, nil
	}}
	req := DocumentRequest{Document: d, Entry: "page", SessionID: "page-load", SourceRevision: "page-load", Sequence: 1, Mode: preparation.Prototype, Providers: map[string]ui.CollectionProvider{"page/tabs/one/tree": p}}
	if err = h.Adopt(req); err != nil {
		t.Fatal(err)
	}
	ctx := <-started
	b := h.Current()
	tabs := b.Session.Snapshot().Tabs["page/tabs"]
	header := b.Controls()["page/tabs"].(*paneHeader)
	reject := true
	if err = b.Session.BindInteraction(tabs.Handle, func(ui.Event) (ui.InteractionReply, error) {
		if reject {
			return ui.InteractionReply{}, errors.New("keep page")
		}
		return ui.InteractionReply{}, nil
	}); err != nil {
		t.Fatal(err)
	}
	header.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	select {
	case <-ctx.Done():
		t.Fatal("failed selection canceled accepted page load")
	default:
	}
	if b.Session.Snapshot().Collections["page/tabs/one/tree"].Request == nil {
		t.Fatal("failed selection revoked request")
	}
	reject = false
	header.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	select {
	case <-ctx.Done():
	default:
		t.Fatal("accepted hide did not cancel provider context")
	}
	replies <- ui.CollectionData{Items: []ui.CollectionItem{{ID: "late", Kind: ui.Row, Label: "Late", ChildrenLoaded: true}}}
	select {
	case f := <-posts:
		f()
	case <-time.After(3 * time.Second):
		t.Fatal("completion not delivered")
	}
	header.TypedKey(&fyne.KeyEvent{Name: fyne.KeyHome})
	state := b.Session.Snapshot().Collections["page/tabs/one/tree"]
	if state.Request != nil || state.AutoLoadPending || len(state.Data.Items) != 0 {
		t.Fatal("reveal revived load or late reply", state)
	}
	select {
	case <-started:
		t.Fatal("reveal restarted canceled load")
	default:
	}
}
