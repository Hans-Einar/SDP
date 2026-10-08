package fynehost

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

type documentMenu struct {
	scope  ui.MenuScope
	native *nativeMenu
	path   string
	rows   map[string]string // native structural row index -> exact source path
}

func (b *Bundle) connectCommands() {
	for path, obj := range b.view.Controls {
		path := path
		switch c := obj.(type) {
		case *contextControl:
			if b.hasContext(path) {
				c.context = func(pos fyne.Position) { b.openContext(path, "", pos) }
			}
		case *commandButton:
			c.shortcut = func(s fyne.Shortcut) { b.routeShortcut(b.canvasFor(path), s) }
			c.focus = func() {
				if !b.livePane() {
					return
				}
				if w, ok := b.Session.Widget(path); ok {
					if err := b.owner.Mutate(func(s *ui.Session) error { return s.Focus(w.Handle) }); err != nil {
						return
					}
					b.owner.ensureWidget(path)
				}
			}
			c.tipChanged = func(show bool) { b.showTooltip(path, c, show) }
			c.escape = func() { b.escapeSurface(path) }
			if b.hasContext(path) {
				c.context = func(pos fyne.Position) { b.openContext(path, "", pos) }
			}
			if presentation, ok := b.Session.Snapshot().Presentations[path]; ok {
				handle := presentation.Handle
				b.view.Actions[path] = func(_, _ string) error { return b.invokeCommand(handle, "button") }
			} else if w, ok := b.Session.Widget(path); ok && w.Handle.Kind == "button" {
				b.view.Actions[path] = func(_, _ string) error {
					if !b.livePane() {
						return nil
					}
					return b.owner.Mutate(func(s *ui.Session) error {
						return s.Dispatch(ui.Event{Handle: w.Handle, ModelRevision: s.Revision, Sequence: s.Sequence() + 1, Kind: ui.Activate})
					})
				}
			}
		case *paneHeader:
			c.shortcut = func(s fyne.Shortcut) { b.routeShortcut(b.canvasFor(path), s) }
			c.escape = func() { b.escapeSurface(path) }
		case *paneDivider:
			c.escape = func() { b.escapeSurface(path) }
			c.shortcut = func(s fyne.Shortcut) { b.routeShortcut(b.canvasFor(path), s) }
		case *Input:
			c.OnShortcut = func(s fyne.Shortcut) { b.routeShortcut(b.canvasFor(path), s) }
			c.OnEscape = func() bool {
				if !b.livePane() {
					return false
				}
				w, ok := b.Session.Widget(path)
				if ok && w.Dirty {
					return false
				}
				return b.escapeSurface(path)
			}
			if b.hasContext(path) {
				c.OnContext = func(pos fyne.Position) { b.openContext(path, "", pos) }
			}
		}
	}
}
func (b *Bundle) invokeCommand(origin ui.Handle, via string) error {
	if !b.livePane() {
		return nil
	}
	before := b.Session.Snapshot().ActiveSurface
	err := b.owner.Mutate(func(s *ui.Session) error {
		event, err := s.CaptureCommand(origin, via, nil)
		if err != nil {
			return err
		}
		result, err := s.DispatchInteraction(event)
		return interactionError(result, err)
	})
	b.focusAfterGesture(before)
	return err
}
func (b *Bundle) dispatchDialog(target ui.SurfaceTarget, kind ui.EventKind) {
	if !b.livePane() {
		return
	}
	before := b.Session.Snapshot().ActiveSurface
	b.owner.Mutate(func(s *ui.Session) error {
		event, err := s.CaptureDialog(target, kind)
		if err != nil {
			return err
		}
		result, err := s.DispatchInteraction(event)
		return interactionError(result, err)
	})
	if before != nil && *before == target {
		b.focusAfterGesture(before)
	}
}
func (b *Bundle) escapeSurface(path string) bool {
	if !b.livePane() {
		return false
	}
	if owner := surfacePath(b.Session.Snapshot(), path); owner != "" {
		if surface, ok := b.Session.Snapshot().Surfaces[owner]; ok && surface.Open {
			b.dispatchDialog(surface.Target, ui.Cancel)
			return true
		}
	}
	return false
}
func (b *Bundle) openContext(path string, item ui.ItemID, position fyne.Position) {
	if !b.livePane() {
		return
	}
	w, ok := b.Session.Widget(path)
	if !ok {
		return
	}
	context := ui.ContextTarget{Widget: w.Handle, ModelRevision: b.Session.Revision}
	if item != "" {
		target, err := b.Session.Target(w.Handle, item)
		if err != nil {
			b.owner.status(err)
			return
		}
		context.Item = &target
	}
	ids, err := parser.ResolveInteractions(b.presentation.snapshot.Root)
	if err != nil {
		b.owner.status(err)
		return
	}
	root := ""
	b.presentation.snapshot.Root.Walk(func(n *parser.Instance) {
		if root == "" && n.Widget == "menu" && n.Argument("mode") == "context" && ids[n.Path].Target == path {
			root = n.Path
		}
	})
	if root != "" {
		b.openMenu(root, context, position)
	}
}
func (b *Bundle) openMenu(path string, context ui.ContextTarget, position fyne.Position) {
	if !b.livePane() {
		return
	}
	handle, ok := b.Session.Menu(path)
	if !ok {
		return
	}
	var scope ui.MenuScope
	if err := b.owner.Mutate(func(s *ui.Session) error { var err error; scope, err = s.OpenMenu(handle, context); return err }); err != nil {
		return
	}
	snapshot := b.Session.Snapshot()
	var source *parser.Instance
	snapshot.Root.Walk(func(n *parser.Instance) {
		if n.Path == path {
			source = n
		}
	})
	if source == nil {
		return
	}
	dm := &documentMenu{scope: scope, path: path, rows: map[string]string{}}
	var build func(*parser.Instance, string) *fyne.Menu
	build = func(node *parser.Instance, index string) *fyne.Menu {
		model := fyne.NewMenu(node.Argument("label"))
		var add func(*parser.Instance)
		add = func(n *parser.Instance) {
			if n.Layout["visible"] == false {
				return
			}
			idx := fmt.Sprintf("%s/%d", index, len(model.Items))
			switch n.Widget {
			case "separator":
				if len(model.Items) > 0 && !model.Items[len(model.Items)-1].IsSeparator {
					model.Items = append(model.Items, fyne.NewMenuItemSeparator())
					dm.rows[idx] = n.Path
				}
			case "menuGroup":
				model.Items = append(model.Items, &fyne.MenuItem{Label: n.Argument("label"), Disabled: true})
				dm.rows[idx] = n.Path
				for _, row := range n.Rows {
					for _, child := range row {
						add(child)
					}
				}
			case "menu":
				model.Items = append(model.Items, &fyne.MenuItem{Label: n.Argument("label"), ChildMenu: build(n, idx), Disabled: n.Layout["enabled"] == false})
				dm.rows[idx] = n.Path
			case "item":
				presentation, ok := snapshot.Presentations[n.Path]
				if !ok {
					return
				}
				item := fyne.NewMenuItem(presentation.Label, func() { b.invokeCommand(presentation.Handle, "menu") })
				item.Icon = b.icons[presentation.Icon]
				for _, cmd := range snapshot.Commands {
					if cmd.Handle == presentation.Command {
						item.Checked = cmd.Checked
						break
					}
				}
				_, err := b.Session.CaptureCommand(presentation.Handle, "menu", nil)
				item.Disabled = err != nil || !presentation.Enabled || !presentation.Visible
				model.Items = append(model.Items, item)
				dm.rows[idx] = n.Path
			}
		}
		for _, row := range node.Rows {
			for _, child := range row {
				add(child)
			}
		}
		if len(model.Items) > 0 && model.Items[len(model.Items)-1].IsSeparator {
			model.Items = model.Items[:len(model.Items)-1]
		}
		return model
	}
	model := build(source, "")
	c := b.canvasFor(path)
	if c == nil {
		return
	}
	dm.native = newNativeMenu(model, c, func() {
		if b.closed || b.owner.current != b {
			return
		}
		if b.menus[path] == dm {
			delete(b.menus, path)
		}
		b.owner.Mutate(func(s *ui.Session) error { return s.CloseMenu(scope) })
	})
	dm.native.nativeTransition = func(fn func()) { muted := b.muted; b.muted = true; defer func() { b.muted = muted }(); fn() }
	b.menus[path] = dm
	if position.IsZero() {
		if control := b.view.Controls[path]; control != nil {
			position = fyne.CurrentApp().Driver().AbsolutePositionForObject(control).Add(fyne.NewPos(0, control.Size().Height))
		}
	}
	muted := b.muted
	b.muted = true
	dm.native.show(position)
	b.muted = muted
}
func (b *Bundle) syncMenus() {
	if b.Session == nil {
		return
	}
	snapshot := b.Session.Snapshot()
	for path, m := range b.menus {
		state, ok := snapshot.Menus[path]
		if !ok || !state.Open || state.Scope != m.scope || snapshot.StateRevision != m.scope.StateRevision {
			// During synchronous native selection the runtime decides whether capture
			// is still valid; don't close it speculatively before Dispatch validates.
			if m.native.scoped {
				continue
			}
			m.native.close()
		}
	}
}

