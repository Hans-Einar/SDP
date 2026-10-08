package fynehost

import (
	"errors"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func textRequest(t *testing.T, source string) DocumentRequest {
	t.Helper()
	d, err := parser.Parse(source)
	if err != nil {
		t.Fatal(err)
	}
	return DocumentRequest{Document: d, Entry: "page", SessionID: "text", Sequence: 1, SourceRevision: "text", Mode: preparation.Prototype}
}
func textHost(t *testing.T, args string) (*DocumentHost, *Bundle, *textEntry) {
	t.Helper()
	h := documentHost(t)
	windowForCanvas(h.canvas).Resize(fyne.NewSize(800, 500))
	h.canvas.SetContent(h.Container)
	if err := h.Adopt(textRequest(t, `sdui 0.3; page=[edit=input("Text", `+args+`) {x=fill,y=fill}] {x=fill,y=fill};`)); err != nil {
		t.Fatal(err)
	}
	b := h.Current()
	return h, b, b.Controls()["page/edit"].(*textControl).entry
}
func TestExtendedTextSubmissionAndHistory(t *testing.T) {
	_, b, e := textHost(t, `multiline=false`)
	calls := countCommits(t, b, "edit")
	changes := 0
	b.Session.ObserveChanges(e.owner.state.Target.Handle, func(ui.FieldChange) { changes++ })
	e.TypedRune('a')
	e.TypedRune('b')
	e.TypedShortcut(&fyne.ShortcutSelectAll{})
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	if *calls != 1 || changes != 2 || e.Text != "ab" || fieldAt(t, b, "edit").Accepted.Text != "ab" {
		t.Fatal("Enter mutated selection/failed commit", e.Text, changes, *calls)
	}
	e.TypedShortcut(&fyne.ShortcutUndo{})
	if e.Text != "" || *calls != 1 {
		t.Fatal("self echo cleared native undo", e.Text, *calls)
	}
	e.TypedShortcut(&fyne.ShortcutRedo{})
	if e.Text != "ab" || *calls != 1 {
		t.Fatal("redo committed", e.Text, *calls)
	}
	e.TypedShortcut(&fyne.ShortcutUndo{})
	e.TypedShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyZ, Modifier: fyne.KeyModifierShortcutDefault | fyne.KeyModifierShift})
	if e.Text != "ab" || *calls != 1 {
		t.Fatal("GLFW PrimaryShiftZ failed redo", e.Text, *calls)
	}
}
func TestExtendedMultilineKeys(t *testing.T) {
	_, b, e := textHost(t, `multiline=true`)
	calls := countCommits(t, b, "edit")
	e.TypedRune('a')
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	e.KeyDown(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	e.KeyUp(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
	if e.Text != "a\n\n" || *calls != 0 || e.AcceptsTab() {
		t.Fatal("multiline key policy", e.Text, *calls)
	}
	e.TypedShortcut(&fyne.ShortcutSelectAll{})
	e.TypedShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyReturn, Modifier: fyne.KeyModifierShortcutDefault})
	if e.Text != "a\n\n" || *calls != 1 {
		t.Fatal("PrimaryEnter erased selection or failed", e.Text, *calls)
	}
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyTab})
	e.FocusLost()
	if *calls != 1 || e.Text != "a\n\n" || e.shift {
		t.Fatal("Tab/blur changed draft")
	}
}

type changingClipboard struct {
	reads int
	value string
}

