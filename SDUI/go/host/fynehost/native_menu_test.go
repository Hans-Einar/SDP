package fynehost

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"testing"
)

func TestNativeMenuSynchronousSelectionScope(t *testing.T) {
	for _, gesture := range []string{"pointer", "enter", "space", "escape", "outside"} {
		t.Run(gesture, func(t *testing.T) {
			app := test.NewApp()
			defer app.Quit()
			w := app.NewWindow("menu scope")
			defer w.Close()
			w.Resize(fyne.NewSize(600, 400))
			w.Show()
			calls, closes := 0, 0
			var m *nativeMenu
			leaf := fyne.NewMenuItem("Run", func() {
				calls++
				if !m.scoped || !m.dismissed || !m.claimed || m.closed {
					t.Fatal("action escaped exact live selection scope")
				}
				if w.Canvas().Overlays().Top() != nil {
					t.Fatal("native menu not hidden before Action")
				}
			})
			model := fyne.NewMenu("Root", &fyne.MenuItem{Label: "Nested", ChildMenu: fyne.NewMenu("Nested", leaf)})
			m = newNativeMenu(model, w.Canvas(), func() { closes++ })
			m.show(fyne.NewPos(50, 40))
			parent := m.popup.Items[0]
			parent.(desktop.Hoverable).MouseIn(&desktop.MouseEvent{})
			child := parent.(nativeMenuItemChildren).Child()
			if child == nil || !child.Visible() {
				t.Fatal("native submenu did not open")
			}
			child.ActivateNext()
			switch gesture {
			case "pointer":
				item := child.Items[0]
				p := app.Driver().AbsolutePositionForObject(item)
				s := item.Size()
				test.TapCanvas(w.Canvas(), p.Add(fyne.NewPos(s.Width/2, s.Height/2)))
			case "enter":
				m.input.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
			case "space":
				m.input.TypedKey(&fyne.KeyEvent{Name: fyne.KeySpace})
			case "escape":
				m.input.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEscape})
			case "outside":
				test.TapCanvas(w.Canvas(), fyne.NewPos(550, 350))
			}
			expected := 1
			if gesture == "escape" || gesture == "outside" {
				expected = 0
			}
			if calls != expected || closes != 1 || !m.closed || m.scoped {
				t.Fatalf("scope did not finalize synchronously: calls=%d close=%d closed=%v scoped=%v", calls, closes, m.closed, m.scoped)
			}
			m.model.Items[0].ChildMenu.Items[0].Action()
			child.Items[0].(fyne.Tappable).Tapped(&fyne.PointEvent{})
			m.input.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
			m.popup.Dismiss()
			if calls != expected || closes != 1 {
				t.Fatal("late native callback reused capture", calls, closes)
			}
		})
	}
}

func TestNativeMenuOldDismissLeavesReplacement(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	w := app.NewWindow("menu replacement")
	defer w.Close()
	w.Resize(fyne.NewSize(600, 400))
	w.Show()
	var replacement *nativeMenu
	old := newNativeMenu(fyne.NewMenu("Old", fyne.NewMenuItem("Replace", func() {
		replacement = newNativeMenu(fyne.NewMenu("New", fyne.NewMenuItem("New", func() {})), w.Canvas(), func() {})
		replacement.show(fyne.NewPos(90, 60))
	})), w.Canvas(), func() {})
	old.show(fyne.NewPos(40, 30))
	old.popup.ActivateNext()
	old.input.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
	old.popup.Dismiss()
	old.close()
	if replacement == nil || replacement.closed || w.Canvas().Overlays().Top() != replacement.overlay {
		t.Fatal("old cleanup removed replacement")
	}
	replacement.close()
}

func TestNativeMenuHomeEndLongerNestedMenu(t *testing.T) {
	for _, key := range []fyne.KeyName{fyne.KeyHome, fyne.KeyEnd} {
		t.Run(string(key), func(t *testing.T) {
			app := test.NewApp()
			defer app.Quit()
			w := app.NewWindow("menu boundary")
			defer w.Close()
			w.Resize(fyne.NewSize(600, 400))
			w.Show()
			selected := -1
			items := make([]*fyne.MenuItem, 7)
			for i := range items {
				i := i
				items[i] = fyne.NewMenuItem(string(rune('A'+i)), func() { selected = i })
			}
			m := newNativeMenu(fyne.NewMenu("Root", &fyne.MenuItem{Label: "Nested", ChildMenu: fyne.NewMenu("Nested", items...)}), w.Canvas(), func() {})
			m.show(fyne.NewPos(30, 20))
			m.popup.ActivateNext()
			m.popup.ActivateLastSubmenu()
			for i := 0; i < 3; i++ {
				m.popup.ActivateNext()
			}
			m.input.TypedKey(&fyne.KeyEvent{Name: key})
			m.input.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
			want := 0
			if key == fyne.KeyEnd {
				want = 6
			}
			if selected != want {
				t.Fatalf("nested boundary got %d want %d", selected, want)
			}
		})
	}
}

func TestNativeMenuFreshOpeningsDoNotDecorateCaller(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()
	w := app.NewWindow("reuse")
	defer w.Close()
	w.Resize(fyne.NewSize(600, 400))
	w.Show()
	calls := 0
	leaf := fyne.NewMenuItem("Run", func() { calls++ })
	model := fyne.NewMenu("Root", leaf)
	for i := 0; i < 2; i++ {
		m := newNativeMenu(model, w.Canvas(), func() {})
		m.show(fyne.NewPos(20, 20))
		m.popup.ActivateNext()
		m.input.TypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
	}
	if calls != 2 {
		t.Fatal("fresh opening retained old Action wrapper", calls)
	}
	leaf.Action()
	if calls != 3 {
		t.Fatal("caller model was mutated")
	}
}
func TestNativeMenuDisabledKeyboardCannotClaim(t *testing.T) {
	for _, key := range []fyne.KeyName{fyne.KeyEnter, fyne.KeySpace} {
		t.Run(string(key), func(t *testing.T) {
			app := test.NewApp()
			defer app.Quit()
			w := app.NewWindow("disabled")
			defer w.Close()
			w.Resize(fyne.NewSize(600, 400))
			w.Show()
			calls := 0
			leaf := fyne.NewMenuItem("Disabled", func() { calls++ })
			leaf.Disabled = true
			m := newNativeMenu(fyne.NewMenu("Root", leaf), w.Canvas(), func() {})
			m.show(fyne.NewPos(20, 20))
			m.popup.ActivateNext()
			m.input.TypedKey(&fyne.KeyEvent{Name: key})
			if calls != 0 || m.claimed || !m.closed {
				t.Fatal("disabled keyboard selection claimed action", calls, m.claimed, m.closed)
			}
		})
	}
}
