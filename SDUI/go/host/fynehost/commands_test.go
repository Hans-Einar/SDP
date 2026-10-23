package fynehost

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"image/color"
	"strings"
	"testing"
)

const commandSource = `sdui 0.3; page=[flag=command("Flag",toggle=true,key="Primary+K");open=command("Settings",effect="open",target="settings");tools=<toggle=button(command="flag"),show=button(command="open")>;actions=menu("Actions")[toggleItem=item(command="flag")];edit=input("Main",value="main");settings=dialog("Settings",modal=true)[name=input("Name",value="accepted");buttons=<ok=button("OK",effect="accept"),cancel=button("Cancel",effect="cancel")>] {scale-x=0.6,scale-y=0.6}] {x=fill,y=fill};`

func commandHost(t *testing.T, nonmodal bool) (*DocumentHost, fyne.Window) {
	t.Helper()
	h, w := lifecycleHost(t)
	source := commandSource
	if nonmodal {
		source = strings.Replace(source, "modal=true", "modal=false", 1)
	}
	doc, err := parser.Parse(source)
	if err != nil {
		t.Fatal(err)
	}
	if err = h.Adopt(DocumentRequest{Document: doc, Entry: "page", SessionID: "command-test", SourceRevision: "source", Sequence: 1, Mode: preparation.Prototype}); err != nil {
		t.Fatal(err)
	}
	w.SetContent(h.Container)
	w.Resize(fyne.NewSize(800, 500))
	w.Show()
	return h, w
}
func TestCommandButtonAndMenuShareAcceptedToggle(t *testing.T) {
	h, _ := commandHost(t, false)
	b := h.Current()
	button := b.Controls()["page/tools/toggle"].(*commandButton)
	button.Tapped(&fyne.PointEvent{})
	command, _ := b.Session.Command("page/flag")
	state, _ := b.Session.CommandState(command)
	if !state.Checked || !strings.HasPrefix(button.Text, "[x]") {
		t.Fatal("button not projected", state, button.Text)
	}
	b.openMenu("page/actions", ui.ContextTarget{}, fyne.NewPos(10, 10))
	menu := b.menus["page/actions"]
	if menu == nil {
		t.Fatal("no native menu")
	}
	menu.native.popup.ActivateNext()
	menu.native.input.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
	state, _ = b.Session.CommandState(command)
	if state.Checked || !strings.HasPrefix(button.Text, "[ ]") || len(b.menus) != 0 {
		t.Fatal("menu did not use canonical toggle/dismiss", state, len(b.menus))
	}
}
func TestDialogPublicationGeometryAndCancel(t *testing.T) {
	for _, nonmodal := range []bool{false, true} {
		t.Run(map[bool]string{false: "modal", true: "nonmodal"}[nonmodal], func(t *testing.T) {
			h, _ := commandHost(t, nonmodal)
			b := h.Current()
			var results []ui.DialogResult
			h.OnDialogResult = func(r ui.DialogResult) { results = append(results, r) }
			var status error
			h.OnStatus = func(e error) { status = e }
			b.Controls()["page/tools/show"].(*commandButton).Tapped(&fyne.PointEvent{})
			native := b.surfaces["page/settings"]
			if native == nil {
				t.Fatal("opening missing", status)
			}
			frame := b.presentation.canvases["page/settings"]
			if native.content.Size() != fyne.NewSize(float32(frame.size.W), float32(frame.size.H)) {
				t.Fatal("native content disagrees with shared canvas", native.content.Size(), frame.size)
			}
			input := b.Controls()["page/settings/name"].(*Input)
			input.SetText("raw draft")
			b.Controls()["page/settings/buttons/cancel"].(*commandButton).Tapped(&fyne.PointEvent{})
			if len(results) != 1 || results[0].Kind != "cancel" || native.retiring == false || len(b.surfaces) != 0 {
				t.Fatal("cancel lifecycle", results, status)
			}
			accepted, _ := b.Session.Widget("page/settings/name")
			if accepted.Dirty || accepted.Draft != "accepted" {
				t.Fatal("cancel retained unaccepted draft", accepted)
			}
		})
	}
}
func TestDialogRejectedAcceptVisibleOutcome(t *testing.T) {
	h, _ := commandHost(t, false)
	b := h.Current()
	target, _ := b.Session.Surface("page/settings")
	if err := b.Session.BindInteraction(target, func(ui.Event) (ui.InteractionReply, error) {
		return ui.InteractionReply{Domain: ui.DomainUnknown}, errors.New("remote outcome unknown")
	}); err != nil {
		t.Fatal(err)
	}
	b.Controls()["page/tools/show"].(*commandButton).Tapped(&fyne.PointEvent{})
	b.Controls()["page/settings/buttons/ok"].(*commandButton).Tapped(&fyne.PointEvent{})
	native := b.surfaces["page/settings"]
	if native == nil {
		t.Fatal("failed acceptance closed dialog")
	}
	if !strings.Contains(native.status.Text, "unknown") || !strings.Contains(native.status.Text, "blocked") {
		t.Fatal("missing visible outcome", native.status.Text)
	}
}

