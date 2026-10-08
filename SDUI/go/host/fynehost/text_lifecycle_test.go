package fynehost

import (
	"errors"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func TestExtendedTextReentrantRejectedEditUsesCurrentDraft(t *testing.T) {
	_, b, e := textHost(t, `multiline=false`)
	armed := false
	err := b.Session.ValidateFieldWith(e.owner.state.Target.Handle, func(f ui.FieldState) ui.FieldValidation {
		if armed {
			armed = false
			if _, err := b.Session.EditField(f.Target.Handle, b.Session.Revision, ui.Text("newer")); err != nil {
				t.Error(err)
			}
		}
		return ui.FieldValidation{}
	})
	if err != nil {
		t.Fatal(err)
	}
	armed = true
	e.TypedRune('x')
	if e.Text != "newer" || fieldAt(t, b, "edit").Proposed.Text != "newer" {
		t.Fatal("restored stale precallback text", e.Text)
	}
	e.TypedShortcut(&fyne.ShortcutUndo{})
	if e.Text != "newer" {
		t.Fatal("rejected outer action survived in history")
	}
}
func TestExtendedTextInvalidDraftRetainsHistory(t *testing.T) {
	_, b, e := textHost(t, `required=true, value="ok"`)
	calls := countCommits(t, b, "edit")
	e.TypedShortcut(&fyne.ShortcutSelectAll{})
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
	f := fieldAt(t, b, "edit")
	if e.Text != "" || !f.Dirty || f.Validation.Code == "" {
		t.Fatal("required draft not retained", f)
	}
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	e.TypedShortcut(&fyne.ShortcutUndo{})
	if e.Text != "ok" || *calls != 0 {
		t.Fatal("invalid draft/failedCommit reset history", e.Text, *calls)
	}
}
func TestExtendedTextPageCollapseAndFailedReloadRetainEntry(t *testing.T) {
	h := documentHost(t)
	windowForCanvas(h.canvas).Resize(fyne.NewSize(800, 500))
	h.canvas.SetContent(h.Container)
	source := strings.Replace(paneSource, `input("First",value="first")`, `input("",value="first",multiline=true)`, 1)
	r := textRequest(t, source)
	if err := h.Adopt(r); err != nil {
		t.Fatal(err)
	}
	b := h.Current()
	path := "page/panes/tabs/one/edit"
	e := b.Controls()[path].(*textControl).entry
	b.canvasFor(path).Focus(e)
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnd})
	e.TypedRune('x')
	e.TypedShortcut(&fyne.ShortcutSelectAll{})
	selected := e.SelectedText()
	row, col := e.CursorRow, e.CursorColumn
	header := b.Controls()["page/panes/tabs"].(*paneHeader)
	header.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnd})
	header.TypedKey(&fyne.KeyEvent{Name: fyne.KeyHome})
	divider := b.Controls()["page/panes"].(*paneDivider)
	divider.TypedShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyHome, Modifier: fyne.KeyModifierControl})
	divider.TypedKey(&fyne.KeyEvent{Name: fyne.KeySpace})
	if b.Controls()[path].(*textControl).entry != e || e.Text != "firstx" || e.SelectedText() != selected || e.CursorRow != row || e.CursorColumn != col {
		t.Fatal("page/collapse lost editing state", e.Text, e.SelectedText(), e.CursorRow, e.CursorColumn)
	}
	r.Sequence = 2
	r.SourceRevision = "rejected"
	r.PrepareResources = func(ui.Snapshot) error { return errors.New("reload-resource") }
	if _, err := h.Prepare(r); err == nil {
		t.Fatal("reload should fail")
	}
	if h.Current() != b || b.Controls()[path].(*textControl).entry != e || e.SelectedText() != selected {
		t.Fatal("failed reload replaced editor")
	}
	e.TypedShortcut(&fyne.ShortcutUndo{})
	if e.Text != "first" {
		t.Fatal("failed reload lost history", e.Text)
	}
	r.Sequence = 3
	r.SourceRevision = "accepted"
	r.PrepareResources = nil
	if err := h.Adopt(r); err != nil {
		t.Fatal(err)
	}
	next := h.Current().Controls()[path].(*textControl).entry
	if next == e || next.Text != "first" {
		t.Fatal("successful reload policy", next.Text)
	}
	e.TypedRune('z')
	if next.Text != "first" {
		t.Fatal("retired edit affected new control")
	}
}
func TestExtendedTextContextPasteAndLateReadonlyGuard(t *testing.T) {
	h, b, e := textHost(t, `multiline=false`)
	e.TypedRune('a')
	clipboard := fyne.CurrentApp().Clipboard()
	clipboard.SetContent("bad\r\npaste")
	selectAction := func(label string) {
		t.Helper()
		menu := e.owner.editMenu
		if menu == nil {
			t.Fatal("no editing menu")
		}
		for i, item := range menu.model.Items {
			if item.Label == label {
				menu.selection(func() { menu.popup.Items[i].(fyne.Tappable).Tapped(&fyne.PointEvent{}) })
				return
			}
		}
		t.Fatal("missing menu action", label)
	}
	e.owner.openEditMenu(fyne.Position{})
	selectAction("Paste")
	if e.Text != "a" {
		t.Fatal("context paste flattened newline", e.Text)
	}
	e.TypedShortcut(&fyne.ShortcutUndo{})
	if e.Text != "" {
		t.Fatal("context refusal reset history")
	}
	e.TypedRune('b')
	e.owner.openEditMenu(fyne.Position{})
	f := fieldAt(t, b, "edit")
	if err := h.Mutate(func(s *ui.Session) error {
		return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: f.Target.Handle, Property: ui.ReadOnly, Value: ui.Bool(true), ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}})
	}); err != nil {
		t.Fatal(err)
	}
	selectAction("Undo")
	if e.Text != "b" {
		t.Fatal("late editing menu bypassed readonly")
	}
	e.owner.openEditMenu(fyne.Position{})
	for _, item := range e.owner.editMenu.model.Items {
		if item.Label != "Copy" && item.Label != "Select all" {
			t.Fatal("readonly mutation affordance", item.Label)
		}
	}
	e.owner.editMenu.close()
}
