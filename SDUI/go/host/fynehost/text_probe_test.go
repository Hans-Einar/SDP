// Regression cases derived from Mendel's independent M2 host review.
package fynehost

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"strings"
	"testing"
)

func TestExtendedTextScrollPreservationAndUserReveal(t *testing.T) {
	h, b, e := textHost(t, `multiline=true`)
	b.canvasFor(e.owner.path).Focus(e)
	e.TypedShortcut(&fyne.ShortcutPaste{Clipboard: &capturedTextClipboard{text: strings.Repeat("line words\n", 80)}})
	sc := e.nativeScroll()
	if sc == nil {
		t.Fatal("native scroll missing")
	}
	sc.ScrollToOffset(fyne.NewPos(0, 100))
	off := sc.Offset
	if err := h.Mutate(func(s *ui.Session) error { return s.Focus(e.owner.state.Target.Handle) }); err != nil {
		t.Fatal(err)
	}
	if sc.Offset != off {
		t.Fatalf("ordinary sync changed scroll %v -> %v", off, sc.Offset)
	}
	sc.ScrollToOffset(fyne.Position{})
	e.TypedRune('z')
	if sc.Offset.Y <= 0 {
		t.Fatal("sync retention suppressed user caret reveal")
	}
	sc.ScrollToOffset(fyne.Position{})
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
	if sc.Offset.Y <= 0 {
		t.Fatal("arrow did not reveal offscreen caret")
	}
}
func TestExtendedTextDetachedProbesPreserveEditing(t *testing.T) {
	_, b, e := textHost(t, `multiline=true`)
	b.canvasFor(e.owner.path).Focus(e)
	e.TypedShortcut(&fyne.ShortcutPaste{Clipboard: &capturedTextClipboard{text: strings.Repeat("line words\n", 80)}})
	e.TypedShortcut(&fyne.ShortcutSelectAll{})
	sc := e.nativeScroll()
	sc.ScrollToOffset(fyne.NewPos(0, 100))
	off, text, sel, row, col := sc.Offset, e.Text, e.SelectedText(), e.CursorRow, e.CursorColumn
	paint, pending := b.presentation, b.pending
	before, _ := json.Marshal(b.Session.Snapshot())
	snap := b.Session.Snapshot()
	if _, err := b.stage(snap); err != nil {
		t.Fatal(err)
	}
	ticket, err := b.preparePresentation(snap)
	if err != nil {
		t.Fatal(err)
	}
	ticket.Discard()
	snap.Root.Walk(func(n *parser.Instance) {
		if n.Path == e.owner.path {
			n.Layout["max-x"] = 0.01
		}
	})
	if _, err = b.stage(snap); err == nil {
		t.Fatal("tiny candidate should reject")
	}
	after, _ := json.Marshal(b.Session.Snapshot())
	if string(before) != string(after) || b.presentation != paint || b.pending != pending || e.Text != text || e.SelectedText() != sel || e.CursorRow != row || e.CursorColumn != col || sc.Offset != off {
		t.Fatal("probe changed published/runtime/native editing state")
	}
	e.TypedShortcut(&fyne.ShortcutUndo{})
	if e.Text != "" {
		t.Fatal("probe reset native history")
	}
}