func TestDeclaredShortcutForwardedFromFocusedAdapters(t *testing.T) {
	h, w := lifecycleHost(t)
	doc, err := parser.Parse(`sdui 0.3; page=[flag=command("Flag",toggle=true,key="Primary+R");tree=tree("Tree") {x=fill,y=fill,overflow-x=scroll,overflow-y=scroll};menu=menu("Context",mode="context",target="tree")[action=item(command="flag")];edit=input("Edit");split=split(axis="horizontal")[a=[leftEdit=input("A")];b=[rightEdit=input("B")]]] {x=fill,y=fill};`)
	if err != nil {
		t.Fatal(err)
	}
	if err = h.Adopt(DocumentRequest{Document: doc, Entry: "page", SessionID: "shortcut", SourceRevision: "source", Sequence: 1, Mode: preparation.Prototype, Providers: map[string]ui.CollectionProvider{"page/tree": eagerProvider()}}); err != nil {
		t.Fatal(err)
	}
	w.SetContent(h.Container)
	w.Resize(fyne.NewSize(800, 500))
	w.Show()
	b := h.Current()
	key, err := nativeShortcut("Primary+R")
	if err != nil {
		t.Fatal(err)
	}
	handle, _ := b.Session.Command("page/flag")
	for _, path := range []string{"page/tree", "page/edit", "page/split"} {
		control := b.Controls()[path]
		w.Canvas().Focus(control.(fyne.Focusable))
		before, _ := b.Session.CommandState(handle)
		control.(fyne.Shortcutable).TypedShortcut(key)
		after, _ := b.Session.CommandState(handle)
		if before.Checked == after.Checked {
			t.Fatalf("%s swallowed declared command", path)
		}
	}
	tree := b.Controls()["page/tree"].(*CollectionControl)
	w.Canvas().Focus(tree)
	tree.KeyDown(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
	tree.TypedKey(&fyne.KeyEvent{Name: fyne.KeyF10})
	menu := b.menus["page/menu"]
	if menu == nil {
		t.Fatal("Shift+F10 TypedKey did not open native context")
	}
	menu.native.input.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEscape})
	tree.KeyUp(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
	input := b.Controls()["page/edit"].(*Input)
	w.Canvas().Focus(input)
	calls := 0
	input.OnShortcut = func(fyne.Shortcut) { calls++ }
	input.TypedShortcut(&fyne.ShortcutSelectAll{})
	if calls != 0 {
		t.Fatal("native editing shortcut forwarded")
	}
}
func TestCommandButtonNativeFocusTracksSurfaceAndMutedPublication(t *testing.T) {
	h, w := commandHost(t, true)
	b := h.Current()
	button := b.Controls()["page/tools/show"].(*commandButton)
	w.Canvas().Focus(button)
	state, _ := b.Session.Widget("page/tools/show")
	if b.Session.Focused() != state.Handle.Path {
		t.Fatal("native button focus absent from runtime", b.Session.Focused())
	}
	revision := b.Session.StateRevision
	b.muted = true
	button.FocusGained()
	b.muted = false
	if b.Session.StateRevision != revision {
		t.Fatal("publication focus synthesized state")
	}
	button.Tapped(&fyne.PointEvent{})
	cancel := b.Controls()["page/settings/buttons/cancel"].(*commandButton)
	b.canvasFor("page/settings/buttons/cancel").Focus(cancel)
	if b.Session.Snapshot().ActiveSurface == nil {
		t.Fatal("auxiliary button focus did not activate surface")
	}
}