// All adapters use this one canonical route after their native editing/navigation
// handlers decline a shortcut. Modal canvases share registrations with parents;
// runtime's active surface, not registration order, selects the command scope.
func (b *Bundle) routeShortcut(c fyne.Canvas, shortcut fyne.Shortcut) {
	nativeEditing := nativeEditingShortcut(shortcut)
	shortcut = commandShortcut(shortcut)
	if !b.livePane() || c == nil {
		return
	}
	context := &desktop.CustomShortcut{KeyName: fyne.KeyF10, Modifier: fyne.KeyModifierShift}
	if shortcut.ShortcutName() == context.ShortcutName() {
		for path, control := range b.view.Controls {
			if focused, ok := controlFocusable(control); ok && focused == c.Focused() {
				item := ui.ItemID("")
				if collection, ok := control.(*CollectionControl); ok {
					item = collection.state.Focused
				}
				b.openContext(path, item, fyne.CurrentApp().Driver().AbsolutePositionForObject(control))
				return
			}
		}
		return
	}
	snapshot := b.Session.Snapshot()
	surface := ""
	if snapshot.ActiveSurface != nil {
		for path, state := range snapshot.Surfaces {
			if state.Open && state.Target == *snapshot.ActiveSurface {
				surface = path
				break
			}
		}
	}
	for _, key := range b.keys {
		if key.surface != surface || key.shortcut.ShortcutName() != shortcut.ShortcutName() {
			continue
		}
		if _, editing := c.Focused().(*Input); editing && nativeEditing {
			return
		}
		if handle, ok := b.Session.Command(key.path); ok {
			b.invokeCommand(handle, "key")
		}
		return
	}
}
func (b *Bundle) installKeys(c fyne.Canvas, _ string) {
	if c == nil {
		return
	}
	if b.keyCanvases == nil {
		b.keyCanvases = map[fyne.Canvas]bool{}
	}
	if b.keyCanvases[c] {
		return
	}
	b.keyCanvases[c] = true
	b.installCanvasKeys(c)
	for _, key := range b.keys {
		c.AddShortcut(key.shortcut, func(s fyne.Shortcut) { b.routeShortcut(c, s) })
	}
	c.AddShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyF10, Modifier: fyne.KeyModifierShift}, func(s fyne.Shortcut) { b.routeShortcut(c, s) })
}
func (b *Bundle) removeKeys(c fyne.Canvas, _ string) {
	if !b.keyCanvases[c] {
		return
	}
	delete(b.keyCanvases, c)
	b.removeCanvasKeys(c)
	for _, key := range b.keys {
		c.RemoveShortcut(key.shortcut)
	}
	c.RemoveShortcut(&desktop.CustomShortcut{KeyName: fyne.KeyF10, Modifier: fyne.KeyModifierShift})
}

func (b *Bundle) hasContext(path string) bool {
	root := b.Session.Snapshot().Root
	ids, err := parser.ResolveInteractions(root)
	if err != nil {
		return false
	}
	found := false
	root.Walk(func(n *parser.Instance) {
		if n.Widget == "menu" && n.Argument("mode") == "context" && ids[n.Path].Target == path {
			found = true
		}
	})
	return found
}
func (b *Bundle) focusAfterGesture(before *ui.SurfaceTarget) {
	if !b.livePane() {
		return
	}
	after := b.Session.Snapshot().ActiveSurface
	if before == nil && after == nil || before != nil && after != nil && *before == *after {
		return
	}
	c := b.owner.canvas
	if after != nil {
		for _, native := range b.surfaces {
			if native.target == *after {
				c = native.canvas
				break
			}
		}
	}
	if w := windowForCanvas(c); w != nil {
		w.RequestFocus()
	}
}