func (c *changingClipboard) Content() string {
	c.reads++
	if c.reads > 1 {
		return "MUTATED\n"
	}
	return c.value
}
func (c *changingClipboard) SetContent(v string) { c.value = v }
func TestExtendedTextClipboardAndReadOnly(t *testing.T) {
	h, b, e := textHost(t, `placeholder=""`)
	e.TypedRune('a')
	clip := &changingClipboard{value: "\r\n"}
	e.TypedShortcut(&fyne.ShortcutPaste{Clipboard: clip})
	if e.Text != "a" || clip.reads != 1 {
		t.Fatal("singleline paste normalized/multiread", e.Text, clip.reads)
	}
	e.TypedShortcut(&fyne.ShortcutUndo{})
	if e.Text != "" {
		t.Fatal("predelegation refusal cleared history", e.Text)
	}
	clip = &changingClipboard{value: "世界"}
	e.TypedShortcut(&fyne.ShortcutPaste{Clipboard: clip})
	if e.Text != "世界" || clip.reads != 1 {
		t.Fatal("validated clipboard reread", e.Text, clip.reads)
	}
	e.TypedShortcut(&fyne.ShortcutPaste{Clipboard: &capturedTextClipboard{text: "!"}}) // Paste adds a separate native undo action.
	e.TypedShortcut(&fyne.ShortcutUndo{})
	if e.Text != "世界" {
		t.Fatal("readonly test setup lost initial paste", e.Text)
	}
	f := fieldAt(t, b, "edit")
	if err := h.Mutate(func(s *ui.Session) error {
		return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: f.Target.Handle, Property: ui.ReadOnly, Value: ui.Bool(true), ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}})
	}); err != nil {
		t.Fatal(err)
	}
	e.TypedShortcut(&fyne.ShortcutSelectAll{})
	copy := &capturedTextClipboard{}
	e.TypedShortcut(&fyne.ShortcutCopy{Clipboard: copy})
	for _, sh := range []fyne.Shortcut{&fyne.ShortcutCut{Clipboard: copy}, &fyne.ShortcutPaste{Clipboard: clip}, &fyne.ShortcutUndo{}, &fyne.ShortcutRedo{}, &desktop.CustomShortcut{KeyName: fyne.KeyZ, Modifier: fyne.KeyModifierShortcutDefault | fyne.KeyModifierShift}, &desktop.CustomShortcut{KeyName: fyne.KeyBackspace, Modifier: fyne.KeyModifierShortcutDefault}} {
		e.TypedShortcut(sh)
	}
	e.TypedRune('x')
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyDelete})
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	if e.Disabled() || e.Text != "世界" || copy.text != "世界" {
		t.Fatal("readonly copy/edit guards", e.Text, copy.text)
	}
}
func TestExtendedTextRejectedEditRestoresDraftResetsHistory(t *testing.T) {
	h := documentHost(t)
	windowForCanvas(h.canvas).Resize(fyne.NewSize(800, 500))
	h.canvas.SetContent(h.Container)
	reject := false
	r := textRequest(t, `sdui 0.3; page=[edit=input("", multiline=false) {x=fill,y=fill}] {x=fill,y=fill};`)
	r.PrepareResources = func(ui.Snapshot) error {
		if reject {
			return errors.New("reject-edit")
		}
		return nil
	}
	if err := h.Adopt(r); err != nil {
		t.Fatal(err)
	}
	b := h.Current()
	e := b.Controls()["page/edit"].(*textControl).entry
	b.canvasFor(e.owner.path).Focus(e)
	calls := countCommits(t, b, "edit")
	changes := 0
	b.Session.ObserveChanges(e.owner.state.Target.Handle, func(ui.FieldChange) { changes++ })
	e.TypedRune('a')
	before := b.Session.StateRevision
	paint := b.presentation
	reject = true
	e.TypedRune('b')
	reject = false
	if e.Text != "a" || fieldAt(t, b, "edit").Proposed.Text != "a" || b.Session.StateRevision != before || b.presentation != paint || *calls != 0 || changes != 1 || b.canvasFor(e.owner.path).Focused() != e {
		t.Fatalf("rejected edit leaked text=%q changes=%d calls=%d state=%d/%d paintSame=%v focus=%T same=%v field=%+v", e.Text, changes, *calls, b.Session.StateRevision, before, b.presentation == paint, b.canvasFor(e.owner.path).Focused(), b.canvasFor(e.owner.path).Focused() == e, fieldAt(t, b, "edit"))
	}
	e.TypedShortcut(&fyne.ShortcutUndo{})
	e.TypedShortcut(&fyne.ShortcutRedo{})
	if e.Text != "a" {
		t.Fatal("rejection retained corrupt history", e.Text)
	}
	e.TypedRune('c')
	e.TypedShortcut(&fyne.ShortcutUndo{})
	if e.Text != "a" {
		t.Fatal("subsequent native edit broken", e.Text)
	}
}
func TestExtendedTextProgrammaticReplacementAndFailedCommit(t *testing.T) {
	h, b, e := textHost(t, `required=false`)
	replace := func(value string) {
		t.Helper()
		f := fieldAt(t, b, "edit")
		if err := h.Mutate(func(s *ui.Session) error {
			return s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: f.Target.Handle, Property: ui.AcceptedValue, Value: ui.Text(value), AcceptDraft: f.Dirty && value == f.Proposed.Text, ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}})
		}); err != nil {
			t.Fatal(err)
		}
	}
	e.TypedRune('a')
	replace("a")
	e.TypedShortcut(&fyne.ShortcutUndo{})
	if e.Text != "" {
		t.Fatal("identical Apply reset history", e.Text)
	}
	e.TypedShortcut(&fyne.ShortcutRedo{})
	replace("new")
	e.TypedShortcut(&fyne.ShortcutUndo{})
	if e.Text != "new" {
		t.Fatal("replacement kept history", e.Text)
	}
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnd})
	e.TypedRune('x')
	b.Session.Bind(e.owner.state.Target.Handle, func(ui.Event) ([]ui.Update, error) { return nil, errors.New("domain-failed") })
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	e.TypedShortcut(&fyne.ShortcutUndo{})
	if e.Text != "new" {
		t.Fatal("failed Commit lost undo", e.Text)
	}
}
func TestExtendedTextNativeMetricsScrollAndSync(t *testing.T) {
	h, b, e := textHost(t, `multiline=true`)
	b.canvasFor(e.owner.path).Focus(e)
	e.TypedShortcut(&fyne.ShortcutPaste{Clipboard: &capturedTextClipboard{text: strings.Repeat("wrapped text and words\n", 80)}})
	scroll := e.nativeScroll()
	if scroll == nil {
		t.Fatal("missing actual native scroll")
	}
	e.TypedKey(&fyne.KeyEvent{Name: fyne.KeyUp})
	scroll.ScrollToOffset(fyne.Position{})
	rect := inspectionRect(e)
	pos := fyne.NewPos(float32(rect.X+rect.W/2), float32(rect.Y+rect.H/2))
	outer := 0
	b.view.controls[e.owner.path].clip.onScroll = func(*fyne.ScrollEvent) { outer++ }
	test.Scroll(b.canvasFor(e.owner.path), pos, 0, -200)
	if scroll.Offset.Y <= 0 || outer != 0 {
		t.Fatal("native wheel target", scroll.Offset, outer)
	}
	offset := scroll.Offset
	e.FocusLost()
	b.muted = true
	e.FocusGained()
	b.muted = false
	if scroll.Offset != offset {
		t.Fatal("focus restoration lost native offset", offset, scroll.Offset)
	}
	if err := h.Mutate(func(s *ui.Session) error { return s.Focus(e.owner.state.Target.Handle) }); err != nil {
		t.Fatal(err)
	}
	if scroll.Offset != offset || b.Controls()[e.owner.path].(*textControl).entry != e {
		t.Fatal("sync lost native offset/object", offset, scroll.Offset)
	}
	for i := 0; i < 4; i++ {
		test.Scroll(b.canvasFor(e.owner.path), pos, 0, -10000)
	}
	if outer != 0 {
		t.Fatal("wheel chained at native limit")
	}
	fields, _ := b.inspectFields()
	if fields[e.owner.path]["text"] != e.Text {
		t.Fatal("inspector reconstructed text")
	}
	m := &collectionMeasure{}
	var n *parser.Instance
	b.presentation.snapshot.Root.Walk(func(x *parser.Instance) {
		if x.Path == e.owner.path {
			n = x
		}
	})
	f := fieldAt(t, b, "edit")
	base, err := m.MeasureField(n, f, 14, layout.Size{})
	if err != nil {
		t.Fatal(err)
	}
	n.Arguments["value"] = parser.Literal{Kind: "string", Value: strings.Repeat("long line\n", 200)}
	after, err := m.MeasureField(n, f, 14, layout.Size{})
	if err != nil {
		t.Fatal(err)
	}
	if after.Minimum != base.Minimum {
		t.Fatal("text grew native minimum", base.Minimum, after.Minimum)
	}
}