func TestEmptyDialogFocusEscapeAndNativeCloseFallback(t *testing.T) {
	for _, nonmodal := range []bool{false, true} {
		t.Run(map[bool]string{false: "modal", true: "nonmodal"}[nonmodal], func(t *testing.T) {
			h, w := lifecycleHost(t)
			mode := "true"
			if nonmodal {
				mode = "false"
			}
			doc, err := parser.Parse(`sdui 0.3; page=[open=button("Open",effect="open",target="empty");empty=dialog("A long but bounded empty title",modal=` + mode + `)[flag=command("Flag",toggle=true,key="Shift+F8")]] {x=fill,y=fill};`)
			if err != nil {
				t.Fatal(err)
			}
			if err = h.Adopt(DocumentRequest{Document: doc, Entry: "page", SessionID: "empty", SourceRevision: "source", Sequence: 1, Mode: preparation.Prototype}); err != nil {
				t.Fatal(err)
			}
			w.SetContent(h.Container)
			w.Resize(fyne.NewSize(800, 500))
			w.Show()
			b := h.Current()
			var results []ui.DialogResult
			h.OnDialogResult = func(r ui.DialogResult) { results = append(results, r) }
			button := b.Controls()["page/open"].(*commandButton)
			button.Tapped(&fyne.PointEvent{})
			native := b.surfaces["page/empty"]
			if native == nil {
				t.Fatal("empty surface absent")
			}
			if native.canvas.Focused() != native.anchor {
				t.Fatal("empty surface has no native focus")
			}
			native.anchor.KeyDown(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
			native.anchor.TypedKey(&fyne.KeyEvent{Name: fyne.KeyF8})
			native.anchor.KeyUp(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
			command, _ := b.Session.Command("page/empty/flag")
			state, _ := b.Session.CommandState(command)
			if !state.Checked {
				t.Fatal("empty dialog swallowed declared Shift+F8")
			}
			frame := b.presentation.canvases["page/empty"]
			if native.content.Size() != fyne.NewSize(float32(frame.size.W), float32(frame.size.H)) {
				t.Fatal("title changed model content geometry", native.content.Size(), frame.size)
			}
			native.anchor.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEscape})
			if len(results) != 1 || results[0].Kind != "cancel" {
				t.Fatal("empty Escape did not Cancel once", results)
			}
			if !nonmodal {
				return
			}
			button.Tapped(&fyne.PointEvent{})
			native = b.surfaces["page/empty"]
			// Irrevocable native close must still clean up after the geometry gate fails.
			b.request.PrepareResources = func(ui.Snapshot) error { return errors.New("resource unavailable") }
			native.window.Close()
			if len(results) != 2 || results[1].Kind != "close" || len(b.surfaces) != 0 {
				t.Fatal("native bypass close stranded surface", results)
			}
			native.nativeClosed()
			if len(results) != 2 {
				t.Fatal("duplicate native terminal result")
			}
		})
	}
}

func TestPreparedThemeIconRetainsNativeColorContract(t *testing.T) {
	h, w := lifecycleHost(t)
	source := `sdui 0.3;page=[run=command("Run",toggle=true,icon="play");invoke=button(command="run")] {x=fill,y=fill};`
	doc, err := parser.Parse(source)
	if err != nil {
		t.Fatal(err)
	}
	if err = h.Adopt(DocumentRequest{Document: doc, Entry: "page", SessionID: "icon", SourceRevision: "icon", Sequence: 1, Mode: preparation.Prototype, Icons: map[string]fyne.Resource{"play": theme.MediaPlayIcon()}}); err != nil {
		t.Fatal(err)
	}
	w.SetContent(h.Container)
	w.Resize(fyne.NewSize(800, 500))
	w.Show()
	icon := h.Current().Controls()["page/invoke"].(*commandButton).Icon
	themed, ok := icon.(fyne.ThemedResource)
	if !ok || themed.ThemeColorName() != theme.ColorNameForeground {
		t.Fatal("prepared icon lost native component-theme recoloring", icon)
	}
}

func TestModalChromeResizeAdmissionKeepsPublishedSurface(t *testing.T) {
	h, _ := commandHost(t, false)
	b := h.Current()
	b.Controls()["page/tools/show"].(*commandButton).Tapped(&fyne.PointEvent{})
	old := b.surfaces["page/settings"]
	before := b.Session.Snapshot()
	presentation := b.presentation
	if err := h.Resize(layout.Size{W: 800, H: 40}); err == nil {
		t.Fatal("accepted viewport smaller than dialog chrome/content")
	}
	after := b.Session.Snapshot()
	if before.StateRevision != after.StateRevision || b.presentation != presentation || b.surfaces["page/settings"] != old || old.retiring {
		t.Fatal("failed resize changed published surface")
	}
	if err := h.Resize(layout.Size{W: 800, H: 500}); err != nil {
		t.Fatal(err)
	}
	b.Controls()["page/settings/buttons/cancel"].(*commandButton).Tapped(&fyne.PointEvent{})
	if len(b.surfaces) != 0 {
		t.Fatal("recovery cancel failed")
	}
}

func TestNativeCommandKeyBoundaries(t *testing.T) {
	calls := 0
	route := func(fyne.Shortcut) { calls++ }
	if !commandKeyEvent(&fyne.KeyEvent{Name: fyne.KeyK}, true, false, route) || calls != 1 {
		t.Fatal("Shift letter command dropped")
	}
	if commandKeyEvent(&fyne.KeyEvent{Name: fyne.KeyR}, true, true, route) || calls != 1 {
		t.Fatal("native collection recovery R stolen")
	}
	if commandKeyEvent(&fyne.KeyEvent{Name: fyne.KeyK}, true, true, route) || calls != 1 {
		t.Fatal("Entry text editing stolen")
	}
	header := newPaneHeader()
	header.shift = true
	header.FocusLost()
	if header.shift {
		t.Fatal("header retained released focus modifier")
	}
	key, _ := nativeShortcut("Primary+C")
	if commandShortcut(&fyne.ShortcutCopy{}).ShortcutName() != key.ShortcutName() {
		t.Fatal("native Copy outside Entry lost declared command key")
	}
}

func TestDeclaredKeysPreserveExactEntryEditingAndCanvasLifecycle(t *testing.T) {
	h, w := commandHost(t, false)
	b := h.Current()
	input := b.Controls()["page/edit"].(*Input)
	calls := 0
	input.OnShortcut = func(fyne.Shortcut) { calls++ }
	for _, modifier := range []fyne.KeyModifier{fyne.KeyModifierControl | fyne.KeyModifierAlt, fyne.KeyModifierControl | fyne.KeyModifierShift} {
		input.TypedShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyA, Modifier: modifier})
	}
	if calls != 2 {
		t.Fatal("extra modifiers mistaken for native Entry editing", calls)
	}
	input.TypedShortcut(&fyne.ShortcutSelectAll{})
	if calls != 2 {
		t.Fatal("native editing leaked to command")
	}
	if _, err := nativeShortcut("Shift+F10"); err == nil {
		t.Fatal("reserved context key admitted")
	}
	// The prior callback belongs to the composition owner and must be restored.
	h.Close()
	originalCalls := 0
	w.Canvas().SetOnTypedKey(func(*fyne.KeyEvent) { originalCalls++ })
	next := NewDocumentHost(w.Canvas(), layout.Size{W: 800, H: 500})
	defer next.Close()
	doc, err := parser.Parse(`sdui 0.3;page=[flag=command("Flag",toggle=true,key="F2")];`)
	if err != nil {
		t.Fatal(err)
	}
	if err = next.Adopt(DocumentRequest{Document: doc, Entry: "page", SessionID: "no-focus", SourceRevision: "source", Sequence: 1, Mode: preparation.Prototype}); err != nil {
		t.Fatal(err)
	}
	w.SetContent(next.Container)
	w.Canvas().Unfocus()
	w.Canvas().OnTypedKey()(&fyne.KeyEvent{Name: fyne.KeyF2})
	handle, _ := next.Current().Session.Command("page/flag")
	state, _ := next.Current().Session.CommandState(handle)
	if !state.Checked || originalCalls != 0 {
		t.Fatal("no-focus command missing or duplicated", state.Checked, originalCalls)
	}
	next.Close()
	w.Canvas().OnTypedKey()(&fyne.KeyEvent{Name: fyne.KeyF2})
	if originalCalls != 1 {
		t.Fatal("composition callback not restored", originalCalls)
	}
}
func TestTooltipUsesActualNativeMinimumBeforeTruncation(t *testing.T) {
	h, _ := commandHost(t, false)
	b := h.Current()
	button := b.Controls()["page/tools/toggle"].(*commandButton)
	button.tooltip = "Run the shared action"
	label := widget.NewLabel(button.tooltip)
	expected := container.NewStack(canvas.NewRectangle(color.Transparent), container.NewPadded(label)).MinSize()
	b.showTooltip("page/tools/toggle", button, true)
	if b.tip == nil || b.tip.Size().Width != expected.Width {
		t.Fatal("tooltip prematurely truncates despite available canvas", b.tip, expected)
	}
}
